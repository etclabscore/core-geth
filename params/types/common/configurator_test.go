package common

import (
	"testing"
)

func TestIsEmpty(t *testing.T) {
	type thing struct {
		num uint64
	}
	thing1 := thing{42}

	if IsEmpty(thing1) {
		t.Error("empty not empty")
	}

	thing2 := thing{}
	if !IsEmpty(thing2) {
		t.Error("not empty empty")
	}

	if !IsEmpty(nil) {
		t.Error("nil not empty")
	}
}
