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


package paramtypes

import "math/big"

// OneOrAllEqOfBlocks returns <one> if it is not nil.
// If <allOf> are not equal, it returns nil.
// Otherwise, it returns the first value of <allOf>.
func OneOrAllEqOfBlocks(one *big.Int, allOf ...*big.Int) *big.Int {
	if one != nil {
		return one
	}
	for _, i := range allOf {
		if i == nil {
			return nil
		}
	}
	if len(allOf) == 0 {
		return nil
	}
	if len(allOf) == 1 {
		return allOf[0]
	}
	for i := 1; i < len(allOf); i++ {
		if allOf[i-1].Cmp(allOf[i]) != 0 {
			return nil
		}
	}
	return allOf[0]
}

func FeatureOrMetaBlock(featureBlock *big.Int, metaBlock *big.Int) *big.Int {
	if featureBlock != nil {
		return featureBlock
	}
	return metaBlock
}