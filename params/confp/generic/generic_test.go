package generic

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/params/types/goethereum"
	"github.com/ethereum/go-ethereum/params/types/multigeth"
	"github.com/ethereum/go-ethereum/params/types/parity"
)

func TestUnmarshalChainConfigurator(t *testing.T) {
	cases := []struct {
		file  string
		wantT interface{}
	}{
		{
			filepath.Join("..", "testdata", "stureby_parity.json"),
			&parity.ParityChainSpec{},
		},
		{
			filepath.Join("..", "testdata", "stureby_geth.json"),
			&goethereum.ChainConfig{},
		},
		{
			filepath.Join("..", "testdata", "stureby_multigeth.json"),
			&multigeth.MultiGethChainConfig{},
		},
	}

	for i, c := range cases {
		b, err := ioutil.ReadFile(c.file)
		if err != nil {
			t.Fatal(err)
		}
		got, err := UnmarshalChainConfigurator(b)
		if err != nil {
			t.Fatal(err)
		}
		if reflect.TypeOf(got) != reflect.TypeOf(c.wantT) {
			gotb, _ := json.MarshalIndent(got, "", "    ")
			t.Errorf(`%d / wrong type
want: (%s)
got: (%s)
---
file:
%s
---
result:
%s
`,
				i,
				reflect.TypeOf(c.wantT).String(),
				reflect.TypeOf(got).String(),
				string(b),
				string(gotb),
			)
		}
	}
}
