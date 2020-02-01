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

	"github.com/ethereum/go-ethereum/params/types/ctypes"
	"github.com/ethereum/go-ethereum/params/types/goethereum"
	"github.com/ethereum/go-ethereum/params/types/multigeth"
	"github.com/ethereum/go-ethereum/params/types/oldmultigeth"
	"github.com/ethereum/go-ethereum/params/types/parity"
	"github.com/ethereum/go-ethereum/params/vars"
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
	if omg, ok := c.ChainConfigurator.(*oldmultigeth.ChainConfig); ok {
		return omg.DAOForkSupport
	}
	if mg, ok := c.ChainConfigurator.(*multigeth.MultiGethChainConfig); ok {
		return mg.GetEthashEIP779Transition() != nil
	}
	if pc, ok := c.ChainConfigurator.(*parity.ParityChainSpec); ok {
		return pc.Engine.Ethash.Params.DaoHardforkTransition != nil &&
			pc.Engine.Ethash.Params.DaoHardforkBeneficiary != nil &&
			*pc.Engine.Ethash.Params.DaoHardforkBeneficiary == vars.DAORefundContract &&
			len(pc.Engine.Ethash.Params.DaoHardforkAccounts) == len(vars.DAODrainList())
	}
	panic(fmt.Sprintf("uimplemented DAO logic, config: %v", c.ChainConfigurator))
}

// Following vars define sufficient JSON schema keys for configurator type inference.
var (
	paritySchemaKeysMust = []string{
		"engine",
		"genesis.seal",
	}
	// These are fields which must differentiate "new" multigeth from "old" multigeth.
	multigethSchemaMust = []string{
		"eip2FBlock", "config.eip2FBlock",
		"eip7FBlock", "config.eip7FBlock",
	}
	// These are fields which differentiate old multigeth from goethereum config.
	oldmultigethSchemaMust = []string{
		"eip160Block", "config.eip160Block",
		"ecip1010PauseBlock", "config.ecip1010PauseBlock",
	}
	goethereumSchemaMust = []string{
		"difficulty",
		"byzantiumBlock", "config.byzantiumBlock",
		"chainId", "config.chainId",
		"homesteadBlock", "config.homesteadBlock",
	}
)

func UnmarshalChainConfigurator(input []byte) (ctypes.ChainConfigurator, error) {
	var cases = []struct {
		cnf ctypes.ChainConfigurator
		fn  []string
	}{
		{&parity.ParityChainSpec{}, paritySchemaKeysMust},
		{&multigeth.MultiGethChainConfig{}, multigethSchemaMust},
		{&oldmultigeth.ChainConfig{}, oldmultigethSchemaMust},
		{&goethereum.ChainConfig{}, goethereumSchemaMust},
	}
	for _, c := range cases {
		ok, err := asMapHasAnyKey(input, c.fn)
		if err != nil {
			return nil, err
		}
		if ok {
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
		if g.Exists() && g.Value() != nil {
			return true, nil
		}
	}
	return false, nil
}
