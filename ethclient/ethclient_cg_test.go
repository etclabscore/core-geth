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

// TestEthGetBlockByNumber_ValidJSONResponse tests that
// JSON RPC API responses to eth_getBlockByNumber meet pattern-based expectations.
// These validations include the null-ness of certain fields for the 'pending' block
// as well existence of all expected keys and values.
func TestEthGetBlockByNumber_ValidJSONResponse(t *testing.T) {
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

	// completeBlockExpectations define expectations for 'earliest' and 'latest' blocks.
	completeBlockExpectations := map[string]*regexp.Regexp{
		"nonce": reHexAnyLen,
		"hash":  reHexHashLen,
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
	}

	// Construct the 'pending' block expectations as a copy of the concrete block
	// expectations.
	pendingBlockExpectations := map[string]*regexp.Regexp{}
	for k, v := range completeBlockExpectations {
		pendingBlockExpectations[k] = v
	}

	// Make 'pending' specific adjustments.
	pendingBlockExpectations["nonce"] = reNull
	pendingBlockExpectations["hash"] = reNull
	pendingBlockExpectations["miner"] = reNull
	pendingBlockExpectations["totalDifficulty"] = reNull

	for blockHeight, cases := range map[string]map[string]*regexp.Regexp{
		"earliest": completeBlockExpectations,
		"latest":   completeBlockExpectations,
		"pending":  pendingBlockExpectations,
	} {
		for _, fullTxes := range []bool{true, false} {
			t.Run(fmt.Sprintf("eth_getBlockByNumber-%s-%v", blockHeight, fullTxes), func(t *testing.T) {
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
			})
		}
	}
}
