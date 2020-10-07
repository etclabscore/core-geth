package eth

import "github.com/ethereum/go-ethereum/common"

type PrivateTxPoolAPI struct {
	eth *Ethereum
}

func NewPrivateTxPoolAPI(eth *Ethereum) *PrivateTxPoolAPI {
	return &PrivateTxPoolAPI{eth: eth}
}

// returns true if a specified TX was in TX pool
func (pool *PrivateTxPoolAPI) RemoveTransaction(hash common.Hash) bool {
	return pool.eth.txPool.RemoveTx(hash)
}
