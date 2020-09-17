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
  AttackerShouldWin bool
  AttackerWon bool
}

func stabilize(nodes *agethSet) {
  badGuy := nodes.indexed(0) // NOTE: Assumes badguy will always be [0]
  goodGuys := nodes.where(func(a *ageth) bool { return a.name != badGuy.name })
  badGuy.truncateHead(goodGuys.headMax())
  minimumPeerCount := int64(2)
  for _, node := range nodes.all() {
    node.stopMining()
    var result interface{}
    node.client.Call(&result, "admin_maxPeers", 20)
    for node.getPeerCount() < minimumPeerCount {
      node.addPeer(nodes.random())
    }
  }
  goodGuys.random().startMining(30) // Start single miner at a slow pace.
  for len(nodes.distinctChains()) > 1 {
    time.Sleep(30)
  }
}

func scenarioGenerator(blockTime int, attackDuration, stabilizeDuration time.Duration, targetTotalDifficultyRatio, miningRatio float64, attackerShouldWin bool) func(*agethSet) {
  return func(nodes *agethSet) {
    // Setup

    // Start all nodes mining at 150% of the blocktime. They will be the long tail of small miners.
    for i, node := range nodes.all() {
      node.startMining(blockTime * 3 / 2)
      if i > int(float64(len(nodes.all())) * miningRatio) { break }
    }
    bigMiners := newAgethSet()
    badGuy := nodes.indexed(0) // NOTE: Assumes badguy will always be [0]
    goodGuys := nodes.where(func(a *ageth) bool { return a.name != badGuy.name })

    // Simulate the small proportion of "whale" miners, whose block times will be nearer the target time.
    hashtimes := []int{blockTime, blockTime * 12/10, blockTime * 12/10, blockTime * 14/10, blockTime * 14/10}
    for _, hashtime := range hashtimes {
      nextMiner := nodes.random()
      for nextMiner.name == badGuy.name || bigMiners.contains(nextMiner) {
        nextMiner = nodes.random()
      }
      bigMiners.push(nextMiner)
      nextMiner.startMining(hashtime)
    }
    log.Info("Started miners")

    // Ensure nodes have SOME view of a network
    minimumPeerCount := int64(2)
    for _, n := range nodes.all() {
      for n.getPeerCount() < minimumPeerCount {
        n.addPeer(nodes.random())
      }
    }

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
    resumePeering := badGuy.refusePeers(100)
    forkBlock := badGuy.block()
    forkBlockTd := badGuy.getTd()
    badGuy.startMining(blockTime / 2)
    attackStartTime := time.Now()

    // Once a second, check to see if the bad guy's block difficulty has
    // reached the target.
    lastChainRatioTD := 0.0
    badGuyBlockTime := blockTime / 2
    for {
      bestPeer := goodGuys.peerMax()
      if bestPeer == nil { continue }

      // Get the total difficulty ratio of the forked chain / majority chain.
      chainRatioTD := float64(badGuy.getTd() - forkBlockTd) / float64(bestPeer.getTd() - forkBlockTd)

      if chainRatioTD != lastChainRatioTD {
        // The ratio has changed, adjust mining power
        if chainRatioTD < targetTotalDifficultyRatio {
          // We're behind the target ratio. We need to mine faster
          badGuyBlockTime--
        } else if chainRatioTD > targetTotalDifficultyRatio {
          // We're above the target ratio, we can mine slower.
          badGuyBlockTime++
        }
        if badGuyBlockTime < 1 {
          // We can't mine that fast.
          badGuyBlockTime = 1
        } else {
          badGuy.startMining(badGuyBlockTime)
        }
        lastChainRatio = chainRatio
      }
      targetDifficulty := uint64(float64(goodGuys.headBlock().Difficulty().Uint64()) * targetDifficultyRatio)
      if time.Since(attackStartTime) > attackDuration {
        log.Info("Attacker reached time limit without reaching target difficulty", "blocks", badGuy.block().number - forkBlock.number, "difficulty", badGuy.block().difficulty, "target", targetDifficulty)
        break
      }
      time.Sleep(time.Second)
    }

    finalAttackBlock := badGuy.block()
    log.Info("Attacker done mining", "mined", finalAttackBlock.number - forkBlock.number)
    badGuy.stopMining()
    // badGuy.setPeerCount(25)
    resumePeering()
    scenarioEnded := false
    defer func() { scenarioEnded = true }()
    attackMiners := make(chan struct{})
    go func() {
      for !scenarioEnded {
        // Sleep until the badguy has found 5 peers
        badGuy.addPeers(bigMiners) // Aggressively try to add the miners. NOTE that this represents an attacker sophisticated enough to identify and target miners.
        if badGuy.getPeerCount() >= 5 { attackMiners <- struct{}{} }
        time.Sleep(1 * time.Second)
      }
    }()
    <- attackMiners
    time.Sleep(stabilizeDuration)



    b, err := badGuy.eclient.BlockByNumber(context.Background(), big.NewInt(int64(finalAttackBlock.number)))
    if err != nil {
      log.Error("Error getting block", "blockno", finalAttackBlock.number, "node", badGuy.name)
    }
    if b.Hash() != finalAttackBlock.hash {
      log.Info("Even the attacker gave up their block")
    }

    forkBlockPlusOne, err := badGuy.eclient.BlockByNumber(context.Background(), big.NewInt(int64(forkBlock.number + 1)))
    if err != nil {
      log.Error("Error getting forkBlock + 1", "err", err)
    }

    convertedNodes := 0
    unconvertedNodes := 0
    for _, node := range goodGuys.all() {
      b, err := node.eclient.BlockByNumber(context.Background(), forkBlockPlusOne.Number())
      if err != nil {
        log.Error("Error getting block", "blockno", finalAttackBlock.number, "node", node.name, "err", err)
        continue
      }
      if b.Hash() == forkBlockPlusOne.Hash() {
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
      AttackerShouldWin: attackerShouldWin,
      AttackerWon: unconvertedNodes == 0,
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
