package cmd

import (
  "fmt"
  "log"
  "os"
  "path/filepath"

  "github.com/ethereum/go-ethereum/cmd/ancientremote/mock-ancient-remote/lib"
  "github.com/ethereum/go-ethereum/rpc"
  "github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
  Use:   "mock-ancient-remote",
  Short: "A brief description of your application",
  Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
  // Uncomment the following line if your bare application
  // has an action associated with it:
  	Run: func(cmd *cobra.Command, args []string) {
  	  ipcPath := args[0]
  	  fi, err := os.Stat(ipcPath)
  	  if fi.IsDir() {
  	    ipcPath = filepath.Join(ipcPath, "mock-freezer.ipc")
      }
  	  listener, server, err := rpc.StartIPCEndpoint(ipcPath, nil)
      if err != nil {
        log.Fatalln(err)
      }
      defer os.Remove(ipcPath)
      mock := lib.NewMockFreezerRemoteServerAPI()
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

func init() {

  // Cobra also supports local flags, which will only run
  // when this action is called directly.
  rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

