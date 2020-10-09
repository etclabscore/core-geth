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
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/internal/ethapi"
	"github.com/ethereum/go-ethereum/rpc"
)

// setConfigTracerToOpenEthereum forces the Tracer to the OpenEthereum one
func setConfigTracerToOpenEthereum(config *TraceConfig) *TraceConfig {
	if config == nil {
		config = &TraceConfig{}
	}

	tracer := "callTracerOpenEthereum"
	config.Tracer = &tracer
	return config
}

// Block returns the structured logs created during the execution of
// EVM and returns them as a JSON object.
// The correct name will be TraceBlockByNumber, though we want to be compatible with OpenEthereum trace module.
func (api *PrivateTraceAPI) Block(ctx context.Context, number rpc.BlockNumber, config *TraceConfig) ([]*txTraceResult, error) {
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
	return traceBlockByNumber(ctx, api.eth, number, config)
}

// Transaction returns the structured logs created during the execution of EVM
// and returns them as a JSON object.
func (api *PrivateTraceAPI) Transaction(ctx context.Context, hash common.Hash, config *TraceConfig) (interface{}, error) {
	config = setConfigTracerToOpenEthereum(config)
	return traceTransaction(ctx, api.eth, hash, config)
}

func (api *PrivateTraceAPI) Filter(ctx context.Context, args ethapi.CallArgs, config *TraceConfig) ([]*txTraceResult, error) {
	config = setConfigTracerToOpenEthereum(config)
	fmt.Printf("args: %#v\n", args)
	return nil, nil
}
