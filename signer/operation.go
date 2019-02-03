package signer

import (
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"math/big"
	"strings"

	"golang.org/x/crypto/blake2b"
)

// Operation parses and validates an arbitrary tz request
// to sign some message before passing to the Signer
type Operation struct {
	hex []byte
}

// Watermark of different operations
// According to: https://gitlab.com/tezos/tezos/blob/master/src/lib_crypto/signature.ml#L523
const (
	opTypeBlock       = 0x01
	opTypeEndorsement = 0x02
	opTypeGeneric     = 0x03
)

// ParseOperation parses a raw byte string into a meaningful tz operation
// and performs simple validation
func ParseOperation(opBytes []byte) (*Operation, error) {

	// Must begin and end with quotes
	opString := strings.TrimSpace(string(opBytes))
	if !strings.HasPrefix(opString, "\"") || !strings.HasSuffix(opString, "\"") {
		return nil, errors.New("A valid operation begins and ends with a quote")
	}
	opString = strings.Trim(opString, "\"")

	// Must be valid hex chars
	parsedHex, err := hex.DecodeString(opString)
	if err != nil {
		return nil, err
	}

	op := Operation{
		hex: parsedHex,
	}

	// Validate and print debug statements
	switch op.Type() {
	case opTypeGeneric:
		debugln("Operation is Generic.  Possibly a Transaction")
	case opTypeBlock:
		debugln("Operation is a Block at level: ", op.Level().String())
	case opTypeEndorsement:
		debugln("Operation is an Endorsement at level: ", op.Level().String())
	default:
		return nil, fmt.Errorf("Operation: Unsupported Operation Type: %v", op.Type())
	}

	return &op, nil
}

// Hex returns a copy of the parsed hex bytes of the operation
func (op *Operation) Hex() []byte {
	hexCopy := make([]byte, len(op.hex))
	copy(hexCopy, op.hex)
	return hexCopy
}

// Type of this tezos operation included in the operation
func (op *Operation) Type() uint8 {
	return op.hex[0]
}

// ChainID to determine what we're running on
func (op *Operation) ChainID() string {
	chainID := op.hex[1:5]
	prefix, _ := hex.DecodeString(tzChainID)
	return b58CheckEncode(prefix, chainID)
}

// Level returns a copy of the level, if one can be parsed from this operation
func (op *Operation) Level() *big.Int {
	if op.Type() == opTypeBlock {
		return new(big.Int).SetBytes(op.hex[5:9])
	} else if op.Type() == opTypeEndorsement {
		return new(big.Int).SetBytes(op.hex[len(op.hex)-4:])
	}
	log.Println("Warn: Requested level for unexpected type", op.Type())
	return nil
}

// TzSign this operation with the provided Signer and Key
func (op *Operation) TzSign(signer Signer, key *Key) (string, error) {
	msg := op.Hex()
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
