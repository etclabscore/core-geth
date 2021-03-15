package ethclient

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"reflect"
	"regexp"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/go-test/deep"
)

func TestHeader_TxesUnclesNotEmpty(t *testing.T) {
	backend, _ := newTestBackend(t)
	client, _ := backend.Attach()
	defer backend.Close()
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	res := make(map[string]interface{})
	err := client.CallContext(ctx, &res, "eth_getBlockByNumber", "latest", false)
	if err != nil {
		log.Fatalln(err)
	}

	// Sanity check response
	if v, ok := res["number"]; !ok {
		t.Fatal("missing 'number' field")
	} else if n, err := hexutil.DecodeBig(v.(string)); err != nil || n == nil {
		t.Fatal(err)
	} else if n.Cmp(big.NewInt(1)) != 0 {
		t.Fatalf("unexpected 'latest' block number: %v", n)
	}
	// 'transactions' key should exist as []
	if v, ok := res["transactions"]; !ok {
		t.Fatal("missing transactions field")
	} else if len(v.([]interface{})) != 0 {
		t.Fatal("'transactions' value not []")
	}
	// 'uncles' key should exist as []
	if v, ok := res["uncles"]; !ok {
		t.Fatal("missing uncles field")
	} else if len(v.([]interface{})) != 0 {
		t.Fatal("'uncles' value not []'")
	}
}

func TestHeader_PendingNull(t *testing.T) {
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

	for _, fullTxes := range []bool{true, false} {
		gotPending := make(map[string]interface{})
		err := client.CallContext(ctx, &gotPending, "eth_getBlockByNumber", "pending", fullTxes)
		if err != nil {
			t.Fatal(err)
		}

		// Iterate expected values, checking validity.
		wantBlockNumber := big.NewInt(2)
		want := map[string]interface{}{
			// Nulls.
			"nonce": nil,
			"hash":  nil,
			"miner": nil,
			// "totalDifficulty":  (*hexutil.Big)(new(big.Int).Mul(vars.MinimumDifficulty, wantBlockNumber)).String(),
			"totalDifficulty": nil,

			// Zero-values.
			"logsBloom":    "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
			"mixHash":      "0x0000000000000000000000000000000000000000000000000000000000000000",
			"transactions": []interface{}{},
			"uncles":       []interface{}{},

			// Filled.
			"number":           (*hexutil.Big)(wantBlockNumber).String(),
			"gasLimit":         "0x47d5cc",
			"gasUsed":          hexutil.Uint64(0).String(),
			"difficulty":       (*hexutil.Big)(parent.Difficulty).String(),
			"size":             "0x21a",
			"parentHash":       parent.Hash().Hex(),
			"extraData":        "0xda83010b1788436f72654765746886676f312e3135856c696e7578",
			"transactionsRoot": "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
			"receiptsRoot":     "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
			"stateRoot":        "0x02189854bc38ea675df81794e54a2676230444d87adc7a51bbba0d4cc6519d43",
			"sha3Uncles":       "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
			"timestamp":        "0x604f76f4", // Incidentally nondeterministic; special case.
		}
		for k, v := range want {
			gotVal, ok := gotPending[k]
			if !ok {
				t.Errorf("%s: missing key", k)
			}
			// Special case (indeterminate time).
			if k == "timestamp" {
				if !regexp.MustCompile(fmt.Sprintf(`^0x[a-zA-Z0-9]{%d}`, len("604f76f4"))).MatchString(gotVal.(string)) {
					t.Errorf("%s: unexpected value: %v", k, gotVal)
				}
				gotVal = want["timestamp"]
				gotPending["timestamp"] = gotVal
			}
			if !reflect.DeepEqual(v, gotVal) {
				t.Errorf("%s: want: %v, got: %v", k, v, gotVal)
			}
		}

		// Slightly redundant, but additionally checks the occurrence of got -> want (supplementing want -> got).
		for _, diff := range deep.Equal(want, gotPending) {
			t.Errorf("[want/got] +/-: %s", diff)
		}
	}
}
