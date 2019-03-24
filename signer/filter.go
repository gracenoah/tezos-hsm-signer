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
	dailyTxMaxKey       string
	dailyTxMaxCounter   *big.Int
	dailyVoteMaxKey     string
	dailyVoteMaxCounter *big.Int
}

// No more than 100 XTZ spent per day voting
const maxDailyVotingValue = 100 * 1000000

// IsAllowed by this filter?
func (filter *OperationFilter) IsAllowed(op *Operation) bool {
	if op.Watermark() == opWatermarkGeneric {
		generic := GetGenericOperation(op)
		if filter.EnableGeneric {
			return true
		}
		if filter.EnableTx && generic.Kind() == opKindTransaction {
			return filter.isWhitelisted(generic) && filter.authorizeTxAmount(generic.TransactionValue())
		}
		if filter.EnableVoting && (generic.Kind() == opKindBallot) {
			return filter.authorizeVoteAmount(generic.BallotValue())
		}
		if filter.EnableVoting && (generic.Kind() == opKindProposals) {
			return filter.authorizeVoteAmount(generic.ProposalValue())
		}
		return false
	}
	return true

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

// AuthorizeVoteAmount for brodcast votes.  Fails if this amount would push us
// over 100 XTZ daily limit.
func (filter *OperationFilter) authorizeVoteAmount(value *big.Int) bool {

	now := time.Now()
	key := fmt.Sprintf("%v-%v", now.Year(), now.YearDay())
	// Reset the counter when we're in a new day
	if filter.dailyVoteMaxKey != key {
		filter.dailyVoteMaxKey = key
		filter.dailyVoteMaxCounter = new(big.Int).SetInt64(0)
	}
	filter.dailyVoteMaxCounter.Add(filter.dailyVoteMaxCounter, value)
	return filter.dailyVoteMaxCounter.Cmp(new(big.Int).SetInt64(maxDailyVotingValue)) == -1
}
