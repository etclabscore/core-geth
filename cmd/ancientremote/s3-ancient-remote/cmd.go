package main

import (
	"io"
	"os"

	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/log"
	"github.com/mattn/go-colorable"
	"github.com/mattn/go-isatty"
	"gopkg.in/urfave/cli.v1"
)

var (
	// LogLevelFlag sets log level for server
	LogLevelFlag = cli.IntFlag{
		Name:  "loglevel",
		Value: 3,
		Usage: "log level to emit to the screen",
	}
	// BucketNameFlag sets namespace for S3 bucket
	BucketNameFlag = cli.StringFlag{
		Name:  "bucket",
		Usage: "S3 bucket name (required)",
	}
	// RPCPortFlag sets port for http operation
	RPCPortFlag = cli.IntFlag{
		Name:  "rpcport",
		Usage: "HTTP-RPC server listening port",
		Value: 9797,
	}
	// IPCPathFlag sets ipc path for an ipc rpc server
	IPCPathFlag = utils.DirectoryFlag{
		Name:  "ipcpath",
		Usage: "Filename for IPC socket/pipe within the datadir (explicit paths escape it)",
	}
	// HTTPListenAddrFlag sets address http address to listen on
	HTTPListenAddrFlag = cli.StringFlag{
		Name:  "http.addr",
		Usage: "HTTP-RPC server listening interface",
		Value: "localhost",
	}
)

func mustBucketName(c *cli.Context) (bucketName string) {
	bucketName = c.GlobalString(BucketNameFlag.Name)
	if bucketName == "" {
		utils.Fatalf("Missing required option --%s", BucketNameFlag.Name)
	}
	return bucketName
}

func setupLogFormat(c *cli.Context) error {
	// Set up the logger to print everything
	logOutput := os.Stdout
	usecolor := (isatty.IsTerminal(os.Stderr.Fd()) || isatty.IsCygwinTerminal(os.Stderr.Fd())) && os.Getenv("TERM") != "dumb"
	output := io.Writer(logOutput)
	if usecolor {
		output = colorable.NewColorable(logOutput)
	}
	log.Root().SetHandler(log.LvlFilterHandler(log.Lvl(c.Int(LogLevelFlag.Name)), log.StreamHandler(output, log.TerminalFormat(usecolor))))

	return nil
}
