package convert

import (
	"errors"
	"log"
	"reflect"
	"strings"

	"github.com/ethereum/go-ethereum/params/types/common"
)

// Automagically translate between [Must|]Setters and Getters.
func Convert(from, to common.Configurator) error {
	if _, ok := from.(common.Configurator); !ok {
		return errors.New("from value not a configurator")
	}
	if _, ok := to.(common.Configurator); !ok {
		return errors.New("to value not a configurator")
	}

	// Order may matter; configuration parameters may be interdependent across data structures, eg EIP1283 and Genesis builtins.
	// Try to order translation sensibly.
	// Set Genesis.
	et := from.GetSealingType()
	switch et {
	case common.BlockSealing_Ethereum:
		k := reflect.TypeOf((*common.GenesisBlocker)(nil)).Elem()
		if err := convert(k, from, to); err != nil {
			return err
		}
	default:
		return common.UnsupportedConfigError(common.ErrUnsupportedConfigFatal, "sealing type", et)
	}

	// Set accounts (genesis).
	if err := from.ForEachAccount(to.UpdateAccount); err != nil {
		return err
	}

	// Set general chain parameters.
	k := reflect.TypeOf((*common.CatHerder)(nil)).Elem()
	if err := convert(k, from, to); err != nil {
		return err
	}

	// Set consensus engine params.
	engineType := from.GetConsensusEngineType()
	if err := to.MustSetConsensusEngineType(engineType); err != nil {
		return common.UnsupportedConfigError(err, "consensus engine", engineType)
	}
	switch engineType {
	case common.ConsensusEngineT_Ethash:
		k := reflect.TypeOf((*common.EthashConfigurator)(nil)).Elem()
		if err := convert(k, from, to); err != nil {
			return err
		}
	case common.ConsensusEngineT_Clique:
		k := reflect.TypeOf((*common.CliqueConfigurator)(nil)).Elem()
		if err := convert(k, from, to); err != nil {
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

