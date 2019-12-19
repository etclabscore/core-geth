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


package tconvert

// FIXME(meowsbits): Requires implementing ChainConfigurator interface for Aleth data type.
//// Tests the go-ethereum to Aleth chainspec conversion for the Stureby testnet.
//func TestAlethSturebyConverter(t *testing.T) {
//	blob, err := ioutil.ReadFile("testdata/stureby_geth.json")
//	if err != nil {
//		t.Fatalf("could not read file: %v", err)
//	}
//	var genesis paramtypes.Genesis
//	if err := json.Unmarshal(blob, &genesis); err != nil {
//		t.Fatalf("failed parsing genesis: %v", err)
//	}
//	spec, err := NewAlethGenesisSpec("stureby", &genesis)
//	if err != nil {
//		t.Fatalf("failed creating chainspec: %v", err)
//	}
//
//	expBlob, err := ioutil.ReadFile("testdata/stureby_aleth.json")
//	if err != nil {
//		t.Fatalf("could not read file: %v", err)
//	}
//	expspec := &aleth.AlethGenesisSpec{}
//	if err := json.Unmarshal(expBlob, expspec); err != nil {
//		t.Fatalf("failed parsing genesis: %v", err)
//	}
//	// Compare the read-in SomeSpec vs. the generated SomeSpec
//	if diffs := deep.Equal(expspec, spec); len(diffs) != 0 {
//		t.Errorf("chainspec mismatch")
//		for _, d := range diffs {
//			t.Log(d)
//		}
//		t.Log(spew.Sprint(spec))
//		bm, _ := json.MarshalIndent(spec, "", "    ")
//		t.Log(string(bm))
//	}
//}
