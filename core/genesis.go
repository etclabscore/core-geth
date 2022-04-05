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
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/params/confp"
	"github.com/ethereum/go-ethereum/params/types/ctypes"
	"github.com/ethereum/go-ethereum/params/types/genesisT"
	"github.com/ethereum/go-ethereum/params/vars"
	"github.com/ethereum/go-ethereum/trie"
)

// SetupGenesisBlock wraps SetupGenesisBlockWithOverride, always using a nil value for the override.
func SetupGenesisBlock(db ethdb.Database, genesis *genesisT.Genesis) (ctypes.ChainConfigurator, common.Hash, error) {
	return SetupGenesisBlockWithOverride(db, genesis, nil, nil)
}

func SetupGenesisBlockWithOverride(db ethdb.Database, genesis *genesisT.Genesis, overrideMystique, overrideTerminalTotalDifficulty *big.Int) (ctypes.ChainConfigurator, common.Hash, error) {
	if genesis != nil && confp.IsEmpty(genesis.Config) {
		return params.AllEthashProtocolChanges, common.Hash{}, genesisT.ErrGenesisNoConfig
	}
	// Just commit the new block if there is no stored genesis block.
	stored := rawdb.ReadCanonicalHash(db, 0)
	if (stored == common.Hash{}) {
		if genesis == nil {
			log.Info("Writing default main-net genesis block")
			log.Warn("Not specifying a chain flag is deprecated and will be removed in the future, please use --mainnet for Ethereum mainnet")
			genesis = params.DefaultGenesisBlock()
		} else {
			log.Info("Writing custom genesis block")
		}

		if overrideMystique != nil {
			n := overrideMystique.Uint64()
			if err := genesis.SetEIP1559Transition(&n); err != nil {
				return genesis, stored, err
			}
			if err := genesis.SetEIP3198Transition(&n); err != nil {
				return genesis, stored, err
			}
			if err := genesis.SetEIP3529Transition(&n); err != nil {
				return genesis, stored, err
			}
			if err := genesis.SetEIP3541Transition(&n); err != nil {
				return genesis, stored, err
			}
			if err := genesis.SetEthashEIP3554Transition(&n); err != nil {
				return genesis, stored, err
			}
		}

		block, err := CommitGenesis(genesis, db)
		if err != nil {
			return genesis.Config, common.Hash{}, err
		}
		log.Info("Wrote custom genesis block OK", "config", genesis.Config)
		return genesis.Config, block.Hash(), nil
	}
	// We have the genesis block in database(perhaps in ancient database)
	// but the corresponding state is missing.
	header := rawdb.ReadHeader(db, stored, 0)
	if _, err := state.New(header.Root, state.NewDatabaseWithConfig(db, nil), nil); err != nil {
		if genesis == nil {
			genesis = params.DefaultGenesisBlock()
		}
		// Ensure the stored genesis matches with the given one.
		hash := GenesisToBlock(genesis, nil).Hash()
		if hash != stored {
			return genesis.Config, hash, &genesisT.GenesisMismatchError{Stored: stored, New: hash}
		}
		block, err := CommitGenesis(genesis, db)
		if err != nil {
			return genesis.Config, hash, err
		}
		return genesis.Config, block.Hash(), nil
	}
	// Check whether the genesis block is already written.
	if genesis != nil {
		hash := GenesisToBlock(genesis, nil).Hash()
		if hash != stored {
			return genesis.Config, hash, &genesisT.GenesisMismatchError{Stored: stored, New: hash}
		}
	}
	// Get the existing chain configuration.
	newcfg := configOrDefault(genesis, stored)

	if overrideTerminalTotalDifficulty != nil {
		newcfg.SetEthashTerminalTotalDifficulty(overrideTerminalTotalDifficulty)
	}

	if overrideMystique != nil {
		n := overrideMystique.Uint64()
		if err := newcfg.SetEIP1559Transition(&n); err != nil {
			return newcfg, stored, err
		}
		if err := newcfg.SetEIP3198Transition(&n); err != nil {
			return newcfg, stored, err
		}
		if err := newcfg.SetEIP3529Transition(&n); err != nil {
			return newcfg, stored, err
		}
		if err := newcfg.SetEIP3541Transition(&n); err != nil {
			return newcfg, stored, err
		}
		if err := newcfg.SetEthashEIP3554Transition(&n); err != nil {
			return newcfg, stored, err
		}
	}

	// TODO (ziogaschr): Add EIPs, after Mystique activations
	// if overrideArrowGlacier != nil {
	// 	newcfg.ArrowGlacierBlock = overrideArrowGlacier
	// }

	storedcfg := rawdb.ReadChainConfig(db, stored)
	if storedcfg == nil {
		log.Warn("Found genesis block without chain config")
		rawdb.WriteChainConfig(db, stored, newcfg)
		return newcfg, stored, nil
	} else {
		log.Info("Found stored genesis block", "config", storedcfg)
	}

	// Special case: don't change the existing config of a non-mainnet chain if no new
	// config is supplied. These chains would get AllProtocolChanges (and a compat error)
	// if we just continued here.
	//
	// (meowsbits): The idea here is to use stored configs when they are not upgrade-able via defaults.
	// Pre-existing logic only upgraded mainnet, when it should upgrade all defaulty chains.
	// New logic (below) checks _inequality_ between a defaulty config and a stored config. If different,
	// the stored config is used. This breaks auto-upgrade magic for defaulty chains.
	if genesis == nil && !confp.Identical(storedcfg, newcfg, []string{"NetworkID", "ChainID"}) {
		// TODO/meowsbits/20220405: ethereum code for this scope follows:
		/*
			// Special case: if a private network is being used (no genesis and also no
			// mainnet hash in the database), we must not apply the `configOrDefault`
			// chain config as that would be AllProtocolChanges (applying any new fork
			// on top of an existing private network genesis block). In that case, only
			// apply the overrides.

			if ... :
			newcfg = storedcfg
			if overrideArrowGlacier != nil {
				newcfg.ArrowGlacierBlock = overrideArrowGlacier
			}
			if overrideTerminalTotalDifficulty != nil {
				newcfg.TerminalTotalDifficulty = overrideTerminalTotalDifficulty
			}
		*/
		// ... and this is ours:
		log.Info("Found non-defaulty stored config, using it.")
		return storedcfg, stored, nil
	}
	// Check config compatibility and write the config. Compatibility errors
	// are returned to the caller unless we're already at block zero.
	height := rawdb.ReadHeaderNumber(db, rawdb.ReadHeadHeaderHash(db))
	if height == nil {
		return newcfg, stored, fmt.Errorf("missing block number for head header hash")
	}
	compatErr := confp.Compatible(height, storedcfg, newcfg)
	if compatErr != nil && *height != 0 && compatErr.RewindTo != 0 {
		return newcfg, stored, compatErr
	}
	rawdb.WriteChainConfig(db, stored, newcfg)
	return newcfg, stored, nil
}

func configOrDefault(g *genesisT.Genesis, ghash common.Hash) ctypes.ChainConfigurator {
	switch {
	case g != nil:
		return g.Config
	case ghash == params.MainnetGenesisHash:
		return params.MainnetChainConfig
	case ghash == params.RinkebyGenesisHash:
		return params.RinkebyChainConfig
	case ghash == params.GoerliGenesisHash:
		return params.GoerliChainConfig
	case ghash == params.KottiGenesisHash:
		return params.KottiChainConfig
	case ghash == params.MordorGenesisHash:
		return params.MordorChainConfig
	case ghash == params.RopstenGenesisHash:
		return params.RopstenChainConfig
	case ghash == params.SepoliaGenesisHash:
		return params.SepoliaChainConfig
	case ghash == params.MintMeGenesisHash:
		return params.MintMeChainConfig
	case ghash == params.KilnGenesisHash:
		return params.DefaultKilnGenesisBlock().Config
	default:
		return params.AllEthashProtocolChanges
	}
}

// Flush adds allocated genesis accounts into a fresh new statedb and
// commit the state changes into the given database handler.
func gaFlush(ga *genesisT.GenesisAlloc, db ethdb.Database) (common.Hash, error) {
	statedb, err := state.New(common.Hash{}, state.NewDatabase(db), nil)
	if err != nil {
		return common.Hash{}, err
	}
	for addr, account := range *ga {
		statedb.AddBalance(addr, account.Balance)
		statedb.SetCode(addr, account.Code)
		statedb.SetNonce(addr, account.Nonce)
		for key, value := range account.Storage {
			statedb.SetState(addr, key, value)
		}
	}
	root, err := statedb.Commit(false)
	if err != nil {
		return common.Hash{}, err
	}
	err = statedb.Database().TrieDB().Commit(root, true, nil)
	if err != nil {
		return common.Hash{}, err
	}
	return root, nil
}

// Write writes the json marshaled genesis state into database
// with the given block hash as the unique identifier.
func gaWrite(ga *genesisT.GenesisAlloc, db ethdb.KeyValueWriter, hash common.Hash) error {
	blob, err := json.Marshal(ga)
	if err != nil {
		return err
	}
	rawdb.WriteGenesisState(db, hash, blob)
	return nil
}

// CommitGenesisState loads the stored genesis state with the given block
// hash and commits them into the given database handler.
func CommitGenesisState(db ethdb.Database, hash common.Hash) error {
	var alloc genesisT.GenesisAlloc
	blob := rawdb.ReadGenesisState(db, hash)
	if len(blob) != 0 {
		if err := alloc.UnmarshalJSON(blob); err != nil {
			return err
		}
	} else {
		// Genesis allocation is missing and there are several possibilities:
		// the node is legacy which doesn't persist the genesis allocation or
		// the persisted allocation is just lost.
		// - supported networks(mainnet, testnets), recover with defined allocations
		// - private network, can't recover
		var genesis *genesisT.Genesis
		switch hash {
		case params.MainnetGenesisHash:
			genesis = params.DefaultGenesisBlock()
			// TODO/meowsbits/20220405: make sure we don't need Classic in here
		case params.RopstenGenesisHash:
			genesis = params.DefaultRopstenGenesisBlock()
		case params.RinkebyGenesisHash:
			genesis = params.DefaultRinkebyGenesisBlock()
		case params.GoerliGenesisHash:
			genesis = params.DefaultGoerliGenesisBlock()
		case params.SepoliaGenesisHash:
			genesis = params.DefaultSepoliaGenesisBlock()
		case params.KottiGenesisHash:
			genesis = params.DefaultKottiGenesisBlock()
		case params.MordorGenesisHash:
			genesis = params.DefaultMordorGenesisBlock()
		case params.MintMeGenesisHash:
			genesis = params.DefaultMintMeGenesisBlock()
		}
		if genesis != nil {
			alloc = genesis.Alloc
		} else {
			return errors.New("not found")
		}
	}
	_, err := gaFlush(&alloc, db)
	return err
}

// GenesisToBlock creates the genesis block and writes state of a genesis specification
// to the given database (or discards it if nil).
func GenesisToBlock(g *genesisT.Genesis, db ethdb.Database) *types.Block {
	if db == nil {
		db = rawdb.NewMemoryDatabase()
	}
	root, err := gaFlush(&g.Alloc, db)
	if err != nil {
		panic(err)
	}
	head := &types.Header{
		Number:     new(big.Int).SetUint64(g.Number),
		Nonce:      types.EncodeNonce(g.Nonce),
		Time:       g.Timestamp,
		ParentHash: g.ParentHash,
		Extra:      g.ExtraData,
		GasLimit:   g.GasLimit,
		GasUsed:    g.GasUsed,
		BaseFee:    g.BaseFee,
		Difficulty: g.Difficulty,
		MixDigest:  g.Mixhash,
		Coinbase:   g.Coinbase,
		Root:       root,
	}
	if g.GasLimit == 0 {
		head.GasLimit = vars.GenesisGasLimit
	}
	// -- meowsbits/202203 go-ethereum has: if g.Difficulty == nil && g.Mixhash == (common.Hash{}) {
	// They also assign the Difficulty field directly.
	if g.Difficulty == nil {
		head.Difficulty = new(big.Int)
		head.Difficulty.Set(vars.GenesisDifficulty)
	}
	if g.Config != nil && g.Config.IsEnabled(g.Config.GetEIP1559Transition, common.Big0) {
		if g.BaseFee != nil {
			head.BaseFee = g.BaseFee
		} else {
			head.BaseFee = new(big.Int).SetUint64(vars.InitialBaseFee)
		}
	}
	return types.NewBlock(head, nil, nil, nil, trie.NewStackTrie(nil))
}

// CommitGenesis writes the block and state of a genesis specification to the database.
// The block is committed as the canonical head block.
func CommitGenesis(g *genesisT.Genesis, db ethdb.Database) (*types.Block, error) {
	block := GenesisToBlock(g, db)
	if block.Number().Sign() != 0 {
		return nil, errors.New("can't commit genesis block with number > 0")
	}
	config := g.Config
	if config == nil {
		config = params.AllEthashProtocolChanges
	}
	if config.GetConsensusEngineType().IsClique() && len(block.Extra()) == 0 {
		return nil, errors.New("can't start clique chain without signers")
	}
	if err := gaWrite(&g.Alloc, db, block.Hash()); err != nil {
		return nil, err
	}
	rawdb.WriteTd(db, block.Hash(), block.NumberU64(), block.Difficulty())
	rawdb.WriteBlock(db, block)
	rawdb.WriteReceipts(db, block.Hash(), block.NumberU64(), nil)
	rawdb.WriteCanonicalHash(db, block.Hash(), block.NumberU64())
	rawdb.WriteHeadBlockHash(db, block.Hash())
	rawdb.WriteHeadFastBlockHash(db, block.Hash())
	rawdb.WriteHeadHeaderHash(db, block.Hash())
	rawdb.WriteChainConfig(db, block.Hash(), config)
	return block, nil
}

// MustCommitGenesis writes the genesis block and state to db, panicking on error.
// The block is committed as the canonical head block.
func MustCommitGenesis(db ethdb.Database, g *genesisT.Genesis) *types.Block {
	block, err := CommitGenesis(g, db)
	if err != nil {
		panic(err)
	}
	return block
}

// GenesisBlockForTesting creates and writes a block in which addr has the given wei balance.
func GenesisBlockForTesting(db ethdb.Database, addr common.Address, balance *big.Int) *types.Block {
	g := genesisT.Genesis{
		Alloc:   genesisT.GenesisAlloc{addr: {Balance: balance}},
		BaseFee: big.NewInt(vars.InitialBaseFee),
	}
	return MustCommitGenesis(db, &g)
}
