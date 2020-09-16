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

func stabilize(nodes *agethSet) {
  for len(nodes.distinctChains()) > 1 {
    time.Sleep(30)
  }
}

func scenarioGenerator(blockTime int, attackBlocks uint64, difficultyRatio float64) func(*agethSet) {
  return func(nodes *agethSet) {
    // Setup

    // Start all nodes mining at 150% of the blocktime. They will be the long tail of small miners.
    for _, node := range nodes.all() {
      node.startMining(blockTime * 3 / 2)
    }
    bigMiners := newAgethSet()
    badGuy := nodes.indexed(0) // NOTE: Assumes badguy will always be [0]

    // Simulate the small proportion of "whale" miners, whose block times will be nearer the target time.
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
    blockNumber := nodes.headMax() // grab current block number after 30s spin up
    for {
      // Allow the chain to have mined at least 5 blocks since.
      if nodes.headMax() > blockNumber + 5 {
        break
      }
      time.Sleep(1 * time.Second)
    }
    log.Info("Starting attacker")
    // badGuy.setPeerCount(0)
    resumePeering := badGuy.refusePeers(10)
    forkBlock := badGuy.block()
    badGuy.startMining(blockTime / 2)

    // Once a second, check to see if the bad guy's block difficulty has
    // reached the target.
    for {
      if badGuy.block().difficulty > uint64(float64(forkBlock.difficulty) * difficultyRatio) {
        log.Info("Attacker reached target relative difficulty ratio", "target ratio", difficultyRatio)
        break
      }
      if badGuy.block().number > forkBlock.number + attackBlocks {
        log.Info("Attacker mined attack blocks allowance", "blocks", attackBlocks)
        break
      }
      time.Sleep(time.Second)
    }

    // The target difficulty or chain length has been reached.
    // If target difficulty was reached before the desired number of attack blocks,
    // mining at blockTime will keep it roughly in place until the desired number of attack blocks is produced.
    for badGuy.block().number < forkBlock.number + attackBlocks {
      if !badGuy.isMining {
        badGuy.startMining(blockTime)
      }
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
    if err != nil { log.Error("Error marshalling final report", "err", err) }
    os.Stdout.Write(data)
    os.Stdout.WriteString("\n")
    badGuy.truncateHead(forkBlock.number) // Reset badGuy to the fork block
  }
}
