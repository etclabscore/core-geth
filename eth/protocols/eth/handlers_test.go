// Copyright 2026 The go-ethereum Authors
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
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/params/vars"
)

// recordingBackend wraps testBackend so that packets surviving the protocol
// level validation can be observed by the tests, and so that transaction
// acceptance can be toggled without pulling in the full eth handler.
type recordingBackend struct {
	*testBackend
	acceptTxs bool
	handled   chan Packet
}

func (b *recordingBackend) AcceptTxs() bool { return b.acceptTxs }

func (b *recordingBackend) Handle(peer *Peer, packet Packet) error {
	b.handled <- packet
	return nil
}

// makeTestTransactions creates n minimal distinct transactions. They are not
// signed: the protocol handlers only ever decode them, they never recover the
// sender.
func makeTestTransactions(n int) []*types.Transaction {
	txs := make([]*types.Transaction, n)
	for i := range txs {
		txs[i] = types.NewTransaction(uint64(i), common.Address{0x01}, big.NewInt(0), vars.TxGas, big.NewInt(1), nil)
	}
	return txs
}

// probePeer round trips a header request to verify that the peer's message
// loop is still alive, and implicitly that all previously sent messages have
// been fully processed.
func probePeer(t *testing.T, peer *testPeer, chain *core.BlockChain, reqID uint64) {
	t.Helper()

	genesis := chain.GetBlockByNumber(0).Header()
	p2p.Send(peer.app, GetBlockHeadersMsg, &GetBlockHeadersPacket{
		RequestId:              reqID,
		GetBlockHeadersRequest: &GetBlockHeadersRequest{Origin: HashOrNumber{Number: 0}, Amount: 1},
	})
	if err := p2p.ExpectMsg(peer.app, BlockHeadersMsg, &BlockHeadersPacket{
		RequestId: reqID,
		List:      encodeRL([]*types.Header{genesis}),
	}); err != nil {
		t.Fatalf("peer no longer serving requests: %v", err)
	}
}

// Tests that broadcast transactions are delivered to the backend only when the
// node is accepting them.
func TestHandleTransactions(t *testing.T) {
	t.Parallel()

	txs := makeTestTransactions(2)

	t.Run("delivered when accepting", func(t *testing.T) {
		backend := &recordingBackend{testBackend: newTestBackend(1), acceptTxs: true, handled: make(chan Packet, 1)}
		defer backend.close()

		peer, _ := newTestPeer("peer", ETH68, backend)
		defer peer.close()

		if err := p2p.Send(peer.app, TransactionsMsg, &TransactionsPacket{RawList: encodeRL(txs)}); err != nil {
			t.Fatalf("failed to send transactions: %v", err)
		}
		select {
		case packet := <-backend.handled:
			res, ok := packet.(*TransactionsPacket)
			if !ok {
				t.Fatalf("unexpected packet type delivered: %T", packet)
			}
			items, err := res.Items()
			if err != nil {
				t.Fatalf("failed to decode delivered transactions: %v", err)
			}
			if len(items) != len(txs) {
				t.Fatalf("transaction count mismatch: have %d, want %d", len(items), len(txs))
			}
			for i, tx := range items {
				if tx.Hash() != txs[i].Hash() {
					t.Fatalf("transaction %d: hash mismatch: have %x, want %x", i, tx.Hash(), txs[i].Hash())
				}
			}
		case <-time.After(2 * time.Second):
			t.Fatalf("transactions not delivered to backend")
		}
	})

	t.Run("dropped when not accepting", func(t *testing.T) {
		backend := &recordingBackend{testBackend: newTestBackend(1), acceptTxs: false, handled: make(chan Packet, 1)}
		defer backend.close()

		peer, errc := newTestPeer("peer", ETH68, backend)
		defer peer.close()

		if err := p2p.Send(peer.app, TransactionsMsg, &TransactionsPacket{RawList: encodeRL(txs)}); err != nil {
			t.Fatalf("failed to send transactions: %v", err)
		}
		// The message must be dropped while keeping the peer: a request sent
		// afterwards on the same peer must still be answered.
		probePeer(t, peer, backend.chain, 1)

		select {
		case packet := <-backend.handled:
			t.Fatalf("transactions delivered while not accepting: %v", packet)
		case err := <-errc:
			t.Fatalf("peer dropped for broadcast while not accepting: %v", err)
		default:
		}
	})
}

// Tests that the transaction broadcast count limit is the one upstream defines,
// and that it is applied where upstream applies it. The tree previously carried
// a stricter pre-decode ceiling of maxHeadersServe*4, so counts between that and
// maxTransactionAnnouncements are the ones that distinguish the two.
func TestHandleTransactionsCountLimit(t *testing.T) {
	t.Parallel()

	all := makeTestTransactions(maxTransactionAnnouncements + 1)

	accepted := func(t *testing.T, count int) {
		backend := &recordingBackend{testBackend: newTestBackend(1), acceptTxs: true, handled: make(chan Packet, 1)}
		defer backend.close()

		peer, errc := newTestPeer("peer", ETH68, backend)
		defer peer.close()

		if err := p2p.Send(peer.app, TransactionsMsg, &TransactionsPacket{RawList: encodeRL(all[:count])}); err != nil {
			t.Fatalf("failed to send transactions: %v", err)
		}
		select {
		case packet := <-backend.handled:
			res, ok := packet.(*TransactionsPacket)
			if !ok {
				t.Fatalf("unexpected packet type delivered: %T", packet)
			}
			if res.Len() != count {
				t.Fatalf("transaction count mismatch: have %d, want %d", res.Len(), count)
			}
		case err := <-errc:
			t.Fatalf("peer dropped for %d transactions: %v", count, err)
		case <-time.After(4 * time.Second):
			t.Fatalf("%d transactions not delivered to backend", count)
		}
	}

	t.Run("above the removed pre-decode ceiling", func(t *testing.T) {
		accepted(t, maxHeadersServe*4+1)
	})

	t.Run("at the upstream limit", func(t *testing.T) {
		accepted(t, maxTransactionAnnouncements)
	})

	t.Run("above the upstream limit drops peer", func(t *testing.T) {
		backend := &recordingBackend{testBackend: newTestBackend(1), acceptTxs: true, handled: make(chan Packet, 1)}
		defer backend.close()

		peer, errc := newTestPeer("peer", ETH68, backend)
		defer peer.close()

		if err := p2p.Send(peer.app, TransactionsMsg, &TransactionsPacket{RawList: encodeRL(all)}); err != nil {
			t.Fatalf("failed to send transactions: %v", err)
		}
		select {
		case err := <-errc:
			if !strings.Contains(err.Error(), "too many transactions") {
				t.Fatalf("error mismatch: have %v, want it to mention too many transactions", err)
			}
		case packet := <-backend.handled:
			t.Fatalf("oversized broadcast delivered: %v", packet)
		case <-time.After(4 * time.Second):
			t.Fatalf("peer not dropped for %d transactions", len(all))
		}
	})

	// With the count checked inside the handler rather than before dispatch, a
	// node that is not accepting transactions returns before the handler allocates
	// another payload buffer, counts its RLP items, or decodes its transactions.
	// An oversized broadcast is therefore not a drop reason while syncing.
	t.Run("above the upstream limit ignored while not accepting", func(t *testing.T) {
		backend := &recordingBackend{testBackend: newTestBackend(1), acceptTxs: false, handled: make(chan Packet, 1)}
		defer backend.close()

		peer, errc := newTestPeer("peer", ETH68, backend)
		defer peer.close()

		if err := p2p.Send(peer.app, TransactionsMsg, &TransactionsPacket{RawList: encodeRL(all)}); err != nil {
			t.Fatalf("failed to send transactions: %v", err)
		}
		probePeer(t, peer, backend.chain, 1)

		select {
		case packet := <-backend.handled:
			t.Fatalf("transactions delivered while not accepting: %v", packet)
		case err := <-errc:
			t.Fatalf("peer dropped for oversized broadcast while not accepting: %v", err)
		default:
		}
	})
}
