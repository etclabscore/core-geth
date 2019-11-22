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

	"github.com/ethereum/go-ethereum/chainspecs/parity"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params/types"
)

// Genesis hashes to enforce below configs on.
var (
	MainnetGenesisHash = common.HexToHash("0xd4e56740f876aef8c010b86a40d5f56745a118d0906a34e69aec8c0db1cb8fa3")
	TestnetGenesisHash = common.HexToHash("0x41941023680923e0fe4d74a34bdac8141f2540e3ae90623718e47d66d1ca4a2d")
	RinkebyGenesisHash = common.HexToHash("0x6341fd3daf94b748c72ced5a5b26028f2474f5f00d824504e4fa37a75767e177")
	GoerliGenesisHash  = common.HexToHash("0xbf7e331f7f7c1dd2e05159666b3bf8bc7a8a3a9eb1d518969eab529dd9b88c1a")
)

// TrustedCheckpoints associates each known checkpoint with the genesis hash of
// the chain it belongs to.
var TrustedCheckpoints = map[common.Hash]*paramtypes.TrustedCheckpoint{
	MainnetGenesisHash: MainnetTrustedCheckpoint,
	TestnetGenesisHash: TestnetTrustedCheckpoint,
	RinkebyGenesisHash: RinkebyTrustedCheckpoint,
	GoerliGenesisHash:  GoerliTrustedCheckpoint,
}

// CheckpointOracles associates each known checkpoint oracles with the genesis hash of
// the chain it belongs to.
var CheckpointOracles = map[common.Hash]*paramtypes.CheckpointOracleConfig{
	MainnetGenesisHash: MainnetCheckpointOracle,
	TestnetGenesisHash: TestnetCheckpointOracle,
	RinkebyGenesisHash: RinkebyCheckpointOracle,
	GoerliGenesisHash:  GoerliCheckpointOracle,
}

var (
	// MainnetChainConfig is the chain parameters to run a node on the main network.
	MainnetChainConfig = &paramtypes.ChainConfig{
		NetworkID:           1,
		ChainID:             big.NewInt(1),
		HomesteadBlock:      big.NewInt(1150000),
		DAOForkBlock:        big.NewInt(1920000),
		DAOForkSupport:      true,
		EIP150Block:         big.NewInt(2463000),
		EIP150Hash:          common.HexToHash("0x2086799aeebeae135c246c65021c82b4e15a2c451340993aacfd2751886514f0"),
		EIP155Block:         big.NewInt(2675000),
		EIP158Block:         big.NewInt(2675000),
		ByzantiumBlock:      big.NewInt(4370000),
		DisposalBlock:       nil,
		SocialBlock:         nil,
		EthersocialBlock:    nil,
		ConstantinopleBlock: big.NewInt(7280000),
		PetersburgBlock:     big.NewInt(7280000),
		IstanbulBlock:       big.NewInt(9069000),
		Ethash:              new(paramtypes.EthashConfig),
		DifficultyBombDelaySchedule: parity.Uint64BigMapEncodesHex{
			uint64(0x42ae50): new(big.Int).SetUint64(uint64(0x2dc6c0)),
			uint64(0x6f1580): new(big.Int).SetUint64(uint64(0x1e8480)),
		},
		BlockRewardSchedule: parity.Uint64BigMapEncodesHex{
			uint64(0x0):      new(big.Int).SetUint64(uint64(0x4563918244f40000)),
			uint64(0x42ae50): new(big.Int).SetUint64(uint64(0x29a2241af62c0000)),
			uint64(0x6f1580): new(big.Int).SetUint64(uint64(0x1bc16d674ec80000)),
		},
	}

	// MainnetTrustedCheckpoint contains the light client trusted checkpoint for the main network.
	MainnetTrustedCheckpoint = &paramtypes.TrustedCheckpoint{
		SectionIndex: 270,
		SectionHead:  common.HexToHash("0xb67c33d838a60c282c2fb49b188fbbac1ef8565ffb4a1c4909b0a05885e72e40"),
		CHTRoot:      common.HexToHash("0x781daa4607782300da85d440df3813ba38a1262585231e35e9480726de81dbfc"),
		BloomRoot:    common.HexToHash("0xfd8951fa6d779cbc981df40dc31056ed1a549db529349d7dfae016f9d96cae72"),
	}

	// MainnetCheckpointOracle contains a set of configs for the main network oracle.
	MainnetCheckpointOracle = &paramtypes.CheckpointOracleConfig{
		Address: common.HexToAddress("0x9a9070028361F7AAbeB3f2F2Dc07F82C4a98A02a"),
		Signers: []common.Address{
			common.HexToAddress("0x1b2C260efc720BE89101890E4Db589b44E950527"), // Peter
			common.HexToAddress("0x78d1aD571A1A09D60D9BBf25894b44e4C8859595"), // Martin
			common.HexToAddress("0x286834935f4A8Cfb4FF4C77D5770C2775aE2b0E7"), // Zsolt
			common.HexToAddress("0xb86e2B0Ab5A4B1373e40c51A7C712c70Ba2f9f8E"), // Gary
			common.HexToAddress("0x0DF8fa387C602AE62559cC4aFa4972A7045d6707"), // Guillaume
		},
	}

	// TestnetChainConfig contains the chain parameters to run a node on the Ropsten test network.
	TestnetChainConfig = &paramtypes.ChainConfig{
		NetworkID:           3,
		ChainID:             big.NewInt(3),
		HomesteadBlock:      big.NewInt(0),
		DAOForkBlock:        nil,
		DAOForkSupport:      true,
		EIP150Block:         big.NewInt(0),
		EIP150Hash:          common.HexToHash("0x41941023680923e0fe4d74a34bdac8141f2540e3ae90623718e47d66d1ca4a2d"),
		EIP155Block:         big.NewInt(10),
		EIP158Block:         big.NewInt(10),
		ByzantiumBlock:      big.NewInt(1700000),
		DisposalBlock:       nil,
		SocialBlock:         nil,
		EthersocialBlock:    nil,
		ConstantinopleBlock: big.NewInt(4230000),
		PetersburgBlock:     big.NewInt(4939394),
		IstanbulBlock:       big.NewInt(6485846),
		Ethash:              new(paramtypes.EthashConfig),
		DifficultyBombDelaySchedule: parity.Uint64BigMapEncodesHex{
			1700000: new(big.Int).SetUint64(uint64(0x2dc6c0)),
			4230000: new(big.Int).SetUint64(uint64(0x1e8480)),
		},
		BlockRewardSchedule: parity.Uint64BigMapEncodesHex{
			uint64(0):       new(big.Int).SetUint64(uint64(0x4563918244f40000)),
			uint64(1700000): new(big.Int).SetUint64(uint64(0x29a2241af62c0000)),
			uint64(4230000): new(big.Int).SetUint64(uint64(0x1bc16d674ec80000)),
		},
	}

	// TestnetTrustedCheckpoint contains the light client trusted checkpoint for the Ropsten test network.
	TestnetTrustedCheckpoint = &paramtypes.TrustedCheckpoint{
		SectionIndex: 204,
		SectionHead:  common.HexToHash("0xa39168b51c3205456f30ce6a91f3590a43295b15a1c8c2ab86bb8c06b8ad1808"),
		CHTRoot:      common.HexToHash("0x9a3654147b79882bfc4e16fbd3421512aa7e4dfadc6c511923980e0877bdf3b4"),
		BloomRoot:    common.HexToHash("0xe72b979522d94fa45c1331639316da234a9bb85062d64d72e13afe1d3f5c17d5"),
	}

	// TestnetCheckpointOracle contains a set of configs for the Ropsten test network oracle.
	TestnetCheckpointOracle = &paramtypes.CheckpointOracleConfig{
		Address: common.HexToAddress("0xEF79475013f154E6A65b54cB2742867791bf0B84"),
		Signers: []common.Address{
			common.HexToAddress("0x32162F3581E88a5f62e8A61892B42C46E2c18f7b"), // Peter
			common.HexToAddress("0x78d1aD571A1A09D60D9BBf25894b44e4C8859595"), // Martin
			common.HexToAddress("0x286834935f4A8Cfb4FF4C77D5770C2775aE2b0E7"), // Zsolt
			common.HexToAddress("0xb86e2B0Ab5A4B1373e40c51A7C712c70Ba2f9f8E"), // Gary
			common.HexToAddress("0x0DF8fa387C602AE62559cC4aFa4972A7045d6707"), // Guillaume
		},
	}

	// RinkebyChainConfig contains the chain parameters to run a node on the Rinkeby test network.
	RinkebyChainConfig = &paramtypes.ChainConfig{
		NetworkID:           4,
		ChainID:             big.NewInt(4),
		HomesteadBlock:      big.NewInt(1),
		DAOForkBlock:        nil,
		DAOForkSupport:      true,
		EIP150Block:         big.NewInt(2),
		EIP150Hash:          common.HexToHash("0x9b095b36c15eaf13044373aef8ee0bd3a382a5abb92e402afa44b8249c3a90e9"),
		EIP155Block:         big.NewInt(3),
		EIP158Block:         big.NewInt(3),
		ByzantiumBlock:      big.NewInt(1035301),
		DisposalBlock:       nil,
		SocialBlock:         nil,
		EthersocialBlock:    nil,
		ConstantinopleBlock: big.NewInt(3660663),
		PetersburgBlock:     big.NewInt(4321234),
		IstanbulBlock:       big.NewInt(5435345),
		Clique: &paramtypes.CliqueConfig{
			Period: 15,
			Epoch:  30000,
		},
		DifficultyBombDelaySchedule: parity.Uint64BigMapEncodesHex{
			uint64(1035301): new(big.Int).SetUint64(uint64(0x2dc6c0)),
			uint64(3660663): new(big.Int).SetUint64(uint64(0x1e8480)),
		},
		BlockRewardSchedule: parity.Uint64BigMapEncodesHex{
			uint64(0x0):     new(big.Int).SetUint64(uint64(0x4563918244f40000)),
			uint64(1035301): new(big.Int).SetUint64(uint64(0x29a2241af62c0000)),
			uint64(3660663): new(big.Int).SetUint64(uint64(0x1bc16d674ec80000)),
		},
	}

	// RinkebyTrustedCheckpoint contains the light client trusted checkpoint for the Rinkeby test network.
	RinkebyTrustedCheckpoint = &paramtypes.TrustedCheckpoint{
		SectionIndex: 163,
		SectionHead:  common.HexToHash("0x36e5deaa46f258bece94b05d8e10f1ef68f422fb62ed47a2b6e616aa26e84997"),
		CHTRoot:      common.HexToHash("0x829b9feca1c2cdf5a4cf3efac554889e438ee4df8718c2ce3e02555a02d9e9e5"),
		BloomRoot:    common.HexToHash("0x58c01de24fdae7c082ebbe7665f189d0aa4d90ee10e72086bf56651c63269e54"),
	}

	// RinkebyCheckpointOracle contains a set of configs for the Rinkeby test network oracle.
	RinkebyCheckpointOracle = &paramtypes.CheckpointOracleConfig{
		Address: common.HexToAddress("0xebe8eFA441B9302A0d7eaECc277c09d20D684540"),
		Signers: []common.Address{
			common.HexToAddress("0xd9c9cd5f6779558b6e0ed4e6acf6b1947e7fa1f3"), // Peter
			common.HexToAddress("0x78d1aD571A1A09D60D9BBf25894b44e4C8859595"), // Martin
			common.HexToAddress("0x286834935f4A8Cfb4FF4C77D5770C2775aE2b0E7"), // Zsolt
			common.HexToAddress("0xb86e2B0Ab5A4B1373e40c51A7C712c70Ba2f9f8E"), // Gary
		},
	}

	// GoerliChainConfig contains the chain parameters to run a node on the Görli test network.
	GoerliChainConfig = &paramtypes.ChainConfig{
		NetworkID:           5,
		ChainID:             big.NewInt(5),
		HomesteadBlock:      big.NewInt(0),
		DAOForkBlock:        nil,
		DAOForkSupport:      true,
		EIP150Block:         big.NewInt(0),
		EIP155Block:         big.NewInt(0),
		EIP158Block:         big.NewInt(0),
		ByzantiumBlock:      big.NewInt(0),
		ConstantinopleBlock: big.NewInt(0),
		PetersburgBlock:     big.NewInt(0),
		IstanbulBlock:       big.NewInt(1561651),
		Clique: &paramtypes.CliqueConfig{
			Period: 15,
			Epoch:  30000,
		},
		DifficultyBombDelaySchedule: parity.Uint64BigMapEncodesHex{
			uint64(0x0): new(big.Int).SetUint64(uint64(0x1e8480)),
		},
		BlockRewardSchedule: parity.Uint64BigMapEncodesHex{
			uint64(0x0): new(big.Int).SetUint64(uint64(0x1bc16d674ec80000)),
		},
	}

	// GoerliTrustedCheckpoint contains the light client trusted checkpoint for the Görli test network.
	GoerliTrustedCheckpoint = &paramtypes.TrustedCheckpoint{
		SectionIndex: 47,
		SectionHead:  common.HexToHash("0x00c5b54c6c9a73660501fd9273ccdb4c5bbdbe5d7b8b650e28f881ec9d2337f6"),
		CHTRoot:      common.HexToHash("0xef35caa155fd659f57167e7d507de2f8132cbb31f771526481211d8a977d704c"),
		BloomRoot:    common.HexToHash("0xbda330402f66008d52e7adc748da28535b1212a7912a21244acd2ba77ff0ff06"),
	}

	// GoerliCheckpointOracle contains a set of configs for the Goerli test network oracle.
	GoerliCheckpointOracle = &paramtypes.CheckpointOracleConfig{
		Address: common.HexToAddress("0x18CA0E045F0D772a851BC7e48357Bcaab0a0795D"),
		Signers: []common.Address{
			common.HexToAddress("0x4769bcaD07e3b938B7f43EB7D278Bc7Cb9efFb38"), // Peter
			common.HexToAddress("0x78d1aD571A1A09D60D9BBf25894b44e4C8859595"), // Martin
			common.HexToAddress("0x286834935f4A8Cfb4FF4C77D5770C2775aE2b0E7"), // Zsolt
			common.HexToAddress("0xb86e2B0Ab5A4B1373e40c51A7C712c70Ba2f9f8E"), // Gary
			common.HexToAddress("0x0DF8fa387C602AE62559cC4aFa4972A7045d6707"), // Guillaume
		},
	}

	// AllEthashProtocolChanges contains every protocol change (EIPs) introduced
	// and accepted by the Ethereum core developers into the Ethash consensus.
	//
	// This configuration is intentionally not using keyed fields to force anyone
	// adding flags to the config to also have to set these fields.
	AllEthashProtocolChanges = &paramtypes.ChainConfig{
		1337,
		big.NewInt(1337), // ChainID

		big.NewInt(0), // HomesteadBlock
		nil,           // EIP2FBlock
		nil,           // EIP7FBlock

		nil,   // DAOForkBlock
		false, // DAOForkSupport

		big.NewInt(0), // EIP150Block
		common.Hash{}, // EIP150Hash
		big.NewInt(0), // EIP155Block
		big.NewInt(0), // EIP158Block
		nil,           // EIP160FBlock
		nil,           // EIP161FBlock
		nil,           // EIP170FBlock

		big.NewInt(0), // ByzantiumBlock
		nil,           // EIP100FBlock
		nil,           // EIP140FBlock
		nil,           // EIP198FBlock
		nil,           // EIP211FBlock
		nil,           // EIP212FBlock
		nil,           // EIP213FBlock
		nil,           // EIP214FBlock
		nil,           // EIP649FBlock
		nil,           // EIP658FBlock

		big.NewInt(0), // ConstantinopleBlock
		nil,           // EIP145FBlock
		nil,           // EIP1014FBlock
		nil,           // EIP1052FBlock
		nil,           // EIP1234FBlock
		nil,           // EIP1283FBlock

		big.NewInt(0), // PetersburgBlock
		big.NewInt(0), // IstanbulBlock

		nil, // EIP152FBlock
		nil, // EIP1108FBlock
		nil, // EIP1344FBlock
		nil, // EIP1884FBlock
		nil, // EIP2028FBlock
		nil, // EIP2200FBlock
		nil, // EWASMBlock

		nil, // ECIP1010PauseBlock
		nil, // ECIP1010Length
		nil, // ECIP1017FBlock
		nil, // ECIP1017EraRounds
		nil, // DisposalBlock
		nil, // SocialBlock
		nil, // EthersocialBlock

		nil, // Musicoin MCIP0Block UBI
		nil, // Musicoin MCIP3Block UBI
		nil, // Musicoin MCIP8Block QT

		new(paramtypes.EthashConfig), // Ethash
		nil,                          // Clique
		nil,
		nil,
		// DifficultyBombDelaySchedule
		parity.Uint64BigMapEncodesHex{
			0: new(big.Int).SetUint64(uint64(0x1e8480)),
		},
		// BlockRewardSchedule
		parity.Uint64BigMapEncodesHex{
			0: new(big.Int).SetUint64(uint64(0x1bc16d674ec80000)),
		},
	}

	// AllCliqueProtocolChanges contains every protocol change (EIPs) introduced
	// and accepted by the Ethereum core developers into the Clique consensus.
	//
	// This configuration is intentionally not using keyed fields to force anyone
	// adding flags to the config to also have to set these fields.
	AllCliqueProtocolChanges = &paramtypes.ChainConfig{
		5,
		big.NewInt(1337), // ChainID

		big.NewInt(0), // HomesteadBlock
		nil,           // EIP2FBlock
		nil,           // EIP7FBlock

		nil,   // DAOForkBlock
		false, // DAOForkSupport

		big.NewInt(0), // EIP150Block
		common.Hash{}, // EIP150Hash
		big.NewInt(0), // EIP155Block
		big.NewInt(0), // EIP158Block
		nil,           // EIP160FBlock
		nil,           // EIP161FBlock
		nil,           // EIP170FBlock

		big.NewInt(0), // ByzantiumBlock
		nil,           // EIP100FBlock
		nil,           // EIP140FBlock
		nil,           // EIP198FBlock
		nil,           // EIP211FBlock
		nil,           // EIP212FBlock
		nil,           // EIP213FBlock
		nil,           // EIP214FBlock
		nil,           // EIP649FBlock
		nil,           // EIP658FBlock

		big.NewInt(0), // ConstantinopleBlock
		nil,           // EIP145FBlock
		nil,           // EIP1014FBlock
		nil,           // EIP1052FBlock
		nil,           // EIP1234FBlock
		nil,           // EIP1283FBlock

		nil, // PetersburgBlock
		nil, // IstanbulBlock
		nil, // EIP152FBlock
		nil, // EIP1108FBlock
		nil, // EIP1344FBlock
		nil, // EIP1884FBlock
		nil, // EIP2028FBlock
		nil, // EIP2200FBlock
		nil, // EWASMBlock

		nil, // ECIP1010PauseBlock
		nil, // ECIP1010Length
		nil, // ECIP1017FBlock
		nil, // ECIP1017EraRounds
		nil, // DisposalBlock
		nil, // SocialBlock
		nil, // EthersocialBlock

		nil, // Musicoin MCIP0Block UBI
		nil, // Musicoin MCIP3Block UBI
		nil, // Musicoin MCIP8Block QT

		nil, // Ethash
		&paramtypes.CliqueConfig{
			Period: 0,
			Epoch:  30000,
		},
		nil,
		nil,
		// DifficultyBombDelaySchedule
		parity.Uint64BigMapEncodesHex{
			0: new(big.Int).SetUint64(uint64(0x1e8480)),
		},
		// BlockRewardSchedule
		parity.Uint64BigMapEncodesHex{
			0: new(big.Int).SetUint64(uint64(0x1bc16d674ec80000)),
		},
	}

	// TestChainConfig is used for tests.
	TestChainConfig = &paramtypes.ChainConfig{
		3,
		big.NewInt(1), // ChainID

		big.NewInt(0), // HomesteadBlock
		nil,           // EIP2FBlock
		nil,           // EIP7FBlock

		nil,   // DAOForkBlock
		false, // DAOForkSupport

		big.NewInt(0), // EIP150Block
		common.Hash{}, // EIP150Hash
		big.NewInt(0), // EIP155Block
		big.NewInt(0), // EIP158Block
		nil,           // EIP160FBlock
		nil,           // EIP161FBlock
		nil,           // EIP170FBlock

		big.NewInt(0), // ByzantiumBlock
		nil,           // EIP100FBlock
		nil,           // EIP140FBlock
		nil,           // EIP198FBlock
		nil,           // EIP211FBlock
		nil,           // EIP212FBlock
		nil,           // EIP213FBlock
		nil,           // EIP214FBlock
		nil,           // EIP649FBlock
		nil,           // EIP658FBlock

		big.NewInt(0), // ConstantinopleBlock
		nil,           // EIP145FBlock
		nil,           // EIP1014FBlock
		nil,           // EIP1052FBlock
		nil,           // EIP1234FBlock
		nil,           // EIP1283FBlock

		nil, // PetersburgBlock
		nil, // IstanbulBlock
		nil, // EIP152FBlock
		nil, // EIP1108FBlock
		nil, // EIP1344FBlock
		nil, // EIP1884FBlock
		nil, // EIP2028FBlock
		nil, // EIP2200FBlock
		nil, // EWASMBlock

		nil, // ECIP1010PauseBlock
		nil, // ECIP1010Length
		nil, // ECIP1017FBlock
		nil, // ECIP1017EraRounds
		nil, // DisposalBlock
		nil, // SocialBlock
		nil, // EthersocialBlock

		nil, // Musicoin MCIP0Block UBI
		nil, // Musicoin MCIP3Block UBI
		nil, // Musicoin MCIP8Block QT

		new(paramtypes.EthashConfig), // Ethash
		nil,                          // Clique
		nil,
		nil,
		// DifficultyBombDelaySchedule
		parity.Uint64BigMapEncodesHex{
			0: new(big.Int).SetUint64(uint64(0x1e8480)),
		},
		// BlockRewardSchedule
		parity.Uint64BigMapEncodesHex{
			0: new(big.Int).SetUint64(uint64(0x1bc16d674ec80000)),
		},
	}

	// TestRules are all rules from TestChainConfig initialized at 0.
	TestRules = TestChainConfig.Rules(new(big.Int))
)


func EthashBlockReward(c *paramtypes.ChainConfig, n *big.Int) *big.Int {
	// if c.Ethash == nil {
	// 	panic("non ethash config called EthashBlockReward")
	// }
	// Select the correct block reward based on chain progression
	blockReward := FrontierBlockReward
	if c == nil || n == nil {
		return blockReward
	}
	// Because the map is not necessarily sorted low-high, we
	// have to ensure that we're walking upwards only.
	var lastActivation *big.Int
	for activation, reward := range c.BlockRewardSchedule {
		activationBig := big.NewInt(int64(activation))
		if paramtypes.IsForked(activationBig, n) {
			if lastActivation == nil {
				lastActivation = new(big.Int).Set(activationBig)
			}
			if activationBig.Cmp(lastActivation) >= 0 {
				blockReward = reward
			}
		}
	}
	return blockReward
}
