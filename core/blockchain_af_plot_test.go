package core

import (
	"fmt"
	"image/color"
	"log"
	"math/big"
	"math/rand"
	"testing"

	emath "github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/params/vars"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

// Some weird constants to avoid constant memory allocs for them.
var (
	bigMinus99 = big.NewInt(-99)
	big1       = big.NewInt(1)
	big0       = big.NewInt(0)
)


func segmentTotalDifficulty (segment []*big.Int) *big.Int {
	out := big.NewInt(0)
	for _, b := range segment {
		out.Add(out, b)
	}
	return out
}

func generateDifficultySegment(initDifficulty, maxDifficulty *big.Int, duration int64) []*big.Int {

	parentDiff := new(big.Int).Set(initDifficulty)

	blockTime := uint64(1) // aggressor
	isBase := maxDifficulty.Cmp(big0) == 0
	if isBase {
		// base
		blockTime = 10
	}

	outset := []*big.Int{}
	for duration > 0 {

		// parent difficulty met the max (allowed hashrate fulfilled)
		// keep block difficulty constant while still creating blocks as quickly as possible
		if !isBase && parentDiff.Cmp(maxDifficulty) >= 0 {
			blockTime = 9
		}
		duration -= int64(blockTime)

		// https://github.com/ethereum/EIPs/issues/100
		// algorithm:
		// diff = (parent_diff +
		//         (parent_diff / 2048 * max((2 if len(parent.uncles) else 1) - ((timestamp - parent.timestamp) // 9), -99))
		//        ) + 2^(periodCount - 2)
		out := new(big.Int)
		out.Div(new(big.Int).SetUint64(blockTime), vars.EIP100FDifficultyIncrementDivisor)

		// if parent.UncleHash == types.EmptyUncleHash {
		// 	out.Sub(big1, out)
		// } else {
		// 	out.Sub(big2, out)
		// }
		out.Sub(big1, out)

		out.Set(emath.BigMax(out, bigMinus99))

		out.Mul(new(big.Int).Div(parentDiff, vars.DifficultyBoundDivisor), out)
		out.Add(out, parentDiff)

		// after adjustment and before bomb
		out.Set(emath.BigMax(out, vars.MinimumDifficulty))

		parentDiff.Set(out) // set for next iteration

		outset = append(outset, out)
	}
	return outset
}

func localSegmentDifficulties(spanSeconds int64) []*big.Int {
	return generateDifficultySegment(
		params.DefaultMessNetGenesisBlock().Difficulty,
		big0,
		spanSeconds)
}

func ecbp1100SegmentAccept(x, localTD, propTD *big.Int) (ok bool) {
	eq := ecbp1100PolynomialV(x)
	want := eq.Mul(eq, localTD)
	got := new(big.Int).Mul(propTD, ecbp1100PolynomialVCurveFunctionDenominator)

	return got.Cmp(want) >= 0 // reject => got < want
}

func eyalSirer(sameDiff, sameLen bool) bool {
	if !sameDiff || ! sameLen {
		return true
	}
	return rand.Float64() < 0.5
}

func testDepth(localSegment []*big.Int, maxDifficulty *big.Int, depthSeconds uint64) (accept bool) {
	proposedSegment := generateDifficultySegment(
		params.DefaultMessNetGenesisBlock().Difficulty,

		// localSegment[len(localSegment)-proposedSpanSeconds], // localBlockOffset = proposedSegment.len / fitBlocks
		maxDifficulty,
		int64(depthSeconds),
	)

	localBlockOffset := depthSeconds / 10
	ls := localSegment[uint64(len(localSegment))-localBlockOffset:]

	localTD := segmentTotalDifficulty(ls)
	propTD := segmentTotalDifficulty(proposedSegment)

	if propTD.Cmp(localTD) < 0 {
		return false
	}

	if eyalSirer(localTD.Cmp(propTD) == 0, len(ls) == len(proposedSegment)) {
		ecbp1100OK := ecbp1100SegmentAccept(new(big.Int).SetUint64(depthSeconds), localTD, propTD)
		if ecbp1100OK {
			accept = true
		}
	}
	return accept
}

func generateDepthPlot(p *plot.Plot) {

	oks, notoks := plotter.XYs{}, plotter.XYs{}

	ls := localSegmentDifficulties(3000*10)

	quo := big.NewInt(3)

	for numerator := uint64(1); numerator <= 15; numerator++ {
		num := new(big.Int).SetUint64(numerator)
		md := new(big.Int).Mul(num, ls[0])
		md.Div(md, quo)
		rat, _ := new(big.Float).Quo(
			new(big.Float).SetInt(num),
			new(big.Float).SetInt(quo),
			).Float64()

		for depth := 0; depth <= 900*10; depth += 10 {
			ok := testDepth(ls, md, uint64(depth))
			if ok {
				oks = append(oks, plotter.XY{X: rat, Y: float64(depth)})
			} else {
				notoks = append(notoks, plotter.XY{X: rat, Y: float64(depth)})
			}
		}
	}

	scatterAccept, _ := plotter.NewScatter(oks)
	scatterSide, _ := plotter.NewScatter(notoks)

	pixelWidth := vg.Length(len(ls) * 110 / 100)

	scatterAccept.Color = color.RGBA{R: 0, G: 200, B: 11, A: 255}
	scatterAccept.Shape = draw.BoxGlyph{}
	scatterAccept.Radius = vg.Length((float64(pixelWidth) / float64(len(ls))) * 2 / 3)
	scatterSide.Color = color.RGBA{R: 220, G: 227, B: 246, A: 255}
	scatterSide.Shape = draw.BoxGlyph{}
	scatterSide.Radius = vg.Length((float64(pixelWidth) / float64(len(ls))) * 2 / 3)

	p.Add(scatterAccept)
	p.Add(scatterSide)

	err := p.Save(600, pixelWidth, fmt.Sprintf("refactored-reorg-depth-%s.png", ""))
	if err != nil {
		log.Panic(err)
	}
}

func TestGenerateMESSPlot_Depth_HashrateM(t *testing.T) {
	p, err := plot.New()
	if err != nil {
		log.Panic(err)
	}
	p.Title.Text = "Depth over BadGuy Hashrate"
	p.Y.Label.Text = "Block Depth in Seconds"
	p.X.Label.Text = "Antagonist Hashrate"

	generateDepthPlot(p)
}
