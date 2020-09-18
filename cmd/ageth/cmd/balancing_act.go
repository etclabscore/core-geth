package cmd

import (
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
)

func nodeTDRatioAB(a, b *ageth, commonBlockTD uint64) (tdA, tdB uint64, ratio float64) {
	tdA, tdB = a.getTd(), b.getTd()
	if tdA == 0 || tdB == 0 {
		return 0, 0, 0
	}
	return tdA, tdB, (float64(tdA) - float64(commonBlockTD)) / (float64(tdB) - float64(commonBlockTD))
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
			a.client.Call(&res, "admin_ecbp1100", hexutil.Uint64(999999999).String())
		})

		nodes.dexedni(0).startMining(13)
		log.Info("Waiting for the dust to settle")
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

		scenarioDone := false
		defer func() {
			scenarioDone = true
		}()
		go func() {
			for !scenarioDone {
				nodes.forEach(func(i int, a *ageth) {
					if a.peers.len() < 11 {
						a.addPeer(nodes.where(func(g *ageth) bool {
							return g.name != a.name
						}).random())
					}
				})
				time.Sleep(60 * time.Second)
			}
		}()

		luke := nodes.indexed(0)
		luke.log.Info("== Luke ==")
		solo := nodes.indexed(1)
		solo.log.Info("== Solo ==")

		solo.mustEtherbases([]common.Address{solo.coinbase})
		luke.mustEtherbases([]common.Address{luke.coinbase}) // double evil
		defer func() {
			solo.mustEtherbases([]common.Address{})
			luke.mustEtherbases([]common.Address{})
		}()
		normalBlockTime := 7

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
		soloPaceMaker := func(solo, luke *ageth, forkedTd uint64) {
			if !eitherNewHead {
				return
			}
			defer func() {
				eitherNewHead = false
			}()
			if luke.sameChainAs(solo) {
				forkedBlock = luke.block()
				forkedBlockTime = luke.latestBlock.Time
				forkedTD = luke.getTd()
				return
			}
			_, _, balance := nodeTDRatioAB(solo, luke, forkedTD)
			tdRatioTolerance := float64(1 / 2048)
			wantRatio := float64(1)
			if followGravity {
				wantRatio = ecbp1100AGSinusoidalA(float64(solo.latestBlock.Time - forkedBlockTime))
			}
			upper, lower := wantRatio+tdRatioTolerance, wantRatio-tdRatioTolerance
			if balance > upper && solo.mining != normalBlockTime+1 {
				log.Warn("Solo above balance", "want", wantRatio, "upper", upper, "got", balance)
				solo.startMining(normalBlockTime + 1)
			} else if balance < lower && solo.mining != normalBlockTime-1 {
				log.Warn("Solo below balance", "want", wantRatio, "lower", lower, "got", balance)
				solo.startMining(normalBlockTime - 1)
			} else if balance >= lower && balance <= upper && solo.mining != normalBlockTime {
				solo.startMining(normalBlockTime)
			}
		}
		go func() {
			for !scenarioDone {
				soloPaceMaker(solo, luke, forkedTD)
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

		go luke.startMining(normalBlockTime)
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
