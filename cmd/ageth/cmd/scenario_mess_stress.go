package cmd

import (
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/log"
)

func generateMESSStressTest() func(nodes *agethSet) {
	return func(nodes *agethSet) {
		log.Info("Running MESS stress test scenario")

		// tabula rasa
		nodes.forEach(func(i int, a *ageth) {
			a.stopMining()
			a.setMaxPeers(100)
			var res bool
			a.client.Call(&res, "admin_ecbp1100", hexutil.Uint64(999999999).String())
		})

		vader := nodes.indexed(0)
		vader.log.Info("Is bad guy")

		log.Info("Running pace car...")
		nodes.dexedni(0).startMining(13)
		time.Sleep(60 * time.Second)
		nodes.dexedni(0).stopMining()

		herdHeadMin := nodes.headMin()
		nodes.forEach(func(i int, a *ageth) {
			if a.block().number > herdHeadMin {
				a.truncateHead(herdHeadMin)
			}
			var res bool
			a.client.Call(&res, "admin_ecbp1100", hexutil.Uint64(1).String())
		})

		// Enforce strong peering relationships
		scenarioDone := false
		defer func() {
			scenarioDone = true
		}()
		go func() {
			for !scenarioDone {
				nodes.forEach(func(i int, a *ageth) {
					if a.peers.len() < 15 {
						a.addPeer(nodes.where(func(g *ageth) bool {
							return g.name != a.name
						}).random())
					}
				})
			}
		}()

	}
}
