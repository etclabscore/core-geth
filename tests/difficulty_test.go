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
	"testing"

	"github.com/ethereum/go-ethereum/params/vars"
)

func TestDifficulty(t *testing.T) {
	t.Parallel()

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

	for k, v := range difficultyChainConfigurations {
		dt.config(k, v)
	}

	dt.walk(t, difficultyTestDir, func(t *testing.T, name string, test *DifficultyTest) {
		cfg, _ := dt.findConfig(name)

		if test.ParentDifficulty.Cmp(vars.MinimumDifficulty) < 0 {
			t.Skip("difficulty below minimum")
			return
		}
		if err := dt.checkFailure(t, name, test.Run(cfg)); err != nil {
			t.Error(err)
		}
	})
}

func TestDifficultyNDJSON(t *testing.T) {
	t.Parallel()

	dt := new(testMatcher)

	// Not NDJSON
	dt.skipLoad(`\\.json$`)

	for k, v := range difficultyChainConfigurations {
		dt.config(k, v)
	}

	dt.walkScanNDJSON(t, difficultyTestDir, func(t *testing.T, name string, test *DifficultyTest) {
		// Kind of ugly reverse lookup from file -> fork name.
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

		cfg, _ := dt.findConfig(forkName)
		if test.ParentDifficulty.Cmp(vars.MinimumDifficulty) < 0 {
			t.Skip("difficulty below minimum")
			return
		}
		if err := dt.checkFailure(t, name, test.Run(cfg)); err != nil {
			t.Error(err)
		}
	})
}
