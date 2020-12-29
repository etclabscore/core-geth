package node

import (
	"reflect"
	"testing"
)

type TestReceiver struct{}

// ReturnSigA thru ReturnSigD should be eligible.
func (r *TestReceiver) ReturnSigA() error {
	return nil
}

func (r *TestReceiver) ReturnSigB() (int, error) {
	return 0, nil
}

func (r *TestReceiver) ReturnSigC() int {
	return 0
}

func (r *TestReceiver) ReturnSigD() {}

// ReturnSigE should not be eligible (any returned error must be last).
func (r *TestReceiver) ReturnSigE() (error, int) {
	return nil, 0
}

// ReturnSigF should not be eligible (> 2 return values).
func (r *TestReceiver) ReturnSigF() (int, int, error) {
	return 0, 0, nil
}

func TestEligibleReturnSignature(t *testing.T) {
	tr := &TestReceiver{}

	cases := []struct {
		methodName string
		ok         bool
	}{
		{"ReturnSigA", true},
		{"ReturnSigB", true},
		{"ReturnSigC", true},
		{"ReturnSigD", true},
		{"ReturnSigE", false},
		{"ReturnSigF", false},
	}

	for _, c := range cases {
		method, ok := reflect.TypeOf(tr).MethodByName(c.methodName)
		if !ok {
			t.Fatalf("missing method: %s", c.methodName)
		}
		if got := eligibleReturnSignature(method); got != c.ok {
			t.Errorf("case: %s, got: %v, want: %v", c.methodName, got, c.ok)
		}
	}
}
