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

package keccak256

import (
	"encoding/binary"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/ethash"
	"github.com/ethereum/go-ethereum/core/types"
	"golang.org/x/crypto/sha3"
)

// TestRemoteSealEndToEnd drives a full Seal -> GetWork -> external-mine ->
// SubmitWork cycle entirely in process and asserts that the engine emits a
// valid sealed block back into the results channel.
func TestRemoteSealEndToEnd(t *testing.T) {
	transition := uint64(0)
	engine := New(ethash.NewFaker(), &transition)
	engine.SetThreads(-1) // disable local mining; remote-only
	defer engine.Close()

	api := &API{keccak: engine}

	header := &types.Header{
		ParentHash: common.Hash{},
		Number:     big.NewInt(1),
		Difficulty: big.NewInt(1 << 16),
		GasLimit:   5_000_000,
		Time:       1_700_000_000,
	}
	block := types.NewBlockWithHeader(header)

	results := make(chan *types.Block, 1)
	stop := make(chan struct{})
	defer close(stop)

	if err := engine.Seal(nil, block, results, stop); err != nil {
		t.Fatalf("Seal: %v", err)
	}

	// Pull work out of the public API.
	var work [3]string
	deadline := time.Now().Add(2 * time.Second)
	for {
		var err error
		work, err = api.GetWork()
		if err == nil {
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("GetWork never returned: %v", err)
		}
		time.Sleep(5 * time.Millisecond)
	}
	sealHash := common.HexToHash(work[0])
	if sealHash != engine.SealHash(header) {
		t.Fatalf("unexpected sealhash: %x vs %x", sealHash, engine.SealHash(header))
	}
	target := new(big.Int).SetBytes(common.HexToHash(work[1]).Bytes())

	// Mine externally: search for a satisfying nonce.
	hasher := sha3.NewLegacyKeccak256()
	var buf [8]byte
	winning := uint64(0)
	for n := uint64(0); n < 1<<22; n++ {
		hasher.Reset()
		hasher.Write(sealHash[:])
		binary.BigEndian.PutUint64(buf[:], n)
		hasher.Write(buf[:])
		out := hasher.Sum(nil)
		if new(big.Int).SetBytes(out).Cmp(target) <= 0 {
			winning = n
			break
		}
	}
	if winning == 0 {
		// Either we got lucky on n=0, or we never found one. Verify the n=0
		// case below; otherwise this test would be flaky at extreme diffs.
	}

	if !api.SubmitWork(types.EncodeNonce(winning), sealHash) {
		t.Fatalf("SubmitWork rejected a valid solution (nonce=%d)", winning)
	}

	select {
	case sealed := <-results:
		if sealed.Nonce() != winning {
			t.Fatalf("sealed block has nonce %d, want %d", sealed.Nonce(), winning)
		}
		if sealed.MixDigest() != (common.Hash{}) {
			t.Fatalf("expected zero MixDigest")
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("engine never emitted a sealed block")
	}
}

// TestRemoteSubmitRejectsBadNonce confirms the remote sealer rejects a
// nonce that fails the keccak256 PoW check.
func TestRemoteSubmitRejectsBadNonce(t *testing.T) {
	transition := uint64(0)
	engine := New(ethash.NewFaker(), &transition)
	engine.SetThreads(-1)
	defer engine.Close()

	api := &API{keccak: engine}

	header := &types.Header{
		Number:     big.NewInt(1),
		Difficulty: big.NewInt(1 << 30), // very strict
		GasLimit:   5_000_000,
	}
	block := types.NewBlockWithHeader(header)

	results := make(chan *types.Block, 1)
	stop := make(chan struct{})
	defer close(stop)
	if err := engine.Seal(nil, block, results, stop); err != nil {
		t.Fatalf("Seal: %v", err)
	}

	var work [3]string
	deadline := time.Now().Add(2 * time.Second)
	for {
		var err error
		work, err = api.GetWork()
		if err == nil {
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("GetWork never returned: %v", err)
		}
		time.Sleep(5 * time.Millisecond)
	}
	sealHash := common.HexToHash(work[0])

	// Nonce 1 is overwhelmingly unlikely to satisfy diff=2^30.
	if api.SubmitWork(types.EncodeNonce(1), sealHash) {
		t.Fatalf("SubmitWork accepted a non-winning nonce at diff=2^30")
	}
}
