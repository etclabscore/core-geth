// Copyright 2026 The core-geth Authors
// This file is part of the core-geth library.
//
// The core-geth library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The core-geth library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the core-geth library. If not, see <http://www.gnu.org/licenses/>.

// Package keccak256 implements the ECIP-1049 proof-of-work consensus engine.
//
// It is a thin wrapper around an *ethash.Ethash engine: blocks whose number is
// strictly less than the configured transition block are verified and sealed
// by the inner Ethash engine (preserving Etchash/Ethash compatibility for
// historic chain segments); blocks at or after the transition use Keccak-256
// of the SealHash concatenated with the big-endian block nonce as their PoW
// hash, with no DAG/cache. The block header schema is unchanged: MixDigest is
// required to be the zero hash on Keccak-256 blocks, and the 8-byte Nonce
// field is reused to carry the PoW nonce.
package keccak256

import (
	"errors"
	"math/big"
	"math/rand"
	"sync"

	"github.com/ethereum/go-ethereum/consensus/ethash"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/metrics"
)

// errKeccakStopped is returned when an RPC call races with engine shutdown.
var errKeccakStopped = errors.New("keccak256 engine stopped")

// Keccak256 is the consensus engine implementing ECIP-1049.
type Keccak256 struct {
	inner      *ethash.Ethash // Inner Ethash engine used for pre-transition blocks.
	transition *uint64        // ECIP-1049 transition block; nil disables Keccak256 entirely.

	// Local mining state (for post-transition blocks).
	threads  int
	rand     *rand.Rand
	hashrate metrics.Meter
	update   chan struct{}
	lock     sync.Mutex

	// Remote (external) mining state. Started lazily when Seal is invoked
	// for a post-transition block; nil otherwise (and on testers/fakers).
	remote *remoteSealer

	logger log.Logger
}

// New constructs a Keccak256 engine wrapping the supplied inner Ethash engine.
//
// If transition is nil, the wrapper still functions but never activates Keccak;
// every call delegates to the inner Ethash engine.
func New(inner *ethash.Ethash, transition *uint64) *Keccak256 {
	return &Keccak256{
		inner:      inner,
		transition: transition,
		hashrate:   metrics.NewMeter(),
		update:     make(chan struct{}),
		logger:     log.Root(),
	}
}

// NewFaker returns a Keccak256 engine with an inner Ethash faker, useful for
// testing without an actual proof-of-work being required.
func NewFaker() *Keccak256 {
	return New(ethash.NewFaker(), nil)
}

// NewTester returns a Keccak256 engine for unit-tests, wrapping an Ethash tester.
func NewTester() *Keccak256 {
	zero := uint64(0)
	return New(ethash.NewTester(nil, false), &zero)
}

// isPostTransition reports whether the supplied block number is at or after
// the configured ECIP-1049 transition block.
func (k *Keccak256) isPostTransition(number *big.Int) bool {
	if k.transition == nil || number == nil {
		return false
	}
	return number.Uint64() >= *k.transition
}

// Hashrate returns the combined PoW hashrate across the inner Ethash engine
// (pre-transition mining) and the Keccak256 mining loops.
func (k *Keccak256) Hashrate() float64 {
	return k.inner.Hashrate() + k.hashrate.Snapshot().Rate1()
}

// Threads returns the number of mining threads currently configured.
func (k *Keccak256) Threads() int {
	k.lock.Lock()
	defer k.lock.Unlock()
	return k.threads
}

// SetThreads updates the number of local mining threads.
func (k *Keccak256) SetThreads(threads int) {
	k.lock.Lock()
	k.threads = threads
	k.lock.Unlock()
	k.inner.SetThreads(threads)
	select {
	case k.update <- struct{}{}:
	default:
	}
}

// Close terminates the underlying Ethash engine and stops the remote sealer
// if one was started.
func (k *Keccak256) Close() error {
	k.lock.Lock()
	if k.remote != nil {
		select {
		case k.remote.requestExit <- struct{}{}:
		case <-k.remote.exitCh:
		}
		<-k.remote.exitCh
		k.remote = nil
	}
	k.lock.Unlock()
	return k.inner.Close()
}

// ensureRemote starts the remote sealer on first use.
func (k *Keccak256) ensureRemote() *remoteSealer {
	k.lock.Lock()
	defer k.lock.Unlock()
	if k.remote == nil {
		k.remote = startRemoteSealerKeccak(k)
	}
	return k.remote
}
