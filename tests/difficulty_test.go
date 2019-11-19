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

package tests

import (
	"bufio"
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/params"
)

func TestDifficulty(t *testing.T) {
	generateTests := os.Getenv(MG_GENERATE_DIFFICULTY_TESTS_KEY) != ""

	if !generateTests {
		t.Parallel()
	} else {
		err := os.MkdirAll(filepath.Join(difficultyTestDir, "generated_difficulty"), os.ModePerm)
		if err != nil {
			t.Fatal(err)
		}
	}

	dt := new(testMatcher)

	// Not difficulty-tests
	dt.skipLoad("hexencodetest.*")
	dt.skipLoad("crypto.*")
	dt.skipLoad("blockgenesistest\\.json")
	dt.skipLoad("genesishashestest\\.json")
	dt.skipLoad("keyaddrtest\\.json")
	dt.skipLoad("txtest\\.json")

	// files are 2 years old, contains strange values
	dt.skipLoad("difficultyCustomHomestead\\.json")
	dt.skipLoad("difficultyMorden\\.json")
	dt.skipLoad("difficultyOlimpic\\.json")
	dt.skipLoad("all_difficulty_tests\\.json")

	for k, v := range difficultyChainConfiguations {
		dt.config(k, v)
	}

	// If there is a generated NDJSON difficulty test file available, then use that for testing.
	// This is a minimum necessary dogfooding for the test generation.
	// Note that if this file is indeed available, the tests will run
	if fi, err := os.Open(filepath.Join(difficultyTestDir, "generated_difficulty", "all_difficulty_tests.json")); err == nil {
		scanner := bufio.NewScanner(fi)
		for scanner.Scan() {
			test := &DifficultyTest{}
			err := json.Unmarshal(scanner.Bytes(), &test)
			if err != nil {
				t.Fatal(err)
			}
			var forkName string
			for k, v := range mapForkNameChainspecFileDifficulty {
				if v == test.Chainspec.Filename {
					forkName = k
					break
				}
			}
			if forkName == "" {
				t.Fatal("missing fork/fileconf name", test, mapForkNameChainspecFileDifficulty)
			}
			conf, ok := difficultyChainConfiguations[forkName]
			if !ok {
				t.Fatal("missing chain configurations for forkname: ", forkName)
			}

			// If tests are using JSON chainspecs, then
			// reverse lookup and verify the chain config given the JSON spec filename.
			if os.Getenv(MG_CHAINCONFIG_CHAINSPEC_KEY) != "" {
				cr, ok := chainspecRefsDifficulty[forkName]
				if !ok {
					t.Fatalf("missing chainconfig: %v", forkName)
				}
				if !bytes.Equal(cr.Sha1Sum, test.Chainspec.Sha1Sum) {
					t.Fatalf("mismatch configs, test: %v, spec config sum: %x", test, cr.Sha1Sum)
				}
			}
			t.Run(test.Name, func(t *testing.T) {
				err = test.Run(&conf)
				if err != nil {
					t.Errorf("test: %v err: %v", test, err)
					return
				}
				t.Logf("OK %v", test)
			})
		}
		return
	}

	dt.walk(t, difficultyTestDir, func(t *testing.T, name string, test *DifficultyTest) {
		cfg, key := dt.findConfig(name)
		test.Chainspec = chainspecRefsDifficulty[key]
		if test.ParentDifficulty.Cmp(params.MinimumDifficulty) < 0 {
			t.Skip("difficulty below minimum")
			return
		}
		if err := dt.checkFailure(t, name, test.Run(cfg)); err != nil {
			t.Error(err)
		} else {
			if !generateTests {
				return
			}
			test.Name = strings.ReplaceAll(name, ".json", "")
			b, _ := json.Marshal(test)
			out := []byte{}
			buf := bytes.NewBuffer(out)
			err = json.Compact(buf, b)
			if err != nil {
				t.Fatal(err)
			}
			buf.Write([]byte("\n"))

			fn := filepath.Join(difficultyTestDir, "generated_difficulty", "all_difficulty_tests.json")
			fi, err := os.OpenFile(fn, os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.ModePerm)
			if err != nil {
				t.Fatal(err)
				return
			}
			_, err = fi.Write(buf.Bytes())
			if err != nil {
				t.Fatal(err)
			}
			err = fi.Close()
			if err != nil {
				t.Fatal(err)
			}
		}
	})

	if !generateTests {
		return
	}
}
