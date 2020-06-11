package rawdb

import (
	"log"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
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

func TestHashRLPCodec(t *testing.T) {
	h1 := common.HexToHash("0xbadface")
	b, err := rlp.EncodeToBytes(h1)
	if err != nil {
		t.Fatal(err)
	}
	h2 := common.Hash{}
	err = rlp.DecodeBytes(b, &h2)
	if err != nil {
		t.Error(err)
	}
	log.Println(h2.Hex())
}

func TestBlockKeying(t *testing.T) {
	m := 32
	n := 17
	mod := n % m
	div := n / m
	t.Log(mod)
	t.Log(div)
	t.Log(n * div + mod)
}

func TestIndexThing(t *testing.T) {
	a := []int{1, 2, 3}
	t.Log(a[:0])
}

func TestSpliceBackwards(t *testing.T) {
	a := []int{1, 2, 3, 4, 5}
	mod := 3
	rem := len(a) % mod
	t.Log(a[len(a)-rem:])
}