// Copyright 2019 The go-ethereum Authors
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

package clique

import (
	"bytes"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/params/types/ctypes"
	"github.com/ethereum/go-ethereum/params/types/genesisT"
	"github.com/ethereum/go-ethereum/params/vars"
)

// This test case is a repro of an annoying bug that took us forever to catch.
// In Clique PoA networks (Rinkeby, Görli, etc), consecutive blocks might have
// the same state root (no block subsidy, empty block). If a node crashes, the
// chain ends up losing the recent state and needs to regenerate it from blocks
// already in the database. The bug was that processing the block *prior* to an
// empty one **also completes** the empty one, ending up in a known-block error.
func TestReimportMirroredState(t *testing.T) {
	// Initialize a Clique chain with a single signer
	var (
		db     = rawdb.NewMemoryDatabase()
		key, _ = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		addr   = crypto.PubkeyToAddress(key.PublicKey)
		engine = New(params.AllCliqueProtocolChanges.Clique, db)
		signer = new(types.HomesteadSigner)
	)
	genspec := &genesisT.Genesis{
		ExtraData: make([]byte, extraVanity+common.AddressLength+extraSeal),
		Alloc: map[common.Address]genesisT.GenesisAccount{
			addr: {Balance: big.NewInt(1)},
		},
	}
	copy(genspec.ExtraData[extraVanity:], addr[:])
	genesis := core.MustCommitGenesis(db, genspec)

	// Generate a batch of blocks, each properly signed
	chain, _ := core.NewBlockChain(db, nil, params.AllCliqueProtocolChanges, engine, vm.Config{}, nil, nil)
	defer chain.Stop()

	blocks, _ := core.GenerateChain(params.AllCliqueProtocolChanges, genesis, engine, db, 3, func(i int, block *core.BlockGen) {
		// The chain maker doesn't have access to a chain, so the difficulty will be
		// lets unset (nil). Set it here to the correct value.
		block.SetDifficulty(diffInTurn)

		// We want to simulate an empty middle block, having the same state as the
		// first one. The last is needs a state change again to force a reorg.
		if i != 1 {
			tx, err := types.SignTx(types.NewTransaction(block.TxNonce(addr), common.Address{0x00}, new(big.Int), vars.TxGas, nil, nil), signer, key)
			if err != nil {
				panic(err)
			}
			block.AddTxWithChain(chain, tx)
		}
	})
	for i, block := range blocks {
		header := block.Header()
		if i > 0 {
			header.ParentHash = blocks[i-1].Hash()
		}
		header.Extra = make([]byte, extraVanity+extraSeal)
		header.Difficulty = diffInTurn

		sig, _ := crypto.Sign(SealHash(header).Bytes(), key)
		copy(header.Extra[len(header.Extra)-extraSeal:], sig)
		blocks[i] = block.WithSeal(header)
	}
	// Insert the first two blocks and make sure the chain is valid
	db = rawdb.NewMemoryDatabase()
	core.MustCommitGenesis(db, genspec)

	chain, _ = core.NewBlockChain(db, nil, params.AllCliqueProtocolChanges, engine, vm.Config{}, nil, nil)
	defer chain.Stop()

	if _, err := chain.InsertChain(blocks[:2]); err != nil {
		t.Fatalf("failed to insert initial blocks: %v", err)
	}
	if head := chain.CurrentBlock().NumberU64(); head != 2 {
		t.Fatalf("chain head mismatch: have %d, want %d", head, 2)
	}

	// Simulate a crash by creating a new chain on top of the database, without
	// flushing the dirty states out. Insert the last block, triggering a sidechain
	// reimport.
	chain, _ = core.NewBlockChain(db, nil, params.AllCliqueProtocolChanges, engine, vm.Config{}, nil, nil)
	defer chain.Stop()

	if _, err := chain.InsertChain(blocks[2:]); err != nil {
		t.Fatalf("failed to insert final block: %v", err)
	}
	if head := chain.CurrentBlock().NumberU64(); head != 3 {
		t.Fatalf("chain head mismatch: have %d, want %d", head, 3)
	}
}

func TestClique_EIP3436_Scenario1(t *testing.T) {

	scenarios := []struct {
		lenSigners   int
		commonBlocks []int
		forks        [][]int
	}{
		{
			// SCENARIO-1
			// signers: A..H (8 count)
			//
			/*
				Step 1
				A fully in-order chain exists and validator 8 has just produced an in-turn block.

				1A x
				2B  x
				3C   x
				4D    x
				5E     x
				6F      x
				7G       x
				8H        x

				Step 2
				... and then validators 5, 7, and 8 go offline.

				1A x
				2B  x
				3C   x
				4D    x
				5E     x   -
				6F      x
				7G       x -
				8H        x-

				Step 3
				Two forks form, one with an in-order block from validator 1
				and then an out of order block from validator 3.

				The second fork forms from validators 2, 4, and 6 in order.
				Both have a net total difficulty of 3 more than the common ancestor.

				1A x        y
				2B  x       z
				3C   x       y
				4D    x      z
				5E     x   -
				6F      x     z
				7G       x -
				8H        x-
			*/
			lenSigners:   8,
			commonBlocks: []int{0, 1, 2, 3, 4, 5, 6, 7},
			forks: [][]int{
				{0, 2},
				{1, 3, 5},
			},
		},
	}

	// validators: 8
	//

	for ii, tt := range scenarios {
		// Create the account pool and generate the initial set of signerAddressesSorted
		accountsPool := newTesterAccountPool()

		db := rawdb.NewMemoryDatabase()

		// Assemble a chain of headers from the cast votes
		config := *params.TestChainConfig
		cliquePeriod := uint64(1)
		config.Clique = &ctypes.CliqueConfig{
			Period: cliquePeriod,
			Epoch:  0,
		}
		engine := New(config.Clique, db)
		engine.fakeDiff = false

		signerAddressesSorted := make([]common.Address, tt.lenSigners)
		for i := 0; i < tt.lenSigners; i++ {
			signerAddressesSorted[i] = accountsPool.address(fmt.Sprintf("%s", []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K"}[i]))
		}
		for j := 0; j < len(signerAddressesSorted); j++ {
			for k := j + 1; k < len(signerAddressesSorted); k++ {
				if bytes.Compare(signerAddressesSorted[j][:], signerAddressesSorted[k][:]) > 0 {
					signerAddressesSorted[j], signerAddressesSorted[k] = signerAddressesSorted[k], signerAddressesSorted[j]
				}
			}
		}

		logSortedSigners := ""
		for j, s := range signerAddressesSorted {
			logSortedSigners += fmt.Sprintf("%d: (%s) %s\n", j, accountsPool.name(s), s.Hex())
		}
		t.Logf("SORTED SIGNERS:\n%s\n", logSortedSigners)

		// Create the genesis block with the initial set of signerAddressesSorted
		genesis := &genesisT.Genesis{
			ExtraData: make([]byte, extraVanity+common.AddressLength*len(signerAddressesSorted)+extraSeal),
		}
		for j, signer := range signerAddressesSorted {
			copy(genesis.ExtraData[extraVanity+j*common.AddressLength:], signer[:])
		}

		genesisBlock := core.MustCommitGenesis(db, genesis)

		// Create a pristine blockchain with the genesis injected
		chain, err := core.NewBlockChain(db, nil, &config, engine, vm.Config{}, nil, nil)
		if err != nil {
			t.Errorf("test %d: failed to create test chain: %v", ii, err)
			continue
		}

		// mustMakeChain := func(parent *types.Block, blocksSigners []string)

		// Build the common segment
		commonSegmentBlocks := []*types.Block{genesisBlock}
		for i := 0; i < len(tt.commonBlocks); i++ {

			signerIndex := tt.commonBlocks[i]
			signerAddress := signerAddressesSorted[signerIndex]
			signerName := accountsPool.name(signerAddress)

			parentBlock := commonSegmentBlocks[len(commonSegmentBlocks)-1]

			generatedBlocks, _ := core.GenerateChain(&config, parentBlock, engine, db, 1, nil)
			block := generatedBlocks[0]

			// Get the header and prepare it for signing
			header := block.Header()
			header.ParentHash = parentBlock.Hash()

			header.Time = parentBlock.Time() + cliquePeriod

			// See if our required signer is in or out of turn and assign difficulty respectively.
			difficulty := diffInTurn

			// If the snapshot reports this signer is out of turn, use out-of-turn difficulty.
			snap, err := engine.snapshot(chain, parentBlock.NumberU64(), parentBlock.Hash(), nil)
			if err != nil {
				t.Fatalf("snap err: %v", err)
			}
			if !snap.inturn(parentBlock.NumberU64()+1, signerAddress) {
				difficulty = diffNoTurn
			}
			header.Difficulty = difficulty

			// Generate the signature, embed it into the header and the block.
			header.Extra = make([]byte, extraVanity+extraSeal) // allocate byte slice

			t.Logf("SIGNING: %d (%s) %s\n", i+1, signerName, signerAddress.Hex())

			// Sign the header with the associated validator's key.
			accountsPool.sign(header, signerName)

			// Double check to see what Clique thinks the signer of this block is.
			author, err := engine.Author(header)
			if err != nil {
				t.Fatalf("author error: %v", err)
			}
			if wantSigner := accountsPool.address(signerName); author != wantSigner {
				t.Fatalf("header author != wanted signer: author: %s, signer: %s", author.Hex(), wantSigner.Hex())
			}

			generatedBlocks[0] = block.WithSeal(header)

			commonSegmentBlocks = append(commonSegmentBlocks, generatedBlocks...) // == generatedBlocks[0]

			if k, err := chain.InsertChain(generatedBlocks); err != nil || k != 1 {
				t.Fatalf("case: %d, failed to import block %d, count: %d, err: %v", ii, i, k, err)
			}
		}

		// for scenarioForkIndex, forkBlockSigners := range tt.forks {
		//
		// }

	}
}
