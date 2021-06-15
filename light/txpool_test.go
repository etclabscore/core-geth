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

package light

import (
	"context"
	"math"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/ethash"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/params/types/genesisT"
	"github.com/ethereum/go-ethereum/params/vars"
)

type testTxRelay struct {
	send, discard, mined chan int
}

func (self *testTxRelay) Send(txs types.Transactions) {
	self.send <- len(txs)
}

func (self *testTxRelay) NewHead(head common.Hash, mined []common.Hash, rollback []common.Hash) {
	m := len(mined)
	if m != 0 {
		self.mined <- m
	}
}

func (self *testTxRelay) Discard(hashes []common.Hash) {
	self.discard <- len(hashes)
}

const poolTestTxsN = 1000
const poolTestBlocksN = 100

// test tx 0..n-1
var testTxSet [poolTestTxsN]*types.Transaction

// txs sent before block i
func sentTx(i int) int {
	return int(math.Pow(float64(i)/float64(poolTestBlocksN), 0.9) * poolTestTxsN)
}

// txs included in block i or before that (minedTx(i) <= sentTx(i))
func minedTx(i int) int {
	return int(math.Pow(float64(i)/float64(poolTestBlocksN), 1.1) * poolTestTxsN)
}

func txPoolTestChainGen(i int, block *core.BlockGen) {
	s := minedTx(i)
	e := minedTx(i + 1)
	for i := s; i < e; i++ {
		block.AddTx(testTxSet[i])
	}
}

func TestTxPool(t *testing.T) {
	for i := range testTxSet {
		testTxSet[i], _ = types.SignTx(types.NewTransaction(uint64(i), acc1Addr, big.NewInt(10000), vars.TxGas, nil, nil), types.HomesteadSigner{}, testBankKey)
	}

	var (
		sdb     = rawdb.NewMemoryDatabase()
		ldb     = rawdb.NewMemoryDatabase()
		gspec   = genesisT.Genesis{Alloc: genesisT.GenesisAlloc{testBankAddress: {Balance: testBankFunds}}}
		genesis = core.MustCommitGenesis(sdb, &gspec)
	)
	core.MustCommitGenesis(ldb, &gspec)
	// Assemble the test environment
	blockchain, _ := core.NewBlockChain(sdb, nil, params.TestChainConfig, ethash.NewFullFaker(), vm.Config{}, nil, nil)
	gchain, _ := core.GenerateChain(params.TestChainConfig, genesis, ethash.NewFaker(), sdb, poolTestBlocksN, txPoolTestChainGen)
	if _, err := blockchain.InsertChain(gchain); err != nil {
		panic(err)
	}

	odr := &testOdr{sdb: sdb, ldb: ldb, indexerConfig: TestClientIndexerConfig}
	relay := &testTxRelay{
		send:    make(chan int, 1),
		discard: make(chan int, 1),
		mined:   make(chan int, 1),
	}
	lightchain, _ := NewLightChain(odr, params.TestChainConfig, ethash.NewFullFaker(), nil)
	txPermanent = 50
	pool := NewTxPool(params.TestChainConfig, lightchain, relay)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	for ii, block := range gchain {
		i := ii + 1
		s := sentTx(i - 1)
		e := sentTx(i)
		for i := s; i < e; i++ {
			pool.Add(ctx, testTxSet[i])
			got := <-relay.send
			exp := 1
			if got != exp {
				t.Errorf("relay.Send expected len = %d, got %d", exp, got)
			}
		}

		// core-geth
		if ii == len(gchain)/4 {
			// Fuck up pool head
			// This is an edge case that I'm not sure could really happen (hopefully not),
			// but checking anyways. Call it sanity.
			t.Log("Setting pool head to empty hash")
			pool.head = common.Hash{}
		}

		// core-geth
		if ii == len(gchain)/2 {
			// Attempt to insert a nil header into the headerchain
			// NOTE(ia)
			t.Log("Inserting nil header into header chain")
			if _, err := lightchain.InsertHeaderChain([]*types.Header{nil}, 1); err == nil {
				t.Fatal("insert nil header error should not be errorless")
			}
		}

		// core-geth
		if ii == len(gchain)/4*3 {
			var h *types.Header
			t.Log("Setting pool head to a nil header", h.Hash().Hex())
			pool.setNewHead(h)
		}

		if _, err := lightchain.InsertHeaderChain([]*types.Header{block.Header()}, 1); err != nil {
			t.Fatal(err)
		}

		got := <-relay.mined
		exp := minedTx(i) - minedTx(i-1)
		if got != exp {
			t.Errorf("relay.NewHead expected len(mined) = %d, got %d", exp, got)
		}

		exp = 0
		if i > int(txPermanent)+1 {
			exp = minedTx(i-int(txPermanent)-1) - minedTx(i-int(txPermanent)-2)
		}
		if exp != 0 {
			got = <-relay.discard
			if got != exp {
				t.Errorf("relay.Discard expected len = %d, got %d", exp, got)
			}
		}
	}
}
