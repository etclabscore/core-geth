package miner

import (
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/consensus/clique"
	"github.com/ethereum/go-ethereum/consensus/ethash"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/params/types/coregeth"
	"github.com/ethereum/go-ethereum/params/types/ctypes"
	"github.com/ethereum/go-ethereum/params/types/goethereum"
	"github.com/ethereum/go-ethereum/triedb"
)

func testGenerateBlockAndImportCG(t *testing.T, chainConfig ctypes.ChainConfigurator, numBlocks int) {
	t.Parallel()
	var (
		engine consensus.Engine
		db     = rawdb.NewMemoryDatabase()
	)
	if chainConfig.GetConsensusEngineType().IsClique() {
		// We need to ensure that the clique period is 0 (blocks created per tx) or 1 to avoid using a "dirty" period
		// with a greater time, which would cause the test to timeout.
		if chainConfig.GetCliquePeriod() != 0 {
			t.Logf("adjusting clique chain config: was: %d, want: %d", chainConfig.GetCliquePeriod(), 0)
			chainConfig.SetCliquePeriod(0)
		}
		engine = clique.New(&ctypes.CliqueConfig{
			Period: chainConfig.GetCliquePeriod(),
			Epoch:  chainConfig.GetCliqueEpoch(),
		}, db)
	} else if chainConfig.GetConsensusEngineType().IsEthash() {
		engine = ethash.NewFaker()
	}

	w, b := newTestWorker(t, chainConfig, engine, db, 0)
	defer w.close()

	// This test chain imports the mined blocks.
	db2 := rawdb.NewMemoryDatabase()
	core.MustCommitGenesis(db2, triedb.NewDatabase(db2, nil), b.genesis)
	chain, _ := core.NewBlockChain(db2, nil, b.genesis, nil, engine, vm.Config{}, nil, nil)
	defer chain.Stop()

	// Ignore empty commit here for less noise.
	w.skipSealHook = func(task *task) bool {
		return len(task.receipts) == 0
	}

	// Wait for mined blocks.
	sub := w.mux.Subscribe(core.NewMinedBlockEvent{})
	defer sub.Unsubscribe()

	// Start mining!
	w.start()

	for i := 0; i < numBlocks; i++ {
		b.txPool.Add([]*types.Transaction{b.newRandomTx(true)}, true, false)
		b.txPool.Add([]*types.Transaction{b.newRandomTx(false)}, true, false)
		w.postSideBlock(core.ChainSideEvent{Block: b.newRandomUncle()})
		w.postSideBlock(core.ChainSideEvent{Block: b.newRandomUncle()})

		select {
		case ev := <-sub.Chan():
			block := ev.Data.(core.NewMinedBlockEvent).Block
			if _, err := chain.InsertChain([]*types.Block{block}); err != nil {
				t.Fatalf("failed to insert new mined block %d: %v", block.NumberU64(), err)
			}
		case <-time.After(3 * time.Second): // Worker needs 1s to include new changes.
			t.Fatalf("timeout: %d", i)
		}
	}
}

func TestGenerateBlockAndImport_CG1(t *testing.T) {
	cases := []struct {
		name   string
		conf   ctypes.ChainConfigurator
		blocks int
	}{
		{
			name:   "all-ethash",
			conf:   params.AllEthashProtocolChanges,
			blocks: 5,
		},
		{
			name:   "all-clique",
			conf:   params.AllCliqueProtocolChanges,
			blocks: 5,
		},
		{
			name: "eth-dao=true-ethash",
			conf: func() ctypes.ChainConfigurator {
				c := &goethereum.ChainConfig{
					NetworkID:               1,
					ChainID:                 big.NewInt(1),
					HomesteadBlock:          big.NewInt(1),
					DAOForkBlock:            big.NewInt(2),
					DAOForkSupport:          true,
					EIP150Block:             big.NewInt(3),
					EIP150Hash:              common.Hash{},
					EIP155Block:             big.NewInt(4),
					EIP158Block:             big.NewInt(4),
					ByzantiumBlock:          big.NewInt(5),
					ConstantinopleBlock:     big.NewInt(6),
					PetersburgBlock:         big.NewInt(6),
					IstanbulBlock:           big.NewInt(7),
					MuirGlacierBlock:        big.NewInt(8),
					EWASMBlock:              nil,
					Ethash:                  new(ctypes.EthashConfig),
					Clique:                  nil,
					TrustedCheckpoint:       nil,
					TrustedCheckpointOracle: nil,
					EIP1706Transition:       nil,
					ECIP1080Transition:      nil,
				}
				return c
			}(),
			blocks: 10,
		},
		{
			name: "eth-dao=false-ethash",
			conf: func() ctypes.ChainConfigurator {
				c := &goethereum.ChainConfig{
					NetworkID:               1,
					ChainID:                 big.NewInt(1),
					HomesteadBlock:          big.NewInt(1),
					DAOForkBlock:            big.NewInt(2),
					DAOForkSupport:          false,
					EIP150Block:             big.NewInt(3),
					EIP150Hash:              common.Hash{},
					EIP155Block:             big.NewInt(4),
					EIP158Block:             big.NewInt(4),
					ByzantiumBlock:          big.NewInt(5),
					ConstantinopleBlock:     big.NewInt(6),
					PetersburgBlock:         big.NewInt(6),
					IstanbulBlock:           big.NewInt(7),
					MuirGlacierBlock:        big.NewInt(8),
					EWASMBlock:              nil,
					Ethash:                  new(ctypes.EthashConfig),
					Clique:                  nil,
					TrustedCheckpoint:       nil,
					TrustedCheckpointOracle: nil,
					EIP1706Transition:       nil,
					ECIP1080Transition:      nil,
				}
				return c
			}(),
			blocks: 10,
		},
		{
			name: "classic-ethash",
			conf: func() ctypes.ChainConfigurator {
				c := &coregeth.CoreGethChainConfig{
					NetworkID: 1,
					Ethash:    new(ctypes.EthashConfig),
					ChainID:   big.NewInt(61),

					EIP2FBlock: big.NewInt(1),
					EIP7FBlock: big.NewInt(1),

					EIP150Block: big.NewInt(2),

					EIP155Block:  big.NewInt(3),
					EIP160FBlock: big.NewInt(3),

					// EIP158~
					EIP161FBlock: big.NewInt(6),
					EIP170FBlock: big.NewInt(6),

					// Byzantium eq
					EIP100FBlock: big.NewInt(6),
					EIP140FBlock: big.NewInt(6),
					EIP198FBlock: big.NewInt(6),
					EIP211FBlock: big.NewInt(6),
					EIP212FBlock: big.NewInt(6),
					EIP213FBlock: big.NewInt(6),
					EIP214FBlock: big.NewInt(6),
					EIP658FBlock: big.NewInt(6),

					// Constantinople eq, aka Agharta
					EIP145FBlock:  big.NewInt(7),
					EIP1014FBlock: big.NewInt(7),
					EIP1052FBlock: big.NewInt(7),

					// Istanbul eq, aka Phoenix
					// ECIP-1088
					EIP152FBlock:  big.NewInt(8),
					EIP1108FBlock: big.NewInt(8),
					EIP1344FBlock: big.NewInt(8),
					EIP1884FBlock: big.NewInt(8),
					EIP2028FBlock: big.NewInt(8),
					EIP2200FBlock: big.NewInt(8), // RePetersburg (=~ re-1283)

					DisposalBlock:      big.NewInt(5),
					ECIP1017FBlock:     big.NewInt(4),
					ECIP1017EraRounds:  big.NewInt(4),
					ECIP1010PauseBlock: big.NewInt(3),
					ECIP1010Length:     big.NewInt(2),
				}
				return c
			}(),
			blocks: 10,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			testGenerateBlockAndImportCG(t, c.conf, c.blocks)
		})
	}
}
