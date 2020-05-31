package params

import (
	"math/big"
	"testing"
)

func TestClassicDAO(t *testing.T) {
	blockNumbers := []*big.Int{
		big.NewInt(0),
		big.NewInt(1_920_000),
		big.NewInt(10_000_000),
	}
	for _, bn := range blockNumbers {
		if ClassicChainConfig.IsEnabled(ClassicChainConfig.GetEthashEIP779Transition, bn) {
			t.Fatal("bad")
		}
	}
}

