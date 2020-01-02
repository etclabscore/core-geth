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

package tconvert

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/ethereum/go-ethereum/params/types/aleth"
	"github.com/ethereum/go-ethereum/params/types/genesisT"
	"github.com/go-test/deep"
)

// FIXME(meowsbits): Requires implementing ChainConfigurator interface for Aleth data type.
// Tests the go-ethereum to Aleth chainspec conversion for the Stureby testnet.
func TestAlethSturebyConverter(t *testing.T) {

	// Read GETH genesis type.
	blob, err := ioutil.ReadFile(filepath.Join("..", "testdata", "stureby_geth.json"))
	if err != nil {
		t.Fatalf("could not read file: %v", err)
	}
	var genesis genesisT.Genesis
	if err := json.Unmarshal(blob, &genesis); err != nil {
		t.Fatalf("failed parsing genesis: %v", err)
	}

	// Convert read-in GETH genesis to Aleth spec.
	convertedSpec, err := NewAlethGenesisSpec("stureby", &genesis)
	if err != nil {
		t.Fatalf("failed creating chainspec: %v", err)
	}

	// Read the aleth JSON spec.
	expBlob, err := ioutil.ReadFile(filepath.Join("..", "testdata", "stureby_aleth.json"))
	if err != nil {
		t.Fatalf("could not read file: %v", err)
	}
	readSpec := &aleth.AlethGenesisSpec{}
	if err := json.Unmarshal(expBlob, readSpec); err != nil {
		t.Fatalf("failed parsing genesis: %v", err)
	}

	// Compare the read-in Aleth spec versus the converted-to Aleth spec.
	if diffs := deep.Equal(convertedSpec, readSpec); len(diffs) != 0 {
		t.Errorf("error: chainspec mismatch")

		t.Log("__Differences: (convertedTo_Spec != readIn_Spec)")
		for _, d := range diffs {
			t.Log(d)
		}

		t.Log("__ReadIn spec (Geth genesis):")
		bmm, _ := json.MarshalIndent(genesis, "", "    ")
		t.Log(string(bmm))
		t.Log()

		t.Log("__ConvertedTo spec (Aleth):")
		//t.Log(spew.Sprint(convertedSpec))
		bm, _ := json.MarshalIndent(convertedSpec, "", "    ")
		t.Log(string(bm))
	}
}
