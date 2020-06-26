package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"

	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/mattn/go-colorable"
	"github.com/mattn/go-isatty"
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
		utils.HTTPListenAddrFlag,
		utils.HTTPVirtualHostsFlag,
		utils.HTTPEnabledFlag,
		utils.HTTPCORSDomainFlag,
	}
	app.Action = remoteAncientStore
}

// splitAndTrim splits input separated by a comma
// and trims excessive white space from the substrings.
func splitAndTrim(input string) []string {
	result := strings.Split(input, ",")
	for i, r := range result {
		result[i] = strings.TrimSpace(r)
	}
	return result
}

func initialize(c *cli.Context) error {
	// Set up the logger to print everything
	logOutput := os.Stdout
	usecolor := (isatty.IsTerminal(os.Stderr.Fd()) || isatty.IsCygwinTerminal(os.Stderr.Fd())) && os.Getenv("TERM") != "dumb"
	output := io.Writer(logOutput)
	if usecolor {
		output = colorable.NewColorable(logOutput)
	}
	log.Root().SetHandler(log.LvlFilterHandler(log.Lvl(c.Int(logLevelFlag.Name)), log.StreamHandler(output, log.TerminalFormat(usecolor))))

	return nil
}

func remoteAncientStore(c *cli.Context) error {
	if args := c.Args(); len(args) > 0 {
		return fmt.Errorf("invalid command: %q", args[0])
	}
	var (
		api *rawdb.FreezerRemoteAPI //rawdb.ExternalFreezerRemoteAPI
	)
	namespace := c.GlobalString(utils.AncientRemoteNamespaceFlag.Name)
	clientOption := c.GlobalString(utils.AncientRemoteFlag.Name)
	if err := initialize(c); err != nil {
		return err
	}
	api, err := rawdb.NewFreezerRemoteAPI(clientOption, namespace)
	if err != nil {
		utils.Fatalf("Could not start freezer: %w", err)
	}
	rpcAPI := []rpc.API{
		{
			Namespace: "freezer",
			Public:    true,
			Service:   api,
			Version:   "1.0"},
	}

	if c.GlobalBool(utils.HTTPEnabledFlag.Name) {
		vhosts := []string{"*"} //splitAndTrim(c.GlobalString(utils.HTTPVirtualHostsFlag.Name))
		cors := []string{"*"}   //splitAndTrim(c.GlobalString(utils.HTTPCORSDomainFlag.Name))

		srv := rpc.NewServer()
		err = node.RegisterApisFromWhitelist(rpcAPI, []string{"freezer"}, srv, false)
		if err != nil {
			utils.Fatalf("Could not register API: %w", err)
		}
		handler := node.NewHTTPHandlerStack(srv, cors, vhosts)

		// start http server
		httpEndpoint := fmt.Sprintf("%s:%d", c.GlobalString(utils.HTTPListenAddrFlag.Name), c.Int(rpcPortFlag.Name))
		httpServer, addr, err := node.StartHTTPEndpoint(httpEndpoint, rpc.DefaultHTTPTimeouts, handler)
		if err != nil {
			utils.Fatalf("Could not start RPC api: %v", err)
		}
		extapiURL := fmt.Sprintf("http://%v/", addr)

		log.Info("HTTP endpoint opened", "url", extapiURL)

		defer func() {
			// Don't bother imposing a timeout here.
			httpServer.Shutdown(context.Background())
			log.Info("HTTP endpoint closed", "url", extapiURL)
		}()
	}

	abortChan := make(chan os.Signal, 1)
	signal.Notify(abortChan, os.Interrupt)

	sig := <-abortChan
	log.Info("Exiting...", "signal", sig)

	return nil
}

func main() {
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println("well we didn't blow up")
}
