package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ancient-store-fs",
	Short: "FS-backed remote ancient store application",
	Long: `Uses a fs-backed map to store ancient data.
`,

	Run: func(cmd *cobra.Command, args []string) {

		var (
			datadir string
			listener net.Listener
			server *rpc.Server
			err error
		)

		listenAndServe := func() {
			go func() {
				log.Println("Serve and listen:", listener.Addr())
				log.Fatal(server.ServeListener(listener))
			}()
		}

		registerFreezer := func() {
			freezer, err := rawdb.NewFreezer(datadir, "")
			if err != nil {
				log.Fatalf("new freezer: %v", err)
			}

			err = server.RegisterName("freezer", freezer)
			if err != nil {
				log.Fatalf("server register freezer: %v", err)
			}
		}

		datadir = cmd.Flag("datadir").Value.String()
		if datadir == "" {
			log.Fatalf("must configure datadir")
		}

		if v := cmd.Flag("ws").Value.String(); v != "" {
			listener, err = net.Listen("tcp", v)
			if err != nil {
				log.Fatalf("net listen: %v", err)
			}
			server = rpc.NewServer()
			registerFreezer()
			listenAndServe()
		} else if v := cmd.Flag("ipc").Value.String(); v != "" {
			listener, server, err = rpc.StartIPCEndpoint(v, nil)
			if err != nil {
				log.Fatalf("ipc endpoint: %v", err)
			}
			defer os.Remove(v)
			registerFreezer()
			listenAndServe()
		} else {
			log.Fatal("must configure listener address")
		}

		// Dumb block
		quit := make(chan bool, 1)
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

func init() {
	rootCmd.Flags().StringP("ws", "w", "", "WS listening endpoint")
	rootCmd.Flags().StringP("ipc", "s", "", "Unix socket listening endpoint")
	rootCmd.Flags().StringP("datadir", "d", "", "Filepath to data directory")
}
