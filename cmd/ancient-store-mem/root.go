// Copyright 2020 The core-geth Authors
// This file is part of the core-geth library.
//
// The core-geth library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The core-geth library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the core-geth library. If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/cmd/ancient-store-mem/lib"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ancient-store-mem",
	Short: "Memory-backed remote ancient store application",
	Long: `Uses a memory-backed map to store ancient data.

This application is intended for testing purposed only.
Ancient data is stored ephemerally.

Expects first and only argument to an IPC path, or, the directory
in which a default 'mock-freezer.ipc' path should be created.

Package 'lib' logic may be imported and used in testing contexts as well.
`,

	Run: func(cmd *cobra.Command, args []string) {
		ipcPath := args[0]
		fi, err := os.Stat(ipcPath)
		if err != nil && !os.IsNotExist(err) {
			log.Fatalln(err)
		}
		if fi != nil && fi.IsDir() {
			ipcPath = filepath.Join(ipcPath, "mock-freezer.ipc")
		}
		listener, server, err := rpc.StartIPCEndpoint(ipcPath, nil)
		if err != nil {
			log.Fatalln(err)
		}
		defer os.Remove(ipcPath)
		mock := lib.NewMemFreezerRemoteServerAPI()
		err = server.RegisterName("freezer", mock)
		if err != nil {
			log.Fatalln(err)
		}
		quit := make(chan bool, 1)
		go func() {
			log.Println("Serving", listener.Addr())
			log.Fatalln(server.ServeListener(listener))
		}()
		<-quit
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
