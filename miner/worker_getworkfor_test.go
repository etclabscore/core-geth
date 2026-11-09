// Copyright 2026 The core-geth Authors
// This file is part of the core-geth library.
//
// The core-geth library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The core-geth library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the core-geth library. If not, see <http://www.gnu.org/licenses/>.

package miner

import (
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/consensus/clique"
	"github.com/ethereum/go-ethereum/consensus/ethash"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params/types/ctypes"
	"github.com/ethereum/go-ethereum/params/types/goethereum"
)

var (
	getWorkForPayeeA = common.HexToAddress("0x00000000000000000000000000000000000a11ce")
	getWorkForPayeeB = common.HexToAddress("0x0000000000000000000000000000000000000b0b")
)

// newGetWorkForTestWorker creates a worker on top of a test-mode ethash engine
// with local CPU sealing disabled, so that all work only goes to the remote
// sealer. It returns the worker, its backend and the engine's ethash RPC API.
func newGetWorkForTestWorker(t *testing.T, chainConfig ctypes.ChainConfigurator) (*worker, *testWorkerBackend, *ethash.API) {
	engine := ethash.NewTester(nil, false)
	engine.SetThreads(-1)
	t.Cleanup(func() { engine.Close() })

	w, b := newTestWorker(t, chainConfig, engine, rawdb.NewMemoryDatabase(), 0)
	t.Cleanup(w.close)

	api, ok := engine.APIs(b.chain)[0].Service.(*ethash.API)
	if !ok {
		t.Fatalf("unexpected ethash API type")
	}
	return w, b, api
}

// waitForRemoteWork waits until the remote sealer hands out regular work.
func waitForRemoteWork(t *testing.T, api *ethash.API) [4]string {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for {
		work, err := api.GetWork()
		if err == nil {
			return work
		}
		if time.Now().After(deadline) {
			t.Fatalf("no remote work available: %v", err)
		}
		time.Sleep(10 * time.Millisecond)
	}
}

// sealWithTester finds a valid proof-of-work for the given block using a
// separate test-mode ethash instance.
func sealWithTester(t *testing.T, block *types.Block) *types.Block {
	t.Helper()
	sealer := ethash.NewTester(nil, false)
	defer sealer.Close()

	var (
		results = make(chan *types.Block, 1)
		stop    = make(chan struct{})
	)
	defer close(stop)
	if err := sealer.Seal(nil, block, results, stop); err != nil {
		t.Fatalf("failed to seal block: %v", err)
	}
	select {
	case sealed := <-results:
		return sealed
	case <-time.After(30 * time.Second):
		t.Fatalf("sealing timeout")
	}
	return nil
}

// waitForMinedBlock waits until the worker announced the given mined block.
func waitForMinedBlock(t *testing.T, sub chan *types.Block, hash common.Hash) {
	t.Helper()
	timeout := time.After(5 * time.Second)
	for {
		select {
		case block := <-sub:
			if block.Hash() == hash {
				return
			}
		case <-timeout:
			t.Fatalf("mined block %x not announced", hash)
		}
	}
}

func pendingTaskFor(w *worker, sealhash string) *task {
	w.pendingMu.RLock()
	defer w.pendingMu.RUnlock()
	return w.pendingTasks[common.HexToHash(sealhash)]
}

// Tests that getWorkFor hands out work packages paying the requested coinbase,
// that it doesn't touch the regular work package, and that solutions for the
// packages are imported as new chain blocks.
func TestGetWorkFor(t *testing.T) {
	t.Parallel()

	w, b, api := newGetWorkForTestWorker(t, ethashChainConfig)

	// Wait until both the empty and the full regular sealing tasks were pushed.
	tasks := make(chan struct{}, 2)
	w.newTaskHook = func(*task) {
		select {
		case tasks <- struct{}{}:
		default:
		}
	}
	w.start()
	for i := 0; i < 2; i++ {
		select {
		case <-tasks:
		case <-time.After(5 * time.Second):
			t.Fatal("regular sealing task timeout")
		}
	}
	waitForRemoteWork(t, api)
	time.Sleep(200 * time.Millisecond) // let the remote sealer process the last task
	regular := waitForRemoteWork(t, api)

	workA, err := w.getWorkFor(getWorkForPayeeA)
	if err != nil {
		t.Fatalf("getWorkFor(A) failed: %v", err)
	}
	workB, err := w.getWorkFor(getWorkForPayeeB)
	if err != nil {
		t.Fatalf("getWorkFor(B) failed: %v", err)
	}
	if workA[0] == workB[0] {
		t.Fatalf("different coinbases share the same sealhash %s", workA[0])
	}
	if workA[0] == regular[0] || workB[0] == regular[0] {
		t.Fatalf("getWorkFor returned the regular work package")
	}
	// The regular work package must be untouched.
	if current, err := api.GetWork(); err != nil || current != regular {
		t.Fatalf("regular work changed: have %v (err %v), want %v", current, err, regular)
	}
	for _, c := range []struct {
		work     [4]string
		coinbase common.Address
	}{{workA, getWorkForPayeeA}, {workB, getWorkForPayeeB}} {
		task := pendingTaskFor(w, c.work[0])
		if task == nil {
			t.Fatalf("no pending task for %s", c.work[0])
		}
		block := task.block
		if block.Coinbase() != c.coinbase {
			t.Errorf("coinbase mismatch: have %x, want %x", block.Coinbase(), c.coinbase)
		}
		if block.NumberU64() != 1 || c.work[3] != hexutil.EncodeUint64(1) {
			t.Errorf("number mismatch: block %d, package %s", block.NumberU64(), c.work[3])
		}
		if block.ParentHash() != b.chain.Genesis().Hash() {
			t.Errorf("parent mismatch: have %x, want %x", block.ParentHash(), b.chain.Genesis().Hash())
		}
		if len(block.Transactions()) != len(pendingTxs) || len(task.receipts) != len(pendingTxs) {
			t.Errorf("pending transactions not included: txs %d, receipts %d", len(block.Transactions()), len(task.receipts))
		}
		if c.work[1] != regular[1] {
			t.Errorf("seed hash mismatch: have %s, want %s", c.work[1], regular[1])
		}
		boundary := common.BytesToHash(new(big.Int).Div(new(big.Int).Lsh(common.Big1, 256), block.Difficulty()).Bytes()).Hex()
		if c.work[2] != boundary {
			t.Errorf("boundary mismatch: have %s, want %s", c.work[2], boundary)
		}
	}

	mined := make(chan *types.Block, 16)
	sub := w.mux.Subscribe(core.NewMinedBlockEvent{})
	defer sub.Unsubscribe()
	go func() {
		for ev := range sub.Chan() {
			mined <- ev.Data.(core.NewMinedBlockEvent).Block
		}
	}()

	// A valid solution for A's package becomes the new head paying A.
	sealedA := sealWithTester(t, pendingTaskFor(w, workA[0]).block)
	if !api.SubmitWork(sealedA.Header().Nonce, common.HexToHash(workA[0]), sealedA.MixDigest()) {
		t.Fatal("valid solution for A rejected")
	}
	waitForMinedBlock(t, mined, sealedA.Hash())
	head := b.chain.CurrentBlock()
	if head.Hash() != sealedA.Hash() {
		t.Fatalf("head mismatch: have %x, want %x", head.Hash(), sealedA.Hash())
	}
	if head.Coinbase != getWorkForPayeeA {
		t.Fatalf("head coinbase mismatch: have %x, want %x", head.Coinbase, getWorkForPayeeA)
	}
	state, err := b.chain.State()
	if err != nil {
		t.Fatal(err)
	}
	if state.GetBalance(getWorkForPayeeA).IsZero() {
		t.Error("coinbase A was not paid")
	}

	// B's package is still within the staleness window and is imported too.
	sealedB := sealWithTester(t, pendingTaskFor(w, workB[0]).block)
	if !api.SubmitWork(sealedB.Header().Nonce, common.HexToHash(workB[0]), sealedB.MixDigest()) {
		t.Fatal("valid solution for B rejected")
	}
	waitForMinedBlock(t, mined, sealedB.Hash())
	if block := b.chain.GetBlockByHash(sealedB.Hash()); block == nil || block.Coinbase() != getWorkForPayeeB {
		t.Fatalf("block for B not imported with coinbase B")
	}
}

// Tests the error cases of getWorkFor.
func TestGetWorkForErrors(t *testing.T) {
	t.Parallel()

	t.Run("ethash", func(t *testing.T) {
		w, _, api := newGetWorkForTestWorker(t, ethashChainConfig)

		if _, err := w.getWorkFor(common.Address{}); err != errZeroCoinbase {
			t.Errorf("zero coinbase: have %v, want %v", err, errZeroCoinbase)
		}
		if _, err := w.getWorkFor(getWorkForPayeeA); err != errMinerNotRunning {
			t.Errorf("miner stopped: have %v, want %v", err, errMinerNotRunning)
		}
		w.start()
		waitForRemoteWork(t, api)

		w.syncing.Store(true)
		if _, err := w.getWorkFor(getWorkForPayeeA); err != errMinerSyncing {
			t.Errorf("syncing: have %v, want %v", err, errMinerSyncing)
		}
		w.syncing.Store(false)

		if _, err := w.getWorkFor(getWorkForPayeeA); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		w.stop()
		if _, err := w.getWorkFor(getWorkForPayeeA); err != errMinerNotRunning {
			t.Errorf("miner stopped: have %v, want %v", err, errMinerNotRunning)
		}
	})

	t.Run("clique", func(t *testing.T) {
		engine := clique.New(&ctypes.CliqueConfig{
			Period: cliqueChainConfig.GetCliquePeriod(),
			Epoch:  cliqueChainConfig.GetCliqueEpoch(),
		}, rawdb.NewMemoryDatabase())
		defer engine.Close()
		w, _ := newTestWorker(t, cliqueChainConfig, engine, rawdb.NewMemoryDatabase(), 0)
		defer w.close()

		if _, err := w.getWorkFor(getWorkForPayeeA); err != errRemoteWorkNotSupported {
			t.Errorf("clique: have %v, want %v", err, errRemoteWorkNotSupported)
		}
	})

	t.Run("ethash without remote sealer", func(t *testing.T) {
		engine := ethash.NewFaker()
		defer engine.Close()
		w, _ := newTestWorker(t, ethashChainConfig, engine, rawdb.NewMemoryDatabase(), 0)
		defer w.close()
		w.skipSealHook = func(*task) bool { return true }
		w.start()

		if _, err := w.getWorkFor(getWorkForPayeeA); err == nil || err.Error() != "not supported" {
			t.Errorf("faker: have %v, want not supported", err)
		}
		w.pendingMu.RLock()
		defer w.pendingMu.RUnlock()
		if len(w.pendingTasks) != 0 {
			t.Errorf("pending task left behind after failure")
		}
	})

	t.Run("post-merge", func(t *testing.T) {
		config := ctypes.ChainConfigurator(new(goethereum.ChainConfig))
		config.SetEthashTerminalTotalDifficulty(big.NewInt(0))
		w, _, _ := newGetWorkForTestWorker(t, config)
		w.start()

		if _, err := w.getWorkFor(getWorkForPayeeA); err != errTTDReached {
			t.Errorf("post-merge: have %v, want %v", err, errTTDReached)
		}
	})
}
