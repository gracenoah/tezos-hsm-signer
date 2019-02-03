package signer

import (
	"log"
	"math/big"
	"strings"
	"sync"
)

// A Key identifies a key preloaded in your HSM
type Key struct {
	Name             string `yaml:"Name"`
	PublicKeyHash    string `yaml:"PublicKeyHash"`
	PublicKey        string `yaml:"PublicKey"`
	HsmSlot          uint   `yaml:"HsmSlot"`
	HsmLabel         string `yaml:"HsmLabel"`
	LastEndorseLevel string `yaml:"LastEndorseLevel"`
	LastBakeLevel    string `yaml:"LastBakeLevel"`
	mux              sync.Mutex
}

// Curve represented by this key
func (key *Key) Curve() int {
	if strings.HasPrefix(key.PublicKeyHash, "tz1") {
		return curveEd25519
	} else if strings.HasPrefix(key.PublicKeyHash, "tz2") {
		return curveSecp256k1
	} else if strings.HasPrefix(key.PublicKeyHash, "tz3") {
		return curveNistP256
	}
	return curveUnknown
}

// IsECDSA Curve or EdDSA Curve
func (key *Key) IsECDSA() bool {
	return key.Curve() == curveNistP256 || key.Curve() == curveSecp256k1
}

// IsSafeToSign at this level
func (key *Key) IsSafeToSign(operation uint8, level *big.Int) bool {
	if operation == opTypeGeneric {
		return true
	} else if operation == opTypeEndorsement {
		return getBigInt(key.LastEndorseLevel).Cmp(level) == -1
	} else if operation == opTypeBlock {
		return getBigInt(key.LastEndorseLevel).Cmp(level) == -1
	}
	log.Println("Warning: Attempting to sign unrecognized optype: ", operation)
	return false
}

// Lock a mutex around this key
func (key *Key) Lock() {
	key.mux.Lock()
}

// Unlock a mutex around this key
func (key *Key) Unlock() {
	key.mux.Unlock()
}

// getBigInt from the provided string
func getBigInt(str string) *big.Int {
	newInt := big.Int{}
	newInt.SetString(str, 10)
	return &newInt
}
