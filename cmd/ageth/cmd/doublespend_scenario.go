package cmd

import (
  "context"
  "encoding/json"
  "math/big"
  "os"
  "time"
  "github.com/ethereum/go-ethereum/common"
  log "github.com/ethereum/go-ethereum/log"
)


type finalReport struct {
  Converted int
  Unconverted int
  DistinctChains []common.Hash
  Nodes map[string]common.Hash
}


func scenarioGenerator(blockTime int, attackBlocks uint64, difficultyRatio float64) func(*agethSet) error {
  return func(nodes *agethSet) error {
    // Setup
    for _, node := range nodes.all() {
      node.startMining(blockTime * 3 / 2)
    }
    bigMiners := newAgethSet()
    badGuy := nodes.all()[0] // NOTE: Assumes badguy will always be [0]
    hashtimes := []int{blockTime, blockTime * 12/10, blockTime * 12/10, blockTime * 14/10, blockTime * 14/10}
    for _, hashtime := range hashtimes {
      nextMiner := nodes.random()
      for nextMiner.name == badGuy.name || bigMiners.contains(nextMiner) {
        nextMiner = nodes.random()
      }
      nextMiner.startMining(hashtime)
    }
    log.Info("Started miners")
    time.Sleep(30 * time.Second)
    blockNumber := nodes.headMax()
    for {
      if nodes.headMax() > blockNumber + 5 {
        break
      }
      time.Sleep(1 * time.Second)
    }
    log.Info("Starting attacker")
    // badGuy.setPeerCount(0)
    resumePeering := badGuy.refusePeers()
    forkBlock := badGuy.block()
    badGuy.startMining(blockTime / 2)

    // Once a second, check to see if the bad guy's block difficulty has
    // reached the target
    for uint64(float64(forkBlock.difficulty) * difficultyRatio) >= badGuy.block().difficulty && badGuy.block().number < forkBlock.number + attackBlocks {
      time.Sleep(1 * time.Second)
    }
    log.Info("Attacker reached target difficulty")
    // The target difficulty has been reached. Mining at blockTime will keep it
    // roughly in place.
    if badGuy.block().number < forkBlock.number + attackBlocks {
      badGuy.startMining(blockTime)
    }
    for badGuy.block().number < forkBlock.number + attackBlocks {
      time.Sleep(1 * time.Second)
    }
    log.Info("Attacker mined %v blocks", attackBlocks)
    finalAttackBlock := badGuy.block()
    badGuy.stopMining()
    // badGuy.setPeerCount(25)
    resumePeering()
    for badGuy.getPeerCount() < 5 {
      // Sleep until the badguy has found 5 peers
      time.Sleep(1 * time.Second)
    }

    // Sleep another 30 seconds for blocks to propagate
    time.Sleep(30 * time.Second)

    b, err := badGuy.eclient.BlockByNumber(context.Background(), big.NewInt(int64(finalAttackBlock.number)))
    if err != nil {
      log.Error("Error getting block", "blockno", finalAttackBlock.number, "node", badGuy.name)
    }
    if b.Hash() != finalAttackBlock.hash {
      log.Info("Even the attacker gave up their block")
    }

    convertedNodes := 0
    unconvertedNodes := 0
    for _, node := range nodes.all()[1:] {
      b, err := node.eclient.BlockByNumber(context.Background(), big.NewInt(int64(finalAttackBlock.number)))
      if err != nil {
        log.Error("Error getting block", "blockno", finalAttackBlock.number, "node", node.name)
        continue
      }
      if b.Hash() == finalAttackBlock.hash {
        convertedNodes++
      } else {
        unconvertedNodes++
      }
    }
    report := &finalReport{
      Converted: convertedNodes,
      Unconverted: unconvertedNodes,
      DistinctChains: nodes.distinctChains(),
      Nodes: make(map[string]common.Hash),
    }
    for _, node := range nodes.all() {
      report.Nodes[node.name] = node.block().hash
    }
    data, err := json.Marshal(report)
    if err != nil { return err }
    os.Stdout.Write(data)
    os.Stdout.WriteString("\n")
    return nil
  }
}
