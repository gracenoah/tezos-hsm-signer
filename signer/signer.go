package signer

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"math"

	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/ed25519"

	cloudkms "cloud.google.com/go/kms/apiv1"
	kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"

	"github.com/miekg/pkcs11"
)

// Signer is a generic interface for a signer
type Signer interface {
	Sign(message []byte, key *Key) ([]byte, error)
}

// PKCS11Signer is responsible for signing an arbitrary byte slice with the given
// Key stored within the HSM
type PKCS11Signer struct {
	UserPin string `yaml:"UserPin"`
	LibPath string `yaml:"LibPath"`
}

var _ Signer = &PKCS11Signer{}

// getPrivateKeyHandle returns the handle of the private key loaded
// into your HSM for the corresponding opened session
func (*PKCS11Signer) getPrivateKeyHandle(context *pkcs11.Ctx, session pkcs11.SessionHandle, tokenLabel string) (pkcs11.ObjectHandle, error) {
	var noKeyFound pkcs11.ObjectHandle = math.MaxUint8

	// Combine attributes to select the correct key
	template := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_PRIVATE_KEY),
	}
	if len(tokenLabel) > 0 {
		// SoftHSM does not support label queries, so make this optional
		template = append(template, pkcs11.NewAttribute(pkcs11.CKA_LABEL, tokenLabel))
	}

	err := context.FindObjectsInit(session, template)
	if err != nil {
		log.Println("Error initializing FindObjects, returning no keys. Error: ", err)
		return noKeyFound, err
	}

	// find a maximum of 2 objects that match the template
	keyHandles, _, err := context.FindObjects(session, 2)
	if err != nil {
		log.Println("Error Finding Objects, returning no keys. Error: ", err)
		return noKeyFound, err
	}

	// complete the object search
	err = context.FindObjectsFinal(session)
	if err != nil {
		log.Println("Error Finalizing FindObjects, returning no keys. Error:", err)
		return noKeyFound, err
	}

	// must have found exactly one key
	if len(keyHandles) == 0 {
		return noKeyFound, errors.New("Querier key not found")
	} else if len(keyHandles) > 1 {
		return noKeyFound, errors.New("Multiple matching keys, unsure how to proceed. Returning no results")
	}

	return keyHandles[0], nil
}

// Is Slot available in the provided slice of slots
func (*PKCS11Signer) isSlotAvailable(slotID uint, slots []uint) bool {
	for _, value := range slots {
		if value == slotID {
			return true
		}
	}
	return false
}

// Sign a transaction request
func (hsm *PKCS11Signer) Sign(message []byte, key *Key) ([]byte, error) {
	context := pkcs11.New(hsm.LibPath)

	err := context.Initialize()
	if err != nil {
		log.Println("Error initializing the shared object.  Are you sure this is available? Error: ", err)
		return nil, err
	}
	defer context.Destroy()
	defer context.Finalize()

	// Get slots where tokens are present
	slots, err := context.GetSlotList(true)
	if err != nil {
		log.Println("Could not get slot list. Error: ", err)
		return nil, err
	}

	// Requested slot must be present
	if !hsm.isSlotAvailable(key.HsmSlot, slots) {
		debugln("Available slots are: ", slots)
		return nil, fmt.Errorf("Slot %v not found", key.HsmSlot)
	}

	session, err := context.OpenSession(key.HsmSlot, pkcs11.CKF_SERIAL_SESSION|pkcs11.CKF_RW_SESSION)
	if err != nil {
		fmt.Println("Error opening session: ", err)
		return nil, err
	}
	defer context.CloseSession(session)

	err = context.Login(session, pkcs11.CKU_USER, hsm.UserPin)
	if err != nil {
		fmt.Println("Error logging into HSM: ", err)
		return nil, err
	}
	defer context.Logout(session)

	// Get a handle to our private key
	privateKey, err := hsm.getPrivateKeyHandle(context, session, key.HsmLabel)
	if err != nil {
		fmt.Println("Error retrieving a handle to our private key: ", err)
		return nil, err
	}

	// Init ECDSA signature with this private key handle
	context.SignInit(session, []*pkcs11.Mechanism{pkcs11.NewMechanism(pkcs11.CKM_ECDSA, nil)}, privateKey)

	// Sign
	signedMsg, err := context.Sign(session, message)
	if err != nil {
		fmt.Println("Error signing the message", err)
		return nil, err
	}

	return signedMsg, nil
}

type inMemorySigner struct {
	privateKey    ed25519.PrivateKey
	publicKeyHash string
}

// NewInMemorySigner creates a signer from a key stored plaintext in memory.
// It is not suitable for production use.
func NewInMemorySigner(privateKey ed25519.PrivateKey) Signer {
	publicKeyHash, err := blake2b.New(20, nil)
	if err != nil {
		panic(err.Error())
	}
	_, err = publicKeyHash.Write(privateKey.Public().(ed25519.PublicKey))
	if err != nil {
		panic(err.Error())
	}
	publicKeyHashBytes := publicKeyHash.Sum([]byte{})
	prefix, _ := hex.DecodeString(tzEd25519PublicKeyHash)
	publicKeyHashString := b58CheckEncode(prefix, publicKeyHashBytes)
	return &inMemorySigner{
		privateKey:    privateKey,
		publicKeyHash: publicKeyHashString,
	}
}

func (i *inMemorySigner) Sign(message []byte, key *Key) ([]byte, error) {
	if key.PublicKeyHash != i.publicKeyHash {
		return nil, fmt.Errorf("unknown key %s, expected %s", key.PublicKeyHash, i.publicKeyHash)
	}
	return ed25519.Sign(i.privateKey, message), nil
}

type googleCloudKMSSigner struct{}

// NewGoogleCloudKMSSigner creates a signer backed by Google Cloud KMS
func NewGoogleCloudKMSSigner() Signer {
	return &googleCloudKMSSigner{}
}

func (g *googleCloudKMSSigner) Sign(message []byte, key *Key) ([]byte, error) {
	ctx := context.Background()
	client, err := cloudkms.NewKeyManagementClient(ctx)
	if err != nil {
		return nil, err
	}
	req := &kmspb.AsymmetricSignRequest{
		Name: key.Name,
		Digest: &kmspb.Digest{
			// It's actually Blake2b.Sum256, not SHA256, but google doesn't know the difference
			Digest: &kmspb.Digest_Sha256{
				Sha256: message,
			},
		},
	}
	response, err := client.AsymmetricSign(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("asymmetric sign request failed: %+v", err)
	}
	return response.Signature, nil
}
