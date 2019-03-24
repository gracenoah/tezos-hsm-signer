package signer

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/btcsuite/btcutil/base58"
)

// Tezos Constants from:
// https://gitlab.com/tezos/tezos/blob/master/src/lib_crypto/base58.ml
const (
	/* Public Key Hashes */
	tzEd25519PublicKeyHash   = "06a19f" // tz1
	tzSecp256k1PublicKeyHash = "06a1a1" // tz2
	tzP256PublicKeyHash      = "06a1a4" // tz3

	/* Public Keys */
	tzEd25519PublicKey   = "0d0f25d9" // edpk
	tzSecp256k1PublicKey = "03fee256" // sppk
	tzP256PublicKey      = "03b28b7f" // p2pk

	/* Secret Keys */
	tzEd25519Seed        = "0d0f3a07" // edsk
	tzSecp256k1SecretKey = "11a2e0c9" // spsk
	tzP256SecretKey      = "1051eebd" // p2sk

	/* Encrypted Secret Keys */
	tzEd25519EncryptedSeed        = "075a3cb329" // edesk
	tzSecp256k1EncryptedSecretKey = "09edf1ae96" // spesk
	tzP256EncryptedSecretKey      = "09303973ab" // p2esk

	/* Signatures */
	tzEd25519Signature   = "09f5cd8612" // edsig (len: 99)
	tzSecp256k1Signature = "0d7365133f" // spsig1 (len: 99)
	tzP256Signature      = "36f02c34"   // p2sig (len: 98)
	tzGenericSignature   = "04822b"     // sig (len: 96)

	/* Chain ID */
	tzChainID = "575200" // Net(15)
)

// getSignaturePrefix for a given key to produce the correct
// base58 check encoded tz result
func getSignaturePrefix(key *Key) ([]byte, error) {
	var prefix []byte
	switch key.Curve() {
	case curveEd25519:
		prefix, _ = hex.DecodeString(tzEd25519Signature)
	case curveSecp256k1:
		prefix, _ = hex.DecodeString(tzSecp256k1Signature)
	case curveNistP256:
		prefix, _ = hex.DecodeString(tzP256Signature)
	default:
		return nil, fmt.Errorf("Unknown pkh type %v", string(key.PublicKeyHash[2]))
	}
	return prefix, nil
}

// isValidSignatureFormat ensures a b58 check endoded signature is formatted
// correctly.  It does *not* verify the cryptographic validity of the signature.
func isValidSignatureFormat(key *Key, sig string) bool {
	switch key.Curve() {
	case curveEd25519:
		return strings.HasPrefix(sig, "edsig") && len(sig) == 99
	case curveSecp256k1:
		return strings.HasPrefix(sig, "spsig1") && len(sig) == 99
	case curveNistP256:
		return strings.HasPrefix(sig, "p2sig") && len(sig) == 98
	default:
		return strings.HasPrefix(sig, "sig") && len(sig) == 96
	}
}

func b58CheckEncode(prefix []byte, bytes []byte) string {
	message := append(prefix, bytes...)
	// SHA^2
	h := sha256.Sum256(message)
	h2 := sha256.Sum256(h[:])
	// Append first four of the hash
	finalMessage := append(message, h2[:4]...)
	// b58 encode the response
	encoded := base58.Encode(finalMessage)
	return encoded
}

// PubkeyHashToByteString strips the prefix and checksum bytes,
// returning only the pubkeyhash bytes
func PubkeyHashToByteString(pubkeyhash string) string {
	return hex.EncodeToString(base58.Decode(pubkeyhash)[3:23])
}
