package vm

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/params/confp"
	"github.com/ethereum/go-ethereum/params/vars"
)

// TestECIP1086 shows that SLOAD constant gas is never bumped to 800 on Kotti the testnet config.
// Despite EIP2200 being installed and it's specification calling for 800 gas, this was incorrectly
// implemented in several clients, and thus, via ECIP1086, made into a "canonical mistake," which is tested
// here.
// This tests shows that although EIP2200 is implemented for a window (sloppyWindow), the gas cost for SLOAD does not change.
func TestECIP1086_Kotti(t *testing.T) {
	kc := params.KottiChainConfig

	sloppyWindowStart := uint64(2_058_191)
	sloppyWindowEnd := uint64(2_208_203)

	// Since "sloppy" EIP2200 did not actually modify SLOAD gas (and ECIP1086 makes this canonical),
	// we want to show that SLOAD=200 gas is always used on Kotti whether EIP2200 is enabled or, afterwards, disabled.
	wantSloadGasAlways := vars.SloadGasEIP150

	if wantSloadGasAlways != vars.NetSstoreDirtyGas {
		// https://corepaper.org/ethereum/fork/istanbul/#aztlan-fix
		/*
			Design Failure

			Aztlan decided to include EIP-2200 for SSTORE net gas metering.
			However, it did not thoroughly consider the facts that EIP-2200 is designed with the full context
			of Istanbul hard fork, whose pre-conditions that Aztlan does not fulfill.

			Net gas metering has a family of EIP specifications, including EIP-1283, EIP-1706 and EIP-2200.
			The specifications contain a parameter of "dirty gas", which is charged when a storage value has
			already been modified in the same transaction context. This value, as in the design intention,
			is expected to always equal to the value of SLOAD gas cost.
			In EIP-2200, this "dirty gas" is set to be 800, so as to accommodate EIP-1884’s gas cost change of SLOAD from 200 to 800.
			However, Aztlan does not include EIP-1884 but only applied EIP-2200,
			resulting in inconsistency in the EIP design intention, and may lead to unknown gas cost issues.

			This has also led to confusions in implementations. For example, a contributor previously merged an
			incorrect version of Aztlan hard fork into Parity Ethereum.
			This, if left unspotted, will lead to consensus split on the whole Ethereum Classic network.

			Recommendation: Remove EIP-2200 from the applied specification list, and add EIP-1283 with EIP-1706 instead.
		*/
		t.Fatal("'dirty gas' parameter should be same as SLOAD cost")
	}

	for _, f := range confp.Forks(kc) {
		for _, ft := range []uint64{f - 1, f, f + 1} {

			jt := instructionSetForConfig(kc, new(big.Int).SetUint64(ft))

			// Yes, this logic could be simplified, but I want to show the window/post-window no-op verbosely.
			if ft >= sloppyWindowStart && ft < sloppyWindowEnd {
				if jt[SLOAD].constantGas != wantSloadGasAlways {
					t.Error(ft, "bad gas", jt[SLOAD].constantGas)
				}

			} else if ft >= sloppyWindowEnd {
				// EIP1283 and EIP1706 are activated, but neither modify the SLOAD constant gas price.
				if jt[SLOAD].constantGas != wantSloadGasAlways {
					t.Error(ft, "bad gas 2", jt[SLOAD].constantGas)
				}
			}
		}
	}
}

// TestECIP1086_ETCMainnet is a variant of TestECIP1086_Kotti,
// showing that SLOAD gas is never 800.
func TestECIP1086_ETCMainnet(t *testing.T) {
	etc := params.ClassicChainConfig

	for _, f := range confp.Forks(etc) {
		for _, ft := range []uint64{f - 1, f, f + 1} {
			head := new(big.Int).SetUint64(ft)
			jt := instructionSetForConfig(etc, head)
			if jt[SLOAD].constantGas == vars.SloadGasEIP2200 && etc.IsEnabled(etc.GetEIP2200DisableTransition, head) {
				t.Errorf("wrong gas")
			}
		}
	}
}
