package chainspec

import (
	"bytes"
	"encoding/json"
	"math/big"
	"reflect"
	"strings"
	"testing"
)

type fakeConfig struct {
	Number Uint64BigValOrMapHex `json:"num"`
}

var uint64bigMaybeNoD []byte = []byte(`
{
	"num": "0x1BC16D674EC80000"
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
}

var testCases = []testCase{
	{
		uint64bigMaybeNoD,
		fakeConfig{Uint64BigValOrMapHex{0: big.NewInt(2000000000000000000)}},
	},
	{
		uint64bigMaybeYesD,
		fakeConfig{Uint64BigValOrMapHex{0: big.NewInt(2000000000000000000), 5: big.NewInt(3000000000000000000)}},
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
		if err := json.Compact(wantb, c.rawjson); err != nil {
			t.Fatal(err)
		}
		gots := strings.ToLower(gotb.String())
		wants := strings.ToLower(wantb.String())
		if gots != wants {
			t.Errorf("case.i=%d got: %s, want: %s", i, gots, wants)
		}
	}
}