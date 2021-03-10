package params

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/params/types/genesisT"
)

// FIXME: This is not yet correct.
var ABCGenesisHash = common.HexToHash("0xa68ebde7932eccb177d38d55dcc6461a019dd795a681e59b5a3e4f3a7259a3f1")

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
