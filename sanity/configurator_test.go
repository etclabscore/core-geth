package sanity

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"path/filepath"
	"testing"

	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/params/convert"
	paramtypes "github.com/ethereum/go-ethereum/params/types"
	"github.com/ethereum/go-ethereum/params/types/common"
	"github.com/ethereum/go-ethereum/params/types/parity"
	"github.com/ethereum/go-ethereum/tests"
)

func TestEquivalent_Features(t *testing.T) {

	mustValidate := func (c common.ChainConfigurator) {
		zero, max := uint64(0), uint64(math.MaxUint64)
		for _, head := range []*uint64{
			nil, &zero, &max,
		} {
			if err := common.IsValid(c, head); err != nil {
				t.Fatalf("invalid config, err: %v", err)
			}
		}
	}

	for name, oconf := range tests.Forks {
		log.Println(name)
		oconf := oconf

		if oconf.GetConsensusEngineType().IsUnknown() {
			oconf.MustSetConsensusEngineType(common.ConsensusEngineT_Ethash)
		}

		mustValidate(oconf)

		// Integration tests: conversion

		mg := &paramtypes.MultiGethChainConfig{}
		err := convert.Convert(oconf, mg)
		if err != nil {
			t.Fatal(err)
		}

		mustValidate(mg)

		if mg.GetConsensusEngineType().IsUnknown() {
			t.Fatal("unknown consensus mg")
		}

		err = common.Equivalent(oconf, mg)
		if err != nil {
			t.Log("--------------------")
			t.Errorf("%s oconf/mg err: %v", name, err)

			//s := spew.ConfigState{DisableMethods:true, DisablePointerAddresses: true, Indent: "    "}
			//t.Log("OG:", s.Sdump(oconf))
			//t.Log("MG:", s.Sdump(mg))

			nicelog := func (n *uint64) interface{} {
				if n == nil {
					return "nil"
				}
				return *n
			}
			t.Log("o 649", nicelog(oconf.GetEthashEIP649Transition()))
			t.Log("m 649", nicelog(mg.GetEthashEIP649Transition()))
			t.Log("o 1234", nicelog(oconf.GetEthashEIP1234Transition()))
			t.Log("m 1234", nicelog(mg.GetEthashEIP1234Transition()))

			t.Log(mg.GetEthashBlockRewardSchedule())

			// this looks right
			if v := oconf.GetEthashEIP649Transition(); v != nil {
				t.Log(name, "649T", *v)
			} else {
				t.Log(name, "649T", v)
			}
		}

		pc := &parity.ParityChainSpec{}
		err = convert.Convert(oconf, pc)
		if err != nil {
			t.Fatal(err)
		}

		mustValidate(pc)

		err = common.Equivalent(mg, pc)
		if err != nil {
			t.Errorf("%s oconf/p err: %v", name, err)
		}

	}
}

func TestEquivalent_ReadParity(t *testing.T) {
	// These configs are tested by tests/ (ethereum/tests) suite.
	// If passing there, the config pairs are equivalent IN THE CONTEXT OF THOSE TESTS,
	// which is what the configs are for.
	// In order to pass those tests, however, configs do not need to be strictly equivalent.
	// For example, one config might specify EIP1234 fork without a prior EIP649 fork, and
	// another may specify both (either simulaneously or in succession).
	// Both configs in this case yield will equivalent results, but
	// are not, strictly speaking, equivalent.
	// I've left this test here for debugging, and to demonstrate this case.
	t.Skip("(meowsbits): Not required.")
	parityP := filepath.Join("..", "params", "parity.json.d")
	for k, v := range tests.MapForkNameChainspecFileState {
		a := tests.Forks[k]

		b := &parity.ParityChainSpec{}
		bs, err := ioutil.ReadFile(filepath.Join(parityP, v))
		if err != nil {
			t.Fatal(err)
		}
		err = json.Unmarshal(bs, b)
		if err != nil {
			t.Fatal(err)
		}
		err = common.Equivalent(a, b)
		if err != nil {
			t.Log("-------------------")
			t.Log(b.Engine.Ethash.Params.BlockReward)
			t.Log(b.Engine.Ethash.Params.DifficultyBombDelays)
			t.Errorf("%s:%s err: %v", k, v, err)
		}
	}
}