// Copyright 2016 The go-ethereum Authors
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

package params

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params/types/ctypes"
	"github.com/ethereum/go-ethereum/params/types/goethereum"
	"github.com/ethereum/go-ethereum/params/vars"
)

// Genesis hashes to enforce below configs on.
var (
	MainnetGenesisHash = common.HexToHash("0xd4e56740f876aef8c010b86a40d5f56745a118d0906a34e69aec8c0db1cb8fa3")
	RopstenGenesisHash = common.HexToHash("0x41941023680923e0fe4d74a34bdac8141f2540e3ae90623718e47d66d1ca4a2d")
	SepoliaGenesisHash = common.HexToHash("0x25a5cc106eea7138acab33231d7160d69cb777ee0c2c553fcddf5138993e6dd9")
	RinkebyGenesisHash = common.HexToHash("0x6341fd3daf94b748c72ced5a5b26028f2474f5f00d824504e4fa37a75767e177")
	GoerliGenesisHash  = common.HexToHash("0xbf7e331f7f7c1dd2e05159666b3bf8bc7a8a3a9eb1d518969eab529dd9b88c1a")
	KilnGenesisHash    = common.HexToHash("0x51c7fe41be669f69c45c33a56982cbde405313342d9e2b00d7c91a7b284dd4f8")
)

// TrustedCheckpoints associates each known checkpoint with the genesis hash of
// the chain it belongs to.
var TrustedCheckpoints = map[common.Hash]*ctypes.TrustedCheckpoint{
	MainnetGenesisHash: MainnetTrustedCheckpoint,
	RopstenGenesisHash: RopstenTrustedCheckpoint,
	SepoliaGenesisHash: SepoliaTrustedCheckpoint,
	RinkebyGenesisHash: RinkebyTrustedCheckpoint,
	GoerliGenesisHash:  GoerliTrustedCheckpoint,
}

// CheckpointOracles associates each known checkpoint oracles with the genesis hash of
// the chain it belongs to.
var CheckpointOracles = map[common.Hash]*ctypes.CheckpointOracleConfig{
	MainnetGenesisHash: MainnetCheckpointOracle,
	RopstenGenesisHash: RopstenCheckpointOracle,
	RinkebyGenesisHash: RinkebyCheckpointOracle,
	GoerliGenesisHash:  GoerliCheckpointOracle,
}

var (
	// MainnetChainConfig is the chain parameters to run a node on the main network.
	MainnetChainConfig = &goethereum.ChainConfig{
		ChainID:                   big.NewInt(1),
		SupportedProtocolVersions: vars.DefaultProtocolVersions,
		HomesteadBlock:            big.NewInt(1_150_000),
		DAOForkBlock:              big.NewInt(1_920_000),
		DAOForkSupport:            true,
		EIP150Block:               big.NewInt(2_463_000),
		EIP150Hash:                common.HexToHash("0x2086799aeebeae135c246c65021c82b4e15a2c451340993aacfd2751886514f0"),
		EIP155Block:               big.NewInt(2_675_000),
		EIP158Block:               big.NewInt(2_675_000),
		ByzantiumBlock:            big.NewInt(4_370_000),
		ConstantinopleBlock:       big.NewInt(7_280_000),
		PetersburgBlock:           big.NewInt(7_280_000),
		IstanbulBlock:             big.NewInt(9_069_000),
		MuirGlacierBlock:          big.NewInt(9_200_000),
		BerlinBlock:               big.NewInt(12_244_000),
		LondonBlock:               big.NewInt(12_965_000),
		ArrowGlacierBlock:         big.NewInt(13_773_000),
		GrayGlacierBlock:          big.NewInt(15_050_000),
		Ethash:                    new(ctypes.EthashConfig),
	}

	// MainnetTrustedCheckpoint contains the light client trusted checkpoint for the main network.
	MainnetTrustedCheckpoint = &ctypes.TrustedCheckpoint{
		SectionIndex: 451,
		SectionHead:  common.HexToHash("0xe47f84b9967eb2ad2afff74d59901b63134660011822fdababaf8fdd18a75aa6"),
		CHTRoot:      common.HexToHash("0xc31e0462ca3d39a46111bb6b63ac4e1cac84089472b7474a319d582f72b3f0c0"),
		BloomRoot:    common.HexToHash("0x7c9f25ce3577a3ab330d52a7343f801899cf9d4980c69f81de31ccc1a055c809"),
	}

	// MainnetCheckpointOracle contains a set of configs for the main network oracle.
	MainnetCheckpointOracle = &ctypes.CheckpointOracleConfig{
		Address: common.HexToAddress("0x9a9070028361F7AAbeB3f2F2Dc07F82C4a98A02a"),
		Signers: []common.Address{
			common.HexToAddress("0x1b2C260efc720BE89101890E4Db589b44E950527"), // Peter
			common.HexToAddress("0x78d1aD571A1A09D60D9BBf25894b44e4C8859595"), // Martin
			common.HexToAddress("0x286834935f4A8Cfb4FF4C77D5770C2775aE2b0E7"), // Zsolt
			common.HexToAddress("0xb86e2B0Ab5A4B1373e40c51A7C712c70Ba2f9f8E"), // Gary
			common.HexToAddress("0x0DF8fa387C602AE62559cC4aFa4972A7045d6707"), // Guillaume
		},
		Threshold: 2,
	}

	// RopstenChainConfig contains the chain parameters to run a node on the Ropsten test network.
	RopstenChainConfig = &goethereum.ChainConfig{
		ChainID:                   big.NewInt(3),
		SupportedProtocolVersions: vars.DefaultProtocolVersions,
		HomesteadBlock:            big.NewInt(0),
		DAOForkBlock:              nil,
		DAOForkSupport:            true,
		EIP150Block:               big.NewInt(0),
		EIP150Hash:                common.HexToHash("0x41941023680923e0fe4d74a34bdac8141f2540e3ae90623718e47d66d1ca4a2d"),
		EIP155Block:               big.NewInt(10),
		EIP158Block:               big.NewInt(10),
		ByzantiumBlock:            big.NewInt(1_700_000),
		ConstantinopleBlock:       big.NewInt(4_230_000),
		PetersburgBlock:           big.NewInt(4_939_394),
		IstanbulBlock:             big.NewInt(6_485_846),
		MuirGlacierBlock:          big.NewInt(7_117_117),
		BerlinBlock:               big.NewInt(9_812_189),
		LondonBlock:               big.NewInt(10_499_401),
		TerminalTotalDifficulty:   big.NewInt(50000000000000000),
		Ethash:                    new(ctypes.EthashConfig),
	}

	// RopstenTrustedCheckpoint contains the light client trusted checkpoint for the Ropsten test network.
	RopstenTrustedCheckpoint = &ctypes.TrustedCheckpoint{
		SectionIndex: 346,
		SectionHead:  common.HexToHash("0xafa0384ebd13a751fb7475aaa7fc08ac308925c8b2e2195bca2d4ab1878a7a84"),
		CHTRoot:      common.HexToHash("0x522ae1f334bfa36033b2315d0b9954052780700b69448ecea8d5877e0f7ee477"),
		BloomRoot:    common.HexToHash("0x4093fd53b0d2cc50181dca353fe66f03ae113e7cb65f869a4dfb5905de6a0493"),
	}

	// RopstenCheckpointOracle contains a set of configs for the Ropsten test network oracle.
	RopstenCheckpointOracle = &ctypes.CheckpointOracleConfig{
		Address: common.HexToAddress("0xEF79475013f154E6A65b54cB2742867791bf0B84"),
		Signers: []common.Address{
			common.HexToAddress("0x32162F3581E88a5f62e8A61892B42C46E2c18f7b"), // Peter
			common.HexToAddress("0x78d1aD571A1A09D60D9BBf25894b44e4C8859595"), // Martin
			common.HexToAddress("0x286834935f4A8Cfb4FF4C77D5770C2775aE2b0E7"), // Zsolt
			common.HexToAddress("0xb86e2B0Ab5A4B1373e40c51A7C712c70Ba2f9f8E"), // Gary
			common.HexToAddress("0x0DF8fa387C602AE62559cC4aFa4972A7045d6707"), // Guillaume
		},
		Threshold: 2,
	}

	// SepoliaChainConfig contains the chain parameters to run a node on the Sepolia test network.
	SepoliaChainConfig = &goethereum.ChainConfig{
		ChainID:             big.NewInt(11155111),
		HomesteadBlock:      big.NewInt(0),
		DAOForkBlock:        nil,
		DAOForkSupport:      true,
		EIP150Block:         big.NewInt(0),
		EIP155Block:         big.NewInt(0),
		EIP158Block:         big.NewInt(0),
		ByzantiumBlock:      big.NewInt(0),
		ConstantinopleBlock: big.NewInt(0),
		PetersburgBlock:     big.NewInt(0),
		IstanbulBlock:       big.NewInt(0),
		MuirGlacierBlock:    big.NewInt(0),
		BerlinBlock:         big.NewInt(0),
		LondonBlock:         big.NewInt(0),
		Ethash:              new(ctypes.EthashConfig),
	}

	// SepoliaTrustedCheckpoint contains the light client trusted checkpoint for the Sepolia test network.
	SepoliaTrustedCheckpoint = &ctypes.TrustedCheckpoint{
		SectionIndex: 34,
		SectionHead:  common.HexToHash("0xe361400fcbc468d641e7bdd0b0946a3548e97c5d2703b124f04a3f1deccec244"),
		CHTRoot:      common.HexToHash("0xea6768fd288dce7d84f590884908ec39e4de78e6e1a38de5c5419b0f49a42f91"),
		BloomRoot:    common.HexToHash("0x06d32f35d5a611bfd0333ad44e39c619449824167d8ef2913edc48a8112be2cd"),
	}

	// RinkebyChainConfig contains the chain parameters to run a node on the Rinkeby test network.
	RinkebyChainConfig = &goethereum.ChainConfig{
		ChainID:                   big.NewInt(4),
		SupportedProtocolVersions: vars.DefaultProtocolVersions,
		HomesteadBlock:            big.NewInt(1),
		DAOForkBlock:              nil,
		DAOForkSupport:            true,
		EIP150Block:               big.NewInt(2),
		EIP150Hash:                common.HexToHash("0x9b095b36c15eaf13044373aef8ee0bd3a382a5abb92e402afa44b8249c3a90e9"),
		EIP155Block:               big.NewInt(3),
		EIP158Block:               big.NewInt(3),
		ByzantiumBlock:            big.NewInt(1_035_301),
		ConstantinopleBlock:       big.NewInt(3_660_663),
		PetersburgBlock:           big.NewInt(4_321_234),
		IstanbulBlock:             big.NewInt(5_435_345),
		MuirGlacierBlock:          nil,
		BerlinBlock:               big.NewInt(8_290_928),
		LondonBlock:               big.NewInt(8_897_988),
		ArrowGlacierBlock:         nil,
		TrustedCheckpoint:         RinkebyTrustedCheckpoint,
		TrustedCheckpointOracle:   RinkebyCheckpointOracle,
		Clique: &ctypes.CliqueConfig{
			Period: 15,
			Epoch:  30000,
		},
	}

	// RinkebyTrustedCheckpoint contains the light client trusted checkpoint for the Rinkeby test network.
	RinkebyTrustedCheckpoint = &ctypes.TrustedCheckpoint{
		SectionIndex: 326,
		SectionHead:  common.HexToHash("0x941a41a153b0e36cb15d9d193d1d0f9715bdb2435efd1c95119b64168667ce00"),
		CHTRoot:      common.HexToHash("0xe2331e00d579cf4093091dee35bef772e63c2341380c276041dc22563c8aba2e"),
		BloomRoot:    common.HexToHash("0x595206febcf118958c2bc1218ea71d01fd04b8f97ad71813df4be0af5b36b0e5"),
	}

	// RinkebyCheckpointOracle contains a set of configs for the Rinkeby test network oracle.
	RinkebyCheckpointOracle = &ctypes.CheckpointOracleConfig{
		Address: common.HexToAddress("0xebe8eFA441B9302A0d7eaECc277c09d20D684540"),
		Signers: []common.Address{
			common.HexToAddress("0xd9c9cd5f6779558b6e0ed4e6acf6b1947e7fa1f3"), // Peter
			common.HexToAddress("0x78d1aD571A1A09D60D9BBf25894b44e4C8859595"), // Martin
			common.HexToAddress("0x286834935f4A8Cfb4FF4C77D5770C2775aE2b0E7"), // Zsolt
			common.HexToAddress("0xb86e2B0Ab5A4B1373e40c51A7C712c70Ba2f9f8E"), // Gary
		},
		Threshold: 2,
	}

	// GoerliChainConfig contains the chain parameters to run a node on the Görli test network.
	GoerliChainConfig = &goethereum.ChainConfig{
		ChainID:                   big.NewInt(5),
		SupportedProtocolVersions: vars.DefaultProtocolVersions,
		HomesteadBlock:            big.NewInt(0),
		DAOForkBlock:              nil,
		DAOForkSupport:            true,
		EIP150Block:               big.NewInt(0),
		EIP155Block:               big.NewInt(0),
		EIP158Block:               big.NewInt(0),
		ByzantiumBlock:            big.NewInt(0),
		ConstantinopleBlock:       big.NewInt(0),
		PetersburgBlock:           big.NewInt(0),
		IstanbulBlock:             big.NewInt(1_561_651),
		MuirGlacierBlock:          nil,
		BerlinBlock:               big.NewInt(4_460_644),
		LondonBlock:               big.NewInt(5_062_605),
		ArrowGlacierBlock:         nil,
		TrustedCheckpoint:         GoerliTrustedCheckpoint,
		TrustedCheckpointOracle:   GoerliCheckpointOracle,
		Clique: &ctypes.CliqueConfig{
			Period: 15,
			Epoch:  30000,
		},
	}

	// GoerliTrustedCheckpoint contains the light client trusted checkpoint for the Görli test network.
	GoerliTrustedCheckpoint = &ctypes.TrustedCheckpoint{
		SectionIndex: 210,
		SectionHead:  common.HexToHash("0xbb11eaf551a6c06f74a6c7bbfe1699cbf64b8f248b64691da916dd443176db2f"),
		CHTRoot:      common.HexToHash("0x9934ae326d00d9c7de2e074c0e51689efb7fa7fcba18929ff4279c27259c45e6"),
		BloomRoot:    common.HexToHash("0x7fe3bd4fd45194aa8a5cfe5ac590edff1f870d3d98d3c310494e7f67613a87ff"),
	}

	// GoerliCheckpointOracle contains a set of configs for the Goerli test network oracle.
	GoerliCheckpointOracle = &ctypes.CheckpointOracleConfig{
		Address: common.HexToAddress("0x18CA0E045F0D772a851BC7e48357Bcaab0a0795D"),
		Signers: []common.Address{
			common.HexToAddress("0x4769bcaD07e3b938B7f43EB7D278Bc7Cb9efFb38"), // Peter
			common.HexToAddress("0x78d1aD571A1A09D60D9BBf25894b44e4C8859595"), // Martin
			common.HexToAddress("0x286834935f4A8Cfb4FF4C77D5770C2775aE2b0E7"), // Zsolt
			common.HexToAddress("0xb86e2B0Ab5A4B1373e40c51A7C712c70Ba2f9f8E"), // Gary
			common.HexToAddress("0x0DF8fa387C602AE62559cC4aFa4972A7045d6707"), // Guillaume
		},
		Threshold: 2,
	}

	// AllEthashProtocolChanges contains every protocol change (EIPs) introduced
	// and accepted by the Ethereum core developers into the Ethash consensus.
	//
	// This configuration is intentionally not using keyed fields to force anyone
	// adding flags to the config to also have to set these fields.
	AllEthashProtocolChanges = &goethereum.ChainConfig{
		ChainID:                   big.NewInt(1337),
		SupportedProtocolVersions: vars.SupportedProtocolVersions,
		HomesteadBlock:            big.NewInt(0),
		DAOForkBlock:              nil,
		DAOForkSupport:            false,
		EIP150Block:               big.NewInt(0),
		EIP150Hash:                common.Hash{},
		EIP155Block:               big.NewInt(0),
		EIP158Block:               big.NewInt(0),
		ByzantiumBlock:            big.NewInt(0),
		ConstantinopleBlock:       big.NewInt(0),
		PetersburgBlock:           big.NewInt(0),
		IstanbulBlock:             big.NewInt(0),
		MuirGlacierBlock:          big.NewInt(0),
		BerlinBlock:               big.NewInt(0),
		LondonBlock:               big.NewInt(0),
		ArrowGlacierBlock:         big.NewInt(0),
		GrayGlacierBlock:          big.NewInt(0),
		EWASMBlock:                nil,
		TerminalTotalDifficulty:   nil,
		Ethash:                    new(ctypes.EthashConfig),
		Clique:                    nil,
		TrustedCheckpoint:         nil,
		TrustedCheckpointOracle:   nil,
	}

	/*
					https://github.com/ethereum/go-ethereum/blob/master/params/config.go#L242

						AllEthashProtocolChanges = &ChainConfig{
						big.NewInt(1337),
						big.NewInt(0),
						nil,
						false,
						big.NewInt(0),
						common.Hash{},
						big.NewInt(0),
						big.NewInt(0),
						big.NewInt(0),
						big.NewInt(0),
						big.NewInt(0),
						big.NewInt(0),
						nil,
						nil,
						nil,
						new(EthashConfig),
						nil
						}


		// ChainConfig is the core config which determines the blockchain settings.
		//
		// ChainConfig is stored in the database on a per block basis. This means
		// that any network, identified by its genesis block, can have its own
		// set of configuration options.
		type ChainConfig struct {
			ChainID *big.Int `json:"chainId"` // chainId identifies the current chain and is used for replay protection

			HomesteadBlock *big.Int `json:"homesteadBlock,omitempty"` // Homestead switch block (nil = no fork, 0 = already homestead)

			DAOForkBlock   *big.Int `json:"daoForkBlock,omitempty"`   // TheDAO hard-fork switch block (nil = no fork)
			DAOForkSupport bool     `json:"daoForkSupport,omitempty"` // Whether the nodes supports or opposes the DAO hard-fork

			// EIP150 implements the Gas price changes (https://github.com/ethereum/EIPs/issues/150)
			EIP150Block *big.Int    `json:"eip150Block,omitempty"` // EIP150 HF block (nil = no fork)
			EIP150Hash  common.Hash `json:"eip150Hash,omitempty"`  // EIP150 HF hash (needed for header only clients as only gas pricing changed)

			EIP155Block *big.Int `json:"eip155Block,omitempty"` // EIP155 HF block
			EIP158Block *big.Int `json:"eip158Block,omitempty"` // EIP158 HF block

			ByzantiumBlock      *big.Int `json:"byzantiumBlock,omitempty"`      // Byzantium switch block (nil = no fork, 0 = already on byzantium)
			ConstantinopleBlock *big.Int `json:"constantinopleBlock,omitempty"` // Constantinople switch block (nil = no fork, 0 = already activated)
			PetersburgBlock     *big.Int `json:"petersburgBlock,omitempty"`     // Petersburg switch block (nil = same as Constantinople)
			IstanbulBlock       *big.Int `json:"istanbulBlock,omitempty"`       // Istanbul switch block (nil = no fork, 0 = already on istanbul)
			MuirGlacierBlock    *big.Int `json:"muirGlacierBlock,omitempty"`    // Eip-2384 (bomb delay) switch block (nil = no fork, 0 = already activated)
			BerlinBlock         *big.Int `json:"berlinBlock,omitempty"`         // Berlin switch block (nil = no fork, 0 = already on berlin)
			LondonBlock         *big.Int `json:"londonBlock,omitempty"`         // London switch block (nil = no fork, 0 = already on london)
			ArrowGlacierBlock   *big.Int `json:"arrowGlacierBlock,omitempty"`   // Eip-4345 (bomb delay) switch block (nil = no fork, 0 = already activated)
			GrayGlacierBlock    *big.Int `json:"grayGlacierBlock,omitempty"`    // Eip-5133 (bomb delay) switch block (nil = no fork, 0 = already activated)
			MergeNetsplitBlock  *big.Int `json:"mergeNetsplitBlock,omitempty"`  // Virtual fork after The Merge to use as a network splitter

			// TerminalTotalDifficulty is the amount of total difficulty reached by
			// the network that triggers the consensus upgrade.
			TerminalTotalDifficulty *big.Int `json:"terminalTotalDifficulty,omitempty"`

			// Various consensus engines
			Ethash *EthashConfig `json:"ethash,omitempty"`
			Clique *CliqueConfig `json:"clique,omitempty"`
		}

	*/

	// AllCliqueProtocolChanges contains every protocol change (EIPs) introduced
	// and accepted by the Ethereum core developers into the Clique consensus.
	//
	// This configuration is intentionally not using keyed fields to force anyone
	// adding flags to the config to also have to set these fields.
	AllCliqueProtocolChanges = &goethereum.ChainConfig{
		ChainID:                   big.NewInt(1337),
		SupportedProtocolVersions: vars.SupportedProtocolVersions,
		HomesteadBlock:            big.NewInt(0),
		DAOForkBlock:              nil,
		DAOForkSupport:            false,
		EIP150Block:               big.NewInt(0),
		EIP150Hash:                common.Hash{},
		EIP155Block:               big.NewInt(0),
		EIP158Block:               big.NewInt(0),
		ByzantiumBlock:            big.NewInt(0),
		ConstantinopleBlock:       big.NewInt(0),
		PetersburgBlock:           big.NewInt(0),
		IstanbulBlock:             big.NewInt(0),
		BerlinBlock:               big.NewInt(0),
		LondonBlock:               big.NewInt(0),
		ArrowGlacierBlock:         nil,
		EWASMBlock:                nil,
		TerminalTotalDifficulty:   nil,
		Ethash:                    nil,
		Clique: &ctypes.CliqueConfig{
			Period: 0,
			Epoch:  30000,
		},
		TrustedCheckpoint:       nil,
		TrustedCheckpointOracle: nil,
	}

	// TestChainConfig is used for tests.
	TestChainConfig = &goethereum.ChainConfig{
		ChainID:                   big.NewInt(1),
		SupportedProtocolVersions: vars.SupportedProtocolVersions,
		HomesteadBlock:            big.NewInt(0),
		DAOForkBlock:              nil,
		DAOForkSupport:            false,
		EIP150Block:               big.NewInt(0),
		EIP150Hash:                common.Hash{},
		EIP155Block:               big.NewInt(0),
		EIP158Block:               big.NewInt(0),
		ByzantiumBlock:            big.NewInt(0),
		ConstantinopleBlock:       big.NewInt(0),
		PetersburgBlock:           big.NewInt(0),
		IstanbulBlock:             big.NewInt(0),
		BerlinBlock:               big.NewInt(0),
		LondonBlock:               big.NewInt(0),
		ArrowGlacierBlock:         big.NewInt(0),
		EWASMBlock:                nil,
		TerminalTotalDifficulty:   nil,
		Ethash:                    new(ctypes.EthashConfig),
		Clique:                    nil,
		TrustedCheckpoint:         nil,
		TrustedCheckpointOracle:   nil,
	}
)
