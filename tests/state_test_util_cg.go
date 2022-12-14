package tests

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/params/types/genesisT"
)

type stPre map[common.Address]stPreAccount

func (p stPre) toGenesisAlloc() genesisT.GenesisAlloc {
	genesisAlloc := make(genesisT.GenesisAlloc)
	for addr, acc := range p {
		genesisAlloc[addr] = genesisT.GenesisAccount{
			Code:    acc.Code,
			Nonce:   acc.Nonce,
			Balance: acc.Balance,
			Storage: acc.Storage,
		}
	}
	return genesisAlloc
}

//go:generate go run github.com/fjl/gencodec -type stPreAccount -field-override stPreAccountMarshaling -out gen_stpreaccount.go

// stPreAccount is structurally equivalent to genesisT.GenesisAccount, but
// removes the 'omitempty' json tags for code, storage, and balance, since
// the ethereum/tests JSON schema defines these annotations as 'required'.
// See ethereum/tests/JSONSchema/st-schema.json.
type stPreAccount struct {
	// GenesisAccount is an account in the state of the genesis block.
	Code       []byte                      `json:"code"`
	Storage    map[common.Hash]common.Hash `json:"storage,omitempty"`
	Balance    *big.Int                    `json:"balance" gencodec:"required"`
	Nonce      uint64                      `json:"nonce"`
	PrivateKey []byte                      `json:"secretKey,omitempty"` // for tests
}

type stPreAccountMarshaling struct {
	Code       hexutil.Bytes
	Balance    *math.HexOrDecimal256
	Nonce      math.HexOrDecimal64
	Storage    map[storageJSON]storageJSON
	PrivateKey hexutil.Bytes
}

// storageJSON represents a 256 bit byte array, but allows less than 256 bits when
// unmarshaling from hex.
type storageJSON common.Hash

func (h *storageJSON) UnmarshalText(text []byte) error {
	text = bytes.TrimPrefix(text, []byte("0x"))
	if len(text) > 64 {
		return fmt.Errorf("too many hex characters in storage key/value %q", text)
	}
	offset := len(h) - len(text)/2 // pad on the left
	if _, err := hex.Decode(h[offset:], text); err != nil {
		fmt.Println(err)
		return fmt.Errorf("invalid hex storage key/value %q", text)
	}
	return nil
}

func (h storageJSON) MarshalText() ([]byte, error) {
	return hexutil.Bytes(h[:]).MarshalText()
}
