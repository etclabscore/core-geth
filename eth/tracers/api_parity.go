// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package tracers

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/internal/ethapi"
	"github.com/ethereum/go-ethereum/params/mutations"
	"github.com/ethereum/go-ethereum/rpc"
)

// TraceFilterArgs represents the arguments for a call.
type TraceFilterArgs struct {
	FromBlock   hexutil.Uint64  `json:"fromBlock,omitempty"`   // Trace from this starting block
	ToBlock     hexutil.Uint64  `json:"toBlock,omitempty"`     // Trace utill this end block
	FromAddress *common.Address `json:"fromAddress,omitempty"` // Sent from these addresses
	ToAddress   *common.Address `json:"toAddress,omitempty"`   // Sent to these addresses
	After       uint64          `json:"after,omitempty"`       // The offset trace number
	Count       uint64          `json:"count,omitempty"`       // Integer number of traces to display in a batch
}

// ParityTrace A trace in the desired format (Parity/OpenEtherum) See: https://Parity.github.io/wiki/JSONRPC-trace-module
type ParityTrace struct {
	Action              TraceRewardAction `json:"action"`
	BlockHash           common.Hash       `json:"blockHash"`
	BlockNumber         uint64            `json:"blockNumber"`
	Error               string            `json:"error,omitempty"`
	Result              interface{}       `json:"result"`
	Subtraces           int               `json:"subtraces"`
	TraceAddress        []int             `json:"traceAddress"`
	TransactionHash     *common.Hash      `json:"transactionHash"`
	TransactionPosition *uint64           `json:"transactionPosition"`
	Type                string            `json:"type"`
}

// TraceRewardAction An Parity formatted trace reward action
type TraceRewardAction struct {
	Value      *hexutil.Big    `json:"value,omitempty"`
	Author     *common.Address `json:"author,omitempty"`
	RewardType string          `json:"rewardType,omitempty"`
}

// setTraceConfigDefaultTracer sets the default tracer to "callTracerParity" if none set
func setTraceConfigDefaultTracer(config *TraceConfig) *TraceConfig {
	if config == nil {
		config = &TraceConfig{}
	}

	if config.Tracer == nil {
		tracer := "callTracerParity"
		config.Tracer = &tracer
	}

	return config
}

// setTraceCallConfigDefaultTracer sets the default tracer to "callTracerParity" if none set
func setTraceCallConfigDefaultTracer(config *TraceCallConfig) *TraceCallConfig {
	if config == nil {
		config = &TraceCallConfig{}
	}

	if config.Tracer == nil {
		tracer := "callTracerParity"
		config.Tracer = &tracer
	}

	return config
}

// TraceAPI is the collection of Ethereum full node APIs exposed over
// the private debugging endpoint.
type TraceAPI struct {
	debugAPI *API
}

// NewTraceAPI creates a new API definition for the full node-related
// private debug methods of the Ethereum service.
func NewTraceAPI(debugAPI *API) *TraceAPI {
	return &TraceAPI{debugAPI: debugAPI}
}

// decorateResponse applies formatting to trace results if needed.
func decorateResponse(res interface{}, config *TraceConfig) (interface{}, error) {
	if config != nil && config.NestedTraceOutput && config.Tracer != nil {
		return decorateNestedTraceResponse(res, *config.Tracer), nil
	}
	return res, nil
}

// decorateNestedTraceResponse formats trace results the way Parity does.
// Docs: https://openethereum.github.io/JSONRPC-trace-module
// Example:
/*
{
  "id": 1,
  "jsonrpc": "2.0",
  "result": {
    "output": "0x",
    "stateDiff": { ... },
    "trace": [ { ... }, ],
    "vmTrace": { ... }
  }
}
*/
func decorateNestedTraceResponse(res interface{}, tracer string) interface{} {
	out := map[string]interface{}{}
	if tracer == "callTracerParity" {
		out["trace"] = res
	} else if tracer == "stateDiffTracer" {
		out["stateDiff"] = res
	} else {
		return res
	}
	return out
}

// traceBlockReward retrieve the block reward for the coinbase address
func (api *TraceAPI) traceBlockReward(ctx context.Context, block *types.Block, config *TraceConfig) (*ParityTrace, error) {
	chainConfig := api.debugAPI.backend.ChainConfig()
	minerReward, _ := mutations.GetRewards(chainConfig, block.Header(), block.Uncles())

	coinbase := block.Coinbase()

	tr := &ParityTrace{
		Type: "reward",
		Action: TraceRewardAction{
			Value:      (*hexutil.Big)(minerReward),
			Author:     &coinbase,
			RewardType: "block",
		},
		TraceAddress: []int{},
		BlockHash:    block.Hash(),
		BlockNumber:  block.NumberU64(),
	}

	return tr, nil
}

// traceBlockUncleRewards retrieve the block rewards for uncle addresses
func (api *TraceAPI) traceBlockUncleRewards(ctx context.Context, block *types.Block, config *TraceConfig) ([]*ParityTrace, error) {
	chainConfig := api.debugAPI.backend.ChainConfig()
	_, uncleRewards := mutations.GetRewards(chainConfig, block.Header(), block.Uncles())

	results := make([]*ParityTrace, len(uncleRewards))
	for i, uncle := range block.Uncles() {
		if i < len(uncleRewards) {
			coinbase := uncle.Coinbase

			results[i] = &ParityTrace{
				Type: "reward",
				Action: TraceRewardAction{
					Value:      (*hexutil.Big)(uncleRewards[i]),
					Author:     &coinbase,
					RewardType: "uncle",
				},
				TraceAddress: []int{},
				BlockNumber:  block.NumberU64(),
				BlockHash:    block.Hash(),
			}
		}
	}

	return results, nil
}

// Block returns the structured logs created during the execution of
// EVM and returns them as a JSON object.
// The correct name will be TraceBlockByNumber, though we want to be compatible with Parity trace module.
func (api *TraceAPI) Block(ctx context.Context, number rpc.BlockNumber, config *TraceConfig) ([]interface{}, error) {
	config = setTraceConfigDefaultTracer(config)

	block, err := api.debugAPI.blockByNumber(ctx, number)
	if err != nil {
		return nil, err
	}

	traceResults, err := api.debugAPI.traceBlock(ctx, block, config)
	if err != nil {
		return nil, err
	}

	traceReward, err := api.traceBlockReward(ctx, block, config)
	if err != nil {
		return nil, err
	}

	traceUncleRewards, err := api.traceBlockUncleRewards(ctx, block, config)
	if err != nil {
		return nil, err
	}

	results := []interface{}{}

	for _, result := range traceResults {
		if result.Error != "" {
			return nil, errors.New(result.Error)
		}
		var tmp interface{}
		if err := json.Unmarshal(result.Result.(json.RawMessage), &tmp); err != nil {
			return nil, err
		}
		if *config.Tracer == "stateDiffTracer" {
			results = append(results, tmp)
		} else {
			results = append(results, tmp.([]interface{})...)
		}
	}

	results = append(results, traceReward)

	for _, uncleReward := range traceUncleRewards {
		results = append(results, uncleReward)
	}

	return results, nil
}

// Transaction returns the structured logs created during the execution of EVM
// and returns them as a JSON object.
func (api *TraceAPI) Transaction(ctx context.Context, hash common.Hash, config *TraceConfig) (interface{}, error) {
	config = setTraceConfigDefaultTracer(config)
	return api.debugAPI.TraceTransaction(ctx, hash, config)
}

// Filter configures a new tracer according to the provided configuration, and
// executes all the transactions contained within. The return value will be one item
// per transaction, dependent on the requested tracer.
func (api *TraceAPI) Filter(ctx context.Context, args TraceFilterArgs, config *TraceConfig) (*rpc.Subscription, error) {
	config = setTraceConfigDefaultTracer(config)

	// Fetch the block interval that we want to trace
	start := rpc.BlockNumber(args.FromBlock)
	end := rpc.BlockNumber(args.ToBlock)

	return api.debugAPI.TraceChain(ctx, start, end, config)
}

// Call lets you trace a given eth_call. It collects the structured logs created during the execution of EVM
// if the given transaction was added on top of the provided block and returns them as a JSON object.
// You can provide -2 as a block number to trace on top of the pending block.
func (api *TraceAPI) Call(ctx context.Context, args ethapi.TransactionArgs, blockNrOrHash rpc.BlockNumberOrHash, config *TraceCallConfig) (interface{}, error) {
	config = setTraceCallConfigDefaultTracer(config)
	res, err := api.debugAPI.TraceCall(ctx, args, blockNrOrHash, config)
	if err != nil {
		return nil, err
	}
	traceConfig := getTraceConfigFromTraceCallConfig(config)
	return decorateResponse(res, traceConfig)
}

// CallMany lets you trace a given eth_call. It collects the structured logs created during the execution of EVM
// if the given transaction was added on top of the provided block and returns them as a JSON object.
// You can provide -2 as a block number to trace on top of the pending block.
func (api *TraceAPI) CallMany(ctx context.Context, txs []ethapi.TransactionArgs, blockNrOrHash rpc.BlockNumberOrHash, config *TraceCallConfig) (interface{}, error) {
	config = setTraceCallConfigDefaultTracer(config)
	return api.debugAPI.TraceCallMany(ctx, txs, blockNrOrHash, config)
}
