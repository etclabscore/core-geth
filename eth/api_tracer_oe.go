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

package eth

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/consensus/ethash"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/internal/ethapi"
	"github.com/ethereum/go-ethereum/rpc"
)

// OpenEthereumTrace A trace in the desired format (Parity/OpenEtherum) See: https://openethereum.github.io/wiki/JSONRPC-trace-module
type OpenEthereumTrace struct {
	Action              TraceRewardAction `json:"action"`
	BlockHash           *common.Hash      `json:"blockHash"`
	BlockNumber         uint64            `json:"blockNumber"`
	Error               string            `json:"error,omitempty"`
	Result              interface{}       `json:"result"`
	Subtraces           int               `json:"subtraces"`
	TraceAddress        []int             `json:"traceAddress"`
	TransactionHash     *common.Hash      `json:"transactionHash"`
	TransactionPosition *uint64           `json:"transactionPosition"`
	Type                string            `json:"type"`
}

// TraceRewardAction An OpenEthereum formatted trace reward action
type TraceRewardAction struct {
	Value      *hexutil.Big    `json:"value,omitempty"`
	Author     *common.Address `json:"author,omitempty"`
	RewardType string          `json:"rewardType,omitempty"`
}

// setConfigTracerToOpenEthereum forces the Tracer to the OpenEthereum one, mutating the argument value.
func setConfigTracerToOpenEthereum(config *TraceConfig) *TraceConfig {
	if config == nil {
		config = &TraceConfig{}
	}

	tracer := "callTracerOpenEthereum"
	config.Tracer = &tracer
	return config
}

func traceBlockReward(ctx context.Context, eth *Ethereum, block *types.Block, config *TraceConfig) (*OpenEthereumTrace, error) {
	chainConfig := eth.blockchain.Config()
	minerReward, _ := ethash.AccumulateRewards(chainConfig, block.Header(), block.Uncles())

	coinbase := block.Coinbase()
	blockHash := block.Hash()

	tr := &OpenEthereumTrace{
		Type: "reward",
		Action: TraceRewardAction{
			Value:      (*hexutil.Big)(minerReward),
			Author:     &coinbase,
			RewardType: "block",
		},
		TraceAddress: []int{},
		BlockHash:    &blockHash,
		BlockNumber:  block.NumberU64(),
	}

	return tr, nil
}

func traceBlockUncleRewards(ctx context.Context, eth *Ethereum, block *types.Block, config *TraceConfig) ([]*OpenEthereumTrace, error) {
	chainConfig := eth.blockchain.Config()
	_, uncleRewards := ethash.AccumulateRewards(chainConfig, block.Header(), block.Uncles())

	blockHash := block.Hash()

	results := make([]*OpenEthereumTrace, len(uncleRewards))
	for i, uncle := range block.Uncles() {
		if i < len(uncleRewards) {
			coinbase := uncle.Coinbase

			results[i] = &OpenEthereumTrace{
				Type: "reward",
				Action: TraceRewardAction{
					Value:      (*hexutil.Big)(uncleRewards[i]),
					Author:     &coinbase,
					RewardType: "uncle",
				},
				TraceAddress: []int{},
				BlockNumber:  block.NumberU64(),
				BlockHash:    &blockHash,
			}
		}
	}

	return results, nil
}

// Block returns the structured logs created during the execution of
// EVM and returns them as a JSON object.
// The correct name will be TraceBlockByNumber, though we want to be compatible with OpenEthereum trace module.
func (api *PrivateTraceAPI) Block(ctx context.Context, number rpc.BlockNumber, config *TraceConfig) ([]interface{}, error) {
	// Fetch the block that we want to trace
	var block *types.Block

	switch number {
	case rpc.PendingBlockNumber:
		block = api.eth.miner.PendingBlock()
	case rpc.LatestBlockNumber:
		block = api.eth.blockchain.CurrentBlock()
	default:
		block = api.eth.blockchain.GetBlockByNumber(uint64(number))
	}
	// Trace the block if it was found
	if block == nil {
		return nil, fmt.Errorf("block #%d not found", number)
	}

	config = setConfigTracerToOpenEthereum(config)

	traceResults, err := traceBlockByNumber(ctx, api.eth, number, config)
	if err != nil {
		return nil, err
	}

	traceReward, err := traceBlockReward(ctx, api.eth, block, config)
	if err != nil {
		return nil, err
	}

	traceUncleRewards, err := traceBlockUncleRewards(ctx, api.eth, block, config)
	if err != nil {
		return nil, err
	}

	results := []interface{}{}

	for _, result := range traceResults {
		var tmp []interface{}
		if err := json.Unmarshal(result.Result.(json.RawMessage), &tmp); err != nil {
			return nil, err
		}
		results = append(results, tmp...)
	}

	results = append(results, traceReward)

	for _, uncleReward := range traceUncleRewards {
		results = append(results, uncleReward)
	}

	return results, nil
}

// Transaction returns the structured logs created during the execution of EVM
// and returns them as a JSON object.
func (api *PrivateTraceAPI) Transaction(ctx context.Context, hash common.Hash, config *TraceConfig) (interface{}, error) {
	setConfigTracerToOpenEthereum(config)
	return traceTransaction(ctx, api.eth, hash, config)
}

func (api *PrivateTraceAPI) Filter(ctx context.Context, args ethapi.CallArgs, config *TraceConfig) ([]*txTraceResult, error) {
	setConfigTracerToOpenEthereum(config)
	fmt.Printf("args: %#v\n", args)
	return nil, nil
}
