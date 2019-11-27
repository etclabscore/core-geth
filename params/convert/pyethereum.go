package convert

import (
	"encoding/binary"
	"errors"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/params/types"
	"github.com/ethereum/go-ethereum/params/types/pyethereum"
)

// NewPyEthereumGenesisSpec converts a go-ethereum genesis block into a Parity specific
// chain specification format.
func NewPyEthereumGenesisSpec(network string, genesis *paramtypes.Genesis) (*pyethereum.PyEthereumGenesisSpec, error) {
	// Only ethash is currently supported between go-ethereum and pyethereum
	if genesis.Config.(*paramtypes.MultiGethChainConfig).Ethash == nil {
		return nil, errors.New("unsupported consensus engine")
	}
	spec := &pyethereum.PyEthereumGenesisSpec{
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
