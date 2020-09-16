package cmd

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/p2p/enode"
)

type Node struct {
	Name      string `json:"name"`
	HeadHash  string `json:"headHash"`
	HeadNum   uint64 `json:"headNum"`
	HeadBase  string `json:"headBase"`
	HeadMiner bool   `json:"headMiner"`
	HeadD     uint64 `json:"headD"` // Difficulty
	HeadTD    uint64 `json:"headTD"`
}

func (n *Node) UnmarshalJSON(bytes []byte) error {
	var im = struct {
		ID   string `json:"id"`
		Name string `json:"name"`

		HeadHash string `json:"headHash"`
		Hash     string `json:"hash"`

		HeadNum uint64 `json:"headNum"`
		Number  uint64 `json:"number"`

		HeadBase  string `json:"headBase"`
		Etherbase string `json:"etherbase"`

		HeadMiner bool `json:"headMiner"`

		HeadD uint64 `json:"headD` // Difficulty

		HeadTD uint64 `json:"headTD"`
		TD     uint64 `json:"td"`
	}{}
	err := json.Unmarshal(bytes, &im)
	if err != nil {
		return err
	}

	n.Name = im.Name
	if im.Name != "" {
		n.Name = im.ID
	}
	n.HeadHash = im.HeadHash
	if im.Hash != "" {
		n.HeadHash = im.Hash
	}

	n.HeadNum = im.HeadNum
	if im.Number != 0 {
		n.HeadNum = im.Number
	}
	n.HeadBase = im.HeadBase
	if im.Etherbase != "" {
		n.HeadBase = im.Etherbase
	}

	if !n.HeadMiner {
		// Python script uses etherbase as id.
		n.HeadMiner = n.HeadBase == n.Name
	}

	n.HeadD = im.HeadD

	n.HeadTD = im.HeadTD
	if im.TD != 0 {
		n.HeadTD = im.TD
	}
	return nil
}

type Link struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

type NetworkGraphData struct {
	Tick  int    `json:"tick,omitempty"`
	Nodes []Node `json:"nodes"`
	Links []Link `json:"links"`
}

func (n *NetworkGraphData) UnmarshalJSON(bytes []byte) error {
	var im = struct {
		Tick   int    `json:"tick,omitempty"`
		Nodes  []Node `json:"nodes"`
		Agents []Node `json:"agents"`
		Links  []Link `json:"links"`
	}{}
	err := json.Unmarshal(bytes, &im)
	if err != nil {
		return err
	}
	n.Tick = im.Tick
	n.Nodes = im.Nodes
	if len(im.Agents) != 0 {
		n.Nodes = im.Agents
	}
	n.Links = im.Links
	return nil
}

// func getWorldView(set *agethSet) NetworkGraphData {
// 	nodes := []Node{}
// 	links := []Link{}
// 	for _, l := range set.ageths {
// 		if l.name == "" {
// 			log.Error("empty name?!", "ageth", l)
// 			continue
// 		}
// 		nodes = append(nodes, Node{
// 			Name:      l.name,
// 			HeadHash:  l.block().hash.Hex()[2:8],
// 			HeadNum:   l.block().number,
// 			HeadBase:  l.block().coinbase.Hex()[2:8],
// 			HeadMiner: l.block().coinbase == l.coinbase,
// 			HeadD:     l.block().difficulty,
// 			HeadTD:    l.td.Uint64(),
// 		})
// 		// make a map of recorded connections; we
// 		// only need to record peers in one direction
// 		unilateralMap := map[string]bool{}
// 		for _, p := range l.peers.all() {
// 			// check if connection already recorded
// 			if _, ok := unilateralMap[p.name+l.name]; ok {
// 				continue
// 			}
// 			unilateralMap[l.name+p.name] = true
// 			links = append(links, Link{
// 				Source: l.name,
// 				Target: p.name,
// 			})
// 		}
// 	}
// 	return NetworkGraphData{Nodes: nodes, Links: links}
// }

func getWorldView(set *agethSet) NetworkGraphData {
	nodes := []Node{}
	links := []Link{}
	for _, l := range set.ageths {
		if l.name == "" {
			log.Error("empty name?!", "ageth", l)
			continue
		}
		name := enode.MustParse(l.enode).IP().String()
		nodes = append(nodes, Node{
			// Name:      l.name,
			Name:      name,
			HeadHash:  l.block().hash.Hex()[2:8],
			HeadNum:   l.block().number,
			HeadBase:  l.block().coinbase.Hex()[2:8],
			HeadMiner: l.block().coinbase == l.coinbase,
			HeadD:     l.block().difficulty,
			HeadTD:    l.td.Uint64(),
		})
		// make a map of recorded connections; we
		// only need to record peers in one direction
		unilateralMap := map[string]bool{}
		for _, p := range l.peers.all() {
			pname := enode.MustParse(p.enode).IP().String()
			// check if connection already recorded
			if _, ok := unilateralMap[pname+name]; ok {
				continue
			}
			unilateralMap[name+pname] = true
			links = append(links, Link{
				Source: name,
				Target: pname,
			})
		}
	}
	return NetworkGraphData{Nodes: nodes, Links: links}
}

/*



 */
