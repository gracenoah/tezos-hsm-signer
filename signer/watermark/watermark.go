package watermark

import (
	"math/big"
)

// Watermark stores the last (key, level, chainID) tuple that has been signed
// and fails if you attempt to sign the same or lesser level for that tuple
type Watermark interface {
	// IsSafeToSign returns true if the provided (key, chainID, opType) tuple has
	// not yet been signed at this or greater levels
	IsSafeToSign(keyHash string, chainID string, opType uint8, level *big.Int) bool
}

// watermarkEntry stores our locks
type watermarkEntry struct {
	KeyHash string `yaml:"Key"`
	ChainID string `yaml:"ChainID"`
	OpType  string `yaml:"OpType"`
	Level   string `yaml:"Level"`
}
