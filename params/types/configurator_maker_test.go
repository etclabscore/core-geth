package paramtypes

import (
	"fmt"
	"math/big"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	common2 "github.com/ethereum/go-ethereum/params/types/common"
)

// File contains development-driven tests to reduce number
// of keystrokes.
// File does not contain any actual tests.

func TestQuickerGetterSetter_Inferface_Maker(t *testing.T) {
	type Ethasher struct {
		MinimumDifficulty           *big.Int
		DifficultyBoundDivisor      *big.Int
		HomesteadTransition         *big.Int
		EIP2Transition              *big.Int
		ECIP1010PauseTransition     *big.Int
		ECIP1010ContinueTransition  *big.Int
		ECIP1017Transition          *big.Int
		ECIP1017EraRounds           *big.Int
		EIP100BTransition           *big.Int
		ECIP1041Transition          *big.Int
		DifficultyBombDelaySchedule common2.Uint64BigMapEncodesHex
		BlockRewardSchedule         common2.Uint64BigMapEncodesHex
	}
	type Cliquer struct {
		Period uint64
		Epoch  uint64
	}
	type BlockSealer struct {
		Nonce   uint64
		MixHash common.Hash
	}
	type GenesisBlocker struct {
		BlockSealer

		Difficulty *big.Int
		Author     common.Address
		Timestamp  uint64
		ParentHash common.Hash
		ExtraData  common.Hash
		GasLimit   uint64
	}
	type CatHerder struct {
		AccountStartNonce         *uint64
		MaximumExtraDataSize      *uint64
		MinGasLimit               *uint64
		GasLimitBoundDivisor      *uint64
		NetworkID                 *uint64
		ChainID                   *uint64
		MaxCodeSize               *uint64
		MaxCodeSizeTransition     *uint64
		EIP7Transition            *uint64
		EIP98Transition           *uint64
		EIP150Transition          *uint64
		EIP160Transition          *uint64
		EIP161abcTransition       *uint64
		EIP161dTransition         *uint64
		EIP155Transition          *uint64
		EIP140Transition          *uint64
		EIP211Transition          *uint64
		EIP214Transition          *uint64
		EIP658Transition          *uint64
		EIP145Transition          *uint64
		EIP1014Transition         *uint64
		EIP1052Transition         *uint64
		EIP1283Transition         *uint64
		EIP1283DisableTransition  *uint64
		EIP1283ReenableTransition *uint64
		EIP1344Transition         *uint64
		EIP1884Transition         *uint64
		EIP2028Transition         *uint64
	}

	for _, tt := range []interface{}{
		Ethasher{},
		Cliquer{},
		BlockSealer{},
		GenesisBlocker{},
		CatHerder{},
	} {
		v := reflect.ValueOf(tt)
		fmt.Printf("type %s interface {\n", v.Type().Name())
		for i := 0; i < v.NumField(); i++ {
			//t.Log(v.Field(i).String())    // <*big.Int Value>
			//t.Log(v.Field(i).Interface()) // <nil>
			//t.Log(v.Field(i).Type())      // *big.Int
			//t.Log(v.Type().Field(i).Name) // Face

			fieldName := v.Type().Field(i).Name
			fieldTypeWord := v.Field(i).Type()

			fmt.Printf("Get%s() %s\n", fieldName, fieldTypeWord)
			fmt.Printf("Set%s(%s)\n", fieldName, fieldTypeWord)
		}
		fmt.Println("}")
	}
}
