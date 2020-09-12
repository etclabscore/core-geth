package cmd

import (
	"math/rand"

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
