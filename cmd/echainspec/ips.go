package main

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/params/confp"
	"gopkg.in/urfave/cli.v1"
)

var ipsCommand = cli.Command{
	Name:   "ips",
	Usage:  "List IP transition names and values",
	Action: ips,
}

func ips(ctx *cli.Context) error {
	fns, names := confp.Transitions(globalChainspecValue)
	for i, fn := range fns {
		name := strings.TrimPrefix(names[i], "Get")
		name = strings.TrimSuffix(name, "Transition")

		var printv interface{}
		v := fn()
		if v != nil {
			printv = *v
		} else {
			printv = "-"
		}

		fmt.Println(name, fmt.Sprintf("%v", printv))
	}
	return nil
}
