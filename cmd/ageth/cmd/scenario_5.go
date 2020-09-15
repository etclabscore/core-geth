package cmd

import (
	"math/rand"
	"time"

	"github.com/ethereum/go-ethereum/log"
)

func scenario5(nodes *agethSet) {

	badGuyWaitsUntilNAboveGoodGuys := uint64(30)

	badGuy := nodes.indexed(0)
	badGuy.behaviorsInterval = time.Second
	badGuy.behaviors = append(badGuy.behaviors, func(self *ageth) {
		if self.isMining && self.peers.len() > 0 {
			self.removeAllPeers()
		}
	})
	badGuy.run()
	badGuy.startMining(5)

	goodGuys := nodes.subset(1, nodes.len())

	setupGoodGuys := func() {
		// iterate number of desired guys, with the index actually tracking
		// the world set though.
		// no safety checks to make sure the world is providing us with enough
		// guys because this whole thing Austin is probably going to overwrite.
		for i := 0; i < goodGuys.len(); i++ {
				guy := goodGuys.indexed(i)
				guy.run()
				if i > 0 {
					// Connect them all to each other.
					for j := 0; j < i; j++ {
						guy.addPeer(goodGuys.indexed(j))
					}
				}
				if i < goodGuys.len()/3 {
					guy.startMining(13)
				}
		}
	}

	// badGuyAttack assumes that the bad guy is disconnected from the network (original sin).
	// Since this is assumed and this function will be called repeatedly, this is how the attack should finish.
	badGuyAttack := func() {
		badGuy.log.Info("Beginning bad guy attack")
		if badGuy.block().number < goodGuys.headMax()+badGuyWaitsUntilNAboveGoodGuys {
			badGuy.log.Warn("Bad guy below desired distance from herd", "head", badGuy.block().number, "goodguys.max", goodGuys.headMax())
			if !badGuy.isMining {
				badGuy.startMining(5)
			}
			return
		}
		defer func() {
			for badGuy.peers.len() > 0 {
				badGuy.removeAllPeers()
			}
		}()

		addGuys := func() {
			// add a bunch of random good guys
			badGuy.addPeers(goodGuys.randomN(0.5))

			// make sure we add the miners
			badGuy.addPeers(goodGuys.where(func(a *ageth) bool { return a.isMining }))
		}

		addGuys()

		timeout := time.NewTimer(60 * time.Second)
		tick := time.NewTicker(5 * time.Second)
		defer func() {
			timeout.Stop()
			tick.Stop()
		}()
		for {
			select {
			case <-timeout.C:
				log.Warn("Good guys win")
				// Since good guys won, this means that the common ancestor will keep falling (deeper),
				// and chances of a good attack will diminish (since attacker is only allowed so far above
				// the herd's head).
				// So let's have the attacker sync with the herd first, then reattempt (shorten the reorged chain).
				for badGuy.peers.len() < 5 || badGuy.block().number > goodGuys.headMax() {
					if badGuy.isMining {
						badGuy.stopMining()
					}
					if badGuy.block().number > goodGuys.headMax() {
						badGuy.truncateHead(goodGuys.headMax())
					}
					badGuy.addPeer(goodGuys.random())
				}
				return
			case <-tick.C:
				addGuys()
			default:
				// If any of the good guys are subverted, consider it a victory.
				if goodGuys.headMax() > badGuy.block().number-(badGuyWaitsUntilNAboveGoodGuys/2) {
					log.Warn("Bad guy wins")
					return
				}
			}
		}
	}
	// time.Sleep(time.Duration(5+rand.Int31n(25)) * time.Second)
	setupGoodGuys()
	for {
		badGuyAttack()
		time.Sleep(time.Duration(rand.Int31n(10)) * time.Second)
	}

}
