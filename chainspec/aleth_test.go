package chainspec

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/params/types"
	"github.com/go-test/deep"
)

// Tests the go-ethereum to Aleth chainspec conversion for the Stureby testnet.
func TestAlethSturebyConverter(t *testing.T) {
	blob, err := ioutil.ReadFile("testdata/stureby_geth.json")
	if err != nil {
		t.Fatalf("could not read file: %v", err)
	}
	var genesis paramtypes.Genesis
	if err := json.Unmarshal(blob, &genesis); err != nil {
		t.Fatalf("failed parsing genesis: %v", err)
	}
	spec, err := NewAlethGenesisSpec("stureby", &genesis)
	if err != nil {
		t.Fatalf("failed creating chainspec: %v", err)
	}

	expBlob, err := ioutil.ReadFile("testdata/stureby_aleth.json")
	if err != nil {
		t.Fatalf("could not read file: %v", err)
	}
	expspec := &AlethGenesisSpec{}
	if err := json.Unmarshal(expBlob, expspec); err != nil {
		t.Fatalf("failed parsing genesis: %v", err)
	}
	// Compare the read-in SomeSpec vs. the generated SomeSpec
	if diffs := deep.Equal(expspec, spec); len(diffs) != 0 {
		t.Errorf("chainspec mismatch")
		for _, d := range diffs {
			t.Log(d)
		}
		spew.Dump(spec)
		bm, _ := json.MarshalIndent(spec, "", "    ")
		t.Log(string(bm))
	}
}
