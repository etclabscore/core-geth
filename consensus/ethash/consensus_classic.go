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
package ethash

import (
	"math/big"

	"github.com/ethereum/go-ethereum/params/types/ctypes"
)

func ecip1010Explosion(config ctypes.ChainConfigurator, next *big.Int, exPeriodRef *big.Int) {
	// https://github.com/ethereumproject/ECIPs/blob/master/ECIPs/ECIP-1010.md

	if next.Uint64() < *config.GetEthashECIP1010ContinueTransition() {
		exPeriodRef.SetUint64(*config.GetEthashECIP1010PauseTransition())
	} else {
		length := new(big.Int).SetUint64(*config.GetEthashECIP1010ContinueTransition() - *config.GetEthashECIP1010PauseTransition())
		exPeriodRef.Sub(exPeriodRef, length)
	}
}
