// Copyright 2021 The go-ethereum Authors
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

// Package ethconfig contains the configuration of the ETH and LES protocols.
package ethconfig

import (
    "github.com/ethereum/go-ethereum/consensus/ethashb3"
    "github.com/ethereum/go-ethereum/core"
    "os"
    "os/user"
    "path/filepath"
    "runtime"
    "time"

    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/consensus"
    "github.com/ethereum/go-ethereum/consensus/beacon"
    "github.com/ethereum/go-ethereum/consensus/clique"
    "github.com/ethereum/go-ethereum/consensus/ethash"
    "github.com/ethereum/go-ethereum/consensus/lyra2"
    "github.com/ethereum/go-ethereum/core/txpool/blobpool"
    "github.com/ethereum/go-ethereum/core/txpool/legacypool"
    "github.com/ethereum/go-ethereum/eth/downloader"
    "github.com/ethereum/go-ethereum/eth/gasprice"
    "github.com/ethereum/go-ethereum/ethdb"
    "github.com/ethereum/go-ethereum/log"
    "github.com/ethereum/go-ethereum/miner"
    "github.com/ethereum/go-ethereum/node"
    "github.com/ethereum/go-ethereum/params/types/ctypes"
    "github.com/ethereum/go-ethereum/params/types/genesisT"
    "github.com/ethereum/go-ethereum/params/vars"
)

// FullNodeGPO contains default gasprice oracle settings for full node.
var FullNodeGPO = gasprice.Config{
    Blocks:           20,
    Percentile:       60,
    MaxHeaderHistory: 1024,
    MaxBlockHistory:  1024,
    MaxPrice:         gasprice.DefaultMaxPrice,
    IgnorePrice:      gasprice.DefaultIgnorePrice,
}

// LightClientGPO contains default gasprice oracle settings for light client.
var LightClientGPO = gasprice.Config{
    Blocks:           2,
    Percentile:       60,
    MaxHeaderHistory: 300,
    MaxBlockHistory:  5,
    MaxPrice:         gasprice.DefaultMaxPrice,
    IgnorePrice:      gasprice.DefaultIgnorePrice,
}

// Defaults contains default settings for use on the Ethereum main net.
var Defaults = Config{
    SyncMode: downloader.SnapSync,
    Ethash: ethash.Config{
        CacheDir:         "ethash",
        CachesInMem:      2,
        CachesOnDisk:     3,
        CachesLockMmap:   false,
        DatasetsInMem:    1,
        DatasetsOnDisk:   2,
        DatasetsLockMmap: false,
    },

    EthashB3: ethashb3.Config{
        CacheDir:         "ethashb3",
        CachesInMem:      2,
        CachesOnDisk:     3,
        CachesLockMmap:   false,
        DatasetsInMem:    1,
        DatasetsOnDisk:   2,
        DatasetsLockMmap: false,
    },

    NetworkId:          0, // enable auto configuration of networkID == chainID
    ProtocolVersions:   vars.DefaultProtocolVersions,
    TxLookupLimit:      2350000,
    TransactionHistory: 2350000,
    StateHistory:       vars.FullImmutabilityThreshold,
    LightPeers:         100,
    UltraLightFraction: 75,
    DatabaseCache:      512,
    TrieCleanCache:     154,
    TrieDirtyCache:     256,
    TrieTimeout:        60 * time.Minute,
    SnapshotCache:      102,
    FilterLogCacheSize: 32,
    Miner:              miner.DefaultConfig,
    TxPool:             legacypool.DefaultConfig,
    BlobPool:           blobpool.DefaultConfig,
    RPCGasCap:          50000000,
    RPCEVMTimeout:      5 * time.Second,
    GPO:                FullNodeGPO,
    RPCTxFeeCap:        1, // 1 ether
}

func init() {
    home := os.Getenv("HOME")
    if home == "" {
        if user, err := user.Current(); err == nil {
            home = user.HomeDir
        }
    }
    if runtime.GOOS == "darwin" {
        Defaults.Ethash.DatasetDir = filepath.Join(home, "Library", "Ethash")
        Defaults.EthashB3.DatasetDir = filepath.Join(home, "Library", "EthashB3")
    } else if runtime.GOOS == "windows" {
        localappdata := os.Getenv("LOCALAPPDATA")
        if localappdata != "" {
            Defaults.Ethash.DatasetDir = filepath.Join(localappdata, "Ethash")
            Defaults.EthashB3.DatasetDir = filepath.Join(localappdata, "EthashB3")
        } else {
            Defaults.Ethash.DatasetDir = filepath.Join(home, "AppData", "Local", "Ethash")
            Defaults.EthashB3.DatasetDir = filepath.Join(home, "AppData", "Local", "EthashB3")
        }
    } else {
        Defaults.Ethash.DatasetDir = filepath.Join(home, ".ethash")
        Defaults.EthashB3.DatasetDir = filepath.Join(home, ".ethashb3")
    }
}

//go:generate go run github.com/fjl/gencodec -type Config -formats toml -out gen_config.go

// Config contains configuration options for ETH and LES protocols.
type Config struct {
    // The genesis block, which is inserted if the database is empty.
    // If nil, the Ethereum main net block is used.
    Genesis *genesisT.Genesis `toml:",omitempty"`

    // Protocol options
    NetworkId        uint64 // Network ID to use for selecting peers to connect to. When 0, chainID is used.
    ProtocolVersions []uint // Protocol versions are the supported versions of the eth protocol (first is primary).
    SyncMode         downloader.SyncMode

    // This can be set to list of enrtree:// URLs which will be queried for
    // for nodes to connect to.
    EthDiscoveryURLs  []string
    SnapDiscoveryURLs []string

    NoPruning  bool // Whether to disable pruning and flush everything to disk
    NoPrefetch bool // Whether to disable prefetching and only load state on demand

    // Deprecated, use 'TransactionHistory' instead.
    TxLookupLimit      uint64 `toml:",omitempty"` // The maximum number of blocks from head whose tx indices are reserved.
    TransactionHistory uint64 `toml:",omitempty"` // The maximum number of blocks from head whose tx indices are reserved.
    StateHistory       uint64 `toml:",omitempty"` // The maximum number of blocks from head whose state histories are reserved.

    // State scheme represents the scheme used to store ethereum states and trie
    // nodes on top. It can be 'hash', 'path', or none which means use the scheme
    // consistent with persistent state.
    StateScheme string `toml:",omitempty"`

    // RequiredBlocks is a set of block number -> hash mappings which must be in the
    // canonical chain of all remote peers. Setting the option makes geth verify the
    // presence of these blocks for every new peer connection.
    RequiredBlocks map[uint64]common.Hash `toml:"-"`

    // Light client options
    LightServ          int  `toml:",omitempty"` // Maximum percentage of time allowed for serving LES requests
    LightIngress       int  `toml:",omitempty"` // Incoming bandwidth limit for light servers
    LightEgress        int  `toml:",omitempty"` // Outgoing bandwidth limit for light servers
    LightPeers         int  `toml:",omitempty"` // Maximum number of LES client peers
    LightNoPrune       bool `toml:",omitempty"` // Whether to disable light chain pruning
    LightNoSyncServe   bool `toml:",omitempty"` // Whether to serve light clients before syncing
    SyncFromCheckpoint bool `toml:",omitempty"` // Whether to sync the header chain from the configured checkpoint

    // Ultra Light client options
    UltraLightServers      []string `toml:",omitempty"` // List of trusted ultra light servers
    UltraLightFraction     int      `toml:",omitempty"` // Percentage of trusted servers to accept an announcement
    UltraLightOnlyAnnounce bool     `toml:",omitempty"` // Whether to only announce headers, or also serve them

    // Database options
    SkipBcVersionCheck    bool `toml:"-"`
    DatabaseHandles       int  `toml:"-"`
    DatabaseCache         int
    DatabaseFreezer       string
    DatabaseFreezerRemote string

    TrieCleanCache int
    TrieDirtyCache int
    TrieTimeout    time.Duration
    SnapshotCache  int
    Preimages      bool

    // This is the number of blocks for which logs will be cached in the filter system.
    FilterLogCacheSize int

    // Mining options
    Miner miner.Config

    // Ethash options
    Ethash ethash.Config

    // EthashB3 options
    EthashB3 ethashb3.Config

    // Transaction pool options
    TxPool   legacypool.Config
    BlobPool blobpool.Config

    // Gas Price Oracle options
    GPO gasprice.Config

    // Enables tracking of SHA3 preimages in the VM
    EnablePreimageRecording bool

    // Miscellaneous options
    DocRoot string `toml:"-"`

    // Type of the EWASM interpreter ("" for default)
    EWASMInterpreter string

    // Type of the EVM interpreter ("" for default)
    EVMInterpreter string

    // RPCGasCap is the global gas cap for eth-call variants.
    RPCGasCap uint64

    // RPCEVMTimeout is the global timeout for eth-call.
    RPCEVMTimeout time.Duration

    // RPCTxFeeCap is the global transaction fee(price * gaslimit) cap for
    // send-transaction variants. The unit is ether.
    RPCTxFeeCap float64

    // Checkpoint is a hardcoded checkpoint which can be nil.
    Checkpoint *ctypes.TrustedCheckpoint `toml:",omitempty"`

    // CheckpointOracle is the configuration for checkpoint oracle.
    CheckpointOracle *ctypes.CheckpointOracleConfig `toml:",omitempty"`

    // Manual configuration field for ECBP1100 activation number. Used for modifying genesis config via CLI flag.
    OverrideECBP1100 *uint64 `toml:",omitempty"`
    // Manual configuration field for ECBP1100's disablement block number. Used for modifying genesis config via CLI flag.
    OverrideECBP1100Deactivate *uint64 `toml:",omitempty"`

    // ECBP1100NoDisable overrides
    // When this value is *true, ECBP100 will not (ever) be disabled; when *false, it will never be enabled.
    ECBP1100NoDisable *bool `toml:",omitempty"`

    // OverrideShanghai (TODO: remove after the fork)
    OverrideShanghai *uint64 `toml:",omitempty"`

    // OverrideCancun (TODO: remove after the fork)
    OverrideCancun *uint64 `toml:",omitempty"`

    // OverrideVerkle (TODO: remove after the fork)
    OverrideVerkle *uint64 `toml:",omitempty"`
}

// CreateConsensusEngine creates a consensus engine for the given chain configuration.
// Option for simplifying this would be the following code, but would require changed to each consensus engine to accept a config struct.
// ```
// type EngineCreator func(*node.Node, *Config, ethdb.Database) consensus.Engine
//
// var engineCreators = map[ctypes.ConsensusEngineT]EngineCreator{
//     ctypes.ConsensusEngineT_Clique:   CreateConsensusEngineClique,
//     ctypes.ConsensusEngineT_Lyra2:    CreateConsensusEngineLyra2,
//     ctypes.ConsensusEngineT_EthashB3: CreateConsensusEngineEthashB3,
//     ctypes.ConsensusEngineT_Ethash:   CreateConsensusEngineEthash,
// }
//
// func CreateConsensusEngine(stack *node.Node, config *Config, db ethdb.Database) consensus.Engine {
//     engineType := config.Genesis.Config.GetConsensusEngineType()
//     if createEngine, ok := engineCreators[engineType]; ok {
//         return createEngine(stack, config, db)
//     }
// }
// ```
func CreateConsensusEngine(stack *node.Node, config *Config, db ethdb.Database) consensus.Engine {
    var engine consensus.Engine
    engineType := config.Genesis.Config.GetConsensusEngineType()
    if engineType == ctypes.ConsensusEngineT_Clique {
        cliqueConfig, _ := core.LoadCliqueConfig(db, config.Genesis)
        engine = clique.New(cliqueConfig, db)
    } else if engineType == ctypes.ConsensusEngineT_Lyra2 {
        engine = lyra2.New(&lyra2.Config{}, config.Miner.Notify, config.Miner.Noverify)
    } else if engineType == ctypes.ConsensusEngineT_EthashB3 {
        engine = CreateConsensusEngineEthashB3(stack, &config.EthashB3, config.Miner.Notify, config.Miner.Noverify)
    } else {
        // If no engine type is found, default to ethash - Beacon chain
        engine = CreateConsensusEngineEthash(stack, &config.Ethash, config.Miner.Notify, config.Miner.Noverify)
    }

    return engine
}

func CreateConsensusEngineEthash(stack *node.Node, ethashConfig *ethash.Config, notify []string, noverify bool) (engine consensus.Engine) {
    switch ethashConfig.PowMode {
    case ethash.ModeFake:
        log.Warn("Ethash used in fake mode")
        engine = ethash.NewFaker()
    case ethash.ModeTest:
        log.Warn("Ethash used in test mode")
        engine = ethash.NewTester(nil, noverify)
    case ethash.ModeShared:
        log.Warn("Ethash used in shared mode")
        engine = ethash.NewShared()
    case ethash.ModePoissonFake:
        log.Warn("Ethash used in fake Poisson mode")
        engine = ethash.NewPoissonFaker()
    default:
        engine = ethash.New(ethash.Config{
            PowMode:          ethashConfig.PowMode,
            CacheDir:         stack.ResolvePath(ethashConfig.CacheDir),
            CachesInMem:      ethashConfig.CachesInMem,
            CachesOnDisk:     ethashConfig.CachesOnDisk,
            CachesLockMmap:   ethashConfig.CachesLockMmap,
            DatasetDir:       ethashConfig.DatasetDir,
            DatasetsInMem:    ethashConfig.DatasetsInMem,
            DatasetsOnDisk:   ethashConfig.DatasetsOnDisk,
            DatasetsLockMmap: ethashConfig.DatasetsLockMmap,
            NotifyFull:       ethashConfig.NotifyFull,
            ECIP1099Block:    ethashConfig.ECIP1099Block,
        }, notify, noverify)
        engine.(*ethash.Ethash).SetThreads(-1) // Disable CPU mining
    }
    engine = beacon.New(engine)
    return
}

func CreateConsensusEngineEthashB3(stack *node.Node, ethashb3Config *ethashb3.Config, notify []string, noverify bool) (engine consensus.Engine) {
    switch ethashb3Config.PowMode {
    case ethashb3.ModeFake:
        log.Warn("EthashB3 used in fake mode")
        engine = ethashb3.NewFaker()
    case ethashb3.ModeTest:
        log.Warn("EthashB3 used in test mode")
        engine = ethashb3.NewTester(nil, noverify)
    case ethashb3.ModeShared:
        log.Warn("EthashB3 used in shared mode")
        engine = ethashb3.NewShared()
    case ethashb3.ModePoissonFake:
        log.Warn("EthashB3 used in fake Poisson mode")
        engine = ethashb3.NewPoissonFaker()
    default:
        engine = ethashb3.New(ethashb3.Config{
            PowMode:          ethashb3Config.PowMode,
            CacheDir:         stack.ResolvePath(ethashb3Config.CacheDir),
            CachesInMem:      ethashb3Config.CachesInMem,
            CachesOnDisk:     ethashb3Config.CachesOnDisk,
            CachesLockMmap:   ethashb3Config.CachesLockMmap,
            DatasetDir:       ethashb3Config.DatasetDir,
            DatasetsInMem:    ethashb3Config.DatasetsInMem,
            DatasetsOnDisk:   ethashb3Config.DatasetsOnDisk,
            DatasetsLockMmap: ethashb3Config.DatasetsLockMmap,
            NotifyFull:       ethashb3Config.NotifyFull,
        }, notify, noverify)
        engine.(*ethashb3.EthashB3).SetThreads(-1) // Disable CPU mining
    }
    return
}

func CreateConsensusEngineClique(cliqueConfig *ctypes.CliqueConfig, db ethdb.Database) (engine consensus.Engine) {
    engine = clique.New(cliqueConfig, db)
    return
}

func CreateConsensusEngineLyra2(lyra2Config *lyra2.Config, notify []string, noverify bool) (engine consensus.Engine) {
    engine = lyra2.New(lyra2Config, notify, noverify)
    return
}
