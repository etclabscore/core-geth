package params

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/params/types/genesisT"
	"github.com/ethereum/go-ethereum/params/vars"
	"github.com/ethereum/go-ethereum/trie"
)

// TestGenesisHashABC tests that ABCGenesisHash is the correct value for the genesis configuration.
func TestGenesisHashABC(t *testing.T) {
	genesis := DefaultABCGenesisBlock()
	block := genesisToBlock(genesis, nil)
	if block.Hash() != ABCGenesisHash {
		t.Errorf("want: %s, got: %s", ABCGenesisHash.Hex(), block.Hash().Hex())
	}
}

// genesisToBlock is a helper function to convert the genesis type to block type,
// which will allow us to generate the block hash.
// This duplicates the logic in core package (which cannot be imported due to import cycles).
func genesisToBlock(g *genesisT.Genesis, db ethdb.Database) *types.Block {
	if db == nil {
		db = rawdb.NewMemoryDatabase()
	}
	statedb, _ := state.New(common.Hash{}, state.NewDatabase(db), nil)
	for addr, account := range g.Alloc {
		statedb.AddBalance(addr, account.Balance)
		statedb.SetCode(addr, account.Code)
		statedb.SetNonce(addr, account.Nonce)
		for key, value := range account.Storage {
			statedb.SetState(addr, key, value)
		}
	}
	root := statedb.IntermediateRoot(false)
	head := &types.Header{
		Number:     new(big.Int).SetUint64(g.Number),
		Nonce:      types.EncodeNonce(g.Nonce),
		Time:       g.Timestamp,
		ParentHash: g.ParentHash,
		Extra:      g.ExtraData,
		GasLimit:   g.GasLimit,
		GasUsed:    g.GasUsed,
		Difficulty: g.Difficulty,
		MixDigest:  g.Mixhash,
		Coinbase:   g.Coinbase,
		Root:       root,
	}
	if g.GasLimit == 0 {
		head.GasLimit = vars.GenesisGasLimit
	}
	if g.Difficulty == nil {
		head.Difficulty = vars.GenesisDifficulty
	}
	statedb.Commit(false)
	statedb.Database().TrieDB().Commit(root, true, nil)

	return types.NewBlock(head, nil, nil, nil, new(trie.Trie))
}
