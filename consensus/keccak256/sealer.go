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
	crand "crypto/rand"
	"math"
	"math/big"
	"math/rand"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/core/types"
)

// Seal generates a new sealing request for the given input block.
//
// Pre-transition blocks are delegated to the inner Ethash engine so that
// historic chain segments continue to be mined the same way. Post-transition
// blocks are served to remote miners over JSON-RPC (eccmine_getWork /
// eccmine_submitWork) and, if local mining threads have been configured via
// SetThreads, additionally mined locally.
func (k *Keccak256) Seal(chain consensus.ChainHeaderReader, block *types.Block, results chan<- *types.Block, stop <-chan struct{}) error {
	if !k.isPostTransition(block.Number()) {
		return k.inner.Seal(chain, block, results, stop)
	}

	// Always queue the work for remote (external) miners.
	remote := k.ensureRemote()
	select {
	case remote.workCh <- &keccakSealTask{block: block, results: results}:
	case <-remote.exitCh:
		return errKeccakStopped
	}

	abort := make(chan struct{})

	k.lock.Lock()
	threads := k.threads
	if k.rand == nil {
		seed, err := crand.Int(crand.Reader, big.NewInt(math.MaxInt64))
		if err != nil {
			k.lock.Unlock()
			return err
		}
		k.rand = rand.New(rand.NewSource(seed.Int64()))
	}
	k.lock.Unlock()
	if threads < 0 {
		// Negative threads disable local mining (remote-only).
		threads = 0
	}

	var (
		pend   sync.WaitGroup
		locals = make(chan *types.Block)
	)
	for i := 0; i < threads; i++ {
		pend.Add(1)
		go func(id int, nonce uint64) {
			defer pend.Done()
			k.mine(block, id, nonce, abort, locals)
		}(i, uint64(k.rand.Int63()))
	}

	go func() {
		var result *types.Block
		select {
		case <-stop:
			close(abort)
		case result = <-locals:
			select {
			case results <- result:
			default:
				k.logger.Warn("Sealing result is not read by miner", "mode", "local", "sealhash", k.SealHash(block.Header()))
			}
			close(abort)
		case <-k.update:
			close(abort)
			if err := k.Seal(chain, block, results, stop); err != nil {
				k.logger.Error("Failed to restart sealing after update", "err", err)
			}
		}
		pend.Wait()
	}()
	return nil
}

// mine is the actual proof-of-work miner that searches for a nonce that
// satisfies the block's difficulty.
func (k *Keccak256) mine(block *types.Block, id int, seed uint64, abort chan struct{}, found chan *types.Block) {
	var (
		header   = block.Header()
		sealhash = k.inner.SealHash(header)
		target   = new(big.Int).Div(two256, header.Difficulty)
	)
	var (
		attempts = int64(0)
		nonce    = seed
	)
	logger := k.logger.New("miner", id)
	logger.Trace("Started keccak256 search for new nonces", "seed", seed)
search:
	for {
		select {
		case <-abort:
			logger.Trace("Keccak256 nonce search aborted", "attempts", nonce-seed)
			k.hashrate.Mark(attempts)
			break search

		default:
			attempts++
			if (attempts % (1 << 15)) == 0 {
				k.hashrate.Mark(attempts)
				attempts = 0
			}
			result := computePoW(sealhash, nonce)
			if new(big.Int).SetBytes(result).Cmp(target) <= 0 {
				sealed := types.CopyHeader(header)
				sealed.Nonce = types.EncodeNonce(nonce)
				sealed.MixDigest = common.Hash{}
				select {
				case found <- block.WithSeal(sealed):
					logger.Trace("Keccak256 nonce found and reported", "attempts", nonce-seed, "nonce", nonce)
				case <-abort:
					logger.Trace("Keccak256 nonce found but discarded", "attempts", nonce-seed, "nonce", nonce)
				}
				break search
			}
			nonce++
		}
	}
}
