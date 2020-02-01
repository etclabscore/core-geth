package generic

import (
	"encoding/json"
	"io/ioutil"
	"math/big"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params/types/ctypes"
	"github.com/ethereum/go-ethereum/params/types/goethereum"
	"github.com/ethereum/go-ethereum/params/types/multigeth"
	"github.com/ethereum/go-ethereum/params/types/oldmultigeth"
	"github.com/ethereum/go-ethereum/params/types/parity"
)

// TestUnmarshalChainConfigurator is a non deterministic test. WTF.
// go test ./params/... -run TestUnmarshalChainConfigurator -count [1|100]
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

	om := &oldmultigeth.ChainConfig{
		ChainID:              big.NewInt(61),
		HomesteadBlock:       big.NewInt(1150000),
		DAOForkBlock:         big.NewInt(1920000),
		DAOForkSupport:       false,
		EIP150Block:          big.NewInt(2500000),
		EIP150Hash:           common.HexToHash("0xca12c63534f565899681965528d536c52cb05b7c48e269c2a6cb77ad864d878a"),
		EIP155Block:          big.NewInt(3000000),
		EIP158Block:          big.NewInt(8772000),
		ByzantiumBlock:       big.NewInt(8772000),
		DisposalBlock:        big.NewInt(5900000),
		SocialBlock:          nil,
		EthersocialBlock:     nil,
		ConstantinopleBlock:  big.NewInt(9573000),
		PetersburgBlock:      big.NewInt(9573000),
		IstanbulBlock:        big.NewInt(10500839),
		EIP1884DisableFBlock: big.NewInt(10500839),
		ECIP1017EraRounds:    big.NewInt(5000000),
		EIP160FBlock:         big.NewInt(3000000),
		ECIP1010PauseBlock:   big.NewInt(3000000),
		ECIP1010Length:       big.NewInt(2000000),
		Ethash:               new(ctypes.EthashConfig),
	}

	b, err := json.MarshalIndent(om, "", "    ")
	if err != nil {
		t.Fatal(err)
	}

	got, err := UnmarshalChainConfigurator(b)
	if err != nil {
		t.Fatal(err)
	}
	if reflect.TypeOf(got) != reflect.TypeOf(&oldmultigeth.ChainConfig{}) {
		t.Fatalf("mismatch, want: %v, got: %v", reflect.TypeOf(&oldmultigeth.ChainConfig{}), reflect.TypeOf(got))
	}

	if tr := got.GetEIP7Transition(); tr == nil || *tr != 1150000 {
		t.Fatal("bad")
	}
	if tr := got.GetEIP7Transition(); tr == nil || *tr != 1150000 {
		t.Fatal("bad", *tr)
	}
	if tr := got.GetEIP150Transition(); tr == nil || *tr != 2500000 {
		t.Fatal("bad", *tr)
	}
	if tr := got.GetEIP160Transition(); tr == nil || *tr != 3000000 {
		t.Fatal("bad", *tr)
	}
	if tr := got.GetEIP155Transition(); tr == nil || *tr != 3000000 {
		t.Fatal("bad", *tr)
	}
	if tr := got.GetEIP170Transition(); tr == nil || *tr != 8772000 {
		t.Fatal("bad", *tr)
	}
}
