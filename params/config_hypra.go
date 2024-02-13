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
package params

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params/types/coregeth"
	"github.com/ethereum/go-ethereum/params/types/ctypes"
	"github.com/ethereum/go-ethereum/params/vars"
)

const HypraChainId = 622277

var (
	// HypraChainConfig is the chain parameters to run a node on the Hypra main network.
	HypraChainConfig = &coregeth.CoreGethChainConfig{
		NetworkID:                 HypraChainId,
		EthashB3:                  new(ctypes.EthashB3Config),
		ChainID:                   big.NewInt(HypraChainId),
		SupportedProtocolVersions: vars.DefaultProtocolVersions,

		EIP2FBlock: big.NewInt(0),
		EIP7FBlock: big.NewInt(0),
		// This is 0 on Hypra but Fork support is nil. Should this be enabled or not?
		// DAOForkSupport isnt in this struct
		DAOForkBlock: big.NewInt(0),
		EIP150Block:  big.NewInt(0),
		EIP155Block:  big.NewInt(0),
		EIP160FBlock: big.NewInt(0),
		EIP161FBlock: big.NewInt(0),
		EIP170FBlock: big.NewInt(0),

		// Byzantium eq -- Enables Smart contracts
		EIP100FBlock: big.NewInt(1001),
		EIP140FBlock: big.NewInt(1001),
		EIP198FBlock: big.NewInt(1001),
		EIP211FBlock: big.NewInt(1001),
		EIP212FBlock: big.NewInt(1001),
		EIP213FBlock: big.NewInt(1001),
		EIP214FBlock: big.NewInt(1001),
		EIP649FBlock: big.NewInt(1001),
		EIP658FBlock: big.NewInt(1001),

		// Constantinople eq
		EIP145FBlock:    big.NewInt(5503),
		EIP1014FBlock:   big.NewInt(5503),
		EIP1052FBlock:   big.NewInt(5503),
		EIP1283FBlock:   big.NewInt(5503),
		PetersburgBlock: big.NewInt(5507),

		// Istanbul eq
		EIP152FBlock:  big.NewInt(5519),
		EIP1108FBlock: big.NewInt(5519),
		EIP1344FBlock: big.NewInt(5519),
		EIP1884FBlock: big.NewInt(5519),
		EIP2028FBlock: big.NewInt(5519),
		EIP2200FBlock: big.NewInt(5519),
		EIP2384FBlock: big.NewInt(5521),

		// Berlin
		EIP2565FBlock: big.NewInt(5527), // ModExp Gas Cost
		EIP2718FBlock: big.NewInt(5527), // Typed Transaction Envelope
		EIP2929FBlock: big.NewInt(5527), // Gas cost increases for state access opcodes
		EIP2930FBlock: big.NewInt(5527), // Optional access lists

		// Veldin fork was used to enable rewards to miners for including uncle blocks on Hypra network.
		// Previously overlooked and unrewarded.
		HIPVeldinFBlock: big.NewInt(500_009),

		// London + shanghai EVM upgrade, aka Gaspar
		EIP3855FBlock: big.NewInt(1_600_957), // PUSH0 instruction
		EIP3860FBlock: big.NewInt(1_600_957), // Limit and meter initcode
		EIP3198FBlock: big.NewInt(1_600_957), // BASEFEE Opcode

		// London chain upgrades, aka Planned Eudora
		// TODO: move block numbers closer once testing has concluded
		EIP1559FBlock: big.NewInt(13_524_557), // EIP-1559 transactions`
		EIP3529FBlock: big.NewInt(13_524_557), // Reduction in refunds

		// Unplanned Upgrade, aka Olantis
		// EIP3651FBlock: big.NewInt(13_524_557), // Warm COINBASE (gas reprice)
		// EIP6049FBlock: big.NewInt(13_524_557), // Deprecate SELFDESTRUCT
		// EIP3541FBlock: big.NewInt(13_524_557), // Reject new contract code starting with the 0xEF byte

		// Spiral, aka Shanghai (partially)
		// EIP4399FBlock: nil, // Supplant DIFFICULTY with PREVRANDAO. Hypra  does not spec 4399 because it's still PoW, and 4399 is only applicable for the PoS system.
		// EIP4895FBlock: nil, // Beacon chain push withdrawals as operations

		// Dummy EIPs, unused by ethashb3
		//EIP3554FBlock: big.NewInt(13_524_557),
		//EIP3541FBlock: big.NewInt(13_524_557),
		//EIP4345FBlock: big.NewInt(27_200_177),
		//EIP5133FBlock: big.NewInt(40_725_107),
		// Define the planned 3 year decreasing rewards.
		BlockRewardSchedule: map[uint64]*big.Int{
			0:          big.NewInt(4 * vars.Ether),
			13_524_557: big.NewInt(3 * vars.Ether),
			27_200_177: big.NewInt(2 * vars.Ether),
			40_725_107: big.NewInt(1 * vars.Ether),
		},

		TrustedCheckpoint: &ctypes.TrustedCheckpoint{
			BloomRoot:    common.HexToHash("0xb09a3dc1215903dc757dd83c960c6dd77ec508e03c185b16df62af069d197fc3"),
			CHTRoot:      common.HexToHash("0x21f6a69cbd097bd2c1d250ed7144888e483d46105b66a99a68e4847c99a863cd"),
			SectionHead:  common.HexToHash("0x2461cbf15f92264332e09ba8a9d6b9106c12b3cf0e7d7bfff78483c1b271e902"),
			SectionIndex: 43,
		},

		RequireBlockHashes: map[uint64]common.Hash{
			156_000:   common.HexToHash("0x2a27bec023108c5f650cb0c9b7aaae7e7fdeefbeb3cd14a8390bb0957043aca2"),
			512_000:   common.HexToHash("0xd547c896967174295c61ea18534d3c69337a6fa6e957909dc29e4bcba873592a"),
			1_878_299: common.HexToHash("0x1834091726bd0890539db04a22b673a6d9262d8b350dffaefdb10abc72e2113d"),
		},
	}
)

func init() {
	// hack to override the default bomb schedule
	eip2384 := HypraChainConfig.EIP2384FBlock.Uint64()
	HypraChainConfig.SetEthashEIP2384Transition(&eip2384)
}
