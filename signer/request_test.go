package signer

import (
	"log"
	"math/big"
	"testing"
)

func TestParseTx(t *testing.T) {
	signingRequest, err := ParseRequest([]byte(testP256Tx.Operation))

	if err != nil {
		log.Println("Error parsing signing request: ", err.Error())
		t.Fail()
	}

	if signingRequest.OpType() != opTypeTx {
		log.Println("Error decoding the tx type")
		t.Fail()
	}
}

func TestParseEndorsement(t *testing.T) {
	signingRequest, err := ParseRequest([]byte(testEndorse.Operation))

	if err != nil {
		log.Println("Error parsing signing request: ", err.Error())
		t.Fail()
	}

	if signingRequest.OpType() != opTypeEndorsement {
		log.Println("Error decoding the op type.  Type: ", signingRequest.OpType())
		t.Fail()
	}
	level, _ := new(big.Int).SetString(testEndorse.Level, 10)
	if signingRequest.Level().Cmp(level) != 0 {
		log.Printf("Incorrectly parsed op level. Received %v, expecting %v\n", signingRequest.Level(), level)
		t.Fail()
	}
}

func TestParseBlock(t *testing.T) {
	signingRequest, err := ParseRequest([]byte(testBlock.Operation))

	if err != nil {
		log.Println("Error parsing signing request: ", err.Error())
		t.Fail()
	}

	if signingRequest.OpType() != opTypeBlock {
		log.Println("Error decoding the op type.  Type: ", signingRequest.OpType())
		t.Fail()
	}
	level, _ := new(big.Int).SetString(testBlock.Level, 10)
	if signingRequest.Level().Cmp(level) != 0 {
		log.Printf("Incorrectly parsed op level. Received %v, expecting %v\n", signingRequest.Level(), level)
		t.Fail()
	}
}
