// Copyright 2014 The go-ethereum Authors
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

package core

import (
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/cmd/ancient-store-mem/lib"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/ethash"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/params/types/genesisT"
	"github.com/ethereum/go-ethereum/params/vars"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/ethereum/go-ethereum/trie"
)

// Tests in this file duplicate select tests from blockchain_test.go,
// replacing the built in file-based ancient db with a remote one over IPC.

var (
	// testRPCFreezerURL defines an optional environment variable that, if set,
	// will be used in lieu of a default ephemeral remote freezer store.
	testRPCFreezerURL = os.Getenv("GETH_ANCIENT_RPC")
)

// testRPCRemoteFreezer provides a configuration option to use an external
// remote freezer server, or to default (with no configured flags) to a built in
// ephemeral in-memory server over a temporary unix socket.
// If an external URL is configured the return value for 'server' will be nil.
func testRPCRemoteFreezer(t *testing.T) (rpcFreezerEndpoint string, server *rpc.Server, ancientDB ethdb.Database) {
	if testRPCFreezerURL == "" {
		// If an external freezer server is not provided, spin up an ephemeral
		// freezer over IPC.
		frdir, err := ioutil.TempDir("", "")
		if err != nil {
			t.Fatalf("failed to create temp freezer dir: %v", err)
		}
		rpcFreezerEndpoint = filepath.Join(frdir, "test.ipc")

		listener, server, err := rpc.StartIPCEndpoint(rpcFreezerEndpoint, nil)
		if err != nil {
			t.Fatal(err)
		}
		mock := lib.NewMemFreezerRemoteServerAPI()
		err = server.RegisterName("freezer", mock)
		if err != nil {
			t.Fatal(err)
		}
		go func() {
			server.ServeListener(listener)
		}()

	} else {
		rpcFreezerEndpoint = testRPCFreezerURL
		t.Log("Using external freezer:", rpcFreezerEndpoint)
	}

	ancientDb, err := rawdb.NewDatabaseWithFreezerRemote(rawdb.NewMemoryDatabase(), rpcFreezerEndpoint)
	if err != nil {
		t.Fatalf("failed to create temp freezer db: %v", err)
	}

	return rpcFreezerEndpoint, server, ancientDb
}

// Tests that fast importing a block chain produces the same chain data as the
// classical full block processing.
// Extra steps have been added to this test compared to its sister TestFastVsFullChains
// that ensure that the ancient store methods are completely exercised; eg. a rollback step is
// used to call TruncateAncients.
func TestFastVsFullChains_RemoteFreezer(t *testing.T) {
	// Configure and generate a sample block chain
	var (
		gendb   = rawdb.NewMemoryDatabase()
		key, _  = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		address = crypto.PubkeyToAddress(key.PublicKey)
		funds   = big.NewInt(1000000000)
		gspec   = &genesisT.Genesis{
			Config: params.TestChainConfig,
			Alloc:  genesisT.GenesisAlloc{address: {Balance: funds}},
		}
		genesis = MustCommitGenesis(gendb, gspec)
		signer  = types.NewEIP155Signer(gspec.Config.GetChainID())
	)
	blocks, receipts := GenerateChain(gspec.Config, genesis, ethash.NewFaker(), gendb, 1024, func(i int, block *BlockGen) {
		block.SetCoinbase(common.Address{0x00})

		// If the block number is multiple of 3, send a few bonus transactions to the miner
		if i%3 == 2 {
			for j := 0; j < i%4+1; j++ {
				tx, err := types.SignTx(types.NewTransaction(block.TxNonce(address), common.Address{0x00}, big.NewInt(1000), vars.TxGas, nil, nil), signer, key)
				if err != nil {
					panic(err)
				}
				block.AddTx(tx)
			}
		}
		// If the block number is a multiple of 5, add a few bonus uncles to the block
		if i%5 == 5 {
			block.AddUncle(&types.Header{ParentHash: block.PrevBlock(i - 1).Hash(), Number: big.NewInt(int64(i - 1))})
		}
	})
	// Import the chain as an archive node for the comparison baseline
	archiveDb := rawdb.NewMemoryDatabase()
	MustCommitGenesis(archiveDb, gspec)
	archive, _ := NewBlockChain(archiveDb, nil, gspec.Config, ethash.NewFaker(), vm.Config{}, nil, nil)
	defer archive.Stop()

	if n, err := archive.InsertChain(blocks); err != nil {
		t.Fatalf("failed to process block %d: %v", n, err)
	}

	// Fast import the chain as a non-archive node to test
	fastDb := rawdb.NewMemoryDatabase()
	MustCommitGenesis(fastDb, gspec)
	fast, _ := NewBlockChain(fastDb, nil, gspec.Config, ethash.NewFaker(), vm.Config{}, nil, nil)
	defer fast.Stop()

	headers := make([]*types.Header, len(blocks))
	for i, block := range blocks {
		headers[i] = block.Header()
	}
	if n, err := fast.InsertHeaderChain(headers, 1); err != nil {
		t.Fatalf("failed to insert header %d: %v", n, err)
	}
	if n, err := fast.InsertReceiptChain(blocks, receipts, 0); err != nil {
		t.Fatalf("failed to insert receipt %d: %v", n, err)
	}

	// Freezer style fast import the chain.
	freezerRPCEndpoint, server, ancientDb := testRPCRemoteFreezer(t)
	if n, err := ancientDb.Ancients(); err != nil {
		t.Fatalf("ancients: %v", err)
	} else if n != 0 {
		t.Logf("truncating pre-existing ancients from: %d (truncating to 0)", n)
		err = ancientDb.TruncateAncients(0)
		if err != nil {
			t.Fatalf("truncate ancients: %v", err)
		}
	}
	if server != nil {
		defer os.RemoveAll(filepath.Dir(freezerRPCEndpoint))
		defer server.Stop()
	}
	defer ancientDb.Close() // Cause the Close method to be called.
	defer func() {
		// A deferred truncation to 0 will allow a single freezer instance to
		// handle multiple tests in serial.
		if err := ancientDb.TruncateAncients(0); err != nil {
			t.Fatalf("deferred truncate ancients error: %v", err)
		}
	}()

	MustCommitGenesis(ancientDb, gspec)
	ancient, _ := NewBlockChain(ancientDb, nil, gspec.Config, ethash.NewFaker(), vm.Config{}, nil, nil)
	defer ancient.Stop()

	ancientLimit := uint64(len(blocks) / 2)

	if n, err := ancient.InsertHeaderChain(headers, 1); err != nil {
		t.Fatalf("failed to insert header %d: %v", n, err)
	}
	if n, err := ancient.InsertReceiptChain(blocks, receipts, ancientLimit); err != nil {
		t.Fatalf("failed to insert receipt %d: %v", n, err)
	}

	// Test a rollback, causing the ancient store to use the TruncateAncient method.
	pinch := len(blocks) / 4
	rollbackHeaders := []common.Hash{}
	for _, v := range headers[pinch:] {
		rollbackHeaders = append(rollbackHeaders, v.Hash())
	}
	ancient.SetHead(headers[pinch].Number.Uint64())

	// Reinsert the rolled-back headers and receipts.
	if n, err := ancient.InsertHeaderChain(headers[pinch:], 1); err != nil {
		t.Fatalf("failed to insert header %d: %v", n, err)
	}
	if n, err := ancient.InsertReceiptChain(blocks[pinch:], receipts, ancientLimit); err != nil {
		t.Fatalf("failed to insert receipt %d: %v", n, err)
	}

	// Explicitly call the HasAncient method.
	// This method doesn't appear to be used in normal geth operation, and I'm not sure
	// why it exists, but we want to make sure that all API-defined methods get used.
	for _, b := range blocks {
		want := b.NumberU64() <= ancientLimit
		if ok, err := ancientDb.HasAncient(rawdb.FreezerRemoteHeaderTable, b.NumberU64()); err != nil || ok != want {
			t.Fatalf("ancientdb !HasAncient #%d: error=%v HasAncient=%v want=%v", b.NumberU64(), err, ok, want)
		}
	}

	// Iterate over all chain data components, and cross reference
	for i := 0; i < len(blocks); i++ {
		num, hash := blocks[i].NumberU64(), blocks[i].Hash()

		if ftd, atd := fast.GetTdByHash(hash), archive.GetTdByHash(hash); ftd.Cmp(atd) != 0 {
			t.Errorf("block #%d [%x]: td mismatch: fastdb %v, archivedb %v", num, hash, ftd, atd)
		}
		if antd, artd := ancient.GetTdByHash(hash), archive.GetTdByHash(hash); antd.Cmp(artd) != 0 {
			t.Errorf("block #%d [%x]: td mismatch: ancientdb %v, archivedb %v", num, hash, antd, artd)
		}
		if fheader, aheader := fast.GetHeaderByHash(hash), archive.GetHeaderByHash(hash); fheader.Hash() != aheader.Hash() {
			t.Errorf("block #%d [%x]: header mismatch: fastdb %v, archivedb %v", num, hash, fheader, aheader)
		}
		if anheader, arheader := ancient.GetHeaderByHash(hash), archive.GetHeaderByHash(hash); anheader.Hash() != arheader.Hash() {
			t.Errorf("block #%d [%x]: header mismatch: ancientdb %v, archivedb %v", num, hash, anheader, arheader)
		}
		if fblock, arblock, anblock := fast.GetBlockByHash(hash), archive.GetBlockByHash(hash), ancient.GetBlockByHash(hash); fblock.Hash() != arblock.Hash() || anblock.Hash() != arblock.Hash() {
			t.Errorf("block #%d [%x]: block mismatch: fastdb %v, ancientdb %v, archivedb %v", num, hash, fblock, anblock, arblock)
		} else if types.DeriveSha(fblock.Transactions(), new(trie.Trie)) != types.DeriveSha(arblock.Transactions(), new(trie.Trie)) || types.DeriveSha(anblock.Transactions(), new(trie.Trie)) != types.DeriveSha(arblock.Transactions(), new(trie.Trie)) {
			t.Errorf("block #%d [%x]: transactions mismatch: fastdb %v, ancientdb %v, archivedb %v", num, hash, fblock.Transactions(), anblock.Transactions(), arblock.Transactions())
		} else if types.CalcUncleHash(fblock.Uncles()) != types.CalcUncleHash(arblock.Uncles()) || types.CalcUncleHash(anblock.Uncles()) != types.CalcUncleHash(arblock.Uncles()) {
			t.Errorf("block #%d [%x]: uncles mismatch: fastdb %v, ancientdb %v, archivedb %v", num, hash, fblock.Uncles(), anblock, arblock.Uncles())
		}
		if freceipts, anreceipts, areceipts := rawdb.ReadReceipts(fastDb, hash, *rawdb.ReadHeaderNumber(fastDb, hash), fast.Config()), rawdb.ReadReceipts(ancientDb, hash, *rawdb.ReadHeaderNumber(ancientDb, hash), fast.Config()), rawdb.ReadReceipts(archiveDb, hash, *rawdb.ReadHeaderNumber(archiveDb, hash), fast.Config()); types.DeriveSha(freceipts, new(trie.Trie)) != types.DeriveSha(areceipts, new(trie.Trie)) {
			t.Errorf("block #%d [%x]: receipts mismatch: fastdb %v, ancientdb %v, archivedb %v", num, hash, freceipts, anreceipts, areceipts)
		}
	}
	// Check that the canonical chains are the same between the databases
	for i := 0; i < len(blocks)+1; i++ {
		if fhash, ahash := rawdb.ReadCanonicalHash(fastDb, uint64(i)), rawdb.ReadCanonicalHash(archiveDb, uint64(i)); fhash != ahash {
			t.Errorf("block #%d: canonical hash mismatch: fastdb %v, archivedb %v", i, fhash, ahash)
		}
		if anhash, arhash := rawdb.ReadCanonicalHash(ancientDb, uint64(i)), rawdb.ReadCanonicalHash(archiveDb, uint64(i)); anhash != arhash {
			t.Errorf("block #%d: canonical hash mismatch: ancientdb %v, archivedb %v", i, anhash, arhash)
		}
	}
}

func TestBlockchainRecovery_RemoteFreezer(t *testing.T) {
	// Configure and generate a sample block chain
	var (
		gendb   = rawdb.NewMemoryDatabase()
		key, _  = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		address = crypto.PubkeyToAddress(key.PublicKey)
		funds   = big.NewInt(1000000000)
		gspec   = &genesisT.Genesis{Config: params.TestChainConfig, Alloc: genesisT.GenesisAlloc{address: {Balance: funds}}}
		genesis = MustCommitGenesis(gendb, gspec)
	)
	height := uint64(1024)
	blocks, receipts := GenerateChain(gspec.Config, genesis, ethash.NewFaker(), gendb, int(height), nil)

	// Import the chain as a ancient-first node and ensure all pointers are updated
	// Freezer style fast import the chain.
	freezerRPCEndpoint, server, ancientDb := testRPCRemoteFreezer(t)
	if n, err := ancientDb.Ancients(); err != nil {
		t.Fatalf("ancients: %v", err)
	} else if n != 0 {
		t.Logf("truncating pre-existing ancients from: %d (truncating to 0)", n)
		err = ancientDb.TruncateAncients(0)
		if err != nil {
			t.Fatalf("truncate ancients: %v", err)
		}
	}
	if server != nil {
		defer os.RemoveAll(filepath.Dir(freezerRPCEndpoint))
		defer server.Stop()
	}
	defer ancientDb.Close() // Cause the Close method to be called.
	defer func() {
		// A deferred truncation to 0 will allow a single freezer instance to
		// handle multiple tests in serial.
		if err := ancientDb.TruncateAncients(0); err != nil {
			t.Fatalf("deferred truncate ancients error: %v", err)
		}
	}()

	MustCommitGenesis(ancientDb, gspec)
	ancient, _ := NewBlockChain(ancientDb, nil, gspec.Config, ethash.NewFaker(), vm.Config{}, nil, nil)

	headers := make([]*types.Header, len(blocks))
	for i, block := range blocks {
		headers[i] = block.Header()
	}
	if n, err := ancient.InsertHeaderChain(headers, 1); err != nil {
		t.Fatalf("failed to insert header %d: %v", n, err)
	}
	if n, err := ancient.InsertReceiptChain(blocks, receipts, uint64(3*len(blocks)/4)); err != nil {
		t.Fatalf("failed to insert receipt %d: %v", n, err)
	}
	ancient.Stop()

	// Destroy head fast block manually
	midBlock := blocks[len(blocks)/2]
	rawdb.WriteHeadFastBlockHash(ancientDb, midBlock.Hash())

	// Reopen broken blockchain again
	ancient, _ = NewBlockChain(ancientDb, nil, gspec.Config, ethash.NewFaker(), vm.Config{}, nil, nil)
	defer ancient.Stop()
	if num := ancient.CurrentBlock().NumberU64(); num != 0 {
		t.Errorf("head block mismatch: have #%v, want #%v", num, 0)
	}
	if num := ancient.CurrentFastBlock().NumberU64(); num != midBlock.NumberU64() {
		t.Errorf("head fast-block mismatch: have #%v, want #%v", num, midBlock.NumberU64())
	}
	if num := ancient.CurrentHeader().Number.Uint64(); num != midBlock.NumberU64() {
		t.Errorf("head header mismatch: have #%v, want #%v", num, midBlock.NumberU64())
	}
}

func TestIncompleteAncientReceiptChainInsertion_RemoteFreezer(t *testing.T) {
	// Configure and generate a sample block chain
	var (
		gendb   = rawdb.NewMemoryDatabase()
		key, _  = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		address = crypto.PubkeyToAddress(key.PublicKey)
		funds   = big.NewInt(1000000000)
		gspec   = &genesisT.Genesis{Config: params.TestChainConfig, Alloc: genesisT.GenesisAlloc{address: {Balance: funds}}}
		genesis = MustCommitGenesis(gendb, gspec)
	)
	height := uint64(1024)
	blocks, receipts := GenerateChain(gspec.Config, genesis, ethash.NewFaker(), gendb, int(height), nil)

	// Import the chain as a ancient-first node and ensure all pointers are updated
	freezerRPCEndpoint, server, ancientDb := testRPCRemoteFreezer(t)
	if n, err := ancientDb.Ancients(); err != nil {
		t.Fatalf("ancients: %v", err)
	} else if n != 0 {
		t.Logf("truncating pre-existing ancients from: %d (truncating to 0)", n)
		err = ancientDb.TruncateAncients(0)
		if err != nil {
			t.Fatalf("truncate ancients: %v", err)
		}
	}
	if server != nil {
		defer os.RemoveAll(filepath.Dir(freezerRPCEndpoint))
		defer server.Stop()
	}
	defer ancientDb.Close() // Cause the Close method to be called.
	defer func() {
		// A deferred truncation to 0 will allow a single freezer instance to
		// handle multiple tests in serial.
		if err := ancientDb.TruncateAncients(0); err != nil {
			t.Fatalf("deferred truncate ancients error: %v", err)
		}
	}()
	MustCommitGenesis(ancientDb, gspec)
	ancient, _ := NewBlockChain(ancientDb, nil, gspec.Config, ethash.NewFaker(), vm.Config{}, nil, nil)
	defer ancient.Stop()

	headers := make([]*types.Header, len(blocks))
	for i, block := range blocks {
		headers[i] = block.Header()
	}
	if n, err := ancient.InsertHeaderChain(headers, 1); err != nil {
		t.Fatalf("failed to insert header %d: %v", n, err)
	}
	// Abort ancient receipt chain insertion deliberately
	ancient.terminateInsert = func(hash common.Hash, number uint64) bool {
		return number == blocks[len(blocks)/2].NumberU64()
	}
	previousFastBlock := ancient.CurrentFastBlock()
	if n, err := ancient.InsertReceiptChain(blocks, receipts, uint64(3*len(blocks)/4)); err == nil {
		t.Fatalf("failed to insert receipt %d: %v", n, err)
	}
	if ancient.CurrentFastBlock().NumberU64() != previousFastBlock.NumberU64() {
		t.Fatalf("failed to rollback ancient data, want %d, have %d", previousFastBlock.NumberU64(), ancient.CurrentFastBlock().NumberU64())
	}
	if frozen, err := ancient.db.Ancients(); err != nil || frozen != 1 {
		t.Fatalf("failed to truncate ancient data")
	}
	ancient.terminateInsert = nil
	if n, err := ancient.InsertReceiptChain(blocks, receipts, uint64(3*len(blocks)/4)); err != nil {
		t.Fatalf("failed to insert receipt %d: %v", n, err)
	}
	if ancient.CurrentFastBlock().NumberU64() != blocks[len(blocks)-1].NumberU64() {
		t.Fatalf("failed to insert ancient recept chain after rollback")
	}
}

func TestTransactionIndices_RemoteFreezer(t *testing.T) {
	// Configure and generate a sample block chain
	var (
		gendb   = rawdb.NewMemoryDatabase()
		key, _  = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		address = crypto.PubkeyToAddress(key.PublicKey)
		funds   = big.NewInt(1000000000)
		gspec   = &genesisT.Genesis{Config: params.TestChainConfig, Alloc: genesisT.GenesisAlloc{address: {Balance: funds}}}
		genesis = MustCommitGenesis(gendb, gspec)
		signer  = types.NewEIP155Signer(gspec.Config.GetChainID())
	)
	height := uint64(128)
	blocks, receipts := GenerateChain(gspec.Config, genesis, ethash.NewFaker(), gendb, int(height), func(i int, block *BlockGen) {
		tx, err := types.SignTx(types.NewTransaction(block.TxNonce(address), common.Address{0x00}, big.NewInt(1000), vars.TxGas, nil, nil), signer, key)
		if err != nil {
			panic(err)
		}
		block.AddTx(tx)
	})
	blocks2, _ := GenerateChain(gspec.Config, blocks[len(blocks)-1], ethash.NewFaker(), gendb, 10, nil)

	check := func(tail *uint64, chain *BlockChain) {
		stored := rawdb.ReadTxIndexTail(chain.db)
		if tail == nil && stored != nil {
			t.Fatalf("Oldest indexded block mismatch, want nil, have %d", *stored)
		}
		if tail != nil && stored == nil {
			t.Fatalf("Oldest indexed block mismatch, want %d, have nil", tail)
		}
		if tail != nil && stored != nil && *stored != *tail {
			t.Fatalf("Oldest indexded block mismatch, want %d, have %d", *tail, *stored)
		}
		if tail != nil {
			for i := *tail; i <= chain.CurrentBlock().NumberU64(); i++ {
				block := rawdb.ReadBlock(chain.db, rawdb.ReadCanonicalHash(chain.db, i), i)
				if block.Transactions().Len() == 0 {
					continue
				}
				for _, tx := range block.Transactions() {
					if index := rawdb.ReadTxLookupEntry(chain.db, tx.Hash()); index == nil {
						t.Fatalf("Miss transaction indice, number %d hash %s", i, tx.Hash().Hex())
					}
				}
			}
			for i := uint64(0); i < *tail; i++ {
				block := rawdb.ReadBlock(chain.db, rawdb.ReadCanonicalHash(chain.db, i), i)
				if block.Transactions().Len() == 0 {
					continue
				}
				for _, tx := range block.Transactions() {
					if index := rawdb.ReadTxLookupEntry(chain.db, tx.Hash()); index != nil {
						t.Fatalf("Transaction indice should be deleted, number %d hash %s", i, tx.Hash().Hex())
					}
				}
			}
		}
	}
	freezerRPCEndpoint, server, ancientDb := testRPCRemoteFreezer(t)
	if n, err := ancientDb.Ancients(); err != nil {
		t.Fatalf("ancients: %v", err)
	} else if n != 0 {
		t.Logf("truncating pre-existing ancients from: %d (truncating to 0)", n)
		err = ancientDb.TruncateAncients(0)
		if err != nil {
			t.Fatalf("truncate ancients: %v", err)
		}
	}
	if server != nil {
		defer os.RemoveAll(filepath.Dir(freezerRPCEndpoint))
		defer server.Stop()
	}
	defer ancientDb.Close() // Cause the Close method to be called.
	defer func() {
		// A deferred truncation to 0 will allow a single freezer instance to
		// handle multiple tests in serial.
		if err := ancientDb.TruncateAncients(0); err != nil {
			t.Fatalf("deferred truncate ancients error: %v", err)
		}
	}()
	MustCommitGenesis(ancientDb, gspec)

	// Import all blocks into ancient db
	l := uint64(0)
	chain, err := NewBlockChain(ancientDb, nil, params.TestChainConfig, ethash.NewFaker(), vm.Config{}, nil, &l)
	if err != nil {
		t.Fatalf("failed to create tester chain: %v", err)
	}
	headers := make([]*types.Header, len(blocks))
	for i, block := range blocks {
		headers[i] = block.Header()
	}
	if n, err := chain.InsertHeaderChain(headers, 0); err != nil {
		t.Fatalf("failed to insert header %d: %v", n, err)
	}
	if n, err := chain.InsertReceiptChain(blocks, receipts, 128); err != nil {
		t.Fatalf("block %d: failed to insert into chain: %v", n, err)
	}
	chain.Stop()
	ancientDb.Close()

	// Init block chain with external ancients, check all needed indices has been indexed.
	limit := []uint64{0, 32, 64, 128}
	for _, l := range limit {
		ancientDb, err := rawdb.NewDatabaseWithFreezerRemote(rawdb.NewMemoryDatabase(), freezerRPCEndpoint)
		if err != nil {
			t.Fatalf("failed to create temp freezer db: %v", err)
		}
		MustCommitGenesis(ancientDb, gspec)
		chain, err = NewBlockChain(ancientDb, nil, params.TestChainConfig, ethash.NewFaker(), vm.Config{}, nil, &l)
		if err != nil {
			t.Fatalf("failed to create tester chain: %v", err)
		}
		time.Sleep(150 * time.Millisecond) // Wait for indices initialisation
		var tail uint64
		if l != 0 {
			tail = uint64(128) - l + 1
		}
		check(&tail, chain)
		chain.Stop()
		ancientDb.Close()
	}

	// Reconstruct a block chain which only reserves HEAD-64 tx indices
	ancientDb, err = rawdb.NewDatabaseWithFreezerRemote(rawdb.NewMemoryDatabase(), freezerRPCEndpoint)
	if err != nil {
		t.Fatalf("failed to create temp freezer db: %v", err)
	}
	MustCommitGenesis(ancientDb, gspec)

	limit = []uint64{0, 64 /* drop stale */, 32 /* shorten history */, 64 /* extend history */, 0 /* restore all */}
	tails := []uint64{0, 67 /* 130 - 64 + 1 */, 100 /* 131 - 32 + 1 */, 69 /* 132 - 64 + 1 */, 0}
	for i, l := range limit {
		chain, err = NewBlockChain(ancientDb, nil, params.TestChainConfig, ethash.NewFaker(), vm.Config{}, nil, &l)
		if err != nil {
			t.Fatalf("failed to create tester chain: %v", err)
		}
		chain.InsertChain(blocks2[i : i+1]) // Feed chain a higher block to trigger indices updater.
		time.Sleep(150 * time.Millisecond)  // Wait for indices initialisation
		check(&tails[i], chain)
		chain.Stop()
	}
}

func TestSkipStaleTxIndicesInFastSync_RemoteFreezer(t *testing.T) {
	// Configure and generate a sample block chain
	var (
		gendb   = rawdb.NewMemoryDatabase()
		key, _  = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		address = crypto.PubkeyToAddress(key.PublicKey)
		funds   = big.NewInt(1000000000)
		gspec   = &genesisT.Genesis{Config: params.TestChainConfig, Alloc: genesisT.GenesisAlloc{address: {Balance: funds}}}
		genesis = MustCommitGenesis(gendb, gspec)
		signer  = types.NewEIP155Signer(gspec.Config.GetChainID())
	)
	height := uint64(128)
	blocks, receipts := GenerateChain(gspec.Config, genesis, ethash.NewFaker(), gendb, int(height), func(i int, block *BlockGen) {
		tx, err := types.SignTx(types.NewTransaction(block.TxNonce(address), common.Address{0x00}, big.NewInt(1000), vars.TxGas, nil, nil), signer, key)
		if err != nil {
			panic(err)
		}
		block.AddTx(tx)
	})

	check := func(tail *uint64, chain *BlockChain) {
		stored := rawdb.ReadTxIndexTail(chain.db)
		if tail == nil && stored != nil {
			t.Fatalf("Oldest indexded block mismatch, want nil, have %d", *stored)
		}
		if tail != nil && *stored != *tail {
			t.Fatalf("Oldest indexded block mismatch, want %d, have %d", *tail, *stored)
		}
		if tail != nil {
			for i := *tail; i <= chain.CurrentBlock().NumberU64(); i++ {
				block := rawdb.ReadBlock(chain.db, rawdb.ReadCanonicalHash(chain.db, i), i)
				if block.Transactions().Len() == 0 {
					continue
				}
				for _, tx := range block.Transactions() {
					if index := rawdb.ReadTxLookupEntry(chain.db, tx.Hash()); index == nil {
						t.Fatalf("Miss transaction indice, number %d hash %s", i, tx.Hash().Hex())
					}
				}
			}
			for i := uint64(0); i < *tail; i++ {
				block := rawdb.ReadBlock(chain.db, rawdb.ReadCanonicalHash(chain.db, i), i)
				if block.Transactions().Len() == 0 {
					continue
				}
				for _, tx := range block.Transactions() {
					if index := rawdb.ReadTxLookupEntry(chain.db, tx.Hash()); index != nil {
						t.Fatalf("Transaction indice should be deleted, number %d hash %s", i, tx.Hash().Hex())
					}
				}
			}
		}
	}

	freezerRPCEndpoint, server, ancientDb := testRPCRemoteFreezer(t)
	if n, err := ancientDb.Ancients(); err != nil {
		t.Fatalf("ancients: %v", err)
	} else if n != 0 {
		t.Logf("truncating pre-existing ancients from: %d (truncating to 0)", n)
		err = ancientDb.TruncateAncients(0)
		if err != nil {
			t.Fatalf("truncate ancients: %v", err)
		}
	}
	if server != nil {
		defer os.RemoveAll(filepath.Dir(freezerRPCEndpoint))
		defer server.Stop()
	}
	defer ancientDb.Close() // Cause the Close method to be called.
	defer func() {
		// A deferred truncation to 0 will allow a single freezer instance to
		// handle multiple tests in serial.
		if err := ancientDb.TruncateAncients(0); err != nil {
			t.Fatalf("deferred truncate ancients error: %v", err)
		}
	}()
	MustCommitGenesis(ancientDb, gspec)

	// Import all blocks into ancient db, only HEAD-32 indices are kept.
	l := uint64(32)
	chain, err := NewBlockChain(ancientDb, nil, params.TestChainConfig, ethash.NewFaker(), vm.Config{}, nil, &l)
	if err != nil {
		t.Fatalf("failed to create tester chain: %v", err)
	}
	headers := make([]*types.Header, len(blocks))
	for i, block := range blocks {
		headers[i] = block.Header()
	}
	if n, err := chain.InsertHeaderChain(headers, 0); err != nil {
		t.Fatalf("failed to insert header %d: %v", n, err)
	}
	// The indices before ancient-N(32) should be ignored. After that all blocks should be indexed.
	if n, err := chain.InsertReceiptChain(blocks, receipts, 64); err != nil {
		t.Fatalf("block %d: failed to insert into chain: %v", n, err)
	}
	tail := uint64(32)
	check(&tail, chain)
}
