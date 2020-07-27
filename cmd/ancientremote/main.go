package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"

	"github.com/ethereum/go-ethereum/cmd/ancientremote/server"
	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/metrics"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/mattn/go-colorable"
	"github.com/mattn/go-isatty"
	"gopkg.in/urfave/cli.v1"
)

var (
	// BucketNameFlag sets namespace for S3 bucket
	BucketNameFlag = cli.StringFlag{
		Name:  "bucket",
		Usage: "S3 bucket name",
	}
	app = cli.NewApp()
)

func init() {
	app.Name = "AncientRemote"
	app.Usage = "Ancient Remote Storage as a service"
	app.Flags = []cli.Flag{
		BucketNameFlag,
		server.RPCPortFlag,
		server.HTTPListenAddrFlag,
	}
	app.Action = remoteAncientStore
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

func checkNamespaceArg(c *cli.Context) (bucketName string) {
	bucketName = c.GlobalString(BucketNameFlag.Name)
	if bucketName == "" {
		utils.Fatalf("Missing namespace please specify a namespace, with --namespace")
	}
	return
}

func setupLogFormat(c *cli.Context) error {
	// Set up the logger to print everything
	logOutput := os.Stdout
	usecolor := (isatty.IsTerminal(os.Stderr.Fd()) || isatty.IsCygwinTerminal(os.Stderr.Fd())) && os.Getenv("TERM") != "dumb"
	output := io.Writer(logOutput)
	if usecolor {
		output = colorable.NewColorable(logOutput)
	}
	log.Root().SetHandler(log.LvlFilterHandler(log.Lvl(c.Int(server.LogLevelFlag.Name)), log.StreamHandler(output, log.TerminalFormat(usecolor))))

	return nil
}

func remoteAncientStore(c *cli.Context) error {

	setupLogFormat(c)
	namespace := checkNamespaceArg(c)
	utils.CheckExclusive(c, server.IPCPathFlag, server.HTTPListenAddrFlag.Name)

	api, quit := createS3FreezerService(namespace)

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

	if c.GlobalIsSet(server.IPCPathFlag.Name) {
		listener, rpcServer, err = rpc.StartIPCEndpoint(c.GlobalString(server.IPCPathFlag.Name), rpcAPIs)
	} else {
		rpcServer = rpc.NewServer()
		err = rpcServer.RegisterName("freezer", api)
		if err != nil {
			return err
		}
		endpoint := fmt.Sprintf("%s:%d", c.GlobalString(utils.HTTPListenAddrFlag.Name), c.Int(server.RPCPortFlag.Name))
		listener, err = net.Listen("tcp", endpoint)
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
