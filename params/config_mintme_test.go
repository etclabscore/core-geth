package params

import (
	"testing"
)

// TestGenesisHashMINTME tests that MINTMEGenesisHash is the correct value for the genesis configuration.
func TestGenesisHashMINTME(t *testing.T) {
	genesis := DefaultMINTMEGenesisBlock()
	block := genesisToBlock(genesis, nil)
	if block.Hash() != MINTMEGenesisHash {
		t.Errorf("want: %s, got: %s", MINTMEGenesisHash.Hex(), block.Hash().Hex())
	}
}
