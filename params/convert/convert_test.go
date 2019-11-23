package convert

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/ethereum/go-ethereum/params/types"
	"github.com/ethereum/go-ethereum/params/types/aleth"
	"github.com/ethereum/go-ethereum/params/types/parity"
)

func Test_UnmarshalJSON(t *testing.T) {
	mustOpenF := func(fabbrev string, into interface{}) {
		b, err := ioutil.ReadFile(filepath.Join("testdata", fmt.Sprintf("stureby_%s.json", fabbrev)))
		if err != nil {
			t.Fatal(err)
		}
		err = json.Unmarshal(b, &into)
		if err != nil {
			t.Fatal(err)
		}
	}
	for _, f := range []string{
		"geth", "parity", "aleth",
	} {
		switch f {
		case "geth":
			c := &paramtypes.Genesis{}
			mustOpenF(f, c)
			if c.Config.NetworkID != 314158 {
				t.Errorf("networkid")
			}

			c.Config.UpgradeToSchedules()

			if len(c.Config.DifficultyBombDelaySchedule) == 0 {
				t.Errorf("no diff sched")
			}
			//t.Log(spew.Sdump(c.Config.DifficultyBombDelaySchedule))
			if len(c.Config.BlockRewardSchedule) == 0 {
				t.Errorf("no block sched")
			}
			//t.Log(spew.Sdump(c.Config.BlockRewardSchedule))
		case "parity":
			p := &parity.ParityChainSpec{}
			mustOpenF(f, p)
			_, err := ParityConfigToMultiGethGenesis(p)
			if err != nil {
				t.Error(err)
			}
		case "aleth":
			a := &aleth.AlethGenesisSpec{}
			mustOpenF(f, a)
		}

	}

}
