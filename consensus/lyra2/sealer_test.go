package lyra2

import (
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Tests whether stale solutions are correctly processed.
func TestStaleSubmission(t *testing.T) {
	lyra2 := NewTester(nil, true)
	defer lyra2.Close()
	api := &API{lyra2}

	fakeNonce, fakeDigest := types.BlockNonce{0x01, 0x02, 0x03}, common.HexToHash("00")

	testcases := []struct {
		headers     []*types.Header
		submitIndex int
		submitRes   bool
	}{
		// Case1: submit solution for the latest mining package
		{
			[]*types.Header{
				{ParentHash: common.BytesToHash([]byte{0xa}), Number: big.NewInt(1), Difficulty: big.NewInt(100000000)},
			},
			0,
			true,
		},
		// Case2: submit solution for the previous package but have same parent.
		{
			[]*types.Header{
				{ParentHash: common.BytesToHash([]byte{0xb}), Number: big.NewInt(2), Difficulty: big.NewInt(100000000)},
				{ParentHash: common.BytesToHash([]byte{0xb}), Number: big.NewInt(2), Difficulty: big.NewInt(100000001)},
			},
			0,
			true,
		},
		// Case3: submit stale but acceptable solution
		{
			[]*types.Header{
				{ParentHash: common.BytesToHash([]byte{0xc}), Number: big.NewInt(3), Difficulty: big.NewInt(100000000)},
				{ParentHash: common.BytesToHash([]byte{0xd}), Number: big.NewInt(9), Difficulty: big.NewInt(100000000)},
			},
			0,
			true,
		},
		// Case4: submit very old solution
		{
			[]*types.Header{
				{ParentHash: common.BytesToHash([]byte{0xe}), Number: big.NewInt(10), Difficulty: big.NewInt(100000000)},
				{ParentHash: common.BytesToHash([]byte{0xf}), Number: big.NewInt(17), Difficulty: big.NewInt(100000000)},
			},
			0,
			false,
		},
	}
	results := make(chan *types.Block, 16)

	for id, c := range testcases {
		for _, h := range c.headers {
			lyra2.Seal(nil, types.NewBlockWithHeader(h), results, nil)
		}
		if res := api.SubmitWork(fakeNonce, lyra2.SealHash(c.headers[c.submitIndex]), fakeDigest); res != c.submitRes {
			t.Errorf("case %d submit result mismatch, want %t, get %t", id+1, c.submitRes, res)
		}
		if !c.submitRes {
			continue
		}
		select {
		case res := <-results:
			if res.Header().Nonce != fakeNonce {
				t.Errorf("case %d block nonce mismatch, want %x, get %x", id+1, fakeNonce, res.Header().Nonce)
			}
			if res.Header().MixDigest != fakeDigest {
				t.Errorf("case %d block digest mismatch, want %x, get %x", id+1, fakeDigest, res.Header().MixDigest)
			}
			if res.Header().Difficulty.Uint64() != c.headers[c.submitIndex].Difficulty.Uint64() {
				t.Errorf("case %d block difficulty mismatch, want %d, get %d", id+1, c.headers[c.submitIndex].Difficulty, res.Header().Difficulty)
			}
			if res.Header().Number.Uint64() != c.headers[c.submitIndex].Number.Uint64() {
				t.Errorf("case %d block number mismatch, want %d, get %d", id+1, c.headers[c.submitIndex].Number.Uint64(), res.Header().Number.Uint64())
			}
			if res.Header().ParentHash != c.headers[c.submitIndex].ParentHash {
				t.Errorf("case %d block parent hash mismatch, want %s, get %s", id+1, c.headers[c.submitIndex].ParentHash.Hex(), res.Header().ParentHash.Hex())
			}
		case <-time.NewTimer(time.Second).C:
			t.Errorf("case %d fetch lyra2 result timeout", id+1)
		}
	}
}
