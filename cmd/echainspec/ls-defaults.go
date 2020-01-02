package main

import (
	"fmt"

	"gopkg.in/urfave/cli.v1"
)

var lsDefaultsCommand = cli.Command{
	Name:   "ls-defaults",
	Usage:  "List default configurations",
	Action: lsDefaults,
}

func lsDefaults(ctx *cli.Context) error {
	for _, name := range defaultChainspecNames {
		fmt.Println(name)
	}
	return nil
}
