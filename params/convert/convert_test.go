package convert

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/params/types"
	"github.com/ethereum/go-ethereum/params/types/aleth"
	"github.com/ethereum/go-ethereum/params/types/common"
	"github.com/ethereum/go-ethereum/params/types/goethereum"
	"github.com/ethereum/go-ethereum/params/types/parity"
)

func mustOpenF(t *testing.T, fabbrev string, into interface{}) {
	b, err := ioutil.ReadFile(filepath.Join("testdata", fmt.Sprintf("stureby_%s.json", fabbrev)))
	if err != nil {
		t.Fatal(err)
	}
	err = json.Unmarshal(b, &into)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_UnmarshalJSON(t *testing.T) {
	for _, f := range []string{
		"geth", "parity", "aleth",
	} {
		switch f {
		case "geth":
			c := &paramtypes.Genesis{}
			mustOpenF(t, f, c)
			if *c.Config.GetNetworkID() != 314158 {
				t.Errorf("networkid")
			}
		case "parity":
			p := &parity.ParityChainSpec{}
			mustOpenF(t, f, p)
			_, err := ParityConfigToMultiGethGenesis(p)
			if err != nil {
				t.Error(err)
			}
		case "aleth":
			a := &aleth.AlethGenesisSpec{}
			mustOpenF(t, f, a)
		}
	}
}

func TestConvert(t *testing.T) {
	spec := parity.ParityChainSpec{}
	mustOpenF(t, "parity", &spec)

	spec2 := parity.ParityChainSpec{}
	err := Convert(&spec, &spec2)
	if err != nil {
		t.Error(err)
	}

	if method, equal := Equal(reflect.TypeOf((*common.Configurator)(nil)), &spec, &spec2); !equal {
		t.Error("not equal", method)
	}
}

// TestConfiguratorImplementationsSatisfied tests that data types expected
// to fulfil certain interfaces do fill them.
func TestConfiguratorImplementationsSatisfied(t *testing.T) {
	for _, ty := range []interface{}{
		&parity.ParityChainSpec{},
	} {
		_ = ty.(common.Configurator)
	}

	for _, ty := range []interface{}{
		&goethereum.ChainConfig{},
		&paramtypes.ChainConfig{},
	} {
		_ = ty.(common.ChainConfigurator)
	}

	for _, ty := range []interface{}{
		&paramtypes.Genesis{},
	} {
		_ = ty.(common.GenesisBlocker)
	}
}

func TestCompatible(t *testing.T) {
	spec := &parity.ParityChainSpec{}
	fns, names := common.Transitions(spec)
	for i, fn := range fns {
		t.Log(names[i], fn())
	}
	t.Log(fns)
}
