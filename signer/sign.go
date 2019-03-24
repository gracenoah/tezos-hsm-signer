package signer

import (
	"encoding/hex"
	"fmt"

	"golang.org/x/crypto/blake2b"
)

// TzSign this operation with the provided Signer and Key
func (op *Operation) TzSign(signer Signer, key *Key) (string, error) {
	msg := op.Hex()
	debugln("Signing for key: ", key.PublicKeyHash)
	debugln("About to sign raw bytes hex.EncodeToString(bytes): ", hex.EncodeToString(msg))

	// Take the 256 bit (32 Byte) Blake2b Hash of the operation
	digest := blake2b.Sum256(msg)

	// Sign
	signedMsg, err := signer.Sign(digest[:], key)
	if err != nil {
		return "", err
	}
	debugln("Signed bytes hex.EncodeToString(bytes): ", hex.EncodeToString(signedMsg))

	// Ensure ECDSA sig's `S` value is modulo their curve's `n` parameter per BIP 62
	if key.IsECDSA() {
		signedMsg = StrictECModN(key, signedMsg)
		debugln("Signed bytes StrictECModN(hex.EncodeToString(bytes)): ", hex.EncodeToString(signedMsg))
	}

	// Get the correct signature prefix
	prefix, err := getSignaturePrefix(key)
	if err != nil {
		return "", err
	}

	// b58 Check Encode the result
	signed := b58CheckEncode(prefix, signedMsg)

	// Result should begin with:  p2sig, edsig, spsig1 or sig
	if !isValidSignatureFormat(key, signed) {
		return "", fmt.Errorf("b58 check encoded result was not correctly formatted: %v", signed)
	}

	return signed, nil
}
