// Copyright 2022 The core-geth Authors
// This file is part of the core-geth library.
//
// The core-geth library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The core-geth library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the core-geth library. If not, see <http://www.gnu.org/licenses/>.

package accounts

// https://github.com/satoshilabs/slips/blob/master/slip-0044.md
const BIP0044CoinTypeTestnet uint32 = 0x1       // 1
const BIP0044CoinTypeEther uint32 = 0x3c        // 60
const BIP0044CoinTypeEtherClassic uint32 = 0x3d // 61

// BIP0044CoinType is the global default value for the chain configured hardware derivation path.
// Its value is set by an init() function and can be modified, along with its dependent global variables
// to change the default coin type by using the function SetCoinTypeConfiguration.
var BIP0044CoinType uint32

// SetCoinTypeConfiguration sets the global coin type configuration to the given value.
func SetCoinTypeConfiguration(coinType uint32) {
	BIP0044CoinType = coinType
	DefaultRootDerivationPath = DerivationPath{0x80000000 + 44, 0x80000000 + coinType, 0x80000000 + 0, 0}
	DefaultBaseDerivationPath = DerivationPath{0x80000000 + 44, 0x80000000 + coinType, 0x80000000 + 0, 0, 0}
	LegacyLedgerBaseDerivationPath = DerivationPath{0x80000000 + 44, 0x80000000 + coinType, 0x80000000 + 0, 0}
}

// init configures the global coin type and root derivation path for Ethereum mainnet.
func init() {
	SetCoinTypeConfiguration(BIP0044CoinTypeEther)
}
