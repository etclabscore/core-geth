package cmd

import (
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
)

func nodeTDRatioAB(a, b *ageth, commonBlockTD *big.Int) (tdA, tdB *big.Int, ratio float64) {
	tdA, tdB = a.getTd(), b.getTd()
	if tdA.Cmp(bigZero) == 0 || tdB.Cmp(bigZero) == 0 {
		return bigZero, bigZero, 0
	}
	r, _ := new(big.Float).Quo(
		new(big.Float).Sub(new(big.Float).SetInt(tdA), new(big.Float).SetInt(commonBlockTD)),
		new(big.Float).Sub(new(big.Float).SetInt(tdB), new(big.Float).SetInt(commonBlockTD)),
	).Float64()
	return tdA, tdB, r
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

func generateScenarioPartitioning(followGravity bool, minDuration, maxDuration time.Duration) func(set *agethSet) {
	return func(nodes *agethSet) {

		log.Info("Running partitioning scenario")

		// tabula rasa
		nodes.forEach(func(i int, a *ageth) {
			a.stopMining()
			a.setMaxPeers(100)
			var res bool
			a.client.Call(&res, "admin_ecbp1100", hexutil.Uint64(999999999).String())
		})

		luke := nodes.indexed(0)
		luke.log.Info("== Luke ==")
		solo := nodes.indexed(1)
		solo.log.Info("== Solo ==")

		if !luke.sameChainAs(solo) {
			log.Warn("Luke and Solo on different chains", "action", "truncate to genesis")
			luke.truncateHead(0)
			solo.truncateHead(0)
		}

		log.Info("Asserting protagonist sync")
		herdMax := nodes.headMax()
		for luke.block().number < herdMax || solo.block().number < herdMax {
			time.Sleep(10 * time.Second)
		}

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

		solo.mustEtherbases([]common.Address{solo.coinbase})
		luke.mustEtherbases([]common.Address{luke.coinbase}) // double evil
		defer func() {
			solo.mustEtherbases([]common.Address{})
			luke.mustEtherbases([]common.Address{})
		}()
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
		soloPaceMaker := func(solo, luke *ageth, forkedTd *big.Int) {
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
			if balance >= lower && balance <= upper {
				return
			}
			mineAtInterval, ivalMax, ivalMin := solo.mining, normalBlockTime+1, normalBlockTime-1
			if balance < lower {
				mineAtInterval--
			} else if balance > upper {
				mineAtInterval++
			}
			if mineAtInterval < ivalMin {
				mineAtInterval = ivalMin
			} else if mineAtInterval > ivalMax {
				mineAtInterval = ivalMax
			}
			if mineAtInterval == solo.mining {
				return
			}
			solo.startMining(mineAtInterval)
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

		start := time.Now()
		for {
			limitRatio := ecbp1100AGSinusoidalA(float64(solo.latestBlock.Time - forkedBlockTime))
			_, _, balance := nodeTDRatioAB(solo, luke, forkedTD)
			tlog := log.Info

			/*
			soloTD := solo.getTd()
			unitTDRat := 1 + (float64(solo.latestBlock.Difficulty.Uint64()) / (float64(soloTD) - float64(forkedTD)))
			*/

			unitTDRat := new(big.Float).Quo(
				new(big.Float).SetInt(solo.latestBlock.Difficulty),
				new(big.Float).SetInt(new(big.Int).Sub(solo.getTd(), forkedTD)),
				)
			unitTDRat.Add(new(big.Float).SetInt(bigOne), unitTDRat)

			// if unitTDRat < limitRatio {
			if unitTDRat.Cmp(new(big.Float).SetFloat64(limitRatio)) < 0 {
				/*
				IF WE COMPARE BALANCE,
				THIS WILL ALWAYS BE TRUE BECAUSE BALANCE IS SUPPOSED TO BE 1:1!
				WHAT WE WANT TO MEASURE IS __MARGINAL__ TOTAL DIFFICULTY: ie.
				the TotalDifficulty ratio of Solo's current block (proposed) over his last block (current).
				Once this delta proportion passes below the antigravity ratio, then
				we can be confident that people are pretty much where they're gonna stay.
				 */
				log.Warn("SoloTD/LukeTD ratio beneath antigravity", "antigravity", limitRatio, "S/L", balance)
				tlog = log.Warn
				if time.Since(start) > minDuration {
					break
				}
			}
			if time.Since(start) > maxDuration {
				break
			}
			tlog("Status", "fork.number", forkedBlock.number, "fork.age", common.PrettyAge(time.Unix(int64(forkedBlockTime), 0)), "unit.rate" ,unitTDRat , "antigravity", limitRatio, "S/L", balance)
			time.Sleep(10 * time.Second)
		}

		time.Sleep(10*time.Minute)

		_, _, resultingTDRatio := nodeTDRatioAB(luke, solo, new(big.Int).SetUint64(forkedBlock.number))

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
