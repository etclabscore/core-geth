package cmd

import (
	"time"

	"github.com/ethereum/go-ethereum/core/types"
)

func balanceLikeSneakyCobra(nodes *agethSet) {

	// tabula rasa
	nodes.forEach(func(i int, a *ageth) {
		a.stopMining()
	})

	luke := nodes.indexed(0)
	vader := nodes.indexed(1)

	vader.registerNewHeadCallback(func(a *ageth, h *types.Block) {
		if a.coinbase != h.Coinbase() {
			a.truncateHead(h.NumberU64()-1)
		}
	})
	defer func() {
		vader.purgeNewHeadCallbacks()
	}()

	darkJourneyDuration := 20*time.Minute

	

	// set up 10 good guy miners.
	// distribute their power via poisson

	// set up 10 bad guy miners
	// distribute their power via poisson

	// bad guys break off into their own clique

	// bad guys mine for
}
