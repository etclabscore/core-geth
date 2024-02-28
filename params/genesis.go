// Copyright 2019 The multi-geth Authors
// This file is part of the multi-geth library.
//
// The multi-geth library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The multi-geth library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the multi-geth library. If not, see <http://www.gnu.org/licenses/>.

package params

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/params/types/coregeth"
	"github.com/ethereum/go-ethereum/params/types/ctypes"
	"github.com/ethereum/go-ethereum/params/types/genesisT"
	"github.com/ethereum/go-ethereum/params/types/goethereum"
	"github.com/ethereum/go-ethereum/params/vars"
)

// DefaultGenesisBlock returns the Ethereum main net genesis block.
func DefaultGenesisBlock() *genesisT.Genesis {
	return &genesisT.Genesis{
		Config:     MainnetChainConfig,
		Nonce:      66,
		ExtraData:  hexutil.MustDecode("0x11bbe8db4e347b4e8c937c1c8370e4b5ed33adb3db69cbdb7a38e1e50b1b82fa"),
		GasLimit:   5000,
		Difficulty: big.NewInt(17179869184),
		Alloc:      genesisT.DecodePreAlloc(mainnetAllocData),
	}
}

// DefaultSepoliaGenesisBlock returns the Sepolia network genesis block.
func DefaultSepoliaGenesisBlock() *genesisT.Genesis {
	return &genesisT.Genesis{
		Config:     SepoliaChainConfig,
		Nonce:      0,
		ExtraData:  []byte("Sepolia, Athens, Attica, Greece!"),
		GasLimit:   30000000,
		Difficulty: big.NewInt(131072),
		Timestamp:  1633267481,
		Alloc:      genesisT.DecodePreAlloc(sepoliaAllocData),
	}
}

// DefaultGoerliGenesisBlock returns the Görli network genesis block.
func DefaultGoerliGenesisBlock() *genesisT.Genesis {
	return &genesisT.Genesis{
		Config:     GoerliChainConfig,
		Timestamp:  1548854791,
		ExtraData:  hexutil.MustDecode("0x22466c6578692069732061207468696e6722202d204166726900000000000000e0a2bd4258d2768837baa26a28fe71dc079f84c70000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
		GasLimit:   10485760,
		Difficulty: big.NewInt(1),
		Alloc:      genesisT.DecodePreAlloc(goerliAllocData),
	}
}

// DefaultHoleskyGenesisBlock returns the Holesky network genesis block.
func DefaultHoleskyGenesisBlock() *genesisT.Genesis {
	return &genesisT.Genesis{
		Config:     HoleskyChainConfig,
		Nonce:      0x1234,
		GasLimit:   0x17d7840,
		Difficulty: big.NewInt(0x01),
		Timestamp:  1695902100,
		Alloc:      genesisT.DecodePreAlloc(holeskyAllocData),
	}
}

// DeveloperGenesisBlock returns the 'geth --dev' genesis block. Note, this must
// be seeded with the
func DeveloperGenesisBlock(gasLimit uint64, faucet *common.Address, useEthash bool) *genesisT.Genesis {
	if !useEthash {
		// Make a copy to avoid unpredicted contamination.
		config := &goethereum.ChainConfig{}
		*config = *AllDevChainProtocolChanges

		// Assemble and return the genesis with the precompiles and faucet pre-funded
		genesis := &genesisT.Genesis{
			Config:     config,
			GasLimit:   gasLimit,
			BaseFee:    big.NewInt(vars.InitialBaseFee),
			Difficulty: big.NewInt(0),
			Alloc: map[common.Address]genesisT.GenesisAccount{
				common.BytesToAddress([]byte{1}): {Balance: big.NewInt(1)}, // ECRecover
				common.BytesToAddress([]byte{2}): {Balance: big.NewInt(1)}, // SHA256
				common.BytesToAddress([]byte{3}): {Balance: big.NewInt(1)}, // RIPEMD
				common.BytesToAddress([]byte{4}): {Balance: big.NewInt(1)}, // Identity
				common.BytesToAddress([]byte{5}): {Balance: big.NewInt(1)}, // ModExp
				common.BytesToAddress([]byte{6}): {Balance: big.NewInt(1)}, // ECAdd
				common.BytesToAddress([]byte{7}): {Balance: big.NewInt(1)}, // ECScalarMul
				common.BytesToAddress([]byte{8}): {Balance: big.NewInt(1)}, // ECPairing
				common.BytesToAddress([]byte{9}): {Balance: big.NewInt(1)}, // BLAKE2b
			},
		}
		if faucet != nil {
			genesis.Alloc[*faucet] = genesisT.GenesisAccount{Balance: new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 256), big.NewInt(9))}
		}
		return genesis
	}

	// Use an ETC equivalent of AllEthashProtocolChanges.
	// This will allow initial permanent disposal of the difficulty bomb,
	// and we'll override the monetary policy block reward schedule to be a non-occurring.
	//
	// This was originally intended to be as follows, but import cycles prevent it.
	// Leaving here to show provenance of initial configuration value.
	// config := &coregeth.CoreGethChainConfig{}
	// *config = *tests.Forks["ETC_Phoenix"].(*coregeth.CoreGethChainConfig)
	config := &coregeth.CoreGethChainConfig{
		NetworkID:          AllCliqueProtocolChanges.GetChainID().Uint64(), // Use network and chain IDs equivalent to Clique configuration, ie 1337.
		Ethash:             new(ctypes.EthashConfig),
		ChainID:            AllCliqueProtocolChanges.GetChainID(),
		EIP2FBlock:         big.NewInt(0),
		EIP7FBlock:         big.NewInt(0),
		EIP150Block:        big.NewInt(0),
		EIP155Block:        big.NewInt(0),
		EIP160FBlock:       big.NewInt(0),
		EIP161FBlock:       big.NewInt(0),
		EIP170FBlock:       big.NewInt(0),
		EIP100FBlock:       big.NewInt(0),
		EIP140FBlock:       big.NewInt(0),
		EIP198FBlock:       big.NewInt(0),
		EIP211FBlock:       big.NewInt(0),
		EIP212FBlock:       big.NewInt(0),
		EIP213FBlock:       big.NewInt(0),
		EIP214FBlock:       big.NewInt(0),
		EIP658FBlock:       big.NewInt(0),
		EIP145FBlock:       big.NewInt(0),
		EIP1014FBlock:      big.NewInt(0),
		EIP1052FBlock:      big.NewInt(0),
		EIP1283FBlock:      big.NewInt(0),
		PetersburgBlock:    big.NewInt(0),
		EIP152FBlock:       big.NewInt(0),
		EIP1108FBlock:      big.NewInt(0),
		EIP1344FBlock:      big.NewInt(0),
		EIP1884FBlock:      big.NewInt(0),
		EIP2028FBlock:      big.NewInt(0),
		EIP2200FBlock:      big.NewInt(0),
		DisposalBlock:      big.NewInt(0),
		ECIP1017FBlock:     nil, // disable block reward disinflation
		ECIP1017EraRounds:  nil, // ^
		ECIP1010PauseBlock: nil, // no need for difficulty bomb delay (see disposal block)
		ECIP1010Length:     nil, // ^
	}

	// Assemble and return the genesis with the precompiles and faucet pre-funded
	genesis := &genesisT.Genesis{
		Config:     config,
		GasLimit:   6283185,
		Difficulty: vars.MinimumDifficulty,
		BaseFee:    big.NewInt(vars.InitialBaseFee),
		Alloc: map[common.Address]genesisT.GenesisAccount{
			common.BytesToAddress([]byte{1}): {Balance: big.NewInt(1)}, // ECRecover
			common.BytesToAddress([]byte{2}): {Balance: big.NewInt(1)}, // SHA256
			common.BytesToAddress([]byte{3}): {Balance: big.NewInt(1)}, // RIPEMD
			common.BytesToAddress([]byte{4}): {Balance: big.NewInt(1)}, // Identity
			common.BytesToAddress([]byte{5}): {Balance: big.NewInt(1)}, // ModExp
			common.BytesToAddress([]byte{6}): {Balance: big.NewInt(1)}, // ECAdd
			common.BytesToAddress([]byte{7}): {Balance: big.NewInt(1)}, // ECScalarMul
			common.BytesToAddress([]byte{8}): {Balance: big.NewInt(1)}, // ECPairing
		},
	}
	if faucet != nil {
		genesis.Alloc[*faucet] = genesisT.GenesisAccount{Balance: new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 256), big.NewInt(9))}
	}
	return genesis
}
