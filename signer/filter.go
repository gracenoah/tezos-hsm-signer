package signer

import (
	"fmt"
	"log"
	"math/big"
	"time"
)

// OperationFilter controls what operations will be signed
type OperationFilter struct {
	EnableGeneric        bool
	EnableTx             bool
	EnableVoting         bool
	TxWhitelistAddresses []string
	TxDailyMax           *big.Int

	// Keep track of daily max withdrawals
	dailyTxMaxKey     string
	dailyTxMaxCounter *big.Int
}

// IsAllowed by this filter?
func (filter *OperationFilter) IsAllowed(op *Operation) bool {
	switch op.MagicByte() {
	case opMagicByteBlock, opMagicByteEndorsement:
		return true
	case opMagicByteGeneric:
		generic := GetGenericOperation(op)
		if filter.EnableGeneric {
			return true
		}
		if filter.EnableTx && generic.Kind() == opKindTransaction {
			return filter.isWhitelisted(generic) && filter.authorizeTxAmount(generic.TransactionValue())
		}
		if filter.EnableVoting && (generic.Kind() == opKindBallot) {
			return true
		}
		if filter.EnableVoting && (generic.Kind() == opKindProposals) {
			return true
		}
		return false
	default:
		return false
	}
}

// Is this address whitelisted? Returns true if whitelistising is disabled
func (filter *OperationFilter) isWhitelisted(generic *GenericOperation) bool {
	if filter.TxWhitelistAddresses == nil {
		return true
	}
	for _, pkh := range filter.TxWhitelistAddresses {
		if generic.TransactionDestination() == PubkeyHashToByteString(pkh) {
			return true
		}
	}
	log.Println("[WARN] Address is not whitelisted")
	return false
}

// AuthorizeTxAmount for withdrawal.  Fails if this amount would push us
// over the daily limit in XTZ.  Returns true if limits are disabled
func (filter *OperationFilter) authorizeTxAmount(value *big.Int) bool {
	if filter.TxDailyMax == nil {
		return true
	}

	now := time.Now()
	key := fmt.Sprintf("%v-%v", now.Year(), now.YearDay())
	// Reset the counter when we're in a new day
	if filter.dailyTxMaxKey != key {
		filter.dailyTxMaxKey = key
		filter.dailyTxMaxCounter = new(big.Int).SetInt64(0)
	}
	filter.dailyTxMaxCounter.Add(filter.dailyTxMaxCounter, value)
	return filter.dailyTxMaxCounter.Cmp(filter.TxDailyMax) == -1
}
