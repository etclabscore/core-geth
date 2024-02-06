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

// MordorBootnodes are the enode URLs of the P2P bootstrap nodes running on
// the Ethereum Classic Mordor network.
// https://github.com/etclabscore/mordor/blob/master/static-nodes.json
var MordorBootnodes = []string{
	"enode://4539a067ae1f6a7ffac509603ba37baf772fc832880ddc67c53f292b6199fb048267f0311c820bc90bfd39ec663bc6b5256bdf787ec38425c82bde6bc2bcfe3c@24.199.107.164:30303", // @etccoop-sfo
}
var MordorDNSNetwork1 = dnsPrefixETC + "all.mordor.blockd.info"
