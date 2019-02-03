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

// Request parses and validates an arbitrary request
// to sign some message before passing to the Signer
type Request struct {
	hex []byte
}

// Watermark of different operations
// According to: https://gitlab.com/tezos/tezos/blob/master/src/lib_crypto/signature.ml#L523
const (
	opTypeBlock       = 0x01
	opTypeEndorsement = 0x02
	opTypeGeneric     = 0x03
)

// ParseRequest parses a raw byte string into a meaningful signing payload
// if the request can be parsed and is valid
func ParseRequest(requestBytes []byte) (*Request, error) {

	// Must begin and end with quotes
	requestString := strings.TrimSpace(string(requestBytes))
	if !strings.HasPrefix(requestString, "\"") || !strings.HasSuffix(requestString, "\"") {
		return nil, errors.New("request: A signing request must begin and end with a quote")
	}
	requestString = strings.Trim(requestString, "\"")

	// Must be valid hex chars
	parsedHex, err := hex.DecodeString(requestString)
	if err != nil {
		return nil, err
	}

	request := Request{
		hex: parsedHex,
	}

	// Validate and print debug statements
	switch request.OpType() {
	case opTypeGeneric:
		debugln("Request is Generic.  Possibly a Transaction")
	case opTypeBlock:
		debugln("Request is a Block at level: ", request.Level().String())
	case opTypeEndorsement:
		debugln("Request is an Endorsement at level: ", request.Level().String())
	default:
		return nil, fmt.Errorf("Unsupported Operation Type: %v", request.OpType())
	}

	return &request, nil
}

// OpType of this tezos operation included in the signing request
func (request *Request) OpType() uint8 {
	return request.hex[0]
}

// Hex returns a copy of the parsed hex bytes of the signing request
func (request *Request) Hex() []byte {
	hexCopy := make([]byte, len(request.hex))
	copy(hexCopy, request.hex)
	return hexCopy
}

// ChainID to determine what we're running on
func (request *Request) ChainID() string {
	chainID := request.hex[1:5]
	return hex.EncodeToString(chainID)
}

// Level returns a copy of the level, if one was parsed from this request
func (request *Request) Level() *big.Int {
	if request.OpType() == opTypeBlock {
		return new(big.Int).SetBytes(request.hex[5:9])
	} else if request.OpType() == opTypeEndorsement {
		return new(big.Int).SetBytes(request.hex[len(request.hex)-4:])
	}
	log.Println("Warn: Requested level for unexpected optype", request.OpType())
	return nil
}

// TzSign this request with the provided Signer and Key
func (request *Request) TzSign(signer Signer, key *Key) (string, error) {
	msg := request.Hex()
	debugln("About to sign raw bytes hex.EncodeToString(bytes): ", hex.EncodeToString(msg))

	// Take the 256 bit (32 Byte) Blake2b Hash of the signing request
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
		return "", fmt.Errorf("request: b58 check encoded result was not correctly formatted: %v", signed)
	}

	return signed, nil
}
