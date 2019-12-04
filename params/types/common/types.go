package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"sort"
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
	if err == nil {
		return false
	}
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
		case int, int64:
			vv = big.NewInt(v.(int64))
		case float64:
			i, err := strconv.ParseUint(fmt.Sprintf("%.0f", v), 10, 64)
			if err != nil {
				panic(err)
			}
			vv = big.NewInt(int64(i))
		default:
			panic(fmt.Sprintf("unknown type: %t %v", v, v))
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

// ExtractHostageSituationN returns the block number for a given hostage situation, ie EIP649 and/or EIP1234.
// This is a reverse lookup to extract EIP-spec'd parameters from difficulty and reward maps implementations.
func ExtractHostageSituationN(difficulties Uint64BigMapEncodesHex, rewards Uint64BigMapEncodesHex, difficultySum, wantedReward *big.Int) *uint64 {
	var diffN *uint64
	var sl = []uint64{}

	// difficulty
	for k, _ := range difficulties {
		sl = append(sl, k)
	}
	sort.Slice(sl, func(i, j int) bool {
		return sl[i] < sl[j]
	})

	var total = new(big.Int)
	for _, s := range sl {
		d :=  difficulties[s]
		if d == nil {
			panic(fmt.Sprintf("dnil difficulties: %v, sl: %v", difficulties, sl))
		}
		total.Add(total, d)
		if total.Cmp(difficultySum) == 0 {
			diffN = &s
			break
		}
	}
	if diffN == nil {
		// difficulty bomb delay not configured,
		// then does not meet eip649/eip1234 spec
		return nil
	}

	reward, ok := rewards[*diffN]
	if !ok {
		return nil
	}
	if reward.Cmp(wantedReward) != 0 {
		return nil
	}

	return diffN
}

type ConsensusEngineT int

const (
	ConsensusEngineT_Unknown = iota
	ConsensusEngineT_Ethash
	ConsensusEngineT_Clique
)

func (c ConsensusEngineT) String() string {
	switch c {
	case ConsensusEngineT_Ethash:
		return "ethash"
	case ConsensusEngineT_Clique:
		return "clique"
	default:
		return "unknown"
	}
}

func (c ConsensusEngineT) IsEthash() bool {
	return c == ConsensusEngineT_Ethash
}

func (c ConsensusEngineT) IsClique() bool {
	return c == ConsensusEngineT_Clique
}

type BlockSealingT int

const (
	BlockSealing_Unknown = iota
	BlockSealing_Ethereum
)

func (b BlockSealingT) String() string {
	switch b {
	case BlockSealing_Ethereum:
		return "ethereum"
	default:
		return "unknown"
	}
}

