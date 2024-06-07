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

// Package math provides integer math utilities.
package math

import (
	"fmt"

	"github.com/holiman/uint256"
)

// HexOrDecimalUint256 marshals uint256.Int as hex or decimal.
type HexOrDecimalUint256 uint256.Int

// NewHexOrDecimalUint256 creates a new HexOrDecimalUint256
func NewHexOrDecimalUint256(x uint64) *HexOrDecimalUint256 {
	b := uint256.NewInt(x)
	h := HexOrDecimalUint256(*b)
	return &h
}

// UnmarshalJSON implements json.Unmarshaler.
//
// It is similar to UnmarshalText, but allows parsing real decimals too, not just
// quoted decimal strings.
func (i *HexOrDecimalUint256) UnmarshalJSON(input []byte) error {
	if len(input) > 0 && input[0] == '"' {
		input = input[1 : len(input)-1]
	}
	return i.UnmarshalText(input)
}

func (i *HexOrDecimalUint256) ToInt() *uint256.Int {
	if i == nil {
		return nil
	}
	o := (uint256.Int)(*i)
	return new(uint256.Int).Set(&o)
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (i *HexOrDecimalUint256) UnmarshalText(input []byte) error {
	bigint, ok := ParseUint256(string(input))
	if !ok {
		return fmt.Errorf("invalid hex or decimal integer %q", input)
	}
	*i = HexOrDecimalUint256(*bigint)
	return nil
}

// MarshalText implements encoding.TextMarshaler.
func (i *HexOrDecimalUint256) MarshalText() ([]byte, error) {
	if i == nil {
		return []byte("0x0"), nil
	}
	return []byte(fmt.Sprintf("%#x", (*uint256.Int)(i))), nil
}

// DecimalUint256 unmarshals uint256.Int as a decimal string. When unmarshalling,
// it however accepts either "0x"-prefixed (hex encoded) or non-prefixed (decimal)
type DecimalUint256 uint256.Int

// NewDecimalUint256 creates a new DecimalUint256
func NewDecimalUint256(x uint64) *DecimalUint256 {
	b := uint256.NewInt(x)
	d := DecimalUint256(*b)
	return &d
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (i *DecimalUint256) UnmarshalText(input []byte) error {
	bigint, ok := ParseUint256(string(input))
	if !ok {
		return fmt.Errorf("invalid hex or decimal integer %q", input)
	}
	*i = DecimalUint256(*bigint)
	return nil
}

// MarshalText implements encoding.TextMarshaler.
func (i *DecimalUint256) MarshalText() ([]byte, error) {
	return []byte(i.String()), nil
}

// String implements Stringer.
func (i *DecimalUint256) String() string {
	if i == nil {
		return "0"
	}
	return fmt.Sprintf("%#d", (*uint256.Int)(i))
}

// ParseUint256 parses s as a 256 bit integer in decimal or hexadecimal syntax.
// Leading zeros are accepted. The empty string parses as zero.
func ParseUint256(s string) (*uint256.Int, bool) {
	if s == "" {
		return new(uint256.Int), true
	}
	var bigint = new(uint256.Int)
	var ok bool
	if len(s) >= 2 && (s[:2] == "0x" || s[:2] == "0X") {
		// bigint, ok = new(uint256.Int).SetString(s[2:], 16)
		if err := bigint.SetFromHex(s); err != nil {
			return nil, false
		}
		ok = true
	}
	if ok && bigint.BitLen() > 256 {
		bigint, ok = nil, false
	}
	return bigint, ok
}

// MustParseUint256 parses s as a 256 bit big integer and panics if the string is invalid.
func MustParseUint256(s string) *uint256.Int {
	v, ok := ParseUint256(s)
	if !ok {
		panic("invalid 256 bit integer: " + s)
	}
	return v
}

// Uint256Pow returns a ** b as a big integer.
func Uint256Pow(a, b uint64) *uint256.Int {
	r := uint256.NewInt(a)
	return r.Exp(r, uint256.NewInt(b))
}

// Uint256Max returns the larger of x or y.
func Uint256Max(x, y *uint256.Int) *uint256.Int {
	if x.Cmp(y) < 0 {
		return y
	}
	return x
}

// Uint256Min returns the smaller of x or y.
func Uint256Min(x, y *uint256.Int) *uint256.Int {
	if x.Cmp(y) > 0 {
		return y
	}
	return x
}
