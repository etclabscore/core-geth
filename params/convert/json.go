package convert

import (
	"encoding/json"
	"errors"
	"strings"

	paramtypes "github.com/ethereum/go-ethereum/params/types"
	"github.com/ethereum/go-ethereum/params/types/common"
	"github.com/ethereum/go-ethereum/params/types/goethereum"
	"github.com/ethereum/go-ethereum/params/types/parity"
)

func UnmarshalChainConfigurator(input []byte) (common.ChainConfigurator, error) {
	ifaces := []interface{}{
		&goethereum.ChainConfig{},
		&paramtypes.MultiGethChainConfig{},
		&parity.ParityChainSpec{},
	}
	errs := []string{}
	for _, iface := range ifaces {
		err := json.Unmarshal(input, iface)
		if err != nil {
			errs = append(errs, err.Error())
		} else {
			return iface.(common.ChainConfigurator), nil
		}
	}
	return nil, errors.New(strings.Join(errs, "\n"))
}
