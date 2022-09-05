package goethereum

import (
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/params/types/ctypes"
)

func TestChainConfig_converting(t *testing.T) {
	var c interface{} = &ChainConfig{}
	fromChainer := c.(ctypes.ChainConfigurator)

	if _, ok := reflect.TypeOf(fromChainer).Elem().FieldByName("Converting"); ok {
		reflect.ValueOf(fromChainer).Elem().Field(0).SetBool(true)
	}
}
