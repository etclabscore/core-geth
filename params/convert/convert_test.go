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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/params/types"
	"github.com/ethereum/go-ethereum/params/types/aleth"
	"github.com/ethereum/go-ethereum/params/types/common"
	"github.com/ethereum/go-ethereum/params/types/goethereum"
	"github.com/ethereum/go-ethereum/params/types/parity"
)

func mustOpenF(t *testing.T, fabbrev string, into interface{}) {
	b, err := ioutil.ReadFile(filepath.Join("testdata", fmt.Sprintf("stureby_%s.json", fabbrev)))
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
		"geth", "parity", "aleth",
	} {
		switch f {
		case "geth":
			c := &paramtypes.Genesis{}
			mustOpenF(t, f, c)
			if *c.Config.GetNetworkID() != 314158 {
				t.Errorf("networkid")
			}
		case "parity":
			p := &parity.ParityChainSpec{}
			mustOpenF(t, f, p)
			_, err := ParityConfigToMultiGethGenesis(p)
			if err != nil {
				t.Error(err)
			}
		case "aleth":
			a := &aleth.AlethGenesisSpec{}
			mustOpenF(t, f, a)
		}
	}
}

func TestConvert(t *testing.T) {
	spec := parity.ParityChainSpec{}
	mustOpenF(t, "parity", &spec)

	spec2 := parity.ParityChainSpec{}
	err := Convert(&spec, &spec2)
	if err != nil {
		t.Error(err)
	}

	if method, equal := Equal(reflect.TypeOf((*common.Configurator)(nil)), &spec, &spec2); !equal {
		t.Error("not equal", method)
	}
}

// TestConfiguratorImplementationsSatisfied tests that data types expected
// to fulfil certain interfaces do fill them.
func TestConfiguratorImplementationsSatisfied(t *testing.T) {
	for _, ty := range []interface{}{
		&parity.ParityChainSpec{},
	} {
		_ = ty.(common.Configurator)
	}

	for _, ty := range []interface{}{
		&goethereum.ChainConfig{},
		&paramtypes.ChainConfig{},
	} {
		_ = ty.(common.ChainConfigurator)
	}

	for _, ty := range []interface{}{
		&paramtypes.Genesis{},
	} {
		_ = ty.(common.GenesisBlocker)
	}
}

func TestCompatible(t *testing.T) {
	spec := &parity.ParityChainSpec{}
	fns, names := common.Transitions(spec)
	for i, fn := range fns {
		t.Log(names[i], fn())
	}
	t.Log(fns)
}
