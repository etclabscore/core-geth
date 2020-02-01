package generic

import (
	"encoding/json"
	"io/ioutil"
	"math/big"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params/confp"
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

// An example of v1.9.6 multigeth config marshaled to JSON.
// Note the fields EIP1108FBlock; these were included accidentally because
// of a typo in the struct field json tags, and because of that, will
// not be omitted when empty, nor "properly" (lowercase) named.
//
// This should be treated as an 'oldmultigeth' data type, since it has values which are
// not present in the contemporary multigeth data type.
var cc_v196_a = `{
  "chainId":61,
  "homesteadBlock":1150000,
  "daoForkBlock":1920000,
  "eip150Block":2500000,
  "eip150Hash":"0xca12c63534f565899681965528d536c52cb05b7c48e269c2a6cb77ad864d878a",
  "eip155Block":3000000,
  "eip158Block":8772000,
  "byzantiumBlock":8772000,
  "constantinopleBlock":9573000,
  "petersburgBlock":9573000,
  "ethash":{

  
},
  "trustedCheckpoint":null,
  "trustedCheckpointOracle":null,
  "networkId":1,
  "eip7FBlock":null,
  "eip160Block":3000000,
  "EIP1108FBlock":null,
  "EIP1344FBlock":null,
  "EIP1884FBlock":null,
  "EIP2028FBlock":null,
  "EIP2200FBlock":null,
  "ecip1010PauseBlock":3000000,
  "ecip1010Length":2000000,
  "ecip1017FBlock":5000000,
  "ecip1017EraRounds":5000000,
  "disposalBlock":5900000
}`

// An example of a v1.9.6 config which has been marshaled as "new" (v1.9.7+) by a v1.9.7+ client.
// Note that the JSON does not contain eip158 or constantinople upgrade;
// these values are not present in the multigeth datatype, and when
// marshaling their values are simply dropped.
// This is a BAD (mangled) config.
// Code commented now, it's unused. But I left it here because it may be of use down the line,
// so serves at least a documentation purpose. Feel free to remove it if you hate it.
//var cc_v196_b = `{
//    "networkId": 1,
//    "chainId": 61,
//    "daoForkBlock": 1920000,
//    "eip150Block": 2500000,
//    "eip155Block": 3000000,
//    "eip160Block": 3000000,
//    "petersburgBlock": 9573000,
//    "ecip1010PauseBlock": 3000000,
//    "ecip1010Length": 2000000,
//    "ecip1017FBlock": 5000000,
//    "ecip1017EraRounds": 5000000,
//    "disposalBlock": 5900000,
//    "ethash": {},
//    "trustedCheckpoint": null,
//    "trustedCheckpointOracle": null,
//    "requireBlockHashes": null
//}`

// An example of a "healthy" multigeth configuration marshaled to JSON.
var cc_v197_a = `{
    "networkId": 1,
    "chainId": 61,
    "eip2FBlock": 1150000,
    "eip7FBlock": 1150000,
    "eip150Block": 2500000,
    "eip155Block": 3000000,
    "eip160Block": 3000000,
    "eip161FBlock": 8772000,
    "eip170FBlock": 8772000,
    "eip100FBlock": 8772000,
    "eip140FBlock": 8772000,
    "eip198FBlock": 8772000,
    "eip211FBlock": 8772000,
    "eip212FBlock": 8772000,
    "eip213FBlock": 8772000,
    "eip214FBlock": 8772000,
    "eip658FBlock": 8772000,
    "eip145FBlock": 9573000,
    "eip1014FBlock": 9573000,
    "eip1052FBlock": 9573000,
    "eip152FBlock": 10500839,
    "eip1108FBlock": 10500839,
    "eip1344FBlock": 10500839,
    "eip2028FBlock": 10500839,
    "eip2200FBlock": 10500839,
    "ecip1010PauseBlock": 3000000,
    "ecip1010Length": 2000000,
    "ecip1017FBlock": 5000000,
    "ecip1017EraRounds": 5000000,
    "disposalBlock": 5900000,
    "ethash": {},
    "trustedCheckpoint": null,
    "trustedCheckpointOracle": null,
    "requireBlockHashes": {
        "1920000": "0x94365e3a8c0b35089c1d1195081fe7489b528a84b22199c916180db8b28ade7f",
        "2500000": "0xca12c63534f565899681965528d536c52cb05b7c48e269c2a6cb77ad864d878a"
    }
}`

func TestUnmarshalChainConfigurator2(t *testing.T) {
	conf, err := UnmarshalChainConfigurator([]byte(cc_v196_a))
	if err != nil {
		t.Fatal(err)
	}
	wantType := reflect.TypeOf(&oldmultigeth.ChainConfig{})
	if reflect.TypeOf(conf) != wantType {
		t.Fatalf("mismatch, want: %v, got: %v", wantType, reflect.TypeOf(conf))
	}

	conf2, err := UnmarshalChainConfigurator([]byte(cc_v197_a))
	if err != nil {
		t.Fatal(err)
	}
	head := uint64(10_000_000)
	compatErr := confp.Compatible(&head, conf, conf2)
	if compatErr != nil {
		t.Error(compatErr)
	}
}
