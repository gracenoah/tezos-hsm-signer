package signer

import (
	"io/ioutil"
	"log"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

// A Key identifies a key preloaded in your HSM
type Key struct {
	Name          string `yaml:"Name"`
	PublicKeyHash string `yaml:"PublicKeyHash"`
	PublicKey     string `yaml:"PublicKey"`
	HsmSlot       uint   `yaml:"HsmSlot"`
	HsmLabel      string `yaml:"HsmLabel"`
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

func loadKeyFile(keyfile string) []Key {
	keys := []Key{}

	yamlFile, err := ioutil.ReadFile(keyfile)
	if err != nil {
		log.Fatalln("Unable to read file: " + keyfile)
	}
	err = yaml.Unmarshal(yamlFile, &keys)
	if err != nil {
		log.Fatalln("Unable to parse yaml file: " + keyfile)
	}
	return keys
}
