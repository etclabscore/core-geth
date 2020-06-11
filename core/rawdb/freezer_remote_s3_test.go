package rawdb

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/rlp"
)

func TestDifficultyRLPDecoding(t *testing.T) {
	td := big.NewInt(42)
	b, err := rlp.EncodeToBytes(td)
	if err != nil {
		t.Fatal(err)
	}
	difficulty := new(big.Int)
	err = rlp.DecodeBytes(b, difficulty)
	if err != nil {
		t.Error(err)
	}
	t.Log(difficulty)
}
