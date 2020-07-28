package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"context"

	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/metrics"
	"github.com/ethereum/go-ethereum/rpc"
	"gopkg.in/urfave/cli.v1"
)

var (

	storjAPIKey string
	storjSatellite string
	storjSecret string
	app = cli.NewApp()
)

func init() {
	app.Name = "StorjAncientRemote"
	app.Usage = "Storj Ancient Remote Storage as a service"
	app.Flags = []cli.Flag{
		BucketNameFlag,
		RPCPortFlag,
		HTTPListenAddrFlag,
		IPCPathFlag,
	}
	app.Action = remoteAncientStore
	storjAPIKey = os.Getenv("STORJ_API_KEY")
	storjSatellite = os.Getenv("STORJ_SATELLITE")
	storjSecret = os.Getenv("STORJ_SECRET")

	if(storjAPIKey == "" ){
		utils.Fatalf("Missing environment variable for STORJ_API_KEY")
	}
	if( storjSecret == ""){
		utils.Fatalf("Missing environment variable for STORJ_SECRET")
	}
	if( storjSatellite == ""){
		utils.Fatalf("Missing one environment variable for STORJ_SATELLITE")
	}

}

func createStorjFreezerService(ctx context.Context, bucketName string) (*freezerRemoteStorj, chan struct{}) {
	var (
		service    *freezerRemoteStorj
		err        error
		mets = logMetrics{
			readMeter: metrics.NewRegisteredMeter("ancient.remote /read", nil), 
			writeMeter: metrics.NewRegisteredMeter("ancient.remote /write", nil), 
			sizeGauge: metrics.NewRegisteredGauge("ancient.remote /size", nil),
		}
		access = storjAccess{
			apiKey: storjAPIKey,
			passphrase: storjSecret,
			satellite: storjSatellite,
		}
	)

	service, err = newFreezerRemoteStorj(ctx, bucketName, access, mets)
	if err != nil {
		utils.Fatalf("Could not initialize Storj service: %w", err)
	}
	return service, service.quit
}

func remoteAncientStore(c *cli.Context) error {

	setupLogFormat(c)
	bucketName := checkBucketArg(c)
	utils.CheckExclusive(c, IPCPathFlag, HTTPListenAddrFlag.Name)

	api, quit := createStorjFreezerService(context.Background(), bucketName)

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
		listener, rpcServer, err = rpc.StartIPCEndpoint(c.GlobalString(IPCPathFlag.Name), rpcAPIs)
	} else {
		rpcServer = rpc.NewServer()
		err = rpcServer.RegisterName("freezer", api)
		if err != nil {
			return err
		}
		endpoint := fmt.Sprintf("%s:%d", c.GlobalString(utils.HTTPListenAddrFlag.Name), c.Int(RPCPortFlag.Name))
		listener, err = net.Listen("tcp", endpoint)
		fmt.Println("listening on",endpoint)
		if err != nil {
			return err
		}
	}

	go func() {
		if err := rpcServer.ServeListener(listener); err != nil {
			log.Crit("exiting", "error", err)
		}
	}()

	abortChan := make(chan os.Signal, 1)
	signal.Notify(abortChan, os.Interrupt)

	defer func() {
		// Don't bother imposing a timeout here.
		select {
		case sig := <-abortChan:
			log.Info("Exiting...", "signal", sig)
			rpcServer.Stop()
		case <-quit:
			log.Info("Storj connection closing")
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
