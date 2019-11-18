package chainspec

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/go-test/deep"
)

func asSpecFilePath(name string) string {
	return filepath.Join("..", "chainspecs", name)
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

		if diffs := deep.Equal(gen1, gen2); len(diffs) != 0 {
			//b, _ := json.MarshalIndent(gen1, "", "    ")
			//t.Log(string(b)[:10000])
			for _, d := range diffs {
				t.Errorf("%s | diff: %s", p, d)
			}
		}
	}
}

var exampleAccountWithBuiltinA = []byte(`{
			"builtin": {
				"name": "modexp",
				"activate_at": "0x85d9a0",
				"pricing": {
					"modexp": {
						"divisor": 20
					}
				}
			}
		}`)
var exampleAccountWithBuiltinB = []byte(`{
			"builtin": {
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
		}`)

func TestParityBuiltinType(t *testing.T) {

}
