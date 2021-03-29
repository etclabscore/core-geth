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
	"enode://534d18fd46c5cd5ba48a68250c47cea27a1376869755ed631c94b91386328039eb607cf10dd8d0aa173f5ec21e3fb45c5d7a7aa904f97bc2557e9cb4ccc703f1@51.158.190.99:30303", // @q9f lyrae
	"enode://15b6ae4e9e18772f297c90d83645b0fbdb56667ce2d747d6d575b21d7b60c2d3cd52b11dec24e418438caf80ddc433232b3685320ed5d0e768e3972596385bfc@51.158.191.43:41235", // @q9f mizar
	"enode://8fa15f5012ac3c47619147220b7772fcc5db0cb7fd132b5d196e7ccacb166ac1fcf83be1dace6cd288e288a85e032423b6e7e9e57f479fe7373edea045caa56b@176.9.51.216:31355",  // @q9f ceibo
	"enode://34c14141b79652afc334dcd2ba4d8047946246b2310dc8e45737ebe3e6f15f9279ca4702b90bc5be12929f6194e2c3ce19a837b7fec7ebffcee9e9fe4693b504@176.9.51.216:31365",  // @q9f ceibo
}

var MordorDNSNetwork1 = dnsPrefixETC + "all.mordor.blockd.info"
