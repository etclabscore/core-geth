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
	"encoding/json"
	"fmt"
	"testing"
)

/*
TODO(meowsbits): Configs for reference.
var (
	mainnetChainConfig = params.ChainConfig{
		ChainID:        big.NewInt(1),
		HomesteadBlock: big.NewInt(1150000),
		DAOForkBlock:   big.NewInt(1920000),
		DAOForkSupport: true,
		EIP150Block:    big.NewInt(2463000),
		EIP155Block:    big.NewInt(2675000),
		EIP158Block:    big.NewInt(2675000),
		ByzantiumBlock: big.NewInt(4370000),
	}

	ropstenChainConfig = params.ChainConfig{
		ChainID:                       big.NewInt(3),
		HomesteadBlock:                big.NewInt(0),
		DAOForkBlock:                  nil,
		DAOForkSupport:                true,
		EIP150Block:                   big.NewInt(0),
		EIP155Block:                   big.NewInt(10),
		EIP158Block:                   big.NewInt(10),
		ByzantiumBlock:                big.NewInt(1_700_000),
		ConstantinopleBlock:           big.NewInt(4_230_000),
		PetersburgBlock:               big.NewInt(4_939_394),
		IstanbulBlock:                 big.NewInt(6_485_846),
		MuirGlacierBlock:              big.NewInt(7_117_117),
		BerlinBlock:                   big.NewInt(9_812_189),
		LondonBlock:                   big.NewInt(10_499_401),
		TerminalTotalDifficulty:       new(big.Int).SetUint64(50_000_000_000_000_000),
		TerminalTotalDifficultyPassed: true,
	}
)
*/

func TestDifficulty(t *testing.T) {
	t.Parallel()

	dt := new(testMatcher)

	for k, v := range difficultyChainConfigurations {
		dt.config(k, v)
	}

	for _, dir := range []string{
		difficultyTestDir,
		difficultyTestDirETC,
	} {
		dt.walk(t, dir, func(t *testing.T, name string, superTest map[string]json.RawMessage) {
			for fork, rawTests := range superTest {
				if fork == "_info" {
					continue
				}
				var tests map[string]DifficultyTest
				if err := json.Unmarshal(rawTests, &tests); err != nil {
					t.Error(err)
					continue
				}

				cfg, ok := Forks[fork]
				if !ok {
					t.Error(UnsupportedForkError{fork})
					continue
				}

				for subname, subtest := range tests {
					key := fmt.Sprintf("%s/%s", fork, subname)
					t.Run(key, func(t *testing.T) {
						if err := dt.checkFailure(t, subtest.Run(cfg)); err != nil {
							t.Error(err)
						}
					})
				}
			}
		})
	}
}
