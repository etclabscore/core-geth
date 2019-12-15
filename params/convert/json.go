package convert

import (
	"encoding/json"

	paramtypes "github.com/ethereum/go-ethereum/params/types"
	"github.com/ethereum/go-ethereum/params/types/common"
	"github.com/ethereum/go-ethereum/params/types/goethereum"
	"github.com/ethereum/go-ethereum/params/types/parity"
)

func UnmarshalChainConfigurator(input []byte) (common.ChainConfigurator, error) {
	var map1 = make(map[string]interface{})
	err := json.Unmarshal(input, &map1)
	if err != nil {
		return nil, err
	}
	if _, ok := map1["params"]; ok {
		pspec := &parity.ParityChainSpec{}
		err = json.Unmarshal(input, pspec)
		if err != nil {
			return nil, err
		}
		return pspec, nil
	}

	if _, ok := map1["networkId"]; ok {
		mspec := &paramtypes.MultiGethChainConfig{}
		err = json.Unmarshal(input, mspec)
		if err != nil {
			return nil, err
		}
		return mspec, nil
	}

	gspec := &goethereum.ChainConfig{}
	err = json.Unmarshal(input,gspec)
	if err != nil {
		return nil, err
	}
	return gspec, nil
}
