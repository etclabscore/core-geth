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
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/vm"
)

// RunSetPost runs the state subtest for a given config, and writes the resulting
// state to the corresponding subtest post field.
func (t *StateTest) RunSetPost(subtest StateSubtest, vmconfig vm.Config) error {
	state, root, err := t.RunNoVerify(subtest, vmconfig, false, rawdb.HashScheme)
	if err != nil {
		return err
	}
	t.json.Post[subtest.Fork][subtest.Index].Root = common.UnprefixedHash(root)
	t.json.Post[subtest.Fork][subtest.Index].Logs = common.UnprefixedHash(rlpHash(state.StateDB.Logs()))
	t.json.Post[subtest.Fork][subtest.Index].filled = true
	return nil
}
