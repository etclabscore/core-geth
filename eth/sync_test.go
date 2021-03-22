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
	"sync/atomic"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/eth/downloader"
	"github.com/ethereum/go-ethereum/eth/protocols/eth"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/ethereum/go-ethereum/params"
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

// Tests that fast sync is disabled after a successful sync cycle.
func TestFastSyncDisabling64(t *testing.T) { testFastSyncDisabling(t, 64) }
func TestFastSyncDisabling65(t *testing.T) { testFastSyncDisabling(t, 65) }

// Tests that fast sync gets disabled as soon as a real block is successfully
// imported into the blockchain.
func testFastSyncDisabling(t *testing.T, protocol uint) {
	t.Parallel()

	// Create an empty handler and ensure it's in fast sync mode
	empty := newTestHandler()
	if atomic.LoadUint32(&empty.handler.fastSync) == 0 {
		t.Fatalf("fast sync disabled on pristine blockchain")
	}
	defer empty.close()

	// Create a full handler and ensure fast sync ends up disabled
	full := newTestHandlerWithBlocks(1024)
	if atomic.LoadUint32(&full.handler.fastSync) == 1 {
		t.Fatalf("fast sync not disabled on non-empty blockchain")
	}
	defer full.close()

	// Sync up the two handlers
	emptyPipe, fullPipe := p2p.MsgPipe()
	defer emptyPipe.Close()
	defer fullPipe.Close()

	emptyPeer := eth.NewPeer(protocol, p2p.NewPeer(enode.ID{1}, "", nil), emptyPipe, empty.txpool)
	fullPeer := eth.NewPeer(protocol, p2p.NewPeer(enode.ID{2}, "", nil), fullPipe, full.txpool)
	defer emptyPeer.Close()
	defer fullPeer.Close()

	go empty.handler.runEthPeer(emptyPeer, func(peer *eth.Peer) error {
		return eth.Handle((*ethHandler)(empty.handler), peer)
	})
	go full.handler.runEthPeer(fullPeer, func(peer *eth.Peer) error {
		return eth.Handle((*ethHandler)(full.handler), peer)
	})
	// Wait a bit for the above handlers to start
	time.Sleep(250 * time.Millisecond)

	// Check that fast sync was disabled
	op := peerToSyncOp(downloader.FastSync, empty.handler.peers.peerWithHighestTD())
	if err := empty.handler.doSync(op); err != nil {
		t.Fatal("sync failed:", err)
	}
	if atomic.LoadUint32(&empty.handler.fastSync) == 1 {
		t.Fatalf("fast sync not disabled after successful synchronisation")
	}
}

func TestArtificialFinalityFeatureEnablingDisabling(t *testing.T) {
	maxBlocksCreated := 1024
	genFunc := blockGenContemporaryTime(int64(maxBlocksCreated))

	// Create a full protocol manager, check that fast sync gets disabled
	a, _ := newTestProtocolManagerMust(t, downloader.FastSync, maxBlocksCreated, genFunc, nil)
	if atomic.LoadUint32(&a.fastSync) == 1 {
		t.Fatalf("fast sync not disabled on non-empty blockchain")
	}

	one := uint64(1)
	a.blockchain.Config().SetECBP1100Transition(&one)

	oMinAFPeers := minArtificialFinalityPeers
	defer func() {
		// Clean up after, resetting global default to original valu.
		minArtificialFinalityPeers = oMinAFPeers
	}()
	minArtificialFinalityPeers = 1

	// Create a full protocol manager, check that fast sync gets disabled
	b, _ := newTestProtocolManagerMust(t, downloader.FastSync, 0, genFunc, nil)
	if atomic.LoadUint32(&b.fastSync) == 0 {
		t.Fatalf("fast sync disabled on pristine blockchain")
	}
	b.blockchain.Config().SetECBP1100Transition(&one)
	// b.chainSync.forced = true

	io1, io2 := p2p.MsgPipe()
	go a.handle(a.newPeer(65, p2p.NewPeer(enode.ID{}, "peer-b", nil), io2, a.txpool.Get))
	go b.handle(b.newPeer(65, p2p.NewPeer(enode.ID{}, "peer-a", nil), io1, b.txpool.Get))
	time.Sleep(250 * time.Millisecond)

	op := peerToSyncOp(downloader.FullSync, b.peers.BestPeer())
	if err := b.doSync(op); err != nil {
		t.Fatalf("sync failed: %v", err)
	}

	b.chainSync.forced = true
	next := b.chainSync.nextSyncOp()
	if next != nil {
		t.Fatal("non-nil next sync op")
	}
	if !b.blockchain.Config().IsEnabled(b.blockchain.Config().GetECBP1100Transition, b.blockchain.CurrentBlock().Number()) {
		t.Error("AF feature not configured")
	}
	if !b.blockchain.IsArtificialFinalityEnabled() {
		t.Error("AF not enabled")
	}

	// Set the value back to default (more than 1).
	minArtificialFinalityPeers = oMinAFPeers

	// Next sync op will unset AF because manager only has 1 peer.
	b.chainSync.forced = true
	next = b.chainSync.nextSyncOp()
	if next != nil {
		t.Fatal("non-nil next sync op")
	}
	if b.blockchain.IsArtificialFinalityEnabled() {
		t.Error("AF not disabled")
	}
}

// TestArtificialFinalityFeatureEnablingDisabling_NoDisable tests that when the nodisable override
// is in place (see NOTE1 below), AF is not disabled at the min peer floor.
func TestArtificialFinalityFeatureEnablingDisabling_NoDisable(t *testing.T) {
	maxBlocksCreated := 1024
	genFunc := blockGenContemporaryTime(int64(maxBlocksCreated))

	// Create a full protocol manager, check that fast sync gets disabled
	a, _ := newTestProtocolManagerMust(t, downloader.FastSync, maxBlocksCreated, genFunc, nil)
	if atomic.LoadUint32(&a.fastSync) == 1 {
		t.Fatalf("fast sync not disabled on non-empty blockchain")
	}

	one := uint64(1)
	a.blockchain.Config().SetECBP1100Transition(&one)

	oMinAFPeers := minArtificialFinalityPeers
	defer func() {
		// Clean up after, resetting global default to original value.
		minArtificialFinalityPeers = oMinAFPeers
	}()
	minArtificialFinalityPeers = 1

	// Create a full protocol manager, check that fast sync gets disabled
	b, _ := newTestProtocolManagerMust(t, downloader.FastSync, 0, genFunc, nil)
	if atomic.LoadUint32(&b.fastSync) == 0 {
		t.Fatalf("fast sync disabled on pristine blockchain")
	}
	b.blockchain.Config().SetECBP1100Transition(&one)

	// NOTE1: Set the nodisable switch to the on position.
	// This prevents the de-activation of AF features once they have been enabled.
	// In geth code, this is set in the backend during blockchain construction.
	b.blockchain.ArtificialFinalityNoDisable(1)

	io1, io2 := p2p.MsgPipe()
	go a.handle(a.newPeer(65, p2p.NewPeer(enode.ID{}, "peer-b", nil), io2, a.txpool.Get))
	go b.handle(b.newPeer(65, p2p.NewPeer(enode.ID{}, "peer-a", nil), io1, b.txpool.Get))
	time.Sleep(250 * time.Millisecond)

	op := peerToSyncOp(downloader.FullSync, b.peers.BestPeer())
	if err := b.doSync(op); err != nil {
		t.Fatalf("sync failed: %v", err)
	}

	b.chainSync.forced = true
	next := b.chainSync.nextSyncOp()

	// Revert safety condition overrides to default values.
	// Set the value back to default (more than 1).
	minArtificialFinalityPeers = oMinAFPeers

	if next != nil {
		t.Fatal("non-nil next sync op")
	}
	if !b.blockchain.Config().IsEnabled(b.blockchain.Config().GetECBP1100Transition, b.blockchain.CurrentBlock().Number()) {
		t.Error("AF feature not configured")
	}
	if !b.blockchain.IsArtificialFinalityEnabled() {
		t.Error("AF not enabled")
	}

	// Next sync op will unset AF because manager only has 1 peer.
	b.chainSync.forced = true
	next = b.chainSync.nextSyncOp()
	if next != nil {
		t.Fatal("non-nil next sync op")
	}
	if !b.blockchain.IsArtificialFinalityEnabled() {
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
	a, _ := newTestProtocolManagerMust(t, downloader.FastSync, maxBlocksCreated, nil, nil)
	if atomic.LoadUint32(&a.fastSync) == 1 {
		t.Fatalf("fast sync not disabled on non-empty blockchain")
	}

	one := uint64(1)
	a.blockchain.Config().SetECBP1100Transition(&one)

	oMinAFPeers := minArtificialFinalityPeers
	defer func() {
		// Clean up after, resetting global default to original value.
		minArtificialFinalityPeers = oMinAFPeers
	}()
	minArtificialFinalityPeers = 1

	// Create a full protocol manager, check that fast sync gets disabled
	b, _ := newTestProtocolManagerMust(t, downloader.FastSync, 0, nil, nil)
	if atomic.LoadUint32(&b.fastSync) == 0 {
		t.Fatalf("fast sync disabled on pristine blockchain")
	}
	b.blockchain.Config().SetECBP1100Transition(&one)

	io1, io2 := p2p.MsgPipe()
	go a.handle(a.newPeer(65, p2p.NewPeer(enode.ID{}, "peer-b", nil), io2, a.txpool.Get))
	go b.handle(b.newPeer(65, p2p.NewPeer(enode.ID{}, "peer-a", nil), io1, b.txpool.Get))
	time.Sleep(250 * time.Millisecond)

	op := peerToSyncOp(downloader.FullSync, b.peers.BestPeer())
	if err := b.doSync(op); err != nil {
		t.Fatalf("sync failed: %v", err)
	}

	b.chainSync.forced = true
	next := b.chainSync.nextSyncOp()

	// Revert safety condition overrides to default values.
	// Set the value back to default (more than 1).
	minArtificialFinalityPeers = oMinAFPeers

	if next != nil {
		t.Fatal("non-nil next sync op")
	}
	if !b.blockchain.Config().IsEnabled(b.blockchain.Config().GetECBP1100Transition, b.blockchain.CurrentBlock().Number()) {
		t.Error("AF feature not configured")
	}
	// Unit test the timestamp. We want to be sure that the blockchain's current header is actually very old
	// (since we expect its age to act as a condition preventing the enabling of AF features).
	d := uint64(time.Now().Unix()) - b.blockchain.CurrentHeader().Time
	if time.Second*time.Duration(d) < artificialFinalitySafetyInterval {
		t.Errorf("Expected blockchain current head to be very old, it is not.")
	}
	if b.blockchain.IsArtificialFinalityEnabled() {
		t.Error("AF is enabled despite a very stale head")
	}
}

func TestArtificialFinalitySafetyLoopTimeComparison(t *testing.T) {
	if !(time.Since(time.Unix(int64(params.DefaultMessNetGenesisBlock().Timestamp), 0)) > artificialFinalitySafetyInterval) {
		t.Fatal("bad unit logic!")
	}
}
