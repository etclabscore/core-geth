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
	"bytes"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/params/confp"
	"github.com/ethereum/go-ethereum/params/types/coregeth"
	"github.com/ethereum/go-ethereum/params/types/ctypes"
	"github.com/ethereum/go-ethereum/params/types/genesisT"
	"github.com/ethereum/go-ethereum/params/types/goethereum"
)

func mustReadTestdataTo(t *testing.T, fabbrev string, into interface{}) {
	b, err := os.ReadFile(filepath.Join("..", "testdata", fmt.Sprintf("%s_foundation.json", fabbrev)))
	if err != nil {
		t.Fatal(err)
	}
	err = json.Unmarshal(b, &into)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_UnmarshalJSON(t *testing.T) {
	for _, f := range []string{
		"geth", "coregeth",
	} {
		switch f {
		case "geth":
			c := &genesisT.Genesis{}
			mustReadTestdataTo(t, f, c)
			if c.Config.GetChainID().Cmp(big.NewInt(1)) != 0 {
				t.Errorf("go-ethereum: wrong chainid")
			}
			if _, ok := c.Config.(*goethereum.ChainConfig); !ok {
				t.Errorf("go-ethereum: wrong type")
			}
		case "coregeth":
			c := &genesisT.Genesis{}
			mustReadTestdataTo(t, f, c)
			if c.Config.GetChainID().Cmp(big.NewInt(1)) != 0 {
				t.Errorf("core-geth: wrong chainid")
			}
			if _, ok := c.Config.(*coregeth.CoreGethChainConfig); !ok {
				t.Errorf("core-geth: wrong type")
			}
		}
	}
}

func testCrush(t *testing.T, aType string, a, b ctypes.ChainConfigurator) {
	ag, bg := &genesisT.Genesis{Config: a}, &genesisT.Genesis{Config: b}
	mustReadTestdataTo(t, aType, ag)
	err := confp.Crush(bg, ag, true)
	if err != nil {
		t.Fatal(err)
	}
	didError := false
	if diffs := confp.Equal(reflect.TypeOf((*ctypes.Configurator)(nil)), ag, bg); len(diffs) != 0 {
		didError = true
		for _, diff := range diffs {
			t.Error("not equal", diff.Field, diff.A, diff.B)
		}
	}
	if didError {
		// Log the json encoding of the b struct for debugging.
		bb, err := json.MarshalIndent(bg, "", "  ")
		if err != nil {
			t.Fatal(err)
		}
		t.Log(string(bb))
		t.Fatal("not equal")
	}
}

func TestCrush(t *testing.T) {
	testCrush(t, "geth", &goethereum.ChainConfig{}, &goethereum.ChainConfig{})
	testCrush(t, "coregeth", &coregeth.CoreGethChainConfig{}, &coregeth.CoreGethChainConfig{})
	testCrush(t, "geth", &goethereum.ChainConfig{}, &coregeth.CoreGethChainConfig{})
	testCrush(t, "coregeth", &coregeth.CoreGethChainConfig{}, &goethereum.ChainConfig{})
}

func TestCrush_SkipZeroValues(t *testing.T) {
	// Fill a destination config with testdata.
	destination := &genesisT.Genesis{Config: &coregeth.CoreGethChainConfig{}}
	mustReadTestdataTo(t, "coregeth", destination)

	// Cache some values from the destination to check they don't get mutated accidentally;
	// 2 non-nil values and 1 nil value.
	destinationEIP7FBlock := destination.Config.GetEIP7Transition()
	if destinationEIP7FBlock == nil {
		t.Fatal("destinationEIP7FBlock is nil")
	}
	destinationEIP100bFBlock := destination.Config.GetEthashEIP100BTransition()
	if destinationEIP100bFBlock == nil {
		t.Fatal("destinationEIP100bFBlock is nil")
	}
	destinationECIP1041FBlock := destination.Config.GetEthashECIP1041Transition()
	if destinationECIP1041FBlock != nil {
		t.Fatal("destinationECIP1041FBlock is NOT nil")
	}

	// Construct some source config to crush into the destination.
	source := &coregeth.CoreGethChainConfig{
		ChainID:     big.NewInt(69),
		EIP155Block: big.NewInt(42),
	}

	// Crush the source into the destination, skipping zero-values.
	err := confp.Crush(destination.Config, source, false)
	if err != nil {
		t.Fatal(err)
	}

	// Check that the destination config was modified as expected.
	if v := destination.Config.GetChainID(); v == nil || v.Cmp(source.ChainID) != 0 {
		t.Errorf("mismatch chainid, want: %v, got: %v", source.ChainID, v)
	}
	if v := destination.Config.GetEIP155Transition(); v == nil || *v != source.EIP155Block.Uint64() {
		t.Errorf("mismatch eip155 block, want: %v, got: %v", source.EIP155Block.Uint64(), v)
	}

	// Using the cached values from the destination,
	// check that the destination was not modified where it should not have been.
	if v := destination.Config.GetEIP7Transition(); v == nil || *v != *destinationEIP7FBlock {
		t.Errorf("mismatch eip7 block, want: %v, got: %v", destinationEIP7FBlock, v)
	}
	if v := destination.Config.GetEthashEIP100BTransition(); v == nil || *v != *destinationEIP100bFBlock {
		t.Errorf("mismatch eip100b block, want: %v, got: %v", destinationEIP100bFBlock, v)
	}
	if v := destination.Config.GetEthashECIP1041Transition(); v != nil {
		t.Errorf("mismatch ecip1041 block, want: nil, got: %v", v)
	}
}

func TestIdentical(t *testing.T) {
	methods := []string{
		"ChainID",
		"NetworkID",
	}
	configs := []ctypes.ChainConfigurator{
		&coregeth.CoreGethChainConfig{},
		&goethereum.ChainConfig{},
		&coregeth.CoreGethChainConfig{}, // Complete combination test set.
	}
	for i := range configs {
		if i == 0 {
			continue
		}
		f42, f43 := uint64(43), big.NewInt(43)
		configs[i-1].SetNetworkID(&f42)
		configs[i].SetNetworkID(&f42)
		configs[i-1].SetChainID(f43)
		configs[i].SetChainID(f43)
		if !confp.Identical(configs[i-1], configs[i], methods) {
			t.Errorf("nonident")
		}
		f24 := uint64(24)
		configs[i-1].SetNetworkID(&f24)
		if confp.Identical(configs[i-1], configs[i], methods) {
			t.Error(i, "ident")
		}
	}
}

// TestConfiguratorImplementationsSatisfied tests that data types expected
// to fulfil certain interfaces do fill them.
func TestConfiguratorImplementationsSatisfied(t *testing.T) {
	for _, ty := range []interface{}{
		&goethereum.ChainConfig{},
		&coregeth.CoreGethChainConfig{},
	} {
		_ = ty.(ctypes.ChainConfigurator)
	}

	for _, ty := range []interface{}{
		&genesisT.Genesis{},
	} {
		_ = ty.(ctypes.GenesisBlocker)
	}
}

func TestCompatible(t *testing.T) {
	spec := &goethereum.ChainConfig{}
	fns, names := confp.Transitions(spec)
	for i, fn := range fns {
		t.Log(names[i], fn())
	}
	t.Log(fns)
}

func isJSONEqual(a, b interface{}) bool {
	aa, err := json.Marshal(a)
	if err != nil {
		return false
	}
	bb, err := json.Marshal(b)
	if err != nil {
		return false
	}
	return bytes.Equal(aa, bb)
}

func TestCloneChainConfigurator(t *testing.T) {
	for _, f := range []string{
		"geth", "coregeth",
	} {
		switch f {
		case "geth":
			c := &genesisT.Genesis{}
			mustReadTestdataTo(t, f, c)
			if c.Config.GetChainID().Cmp(big.NewInt(1)) != 0 {
				t.Errorf("go-ethereum: wrong chainid")
			}
			if _, ok := c.Config.(*goethereum.ChainConfig); !ok {
				t.Errorf("go-ethereum: wrong type")
			}
			cloned, err := confp.CloneChainConfigurator(c.Config)
			if err != nil {
				t.Fatal(err)
			}
			if diffs := confp.Equal(reflect.TypeOf((*ctypes.ChainConfigurator)(nil)), c.Config, cloned); len(diffs) != 0 {
				for _, diff := range diffs {
					t.Error("not equal", diff.Field, diff.A, diff.B)
				}
			}
			if !isJSONEqual(c.Config, cloned) {
				t.Error("not json equal")
			}

		case "coregeth":
			c := &genesisT.Genesis{}
			mustReadTestdataTo(t, f, c)
			if c.Config.GetChainID().Cmp(big.NewInt(1)) != 0 {
				t.Errorf("core-geth: wrong chainid")
			}
			if _, ok := c.Config.(*coregeth.CoreGethChainConfig); !ok {
				t.Errorf("core-geth: wrong type")
			}
			cloned, err := confp.CloneChainConfigurator(c.Config)
			if err != nil {
				t.Fatal(err)
			}
			if diffs := confp.Equal(reflect.TypeOf((*ctypes.ChainConfigurator)(nil)), c.Config, cloned); len(diffs) != 0 {
				for _, diff := range diffs {
					t.Error("not equal", diff.Field, diff.A, diff.B)
				}
			}
			if !isJSONEqual(c.Config, cloned) {
				t.Error("not json equal")
			}
		}
	}
}
