package core

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"math/rand"
	"testing"

	"github.com/ethereum/go-ethereum/consensus/ethash"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/params"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

var yuckyGlobalTestEnableMess = false

func runMESSTest(t *testing.T, easyL, hardL, caN int, easyT, hardT int64) (hardHead bool, err error) {
	// Generate the original common chain segment and the two competing forks
	engine := ethash.NewFaker()

	db := rawdb.NewMemoryDatabase()
	genesis := params.DefaultMessNetGenesisBlock()
	genesisB := MustCommitGenesis(db, genesis)

	chain, err := NewBlockChain(db, nil, genesis.Config, engine, vm.Config{}, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer chain.Stop()
	chain.EnableArtificialFinality(yuckyGlobalTestEnableMess)

	easy, _ := GenerateChain(genesis.Config, genesisB, engine, db, easyL, func(i int, b *BlockGen) {
		b.SetNonce(types.EncodeNonce(uint64(rand.Int63n(math.MaxInt64))))
		b.OffsetTime(easyT)
	})
	commonAncestor := easy[caN-1]
	hard, _ := GenerateChain(genesis.Config, commonAncestor, engine, db, hardL, func(i int, b *BlockGen) {
		b.SetNonce(types.EncodeNonce(uint64(rand.Int63n(math.MaxInt64))))
		b.OffsetTime(hardT)
	})

	if _, err := chain.InsertChain(easy); err != nil {
		t.Fatal(err)
	}
	_, err = chain.InsertChain(hard)
	hardHead = chain.CurrentBlock().Hash() == hard[len(hard)-1].Hash()
	return
}

func TestBlockChain_AF_ECBP11355(t *testing.T) {
	yuckyGlobalTestEnableMess = true
	defer func() {
		yuckyGlobalTestEnableMess = false
	}()

	cases := []struct {
		easyLen, hardLen, commonAncestorN int
		easyOffset, hardOffset            int64
		hardGetsHead, accepted            bool
	}{
		// Hard has insufficient total difficulty / length and is rejected.
		{
			5000, 7500, 2500,
			50, -9,
			false, false,
		},
		// Hard has sufficient total difficulty / length and is accepted.
		{
			1000, 7, 995,
			60, 0,
			true, true,
		},
		// Hard has sufficient total difficulty / length and is accepted.
		{
			1000, 7, 995,
			60, 7,
			true, true,
		},
		// Hard has sufficient total difficulty / length and is accepted.
		{
			1000, 1, 999,
			30, 1,
			true, true,
		},
		// Hard has sufficient total difficulty / length and is accepted.
		{
			500, 3, 497,
			0, -8,
			true, true,
		},
		// Hard has sufficient total difficulty / length and is accepted.
		{
			500, 4, 496,
			0, -9,
			true, true,
		},
		// Hard has sufficient total difficulty / length and is accepted.
		{
			500, 5, 495,
			0, -9,
			true, true,
		},
		// Hard has sufficient total difficulty / length and is accepted.
		{
			500, 6, 494,
			0, -9,
			true, true,
		},
		// Hard has sufficient total difficulty / length and is accepted.
		{
			500, 7, 493,
			0, -9,
			true, true,
		},
		// Hard has sufficient total difficulty / length and is accepted.
		{
			500, 8, 492,
			0, -9,
			true, true,
		},
		// Hard has sufficient total difficulty / length and is accepted.
		{
			500, 9, 491,
			0, -9,
			true, true,
		},
		// Hard has sufficient total difficulty / length and is accepted.
		{
			500, 12, 488,
			0, -9,
			true, true,
		},
		// Hard has sufficient total difficulty / length and is accepted.
		{
			500, 20, 480,
			0, -9,
			true, true,
		},
		// Hard has sufficient total difficulty / length and is accepted.
		{
			500, 40, 460,
			0, -9,
			true, true,
		},
		// Hard has sufficient total difficulty / length and is accepted.
		{
			500, 60, 440,
			0, -9,
			true, true,
		},
		// Hard has insufficient total difficulty / length and is rejected.
		{
			500, 200, 300,
			0, -9,
			false, false,
		},
		// Hard has insufficient total difficulty / length and is rejected.
		{
			500, 200, 300,
			7, -9,
			false, false,
		},
		// Hard has insufficient total difficulty / length and is rejected.
		{
			500, 200, 300,
			17, -9,
			false, false,
		},
		// Hard has sufficient total difficulty / length and is accepted.
		{
			500, 200, 300,
			47, -9,
			true, true,
		},
		// Hard has insufficient total difficulty / length and is rejected.
		{
			500, 200, 300,
			47, -8,
			false, false,
		},
		// Hard has insufficient total difficulty / length and is rejected.
		{
			500, 200, 300,
			17, -8,
			false, false,
		},
		// Hard has insufficient total difficulty / length and is rejected.
		{
			500, 200, 300,
			7, -8,
			false, false,
		},
		// Hard has insufficient total difficulty / length and is rejected.
		{
			500, 200, 300,
			0, -8,
			false, false,
		},
		// Hard has insufficient total difficulty / length and is rejected.
		{
			500, 100, 400,
			0, -7,
			false, false,
		},
		// Hard is accepted, but does not have greater total difficulty,
		// and is not set as the chain head.
		{
			1000, 1, 900,
			60, -9,
			false, true,
		},
		// Hard is shorter, but sufficiently heavier chain, is accepted.
		{
			500, 100, 390,
			60, -9,
			true, true,
		},
	}

	for i, c := range cases {
		hardHead, err := runMESSTest(t, c.easyLen, c.hardLen, c.commonAncestorN, c.easyOffset, c.hardOffset)
		if (err != nil && c.accepted) || (err == nil && !c.accepted) || (hardHead != c.hardGetsHead) {
			t.Errorf("case=%d [easy=%d hard=%d ca=%d eo=%d ho=%d] want.accepted=%v want.hardHead=%v got.hardHead=%v err=%v",
				i,
				c.easyLen, c.hardLen, c.commonAncestorN, c.easyOffset, c.hardOffset,
				c.accepted, c.hardGetsHead, hardHead, err)
		}
	}
}

func TestBlockChain_GenerateMESSPlot(t *testing.T) {
	t.Skip("This test plots graph of chain acceptance for visualization.")

	easyLen := 200
	maxHardLen := 100

	generatePlot := func(title, fileName string) {
		p, err := plot.New()
		if err != nil {
			log.Panic(err)
		}
		p.Title.Text = title
		p.X.Label.Text = "Block Depth"
		p.Y.Label.Text = "Relative Block Time Delta (10 seconds + y)"

		accepteds := plotter.XYs{}
		rejecteds := plotter.XYs{}
		sides := plotter.XYs{}

		for i := 1; i <= maxHardLen; i++ {
			for j := -9; j <= 8; j++ {
				fmt.Println("running", i, j)
				hardHead, err := runMESSTest(t, easyLen, i, easyLen-i, 0, int64(j))
				point := plotter.XY{X: float64(i), Y: float64(j)}
				if err == nil && hardHead {
					accepteds = append(accepteds, point)
				} else if err == nil && !hardHead {
					sides = append(sides, point)
				} else if err != nil {
					rejecteds = append(rejecteds, point)
				}

				if err != nil {
					t.Log(err)
				}
			}
		}

		scatterAccept, _ := plotter.NewScatter(accepteds)
		scatterReject, _ := plotter.NewScatter(rejecteds)
		scatterSide, _ := plotter.NewScatter(sides)

		pixelWidth := vg.Length(1000)

		scatterAccept.Color = color.RGBA{R: 152, G: 236, B: 161, A: 255}
		scatterAccept.Shape = draw.BoxGlyph{}
		scatterAccept.Radius = vg.Length((float64(pixelWidth) / float64(maxHardLen)) * 2 / 3)
		scatterReject.Color = color.RGBA{R: 236, G: 106, B: 94, A: 255}
		scatterReject.Shape = draw.BoxGlyph{}
		scatterReject.Radius = vg.Length((float64(pixelWidth) / float64(maxHardLen)) * 2 / 3)
		scatterSide.Color = color.RGBA{R: 190, G: 197, B: 236, A: 255}
		scatterSide.Shape = draw.BoxGlyph{}
		scatterSide.Radius = vg.Length((float64(pixelWidth) / float64(maxHardLen)) * 2 / 3)

		p.Add(scatterAccept)
		p.Legend.Add("Accepted", scatterAccept)
		p.Add(scatterReject)
		p.Legend.Add("Rejected", scatterReject)
		p.Add(scatterSide)
		p.Legend.Add("Sidechained", scatterSide)

		p.Legend.YOffs = -30

		err = p.Save(pixelWidth, 300, fileName)
		if err != nil {
			log.Panic(err)
		}
	}
	yuckyGlobalTestEnableMess = true
	defer func() {
		yuckyGlobalTestEnableMess = false
	}()
	baseTitle := fmt.Sprintf("Accept/Reject Reorgs: Relative Time (Difficulty) over Proposed Segment Length (%d-block original chain)", easyLen)
	generatePlot(baseTitle, "reorgs-MESS.png")
	yuckyGlobalTestEnableMess = false
	generatePlot("WITHOUT MESS: "+baseTitle, "reorgs-noMESS.png")
}

func TestEcbp11355AGSinusoidalA(t *testing.T) {
	cases := []struct{
		in, out float64
	}{
		{0, 1},
		{25132, 31},
	}
	tolerance := 0.0000001
	for i, c := range cases {
		if got := ecbp11355AGSinusoidalA(c.in); got < c.out - tolerance || got > c.out + tolerance {
			t.Fatalf("%d: in: %0.6f want: %0.6f got: %0.6f", i, c.in, c.out, got)
		}
	}
}
