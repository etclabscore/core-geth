// Copyright 2017 The go-ethereum Authors
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

package keccak

import (
	"encoding/binary"

	"github.com/ethereum/go-ethereum/crypto"
)

// keccakHasher produces a hash from combining the hash of the block header contents
// combined with a nonce.
func keccakHasher(hash []byte, nonce uint64) ([]byte, []byte) {

	// Combine hash+nonce into a 64 byte seed
	seed := make([]byte, 40)
	copy(seed, hash)
	binary.BigEndian.PutUint64(seed[32:], nonce)

	solution := crypto.Keccak256(seed)
	return solution, solution
}
