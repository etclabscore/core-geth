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

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/ethash"
	"github.com/ethereum/go-ethereum/core/types"
	"golang.org/x/crypto/sha3"
)

// newTestEngine returns a Keccak256 engine wrapping an ethash faker, with the
// supplied transition block (a value of 0 means "active from genesis").
func newTestEngine(transition uint64) *Keccak256 {
	t := transition
	return New(ethash.NewFaker(), &t)
}

// TestComputePoWMatchesPlainKeccak verifies that the engine's computePoW
// produces the literal keccak256( sealHash || be8(nonce) ) digest a third
// party would expect.
func TestComputePoWMatchesPlainKeccak(t *testing.T) {
	sealHash := common.HexToHash("0x1122334455667788990011223344556677889900112233445566778899001122")
	nonce := uint64(0xdeadbeefcafef00d)

	got := computePoW(sealHash, nonce)

	hasher := sha3.NewLegacyKeccak256()
	hasher.Write(sealHash[:])
	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], nonce)
	hasher.Write(buf[:])
	want := hasher.Sum(nil)

	if !bytesEq(got, want) {
		t.Fatalf("computePoW mismatch:\n  got  %x\n  want %x", got, want)
	}
}

func bytesEq(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// TestIsPostTransition exercises the activation predicate.
func TestIsPostTransition(t *testing.T) {
	cases := []struct {
		transition *uint64
		number     int64
		want       bool
	}{
		{nil, 0, false},
		{nil, 1_000_000, false},
		{u64ptr(0), 0, true},
		{u64ptr(0), 1, true},
		{u64ptr(100), 99, false},
		{u64ptr(100), 100, true},
		{u64ptr(100), 101, true},
	}
	for i, c := range cases {
		k := New(ethash.NewFaker(), c.transition)
		got := k.isPostTransition(big.NewInt(c.number))
		if got != c.want {
			t.Errorf("case %d: isPostTransition(%v, %d) = %v, want %v", i, c.transition, c.number, got, c.want)
		}
	}
}

func u64ptr(v uint64) *uint64 { return &v }

// TestVerifySealRejectsNonZeroMixDigest ensures the MixDigest invariant on
// post-transition blocks is enforced.
func TestVerifySealRejectsNonZeroMixDigest(t *testing.T) {
	k := newTestEngine(0)

	header := &types.Header{
		Number:     big.NewInt(1),
		Difficulty: big.NewInt(1),
		MixDigest:  common.HexToHash("0x01"),
	}
	if err := k.verifySeal(header); err != errInvalidMixDigest {
		t.Fatalf("expected errInvalidMixDigest, got %v", err)
	}
}

// TestVerifySealRejectsZeroDifficulty checks the difficulty sanity guard.
func TestVerifySealRejectsZeroDifficulty(t *testing.T) {
	k := newTestEngine(0)

	header := &types.Header{
		Number:     big.NewInt(1),
		Difficulty: big.NewInt(0),
	}
	if err := k.verifySeal(header); err != errInvalidDifficulty {
		t.Fatalf("expected errInvalidDifficulty, got %v", err)
	}
}

// TestMineThenVerify mines a block at minimum difficulty using the engine's
// own mining loop and verifies that the resulting header passes verifySeal.
func TestMineThenVerify(t *testing.T) {
	k := newTestEngine(0)
	k.SetThreads(1)

	header := &types.Header{
		ParentHash: common.Hash{},
		Number:     big.NewInt(1),
		Difficulty: big.NewInt(1), // trivial target -> any nonce satisfies it
		GasLimit:   5_000_000,
	}

	// Drive the mine() helper directly so we avoid needing a full chain.
	abort := make(chan struct{})
	defer close(abort)
	found := make(chan *types.Block, 1)

	go k.mine(types.NewBlockWithHeader(header), 0, 0, abort, found)

	select {
	case sealed := <-found:
		got := sealed.Header()
		if got.MixDigest != (common.Hash{}) {
			t.Fatalf("expected zero MixDigest, got %x", got.MixDigest)
		}
		if err := k.verifySeal(got); err != nil {
			t.Fatalf("mined header failed verifySeal: %v", err)
		}
	case <-make(chan struct{}):
		// Unreachable; included so the select compiles even if found is empty.
	}
}

// TestVerifySealRejectsTamperedNonce mines a valid solution then bumps the
// nonce by one — the resulting hash must (with overwhelming probability) miss
// the target at sufficiently high difficulty.
func TestVerifySealRejectsTamperedNonce(t *testing.T) {
	k := newTestEngine(0)

	header := &types.Header{
		ParentHash: common.Hash{},
		Number:     big.NewInt(1),
		Difficulty: big.NewInt(1 << 20),
		GasLimit:   5_000_000,
	}

	// Mine until valid.
	abort := make(chan struct{})
	defer close(abort)
	found := make(chan *types.Block, 1)
	go k.mine(types.NewBlockWithHeader(header), 0, 0, abort, found)
	sealed := <-found

	good := sealed.Header()
	if err := k.verifySeal(good); err != nil {
		t.Fatalf("freshly mined header rejected: %v", err)
	}

	bad := types.CopyHeader(good)
	bad.Nonce = types.EncodeNonce(good.Nonce.Uint64() + 1)
	if err := k.verifySeal(bad); err == nil {
		// Extremely improbable but possible at low difficulty; tolerate by
		// re-tampering with another nonce until we see a rejection. With diff
		// 2^20 the probability of any single nonce satisfying is ~2^-20, so
		// after a handful of mutations we are certain to find a bad one.
		for i := uint64(2); i < 256; i++ {
			bad.Nonce = types.EncodeNonce(good.Nonce.Uint64() + i)
			if err := k.verifySeal(bad); err == errInvalidPoW {
				return
			}
		}
		t.Fatalf("tampered header unexpectedly always validated")
	} else if err != errInvalidPoW {
		t.Fatalf("expected errInvalidPoW for tampered nonce, got %v", err)
	}
}
