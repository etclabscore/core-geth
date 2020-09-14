package cmd

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	llog "log"
	"math/big"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	log "github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/ethereum/go-ethereum/rpc"
	"syreclabs.com/go/faker"
)

type ageth struct {
	name              string
	ipcpath           string
	command           *exec.Cmd
	proc              *os.Process
	log               log.Logger
	logr              io.ReadCloser
	client            *rpc.Client
	eclient           *ethclient.Client
	newheadChan       chan *types.Header
	latestBlock       *types.Block
	headsub           ethereum.Subscription
	errChan           chan error
	eventChan         chan interface{}
	enode             string
	mining            int
	peers             *agethSet
	coinbase          common.Address
	behaviorsInterval time.Duration
	behaviors         []func()
	td                *big.Int
	isMining          bool
	quitChan          chan struct{}
}

func newAgeth() *ageth {

	// ID
	var name = faker.Name().FirstName()
	if len(runningRegistry) == 0 {
		name = "Aarchimedes"
	}
	for !nameIsValid(name) {
		name = faker.Name().FirstName()
	}

	ipcpath := filepath.Join(os.TempDir(), fmt.Sprintf("ageth-%d.ipc", rand.Int()))

	datadir := filepath.Join(os.TempDir(), "ageth", name)
	os.MkdirAll(datadir, os.ModePerm)

	ks := filepath.Join(datadir, "keystore")
	os.MkdirAll(ks, os.ModePerm)

	passphraseFile := filepath.Join(ks, "pass.txt")
	ioutil.WriteFile(passphraseFile, []byte("foo"), os.ModePerm)

	createEtherbaseKey := exec.Command(gethPath, "--keystore", ks, "--password", passphraseFile, "account", "new")
	if err := createEtherbaseKey.Run(); err != nil {
		llog.Fatal(err)
	}

	gethArgs := []string{
		"--messnet",
		// "--ecbp1100", "9999",
		"--datadir", datadir,
		"--keystore", ks,
		"--fakepow",
		"--syncmode", "full",
		"--ipcpath", ipcpath,
		"--port", "0",
		"--maxpeers", "25",
		"--debug",
		"--nousb",
		"--ethash.dagsinmem", "0",
		"--ethash.dagsondisk", "0",
		"--ethash.cachesinmem", "0",
		"--ethash.cachesondisk", "0",

		"--metrics",
		"--metrics.influxdb",
		"--metrics.influxdb.database", "db0",

		// "--nodiscover",
		// "--mine", "--miner.threads", "0",
		// "--vmodule=eth/*=5,p2p=5,core/*=5",
		// "--verbosity", "5",

	}
	geth := exec.Command(gethPath, gethArgs...)
	p, err := geth.StderrPipe()
	if err != nil {
		log.Crit("stderr pipe", "error", err)
	}

	enodeNamesMu.Lock()
	enodeNames[name] = ""
	enodeNamesMu.Unlock()
	if len(name) > longestName {
		longestName = len(name)
	}
	a := &ageth{
		name:        name,
		command:     geth,
		ipcpath:     ipcpath,
		logr:        p,
		newheadChan: make(chan *types.Header),
		errChan:     make(chan error),
		peers:       newAgethSet(),
		td:          big.NewInt(0),
		quitChan:    make(chan struct{}, 100),
	}
	a.log = log.Root().New("source", a.name)

	return a
}

type block struct {
	number     uint64
	hash       common.Hash
	coinbase   common.Address
	difficulty uint64
	td         *big.Int
	parentHash common.Hash
}

var big0 = big.NewInt(0)

func (a *ageth) block() block {
	if a.latestBlock == nil {
		return block{
			number: 0, hash: common.Hash{}, coinbase: common.Address{}, td: big0, parentHash: common.Hash{}, difficulty: 0,
		}
	}
	return block{
		number:     a.latestBlock.NumberU64(),
		hash:       a.latestBlock.Hash(),
		coinbase:   a.latestBlock.Coinbase(),
		difficulty: a.latestBlock.Difficulty().Uint64(),
		td:         a.td,
		parentHash: a.latestBlock.ParentHash(),
	}
}

func (a *ageth) start() {
	err := a.command.Start()
	if err != nil {
		a.log.Crit("start geth", "error", err)
	}
	a.proc = a.command.Process
}

func (a *ageth) stop() {
	a.log.Info("Stopping ageth", "name", a.name)
	if a.eventChan != nil {
		a.eventChan <- eventNode{
			Node: Node{
				Name:      a.name,
				HeadHash:  a.block().hash.Hex()[2:8],
				HeadNum:   a.block().number,
				HeadMiner: a.block().coinbase == a.coinbase,
				HeadD:     a.block().difficulty,
				HeadTD:    a.td.Uint64(),
			},
			Up: false,
		}
	}
	err := a.proc.Kill()
	if err != nil {
		a.log.Crit("stop geth", "error", err)
	}
	delete(enodeNames, a.name)
	os.Remove(a.ipcpath)
	for i := 0; i < 100; i++ {
		a.quitChan <- struct{}{}
	}
}

func (a *ageth) run() {
	a.log.Info("Running ageth", "name", a.name)
	var ready bool
	a.start()
	go func() {
		buf := bufio.NewScanner(a.logr)
		for buf.Scan() {
			text := buf.Text()
			// Un/comment me to stream verbose geth logs on stdout.
			// String repeater is poor mans columnar alignment.
			fmt.Printf("%s%s: %s\n", a.name, strings.Repeat(" ", longestName-len(a.name)), text)

			// Wait for geth's IPC to initialize.
			if strings.Contains(text, "IPC endpoint opened") {
				ready = true
			}

			// Parse and assign ageth's enode from logs.
			if strings.Contains(text, "self=") {
				ee := regexp.MustCompile(`enode[a-zA-Z0-9\.\?=\:@\/]+`).FindString(text)
				a.log.Info("Found enode", "enode", ee)

				n := enode.MustParse(ee)
				a.enode = n.URLv4()

				enodeNamesMu.Lock()
				enodeNames[a.name] = a.enode
				enodeNamesMu.Unlock()

				regMu.Lock()
				runningRegistry[a.enode] = a
				regMu.Unlock()
			}
		}
	}()
	go func() {
		for {
			select {
			case h := <-a.newheadChan:
				b, err := a.eclient.BlockByHash(context.Background(), h.Hash())
				if err != nil {
					a.log.Error("get block by hash", "error", err)
				}
				if b == nil || h.Hash() == (common.Hash{}) {
					continue
				}
				if b.Hash() == a.block().hash {
					continue
				}
				a.setHead(b)
			case <-a.quitChan:
				return
			}
		}
	}()

	a.log.Info("Waiting for IPC to start")
	for !ready {
	}
	a.log.Info("IPC started")

	cl, err := rpc.DialIPC(context.Background(), a.ipcpath)
	if err != nil {
		log.Crit("rpc client", "error", err)
	}
	a.client = cl

	ecl, err := ethclient.Dial(a.ipcpath)
	if err != nil {
		log.Crit("dial ethclient", "error", err)
	}
	a.eclient = ecl
	sub, err := a.eclient.SubscribeNewHead(context.Background(), a.newheadChan)
	if err != nil {
		a.log.Crit("subscribe new head", "error", err)
	}
	a.headsub = sub
	go func() {
		for {
			select {
			case err := <-sub.Err():
				a.errChan <- err
			case <-a.quitChan:
				return
			default:
			}
		}
	}()
	go func() {
		for {
			select {
			case e := <-a.errChan:
				a.log.Error("errored", "error", e)
			case <-a.quitChan:
				return
			default:
			}
		}
	}()
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		for range ticker.C {
			select {
			case <-a.quitChan:
				ticker.Stop()
				return
			default:
				a.updateSelfStats()
			}
		}
	}()

	go func() {
		if len(a.behaviors) == 0 {
			return
		}
		for {
			select {
			case <-a.quitChan:
				return
			default:
			}
			time.Sleep(a.behaviorsInterval)
			for _, f := range a.behaviors {
				f()
			}
		}
	}()
}

func (a *ageth) withStandardPeerChurn(targetPeers int, peerSet *agethSet) {
	go func() {
		for {
			select {
			case <-a.quitChan:
				return
			default:
				time.Sleep(time.Second)
				if peerSet.len() == 0 {
					continue
				}
				// Try to keep the good guys at about 1/4 total network size. This is probably a little higher
				// proportionally than the mainnet, but close-ish, I guess. (Assuming ~500-800 total nodes, around 25-50 peer limits).
				// The important related variable here is DefaultSyncMinPeers, which Permapoint uses
				// to disable itself if the node's peer count falls below it.
				// The normal value for DefaultSyncMinPeers is 5, and the normal MaxPeers value (which geth
				// tries to fill) is 25.
				if rand.Float64() < 0.95 {
					continue
				}
				offset := a.peers.len() - targetPeers
				if offset > 0 {
					a.removePeer(a.peers.random()) // not exact, but closeish
				} else if offset < 0 {
					a.addPeer(peerSet.random())
				} else {
					// simulate some small churn
					if rand.Float32() < 0.5 {
						a.addPeer(peerSet.random())
					} else {
						a.removePeer(a.peers.random())
					}
				}
			}
		}
	}()
}

func (a *ageth) startMining(n int) {
	a.log.Info("Start mining", "threads", n)
	var ok bool
	err := a.client.Call(&ok, "miner_start", n)
	if err != nil {
		a.log.Crit("miner_start", "error", err)
	}
	a.mining = n
	a.isMining = true
	err = a.client.Call(&a.coinbase, "eth_coinbase")
	if err != nil {
		a.log.Crit("eth_coinbase", "error", err)
	}
	a.log.Info("Start mining")
}

func (a *ageth) stopMining() {
	a.log.Info("Stop mining")
	var ok bool
	err := a.client.Call(&ok, "miner_stop")
	if err != nil {
		a.log.Crit("miner_stop", "error", err)
	}
	a.log.Info("Stop mining")
	a.mining = 0
	a.isMining = false
}

func lookupNameByEnode(enode string) string {
	enodeNamesMu.Lock()
	defer enodeNamesMu.Unlock()
	for k, v := range enodeNames {
		if v == enode {
			return k
		}
	}
	return ""
}

func (a *ageth) addPeer(b *ageth) {
	var ok bool
	err := a.client.Call(&ok, "admin_addPeer", b.enode)
	if err != nil {
		a.log.Error("admin_addPeer", "error", err)
		return
	}
	if !a.peers.push(b) {
		return
	}
	a.log.Debug("Add peer", "target", b.name, "status", ok)
	if a.eventChan != nil {
		a.eventChan <- eventPeer{}
	}
}

func (a *ageth) removePeer(b *ageth) {
	var ok bool
	err := a.client.Call(&ok, "admin_removePeer", b.enode)
	if err != nil {
		a.log.Error("admin_removePeer", "error", err)
		return
	}
	a.peers.remove(b)
	a.log.Debug("Remove peer", "target", b.name, "status", ok)
	if a.eventChan != nil {
		a.eventChan <- eventPeer{}
	}
}

func (a *ageth) removeAllPeers() {
	for _, p := range a.peers.all() {
		a.removePeer(p)
	}
}

func (a *ageth) refreshPeers() {
	res := []p2p.NodeInfo{}
	err := a.client.Call(&res, "admin_peers")
	if err != nil {
		a.log.Error("admin_peers", "error", err)
		return
	}
	// a.peers.ageths = []*ageth{} // clear
	incomingSet := newAgethSet()
	for _, r := range res {
		n := enode.MustParse(r.Enode)
		if n.IP().IsLoopback() || n.IP().IsLinkLocalUnicast() || n.IP().IsLinkLocalMulticast() {
			nn := getAgethByEnode(n.URLv4())
			incomingSet.push(nn)
			if nn == nil {
				s := []string{}
				for k := range runningRegistry {
					s = append(s, k)
				}
				ss := strings.Join(s, "\n")
				log.Crit("Bad enode parsing (nonexist peer)", "enode", n, "running", ss)
			}
			// if !a.peers.contains(nn) {
			// 	a.peers.push(nn)
			// }
		}
	}
	a.peers = incomingSet
	// for _, p := range a.peers.all() {
	// 	nn := getAgethByEnode(p.enode)
	// 	if !incomingSet.contains(nn) {
	// 		a.peers.remove(nn)
	// 	}
	// }
	if a.eventChan != nil {
		a.eventChan <- eventPeer{}
	}
}

func (a *ageth) updateSelfStats() {
	a.refreshPeers()
	// a.getHeadManually()
	// err := a.client.Call(&a.isMining, "eth_mining")
	// if err != nil {
	// 	a.log.Error("eth_mining", "error", err)
	// 	return
	// }
	// a.log.Info("Status", "is_mining", a.isMining)
}

func (a *ageth) getHeadManually() {
	b, err := a.eclient.BlockByNumber(context.Background(), nil)
	if err != nil {
		a.log.Error("get block by number [nil=latest]", "error", err)
	}
	if b == nil {
		a.log.Warn("No latest block")
		return
	}
	if b.Hash() == a.block().hash {
		return
	}
	a.setHead(b)
}

func (a *ageth) truncateHead(n uint64) {
	a.log.Warn("Truncating head", "number", n)
	var res bool
	err := a.client.Call(&res, "debug_setHead", hexutil.EncodeUint64(n)) // have to pass as string, annoying
	if err != nil {
		a.log.Error("debug_setHead", "error", err)
	}
}

func (a *ageth) setHead(head *types.Block) {
	a.latestBlock = head
	if a.eventChan != nil {
		select {
		case a.eventChan <- eventNode{
			Node: Node{
				Name:      a.name,
				HeadHash:  a.block().hash.Hex()[2:8],
				HeadNum:   a.block().number,
				HeadMiner: a.block().coinbase == a.coinbase,
				HeadD:     a.block().difficulty,
				HeadTD:    a.td.Uint64(),
			},
			Up: true,
		}:
		default:
		}

	}
}

func (a *ageth) onNewHead(head *types.Block) {

}
