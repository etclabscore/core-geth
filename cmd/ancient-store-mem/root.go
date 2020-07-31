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
