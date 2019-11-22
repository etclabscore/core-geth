package chainspec

import (
	"encoding/binary"
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/params"
)

// AlethGenesisSpec represents the genesis specification format used by the
// C++ Ethereum implementation.
type AlethGenesisSpec struct {
	SealEngine string `json:"sealEngine"`
	Params     struct {
		AccountStartNonce          math.HexOrDecimal64   `json:"accountStartNonce"`
		MaximumExtraDataSize       hexutil.Uint64        `json:"maximumExtraDataSize"`
		HomesteadForkBlock         *hexutil.Big          `json:"homesteadForkBlock,omitempty"`
		DaoHardforkBlock           math.HexOrDecimal64   `json:"daoHardforkBlock"`
		EIP150ForkBlock            *hexutil.Big          `json:"EIP150ForkBlock,omitempty"`
		EIP158ForkBlock            *hexutil.Big          `json:"EIP158ForkBlock,omitempty"`
		ByzantiumForkBlock         *hexutil.Big          `json:"byzantiumForkBlock,omitempty"`
		ConstantinopleForkBlock    *hexutil.Big          `json:"constantinopleForkBlock,omitempty"`
		ConstantinopleFixForkBlock *hexutil.Big          `json:"constantinopleFixForkBlock,omitempty"`
		IstanbulForkBlock          *hexutil.Big          `json:"istanbulForkBlock,omitempty"`
		MinGasLimit                hexutil.Uint64        `json:"minGasLimit"`
		MaxGasLimit                hexutil.Uint64        `json:"maxGasLimit"`
		TieBreakingGas             bool                  `json:"tieBreakingGas"`
		GasLimitBoundDivisor       math.HexOrDecimal64   `json:"gasLimitBoundDivisor"`
		MinimumDifficulty          *hexutil.Big          `json:"minimumDifficulty"`
		DifficultyBoundDivisor     *math.HexOrDecimal256 `json:"difficultyBoundDivisor"`
		DurationLimit              *math.HexOrDecimal256 `json:"durationLimit"`
		BlockReward                *hexutil.Big          `json:"blockReward"`
		NetworkID                  hexutil.Uint64        `json:"networkID"`
		ChainID                    hexutil.Uint64        `json:"chainID"`
		AllowFutureBlocks          bool                  `json:"allowFutureBlocks"`
	} `json:"params"`
	Genesis struct {
		Nonce      hexutil.Bytes  `json:"nonce"`
		Difficulty *hexutil.Big   `json:"difficulty"`
		MixHash    common.Hash    `json:"mixHash"`
		Author     common.Address `json:"author"`
		Timestamp  hexutil.Uint64 `json:"timestamp"`
		ParentHash common.Hash    `json:"parentHash"`
		ExtraData  hexutil.Bytes  `json:"extraData"`
		GasLimit   hexutil.Uint64 `json:"gasLimit"`
	} `json:"genesis"`

	Accounts map[common.UnprefixedAddress]*alethGenesisSpecAccount `json:"accounts"`
}


// alethGenesisSpecAccount is the prefunded genesis account and/or precompiled
// contract definition.
type alethGenesisSpecAccount struct {
	Balance     *math.HexOrDecimal256   `json:"balance,omitempty"`
	Nonce       uint64                   `json:"nonce,omitempty"`
	Precompiled *alethGenesisSpecBuiltin `json:"precompiled,omitempty"`
}

// alethGenesisSpecBuiltin is the precompiled contract definition.
type alethGenesisSpecBuiltin struct {
	Name          string                         `json:"name,omitempty"`
	StartingBlock *hexutil.Big                   `json:"startingBlock,omitempty"`
	Linear        *alethGenesisSpecLinearPricing `json:"linear,omitempty"`
}

type alethGenesisSpecLinearPricing struct {
	Base uint64 `json:"base"`
	Word uint64 `json:"word"`
}

// NewAlethGenesisSpec converts a go-ethereum genesis block into a Aleth-specific
// chain specification format.
func NewAlethGenesisSpec(network string, genesis *params.Genesis) (*AlethGenesisSpec, error) {
	// Only ethash is currently supported between go-ethereum and aleth
	if genesis.Config.Ethash == nil {
		return nil, errors.New("unsupported consensus engine")
	}
	// Reconstruct the chain spec in Aleth format
	spec := &AlethGenesisSpec{
		SealEngine: "Ethash",
	}
	// Some defaults
	spec.Params.AccountStartNonce = 0
	spec.Params.TieBreakingGas = false
	spec.Params.AllowFutureBlocks = false

	// Dao hardfork block is a special one. The fork block is listed as 0 in the
	// config but aleth will sync with ETC clients up until the actual dao hard
	// fork block.
	spec.Params.DaoHardforkBlock = 0

	if num := genesis.Config.HomesteadBlock; num != nil {
		spec.Params.HomesteadForkBlock = (*hexutil.Big)(num)
	}
	if num := genesis.Config.EIP150Block; num != nil {
		spec.Params.EIP150ForkBlock = (*hexutil.Big)(num)
	}
	if num := genesis.Config.EIP158Block; num != nil {
		spec.Params.EIP158ForkBlock = (*hexutil.Big)(num)
	}
	if num := genesis.Config.ByzantiumBlock; num != nil {
		spec.Params.ByzantiumForkBlock = (*hexutil.Big)(num)
	}
	if num := genesis.Config.ConstantinopleBlock; num != nil {
		spec.Params.ConstantinopleForkBlock = (*hexutil.Big)(num)
	}
	if num := genesis.Config.PetersburgBlock; num != nil {
		spec.Params.ConstantinopleFixForkBlock = (*hexutil.Big)(num)
	}
	if num := genesis.Config.IstanbulBlock; num != nil {
		spec.Params.IstanbulForkBlock = (*hexutil.Big)(num)
	}
	spec.Params.NetworkID = (hexutil.Uint64)(genesis.Config.ChainID.Uint64())
	spec.Params.ChainID = (hexutil.Uint64)(genesis.Config.ChainID.Uint64())
	spec.Params.MaximumExtraDataSize = (hexutil.Uint64)(params.MaximumExtraDataSize)
	spec.Params.MinGasLimit = (hexutil.Uint64)(params.MinGasLimit)
	spec.Params.MaxGasLimit = (hexutil.Uint64)(math.MaxInt64)
	spec.Params.MinimumDifficulty = (*hexutil.Big)(params.MinimumDifficulty)
	spec.Params.DifficultyBoundDivisor = (*math.HexOrDecimal256)(params.DifficultyBoundDivisor)
	spec.Params.GasLimitBoundDivisor = (math.HexOrDecimal64)(params.GasLimitBoundDivisor)
	spec.Params.DurationLimit = (*math.HexOrDecimal256)(params.DurationLimit)
	spec.Params.BlockReward = (*hexutil.Big)(params.FrontierBlockReward)

	spec.Genesis.Nonce = (hexutil.Bytes)(make([]byte, 8))
	binary.LittleEndian.PutUint64(spec.Genesis.Nonce[:], genesis.Nonce)

	spec.Genesis.MixHash = genesis.Mixhash
	spec.Genesis.Difficulty = (*hexutil.Big)(genesis.Difficulty)
	spec.Genesis.Author = genesis.Coinbase
	spec.Genesis.Timestamp = (hexutil.Uint64)(genesis.Timestamp)
	spec.Genesis.ParentHash = genesis.ParentHash
	spec.Genesis.ExtraData = (hexutil.Bytes)(genesis.ExtraData)
	spec.Genesis.GasLimit = (hexutil.Uint64)(genesis.GasLimit)

	for address, account := range genesis.Alloc {
		spec.setAccount(address, account)
	}

	spec.setPrecompile(1, &alethGenesisSpecBuiltin{Name: "ecrecover",
		Linear: &alethGenesisSpecLinearPricing{Base: 3000}})
	spec.setPrecompile(2, &alethGenesisSpecBuiltin{Name: "sha256",
		Linear: &alethGenesisSpecLinearPricing{Base: 60, Word: 12}})
	spec.setPrecompile(3, &alethGenesisSpecBuiltin{Name: "ripemd160",
		Linear: &alethGenesisSpecLinearPricing{Base: 600, Word: 120}})
	spec.setPrecompile(4, &alethGenesisSpecBuiltin{Name: "identity",
		Linear: &alethGenesisSpecLinearPricing{Base: 15, Word: 3}})
	if genesis.Config.ByzantiumBlock != nil {
		spec.setPrecompile(5, &alethGenesisSpecBuiltin{Name: "modexp",
			StartingBlock: (*hexutil.Big)(genesis.Config.ByzantiumBlock)})
		spec.setPrecompile(6, &alethGenesisSpecBuiltin{Name: "alt_bn128_G1_add",
			StartingBlock: (*hexutil.Big)(genesis.Config.ByzantiumBlock),
			Linear:        &alethGenesisSpecLinearPricing{Base: 500}})
		spec.setPrecompile(7, &alethGenesisSpecBuiltin{Name: "alt_bn128_G1_mul",
			StartingBlock: (*hexutil.Big)(genesis.Config.ByzantiumBlock),
			Linear:        &alethGenesisSpecLinearPricing{Base: 40000}})
		spec.setPrecompile(8, &alethGenesisSpecBuiltin{Name: "alt_bn128_pairing_product",
			StartingBlock: (*hexutil.Big)(genesis.Config.ByzantiumBlock)})
	}
	if genesis.Config.IstanbulBlock != nil {
		if genesis.Config.ByzantiumBlock == nil {
			return nil, errors.New("invalid genesis, istanbul fork is enabled while byzantium is not")
		}
		spec.setPrecompile(6, &alethGenesisSpecBuiltin{
			Name:          "alt_bn128_G1_add",
			StartingBlock: (*hexutil.Big)(genesis.Config.ByzantiumBlock),
		}) // Aleth hardcoded the gas policy
		spec.setPrecompile(7, &alethGenesisSpecBuiltin{
			Name:          "alt_bn128_G1_mul",
			StartingBlock: (*hexutil.Big)(genesis.Config.ByzantiumBlock),
		}) // Aleth hardcoded the gas policy
		spec.setPrecompile(9, &alethGenesisSpecBuiltin{
			Name:          "blake2_compression",
			StartingBlock: (*hexutil.Big)(genesis.Config.IstanbulBlock),
		})
	}
	return spec, nil
}

func (spec *AlethGenesisSpec) setPrecompile(address byte, data *alethGenesisSpecBuiltin) {
	if spec.Accounts == nil {
		spec.Accounts = make(map[common.UnprefixedAddress]*alethGenesisSpecAccount)
	}
	addr := common.UnprefixedAddress(common.BytesToAddress([]byte{address}))
	if _, exist := spec.Accounts[addr]; !exist {
		spec.Accounts[addr] = &alethGenesisSpecAccount{}
	}
	spec.Accounts[addr].Precompiled = data
}

func (spec *AlethGenesisSpec) setAccount(address common.Address, account params.GenesisAccount) {
	if spec.Accounts == nil {
		spec.Accounts = make(map[common.UnprefixedAddress]*alethGenesisSpecAccount)
	}

	a, exist := spec.Accounts[common.UnprefixedAddress(address)]
	if !exist {
		a = &alethGenesisSpecAccount{}
		spec.Accounts[common.UnprefixedAddress(address)] = a
	}
	a.Balance = (*math.HexOrDecimal256)(account.Balance)
	a.Nonce = account.Nonce

}