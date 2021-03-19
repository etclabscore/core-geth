package main

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/ethereum/go-ethereum/params"
)

// TestConsoleCmdNetworkIdentities tests network identity variables at runtime for a geth instance.
// This provides a "production equivalent" integration test for consensus-relevant chain identity values which
// cannot be adequately unit tested because of reliance on cli context variables.
// These tests should cover expected default values and possible flag-interacting values, like --<chain> with --networkid=n.
func TestConsoleCmdNetworkIdentities(t *testing.T) {
	chainIdentityCases := []struct {
		flags       []string
		networkId   int
		chainId     int
		genesisHash string
	}{
		// Default chain value, without and with --networkid flag set.
		{[]string{}, 1, 1, params.MainnetGenesisHash.Hex()},
		{[]string{"--networkid", "42"}, 42, 1, params.MainnetGenesisHash.Hex()},

		// Non-default chain value, without and with --networkid flag set.
		{[]string{"--classic"}, 1, 61, params.MainnetGenesisHash.Hex()},
		{[]string{"--classic", "--networkid", "42"}, 42, 61, params.MainnetGenesisHash.Hex()},

		// All other possible --<chain> values.
		{[]string{"--testnet"}, 3, 3, params.RopstenGenesisHash.Hex()},
		{[]string{"--ropsten"}, 3, 3, params.RopstenGenesisHash.Hex()},
		{[]string{"--rinkeby"}, 4, 4, params.RinkebyGenesisHash.Hex()},
		{[]string{"--goerli"}, 5, 5, params.GoerliGenesisHash.Hex()},
		{[]string{"--kotti"}, 6, 6, params.KottiGenesisHash.Hex()},
		{[]string{"--mordor"}, 7, 63, params.MordorGenesisHash.Hex()},
		{[]string{"--social"}, 28, 28, params.SocialGenesisHash.Hex()},
		{[]string{"--ethersocial"}, 1, 31102, params.EthersocialGenesisHash.Hex()},
		{[]string{"--yolov2"}, 133519467574834, 133519467574834, params.YoloV2GenesisHash.Hex()},
	}
	for i, p := range chainIdentityCases {

		// Disable networking, preventing false-negatives if in an environment without networking service
		// or collisions with an existing geth service.
		p.flags = append(p.flags, "--port", "0", "--maxpeers", "0", "--nodiscover", "--nat", "none")

		t.Run(fmt.Sprintf("%d/%v/networkid", i, p.flags),
			consoleCmdStdoutTest(p.flags, "admin.nodeInfo.protocols.eth.network", p.networkId))
		t.Run(fmt.Sprintf("%d/%v/chainid", i, p.flags),
			consoleCmdStdoutTest(p.flags, "admin.nodeInfo.protocols.eth.config.chainId", p.chainId))
		t.Run(fmt.Sprintf("%d/%v/genesis_hash", i, p.flags),
			consoleCmdStdoutTest(p.flags, "eth.getBlock(0, false).hash", strconv.Quote(p.genesisHash)))
	}
}

func consoleCmdStdoutTest(flags []string, execCmd string, want interface{}) func(t *testing.T) {
	return func(t *testing.T) {
		flags = append(flags, "--exec", execCmd, "console")
		geth := runGeth(t, flags...)
		geth.Expect(fmt.Sprintf(`%v
`, want))
		geth.ExpectExit()
	}
}
