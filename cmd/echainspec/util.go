package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/params/types/ctypes"
	"github.com/ethereum/go-ethereum/params/types/genesisT"
	"gopkg.in/urfave/cli.v1"
)

func readInputData(ctx *cli.Context) ([]byte, error) {
	if !ctx.GlobalIsSet(fileInFlag.Name) {
		return io.ReadAll(os.Stdin)
	}
	return os.ReadFile(ctx.GlobalString(fileInFlag.Name))
}

func unmarshalChainSpec(format string, data []byte) (conf ctypes.Configurator, err error) {
	conf, ok := chainspecFormatTypes[format]
	if !ok {
		return nil, errInvalidChainspecValue
	}
	err = json.Unmarshal(data, conf)
	if err != nil {
		return conf, err
	}
	if !strings.Contains(format, "geth") {
		return
	}
	// Logic in params/types/gen_genesis.go already "auto-magically"
	// handles genesis Config unmarshaling, and IT PREFERS COREGETH,
	// and the two data types are not mutually exclusive (are overlapping).
	// So we need to redo custom unmarshaling logic to enforce data type
	// preference based on passed format value.
	type dec struct {
		Config ctypes.ChainConfigurator `json:"config"`
	}
	var d dec
	configurator, ok := chainspecFormatTypes[format]
	if !ok {
		return nil, fmt.Errorf("unknown chainspec format type: %v", format)
	}
	switch t := configurator.(type) {
	case *genesisT.Genesis:
		d.Config = t.Config
	default:
		return nil, fmt.Errorf("unhandled chainspec type: %v %v", format, t)
	}
	t := chainspecFormatTypes[format].(*genesisT.Genesis)
	err = json.Unmarshal(data, &d)
	if err != nil {
		return conf, err
	}
	t.Config = d.Config
	conf = t
	return
}

func jsonMarshalPretty(i interface{}) ([]byte, error) {
	return json.MarshalIndent(i, "", "    ")
}
