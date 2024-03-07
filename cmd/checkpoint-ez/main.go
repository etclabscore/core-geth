package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/internal/flags"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params/types/ctypes"
	"github.com/ethereum/go-ethereum/params/vars"
	"github.com/urfave/cli/v2"
)

var app = flags.NewApp("ethereum checkpoint helper tool")

var commandCreate = &cli.Command{
	Name:  "create",
	Usage: "Prints the latest eligible checkpoint value, JSON encoded to stdout, based on the current chain state of the remote node",
	Flags: []cli.Flag{
		nodeURLFlag,
	},
	Action: create,
}

func init() {
	app.Commands = []*cli.Command{
		commandCreate,
	}
	app.Flags = []cli.Flag{
		nodeURLFlag,
	}
}

// Command line flags.
var (
	nodeURLFlag = &cli.StringFlag{
		Name:  "rpc",
		Value: "http://localhost:8545",
		Usage: "The rpc endpoint of a local or remote geth node",
	}
)

func create(ctx *cli.Context) error {
	client, err := ethclient.Dial(ctx.String(nodeURLFlag.Name))
	if err != nil {
		return err
	}
	defer client.Close()

	// Check if the node is still syncing.
	// Only fully-synced remote nodes should be used to derive a checkpoint.
	progress, err := client.SyncProgress(context.Background())
	if err != nil {
		return err
	}
	if progress != nil {
		return fmt.Errorf("node is still syncing")
	}

	var result [4]string
	err = client.Client().Call(&result, "les_latestCheckpoint")
	if err == nil {
		log.Info("got remote checkpoint")
		index, err := strconv.ParseUint(result[0], 0, 64)
		if err != nil {
			return fmt.Errorf("failed to parse checkpoint index %v", err)
		}
		checkpoint := &ctypes.TrustedCheckpoint{
			SectionIndex: index,
			SectionHead:  common.HexToHash(result[1]),
			CHTRoot:      common.HexToHash(result[2]),
			BloomRoot:    common.HexToHash(result[3]),
		}
		b, err := json.MarshalIndent(checkpoint, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(b))
		return nil
	}

	log.Warn("failed to get remote checkpoint (is the les API exposed?)", "error", err)
	log.Warn("falling back to manual checkpoint creation")

	blockNumber, err := client.BlockNumber(context.Background())
	if err != nil {
		return err
	}
	log.Info("Got latest block number", "blockNumber", blockNumber)
	if blockNumber <= vars.FullImmutabilityThreshold {
		return fmt.Errorf("chain is not yet mature enough for checkpoint creation")
	}

	immutableBlockNumber := blockNumber - vars.FullImmutabilityThreshold
	sectionIndex := immutableBlockNumber / vars.CheckpointFrequency // uint64s => floor division

	// The checkpoint is the last block of the last section.
	checkpointNumber := sectionIndex*vars.CheckpointFrequency - 1

	header, err := client.HeaderByNumber(context.Background(), big.NewInt(int64(checkpointNumber)))
	if err != nil {
		return err
	}

	checkpoint := &ctypes.TrustedCheckpoint{
		SectionIndex: sectionIndex - 1, // zero-indexed
		SectionHead:  header.Hash(),
		CHTRoot:      header.Root,

		// BloomRoot can't be derived remotely, and is unnecessary.
		BloomRoot: common.Hash{},
	}
	b, err := json.MarshalIndent(checkpoint, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	return nil
}

func main() {
	log.Root().SetHandler(log.LvlFilterHandler(log.LvlInfo, log.StreamHandler(os.Stderr, log.TerminalFormat(true))))

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
