package cmd

import (
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params/vars"
)

func nodeTDRatioAB(a, b *ageth, commonBlockTD uint64) (tdA, tdB uint64, ratio float64) {
	tdA, tdB = a.getTd(), b.getTd()
	if tdB == 0 {
		return tdA, 0, 0
	}
	return tdA, tdB, (float64(tdA) - float64(commonBlockTD)) / (float64(tdB) / float64(commonBlockTD))
}

/*
ecbp1100AGSinusoidalA is a sinusoidal function.

OPTION 3: Yet slower takeoff, yet steeper eventual ascent. Has a differentiable ceiling transition.
h(x)=15 sin((x+12000 π)/(8000))+15+1

*/
func ecbp1100AGSinusoidalA(x float64) (antiGravity float64) {
	ampl := float64(15)   // amplitude
	pDiv := float64(8000) // period divisor
	phaseShift := math.Pi * (pDiv * 1.5)
	peakX := math.Pi * pDiv // x value of first sin peak where x > 0
	if x > peakX {
		// Cause the x value to limit to the x value of the first peak of the sin wave (ceiling).
		x = peakX
	}
	return (ampl * math.Sin((x+phaseShift)/pDiv)) + ampl + 1
}

func generateScenarioPartitioning(followGravity bool, duration time.Duration) func(set *agethSet) {
	return func(nodes *agethSet) {

		log.Info("Running partitioning scenario")

		// tabula rasa
		nodes.forEach(func(i int, a *ageth) {
			a.stopMining()
			a.setMaxPeers(100)
			var res bool
			a.client.Call(&res, "admin_ecbp1100", 9999999999)
		})
		log.Info("Waiting for the dust to settle")
		time.Sleep(1 * time.Minute)
		// outer:
		// 	for {
		// 		for _, n := range nodes.all() {
		// 			if n.block().number != nodes.headMax() {
		// 				time.Sleep(10 * time.Second)
		// 				continue outer
		// 			}
		// 		}
		// 		break
		// 	}

		nodes.forEach(func(i int, a *ageth) {
			var res bool
			a.client.Call(&res, "admin_ecbp1100", 1)
		})

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

		luke := nodes.indexed(0)
		luke.log.Info("I'm Luke")
		solo := nodes.indexed(1)
		solo.log.Info("I'm Solo")

		solo.mustEtherbases([]common.Address{solo.coinbase})
		normalBlockTime := 13

		// Toes on the same line
		if luke.block().number > solo.block().number {
			luke.truncateHead(solo.block().number)
		} else if solo.block().number > luke.block().number {
			solo.truncateHead(luke.block().number)
		}

		forkedBlock := luke.block()
		forkedBlockTime := luke.latestBlock.Time
		forkedTD := luke.getTd()
		if forkedTD != solo.getTd() || forkedBlock.number != solo.block().number {
			log.Error("The force is strong but not impossibly strong")
			return
		}

		eitherNewHead := false
		soloPaceMaker := func() {
			if !eitherNewHead {
				return
			}
			soloTD, _, balance := nodeTDRatioAB(solo, luke, forkedTD)
			tdRatioTolerance := (float64(soloTD) / float64(vars.DifficultyBoundDivisor.Int64()))
			wantRatio := float64(1)
			if followGravity {
				wantRatio = ecbp1100AGSinusoidalA(float64(solo.latestBlock.Time - forkedBlockTime))
			}
			if balance > wantRatio+tdRatioTolerance {
				solo.startMining(42)
			} else if balance < wantRatio-tdRatioTolerance {
				solo.startMining(2)
			} else {
				solo.startMining(normalBlockTime)
			}
			eitherNewHead = false
		}
		go func() {
			for !scenarioDone {
				soloPaceMaker()
				time.Sleep(1 * time.Second)
			}
		}()

		solo.registerNewHeadCallback(func(self *ageth, h *types.Header) {
			eitherNewHead = true
		})
		luke.registerNewHeadCallback(func(self *ageth, h *types.Header) {
			eitherNewHead = true
		})
		defer func() {
			solo.purgeNewHeadCallbacks()
			luke.purgeNewHeadCallbacks()
		}()

		luke.startMining(normalBlockTime)
		solo.startMining(normalBlockTime)

		time.Sleep(duration)

		_, _, resultingTDRatio := nodeTDRatioAB(luke, solo, forkedBlock.number)

		if luke.sameChainAs(solo) {
			log.Error("Test failed; Luke is on the dark the side :(")
		}

		var lightSideNodesCount, darkSideNodesCount int

		nodes.forEach(func(i int, a *ageth) {
			if a.sameChainAs(luke) {
				lightSideNodesCount++
			} else if a.sameChainAs(solo) {
				darkSideNodesCount++
			}
		})

		distinctChains := nodes.distinctChains()

		report := &finalReport{
			Converted:             darkSideNodesCount,
			Unconverted:           lightSideNodesCount,
			DistinctChains:        distinctChains,
			Nodes:                 make(map[string]common.Hash),
			DifficultyRatio:       resultingTDRatio,
			TargetDifficultyRatio: 1,
			AttackerShouldWin:     true,
			AttackerWon:           len(distinctChains) == 2,
		}
		for _, node := range nodes.all() {
			report.Nodes[node.name] = node.block().hash
		}
		data, _ := json.Marshal(report)
		fmt.Println(string(data))
	}
}
