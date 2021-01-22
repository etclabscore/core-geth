package eth

import (
	"bytes"
	"testing"
)

func TestMakeExtraDataDefault(t *testing.T) {
	if !bytes.Contains(makeExtraData(nil), []byte("CoreGeth")) {
		t.Error("missing extra data default client identifier")
	}
}
