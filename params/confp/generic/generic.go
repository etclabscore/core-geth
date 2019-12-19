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
	"fmt"

	"github.com/ethereum/go-ethereum/params/types/ctypes"
	"github.com/ethereum/go-ethereum/params/types/goethereum"
	"github.com/ethereum/go-ethereum/params/types/multigeth"
	"github.com/ethereum/go-ethereum/params/types/parity"
	"github.com/ethereum/go-ethereum/params/vars"
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

func UnmarshalChainConfigurator(input []byte) (ctypes.ChainConfigurator, error) {
	var map1 = make(map[string]interface{})
	err := json.Unmarshal(input, &map1)
	if err != nil {
		return nil, err
	}
	if _, ok := map1["params"]; ok {
		pspec := &parity.ParityChainSpec{}
		err = json.Unmarshal(input, pspec)
		if err != nil {
			return nil, err
		}
		return pspec, nil
	}

	if _, ok := map1["networkId"]; ok {
		mspec := &multigeth.MultiGethChainConfig{}
		err = json.Unmarshal(input, mspec)
		if err != nil {
			return nil, err
		}
		return mspec, nil
	}

	gspec := &goethereum.ChainConfig{}
	err = json.Unmarshal(input,gspec)
	if err != nil {
		return nil, err
	}
	return gspec, nil
}
