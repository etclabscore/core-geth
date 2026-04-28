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
	"errors"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
)

// staleThreshold is the maximum age (in blocks) at which the remote sealer
// will still accept a submitted PoW solution.
const remoteStaleThreshold = 7

// remoteSealerTimeout bounds blocking RPC operations against the sealer loop.
const remoteSealerTimeoutKeccak = 1 * time.Second

// errNoMiningWorkKeccak is returned when GetWork is invoked without any
// pending block having been queued by the engine.
var errNoMiningWorkKeccak = errors.New("no mining work available yet")

// errInvalidSealResultKeccak is returned when a SubmitWork solution either
// references an unknown block or fails proof-of-work verification.
var errInvalidSealResultKeccak = errors.New("invalid or stale proof-of-work solution")

// keccakSealTask wraps a sealing request with the channel that wants the
// solution.
type keccakSealTask struct {
	block   *types.Block
	results chan<- *types.Block
}

// keccakMineResult wraps a submitted PoW solution.
type keccakMineResult struct {
	nonce types.BlockNonce
	hash  common.Hash // SealHash of the worked-on block

	errc chan error
}

// keccakSealWork wraps a request for fresh mining work.
type keccakSealWork struct {
	errc chan error
	res  chan [3]string
}

// keccakHashrate carries a remote miner's reported hashrate.
type keccakHashrate struct {
	id   common.Hash
	ping time.Time
	rate uint64
	done chan struct{}
}

// remoteSealer holds the queue of blocks awaiting external proof-of-work and
// dispatches between the engine's Seal call sites and JSON-RPC clients.
//
// It deliberately does NOT push HTTP notifications to miner endpoints — the
// engine's external miner is expected to poll GetWork. This keeps the
// implementation small; users who want push semantics can layer them on top.
type remoteSealer struct {
	engine *Keccak256

	works        map[common.Hash]*types.Block
	rates        map[common.Hash]keccakHashrate
	currentBlock *types.Block
	currentWork  [3]string

	results chan<- *types.Block

	workCh       chan *keccakSealTask
	fetchWorkCh  chan *keccakSealWork
	submitWorkCh chan *keccakMineResult
	fetchRateCh  chan chan uint64
	submitRateCh chan *keccakHashrate
	requestExit  chan struct{}
	exitCh       chan struct{}
}

func startRemoteSealerKeccak(engine *Keccak256) *remoteSealer {
	s := &remoteSealer{
		engine:       engine,
		works:        make(map[common.Hash]*types.Block),
		rates:        make(map[common.Hash]keccakHashrate),
		workCh:       make(chan *keccakSealTask),
		fetchWorkCh:  make(chan *keccakSealWork),
		submitWorkCh: make(chan *keccakMineResult),
		fetchRateCh:  make(chan chan uint64),
		submitRateCh: make(chan *keccakHashrate),
		requestExit:  make(chan struct{}),
		exitCh:       make(chan struct{}),
	}
	go s.loop()
	return s
}

func (s *remoteSealer) loop() {
	defer func() {
		close(s.exitCh)
	}()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case work := <-s.workCh:
			s.results = work.results
			s.makeWork(work.block)

		case work := <-s.fetchWorkCh:
			if s.currentBlock == nil {
				work.errc <- errNoMiningWorkKeccak
			} else {
				work.res <- s.currentWork
			}

		case result := <-s.submitWorkCh:
			if s.submitWork(result.nonce, result.hash) {
				result.errc <- nil
			} else {
				result.errc <- errInvalidSealResultKeccak
			}

		case rate := <-s.submitRateCh:
			s.rates[rate.id] = keccakHashrate{rate: rate.rate, ping: time.Now()}
			close(rate.done)

		case req := <-s.fetchRateCh:
			var total uint64
			for _, r := range s.rates {
				total += r.rate
			}
			req <- total

		case <-ticker.C:
			for id, r := range s.rates {
				if time.Since(r.ping) > 10*time.Second {
					delete(s.rates, id)
				}
			}
			if s.currentBlock != nil {
				for hash, block := range s.works {
					if block.NumberU64()+remoteStaleThreshold <= s.currentBlock.NumberU64() {
						delete(s.works, hash)
					}
				}
			}

		case <-s.requestExit:
			return
		}
	}
}

// makeWork formats the work tuple for the supplied block. The tuple is:
//
//	work[0] = SealHash(header) (32-byte hex)
//	work[1] = target = 2^256 / difficulty (32-byte hex, big-endian)
//	work[2] = block number (hex-encoded)
func (s *remoteSealer) makeWork(block *types.Block) {
	hash := s.engine.inner.SealHash(block.Header())
	s.currentWork[0] = hash.Hex()
	s.currentWork[1] = common.BytesToHash(new(big.Int).Div(two256, block.Difficulty()).Bytes()).Hex()
	s.currentWork[2] = hexutil.EncodeBig(block.Number())
	s.currentBlock = block
	s.works[hash] = block
}

// submitWork validates a SubmitWork payload against the queued blocks and,
// on success, forwards the sealed block to the engine results channel.
func (s *remoteSealer) submitWork(nonce types.BlockNonce, sealhash common.Hash) bool {
	if s.currentBlock == nil {
		return false
	}
	block := s.works[sealhash]
	if block == nil {
		return false
	}
	header := types.CopyHeader(block.Header())
	header.Nonce = nonce
	header.MixDigest = common.Hash{}

	if err := s.engine.verifySeal(header); err != nil {
		s.engine.logger.Warn("Invalid keccak256 PoW submitted", "sealhash", sealhash, "err", err)
		return false
	}
	if s.results == nil {
		s.engine.logger.Warn("Keccak256 result channel is empty, submitted mining result is rejected")
		return false
	}
	solution := block.WithSeal(header)
	if solution.NumberU64()+remoteStaleThreshold > s.currentBlock.NumberU64() {
		select {
		case s.results <- solution:
			s.engine.logger.Debug("Work submitted is acceptable", "number", solution.NumberU64(), "sealhash", sealhash, "hash", solution.Hash())
			return true
		default:
			s.engine.logger.Warn("Sealing result is not read by miner", "mode", "remote", "sealhash", sealhash)
			return false
		}
	}
	s.engine.logger.Warn("Work submitted is too old", "number", solution.NumberU64(), "sealhash", sealhash, "hash", solution.Hash())
	return false
}
