package signer

import (
	"log"
	"math/big"
	"testing"
)

func TestParseTx(t *testing.T) {
	request, err := ParseRequest([]byte(testP256Tx.Operation))

	if err != nil {
		log.Println("Error parsing request: ", err.Error())
		t.Fail()
	}

	if request.OpType() != opTypeGeneric {
		log.Println("Error decoding the tx type")
		t.Fail()
	}
}

func testParse(t *testing.T, op testOperation, id string) {
	request, err := ParseRequest([]byte(op.Operation))

	if err != nil {
		log.Printf("%v: Error parsing signing request: %v\n", id, err.Error())
		t.Fail()
	}

	if request.OpType() != op.OpType {
		log.Printf("%v: Error decoding the op type.  Type: %v\n", id, request.OpType())
		t.Fail()
	}
	level, _ := new(big.Int).SetString(op.Level, 10)
	if request.Level().Cmp(level) != 0 {
		log.Printf("%v: Incorrectly parsed op level. Received %v, expecting %v\n", id, request.Level(), level)
		t.Fail()
	}

	if request.Level().Cmp(level) != 0 {
		log.Printf("%v: Incorrectly parsed op level. Received %v, expecting %v\n", id, request.Level(), level)
		t.Fail()
	}

	if request.ChainID() != op.ChainID {
		log.Printf("%v: Incorrectly parsed Chain ID. Received %v, expecting %v\n", id, request.ChainID(), op.ChainID)
		t.Fail()
	}

}

func TestParseEndorsement(t *testing.T) {
	testParse(t, testEndorse, "Endorse")
}

func TestParseBlock(t *testing.T) {
	testParse(t, testBlock, "Block")
}
