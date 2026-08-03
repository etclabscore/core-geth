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
	"errors"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/tracker"
	"github.com/ethereum/go-ethereum/params/vars"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/trie"
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

// testChainGenerator fills blocks with a few value transfers and an uncle so
// that the derived transaction, uncle and receipt roots are all non-trivial.
func testChainGenerator(i int, block *core.BlockGen) {
	signer := types.HomesteadSigner{}
	switch i {
	case 0:
		tx, _ := types.SignTx(types.NewTransaction(block.TxNonce(testAddr), common.Address{0x01}, big.NewInt(1000), vars.TxGas, block.BaseFee(), nil), signer, testKey)
		block.AddTx(tx)
	case 1:
		tx1, _ := types.SignTx(types.NewTransaction(block.TxNonce(testAddr), common.Address{0x02}, big.NewInt(1000), vars.TxGas, block.BaseFee(), nil), signer, testKey)
		block.AddTx(tx1)
		tx2, _ := types.SignTx(types.NewTransaction(block.TxNonce(testAddr), common.Address{0x03}, big.NewInt(1000), vars.TxGas, block.BaseFee(), nil), signer, testKey)
		block.AddTx(tx2)
	case 3:
		b2 := block.PrevBlock(1).Header()
		b2.Extra = []byte("foo")
		block.AddUncle(b2)
	}
}

// Tests that inbound block body responses are checked against their pending
// request before any of their content is decoded and dispatched.
func TestHandleBlockBodies(t *testing.T) {
	t.Parallel()

	backend := newTestBackendWithGenerator(4, false, testChainGenerator)
	defer backend.close()

	var (
		hashes  []common.Hash
		headers []*types.Header
		bodies  []BlockBody
	)
	for i := uint64(1); i <= backend.chain.CurrentBlock().Number.Uint64(); i++ {
		block := backend.chain.GetBlockByNumber(i)
		hashes = append(hashes, block.Hash())
		headers = append(headers, block.Header())
		bodies = append(bodies, encodeBody(block))
	}

	t.Run("matching response delivered", func(t *testing.T) {
		peer, _ := newTestPeer("peer", ETH68, backend)
		defer peer.close()

		sink := make(chan *Response, 1)
		sent := make(chan error, 1)
		go func() {
			_, err := peer.RequestBodies(hashes, sink)
			sent <- err
		}()
		// Read the query off the network and answer it with the real bodies.
		msg, err := peer.app.ReadMsg()
		if err != nil {
			t.Fatalf("failed to read bodies query: %v", err)
		}
		if msg.Code != GetBlockBodiesMsg {
			t.Fatalf("query code mismatch: have %d, want %d", msg.Code, GetBlockBodiesMsg)
		}
		var query GetBlockBodiesPacket
		if err := msg.Decode(&query); err != nil {
			t.Fatalf("failed to decode bodies query: %v", err)
		}
		if err := <-sent; err != nil {
			t.Fatalf("failed to request bodies: %v", err)
		}
		if err := p2p.Send(peer.app, BlockBodiesMsg, &BlockBodiesPacket{
			RequestId: query.RequestId,
			List:      encodeRL(bodies),
		}); err != nil {
			t.Fatalf("failed to send bodies response: %v", err)
		}
		select {
		case res := <-sink:
			items := res.Res.(*BlockBodiesResponse)
			if len(*items) != len(hashes) {
				t.Errorf("body count mismatch: have %d, want %d", len(*items), len(hashes))
			}
			meta := res.Meta.(BlockBodyHashes)
			if len(meta.TransactionRoots) != len(headers) {
				t.Fatalf("metadata count mismatch: have %d, want %d", len(meta.TransactionRoots), len(headers))
			}
			for i, header := range headers {
				if meta.TransactionRoots[i] != header.TxHash {
					t.Errorf("body %d: transaction root mismatch: have %x, want %x", i, meta.TransactionRoots[i], header.TxHash)
				}
				if meta.UncleHashes[i] != header.UncleHash {
					t.Errorf("body %d: uncle hash mismatch: have %x, want %x", i, meta.UncleHashes[i], header.UncleHash)
				}
				if meta.WithdrawalRoots[i] != (common.Hash{}) {
					t.Errorf("body %d: unexpected withdrawal root: %x", i, meta.WithdrawalRoots[i])
				}
			}
			res.Done <- nil
		case <-time.After(2 * time.Second):
			t.Fatalf("matching bodies response not delivered")
		}
	})

	t.Run("oversized response drops peer", func(t *testing.T) {
		peer, errc := newTestPeer("peer", ETH68, backend)
		defer peer.close()

		sink := make(chan *Response, 1)
		sent := make(chan error, 1)
		go func() {
			_, err := peer.RequestBodies(hashes[:1], sink)
			sent <- err
		}()
		msg, err := peer.app.ReadMsg()
		if err != nil {
			t.Fatalf("failed to read bodies query: %v", err)
		}
		var query GetBlockBodiesPacket
		if err := msg.Decode(&query); err != nil {
			t.Fatalf("failed to decode bodies query: %v", err)
		}
		if err := <-sent; err != nil {
			t.Fatalf("failed to request bodies: %v", err)
		}
		// Respond with more bodies than were requested.
		if err := p2p.Send(peer.app, BlockBodiesMsg, &BlockBodiesPacket{
			RequestId: query.RequestId,
			List:      encodeRL(bodies[:2]),
		}); err != nil {
			t.Fatalf("failed to send bodies response: %v", err)
		}
		select {
		case err := <-errc:
			if !errors.Is(err, tracker.ErrTooManyItems) {
				t.Fatalf("error mismatch: have %v, want %v", err, tracker.ErrTooManyItems)
			}
		case <-sink:
			t.Fatalf("oversized bodies response delivered")
		case <-time.After(2 * time.Second):
			t.Fatalf("peer not dropped for oversized response")
		}
	})

	t.Run("undecodable response drops peer", func(t *testing.T) {
		peer, errc := newTestPeer("peer", ETH68, backend)
		defer peer.close()

		sink := make(chan *Response, 1)
		sent := make(chan error, 1)
		go func() {
			_, err := peer.RequestBodies(hashes[:1], sink)
			sent <- err
		}()
		msg, err := peer.app.ReadMsg()
		if err != nil {
			t.Fatalf("failed to read bodies query: %v", err)
		}
		var query GetBlockBodiesPacket
		if err := msg.Decode(&query); err != nil {
			t.Fatalf("failed to decode bodies query: %v", err)
		}
		if err := <-sent; err != nil {
			t.Fatalf("failed to request bodies: %v", err)
		}
		// Respond with the requested item count, but an item that is not a
		// decodable block body.
		var junk rlp.RawList[BlockBody]
		if err := junk.AppendRaw([]byte{0x05}); err != nil {
			t.Fatal(err)
		}
		if err := p2p.Send(peer.app, BlockBodiesMsg, &BlockBodiesPacket{
			RequestId: query.RequestId,
			List:      junk,
		}); err != nil {
			t.Fatalf("failed to send bodies response: %v", err)
		}
		select {
		case err := <-errc:
			if err == nil {
				t.Fatalf("expected peer to be dropped with an error")
			}
		case <-sink:
			t.Fatalf("undecodable bodies response delivered")
		case <-time.After(2 * time.Second):
			t.Fatalf("peer not dropped for undecodable response")
		}
	})

	t.Run("unsolicited response drops peer", func(t *testing.T) {
		peer, errc := newTestPeer("peer", ETH68, backend)
		defer peer.close()

		if err := p2p.Send(peer.app, BlockBodiesMsg, &BlockBodiesPacket{
			RequestId: 0x1337beef,
			List:      encodeRL(bodies[:1]),
		}); err != nil {
			t.Fatalf("failed to send bodies response: %v", err)
		}
		select {
		case err := <-errc:
			if !errors.Is(err, tracker.ErrNoMatchingRequest) {
				t.Fatalf("error mismatch: have %v, want %v", err, tracker.ErrNoMatchingRequest)
			}
		case <-time.After(2 * time.Second):
			t.Fatalf("peer not dropped for unsolicited response")
		}
	})
}

// Tests that inbound receipt responses are checked against their pending
// request before any of their content is decoded and dispatched.
func TestHandleReceipts(t *testing.T) {
	t.Parallel()

	backend := newTestBackendWithGenerator(4, false, testChainGenerator)
	defer backend.close()

	var (
		hashes  []common.Hash
		headers []*types.Header
	)
	for i := uint64(1); i <= backend.chain.CurrentBlock().Number.Uint64(); i++ {
		block := backend.chain.GetBlockByNumber(i)
		hashes = append(hashes, block.Hash())
		headers = append(headers, block.Header())
	}
	receiptList := func(hash common.Hash) ReceiptList {
		return encodeRL([]*types.Receipt(backend.chain.GetReceiptsByHash(hash)))
	}

	t.Run("matching response delivered", func(t *testing.T) {
		peer, _ := newTestPeer("peer", ETH68, backend)
		defer peer.close()

		sink := make(chan *Response, 1)
		sent := make(chan error, 1)
		go func() {
			_, err := peer.RequestReceipts(hashes, sink)
			sent <- err
		}()
		msg, err := peer.app.ReadMsg()
		if err != nil {
			t.Fatalf("failed to read receipts query: %v", err)
		}
		if msg.Code != GetReceiptsMsg {
			t.Fatalf("query code mismatch: have %d, want %d", msg.Code, GetReceiptsMsg)
		}
		var query GetReceiptsPacket
		if err := msg.Decode(&query); err != nil {
			t.Fatalf("failed to decode receipts query: %v", err)
		}
		if err := <-sent; err != nil {
			t.Fatalf("failed to request receipts: %v", err)
		}
		var receipts rlp.RawList[ReceiptList]
		for _, hash := range query.GetReceiptsRequest {
			receipts.Append(receiptList(hash))
		}
		if err := p2p.Send(peer.app, ReceiptsMsg, &ReceiptsPacket{
			RequestId: query.RequestId,
			List:      receipts,
		}); err != nil {
			t.Fatalf("failed to send receipts response: %v", err)
		}
		select {
		case res := <-sink:
			items := res.Res.(*ReceiptsRLPResponse)
			if len(*items) != len(hashes) {
				t.Errorf("receipt list count mismatch: have %d, want %d", len(*items), len(hashes))
			}
			meta := res.Meta.([]common.Hash)
			if len(meta) != len(headers) {
				t.Fatalf("metadata count mismatch: have %d, want %d", len(meta), len(headers))
			}
			for i, header := range headers {
				if meta[i] != header.ReceiptHash {
					t.Errorf("list %d: receipt root mismatch: have %x, want %x", i, meta[i], header.ReceiptHash)
				}
			}
			res.Done <- nil
		case <-time.After(2 * time.Second):
			t.Fatalf("matching receipts response not delivered")
		}
	})

	t.Run("oversized response drops peer", func(t *testing.T) {
		peer, errc := newTestPeer("peer", ETH68, backend)
		defer peer.close()

		sink := make(chan *Response, 1)
		sent := make(chan error, 1)
		go func() {
			_, err := peer.RequestReceipts(hashes[:1], sink)
			sent <- err
		}()
		msg, err := peer.app.ReadMsg()
		if err != nil {
			t.Fatalf("failed to read receipts query: %v", err)
		}
		var query GetReceiptsPacket
		if err := msg.Decode(&query); err != nil {
			t.Fatalf("failed to decode receipts query: %v", err)
		}
		if err := <-sent; err != nil {
			t.Fatalf("failed to request receipts: %v", err)
		}
		// Respond with more receipt lists than were requested.
		var receipts rlp.RawList[ReceiptList]
		for _, hash := range hashes[:2] {
			receipts.Append(receiptList(hash))
		}
		if err := p2p.Send(peer.app, ReceiptsMsg, &ReceiptsPacket{
			RequestId: query.RequestId,
			List:      receipts,
		}); err != nil {
			t.Fatalf("failed to send receipts response: %v", err)
		}
		select {
		case err := <-errc:
			if !errors.Is(err, tracker.ErrTooManyItems) {
				t.Fatalf("error mismatch: have %v, want %v", err, tracker.ErrTooManyItems)
			}
		case <-sink:
			t.Fatalf("oversized receipts response delivered")
		case <-time.After(2 * time.Second):
			t.Fatalf("peer not dropped for oversized response")
		}
	})

	t.Run("undecodable response drops peer", func(t *testing.T) {
		peer, errc := newTestPeer("peer", ETH68, backend)
		defer peer.close()

		sink := make(chan *Response, 1)
		sent := make(chan error, 1)
		go func() {
			_, err := peer.RequestReceipts(hashes[:1], sink)
			sent <- err
		}()
		msg, err := peer.app.ReadMsg()
		if err != nil {
			t.Fatalf("failed to read receipts query: %v", err)
		}
		var query GetReceiptsPacket
		if err := msg.Decode(&query); err != nil {
			t.Fatalf("failed to decode receipts query: %v", err)
		}
		if err := <-sent; err != nil {
			t.Fatalf("failed to request receipts: %v", err)
		}
		// Respond with the requested item count, but an item that is not a
		// decodable receipt list.
		var junk rlp.RawList[ReceiptList]
		if err := junk.AppendRaw([]byte{0x05}); err != nil {
			t.Fatal(err)
		}
		if err := p2p.Send(peer.app, ReceiptsMsg, &ReceiptsPacket{
			RequestId: query.RequestId,
			List:      junk,
		}); err != nil {
			t.Fatalf("failed to send receipts response: %v", err)
		}
		select {
		case err := <-errc:
			if err == nil {
				t.Fatalf("expected peer to be dropped with an error")
			}
		case <-sink:
			t.Fatalf("undecodable receipts response delivered")
		case <-time.After(2 * time.Second):
			t.Fatalf("peer not dropped for undecodable response")
		}
	})

	t.Run("unsolicited response drops peer", func(t *testing.T) {
		peer, errc := newTestPeer("peer", ETH68, backend)
		defer peer.close()

		var receipts rlp.RawList[ReceiptList]
		receipts.Append(receiptList(hashes[0]))

		if err := p2p.Send(peer.app, ReceiptsMsg, &ReceiptsPacket{
			RequestId: 0x1337beef,
			List:      receipts,
		}); err != nil {
			t.Fatalf("failed to send receipts response: %v", err)
		}
		select {
		case err := <-errc:
			if !errors.Is(err, tracker.ErrNoMatchingRequest) {
				t.Fatalf("error mismatch: have %v, want %v", err, tracker.ErrNoMatchingRequest)
			}
		case <-time.After(2 * time.Second):
			t.Fatalf("peer not dropped for unsolicited response")
		}
	})
}

// makeTestUncles creates n minimal distinct uncle headers.
func makeTestUncles(n int) []*types.Header {
	uncles := make([]*types.Header, n)
	for i := range uncles {
		uncles[i] = &types.Header{
			Number:     big.NewInt(int64(i)),
			Difficulty: big.NewInt(131072),
			Extra:      []byte{byte(i)},
		}
	}
	return uncles
}

// newPropagatedBlock assembles a block with the given body parts and header
// roots matching them, the way a remote peer would propagate it.
func newPropagatedBlock(txs []*types.Transaction, uncles []*types.Header) *types.Block {
	header := &types.Header{
		Number:     big.NewInt(1),
		Difficulty: big.NewInt(131072),
		GasLimit:   8_000_000,
		Time:       1,
	}
	return types.NewBlock(header, txs, uncles, nil, trie.NewStackTrie(nil))
}

// Tests that a propagated block carrying exactly the allowed number of
// transactions and uncles is accepted and delivered to the backend.
func TestNewBlockAtLimits(t *testing.T) {
	t.Parallel()

	backend := &recordingBackend{testBackend: newTestBackend(1), handled: make(chan Packet, 4)}
	defer backend.close()

	peer, _ := newTestPeer("peer", ETH68, backend)
	defer peer.close()

	td := big.NewInt(131136)
	for i, block := range []*types.Block{
		newPropagatedBlock(makeTestTransactions(maxBlockTransactions), nil),
		newPropagatedBlock(nil, makeTestUncles(maxBlockUncles)),
	} {
		if err := p2p.Send(peer.app, NewBlockMsg, &NewBlockPacket{Block: block, TD: td}); err != nil {
			t.Fatalf("test %d: failed to send block: %v", i, err)
		}
		select {
		case packet := <-backend.handled:
			res, ok := packet.(*NewBlockPacket)
			if !ok {
				t.Fatalf("test %d: unexpected packet type delivered: %T", i, packet)
			}
			if res.Block.Hash() != block.Hash() {
				t.Fatalf("test %d: block hash mismatch: have %x, want %x", i, res.Block.Hash(), block.Hash())
			}
			if res.TD.Cmp(td) != 0 {
				t.Fatalf("test %d: block TD mismatch: have %v, want %v", i, res.TD, td)
			}
		case <-time.After(2 * time.Second):
			t.Fatalf("test %d: block at limit not delivered to backend", i)
		}
		if !peer.KnownBlock(block.Hash()) {
			t.Fatalf("test %d: block not marked as known to the peer", i)
		}
	}
}

// Tests that a propagated block exceeding the transaction or uncle bounds, or
// carrying a body element which proof-of-work blocks cannot have, is rejected
// with an error, dropping the sending peer.
func TestNewBlockExceedsLimits(t *testing.T) {
	t.Parallel()

	// withdrawalsBlock mirrors the block encoding with a fourth body element
	// appended, which rawBlock deliberately does not admit.
	type withdrawalsBlock struct {
		Header      *types.Header
		Txs         []*types.Transaction
		Uncles      []*types.Header
		Withdrawals []*types.Withdrawal
	}
	td := big.NewInt(131136)

	// The header of the withdrawals-carrying block commits to its (empty)
	// transaction and uncle lists, so only the extra body element is wrong
	// with it.
	header := &types.Header{
		Number:      big.NewInt(1),
		Difficulty:  big.NewInt(131072),
		GasLimit:    8_000_000,
		TxHash:      types.EmptyTxsHash,
		UncleHash:   types.EmptyUncleHash,
		ReceiptHash: types.EmptyReceiptsHash,
	}
	tests := []struct {
		name string
		msg  interface{}
	}{
		{
			name: "too many transactions",
			msg:  &NewBlockPacket{Block: newPropagatedBlock(makeTestTransactions(maxBlockTransactions+1), nil), TD: td},
		},
		{
			name: "too many uncles",
			msg:  &NewBlockPacket{Block: newPropagatedBlock(nil, makeTestUncles(maxBlockUncles+1)), TD: td},
		},
		{
			name: "withdrawals",
			msg: &struct {
				Block withdrawalsBlock
				TD    *big.Int
			}{
				Block: withdrawalsBlock{
					Header:      header,
					Withdrawals: []*types.Withdrawal{{Index: 1, Validator: 2, Address: common.Address{0x03}, Amount: 4}},
				},
				TD: td,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			backend := &recordingBackend{testBackend: newTestBackend(1), handled: make(chan Packet, 1)}
			defer backend.close()

			peer, errc := newTestPeer("peer", ETH68, backend)
			defer peer.close()

			if err := p2p.Send(peer.app, NewBlockMsg, tt.msg); err != nil {
				t.Fatalf("failed to send block: %v", err)
			}
			select {
			case err := <-errc:
				if err == nil {
					t.Fatalf("expected block to be rejected with an error")
				}
			case packet := <-backend.handled:
				t.Fatalf("out of bounds block delivered to backend: %v", packet)
			case <-time.After(2 * time.Second):
				t.Fatalf("peer not dropped for out of bounds block")
			}
		})
	}
}

// Tests that a propagated block whose encoded body matches the roots committed
// in its header, but does not decode into transactions or uncles, is rejected
// with an error when materialized, dropping the sending peer.
func TestNewBlockUndecodableBody(t *testing.T) {
	t.Parallel()

	// rawBlockOut mirrors the block encoding while keeping the body lists in
	// their encoded form, so undecodable items can be injected.
	type rawBlockOut struct {
		Header *types.Header
		Txs    rlp.RawList[*types.Transaction]
		Uncles rlp.RawList[*types.Header]
	}
	newHeader := func() *types.Header {
		return &types.Header{
			Number:      big.NewInt(1),
			Difficulty:  big.NewInt(131072),
			GasLimit:    8_000_000,
			TxHash:      types.EmptyTxsHash,
			UncleHash:   types.EmptyUncleHash,
			ReceiptHash: types.EmptyReceiptsHash,
		}
	}
	// A structurally valid RLP item that decodes into neither a transaction
	// nor a header.
	badItem := []byte{0x05}

	var badTxs rlp.RawList[*types.Transaction]
	if err := badTxs.AppendRaw(badItem); err != nil {
		t.Fatal(err)
	}
	var badUncles rlp.RawList[*types.Header]
	if err := badUncles.AppendRaw(badItem); err != nil {
		t.Fatal(err)
	}
	// Commit the junk lists into the headers, so that only materialization is
	// left to reject the blocks.
	txBlock := rawBlockOut{Header: newHeader(), Txs: badTxs}
	txBlock.Header.TxHash = types.DeriveSha(newDerivableRawList(&txBlock.Txs, writeTxForHash), trie.NewStackTrie(nil))

	uncleBlock := rawBlockOut{Header: newHeader(), Uncles: badUncles}
	uncleBlock.Header.UncleHash = crypto.Keccak256Hash(uncleBlock.Uncles.Bytes())

	for _, tt := range []struct {
		name  string
		block rawBlockOut
	}{
		{"transactions", txBlock},
		{"uncles", uncleBlock},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			backend := &recordingBackend{testBackend: newTestBackend(1), handled: make(chan Packet, 1)}
			defer backend.close()

			peer, errc := newTestPeer("peer", ETH68, backend)
			defer peer.close()

			msg := &struct {
				Block rawBlockOut
				TD    *big.Int
			}{Block: tt.block, TD: big.NewInt(131136)}

			if err := p2p.Send(peer.app, NewBlockMsg, msg); err != nil {
				t.Fatalf("failed to send block: %v", err)
			}
			select {
			case err := <-errc:
				if err == nil {
					t.Fatalf("expected block to be rejected with an error")
				}
			case packet := <-backend.handled:
				t.Fatalf("undecodable block delivered to backend: %v", packet)
			case <-time.After(2 * time.Second):
				t.Fatalf("peer not dropped for undecodable block")
			}
		})
	}
}

// Tests that a propagated block whose body does not match the roots committed
// in its header is discarded without dropping the peer. Keeping the peer is
// deliberate (see handleNewBlock), so this pins the warn-and-ignore behaviour
// rather than an error.
func TestNewBlockInvalidBodyKeepsPeer(t *testing.T) {
	t.Parallel()

	backend := &recordingBackend{testBackend: newTestBackend(1), handled: make(chan Packet, 1)}
	defer backend.close()

	peer, errc := newTestPeer("peer", ETH68, backend)
	defer peer.close()

	base := newPropagatedBlock(makeTestTransactions(2), makeTestUncles(1))

	malformedTxs := base.Header()
	malformedTxs.TxHash[0]++
	malformedUncles := base.Header()
	malformedUncles.UncleHash[0]++

	for i, header := range []*types.Header{malformedTxs, malformedUncles} {
		block := types.NewBlockWithHeader(header).WithBody(base.Transactions(), base.Uncles())
		if err := p2p.Send(peer.app, NewBlockMsg, &NewBlockPacket{Block: block, TD: big.NewInt(131136)}); err != nil {
			t.Fatalf("test %d: failed to send block: %v", i, err)
		}
		// The block must be dropped without killing the connection: a request
		// sent afterwards on the same peer must still be answered.
		probePeer(t, peer, backend.chain, uint64(i))

		select {
		case packet := <-backend.handled:
			t.Fatalf("test %d: mismatching block delivered to backend: %v", i, packet)
		case err := <-errc:
			t.Fatalf("test %d: peer dropped after mismatching block: %v", i, err)
		default:
		}
	}
}

// Tests the propagated block pre-decoding sanity checks.
func TestNewBlockSanityCheck(t *testing.T) {
	t.Parallel()

	head := func() *types.Header {
		return &types.Header{Number: big.NewInt(1), Difficulty: big.NewInt(131072)}
	}
	hugeExtra := head()
	hugeExtra.Extra = make([]byte, 100*1024+1)

	tests := []struct {
		name string
		pkt  *rawNewBlockPacket
		fail bool
	}{
		{"valid", &rawNewBlockPacket{Block: rawBlock{Header: head()}, TD: big.NewInt(131136)}, false},
		{"oversized td", &rawNewBlockPacket{Block: rawBlock{Header: head()}, TD: new(big.Int).Lsh(big.NewInt(1), 100)}, true},
		{"oversized extradata", &rawNewBlockPacket{Block: rawBlock{Header: hugeExtra}, TD: big.NewInt(131136)}, true},
	}
	for _, tt := range tests {
		if err := tt.pkt.sanityCheck(); (err != nil) != tt.fail {
			t.Errorf("test %s: sanity check mismatch: have error %v, want failure %v", tt.name, err, tt.fail)
		}
	}
}
