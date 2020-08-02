// Copyright 2017 The go-ethereum Authors
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
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params/confp"
	"github.com/ethereum/go-ethereum/params/types/coregeth"
	"github.com/ethereum/go-ethereum/params/types/ctypes"
	"github.com/ethereum/go-ethereum/params/types/goethereum"
	"github.com/ethereum/go-ethereum/params/types/multigethv0"
)

func uint64P(n uint64) *uint64 {
	return &n
}

func TestCheckCompatible(t *testing.T) {
	type test struct {
		stored, new ctypes.ChainConfigurator
		head        uint64
		wantErr     *confp.ConfigCompatError
	}
	tests := []test{
		{stored: AllEthashProtocolChanges, new: AllEthashProtocolChanges, head: 0, wantErr: nil},
		{stored: AllEthashProtocolChanges, new: AllEthashProtocolChanges, head: 100, wantErr: nil},
		{
			stored:  &goethereum.ChainConfig{Ethash: new(ctypes.EthashConfig), EIP150Block: big.NewInt(10)},
			new:     &goethereum.ChainConfig{Ethash: new(ctypes.EthashConfig), EIP150Block: big.NewInt(20)},
			head:    9,
			wantErr: nil,
		},
		{
			stored: AllEthashProtocolChanges,
			new:    &goethereum.ChainConfig{Ethash: new(ctypes.EthashConfig), HomesteadBlock: nil},
			head:   3,
			wantErr: &confp.ConfigCompatError{
				What:         "Homestead fork block",
				StoredConfig: uint64P(0),
				NewConfig:    nil,
				RewindTo:     0,
			},
		},
		{
			stored: AllEthashProtocolChanges,
			new:    &goethereum.ChainConfig{Ethash: new(ctypes.EthashConfig), HomesteadBlock: big.NewInt(1)},
			head:   3,
			wantErr: &confp.ConfigCompatError{
				What:         "Homestead fork block",
				StoredConfig: uint64P(0),
				NewConfig:    uint64P(1),
				RewindTo:     0,
			},
		},
		{
			stored: &goethereum.ChainConfig{Ethash: new(ctypes.EthashConfig), HomesteadBlock: big.NewInt(30), EIP150Block: big.NewInt(10)},
			new:    &goethereum.ChainConfig{Ethash: new(ctypes.EthashConfig), HomesteadBlock: big.NewInt(25), EIP150Block: big.NewInt(20)},
			head:   25,
			wantErr: &confp.ConfigCompatError{
				What:         "EIP150 fork block",
				StoredConfig: uint64P(10),
				NewConfig:    uint64P(20),
				RewindTo:     9,
			},
		},
		{
			stored: &coregeth.CoreGethChainConfig{Ethash: new(ctypes.EthashConfig), EIP100FBlock: big.NewInt(30), EIP649FBlock: big.NewInt(30)},
			new:    &coregeth.CoreGethChainConfig{Ethash: new(ctypes.EthashConfig), EIP100FBlock: big.NewInt(24), EIP649FBlock: big.NewInt(24)},
			head:   25,
			wantErr: &confp.ConfigCompatError{
				What:         "EIP100F fork block",
				StoredConfig: uint64P(30),
				NewConfig:    uint64P(24),
				RewindTo:     23,
			},
		},
		{
			stored:  &goethereum.ChainConfig{Ethash: new(ctypes.EthashConfig), ByzantiumBlock: big.NewInt(30)},
			new:     &coregeth.CoreGethChainConfig{Ethash: new(ctypes.EthashConfig), EIP211FBlock: big.NewInt(26)},
			head:    25,
			wantErr: nil,
		},
		{
			stored:  &goethereum.ChainConfig{Ethash: new(ctypes.EthashConfig), ByzantiumBlock: big.NewInt(30)},
			new:     &coregeth.CoreGethChainConfig{Ethash: new(ctypes.EthashConfig), EIP100FBlock: big.NewInt(26), EIP649FBlock: big.NewInt(26)},
			head:    25,
			wantErr: nil,
		},
		{
			stored: MainnetChainConfig,
			new: func() ctypes.ChainConfigurator {
				c := &goethereum.ChainConfig{}
				confp.Convert(MainnetChainConfig, c)
				c.SetEthashEIP779Transition(uint64P(1900000))
				return c
			}(),
			head: MainnetChainConfig.DAOForkBlock.Uint64(),
			wantErr: &confp.ConfigCompatError{
				What:         "DAO fork support flag",
				StoredConfig: uint64P(MainnetChainConfig.DAOForkBlock.Uint64()),
				NewConfig:    uint64P(1900000),
				RewindTo:     1900000 - 1,
			},
		},
		{
			stored: MainnetChainConfig,
			new: func() ctypes.ChainConfigurator {
				c := &goethereum.ChainConfig{}
				confp.Convert(MainnetChainConfig, c)
				c.SetEthashEIP779Transition(nil)
				return c
			}(),
			head: MainnetChainConfig.DAOForkBlock.Uint64(),
			wantErr: &confp.ConfigCompatError{
				What:         "DAO fork support flag",
				StoredConfig: uint64P(MainnetChainConfig.DAOForkBlock.Uint64()),
				NewConfig:    nil,
				RewindTo:     1920000 - 1,
			},
		},
		{
			stored: MainnetChainConfig,
			new: func() ctypes.ChainConfigurator {
				c := &goethereum.ChainConfig{}
				*c = *MainnetChainConfig
				c.SetChainID(new(big.Int).Sub(MainnetChainConfig.EIP155Block, common.Big1))
				return c
			}(),
			head: MainnetChainConfig.EIP158Block.Uint64(),
			wantErr: &confp.ConfigCompatError{
				What:         "EIP155 chain ID",
				StoredConfig: uint64P(MainnetChainConfig.EIP155Block.Uint64()),
				NewConfig:    uint64P(MainnetChainConfig.EIP155Block.Uint64()),
				RewindTo:     new(big.Int).Sub(MainnetChainConfig.EIP158Block, common.Big1).Uint64(),
			},
		},
		{
			stored: func() ctypes.ChainConfigurator {
				c := &goethereum.ChainConfig{
					Ethash:         new(ctypes.EthashConfig),
					DAOForkBlock:   big.NewInt(3),
					DAOForkSupport: false,
				}
				return c
			}(),
			new: func() ctypes.ChainConfigurator {
				c := &coregeth.CoreGethChainConfig{
					Ethash:       new(ctypes.EthashConfig),
					DAOForkBlock: nil,
				}
				return c
			}(),
			head:    5,
			wantErr: nil,
		},
		{
			// v1.9.5 -> v1.9.7
			stored: func() ctypes.ChainConfigurator {
				c := &coregeth.CoreGethChainConfig{}
				*c = *ClassicChainConfig
				c.SetEIP145Transition(nil)
				c.SetEIP1014Transition(nil)
				c.SetEIP1052Transition(nil)
				c.SetEIP152Transition(nil)
				c.SetEIP1108Transition(nil)
				c.SetEIP1344Transition(nil)
				//c.SetEIP1884Transition(nil)
				c.SetEIP2028Transition(nil)
				c.SetEIP2200Transition(nil)
				return c
			}(),
			new:     ClassicChainConfig,
			head:    9550000,
			wantErr: nil,
		},
		{
			// v1.9.6 -> v1.9.7
			stored: func() ctypes.ChainConfigurator {
				c := &coregeth.CoreGethChainConfig{}
				*c = *ClassicChainConfig
				c.SetEIP152Transition(nil)
				c.SetEIP1108Transition(nil)
				c.SetEIP1344Transition(nil)
				//c.SetEIP1884Transition(nil)
				c.SetEIP2028Transition(nil)
				c.SetEIP2200Transition(nil)
				return c
			}(),
			new:     ClassicChainConfig,
			head:    9550000,
			wantErr: nil,
		},
		{
			stored: MainnetChainConfig,
			new: func() ctypes.ChainConfigurator {
				c := &coregeth.CoreGethChainConfig{}
				err := confp.Convert(MainnetChainConfig, c)
				if err != nil {
					panic(err)
				}
				return c
			}(),
		},
		{
			stored: func() ctypes.ChainConfigurator {
				// ClassicChainConfig is the chain parameters to run a node on the Classic main network.
				c := &multigethv0.ChainConfig{
					ChainID:             big.NewInt(61),
					HomesteadBlock:      big.NewInt(1150000),
					DAOForkBlock:        big.NewInt(1920000),
					DAOForkSupport:      false,
					EIP150Block:         big.NewInt(2500000),
					EIP150Hash:          common.HexToHash("0xca12c63534f565899681965528d536c52cb05b7c48e269c2a6cb77ad864d878a"),
					EIP155Block:         big.NewInt(3000000),
					EIP158Block:         big.NewInt(8772000),
					ByzantiumBlock:      big.NewInt(8772000),
					DisposalBlock:       big.NewInt(5900000),
					ConstantinopleBlock: big.NewInt(9573000),
					PetersburgBlock:     big.NewInt(9573000),
					// As if client hasn't upgraded config to latest fork.
					//IstanbulBlock:       big.NewInt(10500839),
					//EIP1884DisableFBlock:big.NewInt(10500839),
					ECIP1017EraBlock:   big.NewInt(5000000),
					EIP160Block:        big.NewInt(3000000),
					ECIP1010PauseBlock: big.NewInt(3000000),
					ECIP1010Length:     big.NewInt(2000000),
					Ethash:             new(ctypes.EthashConfig),
				}
				return c
			}(),
			new:     ClassicChainConfig,
			head:    9700000,
			wantErr: nil,
		},
		{
			stored: func() ctypes.ChainConfigurator {
				// ClassicChainConfig is the chain parameters to run a node on the Classic main network.
				c := &multigethv0.ChainConfig{
					ChainID:             big.NewInt(61),
					HomesteadBlock:      big.NewInt(1150000),
					DAOForkBlock:        big.NewInt(1920000),
					DAOForkSupport:      false,
					EIP150Block:         big.NewInt(2500000),
					EIP150Hash:          common.HexToHash("0xca12c63534f565899681965528d536c52cb05b7c48e269c2a6cb77ad864d878a"),
					EIP155Block:         big.NewInt(3000000),
					EIP158Block:         big.NewInt(8772000),
					ByzantiumBlock:      big.NewInt(8772000),
					DisposalBlock:       big.NewInt(5900000),
					ConstantinopleBlock: big.NewInt(9573000),
					PetersburgBlock:     big.NewInt(9573000),
					IstanbulBlock:       big.NewInt(10500839),
					ECIP1017EraBlock:    big.NewInt(5000000),
					EIP160Block:         big.NewInt(3000000),
					ECIP1010PauseBlock:  big.NewInt(3000000),
					ECIP1010Length:      big.NewInt(2000000),
					Ethash:              new(ctypes.EthashConfig),
				}
				return c
			}(),
			new:     ClassicChainConfig,
			head:    9700000,
			wantErr: nil,
		},
	}

	for _, test := range tests {
		err := confp.Compatible(&test.head, test.stored, test.new)
		if (err == nil && test.wantErr != nil) || (err != nil && test.wantErr == nil) {
			t.Errorf("nil/nonnil, error mismatch:\nstored: %v\nnew: %v\nhead: %v\nerr: %v\nwant: %v", test.stored, test.new, test.head, err, test.wantErr)
		} else if err != nil && (err.RewindTo != test.wantErr.RewindTo) {
			//if !reflect.DeepEqual(err, test.wantErr) {
			t.Errorf("error mismatch:\nstored: %v\nnew: %v\nhead: %v\nerr: %v\nwant: %v", test.stored, test.new, test.head, err, test.wantErr)
		}
	}
}

func TestFoundationIsForked(t *testing.T) {
	c := MainnetChainConfig
	if !c.IsEnabled(c.GetEthashEIP2384Transition, big.NewInt(9200001)) {
		t.Fatal("nofork muir bad")
	}
}

func TestClassicIs649(t *testing.T) {
	c := ClassicChainConfig
	got := c.GetEthashEIP649Transition()
	if got != nil {
		t.Fatal("classic config doesn't support 649; difficulty bomb was disposed of")
	}
}

func TestFoundationIsEIP779(t *testing.T) {
	blockNumbers := []*big.Int{
		big.NewInt(0),
		big.NewInt(1_920_000),
		big.NewInt(10_000_000),
	}
	for _, bn := range blockNumbers {
		if bn.Cmp(big.NewInt(0)) > 0 && !MainnetChainConfig.IsEnabled(MainnetChainConfig.GetEthashEIP779Transition, bn) {
			t.Fatal("bad")
		}
		if *MainnetChainConfig.GetEthashEIP779Transition() != 1_920_000 {
			t.Fatal("bad")
		}
	}
}
