package parity

import (
	"testing"

	"github.com/ethereum/go-ethereum/common/math"
	paramtypes "github.com/ethereum/go-ethereum/params/types"
)

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
