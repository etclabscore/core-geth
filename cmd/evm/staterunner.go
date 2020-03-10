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
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/tests"

	cli "gopkg.in/urfave/cli.v1"
)

var stateTestCommand = cli.Command{
	Action:    stateTestCmd,
	Name:      "statetest",
	Usage:     "executes the given state tests",
	ArgsUsage: "<file>",
	Flags:     []cli.Flag{stateTestEVMCEWASMFlag},
}

var stateTestEVMCEWASMFlag = cli.StringFlag{
	Name:  "evmc.ewasm",
	Usage: "EVMC EWASM configuration",
}

// StatetestResult contains the execution status after running a state test, any
// error that might have occurred and a dump of the final state if requested.
type StatetestResult struct {
	Name  string      `json:"name"`
	Pass  bool        `json:"pass"`
	Fork  string      `json:"fork"`
	Error string      `json:"error,omitempty"`
	State *state.Dump `json:"state,omitempty"`
}

func stateTestCmd(ctx *cli.Context) error {
	if len(ctx.Args().First()) == 0 {
		return errors.New("path-to-test argument required")
	}
	// Configure the go-ethereum logger
	glogger := log.NewGlogHandler(log.StreamHandler(os.Stderr, log.TerminalFormat(false)))
	glogger.Verbosity(log.Lvl(ctx.GlobalInt(VerbosityFlag.Name)))
	log.Root().SetHandler(glogger)

	if s := ctx.String(stateTestEVMCEWASMFlag.Name); s != "" {
		log.Info("Running tests with %s=%s", "evmc.ewasm", s)
		vm.InitEVMCEwasm(s)
	}

	// Configure the EVM logger
	config := &vm.LogConfig{
		DisableMemory: ctx.GlobalBool(DisableMemoryFlag.Name),
		DisableStack:  ctx.GlobalBool(DisableStackFlag.Name),
	}
	var (
		tracer   vm.Tracer
		debugger *vm.StructLogger
	)
	switch {
	case ctx.GlobalBool(MachineFlag.Name):
		tracer = vm.NewJSONLogger(config, os.Stderr)

	case ctx.GlobalBool(DebugFlag.Name):
		debugger = vm.NewStructLogger(config)
		tracer = debugger

	default:
		debugger = vm.NewStructLogger(config)
	}

	// Iterate over all the tests, run them and aggregate the results
	cfg := vm.Config{
		Tracer:           tracer,
		Debug:            ctx.GlobalBool(DebugFlag.Name) || ctx.GlobalBool(MachineFlag.Name),
		EWASMInterpreter: ctx.String(stateTestEVMCEWASMFlag.Name),
	}

	results := []StatetestResult{}
	hadFailure := false

	runFile := func(filename string) error {
		if filepath.Ext(filename) != ".json" {
			log.Warn("Skipping non-json file", "file", filename)
			return nil
		}
		// Load the test content from the input file
		src, err := ioutil.ReadFile(filename)
		if err != nil {
			return err
		}

		var tests map[string]tests.StateTest
		if err = json.Unmarshal(src, &tests); err != nil {
			return fmt.Errorf("err=%v filename=%v", err, filename)
		}

		for key, test := range tests {
			for _, st := range test.Subtests() {
				// Run the test and aggregate the result
				result := &StatetestResult{Name: key, Fork: st.Fork, Pass: true}
				state, err := test.Run(st, cfg)
				// print state root for evmlab tracing
				if ctx.GlobalBool(MachineFlag.Name) && state != nil {
					fmt.Fprintf(os.Stderr, "{\"stateRoot\": \"%x\"}\n", state.IntermediateRoot(false))
				}
				if err != nil {
					// Test failed, mark as so and dump any state to aid debugging
					result.Pass, result.Error = false, err.Error()
					hadFailure = true
					if ctx.GlobalBool(DumpFlag.Name) && state != nil {
						dump := state.RawDump(false, false, true)
						result.State = &dump
					}
				}

				results = append(results, *result)

				// Print any structured logs collected
				if ctx.GlobalBool(DebugFlag.Name) {
					if debugger != nil {
						fmt.Fprintln(os.Stderr, "#### TRACE ####")
						vm.WriteTrace(os.Stderr, debugger.StructLogs())
					}
				}
			}
		}
		return nil
	}
	runDir := func(filename string) error {
		return filepath.Walk(filename, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			log.Info(path)
			if info.IsDir() {
				return nil
			}
			err = runFile(path)
			if err != nil {
				return err
			}
			return nil
		})
	}
	target := ctx.Args().First()
	if info, err := os.Stat(target); err != nil {
		log.Error("err stat", "error", err.Error())
		return err
	} else if info.IsDir() {
		log.Info("rundir", target)
		if err := runDir(target); err != nil {
			log.Error("err dir", "error", err.Error())
			return err
		}
	} else {
		log.Info("runfile", target)
		if err := runFile(target); err != nil {
			log.Error("err file", "error", err.Error())
			return err
		}
	}

	out, _ := json.MarshalIndent(results, "", "  ")
	fmt.Println(string(out))
	if hadFailure {
		os.Exit(1)
	}
	return nil
}
