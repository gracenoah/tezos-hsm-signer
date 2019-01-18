package signer

import (
	"log"
	"math/big"
	"testing"
)

func TestStripQuotes(t *testing.T) {
	var result []byte
	var err error

	expectError := func(err error, result []byte, msg string) {
		if err == nil || result != nil {
			log.Printf("%v, result: %v, error: %v", msg, string(result), err)
			t.Fail()
		}
	}
	expectSuccess := func(err error, result []byte, msg string) {
		if err != nil || string(result) != "testing" {
			log.Printf("%v, result: %v, error: %v", msg, string(result), err)
			t.Fail()
		}
	}

	result, err = stripQuotes([]byte(""))
	expectError(err, result, "Expect error on empty string")
	result, err = stripQuotes([]byte("\"testing"))
	expectError(err, result, "Expect error with only one quote")
	result, err = stripQuotes([]byte("\"testing\"abc"))
	expectError(err, result, "No non-space characters allowed after quotes")
	result, err = stripQuotes([]byte("abc\"testing\""))
	expectError(err, result, "No non-space characters allowed before quotes")

	result, err = stripQuotes([]byte("\"testing\""))
	expectSuccess(err, result, "Result should match text within quotes")
	result, err = stripQuotes([]byte("  \n \"testing\""))
	expectSuccess(err, result, "Whitespace allowed before quotes")
	result, err = stripQuotes([]byte("\"testing\"\n"))
	expectSuccess(err, result, "Whitespace allowed after quotes")
}

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
		log.Println("Error decoding the tx type.  Type: ", signingRequest.OpType())
		t.Fail()
	}
	level, _ := new(big.Int).SetString(testEndorse.Level, 10)
	if signingRequest.Level().Cmp(level) != 0 {
		log.Printf("Incorrectly parsed Tx level. Received %v, expecting %v\n", signingRequest.Level(), level)
		t.Fail()
	}
}
