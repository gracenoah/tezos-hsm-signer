package watermark

import (
	"math/big"
	"strconv"
	"sync"
)

// SessionWatermark stores the last-signed level in memory
type SessionWatermark struct {
	watermarkEntries []*watermarkEntry
	mux              sync.Mutex
}

// GetSessionWatermark returns a new in-memory watermark manager
func GetSessionWatermark() *SessionWatermark {
	// Initialize with an empty watermark entry
	return &SessionWatermark{
		watermarkEntries: []*watermarkEntry{},
		mux:              sync.Mutex{},
	}
}

// IsSafeToSign returns true if the provided (key, chainID, opType) tuple has
// not yet been signed at this or greater levels
func (mw *SessionWatermark) IsSafeToSign(keyHash string, chainID string, opType uint8, level *big.Int) bool {
	mw.mux.Lock()
	defer mw.mux.Unlock()

	sOpType := strconv.Itoa(int(opType))

	for _, entry := range mw.watermarkEntries {
		if entry.KeyHash == keyHash && entry.ChainID == chainID && entry.OpType == sOpType {
			iLevel, ok := new(big.Int).SetString(entry.Level, 10)
			if !ok {
				return false
			}
			// If the new level is > last level, update level and return true
			if level.Cmp(iLevel) == 1 {
				entry.Level = level.String()
				return true
			}
			return false
		}
	}
	mw.watermarkEntries = append(mw.watermarkEntries, &watermarkEntry{
		KeyHash: keyHash,
		ChainID: chainID,
		OpType:  strconv.Itoa(int(opType)),
		Level:   level.String(),
	})
	return true
}
