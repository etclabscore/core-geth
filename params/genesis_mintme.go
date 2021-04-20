package params

import (
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/params/types/genesisT"
	"github.com/ethereum/go-ethereum/rlp"
)

var MINTMEGenesisHash = common.HexToHash("0x5f32ce1ed875a04d74361164bcfcc24af721df8e616642486338300a691fe582") // TODO: lumat

func DecodeMintmeAlloc(data string) genesisT.GenesisAlloc {
	var p []struct{
		Addr *big.Int
		Balance *big.Int
		Code []byte
		Storage []string
		//Nonce uint64
	}
	if err := rlp.NewStream(strings.NewReader(data), 0).Decode(&p); err != nil {
		panic(err)
	}
	ga := make(genesisT.GenesisAlloc, len(p))
	for _, account := range p {
		storage := make(map[common.Hash]common.Hash)
		for i := 0; i < len(account.Storage); i += 2 {
			storage[common.HexToHash(account.Storage[i])] = common.HexToHash(account.Storage[i+1])
		}

		ga[common.BigToAddress(account.Addr)] = genesisT.GenesisAccount{
			Balance: account.Balance,
			Code: account.Code,
			Storage: storage,
			//Nonce: account.Nonce,
		}
	}
	return ga
}

func DefaultMINTMEGenesisBlock() *genesisT.Genesis {
	// TODO: lumat
	return &genesisT.Genesis{
		Config:     MINTMEChainConfig,
		Nonce:      hexutil.MustDecodeUint64("0x0"),
		ExtraData:  hexutil.MustDecode("0x42"),
		GasLimit:   hexutil.MustDecodeUint64("0x2fefd8"),
		Difficulty: hexutil.MustDecodeBig("0x20000"),
		Timestamp:  1615385980,
		Alloc: DecodeMintmeAlloc(allocMintme),
	}
}
