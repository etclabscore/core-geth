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

package tests

import (
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

// This file holds variable and type relating specifically
// to the task of generating tests.

var (
	CG_GENERATE_STATE_TESTS_KEY             = "COREGETH_TESTS_GENERATE_STATE_TESTS"
	CG_GENERATE_DIFFICULTY_TESTS_KEY        = "COREGETH_TESTS_GENERATE_DIFFICULTY_TESTS"
	CG_GENERATE_DIFFICULTY_TEST_CONFIGS_KEY = "COREGETH_TESTS_GENERATE_DIFFICULTY_TESTS_CONFIGS"

	// Feature Equivalence tests use convert.Convert to
	// run tests using alternating ChainConfig data type implementations.
	CG_CHAINCONFIG_FEATURE_EQ_COREGETH_KEY = "COREGETH_TESTS_CHAINCONFIG_FEATURE_EQUIVALENCE_COREGETH"

	CG_CHAINCONFIG_CONSENSUS_EQ_CLIQUE = "COREGETH_TESTS_CHAINCONFIG_CONSENSUS_EQUIVALENCE_CLIQUE"

	// CG_CHAINCONFIG_CHAINSPECS_COREGETH_KEY uses core-geth config data (in params/core-geth.json.d)
	// when applicable as equivalent config implementations for the default Go data type
	// configs.
	CG_CHAINCONFIG_CHAINSPECS_COREGETH_KEY = "COREGETH_TESTS_CHAINCONFIG_COREGETH_SPECS"
)

type chainspecRefsT map[string]chainspecRef

var chainspecRefsState = chainspecRefsT{}
var chainspecRefsDifficulty = chainspecRefsT{}

type chainspecRef struct {
	Filename string `json:"filename"`
	Sha1Sum  []byte `json:"sha1sum"`
}

func (c chainspecRef) String() string {
	return fmt.Sprintf("file: %s, file.sha1sum: %x", c.Filename, c.Sha1Sum)
}

func (c *chainspecRef) UnmarshalJSON(input []byte) error {
	type xT struct {
		F string `json:"filename"`
		S string `json:"sha1sum"`
	}
	var x = xT{}
	err := json.Unmarshal(input, &x)
	if err != nil {
		return err
	}
	c.Filename = x.F
	c.Sha1Sum = common.Hex2Bytes(x.S)
	return nil
}

func (c chainspecRef) MarshalJSON() ([]byte, error) {
	var x = struct {
		F string `json:"filename"`
		S string `json:"sha1sum"`
	}{
		F: c.Filename,
		S: common.Bytes2Hex(c.Sha1Sum[:]),
	}

	return json.MarshalIndent(x, "", "    ")
}
