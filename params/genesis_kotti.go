package params

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

// DefaultKottiGenesisBlock returns the Kotti network genesis block.
func DefaultKottiGenesisBlock() *Genesis {
	return &Genesis{
		Config:     KottiChainConfig,
		Timestamp:  1546461831,
		ExtraData:  hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000000025b7955e43adf9c2a01a9475908702cce67f302a6aaf8cba3c9255a2b863415d4db7bae4f4bbca020000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
		GasLimit:   10485760,
		Difficulty: big.NewInt(1),
		Alloc:      decodePrealloc(KottiAllocData),
	}
}
