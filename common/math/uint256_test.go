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

package math

import (
	"testing"

	"github.com/holiman/uint256"
)

func TestHexOrDecimalUint256(t *testing.T) {
	tests := []struct {
		input string
		num   *uint256.Int
		ok    bool
	}{
		{"", uint256.NewInt(0), true},
		{"0", uint256.NewInt(0), false},
		{"0x0", uint256.NewInt(0), true},
		{"12345678", uint256.NewInt(0), false},
		{"0x12345678", uint256.NewInt(0x12345678), true},
		{"0X12345678", uint256.NewInt(0x12345678), true},
		// Tests for leading zero behaviour:
		{"0123456789", uint256.NewInt(0), false}, // note: not octal
		{"00", uint256.NewInt(0), false},
		{"0x00", uint256.NewInt(0), false},
		{"0x012345678abc", uint256.NewInt(0), false},
		// Invalid syntax:
		{"abcdef", nil, false},
		{"0xgg", nil, false},
		// Larger than 256 bits:
		{"115792089237316195423570985008687907853269984665640564039457584007913129639936", nil, false},
	}
	for _, test := range tests {
		t.Logf("Unmarshaling %q", test.input)
		var num HexOrDecimalUint256
		err := num.UnmarshalText([]byte(test.input))
		if (err == nil) != test.ok {
			t.Errorf("ParseBig(%q) -> (err == nil) == %t, want %t", test.input, err == nil, test.ok)
			continue
		}
		if test.num != nil && (*uint256.Int)(&num).Cmp(test.num) != 0 {
			t.Errorf("ParseBig(%q) -> %d, want %d", test.input, (*uint256.Int)(&num), test.num)
		}
	}
}

func TestMustParseUint256(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("MustParseBig should've panicked")
		}
	}()
	MustParseUint256("ggg")
}

func TestUint64Max(t *testing.T) {
	a := uint256.NewInt(10)
	b := uint256.NewInt(5)

	max1 := Uint256Max(a, b)
	if max1 != a {
		t.Errorf("Expected %d got %d", a, max1)
	}

	max2 := Uint256Max(b, a)
	if max2 != a {
		t.Errorf("Expected %d got %d", a, max2)
	}
}

func TestUint64Min(t *testing.T) {
	a := uint256.NewInt(10)
	b := uint256.NewInt(5)

	min1 := Uint256Min(a, b)
	if min1 != b {
		t.Errorf("Expected %d got %d", b, min1)
	}

	min2 := Uint256Min(b, a)
	if min2 != b {
		t.Errorf("Expected %d got %d", b, min2)
	}
}
