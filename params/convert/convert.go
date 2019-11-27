package convert

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/ethereum/go-ethereum/params/types/common"
)

// Automagically translate between [Must|]Setters and Getters.
func Convert(from, to interface{}) error {
	// Interfaces must be either ChainConfigurator or GenesisBlocker.
	for i, v := range []interface{}{
		from, to,
	}{
		_, genesiser := v.(common.GenesisBlocker)
		_, chainconfer := v.(common.ChainConfigurator)
		if !genesiser && !chainconfer {
			return fmt.Errorf("%d value neither chain nor genesis configurator", i)
		}
	}

	// Order may matter; configuration parameters may be interdependent across data structures, eg EIP1283 and Genesis builtins.
	// Try to order translation sensibly

	fromGener, fromGenerOk := from.(common.GenesisBlocker)
	toGener, toGenerOk := to.(common.GenesisBlocker)

	// Set Genesis.
	if fromGenerOk && toGenerOk {
		et := fromGener.GetSealingType()
		switch et {
		case common.BlockSealing_Ethereum:
			k := reflect.TypeOf((*common.GenesisBlocker)(nil)).Elem()
			if err := convert(k, fromGener, toGener); err != nil {
				return err
			}
		default:
			return common.UnsupportedConfigError(common.ErrUnsupportedConfigFatal, "sealing type", et)
		}

		// Set accounts (genesis).
		if err := fromGener.ForEachAccount(toGener.UpdateAccount); err != nil {
			return err
		}
	}

	fromChainer, fromChainerOk := from.(common.ChainConfigurator)
	toChainer, toChainerOk := to.(common.ChainConfigurator)

	if !fromChainerOk || !toChainerOk {
		return nil
	}

	// Set general chain parameters.
	k := reflect.TypeOf((*common.CatHerder)(nil)).Elem()
	if err := convert(k, fromChainer, toChainer); err != nil {
		return err
	}

	// Set consensus engine params.
	engineType := fromChainer.GetConsensusEngineType()
	if err := toChainer.MustSetConsensusEngineType(engineType); err != nil {
		return common.UnsupportedConfigError(err, "consensus engine", engineType)
	}
	switch engineType {
	case common.ConsensusEngineT_Ethash:
		k := reflect.TypeOf((*common.EthashConfigurator)(nil)).Elem()
		if err := convert(k, fromChainer, toChainer); err != nil {
			return err
		}
	case common.ConsensusEngineT_Clique:
		k := reflect.TypeOf((*common.CliqueConfigurator)(nil)).Elem()
		if err := convert(k, fromChainer, toChainer); err != nil {
			return err
		}
	default:
		return common.UnsupportedConfigError(common.ErrUnsupportedConfigFatal, "consensus engine", common.ConsensusEngineT_Unknown)
	}

	return nil
}

func convert(k reflect.Type, source, target interface{}) error {
	for i := 0; i < k.NumMethod(); i++ {
		method := k.Method(i)

		var setterName string
		if strings.HasPrefix(method.Name, "Get") {
			setterName = strings.Replace(method.Name, "Get", "Set", 1)
		}
		if _, ok := k.MethodByName(setterName); !ok {
			continue
		}

		callGetIn := []reflect.Value{}
		response := reflect.ValueOf(source).MethodByName(method.Name).Call(callGetIn)
		setResponse := reflect.ValueOf(target).MethodByName(setterName).Call(response)

		if !setResponse[0].IsNil() {
			err := setResponse[0].Interface().(error)
			e := common.UnsupportedConfigError(err, strings.TrimPrefix(method.Name, "Get"), response[0].Interface())
			if common.IsFatalUnsupportedErr(err) {
				return e
			}
			log.Println(e) // FIXME?
		}
	}
	return nil
}
func Equal(k reflect.Type, a, b common.Configurator) (string, bool) {
	m, _, err := compare(k.Elem(), a, b) // TODO: maybe return a value, or even a dedicated type, for better debugging
	if err == nil {
		return "", true
	}
	return m, false
}

func compare(k reflect.Type, source, target interface{}) (method string, value interface{}, err error) {
	if k.NumMethod() == 0 {
		return "", "", errors.New("empty comparison")
	}
	for i := 0; i < k.NumMethod(); i++ {
		method := k.Method(i)

		if !strings.HasPrefix(method.Name, "Get") {
			continue
		}

		callGetIn := []reflect.Value{}
		response := reflect.ValueOf(source).MethodByName(method.Name).Call(callGetIn)
		response2 := reflect.ValueOf(target).MethodByName(method.Name).Call(callGetIn)

		if !reflect.DeepEqual(response[0].Interface(), response2[0].Interface()) {
			return method.Name, struct {
				source, target interface{}
			}{
				response, response2,
			}, errors.New("reflect.DeepEqual")
		}
	}
	return "", nil, nil
}


//func UnmarshalConfigurator(data []byte) (Configurator, error) {
//	var unMarshalErr error
//	for _, t := range []Configurator{
//		&parity.ParityChainSpec{},
//	}{
//		unMarshalErr = json.Unmarshal(data, t)
//		if unMarshalErr == nil {
//			return t, nil
//		}
//	}
//	return nil, unMarshalErr
//}

//func UnmarshalChainConfigurator(data []byte) (common.ChainConfigurator, error) {
//	var unMarshalErr error
//	for _, t := range []interface{}{
//		&paramtypes.ChainConfig{},
//		&goethereum.ChainConfig{},
//		&parity.ParityChainSpec{},
//		&paramtypes.Genesis{
//			Config: &paramtypes.ChainConfig{},
//		},
//		&paramtypes.Genesis{
//			Config: &goethereum.ChainConfig{},
//		},
//	}{
//		unMarshalErr = json.Unmarshal(data, t)
//		if unMarshalErr == nil {
//			rt := t
//			if d, ok := t.(*paramtypes.Genesis); ok {
//				if e, ok := d.Config.(*paramtypes.ChainConfig); ok {
//					rt = e
//				} else if e, ok := d.Config.(*goethereum.ChainConfig); ok {
//					rt = e
//				}
//			}
//			return rt.(common.ChainConfigurator), nil
//		}
//	}
//	return nil, unMarshalErr
//}