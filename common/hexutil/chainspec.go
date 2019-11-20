package hexutil

import (
	"encoding/json"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/common/math"
)

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
