// Copyright 2016 The go-ethereum Authors
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

package mutations

import (
	"bytes"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params/types/ctypes"
	"github.com/ethereum/go-ethereum/params/vars"
	"github.com/holiman/uint256"
)

var (
	// ErrBadProDAOExtra is returned if a header doesn't support the DAO fork on a
	// pro-fork client.
	ErrBadProDAOExtra = errors.New("bad DAO pro-fork extra-data")

	// ErrBadNoDAOExtra is returned if a header does support the DAO fork on a no-
	// fork client.
	ErrBadNoDAOExtra = errors.New("bad DAO no-fork extra-data")
)

// VerifyDAOHeaderExtraData validates the extra-data field of a block header to
// ensure it conforms to DAO hard-fork rules.
//
// DAO hard-fork extension to the header validity:
//
//   - if the node is no-fork, do not accept blocks in the [fork, fork+10) range
//     with the fork specific extra-data set.
//   - if the node is pro-fork, require blocks in the specific range to have the
//     unique extra-data set.
func VerifyDAOHeaderExtraData(config ctypes.ChainConfigurator, header *types.Header) error {
	// If the config wants the DAO fork, it should validate the extra data.
	// Otherwise, like any other block or any other config, it should not care.
	daoForkBlock := config.GetEthashEIP779Transition()
	if daoForkBlock == nil {
		return nil
	}
	daoForkBlockB := new(big.Int).SetUint64(*daoForkBlock)
	// Make sure the block is within the fork's modified extra-data range
	limit := new(big.Int).Add(daoForkBlockB, vars.DAOForkExtraRange)
	if header.Number.Cmp(daoForkBlockB) < 0 || header.Number.Cmp(limit) >= 0 {
		return nil
	}
	if !bytes.Equal(header.Extra, vars.DAOForkBlockExtra) {
		return ErrBadProDAOExtra
	}
	return nil

	// Leaving the "old" code in as dead commented code for reference.
	//
	// // Short circuit validation if the node doesn't care about the DAO fork
	// daoForkBlock := config.GetEthashEIP779Transition()
	// // Second clause catches test configs with nil fork blocks (maybe set dynamically or
	// // testing agnostic of chain config).
	// if daoForkBlock == nil && !generic.AsGenericCC(config).DAOSupport() {
	//	return nil
	// }
	//
	// if daoForkBlock == nil {
	//
	// }
	//
	// daoForkBlockB := new(big.Int).SetUint64(*daoForkBlock)
	//
	// // Make sure the block is within the fork's modified extra-data range
	// limit := new(big.Int).Add(daoForkBlockB, vars.DAOForkExtraRange)
	// if header.Number.Cmp(daoForkBlockB) < 0 || header.Number.Cmp(limit) >= 0 {
	//	return nil
	// }
	// // Depending on whether we support or oppose the fork, validate the extra-data contents
	// if generic.AsGenericCC(config).DAOSupport() {
	//	if !bytes.Equal(header.Extra, vars.DAOForkBlockExtra) {
	//		return ErrBadProDAOExtra
	//	}
	// } else {
	//	if bytes.Equal(header.Extra, vars.DAOForkBlockExtra) {
	//		return ErrBadNoDAOExtra
	//	}
	// }
	// // All ok, header has the same extra-data we expect
	// return nil
}

// ApplyDAOHardFork modifies the state database according to the DAO hard-fork
// rules, transferring all balances of a set of DAO accounts to a single refund
// contract.
func ApplyDAOHardFork(statedb *state.StateDB) {
	// Retrieve the contract to refund balances into
	if !statedb.Exist(vars.DAORefundContract) {
		statedb.CreateAccount(vars.DAORefundContract)
	}

	// Move every DAO account and extra-balance account funds into the refund contract
	for _, addr := range vars.DAODrainList() {
		statedb.AddBalance(vars.DAORefundContract, statedb.GetBalance(addr))
		statedb.SetBalance(addr, new(uint256.Int))
	}
}
