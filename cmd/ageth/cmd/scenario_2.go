package cmd

import (
	"math/rand"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
)

/*

This test program stress the "fray-ability" of the edge of the
network.

In this test, the goal of the attacker is to partition the network.
He's not trying to do a massive reorg, he just wants to cause mayhem.

To accomplish this devious end, he will attempt to exploit the permapoint/flashpoint
system to cause nodes to behave incongruently in regard to small(ish) reorgs.

He will act immediately on the publication of a flashpoint (his learning of it, at
least), and thereon will add peers like crazy (broadcasting them his
head), and trying to convince SOME of them to take the candy.

In this scenario, as programmed, the badguy has access to the "God mode" state view,
and can target his attack accordingly. As such, he should be able pull off the MAXIMUM number
of subvertible peers.

*/
func badGuyAttack(badGuy *ageth, goodGuys, minions *agethSet, badGuyMiningPower, numberOfGoodGuyMiners int, done chan struct{}) {

	log.Info("Beginning bad guy attack")
	for badGuy.peers.len() > 0 {
		badGuy.removeAllPeers()
	}
	for badGuy.block().number <= goodGuys.headMax() {
		if !badGuy.isMining {
			badGuy.startMining(12)
		}
		log.Info("Bad guy waiting ahead head")
		badGuy.removeAllPeers()
		time.Sleep(5 * time.Second)
	}
	// if badGuy.block().number < goodGuys.headMax()+badGuyWaitsUntilNAboveGoodGuys {
	// 	log.Warn("Bad guy below desired distance from herd", "bad guy", badGuy.block().number, "good guy", goodGuys.headMax())
	// 	return
	// }
	log.Info("Bad guy waiting permapoint")
	var perma *types.Block
	for perma == nil {
		head := minions.headHeader()
		if head == nil {
			continue
		}
		// if core.IsPermapoint(nil, head.Header()) {
		// 	perma = head
		// }
	}
	log.Info("Bad guy waiting flashpoint")
	var flash *types.Block
	for flash == nil {
		head := minions.headHeader()
		if head == nil {
			continue
		}
		// if !(head.NumberU64()-vars.FlashpointThreshold >= perma.NumberU64()) {
		// 	continue
		// }
		// if core.IsFlashpoint(nil, head.Header()) {
		// 	flash = head
		// }
	}
	log.Info("Bad guy attack!")
	// make sure we add the miners
	for i := 0; i < numberOfGoodGuyMiners; i++ {
		badGuy.addPeer(goodGuys.ageths[i])
	}
	// add a bunch of good guys
	for _, peer := range goodGuys.all() {
		if minions.contains(peer) {
			continue
		}
		badGuy.addPeer(peer)
	}
	timeout := time.NewTimer(30 * time.Second)
	defer timeout.Stop()
	checkTime := time.NewTimer(5*time.Second)
	defer checkTime.Stop()
	defer func() {
		done <- struct{}{}
	}()
	for {
		select {
		case <-checkTime.C:
			if badGuy.peers.len() == 0 {
				// got rejected!...
				log.Info("Bad guy attack failed")
				return
			}
		case <-timeout.C:

			// Bad guy has mounted attack.
			// Now allow bad guy to fall back in with the herd.
			log.Info("Bad guy done")
			badGuy.stopMining()
			if goodGuysMode := goodGuys.headMode(); goodGuysMode > 0 {
				badGuy.truncateHead(goodGuysMode)
			}
			for badGuy.block().number != goodGuys.headMax() {
				time.Sleep(1 * time.Second)
				badGuy.addPeer(goodGuys.random())
			}

			// Set bad guy up for next attack.
			badGuy.removeAllPeers()
			return
		default:
		}
	}
	// TODO improve this; check for network partitions
}

func goodGuyBasicChurn(goodGuys *agethSet, goodGuysPeerTarget int) {
	for _, g := range goodGuys.all() {
		time.Sleep(time.Second)
		// Try to keep the good guys at about 1/4 total network size. This is probably a little higher
		// proportionally than the mainnet, but close-ish, I guess. (Assuming ~500-800 total nodes, around 25-50 peer limits).
		// The important related variable here is DefaultSyncMinPeers, which Permapoint uses
		// to disable itself if the node's peer count falls below it.
		// The normal value for DefaultSyncMinPeers is 5, and the normal MaxPeers value (which geth
		// tries to fill) is 25.
		if rand.Float64() < 0.5 {
			continue
		}
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


func scenario2(world *agethSet) {
	goodGuys := newAgethSet()
	numberOfGoodGuys := 20
	numberOfGoodGuyMiners := 3
	eachGoodGuyMiningPower := 1

	minions := newAgethSet() // these will be selected from good guys, but will "report" to badGuy
	numberOfMinions := 8

	badGuyMiningPower := 12
	goodGuysPeerTarget := (numberOfGoodGuys / 2)
	// badGuyWaitsUntilNAboveGoodGuys := uint64(10)

	badGuy := world.all()[0]
	badGuy.run()
	badGuy.startMining(badGuyMiningPower)
	world.push(badGuy)

	for i := 1; i <= numberOfGoodGuys; i++ {
		go func(i int) {
			guy := world.all()[i]
			guy.run()
			if goodGuys.len() > 0 {
				for j := 0; j < goodGuys.len()/2; j++ {
					guy.addPeer(goodGuys.random())
				}
			}
			if i < numberOfGoodGuyMiners {
				log.Info("Start mining", guy.name)
				guy.startMining(eachGoodGuyMiningPower)
				guy.behaviorsInterval = 5*time.Second
				guy.behaviors = append(guy.behaviors, func(self *ageth) {
					if self.peers.len() < 3 {
						self.stopMining()
					} else if !self.isMining {
						self.startMining(eachGoodGuyMiningPower)
					}
				})
			}

			if i > numberOfGoodGuys-numberOfMinions {
				minions.push(guy)
			}
			goodGuys.push(guy)
			world.push(guy)
		}(i)
	}

	// 		time.Sleep(time.Duration(rand.Int31n(10)) * time.Second)
	quit := false
	go func() {
		for {
			if quit {
				return
			}
			// wait until everyone gets fired up ok
			if goodGuys.len() < numberOfGoodGuys {
				time.Sleep(5 * time.Second)
				continue
			}
			goodGuyBasicChurn(goodGuys, goodGuysPeerTarget)
			time.Sleep(time.Duration(5+rand.Int31n(25)) * time.Second)
		}
	}()

	for i := 0; i < 3; i++ {
		// wait until everyone gets fired up ok
		for goodGuys.len() < numberOfGoodGuys / 2 {
			time.Sleep(5 * time.Second)
			continue
		}
		done := make(chan struct{})
		badGuyAttack(badGuy, goodGuys, minions, badGuyMiningPower, numberOfGoodGuyMiners, done)
		close(done)
	}
	quit = true
}
