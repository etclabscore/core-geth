package cmd

import (
	"math/rand"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
)

func scenario4(world *agethSet) {
	goodGuys := newAgethSet()
	badGuy := world.all()[0]
		badGuy.run()
		badGuy.startMining(10)
		world.push(badGuy)

	numberOfGoodGuys := 20
	numberOfGoodGuyMiners := 3
	eachGoodGuyMiningPower := 4
	goodGuysPeerTarget := (numberOfGoodGuys / 2)

	minions := newAgethSet() // these will be selected from good goodGuys, but will "report" to badGuy
	numberOfMinions := 20

	for i := 1; i < numberOfGoodGuys; i++ {
		go func(i int) {
			guy := world.all()[i]
			guy.run()
			if world.len() > 0 {
				for j := 0; j < world.len()/2; j++ {
					guy.addPeer(world.random())
				}
			}
			goodGuys.push(guy)
			world.push(guy)
			if i < numberOfGoodGuyMiners {
				guy.startMining(eachGoodGuyMiningPower)
			}
			if i < numberOfMinions {
				minions.push(guy)
			}
		}(i)
	}
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
			goodGuyBasicChurn(world, goodGuysPeerTarget)
			time.Sleep(time.Duration(5+rand.Int31n(25)) * time.Second)
		}
	}()
	for i := 0; i < 2; i++ {
		// wait until everyone gets fired up ok
		for goodGuys.len() < numberOfGoodGuys {
			time.Sleep(5 * time.Second)
			continue
		}
		done := make(chan struct{})
		go badGuyAttack4(badGuy, goodGuys, minions, numberOfGoodGuyMiners, done)
		<-done
		close(done)
	}
	quit = true
}

func badGuyAttack4(badGuy *ageth, goodGuys, minions *agethSet, numberOfGoodGuyMiners int, done chan struct{}) {

	addGuys := func() {
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
	}

	// guy is connected normally.
	// waits until permapoint occurs.
	log.Info("Bad guy begins plotting his attack...")
	for minions.headMax() < 20 {
		continue
	}

	badgers:
	for {
		switch {
		case badGuy.block().number < goodGuys.headMax() - 2:
			log.Info("Bad guy behind herd, adding peers")
			addGuys()
			time.Sleep(5*time.Second)
		default:
			break badgers
		}
	}

	log.Info("Bad guy waiting permapoint")
	var permas *types.Block
	for permas == nil {
		head := goodGuys.headBlock()
		if head == nil {
			continue
		}
		// if core.IsPermapoint(nil, head.Header()) {
		// 	permas = head
		// }
	}


		i := 0
		for badGuy.peers.len() > 0 || i < 4 {
			badGuy.removeAllPeers()
			i++
			time.Sleep(time.Second)
		}

	// starts (keeps) mining.
	// waits until flashpoints occurs
	log.Info("Bad guy waiting flashpoint")
	var flash *types.Block
	for flash == nil {
		head := minions.headBlock()
		if head == nil {
			continue
		}
		// if !(head.NumberU64()-vars.FlashpointThreshold >= permas.NumberU64()) {
		// 	continue
		// }
		// if core.IsFlashpoint(nil, head.Header()) {
		// 	flash = head
		// }
	}
	log.Info("Bad guy attack!")
	// connects.
	addGuys()
	// waits 60 seconds, continually reconnecting.
	time.Sleep(60*time.Second)
	done <- struct{}{}
}
