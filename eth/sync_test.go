// Copyright 2015 The go-ethereum Authors
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

package eth

import (
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/consensus/ethash"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/eth/downloader"
	"github.com/ethereum/go-ethereum/eth/protocols/eth"
	"github.com/ethereum/go-ethereum/eth/protocols/snap"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/params/types/genesisT"
	"github.com/ethereum/go-ethereum/triedb"
)

// blockGenContemporaryTime creates a block gen function that will bump the block times to within throwing
// distance of the current time.
// Without this, block times are built on a genesis with a 0 (or 10) timestamp, which is in 1970.
// This causes the apparent age of the chain to exceed the safety interval and thus interferes with
// reasonable expectations of staleness.
func blockGenContemporaryTime(numBlocks int64) func(i int, gen *core.BlockGen) {
	startTimeUnix := time.Now().Unix()
	return func(i int, gen *core.BlockGen) {
		if i == 0 {
			// Here, 1024 is the max number of block we'll create.
			// 10 is the magic number that the test helper block gen will use for the block time offsets.
			gen.OffsetTime(startTimeUnix - (numBlocks * 10))
		}
	}
}

// newTestHandlerWithBlocks creates a new handler for testing purposes, with a
// given number of initial blocks.
func newTestHandlerWithBlocksWithOpts(blocks int, mode downloader.SyncMode, gen func(int, *core.BlockGen)) *testHandler {
	// Create a database pre-initialize with a genesis block
	db := rawdb.NewMemoryDatabase()
	gspec := &genesisT.Genesis{
		Config: params.TestChainConfig,
		Alloc:  genesisT.GenesisAlloc{testAddr: {Balance: big.NewInt(1000000)}},
	}
	core.MustCommitGenesis(db, triedb.NewDatabase(db, nil), gspec)

	chain, _ := core.NewBlockChain(db, nil, gspec, nil, ethash.NewFaker(), vm.Config{}, nil, nil)

	bs, _ := core.GenerateChain(params.TestChainConfig, chain.Genesis(), ethash.NewFaker(), db, blocks, gen)
	if _, err := chain.InsertChain(bs); err != nil {
		panic(err)
	}
	txpool := newTestTxPool()

	handler, _ := newHandler(&handlerConfig{
		Database:   db,
		Merger:     consensus.NewMerger(db),
		Chain:      chain,
		TxPool:     txpool,
		Network:    1,
		Sync:       mode,
		BloomCache: 1,
	})
	handler.Start(1000)

	return &testHandler{
		db:      db,
		chain:   chain,
		txpool:  txpool,
		handler: handler,
	}
}

// Tests that snap sync is disabled after a successful sync cycle.
func TestSnapSyncDisabling68(t *testing.T) { testSnapSyncDisabling(t, eth.ETH68, snap.SNAP1) }

// Tests that snap sync gets disabled as soon as a real block is successfully
// imported into the blockchain.
func testSnapSyncDisabling(t *testing.T, ethVer uint, snapVer uint) {
	t.Parallel()

	// Create an empty handler and ensure it's in snap sync mode
	empty := newTestHandler()
	if !empty.handler.snapSync.Load() {
		t.Fatalf("snap sync disabled on pristine blockchain")
	}
	defer empty.close()

	// Create a full handler and ensure snap sync ends up disabled
	full := newTestHandlerWithBlocks(1024)
	if full.handler.snapSync.Load() {
		t.Fatalf("snap sync not disabled on non-empty blockchain")
	}
	defer full.close()

	// Sync up the two handlers via both `eth` and `snap`
	caps := []p2p.Cap{{Name: "eth", Version: ethVer}, {Name: "snap", Version: snapVer}}

	emptyPipeEth, fullPipeEth := p2p.MsgPipe()
	defer emptyPipeEth.Close()
	defer fullPipeEth.Close()

	emptyPeerEth := eth.NewPeer(ethVer, p2p.NewPeer(enode.ID{1}, "", caps), emptyPipeEth, empty.txpool)
	fullPeerEth := eth.NewPeer(ethVer, p2p.NewPeer(enode.ID{2}, "", caps), fullPipeEth, full.txpool)
	defer emptyPeerEth.Close()
	defer fullPeerEth.Close()

	go empty.handler.runEthPeer(emptyPeerEth, func(peer *eth.Peer) error {
		return eth.Handle((*ethHandler)(empty.handler), peer)
	})
	go full.handler.runEthPeer(fullPeerEth, func(peer *eth.Peer) error {
		return eth.Handle((*ethHandler)(full.handler), peer)
	})

	emptyPipeSnap, fullPipeSnap := p2p.MsgPipe()
	defer emptyPipeSnap.Close()
	defer fullPipeSnap.Close()

	emptyPeerSnap := snap.NewPeer(snapVer, p2p.NewPeer(enode.ID{1}, "", caps), emptyPipeSnap)
	fullPeerSnap := snap.NewPeer(snapVer, p2p.NewPeer(enode.ID{2}, "", caps), fullPipeSnap)

	go empty.handler.runSnapExtension(emptyPeerSnap, func(peer *snap.Peer) error {
		return snap.Handle((*snapHandler)(empty.handler), peer)
	})
	go full.handler.runSnapExtension(fullPeerSnap, func(peer *snap.Peer) error {
		return snap.Handle((*snapHandler)(full.handler), peer)
	})
	// Wait a bit for the above handlers to start
	time.Sleep(250 * time.Millisecond)

	// Check that snap sync was disabled
	op := peerToSyncOp(downloader.SnapSync, empty.handler.peers.peerWithHighestTD())
	if err := empty.handler.doSync(op); err != nil {
		t.Fatal("sync failed:", err)
	}
	if empty.handler.snapSync.Load() {
		t.Fatalf("snap sync not disabled after successful synchronisation")
	}
}

func TestArtificialFinalityFeatureEnablingDisabling(t *testing.T) {
	maxBlocksCreated := 1024
	genFunc := blockGenContemporaryTime(int64(maxBlocksCreated))

	// Create a full protocol manager, check that fast sync gets disabled
	a := newTestHandlerWithBlocksWithOpts(1024, downloader.FullSync, genFunc)
	if a.handler.snapSync.Load() {
		t.Fatalf("snap sync not disabled on non-empty blockchain")
	}
	defer a.close()

	one := uint64(1)
	a.chain.Config().SetECBP1100Transition(&one)

	oMinAFPeers := minArtificialFinalityPeers
	defer func() {
		// Clean up after, resetting global default to original value.
		minArtificialFinalityPeers = oMinAFPeers
	}()
	minArtificialFinalityPeers = 1

	// Create a full protocol manager, check that fast sync gets disabled
	b := newTestHandlerWithBlocksWithOpts(0, downloader.FullSync, genFunc)
	if b.handler.snapSync.Load() {
		t.Fatalf("snap sync not disabled on non-empty blockchain")
	}
	defer b.close()
	b.chain.Config().SetECBP1100Transition(&one)
	// b.chainSync.forced = true

	// Sync up the two handlers
	emptyPipe, fullPipe := p2p.MsgPipe()
	defer emptyPipe.Close()
	defer fullPipe.Close()

	fullPeer := eth.NewPeer(66, p2p.NewPeer(enode.ID{2}, "", nil), fullPipe, a.txpool)
	emptyPeer := eth.NewPeer(66, p2p.NewPeer(enode.ID{1}, "", nil), emptyPipe, b.txpool)
	defer emptyPeer.Close()
	defer fullPeer.Close()

	go b.handler.runEthPeer(emptyPeer, func(peer *eth.Peer) error {
		return eth.Handle((*ethHandler)(b.handler), peer)
	})
	go a.handler.runEthPeer(fullPeer, func(peer *eth.Peer) error {
		return eth.Handle((*ethHandler)(a.handler), peer)
	})
	// Wait a bit for the above handlers to start
	time.Sleep(250 * time.Millisecond)

	op := peerToSyncOp(downloader.FullSync, b.handler.peers.peerWithHighestTD())
	if err := b.handler.doSync(op); err != nil {
		t.Fatalf("sync failed: %v", err)
	}

	b.handler.chainSync.forced = true
	next := b.handler.chainSync.nextSyncOp()
	if next != nil {
		t.Fatal("non-nil next sync op")
	}
	if !b.chain.Config().IsEnabled(b.chain.Config().GetECBP1100Transition, b.chain.CurrentBlock().Number) {
		t.Error("AF feature not configured")
	}
	if !b.chain.IsArtificialFinalityEnabled() {
		t.Error("AF not enabled")
	}

	// Set the value back to default (more than 1).
	minArtificialFinalityPeers = oMinAFPeers

	// Next sync op will unset AF because manager only has 1 peer.
	b.handler.chainSync.forced = true
	next = b.handler.chainSync.nextSyncOp()
	if next != nil {
		t.Fatal("non-nil next sync op")
	}
	if b.chain.IsArtificialFinalityEnabled() {
		t.Error("AF not disabled")
	}
}

// TestArtificialFinalityFeatureEnablingDisabling_NoDisable tests that when the nodisable override
// is in place (see NOTE1 below), AF is not disabled at the min peer floor.
func TestArtificialFinalityFeatureEnablingDisabling_NoDisable(t *testing.T) {
	maxBlocksCreated := 1024
	genFunc := blockGenContemporaryTime(int64(maxBlocksCreated))

	// Create a full protocol manager, check that fast sync gets disabled
	a := newTestHandlerWithBlocksWithOpts(1024, downloader.FullSync, genFunc)
	defer a.close()

	one := uint64(1)
	a.chain.Config().SetECBP1100Transition(&one)

	oMinAFPeers := minArtificialFinalityPeers
	defer func() {
		// Clean up after, resetting global default to original value.
		minArtificialFinalityPeers = oMinAFPeers
	}()
	minArtificialFinalityPeers = 1

	// Create a full protocol manager, check that fast sync gets disabled
	b := newTestHandlerWithBlocksWithOpts(0, downloader.FullSync, genFunc)
	defer b.close()
	b.chain.Config().SetECBP1100Transition(&one)

	// NOTE1: Set the nodisable switch to the on position.
	// This prevents the de-activation of AF features once they have been enabled.
	// In geth code, this is set in the backend during blockchain construction.
	b.chain.ArtificialFinalityNoDisable(1)

	// Sync up the two handlers
	emptyPipe, fullPipe := p2p.MsgPipe()
	defer emptyPipe.Close()
	defer fullPipe.Close()

	fullPeer := eth.NewPeer(66, p2p.NewPeer(enode.ID{2}, "", nil), fullPipe, a.txpool)
	emptyPeer := eth.NewPeer(66, p2p.NewPeer(enode.ID{1}, "", nil), emptyPipe, b.txpool)
	defer emptyPeer.Close()
	defer fullPeer.Close()

	go b.handler.runEthPeer(emptyPeer, func(peer *eth.Peer) error {
		return eth.Handle((*ethHandler)(b.handler), peer)
	})
	go a.handler.runEthPeer(fullPeer, func(peer *eth.Peer) error {
		return eth.Handle((*ethHandler)(a.handler), peer)
	})
	// Wait a bit for the above handlers to start
	time.Sleep(250 * time.Millisecond)

	op := peerToSyncOp(downloader.FullSync, b.handler.peers.peerWithHighestTD())
	if err := b.handler.doSync(op); err != nil {
		t.Fatalf("sync failed: %v", err)
	}

	b.handler.chainSync.forced = true
	next := b.handler.chainSync.nextSyncOp()

	// Revert safety condition overrides to default values.
	// Set the value back to default (more than 1).
	minArtificialFinalityPeers = oMinAFPeers

	if next != nil {
		t.Fatal("non-nil next sync op")
	}
	if !b.chain.Config().IsEnabled(b.chain.Config().GetECBP1100Transition, b.chain.CurrentBlock().Number) {
		t.Error("AF feature not configured")
	}
	if !b.chain.IsArtificialFinalityEnabled() {
		t.Error("AF not enabled")
	}

	// Next sync op will unset AF because manager only has 1 peer.
	b.handler.chainSync.forced = true
	next = b.handler.chainSync.nextSyncOp()
	if next != nil {
		t.Fatal("non-nil next sync op")
	}
	if !b.chain.IsArtificialFinalityEnabled() {
		t.Error(`AF not enabled;
The AF disable mechanism triggered by the minimum peers floor should have been short-circuited,
preventing AF disablement on this sync op (with the minArtificialFinalityPeers value set to > 1 (defaultMinSyncPeers = 5)),
and the number of peers 'a' is connected with being only 1.)`)
	}
}

// TestArtificialFinalityFeatureEnablingDisabling_StaleHead tests that the stale head condition is respected.
// The block generator function will yield a chain with a head, which because of it time=0 (aka year 1970) genesis, will
// be very old (far exceeding the auto-disable stale limit).
// In this case, AF features should NOT be enabled.
func TestArtificialFinalityFeatureEnablingDisabling_StaleHead(t *testing.T) {
	maxBlocksCreated := 1024

	// Create a full protocol manager, check that fast sync gets disabled
	a := newTestHandlerWithBlocksWithOpts(maxBlocksCreated, downloader.FullSync, nil)
	defer a.close()
	one := uint64(1)
	a.chain.Config().SetECBP1100Transition(&one)

	oMinAFPeers := minArtificialFinalityPeers
	defer func() {
		// Clean up after, resetting global default to original value.
		minArtificialFinalityPeers = oMinAFPeers
	}()
	minArtificialFinalityPeers = 1

	// Create a full protocol manager, check that fast sync gets disabled
	b := newTestHandlerWithBlocksWithOpts(0, downloader.FullSync, nil)
	defer b.close()
	b.chain.Config().SetECBP1100Transition(&one)

	// Sync up the two handlers
	emptyPipe, fullPipe := p2p.MsgPipe()
	defer emptyPipe.Close()
	defer fullPipe.Close()

	fullPeer := eth.NewPeer(66, p2p.NewPeer(enode.ID{2}, "", nil), fullPipe, a.txpool)
	emptyPeer := eth.NewPeer(66, p2p.NewPeer(enode.ID{1}, "", nil), emptyPipe, b.txpool)
	defer emptyPeer.Close()
	defer fullPeer.Close()

	go b.handler.runEthPeer(emptyPeer, func(peer *eth.Peer) error {
		return eth.Handle((*ethHandler)(b.handler), peer)
	})
	go a.handler.runEthPeer(fullPeer, func(peer *eth.Peer) error {
		return eth.Handle((*ethHandler)(a.handler), peer)
	})
	// Wait a bit for the above handlers to start
	time.Sleep(250 * time.Millisecond)

	op := peerToSyncOp(downloader.FullSync, b.handler.peers.peerWithHighestTD())
	if err := b.handler.doSync(op); err != nil {
		t.Fatalf("sync failed: %v", err)
	}

	b.handler.chainSync.forced = true
	next := b.handler.chainSync.nextSyncOp()

	// Revert safety condition overrides to default values.
	// Set the value back to default (more than 1).
	minArtificialFinalityPeers = oMinAFPeers

	if next != nil {
		t.Fatal("non-nil next sync op")
	}
	if !b.chain.Config().IsEnabled(b.chain.Config().GetECBP1100Transition, b.chain.CurrentBlock().Number) {
		t.Error("AF feature not configured")
	}
	// Unit test the timestamp. We want to be sure that the blockchain's current header is actually very old
	// (since we expect its age to act as a condition preventing the enabling of AF features).
	d := uint64(time.Now().Unix()) - b.chain.CurrentHeader().Time
	if time.Second*time.Duration(d) < artificialFinalitySafetyInterval {
		t.Errorf("Expected blockchain current head to be very old, it is not.")
	}
	if b.chain.IsArtificialFinalityEnabled() {
		t.Error("AF is enabled despite a very stale head")
	}
}

func TestArtificialFinalitySafetyLoopTimeComparison(t *testing.T) {
	if !(time.Since(time.Unix(int64(params.DefaultMessNetGenesisBlock().Timestamp), 0)) > artificialFinalitySafetyInterval) {
		t.Fatal("bad unit logic!")
	}
}
