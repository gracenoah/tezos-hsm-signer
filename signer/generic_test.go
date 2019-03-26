package signer

import (
	"encoding/hex"
	"log"
	"math/big"
	"testing"
)

type testGenericOperation struct {
	Name         string
	Operation    string
	Kind         uint8
	Source       string
	Fee          *big.Int
	Counter      *big.Int
	GasLimit     *big.Int
	StorageLimit *big.Int
	Amount       *big.Int
	Destination  string
}

func TestParseKind(t *testing.T) {
	var op *Operation
	op, _ = ParseOperation([]byte(testP256Tx.Operation))
	generic := GetGenericOperation(op)
	if generic.Kind() != opKindTransaction {
		log.Println("Tx was not parsed as a generic transaction")
		t.Fail()
	}
	op, _ = ParseOperation([]byte(testSecp256k1Tx.Operation))
	if generic.Kind() != opKindTransaction {
		log.Println("Tx was not parsed as a generic transaction")
		t.Fail()
	}
}

func testParseGenericOperation(t *testing.T, test *testGenericOperation) {
	op, _ := ParseOperation([]byte(test.Operation))
	generic := GetGenericOperation(op)
	// Verify Each Field
	if generic.Kind() != test.Kind {
		log.Printf("[Generic Test - %v] Kind mismatch. Expected %v but received %v\n", test.Name, test.Kind, generic.Kind())
		t.Fail()
	}

	if generic.TransactionSource() != PubkeyHashToByteString(test.Source) {
		log.Printf("[Generic Test - %v] Source mismatch. Expected %v but received %v\n", test.Name, test.Source, generic.TransactionSource())
		t.Fail()
	}
	if generic.TransactionFee().Cmp(test.Fee) != 0 {
		log.Printf("[Generic Test - %v] Fee mismatch. Expected %v but received %v\n", test.Name, test.Fee, generic.TransactionFee())
		t.Fail()
	}
	if generic.TransactionCounter().Cmp(test.Counter) != 0 {
		log.Printf("[Generic Test - %v] Counter mismatch. Expected %v but received %v\n", test.Name, test.Counter, generic.TransactionCounter())
		t.Fail()
	}
	if generic.TransactionGasLimit().Cmp(test.GasLimit) != 0 {
		log.Printf("[Generic Test - %v] GasLimit mismatch. Expected %v but received %v\n", test.Name, test.GasLimit, generic.TransactionGasLimit())
		t.Fail()
	}
	if generic.TransactionStorageLimit().Cmp(test.StorageLimit) != 0 {
		log.Printf("[Generic Test - %v] StorageLimit mismatch. Expected %v but received %v\n", test.Name, test.StorageLimit, generic.TransactionStorageLimit())
		t.Fail()
	}
	if generic.TransactionAmount().Cmp(test.Amount) != 0 {
		log.Printf("[Generic Test - %v] Amount mismatch. Expected %v but received %v\n", test.Name, test.Amount, generic.TransactionAmount())
		t.Fail()
	}
	if generic.TransactionDestination() != PubkeyHashToByteString(test.Destination) {
		log.Printf("[Generic Test - %v] Destination mismatch. Expected %v but received %v\n", test.Name, test.Destination, generic.TransactionDestination())
		t.Fail()
	}
}

func TestParseTransactions(t *testing.T) {
	testParseGenericOperation(t, &testGenericOperation{
		Name:         "Small Values",
		Kind:         opKindTransaction,
		Operation:    "\"03ce69c5713dac3537254e7be59759cf59c15abd530d10501ccf9028a5786314cf08000002298c03ed7d454a101eb7022bc95f7e5f41ac780102030405000202298c03ed7d454a101eb7022bc95f7e5f41ac7800\"",
		Source:       "tz1KqTpEZ7Yob7QbPE4Hy4Wo8fHG8LhKxZSx",
		Fee:          new(big.Int).SetInt64(1),
		Counter:      new(big.Int).SetInt64(2),
		GasLimit:     new(big.Int).SetInt64(3),
		StorageLimit: new(big.Int).SetInt64(4),
		Amount:       new(big.Int).SetInt64(5),
		Destination:  "tz1KqTpEZ7Yob7QbPE4Hy4Wo8fHG8LhKxZSx",
	})
	testParseGenericOperation(t, &testGenericOperation{
		Name: "Large Values",
		Kind: opKindTransaction,

		Operation:    "\"03ce69c5713dac3537254e7be59759cf59c15abd530d10501ccf9028a5786314cf08000002298c03ed7d454a101eb7022bc95f7e5f41ac787f80018101ffff03808004000202298c03ed7d454a101eb7022bc95f7e5f41ac7800\"",
		Source:       "tz1KqTpEZ7Yob7QbPE4Hy4Wo8fHG8LhKxZSx",
		Fee:          new(big.Int).SetInt64(127),
		Counter:      new(big.Int).SetInt64(128),
		GasLimit:     new(big.Int).SetInt64(129),
		StorageLimit: new(big.Int).SetInt64(65535),
		Amount:       new(big.Int).SetInt64(65536),
		Destination:  "tz1KqTpEZ7Yob7QbPE4Hy4Wo8fHG8LhKxZSx",
	})
	testParseGenericOperation(t, &testGenericOperation{
		Name:         "Zero Values",
		Kind:         opKindTransaction,
		Operation:    "\"03ce69c5713dac3537254e7be59759cf59c15abd530d10501ccf9028a5786314cf08000002298c03ed7d454a101eb7022bc95f7e5f41ac786404020000000002298c03ed7d454a101eb7022bc95f7e5f41ac7800\"",
		Source:       "tz1KqTpEZ7Yob7QbPE4Hy4Wo8fHG8LhKxZSx",
		Fee:          new(big.Int).SetInt64(100),
		Counter:      new(big.Int).SetInt64(4),
		GasLimit:     new(big.Int).SetInt64(2),
		StorageLimit: new(big.Int).SetInt64(0),
		Amount:       new(big.Int).SetInt64(0),
		Destination:  "tz1KqTpEZ7Yob7QbPE4Hy4Wo8fHG8LhKxZSx",
	})
	testParseGenericOperation(t, &testGenericOperation{
		Name:         "KT Address",
		Kind:         opKindTransaction,
		Operation:    "\"037072fa916732ed788ab030ca81714956fea286521fc0ead7ad533868e0f0308408000154f5d8f71ce18f9f05bb885a4120e64c667bc1b4f809ca69d84f00c0843d016e7c23cc06c7b0743256f65e34d5b0f7c91e4eb20000\"",
		Source:       "tz2G4TwEbsdFrJmApAxJ1vdQGmADnBp95n9m",
		Fee:          new(big.Int).SetInt64(1272),
		Counter:      new(big.Int).SetInt64(13514),
		GasLimit:     new(big.Int).SetInt64(10200),
		StorageLimit: new(big.Int).SetInt64(0),
		Amount:       new(big.Int).SetInt64(1000000),
		Destination:  "KT1JexcFezMnUAaWmvUGY99jwTA4jcKiUgFp",
	})

}

func TestParseProposal(t *testing.T) {
	op, _ := ParseOperation([]byte("\"03ce69c5713dac3537254e7be59759cf59c15abd530d10501ccf9028a5786314cf05008fb5cea62d147c696afd9a93dbce962f4c8a9c910000000a00000020ab22e46e7872aa13e366e455bb4f5dbede856ab0864e1da7e122554579ee71f8\""))
	generic := GetGenericOperation(op)
	// Verify Each Field
	if generic.Kind() != opKindProposals {
		log.Printf("[Proposal Test] Kind mismatch. Expected %v but received %v\n", opKindProposals, generic.Kind())
		t.Fail()
	}
}
func TestParseBallot(t *testing.T) {
	op, _ := ParseOperation([]byte("\"03ce69c5713dac3537254e7be59759cf59c15abd530d10501ccf9028a5786314cf0600531ab5764a29f77c5d40b80a5da45c84468f08a10000000bab22e46e7872aa13e366e455bb4f5dbede856ab0864e1da7e122554579ee71f800\""))
	generic := GetGenericOperation(op)
	// Verify Each Field
	if generic.Kind() != opKindBallot {
		log.Printf("[Proposal Test] Kind mismatch. Expected %v but received %v\n", opKindProposals, generic.Kind())
		t.Fail()
	}
}

func testParseBytes(t *testing.T, bytes string, expect int64) {
	var op GenericOperation
	hex, _ := hex.DecodeString(bytes)
	op = GenericOperation{hex: hex}

	num, _ := op.parseSerializedNumber(0)
	if num.Int64() != expect {
		log.Printf("Expecting %v, received %v\n", expect, num.String())
		t.Fail()
	}
}
func TestParseBytes(t *testing.T) {
	testParseBytes(t, "8001", 128)

	testParseBytes(t, "ff7f", 16383)
	testParseBytes(t, "808001", 16384)
	testParseBytes(t, "818001", 16385)

	testParseBytes(t, "ffff01", 32767)
	testParseBytes(t, "808002", 32768)
	testParseBytes(t, "818002", 32769)

	testParseBytes(t, "ff8002", 32895)
	testParseBytes(t, "808102", 32896)

	testParseBytes(t, "ffff03", 65535)
	testParseBytes(t, "808004", 65536)
}
