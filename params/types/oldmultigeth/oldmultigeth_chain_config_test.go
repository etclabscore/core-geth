package oldmultigeth

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params/types/ctypes"
)

func TestOldMultigethIsEIP7(t *testing.T) {
	// ClassicChainConfig is the chain parameters to run a node on the Classic main network.
	c := &ChainConfig{
		ChainID:              big.NewInt(61),
		HomesteadBlock:       big.NewInt(1150000),
		DAOForkBlock:         big.NewInt(1920000),
		DAOForkSupport:       false,
		EIP150Block:          big.NewInt(2500000),
		EIP150Hash:           common.HexToHash("0xca12c63534f565899681965528d536c52cb05b7c48e269c2a6cb77ad864d878a"),
		EIP155Block:          big.NewInt(3000000),
		EIP158Block:          big.NewInt(8772000),
		ByzantiumBlock:       big.NewInt(8772000),
		DisposalBlock:        big.NewInt(5900000),
		SocialBlock:          nil,
		EthersocialBlock:     nil,
		ConstantinopleBlock:  big.NewInt(9573000),
		PetersburgBlock:      big.NewInt(9573000),
		IstanbulBlock:        big.NewInt(10500839),
		EIP1884DisableFBlock: big.NewInt(10500839),
		ECIP1017EraRounds:    big.NewInt(5000000),
		EIP160FBlock:         big.NewInt(3000000),
		ECIP1010PauseBlock:   big.NewInt(3000000),
		ECIP1010Length:       big.NewInt(2000000),
		Ethash:               new(ctypes.EthashConfig),
	}

	if got := c.GetEIP7Transition(); got == nil || *got != 1150000 {
		t.Fatal("bad", *got)
	}
	if got := c.GetEIP150Transition(); got == nil || *got != 2500000 {
		t.Fatal("bad", *got)
	}
	if got := c.GetEIP160Transition(); got == nil || *got != 3000000 {
		t.Fatal("bad", *got)
	}
	if got := c.GetEIP155Transition(); got == nil || *got != 3000000 {
		t.Fatal("bad", *got)
	}
	if got := c.GetEIP170Transition(); got == nil || *got != 8772000 {
		t.Fatal("bad", *got)
	}
}