package main

import (
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/node"
	"gopkg.in/urfave/cli.v1"
)

var (
	logLevelFlag = cli.IntFlag{
		Name:  "loglevel",
		Value: 4,
		Usage: "log level to emit to the screen",
	}
	rpcPortFlag = cli.IntFlag{
		Name:  "rpcport",
		Usage: "HTTP-RPC server listening port",
		Value: node.DefaultHTTPPort + 5,
	}

	app = cli.NewApp()
)

func init() {
	app.Name = "AncientRemote"
	app.Usage = "Ancient Remote Storage as a service"
	app.Flags = []cli.Flag{
		rpcPortFlag,
		logLevelFlag,
		utils.AncientRemoteFlag,
		utils.AncientRemoteNamespaceFlag,
	}
	app.Action = remoteAncientStore
}

func remoteAncientStore(c *cli.Context) error {
	if args := c.Args(); len(args) > 0 {
		return fmt.Errorf("invalid command: %q", args[0])
	}
	var (
	//	api *rawdb.FreezerRemoteAPI //rawdb.ExternalFreezerRemoteAPI
	)
	namespace := c.GlobalString(utils.AncientRemoteNamespaceFlag.Name)
	clientOption := c.GlobalString(utils.AncientRemoteFlag.Name)

	_, err := rawdb.NewFreezerRemoteAPI(clientOption, namespace)

	return err
}

func main() {
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println("well we didn't blow up")
}
