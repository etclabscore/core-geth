package tests

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
)

// RunSetPost runs the state subtest for a given config, and writes the resulting
// state to the corresponding subtest post field.
func (t *StateTest) RunSetPost(subtest StateSubtest, vmconfig vm.Config) error {
	statedb, root, err := t.RunNoVerify(subtest, vmconfig)
	if err != nil {
		return err
	}
	t.json.Post[subtest.Fork][subtest.Index].Root = common.UnprefixedHash(root)
	t.json.Post[subtest.Fork][subtest.Index].Logs = common.UnprefixedHash(rlpHash(statedb.Logs()))
	return nil
}

// filledPostStates returns true if all poststate elems are filled (non zero-valued)
func filledPostStates(s []stPostState) bool {
	for _, l := range s {
		if common.Hash(l.Root) == (common.Hash{}) {
			return false
		}
	}
	return true
}

