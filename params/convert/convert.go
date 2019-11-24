package convert

import (
	"errors"
	"log"
	"reflect"
	"strings"

	paramtypes "github.com/ethereum/go-ethereum/params/types"
	"github.com/ethereum/go-ethereum/params/types/common"
)

func Convert(from, to *paramtypes.ChainConfigurator) error {
	if _, ok := (*from).(paramtypes.ChainConfigurator); !ok {
		return errors.New("from value not a configurator")
	}
	if _, ok := (*to).(paramtypes.ChainConfigurator); !ok {
		return errors.New("to value not a configurator")
	}

	// Automagically translate between [Must|]Setters and Getters.
	magic := func(k reflect.Type) error {
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
			response := reflect.ValueOf(from).MethodByName(method.Name).Call(callGetIn)
			setResponse := reflect.ValueOf(to).MethodByName(setterName).Call(response)

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

	// Order may matter; configuration parameters may be interdependent across data structures, eg EIP1283 and Genesis builtins.
	// Try to order translation sensibly.
	// Set Genesis.
	switch (*from).GetSealingType() {
	case paramtypes.BlockSealing_Ethereum:
		k := reflect.TypeOf((*paramtypes.GenesisBlocker)(nil)).Elem()
		if err := magic(k); err != nil {
			return err
		}
	default:
		return common.UnsupportedConfigError(common.ErrUnsupportedConfigFatal, "sealing type", paramtypes.BlockSealing_Unknown)
	}

	// Set accounts (genesis).
	if err := (*from).ForEachAccount((*to).SetPlainAccount); err != nil {
		return err
	}

	// Set general chain parameters.
	k := reflect.TypeOf((*paramtypes.CatHerder)(nil)).Elem()
	if err := magic(k); err != nil {
		return err
	}

	// Set consensus engine params.
	engineType := (*from).GetConsensusEngineType()
	if err := (*to).MustSetConsensusEngineType(engineType); err != nil {
		return common.UnsupportedConfigError(err, "consensus engine", engineType)
	}
	switch engineType {
	case paramtypes.ConsensusEngineT_Ethash:
		k := reflect.TypeOf((*paramtypes.EthashConfigurator)(nil)).Elem()
		if err := magic(k); err != nil {
			return err
		}
	case paramtypes.ConsensusEngineT_Clique:
		k := reflect.TypeOf((*paramtypes.CliqueConfigurator)(nil)).Elem()
		if err := magic(k); err != nil {
			return err
		}
	default:
		return common.UnsupportedConfigError(common.ErrUnsupportedConfigFatal, "consensus engine", paramtypes.ConsensusEngineT_Unknown)
	}

	return nil
}
