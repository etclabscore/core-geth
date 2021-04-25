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
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/params/types/ctypes"
	"github.com/ethereum/go-ethereum/params/types/genesisT"
	"github.com/ethereum/go-ethereum/params/vars"
	"github.com/holiman/uint256"
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

type cliqueEIP3436TestCase struct {
	// The number of signers (aka validators).
	// These addresses will be randomly generated on demand for each scenario
	// and sorted per status quo Clique spec.
	// They are referenced by their cardinality in this sorted list.
	lenSigners int

	// commonBlocks defines blocks by which validator should sign them.
	// Validators are referenced by their index in the sorted signers list,
	// which sorts by address.
	commonBlocks []int

	// forks defines blocks by which validator should sign them.
	// These forks will be generated and attempted to be imported into
	// the chain. We expect that import should always succeed for all blocks,
	// with this test measuring only if the expected fork head actually gets
	// canonical preference.
	forks [][]int

	// assertions are functions used to assert the expectations of
	// fork heads in chain context against the specification's requirements, eg.
	// make sure that the forks really do have equal total difficulty, or equal block numbers.
	// This way we know that the deciding condition for achieving canonical status
	// really is what we think it is (and not, eg. total difficulty).

	assertions []func(t *testing.T, chain *core.BlockChain, forkHeads ...*types.Header)

	// canonicalForkIndex takes variadic fork heads and tells us which
	// should get canonical preference (by index).
	canonicalForkIndex func(t *testing.T, forkHeads ...*types.Header) int

	cliqueConfig *ctypes.CliqueConfig
}

var cliqueEIP3436TestPeriod = uint64(1)
var cliqueConfigEIP3436 = &ctypes.CliqueConfig{
	Period:            cliqueEIP3436TestPeriod,
	Epoch:             0,
	EIP3436Transition: big.NewInt(0),
}
var cliqueConfigNoEIP3436 = &ctypes.CliqueConfig{
	Period:            cliqueEIP3436TestPeriod,
	Epoch:             0,
	EIP3436Transition: nil,
}

func TestCliqueEIP3436_Scenario1_positive(t *testing.T) {
	testCliqueEIP3436(t, cliqueEIP3436TestCase{
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
		commonBlocks: []int{1, 2, 3, 4, 5, 6, 7, 0, 1, 2, 3, 4, 5, 6, 7},
		forks: [][]int{
			{1, 3, 5}, // 2, 4, 6
			{0, 2},    // 1, 3
		},
		canonicalForkIndex: shorterFork,
		assertions: []func(t *testing.T, chain *core.BlockChain, forkHeads ...*types.Header){
			assertEqualTotalDifficulties,
		},
		cliqueConfig: cliqueConfigEIP3436,
	})
}

func TestCliqueEIP3436_Scenario1_negative(t *testing.T) {
	testCliqueEIP3436(t, cliqueEIP3436TestCase{
		// SCENARIO-1 without EIP-3436 (negative test)
		lenSigners:   8,
		commonBlocks: []int{1, 2, 3, 4, 5, 6, 7, 0, 1, 2, 3, 4, 5, 6, 7},
		forks: [][]int{
			{1, 3, 5}, // 2, 4, 6
			{0, 2},    // 1, 3
		},
		canonicalForkIndex: shorterFork,
		assertions: []func(t *testing.T, chain *core.BlockChain, forkHeads ...*types.Header){
			assertEqualTotalDifficulties,
		},
		cliqueConfig: cliqueConfigNoEIP3436,
	})
}

func TestCliqueEIP3436_Scenario2_positive(t *testing.T) {
	testCliqueEIP3436(t, cliqueEIP3436TestCase{
		// SCENARIO-2

		/*
				Step 1
				For the second scenario with the same validator set and in-order chain with
				validator 7 having just produced an in order block, then validators 7 and 8 go offline.

				1A x
				2B  x
				3C   x
				4D    x
				5E     x
				6F      x
				7G       x
				8H

				1A x
				2B  x
				3C   x
				4D    x
				5E     x
				6F      x
				7G       x-
				8H        -

				Two forks form, 1,3,5 on one side and 2,4,6 on the other.
				Both forks become aware of the other fork after producing their third block.
				In this case both forks have equal total difficulty and equal length.

				1A x       x
				2B  x      y
				3C   x      x
				4D    x     y
				5E     x     x
				6F      x    y
				7G       x-
				8H        -

			FIXME(meowsbits): This scenario yields a "recently signed" error
			when attempting to import Signer 5 (really #6 b/c zero-indexing) into the
			second fork.

			On that fork, the sequence of signers is specified to be
			... 7, 0, 1, 2, 3, 4, 5, 6, 1, 3, 5

			(vs. the other fork)
			... 7, 0, 1, 2, 3, 4, 5, 6, 0, 2, 4

			The condition for "recently signed" is (from *Clique#verifySeal):

			// Signer is among recents, only fail if the current block doesn't shift it out
			if limit := uint64(len(snap.Signers)/2 + 1); seen > number-limit {
				return errRecentlySigned
			}

			Evaluated, this yields

			=> recently signed: limit=(8/2+1)=5 seen=13 number=17 number-limit=12

			RESOLUTION (tentative): Use forks with length 2 instead of 3.

		*/
		lenSigners:   8,
		commonBlocks: []int{1, 2, 3, 4, 5, 6, 7, 0, 1, 2, 3, 4, 5, 6},
		forks: [][]int{
			// These expected values have been truncated to 2 block each instead
			// of the EIP's described 3 blocks (per inline comments).
			{0, 2}, // 1, 3, 5
			{1, 3}, // 2, 4, 6
		},
		canonicalForkIndex: lowerHash,
		assertions: []func(t *testing.T, chain *core.BlockChain, forkHeads ...*types.Header){
			assertEqualTotalDifficulties,
			assertEqualNumbers,
		},
		cliqueConfig: cliqueConfigEIP3436,
	})
}

func TestCliqueEIP3436_Scenario2_negative(t *testing.T) {
	testCliqueEIP3436(t, cliqueEIP3436TestCase{
		// SCENARIO-2 without EIP-3436 (negative test)
		lenSigners:   8,
		commonBlocks: []int{1, 2, 3, 4, 5, 6, 7, 0, 1, 2, 3, 4, 5, 6},
		forks: [][]int{
			// These expected values have been truncated to 2 block each instead
			// of the EIP's described 3 blocks (per inline comments).
			{0, 2}, // 1, 3, 5
			{1, 3}, // 2, 4, 6
		},
		canonicalForkIndex: func(t *testing.T, forkHeads ...*types.Header) int {
			// Prefer the first-seen fork.
			return 0
		},
		assertions: []func(t *testing.T, chain *core.BlockChain, forkHeads ...*types.Header){
			assertEqualTotalDifficulties,
			assertEqualNumbers,
		},
		cliqueConfig: cliqueConfigNoEIP3436,
	})
}

func TestCliqueEIP3436_Scenario3_positive(t *testing.T) {
	testCliqueEIP3436(t, cliqueEIP3436TestCase{
		// SCENARIO-3
		// https://ethereum-magicians.org/t/eip-3436-expanded-clique-block-choice-rule/5809/3?u=meowsbits

		/*
			Step 1
			For the second scenario with the same validator set and in-order chain with
			validator 7 having just produced an in order block, then validator 8 goes offline.

			1A x
			2B  x
			3C   x
			4D    x
			5E     x
			6F      x
			7G       x
			8H

			1A x
			2B  x
			3C   x
			4D    x
			5E     x
			6F      x
			7G       x
			8H        -

			Here’s a revised one. 8 nodes, zero based. 0-6 all produce in-order blocks,
			then a netsplit. 0, 2, and 3 on the first fork and 1, 4, 6, 7 on the second fork, and 5 goes offline.
			7, 0, and 1 all missed an important in-turn block.

			1A x       x
			2B  x      y
			3C   x       x
			4D    x     x
			5E     x     y
			6F      x -
			7G       x
			8H          y

		*/
		lenSigners:   8,
		commonBlocks: []int{1, 2, 3, 4, 5, 6, 7, 0, 1, 2, 3, 4, 5, 6},
		forks: [][]int{
			{0, 3, 2},
			{1, 7, 4},
		},
		canonicalForkIndex: lowerHash,
		assertions: []func(t *testing.T, chain *core.BlockChain, forkHeads ...*types.Header){
			assertEqualTotalDifficulties,
			assertEqualNumbers,
		},
		cliqueConfig: cliqueConfigEIP3436,
	})
}

func testCliqueEIP3436(t *testing.T, testConfig cliqueEIP3436TestCase) {
	// Create the account pool and generate the initial set of signerAddressesSorted
	accountsPool := newTesterAccountPool()

	db := rawdb.NewMemoryDatabase()

	// Assemble a chain of headers from the cast votes
	chainConfig := *params.TestChainConfig

	chainConfig.Ethash = nil
	chainConfig.Clique = testConfig.cliqueConfig
	engine := New(chainConfig.Clique, db)
	engine.fakeDiff = false

	signerAddressesSorted := make([]common.Address, testConfig.lenSigners)
	for i := 0; i < testConfig.lenSigners; i++ {
		signerAddressesSorted[i] = accountsPool.address(fmt.Sprintf("%s", []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K"}[i]))
	}
	for j := 0; j < len(signerAddressesSorted); j++ {
		for k := j + 1; k < len(signerAddressesSorted); k++ {
			if bytes.Compare(signerAddressesSorted[j][:], signerAddressesSorted[k][:]) > 0 {
				signerAddressesSorted[j], signerAddressesSorted[k] = signerAddressesSorted[k], signerAddressesSorted[j]
			}
		}
	}

	// Pretty logging of the sorted signers list.
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
	chain, err := core.NewBlockChain(db, nil, &chainConfig, engine, vm.Config{}, nil, nil)
	if err != nil {
		t.Errorf("failed to create test chain: %v", err)
		return
	}

	// getNextBlockWithSigner generates a block given a parent block and the
	// signer index of the validator that should sign it.
	// It will use the Clique snapshot function to see if the signer is in turn or not,
	// and will assign difficulty appropriately.
	// After signing, it sanity-checks the Engine.Author method against the newly-signed
	// block value to make sure that signing is done properly.
	getNextBlockWithSigner := func(parentBlock *types.Block, signerIndex int) *types.Block {
		signerAddress := signerAddressesSorted[signerIndex]
		signerName := accountsPool.name(signerAddress)

		generatedBlocks, _ := core.GenerateChain(&chainConfig, parentBlock, engine, db, 1, nil)
		block := generatedBlocks[0]

		// Get the header and prepare it for signing
		header := block.Header()
		header.ParentHash = parentBlock.Hash()

		header.Time = parentBlock.Time() + cliqueEIP3436TestPeriod

		// See if our required signer is in or out of turn and assign difficulty respectively.
		difficulty := diffInTurn

		// If the snapshot reports this signer is out of turn, use out-of-turn difficulty.
		snap, err := engine.snapshot(chain, parentBlock.NumberU64(), parentBlock.Hash(), nil)
		if err != nil {
			t.Fatalf("snap err: %v", err)
		}
		inturn := snap.inturn(parentBlock.NumberU64()+1, signerAddress)
		if !inturn {
			difficulty = diffNoTurn
		}
		header.Difficulty = difficulty

		// Generate the signature, embed it into the header and the block.
		header.Extra = make([]byte, extraVanity+extraSeal) // allocate byte slice

		t.Logf("SIGNING: %d (%s) %v %s\n", signerIndex, signerName, inturn, signerAddress.Hex())

		// Sign the header with the associated validator's key.
		accountsPool.sign(header, signerName)

		// Double check to see what the Clique Engine thinks the signer of this block is.
		// It obviously should be the address we just used.
		author, err := engine.Author(header)
		if err != nil {
			t.Fatalf("author error: %v", err)
		}
		if wantSigner := accountsPool.address(signerName); author != wantSigner {
			t.Fatalf("header author != wanted signer: author: %s, signer: %s", author.Hex(), wantSigner.Hex())
		}

		out := block.WithSeal(header)
		t.Logf("BLOCK:   %x", out.Hash().Bytes()[:8])
		return out
	}

	// Build the common segment
	commonSegmentBlocks := []*types.Block{genesisBlock}
	for i := 0; i < len(testConfig.commonBlocks); i++ {

		signerIndex := testConfig.commonBlocks[i]
		parentBlock := commonSegmentBlocks[len(commonSegmentBlocks)-1]

		bl := getNextBlockWithSigner(parentBlock, signerIndex)

		if k, err := chain.InsertChain([]*types.Block{bl}); err != nil || k != 1 {
			t.Fatalf("failed to import block %d, count: %d, err: %v", i, k, err)
		}

		commonSegmentBlocks = append(commonSegmentBlocks, bl)
	}

	t.Logf("--- COMMON SEGMENT, td=%v", chain.GetTd(chain.CurrentHeader().Hash(), chain.CurrentHeader().Number.Uint64()))

	forkHeads := make([]*types.Header, len(testConfig.forks))
	forkTDs := make([]*big.Int, len(testConfig.forks))

	// Create and import blocks for all the scenario's forks.
	for scenarioForkIndex, forkBlockSigners := range testConfig.forks {

		forkBlocks := make([]*types.Block, len(commonSegmentBlocks))
		for i, b := range commonSegmentBlocks {
			bcopy := &types.Block{}
			*bcopy = *b
			forkBlocks[i] = bcopy
		}

		for si, signerInt := range forkBlockSigners {

			parent := forkBlocks[len(forkBlocks)-1]
			bl := getNextBlockWithSigner(parent, signerInt)

			if k, err := chain.InsertChain([]*types.Block{bl}); err != nil || k != 1 {
				t.Fatalf("failed to import block %d, count: %d, err: %v", si, k, err)
			} else {
				t.Logf("INSERTED block.n=%d TD=%v", bl.NumberU64(), chain.GetTd(bl.Hash(), bl.NumberU64()))
			}
			forkBlocks = append(forkBlocks, bl)

		} // End fork block imports.

		forkHeads[scenarioForkIndex] = forkBlocks[len(forkBlocks)-1].Header()
		forkTDs[scenarioForkIndex] = chain.GetTd(chain.CurrentHeader().Hash(), chain.CurrentHeader().Number.Uint64())

	} // End scenario fork imports.

	// Run arbitrary assertion tests, ie. make sure that we've created forks that meet
	// the expected scenario characteristics.
	for _, f := range testConfig.assertions {
		f(t, chain, forkHeads...)
	}

	// Finally, check that the current chain head matches the
	// head of the wanted fork index.
	forkHeadHashes := ""
	for i, fh := range forkHeads {
		forkHeadHashes += fmt.Sprintf("%d: %s td=%v\n", i, fh.Hash().Hex(), forkTDs[i])
	}
	if chain.CurrentHeader().Hash() != forkHeads[testConfig.canonicalForkIndex(t, forkHeads...)].Hash() {
		t.Errorf("wrong fork index head: got: %s\nFork heads:\n%s", chain.CurrentHeader().Hash().Hex(), forkHeadHashes)
	} else {
		t.Logf("CHAIN CURRENT HEAD: %x", chain.CurrentHeader().Hash().Bytes()[:8])
		t.Logf("Heads: \n%s", forkHeadHashes)
	}
}

// shorterFork defines logic returning the fork head the lower head number.
func shorterFork(t *testing.T, forkHeads ...*types.Header) int {
	// Prefer the shorter fork.
	minHeight := math.MaxBig63.Uint64()
	n := -1
	for i, head := range forkHeads {
		if h := head.Number.Uint64(); h < minHeight {
			n = i
			minHeight = h
		}
	}
	return n
}

// lowerHash defines logic returning the fork head having the lesser block hash.
func lowerHash(t *testing.T, forkHeads ...*types.Header) int {
	// Prefer the lowest hash.
	minHashV, overflow := uint256.FromBig(big.NewInt(0))
	if overflow {
		t.Fatalf("uint256 overflowed: 0")
	}
	n := -1
	for i, head := range forkHeads {
		hv, err := hashToUint256(head.Hash())
		if err != nil {
			t.Fatalf("uint256 err: %v, head.hex: %s", err, head.Hash().Hex())
		}
		if n == -1 || hv.Cmp(minHashV) < 0 {
			minHashV.Set(hv)
			n = i
		}
	}
	return n
}

// assertEqualTotalDifficulties fatals if the fork heads do not have equal total difficulties.
func assertEqualTotalDifficulties(t *testing.T, chain *core.BlockChain, forkHeads ...*types.Header) {
	d := new(big.Int)
	for i, head := range forkHeads {
		td := chain.GetTd(head.Hash(), head.Number.Uint64())
		if i == 0 {
			d.Set(td)
			continue
		}
		if d.Cmp(td) != 0 {
			t.Fatalf("want equal fork heads total difficulty")
		}
	}
}

// assertEqualNumbers fatals if the fork heads do not have equal numbers.
func assertEqualNumbers(t *testing.T, chain *core.BlockChain, forkHeads ...*types.Header) {
	n := new(big.Int)
	for i, head := range forkHeads {
		if i == 0 {
			n.Set(head.Number)
			continue
		}
		if n.Cmp(head.Number) != 0 {
			t.Fatalf("want equal fork head numbers")
		}
	}
}

func TestClique_Eip3436Rule4_hashToUint256(t *testing.T) {
	// "0x0f73844a381e7a0e002f128376a848d5f1193d07dc0570f1531d1c95f9e95a93"
	cases := []struct {
		hash common.Hash
		want *uint256.Int
	}{
		// This batch of cases is primarily concerned with checking
		// that hashes containing leading 0's don't return errors
		// (and that we're allowed to pass odd-length hex strings).
		{common.HexToHash("0xff73844a381e7a0e002f128376a848d5f1193d07dc0570f1531d1c95f9e95a93"), nil},
		{common.HexToHash("0x0f73844a381e7a0e002f128376a848d5f1193d07dc0570f1531d1c95f9e95a93"), nil},
		{common.HexToHash("0xf73844a381e7a0e002f128376a848d5f1193d07dc0570f1531d1c95f9e95a93"), nil},
	}
	for i, c := range cases {
		got, err := hashToUint256(c.hash)
		if err != nil {
			t.Errorf("case %d, error: %v (hash.hex=%s)", i, err, c.hash.Hex())
		}
		// Skip tests with want=nil
		if c.want != nil && ((got == nil && c.want != nil) || (got.Cmp(c.want) != 0)) {
			t.Errorf("case: %d, want: %v, got: %v", i, c.want, got)
		}
	}
}

//
// func TestClique_EIP3436_Scenario1(t *testing.T) {
//
// 	scenarios := []struct {
// 		// The number of signers (aka validators).
// 		// These addresses will be randomly generated on demand for each scenario
// 		// and sorted per status quo Clique spec.
// 		// They are referenced by their cardinality in this sorted list.
// 		lenSigners int
//
// 		// commonBlocks defines blocks by which validator should sign them.
// 		// Validators are referenced by their index in the sorted signers list,
// 		// which sorts by address.
// 		commonBlocks []int
//
// 		// forks defines blocks by which validator should sign them.
// 		// These forks will be generated and attempted to be imported into
// 		// the chain. We expect that import should always succeed for all blocks,
// 		// with this test measuring only if the expected fork head actually gets
// 		// canonical preference.
// 		forks [][]int
//
// 		// assertions are functions used to assert the expectations of
// 		// fork heads in chain context against the specification's requirements, eg.
// 		// make sure that the forks really do have equal total difficulty, or equal block numbers.
// 		// This way we know that the deciding condition for achieving canonical status
// 		// really is what we think it is (and not, eg. total difficulty).
//
// 		assertions []func(t *testing.T, chain *core.BlockChain, forkHeads ...*types.Header)
//
// 		// canonicalForkIndex takes variadic fork heads and tells us which
// 		// should get canonical preference (by index).
// 		canonicalForkIndex func(t *testing.T, forkHeads ...*types.Header) int
//
// 		cliqueConfig *ctypes.CliqueConfig
// 	}{
// 		{
// 			// SCENARIO-1
// 			// signers: A..H (8 count)
// 			//
// 			/*
// 				Step 1
// 				A fully in-order chain exists and validator 8 has just produced an in-turn block.
//
// 				1A x
// 				2B  x
// 				3C   x
// 				4D    x
// 				5E     x
// 				6F      x
// 				7G       x
// 				8H        x
//
// 				Step 2
// 				... and then validators 5, 7, and 8 go offline.
//
// 				1A x
// 				2B  x
// 				3C   x
// 				4D    x
// 				5E     x   -
// 				6F      x
// 				7G       x -
// 				8H        x-
//
// 				Step 3
// 				Two forks form, one with an in-order block from validator 1
// 				and then an out of order block from validator 3.
//
// 				The second fork forms from validators 2, 4, and 6 in order.
// 				Both have a net total difficulty of 3 more than the common ancestor.
//
// 				1A x        y
// 				2B  x       z
// 				3C   x       y
// 				4D    x      z
// 				5E     x   -
// 				6F      x     z
// 				7G       x -
// 				8H        x-
// 			*/
// 			lenSigners:   8,
// 			commonBlocks: []int{1, 2, 3, 4, 5, 6, 7, 0, 1, 2, 3, 4, 5, 6, 7},
// 			forks: [][]int{
// 				{1, 3, 5}, // 2, 4, 6
// 				{0, 2},    // 1, 3
// 			},
// 			canonicalForkIndex: shorterFork,
// 			assertions: []func(t *testing.T, chain *core.BlockChain, forkHeads ...*types.Header){
// 				assertEqualTotalDifficulties,
// 			},
// 			cliqueConfig: cliqueConfigEIP3436,
// 		},
// 		{
// 			// SCENARIO-1 without EIP-3436 (negative test)
// 			lenSigners:   8,
// 			commonBlocks: []int{1, 2, 3, 4, 5, 6, 7, 0, 1, 2, 3, 4, 5, 6, 7},
// 			forks: [][]int{
// 				{1, 3, 5}, // 2, 4, 6
// 				{0, 2},    // 1, 3
// 			},
// 			canonicalForkIndex: shorterFork,
// 			assertions: []func(t *testing.T, chain *core.BlockChain, forkHeads ...*types.Header){
// 				assertEqualTotalDifficulties,
// 			},
// 			cliqueConfig: cliqueConfigNoEIP3436,
// 		},
// 		{
// 			// SCENARIO-2
//
// 			/*
// 					Step 1
// 					For the second scenario with the same validator set and in-order chain with
// 					validator 7 having just produced an in order block, then validators 7 and 8 go offline.
//
// 					1A x
// 					2B  x
// 					3C   x
// 					4D    x
// 					5E     x
// 					6F      x
// 					7G       x
// 					8H
//
// 					1A x
// 					2B  x
// 					3C   x
// 					4D    x
// 					5E     x
// 					6F      x
// 					7G       x-
// 					8H        -
//
// 					Two forks form, 1,3,5 on one side and 2,4,6 on the other.
// 					Both forks become aware of the other fork after producing their third block.
// 					In this case both forks have equal total difficulty and equal length.
//
// 					1A x       x
// 					2B  x      y
// 					3C   x      x
// 					4D    x     y
// 					5E     x     x
// 					6F      x    y
// 					7G       x-
// 					8H        -
//
// 				FIXME(meowsbits): This scenario yields a "recently signed" error
// 				when attempting to import Signer 5 (really #6 b/c zero-indexing) into the
// 				second fork.
//
// 				On that fork, the sequence of signers is specified to be
// 				... 7, 0, 1, 2, 3, 4, 5, 6, 1, 3, 5
//
// 				(vs. the other fork)
// 				... 7, 0, 1, 2, 3, 4, 5, 6, 0, 2, 4
//
// 				The condition for "recently signed" is (from *Clique#verifySeal):
//
// 				// Signer is among recents, only fail if the current block doesn't shift it out
// 				if limit := uint64(len(snap.Signers)/2 + 1); seen > number-limit {
// 					return errRecentlySigned
// 				}
//
// 				Evaluated, this yields
//
// 				=> recently signed: limit=(8/2+1)=5 seen=13 number=17 number-limit=12
//
// 				RESOLUTION (tentative): Use forks with length 2 instead of 3.
//
// 			*/
// 			lenSigners:   8,
// 			commonBlocks: []int{1, 2, 3, 4, 5, 6, 7, 0, 1, 2, 3, 4, 5, 6},
// 			forks: [][]int{
// 				// These expected values have been truncated to 2 block each instead
// 				// of the EIP's described 3 blocks (per inline comments).
// 				{0, 2}, // 1, 3, 5
// 				{1, 3}, // 2, 4, 6
// 			},
// 			canonicalForkIndex: lowerHash,
// 			assertions: []func(t *testing.T, chain *core.BlockChain, forkHeads ...*types.Header){
// 				assertEqualTotalDifficulties,
// 				assertEqualNumbers,
// 			},
// 			cliqueConfig: cliqueConfigEIP3436,
// 		},
// 		{
// 			// SCENARIO-2 without EIP-3436 (negative test)
// 			lenSigners:   8,
// 			commonBlocks: []int{1, 2, 3, 4, 5, 6, 7, 0, 1, 2, 3, 4, 5, 6},
// 			forks: [][]int{
// 				// These expected values have been truncated to 2 block each instead
// 				// of the EIP's described 3 blocks (per inline comments).
// 				{0, 2}, // 1, 3, 5
// 				{1, 3}, // 2, 4, 6
// 			},
// 			canonicalForkIndex: func(t *testing.T, forkHeads ...*types.Header) int {
// 				// Prefer the first-seen fork.
// 				return 0
// 			},
// 			assertions: []func(t *testing.T, chain *core.BlockChain, forkHeads ...*types.Header){
// 				assertEqualTotalDifficulties,
// 				assertEqualNumbers,
// 			},
// 			cliqueConfig: cliqueConfigNoEIP3436,
// 		},
// 		{
// 			// SCENARIO-3
// 			// https://ethereum-magicians.org/t/eip-3436-expanded-clique-block-choice-rule/5809/3?u=meowsbits
//
// 			/*
// 				Step 1
// 				For the second scenario with the same validator set and in-order chain with
// 				validator 7 having just produced an in order block, then validator 8 goes offline.
//
// 				1A x
// 				2B  x
// 				3C   x
// 				4D    x
// 				5E     x
// 				6F      x
// 				7G       x
// 				8H
//
// 				1A x
// 				2B  x
// 				3C   x
// 				4D    x
// 				5E     x
// 				6F      x
// 				7G       x
// 				8H        -
//
// 				Here’s a revised one. 8 nodes, zero based. 0-6 all produce in-order blocks,
// 				then a netsplit. 0, 2, and 3 on the first fork and 1, 4, 6, 7 on the second fork, and 5 goes offline.
// 				7, 0, and 1 all missed an important in-turn block.
//
// 				1A x       x
// 				2B  x      y
// 				3C   x       x
// 				4D    x     x
// 				5E     x     y
// 				6F      x -
// 				7G       x
// 				8H          y
//
// 			*/
// 			lenSigners:   8,
// 			commonBlocks: []int{1, 2, 3, 4, 5, 6, 7, 0, 1, 2, 3, 4, 5, 6},
// 			forks: [][]int{
// 				{0, 3, 2},
// 				{1, 7, 4},
// 			},
// 			canonicalForkIndex: lowerHash,
// 			assertions: []func(t *testing.T, chain *core.BlockChain, forkHeads ...*types.Header){
// 				assertEqualTotalDifficulties,
// 				assertEqualNumbers,
// 			},
// 			cliqueConfig: cliqueConfigEIP3436,
// 		},
// 	}
//
// 	for ii, tt := range scenarios {
// 		t.Logf("************************ SCENARIO %d ************************", ii)
// 		// Create the account pool and generate the initial set of signerAddressesSorted
// 		accountsPool := newTesterAccountPool()
//
// 		db := rawdb.NewMemoryDatabase()
//
// 		// Assemble a chain of headers from the cast votes
// 		config := *params.TestChainConfig
//
// 		config.Ethash = nil
// 		config.Clique = tt.cliqueConfig
// 		engine := New(config.Clique, db)
// 		engine.fakeDiff = false
//
// 		signerAddressesSorted := make([]common.Address, tt.lenSigners)
// 		for i := 0; i < tt.lenSigners; i++ {
// 			signerAddressesSorted[i] = accountsPool.address(fmt.Sprintf("%s", []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K"}[i]))
// 		}
// 		for j := 0; j < len(signerAddressesSorted); j++ {
// 			for k := j + 1; k < len(signerAddressesSorted); k++ {
// 				if bytes.Compare(signerAddressesSorted[j][:], signerAddressesSorted[k][:]) > 0 {
// 					signerAddressesSorted[j], signerAddressesSorted[k] = signerAddressesSorted[k], signerAddressesSorted[j]
// 				}
// 			}
// 		}
//
// 		// Pretty logging of the sorted signers list.
// 		logSortedSigners := ""
// 		for j, s := range signerAddressesSorted {
// 			logSortedSigners += fmt.Sprintf("%d: (%s) %s\n", j, accountsPool.name(s), s.Hex())
// 		}
// 		t.Logf("SORTED SIGNERS:\n%s\n", logSortedSigners)
//
// 		// Create the genesis block with the initial set of signerAddressesSorted
// 		genesis := &genesisT.Genesis{
// 			ExtraData: make([]byte, extraVanity+common.AddressLength*len(signerAddressesSorted)+extraSeal),
// 		}
// 		for j, signer := range signerAddressesSorted {
// 			copy(genesis.ExtraData[extraVanity+j*common.AddressLength:], signer[:])
// 		}
//
// 		genesisBlock := core.MustCommitGenesis(db, genesis)
//
// 		// Create a pristine blockchain with the genesis injected
// 		chain, err := core.NewBlockChain(db, nil, &config, engine, vm.Config{}, nil, nil)
// 		if err != nil {
// 			t.Errorf("test %d: failed to create test chain: %v", ii, err)
// 			continue
// 		}
//
// 		// getNextBlockWithSigner generates a block given a parent block and the
// 		// signer index of the validator that should sign it.
// 		// It will use the Clique snapshot function to see if the signer is in turn or not,
// 		// and will assign difficulty appropriately.
// 		// After signing, it sanity-checks the Engine.Author method against the newly-signed
// 		// block value to make sure that signing is done properly.
// 		getNextBlockWithSigner := func(parentBlock *types.Block, signerIndex int) *types.Block {
// 			signerAddress := signerAddressesSorted[signerIndex]
// 			signerName := accountsPool.name(signerAddress)
//
// 			generatedBlocks, _ := core.GenerateChain(&config, parentBlock, engine, db, 1, nil)
// 			block := generatedBlocks[0]
//
// 			// Get the header and prepare it for signing
// 			header := block.Header()
// 			header.ParentHash = parentBlock.Hash()
//
// 			header.Time = parentBlock.Time() + cliqueEIP3436TestPeriod
//
// 			// See if our required signer is in or out of turn and assign difficulty respectively.
// 			difficulty := diffInTurn
//
// 			// If the snapshot reports this signer is out of turn, use out-of-turn difficulty.
// 			snap, err := engine.snapshot(chain, parentBlock.NumberU64(), parentBlock.Hash(), nil)
// 			if err != nil {
// 				t.Fatalf("snap err: %v", err)
// 			}
// 			inturn := snap.inturn(parentBlock.NumberU64()+1, signerAddress)
// 			if !inturn {
// 				difficulty = diffNoTurn
// 			}
// 			header.Difficulty = difficulty
//
// 			// Generate the signature, embed it into the header and the block.
// 			header.Extra = make([]byte, extraVanity+extraSeal) // allocate byte slice
//
// 			t.Logf("SIGNING: %d (%s) %v %s\n", signerIndex, signerName, inturn, signerAddress.Hex())
//
// 			// Sign the header with the associated validator's key.
// 			accountsPool.sign(header, signerName)
//
// 			// Double check to see what the Clique Engine thinks the signer of this block is.
// 			// It obviously should be the address we just used.
// 			author, err := engine.Author(header)
// 			if err != nil {
// 				t.Fatalf("author error: %v", err)
// 			}
// 			if wantSigner := accountsPool.address(signerName); author != wantSigner {
// 				t.Fatalf("header author != wanted signer: author: %s, signer: %s", author.Hex(), wantSigner.Hex())
// 			}
//
// 			out := block.WithSeal(header)
// 			t.Logf("BLOCK:   %x", out.Hash().Bytes()[:8])
// 			return out
// 		}
//
// 		// Build the common segment
// 		commonSegmentBlocks := []*types.Block{genesisBlock}
// 		for i := 0; i < len(tt.commonBlocks); i++ {
//
// 			signerIndex := tt.commonBlocks[i]
// 			parentBlock := commonSegmentBlocks[len(commonSegmentBlocks)-1]
//
// 			bl := getNextBlockWithSigner(parentBlock, signerIndex)
//
// 			if k, err := chain.InsertChain([]*types.Block{bl}); err != nil || k != 1 {
// 				t.Fatalf("case: %d, failed to import block %d, count: %d, err: %v", ii, i, k, err)
// 			}
//
// 			commonSegmentBlocks = append(commonSegmentBlocks, bl)
// 		}
//
// 		t.Logf("--- COMMON SEGMENT, td=%v", chain.GetTd(chain.CurrentHeader().Hash(), chain.CurrentHeader().Number.Uint64()))
//
// 		forkHeads := make([]*types.Header, len(tt.forks))
// 		forkTDs := make([]*big.Int, len(tt.forks))
//
// 		// Create and import blocks for all the scenario's forks.
// 		for scenarioForkIndex, forkBlockSigners := range tt.forks {
//
// 			forkBlocks := make([]*types.Block, len(commonSegmentBlocks))
// 			for i, b := range commonSegmentBlocks {
// 				bcopy := &types.Block{}
// 				*bcopy = *b
// 				forkBlocks[i] = bcopy
// 			}
//
// 			for si, signerInt := range forkBlockSigners {
//
// 				parent := forkBlocks[len(forkBlocks)-1]
// 				bl := getNextBlockWithSigner(parent, signerInt)
//
// 				if k, err := chain.InsertChain([]*types.Block{bl}); err != nil || k != 1 {
// 					t.Fatalf("case: %d, failed to import block %d, count: %d, err: %v", ii, si, k, err)
// 				} else {
// 					t.Logf("INSERTED block.n=%d TD=%v", bl.NumberU64(), chain.GetTd(bl.Hash(), bl.NumberU64()))
// 				}
// 				forkBlocks = append(forkBlocks, bl)
//
// 			} // End fork block imports.
//
// 			forkHeads[scenarioForkIndex] = forkBlocks[len(forkBlocks)-1].Header()
// 			forkTDs[scenarioForkIndex] = chain.GetTd(chain.CurrentHeader().Hash(), chain.CurrentHeader().Number.Uint64())
//
// 		} // End scenario fork imports.
//
// 		// Run arbitrary assertion tests, ie. make sure that we've created forks that meet
// 		// the expected scenario characteristics.
// 		for _, f := range tt.assertions {
// 			f(t, chain, forkHeads...)
// 		}
//
// 		// Finally, check that the current chain head matches the
// 		// head of the wanted fork index.
// 		forkHeadHashes := ""
// 		for i, fh := range forkHeads {
// 			forkHeadHashes += fmt.Sprintf("%d: %s td=%v\n", i, fh.Hash().Hex(), forkTDs[i])
// 		}
// 		if chain.CurrentHeader().Hash() != forkHeads[tt.canonicalForkIndex(t, forkHeads...)].Hash() {
// 			t.Errorf("wrong fork index head: got: %s\nFork heads:\n%s", chain.CurrentHeader().Hash().Hex(), forkHeadHashes)
// 		} else {
// 			t.Logf("CHAIN CURRENT HEAD: %x", chain.CurrentHeader().Hash().Bytes()[:8])
// 			t.Logf("Heads: \n%s", forkHeadHashes)
// 		}
// 	}
// }
