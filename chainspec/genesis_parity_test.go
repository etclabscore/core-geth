package chainspec

import (
	"encoding/json"
	"io/ioutil"
	"math/big"
	"path/filepath"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/common"
	math2 "github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core"
)

func asSpecFilePath(name string) string {
	return filepath.Join("..", "chainspecs", "parity", name)
}

var chainSpecEquivs = map[string]*core.Genesis{
	"classic.json":    core.DefaultClassicGenesisBlock(),
	"foundation.json": core.DefaultGenesisBlock(),
}

func TestParityConfigToMultiGethGenesis(t *testing.T) {
	var gen1, gen2 *core.Genesis

	for p, gen := range chainSpecEquivs {
		gen1 = gen
		gen1.Config.EIP150Hash = common.Hash{}
		b, err := ioutil.ReadFile(asSpecFilePath(p))
		if err != nil {
			t.Fatal(err)
		}
		paritySpec := ParityChainSpec{}
		err = json.Unmarshal(b, &paritySpec)
		if err != nil {
			t.Fatal(err)
		}

		gen2, err = ParityConfigToMultiGethGenesis(&paritySpec)
		if err != nil {
			t.Fatal(err)
		}

		for _, bl := range []int64{
			0, 10000, 40000, 50000, 100000,
			2000000, 4000000, 6000000, 8000000, 10000000,
		} {
			b := big.NewInt(bl)
			fns := []func(b *big.Int) bool{
				gen1.Config.IsEIP2F, gen2.Config.IsEIP2F,
				gen1.Config.IsEIP100F, gen2.Config.IsEIP100F,
				gen1.Config.IsEIP213F, gen2.Config.IsEIP213F,
				gen1.Config.IsEIP1052F, gen2.Config.IsEIP1052F,
				gen1.Config.IsEIP140F, gen2.Config.IsEIP140F,
				gen1.Config.IsEIP161F, gen2.Config.IsEIP161F,
			}
			for i, f := range fns {
				if i == 0 || i%2 == 0 {
					continue
				}
				if (f(b) && !fns[i-1](b)) || (!f(b) && fns[i-1](b)) {
					t.Errorf("%d mismatch", i)
				}
			}
		}
	}
}

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
	b := parityChainSpecBuiltin{}
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