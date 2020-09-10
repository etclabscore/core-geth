package core

import (
	"math"
	"math/rand"
	"testing"

	"github.com/ethereum/go-ethereum/consensus/ethash"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/params"
)

func TestBlockChain_AF_ECBP11355(t *testing.T) {

	cases := []struct {
		easyLen, hardLen, commonAncestorN int
		easyOffset, hardOffset            int64
		hardGetsHead, accepted            bool
	}{
		// Hard has insufficient total difficulty / length and is rejected.
		{
			5000, 7500, 2500,
			50, -9,
			false, false,
		},
		// Hard has sufficient total difficulty / length and is accepted.
		{
			1000, 7, 995,
			60, 0,
			true, true,
		},
		// Hard has sufficient total difficulty / length and is accepted.
		{
			1000, 7, 995,
			60, 7,
			true, true,
		},
		// Hard has sufficient total difficulty / length and is accepted.
		{
			1000, 1, 999,
			30, 1,
			true, true,
		},
		// Hard has sufficient total difficulty / length and is accepted.
		{
			500, 3, 497,
			0, -8,
			true, true,
		},
		// Hard has sufficient total difficulty / length and is accepted.
		{
			500, 4, 496,
			0, -9,
			true, true,
		},
		// Hard has sufficient total difficulty / length and is accepted.
		{
			500, 5, 495,
			0, -9,
			true, true,
		},
		// Hard has sufficient total difficulty / length and is accepted.
		{
			500, 6, 494,
			0, -9,
			true, true,
		},
		// Hard has sufficient total difficulty / length and is accepted.
		{
			500, 7, 493,
			0, -9,
			true, true,
		},
		// Hard has sufficient total difficulty / length and is accepted.
		{
			500, 8, 492,
			0, -9,
			true, true,
		},
		// Hard has sufficient total difficulty / length and is accepted.
		{
			500, 9, 491,
			0, -9,
			true, true,
		},
		// Hard has sufficient total difficulty / length and is accepted.
		{
			500, 12, 488,
			0, -9,
			true, true,
		},
		// Hard has sufficient total difficulty / length and is accepted.
		{
			500, 20, 480,
			0, -9,
			true, true,
		},
		// Hard has sufficient total difficulty / length and is accepted.
		{
			500, 40, 460,
			0, -9,
			true, true,
		},
		// Hard has sufficient total difficulty / length and is accepted.
		{
			500, 60, 440,
			0, -9,
			true, true,
		},
		// Hard has insufficient total difficulty / length and is rejected.
		{
			500, 80, 420,
			0, -9,
			false, false,
		},
		// Hard has insufficient total difficulty / length and is rejected.
		{
			500, 80, 420,
			7, -9,
			false, false,
		},
		// Hard has insufficient total difficulty / length and is rejected.
		{
			500, 80, 420,
			17, -9,
			false, false,
		},
		// Hard has sufficient total difficulty / length and is accepted.
		{
			500, 80, 420,
			47, -9,
			true, true,
		},
		// Hard has insufficient total difficulty / length and is rejected.
		{
			500, 80, 420,
			47, -8,
			false, false,
		},
		// Hard has insufficient total difficulty / length and is rejected.
		{
			500, 80, 420,
			17, -8,
			false, false,
		},
		// Hard has insufficient total difficulty / length and is rejected.
		{
			500, 80, 420,
			7, -8,
			false, false,
		},
		// Hard has insufficient total difficulty / length and is rejected.
		{
			500, 80, 420,
			0, -8,
			false, false,
		},
		// Hard has insufficient total difficulty / length and is rejected.
		{
			500, 40, 460,
			0, -7,
			false, false,
		},
		// Hard has insufficient total difficulty / length and is rejected.
		{
			500, 14, 486,
			0, -7,
			false, false,
		},
		// Hard is accepted, but does not have greater total difficulty,
		// and is not set as the chain head.
		{
			1000, 1, 900,
			60, -9,
			false, true,
		},
	}

	runTest := func(easyL, hardL, caN int, easyT, hardT int64) (hardHead bool, err error) {
		// Generate the original common chain segment and the two competing forks
		engine := ethash.NewFaker()

		db := rawdb.NewMemoryDatabase()
		genesis := params.DefaultMessNetGenesisBlock()
		genesisB := MustCommitGenesis(db, genesis)

		chain, err := NewBlockChain(db, nil, genesis.Config, engine, vm.Config{}, nil, nil)
		if err != nil {
			t.Fatal(err)
		}
		defer chain.Stop()
		chain.EnableArtificialFinality(true)

		easy, _ := GenerateChain(genesis.Config, genesisB, engine, db, easyL, func(i int, b *BlockGen) {
			b.SetNonce(types.EncodeNonce(uint64(rand.Int63n(math.MaxInt64))))
			b.OffsetTime(easyT)
		})
		commonAncestor := easy[caN-1]
		hard, _ := GenerateChain(genesis.Config, commonAncestor, engine, db, hardL, func(i int, b *BlockGen) {
			b.SetNonce(types.EncodeNonce(uint64(rand.Int63n(math.MaxInt64))))
			b.OffsetTime(hardT)
		})

		if _, err := chain.InsertChain(easy); err != nil {
			t.Fatal(err)
		}
		_, err = chain.InsertChain(hard)
		hardHead = chain.CurrentBlock().Hash() == hard[len(hard)-1].Hash()
		return
	}

	for i, c := range cases {
		hardHead, err := runTest(c.easyLen, c.hardLen, c.commonAncestorN, c.easyOffset, c.hardOffset)
		if (err != nil && c.accepted) || (err == nil && !c.accepted) || (hardHead != c.hardGetsHead) {
			t.Errorf("case=%d want.accepted=%v want.hardHead=%v got.hardHead=%v err=%v",
				i, c.accepted, c.hardGetsHead, hardHead, err)
		}
	}
}
