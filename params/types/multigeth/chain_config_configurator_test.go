package multigeth

import (
	"math"
	"math/big"
	"runtime"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params/types/ctypes"
	"github.com/ethereum/go-ethereum/params/vars"
)

func TestMultiGethChainConfig_SetDifficultyBombDelays(t *testing.T) {
	newMG := func() * MultiGethChainConfig {
		v := &MultiGethChainConfig{}
		v.MustSetConsensusEngineType(ctypes.ConsensusEngineT_Ethash)
		return v
	}
	byzaBlock := big.NewInt(4370000).Uint64()
	consBlock := big.NewInt(7280000).Uint64()
	muirBlock := big.NewInt(9200000).Uint64()

	max := uint64(math.MaxUint64)

	check := func(mg *MultiGethChainConfig, got, want uint64) {
		if got != want {
			t.Log(runtime.Caller(1))
			t.Log(runtime.Caller(2))
			t.Errorf("got: %d, want: %d", got, want)
			//t.Log(spew.Sprintln(mg.DifficultyBombDelaySchedule))
			for k, v := range mg.DifficultyBombDelaySchedule {
				t.Logf("%d: %d", k, v.Uint64())
			}
			t.Log("---")
		}
	}

	// Test set one (latest) value, eg for testing or new chains.
	mgSoloOrdered := newMG()
	mgSoloOrdered.SetEthashEIP2384Transition(&muirBlock)
	check(mgSoloOrdered, mgSoloOrdered.DifficultyBombDelaySchedule.SumValues(&muirBlock), vars.EIP2384DifficultyBombDelay.Uint64())
	check(mgSoloOrdered, mgSoloOrdered.DifficultyBombDelaySchedule.SumValues(&max), vars.EIP2384DifficultyBombDelay.Uint64())


	checkFinal := func(mg *MultiGethChainConfig) {
		check(mg, mg.DifficultyBombDelaySchedule.SumValues(&byzaBlock), vars.EIP649DifficultyBombDelay.Uint64())
		check(mg, mg.DifficultyBombDelaySchedule.SumValues(&consBlock), vars.EIP1234DifficultyBombDelay.Uint64())
		check(mg, mg.DifficultyBombDelaySchedule.SumValues(&muirBlock), vars.EIP2384DifficultyBombDelay.Uint64())
		check(mg, mg.DifficultyBombDelaySchedule.SumValues(&max), vars.EIP2384DifficultyBombDelay.Uint64())
	}

	// Test set fork values in chronological order.
	mgBasicOrdered := newMG()
	mgBasicOrdered.SetEthashEIP649Transition(&byzaBlock)
	check(mgBasicOrdered, mgBasicOrdered.DifficultyBombDelaySchedule.SumValues(&max), vars.EIP649DifficultyBombDelay.Uint64())
	mgBasicOrdered.SetEthashEIP1234Transition(&consBlock)
	mgBasicOrdered.SetEthashEIP2384Transition(&muirBlock)

	checkFinal(mgBasicOrdered)

	// Test set all features, unordered, 1
	mgBasicUnordered := newMG()
	mgBasicUnordered.SetEthashEIP2384Transition(&muirBlock)
	mgBasicUnordered.SetEthashEIP1234Transition(&consBlock)
	mgBasicUnordered.SetEthashEIP649Transition(&byzaBlock)

	checkFinal(mgBasicUnordered)

	// Same, but again, more.
	mgBasicUnordered2 := newMG()
	mgBasicUnordered2.SetEthashEIP2384Transition(&muirBlock)
	mgBasicUnordered2.SetEthashEIP649Transition(&byzaBlock)
	mgBasicUnordered2.SetEthashEIP1234Transition(&consBlock)

	checkFinal(mgBasicUnordered2)

	// Test set all features, unordered, and with edge cases, 2
	mgWildUnordered2 := newMG()
	mgWildUnordered2.SetEthashEIP2384Transition(&muirBlock)
	mgWildUnordered2.SetEthashEIP649Transition(&byzaBlock)
	mgWildUnordered2.SetEthashEIP1234Transition(&consBlock)

	// Set a dupe.
	mgWildUnordered2.SetEthashEIP649Transition(&byzaBlock)

	// Set a random.
	randoK := new(big.Int).Div(new(big.Int).Add(big.NewInt(int64(byzaBlock)), big.NewInt(int64(consBlock))), common.Big2).Uint64()
	randoV := new(big.Int).Div(new(big.Int).Add(vars.EIP649DifficultyBombDelay, vars.EIP1234DifficultyBombDelay), common.Big2)
	mgWildUnordered2.DifficultyBombDelaySchedule.SetValueTotalForHeight(&randoK, randoV)

	checkFinal(mgWildUnordered2)

	// Test repetitious set's.
	mgRepetitious := newMG()
	mgRepetitious.SetEthashEIP649Transition(&byzaBlock)
	mgRepetitious.SetEthashEIP1234Transition(&consBlock)
	mgRepetitious.SetEthashEIP2384Transition(&muirBlock)
	mgRepetitious.SetEthashEIP1234Transition(&consBlock)
	mgRepetitious.SetEthashEIP1234Transition(&consBlock)
	mgRepetitious.SetEthashEIP649Transition(&byzaBlock)
	mgRepetitious.SetEthashEIP649Transition(&byzaBlock)

	checkFinal(mgRepetitious)
}
