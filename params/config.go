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
)

// Genesis hashes to enforce below configs on.
var (
	MainnetGenesisHash = common.HexToHash("0xd4e56740f876aef8c010b86a40d5f56745a118d0906a34e69aec8c0db1cb8fa3")
	RopstenGenesisHash = common.HexToHash("0x41941023680923e0fe4d74a34bdac8141f2540e3ae90623718e47d66d1ca4a2d")
	RinkebyGenesisHash = common.HexToHash("0x6341fd3daf94b748c72ced5a5b26028f2474f5f00d824504e4fa37a75767e177")
	GoerliGenesisHash  = common.HexToHash("0xbf7e331f7f7c1dd2e05159666b3bf8bc7a8a3a9eb1d518969eab529dd9b88c1a")
	YoloV2GenesisHash  = common.HexToHash("0xc3fd235071f24f93865b0850bd2a2119b30f7224d18a0e34c7bbf549ad7e3d36") // TODO: check if yolov2 is correct
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
		ChainID:                 big.NewInt(1),
		HomesteadBlock:          big.NewInt(1150000),
		DAOForkBlock:            big.NewInt(1920000),
		DAOForkSupport:          true,
		EIP150Block:             big.NewInt(2463000),
		EIP150Hash:              common.HexToHash("0x2086799aeebeae135c246c65021c82b4e15a2c451340993aacfd2751886514f0"),
		EIP155Block:             big.NewInt(2675000),
		EIP158Block:             big.NewInt(2675000),
		ByzantiumBlock:          big.NewInt(4370000),
		ConstantinopleBlock:     big.NewInt(7280000),
		PetersburgBlock:         big.NewInt(7280000),
		IstanbulBlock:           big.NewInt(9069000),
		MuirGlacierBlock:        big.NewInt(9200000),
		Ethash:                  new(ctypes.EthashConfig),
		TrustedCheckpoint:       MainnetTrustedCheckpoint,
		TrustedCheckpointOracle: MainnetCheckpointOracle,
	}

	// MainnetTrustedCheckpoint contains the light client trusted checkpoint for the main network.
	MainnetTrustedCheckpoint = &ctypes.TrustedCheckpoint{
		SectionIndex: 336,
		SectionHead:  common.HexToHash("0xd42b78902b6527a80337bf1bc372a3ccc3db97e9cc7cf421ca047ae9076c716b"),
		CHTRoot:      common.HexToHash("0xd97f3b30f7e0cb958e4c67c53ec27745e5a165e33e56821b86523dfee62b783a"),
		BloomRoot:    common.HexToHash("0xf3cbfd070fababfe2adc9b23fc02c731f6ca2cce6646b3ede4ef2db06092ccce"),
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
		ChainID:             big.NewInt(3),
		HomesteadBlock:      big.NewInt(0),
		DAOForkBlock:        nil,
		DAOForkSupport:      true,
		EIP150Block:         big.NewInt(0),
		EIP150Hash:          common.HexToHash("0x41941023680923e0fe4d74a34bdac8141f2540e3ae90623718e47d66d1ca4a2d"),
		EIP155Block:         big.NewInt(10),
		EIP158Block:         big.NewInt(10),
		ByzantiumBlock:      big.NewInt(1700000),
		ConstantinopleBlock: big.NewInt(4230000),
		PetersburgBlock:     big.NewInt(4939394),
		IstanbulBlock:       big.NewInt(6485846),
		MuirGlacierBlock:    big.NewInt(7117117),
		Ethash:              new(ctypes.EthashConfig),
	}

	// RopstenTrustedCheckpoint contains the light client trusted checkpoint for the Ropsten test network.
	RopstenTrustedCheckpoint = &ctypes.TrustedCheckpoint{
		SectionIndex: 269,
		SectionHead:  common.HexToHash("0x290a9eb65e65c64601d1b05522533ed502098a246736b348502a170818a33d64"),
		CHTRoot:      common.HexToHash("0x530ebac02264227277d0a16b0819ef96a2011a6e1e66523ebff8040f4a3437ca"),
		BloomRoot:    common.HexToHash("0x480cd5b3198a0767022902130546854a2e8867cce573c1cf0ce54e67a7bf5efb"),
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
		ChainID:                 big.NewInt(4),
		HomesteadBlock:          big.NewInt(1),
		DAOForkBlock:            nil,
		DAOForkSupport:          true,
		EIP150Block:             big.NewInt(2),
		EIP150Hash:              common.HexToHash("0x9b095b36c15eaf13044373aef8ee0bd3a382a5abb92e402afa44b8249c3a90e9"),
		EIP155Block:             big.NewInt(3),
		EIP158Block:             big.NewInt(3),
		ByzantiumBlock:          big.NewInt(1035301),
		ConstantinopleBlock:     big.NewInt(3660663),
		PetersburgBlock:         big.NewInt(4321234),
		IstanbulBlock:           big.NewInt(5435345),
		MuirGlacierBlock:        nil,
		TrustedCheckpoint:       RinkebyTrustedCheckpoint,
		TrustedCheckpointOracle: RinkebyCheckpointOracle,
		Clique: &ctypes.CliqueConfig{
			Period: 15,
			Epoch:  30000,
		},
	}

	// RinkebyTrustedCheckpoint contains the light client trusted checkpoint for the Rinkeby test network.
	RinkebyTrustedCheckpoint = &ctypes.TrustedCheckpoint{
		SectionIndex: 223,
		SectionHead:  common.HexToHash("0x03ca0d5e3a931c77cd7a97bbaa2d9e4edc4549c621dc1d223a29f10c86a4a16a"),
		CHTRoot:      common.HexToHash("0x6573dbdd91b2958b446bd04d67c23e5f14b2510ac96e8df1b6a894dc49e37c6c"),
		BloomRoot:    common.HexToHash("0x28a35042a4e88efbac55fe566faf7fce000dc436f17fd4cb4b081c9cd793e1a7"),
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
		ChainID:                 big.NewInt(5),
		HomesteadBlock:          big.NewInt(0),
		DAOForkBlock:            nil,
		DAOForkSupport:          true,
		EIP150Block:             big.NewInt(0),
		EIP155Block:             big.NewInt(0),
		EIP158Block:             big.NewInt(0),
		ByzantiumBlock:          big.NewInt(0),
		ConstantinopleBlock:     big.NewInt(0),
		PetersburgBlock:         big.NewInt(0),
		IstanbulBlock:           big.NewInt(1561651),
		MuirGlacierBlock:        nil,
		TrustedCheckpoint:       GoerliTrustedCheckpoint,
		TrustedCheckpointOracle: GoerliCheckpointOracle,
		Clique: &ctypes.CliqueConfig{
			Period: 15,
			Epoch:  30000,
		},
	}

	// GoerliTrustedCheckpoint contains the light client trusted checkpoint for the Görli test network.
	GoerliTrustedCheckpoint = &ctypes.TrustedCheckpoint{
		SectionIndex: 107,
		SectionHead:  common.HexToHash("0xff3ae39199fa191894de419e7f673c8627aa8cc7af924b90f36635b6add375f2"),
		CHTRoot:      common.HexToHash("0x27d59d60c652425b6b593a882f55a4ff57f24e470a810a6e3c8ba71833a20220"),
		BloomRoot:    common.HexToHash("0x3c14066d8bb3733780c06b8165768dbb9dd23b75f56012fe5f2fb3c2fb70cadb"),
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

	// YoloV2ChainConfig contains the chain parameters to run a node on the YOLOv1 test network.
	YoloV2ChainConfig = &goethereum.ChainConfig{
		ChainID:             big.NewInt(133519467574833),
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
		MuirGlacierBlock:    nil,
		YoloV2Block:         big.NewInt(0),
		Clique: &ctypes.CliqueConfig{
			Period: 15,
			Epoch:  30000,
		},
	}

	// AllEthashProtocolChanges contains every protocol change (EIPs) introduced
	// and accepted by the Ethereum core developers into the Ethash consensus.
	//
	// This configuration is intentionally not using keyed fields to force anyone
	// adding flags to the config to also have to set these fields.
	AllEthashProtocolChanges = &goethereum.ChainConfig{
		ChainID:                 big.NewInt(1137),
		HomesteadBlock:          big.NewInt(0),
		DAOForkBlock:            nil,
		DAOForkSupport:          false,
		EIP150Block:             big.NewInt(0),
		EIP150Hash:              common.Hash{},
		EIP155Block:             big.NewInt(0),
		EIP158Block:             big.NewInt(0),
		ByzantiumBlock:          big.NewInt(0),
		ConstantinopleBlock:     big.NewInt(0),
		PetersburgBlock:         big.NewInt(0),
		IstanbulBlock:           big.NewInt(0),
		MuirGlacierBlock:        nil,
		YoloV2Block:             nil,
		EWASMBlock:              nil,
		Ethash:                  new(ctypes.EthashConfig),
		Clique:                  nil,
		TrustedCheckpoint:       nil,
		TrustedCheckpointOracle: nil,
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
		ChainID:             big.NewInt(1337),
		HomesteadBlock:      big.NewInt(0),
		DAOForkBlock:        nil,
		DAOForkSupport:      false,
		EIP150Block:         big.NewInt(0),
		EIP150Hash:          common.Hash{},
		EIP155Block:         big.NewInt(0),
		EIP158Block:         big.NewInt(0),
		ByzantiumBlock:      big.NewInt(0),
		ConstantinopleBlock: big.NewInt(0),
		PetersburgBlock:     big.NewInt(0),
		IstanbulBlock:       big.NewInt(0),
		EWASMBlock:          nil,
		Ethash:              nil,
		Clique: &ctypes.CliqueConfig{
			Period: 0,
			Epoch:  30000,
		},
		TrustedCheckpoint:       nil,
		TrustedCheckpointOracle: nil,
	}

	// TestChainConfig is used for tests.
	TestChainConfig = &goethereum.ChainConfig{
		ChainID:                 big.NewInt(1),
		HomesteadBlock:          big.NewInt(0),
		DAOForkBlock:            nil,
		DAOForkSupport:          false,
		EIP150Block:             big.NewInt(0),
		EIP150Hash:              common.Hash{},
		EIP155Block:             big.NewInt(0),
		EIP158Block:             big.NewInt(0),
		ByzantiumBlock:          big.NewInt(0),
		ConstantinopleBlock:     big.NewInt(0),
		PetersburgBlock:         big.NewInt(0),
		IstanbulBlock:           big.NewInt(0),
		EWASMBlock:              nil,
		Ethash:                  new(ctypes.EthashConfig),
		Clique:                  nil,
		TrustedCheckpoint:       nil,
		TrustedCheckpointOracle: nil,
	}
)
