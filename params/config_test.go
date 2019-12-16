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
	"github.com/ethereum/go-ethereum/params/convert"
	"github.com/ethereum/go-ethereum/params/types"
	"github.com/ethereum/go-ethereum/params/types/ctypes"
	"github.com/ethereum/go-ethereum/params/types/goethereum"
)

func uint64P(n uint64) *uint64 {
	return &n
}

func TestCheckCompatible(t *testing.T) {
	type test struct {
		stored, new ctypes.ChainConfigurator
		head        uint64
		wantErr     *ctypes.ConfigCompatError
	}
	tests := []test{
		{stored: AllEthashProtocolChanges, new: AllEthashProtocolChanges, head: 0, wantErr: nil},
		{stored: AllEthashProtocolChanges, new: AllEthashProtocolChanges, head: 100, wantErr: nil},
		{
			stored:  &goethereum.ChainConfig{EIP150Block: big.NewInt(10)},
			new:     &goethereum.ChainConfig{EIP150Block: big.NewInt(20)},
			head:    9,
			wantErr: nil,
		},
		{
			stored: AllEthashProtocolChanges,
			new:    &goethereum.ChainConfig{HomesteadBlock: nil},
			head:   3,
			wantErr: &ctypes.ConfigCompatError{
				What:         "Homestead fork block",
				StoredConfig: uint64P(0),
				NewConfig:    nil,
				RewindTo:     0,
			},
		},
		{
			stored: AllEthashProtocolChanges,
			new:    &goethereum.ChainConfig{HomesteadBlock: big.NewInt(1)},
			head:   3,
			wantErr: &ctypes.ConfigCompatError{
				What:         "Homestead fork block",
				StoredConfig: uint64P(0),
				NewConfig:    uint64P(1),
				RewindTo:     0,
			},
		},
		{
			stored: &goethereum.ChainConfig{HomesteadBlock: big.NewInt(30), EIP150Block: big.NewInt(10)},
			new:    &goethereum.ChainConfig{HomesteadBlock: big.NewInt(25), EIP150Block: big.NewInt(20)},
			head:   25,
			wantErr: &ctypes.ConfigCompatError{
				What:         "EIP150 fork block",
				StoredConfig: uint64P(10),
				NewConfig:    uint64P(20),
				RewindTo:     9,
			},
		},
		{
			stored: &paramtypes.MultiGethChainConfig{EIP100FBlock: big.NewInt(30), EIP649FBlock: big.NewInt(30)},
			new:    &paramtypes.MultiGethChainConfig{EIP100FBlock: big.NewInt(24), EIP649FBlock: big.NewInt(24)},
			head:   25,
			wantErr: &ctypes.ConfigCompatError{
				What:         "EIP100F fork block",
				StoredConfig: uint64P(30),
				NewConfig:    uint64P(24),
				RewindTo:     23,
			},
		},
		{
			stored:  &goethereum.ChainConfig{ByzantiumBlock: big.NewInt(30)},
			new:     &paramtypes.MultiGethChainConfig{EIP211FBlock: big.NewInt(26)},
			head:    25,
			wantErr: nil,
		},
		{
			stored:  &goethereum.ChainConfig{ByzantiumBlock: big.NewInt(30)},
			new:     &paramtypes.MultiGethChainConfig{EIP100FBlock: big.NewInt(26), EIP649FBlock: big.NewInt(26)},
			head:    25,
			wantErr: nil,
		},
		{
			stored: MainnetChainConfig,
			new: func() ctypes.ChainConfigurator {
				c := &goethereum.ChainConfig{}
				convert.Convert(MainnetChainConfig, c)
				c.SetEthashEIP779Transition(uint64P(1900000))
				return c
			}(),
			head: MainnetChainConfig.DAOForkBlock.Uint64(),
			wantErr: &ctypes.ConfigCompatError{
				What:         "DAO fork support flag",
				StoredConfig: uint64P(MainnetChainConfig.DAOForkBlock.Uint64()),
				NewConfig:    uint64P(1900000),
				RewindTo:     1900000-1,
			},
		},
		{
			stored: MainnetChainConfig,
			new: func() ctypes.ChainConfigurator {
				c := &goethereum.ChainConfig{}
				convert.Convert(MainnetChainConfig, c)
				c.SetEthashEIP779Transition(nil)
				return c
			}(),
			head: MainnetChainConfig.DAOForkBlock.Uint64(),
			wantErr: &ctypes.ConfigCompatError{
				What:         "DAO fork support flag",
				StoredConfig: uint64P(MainnetChainConfig.DAOForkBlock.Uint64()),
				NewConfig:    nil,
				RewindTo:     1920000-1,
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
			wantErr: &ctypes.ConfigCompatError{
				What:         "EIP155 chain ID",
				StoredConfig: uint64P(MainnetChainConfig.EIP155Block.Uint64()),
				NewConfig:    uint64P(MainnetChainConfig.EIP155Block.Uint64()),
				RewindTo:     new(big.Int).Sub(MainnetChainConfig.EIP158Block, common.Big1).Uint64(),
			},
		},
	}

	for _, test := range tests {
		err := ctypes.Compatible(&test.head, test.stored, test.new)
		if (err == nil && test.wantErr != nil) || (err != nil && test.wantErr == nil) {
			t.Errorf("nil/nonnil, error mismatch:\nstored: %v\nnew: %v\nhead: %v\nerr: %v\nwant: %v", test.stored, test.new, test.head, err, test.wantErr)
		} else if err != nil && (err.RewindTo != test.wantErr.RewindTo) {
		//if !reflect.DeepEqual(err, test.wantErr) {
			t.Errorf("error mismatch:\nstored: %v\nnew: %v\nhead: %v\nerr: %v\nwant: %v", test.stored, test.new, test.head, err, test.wantErr)
		}
	}
}
