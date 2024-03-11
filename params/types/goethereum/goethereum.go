// Copyright 2019 The multi-geth Authors
// This file is part of the multi-geth library.
//
// The multi-geth library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The multi-geth library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the multi-geth library. If not, see <http://www.gnu.org/licenses/>.

package goethereum

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params/confp"
	"github.com/ethereum/go-ethereum/params/types/ctypes"
)

type ChainConfig struct {
	NetworkID                 uint64   `json:"-"`
	ChainID                   *big.Int `json:"chainId"`                             // chainId identifies the current chain and is used for replay protection
	SupportedProtocolVersions []uint   `json:"supportedProtocolVersions,omitempty"` // supportedProtocolVersions identifies the supported eth protocol versions for the current chain

	// HF: Homestead
	HomesteadBlock *big.Int `json:"homesteadBlock,omitempty"` // Homestead switch block (nil = no fork, 0 = already homestead)

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

	// HF: Byzantium
	ByzantiumBlock *big.Int `json:"byzantiumBlock,omitempty"` // Byzantium switch block (nil = no fork, 0 = already on byzantium)

	// HF: Constantinople
	ConstantinopleBlock *big.Int `json:"constantinopleBlock,omitempty"` // Constantinople switch block (nil = no fork, 0 = already activated)

	PetersburgBlock *big.Int `json:"petersburgBlock,omitempty"` // Petersburg switch block (nil = same as Constantinople)

	// HF: Istanbul
	IstanbulBlock    *big.Int `json:"istanbulBlock,omitempty"`    // Istanbul switch block (nil = no fork, 0 = already on istanbul)
	MuirGlacierBlock *big.Int `json:"muirGlacierBlock,omitempty"` // Eip-2384 (bomb delay) switch block (nil = no fork, 0 = already activated)

	BerlinBlock        *big.Int `json:"berlinBlock,omitempty"`        // Berlin switch block
	LondonBlock        *big.Int `json:"londonBlock,omitempty"`        // London switch block
	ArrowGlacierBlock  *big.Int `json:"arrowGlacierBlock,omitempty"`  // ArrowGlacier switch block
	GrayGlacierBlock   *big.Int `json:"grayGlacierBlock,omitempty"`   // Eip-5133 (bomb delay) switch block (nil = no fork, 0 = already activated)
	MergeNetsplitBlock *big.Int `json:"mergeNetsplitBlock,omitempty"` // Virtual fork after The Merge to use as a network splitter

	ShanghaiBlock *big.Int `json:"shanghaiBlock,omitempty"` // Shanghai switch block (nil = no fork, 0 = already on shanghai)
	CancunBlock   *big.Int `json:"cancunBlock,omitempty"`   // Cancun switch block (nil = no fork, 0 = already on cancun)

	// Fork scheduling was switched from blocks to timestamps here

	ShanghaiTime *uint64 `json:"shanghaiTime,omitempty"` // Shanghai switch time (nil = no fork, 0 = already on shanghai)
	CancunTime   *uint64 `json:"cancunTime,omitempty"`   // Cancun switch time (nil = no fork, 0 = already on cancun)
	PragueTime   *uint64 `json:"pragueTime,omitempty"`   // Prague switch time (nil = no fork, 0 = already on prague)
	VerkleTime   *uint64 `json:"verkleTime,omitempty"`   // Verkle switch time (nil = no fork, 0 = already on verkle)

	// TerminalTotalDifficulty is the amount of total difficulty reached by
	// the network that triggers the consensus upgrade.
	TerminalTotalDifficulty *big.Int `json:"terminalTotalDifficulty,omitempty"`

	// TerminalTotalDifficultyPassed is a flag specifying that the network already
	// passed the terminal total difficulty. Its purpose is to disable legacy sync
	// even without having seen the TTD locally (safer long term).
	TerminalTotalDifficultyPassed bool `json:"terminalTotalDifficultyPassed,omitempty"`

	EWASMBlock *big.Int `json:"ewasmBlock,omitempty"` // EWASM switch block (nil = no fork, 0 = already activated)

	// Various consensus engines
	Ethash    *ctypes.EthashConfig `json:"ethash,omitempty"`
	Clique    *ctypes.CliqueConfig `json:"clique,omitempty"`
	Lyra2     *ctypes.Lyra2Config  `json:"lyra2,omitempty"`
	IsDevMode bool                 `json:"isDev,omitempty"`

	// NOTE: These are not included in this type upstream.
	TrustedCheckpoint       *ctypes.TrustedCheckpoint      `json:"trustedCheckpoint"`
	TrustedCheckpointOracle *ctypes.CheckpointOracleConfig `json:"trustedCheckpointOracle"`

	EIP1706Transition  *big.Int `json:"-"`
	ECIP1080Transition *big.Int `json:"-"`

	// Cache types for use with testing, but will not show up in config API.
	ecbp1100Transition           *big.Int
	ecbp1100DeactivateTransition *big.Int

	Lyra2NonceTransitionBlock *big.Int `json:"lyra2NonceTransitionBlock,omitempty"`
}

// networkNames are user friendly names to use in the chain spec banner.
// The key values are Chain IDs. These are not imported to avoid circular dependency issues and because
// these values are expected to be eternally constant.
var networkNames = map[string]string{
	"1":        "mainnet",
	"3":        "ropsten",
	"4":        "rinkeby",
	"5":        "goerli",
	"11155111": "sepolia",
}

// String implements the fmt.Stringer interface.
func (c *ChainConfig) String() string {
	var banner string

	// Create some basinc network config output
	network := networkNames[c.ChainID.String()]
	if network == "" {
		network = "unknown"
	}
	banner += fmt.Sprintf("Chain ID:  %v (%s)\n", c.ChainID, network)
	switch {
	case c.Ethash != nil:
		if c.TerminalTotalDifficulty == nil {
			banner += "Consensus: Ethash (proof-of-work)\n"
		} else {
			banner += "Consensus: Beacon (proof-of-stake), merged from Ethash (proof-of-work)\n"
		}
	case c.Clique != nil:
		if c.TerminalTotalDifficulty == nil {
			banner += "Consensus: Clique (proof-of-authority)\n"
		} else {
			banner += "Consensus: Beacon (proof-of-stake), merged from Clique (proof-of-authority)\n"
		}
	default:
		banner += "Consensus: unknown\n"
	}
	banner += "\n"
	banner += fmt.Sprintf(`_ Block-based Forks: %v`, confp.BlockForks(c))
	banner += fmt.Sprintf(`_ Time-based Forks: %v`, confp.TimeForks(c, 0))
	banner += fmt.Sprintf(`_ TTD: %v`, c.GetEthashTerminalTotalDifficulty())

	// // Create a list of forks with a short description of them. Forks that only
	// // makes sense for mainnet should be optional at printing to avoid bloating
	// // the output for testnets and private networks.
	// banner += "Pre-Merge hard forks:\n"
	// banner += fmt.Sprintf(" - Homestead:                   %-8v (https://github.com/ethereum/execution-specs/blob/master/network-upgrades/mainnet-upgrades/homestead.md)\n", c.HomesteadBlock)
	// if c.DAOForkBlock != nil {
	// 	banner += fmt.Sprintf(" - DAO Fork:                    %-8v (https://github.com/ethereum/execution-specs/blob/master/network-upgrades/mainnet-upgrades/dao-fork.md)\n", c.DAOForkBlock)
	// }
	// banner += fmt.Sprintf(" - Tangerine Whistle (EIP 150): %-8v (https://github.com/ethereum/execution-specs/blob/master/network-upgrades/mainnet-upgrades/tangerine-whistle.md)\n", c.EIP150Block)
	// banner += fmt.Sprintf(" - Spurious Dragon/1 (EIP 155): %-8v (https://github.com/ethereum/execution-specs/blob/master/network-upgrades/mainnet-upgrades/spurious-dragon.md)\n", c.EIP155Block)
	// banner += fmt.Sprintf(" - Spurious Dragon/2 (EIP 158): %-8v (https://github.com/ethereum/execution-specs/blob/master/network-upgrades/mainnet-upgrades/spurious-dragon.md)\n", c.EIP155Block)
	// banner += fmt.Sprintf(" - Byzantium:                   %-8v (https://github.com/ethereum/execution-specs/blob/master/network-upgrades/mainnet-upgrades/byzantium.md)\n", c.ByzantiumBlock)
	// banner += fmt.Sprintf(" - Constantinople:              %-8v (https://github.com/ethereum/execution-specs/blob/master/network-upgrades/mainnet-upgrades/constantinople.md)\n", c.ConstantinopleBlock)
	// banner += fmt.Sprintf(" - Petersburg:                  %-8v (https://github.com/ethereum/execution-specs/blob/master/network-upgrades/mainnet-upgrades/petersburg.md)\n", c.PetersburgBlock)
	// banner += fmt.Sprintf(" - Istanbul:                    %-8v (https://github.com/ethereum/execution-specs/blob/master/network-upgrades/mainnet-upgrades/istanbul.md)\n", c.IstanbulBlock)
	// if c.MuirGlacierBlock != nil {
	// 	banner += fmt.Sprintf(" - Muir Glacier:                %-8v (https://github.com/ethereum/execution-specs/blob/master/network-upgrades/mainnet-upgrades/muir-glacier.md)\n", c.MuirGlacierBlock)
	// }
	// banner += fmt.Sprintf(" - Berlin:                      %-8v (https://github.com/ethereum/execution-specs/blob/master/network-upgrades/mainnet-upgrades/berlin.md)\n", c.BerlinBlock)
	// banner += fmt.Sprintf(" - London:                      %-8v (https://github.com/ethereum/execution-specs/blob/master/network-upgrades/mainnet-upgrades/london.md)\n", c.LondonBlock)
	// if c.ArrowGlacierBlock != nil {
	// 	banner += fmt.Sprintf(" - Arrow Glacier:               %-8v (https://github.com/ethereum/execution-specs/blob/master/network-upgrades/mainnet-upgrades/arrow-glacier.md)\n", c.ArrowGlacierBlock)
	// }
	// if c.GrayGlacierBlock != nil {
	// 	banner += fmt.Sprintf(" - Gray Glacier:                %-8v (https://github.com/ethereum/execution-specs/blob/master/network-upgrades/mainnet-upgrades/gray-glacier.md)\n", c.GrayGlacierBlock)
	// }
	// banner += "\n"
	//
	// // Add a special section for the merge as it's non-obvious
	// if c.TerminalTotalDifficulty == nil {
	// 	banner += "Merge not configured!\n"
	// 	banner += " - Hard-fork specification: https://github.com/ethereum/execution-specs/blob/master/network-upgrades/mainnet-upgrades/paris.md)"
	// } else {
	// 	banner += "Merge configured:\n"
	// 	banner += " - Hard-fork specification:   https://github.com/ethereum/execution-specs/blob/master/network-upgrades/mainnet-upgrades/paris.md)\n"
	// 	banner += fmt.Sprintf(" - Total terminal difficulty: %v\n", c.TerminalTotalDifficulty)
	// 	banner += fmt.Sprintf(" - Merge netsplit block:      %-8v", c.MergeNetsplitBlock)
	// }
	return banner
}
