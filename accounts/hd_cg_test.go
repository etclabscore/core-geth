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

import (
	"fmt"
	"reflect"
	"testing"
)

// Tests that HD derivation paths can be correctly parsed into our internal binary
// representation.
func TestHDPathParsing_CG(t *testing.T) {
	t.Run("testParseDerivationPath_Default", testHDPathParsing(BIP0044CoinType))
	t.Run("testParseDerivationPath_Ether", testHDPathParsing(BIP0044CoinTypeEther))
	t.Run("testParseDerivationPath_EtherClassic", testHDPathParsing(BIP0044CoinTypeEtherClassic))
	t.Run("testParseDerivationPath_Testnet", testHDPathParsing(BIP0044CoinTypeTestnet))
}

func testHDPathParsing(coinType uint32) func(t *testing.T) {
	return func(t *testing.T) {
		// Set the coin type to the one we're testing.
		// This assigns the default derivations and iterator to the coin type.
		SetCoinTypeConfiguration(coinType)
		defer SetCoinTypeConfiguration(BIP0044CoinTypeEther)

		tests := []struct {
			input  string
			output DerivationPath
		}{
			// Plain absolute derivation paths
			{"m/44'/" + mustStr(coinType) + "'/0'/0", DerivationPath{0x80000000 + 44, 0x80000000 + coinType, 0x80000000 + 0, 0}},
			{"m/44'/" + mustStr(coinType) + "'/0'/128", DerivationPath{0x80000000 + 44, 0x80000000 + coinType, 0x80000000 + 0, 128}},
			{"m/44'/" + mustStr(coinType) + "'/0'/0'", DerivationPath{0x80000000 + 44, 0x80000000 + coinType, 0x80000000 + 0, 0x80000000 + 0}},
			{"m/44'/" + mustStr(coinType) + "'/0'/128'", DerivationPath{0x80000000 + 44, 0x80000000 + coinType, 0x80000000 + 0, 0x80000000 + 128}},
			{"m/2147483692/" + mustStr(0x80000000+coinType) + "/2147483648/0", DerivationPath{0x80000000 + 44, 0x80000000 + coinType, 0x80000000 + 0, 0}},
			{"m/2147483692/" + mustStr(0x80000000+coinType) + "/2147483648/2147483648", DerivationPath{0x80000000 + 44, 0x80000000 + coinType, 0x80000000 + 0, 0x80000000 + 0}},

			// Plain relative derivation paths
			{"0", DerivationPath{0x80000000 + 44, 0x80000000 + coinType, 0x80000000 + 0, 0, 0}},
			{"128", DerivationPath{0x80000000 + 44, 0x80000000 + coinType, 0x80000000 + 0, 0, 128}},
			{"0'", DerivationPath{0x80000000 + 44, 0x80000000 + coinType, 0x80000000 + 0, 0, 0x80000000 + 0}},
			{"128'", DerivationPath{0x80000000 + 44, 0x80000000 + coinType, 0x80000000 + 0, 0, 0x80000000 + 128}},
			{"2147483648", DerivationPath{0x80000000 + 44, 0x80000000 + coinType, 0x80000000 + 0, 0, 0x80000000 + 0}},

			// Hexadecimal absolute derivation paths
			{"m/0x2C'/" + mustHex(coinType) + "'/0x00'/0x00", DerivationPath{0x80000000 + 44, 0x80000000 + coinType, 0x80000000 + 0, 0}},
			{"m/0x2C'/" + mustHex(coinType) + "'/0x00'/0x80", DerivationPath{0x80000000 + 44, 0x80000000 + coinType, 0x80000000 + 0, 128}},
			{"m/0x2C'/" + mustHex(coinType) + "'/0x00'/0x00'", DerivationPath{0x80000000 + 44, 0x80000000 + coinType, 0x80000000 + 0, 0x80000000 + 0}},
			{"m/0x2C'/" + mustHex(coinType) + "'/0x00'/0x80'", DerivationPath{0x80000000 + 44, 0x80000000 + coinType, 0x80000000 + 0, 0x80000000 + 128}},
			{"m/0x8000002C/" + mustHex(0x80000000+coinType) + "/0x80000000/0x00", DerivationPath{0x80000000 + 44, 0x80000000 + coinType, 0x80000000 + 0, 0}},
			{"m/0x8000002C/" + mustHex(0x80000000+coinType) + "/0x80000000/0x80000000", DerivationPath{0x80000000 + 44, 0x80000000 + coinType, 0x80000000 + 0, 0x80000000 + 0}},

			// Hexadecimal relative derivation paths
			{"0x00", DerivationPath{0x80000000 + 44, 0x80000000 + coinType, 0x80000000 + 0, 0, 0}},
			{"0x80", DerivationPath{0x80000000 + 44, 0x80000000 + coinType, 0x80000000 + 0, 0, 128}},
			{"0x00'", DerivationPath{0x80000000 + 44, 0x80000000 + coinType, 0x80000000 + 0, 0, 0x80000000 + 0}},
			{"0x80'", DerivationPath{0x80000000 + 44, 0x80000000 + coinType, 0x80000000 + 0, 0, 0x80000000 + 128}},
			{"0x80000000", DerivationPath{0x80000000 + 44, 0x80000000 + coinType, 0x80000000 + 0, 0, 0x80000000 + 0}},

			// Weird inputs just to ensure they work
			{"	m  /   44			'\n/\n   " + mustStr(coinType) + "	\n\n\t'   /\n0 ' /\t\t	0", DerivationPath{0x80000000 + 44, 0x80000000 + coinType, 0x80000000 + 0, 0}},

			// Invalid derivation paths
			{"", nil},   // Empty relative derivation path
			{"m", nil},  // Empty absolute derivation path
			{"m/", nil}, // Missing last derivation component
			{"/44'/" + mustStr(coinType) + "'/0'/0", nil}, // Absolute path without m prefix, might be user error
			{"m/2147483648'", nil},                        // Overflows 32 bit integer
			{"m/-1'", nil},                                // Cannot contain negative number
		}
		for i, tt := range tests {
			if path, err := ParseDerivationPath(tt.input); !reflect.DeepEqual(path, tt.output) {
				t.Errorf("test %d: parse mismatch: have %v (%v), want %v", i, path, err, tt.output)
			} else if path == nil && err == nil {
				t.Errorf("test %d: nil path and error: %v", i, err)
			}
		}
	}
}

func testHdPathIteration(coinType uint32) func(t *testing.T) {
	return func(t *testing.T) {
		SetCoinTypeConfiguration(coinType)
		testDerive(t, DefaultIterator(DefaultBaseDerivationPath),
			[]string{
				"m/44'/" + mustStr(BIP0044CoinType) + "'/0'/0/0", "m/44'/" + mustStr(BIP0044CoinType) + "'/0'/0/1",
				"m/44'/" + mustStr(BIP0044CoinType) + "'/0'/0/2", "m/44'/" + mustStr(BIP0044CoinType) + "'/0'/0/3",
				"m/44'/" + mustStr(BIP0044CoinType) + "'/0'/0/4", "m/44'/" + mustStr(BIP0044CoinType) + "'/0'/0/5",
				"m/44'/" + mustStr(BIP0044CoinType) + "'/0'/0/6", "m/44'/" + mustStr(BIP0044CoinType) + "'/0'/0/7",
				"m/44'/" + mustStr(BIP0044CoinType) + "'/0'/0/8", "m/44'/" + mustStr(BIP0044CoinType) + "'/0'/0/9",
			})

		testDerive(t, DefaultIterator(LegacyLedgerBaseDerivationPath),
			[]string{
				"m/44'/" + mustStr(BIP0044CoinType) + "'/0'/0", "m/44'/" + mustStr(BIP0044CoinType) + "'/0'/1",
				"m/44'/" + mustStr(BIP0044CoinType) + "'/0'/2", "m/44'/" + mustStr(BIP0044CoinType) + "'/0'/3",
				"m/44'/" + mustStr(BIP0044CoinType) + "'/0'/4", "m/44'/" + mustStr(BIP0044CoinType) + "'/0'/5",
				"m/44'/" + mustStr(BIP0044CoinType) + "'/0'/6", "m/44'/" + mustStr(BIP0044CoinType) + "'/0'/7",
				"m/44'/" + mustStr(BIP0044CoinType) + "'/0'/8", "m/44'/" + mustStr(BIP0044CoinType) + "'/0'/9",
			})

		testDerive(t, LedgerLiveIterator(DefaultBaseDerivationPath),
			[]string{
				"m/44'/" + mustStr(BIP0044CoinType) + "'/0'/0/0", "m/44'/" + mustStr(BIP0044CoinType) + "'/1'/0/0",
				"m/44'/" + mustStr(BIP0044CoinType) + "'/2'/0/0", "m/44'/" + mustStr(BIP0044CoinType) + "'/3'/0/0",
				"m/44'/" + mustStr(BIP0044CoinType) + "'/4'/0/0", "m/44'/" + mustStr(BIP0044CoinType) + "'/5'/0/0",
				"m/44'/" + mustStr(BIP0044CoinType) + "'/6'/0/0", "m/44'/" + mustStr(BIP0044CoinType) + "'/7'/0/0",
				"m/44'/" + mustStr(BIP0044CoinType) + "'/8'/0/0", "m/44'/" + mustStr(BIP0044CoinType) + "'/9'/0/0",
			})
	}
}

func TestHdPathIteration_CG(t *testing.T) {
	t.Run("TestHdPathIteration_Default", testHdPathIteration(BIP0044CoinType))
	t.Run("TestHdPathIteration_Ether", testHdPathIteration(BIP0044CoinTypeEther))
	t.Run("TestHdPathIteration_EtherClassic", testHdPathIteration(BIP0044CoinTypeEtherClassic))
	t.Run("TestHdPathIteration_Testnet", testHdPathIteration(BIP0044CoinTypeTestnet))
}

func mustStr(i uint32) string {
	return fmt.Sprintf("%d", i)
}

func mustHex(i uint32) string {
	return fmt.Sprintf("0x%x", i)
}
