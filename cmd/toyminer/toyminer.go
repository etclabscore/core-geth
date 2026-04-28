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

// toyminer is an external miner for the ECIP-1049 keccak256 proof-of-work.
//
// It speaks JSON-RPC against a running geth node (any chain post-ECIP-1049:
// Greenpoint, Mordor after its switchover, mainnet ETC after its switchover)
// and produces sealed blocks via the engine's "eccmine" namespace:
//
//	eccmine_getWork           -> [sealHash, target, blockNumber]
//	eccmine_submitWork(n, h)  -> bool
//	eccmine_submitHashrate(r, id) -> bool
//
// PoW: keccak256(sealHash || be8(nonce)) <= target  (target = 2^256 / difficulty)
//
// Usage:
//
//	geth --greenpoint --mine --miner.threads=0 --http --http.api=eth,net,web3,eccmine
//	go run ./cmd/toyminer --rpc http://127.0.0.1:8545 --threads 4
package main

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"flag"
	"fmt"
	"math/big"
	mrand "math/rand"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"golang.org/x/crypto/sha3"
)

func main() {
	var (
		endpoint  = flag.String("rpc", "http://127.0.0.1:8545", "JSON-RPC endpoint of the keccak256 node")
		threads   = flag.Int("threads", 1, "number of mining threads")
		pollEvery = flag.Duration("poll", 1*time.Second, "interval between getWork polls when no work is available")
		reportFor = flag.Duration("report", 10*time.Second, "interval between hashrate reports / progress logs")
		minerID   = flag.String("id", "", "stable identifier for hashrate reporting (default: random)")
	)
	flag.Parse()

	if *threads <= 0 {
		fatalf("--threads must be > 0")
	}

	id := minerIDOrRandom(*minerID)
	fmt.Printf("toyminer: connecting to %s (id=%s, threads=%d)\n", *endpoint, id.Hex(), *threads)

	client, err := rpc.Dial(*endpoint)
	if err != nil {
		fatalf("rpc dial: %v", err)
	}
	defer client.Close()

	ctx, cancel := signalContext()
	defer cancel()

	m := &miner{
		client:    client,
		threads:   *threads,
		minerID:   id,
		pollEvery: *pollEvery,
		reportFor: *reportFor,
	}
	if err := m.run(ctx); err != nil && !errors.Is(err, context.Canceled) {
		fatalf("%v", err)
	}
}

// miner orchestrates the polling/mining/submitting loop.
type miner struct {
	client    *rpc.Client
	threads   int
	minerID   common.Hash
	pollEvery time.Duration
	reportFor time.Duration

	hashes uint64 // total hashes performed (atomic)
	found  uint64 // total accepted shares (atomic)
}

func (m *miner) run(ctx context.Context) error {
	go m.reportLoop(ctx)

	var (
		curHash common.Hash
		stopCh  chan struct{}
		wg      sync.WaitGroup
	)
	stopWorkers := func() {
		if stopCh != nil {
			close(stopCh)
			wg.Wait()
			stopCh = nil
		}
	}
	defer stopWorkers()

	ticker := time.NewTicker(m.pollEvery)
	defer ticker.Stop()

	for {
		work, err := m.getWork(ctx)
		if err != nil {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-ticker.C:
				continue
			}
		}
		hash := common.HexToHash(work[0])
		if hash != curHash {
			stopWorkers()
			target, ok := decodeTarget(work[1])
			if !ok {
				return fmt.Errorf("malformed target %q", work[1])
			}
			curHash = hash
			fmt.Printf("toyminer: new work hash=%s %s number=%s\n",
				hash.Hex(), targetHumanString(target), work[2])
			stopCh = make(chan struct{})
			for i := 0; i < m.threads; i++ {
				wg.Add(1)
				go m.workerLoop(ctx, &wg, stopCh, hash, target)
			}
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}
	}
}

// workerLoop searches for a winning nonce on the given (sealHash, target)
// pair, submitting any solutions back to the node.
func (m *miner) workerLoop(ctx context.Context, wg *sync.WaitGroup, stop <-chan struct{}, sealHash common.Hash, target *big.Int) {
	defer wg.Done()

	src := mrand.New(mrand.NewSource(time.Now().UnixNano() ^ int64(seedNonce(sealHash))))
	nonce := uint64(src.Int63())

	hasher := sha3.NewLegacyKeccak256()
	var buf [8]byte
	out := make([]byte, 0, 32)

	for {
		select {
		case <-ctx.Done():
			return
		case <-stop:
			return
		default:
		}
		const batch = 1024
		for i := 0; i < batch; i++ {
			hasher.Reset()
			hasher.Write(sealHash[:])
			binary.BigEndian.PutUint64(buf[:], nonce)
			hasher.Write(buf[:])
			out = hasher.Sum(out[:0])

			if leqBytes(out, target) {
				m.submit(ctx, sealHash, nonce)
			}
			nonce++
		}
		atomic.AddUint64(&m.hashes, batch)
	}
}

// submit forwards a winning nonce to the node and reports the outcome.
func (m *miner) submit(ctx context.Context, sealHash common.Hash, nonce uint64) {
	enc := types.EncodeNonce(nonce)
	var ok bool
	if err := m.client.CallContext(ctx, &ok, "eccmine_submitWork", enc, sealHash); err != nil {
		fmt.Printf("toyminer: submitWork rpc error: %v\n", err)
		return
	}
	if ok {
		atomic.AddUint64(&m.found, 1)
		fmt.Printf("toyminer: ACCEPTED nonce=0x%016x sealhash=%s\n", nonce, sealHash.Hex())
	} else {
		fmt.Printf("toyminer: rejected nonce=0x%016x sealhash=%s (stale or invalid)\n", nonce, sealHash.Hex())
	}
}

// getWork polls the node for the next work package.
func (m *miner) getWork(ctx context.Context) ([3]string, error) {
	var work [3]string
	if err := m.client.CallContext(ctx, &work, "eccmine_getWork"); err != nil {
		return work, err
	}
	if work[0] == "" {
		return work, errors.New("empty work")
	}
	return work, nil
}

// reportLoop periodically logs aggregate hashrate and submits the same value
// back to the node so geth's miner UI / metrics see the external miner.
func (m *miner) reportLoop(ctx context.Context) {
	ticker := time.NewTicker(m.reportFor)
	defer ticker.Stop()

	var lastHashes uint64
	lastT := time.Now()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
		now := time.Now()
		total := atomic.LoadUint64(&m.hashes)
		dh := total - lastHashes
		dt := now.Sub(lastT).Seconds()
		lastHashes = total
		lastT = now

		var rate uint64
		if dt > 0 {
			rate = uint64(float64(dh) / dt)
		}
		fmt.Printf("toyminer: hashrate=%d H/s   total_hashes=%d   accepted=%d\n",
			rate, total, atomic.LoadUint64(&m.found))

		var ok bool
		_ = m.client.CallContext(ctx, &ok, "eccmine_submitHashrate",
			hexutil.Uint64(rate), m.minerID)
	}
}

// --- helpers ----------------------------------------------------------------

func decodeTarget(s string) (*big.Int, bool) {
	b, err := hexutil.Decode(s)
	if err != nil {
		return nil, false
	}
	return new(big.Int).SetBytes(b), true
}

func leqBytes(out []byte, target *big.Int) bool {
	v := new(big.Int).SetBytes(out)
	return v.Cmp(target) <= 0
}

func targetHumanString(target *big.Int) string {
	if target.Sign() == 0 {
		return "diff=0"
	}
	twoToThe256 := new(big.Int).Lsh(big.NewInt(1), 256)
	d := new(big.Int).Div(twoToThe256, target)
	return fmt.Sprintf("diff=%s", d.String())
}

func seedNonce(h common.Hash) uint64 {
	return binary.BigEndian.Uint64(h[:8])
}

func minerIDOrRandom(s string) common.Hash {
	if s != "" {
		return common.HexToHash(s)
	}
	var h common.Hash
	_, _ = rand.Read(h[:])
	return h
}

func signalContext() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		fmt.Println("\ntoyminer: shutdown requested")
		cancel()
	}()
	return ctx, cancel
}

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "toyminer: "+format+"\n", args...)
	os.Exit(1)
}
