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
	var err1, err2, err3 error
	var config goethereum.ChainConfig
	err2 = json.Unmarshal(input, &config)
	if err2 == nil {
		return &config, nil
	}
	var config2 paramtypes.MultiGethChainConfig
	err1 = json.Unmarshal(input, &config2)
	if err1 == nil {
		return &config2, nil
	}
	var config3 parity.ParityChainSpec
	err3 = json.Unmarshal(input, &config3)
	if err3 == nil {
		return &config3, nil
	}
	return nil, errors.New(strings.Join([]string{err1.Error(), err2.Error(), err3.Error()}, "\n"))
}
