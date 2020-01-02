package main

import (
	"fmt"

	"gopkg.in/urfave/cli.v1"
)

var lsFormatsCommand = cli.Command{
	Name:   "ls-formats",
	Usage:  "List client configuration formats",
	Action: lsFormats,
}

func lsFormats(ctx *cli.Context) error {
	for _, name := range chainspecFormats {
		fmt.Println(name)
	}
	return nil
}
