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
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
)

// API exposes the keccak256 ECIP-1049 remote-mining JSON-RPC under the
// "eccmine" namespace. Wallet integrations and miners that already speak
// the eth_getWork/eth_submitWork protocol should consult --eccmine-* methods
// post-transition since eth_getWork (served by the inner Ethash engine) has
// no useful semantics for keccak256 blocks.
type API struct {
	keccak *Keccak256
}

// GetWork returns a fresh work package for an external miner.
//
// The returned tuple is:
//
//	result[0] - 32-byte hex-encoded SealHash of the pending block header
//	result[1] - 32-byte hex-encoded boundary ("target") = 2^256 / difficulty
//	result[2] - hex-encoded block number
//
// Unlike Ethash's GetWork, no DAG seed hash is included because Keccak-256
// proof-of-work has no epoch / dataset.
func (api *API) GetWork() ([3]string, error) {
	if api.keccak.remote == nil {
		return [3]string{}, errNoMiningWorkKeccak
	}
	var (
		workCh = make(chan [3]string, 1)
		errc   = make(chan error, 1)
	)
	select {
	case api.keccak.remote.fetchWorkCh <- &keccakSealWork{errc: errc, res: workCh}:
	case <-api.keccak.remote.exitCh:
		return [3]string{}, errKeccakStopped
	}
	select {
	case work := <-workCh:
		return work, nil
	case err := <-errc:
		return [3]string{}, err
	}
}

// SubmitWork submits a proof-of-work solution. The submitted nonce must
// satisfy keccak256(SealHash || nonce) <= 2^256 / difficulty for the block
// previously returned by GetWork with the matching SealHash.
//
// Returns true if the solution was accepted and the block successfully
// sealed; false if it was stale, malformed, or the PoW was invalid.
func (api *API) SubmitWork(nonce types.BlockNonce, hash common.Hash) bool {
	if api.keccak.remote == nil {
		return false
	}
	errc := make(chan error, 1)
	select {
	case api.keccak.remote.submitWorkCh <- &keccakMineResult{nonce: nonce, hash: hash, errc: errc}:
	case <-api.keccak.remote.exitCh:
		return false
	}
	err := <-errc
	return err == nil
}

// SubmitHashrate submits a self-reported hashrate measurement from a remote
// miner. The id is an opaque caller-chosen identifier so that multiple miners
// reporting under one node don't clobber each other.
func (api *API) SubmitHashrate(rate hexutil.Uint64, id common.Hash) bool {
	if api.keccak.remote == nil {
		return false
	}
	done := make(chan struct{}, 1)
	select {
	case api.keccak.remote.submitRateCh <- &keccakHashrate{done: done, rate: uint64(rate), id: id}:
	case <-api.keccak.remote.exitCh:
		return false
	}
	<-done
	return true
}

// GetHashrate returns the aggregated hashrate reported by all remote miners.
func (api *API) GetHashrate() uint64 {
	if api.keccak.remote == nil {
		return 0
	}
	res := make(chan uint64, 1)
	select {
	case api.keccak.remote.fetchRateCh <- res:
		return <-res
	case <-api.keccak.remote.exitCh:
		return 0
	}
}
