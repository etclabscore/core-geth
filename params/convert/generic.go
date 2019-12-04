package convert

import (
	"fmt"

	paramtypes "github.com/ethereum/go-ethereum/params/types"
	"github.com/ethereum/go-ethereum/params/types/common"
	"github.com/ethereum/go-ethereum/params/types/goethereum"
	"github.com/ethereum/go-ethereum/params/types/parity"
	"github.com/ethereum/go-ethereum/params/vars"
)

// GenericCC is a generic-y struct type used to expose some meta-logic methods
// shared by all ChainConfigurator implementations but not existing in that interface.
// These logics differentiate from the logics present in the ChainConfigurator interface
// itself because they are chain-aware, or fit nuanced, or adhoc, use cases and should
// not be demanded of EVM-based ecosystem logic as a whole. Debatable. NOTE.
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
	if mg, ok := c.ChainConfigurator.(*paramtypes.MultiGethChainConfig); ok {
		return mg.GetEthashEIP779Transition() != nil
	}
	if pc, ok := c.ChainConfigurator.(*parity.ParityChainSpec); ok {
		return pc.Engine.Ethash.Params.DaoHardforkTransition != nil &&
			pc.Engine.Ethash.Params.DaoHardforkBeneficiary != nil &&
			*pc.Engine.Ethash.Params.DaoHardforkBeneficiary == vars.DAORefundContract &&
			len(pc.Engine.Ethash.Params.DaoHardforkAccounts) == len(vars.DAODrainList())
	}
	panic(fmt.Sprintf("uimplemented DAO logic, config: %v", c.ChainConfigurator))
}
