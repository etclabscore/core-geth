package chainspec

import (
	"encoding/json"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

// Uint64BigValOrMapHex is an encoding type for Parity's chain config,
// used for their 'blockReward' field.
// When only an intial value, eg 0:0x42 is set, the type is a hex-encoded string.
// When multiple values are set, eg modified block rewards, the type is a map of hex-encoded strings.
type Uint64BigValOrMapHex map[uint64]*big.Int

// UnmarshalJSON implements the json Unmarshaler interface.
func (m *Uint64BigValOrMapHex) UnmarshalJSON(input []byte) error {
	mm := make(map[string]string)
	err := json.Unmarshal(input, &mm)
	if err == nil {
		mp := Uint64BigValOrMapHex{}
		for k, v := range mm {
			u, err := hexutil.DecodeUint64(k)
			if err != nil {
				return err
			}
			b, err := hexutil.DecodeBig(v)
			if err != nil {
				return err
			}
			mp[u] = b
		}
		*m = mp
		return nil
	}
	uq, err := strconv.Unquote(string(input))
	if err != nil {
		return err
	}
	input = []byte(uq)
	b, err := hexutil.DecodeBig(string(input))
	if err != nil {
		return err
	}
	*m = Uint64BigValOrMapHex{0: b}
	return nil
}

// MarshalJSON implements the json Marshaler interface.
func (m Uint64BigValOrMapHex) MarshalJSON() (output []byte, err error) {
	if v, ok := m[0]; ok && len(m) == 1 {
		return []byte(strconv.Quote(hexutil.EncodeBig(v))), nil
	}

	mm := make(map[hexutil.Uint64]hexutil.Big)
	for k, v := range m {
		mm[hexutil.Uint64(k)] = hexutil.Big(*v)
	}
	return json.Marshal(mm)
}
