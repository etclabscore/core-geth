// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package ethash

import (
	crand "crypto/rand"
	"encoding/binary"
	"encoding/json"
	"math/big"
	"math/rand"
	"os"
	"path/filepath"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params/types/goethereum"
	"github.com/ethereum/go-ethereum/params/vars"
)

type diffTest struct {
	ParentTimestamp    uint64
	ParentDifficulty   *big.Int
	CurrentTimestamp   uint64
	CurrentBlocknumber *big.Int
	CurrentDifficulty  *big.Int
}

func (d *diffTest) UnmarshalJSON(b []byte) (err error) {
	var ext struct {
		ParentTimestamp    string
		ParentDifficulty   string
		CurrentTimestamp   string
		CurrentBlocknumber string
		CurrentDifficulty  string
	}
	if err := json.Unmarshal(b, &ext); err != nil {
		return err
	}

	d.ParentTimestamp = math.MustParseUint64(ext.ParentTimestamp)
	d.ParentDifficulty = math.MustParseBig256(ext.ParentDifficulty)
	d.CurrentTimestamp = math.MustParseUint64(ext.CurrentTimestamp)
	d.CurrentBlocknumber = math.MustParseBig256(ext.CurrentBlocknumber)
	d.CurrentDifficulty = math.MustParseBig256(ext.CurrentDifficulty)

	return nil
}

func TestCalcDifficulty(t *testing.T) {
	file, err := os.Open(filepath.Join("..", "..", "tests", "testdata", "BasicTests", "difficulty.json"))
	if err != nil {
		t.Skip(err)
	}
	defer file.Close()

	tests := make(map[string]diffTest)
	err = json.NewDecoder(file).Decode(&tests)
	if err != nil {
		t.Fatal(err)
	}

	config := &goethereum.ChainConfig{HomesteadBlock: big.NewInt(1150000)}

	for name, test := range tests {
		number := new(big.Int).Sub(test.CurrentBlocknumber, big.NewInt(1))
		diff := CalcDifficulty(config, test.CurrentTimestamp, &types.Header{
			Number:     number,
			Time:       test.ParentTimestamp,
			Difficulty: test.ParentDifficulty,
		})
		if diff.Cmp(test.CurrentDifficulty) != 0 {
			t.Error(name, "failed. Expected", test.CurrentDifficulty, "and calculated", diff)
		}
	}
}

func randSlice(min, max uint32) []byte {
	var b = make([]byte, 4)
	crand.Read(b)
	a := binary.LittleEndian.Uint32(b)
	size := min + a%(max-min)
	out := make([]byte, size)
	crand.Read(out)
	return out
}

func TestDifficultyCalculators(t *testing.T) {
	for i := 0; i < 5000; i++ {
		// 1 to 300 seconds diff
		var timeDelta = uint64(1 + rand.Uint32()%3000)
		diffBig := new(big.Int).SetBytes(randSlice(2, 10))
		if diffBig.Cmp(vars.MinimumDifficulty) < 0 {
			diffBig.Set(vars.MinimumDifficulty)
		}
		// rand.Read(difficulty)
		header := &types.Header{
			Difficulty: diffBig,
			Number:     new(big.Int).SetUint64(rand.Uint64() % 50_000_000),
			Time:       rand.Uint64() - timeDelta,
		}
		if rand.Uint32()&1 == 0 {
			header.UncleHash = types.EmptyUncleHash
		}
		bombDelay := new(big.Int).SetUint64(rand.Uint64() % 50_000_000)
		for i, pair := range []struct {
			bigFn  func(time uint64, parent *types.Header) *big.Int
			u256Fn func(time uint64, parent *types.Header) *big.Int
		}{
			{FrontierDifficultyCalculator, CalcDifficultyFrontierU256},
			{HomesteadDifficultyCalculator, CalcDifficultyHomesteadU256},
			{DynamicDifficultyCalculator(bombDelay), MakeDifficultyCalculatorU256(bombDelay)},
		} {
			time := header.Time + timeDelta
			want := pair.bigFn(time, header)
			have := pair.u256Fn(time, header)
			if want.BitLen() > 256 {
				continue
			}
			if want.Cmp(have) != 0 {
				t.Fatalf("pair %d: want %x have %x\nparent.Number: %x\np.Time: %x\nc.Time: %x\nBombdelay: %v\n", i, want, have,
					header.Number, header.Time, time, bombDelay)
			}
		}
	}
}

func BenchmarkDifficultyCalculator(b *testing.B) {
	x1 := MakeDifficultyCalculatorU256(big.NewInt(1000000))
	x2 := MakeDifficultyCalculatorU256(big.NewInt(1000000))
	h := &types.Header{
		ParentHash: common.Hash{},
		UncleHash:  types.EmptyUncleHash,
		Difficulty: big.NewInt(0xffffff),
		Number:     big.NewInt(500000),
		Time:       1000000,
	}
	b.Run("big-frontier", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			CalcDifficultyFrontierU256(1000014, h)
		}
	})
	b.Run("u256-frontier", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			CalcDifficultyFrontierU256(1000014, h)
		}
	})
	b.Run("big-homestead", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			CalcDifficultyHomesteadU256(1000014, h)
		}
	})
	b.Run("u256-homestead", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			CalcDifficultyHomesteadU256(1000014, h)
		}
	})
	b.Run("big-generic", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			x1(1000014, h)
		}
	})
	b.Run("u256-generic", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			x2(1000014, h)
		}
	})
}

// TestEtchash_11700000 tests that Etchash can verify the block at 11700000.
// EtcHash is a fork of Ethash developed by Ethereum Classic which doubles the epoch length,
// approximately halving the DAG size over time.
// Block 11_700_000 is the block at which ETC implemented this hard-fork feature.
// This test takes a long time to run; about 2 minutes on my machine.
// It builds a real DAG and verifies the canonical block header.
func TestEtchash_11700000(t *testing.T) {
	// canon11700000 is a JSONRPC response to a getBlockByNumber request for block 11700000 for the ETC chain from rivet.cloud.
	// Unmarshal this data into a block header.
	canon11700000 := `{"jsonrpc":"2.0","id":3019187,"result":{"difficulty":"0x317df020bde4","extraData":"0x7374726174756d2d65752d31","gasLimit":"0x7a6d21","gasUsed":"0x0","hash":"0x3cdbd6d80f7fd983ac2a6b13230f56c61dc58cef3deeaeaffda36ca5efc6cb31","logsBloom":"0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000","miner":"0xdf7d7e053933b5cc24372f878c90e62dadad5d42","mixHash":"0x5f688fe8043392fd60d661edd5744dcb3d557042d3e472e7a708294c34216fb8","nonce":"0x69710a743181b923","number":"0xb28720","parentHash":"0xadaa50fd7a72baa19df873163b563c7cfa5ad254ca8dcf3b019076499e219e2e","receiptsRoot":"0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421","sha3Uncles":"0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347","size":"0x211","stateRoot":"0x373bd9a0991eff6c03f5ead0f25a0466a2ae10cf80cb37bd51d144f68668af29","timestamp":"0x5fc2a999","totalDifficulty":"0x388236832f6f1333ed","transactions":[],"transactionsRoot":"0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421","uncles":[]}}`
	tmpRPCRes := struct {
		Result types.Header
	}{}
	if err := json.Unmarshal([]byte(canon11700000), &tmpRPCRes); err != nil {
		t.Fatal(err)
	}
	canonHeader := tmpRPCRes.Result

	// The block number is 11700000.
	forkBlock := big.NewInt(11700000)
	forkBlockU := forkBlock.Uint64()
	if canonHeader.Number.Cmp(forkBlock) != 0 {
		t.Fatalf("wrong block number: %v", canonHeader.Number)
	}

	// Set up the ethash config, defining
	// the cache and DAG as namespaced directories within the os temp dir.
	temp := filepath.Join(os.TempDir(), "/geth-test-etchash-11700000")
	ethashConfig := Config{
		CacheDir:         filepath.Join(temp, "etcash-cache"),
		CachesInMem:      2,
		CachesOnDisk:     3,
		CachesLockMmap:   false,
		DatasetDir:       filepath.Join(temp, "etcash-dag"),
		DatasetsInMem:    1,
		DatasetsOnDisk:   2,
		DatasetsLockMmap: false,
		ECIP1099Block:    &forkBlockU,
	}
	// Clean up ahead of ourselves to prevent interference.
	os.RemoveAll(ethashConfig.CacheDir)
	os.RemoveAll(ethashConfig.DatasetDir)

	// Create a new instance of *Ethash and use it to
	// test that verifySeal confirms this is a valid header.
	ethash := New(ethashConfig, nil, false)
	defer ethash.Close()
	if err := ethash.verifySeal(nil, &canonHeader, true); err != nil {
		t.Fatalf("verifySeal failed: %v", err)
	}
}
