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

	cases := []struct{
		easyLen, hardLen, commonAncestorN int
		easyOffset, hardOffset int64
		hardGetsHead, accepted bool
	}{
		// Hard has insufficient total difficulty / length.
		{
			5000, 7500, 2500,
			60, 1,
			false, false,
		},
		// Hard has insufficient total difficulty / length.
		{
			1000, 7, 995,
			60, 9,
			false, false,
		},
		// Hard has sufficient total difficulty / length to be accepted and set as head.
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
		// Hard is accepted, but does not have greater total difficulty,
		// and is not set as the chain head.
		{
			1000, 1, 900,
			60, 1,
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
