package multigeth

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params/types/ctypes"
)

// ChainConfig is the core config which determines the blockchain settings.
//
// ChainConfig is stored in the database on a per block basis. This means
// that any network, identified by its genesis block, can have its own
// set of configuration options.
type ChainConfig struct {
	NetworkID                 uint64   `json:"-"`
	ChainID                   *big.Int `json:"chainId"`                             // chainId identifies the current chain and is used for replay protection
	SupportedProtocolVersions []uint   `json:"supportedProtocolVersions,omitempty"` // supportedProtocolVersions identifies the supported eth protocol versions for the current chain

	HomesteadBlock *big.Int `json:"homesteadBlock,omitempty"` // Homestead switch block (nil = no fork, 0 = already homestead)

	DAOForkBlock   *big.Int `json:"daoForkBlock,omitempty"`   // TheDAO hard-fork switch block (nil = no fork)
	DAOForkSupport bool     `json:"daoForkSupport,omitempty"` // Whether the nodes supports or opposes the DAO hard-fork

	// EIP150 implements the Gas price changes (https://github.com/ethereum/EIPs/issues/150)
	EIP150Block *big.Int    `json:"eip150Block,omitempty"` // EIP150 HF block (nil = no fork)
	EIP150Hash  common.Hash `json:"eip150Hash,omitempty"`  // EIP150 HF hash (needed for header only clients as only gas pricing changed)

	EIP155Block *big.Int `json:"eip155Block,omitempty"` // EIP155 HF block
	EIP158Block *big.Int `json:"eip158Block,omitempty"` // EIP158 HF block

	ByzantiumBlock      *big.Int `json:"byzantiumBlock,omitempty"`      // Byzantium switch block (nil = no fork, 0 = already on byzantium)
	ConstantinopleBlock *big.Int `json:"constantinopleBlock,omitempty"` // Constantinople switch block (nil = no fork, 0 = already activated)
	PetersburgBlock     *big.Int `json:"petersburgBlock,omitempty"`     // Petersburg switch block (nil = same as Constantinople)
	IstanbulBlock       *big.Int `json:"istanbulBlock,omitempty"`       // Istanbul switch block (nil = no fork, 0 = already on istanbul)
	MuirGlacierBlock    *big.Int `json:"muirGlacierBlock,omitempty"`    // Eip-2384 (bomb delay) switch block (nil = no fork, 0 = already activated)

	BerlinBlock       *big.Int `json:"berlinBlock,omitempty"`       // Berlin switch block
	LondonBlock       *big.Int `json:"londonBlock,omitempty"`       // London switch block
	ArrowGlacierBlock *big.Int `json:"arrowGlacierBlock,omitempty"` // ArrowGlacier switch block

	MergeForkBlock *big.Int `json:"mergeForkBlock,omitempty"` // EIP-3675 (TheMerge) switch block (nil = no fork, 0 = already in merge proceedings)

	EWASMBlock *big.Int `json:"ewasmBlock,omitempty"` // EWASM switch block (nil = no fork, 0 = already activated)

	// TerminalTotalDifficulty is the amount of total difficulty reached by
	// the network that triggers the consensus upgrade.
	TerminalTotalDifficulty *big.Int `json:"terminalTotalDifficulty,omitempty"`

	//
	// EXP cost increase
	// https://github.com/ethereum/EIPs/blob/master/EIPS/eip-160.md
	// NOTE: this json tag:
	// (a.) varies from it's 'siblings', which have 'F's in them
	// (b.) without the 'F' will vary from ETH implementations if they choose to accept the proposed changes
	// with corresponding refactoring (https://github.com/ethereum/go-ethereum/pull/18401)
	EIP160Block         *big.Int `json:"eip160Block,omitempty"`
	EIP161DisableBlock  *big.Int `json:"eip161DisableBlock,omitempty"`
	EIP161ReenableBlock *big.Int `json:"eip161ReenableBlock,omitempty"`
	ECIP1010PauseBlock  *big.Int `json:"ecip1010PauseBlock,omitempty"` // ECIP1010 pause HF block
	ECIP1010Length      *big.Int `json:"ecip1010Length,omitempty"`     // ECIP1010 length
	ECIP1017EraBlock    *big.Int `json:"ecip1017EraBlock,omitempty"`   // ECIP1017 era rounds
	DisposalBlock       *big.Int `json:"disposalBlock,omitempty"`      // Bomb disposal HF block

	MCIP0Block *big.Int `json:"mcip0Block,omitempty"` // Musicoin default block; no MCIP, just denotes chain pref
	MCIP3Block *big.Int `json:"mcip3Block,omitempty"` // Musicoin 'UBI Fork' block
	MCIP8Block *big.Int `json:"mcip8Block,omitempty"` // Musicoin 'QT For' block

	// Various consensus engines
	Ethash *ctypes.EthashConfig `json:"ethash,omitempty"`
	Clique *ctypes.CliqueConfig `json:"clique,omitempty"`
	Lyra2  *ctypes.Lyra2Config  `json:"lyra2,omitempty"`

	Lyra2NonceTransitionBlock *big.Int `json:"lyra2NonceTransitionBlock,omitempty"`
}
