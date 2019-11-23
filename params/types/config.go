package paramtypes

import (
	"encoding/binary"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	common2 "github.com/ethereum/go-ethereum/params/types/common"
)

// TrustedCheckpoint represents a set of post-processed trie roots (CHT and
// BloomTrie) associated with the appropriate section index and head hash. It is
// used to start light syncing from this checkpoint and avoid downloading the
// entire header chain while still being able to securely access old headers/logs.
type TrustedCheckpoint struct {
	SectionIndex uint64      `json:"sectionIndex"`
	SectionHead  common.Hash `json:"sectionHead"`
	CHTRoot      common.Hash `json:"chtRoot"`
	BloomRoot    common.Hash `json:"bloomRoot"`
}

// HashEqual returns an indicator comparing the itself hash with given one.
func (c *TrustedCheckpoint) HashEqual(hash common.Hash) bool {
	if c.Empty() {
		return hash == common.Hash{}
	}
	return c.Hash() == hash
}

// Hash returns the hash of checkpoint's four key fields(index, sectionHead, chtRoot and bloomTrieRoot).
func (c *TrustedCheckpoint) Hash() common.Hash {
	buf := make([]byte, 8+3*common.HashLength)
	binary.BigEndian.PutUint64(buf, c.SectionIndex)
	copy(buf[8:], c.SectionHead.Bytes())
	copy(buf[8+common.HashLength:], c.CHTRoot.Bytes())
	copy(buf[8+2*common.HashLength:], c.BloomRoot.Bytes())
	return crypto.Keccak256Hash(buf)
}

// Empty returns an indicator whether the checkpoint is regarded as empty.
func (c *TrustedCheckpoint) Empty() bool {
	return c.SectionHead == (common.Hash{}) || c.CHTRoot == (common.Hash{}) || c.BloomRoot == (common.Hash{})
}

// CheckpointOracleConfig represents a set of checkpoint contract(which acts as an oracle)
// config which used for light client checkpoint syncing.
type CheckpointOracleConfig struct {
	Address   common.Address   `json:"address"`
	Signers   []common.Address `json:"signers"`
	Threshold uint64           `json:"threshold"`
}

// ChainConfig is the core config which determines the blockchain settings.
//
// ChainConfig is stored in the database on a per block basis. This means
// that any network, identified by its genesis block, can have its own
// set of configuration options.
type ChainConfig struct {
	NetworkID uint64   `json:"networkId"`
	ChainID   *big.Int `json:"chainId"` // chainId identifies the current chain and is used for replay protection

	// HF: Homestead
	HomesteadBlock *big.Int `json:"homesteadBlock,omitempty"` // Homestead switch block (nil = no fork, 0 = already homestead)
	// "Homestead Hard-fork Changes"
	// https://github.com/ethereum/EIPs/blob/master/EIPS/eip-2.md
	EIP2FBlock *big.Int `json:"eip2FBlock,omitempty"`
	// DELEGATECALL
	// https://eips.ethereum.org/EIPS/eip-7
	EIP7FBlock *big.Int `json:"eip7FBlock,omitempy"`
	// Note: EIP 8 was also included in this fork, but was not backwards-incompatible

	// HF: DAO
	DAOForkBlock   *big.Int `json:"daoForkBlock,omitempty"`   // TheDAO hard-fork switch block (nil = no fork)
	DAOForkSupport bool     `json:"daoForkSupport,omitempty"` // Whether the nodes supports or opposes the DAO hard-fork

	// HF: Tangerine Whistle
	// EIP150 implements the Gas price changes (https://github.com/ethereum/EIPs/issues/150)
	EIP150Block *big.Int    `json:"eip150Block,omitempty"` // EIP150 HF block (nil = no fork)
	EIP150Hash  common.Hash `json:"eip150Hash,omitempty"`  // EIP150 HF hash (needed for header only clients as only gas pricing changed)

	// HF: Spurious Dragon
	EIP155Block *big.Int `json:"eip155Block,omitempty"` // EIP155 HF block
	EIP158Block *big.Int `json:"eip158Block,omitempty"` // EIP158 HF block, includes implementations of 158/161, 160, and 170
	//
	// EXP cost increase
	// https://github.com/ethereum/EIPs/blob/master/EIPS/eip-160.md
	// NOTE: this json tag:
	// (a.) varies from it's 'siblings', which have 'F's in them
	// (b.) without the 'F' will vary from ETH implementations if they choose to accept the proposed changes
	// with corresponding refactoring (https://github.com/ethereum/go-ethereum/pull/18401)
	EIP160FBlock *big.Int `json:"eip160Block,omitempty"`
	// State trie clearing (== EIP158 proper)
	// https://github.com/ethereum/EIPs/blob/master/EIPS/eip-161.md
	EIP161FBlock *big.Int `json:"eip161FBlock,omitempty"`
	// Contract code size limit
	// https://github.com/ethereum/EIPs/blob/master/EIPS/eip-170.md
	EIP170FBlock *big.Int `json:"eip170FBlock,omitempty"`

	// HF: Byzantium
	ByzantiumBlock *big.Int `json:"byzantiumBlock,omitempty"` // Byzantium switch block (nil = no fork, 0 = already on byzantium)
	//
	// Difficulty adjustment to target mean block time including uncles
	// https://github.com/ethereum/EIPs/issues/100
	EIP100FBlock *big.Int `json:"eip100FBlock,omitempty"`
	// Opcode REVERT
	// https://eips.ethereum.org/EIPS/eip-140
	EIP140FBlock *big.Int `json:"eip140FBlock,omitempty"`
	// Precompiled contract for bigint_modexp
	// https://github.com/ethereum/EIPs/issues/198
	EIP198FBlock *big.Int `json:"eip198FBlock,omitempty"`
	// Opcodes RETURNDATACOPY, RETURNDATASIZE
	// https://github.com/ethereum/EIPs/issues/211
	EIP211FBlock *big.Int `json:"eip211FBlock,omitempty"`
	// Precompiled contract for pairing check
	// https://github.com/ethereum/EIPs/issues/212
	EIP212FBlock *big.Int `json:"eip212FBlock,omitempty"`
	// Precompiled contracts for addition and scalar multiplication on the elliptic curve alt_bn128
	// https://github.com/ethereum/EIPs/issues/213
	EIP213FBlock *big.Int `json:"eip213FBlock,omitempty"`
	// Opcode STATICCALL
	// https://github.com/ethereum/EIPs/issues/214
	EIP214FBlock *big.Int `json:"eip214FBlock,omitempty"`
	// Metropolis diff bomb delay and reducing block reward
	// https://github.com/ethereum/EIPs/issues/649
	// note that this is closely related to EIP100.
	// In fact, EIP100 is bundled in
	EIP649FBlock *big.Int `json:"eip649FBlock,omitempty"`
	// Transaction receipt status
	// https://github.com/ethereum/EIPs/issues/658
	EIP658FBlock *big.Int `json:"eip658FBlock,omitempty"`
	// NOT CONFIGURABLE: prevent overwriting contracts
	// https://github.com/ethereum/EIPs/issues/684
	// EIP684FBlock *big.Int `json:"eip684BFlock,omitempty"`

	// HF: Constantinople
	ConstantinopleBlock *big.Int `json:"constantinopleBlock,omitempty"` // Constantinople switch block (nil = no fork, 0 = already activated)
	//
	// Opcodes SHR, SHL, SAR
	// https://eips.ethereum.org/EIPS/eip-145
	EIP145FBlock *big.Int `json:"eip145FBlock,omitempty"`
	// Opcode CREATE2
	// https://eips.ethereum.org/EIPS/eip-1014
	EIP1014FBlock *big.Int `json:"eip1014FBlock,omitempty"`
	// Opcode EXTCODEHASH
	// https://eips.ethereum.org/EIPS/eip-1052
	EIP1052FBlock *big.Int `json:"eip1052FBlock,omitempty"`
	// Constantinople difficulty bomb delay and block reward adjustment
	// https://eips.ethereum.org/EIPS/eip-1234
	EIP1234FBlock *big.Int `json:"eip1234FBlock,omitempty"`
	// Net gas metering
	// https://eips.ethereum.org/EIPS/eip-1283
	EIP1283FBlock *big.Int `json:"eip1283FBlock,omitempty"`

	PetersburgBlock *big.Int `json:"petersburgBlock,omitempty"` // Petersburg switch block (nil = same as Constantinople)

	// HF: Istanbul
	IstanbulBlock *big.Int `json:"istanbulBlock,omitempty"` // Istanbul switch block (nil = no fork, 0 = already on istanbul)
	//
	// EIP-152: Add Blake2 compression function F precompile
	EIP152FBlock *big.Int `json:"eip152FBlock,omitempty"`
	// EIP-1108: Reduce alt_bn128 precompile gas costs
	EIP1108FBlock *big.Int `json="eip1108FBlock,omitempty"`
	// EIP-1344: Add ChainID opcode
	EIP1344FBlock *big.Int `json="eip1344FBlock,omitempty"`
	// EIP-1884: Repricing for trie-size-dependent opcodes
	EIP1884FBlock *big.Int `json="eip1884FBlock,omitempty"`
	// EIP-2028: Calldata gas cost reduction
	EIP2028FBlock *big.Int `json="eip2028FBlock,omitempty"`
	// EIP-2200: Rebalance net-metered SSTORE gas cost with consideration of SLOAD gas cost change
	EIP2200FBlock *big.Int `json="eip2200FBlock,omitempty"`

	EWASMBlock *big.Int `json:"ewasmBlock,omitempty"` // EWASM switch block (nil = no fork, 0 = already activated)

	ECIP1010PauseBlock *big.Int `json:"ecip1010PauseBlock,omitempty"` // ECIP1010 pause HF block
	ECIP1010Length     *big.Int `json:"ecip1010Length,omitempty"`     // ECIP1010 length
	ECIP1017FBlock     *big.Int `json:"ecip1017FBlock,omitempty"`
	ECIP1017EraRounds  *big.Int `json:"ecip1017EraRounds,omitempty"` // ECIP1017 era rounds
	DisposalBlock      *big.Int `json:"disposalBlock,omitempty"`     // Bomb disposal HF block
	SocialBlock        *big.Int `json:"socialBlock,omitempty"`       // Ethereum Social Reward block
	EthersocialBlock   *big.Int `json:"ethersocialBlock,omitempty"`  // Ethersocial Reward block

	MCIP0Block *big.Int `json:"mcip0Block,omitempty"` // Musicoin default block; no MCIP, just denotes chain pref
	MCIP3Block *big.Int `json:"mcip3Block,omitempty"` // Musicoin 'UBI Fork' block
	MCIP8Block *big.Int `json:"mcip8Block,omitempty"` // Musicoin 'QT For' block

	// Various consensus engines
	Ethash *EthashConfig `json:"ethash,omitempty"`
	Clique *CliqueConfig `json:"clique,omitempty"`

	TrustedCheckpoint       *TrustedCheckpoint      `json:"trustedCheckpoint"`
	TrustedCheckpointOracle *CheckpointOracleConfig `json:"trustedCheckpointOracle"`

	DifficultyBombDelaySchedule common2.Uint64BigMapEncodesHex `json:"difficultyBombDelays,omitempty"'` // JSON tag matches Parity's
	BlockRewardSchedule         common2.Uint64BigMapEncodesHex `json:"blockReward,omitempty"`           // JSON tag matches Parity's
}

// EthashConfig is the consensus engine configs for proof-of-work based sealing.
type EthashConfig struct{}

// String implements the stringer interface, returning the consensus engine details.
func (c *EthashConfig) String() string {
	return "ethash"
}

// CliqueConfig is the consensus engine configs for proof-of-authority based sealing.
type CliqueConfig struct {
	Period uint64 `json:"period"` // Number of seconds between blocks to enforce
	Epoch  uint64 `json:"epoch"`  // Epoch length to reset votes and checkpoint
}

// String implements the stringer interface, returning the consensus engine details.
func (c *CliqueConfig) String() string {
	return "clique"
}

// String implements the fmt.Stringer interface.
func (c *ChainConfig) String() string {
	var engine interface{}
	switch {
	case c.Ethash != nil:
		engine = c.Ethash
	case c.Clique != nil:
		engine = c.Clique
	default:
		engine = "unknown"
	}
	return fmt.Sprintf("{NetworkID: %v, ChainID: %v Homestead: %v DAO: %v DAOSupport: %v EIP150: %v EIP155: %v EIP158: %v Byzantium: %v Disposal: %v Social: %v Ethersocial: %v ECIP1017: %v EIP160: %v ECIP1010PauseBlock: %v ECIP1010Length: %v Constantinople: %v ConstantinopleFix: %v Istanbul: %v Engine: %v}",
		c.NetworkID,
		c.ChainID,
		c.HomesteadBlock,
		c.DAOForkBlock,
		c.DAOForkSupport,
		c.EIP150Block,
		c.EIP155Block,
		c.EIP158Block,
		c.ByzantiumBlock,
		c.DisposalBlock,
		c.SocialBlock,
		c.EthersocialBlock,
		c.ECIP1017EraRounds,
		c.EIP160FBlock,
		c.ECIP1010PauseBlock,
		c.ECIP1010Length,
		c.ConstantinopleBlock,
		c.PetersburgBlock,
		c.IstanbulBlock,
		engine,
	)
}

// IsECIP1017F returns whether the chain is configured with ECIP1017.
func (c *ChainConfig) IsECIP1017F(num *big.Int) bool {
	return IsForked(c.ECIP1017FBlock, num) || c.ECIP1017EraRounds != nil
}

// IsEIP2F returns whether num is equal to or greater than the Homestead or EIP2 block.
func (c *ChainConfig) IsEIP2F(num *big.Int) bool {
	return IsForked(c.HomesteadBlock, num) || IsForked(c.EIP2FBlock, num)
}

// IsEIP7F returns whether num is equal to or greater than the Homestead or EIP7 block.
func (c *ChainConfig) IsEIP7F(num *big.Int) bool {
	return IsForked(c.HomesteadBlock, num) || IsForked(c.EIP7FBlock, num)
}

// IsDAOFork returns whether num is either equal to the DAO fork block or greater.
func (c *ChainConfig) IsDAOFork(num *big.Int) bool {
	return IsForked(c.DAOForkBlock, num)
}

// IsEIP150 returns whether num is either equal to the EIP150 fork block or greater.
func (c *ChainConfig) IsEIP150(num *big.Int) bool {
	return IsForked(c.EIP150Block, num)
}

// IsEIP155 returns whether num is either equal to the EIP155 fork block or greater.
func (c *ChainConfig) IsEIP155(num *big.Int) bool {
	return IsForked(c.EIP155Block, num)
}

// IsEIP160F returns whether num is either equal to or greater than the "EIP158HF" Block or EIP160 block.
func (c *ChainConfig) IsEIP160F(num *big.Int) bool {
	return IsForked(c.EIP158Block, num) || IsForked(c.EIP160FBlock, num)
}

// IsEIP161F returns whether num is either equal to or greater than the "EIP158HF" Block or EIP161 block.
func (c *ChainConfig) IsEIP161F(num *big.Int) bool {
	return IsForked(c.EIP158Block, num) || IsForked(c.EIP161FBlock, num)
}

// IsEIP170F returns whether num is either equal to or greater than the "EIP158HF" Block or EIP170 block.
func (c *ChainConfig) IsEIP170F(num *big.Int) bool {
	return IsForked(c.EIP158Block, num) || IsForked(c.EIP170FBlock, num)
}

// IsEIP100F returns whether num is equal to or greater than the Byzantium or EIP100 block.
func (c *ChainConfig) IsEIP100F(num *big.Int) bool {
	return IsForked(c.ByzantiumBlock, num) || IsForked(c.ConstantinopleBlock, num) || IsForked(c.EIP100FBlock, num)
}

// IsEIP140F returns whether num is equal to or greater than the Byzantium or EIP140 block.
func (c *ChainConfig) IsEIP140F(num *big.Int) bool {
	return IsForked(c.ByzantiumBlock, num) || IsForked(c.EIP140FBlock, num)
}

// IsEIP198F returns whether num is equal to or greater than the Byzantium or EIP198 block.
func (c *ChainConfig) IsEIP198F(num *big.Int) bool {
	return IsForked(c.ByzantiumBlock, num) || IsForked(c.EIP198FBlock, num)
}

// IsEIP211F returns whether num is equal to or greater than the Byzantium or EIP211 block.
func (c *ChainConfig) IsEIP211F(num *big.Int) bool {
	return IsForked(c.ByzantiumBlock, num) || IsForked(c.EIP211FBlock, num)
}

// IsEIP212F returns whether num is equal to or greater than the Byzantium or EIP212 block.
func (c *ChainConfig) IsEIP212F(num *big.Int) bool {
	return IsForked(c.ByzantiumBlock, num) || IsForked(c.EIP212FBlock, num)
}

// IsEIP213F returns whether num is equal to or greater than the Byzantium or EIP213 block.
func (c *ChainConfig) IsEIP213F(num *big.Int) bool {
	return IsForked(c.ByzantiumBlock, num) || IsForked(c.EIP213FBlock, num)
}

// IsEIP214F returns whether num is equal to or greater than the Byzantium or EIP214 block.
func (c *ChainConfig) IsEIP214F(num *big.Int) bool {
	return IsForked(c.ByzantiumBlock, num) || IsForked(c.EIP214FBlock, num)
}

// IsEIP649F returns whether num is equal to or greater than the Byzantium or EIP649 block.
func (c *ChainConfig) IsEIP649F(num *big.Int) bool {
	return IsForked(c.ByzantiumBlock, num) || IsForked(c.EIP649FBlock, num)
}

// IsEIP658F returns whether num is equal to or greater than the Byzantium or EIP658 block.
func (c *ChainConfig) IsEIP658F(num *big.Int) bool {
	return IsForked(c.ByzantiumBlock, num) || IsForked(c.EIP658FBlock, num)
}

// IsEIP145F returns whether num is equal to or greater than the Constantinople or EIP145 block.
func (c *ChainConfig) IsEIP145F(num *big.Int) bool {
	return IsForked(c.ConstantinopleBlock, num) || IsForked(c.EIP145FBlock, num)
}

// IsEIP1014F returns whether num is equal to or greater than the Constantinople or EIP1014 block.
func (c *ChainConfig) IsEIP1014F(num *big.Int) bool {
	return IsForked(c.ConstantinopleBlock, num) || IsForked(c.EIP1014FBlock, num)
}

// IsEIP1052F returns whether num is equal to or greater than the Constantinople or EIP1052 block.
func (c *ChainConfig) IsEIP1052F(num *big.Int) bool {
	return IsForked(c.ConstantinopleBlock, num) || IsForked(c.EIP1052FBlock, num)
}

// IsEIP1234F returns whether num is equal to or greater than the Constantinople or EIP1234 block.
func (c *ChainConfig) IsEIP1234F(num *big.Int) bool {
	return IsForked(c.ConstantinopleBlock, num) || IsForked(c.EIP1234FBlock, num)
}

// IsEIP1283F returns whether num is equal to or greater than the Constantinople or EIP1283 block.
func (c *ChainConfig) IsEIP1283F(num *big.Int) bool {
	return !c.IsPetersburg(num) && (IsForked(c.ConstantinopleBlock, num) || IsForked(c.EIP1283FBlock, num))
}

// IsEIP152F returns whether num is equal to or greater than the Istanbul block.
func (c *ChainConfig) IsEIP152F(num *big.Int) bool {
	return IsForked(c.IstanbulBlock, num) || IsForked(c.EIP152FBlock, num)
}

// IsEIP1108F returns whether num is equal to or greater than the Istanbul block.
func (c *ChainConfig) IsEIP1108F(num *big.Int) bool {
	return IsForked(c.IstanbulBlock, num) || IsForked(c.EIP1108FBlock, num)
}

// IsEIP1344F returns whether num is equal to or greater than the Istanbul block.
func (c *ChainConfig) IsEIP1344F(num *big.Int) bool {
	return IsForked(c.IstanbulBlock, num) || IsForked(c.EIP1344FBlock, num)
}

// IsEIP1884F returns whether num is equal to or greater than the Istanbul block.
func (c *ChainConfig) IsEIP1884F(num *big.Int) bool {
	return IsForked(c.IstanbulBlock, num) || IsForked(c.EIP1884FBlock, num)
}

// IsEIP2028F returns whether num is equal to or greater than the Istanbul block.
func (c *ChainConfig) IsEIP2028F(num *big.Int) bool {
	return IsForked(c.IstanbulBlock, num) || IsForked(c.EIP2028FBlock, num)
}

// IsEIP2200F returns whether num is equal to or greater than the Istanbul block.
func (c *ChainConfig) IsEIP2200F(num *big.Int) bool {
	return IsForked(c.IstanbulBlock, num) || IsForked(c.EIP2200FBlock, num)
}

func (c *ChainConfig) IsBombDisposal(num *big.Int) bool {
	return IsForked(c.DisposalBlock, num)
}

func (c *ChainConfig) IsECIP1010(num *big.Int) bool {
	return IsForked(c.ECIP1010PauseBlock, num)
}

// IsPetersburg returns whether num is either
// - equal to or greater than the PetersburgBlock fork block,
// - OR is nil, and Constantinople is active
func (c *ChainConfig) IsPetersburg(num *big.Int) bool {
	return IsForked(c.PetersburgBlock, num) || c.PetersburgBlock == nil && IsForked(c.ConstantinopleBlock, num)
}

// IsIstanbul returns whether num is either equal to the Istanbul fork block or greater.
func (c *ChainConfig) IsIstanbul(num *big.Int) bool {
	return IsForked(c.IstanbulBlock, num)
}

// IsEWASM returns whether num represents a block number after the EWASM fork
func (c *ChainConfig) IsEWASM(num *big.Int) bool {
	return IsForked(c.EWASMBlock, num)
}

// CheckCompatible checks whether scheduled fork transitions have been imported
// with a mismatching chain configuration.
func (c *ChainConfig) CheckCompatible(newcfg *ChainConfig, height uint64) *ConfigCompatError {
	bhead := new(big.Int).SetUint64(height)

	// Iterate checkCompatible to find the lowest conflict.
	var lasterr *ConfigCompatError
	for {
		err := c.checkCompatible(newcfg, bhead)
		if err == nil || (lasterr != nil && err.RewindTo == lasterr.RewindTo) {
			break
		}
		lasterr = err
		bhead.SetUint64(err.RewindTo)
	}
	return lasterr
}

// CheckConfigForkOrder checks that we don't "skip" any forks, geth isn't pluggable enough
// to guarantee that forks can be implemented in a different order than on official networks
func (c *ChainConfig) CheckConfigForkOrder() error {
	// Multi-geth does not set or rely on these values in the same way
	// that ethereum/go-ethereum does; so multi-geth IS "pluggable" in this way -
	// (with chain configs broken into independent features, fork feature implementations
	// do not need to make assumptions about preceding forks or their associated
	// features).
	// This method was added in response to an issue, summarized here:
	// https://github.com/ethereum/go-ethereum/issues/20136#issuecomment-541895855
	return nil
}

func (c *ChainConfig) checkCompatible(newcfg *ChainConfig, head *big.Int) *ConfigCompatError {
	for _, ch := range []struct {
		name   string
		c1, c2 *big.Int
	}{
		{"Homestead", c.HomesteadBlock, newcfg.HomesteadBlock},
		{"EIP7F", c.EIP7FBlock, newcfg.EIP7FBlock},
		{"DAO", c.DAOForkBlock, newcfg.DAOForkBlock},
		{"EIP150", c.EIP150Block, newcfg.EIP150Block},
		{"EIP155", c.EIP155Block, newcfg.EIP155Block},
		{"EIP158", c.EIP158Block, newcfg.EIP158Block},
		{"EIP160F", c.EIP160FBlock, newcfg.EIP160FBlock},
		{"EIP161F", c.EIP161FBlock, newcfg.EIP161FBlock},
		{"EIP170F", c.EIP170FBlock, newcfg.EIP170FBlock},
		{"Byzantium", c.ByzantiumBlock, newcfg.ByzantiumBlock},
		{"EIP100F", c.EIP100FBlock, newcfg.EIP100FBlock},
		{"EIP140F", c.EIP140FBlock, newcfg.EIP140FBlock},
		{"EIP198F", c.EIP198FBlock, newcfg.EIP198FBlock},
		{"EIP211F", c.EIP211FBlock, newcfg.EIP211FBlock},
		{"EIP212F", c.EIP212FBlock, newcfg.EIP212FBlock},
		{"EIP213F", c.EIP213FBlock, newcfg.EIP213FBlock},
		{"EIP214F", c.EIP214FBlock, newcfg.EIP214FBlock},
		{"EIP649F", c.EIP649FBlock, newcfg.EIP649FBlock},
		{"EIP658F", c.EIP658FBlock, newcfg.EIP658FBlock},
		{"Constantinople", c.ConstantinopleBlock, newcfg.ConstantinopleBlock},
		{"EIP145F", c.EIP145FBlock, newcfg.EIP145FBlock},
		{"EIP1014F", c.EIP1014FBlock, newcfg.EIP1014FBlock},
		{"EIP1052F", c.EIP1052FBlock, newcfg.EIP1052FBlock},
		{"EIP1234F", c.EIP1234FBlock, newcfg.EIP1234FBlock},
		{"EIP1283F", c.EIP1283FBlock, newcfg.EIP1283FBlock},
		{"EWASM", c.EWASMBlock, newcfg.EWASMBlock},
	} {
		if err := func(c1, c2, head *big.Int) *ConfigCompatError {
			if isForkIncompatible(ch.c1, ch.c2, head) {
				return newCompatError(ch.name+" fork block", ch.c1, ch.c2)
			}
			return nil
		}(ch.c1, ch.c2, head); err != nil {
			return err
		}
	}

	if c.IsDAOFork(head) && c.DAOForkSupport != newcfg.DAOForkSupport {
		return newCompatError("DAO fork support flag", c.DAOForkBlock, newcfg.DAOForkBlock)
	}
	if c.IsEIP155(head) && !configNumEqual(c.ChainID, newcfg.ChainID) {
		return newCompatError("EIP155 chain ID", c.EIP155Block, newcfg.EIP155Block)
	}
	// Either Byzantium block must be set OR EIP100 and EIP649 must be equivalent
	if newcfg.ByzantiumBlock == nil {
		if !configNumEqual(newcfg.EIP100FBlock, newcfg.EIP649FBlock) {
			return newCompatError("EIP100F/EIP649F not equal", newcfg.EIP100FBlock, newcfg.EIP649FBlock)
		}
		if isForkIncompatible(c.EIP100FBlock, newcfg.EIP649FBlock, head) {
			return newCompatError("EIP100F/EIP649F fork block", c.EIP100FBlock, newcfg.EIP649FBlock)
		}
		if isForkIncompatible(c.EIP649FBlock, newcfg.EIP100FBlock, head) {
			return newCompatError("EIP649F/EIP100F fork block", c.EIP649FBlock, newcfg.EIP100FBlock)
		}
	}
	if isForkIncompatible(c.PetersburgBlock, newcfg.PetersburgBlock, head) {
		return newCompatError("ConstantinopleFix fork block", c.PetersburgBlock, newcfg.PetersburgBlock)
	}
	if isForkIncompatible(c.PetersburgBlock, newcfg.PetersburgBlock, head) {
		return newCompatError("Petersburg fork block", c.PetersburgBlock, newcfg.PetersburgBlock)
	}
	if isForkIncompatible(c.IstanbulBlock, newcfg.IstanbulBlock, head) {
		return newCompatError("Istanbul fork block", c.IstanbulBlock, newcfg.IstanbulBlock)
	}
	if isForkIncompatible(c.EWASMBlock, newcfg.EWASMBlock, head) {
		return newCompatError("ewasm fork block", c.EWASMBlock, newcfg.EWASMBlock)
	}

	return nil
}

// isForkIncompatible returns true if a fork scheduled at s1 cannot be rescheduled to
// block s2 because head is already past the fork.
func isForkIncompatible(s1, s2, head *big.Int) bool {
	return (IsForked(s1, head) || IsForked(s2, head)) && !configNumEqual(s1, s2)
}

// IsForked returns whether a fork scheduled at block s is active at the given head block.
func IsForked(s, head *big.Int) bool {
	if s == nil || head == nil {
		return false
	}
	return s.Cmp(head) <= 0
}

func configNumEqual(x, y *big.Int) bool {
	if x == nil {
		return y == nil
	}
	if y == nil {
		return x == nil
	}
	return x.Cmp(y) == 0
}

// ConfigCompatError is raised if the locally-stored blockchain is initialised with a
// ChainConfig that would alter the past.
type ConfigCompatError struct {
	What string
	// block numbers of the stored and new configurations
	StoredConfig, NewConfig *big.Int
	// the block number to which the local chain must be rewound to correct the error
	RewindTo uint64
}

func newCompatError(what string, storedblock, newblock *big.Int) *ConfigCompatError {
	var rew *big.Int
	switch {
	case storedblock == nil:
		rew = newblock
	case newblock == nil || storedblock.Cmp(newblock) < 0:
		rew = storedblock
	default:
		rew = newblock
	}
	err := &ConfigCompatError{what, storedblock, newblock, 0}
	if rew != nil && rew.Sign() > 0 {
		err.RewindTo = rew.Uint64() - 1
	}
	return err
}

func (err *ConfigCompatError) Error() string {
	return fmt.Sprintf("mismatching %s in database (have %d, want %d, rewindto %d)", err.What, err.StoredConfig, err.NewConfig, err.RewindTo)
}

// Rules wraps ChainConfig and is merely syntactic sugar or can be used for functions
// that do not have or require information about the block.
//
// Rules is a one time interface meaning that it shouldn't be used in between transition
// phases.
type Rules struct {
	ChainID                       *big.Int
	IsHomestead, IsEIP2F, IsEIP7F bool
	IsEIP150                      bool
	IsEIP155                      bool
	// EIP158HF - Tangerine Whistle
	IsEIP160F, IsEIP161F, IsEIP170F bool
	// Byzantium
	IsEIP100F, IsEIP140F, IsEIP198F, IsEIP211F, IsEIP212F, IsEIP213F, IsEIP214F, IsEIP649F, IsEIP658F bool
	// Constantinople
	IsEIP145F, IsEIP1014F, IsEIP1052F, IsEIP1283F, IsEIP1234F bool
	/// Istanbul
	IsEIP152F, IsEIP1108F, IsEIP1344F, IsEIP1884F, IsEIP2028F, IsEIP2200F bool
	IsPetersburg, IsIstanbul                                              bool
	IsBombDisposal, IsECIP1010, IsECIP1017F                               bool
	IsMCIP0, IsMCIP3, IsMCIP8                                             bool
}

// Rules ensures c's ChainID is not nil.
func (c *ChainConfig) Rules(num *big.Int) Rules {
	chainID := c.ChainID
	if chainID == nil {
		chainID = new(big.Int)
	}
	return Rules{
		ChainID: new(big.Int).Set(chainID),

		IsEIP2F: c.IsEIP2F(num),
		IsEIP7F: c.IsEIP7F(num),

		IsEIP150:  c.IsEIP150(num),
		IsEIP155:  c.IsEIP155(num),
		IsEIP160F: c.IsEIP160F(num),
		IsEIP161F: c.IsEIP161F(num),
		IsEIP170F: c.IsEIP170F(num),

		IsEIP100F: c.IsEIP100F(num),
		IsEIP140F: c.IsEIP140F(num),
		IsEIP198F: c.IsEIP198F(num),
		IsEIP211F: c.IsEIP211F(num),
		IsEIP212F: c.IsEIP212F(num),
		IsEIP213F: c.IsEIP213F(num),
		IsEIP214F: c.IsEIP214F(num),
		IsEIP649F: c.IsEIP649F(num),
		IsEIP658F: c.IsEIP658F(num),

		IsEIP145F:  c.IsEIP145F(num),
		IsEIP1014F: c.IsEIP1014F(num),
		IsEIP1052F: c.IsEIP1052F(num),
		IsEIP1234F: c.IsEIP1234F(num),
		IsEIP1283F: c.IsEIP1283F(num),

		IsEIP152F:  c.IsEIP152F(num),
		IsEIP1108F: c.IsEIP1108F(num),
		IsEIP1344F: c.IsEIP1344F(num),
		IsEIP1884F: c.IsEIP1884F(num),
		IsEIP2028F: c.IsEIP2028F(num),
		IsEIP2200F: c.IsEIP2200F(num),

		IsPetersburg: c.IsPetersburg(num),
		IsIstanbul:   c.IsIstanbul(num),

		IsBombDisposal: c.IsBombDisposal(num),
		IsECIP1010:     c.IsECIP1010(num),
		IsECIP1017F:    c.IsECIP1017F(num),

		IsMCIP0: c.IsMCIP0(num),
		IsMCIP3: c.IsMCIP3(num),
		IsMCIP8: c.IsMCIP8(num),
	}
}

// IsMCIP0 returns whether MCIP0 block is engaged; this is equivalent to 'IsMusicoin'.
// (There is no MCIP-0).
func (c *ChainConfig) IsMCIP0(num *big.Int) bool {
	return IsForked(c.MCIP0Block, num)
}

// IsMCIP3 returns whether MCIP3-UBI block is engaged.
func (c *ChainConfig) IsMCIP3(num *big.Int) bool {
	return IsForked(c.MCIP3Block, num)
}

// IsMCIP8 returns whether MCIP3-QT block is engaged.
func (c *ChainConfig) IsMCIP8(num *big.Int) bool {
	return IsForked(c.MCIP8Block, num)
}
