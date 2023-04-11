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

package generic

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/params/types/coregeth"
	"github.com/ethereum/go-ethereum/params/types/ctypes"
	"github.com/ethereum/go-ethereum/params/types/goethereum"
	"github.com/tidwall/gjson"
)

// GenericCC is a generic-y struct type used to expose some meta-logic methods
// shared by all ChainConfigurator implementations but not existing in that interface.
// These logics differentiate from the logics present in the ChainConfigurator interface
// itself because they are chain-aware, or fit nuanced, or adhoc, use cases and should
// not be demanded of EVM-based ecosystem logic as a whole. Debatable. NOTE.
type GenericCC struct {
	ctypes.ChainConfigurator
}

func AsGenericCC(c ctypes.ChainConfigurator) GenericCC {
	return GenericCC{c}
}

func (c GenericCC) DAOSupport() bool {
	if gc, ok := c.ChainConfigurator.(*goethereum.ChainConfig); ok {
		return gc.DAOForkSupport
	}
	if mg, ok := c.ChainConfigurator.(*coregeth.CoreGethChainConfig); ok {
		return mg.GetEthashEIP779Transition() != nil
	}
	panic(fmt.Sprintf("uimplemented DAO logic, config: %v", c.ChainConfigurator))
}

// Following vars define sufficient JSON schema keys for configurator type inference.
var (
	// Fields known (and unique, if possible) to etclabscore/core-geth.
	coregethSchemaSuffice = []string{
		"networkId", "config.networkId",
		"requireBlockHashes", "config.requireBlockHashes",
		"eip2FBlock", "config.eip2FBlock",
		"supportedProtocolVersions", "config.supportedProtocolVersions",
	}
	// Fields unknown for etclabscore/core-geth.
	coregethSchemaMustNot = []string{
		"EIP1108FBlock", "config.EIP1108FBlock",
		"eip158Block", "config.eip158Block",
		"daoForkSupport", "config.daoForkSupport",
	}

	// Fields known to ethereum/go-ethereum.
	goethereumSchemaSuffice = []string{
		"difficulty",
		"byzantiumBlock", "config.byzantiumBlock",
		"chainId", "config.chainId",
		"homesteadBlock", "config.homesteadBlock",
	}
	// Fields unknown to ethereum/go-ethereum.
	goethereumSchemaMustNot = []string{
		"engine",
		"genesis.seal",
		"networkId", "config.networkId",
	}
)

func UnmarshalChainConfigurator(input []byte) (ctypes.ChainConfigurator, error) {
	var cases = []struct {
		cnf        ctypes.ChainConfigurator
		sufficient []string
		negates    []string
	}{
		{&coregeth.CoreGethChainConfig{}, coregethSchemaSuffice, coregethSchemaMustNot},
		{&goethereum.ChainConfig{}, goethereumSchemaSuffice, goethereumSchemaMustNot},
	}
	for _, c := range cases {
		ok, err := asMapHasAnyKey(input, c.sufficient)
		if err != nil {
			return nil, err
		}
		negated := false
		if len(c.negates) > 0 {
			negated, err = asMapHasAnyKey(input, c.negates)
			if err != nil {
				return nil, err
			}
		}
		if !negated && ok {
			if err := json.Unmarshal(input, c.cnf); err != nil {
				return nil, err
			}
			return c.cnf, nil
		}
	}
	return nil, errors.New("invalid configurator schema")
}

func asMapHasAnyKey(input []byte, keys []string) (bool, error) {
	results := gjson.GetManyBytes(input, keys...)
	for _, g := range results {
		if g.Exists() {
			return true, nil
		}
	}
	return false, nil
}
