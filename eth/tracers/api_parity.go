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
	"github.com/ethereum/go-ethereum/consensus/ethash"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/internal/ethapi"
	"github.com/ethereum/go-ethereum/rpc"
)

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

// setConfigTracerToParity forces the Tracer to the Parity one
func setConfigTracerToParity(config *TraceConfig) *TraceConfig {
	if config == nil {
		config = &TraceConfig{}
	}

	tracer := "callTracerParity"
	config.Tracer = &tracer
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

// TraceBlockReward retrieve the block reward for the coinbase address
func (api *TraceAPI) TraceBlockReward(ctx context.Context, block *types.Block, config *TraceConfig) (*ParityTrace, error) {
	chainConfig := api.debugAPI.backend.ChainConfig()
	minerReward, _ := ethash.GetRewards(chainConfig, block.Header(), block.Uncles())

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

// TraceBlockUncleRewards retrieve the block rewards for uncle addresses
func (api *TraceAPI) TraceBlockUncleRewards(ctx context.Context, block *types.Block, config *TraceConfig) ([]*ParityTrace, error) {
	chainConfig := api.debugAPI.backend.ChainConfig()
	_, uncleRewards := ethash.GetRewards(chainConfig, block.Header(), block.Uncles())

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
	config = setConfigTracerToParity(config)

	block, err := api.debugAPI.blockByNumber(ctx, number)
	if err != nil {
		return nil, err
	}

	traceResults, err := api.debugAPI.traceBlock(ctx, block, config)
	if err != nil {
		return nil, err
	}

	traceReward, err := api.TraceBlockReward(ctx, block, config)
	if err != nil {
		return nil, err
	}

	traceUncleRewards, err := api.TraceBlockUncleRewards(ctx, block, config)
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
func (api *TraceAPI) Transaction(ctx context.Context, hash common.Hash, config *TraceConfig) (interface{}, error) {
	config = setConfigTracerToParity(config)
	return api.debugAPI.TraceTransaction(ctx, hash, config)
}

func (api *TraceAPI) Filter(ctx context.Context, args ethapi.CallArgs, config *TraceConfig) ([]*txTraceResult, error) {
	return nil, errors.New("not implemented")
}
