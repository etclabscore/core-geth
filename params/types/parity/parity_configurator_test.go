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

package parity

import (
	"testing"

	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/params/types/ctypes"
)

// This file contains a few unit tests for the parity-specific configuration interface.
// It does not contain integration tests, since this logic is covered by the test in convert_test.go,
// where specs are read, filled (converted), and verified equivalent.
//   Those tests cannot pass if the logic here is not sound.

func TestParityChainSpec_GetConsensusEngineType(t *testing.T) {
	spec := new(ParityChainSpec)

	if engine := (*spec).GetConsensusEngineType(); engine != ctypes.ConsensusEngineT_Unknown {
		t.Error("unwanted engine type", engine)
	}

	spec.Engine.Ethash.Params.MinimumDifficulty = math.NewHexOrDecimal256(42)
	if engine := (*spec).GetConsensusEngineType(); engine != ctypes.ConsensusEngineT_Ethash {
		t.Error("mismatch engine", engine)
	}
}

func TestParityChainSpec_GetSetUint64(t *testing.T) {
	spec := &ParityChainSpec{}
	if spec.GetEthashHomesteadTransition() != nil {
		t.Error("not empty")
	}
	spec.SetEthashHomesteadTransition(nil)
	if spec.GetEthashHomesteadTransition() != nil {
		t.Error("not nil")
	}
	fortyTwo := uint64(42)
	spec.SetEthashHomesteadTransition(&fortyTwo)
	if *spec.GetEthashHomesteadTransition() != fortyTwo {
		t.Error("not right answer")
	}
}
