package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/ethash"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/params/types/genesisT"
	"github.com/ethereum/go-ethereum/params/vars"
)

func TestBlockChain_EIP161_Export(t *testing.T) {
	var (
		// Our signer
		key1, _ = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		// Empty 1
		key2, _ = crypto.HexToECDSA("8a1f9a8f95be41cd7ccb6168179afb4504aefe388d1e14474d32c45c72ce7b7a")
		// Empty 2
		key3, _ = crypto.HexToECDSA("49a7b37aa6f6645917e7b807e9d1c00d4fa71f18343b0d4122a4d2df64dd6fee")

		addr1 = crypto.PubkeyToAddress(key1.PublicKey)
		addr2 = crypto.PubkeyToAddress(key2.PublicKey)
		addr3 = crypto.PubkeyToAddress(key3.PublicKey)
		db    = rawdb.NewMemoryDatabase()
	)

	printPrettyJSON := func(any interface{}) {
		b, err := json.MarshalIndent(any, "", "    ")
		if err != nil {
			t.Fatal(err)
		}
		t.Log(string(b))
	}

	config := params.TestChainConfig
	config.ChainID.Set(big.NewInt(1111))
	config.EIP158Block.Set(big.NewInt(2))
	config.ByzantiumBlock.Set(big.NewInt(2))
	config.ConstantinopleBlock.Set(big.NewInt(2))
	config.PetersburgBlock.Set(big.NewInt(2))
	config.IstanbulBlock.Set(big.NewInt(2))
	config.BerlinBlock.Set(big.NewInt(2))
	config.LondonBlock.Set(big.NewInt(2))
	config.ArrowGlacierBlock.Set(big.NewInt(2))
	printPrettyJSON(config)

	genesis := params.DefaultMessNetGenesisBlock()
	genesis.Config = config

	// emptyWithStorage
	genesis.Alloc[addr1] = genesisT.GenesisAccount{
		Balance: new(big.Int).Mul(big.NewInt(vars.Ether), big.NewInt(10)),
	}
	genesis.Alloc[addr2] = genesisT.GenesisAccount{
		Balance: common.Big0,
		Storage: map[common.Hash]common.Hash{
			common.BytesToHash([]byte("a")): common.BytesToHash([]byte("myvalue1")),
		}}
	// emptyWithoutStorage
	genesis.Alloc[addr3] = genesisT.GenesisAccount{
		Balance: common.Big0,
		Storage: map[common.Hash]common.Hash{
			common.BytesToHash([]byte("b")): common.BytesToHash([]byte("myvalue2")),
		}}

	printPrettyJSON(genesis)

	genesisBlock := MustCommitGenesis(db, genesis)

	// This call generates a chain of 5 blocks. The function runs for
	// each block and adds different features to gen based on the
	// block index.
	chain, _ := GenerateChain(genesis.Config, genesisBlock, ethash.NewFaker(), db, 3, func(i int, gen *BlockGen) {
		switch i {
		case 0:
			signer := types.NewEIP155Signer(config.ChainID)
			// In block 1, addr1 sends addr2 some ether.
			signedTx, _ := types.SignTx(types.NewTransaction(gen.TxNonce(addr1), addr2, big.NewInt(0), vars.TxGas, big.NewInt(20), nil), signer, key1)
			gen.AddTx(signedTx)

		case 1:
			signer := types.NewEIP1559Signer(config.ChainID)
			baseTx := &types.DynamicFeeTx{
				To:        &addr2,
				Nonce:     gen.TxNonce(addr1),
				GasFeeCap: big.NewInt(vars.InitialBaseFee),
				GasTipCap: big.NewInt(20),
				Gas:       vars.TxGas,
				Value:     big.NewInt(0),
				Data:      nil,
			}
			tx := types.NewTx(baseTx)
			signedTx, _ := types.SignTx(tx, signer, key1)
			gen.AddTx(signedTx)
		}
	})

	// Import the chain. This runs all block validation rules.
	blockchain, _ := NewBlockChain(db, nil, genesis.Config, ethash.NewFaker(), vm.Config{}, nil, nil)
	defer blockchain.Stop()

	if i, err := blockchain.InsertChain(chain); err != nil {
		fmt.Printf("insert error (block %d): %v\n", chain[i].NumberU64(), err)
		return
	}
	t.Logf("last block: #%d\n", blockchain.CurrentBlock().Number())

	state, _ := blockchain.StateAt(blockchain.GetBlockByNumber(1).Root())
	t.Log(1)
	t.Log(state.Exist(addr1))
	t.Log(state.Exist(addr2))
	t.Log(state.Exist(addr3))

	dump := state.Dump(nil)
	t.Log(string(dump))

	state, _ = blockchain.StateAt(blockchain.GetBlockByNumber(2).Root())
	t.Log(2)
	t.Log(state.Exist(addr1))
	t.Log(state.Exist(addr2))
	t.Log(state.Exist(addr3))

	dump = state.Dump(nil)
	t.Log(string(dump))

	// All files will be written to the /tmp directory.
	// Write the genesis file.
	genesisJSON, err := json.MarshalIndent(genesis, "", "    ")
	if err != nil {
		t.Fatal(err)
	}
	err = ioutil.WriteFile("/tmp/genesis.json", genesisJSON, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	// Establish the export file writer.
	f, err := os.OpenFile("/tmp/export.rlp", os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	err = blockchain.Export(f)
	if err != nil {
		t.Fatal(err)
	}
}
