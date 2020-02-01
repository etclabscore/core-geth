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
	paritySchemaKeysSuffice = []string{
		"engine",
		"genesis.seal",
	}
	paritySchemaKeysMustNot = []string{}

	// These are fields which must differentiate "new" multigeth from "old" multigeth.
	multigethSchemaSuffice = []string{
		"networkId", "config.networkId",
		"requireBlockHashes", "config.requireBlockHashes",
	}
	multigethSchemaMustNot = []string{
		"EIP1108FBlock", "config.EIP1108FBlock",
		"eip158Block", "config.eip158Block",
		"daoForkSupport", "config.daoForkSupport",
	}

	// These are fields which differentiate old multigeth from goethereum config.
	oldmultigethSchemaSuffice = []string{
		"EIP1108FBlock", "config.EIP1108FBlock",
		"eip7FBlock", "config.eip7FBlock",
		"eip2FBlock", "config.eip2FBlock",
		"ecip1010PauseBlock", "config.ecip1010PauseBlock",
		"disposalBlock", "config.disposalBlock",
	}
	oldmultigethSchemaMustNot = []string{
		"requireBlockHashes", "config.requireBlockHashes",
	}

	goethereumSchemaSuffice = []string{
		"difficulty",
		"byzantiumBlock", "config.byzantiumBlock",
		"chainId", "config.chainId",
		"homesteadBlock", "config.homesteadBlock",
	}
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
		{&parity.ParityChainSpec{}, paritySchemaKeysSuffice, paritySchemaKeysMustNot},
		{&multigeth.MultiGethChainConfig{}, multigethSchemaSuffice, multigethSchemaMustNot},
		{&oldmultigeth.ChainConfig{}, oldmultigethSchemaSuffice, oldmultigethSchemaMustNot},
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

// asMapHasValueForAnyKey extends the logic from asMapHasAnyKey, but
// requires that the fields also denote non-nil values.
// This was one path of logic, and may be removed later if not used.
//func asMapHasValueForAnyKey(input []byte, keys []string) (bool, error) {
//	results := gjson.GetManyBytes(input, keys...)
//	for _, g := range results {
//		if g.Exists() && g.Value() != nil {
//			return true, nil
//		}
//	}
//	return false, nil
//}
