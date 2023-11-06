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

// ClassicBootnodes are the enode URLs of the P2P bootstrap nodes running on
// the Ethereum Classic network.
var ClassicBootnodes = []string{

	"enode://6b6ea53a498f0895c10269a3a74b777286bd467de6425c3b512740fcc7fbc8cd281dca4ab041dd97d62b38f3d0b5b05e71f48d28a3a2f4b5de40fe1f6bf05531@157.245.77.211:30303", // AMS
	"enode://16264d48df59c3492972d96bf8a39dd38bab165809a3a4bb161859a337de38b2959cc98efea94355c7a7177cd020867c683aed934dbd6bc937d9e6b61d94d8d9@64.225.0.245:30303",   // NYC

	"enode://55bbc7f0ffa2af2ceca997ec195a98768144a163d389ae87b808dff8a861618405c2582451bbb6022e429e4bcd6b0e895e86160db6e93cdadbcfd80faacf6f06@164.90.144.106:30303", // SFO
}

var dnsPrefixETC = "enrtree://AJE62Q4DUX4QMMXEHCSSCSC65TDHZYSMONSD64P3WULVLSF6MRQ3K@"

var ClassicDNSNetwork1 = dnsPrefixETC + "all.classic.blockd.info"
