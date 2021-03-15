package ethclient

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

/*
The '_CanCompareGoEthereum' denotes tests that can be run against ethereum/go-ethereum.
See .github/workflows/geth_1to1.yml for exemplary testing.
*/

// TestEthGetBlockByNumber_ValidJSONResponse_CanCompareGoEthereum tests that
// JSON RPC API responses to eth_getBlockByNumber meet pattern-based expectations.
// These validations include the null-ness of certain fields for the 'pending' block
// as well existence of all expected keys and values.
func TestEthGetBlockByNumber_ValidJSONResponse_CanCompareGoEthereum(t *testing.T) {
	backend, _ := newTestBackend(t)
	client, _ := backend.Attach()
	defer backend.Close()
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	// Have to sleep a little to make sure miner has time to set pending.
	time.Sleep(time.Millisecond * 100)

	// Get a reference block.
	parent, err := NewClient(client).HeaderByNumber(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	if parent == nil {
		t.Fatal("bad test")
	}

	reNull := regexp.MustCompile(`^null$`)
	reHexAnyLen := regexp.MustCompile(`^"0x[a-zA-Z0-9]+"$`)
	reHexHashLen := regexp.MustCompile(fmt.Sprintf(`^"0x[a-zA-Z0-9]{%d}"$`, common.HashLength*2))

	for blockHeight, cases := range map[string]map[string]*regexp.Regexp{
		"earliest": {
			"nonce": reHexAnyLen,
			"hash":  reHexAnyLen,
			"miner": regexp.MustCompile(fmt.Sprintf(`^"0x[a-zA-Z0-9]{%d}"$`, common.AddressLength*2)),

			"totalDifficulty": reHexAnyLen,

			"mixHash":   regexp.MustCompile(fmt.Sprintf(`^"0x[0]{%d}"$`, common.HashLength*2)),
			"logsBloom": regexp.MustCompile(fmt.Sprintf(`^"0x[0]{%d}"$`, types.BloomByteLength*2)),

			"number":     reHexAnyLen,
			"difficulty": reHexAnyLen,
			"gasLimit":   reHexAnyLen,
			"gasUsed":    reHexAnyLen,
			"size":       reHexAnyLen,
			"timestamp":  reHexAnyLen,
			"extraData":  reHexAnyLen,

			"parentHash":       reHexHashLen,
			"transactionsRoot": reHexHashLen,
			"stateRoot":        reHexHashLen,
			"receiptsRoot":     reHexHashLen,
			"sha3Uncles":       reHexHashLen,

			"uncles":       regexp.MustCompile(`^\[\]$`),
			"transactions": regexp.MustCompile(`^\[\]$`),
		},
		"latest": {
			"nonce": reHexAnyLen,
			"hash":  reHexAnyLen,
			"miner": regexp.MustCompile(fmt.Sprintf(`^"0x[a-zA-Z0-9]{%d}"$`, common.AddressLength*2)),

			"totalDifficulty": reHexAnyLen,

			"mixHash":   regexp.MustCompile(fmt.Sprintf(`^"0x[0]{%d}"$`, common.HashLength*2)),
			"logsBloom": regexp.MustCompile(fmt.Sprintf(`^"0x[0]{%d}"$`, types.BloomByteLength*2)),

			"number":     reHexAnyLen,
			"difficulty": reHexAnyLen,
			"gasLimit":   reHexAnyLen,
			"gasUsed":    reHexAnyLen,
			"size":       reHexAnyLen,
			"timestamp":  reHexAnyLen,
			"extraData":  reHexAnyLen,

			"parentHash":       reHexHashLen,
			"transactionsRoot": reHexHashLen,
			"stateRoot":        reHexHashLen,
			"receiptsRoot":     reHexHashLen,
			"sha3Uncles":       reHexHashLen,

			"uncles":       regexp.MustCompile(`^\[\]$`),
			"transactions": regexp.MustCompile(`^\[\]$`),
		},
		"pending": {
			"nonce": reNull,
			"hash":  reNull,
			"miner": reNull,

			"totalDifficulty": reNull,

			"mixHash":   regexp.MustCompile(fmt.Sprintf(`^"0x[0]{%d}"$`, common.HashLength*2)),
			"logsBloom": regexp.MustCompile(fmt.Sprintf(`^"0x[0]{%d}"$`, types.BloomByteLength*2)),

			"number":     reHexAnyLen,
			"difficulty": reHexAnyLen,
			"gasLimit":   reHexAnyLen,
			"gasUsed":    reHexAnyLen,
			"size":       reHexAnyLen,
			"timestamp":  reHexAnyLen,
			"extraData":  reHexAnyLen,

			"parentHash":       reHexHashLen,
			"transactionsRoot": reHexHashLen,
			"stateRoot":        reHexHashLen,
			"receiptsRoot":     reHexHashLen,
			"sha3Uncles":       reHexHashLen,

			"uncles":       regexp.MustCompile(`^\[\]$`),
			"transactions": regexp.MustCompile(`^\[\]$`),
		},
	} {
		for _, fullTxes := range []bool{true, false} {
			gotPending := make(map[string]json.RawMessage)
			err := client.CallContext(ctx, &gotPending, "eth_getBlockByNumber", blockHeight, fullTxes)
			if err != nil {
				t.Fatal(err)
			}

			for key, re := range cases {
				gotVal, ok := gotPending[key]
				if !ok {
					t.Errorf("%s: missing key", key)
				}
				if !re.Match(gotVal) {
					t.Errorf("%s want: %v, got: %v", key, re, string(gotVal))
				}
			}

			for k, v := range gotPending {
				if _, ok := cases[k]; !ok {
					t.Errorf("%s: missing key (value: %v)", k, string(v))
				}
			}
		}
	}
}
