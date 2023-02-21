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
	"github.com/ethereum/go-ethereum/params/types/parity"
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

func TestConvert(t *testing.T) {
	spec := goethereum.ChainConfig{}
	mustReadTestdataTo(t, "geth", &spec)

	spec2 := goethereum.ChainConfig{}
	err := confp.Convert(&spec, &spec2)
	if err != nil {
		t.Error(err)
	}

	if diffs := confp.Equal(reflect.TypeOf((*ctypes.Configurator)(nil)), &spec, &spec2); len(diffs) != 0 {
		for _, diff := range diffs {
			t.Error("not equal", diff.Field, diff.A, diff.B)
		}
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
		&parity.ParityChainSpec{},
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
		&parity.ParityChainSpec{},
	} {
		_ = ty.(ctypes.Configurator)
	}

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
	spec := &parity.ParityChainSpec{}
	fns, names := confp.Transitions(spec)
	for i, fn := range fns {
		t.Log(names[i], fn())
	}
	t.Log(fns)
}
