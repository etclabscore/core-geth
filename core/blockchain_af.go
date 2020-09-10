package core

import (
	"errors"
	"fmt"
	"math"
	"math/big"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
)

// errReorgFinality represents an error caused by artificial finality mechanisms.
var errReorgFinality = errors.New("finality-enforced invalid new chain")

// EnableArtificialFinality enables and disable artificial finality features for the blockchain.
// Currently toggled features include:
// - ECBP11355-MESS: modified exponential subject scoring
//
// This level of activation works BELOW the chain configuration for any of the
// potential features. eg. If ECBP11355 is not activated at the chain config x block number,
// then calling bc.EnableArtificialFinality(true) will be a noop.
// The method is idempotent.
func (bc *BlockChain) EnableArtificialFinality(enable bool, logValues ...interface{}) {
	// Store enable/disable value regardless of config activation.
	var statusLog string
	if enable {
		statusLog = "Enabled"
		atomic.StoreInt32(&bc.artificialFinalityEnabled, 1)
	} else {
		statusLog = "Disabled"
		atomic.StoreInt32(&bc.artificialFinalityEnabled, 0)
	}
	configActivated := bc.chainConfig.IsEnabled(bc.chainConfig.GetECBP11355Transition, bc.CurrentHeader().Number)
	logFn := log.Debug // Deactivated
	if configActivated && enable {
		logFn = log.Info // Activated and enabled
	} else if configActivated && !enable {
		logFn = log.Warn // Activated and disabled
	}
	logFn(fmt.Sprintf("%s artificial finality features", statusLog), logValues...)
}

// IsArtificialFinalityEnabled returns the status of the blockchain's artificial
// finality feature setting.
// This status is agnostic of feature activation by chain configuration.
func (bc *BlockChain) IsArtificialFinalityEnabled() bool {
	return atomic.LoadInt32(&bc.artificialFinalityEnabled) == 1
}

// ecpb11355 implements the "MESS" artificial finality mechanism
// "Modified Exponential Subjective Scoring" used to prefer known chain segments
// over later-to-come counterparts, especially proposed segments stretching far into the past.
func (bc *BlockChain) ecbp11355(commonAncestor, current, proposed *types.Header) error {
	commonAncestorTD := bc.GetTd(commonAncestor.Hash(), commonAncestor.Number.Uint64())

	proposedParentTD := bc.GetTd(proposed.ParentHash, proposed.Number.Uint64()-1)
	proposedTD := new(big.Int).Add(proposed.Difficulty, proposedParentTD)

	localTD := bc.GetTd(current.Hash(), current.Number.Uint64())

	tdRatio, _ := new(big.Float).Quo(
		new(big.Float).SetInt(new(big.Int).Sub(proposedTD, commonAncestorTD)),
		new(big.Float).SetInt(new(big.Int).Sub(localTD, commonAncestorTD)),
	).Float64()

	antiGravity := math.Pow(1.0001, float64(proposed.Time-commonAncestor.Time))

	if tdRatio < antiGravity {
		// Using "b/a" here as "'B' chain vs. 'A' chain", where A is original (current), and B is proposed (new).
		return fmt.Errorf("%w: ECPB11355-MESS: td.b/a(%0.4f) < antigravity(%0.4f)", errReorgFinality, tdRatio, antiGravity)
	}
	return nil
}
