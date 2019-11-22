package pyethereum

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/params/types"
)

// PyEthereumGenesisSpec represents the genesis specification format used by the
// Python Ethereum implementation.
type PyEthereumGenesisSpec struct {
	Nonce      hexutil.Bytes           `json:"nonce"`
	Timestamp  hexutil.Uint64          `json:"timestamp"`
	ExtraData  hexutil.Bytes           `json:"extraData"`
	GasLimit   hexutil.Uint64          `json:"gasLimit"`
	Difficulty *hexutil.Big            `json:"difficulty"`
	Mixhash    common.Hash             `json:"mixhash"`
	Coinbase   common.Address          `json:"coinbase"`
	Alloc      paramtypes.GenesisAlloc `json:"alloc"`
	ParentHash common.Hash             `json:"parentHash"`
}

