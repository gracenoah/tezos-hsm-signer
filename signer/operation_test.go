package signer

import (
	"log"
	"math/big"
	"testing"
)

func TestParseTx(t *testing.T) {
	op, err := ParseOperation([]byte(testP256Tx.Operation))

	if err != nil {
		log.Println("Error parsing operation: ", err.Error())
		t.Fail()
	}

	if op.Type() != opTypeGeneric {
		log.Println("Error decoding the tx type")
		t.Fail()
	}
}

func testParse(t *testing.T, test testOperation, id string) {
	op, err := ParseOperation([]byte(test.Operation))

	if err != nil {
		log.Printf("%v: Error parsing operation: %v\n", id, err.Error())
		t.Fail()
	}

	if op.Type() != test.OpType {
		log.Printf("%v: Error decoding the op type.  Type: %v\n", id, op.Type())
		t.Fail()
	}
	level, _ := new(big.Int).SetString(test.Level, 10)
	if op.Level().Cmp(level) != 0 {
		log.Printf("%v: Incorrectly parsed op level. Received %v, expecting %v\n", id, op.Level(), level)
		t.Fail()
	}

	if op.Level().Cmp(level) != 0 {
		log.Printf("%v: Incorrectly parsed op level. Received %v, expecting %v\n", id, op.Level(), level)
		t.Fail()
	}

	if op.ChainID() != test.ChainID {
		log.Printf("%v: Incorrectly parsed Chain ID. Received %v, expecting %v\n", id, op.ChainID(), test.ChainID)
		t.Fail()
	}
}

func TestParseEndorsement(t *testing.T) {
	testParse(t, testEndorse, "Endorse")

}

func TestParseBlock(t *testing.T) {
	testParse(t, testBlock, "Block")
}
