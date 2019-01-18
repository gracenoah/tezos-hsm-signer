package signer

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"golang.org/x/crypto/blake2b"
)

// Request parses and validates an arbitrary request
// to sign some message before passing to the Signer
type Request struct {
	raw       []byte
	parsedHex []byte
	opType    uint8
	level     *big.Int
}

// Types of Requests that can be signed
const (
	opTypeBlock       = 0x01
	opTypeEndorsement = 0x02
	opTypeTx          = 0x03
)

// ParseRequest parses a raw byte string into a meaningful signing payload
// if the request can be parsed and is valid
func ParseRequest(raw []byte) (*Request, error) {

	// Must begin and end with quotes
	parsed, err := stripQuotes(raw)
	if err != nil {
		return nil, err
	}

	// Must be valid hex chars
	parsedHex := make([]byte, hex.DecodedLen(len(parsed)))
	_, err = hex.Decode(parsedHex, parsed)
	if err != nil {
		return nil, err
	}

	request := Request{
		raw:       raw,
		parsedHex: parsedHex,
		opType:    parsedHex[0],
	}

	// Parse operation specific fields
	switch request.opType {
	case opTypeTx:
		debugln("Request is a Transaction")
	case opTypeBlock:
		request.level = new(big.Int).SetBytes(parsedHex[5:9])
		debugln("Request is a Block")
	case opTypeEndorsement:
		request.level = new(big.Int).SetBytes(parsedHex[len(parsedHex)-4:])
		debugln("Request is an Endorsement at level: ", request.level.String())
	default:
		return nil, fmt.Errorf("Unsupported Tx Type: %v", request.opType)
	}

	return &request, nil
}

// Hex returns a copy of the parsed hex bytes of the signing request
func (request *Request) Hex() []byte {
	hexCopy := make([]byte, len(request.parsedHex))
	copy(hexCopy, request.parsedHex)
	return hexCopy
}

// Level returns a copy of the level, if one was parsed from this request
func (request *Request) Level() *big.Int {
	if request.level != nil {
		return new(big.Int).Set(request.level)
	}
	return nil
}

// OpType of this tezos operation included in the signing request
func (request *Request) OpType() uint8 {
	return request.opType
}

// stripQuotes on either end of the byte array
func stripQuotes(input []byte) ([]byte, error) {
	firstQuote, secondQuote := getIndicesOfQuotes(input)

	if firstQuote == -1 || strings.TrimSpace(string(input[:firstQuote])) != "" {
		return nil, errors.New("request: A signing request must begin with a quote")
	}
	if secondQuote == -1 || strings.TrimSpace(string(input[secondQuote+1:])) != "" {
		return nil, errors.New("request: A signing request must end with a quote")
	}
	return input[firstQuote+1 : secondQuote], nil
}

// getIndices of first and second quotes in the given input byte slice
func getIndicesOfQuotes(input []byte) (int, int) {
	quoteByte := []byte("\"")[0]
	firstQuoteIndex := bytes.IndexByte(input, quoteByte)
	secondQuoteIndex := bytes.IndexByte(input[firstQuoteIndex+1:], quoteByte)
	return firstQuoteIndex, firstQuoteIndex + secondQuoteIndex + 1
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
	}
	debugln("Signed bytes hex.EncodeToString(bytes): ", hex.EncodeToString(signedMsg))

	// Get the correct signature prefix
	prefix, err := getSignaturePrefix(key)
	if err != nil {
		return "", err
	}

	// b58 Check Encode the result
	signed := b58CheckEncode(prefix, signedMsg)

	// Result should begin with:  p2sig, edsig, spsig1 or sig
	if signed[2:5] != "sig" && signed[0:3] != "sig" {
		return "", fmt.Errorf("request: b58 check encoded result was not correctly formatted: %v", signed)
	}

	return signed, nil
}
