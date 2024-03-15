// Copyright 2017 The go-ethereum Authors
// This file is part of go-ethereum.
//
// go-ethereum is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// go-ethereum is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with go-ethereum. If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/state/snapshot"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/eth/tracers/logger"
	"github.com/ethereum/go-ethereum/internal/flags"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/tests"
	"github.com/urfave/cli/v2"
)

var stateTestCommand = &cli.Command{
	Action:    stateTestCmd,
	Name:      "statetest",
	Usage:     "Executes the given state tests. Filenames can be fed via standard input (batch mode) or as an argument (one-off execution).",
	ArgsUsage: "<file>",
	Flags: []cli.Flag{
		stateTestForkFlag,
		stateTestEVMCEWASMFlag,
		utils.EVMInterpreterFlag,
	},
	Category: flags.DevCategory,
}

var stateTestEVMCEWASMFlag = &cli.StringFlag{
	Name:     "evmc.ewasm",
	Usage:    "EVMC EWASM configuration",
	Category: flags.DevCategory,
}

var stateTestForkFlag = &cli.StringFlag{
	Name:     "fork",
	Usage:    "Fork to use for the test",
	Category: flags.DevCategory,
}

// StatetestResult contains the execution status after running a state test, any
// error that might have occurred and a dump of the final state if requested.
type StatetestResult struct {
	Name  string       `json:"name"`
	Pass  bool         `json:"pass"`
	Root  *common.Hash `json:"stateRoot,omitempty"`
	Fork  string       `json:"fork"`
	Error string       `json:"error,omitempty"`
	State *state.Dump  `json:"state,omitempty"`
}

func stateTestCmd(ctx *cli.Context) error {
	// Configure the EVM logger
	config := &logger.Config{
		EnableMemory:     !ctx.Bool(DisableMemoryFlag.Name),
		DisableStack:     ctx.Bool(DisableStackFlag.Name),
		DisableStorage:   ctx.Bool(DisableStorageFlag.Name),
		EnableReturnData: !ctx.Bool(DisableReturnDataFlag.Name),
	}
	var cfg vm.Config
	switch {
	case ctx.Bool(MachineFlag.Name):
		cfg.Tracer = logger.NewJSONLogger(config, os.Stderr)

	case ctx.Bool(DebugFlag.Name):
		cfg.Tracer = logger.NewStructLogger(config)
	}

	cfg.EWASMInterpreter = ctx.String(stateTestEVMCEWASMFlag.Name)
	cfg.EVMInterpreter = ctx.String(utils.EVMInterpreterFlag.Name)

	if cfg.EVMInterpreter != "" {
		log.Info("Running tests with %s=%s", "evmc.evm", cfg.EVMInterpreter)
		vm.InitEVMCEVM(cfg.EVMInterpreter)
	}
	if cfg.EWASMInterpreter != "" {
		log.Info("Running tests with %s=%s", "evmc.ewasm", cfg.EWASMInterpreter)
		vm.InitEVMCEwasm(cfg.EWASMInterpreter)
	}

	// Load the test content from the input file
	if len(ctx.Args().First()) != 0 {
		return runStateTest(ctx.Args().First(), cfg, ctx.Bool(MachineFlag.Name), ctx.Bool(DumpFlag.Name), ctx.String(stateTestForkFlag.Name))
	}
	// Read filenames from stdin and execute back-to-back
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		fname := scanner.Text()
		if len(fname) == 0 {
			return nil
		}
		if err := runStateTest(fname, cfg, ctx.Bool(MachineFlag.Name), ctx.Bool(DumpFlag.Name), ctx.String(stateTestForkFlag.Name)); err != nil {
			return err
		}
	}
	return nil
}

// runStateTest loads the state-test given by fname, and executes the test.
func runStateTest(fname string, cfg vm.Config, jsonOut, dump bool, testFork string) error {
	src, err := os.ReadFile(fname)
	if err != nil {
		return err
	}
	var tests map[string]tests.StateTest
	if err := json.Unmarshal(src, &tests); err != nil {
		return err
	}
	// Iterate over all the tests, run them and aggregate the results
	results := make([]StatetestResult, 0, len(tests))
	for key, test := range tests {
		for _, st := range test.Subtests(nil) {
			if testFork != "" && testFork != st.Fork {
				continue
			}
			// Run the test and aggregate the result
			result := &StatetestResult{Name: key, Fork: st.Fork, Pass: true}
			test.Run(st, cfg, false, rawdb.HashScheme, func(err error, snaps *snapshot.Tree, state *state.StateDB) {
				if state != nil {
					root := state.IntermediateRoot(false)
					result.Root = &root
					if jsonOut {
						fmt.Fprintf(os.Stderr, "{\"stateRoot\": \"%#x\"}\n", root)
					}
				}
				// Dump any state to aid debugging
				if dump {
					dump := state.RawDump(nil)
					result.State = &dump
				}
				if err != nil {
					// Test failed, mark as so
					result.Pass, result.Error = false, err.Error()
				}
			})
			results = append(results, *result)
		}
	}
	out, _ := json.MarshalIndent(results, "", "  ")
	fmt.Println(string(out))
	return nil
}
