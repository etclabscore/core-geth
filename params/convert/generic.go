package convert

import (
	"github.com/ethereum/go-ethereum/params/types/common"
	"github.com/ethereum/go-ethereum/params/types/goethereum"
	"github.com/ethereum/go-ethereum/params/types/parity"
)

type GenericCC struct {
	common.ChainConfigurator
}

func AsGenericCC(c common.ChainConfigurator) GenericCC {
	return GenericCC{c}
}

func (c GenericCC) DAOSupport() bool {
	if gc, ok := c.ChainConfigurator.(*goethereum.ChainConfig); ok {
		return gc.DAOForkSupport
	}
	if pc, ok := c.ChainConfigurator.(*parity.ParityChainSpec); ok {
		return pc.Engine.Ethash.Params.DaoHardforkBeneficiary != nil
	}
	return c.GetEthashEIP779Transition() != nil
}
