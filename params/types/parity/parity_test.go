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

package parity

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/davecgh/go-spew/spew"
	math2 "github.com/ethereum/go-ethereum/common/math"
)

var exampleAccountWithBuiltinA = []byte(`
			{
				"name": "modexp",
				"activate_at": "0x85d9a0",
				"pricing": {
					"modexp": {
						"divisor": 20
					}
				}
			}
		`)

var exampleAccountWithBuiltinB = []byte(`
			{
				"name": "alt_bn128_add",
				"pricing": {
					"0x85d9a0": {
						"price": { "alt_bn128_const_operations": { "price": 500 }}
					},
					"0x7fffffffffffff": {
						"price": { "alt_bn128_const_operations": { "price": 150 }}
					}
				}
			}
		`)

func TestParityBuiltinType(t *testing.T) {
	b := ParityChainSpecBuiltin{}
	err := json.Unmarshal(exampleAccountWithBuiltinA, &b)
	if err != nil {
		t.Fatal(err)
	}
	if b.Pricing.Pricing == nil {
		t.Errorf("pricing nil")
	}
	if b.Pricing.Pricing.ModExp.Divisor != 20 {
		t.Errorf("wrong price")
	}

	err = json.Unmarshal(exampleAccountWithBuiltinB, &b)
	if err != nil {
		t.Fatal(err)
	}
	if b.Pricing.Map == nil {
		t.Errorf("no map")
	}
	mi := math2.NewHexOrDecimal256(0x85d9a0)
	if len(b.Pricing.Map) == 0 {
		t.Fatal("0 map")
	}
	for k, v := range b.Pricing.Map {
		if k.ToInt().Cmp(mi.ToInt()) == 0 {
			if v.AltBnConstOperation.Price != 500 {
				t.Errorf("wrong map: %v", spew.Sdump(b.Pricing))
			}
		}
	}
}

// TestParityChainSpec_UnmarshalJSON shows that the data structure
// is valid for all included (whitelisty) parity json specs.
func TestParityChainSpec_UnmarshalJSON(t *testing.T) {
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
		})
		return nil
	})

	if err != nil {
		t.Fatal(err)
	}
}
