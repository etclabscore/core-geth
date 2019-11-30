// Copyright 2019 The multi-geth Authors
// This file is part of the multi-geth library.
//
// The multi-geth library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The multi-geth library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the multi-geth library. If not, see <http://www.gnu.org/licenses/>.


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

//func TestParityConfigToMultiGethGenesis(t *testing.T) {
//	var gen1, gen2 *paramtypes.Genesis
//
//	for p, gen := range chainSpecEquivs {
//		gen1 = gen
//		gen1.Config.EIP150Hash = common.Hash{}
//		b, err := ioutil.ReadFile(asSpecFilePath(p))
//		if err != nil {
//			t.Fatal(err)
//		}
//		paritySpec := parity.ParityChainSpec{}
//		err = json.UnmarshalChainConfigurator(b, &paritySpec)
//		if err != nil {
//			t.Fatal(err)
//		}
//
//		gen2, err = ParityConfigToMultiGethGenesis(&paritySpec)
//		if err != nil {
//			t.Fatal(err)
//		}
//
//		for _, bl := range []int64{
//			0, 10000, 40000, 50000, 100000,
//			2000000, 4000000, 6000000, 8000000, 10000000,
//		} {
//			b := big.NewInt(bl)
//			fns := []func(b *big.Int) bool{
//				gen1.Config.IsEIP2F, gen2.Config.IsEIP2F,
//				gen1.Config.IsEIP100F, gen2.Config.IsEIP100F,
//				gen1.Config.IsEIP213F, gen2.Config.IsEIP213F,
//				gen1.Config.IsEIP1052F, gen2.Config.IsEIP1052F,
//				gen1.Config.IsEIP140F, gen2.Config.IsEIP140F,
//				gen1.Config.IsEIP161F, gen2.Config.IsEIP161F,
//			}
//			for i, f := range fns {
//				if i == 0 || i%2 == 0 {
//					continue
//				}
//				if (f(b) && !fns[i-1](b)) || (!f(b) && fns[i-1](b)) {
//					t.Errorf("%d mismatch", i)
//				}
//			}
//		}
//	}
//}

