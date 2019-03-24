package signer

import (
	"encoding/hex"
	"log"
	"math/big"
)

// GenericOperation parses an operation with a generic watermark byte
type GenericOperation struct {
	hex []byte
}

// Kind of different types of generic operations
const (
	// opKindActivateAccount           = 0x00
	// opKindProposals                 = 0x00
	// opKindBallot                    = 0x00
	// opKindReveal                    = 0x00
	opKindTransaction = 0x08
	// opKindOrigination               = 0x00
	// opKindDelegation                = 0x00
)

// GetGenericOperation to parse specific Generic fields
func GetGenericOperation(op *Operation) *GenericOperation {
	if op.Watermark() != opWatermarkGeneric {
		return nil
	}
	return &GenericOperation{
		hex: op.Hex(),
	}
}

// Kind of the generic operation
func (op *GenericOperation) Kind() uint8 {
	return op.hex[33]
}

// TransactionSource address that funds are being moved from
func (op *GenericOperation) TransactionSource() string {
	if op.Kind() != opKindTransaction {
		return ""
	}
	return hex.EncodeToString(op.hex[36:56])
}

// TransactionFee that's being paid along with this tx
func (op *GenericOperation) TransactionFee() *big.Int {
	if op.Kind() != opKindTransaction {
		return nil
	}
	return op.parseSerializedNumberOffset(0)
}

// TransactionCounter ensuring idempotency of this tx
func (op *GenericOperation) TransactionCounter() *big.Int {
	if op.Kind() != opKindTransaction {
		return nil
	}
	return op.parseSerializedNumberOffset(1)
}

// TransactionGasLimit of this tx
func (op *GenericOperation) TransactionGasLimit() *big.Int {
	if op.Kind() != opKindTransaction {
		return nil
	}
	return op.parseSerializedNumberOffset(2)
}

// TransactionStorageLimit of this tx
func (op *GenericOperation) TransactionStorageLimit() *big.Int {
	if op.Kind() != opKindTransaction {
		return nil
	}
	return op.parseSerializedNumberOffset(3)
}

// TransactionAmount that's moving with this tx
func (op *GenericOperation) TransactionAmount() *big.Int {
	if op.Kind() != opKindTransaction {
		return nil
	}
	return op.parseSerializedNumberOffset(4)
}

// TransactionDestination address we're sending funds to
func (op *GenericOperation) TransactionDestination() string {
	if op.Kind() != opKindTransaction {
		return ""
	}
	// Verify these indices align with the end_index of transaction amount
	numberIndex := 56
	for i := 0; i <= 4; i++ {
		_, numberIndex = op.parseSerializedNumber(numberIndex)
	}

	start := len(op.hex) - 21
	end := len(op.hex) - 1
	if start-numberIndex != 2 {
		log.Println("Warning: Incorrect offset between numbers and destination.  Unsure where we're sending.")
		return ""
	}
	return hex.EncodeToString(op.hex[start:end])
}

// Private funcs to parse sequentially serialized numbers in the operation's hex
func (op *GenericOperation) parseSerializedNumberOffset(offset int) *big.Int {
	var num *big.Int
	// Numbers always begin at this index
	index := 56
	for i := 0; i <= offset; i++ {
		num, index = op.parseSerializedNumber(index)
	}
	return num
}

// Parse a numbers starting at the provided index.  Return the number and
// the index of the next byte in the operation.
func (op *GenericOperation) parseSerializedNumber(startIndex int) (*big.Int, int) {
	b := op.hex[startIndex]
	nextIndex := startIndex + 1

	base := new(big.Int).SetInt64(int64(0))
	top := new(big.Int).SetInt64(int64(b))
	top.Mod(top, new(big.Int).SetInt64(0x80))
	if b >= 0x80 {
		var result *big.Int
		result, nextIndex = op.parseSerializedNumber(nextIndex)
		base.Mul(new(big.Int).SetInt64(0x80), result)
	}
	return top.Add(top, base), nextIndex
}
