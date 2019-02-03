package signer

import (
	"encoding/hex"
	"log"
	"reflect"
	"testing"
)

func TestLeftPad(t *testing.T) {
	// A string leading with 0's will be truncted by big.int.Bytes()
	bytes, _ := hex.DecodeString("00fed44b952c94d6be989f684c02b9253c9cf306eb8bc00648d49cfe5da074f90ef2845961911669b8c6d69c35ff24a0f42888bd19e6467f653735049a34702c")
	modded := StrictECModN(&Key{}, bytes)

	//
	if !reflect.DeepEqual(bytes, modded) {
		log.Print("")
		t.Fail()
	}
}
