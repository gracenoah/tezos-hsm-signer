package signer

import (
	"log"
	"math/big"
)

const (
	curveEd25519 = iota + 1
	curveSecp256k1
	curveNistP256
	curveUnknown
)

var (
	// From: http://www.secg.org/SEC2-Ver-1.0.pdf section 2.7.1
	secp256k1Order, _  = new(big.Int).SetString("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141", 16)
	secp256k1HalfOrder = new(big.Int).Rsh(secp256k1Order, 1)
	// From: http://www.secg.org/SEC2-Ver-1.0.pdf section 2.7.2
	nistP256r1Order, _  = new(big.Int).SetString("FFFFFFFF00000000FFFFFFFFFFFFFFFFBCE6FAADA7179E84F3B9CAC2FC632551", 16)
	nistP256r1HalfOrder = new(big.Int).Rsh(nistP256r1Order, 1)
)

// StrictECModN ensures strict compliance with the EC spec by returning S mod n
// for the appropriate keys curve.
//
// Details:
//   Step #6 of the ECDSA algorithm [x] defines an `S` value mod n[0],
//   but most signers (OpenSSL, SoftHSM, YubiHSM) don't return a strict modulo.
//   This variability was exploited with transaction malleability in Bitcoin,
//   leading to BIP#62.  BIP#62 Rule #5[1] requires that signatures return a
//   strict S = ... mod n which this function forces implemented in btcd here [2]
//   [0]: https://en.wikipedia.org/wiki/Elliptic_Curve_Digital_Signature_Algorithm
//   [1]: https://github.com/bitcoin/bips/blob/master/bip-0062.mediawiki#new-rules
//   [2]: https://github.com/btcsuite/btcd/blob/master/btcec/signature.go#L49
func StrictECModN(key *Key, sig []byte) []byte {
	R := new(big.Int).SetBytes(sig[0:32])
	S := new(big.Int).SetBytes(sig[32:])
	if key.Curve() == curveSecp256k1 {
		if S.Cmp(secp256k1HalfOrder) == 1 {
			S.Sub(secp256k1Order, S)
		}
	} else if key.Curve() == curveNistP256 {
		if S.Cmp(nistP256r1HalfOrder) == 1 {
			S.Sub(nistP256r1Order, S)
		}
	}

	// Leftpad to 32 bytes
	rBytes := leftPad(R.Bytes(), 32)
	sBytes := leftPad(S.Bytes(), 32)
	return append(rBytes, sBytes...)
}

// leftPad a byte slice with 0x00 to fit into the specified size
func leftPad(bytes []byte, length int) []byte {
	if len(bytes) > length {
		log.Fatalf("Fatal: Tried to leftpad a %v byte string down to %v bytes\n", len(bytes), length)
	}
	return append(make([]byte, length-len(bytes)), bytes...)
}
