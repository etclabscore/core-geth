package generic

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/params/types/coregeth"
	"github.com/ethereum/go-ethereum/params/types/goethereum"
)

func TestUnmarshalChainConfigurator(t *testing.T) {
	cases := []struct {
		file  string
		wantT interface{}
	}{
		{
			filepath.Join("..", "testdata", "geth_foundation.json"),
			&goethereum.ChainConfig{},
		},
		{
			filepath.Join("..", "testdata", "coregeth_foundation.json"),
			&coregeth.CoreGethChainConfig{},
		},
	}

	for i, c := range cases {
		b, err := os.ReadFile(c.file)
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
