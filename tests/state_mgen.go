package tests

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
)

// writeStateTestsReferencePairs defines reference pairs for use when writing tests.
// The reference (key) is used to define the environment and parameters, while the
// output from these tests run against the <value> state is actually written.
var writeStateTestsReferencePairs = map[string]string{
	"Byzantium":         "ETC_Atlantis",
	"ConstantinopleFix": "ETC_Agharta",
}

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

