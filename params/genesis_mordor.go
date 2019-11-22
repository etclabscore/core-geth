package params

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/params/types"
)

func DefaultMordorGenesisBlock() *paramtypes.Genesis {
	return &paramtypes.Genesis{
		Config:     MordorChainConfig,
		Nonce:      hexutil.MustDecodeUint64("0x0"),
		ExtraData:  hexutil.MustDecode("0x70686f656e697820636869636b656e206162737572642062616e616e61"),
		GasLimit:   hexutil.MustDecodeUint64("0x2fefd8"),
		Difficulty: hexutil.MustDecodeBig("0x20000"),
		Timestamp:  hexutil.MustDecodeUint64("0x5d9676db"),
		Alloc:      paramtypes.GenesisAlloc{},
	}
}
