package signer

import (
	"errors"
	"fmt"
	"log"
	"math"

	"github.com/miekg/pkcs11"
)

// Signer is a generic interface that HSM implements
type Signer interface {
	Sign(message []byte, key *Key) ([]byte, error)
}

// Hsm is responsible for signing an arbitrary byte slice with the given
// Key stored within the HSM
type Hsm struct {
	UserPin string `yaml:"UserPin"`
	LibPath string `yaml:"LibPath"`
}

// getPrivateKeyHandle returns the handle of the private key loaded
// into your HSM for the corresponding opened session
func getPrivateKeyHandle(context *pkcs11.Ctx, session pkcs11.SessionHandle, tokenLabel string) (pkcs11.ObjectHandle, error) {
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
func isSlotAvailable(slotID uint, slots []uint) bool {
	for _, value := range slots {
		if value == slotID {
			return true
		}
	}
	return false
}

// Sign a transaction request
func (hsm *Hsm) Sign(message []byte, key *Key) ([]byte, error) {
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
	if !isSlotAvailable(key.HsmSlot, slots) {
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
	privateKey, err := getPrivateKeyHandle(context, session, key.HsmLabel)
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
