// Copyright 2026 The core-geth Authors
// This file is part of the core-geth library.
//
// The core-geth library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The core-geth library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the core-geth library. If not, see <http://www.gnu.org/licenses/>.

package params

import (
	"math/big"

	"github.com/ethereum/go-ethereum/params/types/coregeth"
	"github.com/ethereum/go-ethereum/params/types/ctypes"
	"github.com/ethereum/go-ethereum/params/vars"
)

var (
	// greenpointECIP1049Transition activates ECIP-1049 (keccak256 PoW) at
	// genesis on the Greenpoint test network.
	greenpointECIP1049Transition uint64 = 0

	// GreenpointChainConfig contains the chain parameters for the Greenpoint
	// test network — a small, EIP/ECIP-loaded test chain that uses the
	// ECIP-1049 keccak256 proof-of-work engine from genesis.
	//
	// All historical Ethereum/ETC fork rules are squashed to block 0 so the
	// testnet behaves like a modern post-Mystique chain. ECIP-1049 supersedes
	// the legacy Ethash/Etchash sealing immediately, meaning blocks are mined
	// with plain keccak256(SealHash || nonce).
	GreenpointChainConfig = &coregeth.CoreGethChainConfig{
		NetworkID:                 7707,
		ChainID:                   big.NewInt(7707),
		SupportedProtocolVersions: vars.DefaultProtocolVersions,
		Ethash:                    new(ctypes.EthashConfig),

		// Squash all historical EIPs to genesis.
		EIP2FBlock:   big.NewInt(0),
		EIP7FBlock:   big.NewInt(0),
		EIP150Block:  big.NewInt(0),
		EIP155Block:  big.NewInt(0),
		EIP160FBlock: big.NewInt(0),
		EIP161FBlock: big.NewInt(0),
		EIP170FBlock: big.NewInt(0),

		// Byzantium eq.
		EIP100FBlock: big.NewInt(0),
		EIP140FBlock: big.NewInt(0),
		EIP198FBlock: big.NewInt(0),
		EIP211FBlock: big.NewInt(0),
		EIP212FBlock: big.NewInt(0),
		EIP213FBlock: big.NewInt(0),
		EIP214FBlock: big.NewInt(0),
		EIP658FBlock: big.NewInt(0),

		// Constantinople / Agharta eq.
		EIP145FBlock:  big.NewInt(0),
		EIP1014FBlock: big.NewInt(0),
		EIP1052FBlock: big.NewInt(0),

		// Istanbul / Phoenix eq.
		EIP152FBlock:  big.NewInt(0),
		EIP1108FBlock: big.NewInt(0),
		EIP1344FBlock: big.NewInt(0),
		EIP1884FBlock: big.NewInt(0),
		EIP2028FBlock: big.NewInt(0),
		EIP2200FBlock: big.NewInt(0),

		// Berlin / Magneto eq.
		EIP2565FBlock: big.NewInt(0),
		EIP2718FBlock: big.NewInt(0),
		EIP2929FBlock: big.NewInt(0),
		EIP2930FBlock: big.NewInt(0),

		// London (partial) / Mystique eq.
		EIP3529FBlock: big.NewInt(0),
		EIP3541FBlock: big.NewInt(0),

		// ECIP-1049: keccak256 proof-of-work from genesis.
		ECIP1049FBlock: big.NewInt(int64(greenpointECIP1049Transition)),

		// Disable ETC supply-emission curve / ECIP-1010 difficulty bomb logic;
		// Greenpoint mirrors a "modern" fork schedule with no bomb.
		DisposalBlock:      big.NewInt(0),
		ECIP1017FBlock:     big.NewInt(0),
		ECIP1017EraRounds:  big.NewInt(2_000_000),
		ECIP1010PauseBlock: nil,
		ECIP1010Length:     nil,
	}
)
