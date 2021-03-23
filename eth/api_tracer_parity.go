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
	"github.com/ethereum/go-ethereum/eth/tracers"
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
func setConfigTracerToParity(config *tracers.TraceConfig) *tracers.TraceConfig {
	if config == nil {
		config = &tracers.TraceConfig{}
	}

	tracer := "callTracerParity"
	config.Tracer = &tracer
	return config
}

func traceBlockReward(ctx context.Context, eth *Ethereum, block *types.Block, config *tracers.TraceConfig) (*ParityTrace, error) {
	chainConfig := eth.blockchain.Config()
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

func traceBlockUncleRewards(ctx context.Context, eth *Ethereum, block *types.Block, config *tracers.TraceConfig) ([]*ParityTrace, error) {
	chainConfig := eth.blockchain.Config()
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
func (api *PrivateTraceAPI) Block(ctx context.Context, number rpc.BlockNumber, config *tracers.TraceConfig) ([]interface{}, error) {
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

	config = setConfigTracerToParity(config)

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
	config = setConfigTracerToParity(config)
	return traceTransaction(ctx, api.eth, hash, config)
}

func (api *PrivateTraceAPI) Filter(ctx context.Context, args ethapi.CallArgs, config *TraceConfig) ([]*txTraceResult, error) {
	// config = setConfigTracerToParity(config)
	return nil, nil
}

// ------------------ DANGER, DRAGONS
//
// // TraceBlockByNumber returns the structured logs created during the execution of
// // EVM and returns them as a JSON object.
// func traceBlockByNumber(ctx context.Context, eth *Ethereum, number rpc.BlockNumber, config *TraceConfig) ([]*txTraceResult, error) {
// 	// Fetch the block that we want to trace
// 	var block *types.Block
//
// 	switch number {
// 	case rpc.PendingBlockNumber:
// 		block = eth.miner.PendingBlock()
// 	case rpc.LatestBlockNumber:
// 		block = eth.blockchain.CurrentBlock()
// 	default:
// 		block = eth.blockchain.GetBlockByNumber(uint64(number))
// 	}
// 	// Trace the block if it was found
// 	if block == nil {
// 		return nil, fmt.Errorf("block #%d not found", number)
// 	}
// 	return traceBlock(ctx, eth, block, config)
// }
//
// // traceBlock configures a new tracer according to the provided configuration, and
// // executes all the transactions contained within. The return value will be one item
// // per transaction, dependent on the requestd tracer.
// func traceBlock(ctx context.Context, eth *Ethereum, block *types.Block, config *TraceConfig) ([]*txTraceResult, error) {
// 	// Create the parent state database
// 	if err := eth.engine.VerifyHeader(eth.blockchain, block.Header(), true); err != nil {
// 		return nil, err
// 	}
// 	parent := eth.blockchain.GetBlock(block.ParentHash(), block.NumberU64()-1)
// 	if parent == nil {
// 		return nil, fmt.Errorf("parent %#x not found", block.ParentHash())
// 	}
// 	reexec := defaultTraceReexec
// 	if config != nil && config.Reexec != nil {
// 		reexec = *config.Reexec
// 	}
// 	statedb, err := computeStateDB(eth, parent, reexec)
// 	if err != nil {
// 		return nil, err
// 	}
// 	// Execute all the transaction contained within the block concurrently
// 	var (
// 		signer = types.MakeSigner(eth.blockchain.Config(), block.Number())
//
// 		txs     = block.Transactions()
// 		results = make([]*txTraceResult, len(txs))
//
// 		pend = new(sync.WaitGroup)
// 		jobs = make(chan *txTraceTask, len(txs))
// 	)
// 	threads := runtime.NumCPU()
// 	if threads > len(txs) {
// 		threads = len(txs)
// 	}
// 	for th := 0; th < threads; th++ {
// 		pend.Add(1)
// 		go func() {
// 			defer pend.Done()
//
// 			// Fetch and execute the next transaction trace tasks
// 			for task := range jobs {
// 				msg, _ := txs[task.index].AsMessage(signer)
// 				vmctx := core.NewEVMContext(msg, block.Header(), eth.blockchain, nil)
//
// 				res, err := traceTx(ctx, eth, msg, vmctx, task.statedb, task.taskExtraContext, config)
// 				if err != nil {
// 					results[task.index] = &txTraceResult{Error: err.Error()}
// 					continue
// 				}
// 				results[task.index] = &txTraceResult{Result: res}
// 			}
// 		}()
// 	}
// 	// Feed the transactions into the tracers and return
// 	var failed error
// 	for i, tx := range txs {
// 		taskExtraContext := map[string]interface{}{
// 			"blockNumber":         block.NumberU64(),
// 			"blockHash":           block.Hash().Hex(),
// 			"transactionHash":     tx.Hash().Hex(),
// 			"transactionPosition": uint64(i),
// 		}
//
// 		// Send the trace task over for execution
// 		jobs <- &txTraceTask{statedb: statedb.Copy(), index: i, taskExtraContext: taskExtraContext}
//
// 		// Generate the next state snapshot fast without tracing
// 		msg, _ := tx.AsMessage(signer)
// 		vmctx := core.NewEVMContext(msg, block.Header(), eth.blockchain, nil)
//
// 		vmenv := vm.NewEVM(vmctx, statedb, eth.blockchain.Config(), vm.Config{})
// 		if _, err := core.ApplyMessage(vmenv, msg, new(core.GasPool).AddGas(msg.Gas())); err != nil {
// 			failed = err
// 			break
// 		}
// 		// Finalize the state so any modifications are written to the trie
// 		// Only delete empty objects if EIP158/161 (a.k.a Spurious Dragon) is in effect
// 		statedb.Finalise(vmenv.ChainConfig().IsEnabled(vmenv.ChainConfig().GetEIP161dTransition, block.Number()))
// 	}
// 	close(jobs)
// 	pend.Wait()
//
// 	// If execution failed in between, abort
// 	if failed != nil {
// 		return nil, failed
// 	}
// 	return results, nil
// }
//
// func (api *PrivateDebugAPI) traceBlock(ctx context.Context, block *types.Block, config *TraceConfig) ([]*txTraceResult, error) {
// 	return traceBlock(ctx, api.eth, block, config)
// }
