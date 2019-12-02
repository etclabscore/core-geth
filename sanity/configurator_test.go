package sanity

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/ethereum/go-ethereum/params/types/common"
	"github.com/ethereum/go-ethereum/params/types/parity"
	"github.com/ethereum/go-ethereum/tests"
)

func TestEquivalent(t *testing.T) {
	t.Skip("TODO(meows)")
	parityP := filepath.Join("..", "params", "parity.json.d")
	for k, v := range tests.MapForkNameChainspecFileState {
		a := tests.Forks[k]

		b := &parity.ParityChainSpec{}
		bs, err := ioutil.ReadFile(filepath.Join(parityP, v))
		if err != nil {
			t.Fatal(err)
		}
		err = json.Unmarshal(bs, b)
		if err != nil {
			t.Fatal(err)
		}
		err = common.Equivalent(a, b)
		if err != nil {
			t.Errorf("%s:%s err: %v", k, v, err)
		}
	}
}