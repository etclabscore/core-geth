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

// KottiBootnodes are the enode URLs of the P2P bootstrap nodes running on the
// Kotti test network.
var KottiBootnodes = []string{
	"enode://93b12383c74c39b67afa99a7ff44ce250fe94295fa1fc087465cc4fe2d0b33b91a8d8cabe03b250104a9096aa0e06bcde5f95665a5bd9f890edd2ab33e16ae47@51.15.41.19:30303", // @q9f ethstats
}

var KottiDNSNetwork1 = dnsPrefixETC + "all.kotti.blockd.info"
