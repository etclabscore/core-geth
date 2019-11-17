// Copyright 2018 The go-ethereum Authors
// This file is part of go-ethereum.
//
// go-ethereum is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// go-ethereum is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with go-ethereum. If not, see <http://www.gnu.org/licenses/>.

package chainspec

import (
	"encoding/json"
	"io/ioutil"
	"math/big"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/params"
	"github.com/go-test/deep"
)

func TestBlockConfig(t *testing.T) {
	frontierCC := &params.ChainConfig{
		ChainID: big.NewInt(1),
		Ethash: new(params.EthashConfig),
	}
	genesis := core.DefaultGenesisBlock()
	genesis.Config = frontierCC
	paritySpec, err := NewParityChainSpec("frontier", genesis, []string{})
	if err != nil {
		t.Fatal(err)
	}
	parityHomestead := paritySpec.Engine.Ethash.Params.HomesteadTransition
	if parityHomestead >= 0 {
		t.Errorf("nonnil parity homestead")
	}
}

// Tests the go-ethereum to Aleth chainspec conversion for the Stureby testnet.
func TestAlethSturebyConverter(t *testing.T) {
	blob, err := ioutil.ReadFile("testdata/stureby_geth.json")
	if err != nil {
		t.Fatalf("could not read file: %v", err)
	}
	var genesis core.Genesis
	if err := json.Unmarshal(blob, &genesis); err != nil {
		t.Fatalf("failed parsing genesis: %v", err)
	}
	spec, err := NewAlethGenesisSpec("stureby", &genesis)
	if err != nil {
		t.Fatalf("failed creating chainspec: %v", err)
	}

	expBlob, err := ioutil.ReadFile("testdata/stureby_aleth.json")
	if err != nil {
		t.Fatalf("could not read file: %v", err)
	}
	expspec := &AlethGenesisSpec{}
	if err := json.Unmarshal(expBlob, expspec); err != nil {
		t.Fatalf("failed parsing genesis: %v", err)
	}
	// Compare the read-in SomeSpec vs. the generated SomeSpec
	if diffs := deep.Equal(expspec, spec); len(diffs) != 0 {
		t.Errorf("chainspec mismatch")
		for _, d := range diffs {
			t.Log(d)
		}
		spew.Dump(spec)
		bm, _ := json.MarshalIndent(spec, "", "    ")
		t.Log(string(bm))
	}
}

// Tests the go-ethereum to Parity chainspec conversion for the Stureby testnet.
func TestParitySturebyConverter(t *testing.T) {
	// Read a native genesis json config
	blob, err := ioutil.ReadFile("testdata/stureby_geth.json")
	if err != nil {
		t.Fatalf("could not read file: %v", err)
	}
	var genesis core.Genesis
	if err := json.Unmarshal(blob, &genesis); err != nil {
		t.Fatalf("failed parsing genesis: %v", err)
	}

	// Marshall this genesis to SomeSpec
	spec, err := NewParityChainSpec("stureby", &genesis, []string{})
	if err != nil {
		t.Fatalf("failed creating chainspec: %v", err)
	}

	// Read comparator json config
	expBlob, err := ioutil.ReadFile("testdata/stureby_parity.json")
	if err != nil {
		t.Fatalf("could not read file: %v", err)
	}
	// Marshal this config to SomeSpec
	expspec := &ParityChainSpec{}
	if err := json.Unmarshal(expBlob, expspec); err != nil {
		t.Fatalf("failed parsing genesis: %v", err)
	}
	expspec.Nodes = []string{}

	// Compare the read-in SomeSpec vs. the generated SomeSpec
	if diffs := deep.Equal(expspec, spec); len(diffs) != 0 {
		t.Errorf("chainspec mismatch")
		for _, d := range diffs {
			t.Log(d)
		}
		spew.Dump(spec)
		bm, _ := json.MarshalIndent(spec, "", "    ")
		t.Log(string(bm))
	}
}
