package convert

import (
	"math/big"
	"path/filepath"
	"testing"

	"github.com/ethereum/go-ethereum/params"
	paramtypes "github.com/ethereum/go-ethereum/params/types"
	"github.com/ethereum/go-ethereum/params/types/goethereum"
)

func asSpecFilePath(name string) string {
	return filepath.Join("..", "parity.json.d", name)
}

var chainSpecEquivs = map[string]*paramtypes.Genesis{
	"classic.json":    params.DefaultClassicGenesisBlock(),
	"foundation.json": params.DefaultGenesisBlock(),
}

func TestBlockConfig(t *testing.T) {
	frontierCC := &goethereum.ChainConfig{
			ChainID: big.NewInt(1),
			Ethash:  new(goethereum.EthashConfig),
	}
	genesis := params.DefaultGenesisBlock()
	genesis.Config = frontierCC
	paritySpec, err := NewParityChainSpec("frontier", genesis, []string{})
	if err != nil {
		t.Fatal(err)
	}
	parityHomestead := paritySpec.Engine.Ethash.Params.HomesteadTransition
	if parityHomestead != nil && *parityHomestead >= 0 {
		t.Errorf("nonnil parity homestead")
	}
}
