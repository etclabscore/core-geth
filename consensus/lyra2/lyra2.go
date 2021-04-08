package lyra2

import (
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/consensus"
    "github.com/ethereum/go-ethereum/consensus/ethash"
    "github.com/ethereum/go-ethereum/core/state"
    "github.com/ethereum/go-ethereum/core/types"
    "github.com/ethereum/go-ethereum/rpc"
    "math/big"
)

type Lyra2 struct {
    ethashEngine *ethash.Ethash
}

func New(ethashConfig ethash.Config) *Lyra2 {
    ethashEngine := ethash.New(ethashConfig, nil, false)

    return &Lyra2{
        ethashEngine: ethashEngine,
    }
}

// Author implements consensus.Engine, returning the Ethereum address recovered
// from the signature in the header's extra-data section.
func (c *Lyra2) Author(header *types.Header) (common.Address, error) {
    return c.ethashEngine.Author(header)
}

// VerifyHeader checks whether a header conforms to the consensus rules.
func (c *Lyra2) VerifyHeader(chain consensus.ChainHeaderReader, header *types.Header, seal bool) error {
    return c.ethashEngine.VerifyHeader(chain, header, seal)
}

// VerifyHeaders is similar to VerifyHeader, but verifies a batch of headers. The
// method returns a quit channel to abort the operations and a results channel to
// retrieve the async verifications (the order is that of the input slice).
func (c *Lyra2) VerifyHeaders(chain consensus.ChainHeaderReader, headers []*types.Header, seals []bool) (chan<- struct{}, <-chan error) {
    return c.ethashEngine.VerifyHeaders(chain, headers, seals)
}

// VerifyUncles implements consensus.Engine, always returning an error for any
// uncles as this consensus mechanism doesn't permit uncles.
func (c *Lyra2) VerifyUncles(chain consensus.ChainReader, block *types.Block) error {
    return c.ethashEngine.VerifyUncles(chain, block)
}

// VerifySeal implements consensus.Engine, checking whether the signature contained
// in the header satisfies the consensus protocol requirements.
func (c *Lyra2) VerifySeal(chain consensus.ChainHeaderReader, header *types.Header) error {
    return c.ethashEngine.VerifySeal(chain, header)
}

// Prepare implements consensus.Engine, preparing all the consensus fields of the
// header for running the transactions on top.
func (c *Lyra2) Prepare(chain consensus.ChainHeaderReader, header *types.Header) error {
    return c.ethashEngine.Prepare(chain, header)
}

// Finalize implements consensus.Engine, ensuring no uncles are set, nor block
// rewards given.
func (c *Lyra2) Finalize(chain consensus.ChainHeaderReader, header *types.Header, state *state.StateDB, txs []*types.Transaction, uncles []*types.Header) {
    c.ethashEngine.Finalize(chain, header, state, txs, uncles)
}

// FinalizeAndAssemble implements consensus.Engine, ensuring no uncles are set,
// nor block rewards given, and returns the final block.
func (c *Lyra2) FinalizeAndAssemble(chain consensus.ChainHeaderReader, header *types.Header, state *state.StateDB, txs []*types.Transaction, uncles []*types.Header, receipts []*types.Receipt) (*types.Block, error) {
    return c.ethashEngine.FinalizeAndAssemble(chain, header, state, txs, uncles, receipts)
}

// Seal implements consensus.Engine, attempting to create a sealed block using
// the local signing credentials.
func (c *Lyra2) Seal(chain consensus.ChainHeaderReader, block *types.Block, results chan<- *types.Block, stop <-chan struct{}) error {
    return c.ethashEngine.Seal(chain, block, results, stop)
}

// CalcDifficulty is the difficulty adjustment algorithm. It returns the difficulty
// that a new block should have:
// * DIFF_NOTURN(2) if BLOCK_NUMBER % SIGNER_COUNT != SIGNER_INDEX
// * DIFF_INTURN(1) if BLOCK_NUMBER % SIGNER_COUNT == SIGNER_INDEX
func (c *Lyra2) CalcDifficulty(chain consensus.ChainHeaderReader, time uint64, parent *types.Header) *big.Int {
    return c.ethashEngine.CalcDifficulty(chain, time, parent)
}

// SealHash returns the hash of a block prior to it being sealed.
func (c *Lyra2) SealHash(header *types.Header) common.Hash {
    return c.ethashEngine.SealHash(header)
}

// Close implements consensus.Engine. It's a noop for clique as there are no background threads.
func (c *Lyra2) Close() error {
    return c.ethashEngine.Close()
}

// APIs implements consensus.Engine, returning the user facing RPC API to allow
// controlling the signer voting.
func (c *Lyra2) APIs(chain consensus.ChainHeaderReader) []rpc.API {
    return c.ethashEngine.APIs(chain)
}