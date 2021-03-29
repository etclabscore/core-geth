package tests

import (
	"testing"

	"github.com/ethereum/go-ethereum/params/confp"
	"github.com/ethereum/go-ethereum/params/types/coregeth"
	"github.com/ethereum/go-ethereum/params/types/parity"
	"github.com/ethereum/go-ethereum/tests"
	"github.com/go-test/deep"
)

func TestParityChainspec_CoreGeth(t *testing.T) {
	coregethA := tests.Forks["ETC_Magneto"]

	pspec := &parity.ParityChainSpec{}
	err := confp.Convert(coregethA, pspec)
	if err != nil {
		t.Fatal(err)
	}

	err = confp.Equivalent(coregethA, pspec)
	if err != nil {
		t.Fatal(err)
	}

	coregethB := &coregeth.CoreGethChainConfig{}
	err = confp.Convert(pspec, coregethB)
	if err != nil {
		t.Fatal(err)
	}

	diff := deep.Equal(coregethA, coregethB)
	for _, line := range diff {
		t.Error(line)
	}

	eip1108a := coregethA.GetEIP1108Transition()
	eip1108b := coregethB.GetEIP1108Transition()

	ne := eip1108a == nil && eip1108b != nil
	ne = ne || eip1108a != nil && eip1108b == nil
	ne = ne || *eip1108a != *eip1108b
	if ne {
		t.Fatal("ne")
	}
	if eip1108a == nil {
		t.Fatal("nil")
	}
	t.Logf("%v %v", *eip1108a, *eip1108b)
}
