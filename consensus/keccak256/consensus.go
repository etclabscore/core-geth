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
	"errors"
	"fmt"
	"math/big"
	"runtime"
	"time"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/consensus/misc"
	"github.com/ethereum/go-ethereum/consensus/misc/eip1559"
	"github.com/ethereum/go-ethereum/consensus/misc/eip4844"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params/mutations"
	"github.com/ethereum/go-ethereum/params/vars"
	"github.com/ethereum/go-ethereum/rpc"
	"golang.org/x/crypto/sha3"
)

// Protocol constants. Mirror the values used by the Ethash engine so that the
// transition is fully transparent at the block-validity level.
var (
	maxUncles              = 2
	allowedFutureBlockTime = 15 * time.Second

	two256 = new(big.Int).Exp(big.NewInt(2), big.NewInt(256), big.NewInt(0))
)

// Engine errors.
var (
	errOlderBlockTime    = errors.New("timestamp older than parent")
	errTooManyUncles     = errors.New("too many uncles")
	errDuplicateUncle    = errors.New("duplicate uncle")
	errUncleIsAncestor   = errors.New("uncle is ancestor")
	errDanglingUncle     = errors.New("uncle's parent is not ancestor")
	errInvalidDifficulty = errors.New("non-positive difficulty")
	errInvalidMixDigest  = errors.New("invalid mix digest")
	errInvalidPoW        = errors.New("invalid proof-of-work")
)

// Author returns the header's coinbase as the verified author.
func (k *Keccak256) Author(header *types.Header) (common.Address, error) {
	return header.Coinbase, nil
}

// VerifyHeader checks whether a header conforms to the consensus rules of
// either the inner Ethash engine (pre-transition) or the Keccak-256 engine
// (post-transition).
func (k *Keccak256) VerifyHeader(chain consensus.ChainHeaderReader, header *types.Header, seal bool) error {
	if !k.isPostTransition(header.Number) {
		return k.inner.VerifyHeader(chain, header, seal)
	}
	number := header.Number.Uint64()
	if chain.GetHeader(header.Hash(), number) != nil {
		return nil
	}
	parent := chain.GetHeader(header.ParentHash, number-1)
	if parent == nil {
		return consensus.ErrUnknownAncestor
	}
	return k.verifyHeader(chain, header, parent, false, seal, time.Now().Unix())
}

// VerifyHeaders is similar to VerifyHeader but verifies a batch of headers
// concurrently. Each header in the batch is dispatched to the appropriate
// engine based on its block number relative to the transition.
func (k *Keccak256) VerifyHeaders(chain consensus.ChainHeaderReader, headers []*types.Header, seals []bool) (chan<- struct{}, <-chan error) {
	if k.transition == nil {
		return k.inner.VerifyHeaders(chain, headers, seals)
	}
	if len(headers) == 0 {
		abort, results := make(chan struct{}), make(chan error, len(headers))
		return abort, results
	}

	workers := runtime.GOMAXPROCS(0)
	if len(headers) < workers {
		workers = len(headers)
	}

	var (
		inputs  = make(chan int)
		done    = make(chan int, workers)
		errs    = make([]error, len(headers))
		abort   = make(chan struct{})
		unixNow = time.Now().Unix()
	)
	for i := 0; i < workers; i++ {
		go func() {
			for index := range inputs {
				errs[index] = k.verifyHeaderWorker(chain, headers, seals, index, unixNow)
				done <- index
			}
		}()
	}

	errorsOut := make(chan error, len(headers))
	go func() {
		defer close(inputs)
		var (
			in, out = 0, 0
			checked = make([]bool, len(headers))
			inputs  = inputs
		)
		for {
			select {
			case inputs <- in:
				if in++; in == len(headers) {
					inputs = nil
				}
			case index := <-done:
				for checked[index] = true; checked[out]; out++ {
					errorsOut <- errs[out]
					if out == len(headers)-1 {
						return
					}
				}
			case <-abort:
				return
			}
		}
	}()
	return abort, errorsOut
}

func (k *Keccak256) verifyHeaderWorker(chain consensus.ChainHeaderReader, headers []*types.Header, seals []bool, index int, unixNow int64) error {
	var parent *types.Header
	if index == 0 {
		parent = chain.GetHeader(headers[0].ParentHash, headers[0].Number.Uint64()-1)
	} else if headers[index-1].Hash() == headers[index].ParentHash {
		parent = headers[index-1]
	}
	if parent == nil {
		return consensus.ErrUnknownAncestor
	}
	header := headers[index]
	if !k.isPostTransition(header.Number) {
		// Delegate the entire single-header verification to the inner engine.
		// The inner engine performs its own ancestry lookup, but we have
		// already established the parent above; the inner VerifyHeader will
		// repeat that check, which is acceptable.
		if chain.GetHeader(header.Hash(), header.Number.Uint64()) != nil {
			return nil
		}
		return k.inner.VerifyHeader(chain, header, seals[index])
	}
	return k.verifyHeader(chain, header, parent, false, seals[index], unixNow)
}

// VerifyUncles verifies that the given block's uncles conform to the consensus
// rules of the wrapping engine, dispatching each uncle to the appropriate
// seal-check based on the uncle's own block number.
func (k *Keccak256) VerifyUncles(chain consensus.ChainReader, block *types.Block) error {
	if !k.isPostTransition(block.Number()) {
		return k.inner.VerifyUncles(chain, block)
	}
	if len(block.Uncles()) > maxUncles {
		return errTooManyUncles
	}
	if len(block.Uncles()) == 0 {
		return nil
	}
	uncles, ancestors := mapset.NewSet[common.Hash](), make(map[common.Hash]*types.Header)

	number, parent := block.NumberU64()-1, block.ParentHash()
	for i := 0; i < 7; i++ {
		ancestorHeader := chain.GetHeader(parent, number)
		if ancestorHeader == nil {
			break
		}
		ancestors[parent] = ancestorHeader
		if ancestorHeader.UncleHash != types.EmptyUncleHash {
			ancestor := chain.GetBlock(parent, number)
			if ancestor == nil {
				break
			}
			for _, uncle := range ancestor.Uncles() {
				uncles.Add(uncle.Hash())
			}
		}
		parent, number = ancestorHeader.ParentHash, number-1
	}
	ancestors[block.Hash()] = block.Header()
	uncles.Add(block.Hash())

	for _, uncle := range block.Uncles() {
		hash := uncle.Hash()
		if uncles.Contains(hash) {
			return errDuplicateUncle
		}
		uncles.Add(hash)

		if ancestors[hash] != nil {
			return errUncleIsAncestor
		}
		if ancestors[uncle.ParentHash] == nil || uncle.ParentHash == block.ParentHash() {
			return errDanglingUncle
		}
		if err := k.verifyHeader(chain, uncle, ancestors[uncle.ParentHash], true, true, time.Now().Unix()); err != nil {
			return err
		}
	}
	return nil
}

// verifyHeader is the post-transition header validity check. It mirrors the
// Ethash implementation in core-geth but performs the seal verification using
// Keccak-256.
func (k *Keccak256) verifyHeader(chain consensus.ChainHeaderReader, header, parent *types.Header, uncle bool, seal bool, unixNow int64) error {
	// Ensure that the header's extra-data section is of a reasonable size.
	if uint64(len(header.Extra)) > vars.MaximumExtraDataSize {
		return fmt.Errorf("extra-data too long: %d > %d", len(header.Extra), vars.MaximumExtraDataSize)
	}
	if !uncle {
		if header.Time > uint64(unixNow+int64(allowedFutureBlockTime.Seconds())) {
			return consensus.ErrFutureBlock
		}
	}
	if header.Time <= parent.Time {
		return errOlderBlockTime
	}
	// Difficulty is delegated to the inner Ethash engine so ETC's existing
	// difficulty rules (Byzantium / EIP-100B / Defuse / etc.) carry over
	// unchanged.
	expected := k.inner.CalcDifficulty(chain, header.Time, parent)
	if expected.Cmp(header.Difficulty) != 0 {
		return fmt.Errorf("invalid difficulty: have %v, want %v", header.Difficulty, expected)
	}
	if header.GasLimit > vars.MaxGasLimit {
		return fmt.Errorf("invalid gasLimit: have %v, max %v", header.GasLimit, vars.MaxGasLimit)
	}
	if header.GasUsed > header.GasLimit {
		return fmt.Errorf("invalid gasUsed: have %d, gasLimit %d", header.GasUsed, header.GasLimit)
	}
	if diff := new(big.Int).Sub(header.Number, parent.Number); diff.Cmp(big.NewInt(1)) != 0 {
		return consensus.ErrInvalidNumber
	}
	if !chain.Config().IsEnabled(chain.Config().GetEIP1559Transition, header.Number) {
		if header.BaseFee != nil {
			return fmt.Errorf("invalid baseFee before fork: have %d, expected 'nil'", header.BaseFee)
		}
		if err := misc.VerifyGaslimit(parent.GasLimit, header.GasLimit); err != nil {
			return err
		}
	} else if err := eip1559.VerifyEIP1559Header(chain.Config(), parent, header); err != nil {
		return err
	}

	// Shanghai / EIP-4895
	eip4895Enabled := chain.Config().IsEnabledByTime(chain.Config().GetEIP4895TransitionTime, &header.Time) || chain.Config().IsEnabled(chain.Config().GetEIP4895Transition, header.Number)
	if !eip4895Enabled {
		if header.WithdrawalsHash != nil {
			return fmt.Errorf("invalid withdrawalsHash: have %x, expected nil", header.WithdrawalsHash)
		}
	} else if header.WithdrawalsHash == nil {
		return errors.New("header is missing withdrawalsHash")
	}

	// Cancun / EIP-4844
	eip4844Enabled := chain.Config().IsEnabledByTime(chain.Config().GetEIP4844TransitionTime, &header.Time) || chain.Config().IsEnabled(chain.Config().GetEIP4844Transition, header.Number)
	if !eip4844Enabled {
		switch {
		case header.ExcessBlobGas != nil:
			return fmt.Errorf("invalid excessBlobGas: have %d, expected nil", header.ExcessBlobGas)
		case header.BlobGasUsed != nil:
			return fmt.Errorf("invalid blobGasUsed: have %d, expected nil", header.BlobGasUsed)
		}
	} else if err := eip4844.VerifyEIP4844Header(parent, header); err != nil {
		return err
	}

	// EIP-4788
	eip4788Enabled := chain.Config().IsEnabledByTime(chain.Config().GetEIP4788TransitionTime, &header.Time) || chain.Config().IsEnabled(chain.Config().GetEIP4788Transition, header.Number)
	if !eip4788Enabled {
		if header.ParentBeaconRoot != nil {
			return fmt.Errorf("invalid parentBeaconRoot, have %#x, expected nil", header.ParentBeaconRoot)
		}
	} else if header.ParentBeaconRoot == nil {
		return errors.New("header is missing beaconRoot")
	}

	if seal {
		if err := k.verifySeal(header); err != nil {
			return err
		}
	}
	if err := mutations.VerifyDAOHeaderExtraData(chain.Config(), header); err != nil {
		return err
	}
	return nil
}

// verifySeal checks the Keccak-256 PoW solution carried by header.
//
//	pow = keccak256( SealHash(header) || big-endian(nonce) )
//	pow as big-endian integer must satisfy pow <= 2^256 / difficulty
//
// MixDigest is unused for Keccak-256 PoW and must be the zero hash.
func (k *Keccak256) verifySeal(header *types.Header) error {
	if header.Difficulty.Sign() <= 0 {
		return errInvalidDifficulty
	}
	if header.MixDigest != (common.Hash{}) {
		return errInvalidMixDigest
	}
	result := computePoW(k.inner.SealHash(header), header.Nonce.Uint64())
	target := new(big.Int).Div(two256, header.Difficulty)
	if new(big.Int).SetBytes(result).Cmp(target) > 0 {
		return errInvalidPoW
	}
	return nil
}

// computePoW returns keccak256( sealHash || big-endian(nonce) ).
func computePoW(sealHash common.Hash, nonce uint64) []byte {
	hasher := sha3.NewLegacyKeccak256()
	hasher.Write(sealHash[:])
	var nonceBytes [8]byte
	binary.BigEndian.PutUint64(nonceBytes[:], nonce)
	hasher.Write(nonceBytes[:])
	return hasher.Sum(nil)
}

// VerifySeal is the (deprecated) external seal-verification entrypoint kept
// for compatibility with code paths that call it directly.
//
//nolint:revive
func (k *Keccak256) VerifySeal(chain consensus.ChainHeaderReader, header *types.Header) error {
	if !k.isPostTransition(header.Number) {
		return k.innerVerifySeal(chain, header)
	}
	return k.verifySeal(header)
}

// innerVerifySeal is a small shim that runs the inner Ethash engine's full
// header verification (which embeds the seal check) for pre-transition headers.
// We re-use VerifyHeader with seal=true to avoid having to expose the inner
// engine's private verifySeal method.
func (k *Keccak256) innerVerifySeal(chain consensus.ChainHeaderReader, header *types.Header) error {
	// To exercise just the seal check we'd need the inner.verifySeal which is
	// private; calling VerifyHeader is the safest equivalent and is already
	// invoked via VerifyHeader/VerifyHeaders during normal block import.
	if chain == nil {
		return errors.New("keccak256: cannot verify pre-transition seal without chain reader")
	}
	return k.inner.VerifyHeader(chain, header, true)
}

// CalcDifficulty defers to the inner Ethash engine, which encapsulates the
// chain-config-aware difficulty rules used by ETC.
func (k *Keccak256) CalcDifficulty(chain consensus.ChainHeaderReader, time uint64, parent *types.Header) *big.Int {
	return k.inner.CalcDifficulty(chain, time, parent)
}

// Prepare initializes the difficulty field of header. Identical for both
// pre- and post-transition blocks.
func (k *Keccak256) Prepare(chain consensus.ChainHeaderReader, header *types.Header) error {
	return k.inner.Prepare(chain, header)
}

// Finalize accumulates block and uncle rewards.
func (k *Keccak256) Finalize(chain consensus.ChainHeaderReader, header *types.Header, statedb *state.StateDB, txs []*types.Transaction, uncles []*types.Header, withdrawals []*types.Withdrawal) {
	k.inner.Finalize(chain, header, statedb, txs, uncles, withdrawals)
}

// FinalizeAndAssemble assembles the final block using the inner engine's logic.
func (k *Keccak256) FinalizeAndAssemble(chain consensus.ChainHeaderReader, header *types.Header, statedb *state.StateDB, txs []*types.Transaction, uncles []*types.Header, receipts []*types.Receipt, withdrawals []*types.Withdrawal) (*types.Block, error) {
	return k.inner.FinalizeAndAssemble(chain, header, statedb, txs, uncles, receipts, withdrawals)
}

// SealHash returns the hash of a block prior to it being sealed. The encoding
// is intentionally identical to Ethash's so that miner-protocol payloads
// (work[0]) are unchanged across the transition.
func (k *Keccak256) SealHash(header *types.Header) common.Hash {
	return k.inner.SealHash(header)
}

// APIs returns the RPC APIs this consensus engine provides.
//
// The inner Ethash engine's "eth" and "ethash" namespaces are passed through
// (they remain the canonical work source for pre-transition blocks). In
// addition, post-transition mining is exposed under the "eccmine" namespace
// (ECIP-1049 mining): eccmine_getWork, eccmine_submitWork, eccmine_submitHashrate,
// eccmine_getHashrate.
func (k *Keccak256) APIs(chain consensus.ChainHeaderReader) []rpc.API {
	apis := k.inner.APIs(chain)
	apis = append(apis, rpc.API{
		Namespace: "eccmine",
		Service:   &API{keccak: k},
	})
	return apis
}

// ensure interface satisfaction at compile time.
var _ consensus.Engine = (*Keccak256)(nil)
