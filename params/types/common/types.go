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


package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/common/math"
)

type UnsupportedConfigErr error

var (
	ErrUnsupportedConfigNoop  UnsupportedConfigErr = errors.New("unsupported config value (noop)")
	ErrUnsupportedConfigFatal UnsupportedConfigErr = errors.New("unsupported config value (fatal)")
)

type ErrUnsupportedConfig struct {
	Err    error
	Method string
	Value  interface{}
}

func (e ErrUnsupportedConfig) Error() string {
	return fmt.Sprintf("%v, field: %s, value: %v", e.Err, e.Method, e.Value)
}

func IsFatalUnsupportedErr(err error) bool {
	return err == ErrUnsupportedConfigFatal
}

func UnsupportedConfigError(err error, method string, value interface{}) ErrUnsupportedConfig {
	return ErrUnsupportedConfig{
		Err:    err,
		Method: method,
		Value:  value,
	}
}

// Uint64BigValOrMapHex is an encoding type for Parity's chain config,
// used for their 'blockReward' field.
// When only an intial value, eg 0:0x42 is set, the type is a hex-encoded string.
// When multiple values are set, eg modified block rewards, the type is a map of hex-encoded strings.
type Uint64BigValOrMapHex map[uint64]*big.Int

// UnmarshalJSON implements the json Unmarshaler interface.
func (m *Uint64BigValOrMapHex) UnmarshalJSON(input []byte) error {
	mm := make(map[math.HexOrDecimal64]math.HexOrDecimal256)
	err := json.Unmarshal(input, &mm)
	if err == nil {
		mp := Uint64BigValOrMapHex{}
		for k, v := range mm {
			u := uint64(k)
			mp[u] = v.ToInt()
		}
		*m = mp
		return nil
	}
	uq, err := strconv.Unquote(string(input))
	if err != nil {
		return err
	}
	input = []byte(uq)
	var b = new(math.HexOrDecimal256)
	err = b.UnmarshalText(input)
	if err != nil {
		return err
	}
	*m = Uint64BigValOrMapHex{0: b.ToInt()}
	return nil
}

// MarshalJSON implements the json Marshaler interface.
func (m Uint64BigValOrMapHex) MarshalJSON() (output []byte, err error) {
	mm := make(map[math.HexOrDecimal64]*math.HexOrDecimal256)
	for k, v := range m {
		mm[math.HexOrDecimal64(k)] = math.NewHexOrDecimal256(v.Int64())
	}
	return json.Marshal(mm)
}

// Uint64BigMapEncodesHex is a map that encodes and decodes w/ JSON hex format.
type Uint64BigMapEncodesHex map[uint64]*big.Int

// UnmarshalJSON implements the json Unmarshaler interface.
func (bb *Uint64BigMapEncodesHex) UnmarshalJSON(input []byte) error {
	// HACK: Parity uses raw numbers here...
	// It would be better to use a consistent format... instead of having to do interface{}-ing
	// and switch on types.
	m := make(map[math.HexOrDecimal64]interface{})
	err := json.Unmarshal(input, &m)
	if err != nil {
		return err
	}
	b := make(map[uint64]*big.Int)
	for k, v := range m {
		var vv *big.Int
		switch v.(type) {
		case string:
			var b = new(math.HexOrDecimal256)
			err = b.UnmarshalText([]byte(v.(string)))
			if err != nil {
				return err
			}
			vv = b.ToInt()
		case int64:
			vv = big.NewInt(v.(int64))
		}
		if vv != nil {
			b[uint64(k)] = vv
		}
	}
	*bb = b
	return nil
}

// MarshalJSON implements the json Marshaler interface.
func (b Uint64BigMapEncodesHex) MarshalJSON() ([]byte, error) {
	mm := make(map[math.HexOrDecimal64]*math.HexOrDecimal256)
	for k, v := range b {
		mm[math.HexOrDecimal64(k)] = math.NewHexOrDecimal256(v.Int64())
	}
	return json.Marshal(mm)
}
