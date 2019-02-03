package signer

import (
	"io/ioutil"
	"log"
	"math/big"
	"sync"

	yaml "gopkg.in/yaml.v2"
)

// KeyManager manages a rest of keys
type KeyManager struct {
	keyfile string
	keys    []Key
}

// loadKeyManager loads a keys yaml file
func loadKeyManager(keyfile string) KeyManager {
	keys := []Key{}

	yamlFile, err := ioutil.ReadFile(keyfile)
	if err != nil {
		log.Fatalln("Unable to read file: " + keyfile)
	}
	err = yaml.Unmarshal(yamlFile, &keys)
	if err != nil {
		log.Fatalln("Unable to parse yaml file: " + keyfile)
	}
	return KeyManager{
		keys:    keys,
		keyfile: keyfile,
	}
}

// GetKeyFromHash returns a full public key from a provided public key *hash*
func (keyManager *KeyManager) GetKeyFromHash(publicKeyHash string) *Key {
	for i := 0; i < len(keyManager.keys); i++ {
		if keyManager.keys[i].PublicKeyHash == publicKeyHash {
			return &keyManager.keys[i]
		}
	}
	return nil
}

// SetLastSignedLevel to this new level
func (keyManager *KeyManager) SetLastSignedLevel(key *Key, operation uint8, level *big.Int) {
	if !key.IsSafeToSign(operation, level) {
		log.Fatal("Signing must always be signing at a higher level.  Exiting")
	}

	if operation == opTypeGeneric {
		return
	} else if operation == opTypeEndorsement {
		key.LastEndorseLevel = level.String()
	} else if operation == opTypeBlock {
		key.LastBakeLevel = level.String()
	}

	// If a file is set, write back to it
	if len(keyManager.keyfile) > 0 {
		lock := sync.Mutex{}
		lock.Lock()

		bytes, err := yaml.Marshal(keyManager.keys)
		if err != nil {
			log.Fatal("Unable to marshall keys: " + keyManager.keyfile)
		}

		err = ioutil.WriteFile(keyManager.keyfile, bytes, 0644)
		if err != nil {
			log.Fatal("Unable to write lockfile: " + keyManager.keyfile)
		}
		lock.Unlock()
	}
}
