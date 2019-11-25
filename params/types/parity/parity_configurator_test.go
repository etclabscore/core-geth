package parity

import (
	"testing"

	"github.com/ethereum/go-ethereum/common/math"
	paramtypes "github.com/ethereum/go-ethereum/params/types"
)

// This file contains a few unit tests for the parity-specific configuration interface.
// It does not contain integration tests, since this logic is covered by the test in convert_test.go,
// where specs are read, filled (converted), and verified equivalent.
//   Those tests cannot pass if the logic here is not sound.

func TestParityChainSpec_GetConsensusEngineType(t *testing.T) {
	spec := new(ParityChainSpec)

	if engine := (*spec).GetConsensusEngineType(); engine != paramtypes.ConsensusEngineT_Unknown {
		t.Error("unwanted engine type", engine)
	}

	spec.Engine.Ethash.Params.MinimumDifficulty = math.NewHexOrDecimal256(42)
	if engine := (*spec).GetConsensusEngineType(); engine != paramtypes.ConsensusEngineT_Ethash {
		t.Error("mismatch engine", engine)
	}
}