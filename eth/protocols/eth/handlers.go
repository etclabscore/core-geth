// Copyright 2021 The go-ethereum Authors
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
	"bytes"
	"encoding/json"
	"fmt"
	"math"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/p2p/tracker"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/trie"
)

func handleGetBlockHeaders(backend Backend, msg Decoder, peer *Peer) error {
	// Decode the complex header query
	var query GetBlockHeadersPacket
	if err := msg.Decode(&query); err != nil {
		return fmt.Errorf("%w: message %v: %v", errDecode, msg, err)
	}
	response := ServiceGetBlockHeadersQuery(backend.Chain(), query.GetBlockHeadersRequest, peer)
	return peer.ReplyBlockHeadersRLP(query.RequestId, response)
}

// ServiceGetBlockHeadersQuery assembles the response to a header query. It is
// exposed to allow external packages to test protocol behavior.
func ServiceGetBlockHeadersQuery(chain *core.BlockChain, query *GetBlockHeadersRequest, peer *Peer) []rlp.RawValue {
	if query.Skip == 0 {
		// The fast path: when the request is for a contiguous segment of headers.
		return serviceContiguousBlockHeaderQuery(chain, query)
	} else {
		return serviceNonContiguousBlockHeaderQuery(chain, query, peer)
	}
}

func serviceNonContiguousBlockHeaderQuery(chain *core.BlockChain, query *GetBlockHeadersRequest, peer *Peer) []rlp.RawValue {
	hashMode := query.Origin.Hash != (common.Hash{})
	first := true
	maxNonCanonical := uint64(100)

	// Gather headers until the fetch or network limits is reached
	var (
		bytes   common.StorageSize
		headers []rlp.RawValue
		unknown bool
		lookups int
	)
	for !unknown && len(headers) < int(query.Amount) && bytes < softResponseLimit &&
		len(headers) < maxHeadersServe && lookups < 2*maxHeadersServe {
		lookups++
		// Retrieve the next header satisfying the query
		var origin *types.Header
		if hashMode {
			if first {
				first = false
				origin = chain.GetHeaderByHash(query.Origin.Hash)
				if origin != nil {
					query.Origin.Number = origin.Number.Uint64()
				}
			} else {
				origin = chain.GetHeader(query.Origin.Hash, query.Origin.Number)
			}
		} else {
			origin = chain.GetHeaderByNumber(query.Origin.Number)
		}
		if origin == nil {
			break
		}
		if rlpData, err := rlp.EncodeToBytes(origin); err != nil {
			log.Crit("Unable to encode our own headers", "err", err)
		} else {
			headers = append(headers, rlp.RawValue(rlpData))
			bytes += common.StorageSize(len(rlpData))
		}
		// Advance to the next header of the query
		switch {
		case hashMode && query.Reverse:
			// Hash based traversal towards the genesis block
			ancestor := query.Skip + 1
			if ancestor == 0 {
				unknown = true
			} else {
				query.Origin.Hash, query.Origin.Number = chain.GetAncestor(query.Origin.Hash, query.Origin.Number, ancestor, &maxNonCanonical)
				unknown = (query.Origin.Hash == common.Hash{})
			}
		case hashMode && !query.Reverse:
			// Hash based traversal towards the leaf block
			var (
				current = origin.Number.Uint64()
				next    = current + query.Skip + 1
			)
			if next <= current {
				infos, _ := json.MarshalIndent(peer.Peer.Info(), "", "  ")
				peer.Log().Warn("GetBlockHeaders skip overflow attack", "current", current, "skip", query.Skip, "next", next, "attacker", infos)
				unknown = true
			} else {
				if header := chain.GetHeaderByNumber(next); header != nil {
					nextHash := header.Hash()
					expOldHash, _ := chain.GetAncestor(nextHash, next, query.Skip+1, &maxNonCanonical)
					if expOldHash == query.Origin.Hash {
						query.Origin.Hash, query.Origin.Number = nextHash, next
					} else {
						unknown = true
					}
				} else {
					unknown = true
				}
			}
		case query.Reverse:
			// Number based traversal towards the genesis block
			if query.Origin.Number >= query.Skip+1 {
				query.Origin.Number -= query.Skip + 1
			} else {
				unknown = true
			}

		case !query.Reverse:
			// Number based traversal towards the leaf block
			query.Origin.Number += query.Skip + 1
		}
	}
	return headers
}

func serviceContiguousBlockHeaderQuery(chain *core.BlockChain, query *GetBlockHeadersRequest) []rlp.RawValue {
	count := query.Amount
	if count > maxHeadersServe {
		count = maxHeadersServe
	}
	if query.Origin.Hash == (common.Hash{}) {
		// Number mode, just return the canon chain segment. The backend
		// delivers in [N, N-1, N-2..] descending order, so we need to
		// accommodate for that.
		from := query.Origin.Number
		if !query.Reverse {
			from = from + count - 1
		}
		headers := chain.GetHeadersFrom(from, count)
		if !query.Reverse {
			for i, j := 0, len(headers)-1; i < j; i, j = i+1, j-1 {
				headers[i], headers[j] = headers[j], headers[i]
			}
		}
		return headers
	}
	// Hash mode.
	var (
		headers []rlp.RawValue
		hash    = query.Origin.Hash
		header  = chain.GetHeaderByHash(hash)
	)
	if header != nil {
		rlpData, _ := rlp.EncodeToBytes(header)
		headers = append(headers, rlpData)
	} else {
		// We don't even have the origin header
		return headers
	}
	num := header.Number.Uint64()
	if !query.Reverse {
		// Theoretically, we are tasked to deliver header by hash H, and onwards.
		// However, if H is not canon, we will be unable to deliver any descendants of
		// H.
		if canonHash := chain.GetCanonicalHash(num); canonHash != hash {
			// Not canon, we can't deliver descendants
			return headers
		}
		descendants := chain.GetHeadersFrom(num+count-1, count-1)
		for i, j := 0, len(descendants)-1; i < j; i, j = i+1, j-1 {
			descendants[i], descendants[j] = descendants[j], descendants[i]
		}
		headers = append(headers, descendants...)
		return headers
	}
	{ // Last mode: deliver ancestors of H
		for i := uint64(1); header != nil && i < count; i++ {
			header = chain.GetHeaderByHash(header.ParentHash)
			if header == nil {
				break
			}
			rlpData, _ := rlp.EncodeToBytes(header)
			headers = append(headers, rlpData)
		}
		return headers
	}
}

func handleGetBlockBodies(backend Backend, msg Decoder, peer *Peer) error {
	// Decode the block body retrieval message
	var query GetBlockBodiesPacket
	if err := msg.Decode(&query); err != nil {
		return fmt.Errorf("%w: message %v: %v", errDecode, msg, err)
	}
	response := ServiceGetBlockBodiesQuery(backend.Chain(), query.GetBlockBodiesRequest)
	return peer.ReplyBlockBodiesRLP(query.RequestId, response)
}

// ServiceGetBlockBodiesQuery assembles the response to a body query. It is
// exposed to allow external packages to test protocol behavior.
func ServiceGetBlockBodiesQuery(chain *core.BlockChain, query GetBlockBodiesRequest) []rlp.RawValue {
	// Gather blocks until the fetch or network limits is reached
	var (
		bytes  int
		bodies []rlp.RawValue
	)
	for lookups, hash := range query {
		if bytes >= softResponseLimit || len(bodies) >= maxBodiesServe ||
			lookups >= 2*maxBodiesServe {
			break
		}
		data := chain.GetBodyRLP(hash)
		if len(data) == 0 {
			break // If we don't have this block's body, stop serving.
		}
		bodies = append(bodies, data)
		bytes += len(data)
	}
	return bodies
}

func handleGetReceipts(backend Backend, msg Decoder, peer *Peer) error {
	// Decode the block receipts retrieval message
	var query GetReceiptsPacket
	if err := msg.Decode(&query); err != nil {
		return fmt.Errorf("%w: message %v: %v", errDecode, msg, err)
	}
	response := ServiceGetReceiptsQuery(backend.Chain(), query.GetReceiptsRequest)
	return peer.ReplyReceiptsRLP(query.RequestId, response)
}

// ServiceGetReceiptsQuery assembles the response to a receipt query. It is
// exposed to allow external packages to test protocol behavior.
func ServiceGetReceiptsQuery(chain *core.BlockChain, query GetReceiptsRequest) []rlp.RawValue {
	// Gather state data until the fetch or network limits is reached
	var (
		bytes    int
		receipts []rlp.RawValue
	)
	for lookups, hash := range query {
		if bytes >= softResponseLimit || len(receipts) >= maxReceiptsServe ||
			lookups >= 2*maxReceiptsServe {
			break
		}
		// Retrieve the requested block's receipts
		results := chain.GetReceiptsByHash(hash)
		if results == nil {
			if header := chain.GetHeaderByHash(hash); header == nil || header.ReceiptHash != types.EmptyRootHash {
				break // Don't have this block's receipts, stop serving.
			}
		}
		// If known, encode and queue for response packet
		if encoded, err := rlp.EncodeToBytes(results); err != nil {
			log.Error("Failed to encode receipt", "err", err)
		} else {
			receipts = append(receipts, encoded)
			bytes += len(encoded)
		}
	}
	return receipts
}

func handleNewBlockhashes(backend Backend, msg Decoder, peer *Peer) error {
	// A batch of new block announcements just arrived
	ann := new(NewBlockHashesPacket)
	if err := msg.Decode(ann); err != nil {
		return fmt.Errorf("%w: message %v: %v", errDecode, msg, err)
	}
	// Mark the hashes as present at the remote node
	for _, block := range *ann {
		peer.markBlock(block.Hash)
	}
	// Deliver them all to the backend for queuing
	return backend.Handle(peer, ann)
}

func handleNewBlock(backend Backend, msg Decoder, peer *Peer) error {
	// Retrieve the propagated block, leaving its body encoded for now
	ann := new(rawNewBlockPacket)
	if err := msg.Decode(ann); err != nil {
		return fmt.Errorf("%w: message %v: %v", errDecode, msg, err)
	}
	if err := ann.sanityCheck(); err != nil {
		return err
	}
	// Bound what the body may expand into. Decoding into rlp.RawList only counts
	// the items, so this runs before any transaction or uncle is materialized.
	if txs := ann.Block.Transactions.Len(); txs > maxBlockTransactions {
		return fmt.Errorf("too many transactions in propagated block: %d > %d", txs, maxBlockTransactions)
	}
	if uncles := ann.Block.Uncles.Len(); uncles > maxBlockUncles {
		return fmt.Errorf("too many uncles in propagated block: %d > %d", uncles, maxBlockUncles)
	}
	// Check the body against the header on the encoded form, so a block that is
	// not what it claims to be is dropped without being materialized either.
	header := ann.Block.Header
	roots := hashBodyParts([]BlockBody{ann.Block.body()})
	if hash := roots.UncleHashes[0]; hash != header.UncleHash {
		log.Warn("Propagated block has invalid uncles", "have", hash, "exp", header.UncleHash)
		return nil // TODO(karalabe): return error eventually, but wait a few releases
	}
	if hash := roots.TransactionRoots[0]; hash != header.TxHash {
		log.Warn("Propagated block has invalid body", "have", hash, "exp", header.TxHash)
		return nil // TODO(karalabe): return error eventually, but wait a few releases
	}
	// Bounded and self-consistent, so it is safe to materialize now.
	txs, err := ann.Block.Transactions.Items()
	if err != nil {
		return fmt.Errorf("NewBlock: %w", err)
	}
	uncles, err := ann.Block.Uncles.Items()
	if err != nil {
		return fmt.Errorf("NewBlock: %w", err)
	}
	block := types.NewBlockWithHeader(header).WithBody(txs, uncles)
	block.ReceivedAt = msg.Time()
	block.ReceivedFrom = peer

	// Mark the peer as owning the block
	peer.markBlock(block.Hash())

	return backend.Handle(peer, &NewBlockPacket{Block: block, TD: ann.TD})
}

func handleBlockHeaders(backend Backend, msg Decoder, peer *Peer) error {
	// A batch of headers arrived to one of our previous requests
	res := new(BlockHeadersPacket)
	if err := msg.Decode(res); err != nil {
		return fmt.Errorf("%w: message %v: %v", errDecode, msg, err)
	}
	tresp := tracker.Response{ID: res.RequestId, MsgCode: BlockHeadersMsg, Size: res.List.Len()}
	if err := peer.tracker.Fulfil(tresp); err != nil {
		return fmt.Errorf("BlockHeaders: %w", err)
	}
	headers, err := res.List.Items()
	if err != nil {
		return fmt.Errorf("BlockHeaders: %w", err)
	}

	metadata := func() interface{} {
		hashes := make([]common.Hash, len(headers))
		for i, header := range headers {
			hashes[i] = header.Hash()
		}
		return hashes
	}
	return peer.dispatchResponse(&Response{
		id:   res.RequestId,
		code: BlockHeadersMsg,
		Res:  (*BlockHeadersRequest)(&headers),
	}, metadata)
}

func handleBlockBodies(backend Backend, msg Decoder, peer *Peer) error {
	// A batch of block bodies arrived to one of our previous requests
	res := new(BlockBodiesPacket)
	if err := msg.Decode(res); err != nil {
		return fmt.Errorf("%w: message %v: %v", errDecode, msg, err)
	}

	// Check against the request.
	length := res.List.Len()
	tresp := tracker.Response{ID: res.RequestId, MsgCode: BlockBodiesMsg, Size: length}
	if err := peer.tracker.Fulfil(tresp); err != nil {
		return fmt.Errorf("BlockBodies: %w", err)
	}

	// Collect items and dispatch.
	items, err := res.List.Items()
	if err != nil {
		return fmt.Errorf("BlockBodies: %w", err)
	}
	metadata := func() any { return hashBodyParts(items) }
	return peer.dispatchResponse(&Response{
		id:   res.RequestId,
		code: BlockBodiesMsg,
		Res:  (*BlockBodiesResponse)(&items),
	}, metadata)
}

// BlockBodyHashes contains the lists of block body part roots for a list of block bodies.
type BlockBodyHashes struct {
	TransactionRoots []common.Hash
	WithdrawalRoots  []common.Hash
	UncleHashes      []common.Hash
}

func hashBodyParts(items []BlockBody) BlockBodyHashes {
	h := BlockBodyHashes{
		TransactionRoots: make([]common.Hash, len(items)),
		WithdrawalRoots:  make([]common.Hash, len(items)),
		UncleHashes:      make([]common.Hash, len(items)),
	}
	hasher := trie.NewStackTrie(nil)
	for i, body := range items {
		// txs
		txsList := newDerivableRawList(&body.Transactions, writeTxForHash)
		h.TransactionRoots[i] = types.DeriveSha(txsList, hasher)
		// uncles
		if body.Uncles.Len() == 0 {
			h.UncleHashes[i] = types.EmptyUncleHash
		} else {
			h.UncleHashes[i] = crypto.Keccak256Hash(body.Uncles.Bytes())
		}
		// withdrawals
		if body.Withdrawals != nil {
			wdlist := newDerivableRawList(body.Withdrawals, nil)
			h.WithdrawalRoots[i] = types.DeriveSha(wdlist, hasher)
		}
	}
	return h
}

// derivableRawList implements types.DerivableList for a serialized RLP list.
type derivableRawList struct {
	data    []byte
	offsets []uint32
	write   func([]byte, *bytes.Buffer)
}

func newDerivableRawList[T any](list *rlp.RawList[T], write func([]byte, *bytes.Buffer)) *derivableRawList {
	dl := derivableRawList{data: list.Content(), write: write}
	if dl.write == nil {
		// default transform is identity
		dl.write = func(b []byte, buf *bytes.Buffer) { buf.Write(b) }
	}
	// Assert to ensure 32-bit offsets are valid. This can never trigger
	// unless a block body component or p2p receipt list is larger than 4GB.
	if uint(len(dl.data)) > math.MaxUint32 {
		panic("list data too big for derivableRawList")
	}
	it := list.ContentIterator()
	dl.offsets = make([]uint32, list.Len())
	for i := 0; it.Next(); i++ {
		dl.offsets[i] = uint32(it.Offset())
	}
	return &dl
}

// Len returns the number of items in the list.
func (dl *derivableRawList) Len() int {
	return len(dl.offsets)
}

// EncodeIndex writes the i'th item to the buffer.
func (dl *derivableRawList) EncodeIndex(i int, buf *bytes.Buffer) {
	start := dl.offsets[i]
	end := uint32(len(dl.data))
	if i != len(dl.offsets)-1 {
		end = dl.offsets[i+1]
	}
	dl.write(dl.data[start:end], buf)
}

// writeTxForHash changes a transaction in 'network encoding' into the format used for
// the transactions MPT.
func writeTxForHash(tx []byte, buf *bytes.Buffer) {
	k, content, _, _ := rlp.Split(tx)
	if k == rlp.List {
		buf.Write(tx) // legacy tx
	} else {
		buf.Write(content) // typed tx
	}
}

func handleReceipts(backend Backend, msg Decoder, peer *Peer) error {
	// A batch of receipts arrived to one of our previous requests
	res := new(ReceiptsPacket)
	if err := msg.Decode(res); err != nil {
		return fmt.Errorf("%w: message %v: %v", errDecode, msg, err)
	}

	// Check against the request.
	tresp := tracker.Response{ID: res.RequestId, MsgCode: ReceiptsMsg, Size: res.List.Len()}
	if err := peer.tracker.Fulfil(tresp); err != nil {
		return fmt.Errorf("Receipts: %w", err)
	}

	receiptLists, err := res.List.Items()
	if err != nil {
		return fmt.Errorf("Receipts: %w", err)
	}
	metadata := func() interface{} {
		hasher := trie.NewStackTrie(nil)
		hashes := make([]common.Hash, len(receiptLists))
		for i := range receiptLists {
			hashes[i] = types.DeriveSha(newDerivableRawList(&receiptLists[i], writeTxForHash), hasher)
		}
		return hashes
	}
	var enc ReceiptsRLPResponse
	for i := range receiptLists {
		enc = append(enc, receiptLists[i].Bytes())
	}
	return peer.dispatchResponse(&Response{
		id:   res.RequestId,
		code: ReceiptsMsg,
		Res:  &enc,
	}, metadata)
}

func handleNewPooledTransactionHashes(backend Backend, msg Decoder, peer *Peer) error {
	// New transaction announcement arrived, make sure we have
	// a valid and fresh chain to handle them
	if !backend.AcceptTxs() {
		return nil
	}
	ann := new(NewPooledTransactionHashesPacket)
	if err := msg.Decode(ann); err != nil {
		return fmt.Errorf("%w: message %v: %v", errDecode, msg, err)
	}
	if len(ann.Hashes) != len(ann.Types) || len(ann.Hashes) != len(ann.Sizes) {
		return fmt.Errorf("%w: message %v: invalid len of fields: %v %v %v", errDecode, msg, len(ann.Hashes), len(ann.Types), len(ann.Sizes))
	}
	// Schedule all the unknown hashes for retrieval
	for _, hash := range ann.Hashes {
		peer.MarkTransaction(hash)
	}
	return backend.Handle(peer, ann)
}

func handleGetPooledTransactions(backend Backend, msg Decoder, peer *Peer) error {
	// Decode the pooled transactions retrieval message
	var query GetPooledTransactionsPacket
	if err := msg.Decode(&query); err != nil {
		return fmt.Errorf("%w: message %v: %v", errDecode, msg, err)
	}
	hashes, txs := answerGetPooledTransactions(backend, query.GetPooledTransactionsRequest)
	return peer.ReplyPooledTransactionsRLP(query.RequestId, hashes, txs)
}

func answerGetPooledTransactions(backend Backend, query GetPooledTransactionsRequest) ([]common.Hash, []rlp.RawValue) {
	// Gather transactions until the fetch or network limits is reached
	var (
		bytes  int
		hashes []common.Hash
		txs    []rlp.RawValue
	)
	for _, hash := range query {
		if bytes >= softResponseLimit {
			break
		}
		// Retrieve the requested transaction, skipping if unknown to us
		tx := backend.TxPool().Get(hash)
		if tx == nil {
			continue
		}
		// If known, encode and queue for response packet
		if encoded, err := rlp.EncodeToBytes(tx); err != nil {
			log.Error("Failed to encode transaction", "err", err)
		} else {
			hashes = append(hashes, hash)
			txs = append(txs, encoded)
			bytes += len(encoded)
		}
	}
	return hashes, txs
}

func handleTransactions(backend Backend, msg Decoder, peer *Peer) error {
	// Transactions arrived, make sure we have a valid and fresh chain to handle them
	if !backend.AcceptTxs() {
		return nil
	}
	// Transactions can be processed, parse all of them and deliver to the pool
	var txs TransactionsPacket
	if err := msg.Decode(&txs); err != nil {
		return fmt.Errorf("%w: message %v: %v", errDecode, msg, err)
	}
	if txs.Len() > maxTransactionAnnouncements {
		return fmt.Errorf("too many transactions")
	}
	return backend.Handle(peer, &txs)
}

func handlePooledTransactions(backend Backend, msg Decoder, peer *Peer) error {
	// Transactions arrived, make sure we have a valid and fresh chain to handle them
	if !backend.AcceptTxs() {
		return nil
	}

	// Check against request and decode.
	var resp PooledTransactionsPacket
	if err := msg.Decode(&resp); err != nil {
		return fmt.Errorf("%w: message %v: %v", errDecode, msg, err)
	}
	tresp := tracker.Response{
		ID:      resp.RequestId,
		MsgCode: PooledTransactionsMsg,
		Size:    resp.List.Len(),
	}
	if err := peer.tracker.Fulfil(tresp); err != nil {
		return fmt.Errorf("PooledTransactions: %w", err)
	}

	return backend.Handle(peer, &resp)
}
