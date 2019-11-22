package params

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/params/types"
)

// SocialGenesisBlock returns the Ethereum Social genesis block.
func DefaultSocialGenesisBlock() *paramtypes.Genesis {
	return &paramtypes.Genesis{
		Config:     SocialChainConfig,
		Nonce:      66,
		ExtraData:  hexutil.MustDecode("0x3230313820457468657265756d20536f6369616c2050726f6a656374"),
		GasLimit:   5000,
		Difficulty: big.NewInt(17179869184),
		Alloc:      paramtypes.DecodePreAlloc(SocialAllocData),
	}
}
