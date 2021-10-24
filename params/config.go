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
	RinkebyGenesisHash = common.HexToHash("0x6341fd3daf94b748c72ced5a5b26028f2474f5f00d824504e4fa37a75767e177")
	GoerliGenesisHash  = common.HexToHash("0xbf7e331f7f7c1dd2e05159666b3bf8bc7a8a3a9eb1d518969eab529dd9b88c1a")
)

// TrustedCheckpoints associates each known checkpoint with the genesis hash of
// the chain it belongs to.
var TrustedCheckpoints = map[common.Hash]*ctypes.TrustedCheckpoint{
	MainnetGenesisHash: MainnetTrustedCheckpoint,
	RopstenGenesisHash: RopstenTrustedCheckpoint,
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
		Ethash:                    new(ctypes.EthashConfig),
	}

	// MainnetTrustedCheckpoint contains the light client trusted checkpoint for the main network.
	MainnetTrustedCheckpoint = &ctypes.TrustedCheckpoint{
		SectionIndex: 395,
		SectionHead:  common.HexToHash("0xbfca95b8c1de014e252288e9c32029825fadbff58285f5b54556525e480dbb5b"),
		CHTRoot:      common.HexToHash("0x2ccf3dbb58eb6375e037fdd981ca5778359e4b8fa0270c2878b14361e64161e7"),
		BloomRoot:    common.HexToHash("0x2d46ec65a6941a2dc1e682f8f81f3d24192021f492fdf6ef0fdd51acb0f4ba0f"),
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
		Ethash:                    new(ctypes.EthashConfig),
	}

	// RopstenTrustedCheckpoint contains the light client trusted checkpoint for the Ropsten test network.
	RopstenTrustedCheckpoint = &ctypes.TrustedCheckpoint{
		SectionIndex: 329,
		SectionHead:  common.HexToHash("0xe66f7038333a01fb95dc9ea03e5a2bdaf4b833cdcb9e393b9127e013bd64d39b"),
		CHTRoot:      common.HexToHash("0x1b0c883338ac0d032122800c155a2e73105fbfebfaa50436893282bc2d9feec5"),
		BloomRoot:    common.HexToHash("0x3cc98c88d283bf002378246f22c653007655cbcea6ed89f98d739f73bd341a01"),
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
		TrustedCheckpoint:         RinkebyTrustedCheckpoint,
		TrustedCheckpointOracle:   RinkebyCheckpointOracle,
		Clique: &ctypes.CliqueConfig{
			Period: 15,
			Epoch:  30000,
		},
	}

	// RinkebyTrustedCheckpoint contains the light client trusted checkpoint for the Rinkeby test network.
	RinkebyTrustedCheckpoint = &ctypes.TrustedCheckpoint{
		SectionIndex: 276,
		SectionHead:  common.HexToHash("0xea89a4b04e3da9bd688e316f8de669396b6d4a38a19d2cd96a00b70d58b836aa"),
		CHTRoot:      common.HexToHash("0xd6889d0bf6673c0d2c1cf6e9098a6fe5b30888a115b6112796aa8ee8efc4a723"),
		BloomRoot:    common.HexToHash("0x6009a9256b34b8bde3a3f094afb647ba5d73237546017b9025d64ac1ff54c47c"),
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
		TrustedCheckpoint:         GoerliTrustedCheckpoint,
		TrustedCheckpointOracle:   GoerliCheckpointOracle,
		Clique: &ctypes.CliqueConfig{
			Period: 15,
			Epoch:  30000,
		},
	}

	// GoerliTrustedCheckpoint contains the light client trusted checkpoint for the Görli test network.
	GoerliTrustedCheckpoint = &ctypes.TrustedCheckpoint{
		SectionIndex: 160,
		SectionHead:  common.HexToHash("0xb5a666c790dc35a5613d04ebba8ba47a850b45a15d9b95ad7745c35ae034b5a5"),
		CHTRoot:      common.HexToHash("0x6b4e00df52bdc38fa6c26c8ef595c2ad6184963ea36ab08ee744af460aa735e1"),
		BloomRoot:    common.HexToHash("0x8fa88f5e50190cb25243aeee262a1a9e4434a06f8d455885dcc1b5fc48c33836"),
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
		EWASMBlock:                nil,
		CatalystBlock:             nil,
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

			YoloV2Block *big.Int `json:"yoloV2Block,omitempty"` // YOLO v1: https://github.com/ethereum/EIPs/pull/2657 (Ephemeral testnet)
					EWASMBlock  *big.Int `json:"ewasmBlock,omitempty"`  // EWASM switch block (nil = no fork, 0 = already activated)

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
		EWASMBlock:                nil,
		CatalystBlock:             nil,
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
		EWASMBlock:                nil,
		CatalystBlock:             nil,
		Ethash:                    new(ctypes.EthashConfig),
		Clique:                    nil,
		TrustedCheckpoint:         nil,
		TrustedCheckpointOracle:   nil,
	}
)
