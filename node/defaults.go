// Copyright 2016 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package node

import (
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/nat"
	"github.com/ethereum/go-ethereum/params/vars"
	"github.com/ethereum/go-ethereum/rpc"
)

const (
	DefaultHTTPHost = "localhost" // Default host interface for the HTTP RPC server
	DefaultHTTPPort = 8545        // Default TCP port for the HTTP RPC server
	DefaultWSHost   = "localhost" // Default host interface for the websocket RPC server
	DefaultWSPort   = 8546        // Default TCP port for the websocket RPC server
	DefaultAuthHost = "localhost" // Default host interface for the authenticated apis
	DefaultAuthPort = 8551        // Default port for the authenticated apis
)

const (
	// Engine API batch limits: these are not configurable by users, and should cover the
	// needs of all CLs.
	engineAPIBatchItemLimit         = 2000
	engineAPIBatchResponseSizeLimit = 250 * 1000 * 1000
)

var (
	DefaultAuthCors    = []string{"localhost"} // Default cors domain for the authenticated apis
	DefaultAuthVhosts  = []string{"localhost"} // Default virtual hosts for the authenticated apis
	DefaultAuthOrigins = []string{"localhost"} // Default origins for the authenticated apis
	DefaultAuthPrefix  = ""                    // Default prefix for the authenticated apis
	DefaultAuthModules = []string{"eth", "engine"}
)

// DefaultConfig contains reasonable default settings.
var DefaultConfig = Config{
	DataDir:              vars.DefaultDataDir(),
	HTTPPort:             DefaultHTTPPort,
	AuthAddr:             DefaultAuthHost,
	AuthPort:             DefaultAuthPort,
	AuthVirtualHosts:     DefaultAuthVhosts,
	HTTPModules:          []string{"net", "web3"},
	HTTPVirtualHosts:     []string{"localhost"},
	HTTPTimeouts:         rpc.DefaultHTTPTimeouts,
	WSPort:               DefaultWSPort,
	WSModules:            []string{"net", "web3"},
	BatchRequestLimit:    1000,
	BatchResponseMaxSize: 25 * 1000 * 1000,
	GraphQLVirtualHosts:  []string{"localhost"},
	P2P: p2p.Config{
		ListenAddr: ":30303",
		MaxPeers:   50,
		NAT:        nat.Any(),
	},
	DBEngine: "", // Use whatever exists, will default to Pebble if non-existent and supported
}
