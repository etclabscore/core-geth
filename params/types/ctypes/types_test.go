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

package ctypes

import (
	"bytes"
	"encoding/json"
	"math"
	"math/big"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params/vars"
	"github.com/holiman/uint256"
)

// Tests for map data types.

type fakeConfig struct {
	Number Uint64Uint256ValOrMapHex `json:"num"`
}

var uint64bigMaybeNoD []byte = []byte(`
{
	"num": "0x1BC16D674EC80000"
}`)

var uint64bigMaybeNoDMarshaledMap []byte = []byte(`
{
	"num": {"0x0": "0x1BC16D674EC80000"}
}`)

var uint64bigMaybeYesD []byte = []byte(`
{
	"num": {
		"0x0": "0x1BC16D674EC80000",
		"0x5": "0x29A2241AF62C0000"
	}
}`)

type testCase struct {
	rawjson       []byte
	dat           fakeConfig
	marshaledWant []byte
}

var testCases = []testCase{
	{
		uint64bigMaybeNoD,
		fakeConfig{Uint64Uint256ValOrMapHex{0: uint256.NewInt(2000000000000000000)}},
		uint64bigMaybeNoDMarshaledMap,
	},
	{
		uint64bigMaybeYesD,
		fakeConfig{Uint64Uint256ValOrMapHex{0: uint256.NewInt(2000000000000000000), 5: uint256.NewInt(3000000000000000000)}},
		uint64bigMaybeYesD,
	},
}

func TestUint64BigMapMaybe_UnmarshalJSON(t *testing.T) {
	var err error
	for i, c := range testCases {
		unMarshaled := fakeConfig{}
		err = json.Unmarshal(c.rawjson, &unMarshaled)
		eq := reflect.DeepEqual(unMarshaled, c.dat)
		if err != nil || !eq {
			t.Log(string(c.rawjson))
		}
		if err != nil {
			t.Errorf("case.i=%d error: %s", i, err)
		}
		if !eq {
			t.Errorf("case.i=%d got: %v, want: %v", i, unMarshaled, c.dat)
		}
	}
}

func TestUint64BigMapMaybe_MarshalJSON(t *testing.T) {
	for i, c := range testCases {
		got, err := json.Marshal(c.dat)
		if err != nil {
			t.Errorf("case.i=%d error: %s", i, err)
		}
		gotb := new(bytes.Buffer)
		wantb := new(bytes.Buffer)
		if err := json.Compact(gotb, got); err != nil {
			t.Fatal(err)
		}
		if err := json.Compact(wantb, c.marshaledWant); err != nil {
			t.Fatal(err)
		}
		gots := strings.ToLower(gotb.String())
		wants := strings.ToLower(wantb.String())
		if gots != wants {
			t.Errorf("case.i=%d got: %s, want: %s", i, gots, wants)
		}
	}
}

func TestBigMapEncodesHex_UnmarshalJSON(t *testing.T) {
	type conf struct {
		Nums Uint64Uint256MapEncodesHex `json:"num"`
	}
	c := conf{}
	err := json.Unmarshal(uint64bigMaybeYesD, &c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Nums[0].Cmp(uint256.NewInt(2000000000000000000)) != 0 {
		t.Error("mismatch")
	}
}

func TestBigMapEncodesHex_MarshalJSON(t *testing.T) {
	type conf struct {
		Nums Uint64Uint256MapEncodesHex `json:"num"`
	}
	c := conf{Uint64Uint256MapEncodesHex{0: uint256.NewInt(2000000000000000000), 5: uint256.NewInt(3000000000000000000)}}
	got, err := json.Marshal(c)
	if err != nil {
		t.Fatal(err)
	}
	gotb := new(bytes.Buffer)
	wantb := new(bytes.Buffer)
	if err := json.Compact(gotb, got); err != nil {
		t.Fatal(err)
	}
	if err := json.Compact(wantb, uint64bigMaybeYesD); err != nil {
		t.Fatal(err)
	}
	gots := strings.ToLower(gotb.String())
	wants := strings.ToLower(wantb.String())
	if gots != wants {
		t.Errorf("got: %s, want: %s", gots, wants)
	}
}

func TestUint64BigMapEncodesHex_SetValueTotalForHeight(t *testing.T) {
	newMG := func() Uint64Uint256MapEncodesHex {
		v := Uint64Uint256MapEncodesHex{}
		return v
	}
	byzaBlock := big.NewInt(4370000).Uint64()
	consBlock := big.NewInt(7280000).Uint64()
	muirBlock := big.NewInt(9200000).Uint64()

	max := uint64(math.MaxUint64)

	check := func(mg Uint64Uint256MapEncodesHex, got, want uint64) {
		if got != want {
			t.Log(runtime.Caller(1))
			t.Log(runtime.Caller(2))
			t.Errorf("got: %d, want: %d", got, want)
			for k, v := range mg {
				t.Logf("%d: %d", k, v.Uint64())
			}
			t.Log("---")
		}
	}

	// Test set one (latest) value, eg for testing or new chains.
	mgSoloOrdered := newMG()
	mgSoloOrdered.SetValueTotalForHeight(&muirBlock, vars.EIP2384DifficultyBombDelay)
	check(mgSoloOrdered, mgSoloOrdered.SumValues(&muirBlock), vars.EIP2384DifficultyBombDelay.Uint64())
	check(mgSoloOrdered, mgSoloOrdered.SumValues(&max), vars.EIP2384DifficultyBombDelay.Uint64())

	checkFinal := func(mg Uint64Uint256MapEncodesHex) {
		check(mg, mg.SumValues(&byzaBlock), vars.EIP649DifficultyBombDelay.Uint64())
		check(mg, mg.SumValues(&consBlock), vars.EIP1234DifficultyBombDelay.Uint64())
		check(mg, mg.SumValues(&muirBlock), vars.EIP2384DifficultyBombDelay.Uint64())
		check(mg, mg.SumValues(&max), vars.EIP2384DifficultyBombDelay.Uint64())
	}

	// Test set fork values in chronological order.
	mgBasicOrdered := newMG()
	mgBasicOrdered.SetValueTotalForHeight(&byzaBlock, vars.EIP649DifficultyBombDelay)
	check(mgBasicOrdered, mgBasicOrdered.SumValues(&max), vars.EIP649DifficultyBombDelay.Uint64())
	mgBasicOrdered.SetValueTotalForHeight(&consBlock, vars.EIP1234DifficultyBombDelay)
	mgBasicOrdered.SetValueTotalForHeight(&muirBlock, vars.EIP2384DifficultyBombDelay)

	checkFinal(mgBasicOrdered)

	// Test set all features, unordered, 1
	mgBasicUnordered := newMG()
	mgBasicUnordered.SetValueTotalForHeight(&muirBlock, vars.EIP2384DifficultyBombDelay)
	mgBasicUnordered.SetValueTotalForHeight(&consBlock, vars.EIP1234DifficultyBombDelay)
	mgBasicUnordered.SetValueTotalForHeight(&byzaBlock, vars.EIP649DifficultyBombDelay)

	checkFinal(mgBasicUnordered)

	// Same, but again, more.
	mgBasicUnordered2 := newMG()
	mgBasicUnordered2.SetValueTotalForHeight(&muirBlock, vars.EIP2384DifficultyBombDelay)
	mgBasicUnordered2.SetValueTotalForHeight(&byzaBlock, vars.EIP649DifficultyBombDelay)
	mgBasicUnordered2.SetValueTotalForHeight(&consBlock, vars.EIP1234DifficultyBombDelay)

	checkFinal(mgBasicUnordered2)

	// Test set all features, unordered, and with edge cases, 2
	mgWildUnordered2 := newMG()
	mgWildUnordered2.SetValueTotalForHeight(&muirBlock, vars.EIP2384DifficultyBombDelay)
	mgWildUnordered2.SetValueTotalForHeight(&byzaBlock, vars.EIP649DifficultyBombDelay)
	mgWildUnordered2.SetValueTotalForHeight(&consBlock, vars.EIP1234DifficultyBombDelay)

	// Set a dupe.
	mgWildUnordered2.SetValueTotalForHeight(&byzaBlock, vars.EIP649DifficultyBombDelay)

	// Set a random.
	randoK := new(big.Int).Div(new(big.Int).Add(big.NewInt(int64(byzaBlock)), big.NewInt(int64(consBlock))), common.Big2).Uint64()
	randoV := new(uint256.Int).Div(new(uint256.Int).Add(vars.EIP649DifficultyBombDelay, vars.EIP1234DifficultyBombDelay), uint256.NewInt(2))
	mgWildUnordered2.SetValueTotalForHeight(&randoK, randoV)

	checkFinal(mgWildUnordered2)

	// Test repetitious set's.
	mgRepetitious := newMG()
	mgRepetitious.SetValueTotalForHeight(&byzaBlock, vars.EIP649DifficultyBombDelay)
	mgRepetitious.SetValueTotalForHeight(&consBlock, vars.EIP1234DifficultyBombDelay)
	mgRepetitious.SetValueTotalForHeight(&muirBlock, vars.EIP2384DifficultyBombDelay)
	mgRepetitious.SetValueTotalForHeight(&consBlock, vars.EIP1234DifficultyBombDelay)
	mgRepetitious.SetValueTotalForHeight(&consBlock, vars.EIP1234DifficultyBombDelay)
	mgRepetitious.SetValueTotalForHeight(&byzaBlock, vars.EIP649DifficultyBombDelay)
	mgRepetitious.SetValueTotalForHeight(&byzaBlock, vars.EIP649DifficultyBombDelay)

	checkFinal(mgRepetitious)

	mgTestlike := newMG()
	zero := uint64(0)
	five := uint64(5)
	mgTestlike.SetValueTotalForHeight(&zero, vars.EIP649DifficultyBombDelay)
	mgTestlike.SetValueTotalForHeight(&five, vars.EIP1234DifficultyBombDelay)
	check(mgTestlike, mgTestlike.SumValues(&zero), vars.EIP649DifficultyBombDelay.Uint64())
}

func TestMapMeetsSpecification_1234(t *testing.T) {
	data := `{
		"difficultyBombDelays": {
            "0x0": "0x4c4b40"
        },
        "blockReward": {
            "0x0": "0x1bc16d674ec80000"
        }}`

	im := struct {
		DifficultyBombDelaySchedule Uint64Uint256MapEncodesHex `json:"difficultyBombDelays,omitempty"` // JSON tag matches Parity's
		BlockRewardSchedule         Uint64Uint256MapEncodesHex `json:"blockReward,omitempty"`          // JSON tag matches Parity's
	}{}
	err := json.Unmarshal([]byte(data), &im)
	if err != nil {
		t.Fatal(err)
	}

	n := MapMeetsSpecification(
		im.DifficultyBombDelaySchedule,
		im.BlockRewardSchedule,
		vars.EIP1234DifficultyBombDelay,
		vars.EIP1234FBlockReward,
	)
	if n == nil || *n != 0 {
		t.Fatal("n should be 0", n)
	}

	t.Logf("%v n=%v", im, n)
}
