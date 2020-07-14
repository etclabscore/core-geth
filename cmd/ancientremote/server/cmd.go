package server

import (
	"fmt"
	"io"
	"os"
	"strings"

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
		Value: 4,
		Usage: "log level to emit to the screen",
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
	// HTTPEnabledFlag sets http enabled for rpc server
	HTTPEnabledFlag = cli.BoolFlag{
		Name:  "http",
		Usage: "Enable the HTTP-RPC server",
	}
	// HTTPListenAddrFlag sets address http address to listen on
	HTTPListenAddrFlag = cli.StringFlag{
		Name:  "http.addr",
		Usage: "HTTP-RPC server listening interface",
		Value: "localhost",
	}
	// HTTPCORSDomainFlag sets corsdomain to accept if needed
	HTTPCORSDomainFlag = cli.StringFlag{
		Name:  "http.corsdomain",
		Usage: "Comma separated list of domains from which to accept cross origin requests (browser enforced)",
		Value: "",
	}
	// HTTPVirtualHostsFlag sets virtual hosts to accept if needed
	HTTPVirtualHostsFlag = cli.StringFlag{
		Name:  "http.vhosts",
		Usage: "Comma separated list of virtual hostnames from which to accept requests (server enforced). Accepts '*' wildcard.",
		Value: strings.Join([]string{"localhost"}, ","),
	}
)

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

func getHTTPEndpoint(c *cli.Context) string {
	return fmt.Sprintf("%s:%d", c.GlobalString(utils.HTTPListenAddrFlag.Name), c.Int(RPCPortFlag.Name))
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
