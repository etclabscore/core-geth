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

	"github.com/ethereum/go-ethereum/params/types/ctypes"
	"github.com/ethereum/go-ethereum/params/types/goethereum"
	"github.com/ethereum/go-ethereum/params/types/multigeth"
	"github.com/ethereum/go-ethereum/params/types/parity"
)

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
