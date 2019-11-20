package parity

import (
	"bytes"
	"encoding/json"
	"math/big"
	"reflect"
	"strings"
	"testing"
)

// Tests for map data types.

type fakeConfig struct {
	Number Uint64BigValOrMapHex `json:"num"`
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
	rawjson []byte
	dat fakeConfig
	marshaledWant []byte
}

var testCases = []testCase{
	{
		uint64bigMaybeNoD,
		fakeConfig{Uint64BigValOrMapHex{0: big.NewInt(2000000000000000000)}},
		uint64bigMaybeNoDMarshaledMap,
	},
	{
		uint64bigMaybeYesD,
		fakeConfig{Uint64BigValOrMapHex{0: big.NewInt(2000000000000000000), 5: big.NewInt(3000000000000000000)}},
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
		Nums Uint64BigMapEncodesHex `json:"num"`
	}
	c := conf{}
	err := json.Unmarshal(uint64bigMaybeYesD, &c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Nums[0].Cmp(big.NewInt(2000000000000000000)) != 0 {
		t.Error("mismatch")
	}
}

func TestBigMapEncodesHex_MarshalJSON(t *testing.T) {
	type conf struct {
		Nums Uint64BigMapEncodesHex `json:"num"`
	}
	c := conf{Uint64BigMapEncodesHex{0: big.NewInt(2000000000000000000), 5: big.NewInt(3000000000000000000)}}
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
