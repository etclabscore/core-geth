package params

import (
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/params/types/genesisT"
	"github.com/ethereum/go-ethereum/rlp"
)

var MintMeGenesisHash = common.HexToHash("0x13d1952f29df2a702b3beae629e5b8297c56e401c4c1094bccb3e5febae099c3")

func DecodeMintmeAlloc(data string) genesisT.GenesisAlloc {
	var p []struct {
		Addr    *big.Int
		Balance *big.Int
		Code    []byte
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
			Code:    account.Code,
			Storage: storage,
			//Nonce: account.Nonce,
		}
	}
	return ga
}

func DefaultMintMeGenesisBlock() *genesisT.Genesis {
	return &genesisT.Genesis{
		Config:     MintMeChainConfig,
		Nonce:      hexutil.MustDecodeUint64("0x0"),
		ExtraData:  hexutil.MustDecode("0x42"),
		GasLimit:   hexutil.MustDecodeUint64("0x2fefd8"),
		Difficulty: hexutil.MustDecodeBig("0x200000"),
		Timestamp:  1623664673,
		Alloc:      DecodeMintmeAlloc(allocMintme),
	}
}
