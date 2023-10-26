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

package integration

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/params/confp"
	"github.com/ethereum/go-ethereum/params/types/coregeth"
	"github.com/ethereum/go-ethereum/params/types/ctypes"
	"github.com/ethereum/go-ethereum/params/vars"
	"github.com/ethereum/go-ethereum/tests"
)

func validateChainConfigurator(c ctypes.ChainConfigurator, t *testing.T) {
	zero, max := uint64(0), uint64(math.MaxUint64)
	for _, head := range []*uint64{nil, &zero, &max} {
		if err := confp.IsValid(c, head); err != nil {
			t.Fatalf("invalid config, err: %v", err)
		}
	}
}

func logChainConfiguratorInfo(c ctypes.ChainConfigurator, name string, t *testing.T) {
	t.Log(name)
	t.Log("Consensus Engine Type:", c.GetConsensusEngineType())
	t.Log("Ethash EIP649 Transition:", nicelog(c.GetEthashEIP649Transition()))
	t.Log("Ethash EIP1234 Transition:", nicelog(c.GetEthashEIP1234Transition()))
	t.Log("Ethash Block Reward Schedule:", c.GetEthashBlockRewardSchedule())
	if v := c.GetEthashEIP649Transition(); v != nil {
		t.Log(name, "649T", *v)
	} else {
		t.Log(name, "649T", v)
	}
}

func testEquivalent(t *testing.T, name string, oconf ctypes.ChainConfigurator, mg *coregeth.CoreGethChainConfig) {
	logChainConfiguratorInfo(oconf, name, t)

	// Integration tests: conversion
	err := confp.Crush(mg, oconf, true)
	if err != nil {
		t.Fatal(err)
	}

	logChainConfiguratorInfo(mg, name, t)
	validateChainConfigurator(mg, t)

	if mg.GetConsensusEngineType().IsUnknown() {
		t.Fatal("Unknown consensus mg")
	}

	err = confp.Equivalent(oconf, mg)
	if err != nil {
		t.Errorf("Equivalence: %s oconf/mg err: %v", name, err) // With error.
	}
}

func TestConstantinopleEquivalence(t *testing.T) {
	conf := tests.Forks["Constantinople"]
	pspec := &coregeth.CoreGethChainConfig{}
	err := confp.Crush(pspec, conf, true)
	if err != nil {
		t.Fatal(err)
	}

	if pspec.GetEthashDifficultyBombDelaySchedule()[*conf.GetEthashEIP1234Transition()].Cmp(vars.EIP1234DifficultyBombDelay) != 0 {
		t.Error("Bad")
	}
}

func TestEquivalent_Features(t *testing.T) {
	for name, oconf := range tests.Forks {
		log.Println(name)
		oconf := oconf

		if oconf.GetConsensusEngineType().IsUnknown() {
			oconf.MustSetConsensusEngineType(ctypes.ConsensusEngineT_Ethash)
		}

		validateChainConfigurator(oconf, t)

		mg := &coregeth.CoreGethChainConfig{}
		testEquivalent(t, name, oconf, mg)
	}
}
