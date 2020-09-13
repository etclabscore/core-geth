package cmd

import (
	"math/rand"
	"time"

	"github.com/ethereum/go-ethereum/log"
)

/*
	 This (test) program:
	- Uses ONE "bad guy" who will be endowed with a variable amount of "hash power," which
	is implemented as --miner.threads. This bad guy will usually be disconnected from everyone,
	but will intermittently (randomly) connect to some number of good guys. The bad guy will stay
	disconnected until he achieves some number of blocks more (measured in block numbers) than
	the highest block of the good guys set.
	- Attacking, the bad guy repeatedly call `addPeer` in half of the peers. He really wants peers.
	- Uses a configurable number of good guys. A configurable number of them will be miners,
	also with hashpower defined as --miner.threads. Remember that they will be connected, so their
	individual hashpowers will be combined, eg. 1+1+1=3.
	- The bad guy's success is decided if he is able subvert ANY of the honest nodes after an arbitrary period.
	- If he's is unable to do so (good guys win), then he will stop mining and wait for the herd to catch
	up, he will (re)sync, and mount the attack again. This prevents his attempted reorg from growing massively
    long without compensating for that length with an equally long out-front chain.
    - This will go on forever.
*/

func scenario1(eventChan chan interface{}) {
	// "Global"s, don't touch.
	goodGuys := newAgethSet()

	// Configurable
	numberOfGoodGuys := 20
	numberOfGoodGuyMiners := 4
	eachGoodGuyMiningPower := 9
	badGuyMiningPower := 3
	badGuyWaitsUntilNAboveGoodGuys := uint64(100)
	goodGuysPeerTarget := (numberOfGoodGuys / 2)

	badGuy := newAgeth()
	badGuy.eventChan = eventChan
	badGuy.behaviorsInterval = 60 * time.Second
	badGuy.behaviors = append(badGuy.behaviors, func() {
		if badGuy.block().number > goodGuys.headMax()+(badGuyWaitsUntilNAboveGoodGuys*3) {
			badGuy.log.Warn("Bad guy too far ahead of herd")
			badGuy.stopMining()
		} else if !badGuy.isMining && badGuy.block().number < goodGuys.headMax()+(badGuyWaitsUntilNAboveGoodGuys*2) {
			badGuy.log.Warn("Bad guy back with herd")
			badGuy.startMining(badGuyMiningPower)
		}
	})
	badGuy.run()
	badGuy.stopMining()
	badGuy.startMining(badGuyMiningPower)
	world.push(badGuy)

	for i := 0; i < numberOfGoodGuys; i++ {
		go func(i int) {
			guy := newAgeth()
			guy.eventChan = eventChan
			guy.run()
			if goodGuys.len() > 0 {
				for j := 0; j < goodGuys.len()/2; j++ {
					guy.addPeer(goodGuys.random())
				}
			}
			if i < numberOfGoodGuyMiners {
				log.Info("Start mining", guy.name)
				guy.startMining(eachGoodGuyMiningPower)
			}
			goodGuys.push(guy)
			world.push(guy)
		}(i)
	}
	// badGuyAttack assumes that the bad guy is disconnected from the network (original sin).
	// Since this is assumed and this function will be called repeatedly, this is how the attack should finish.
	badGuyAttack := func() {
		log.Info("Beginning bad guy attack")
		for badGuy.peers.len() > 0 {
			badGuy.removeAllPeers()
		}
		if badGuy.block().number < goodGuys.headMax()+badGuyWaitsUntilNAboveGoodGuys {
			log.Warn("Bad guy below desired distance from herd", "bad guy", badGuy.block().number, "good guy", goodGuys.headMax())
			return
		}
		defer func() {
			for badGuy.peers.len() > 0 {
				badGuy.removeAllPeers()
			}
		}()

		addGuys := func() {
			// add a bunch of random good guys
			for i := 0; i < goodGuys.len(); i++ {
				p := goodGuys.random()
				if rand.Float64() < 0.5 {
					badGuy.addPeer(p)
				}
			}
			// make sure we add the miners
			for i := 0; i < numberOfGoodGuyMiners; i++ {
				badGuy.addPeer(goodGuys.ageths[i])
			}

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
				for badGuy.peers.len() < goodGuysPeerTarget || badGuy.block().number > goodGuys.headMax() {
					if badGuy.isMining {
						badGuy.stopMining()
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
	goodGuyBasicChurn := func() {
		for _, g := range goodGuys.all() {
			// Try to keep the good guys at about 1/4 total network size. This is probably a little higher
			// proportionally than the mainnet, but close-ish, I guess. (Assuming ~500-800 total nodes, around 25-50 peer limits).
			// The important related variable here is DefaultSyncMinPeers, which Permapoint uses
			// to disable itself if the node's peer count falls below it.
			// The normal value for DefaultSyncMinPeers is 5, and the normal MaxPeers value (which geth
			// tries to fill) is 25.
			offset := g.peers.len() - goodGuysPeerTarget
			if offset > 0 {
				g.removePeer(g.peers.random()) // not exact, but closeish
			} else if offset < 0 {
				g.addPeer(goodGuys.random())
			} else {
				// simulate some small churn
				if rand.Float32() < 0.5 {
					g.addPeer(goodGuys.random())
				} else {
					g.removePeer(g.peers.random())
				}
			}
		}
	}
	time.Sleep(time.Duration(5+rand.Int31n(25)) * time.Second)
	for {
		// wait until everyone gets fired up ok
		if goodGuys.len() < 2 {
			time.Sleep(5 * time.Second)
			continue
		}
		goodGuyBasicChurn()
		time.Sleep(5 * time.Second)
		badGuyAttack()
		time.Sleep(time.Duration(rand.Int31n(10)) * time.Second)
	}
}
