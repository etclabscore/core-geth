package params

import (
	"testing"
)

// TestGenesisHashMintMe tests that MintMeGenesisHash is the correct value for the genesis configuration.
func TestGenesisHashMintMe(t *testing.T) {
	genesis := DefaultMintMeGenesisBlock()
	block := genesisToBlock(genesis, nil)
	if block.Hash() != MintMeGenesisHash {
		t.Errorf("want: %s, got: %s", MintMeGenesisHash.Hex(), block.Hash().Hex())
	}
}
