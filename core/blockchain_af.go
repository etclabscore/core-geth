package core

import (
	"errors"
	"fmt"
	"math"
	"math/big"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
)

// errReorgFinality represents an error caused by artificial finality mechanisms.
var errReorgFinality = errors.New("finality-enforced invalid new chain")

// EnableArtificialFinality enables and disable artificial finality features for the blockchain.
// Currently toggled features include:
// - ECBP1100-MESS: modified exponential subject scoring
//
// This level of activation works BELOW the chain configuration for any of the
// potential features. eg. If ECBP1100 is not activated at the chain config x block number,
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
	if !bc.chainConfig.IsEnabled(bc.chainConfig.GetECBP1100Transition, bc.CurrentHeader().Number) {
		// Don't log anything if the config hasn't enabled it yet.
		return
	}
	logFn := log.Warn // Deactivated and enabled
	if enable {
		logFn = log.Info // Activated and enabled
	}
	logFn(fmt.Sprintf("%s artificial finality features", statusLog), logValues...)
}

// IsArtificialFinalityEnabled returns the status of the blockchain's artificial
// finality feature setting.
// This status is agnostic of feature activation by chain configuration.
func (bc *BlockChain) IsArtificialFinalityEnabled() bool {
	return atomic.LoadInt32(&bc.artificialFinalityEnabled) == 1
}

// getTDRatio is a helper function returning the total difficulty ratio of
// proposed over current chain segments.
func (bc *BlockChain) getTDRatio(commonAncestor, current, proposed *types.Header) float64 {
	// Get the total difficulty ratio of the proposed chain segment over the existing one.
	commonAncestorTD := bc.GetTd(commonAncestor.Hash(), commonAncestor.Number.Uint64())

	proposedParentTD := bc.GetTd(proposed.ParentHash, proposed.Number.Uint64()-1)
	proposedTD := new(big.Int).Add(proposed.Difficulty, proposedParentTD)

	localTD := bc.GetTd(current.Hash(), current.Number.Uint64())

	tdRatio, _ := new(big.Float).Quo(
		new(big.Float).SetInt(new(big.Int).Sub(proposedTD, commonAncestorTD)),
		new(big.Float).SetInt(new(big.Int).Sub(localTD, commonAncestorTD)),
	).Float64()
	return tdRatio
}

// ecbp1100 implements the "MESS" artificial finality mechanism
// "Modified Exponential Subjective Scoring" used to prefer known chain segments
// over later-to-come counterparts, especially proposed segments stretching far into the past.
func (bc *BlockChain) ecbp1100(commonAncestor, current, proposed *types.Header) error {

	tdRatio := bc.getTDRatio(commonAncestor, current, proposed)

	// Time span diff.
	// The minimum value is 1.
	x := float64(proposed.Time - commonAncestor.Time)

	// Commented now is a potential way to "soften" the acceptance while
	// still avoiding discrete acceptance boundaries. In the case that ecbp1100 introduces
	// unacceptable network inefficiency, this (or something similar) may be an option.
	// // Accept with diminishing probability in the case of equivalent total difficulty.
	// // Remember that the equivalent total difficulty case has ALREADY
	// // passed one coin toss.
	// if tdRatio == 1 && rand.Float64() < (1/x) {
	// 	return nil
	// }

	antiGravity := ecbp1100AGSinusoidalA(x)

	if tdRatio < antiGravity {
		// Using "b/a" here as "'B' chain vs. 'A' chain", where A is original (current), and B is proposed (new).
		return fmt.Errorf(`%w: ECBP1100-MESS 🔒 status=rejected age=%v blocks=%d td.B/A=%0.6f < antigravity=%0.6f`,
			errReorgFinality,
			common.PrettyAge(time.Unix(int64(commonAncestor.Time), 0)), proposed.Number.Uint64()-commonAncestor.Number.Uint64(),
			tdRatio, antiGravity,
		)
	}
	log.Info("ECBP1100-MESS 🔓",
		"status", "accepted",
		"age", common.PrettyAge(time.Unix(int64(commonAncestor.Time), 0)),
		"blocks", proposed.Number.Uint64()-commonAncestor.Number.Uint64(),
		"td.B/A", tdRatio,
		"antigravity", antiGravity,
	)
	return nil
}

/*
ecbp1100AGSinusoidalA is a sinusoidal function.

OPTION 3: Yet slower takeoff, yet steeper eventual ascent. Has a differentiable ceiling transition.
h(x)=15 sin((x+12000 π)/(8000))+15+1

*/
func ecbp1100AGSinusoidalA(x float64) (antiGravity float64) {
	ampl := float64(15)   // amplitude
	pDiv := float64(8000) // period divisor
	phaseShift := math.Pi * (pDiv * 1.5)
	peakX := math.Pi * pDiv // x value of first sin peak where x > 0
	if x > peakX {
		// Cause the x value to limit to the x value of the first peak of the sin wave (ceiling).
		x = peakX
	}
	return (ampl * math.Sin((x+phaseShift)/pDiv)) + ampl + 1
}

/*
ecbp1100AGExpB is an exponential function with x as a base (and rationalized exponent).

OPTION 2: Slightly slower takeoff, steeper eventual ascent
g(x)=x^(x*0.00002)
*/
func ecbp1100AGExpB(x float64) (antiGravity float64) {
	return math.Pow(x, x*0.00002)
}

/*
ecbp1100AGExpA is an exponential function with x as exponent.

This was (one of?) Vitalik's "original" specs:
> 1.0001 ** (number of seconds between when S1 was received and when S2 was received)
- https://bitcointalk.org/index.php?topic=865169.msg16349234#msg16349234
> gravity(B') = gravity(B) * 0.99 ^ n
- https://blog.ethereum.org/2014/11/25/proof-stake-learned-love-weak-subjectivity/

OPTION 1 (Original ESS)
f(x)=1.0001^(x)
*/
func ecbp1100AGExpA(x float64) (antiGravity float64) {
	return math.Pow(1.0001, x)
}
