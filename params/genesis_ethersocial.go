package params

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/params/types"
)

// EthersocialGenesisBlock returns the Ethersocial main net genesis block.
func DefaultEthersocialGenesisBlock() *paramtypes.Genesis {
	return &paramtypes.Genesis{
		Config:     EthersocialChainConfig,
		Nonce:      66,
		ExtraData:  hexutil.MustDecode("0x"),
		GasLimit:   3141592,
		Difficulty: big.NewInt(131072),
		Alloc:      paramtypes.DecodePreAlloc(EthersocialAllocData),
	}
}
