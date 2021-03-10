package params

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/params/types/genesisT"
)

var ABCGenesisHash = common.HexToHash("0x5f32ce1ed875a04d74361164bcfcc24af721df8e616642486338300a691fe582")

func DefaultABCGenesisBlock() *genesisT.Genesis {
	return &genesisT.Genesis{
		Config:     ABCChainConfig,
		Nonce:      hexutil.MustDecodeUint64("0x0"),
		ExtraData:  hexutil.MustDecode("0x42"),
		GasLimit:   hexutil.MustDecodeUint64("0x2fefd8"),
		Difficulty: hexutil.MustDecodeBig("0x20000"),
		Timestamp:  1615385980,
		Alloc: genesisT.GenesisAlloc{
			common.HexToAddress("366ae7da62294427c764870bd2a460d7ded29d30"): genesisT.GenesisAccount{
				Balance: big.NewInt(42),
			},
		},
	}
}
