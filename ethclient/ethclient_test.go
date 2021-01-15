// Copyright 2016 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package ethclient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"reflect"
	"strings"
	"testing"
	"text/template"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/consensus/ethash"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/eth"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/params/types/genesisT"
	meta_schema "github.com/open-rpc/meta-schema"
)

// Verify that Client implements the ethereum interfaces.
var (
	_ = ethereum.ChainReader(&Client{})
	_ = ethereum.TransactionReader(&Client{})
	_ = ethereum.ChainStateReader(&Client{})
	_ = ethereum.ChainSyncReader(&Client{})
	_ = ethereum.ContractCaller(&Client{})
	_ = ethereum.GasEstimator(&Client{})
	_ = ethereum.GasPricer(&Client{})
	_ = ethereum.LogFilterer(&Client{})
	_ = ethereum.PendingStateReader(&Client{})
	// _ = ethereum.PendingStateEventer(&Client{})
	_ = ethereum.PendingContractCaller(&Client{})
)

func TestToFilterArg(t *testing.T) {
	blockHashErr := fmt.Errorf("cannot specify both BlockHash and FromBlock/ToBlock")
	addresses := []common.Address{
		common.HexToAddress("0xD36722ADeC3EdCB29c8e7b5a47f352D701393462"),
	}
	blockHash := common.HexToHash(
		"0xeb94bb7d78b73657a9d7a99792413f50c0a45c51fc62bdcb08a53f18e9a2b4eb",
	)

	for _, testCase := range []struct {
		name   string
		input  ethereum.FilterQuery
		output interface{}
		err    error
	}{
		{
			"without BlockHash",
			ethereum.FilterQuery{
				Addresses: addresses,
				FromBlock: big.NewInt(1),
				ToBlock:   big.NewInt(2),
				Topics:    [][]common.Hash{},
			},
			map[string]interface{}{
				"address":   addresses,
				"fromBlock": "0x1",
				"toBlock":   "0x2",
				"topics":    [][]common.Hash{},
			},
			nil,
		},
		{
			"with nil fromBlock and nil toBlock",
			ethereum.FilterQuery{
				Addresses: addresses,
				Topics:    [][]common.Hash{},
			},
			map[string]interface{}{
				"address":   addresses,
				"fromBlock": "0x0",
				"toBlock":   "latest",
				"topics":    [][]common.Hash{},
			},
			nil,
		},
		{
			"with negative fromBlock and negative toBlock",
			ethereum.FilterQuery{
				Addresses: addresses,
				FromBlock: big.NewInt(-1),
				ToBlock:   big.NewInt(-1),
				Topics:    [][]common.Hash{},
			},
			map[string]interface{}{
				"address":   addresses,
				"fromBlock": "pending",
				"toBlock":   "pending",
				"topics":    [][]common.Hash{},
			},
			nil,
		},
		{
			"with blockhash",
			ethereum.FilterQuery{
				Addresses: addresses,
				BlockHash: &blockHash,
				Topics:    [][]common.Hash{},
			},
			map[string]interface{}{
				"address":   addresses,
				"blockHash": blockHash,
				"topics":    [][]common.Hash{},
			},
			nil,
		},
		{
			"with blockhash and from block",
			ethereum.FilterQuery{
				Addresses: addresses,
				BlockHash: &blockHash,
				FromBlock: big.NewInt(1),
				Topics:    [][]common.Hash{},
			},
			nil,
			blockHashErr,
		},
		{
			"with blockhash and to block",
			ethereum.FilterQuery{
				Addresses: addresses,
				BlockHash: &blockHash,
				ToBlock:   big.NewInt(1),
				Topics:    [][]common.Hash{},
			},
			nil,
			blockHashErr,
		},
		{
			"with blockhash and both from / to block",
			ethereum.FilterQuery{
				Addresses: addresses,
				BlockHash: &blockHash,
				FromBlock: big.NewInt(1),
				ToBlock:   big.NewInt(2),
				Topics:    [][]common.Hash{},
			},
			nil,
			blockHashErr,
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			output, err := toFilterArg(testCase.input)
			if (testCase.err == nil) != (err == nil) {
				t.Fatalf("expected error %v but got %v", testCase.err, err)
			}
			if testCase.err != nil {
				if testCase.err.Error() != err.Error() {
					t.Fatalf("expected error %v but got %v", testCase.err, err)
				}
			} else if !reflect.DeepEqual(testCase.output, output) {
				t.Fatalf("expected filter arg %v but got %v", testCase.output, output)
			}
		})
	}
}

var (
	testKey, _  = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
	testAddr    = crypto.PubkeyToAddress(testKey.PublicKey)
	testBalance = big.NewInt(2e10)
)

func newTestBackend(t *testing.T) (*node.Node, []*types.Block) {
	// Generate test chain.
	genesis, blocks := generateTestChain()
	// Create node
	n, err := node.New(&node.Config{})
	if err != nil {
		t.Fatalf("can't create new node: %v", err)
	}
	// Create Ethereum Service
	config := &eth.Config{Genesis: genesis}
	config.Ethash.PowMode = ethash.ModeFake
	ethservice, err := eth.New(n, config)
	if err != nil {
		t.Fatalf("can't create new ethereum service: %v", err)
	}
	// Import the test chain.
	if err := n.Start(); err != nil {
		t.Fatalf("can't start test node: %v", err)
	}
	if _, err := ethservice.BlockChain().InsertChain(blocks[1:]); err != nil {
		t.Fatalf("can't import test blocks: %v", err)
	}
	return n, blocks
}

func generateTestChain() (*genesisT.Genesis, []*types.Block) {
	db := rawdb.NewMemoryDatabase()
	config := params.AllEthashProtocolChanges
	genesis := &genesisT.Genesis{
		Config:    config,
		Alloc:     genesisT.GenesisAlloc{testAddr: {Balance: testBalance}},
		ExtraData: []byte("test genesis"),
		Timestamp: 9000,
	}
	generate := func(i int, g *core.BlockGen) {
		g.OffsetTime(5)
		g.SetExtra([]byte("test"))
	}
	gblock := core.GenesisToBlock(genesis, db)
	engine := ethash.NewFaker()
	blocks, _ := core.GenerateChain(config, gblock, engine, db, 1, generate)
	blocks = append([]*types.Block{gblock}, blocks...)
	return genesis, blocks
}

func TestHeader(t *testing.T) {
	backend, chain := newTestBackend(t)
	client, _ := backend.Attach()
	defer backend.Close()
	defer client.Close()

	tests := map[string]struct {
		block   *big.Int
		want    *types.Header
		wantErr error
	}{
		"genesis": {
			block: big.NewInt(0),
			want:  chain[0].Header(),
		},
		"first_block": {
			block: big.NewInt(1),
			want:  chain[1].Header(),
		},
		"future_block": {
			block: big.NewInt(1000000000),
			want:  nil,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ec := NewClient(client)
			ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
			defer cancel()

			got, err := ec.HeaderByNumber(ctx, tt.block)
			if tt.wantErr != nil && (err == nil || err.Error() != tt.wantErr.Error()) {
				t.Fatalf("HeaderByNumber(%v) error = %q, want %q", tt.block, err, tt.wantErr)
			}
			if got != nil && got.Number.Sign() == 0 {
				got.Number = big.NewInt(0) // hack to make DeepEqual work
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("HeaderByNumber(%v)\n   = %v\nwant %v", tt.block, got, tt.want)
			}
		})
	}
}

func TestBalanceAt(t *testing.T) {
	backend, _ := newTestBackend(t)
	client, _ := backend.Attach()
	defer backend.Close()
	defer client.Close()

	tests := map[string]struct {
		account common.Address
		block   *big.Int
		want    *big.Int
		wantErr error
	}{
		"valid_account": {
			account: testAddr,
			block:   big.NewInt(1),
			want:    testBalance,
		},
		"non_existent_account": {
			account: common.Address{1},
			block:   big.NewInt(1),
			want:    big.NewInt(0),
		},
		"future_block": {
			account: testAddr,
			block:   big.NewInt(1000000000),
			want:    big.NewInt(0),
			wantErr: errors.New("header not found"),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ec := NewClient(client)
			ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
			defer cancel()

			got, err := ec.BalanceAt(ctx, tt.account, tt.block)
			if tt.wantErr != nil && (err == nil || err.Error() != tt.wantErr.Error()) {
				t.Fatalf("BalanceAt(%x, %v) error = %q, want %q", tt.account, tt.block, err, tt.wantErr)
			}
			if got.Cmp(tt.want) != 0 {
				t.Fatalf("BalanceAt(%x, %v) = %v, want %v", tt.account, tt.block, got, tt.want)
			}
		})
	}
}

func TestTransactionInBlockInterrupted(t *testing.T) {
	backend, _ := newTestBackend(t)
	client, _ := backend.Attach()
	defer backend.Close()
	defer client.Close()

	ec := NewClient(client)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	tx, err := ec.TransactionInBlock(ctx, common.Hash{1}, 1)
	if tx != nil {
		t.Fatal("transaction should be nil")
	}
	if err == nil {
		t.Fatal("error should not be nil")
	}
}

func TestChainID(t *testing.T) {
	backend, _ := newTestBackend(t)
	client, _ := backend.Attach()
	defer backend.Close()
	defer client.Close()
	ec := NewClient(client)

	id, err := ec.ChainID(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id == nil || id.Cmp(params.AllEthashProtocolChanges.ChainID) != 0 {
		t.Fatalf("ChainID returned wrong number: %+v", id)
	}
}

func TestBlockNumber(t *testing.T) {
	backend, _ := newTestBackend(t)
	client, _ := backend.Attach()
	defer backend.Close()
	defer client.Close()
	ec := NewClient(client)

	blockNumber, err := ec.BlockNumber(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if blockNumber != 1 {
		t.Fatalf("BlockNumber returned wrong number: %d", blockNumber)
	}
}

func TestRPCDiscover(t *testing.T) {
	backend, _ := newTestBackend(t)
	client, _ := backend.Attach()
	defer backend.Close()
	defer client.Close()

	var res meta_schema.OpenrpcDocument
	err := client.Call(&res, "rpc.discover")
	if err != nil {
		t.Fatal(err)
	}

	sliceContains := func(sl []string, str string) bool {
		for _, s := range sl {
			if str == s {
				return true
			}
		}
		return false
	}

	methodNamesSlice := func() (names []string) {
		for _, m := range *res.Methods {
			names = append(names, string(*m.Name))
		}
		return
	}()

	over, under := []string{}, []string{}

	for _, name := range methodNamesSlice {
		if !sliceContains(allRPCMethods, name) {
			under = append(under, name)
		}
	}
	for _, name := range allRPCMethods {
		if !sliceContains(methodNamesSlice, name) {
			over = append(over, name)
		}
	}

	if len(over) > 0 || len(under) > 0 {
		t.Fatalf("over: %v, under: %v", over, under)
	}
}

func printBullet(any interface{}, depth int) (out string) {
	defer func() {
		out += "\n"
	}()
	switch any.(type) {
	case map[string]interface{}:
		out += "\n"
		for k, v := range any.(map[string]interface{}) {
			out += fmt.Sprintf("%s- %s: %s", strings.Repeat("\t", depth), k, printBullet(v, depth+1))
		}
	case []interface{}:

		// Don't know why this isn't working. Doesn't fire.
		// if c, ok := any.([]string); ok {
		// 	return strings.Join(c, ",")
		// }

		stringSet := []string{}
		for _, vv := range any.([]interface{}) {
			if s, ok := vv.(string); ok {
				stringSet = append(stringSet, s)
			}
		}
		if len(stringSet) == len(any.([]interface{})) {
			return strings.Join(stringSet, ", ")
		}

		out += "\n"
		for _, vv := range any.([]interface{}) {
			out += printBullet(vv, depth+1)
		}
	default:
		return fmt.Sprintf("`%v`", any)
	}
	return
}

// TestRPCDiscover_BuildStatic puts the OpenRPC document in build/static/openrpc.json.
// This is intended to be run as a documentation development tool (as opposed to an actual _test_).
// NOTE that Go maps don't guarantee order, so the diff between runs can be noisy.
func TestRPCDiscover_BuildStatic(t *testing.T) {
	if os.Getenv("COREGETH_GEN_OPENRPC_DOCS") == "" {
		return
	}
	err := os.MkdirAll("../build/static", os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

	backend, _ := newTestBackend(t)
	client, _ := backend.Attach()
	defer backend.Close()
	defer client.Close()

	// Current workaround for https://github.com/open-rpc/meta-schema/issues/356.
	res, err := backend.InprocDiscovery().Discover()
	if err != nil {
		t.Fatal(err)
	}

	// Should do it this way.
	// res := &meta_schema.OpenrpcDocument{}
	// err = client.Call(res, "rpc.discover")
	// if err != nil {
	// 	t.Fatal(err)
	// }

	data, err := json.MarshalIndent(res, "", "    ")
	if err != nil {
		t.Fatal(err)
	}
	err = ioutil.WriteFile("../build/static/openrpc.json", data, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

	if res.Methods == nil {
		return
	}

	tpl := template.New("openrpc_doc")

	tpl.Funcs(template.FuncMap{
		// tomap gives a plain map from a JSON representation of a given value.
		// This is useful because the meta_schema data types, being generated, and conforming to a pretty
		// complex data type in the first place, are not super fun to interact with directly.
		"tomap": func(any interface{}) map[string]interface{} {
			out := make(map[string]interface{})
			data, _ := json.Marshal(any)
			json.Unmarshal(data, &out)
			return out
		},
		// asjson returns indented JSON.
		"asjson": func(any interface{}, prefix, indent int) string {
			by, _ := json.MarshalIndent(any, strings.Repeat("    ", prefix), strings.Repeat("    ", indent))
			return string(by)
		},
		// bulletJSON handles transforming a JSON JSON schema into bullet points, which I think are more legible.
		"bulletJSON": printBullet,
		"sum":        func(a, b int) int { return a + b },
		// trimNameSpecialChars removes characters that the app-specific content descriptor naming
		// method will also remove, eg '*hexutil.Uint64' -> 'hexutilUint64'.
		// These "special" characters were removed because of concerns about by-name arguments
		// and the use of titles for keys.
		"trimNameSpecialChars": func(s string) string {
			remove := []string{".", "*"}
			for _, r := range remove {
				s = strings.ReplaceAll(s, r, "")
			}
			return s
		},
		// methodFormatJSConsole is a pretty-printer that returns the JS console use example for a method.
		"methodFormatJSConsole": func(m *meta_schema.MethodObject) string {
			name := string(*m.Name)
			formattedName := strings.Replace(name, "_", ".", 1)
			getParamName := func(cd *meta_schema.ContentDescriptorObject) string {
				if cd.Name != nil {
					return string(*cd.Name)
				}
				return string(*cd.Description)
			}
			paramNames := func() (paramNames []string) {
				if m.Params == nil {
					return nil
				}
				for _, n := range *m.Params {
					if n.ContentDescriptorObject == nil {
						continue // Should never happen in our implementation; never uses refs.
					}
					paramNames = append(paramNames, getParamName(n.ContentDescriptorObject))
				}
				return
			}()
			return fmt.Sprintf("%s(%s);", formattedName, strings.Join(paramNames, ","))
		},
	})

	tpl, err = tpl.Parse(`
{{ define "schemaTpl" }}
	` + "```" + `
{{ bulletJSON . 1 }}
	` + "```" + `

	<details class="cite"><summary>View Raw</summary>
	` + "```" + `
    {{ asjson . 1 1 }}
	` + "```" + `
	</details>
{{ end }}

{{ define "contentDescTpl" -}}
{{ $nameyDescription := trimNameSpecialChars .description }}
{{ if eq .name $nameyDescription }}
<code>{{ .description }}</code> {{ if .summary }}_{{ .summary }}_{{- end }}
{{- else -}}
{{ .name }} <code>{{ .description }}</code> {{ if .summary }}_{{ .summary }}_{{- end }}
{{- end }}

  + Required: {{ if .required }}✓ Yes{{ else }}No{{- end}}
{{ if .deprecated }}  + Deprecated: :warning: Yes{{- end}}
{{ if (or (gt (len .schema) 1) .schema.properties) }} 
=== "Schema"

	` + "```" + ` Schema
	{{ bulletJSON .schema 1 }}
	` + "```" + `

=== "Raw"

	` + "```" + ` Raw
	{{ asjson .schema 1 1 }}
	` + "```" + `
{{ end }}
{{ end }}

{{ define "methodTpl" }}
{{ $methodmap := tomap . }}
### {{ .Name }}

{{ .Summary }}

__Params ({{ .Params | len }})__
{{ if gt (.Params | len) 0 }}
{{ if eq $methodmap.paramStructure "by-position" }}Parameters must be given _by position_.{{ else if eq $methodmap.paramStructure "by-name" }}Parameters must be given _by name_.{{ end }}  
{{ range $index, $param := .Params }}
{{ $parammap := . | tomap }}
__{{ sum $index 1 }}:__ {{ template "contentDescTpl" $parammap }}
{{ end }}
{{ else }}
_None_
{{- end}}

__Result__

{{ if .Result -}}
{{ $result := .Result | tomap }}
{{- if ne $result.name "Null" }}
{{ template "contentDescTpl" $result }}
{{- else -}}
_None_
{{- end }}
{{- end }}

__Client Method Invocation Examples__

=== "Shell"

	` + "```" + ` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "{{ .Name }}", "params": []}'
	` + "```" + `

=== "Javascript Console"

	` + "```" + ` js
	{{ methodFormatJSConsole . }}
	` + "```" + `

{{ $docs := .ExternalDocs | tomap }}
<details><summary>Source code</summary>
<p>
{{ .Description }}
<a href="{{ $docs.url }}" target="_">View on GitHub →</a>
</p>
</details>

---
{{- end }}

| Entity | Version |
| --- | --- |
| Source | <code>{{ .Info.Version }}</code> |
| OpenRPC | <code>{{ .Openrpc }}</code> |

---

{{ range .Methods }}
{{ template "methodTpl" . }}
{{ end }}
`)
	if err != nil {
		t.Fatal(err)
	}

	moduleMethods := func() (grouped map[string][]meta_schema.MethodObject) {
		if res.Methods == nil {
			return
		}
		grouped = make(map[string][]meta_schema.MethodObject)
		for _, m := range *res.Methods {
			moduleName := strings.Split(string(*m.Name), "_")[0]
			group, ok := grouped[moduleName]
			if !ok {
				group = []meta_schema.MethodObject{}
			}
			group = append(group, m)
			grouped[moduleName] = group
		}
		return
	}()

	_ = os.MkdirAll("../docs/JSON-RPC-API/modules", os.ModePerm)

	for module, group := range moduleMethods {
		fname := fmt.Sprintf("../docs/JSON-RPC-API/modules/%s.md", module)
		fi, err := os.OpenFile(fname, os.O_CREATE|os.O_RDWR, os.ModePerm)
		if err != nil {
			t.Fatal(err)
		}
		fi.Truncate(0)

		nDoc := &meta_schema.OpenrpcDocument{}
		*nDoc = *res
		nDoc.Methods = (*meta_schema.Methods)(&group)
		err = tpl.Execute(fi, nDoc)
		if err != nil {
			t.Fatal(err)
		}
		fi.Sync()
		fi.Close()
	}
}

func subscriptionTestSetup(t *testing.T) (genesisBlock *genesisT.Genesis, backend *node.Node) {
	// Generate test chain.
	// Code largely taken from generateTestChain()
	chainConfig := params.TestChainConfig
	genesis := &genesisT.Genesis{
		Config:    chainConfig,
		Alloc:     genesisT.GenesisAlloc{testAddr: {Balance: testBalance}},
		ExtraData: []byte("test genesis"),
		Timestamp: 9000,
	}

	// Create node
	// Code largely taken from newTestBackend(t)
	backend, err := node.New(&node.Config{})
	if err != nil {
		t.Fatalf("can't create new node: %v", err)
	}

	// Create Ethereum Service
	config := &eth.Config{Genesis: genesis}
	config.Ethash.PowMode = ethash.ModeFake

	return genesis, backend
}

func TestEthSubscribeNewSideHeads(t *testing.T) {

	genesis, backend := subscriptionTestSetup(t)

	db := rawdb.NewMemoryDatabase()
	chainConfig := genesis.Config

	gblock := core.GenesisToBlock(genesis, db)
	engine := ethash.NewFaker()
	originalBlocks, _ := core.GenerateChain(chainConfig, gblock, engine, db, 10, func(i int, gen *core.BlockGen) {
		gen.OffsetTime(5)
		gen.SetExtra([]byte("test"))
	})
	originalBlocks = append([]*types.Block{gblock}, originalBlocks...)

	// Create Ethereum Service
	config := &eth.Config{Genesis: genesis}
	config.Ethash.PowMode = ethash.ModeFake
	ethservice, err := eth.New(backend, config)
	if err != nil {
		t.Fatalf("can't create new ethereum service: %v", err)
	}

	// Import the test chain.
	if err := backend.Start(); err != nil {
		t.Fatalf("can't start test node: %v", err)
	}
	if _, err := ethservice.BlockChain().InsertChain(originalBlocks[1:]); err != nil {
		t.Fatalf("can't import test blocks: %v", err)
	}

	// Create the client and newSideHeads subscription.
	client, err := backend.Attach()
	defer backend.Close()
	defer client.Close()
	if err != nil {
		t.Fatal(err)
	}
	ec := NewClient(client)
	defer ec.Close()

	sideHeadCh := make(chan *types.Header)
	sub, err := ec.SubscribeNewSideHead(context.Background(), sideHeadCh)
	if err != nil {
		t.Fatal(err)
	}
	defer sub.Unsubscribe()

	// Create and import the second-seen chain.
	replacementBlocks, _ := core.GenerateChain(chainConfig, originalBlocks[len(originalBlocks)-5], ethservice.Engine(), db, 5, func(i int, gen *core.BlockGen) {
		gen.OffsetTime(-9) // difficulty++
	})
	if _, err := ethservice.BlockChain().InsertChain(replacementBlocks); err != nil {
		t.Fatalf("can't import test blocks: %v", err)
	}

	headersOf := func(bs []*types.Block) (headers []*types.Header) {
		for _, b := range bs {
			headers = append(headers, b.Header())
		}
		return
	}

	expectations := []*types.Header{}

	// Why do we expect the replacement (second-seen) blocks reported as side events?
	// Because they'll be inserted in ascending order, and until their segment exceeds the total difficulty
	// of the incumbent chain, they won't achieve canonical status, despite having greater difficulty per block
	// (see the time offset in the block generator function above).
	expectations = append(expectations, headersOf(replacementBlocks[:3])...)

	// Once the replacement blocks exceed the total difficulty of the original chain, the
	// blocks they replace will be reported as side chain events.
	expectations = append(expectations, headersOf(originalBlocks[7:])...)

	// This is illustrated in the logs called below.
	for i, b := range originalBlocks {
		t.Log("incumbent", i, b.NumberU64(), b.Hash().Hex()[:8])
	}
	for i, b := range replacementBlocks {
		t.Log("replacement", i, b.NumberU64(), b.Hash().Hex()[:8])
	}

	const timeoutDura = 5 * time.Second
	timeout := time.NewTimer(timeoutDura)

	got := []*types.Header{}
waiting:
	for {
		select {
		case head := <-sideHeadCh:
			t.Log("<-newSideHeads", head.Number.Uint64(), head.Hash().Hex()[:8])
			got = append(got, head)
			if len(got) == len(expectations) {
				timeout.Stop()
				break waiting
			}
			timeout.Reset(timeoutDura)
		case err := <-sub.Err():
			t.Fatal(err)
		case <-timeout.C:
			t.Fatal("timed out")
		}
	}
	for i, b := range expectations {
		if got[i] == nil {
			t.Error("missing expected header (test will improvise a fake value)")
			// Set a nonzero value so I don't have to refactor this...
			got[i] = &types.Header{Number: big.NewInt(math.MaxInt64)}
		}
		if got[i].Number.Uint64() != b.Number.Uint64() {
			t.Errorf("number: want: %d, got: %d", b.Number.Uint64(), got[i].Number.Uint64())
		} else if got[i].Hash() != b.Hash() {
			t.Errorf("hash: want: %s, got: %s", b.Hash().Hex()[:8], got[i].Hash().Hex()[:8])
		}
	}
}

// mustNewTestBackend is the same logic as newTestBackend(t *testing.T) but without the testing.T argument.
// This function is used exclusively for the benchmarking tests, and will panic if it encounters an error.
func mustNewTestBackend() (*node.Node, []*types.Block) {
	// Generate test chain.
	genesis, blocks := generateTestChain()
	// Create node
	n, err := node.New(&node.Config{})
	if err != nil {
		panic(fmt.Sprintf("can't create new node: %v", err))
	}
	// Create Ethereum Service
	config := &eth.Config{Genesis: genesis}
	config.Ethash.PowMode = ethash.ModeFake
	ethservice, err := eth.New(n, config)
	if err != nil {
		panic(fmt.Sprintf("can't create new ethereum service: %v", err))
	}
	// Import the test chain.
	if err := n.Start(); err != nil {
		panic(fmt.Sprintf("can't start test node: %v", err))
	}
	if _, err := ethservice.BlockChain().InsertChain(blocks[1:]); err != nil {
		panic(fmt.Sprintf("can't import test blocks: %v", err))
	}
	return n, blocks
}

// BenchmarkRPC_Discover shows that rpc.discover by reflection is slow.
func BenchmarkRPC_Discover(b *testing.B) {
	backend, _ := mustNewTestBackend()
	client, _ := backend.Attach()
	defer backend.Close()
	defer client.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var res meta_schema.OpenrpcDocument
		err := client.Call(&res, "rpc.discover")
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkRPC_BlockNumber shows that eth_blockNumber is a lot faster than rpc.discover.
func BenchmarkRPC_BlockNumber(b *testing.B) {
	backend, _ := mustNewTestBackend()
	client, _ := backend.Attach()
	defer backend.Close()
	defer client.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var res hexutil.Uint64
		err := client.Call(&res, "eth_blockNumber")
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
	}
}

// allRPCMethods lists all methods exposed over JSONRPC.
var allRPCMethods = []string{
	"admin_addPeer",
	"admin_addTrustedPeer",
	"admin_datadir",
	"admin_ecbp1100",
	"admin_exportChain",
	"admin_importChain",
	"admin_maxPeers",
	"admin_nodeInfo",
	"admin_peers",
	"admin_removePeer",
	"admin_removeTrustedPeer",
	"admin_startRPC",
	"admin_startWS",
	"admin_stopRPC",
	"admin_stopWS",
	"debug_accountRange",
	"debug_backtraceAt",
	"debug_blockProfile",
	"debug_chaindbCompact",
	"debug_chaindbProperty",
	"debug_cpuProfile",
	"debug_dumpBlock",
	"debug_freeOSMemory",
	"debug_gcStats",
	"debug_getBadBlocks",
	"debug_getBlockRlp",
	"debug_getModifiedAccountsByHash",
	"debug_getModifiedAccountsByNumber",
	"debug_goTrace",
	"debug_memStats",
	"debug_mutexProfile",
	"debug_preimage",
	"debug_printBlock",
	"debug_removePendingTransaction",
	"debug_seedHash",
	"debug_setBlockProfileRate",
	"debug_setGCPercent",
	"debug_setHead",
	"debug_setMutexProfileFraction",
	"debug_stacks",
	"debug_standardTraceBadBlockToFile",
	"debug_standardTraceBlockToFile",
	"debug_startCPUProfile",
	"debug_startGoTrace",
	"debug_stopCPUProfile",
	"debug_stopGoTrace",
	"debug_storageRangeAt",
	"debug_testSignCliqueBlock",
	"debug_traceBadBlock",
	"debug_traceBlock",
	"debug_traceBlockByHash",
	"debug_traceBlockByNumber",
	"debug_traceBlockFromFile",
	"debug_traceCall",
	"debug_traceTransaction",
	"debug_verbosity",
	"debug_vmodule",
	"debug_writeBlockProfile",
	"debug_writeMemProfile",
	"debug_writeMutexProfile",
	"eth_accounts",
	"eth_blockNumber",
	"eth_call",
	"eth_chainId",
	"eth_chainId",
	"eth_coinbase",
	"eth_estimateGas",
	"eth_etherbase",
	"eth_fillTransaction",
	"eth_gasPrice",
	"eth_getBalance",
	"eth_getBlockByHash",
	"eth_getBlockByNumber",
	"eth_getBlockTransactionCountByHash",
	"eth_getBlockTransactionCountByNumber",
	"eth_getCode",
	"eth_getFilterChanges",
	"eth_getFilterLogs",
	"eth_getHashrate",
	"eth_getHeaderByHash",
	"eth_getHeaderByNumber",
	"eth_getLogs",
	"eth_getProof",
	"eth_getRawTransactionByBlockHashAndIndex",
	"eth_getRawTransactionByBlockNumberAndIndex",
	"eth_getRawTransactionByHash",
	"eth_getStorageAt",
	"eth_getTransactionByBlockHashAndIndex",
	"eth_getTransactionByBlockNumberAndIndex",
	"eth_getTransactionByHash",
	"eth_getTransactionCount",
	"eth_getTransactionReceipt",
	"eth_getUncleByBlockHashAndIndex",
	"eth_getUncleByBlockNumberAndIndex",
	"eth_getUncleCountByBlockHash",
	"eth_getUncleCountByBlockNumber",
	"eth_getWork",
	"eth_hashrate",
	"eth_mining",
	"eth_newBlockFilter",
	"eth_newSideBlockFilter",
	"eth_newFilter",
	"eth_newPendingTransactionFilter",
	"eth_pendingTransactions",
	"eth_protocolVersion",
	"eth_resend",
	"eth_sendRawTransaction",
	"eth_sendTransaction",
	"eth_sign",
	"eth_signTransaction",
	"eth_submitHashRate",
	"eth_submitWork",
	"eth_subscribe",
	"eth_syncing",
	"eth_uninstallFilter",
	"eth_unsubscribe",
	"ethash_getHashrate",
	"ethash_getWork",
	"ethash_submitHashRate",
	"ethash_submitWork",
	"miner_getHashrate",
	"miner_setEtherbase",
	"miner_setExtra",
	"miner_setGasPrice",
	"miner_setRecommitInterval",
	"miner_start",
	"miner_stop",
	"net_listening",
	"net_peerCount",
	"net_version",
	"personal_deriveAccount",
	"personal_ecRecover",
	"personal_importRawKey",
	"personal_initializeWallet",
	"personal_listAccounts",
	"personal_listWallets",
	"personal_lockAccount",
	"personal_newAccount",
	"personal_openWallet",
	"personal_sendTransaction",
	"personal_sign",
	"personal_signAndSendTransaction",
	"personal_signTransaction",
	"personal_unlockAccount",
	"personal_unpair",
	"trace_block",
	"trace_filter",
	"trace_transaction",
	"txpool_content",
	"txpool_inspect",
	"txpool_status",
	"web3_clientVersion",
	"web3_sha3",
}
