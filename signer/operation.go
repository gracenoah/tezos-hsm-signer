package signer

import (
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"math/big"
	"strings"
)

// Operation parses and validates an arbitrary tz request
// to sign some message before passing to the Signer
type Operation struct {
	hex []byte
}

// Watermark of different operations
// According to: https://gitlab.com/tezos/tezos/blob/master/src/lib_crypto/signature.ml#L525
const (
	opWatermarkBlock       = 0x01
	opWatermarkEndorsement = 0x02
	opWatermarkGeneric     = 0x03
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
	switch op.Watermark() {
	case opWatermarkGeneric:
		debugln("Operation is Generic.  Possibly a Transaction")
	case opWatermarkBlock:
		debugln("Operation is a Block at level: ", op.Level().String())
	case opWatermarkEndorsement:
		debugln("Operation is an Endorsement at level: ", op.Level().String())
	default:
		return nil, fmt.Errorf("Operation: Unsupported Operation Watermark: %v", op.Watermark())
	}

	return &op, nil
}

// Hex returns a copy of the parsed hex bytes of the operation
func (op *Operation) Hex() []byte {
	hexCopy := make([]byte, len(op.hex))
	copy(hexCopy, op.hex)
	return hexCopy
}

// Watermark of this tezos operation included in the operation
func (op *Operation) Watermark() uint8 {
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
	if op.Watermark() == opWatermarkBlock {
		return new(big.Int).SetBytes(op.hex[5:9])
	} else if op.Watermark() == opWatermarkEndorsement {
		return new(big.Int).SetBytes(op.hex[len(op.hex)-4:])
	}
	log.Println("Warn: Requested level for unexpected watermark", op.Watermark())
	return nil
}
