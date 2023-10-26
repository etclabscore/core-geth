// Copyright 2022 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package merkle

import (
	"crypto/sha256"
	"errors"
	"reflect"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// Value represents either a 32 byte leaf value or hash node in a binary Merkle tree/partial proof.
type Value [32]byte

// Values represent a series of Merkle tree leaves/nodes.
type Values []Value

var valueT = reflect.TypeOf(Value{})

// UnmarshalJSON parses a Merkle value in hex syntax.
func (m *Value) UnmarshalJSON(input []byte) error {
	return hexutil.UnmarshalFixedJSON(valueT, input, m[:])
}

// VerifyProof verifies a Merkle proof branch for a single value in a binary Merkle tree (index is a generalized tree index).
func VerifyProof(root common.Hash, index uint64, branch Values, value Value) error {
	if len(branch) == 0 {
		return errors.New("proof branch is empty")
	}
	if index == 0 || index > uint64(len(branch)) {
		return errors.New("index is out of bounds")
	}

	hasher := sha256.New()
	var wg sync.WaitGroup
	var errMutex sync.Mutex
	var verifyErr error

	// Parallelize proof verification for better performance
	for i, sibling := range branch {
		wg.Add(1)
		go func(i int, sibling Value) {
			defer wg.Done()

			hasher.Reset()
			if index&(1<<i) == 0 {
				hasher.Write(value[:])
				hasher.Write(sibling[:])
			} else {
				hasher.Write(sibling[:])
				hasher.Write(value[:])
			}
			hasher.Sum(value[:0])

			// Use a mutex to safely handle errors
			errMutex.Lock()
			if index>>i == 1 && common.Hash(value) != root {
				verifyErr = errors.New("root mismatch")
			}
			errMutex.Unlock()
		}(i, sibling)
	}

	wg.Wait()
	if verifyErr != nil {
		return verifyErr
	}

	if index>>len(branch) != 1 {
		return errors.New("branch has extra items")
	}

	return nil
}

	
