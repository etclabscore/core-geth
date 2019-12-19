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


package tconvert

import (
	"strings"

	"github.com/ethereum/go-ethereum/params/confp"
	"github.com/ethereum/go-ethereum/params/types"
	"github.com/ethereum/go-ethereum/params/types/multigeth"
	"github.com/ethereum/go-ethereum/params/types/parity"
)

// NewParityChainSpec converts a go-ethereum genesis block into a Parity specific
// chain specification format.
func NewParityChainSpec(network string, genesis *paramtypes.Genesis, bootnodes []string) (*parity.ParityChainSpec, error) {
	spec := &parity.ParityChainSpec{
		Name:    network,
		Nodes:   bootnodes,
		Datadir: strings.ToLower(network),
	}
	if err := confp.Convert(genesis, spec); err != nil {
		return nil, err
	}
	if err := confp.Convert(genesis.Config, spec); err != nil {
		return nil, err
	}
	return spec, nil
}

// ToMultiGethGenesis converts a Parity chainspec to the corresponding MultiGeth datastructure.
// Note that the return value 'core.Genesis' includes the respective 'params.MultiGethChainConfig' values.
func ParityConfigToMultiGethGenesis(c *parity.ParityChainSpec) (*paramtypes.Genesis, error) {
	mg := &paramtypes.Genesis{
		Config: &multigeth.MultiGethChainConfig{},
	}

	if err := confp.Convert(c, mg); err != nil {
		return nil,err
	}
	if err := confp.Convert(c, mg.Config); err != nil {
		return nil, err
	}
	return mg, nil
}