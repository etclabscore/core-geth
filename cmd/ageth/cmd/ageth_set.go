package cmd

import (
	"context"
	"math/rand"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/montanaflynn/stats"
)

type agethSet struct {
	ageths []*ageth
}

func newAgethSet() *agethSet {
	return &agethSet{
		ageths: []*ageth{},
	}
}

func (s *agethSet) get(name string) *ageth {
	for _, g := range s.ageths {
		if g.name == name {
			return g
		}
	}
	return nil
}

func (a *agethSet) ready() bool {
	for _, g := range a.all() {
		if !g.online {
			return false
		}
	}
	return true
}

func (s *agethSet) contains(a *ageth) bool {
	for _, g := range s.ageths {
		if g.name == a.name {
			return true
		}
	}
	return false
}

func (s *agethSet) all() []*ageth {
	ret := []*ageth{}
	for _, a := range s.ageths {
		if a != nil {
			ret = append(ret, a)
		}
	}
	return ret
}

func (s *agethSet) len() int {
	return len(s.ageths)
}

func (s *agethSet) push(a *ageth) bool {
	if a == nil || s.contains(a) {
		return false
	}
	s.ageths = append(s.ageths, a)
	return true
}

func (s *agethSet) union(s2 *agethSet) *agethSet {
	ns := &agethSet{ ageths: s.ageths[:] }
	for _, a := range s2.all() { ns.push(a) }
	return ns
}

func (s *agethSet) remove(a *ageth) (ok bool) {
	for i, g := range s.ageths {
		if g.name == a.name {
			ok = true
			if i < len(s.ageths)-1 {
				s.ageths = append(s.ageths[:i], s.ageths[i+1:]...)
			} else {
				s.ageths = s.ageths[:i]
			}
		}
	}
	return ok
}

func (s *agethSet) random() *ageth {
	if len(s.ageths) == 0 {
		return nil
	}
	return s.ageths[rand.Intn(len(s.ageths))]
}

func (s *agethSet) randomN(rate float64) *agethSet {
	ret := newAgethSet()
	for _, a := range s.all() {
		if rand.Float64() < rate {
			ret.push(a)
		}
	}
	return ret
}

func (s *agethSet) headMax() uint64 {
	m := uint64(0)
	for _, g := range s.ageths {
		if g.block().number > m {
			m = g.block().number
		}
	}
	return m
}

func (s *agethSet) headMode() uint64 {
	coll := []float64{}
	for _, g := range s.ageths {
		coll = append(coll, float64(g.block().number))
	}
	mode, err := stats.Mode(coll)
	if err != nil {
		log.Warn("Ageth set head mode errored", "error", err)
		return 0
	}
	if len(mode) == 0 {
		return 0
	}
	return uint64(mode[0])
}

func (s *agethSet) headBlock() *types.Block {
	greatestN := s.headMax()
	for _, g := range s.ageths {
		if g.block().number == greatestN {
			return g.latestBlock
		}
	}
	return nil
}

func (s *agethSet) distinctChains() []common.Hash {
	headNodeMap := make(map[common.Hash]*agethSet)
	for _, a := range s.ageths {
		head := a.block()
		if _, ok := headNodeMap[head.hash]; !ok {
			headNodeMap[head.hash] = newAgethSet()
		}
		headNodeMap[head.hash].push(a)
	}
	knownHeads := make([]common.Hash, len(headNodeMap))
	i := 0
	for hash := range headNodeMap {
		knownHeads[i] = hash
		i++
	}
	distinctHeads := map[common.Hash]struct{}{}
	for i, hasha := range knownHeads {
		if i == len(knownHeads) - 1 { break }
		for _, hashb := range knownHeads[i+1:] {
			isReorg, latest := isReorg(hasha, hashb, headNodeMap[hasha], headNodeMap[hashb])
			if isReorg {
				distinctHeads[hasha] = struct{}{}
				distinctHeads[hashb] = struct{}{}
			} else {
				distinctHeads[latest] = struct{}{}
			}
		}
	}
	distinctHeadList := make([]common.Hash, len(distinctHeads))
	j := 0
	for hash := range distinctHeads {
		distinctHeadList[j] = hash
		j++
	}
	return distinctHeadList
}


func isReorg(a, b common.Hash, aSet, bSet *agethSet) (bool, common.Hash) {
	if a == b { return false, a }
	blocka, err := aSet.random().eclient.BlockByHash(context.Background(), a)
	for err != nil { blocka, err = aSet.random().eclient.BlockByHash(context.Background(), a) }
	blockb, err := bSet.random().eclient.BlockByHash(context.Background(), b)
	for err != nil { blockb, err = bSet.random().eclient.BlockByHash(context.Background(), b) }

	if blocka.NumberU64() == blockb.NumberU64() { return true, blocka.Hash() }

	var higherBlock, lowerBlock *types.Block
	var higherSet *agethSet

	if blocka.NumberU64() > blockb.NumberU64() {
		higherBlock = blocka
		lowerBlock = blockb
		higherSet = aSet
	} else {
		higherBlock = blockb
		lowerBlock = blocka
		higherSet = bSet
	}
	block, err := higherSet.random().eclient.BlockByNumber(context.Background(), lowerBlock.Number())
	for err != nil { block, err = higherSet.random().eclient.BlockByNumber(context.Background(), lowerBlock.Number()) }
	return block.Hash() != lowerBlock.Hash(), higherBlock.Hash()
}

func (s *agethSet) indexed(i int) *ageth {
	if i >= len(s.all()) {
		return nil
	}
	return s.all()[i]
}

func (s *agethSet) subset(inclusiveStartIndex, nonInclusiveEndIndex int) *agethSet {
	newSet := newAgethSet()
	for i := inclusiveStartIndex; i < nonInclusiveEndIndex; i++ {
		newSet.push(s.indexed(i))
	}
	return newSet
}

func (s *agethSet) where(cond func(g *ageth) bool) *agethSet {
	ret := newAgethSet()
	for _, a := range s.ageths {
		if cond(a) {
			ret.push(a)
		}
	}
	return ret
}
