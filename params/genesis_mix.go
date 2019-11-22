package params

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/params/types"
)

// MixGenesisBlock returns the Mix genesis block.
func DefaultMixGenesisBlock() *paramtypes.Genesis {
	return &paramtypes.Genesis{
		Config:     MixChainConfig,
		Nonce:      0x1391eaa92b871f91,
		ExtraData:  hexutil.MustDecode("0x77656c636f6d65746f7468656c696e6b6564776f726c64000000000000000000"),
		GasLimit:   3000000,
		Difficulty: big.NewInt(1048576),
		Alloc:      paramtypes.DecodePreAlloc(MixAllocData),
	}
}
