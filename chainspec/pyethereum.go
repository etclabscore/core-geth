package chainspec

import (
	"encoding/binary"
	"errors"

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

// NewPyEthereumGenesisSpec converts a go-ethereum genesis block into a Parity specific
// chain specification format.
func NewPyEthereumGenesisSpec(network string, genesis *paramtypes.Genesis) (*PyEthereumGenesisSpec, error) {
	// Only ethash is currently supported between go-ethereum and pyethereum
	if genesis.Config.Ethash == nil {
		return nil, errors.New("unsupported consensus engine")
	}
	spec := &PyEthereumGenesisSpec{
		Timestamp:  (hexutil.Uint64)(genesis.Timestamp),
		ExtraData:  genesis.ExtraData,
		GasLimit:   (hexutil.Uint64)(genesis.GasLimit),
		Difficulty: (*hexutil.Big)(genesis.Difficulty),
		Mixhash:    genesis.Mixhash,
		Coinbase:   genesis.Coinbase,
		Alloc:      genesis.Alloc,
		ParentHash: genesis.ParentHash,
	}
	spec.Nonce = (hexutil.Bytes)(make([]byte, 8))
	binary.LittleEndian.PutUint64(spec.Nonce[:], genesis.Nonce)

	return spec, nil
}
