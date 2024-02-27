package core

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/ethash"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/triedb"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

func runMESSTest2(t *testing.T, enableMess bool, easyL, hardL, caN int, easyT, hardT int64) (hardHead bool, err error, hard, easy []*types.Block) {
	// Generate the original common chain segment and the two competing forks
	engine := ethash.NewFaker()

	db := rawdb.NewMemoryDatabase()
	genesis := params.DefaultMessNetGenesisBlock()
	genesisB := MustCommitGenesis(db, triedb.NewDatabase(db, nil), genesis)

	chain, err := NewBlockChain(db, nil, genesis, nil, engine, vm.Config{}, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer chain.Stop()
	chain.EnableArtificialFinality(enableMess)

	easy, _ = GenerateChain(genesis.Config, genesisB, engine, db, easyL, func(i int, b *BlockGen) {
		b.SetNonce(types.EncodeNonce(uint64(rand.Int63n(math.MaxInt64))))
		b.OffsetTime(easyT)
	})
	commonAncestor := easy[caN-1]
	hard, _ = GenerateChain(genesis.Config, commonAncestor, engine, db, hardL, func(i int, b *BlockGen) {
		b.SetNonce(types.EncodeNonce(uint64(rand.Int63n(math.MaxInt64))))
		b.OffsetTime(hardT)
	})

	if _, err := chain.InsertChain(easy); err != nil {
		t.Fatal(err)
	}
	_, err = chain.InsertChain(hard)
	if err != nil {
		t.Logf("insert hard chain error = %v", err)
	}
	hardHead = chain.CurrentBlock().Hash() == hard[len(hard)-1].Hash()
	return
}

func TestBlockChain_AF_ECBP1100_2(t *testing.T) {
	offsetGreaterDifficulty := int64(-2) // 1..8 = -9..-2
	offsetSameDifficulty := int64(0)     // 9..17 = -1..8
	offsetWorseDifficulty := int64(8)    // 18..

	cases := []struct {
		easyLen, hardLen, commonAncestorN int
		easyOffset, hardOffset            int64
		hardGetsHead, accepted            bool
	}{
		// NOTE: Random coin tosses involved for equivalent difficulty.
		// Short trials for those are skipped.

		{
			1000, 30, 970,
			0, offsetSameDifficulty, // same difficulty
			false, true,
		},

		{
			1000, 1, 999,
			0, offsetWorseDifficulty, // worse! difficulty
			false, true,
		},
		{
			1000, 1, 999,
			0, offsetGreaterDifficulty, // better difficulty
			true, true,
		},
		{
			1000, 5, 995,
			0, offsetGreaterDifficulty,
			true, true,
		},
		{
			1000, 25, 975,
			0, offsetGreaterDifficulty,
			false, true,
		},
		{
			1000, 30, 970,
			0, offsetGreaterDifficulty,
			false, true,
		},
		{
			1000, 50, 950,
			0, offsetGreaterDifficulty,
			false, true,
		},
		{
			1000, 50, 950,
			0, offsetGreaterDifficulty,
			false, true,
		},
		{
			1000, 1000, 900,
			0, offsetGreaterDifficulty,
			true, true,
		},
		{
			1000, 2000, 800,
			0, offsetGreaterDifficulty,
			true, true,
		},
		{
			1000, 2000, 700,
			0, offsetGreaterDifficulty,
			true, true,
		},
		{
			1000, 2000, 700,
			0, offsetGreaterDifficulty,
			true, true,
		},
		{
			1000, 999, 1,
			0, offsetGreaterDifficulty,
			false, true,
		},
		{
			1000, 999, 1,
			0, offsetGreaterDifficulty,
			false, true,
		},
		{
			1000, 500, 500,
			0, offsetGreaterDifficulty,
			false, true,
		},
		{
			1000, 500, 500,
			0, offsetGreaterDifficulty,
			false, true,
		},
		{
			1000, 300, 700,
			0, offsetGreaterDifficulty,
			false, true,
		},
		{
			1000, 600, 700,
			0, offsetGreaterDifficulty,
			true, true,
		},
		// Will pass, takes a long time.
		// {
		// 	5000, 4000, 1000,
		// 	0, -2,
		// 	true, true,
		// },
	}

	for i, c := range cases {
		hardHead, err, hard, easy := runMESSTest2(t, true, c.easyLen, c.hardLen, c.commonAncestorN, c.easyOffset, c.hardOffset)

		ee, hh := easy[len(easy)-1], hard[len(hard)-1]
		rat, _ := new(big.Float).Quo(
			new(big.Float).SetInt(hh.Difficulty()),
			new(big.Float).SetInt(ee.Difficulty()),
		).Float64()

		logf := fmt.Sprintf("case=%d [easy=%d hard=%d ca=%d eo=%d ho=%d] drat=%0.6f span=%v hardHead(w|g)=%v|%v err=%v",
			i,
			c.easyLen, c.hardLen, c.commonAncestorN, c.easyOffset, c.hardOffset,
			rat,
			common.PrettyDuration(time.Second*time.Duration(10*(c.easyLen-c.commonAncestorN))),
			c.hardGetsHead, hardHead, err)

		if (err != nil && c.accepted) || (err == nil && !c.accepted) || (hardHead != c.hardGetsHead) {
			t.Error("FAIL", logf)
		} else {
			t.Log("PASS", logf)
		}
	}
}

/*
TestAFKnownBlock tests that AF functionality works for chain re-insertions.

Chain re-insertions use BlockChain.writeKnownBlockAsHead, where first-pass insertions
will hit writeBlockWithState.

AF needs to be implemented at both sites to prevent re-proposed chains from sidestepping
the AF criteria.
*/
func TestAFKnownBlock(t *testing.T) {
	engine := ethash.NewFaker()

	db := rawdb.NewMemoryDatabase()
	genesis := params.DefaultMessNetGenesisBlock()
	// genesis.Timestamp = 1
	genesisB := MustCommitGenesis(db, triedb.NewDatabase(db, nil), genesis)

	chain, err := NewBlockChain(db, nil, genesis, nil, engine, vm.Config{}, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer chain.Stop()
	chain.EnableArtificialFinality(true)

	easy, _ := GenerateChain(genesis.Config, genesisB, engine, db, 1000, func(i int, gen *BlockGen) {
		gen.OffsetTime(0)
	})
	easyN, err := chain.InsertChain(easy)
	if err != nil {
		t.Fatal(err)
	}
	hard, _ := GenerateChain(genesis.Config, easy[easyN-300], engine, db, 300, func(i int, gen *BlockGen) {
		gen.OffsetTime(-7)
	})
	// writeBlockWithState
	if _, err := chain.InsertChain(hard); err != nil {
		t.Error("hard 1 not inserted (should be side)")
	}
	// writeKnownBlockAsHead
	if _, err := chain.InsertChain(hard); err != nil {
		t.Error("hard 2 inserted (will have 'ignored' known blocks, and never tried a reorg)")
	}
	hardHeadHash := hard[len(hard)-1].Hash()
	if chain.CurrentBlock().Hash() == hardHeadHash {
		t.Fatal("hard block got chain head, should be side")
	}
	if h := chain.GetHeaderByHash(hardHeadHash); h == nil {
		t.Fatal("missing hard block (should be imported as side, but still available)")
	}
}

// TestEcbp1100PolynomialV tests the general shape and return values of the ECBP1100 polynomial curve.
// It makes sure domain values above the 'cap' do indeed get limited, as well
// as sanity check some normal domain values.
func TestEcbp1100PolynomialV(t *testing.T) {
	cases := []struct {
		block, ag int64
	}{
		{100, 1},
		{300, 2},
		{500, 5},
		{1000, 16},
		{2000, 31},
		{10000, 31},
		{1e9, 31},
	}
	for i, c := range cases {
		y := ecbp1100PolynomialV(big.NewInt(c.block * 13))
		y.Div(y, ecbp1100PolynomialVCurveFunctionDenominator)
		if c.ag != y.Int64() {
			t.Fatal("mismatch", i)
		}
	}
}

func TestPlot_ecbp1100PolynomialV(t *testing.T) {
	t.Skip("This test plots a graph of the ECBP1100 polynomial curve.")
	p := plot.New()
	p.Title.Text = "ECBP1100 Polynomial Curve Function"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	poly := plotter.NewFunction(func(f float64) float64 {
		n := big.NewInt(int64(f))
		y := ecbp1100PolynomialV(n)
		ff, _ := new(big.Float).SetInt(y).Float64()
		return ff
	})
	p.Add(poly)

	p.X.Min = 0
	p.X.Max = 30000
	p.Y.Min = 0
	p.Y.Max = 5000

	p.Y.Label.Text = "Antigravity imposition"
	p.X.Label.Text = "Seconds difference between local head and proposed common ancestor"

	if err := p.Save(1000, 1000, "ecbp1100-polynomial.png"); err != nil {
		t.Fatal(err)
	}
}

func TestEcbp1100AGSinusoidalA(t *testing.T) {
	cases := []struct {
		in, out float64
	}{
		{0, 1},
		{25132, 31},
	}
	tolerance := 0.0000001
	for i, c := range cases {
		if got := ecbp1100AGSinusoidalA(c.in); got < c.out-tolerance || got > c.out+tolerance {
			t.Fatalf("%d: in: %0.6f want: %0.6f got: %0.6f", i, c.in, c.out, got)
		}
	}
}

func TestDifficultyDelta(t *testing.T) {
	t.Skip("A development test to play with difficulty steps.")
	parent := &types.Header{
		Number:     big.NewInt(1_000_000),
		Difficulty: params.DefaultMessNetGenesisBlock().Difficulty,
		Time:       uint64(time.Now().Unix()),
		UncleHash:  types.EmptyUncleHash,
	}

	data := plotter.XYs{}

	for i := uint64(1); i <= 60; i++ {
		nextTime := parent.Time + i
		d := ethash.CalcDifficulty(params.MessNetConfig, nextTime, parent)

		rat, _ := new(big.Float).Quo(
			new(big.Float).SetInt(d),
			new(big.Float).SetInt(parent.Difficulty),
		).Float64()

		t.Log(i, rat)
		data = append(data, plotter.XY{X: float64(i), Y: rat})
	}

	p := plot.New()
	p.Title.Text = "Block Difficulty Delta by Timestamp Offset"
	p.X.Label.Text = "Timestamp Offset"
	p.Y.Label.Text = "Relative Difficulty (child/parent)"

	dataScatter, _ := plotter.NewScatter(data)
	p.Add(dataScatter)

	if err := p.Save(800, 600, "difficulty-adjustments.png"); err != nil {
		t.Fatal(err)
	}
}

func TestGenerateChainTargetingHashrate(t *testing.T) {
	t.Skip("A development test to play with difficulty steps.")
	engine := ethash.NewFaker()

	db := rawdb.NewMemoryDatabase()
	genesis := params.DefaultMessNetGenesisBlock()
	// genesis.Timestamp = 1
	genesisB := MustCommitGenesis(db, triedb.NewDatabase(db, nil), genesis)

	chain, err := NewBlockChain(db, nil, genesis, nil, engine, vm.Config{}, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer chain.Stop()
	chain.EnableArtificialFinality(true)

	easy, _ := GenerateChain(genesis.Config, genesisB, engine, db, 1000, func(i int, gen *BlockGen) {
		gen.OffsetTime(0)
	})
	if _, err := chain.InsertChain(easy); err != nil {
		t.Fatal(err)
	}

	baseDifficulty := chain.CurrentHeader().Difficulty
	targetDifficultyRatio := big.NewInt(4)
	targetDifficulty := new(big.Int).Mul(baseDifficulty, targetDifficultyRatio)

	data := plotter.XYs{}

	for chain.CurrentHeader().Difficulty.Cmp(targetDifficulty) < 0 {
		bl := chain.GetBlock(chain.CurrentHeader().Hash(), chain.CurrentHeader().Number.Uint64())
		next, _ := GenerateChain(genesis.Config, bl, engine, db, 1, func(i int, gen *BlockGen) {
			gen.OffsetTime(-9) // 8: (=10+8=18>(13+4=17).. // minimum value over stable range
		})
		if _, err := chain.InsertChain(next); err != nil {
			t.Fatal(err)
		}

		// f, _ := new(big.Float).SetInt(next[0].Difficulty()).Float64()
		// data = append(data, plotter.XY{X: float64(next[0].NumberU64()), Y: f})

		rat1, _ := new(big.Float).Quo(
			new(big.Float).SetInt(next[0].Difficulty()),
			new(big.Float).SetInt(targetDifficulty),
		).Float64()

		// rat, _ := new(big.Float).Quo(
		// 	new(big.Float).SetInt(next[0].Difficulty()),
		// 	new(big.Float).SetInt(targetDifficultyRatio),
		// ).Float64()

		data = append(data, plotter.XY{X: float64(next[0].NumberU64()), Y: rat1})
	}
	t.Log(chain.CurrentBlock().Number)

	p := plot.New()
	p.Title.Text = fmt.Sprintf("Block Difficulty Toward Target: %dx", targetDifficultyRatio.Uint64())
	p.X.Label.Text = "Block Number"
	p.Y.Label.Text = "Difficulty"

	dataScatter, _ := plotter.NewScatter(data)
	p.Add(dataScatter)

	if err := p.Save(800, 600, "difficulty-toward-target.png"); err != nil {
		t.Fatal(err)
	}
}

func runMESSTest(t *testing.T, easyL, hardL, caN int, easyT, hardT int64) (hardHead bool, err error) {
	// Generate the original common chain segment and the two competing forks
	engine := ethash.NewFaker()

	db := rawdb.NewMemoryDatabase()
	genesis := params.DefaultMessNetGenesisBlock()
	genesisB := MustCommitGenesis(db, triedb.NewDatabase(db, nil), genesis)

	chain, err := NewBlockChain(db, nil, genesis, nil, engine, vm.Config{}, nil, nil)
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

var yuckyGlobalTestEnableMess = false

func TestBlockChain_GenerateMESSPlot(t *testing.T) {
	t.Skip("This test plots graph of chain acceptance for visualization.")
	easyLen := 500
	maxHardLen := 400

	generatePlot := func(title, fileName string) {
		p := plot.New()
		p.Title.Text = title
		p.X.Label.Text = "Block Depth"
		p.Y.Label.Text = "Mode Block Time Offset (10 seconds + y)"

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

		err := p.Save(pixelWidth, 300, fileName)
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
	// generatePlot("WITHOUT MESS: "+baseTitle, "reorgs-noMESS.png")
}

func TestBlockChain_AF_ECBP1100(t *testing.T) {
	t.Skip("These have been disused as of the sinusoidal -> cubic change.")
	yuckyGlobalTestEnableMess = true
	defer func() {
		yuckyGlobalTestEnableMess = false
	}()

	cases := []struct {
		easyLen, hardLen, commonAncestorN int
		easyOffset, hardOffset            int64
		hardGetsHead, accepted            bool
	}{
		// INDEX=0
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
		// INDEX=5
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
		// INDEX=10
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
		// // INDEX=15
		// Hard has insufficient total difficulty / length and is rejected.
		{
			500, 250, 250,
			0, -9,
			false, false,
		},
		// Hard has insufficient total difficulty / length and is rejected.
		{
			500, 250, 250,
			7, -9,
			false, false,
		},
		// Hard has insufficient total difficulty / length and is rejected.
		{
			500, 300, 200,
			13, -9,
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
		// // INDEX=20
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
		// INDEX=25
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

func TestBlockChain_AF_Difficulty_Develop(t *testing.T) {
	t.Skip("Development version of tests with plotter")
	// Generate the original common chain segment and the two competing forks
	engine := ethash.NewFaker()

	db := rawdb.NewMemoryDatabase()
	genesis := params.DefaultMessNetGenesisBlock()
	// genesis.Timestamp = 1
	genesisB := MustCommitGenesis(db, triedb.NewDatabase(db, nil), genesis)

	chain, err := NewBlockChain(db, nil, genesis, nil, engine, vm.Config{}, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer chain.Stop()
	chain.EnableArtificialFinality(true)

	cases := []struct {
		easyLen, hardLen, commonAncestorN int
		easyOffset, hardOffset            int64
		hardGetsHead, accepted            bool
	}{
		// {
		// 	1000, 800, 200,
		// 	10, 1,
		// 	true, true,
		// },
		// {
		// 	1000, 800, 200,
		// 	60, 1,
		// 	true, true,
		// },
		// {
		// 	10000, 8000, 2000,
		// 	60, 1,
		// 	true, true,
		// },
		// {
		// 	20000, 18000, 2000,
		// 	10, 1,
		// 	true, true,
		// },
		// {
		// 	20000, 18000, 2000,
		// 	60, 1,
		// 	true, true,
		// },
		// {
		// 	10000, 8000, 2000,
		// 	10, 20,
		// 	true, true,
		// },

		// {
		// 	1000, 1, 999,
		// 	10, 1,
		// 	true, true,
		// },
		// {
		// 	1000, 10, 990,
		// 	10, 1,
		// 	true, true,
		// },
		// {
		// 	1000, 100, 900,
		// 	10, 1,
		// 	true, true,
		// },
		// {
		// 	1000, 200, 800,
		// 	10, 1,
		// 	true, true,
		// },
		// {
		// 	1000, 500, 500,
		// 	10, 1,
		// 	true, true,
		// },
		// {
		// 	1000, 999, 1,
		// 	10, 1,
		// 	true, true,
		// },
		// {
		// 	5000, 4000, 1000,
		// 	10, 1,
		// 	true, true,
		// },

		// {
		// 	10000, 9000, 1000,
		// 	10, 1,
		// 	true, true,
		// },
		//
		// {
		// 	7000, 6500, 500,
		// 	10, 1,
		// 	true, true,
		// },

		// {
		// 	100, 90, 10,
		// 	10, 1,
		// 	true, true,
		// },

		// {
		// 	1000, 1, 999,
		// 	10, 1,
		// 	true, true,
		// },
		// {
		// 	1000, 2, 998,
		// 	10, 1,
		// 	true, true,
		// },
		// {
		// 	1000, 3, 997,
		// 	10, 1,
		// 	true, true,
		// },
		// {
		// 	1000, 1, 999,
		// 	10, 8,
		// 	true, true,
		// },

		{
			1000, 50, 950,
			10, 9,
			false, false,
		},
		{
			1000, 100, 900,
			10, 8,
			false, false,
		},
		{
			1000, 100, 900,
			10, 7,
			false, false,
		},
		{
			1000, 50, 950,
			10, 5,
			true, true,
		},
		{
			1000, 50, 950,
			10, 3,
			true, true,
		},
		// 5
		{
			1000, 100, 900,
			10, 3,
			false, false,
		},
		{
			1000, 200, 800,
			10, 3,
			false, false,
		},
		{
			1000, 200, 800,
			10, 1,
			false, false,
		},
	}

	// poissonTime := func(b *BlockGen, seconds int64) {
	// 	poisson := distuv.Poisson{Lambda: float64(seconds)}
	// 	r := poisson.Rand()
	// 	if r < 1 {
	// 		r = 1
	// 	}
	// 	if r > float64(seconds) * 1.5 {
	// 		r = float64(seconds)
	// 	}
	// 	chainreader := &fakeChainReader{config: b.config}
	// 	b.header.Time = b.parent.Time() + uint64(r)
	// 	b.header.Difficulty = b.engine.CalcDifficulty(chainreader, b.header.Time, b.parent.Header())
	// 	for err := b.engine.VerifyHeader(chainreader, b.header, false);
	// 		err != nil && err != consensus.ErrUnknownAncestor && b.header.Time > b.parent.Header().Time; {
	// 		t.Log(err)
	// 		r -= 1
	// 		b.header.Time = b.parent.Time() + uint64(r)
	// 		b.header.Difficulty = b.engine.CalcDifficulty(chainreader, b.header.Time, b.parent.Header())
	// 	}
	// }

	// nolint:unused
	type ratioComparison struct {
		tdRatio float64
		penalty float64
	}
	gotRatioComparisons := []ratioComparison{}

	for i, c := range cases {
		if err := chain.Reset(); err != nil {
			t.Fatal(err)
		}
		easy, _ := GenerateChain(genesis.Config, genesisB, engine, db, c.easyLen, func(i int, b *BlockGen) {
			b.SetNonce(types.EncodeNonce(uint64(rand.Int63n(math.MaxInt64))))
			// poissonTime(b, c.easyOffset)
			b.OffsetTime(c.easyOffset - 10)
		})
		commonAncestor := easy[c.commonAncestorN-1]
		hard, _ := GenerateChain(genesis.Config, commonAncestor, engine, db, c.hardLen, func(i int, b *BlockGen) {
			b.SetNonce(types.EncodeNonce(uint64(rand.Int63n(math.MaxInt64))))
			// poissonTime(b, c.hardOffset)
			b.OffsetTime(c.hardOffset - 10)
		})
		if _, err := chain.InsertChain(easy); err != nil {
			t.Fatal(err)
		}
		n, err := chain.InsertChain(hard)
		hardHead := chain.CurrentBlock().Hash() == hard[len(hard)-1].Hash()

		commons := plotter.XYs{}
		easys := plotter.XYs{}
		hards := plotter.XYs{}
		tdrs := plotter.XYs{}
		antigravities := plotter.XYs{}
		antigravities2 := plotter.XYs{}

		balance := plotter.XYs{}

		for i := 0; i < c.easyLen; i++ {
			td := chain.GetTd(easy[i].Hash(), easy[i].NumberU64())
			point := plotter.XY{X: float64(easy[i].NumberU64()), Y: float64(td.Uint64())}
			if i <= c.commonAncestorN {
				commons = append(commons, point)
			} else {
				easys = append(easys, point)
			}
		}
		// td ratios
		// for j := 0; j < c.hardLen; j++ {
		for j := 0; j < n; j++ {
			td := chain.GetTd(hard[j].Hash(), hard[j].NumberU64())
			if td != nil {
				point := plotter.XY{X: float64(hard[j].NumberU64()), Y: float64(td.Uint64())}
				hards = append(hards, point)
			}

			if commonAncestor.NumberU64() != uint64(c.commonAncestorN) {
				t.Fatalf("bad test common=%d easy=%d can=%d", commonAncestor.NumberU64(), c.easyLen, c.commonAncestorN)
			}

			ee := c.commonAncestorN + j
			easyHeader := easy[ee].Header()
			hardHeader := hard[j].Header()
			if easyHeader.Number.Uint64() != hardHeader.Number.Uint64() {
				t.Fatalf("bad test easyheader=%d hardheader=%d", easyHeader.Number.Uint64(), hardHeader.Number.Uint64())
			}

			/*
				HERE LIES THE RUB (IN MY GRAPHS).


			*/
			// y := chain.getTDRatio(commonAncestor.Header(), easyHeader, hardHeader) // <- unit x unit

			// y := chain.getTDRatio(commonAncestor.Header(), easy[c.easyLen-1].Header(), hardHeader)

			y := chain.getTDRatio(commonAncestor.Header(), chain.CurrentHeader(), hardHeader)

			if j == 0 {
				t.Logf("case=%d first.hard.tdr=%v", i, y)
			}

			ecbp := ecbp1100AGSinusoidalA(float64(hardHeader.Time - commonAncestor.Header().Time))

			if j == n-1 {
				gotRatioComparisons = append(gotRatioComparisons, ratioComparison{
					tdRatio: y, penalty: ecbp,
				})
			}

			// Exploring alternative penalty functions.
			ecbp2 := ecbp1100AGExpA(float64(hardHeader.Time - commonAncestor.Header().Time))
			// t.Log(y, ecbp, ecbp2)

			tdrs = append(tdrs, plotter.XY{X: float64(hard[j].NumberU64()), Y: y})
			antigravities = append(antigravities, plotter.XY{X: float64(hard[j].NumberU64()), Y: ecbp})
			antigravities2 = append(antigravities2, plotter.XY{X: float64(hard[j].NumberU64()), Y: ecbp2})

			balance = append(balance, plotter.XY{X: float64(hardHeader.Number.Uint64()), Y: y - ecbp})
		}
		scatterCommons, _ := plotter.NewScatter(commons)
		scatterEasys, _ := plotter.NewScatter(easys)
		scatterHards, _ := plotter.NewScatter(hards)

		scatterTDRs, _ := plotter.NewScatter(tdrs)
		scatterAntigravities, _ := plotter.NewScatter(antigravities)
		scatterAntigravities2, _ := plotter.NewScatter(antigravities2)
		balanceScatter, _ := plotter.NewScatter(balance)

		scatterCommons.Color = color.RGBA{R: 190, G: 197, B: 236, A: 255}
		scatterCommons.Shape = draw.CircleGlyph{}
		scatterCommons.Radius = 2
		scatterEasys.Color = color.RGBA{R: 152, G: 236, B: 161, A: 255} // green
		scatterEasys.Shape = draw.CircleGlyph{}
		scatterEasys.Radius = 2
		scatterHards.Color = color.RGBA{R: 236, G: 106, B: 94, A: 255}
		scatterHards.Shape = draw.CircleGlyph{}
		scatterHards.Radius = 2

		p := plot.New()
		p.Add(scatterCommons)
		p.Legend.Add("Commons", scatterCommons)
		p.Add(scatterEasys)
		p.Legend.Add("Easys", scatterEasys)
		p.Add(scatterHards)
		p.Legend.Add("Hards", scatterHards)
		p.Title.Text = fmt.Sprintf("TD easy=%d hard=%d", c.easyOffset, c.hardOffset)
		p.Save(1000, 600, fmt.Sprintf("plot-td-%d-%d-%d-%d-%d.png", c.easyLen, c.commonAncestorN, c.hardLen, c.easyOffset, c.hardOffset))

		p = plot.New()

		scatterTDRs.Color = color.RGBA{R: 236, G: 106, B: 94, A: 255} // red
		scatterTDRs.Radius = 3
		scatterTDRs.Shape = draw.PyramidGlyph{}
		p.Add(scatterTDRs)
		p.Legend.Add("TD Ratio", scatterTDRs)

		scatterAntigravities.Color = color.RGBA{R: 190, G: 197, B: 236, A: 255} // blue
		scatterAntigravities.Radius = 3
		scatterAntigravities.Shape = draw.PlusGlyph{}
		p.Add(scatterAntigravities)
		p.Legend.Add("(Anti)Gravity Penalty", scatterAntigravities)

		scatterAntigravities2.Color = color.RGBA{R: 152, G: 236, B: 161, A: 255} // green
		scatterAntigravities2.Radius = 3
		scatterAntigravities2.Shape = draw.PlusGlyph{}
		// p.Add(scatterAntigravities2)
		// p.Legend.Add("(Anti)Gravity Penalty (Alternate)", scatterAntigravities2)

		p.Title.Text = fmt.Sprintf("TD Ratio easy=%d hard=%d", c.easyOffset, c.hardOffset)
		p.Save(1000, 600, fmt.Sprintf("plot-td-ratio-%d-%d-%d-%d-%d.png", c.easyLen, c.commonAncestorN, c.hardLen, c.easyOffset, c.hardOffset))

		p = plot.New()
		p.Title.Text = fmt.Sprintf("TD Ratio - Antigravity Penalty easy=%d hard=%d", c.easyOffset, c.hardOffset)
		balanceScatter.Color = color.RGBA{R: 235, G: 92, B: 236, A: 255} // purple
		balanceScatter.Radius = 3
		balanceScatter.Shape = draw.PlusGlyph{}
		p.Add(balanceScatter)
		p.Legend.Add("TDR - Penalty", balanceScatter)
		p.Save(1000, 600, fmt.Sprintf("plot-td-ratio-diff-%d-%d-%d-%d-%d.png", c.easyLen, c.commonAncestorN, c.hardLen, c.easyOffset, c.hardOffset))

		if (err != nil && c.accepted) || (err == nil && !c.accepted) || (hardHead != c.hardGetsHead) {
			compared := gotRatioComparisons[i]
			t.Errorf(`case=%d [easy=%d hard=%d ca=%d eo=%d ho=%d] want.accepted=%v want.hardHead=%v got.hardHead=%v err=%v
got.tdr=%v got.pen=%v`,
				i,
				c.easyLen, c.hardLen, c.commonAncestorN, c.easyOffset, c.hardOffset,
				c.accepted, c.hardGetsHead, hardHead, err, compared.tdRatio, compared.penalty)
		}
	}
}
