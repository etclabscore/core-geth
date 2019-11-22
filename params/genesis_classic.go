package params

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

// ClassicGenesisBlock returns the Ethereum Classic genesis block.
func DefaultClassicGenesisBlock() *Genesis {
	return &Genesis{
		Config:     ClassicChainConfig,
		Nonce:      66,
		ExtraData:  hexutil.MustDecode("0x11bbe8db4e347b4e8c937c1c8370e4b5ed33adb3db69cbdb7a38e1e50b1b82fa"),
		GasLimit:   5000,
		Difficulty: big.NewInt(17179869184),
		Alloc:      decodePrealloc(MainnetAllocData),
	}
}
