package paramtypes

import (
	"errors"
	"fmt"
	"math/big"
	"reflect"
	"strings"
	"testing"
)

// File contains proof of concept for chain configurator and
// conversion pattern.
// File does not contain any actual tests.

type CatHerderTester interface {
	GetActivationOfEIP1234() *big.Int
	SetActivationOfEIP1234(*big.Int)
}

type ForkerTester interface {
	Forked(func(*big.Int) bool, *big.Int) bool
}

type ChainConfiguratorTester interface {
	CatHerderTester
	ForkerTester
}

type testChainConfig struct {
	EIP1234Block *big.Int
}

func (t *testChainConfig) String() string {
	return fmt.Sprintf("eip1234: %v", t.EIP1234Block)
}

func (t *testChainConfig) GetActivationOfEIP1234() *big.Int {
	return t.EIP1234Block
}

func (t *testChainConfig) SetActivationOfEIP1234(i *big.Int) error {
	t.EIP1234Block = i
	return errors.New("fake unsupported config error")
	//return nil
}

func (t *testChainConfig) Forked(f func() *big.Int, n *big.Int) bool {
	return f().Cmp(n) >= 0
}

func TestChainConfigurator_Interface1(t *testing.T) {
	tcc := testChainConfig{big.NewInt(14)}

	if !tcc.Forked(tcc.GetActivationOfEIP1234, big.NewInt(14)) {
		t.Errorf("un-forking believable")
	}
}

func TestChainConfigurator_Interface_Discover(t *testing.T) {
	source := testChainConfig{big.NewInt(14)}
	dest := new(testChainConfig)

	k := reflect.TypeOf((*ChainConfiguratorTester)(nil)).Elem()
	for i := 0; i < k.NumMethod(); i++ {

		method := k.Method(i)

		var setterName string
		if strings.HasPrefix(method.Name, "Get") {
			setterName = strings.Replace(method.Name, "Get", "Set", 1)
		}
		if _, ok := k.MethodByName(setterName); !ok {
			continue
		}
		t.Log("Converting method:", method.Name)

		callGetIn := []reflect.Value{}
		response := reflect.ValueOf(&source).MethodByName(method.Name).Call(callGetIn)
		setResponse := reflect.ValueOf(dest).MethodByName(setterName).Call(response)

		t.Log(setResponse[0])
		t.Log("with error?", !setResponse[0].IsNil())
		if !setResponse[0].IsNil() {
			t.Logf("%v", setResponse[0].Interface().(error))
		}
	}
	t.Log("source", &source)
	t.Log("dest", dest)
}

