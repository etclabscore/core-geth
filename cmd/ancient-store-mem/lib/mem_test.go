package lib

import (
	"bytes"
	"math/big"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/go-test/deep"
)

func TestMemFreezerRemoteServerAPI_Append(t *testing.T) {
	m := NewMemFreezerRemoteServerAPI()

	header := &types.Header{Number: big.NewInt(1), Difficulty: big.NewInt(42), Extra: []byte{}}
	bs, err := rlp.EncodeToBytes(header)
	if err != nil {
		t.Fatal(err)
	}
	if err := m.Append(freezerRemoteHeaderTable, 0, common.Bytes2Hex(bs)); err != nil {
		t.Fatal(err)
	}

	got, err := m.Ancient(freezerRemoteHeaderTable, 0)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(bs, got) {
		t.Fatal("!=")
	}
	decodeHeader := &types.Header{}
	if err := rlp.DecodeBytes(got, decodeHeader); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(header, decodeHeader) {
		t.Logf("decodeHeader: %v", decodeHeader)
		t.Error("! deepEqual")
		if diffs := deep.Equal(header, decodeHeader); len(diffs) != 0 {
			for i, diff := range diffs {
				t.Logf("diff#%d: %s", i, diff)
			}
		}
	}
}
