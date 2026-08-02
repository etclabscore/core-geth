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
	"bytes"
	"testing"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/rlp"
)

// TestCheckResponseItems_CVE_2026_26313 verifies that the pre-decode item count
// validation rejects messages with too many items, preventing memory amplification
// attacks (CVE-2026-26313).
func TestCheckResponseItems_CVE_2026_26313(t *testing.T) {
	// Build a wrapped packet (BlockHeadersPacket) with N minimal headers.
	buildHeadersMsg := func(n int) p2p.Msg {
		headers := make([]*types.Header, n)
		for i := range headers {
			headers[i] = &types.Header{}
		}
		pkt := &BlockHeadersPacket{
			RequestId: 1,
			List:      encodeRL(headers),
		}
		payload, err := rlp.EncodeToBytes(pkt)
		if err != nil {
			t.Fatal(err)
		}
		return p2p.Msg{
			Code:    BlockHeadersMsg,
			Size:    uint32(len(payload)),
			Payload: bytes.NewReader(payload),
		}
	}

	// Build a bare packet (TransactionsPacket) with N minimal transactions.
	buildTxsMsg := func(n int) p2p.Msg {
		txs := make(TransactionsPacket, n)
		for i := range txs {
			txs[i] = types.NewTx(&types.LegacyTx{})
		}
		payload, err := rlp.EncodeToBytes(txs)
		if err != nil {
			t.Fatal(err)
		}
		return p2p.Msg{
			Code:    TransactionsMsg,
			Size:    uint32(len(payload)),
			Payload: bytes.NewReader(payload),
		}
	}

	t.Run("wrapped packet within limit passes", func(t *testing.T) {
		msg := buildHeadersMsg(maxHeadersServe)
		limit := responseItemLimits[BlockHeadersMsg]
		if err := checkResponseItems(&msg, limit); err != nil {
			t.Fatalf("expected no error for %d headers, got: %v", maxHeadersServe, err)
		}
	})

	t.Run("wrapped packet exceeding limit rejected", func(t *testing.T) {
		msg := buildHeadersMsg(maxHeadersServe + 1)
		limit := responseItemLimits[BlockHeadersMsg]
		if err := checkResponseItems(&msg, limit); err == nil {
			t.Fatalf("expected error for %d headers (limit %d), got nil", maxHeadersServe+1, maxHeadersServe)
		}
	})

	t.Run("bare packet within limit passes", func(t *testing.T) {
		msg := buildTxsMsg(100)
		limit := responseItemLimits[TransactionsMsg]
		if err := checkResponseItems(&msg, limit); err != nil {
			t.Fatalf("expected no error for 100 txs, got: %v", err)
		}
	})

	t.Run("bare packet exceeding limit rejected", func(t *testing.T) {
		msg := buildTxsMsg(maxHeadersServe*4 + 1)
		limit := responseItemLimits[TransactionsMsg]
		if err := checkResponseItems(&msg, limit); err == nil {
			t.Fatalf("expected error for %d txs (limit %d), got nil", maxHeadersServe*4+1, maxHeadersServe*4)
		}
	})

	t.Run("payload still decodable after check", func(t *testing.T) {
		msg := buildHeadersMsg(10)
		limit := responseItemLimits[BlockHeadersMsg]
		if err := checkResponseItems(&msg, limit); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		// Verify the payload can still be decoded after the check
		res := new(BlockHeadersPacket)
		if err := msg.Decode(res); err != nil {
			t.Fatalf("failed to decode after check: %v", err)
		}
		if res.List.Len() != 10 {
			t.Fatalf("expected 10 headers, got %d", res.List.Len())
		}
	})
}
