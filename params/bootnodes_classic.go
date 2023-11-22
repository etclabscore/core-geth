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

	"enode://2b1ef75e8b7119b6e0294f2e51ead2cf1a5400472452c199e9587727ada99e7e2b1199e36adcad6cbae65dce2410559546e4d83d8c93d45a559e723e56444c03@67.207.93.100:30303",
	"enode://8e73168affd8d445edda09c561d607081ca5d7963317caae2702f701eb6546b06948b7f8687a795de576f6a5f33c44828e25a90aa63de18db380a11e660dd06f@159.203.37.80:30303",
	"enode://40b37d729fdc2c29620bbecc61264f450c3e849b5e60d1a362e6af04d84ae36e15120573498bcf56f97aeec18e530113a4d0b55b18e6a6523de8cf5a45e1db23@128.199.160.131:30303",
	"enode://84e70a981a436efd7d6ae80b5d7ecb82636de77029635a2fb25ccf5d02106172d1cafdd7d8171a136d3fd545e5df94f4d63299eda7f8aec79967aa568ddaa71e@165.22.202.207:30303",
	"enode://6ec7ac618a7147d8b0e41bc1e63080abfd56a48f50f06095c112e30f72cd5eee262b06aa168cb46ab470d56885f842a1ae44a3714bfb029ced83dab461852062@164.90.218.200:30303",
	"enode://feba6c4bd757efce4fec0fa5775780d2c209e86231553a72738c2c4e95d0cfd7ac2189f6d30d3c059dd25d93ea060ca28f9936b564f7d04c4270bb60e8e1487b@178.128.171.230:30303",
	"enode://c0b698f9fb1c5dd3b348c46dffbc524cebbfa48fd03c52523892d43d31d5fabacf4e27f215f5647c96190567975e423a5708cbdcd6441da2705bd85533e78eca@209.97.136.89:30303",
	"enode://37b20dc1ca0ba428dc07f4bf87c548aa57b53e6883c607a041c88cc0f72c27ef3487b0e94f50e787fff8af5e2baf029cd6105461aa8ec250ccae93d6972107a9@134.209.30.110:30303",
}

var dnsPrefixETC = "enrtree://AJE62Q4DUX4QMMXEHCSSCSC65TDHZYSMONSD64P3WULVLSF6MRQ3K@"

var ClassicDNSNetwork1 = dnsPrefixETC + "all.classic.blockd.info"
