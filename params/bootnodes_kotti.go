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
	"enode://06333009fc9ef3c9e174768e495722a7f98fe7afd4660542e983005f85e556028410fd03278944f44cfe5437b1750b5e6bd1738f700fe7da3626d52010d2954c@51.141.15.254:30303", // @paritytech authority
	"enode://ae8658da8d255d1992c3ec6e62e11d6e1c5899aa1566504bc1ff96a0c9c8bd44838372be643342553817f5cc7d78f1c83a8093dee13d77b3b0a583c050c81940@18.232.185.151:30303",
	"enode://67913271d14f445689e8310270c304d42f268428f2de7a4ac0275bea97690e021df6f549f462503ff4c7a81d9dd27288867bbfa2271477d0911378b8944fae55@157.230.239.163:30303",
	"enode://e8a786a894db053fe6886e283fc4385389ad034e04a692a26335f30b714059efd5cead0e410ecd783ce095888fdafcc21a685f13501594e969d6f5ac7ba0388c@86.103.236.55:63384",
	"enode://4956f6924335c951cb70cbc169a85c081f6ff0b374aa2815453b8a3132b49613f38a1a6b8e103f878dbec86364f60091e92a376d7cd3aca9d82d2f2554794e63@51.15.97.240:41235",  // @q9f core-geth ginan
	"enode://6c9a052c01bb9995fa53bebfcdbc17733fe90708270d0e6d8e38dc57b32e1dbe8c287590b634ee9753b94ba302f411c96519c7fa07df0df6a6848149d819b2c5@51.15.70.7:41235",    // @q9f core-geth polis
	"enode://95a7302fd8f35d9ad716a591b90dfe839dbf2d3d6a6737ef39e07899f844ad82f1659dd6212492922dd36705fb0a1e984c1d5d5c42069d5bd329388831e820c1@51.15.97.240:45678",  // @q9f parity ginan
	"enode://8c5c4dec9a0728c7058d3c29698bae888adc244e443cebc21f3d708a20751511acbf98a52b5e5ea472f8715c658913e8084123461fd187a4980a0600420b0791@51.15.70.7:45678",    // @q9f parity polis
	"enode://efd7391a3bed73ad74ae5760319bb48f9c9f1983ff22964422688cdb426c5d681004ece26c47121396653cf9bafe7104aa4ecff70e24cc5b11fd76be8e5afce0@51.158.191.43:45678", // @q9f besu mizar
}

var KottiDNSNetwork1 = dnsPrefixETC + "all.kotti.blockd.info"
