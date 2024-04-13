package params

import (
	"testing"
)

// TestGenesisHashHypra tests that HypraGenesisHash is the correct value for the genesis configuration.
func TestGenesisHashHypra(t *testing.T) {
	genesis := DefaultHypraGenesisBlock()
	block := genesisToBlock(genesis, nil)
	if block.Hash() != HypraGenesisHash {
		t.Errorf("want: %s, got: %s", HypraGenesisHash.Hex(), block.Hash().Hex())
	}
}
