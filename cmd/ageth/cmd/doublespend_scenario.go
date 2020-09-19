package cmd

import (
  "context"
  "encoding/json"
  "math/big"
  "os"
  "time"

  "github.com/ethereum/go-ethereum"
  "github.com/ethereum/go-ethereum/common"
  log "github.com/ethereum/go-ethereum/log"
)


type finalReport struct {
  Converted int
  Unconverted int
  DistinctChains []common.Hash
  Nodes map[string]common.Hash
  DifficultyRatio float64
  TargetDifficultyRatio float64
  AttackerShouldWin bool
  AttackerWon bool
  AttackerHash common.Hash
}

func stabilize(nodes *agethSet) {
  badGuy := nodes.indexed(0) // NOTE: Assumes badguy will always be [0]
  goodGuys := nodes.where(func(a *ageth) bool { return a.name != badGuy.name })
	badGuy.truncateHead(goodGuys.headMin())
  if !badGuy.sameChainAs(goodGuys.random()) {
		badGuy.truncateHead(0)
	}
  minimumPeerCount := int64(2)
  nodes.eachParallel(func (node *ageth) {
    node.stopMining()
    var result interface{}
    node.client.Call(&result, "admin_maxPeers", 20)
    for node.getPeerCount() < minimumPeerCount {
      node.addPeer(nodes.random())
    }
  })
  done := make(chan struct{})
  go func() {
    for {
      select {
      case <-done:
        return
      case <-time.NewTimer(30 * time.Second).C:
        log.Info("Still stabilizing", "distinctChains", len(nodes.distinctChains()), "badGuyBlock", badGuy.block().number, "goodGuysBlock", goodGuys.headMax() )
      }
    }
  }()
  goodGuys.random().startMining(13)
  for len(nodes.distinctChains()) > 1 {
    time.Sleep(30)
  }
  for badGuy.block().number < goodGuys.headMax() {
    time.Sleep(5)
  }
  done <- struct{}{}
}

func scenarioGenerator(blockTime int, attackDuration, stabilizeDuration time.Duration, targetDifficultyRatio, miningRatio, ecbp1100ratio float64, attackerShouldWin bool) func(*agethSet) {
  return func(nodes *agethSet) {
    // Setup

    // Start all nodes mining at 150% of the blocktime. They will be the long tail of small miners.
    for i, node := range nodes.all() {
      node.startMining(blockTime * 3 / 2)
      if i > int(float64(len(nodes.all())) * miningRatio) { break }
    }
    bigMiners := newAgethSet()
    healthyNodes := nodes.where(func(a *ageth) bool { return a.peers.len() > 12 })
    badGuy := nodes.indexed(0) // NOTE: Assumes badguy will always be [0]
    goodGuys := nodes.where(func(a *ageth) bool { return a.name != badGuy.name })

    // Simulate the small proportion of "whale" miners, whose block times will be nearer the target time.
    hashtimes := []int{blockTime, blockTime * 12/10, blockTime * 12/10, blockTime * 14/10, blockTime * 14/10}
    for _, hashtime := range hashtimes {
      nextMiner := healthyNodes.random()
      for nextMiner.name == badGuy.name || bigMiners.contains(nextMiner) {
        nextMiner = healthyNodes.random()
      }
      bigMiners.push(nextMiner)
      nextMiner.startMining(hashtime)
    }
    log.Info("Started miners")

    // Ensure nodes have SOME view of a network
    minimumPeerCount := int64(2)
    nodeCount := len(goodGuys.all())
    for i, n := range nodes.all() {
      var result interface{}
      var err error
      if i < int(float64(nodeCount) * ecbp1100ratio) {
        err = n.client.Call(&result, "admin_ecbp1100", "0x0")
      } else {
        err = n.client.Call(&result, "admin_ecbp1100", "0x999999999")
      }
      if err != nil { log.Error("Error setting ecbp110", "node", n.name, "err", err)}
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
    lastChainRatio := big.NewFloat(0)
    badGuyBlockTime := blockTime / 2
    for {
      bestPeer := goodGuys.peerMax()
      if bestPeer == nil { continue }
      chainRatio := big.NewFloat(0).Quo(big.NewFloat(0).SetInt(big.NewInt(0).Sub(badGuy.getTd(), forkBlockTd)), big.NewFloat(0).SetInt(big.NewInt(0).Sub(bestPeer.getTd(), forkBlockTd)))
      if chainRatio.Cmp(lastChainRatio) != 0 {
        // The ratio has changed, adjust mining power
        if chainRatio.Cmp(big.NewFloat(targetDifficultyRatio)) < 0  {
          // We're behind the target ratio. We need to mine faster
          badGuyBlockTime--
          if badGuyBlockTime > blockTime {
            // The tides have turned and we're behind where we should be. Get
            // back in line quickly, rather than one block at a time.
            badGuyBlockTime = blockTime
          }
        } else if chainRatio.Cmp(big.NewFloat(targetDifficultyRatio)) > 0 {
          // We're above the target ratio, we can mine slower.
          badGuyBlockTime++
          if badGuyBlockTime < blockTime / 2 {
            // We're mining way too fast. Slam on the breaks
            badGuyBlockTime = blockTime / 2
          }
        }
        if badGuyBlockTime < 1 {
          // We can't mine that fast.
          badGuyBlockTime = 1
        } else {
          badGuy.startMining(badGuyBlockTime)
        }
        lastChainRatio = chainRatio
      }
      if time.Since(attackStartTime) > attackDuration {
        log.Info("Attacker reached time limit", "blocks", badGuy.block().number - forkBlock.number, "difficulty", badGuy.block().difficulty, "chainRatio", chainRatio, "targetRatio", targetDifficultyRatio)
        break
      }
      time.Sleep(time.Second)
    }


    badGuy.stopMining()
    finalAttackBlock := badGuy.block()
    forkTestBlock, err := badGuy.eclient.BlockByNumber(context.Background(), big.NewInt(int64((forkBlock.number + finalAttackBlock.number) / 2)))
    if err != nil {
      log.Error("Error getting forkBlock + 1", "err", err)
    }
    log.Info("Attacker done mining", "mined", finalAttackBlock.number - forkBlock.number, "attackerhash", finalAttackBlock.hash, "attackernumber", finalAttackBlock.number)


    // Control: make sure bigMiners (at least) have sufficient peers
    // for MESS to be activated (>=10).
    go func() {
      log.Info("Controlling for good guys big miners to have a minimum num of peers", "minimum", 12)
      for _, bigMiner := range bigMiners.all() {
        for bigMiner.peers.len() < 12 {
          friend := goodGuys.random()
          bigMiner.log.Warn("Big miner low on peers", "count", bigMiner.peers.len(), "add", friend.name)
          bigMiner.addPeer(friend)
        }
      }
    }()

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

    convertedNodes := 0
    unconvertedNodes := 0
    for _, node := range goodGuys.all() {
      b, err := node.eclient.BlockByNumber(context.Background(), forkTestBlock.Number())
      if err != nil {
        if err == ethereum.NotFound {
          unconvertedNodes++
        } else {
          log.Error("RPC Error", "blockno", finalAttackBlock.number, "node", node.name, "err", err)
        }
        continue
      }
      if b.Hash() == forkTestBlock.Hash() {
        convertedNodes++
      } else {
        unconvertedNodes++
      }
    }
    difficultyRatio, _ := lastChainRatio.Float64()
    report := &finalReport{
      Converted: convertedNodes,
      Unconverted: unconvertedNodes,
      DistinctChains: nodes.distinctChains(),
      Nodes: make(map[string]common.Hash),
      DifficultyRatio: difficultyRatio,
      TargetDifficultyRatio: targetDifficultyRatio,
      AttackerShouldWin: attackerShouldWin,
      AttackerWon: unconvertedNodes == 0,
      AttackerHash: finalAttackBlock.hash,
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
