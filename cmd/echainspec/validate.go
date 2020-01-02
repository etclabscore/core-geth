package main

import (
	"log"
	"os"

	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/params/confp"
	"gopkg.in/urfave/cli.v1"
)

var validateCommand = cli.Command{
	Name:        "validate",
	Aliases:     []string{"valid"},
	Description: "Exits 0 if valid, 1 if not.",
	Usage:       "Tests whether a configuration is valid",
	ArgsUsage:   "[|0x042|0x42|42]",
	Action:      validate,
}

func validate(ctx *cli.Context) error {
	var h *uint64
	if ctx.Args().Present() {
		var head math.HexOrDecimal64
		err := head.UnmarshalText([]byte(ctx.Args().First()))
		if err != nil {
			return err
		}
		var hh = uint64(head)
		h = &hh
	}
	err := confp.IsValid(globalChainspecValue, h)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	log.Println("Valid")
	os.Exit(0)
	return nil
}
