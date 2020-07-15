package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/ethereum/go-ethereum/cmd/ancientremote/server"
	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/metrics"
	"github.com/ethereum/go-ethereum/rpc"
	"gopkg.in/urfave/cli.v1"
)

var (
	// NamespaceFlag sets namespace for S3 bucket
	NamespaceFlag = cli.StringFlag{
		Name:  "namespace",
		Usage: "Namespace for remote storage, eg. S3 bucket name. Use will vary by remote provider.",
	}
	app = cli.NewApp()
)

func init() {
	app.Name = "AncientRemote"
	app.Usage = "Ancient Remote Storage as a service"
	app.Flags = []cli.Flag{
		NamespaceFlag,
		server.RPCPortFlag,
		server.LogLevelFlag,
		server.IPCPathFlag,
		server.HTTPListenAddrFlag,
		server.HTTPVirtualHostsFlag,
		server.HTTPEnabledFlag,
		server.HTTPCORSDomainFlag,
	}
	app.Action = remoteAncientStore
}

func createS3FreezerService(namespace string) (*rawdb.FreezerRemoteAPI, chan struct{}) {
	var (
		api        *rawdb.FreezerRemoteAPI
		service    *freezerRemoteS3
		err        error
		readMeter  = metrics.NewRegisteredMeter("ancient.remote /read", nil)
		writeMeter = metrics.NewRegisteredMeter("ancient.remote /write", nil)
		sizeGauge  = metrics.NewRegisteredGauge("ancient.remote /size", nil)
	)

	service, err = newFreezerRemoteS3(namespace, readMeter, writeMeter, sizeGauge)
	if err != nil {
		utils.Fatalf("Could not initialize S3 service: %w", err)
	}
	api, err = rawdb.NewFreezerRemoteAPI(service)
	if err != nil {
		utils.Fatalf("Could not start freezer: %w", err)
	}
	return api, service.quit
}

func checkNamespaceArg(c *cli.Context) (namespace string) {
	namespace = c.GlobalString(NamespaceFlag.Name)
	if namespace == "" {
		utils.Fatalf("Missing namespace please specify a namespace, with --namespace")
	}
	return
}

func remoteAncientStore(c *cli.Context) error {

	cfg := server.MakeServerConfig(c)
	namespace := checkNamespaceArg(c)
	api, quit := createS3FreezerService(namespace)
	rpcAPI := []rpc.API{
		{
			Namespace: "freezer",
			Public:    true,
			Service:   api,
			Version:   "1.0"},
	}

	srv := server.NewServer(cfg, rpcAPI, []string{"freezer"})
	srv.Start()

	abortChan := make(chan os.Signal, 1)
	signal.Notify(abortChan, os.Interrupt)

	defer func() {
		// Don't bother imposing a timeout here.
		select {
		case sig := <-abortChan:
			log.Info("Exiting...", "signal", sig)
			srv.Stop()
		case <-quit:
			log.Info("S3 connection closing")
			srv.Stop()
		}
	}()
	return nil
}

func main() {
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
