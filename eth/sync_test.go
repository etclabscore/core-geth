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
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/eth/downloader"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/enode"
)

func TestFastSyncDisabling63(t *testing.T) { testFastSyncDisabling(t, 63) }
func TestFastSyncDisabling64(t *testing.T) { testFastSyncDisabling(t, 64) }
func TestFastSyncDisabling65(t *testing.T) { testFastSyncDisabling(t, 65) }

// Tests that fast sync gets disabled as soon as a real block is successfully
// imported into the blockchain.
func testFastSyncDisabling(t *testing.T, protocol int) {
	t.Parallel()

	// Create a pristine protocol manager, check that fast sync is left enabled
	pmEmpty, _ := newTestProtocolManagerMust(t, downloader.FastSync, 0, nil, nil)
	if atomic.LoadUint32(&pmEmpty.fastSync) == 0 {
		t.Fatalf("fast sync disabled on pristine blockchain")
	}
	// Create a full protocol manager, check that fast sync gets disabled
	pmFull, _ := newTestProtocolManagerMust(t, downloader.FastSync, 1024, nil, nil)
	if atomic.LoadUint32(&pmFull.fastSync) == 1 {
		t.Fatalf("fast sync not disabled on non-empty blockchain")
	}

	// Sync up the two peers
	io1, io2 := p2p.MsgPipe()
	go pmFull.handle(pmFull.newPeer(protocol, p2p.NewPeer(enode.ID{}, "empty", nil), io2, pmFull.txpool.Get))
	go pmEmpty.handle(pmEmpty.newPeer(protocol, p2p.NewPeer(enode.ID{}, "full", nil), io1, pmEmpty.txpool.Get))

	time.Sleep(250 * time.Millisecond)
	op := peerToSyncOp(downloader.FastSync, pmEmpty.peers.BestPeer())
	if err := pmEmpty.doSync(op); err != nil {
		t.Fatal("sync failed:", err)
	}

	// Check that fast sync was disabled
	if atomic.LoadUint32(&pmEmpty.fastSync) == 1 {
		t.Fatalf("fast sync not disabled after successful synchronisation")
	}
}

func TestArtificialFinalityFeatureEnablingDisabling(t *testing.T) {
	// Create a full protocol manager, check that fast sync gets disabled
	a, _ := newTestProtocolManagerMust(t, downloader.FastSync, 1024, nil, nil)
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
	b, _ := newTestProtocolManagerMust(t, downloader.FastSync, 0, nil, nil)
	if atomic.LoadUint32(&b.fastSync) == 0 {
		t.Fatalf("fast sync disabled on pristine blockchain")
	}
	b.blockchain.Config().SetECBP1100Transition(&one)
	// b.chainSync.forced = true

	io1, io2 := p2p.MsgPipe()
	go a.handle(a.newPeer(65, p2p.NewPeer(enode.ID{}, fmt.Sprintf("peer-b"), nil), io2, a.txpool.Get))
	go b.handle(b.newPeer(65, p2p.NewPeer(enode.ID{}, fmt.Sprintf("peer-a"), nil), io1, b.txpool.Get))
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
