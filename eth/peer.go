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
	"math/big"

	"github.com/ethereum/go-ethereum/eth/protocols/eth"
	"github.com/ethereum/go-ethereum/eth/protocols/snap"
)

/* PTAL(meowsbits)
https://github.com/ethereum/go-ethereum/pull/26804
Difficulty, Head, and ForkID are removed, citing (from the PR):

	> Post-merge there is no more block broadcasts and announcements.
	> As such, we cannot maintain the head infos for our peers.
	> This PR unexposes those infos on the admin.peers API to avoid confusion thinking all our peers are unsynced.
	> Fixes #26733

*/

// ethPeerInfo represents a short summary of the `eth` sub-protocol metadata known
// about a connected peer.
type ethPeerInfo struct {
	Version    uint              `json:"version"`          // Ethereum protocol version negotiated
	Difficulty *big.Int          `json:"difficulty"`       // Total difficulty of the peer's blockchain
	Head       string            `json:"head"`             // Hex hash of the peer's best owned block
	ForkID     ethPeerInfoForkID `json:"forkId,omitempty"` // ForkID from handshake. The JSON tag casing follows the pattern established by chainId elsewhere in APIs.
}

type ethPeerInfoForkID struct {
	Next uint64 `json:"next"`
	Hash string `json:"hash"`
}

// ethPeer is a wrapper around eth.Peer to maintain a few extra metadata.
type ethPeer struct {
	*eth.Peer
	snapExt *snapPeer // Satellite `snap` connection
}

// info gathers and returns some `eth` protocol metadata known about a peer.
func (p *ethPeer) info() *ethPeerInfo {
	hash, td, _ := p.Head()

	info := &ethPeerInfo{
		Version:    p.Version(),
		Difficulty: td,
		Head:       hash.Hex(),
	}
	// ForkID was introduced with eth/64
	if p.Version() >= 64 {
		info.ForkID = ethPeerInfoForkID{
			Next: p.ForkID().Next,
			Hash: fmt.Sprintf("0x%x", p.ForkID().Hash),
		}
	}
	return info
}

// snapPeerInfo represents a short summary of the `snap` sub-protocol metadata known
// about a connected peer.
type snapPeerInfo struct {
	Version uint `json:"version"` // Snapshot protocol version negotiated
}

// snapPeer is a wrapper around snap.Peer to maintain a few extra metadata.
type snapPeer struct {
	*snap.Peer
}

// info gathers and returns some `snap` protocol metadata known about a peer.
func (p *snapPeer) info() *snapPeerInfo {
	return &snapPeerInfo{
		Version: p.Version(),
	}
}
