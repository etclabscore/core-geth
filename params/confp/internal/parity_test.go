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

package convert_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/params/confp/tconvert"
	"github.com/ethereum/go-ethereum/params/types/ctypes"
	"github.com/ethereum/go-ethereum/params/types/goethereum"
)

func TestBlockConfig(t *testing.T) {
	frontierCC := &goethereum.ChainConfig{
		ChainID: big.NewInt(1),
		Ethash:  new(ctypes.EthashConfig),
	}
	genesis := params.DefaultGenesisBlock()
	genesis.Config = frontierCC
	paritySpec, err := tconvert.NewParityChainSpec("frontier", genesis, []string{})
	if err != nil {
		t.Fatal(err)
	}
	parityHomestead := paritySpec.Engine.Ethash.Params.HomesteadTransition
	if parityHomestead != nil && *parityHomestead >= 0 {
		t.Errorf("nonnil parity homestead")
	}
}
