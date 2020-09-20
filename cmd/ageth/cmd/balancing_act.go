package cmd

import (
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/consensus/ethash"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
)

func nodeTDRatioAB(a, b *ageth, commonBlockTD *big.Int) (segmentTDA, segmentTDB *big.Int, ratio float64) {
	tdA, tdB := a.getTd(), b.getTd()
	if tdA.Cmp(bigZero) == 0 || tdB.Cmp(bigZero) == 0 {
		return bigOne, bigOne, 0
	}
	segmentTDA = new(big.Int).Sub(tdA, commonBlockTD)
	segmentTDB = new(big.Int).Sub(tdB, commonBlockTD)
	ratio, _ = new(big.Float).Quo(
		new(big.Float).SetInt(segmentTDA),
		new(big.Float).SetInt(segmentTDB),
	).Float64()
	return
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

var badGuyPacemakerLast common.Hash
var goodGuyPacemakerLast common.Hash
func badGuyPacemaker(badGuy, goodGuy *ageth, wantRatio float64, forkedTD *big.Int, forkedTime uint64) {
	if badGuy.block().hash == badGuyPacemakerLast && goodGuy.block().hash == goodGuyPacemakerLast {
		return
	}
	badGuyPacemakerLast = badGuy.block().hash
	goodGuyPacemakerLast = goodGuy.block().hash

	segA, segB, balance := nodeTDRatioAB(badGuy, goodGuy, forkedTD)

	tdRatioTolerance := float64(1 / 100)
	upper, lower := wantRatio+tdRatioTolerance, wantRatio-tdRatioTolerance

	constant := 13
	if balance >= lower && balance <= upper {
		if badGuy.mining != constant {
			badGuy.startMining(constant)
		}
		return
	}

	// Good guys' segment * targetRatio, eg. 42000000 * 1.02
	want := new(big.Float).Mul(new(big.Float).SetInt(segB), new(big.Float).SetFloat64(wantRatio))
	wantBig, _ := want.Int(nil) // as big

	// Difference between target segment difficulty and current segment difficulty
	wantUnitDifficulty := new(big.Int).Sub(wantBig, segA)

	bestDifference := big.NewInt(math.MaxInt64)
	bestTimeOffset := uint64(0)

	for i := uint64(42); i >= 3; i-- {
		// Sample the would-be difficulty for a range of block times.
		difficulty := ethash.CalcDifficulty(params.MessNetConfig, badGuy.block().time + i, badGuy.latestBlock)

		// Compare this sample with desired.
		difference := difficulty.Sub(difficulty, wantUnitDifficulty)
		difference.Abs(difference) // As absolute value.

		// If it's the closest to the target unit block difficulty that would yield
		// the closest to the overall chain segment ratio, peg it.
		if difference.Cmp(bestDifference) <= 0 {
			bestTimeOffset = i
			bestDifference.Set(difference)
		}
	}
	// Nudge the block time toward that.
	bestTimeOffset = (uint64(badGuy.mining) + bestTimeOffset) / 2
	if int(bestTimeOffset) == badGuy.mining {
		return
	}
	badGuy.startMining(int(bestTimeOffset))
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
		go func() {
			for !scenarioDone {
				if eitherNewHead {
					if luke.sameChainAs(solo) {
						forkedBlock = luke.block()
						forkedBlockTime = luke.latestBlock.Time
						forkedTD = luke.getTd()
					} else {
						wantRatio := float64(1)
						if followGravity {
							wantRatio = ecbp1100AGSinusoidalA(float64(solo.latestBlock.Time - forkedBlockTime))
						}
						badGuyPacemaker(solo, luke, wantRatio, forkedTD, forkedBlockTime)
						eitherNewHead = false
					}
				}
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
				Once this rate of growth passes below the antigravity ratio, then
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
