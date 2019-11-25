package parity

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
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

// TestParityChainSpec_UnmarshalJSON shows that the data structure
// is valid for all included (whitelisty) parity json specs.
func TestParityChainSpec_Configurator_UnmarshalJSON(t *testing.T) {
	err := filepath.Walk(filepath.Join("..", "..", "parity.json.d"), func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(info.Name()) != ".json" {
			return nil
		}
		t.Run(info.Name(), func(t *testing.T) {
			b, err := ioutil.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			spec := ParityChainSpec{}
			err = json.Unmarshal(b, &spec)
			if err != nil {
				t.Errorf("%s, err: %v", info.Name(), err)
			}

			f2 := uint64(42)
			if err := spec.SetChainID(&f2); err != nil {
				t.Error("set chainid", err)
			}
		})
		return nil
	})

	if err != nil {
		t.Fatal(err)
	}
}
