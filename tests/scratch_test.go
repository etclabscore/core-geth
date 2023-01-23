package tests

import (
	"fmt"
	"testing"
)

func TestHexMarshaling(t *testing.T) {
	out := fmt.Sprintf("%#x", 1)
	t.Log(out)
}
