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
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/big"
	"reflect"
	"testing"
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
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/ethereum/go-ethereum/eth/filters"
	"github.com/ethereum/go-ethereum/eth/tracers"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/params/types/genesisT"
	"github.com/ethereum/go-ethereum/params/vars"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/ethereum/go-ethereum/triedb"
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
	blockHashErr := errors.New("cannot specify both BlockHash and FromBlock/ToBlock")
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
	testBalance = big.NewInt(2e15)
)

var genesis = &genesisT.Genesis{
	Config:    params.AllEthashProtocolChanges,
	Alloc:     genesisT.GenesisAlloc{testAddr: {Balance: testBalance}},
	ExtraData: []byte("test genesis"),
	Timestamp: 9000,
	BaseFee:   big.NewInt(vars.InitialBaseFee),
}

var testTx1 = types.MustSignNewTx(testKey, types.LatestSigner(genesis.Config), &types.LegacyTx{
	Nonce:    0,
	Value:    big.NewInt(12),
	GasPrice: big.NewInt(vars.InitialBaseFee),
	Gas:      vars.TxGas,
	To:       &common.Address{2},
})

var testTx2 = types.MustSignNewTx(testKey, types.LatestSigner(genesis.Config), &types.LegacyTx{
	Nonce:    1,
	Value:    big.NewInt(8),
	GasPrice: big.NewInt(vars.InitialBaseFee),
	Gas:      vars.TxGas,
	To:       &common.Address{2},
})

func newTestBackend(t *testing.T) (*node.Node, []*types.Block) {
	// Generate test chain.
	blocks := generateTestChain()

	// Create node
	n, err := node.New(&node.Config{})
	if err != nil {
		t.Fatalf("can't create new node: %v", err)
	}
	// Create Ethereum Service
	config := &ethconfig.Config{Genesis: genesis}
	config.Ethash.PowMode = ethash.ModeFake
	ethservice, err := eth.New(n, config)
	if err != nil {
		t.Fatalf("can't create new ethereum service: %v", err)
	}
	n.RegisterAPIs(tracers.APIs(ethservice.APIBackend))

	filterSystem := filters.NewFilterSystem(ethservice.APIBackend, filters.Config{})
	n.RegisterAPIs([]rpc.API{{
		Namespace: "eth",
		Service:   filters.NewFilterAPI(filterSystem, false),
	}})
	// Import the test chain.
	if err := n.Start(); err != nil {
		t.Fatalf("can't start test node: %v", err)
	}
	if _, err := ethservice.BlockChain().InsertChain(blocks[1:]); err != nil {
		t.Fatalf("can't import test blocks: %v", err)
	}
	// Ensure the tx indexing is fully generated
	for ; ; time.Sleep(time.Millisecond * 100) {
		progress, err := ethservice.BlockChain().TxIndexProgress()
		if err == nil && progress.Done() {
			break
		}
	}
	return n, blocks
}

// generateTestChain generates 2 blocks. The first block contains 2 transactions.
func generateTestChain() []*types.Block {
	generate := func(i int, g *core.BlockGen) {
		g.OffsetTime(5)
		g.SetExtra([]byte("test"))
		if i == 1 {
			// Test transactions are included in block #2.
			g.AddTx(testTx1)
			g.AddTx(testTx2)
		}
	}
	_, blocks, _ := core.GenerateChainWithGenesis(genesis, ethash.NewFaker(), 2, generate)
	mem := rawdb.NewMemoryDatabase()
	genesisBlock := core.MustCommitGenesis(mem, triedb.NewDatabase(mem, nil), genesis)
	return append([]*types.Block{genesisBlock}, blocks...)
}

func TestEthClient(t *testing.T) {
	backend, chain := newTestBackend(t)
	client := backend.Attach()
	defer backend.Close()
	defer client.Close()

	tests := map[string]struct {
		test func(t *testing.T)
	}{
		"Header": {
			func(t *testing.T) { testHeader(t, chain, client) },
		},
		"BalanceAt": {
			func(t *testing.T) { testBalanceAt(t, client) },
		},
		"TxInBlockInterrupted": {
			func(t *testing.T) { testTransactionInBlock(t, client) },
		},
		"ChainID": {
			func(t *testing.T) { testChainID(t, client) },
		},
		"GetBlock": {
			func(t *testing.T) { testGetBlock(t, client) },
		},
		"StatusFunctions": {
			func(t *testing.T) { testStatusFunctions(t, client) },
		},
		"CallContract": {
			func(t *testing.T) { testCallContract(t, client) },
		},
		"CallContractAtHash": {
			func(t *testing.T) { testCallContractAtHash(t, client) },
		},
		"AtFunctions": {
			func(t *testing.T) { testAtFunctions(t, client) },
		},
		"TransactionSender": {
			func(t *testing.T) { testTransactionSender(t, client) },
		},
	}

	t.Parallel()
	for name, tt := range tests {
		t.Run(name, tt.test)
	}
}

func testHeader(t *testing.T, chain []*types.Block, client *rpc.Client) {
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
			block:   big.NewInt(1000000000),
			want:    nil,
			wantErr: ethereum.NotFound,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ec := NewClient(client)
			ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
			defer cancel()

			got, err := ec.HeaderByNumber(ctx, tt.block)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("HeaderByNumber(%v) error = %q, want %q", tt.block, err, tt.wantErr)
			}
			if got != nil && got.Number != nil && got.Number.Sign() == 0 {
				got.Number = big.NewInt(0) // hack to make DeepEqual work
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("HeaderByNumber(%v) got = %v, want %v", tt.block, got, tt.want)
			}
		})
	}
}

func testBalanceAt(t *testing.T, client *rpc.Client) {
	tests := map[string]struct {
		account common.Address
		block   *big.Int
		want    *big.Int
		wantErr error
	}{
		"valid_account_genesis": {
			account: testAddr,
			block:   big.NewInt(0),
			want:    testBalance,
		},
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

func TestHeader_TxesUnclesNotEmpty(t *testing.T) {
	backend, blocks := newTestBackend(t)
	client := backend.Attach()
	defer backend.Close()
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	res := make(map[string]interface{})
	err := client.CallContext(ctx, &res, "eth_getBlockByNumber", "latest", false)
	if err != nil {
		log.Fatalln(err)
	}

	// Sanity check response
	wantBlocksN := blocks[len(blocks)-1].Number()
	if v, ok := res["number"]; !ok {
		t.Fatal("missing 'number' field")
	} else if n, err := hexutil.DecodeBig(v.(string)); err != nil || n == nil {
		t.Fatal(err)
	} else if n.Cmp(wantBlocksN) != 0 {
		t.Fatalf("unexpected 'latest' block number: %v, want: %d", n, wantBlocksN)
	}
	// 'transactions' key should exist as []
	if v, ok := res["transactions"]; !ok {
		t.Fatal("missing transactions field")
	} else if len(v.([]interface{})) != 2 {
		t.Fatalf("'transactions' value not [], got: %v", len(v.([]interface{})))
	}
	// 'uncles' key should exist as []
	if v, ok := res["uncles"]; !ok {
		t.Fatal("missing uncles field")
	} else if len(v.([]interface{})) != 0 {
		t.Fatal("'uncles' value not []'")
	}
}

func testTransactionInBlock(t *testing.T, client *rpc.Client) {
	ec := NewClient(client)

	// Get current block by number.
	block, err := ec.BlockByNumber(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Test tx in block not found.
	if _, err := ec.TransactionInBlock(context.Background(), block.Hash(), 20); err != ethereum.NotFound {
		t.Fatal("error should be ethereum.NotFound")
	}

	// Test tx in block found.
	tx, err := ec.TransactionInBlock(context.Background(), block.Hash(), 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tx.Hash() != testTx1.Hash() {
		t.Fatalf("unexpected transaction: %v", tx)
	}

	tx, err = ec.TransactionInBlock(context.Background(), block.Hash(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tx.Hash() != testTx2.Hash() {
		t.Fatalf("unexpected transaction: %v", tx)
	}
}

func testChainID(t *testing.T, client *rpc.Client) {
	ec := NewClient(client)
	id, err := ec.ChainID(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id == nil || id.Cmp(params.AllEthashProtocolChanges.ChainID) != 0 {
		t.Fatalf("ChainID returned wrong number: %+v", id)
	}
}

func testGetBlock(t *testing.T, client *rpc.Client) {
	ec := NewClient(client)

	// Get current block number
	blockNumber, err := ec.BlockNumber(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if blockNumber != 2 {
		t.Fatalf("BlockNumber returned wrong number: %d", blockNumber)
	}
	// Get current block by number
	block, err := ec.BlockByNumber(context.Background(), new(big.Int).SetUint64(blockNumber))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if block.NumberU64() != blockNumber {
		t.Fatalf("BlockByNumber returned wrong block: want %d got %d", blockNumber, block.NumberU64())
	}
	// Get current block by hash
	blockH, err := ec.BlockByHash(context.Background(), block.Hash())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if block.Hash() != blockH.Hash() {
		t.Fatalf("BlockByHash returned wrong block: want %v got %v", block.Hash().Hex(), blockH.Hash().Hex())
	}
	// Get header by number
	header, err := ec.HeaderByNumber(context.Background(), new(big.Int).SetUint64(blockNumber))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if block.Header().Hash() != header.Hash() {
		t.Fatalf("HeaderByNumber returned wrong header: want %v got %v", block.Header().Hash().Hex(), header.Hash().Hex())
	}
	// Get header by hash
	headerH, err := ec.HeaderByHash(context.Background(), block.Hash())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if block.Header().Hash() != headerH.Hash() {
		t.Fatalf("HeaderByHash returned wrong header: want %v got %v", block.Header().Hash().Hex(), headerH.Hash().Hex())
	}
}

func testStatusFunctions(t *testing.T, client *rpc.Client) {
	ec := NewClient(client)

	// Sync progress
	progress, err := ec.SyncProgress(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if progress != nil {
		t.Fatalf("unexpected progress: %v", progress)
	}

	// NetworkID
	networkID, err := ec.NetworkID(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if networkID.Cmp(big.NewInt(1337)) != 0 {
		t.Fatalf("unexpected networkID: %v", networkID)
	}

	// SuggestGasPrice
	gasPrice, err := ec.SuggestGasPrice(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gasPrice.Cmp(big.NewInt(1000000000)) != 0 {
		t.Fatalf("unexpected gas price: %v", gasPrice)
	}

	// SuggestGasTipCap
	gasTipCap, err := ec.SuggestGasTipCap(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gasTipCap.Cmp(big.NewInt(234375000)) != 0 {
		t.Fatalf("unexpected gas tip cap: %v", gasTipCap)
	}

	// FeeHistory
	history, err := ec.FeeHistory(context.Background(), 1, big.NewInt(2), []float64{95, 99})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := &ethereum.FeeHistory{
		OldestBlock: big.NewInt(2),
		Reward: [][]*big.Int{
			{
				big.NewInt(234375000),
				big.NewInt(234375000),
			},
		},
		BaseFee: []*big.Int{
			big.NewInt(765625000),
			big.NewInt(671627818),
		},
		GasUsedRatio: []float64{0.008912678667376286},
	}
	if !reflect.DeepEqual(history, want) {
		t.Fatalf("FeeHistory result doesn't match expected: (got: %v, want: %v)", history, want)
	}
}

func testCallContractAtHash(t *testing.T, client *rpc.Client) {
	ec := NewClient(client)

	// EstimateGas
	msg := ethereum.CallMsg{
		From:  testAddr,
		To:    &common.Address{},
		Gas:   21000,
		Value: big.NewInt(1),
	}
	gas, err := ec.EstimateGas(context.Background(), msg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gas != 21000 {
		t.Fatalf("unexpected gas price: %v", gas)
	}
	block, err := ec.HeaderByNumber(context.Background(), big.NewInt(1))
	if err != nil {
		t.Fatalf("BlockByNumber error: %v", err)
	}
	// CallContract
	if _, err := ec.CallContractAtHash(context.Background(), msg, block.Hash()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func testCallContract(t *testing.T, client *rpc.Client) {
	ec := NewClient(client)

	// EstimateGas
	msg := ethereum.CallMsg{
		From:  testAddr,
		To:    &common.Address{},
		Gas:   21000,
		Value: big.NewInt(1),
	}
	gas, err := ec.EstimateGas(context.Background(), msg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gas != 21000 {
		t.Fatalf("unexpected gas price: %v", gas)
	}
	// CallContract
	if _, err := ec.CallContract(context.Background(), msg, big.NewInt(1)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// PendingCallContract
	if _, err := ec.PendingCallContract(context.Background(), msg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func testAtFunctions(t *testing.T, client *rpc.Client) {
	ec := NewClient(client)

	block, err := ec.HeaderByNumber(context.Background(), big.NewInt(1))
	if err != nil {
		t.Fatalf("BlockByNumber error: %v", err)
	}

	// send a transaction for some interesting pending status
	sendTransaction(ec)
	time.Sleep(100 * time.Millisecond)

	// Check pending transaction count
	pending, err := ec.PendingTransactionCount(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pending != 1 {
		t.Fatalf("unexpected pending, wanted 1 got: %v", pending)
	}
	// Query balance
	balance, err := ec.BalanceAt(context.Background(), testAddr, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	hashBalance, err := ec.BalanceAtHash(context.Background(), testAddr, block.Hash())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if balance.Cmp(hashBalance) == 0 {
		t.Fatalf("unexpected balance at hash: %v %v", balance, hashBalance)
	}
	penBalance, err := ec.PendingBalanceAt(context.Background(), testAddr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if balance.Cmp(penBalance) == 0 {
		t.Fatalf("unexpected balance: %v %v", balance, penBalance)
	}
	// NonceAt
	nonce, err := ec.NonceAt(context.Background(), testAddr, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	hashNonce, err := ec.NonceAtHash(context.Background(), testAddr, block.Hash())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if hashNonce == nonce {
		t.Fatalf("unexpected nonce at hash: %v %v", nonce, hashNonce)
	}
	penNonce, err := ec.PendingNonceAt(context.Background(), testAddr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if penNonce != nonce+1 {
		t.Fatalf("unexpected nonce: %v %v", nonce, penNonce)
	}
	// StorageAt
	storage, err := ec.StorageAt(context.Background(), testAddr, common.Hash{}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	hashStorage, err := ec.StorageAtHash(context.Background(), testAddr, common.Hash{}, block.Hash())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !bytes.Equal(storage, hashStorage) {
		t.Fatalf("unexpected storage at hash: %v %v", storage, hashStorage)
	}
	penStorage, err := ec.PendingStorageAt(context.Background(), testAddr, common.Hash{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !bytes.Equal(storage, penStorage) {
		t.Fatalf("unexpected storage: %v %v", storage, penStorage)
	}
	// CodeAt
	code, err := ec.CodeAt(context.Background(), testAddr, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	hashCode, err := ec.CodeAtHash(context.Background(), common.Address{}, block.Hash())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !bytes.Equal(code, hashCode) {
		t.Fatalf("unexpected code at hash: %v %v", code, hashCode)
	}
	penCode, err := ec.PendingCodeAt(context.Background(), testAddr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !bytes.Equal(code, penCode) {
		t.Fatalf("unexpected code: %v %v", code, penCode)
	}
}

func testTransactionSender(t *testing.T, client *rpc.Client) {
	ec := NewClient(client)
	ctx := context.Background()

	// Retrieve testTx1 via RPC.
	block2, err := ec.HeaderByNumber(ctx, big.NewInt(2))
	if err != nil {
		t.Fatal("can't get block 1:", err)
	}
	tx1, err := ec.TransactionInBlock(ctx, block2.Hash(), 0)
	if err != nil {
		t.Fatal("can't get tx:", err)
	}
	if tx1.Hash() != testTx1.Hash() {
		t.Fatalf("wrong tx hash %v, want %v", tx1.Hash(), testTx1.Hash())
	}

	// The sender address is cached in tx1, so no additional RPC should be required in
	// TransactionSender. Ensure the server is not asked by canceling the context here.
	canceledCtx, cancel := context.WithCancel(context.Background())
	cancel()
	<-canceledCtx.Done() // Ensure the close of the Done channel
	sender1, err := ec.TransactionSender(canceledCtx, tx1, block2.Hash(), 0)
	if err != nil {
		t.Fatal(err)
	}
	if sender1 != testAddr {
		t.Fatal("wrong sender:", sender1)
	}

	// Now try to get the sender of testTx2, which was not fetched through RPC.
	// TransactionSender should query the server here.
	sender2, err := ec.TransactionSender(ctx, testTx2, block2.Hash(), 1)
	if err != nil {
		t.Fatal(err)
	}
	if sender2 != testAddr {
		t.Fatal("wrong sender:", sender2)
	}
}

func sendTransaction(ec *Client) error {
	chainID, err := ec.ChainID(context.Background())
	if err != nil {
		return err
	}
	nonce, err := ec.PendingNonceAt(context.Background(), testAddr)
	if err != nil {
		return err
	}

	signer := types.LatestSignerForChainID(chainID)
	tx, err := types.SignNewTx(testKey, signer, &types.LegacyTx{
		Nonce:    nonce,
		To:       &common.Address{2},
		Value:    big.NewInt(1),
		Gas:      22000,
		GasPrice: big.NewInt(vars.InitialBaseFee),
	})
	if err != nil {
		return err
	}
	return ec.SendTransaction(context.Background(), tx)
}

func sliceContains(sl []string, str string) bool {
	for _, s := range sl {
		if str == s {
			return true
		}
	}
	return false
}

func TestRPCDiscover(t *testing.T) {
	check := func(r meta_schema.OpenrpcDocument) {
		responseMethods := func() (names []string) {
			for _, m := range *r.Methods {
				names = append(names, string(*m.Name))
			}
			return
		}()

		over, under := []string{}, []string{}

		// under: methods which exist in the response document,
		// but are not contained in the canonical hardcoded list below
		for _, name := range responseMethods {
			if !sliceContains(allRPCMethods, name) {
				under = append(under, name)
			}
		}

		// over: methods which DO NOT exist in the response document,
		// but ARE contained in the canonical hardcoded list below
		for _, name := range allRPCMethods {
			if !sliceContains(responseMethods, name) {
				over = append(over, name)
			}
		}

		if len(over) > 0 || len(under) > 0 {
			printList := func(list []string) string {
				if len(list) == 0 {
					return "∅" // empty set
				}
				var str string
				for _, s := range list {
					str += "-" + s + "\n"
				}
				return str
			}

			responseDocument, _ := json.MarshalIndent(r, "", "    ")
			t.Logf(`Response Document:

%s`, string(responseDocument))
			t.Fatalf(`OVER (methods which do not appear in the current API, but exist in the hardcoded response document):):
%v

UNDER (methods which appear in the current API, but do not appear in the hardcoded response document):):
%v
`, printList(over), printList(under))
		}
	}

	backend, _ := newTestBackend(t)
	client := backend.Attach()
	defer backend.Close()
	defer client.Close()

	var res meta_schema.OpenrpcDocument
	err := client.Call(&res, "rpc.discover")
	if err != nil {
		t.Fatal(err)
	}

	check(res)
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

	filterSystem := filters.NewFilterSystem(ethservice.APIBackend, filters.Config{})
	backend.RegisterAPIs([]rpc.API{{
		Namespace: "eth",
		Service:   filters.NewFilterAPI(filterSystem, false),
	}})

	// Import the test chain.
	if err := backend.Start(); err != nil {
		t.Fatalf("can't start test node: %v", err)
	}
	if _, err := ethservice.BlockChain().InsertChain(originalBlocks[1:]); err != nil {
		t.Fatalf("can't import test blocks: %v", err)
	}

	// Create the client and newSideHeads subscription.
	client := backend.Attach()
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
		t.Error(err)
	} else {
		defer sub.Unsubscribe()
	}

	headCh := make(chan *types.Header)
	sub2, err2 := ec.SubscribeNewHead(context.Background(), headCh)
	if err2 != nil {
		t.Error(err2)
	} else {
		defer sub2.Unsubscribe()
	}

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
	blocks := generateTestChain()
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
	client := backend.Attach()
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
	client := backend.Attach()
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

/*
--- FAIL: TestRPCDiscover (0.39s)
    ethclient_test.go:798: over: [
personal_signAndSendTransaction
trace_block
trace_call
trace_callMany
trace_filter
trace_subscribe
trace_transaction
trace_unsubscribe
],

under: [
trace_intermediateRoots
trace_standardTraceBadBlockToFile
trace_standardTraceBlockToFile
trace_traceBadBlock
trace_traceBlock
trace_traceBlockByHash
trace_traceBlockByNumber
trace_traceBlockFromFile
trace_traceCall
trace_traceCallMany
trace_traceChain
trace_traceTransaction
]

*/

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
	"admin_peerEvents",
	"admin_removePeer",
	"admin_removeTrustedPeer",
	"admin_startHTTP",
	"admin_startRPC",
	"admin_startWS",
	"admin_stopHTTP",
	"admin_stopRPC",
	"admin_stopWS",
	"debug_accountRange",
	"debug_blockProfile",
	"debug_chaindbCompact",
	"debug_chaindbProperty",
	"debug_cpuProfile",
	"debug_dbAncient",
	"debug_dbAncients",
	"debug_dbGet",
	"debug_discoveryV4Table",
	"debug_dumpBlock",
	"debug_freeOSMemory",
	"debug_gcStats",
	"debug_getAccessibleState",
	"debug_getBadBlocks",
	"debug_getModifiedAccountsByHash",
	"debug_getModifiedAccountsByNumber",
	"debug_getTrieFlushInterval",
	"debug_getRawBlock",
	"debug_getRawHeader",
	"debug_getRawReceipts",
	"debug_getRawTransaction",
	"debug_goTrace",
	"debug_intermediateRoots",
	"debug_memStats",
	"debug_mutexProfile",
	"debug_preimage",
	"debug_printBlock",
	"debug_seedHash",
	"debug_setBlockProfileRate",
	"debug_setGCPercent",
	"debug_setHead",
	"debug_setMutexProfileFraction",
	"debug_setTrieFlushInterval",
	"debug_stacks",
	"debug_standardTraceBadBlockToFile",
	"debug_standardTraceBlockToFile",
	"debug_startCPUProfile",
	"debug_startGoTrace",
	"debug_stopCPUProfile",
	"debug_stopGoTrace",
	"debug_storageRangeAt",
	"debug_subscribe",
	"debug_traceBadBlock",
	"debug_traceBlock",
	"debug_traceBlockByHash",
	"debug_traceBlockByNumber",
	"debug_traceBlockFromFile",
	"debug_traceCall",
	"debug_traceCallMany",
	"debug_traceChain",
	"debug_traceTransaction",
	"debug_unsubscribe",
	"debug_verbosity",
	"debug_vmodule",
	"debug_writeBlockProfile",
	"debug_writeMemProfile",
	"debug_writeMutexProfile",
	"eth_accounts",
	"eth_blockNumber",
	"eth_call",
	"eth_chainId",
	"eth_coinbase",
	"eth_createAccessList",
	"eth_estimateGas",
	"eth_etherbase",
	"eth_feeHistory",
	"eth_fillTransaction",
	"eth_gasPrice",
	"eth_getBalance",
	"eth_getBlockByHash",
	"eth_getBlockByNumber",
	"eth_getBlockReceipts",
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
	"eth_logs",
	"eth_maxPriorityFeePerGas",
	"eth_mining",
	"eth_newBlockFilter",
	"eth_newHeads",
	"eth_newSideBlockFilter",
	"eth_newSideHeads",
	"eth_newFilter",
	"eth_newPendingTransactionFilter",
	"eth_newPendingTransactions",
	"eth_pendingTransactions",
	"eth_resend",
	"eth_sendRawTransaction",
	"eth_sendTransaction",
	"eth_sign",
	"eth_signTransaction",
	"eth_submitHashrate",
	"eth_submitWork",
	"eth_subscribe",
	"eth_syncing",
	"eth_uninstallFilter",
	"eth_unsubscribe",
	"ethash_getHashrate",
	"ethash_getWork",
	"ethash_submitHashrate",
	"ethash_submitWork",
	"miner_setEtherbase",
	"miner_setExtra",
	"miner_setGasLimit",
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
	"personal_signTransaction",
	"personal_unlockAccount",
	"personal_unpair",
	"trace_block",
	"trace_call",
	"trace_callMany",
	"trace_filter",
	"trace_subscribe",
	"trace_transaction",
	"trace_unsubscribe",
	"txpool_content",
	"txpool_contentFrom",
	"txpool_inspect",
	"txpool_status",
	"web3_clientVersion",
	"web3_sha3",
}
