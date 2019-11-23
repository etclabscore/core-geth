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