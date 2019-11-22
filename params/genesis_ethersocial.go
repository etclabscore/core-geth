package params

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

// EthersocialGenesisBlock returns the Ethersocial main net genesis block.
func DefaultEthersocialGenesisBlock() *Genesis {
	return &Genesis{
		Config:     EthersocialChainConfig,
		Nonce:      66,
		ExtraData:  hexutil.MustDecode("0x"),
		GasLimit:   3141592,
		Difficulty: big.NewInt(131072),
		Alloc:      decodePrealloc(EthersocialAllocData),
	}
}
