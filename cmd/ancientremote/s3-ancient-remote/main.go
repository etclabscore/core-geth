package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"

	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/metrics"
	"github.com/ethereum/go-ethereum/rpc"
	"gopkg.in/urfave/cli.v1"
)

var (
	app = cli.NewApp()
)

func init() {
	app.Name = "S3AncientRemote"
	app.Usage = "S3 Ancient Remote Storage as a service"
	app.Flags = []cli.Flag{
		LogLevelFlag,
		BucketNameFlag,
		IPCPathFlag,
		RPCPortFlag,
		HTTPListenAddrFlag,
	}
	app.Action = run
}

func createS3FreezerService(bucketName string) (*freezerRemoteS3, chan struct{}) {
	var (
		service    *freezerRemoteS3
		err        error
		readMeter  = metrics.NewRegisteredMeter("ancient.remote /read", nil)
		writeMeter = metrics.NewRegisteredMeter("ancient.remote /write", nil)
		sizeGauge  = metrics.NewRegisteredGauge("ancient.remote /size", nil)
	)

	service, err = newFreezerRemoteS3(bucketName, readMeter, writeMeter, sizeGauge)
	if err != nil {
		utils.Fatalf("Could not initialize S3 service: %w", err)
	}
	return service, service.quit
}

func run(c *cli.Context) error {

	if err := setupLogFormat(c); err != nil {
		return err
	}
	bucketName := mustBucketName(c)
	utils.CheckExclusive(c, IPCPathFlag, HTTPListenAddrFlag.Name)

	log.Info("Creating freezer service", "bucket", bucketName)
	api, quit := createS3FreezerService(bucketName)

	var (
		rpcServer *rpc.Server
		listener  net.Listener
		err       error
	)
	rpcAPIs := []rpc.API{
		{
			Namespace: "freezer",
			Public:    true,
			Service:   api,
			Version:   "1.0",
		},
	}

	if c.GlobalIsSet(IPCPathFlag.Name) {
		endpoint := c.GlobalString(IPCPathFlag.Name)
		log.Info("Using IPC", "endpoint", endpoint)
		listener, rpcServer, err = rpc.StartIPCEndpoint(endpoint, rpcAPIs)
	} else {
		rpcServer = rpc.NewServer()
		err = rpcServer.RegisterName("freezer", api)
		if err != nil {
			return err
		}
		endpoint := fmt.Sprintf("%s:%d", c.GlobalString(utils.HTTPListenAddrFlag.Name), c.Int(RPCPortFlag.Name))
		log.Info("Using TCP", "endpoint", endpoint)
		listener, err = net.Listen("tcp", endpoint)
		if err != nil {
			return err
		}
	}

	go func() {
		log.Info("Serving listener", "listener", listener.Addr().String())
		if err := rpcServer.ServeListener(listener); err != nil {
			log.Crit("Exiting", "error", err)
		}
	}()

	abortChan := make(chan os.Signal, 1)
	signal.Notify(abortChan, os.Interrupt)

	defer func() {
		// Don't bother imposing a timeout here.
		select {
		case sig := <-abortChan:
			log.Info("Got abort...", "signal", sig)
			rpcServer.Stop()
		case <-quit:
			log.Info("S3 connection closing")
			rpcServer.Stop()
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
