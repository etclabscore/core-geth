package convert

import (
	"strings"

	"github.com/ethereum/go-ethereum/params/types"
	"github.com/ethereum/go-ethereum/params/types/parity"
)

// NewParityChainSpec converts a go-ethereum genesis block into a Parity specific
// chain specification format.
func NewParityChainSpec(network string, genesis *paramtypes.Genesis, bootnodes []string) (*parity.ParityChainSpec, error) {
	spec := &parity.ParityChainSpec{
		Name:    network,
		Nodes:   bootnodes,
		Datadir: strings.ToLower(network),
	}
	if err := Convert(genesis, spec); err != nil {
		return nil, err
	}
	if err := Convert(genesis.Config, spec); err != nil {
		return nil, err
	}
	return spec, nil
}

// ToMultiGethGenesis converts a Parity chainspec to the corresponding MultiGeth datastructure.
// Note that the return value 'core.Genesis' includes the respective 'params.MultiGethChainConfig' values.
func ParityConfigToMultiGethGenesis(c *parity.ParityChainSpec) (*paramtypes.Genesis, error) {
	mg := &paramtypes.Genesis{
		Config: &paramtypes.MultiGethChainConfig{},
	}

	if err := Convert(c, mg); err != nil {
		return nil,err
	}
	if err := Convert(c, mg.Config); err != nil {
		return nil, err
	}
	return mg, nil
}