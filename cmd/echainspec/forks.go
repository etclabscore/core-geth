package main

import (
	"fmt"

	"github.com/ethereum/go-ethereum/params/confp"
	"gopkg.in/urfave/cli.v1"
)

var forksCommand = cli.Command{
	Name:               "forks",
	Usage:              "List unique and non-zero fork numbers",
	Action: forks,
}

func forks(ctx *cli.Context) error {
	for _, f := range confp.Forks(globalChainspecValue) {
		fmt.Println(f)
	}
	return nil
}
