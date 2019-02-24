package watermark

import (
	"fmt"
	"math/big"
	"testing"
)

func assert(t *testing.T, condition bool, errorMessage string) {
	if !condition {
		fmt.Println("Test Failure: ", errorMessage)
		t.Fail()
	}
}

func SameLevel(t *testing.T) {
	wm := GetSessionWatermark()

	// Vars
	keyHash := "tz2..."
	chainIDMainnet := "NetXdQprcVkpaWU"
	chainIDAlphanet := "NetXgtSLGNJvNye"
	// OpTypes
	opTypeBlock := uint8(0x01)
	opTypeEndorsement := uint8(0x02)
	// Levels
	lvl1 := big.NewInt(1)
	lvl2 := big.NewInt(2)

	// Initial operation should be considered safe
	assert(t, wm.IsSafeToSign(keyHash, chainIDMainnet, opTypeBlock, lvl1), "Mainnent:Block:1 Should be safe to sign")
	assert(t, wm.IsSafeToSign(keyHash, chainIDMainnet, opTypeEndorsement, lvl1), "Mainnent:Endorsement:1 Should be safe to sign")
	assert(t, wm.IsSafeToSign(keyHash, chainIDAlphanet, opTypeBlock, lvl1), "Testnet:Block:1 Should be safe to sign")
	assert(t, wm.IsSafeToSign(keyHash, chainIDAlphanet, opTypeEndorsement, lvl1), "Testnet:Endorsement:1 Should be safe to sign")

	// Subsequent levels should be considered safe
	assert(t, wm.IsSafeToSign(keyHash, chainIDMainnet, opTypeBlock, lvl2), "Mainnent:Block:2 Should be safe to sign")
	assert(t, wm.IsSafeToSign(keyHash, chainIDMainnet, opTypeEndorsement, lvl2), "Mainnent:Endorsement:2 Should be safe to sign")
	assert(t, wm.IsSafeToSign(keyHash, chainIDAlphanet, opTypeBlock, lvl2), "Testnet:Block:2 Should be safe to sign")
	assert(t, wm.IsSafeToSign(keyHash, chainIDAlphanet, opTypeEndorsement, lvl2), "Testnet:Endorsement:2 Should be safe to sign")

	// The same level should fail
	assert(t, !wm.IsSafeToSign(keyHash, chainIDMainnet, opTypeBlock, lvl2), "Mainnent:Block:2 at the same level should fail")
	assert(t, !wm.IsSafeToSign(keyHash, chainIDMainnet, opTypeEndorsement, lvl2), "Mainnent:Endorsement:2 at the same level should fail")
	assert(t, !wm.IsSafeToSign(keyHash, chainIDAlphanet, opTypeBlock, lvl2), "Testnet:Block:2 at the same level should fail")
	assert(t, !wm.IsSafeToSign(keyHash, chainIDAlphanet, opTypeEndorsement, lvl2), "Testnet:Endorsement:2 at the same level should fail")

	// Lower levels should fail
	assert(t, !wm.IsSafeToSign(keyHash, chainIDMainnet, opTypeBlock, lvl1), "Mainnent:Block:1 at lower levels should fail")
	assert(t, !wm.IsSafeToSign(keyHash, chainIDMainnet, opTypeEndorsement, lvl1), "Mainnent:Endorsement:1 at lower levels should fail")
	assert(t, !wm.IsSafeToSign(keyHash, chainIDAlphanet, opTypeBlock, lvl1), "Testnet:Block:1 at lower levels should fail")
	assert(t, !wm.IsSafeToSign(keyHash, chainIDAlphanet, opTypeEndorsement, lvl1), "Testnet:Endorsement:1 at lower levels should fail")
}
