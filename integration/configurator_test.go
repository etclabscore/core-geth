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

func TestConstantinopleEquivalence(t *testing.T) {
	conf := tests.Forks["Constantinople"]
	pspec := &coregeth.CoreGethChainConfig{}
	err := confp.Crush(pspec, conf, true)
	if err != nil {
		t.Fatal(err)
	}
	// This test's config will set Byz delay (3m) at 0, and Const delay (5m) at 0.
	// This check ensures that 5m delay being greater than 3m takes precedence at simultaneous blocks.
	if pspec.GetEthashDifficultyBombDelaySchedule()[*conf.GetEthashEIP1234Transition()].Cmp(vars.EIP1234DifficultyBombDelay) != 0 {
		t.Error("bad")
	}
}

func TestEquivalent_Features(t *testing.T) {
	mustValidate := func(c ctypes.ChainConfigurator) {
		zero, max := uint64(0), uint64(math.MaxUint64)
		for _, head := range []*uint64{
			nil, &zero, &max,
		} {
			if err := confp.IsValid(c, head); err != nil {
				t.Fatalf("invalid config, err: %v", err)
			}
		}
	}

	for name, oconf := range tests.Forks {
		log.Println(name)
		oconf := oconf

		if oconf.GetConsensusEngineType().IsUnknown() {
			oconf.MustSetConsensusEngineType(ctypes.ConsensusEngineT_Ethash)
		}

		mustValidate(oconf)

		// Integration tests: conversion

		mg := &coregeth.CoreGethChainConfig{}
		err := confp.Crush(mg, oconf, true)
		if err != nil {
			t.Fatal(err)
		}

		mustValidate(mg)

		if mg.GetConsensusEngineType().IsUnknown() {
			t.Fatal("unknown consensus mg")
		}

		nicelog := func(n *uint64) interface{} {
			if n == nil {
				return "nil"
			}
			return *n
		}
		debuglog := func(a, b ctypes.ChainConfigurator) {
			// Debugging log lines.
			t.Log("o", oconf.GetConsensusEngineType())
			t.Log("m", mg.GetConsensusEngineType())

			t.Log("o 649", nicelog(oconf.GetEthashEIP649Transition()))
			t.Log("m 649", nicelog(mg.GetEthashEIP649Transition()))
			t.Log("o 1234", nicelog(oconf.GetEthashEIP1234Transition()))
			t.Log("m 1234", nicelog(mg.GetEthashEIP1234Transition()))

			t.Log(mg.GetEthashBlockRewardSchedule())

			if v := oconf.GetEthashEIP649Transition(); v != nil {
				t.Log(name, "649T", *v)
			} else {
				t.Log(name, "649T", v)
			}

			t.Log("--------------------")
			j, _ := json.MarshalIndent(oconf, "", "    ")
			t.Log(string(j))
			j, _ = json.MarshalIndent(mg, "", "    ")
			t.Log(string(j))
		}

		err = confp.Equivalent(oconf, mg)
		if err != nil {
			t.Errorf("Equivalence: %s oconf/mg err: %v", name, err) // With error.
			debuglog(oconf, mg)
		}
	}
}
