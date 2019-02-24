package watermark

import (
	"math/big"
)

// IgnoreWatermark always retruns true
type IgnoreWatermark struct{}

// GetIgnoreWatermark that will always return true
func GetIgnoreWatermark() Watermark {
	return &IgnoreWatermark{}
}

// IsSafeToSign is always true when we're ignoring the watermark
func (mw *IgnoreWatermark) IsSafeToSign(keyHash string, chainID string, opType uint8, level *big.Int) bool {
	return true
}
