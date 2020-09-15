package cmd

import (
	"time"

	"github.com/ethereum/go-ethereum/log"
)

func scenario3(world *agethSet) {
	goodGuys := newAgethSet()
	badGuys := newAgethSet()

	numberOfGoodGuys := 30
	numberOfGoodGuyMiners := 5
	eachGoodGuyMiningPower := 1
	// goodGuysPeerTarget := (numberOfGoodGuys / 3)

	numberOfBadGuys := 3
	eachBadGuyMiningPower := 6

	minions := newAgethSet() // these will be selected from good guys, but will "report" to badGuy
	numberOfMinions := 20

	for i := 0; i < numberOfBadGuys; i++ {
		badGuy := world.all()[i]
		badGuy.run()
		badGuy.startMining(eachBadGuyMiningPower)
		world.push(badGuy)
		badGuys.push(badGuy)
	}
	for i := 0; i < numberOfGoodGuys; i++ {
		go func(i int) {
			guy := world.all()[i + numberOfBadGuys] // just offset the index given the afore badguys looper
			guy.run()
			// guy.withStandardPeerChurn(goodGuysPeerTarget, goodGuys)
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
	// quit := false
	// go func() {
	// 	for {
	// 		if quit {
	// 			return
	// 		}
	// 		// wait until everyone gets fired up ok
	// 		if goodGuys.len() < numberOfGoodGuys {
	// 			time.Sleep(5 * time.Second)
	// 			continue
	// 		}
	// 		goodGuyBasicChurn(goodGuys, goodGuysPeerTarget)
	// 	}
	// }()
	for i := 0; i < 2; i++ {
		// wait until everyone gets fired up ok
		for goodGuys.len() < numberOfGoodGuys {
			time.Sleep(5 * time.Second)
			continue
		}
		done := make(chan struct{})
		for _, badGuy := range badGuys.all() {
			badGuy := badGuy
			go badGuyAttack(badGuy, goodGuys, minions, eachBadGuyMiningPower, numberOfGoodGuyMiners, done)
		}
		for i := 0; i < badGuys.len(); i++ {
			<-done
		}
		close(done)
	}
	// quit = true
}
