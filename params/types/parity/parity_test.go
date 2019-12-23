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
	"github.com/ethereum/go-ethereum/common"
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

// TestParityChainSpec_GetPrecompile checks lexographical unmarshaling for maps which can
// have duplicate keys when unmarshaling builtin pricing.
func TestParityChainSpec_GetPrecompile(t *testing.T) {
	pspec := &ParityChainSpec{}
	err := json.Unmarshal([]byte(`
{
  "name": "Byzantium (Test)",
  "engine": {
    "Ethash": {
      "params": {
        "minimumDifficulty": "0x020000",
        "difficultyBoundDivisor": "0x0800",
        "durationLimit": "0x0d",
        "blockReward": "0x29A2241AF62C0000",
        "homesteadTransition": "0x0",
        "eip100bTransition": "0x0",
        "difficultyBombDelays": {
          "0": 3000000
        }
      }
    }
  },
  "params": {
    "gasLimitBoundDivisor": "0x0400",
    "registrar" : "0xc6d9d2cd449a754c494264e1809c50e34d64562b",
    "accountStartNonce": "0x00",
    "maximumExtraDataSize": "0x20",
    "minGasLimit": "0x1388",
    "networkID" : "0x1",
    "maxCodeSize": 24576,
    "maxCodeSizeTransition": "0x0",
    "eip150Transition": "0x0",
    "eip160Transition": "0x0",
    "eip161abcTransition": "0x0",
    "eip161dTransition": "0x0",
    "eip140Transition": "0x0",
    "eip211Transition": "0x0",
    "eip214Transition": "0x0",
    "eip155Transition": "0x0",
    "eip658Transition": "0x0"
  },
  "genesis": {
    "seal": {
      "ethereum": {
        "nonce": "0x0000000000000042",
        "mixHash": "0x0000000000000000000000000000000000000000000000000000000000000000"
      }
    },
    "difficulty": "0x400000000",
    "author": "0x0000000000000000000000000000000000000000",
    "timestamp": "0x00",
    "parentHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
    "extraData": "0x11bbe8db4e347b4e8c937c1c8370e4b5ed33adb3db69cbdb7a38e1e50b1b82fa",
    "gasLimit": "0x1388"
  },
  "accounts": {
    "0000000000000000000000000000000000000001": { "balance": "1", "builtin": { "name": "ecrecover", "pricing": { "linear": { "base": 3000, "word": 0 } } } },
    "0000000000000000000000000000000000000002": { "balance": "1", "builtin": { "name": "sha256", "pricing": { "linear": { "base": 60, "word": 12 } } } },
    "0000000000000000000000000000000000000003": { "balance": "1", "builtin": { "name": "ripemd160", "pricing": { "linear": { "base": 600, "word": 120 } } } },
    "0000000000000000000000000000000000000004": { "balance": "1", "builtin": { "name": "identity", "pricing": { "linear": { "base": 15, "word": 3 } } } },
    "0000000000000000000000000000000000000005": { "builtin": { "name": "modexp", "activate_at": "0x00", "pricing": { "modexp": { "divisor": 20 } } } },
    "0000000000000000000000000000000000000006": {
      "builtin": {
        "name": "alt_bn128_add",
        "pricing": {
          "0": {
            "price": { "alt_bn128_const_operations": { "price": 500 }}
          },
          "0x7fffffffffffff": {
            "info": "EIP 1108 transition",
            "price": { "alt_bn128_const_operations": { "price": 150 }}
          }
        }
      }
    },
    "0000000000000000000000000000000000000007": {
      "builtin": {
        "name": "alt_bn128_mul",
        "pricing": {
          "42": {
            "price": { "alt_bn128_const_operations": { "price": 50000 }}
          },
          "42": {
            "price": { "alt_bn128_const_operations": { "price": 40000 }}
          },
          "0x2a": {
            "price": { "alt_bn128_const_operations": { "price": 30000 }}
          },
          "0x02a": {
            "price": { "alt_bn128_const_operations": { "price": 10000 }}
          }
        }
      }
    },
    "0000000000000000000000000000000000000008": {
      "builtin": {
        "name": "alt_bn128_pairing",
        "pricing": {
          "0": {
            "price": { "alt_bn128_pairing": { "base": 100000, "pair": 80000 }}
          },
          "0x7fffffffffffff": {
            "info": "EIP 1108 transition",
            "price": { "alt_bn128_pairing": { "base": 45000, "pair": 34000 }}
          }
        }
      }
    }
  }
}
`), pspec)
	if err != nil {
		t.Fatal(err)
	}

	// Use 'alt_bn128_mul' as our test case.
	got := pspec.GetPrecompile(common.BytesToAddress([]byte{7}),
		ParityChainSpecPricing{
			AltBnConstOperation: &ParityChainSpecAltBnConstOperationPricing{
				// Want 40000 because "42" lexically comes after "0x2a" and "0x02a".
				// And then we cross our fingers and hope that literal order is preserved
				// for the map during JSON unmarshaling with the keys as strings, since "42" == "42".
				// We want bottom-wins parsing.
				Price: 40000,
			},
		}).Uint64P()

	want := uint64(42)

	if got == nil || *got != want {
		t.Errorf("got: %v, want: %v", got, want)
	}
}
