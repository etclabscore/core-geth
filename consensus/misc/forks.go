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

package misc

import (
	"github.com/ethereum/go-ethereum/core/types"
	common2 "github.com/ethereum/go-ethereum/params/types/common"
)

// VerifyForkHashes verifies that blocks conforming to network hard-forks do have
// the correct hashes, to avoid clients going off on different chains. This is an
// optional feature.
func VerifyForkHashes(config common2.ChainConfigurator, header *types.Header, uncle bool) error {
	// We don't care about uncles
	if uncle {
		return nil
	}
	// If the homestead reprice hash is set, validate it
	if config.GetEIP150Transition() != nil && *config.GetEIP150Transition() == header.Number.Uint64() {
		// TODO(meows)
		//if config.EIP150Hash != (common.Hash{}) && config.EIP150Hash != header.Hash() {
		//	return fmt.Errorf("homestead gas reprice fork: have 0x%x, want 0x%x", header.Hash(), config.EIP150Hash)
		//}
	}
	// All ok, return
	return nil
}
