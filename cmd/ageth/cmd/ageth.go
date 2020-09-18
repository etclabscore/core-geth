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
	"net"
	"net/url"
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
	"github.com/ethereum/go-ethereum/p2p/nat"
	"github.com/ethereum/go-ethereum/rpc"
	"syreclabs.com/go/faker"
)

type ageth struct {
	name                        string
	rpcEndpointOrExecutablePath string
	command                     *exec.Cmd
	proc                        *os.Process
	log                         log.Logger
	logr                        io.ReadCloser
	client                      *rpc.Client
	eclient                     *ethclient.Client
	newheadChan                 chan *types.Header
	latestBlock                 *types.Header
	latestBlockUpdatedAt        time.Time
	headsub                     ethereum.Subscription
	errChan                     chan error
	eventChan                   chan interface{}
	enode                       string
	mining                      int
	peers                       *agethSet
	coinbase                    common.Address
	behaviorsInterval           time.Duration
	behaviors                   []func(self *ageth)
	td                          uint64
	tdhash                      common.Hash
	isMining                    bool
	quitChan                    chan struct{}
	online                      bool
	onHeadCallbacks             []func(*ageth, *types.Header)
}

func mustStartGethInstance(gethPath, id string) (*exec.Cmd, io.ReadCloser, string) {
	datadir := filepath.Join(os.TempDir(), "ageth", id)
	os.MkdirAll(datadir, os.ModePerm)

	ks := filepath.Join(datadir, "keystore")
	os.MkdirAll(ks, os.ModePerm)

	passphraseFile := filepath.Join(ks, "pass.txt")
	ioutil.WriteFile(passphraseFile, []byte("foo"), os.ModePerm)

	createEtherbaseKey := exec.Command(gethPath, "--keystore", ks, "--password", passphraseFile, "account", "new")
	if err := createEtherbaseKey.Run(); err != nil {
		llog.Fatal(err)
	}

	ipcpath := filepath.Join(os.TempDir(), fmt.Sprintf("ageth-%d.ipc", rand.Int()))

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

		// "--nodiscover",

		"--metrics",
		"--metrics.influxdb",
		"--metrics.influxdb.database", "db0",

		"--nodiscover",
		// "--mine", "--miner.threads", "0",
		// "--vmodule=eth/*=5,p2p=5,core/*=5",
		"--verbosity", "3",
	}
	geth := exec.Command(gethPath, gethArgs...)
	p, err := geth.StderrPipe()
	if err != nil {
		log.Crit("stderr pipe", "error", err)
	}
	return geth, p, ipcpath
}

// newAgeth returns a wrapped geth "ageth".
// If the rpcEndpointOrExecutablePath parameter is empty, a geth instance will be started.
func newAgeth(rpcEndpoint string) *ageth {

	// ID
	var name = faker.Name().FirstName()

	if len(runningRegistry) == 0 {
		name = "Aarchimedes"
	}
	for !nameIsValid(name) {
		name = faker.Name().FirstName()
	}

	enodeNamesMu.Lock()
	enodeNames[name] = ""
	enodeNamesMu.Unlock()
	if len(name) > longestName {
		longestName = len(name)
	}
	a := &ageth{
		name: name,
		// These variables will be set depending if the node is local or remote (currently conventionalized as ipc or not).
		// command:     geth,
		// logr:        p,
		rpcEndpointOrExecutablePath: rpcEndpoint,
		newheadChan:                 make(chan *types.Header),
		errChan:                     make(chan error),
		peers:                       newAgethSet(),
		td:                          0,
		quitChan:                    make(chan struct{}, 100),
	}

	u, _ := url.Parse(rpcEndpoint)
	isLocal := u.Scheme == ""
	if isLocal {
		log.Info("Starting local geth", "name", name, "executable", rpcEndpoint)
		a.command, a.logr, a.rpcEndpointOrExecutablePath = mustStartGethInstance(rpcEndpoint, name)
	} else {
		log.Info("Will connect with remote geth", "name", name, "endpoint", rpcEndpoint)
		a.rpcEndpointOrExecutablePath = rpcEndpoint
	}
	a.log = log.Root().New("source", a.name, "at", a.rpcEndpointOrExecutablePath)

	return a
}

type block struct {
	number     uint64
	hash       common.Hash
	coinbase   common.Address
	difficulty uint64
	td         uint64
	parentHash common.Hash
}

func (a *ageth) isLocal() bool {
	u, _ := url.Parse(a.rpcEndpointOrExecutablePath)
	return u.Scheme == ""
}

func (a *ageth) isRunning() bool {
	return a.client != nil
}

func (a *ageth) block() block {
	if a.latestBlock == nil {
		return block{
			number: 0, hash: common.Hash{}, coinbase: common.Address{}, td: 0, parentHash: common.Hash{}, difficulty: 0,
		}
	}
	return block{
		number:     a.latestBlock.Number.Uint64(),
		hash:       a.latestBlock.Hash(),
		coinbase:   a.latestBlock.Coinbase,
		difficulty: a.latestBlock.Difficulty.Uint64(),
		td:         a.td,
		parentHash: a.latestBlock.ParentHash,
	}
}

// startLocal starts a local instance.
// Scenarios should not call this themselves, it is handled exclusively by run().
// If (accidentally) called on a remote instance, nothing happens.
func (a *ageth) startLocal() {
	if !a.isLocal() {
		a.log.Warn("Noop to startLocal remote instance")
		return
	}
	err := a.command.Start()
	if err != nil {
		a.log.Crit("startLocal geth", "error", err)
	}
	a.proc = a.command.Process
}

// stop shuts down a local instance.
// If called on a remote instance, nothing happens.
func (a *ageth) stop() {
	if !a.isLocal() {
		a.log.Error("Can't stop remote geth instance")
		return
	}
	a.log.Info("Stopping ageth", "name", a.name)
	if a.eventChan != nil {
		a.eventChan <- eventNode{
			Node: Node{
				Name:      a.name,
				HeadHash:  a.block().hash.Hex()[2:8],
				HeadNum:   a.block().number,
				HeadMiner: a.block().coinbase == a.coinbase,
				HeadD:     a.block().difficulty,
				HeadTD:    a.td,
			},
			Up: false,
		}
	}
	err := a.proc.Kill()
	if err != nil {
		a.log.Crit("stop geth", "error", err)
	}
	delete(enodeNames, a.name)
	os.Remove(a.rpcEndpointOrExecutablePath) // this will fail if the path isn't an FS path. so ignore the error
	// Send a stupid number of quits to the quit channel to close all goroutines.
	//
	for i := 0; i < 100; i++ {
		a.quitChan <- struct{}{}
	}
	close(a.quitChan)
}

// run is idempotent.
// It should be called whenever you first want to startLocal using the ageth.
// It will startLocal the instance if it's not already started.
func (a *ageth) run() {

	if a.isRunning() {
		a.log.Warn("Already running")
		return
	}

	defer func() {
		a.online = true
	}()

	var ready bool

	if a.isLocal() {
		a.log.Info("Running ageth", "name", a.name)
		a.startLocal()
		go func() {
			buf := bufio.NewScanner(a.logr)
			for buf.Scan() {
				text := buf.Text()

				// Un/comment me to stream verbose geth logs on stdout.
				// String repeater is poor mans columnar alignment.
				fmt.Printf("%s%s: %s\n", a.name, strings.Repeat(" ", longestName-len(a.name)), text)

				// Wait for geth's IPC to initialize.
				if strings.Contains(text, "IPC endpoint opened") {
					time.Sleep(3 * time.Second)
					ready = true
				}
			}
		}()
	} else {
		ready = true
	}

	go func() {
		for {
			select {
			case h := <-a.newheadChan:
				a.setHead(h)
			case <-a.quitChan:
				return
			}
		}
	}()

	a.log.Info("Waiting for things to be ready")
	for !ready {
	}
	a.log.Info("Ageth ready")

	// Set up RPC clients.
	cl, err := rpc.Dial(a.rpcEndpointOrExecutablePath)
	if err != nil {
		log.Crit("rpc client", "error", err)
	}
	a.client = cl

	ecl, err := ethclient.Dial(a.rpcEndpointOrExecutablePath)
	if err != nil {
		log.Crit("dial ethclient", "error", err)
	}
	a.eclient = ecl

	// Get self enode information
	nodeInfoRes := p2p.NodeInfo{}
	err = a.client.Call(&nodeInfoRes, "admin_nodeInfo")
	if err != nil {
		a.log.Crit("admin_nodeInfo errored", "error", err)
	}

	n := enode.MustParse(nodeInfoRes.Enode)
	nv4 := n.URLv4()
	a.log.Info("Assigned self enode", "enode", nv4)
	a.enode = nv4

	etherbase := a.getEtherbase()
	if etherbase != (common.Address{}) {
		a.coinbase = etherbase
		a.log.Info("Assigned self etherbase", "etherbase", a.coinbase)
	}

	// Register self enode information to global map
	enodeNamesMu.Lock()
	enodeNames[a.name] = a.enode
	enodeNamesMu.Unlock()

	regMu.Lock()
	runningRegistry[a.enode] = a
	regMu.Unlock()

	// Subscribe to the client's new head
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
				// default:
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
				// default:
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
		for {
			select {
			case <-a.quitChan:
				return
			default:
			}
			if len(a.behaviors) == 0 {
				time.Sleep(time.Second)
				continue
			}
			time.Sleep(a.behaviorsInterval)
			for _, f := range a.behaviors {
				f(a)
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

func (a *ageth) registerNewHeadCallback(fn func(*ageth, *types.Header)) {
	a.onHeadCallbacks = append(a.onHeadCallbacks, fn)
}

func (a *ageth) purgeNewHeadCallbacks() {
	a.onHeadCallbacks = nil
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
	if a.coinbase == (common.Address{}) {
		err = a.client.Call(&a.coinbase, "eth_coinbase")
		if err != nil {
			a.log.Crit("eth_coinbase", "error", err)
		}
	}
}

func (a *ageth) stopMining() {
	a.log.Info("Stop mining")
	var ok bool
	err := a.client.Call(&ok, "miner_stop")
	if err != nil {
		a.log.Crit("miner_stop", "error", err)
	}
	a.mining = 0
	a.isMining = false
}

type tdstruct struct {
	types.Header
	TotalDifficulty hexutil.Uint64 `json:"totalDifficulty"`
}

func (a *ageth) getTd() uint64 {
	if a.tdhash == a.latestBlock.Hash() {
		return a.td
	}
	var td *tdstruct
	if err := a.client.Call(td, "eth_getBlockByHash", a.latestBlock.Hash(), false); err != nil {
		a.log.Error("error getting td", "err", err)
	}
	a.td = uint64(td.TotalDifficulty)
	a.tdhash = a.latestBlock.Hash()
	a.setHead(&td.Header)
	return a.td

}

var machineNATExtIP *net.IP

func translateEnodeIPIfLocal(en string) string {
	if machineNATExtIP == nil {
		iface, err := nat.Parse("pmp")
		if err != nil {
			log.Crit("Parse gateway errored", "error", err)
		}
		extIp, err := iface.ExternalIP()
		if err != nil {
			log.Crit("External IP errored", "error", err)
		}
		machineNATExtIP = &extIp
	}
	e := enode.MustParse(en)
	if (*machineNATExtIP).Equal(e.IP()) {
		b := []byte(en)
		b = regexp.MustCompile(`[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}`).ReplaceAll(b, []byte("127.0.0.1"))
		en = string(b)
	}
	return en
}

func (a *ageth) mustEtherbases(addresses []common.Address) {
	var ok bool
	err := a.client.Call(&ok, "admin_mustEtherbases", addresses)
	if err != nil {
		a.log.Error("admin_mustEtherbases errored", "error", err)
		return
	}
	a.log.Info("admin_mustEtherbases", "n", len(addresses))
}

func (a *ageth) getEtherbase() common.Address {
	ad := common.Address{}
	err := a.client.Call(&ad, "eth_coinbase")
	if err != nil {
		a.log.Error("eth_coinbase errored", "error", err)
	}
	a.log.Info("eth_coinbase", "coinbase", ad)
	return ad
}

func (a *ageth) addPeer(b *ageth) {
	var ok bool
	translatedIP := translateEnodeIPIfLocal(b.enode)
	err := a.client.Call(&ok, "admin_addPeer", translatedIP)
	if err != nil {
		a.log.Error("admin_addPeer", "error", err, "enode", translatedIP)
		return
	}
	if !a.peers.push(b) {
		return
	}
	a.log.Debug("Add peer", "target", b.name, "enode", b.enode, "status", ok)
	if a.eventChan != nil {
		a.eventChan <- eventPeer{}
	}
}

func (a *ageth) addPeers(s *agethSet) {
	for i := 0; i < s.len(); i++ {
		a.addPeer(s.indexed(i))
	}
}

func (a *ageth) removePeer(b *ageth) {
	var ok bool
	translatedIP := translateEnodeIPIfLocal(b.enode)
	err := a.client.Call(&ok, "admin_removePeer", translatedIP)
	if err != nil {
		a.log.Error("admin_removePeer", "error", err, "enode", translatedIP)
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
	incomingSet := newAgethSet()
	for _, r := range res {
		n := enode.MustParse(r.Enode)
		nn := getAgethByEnode(n.URLv4())
		incomingSet.push(nn)
	}
	a.peers = incomingSet
	if a.eventChan != nil {
		a.eventChan <- eventPeer{}
	}
}

func (a *ageth) updateSelfStats() {
	a.refreshPeers()

	if time.Since(a.latestBlockUpdatedAt) > 1*time.Minute {
		a.getHeadManually()
	}
}

func (a *ageth) getHeadManually() {
	head := types.Header{}
	err := a.client.Call(&head, "eth_getBlockByNumber", "latest", false)
	if err != nil {
		a.log.Error("get block by number [nil=latest]", "error", err)
	}
	if head.Number == nil {
		a.log.Warn("No latest block")
		return
	}
	if head.Hash() == a.block().hash {
		return
	}
	a.setHead(&head)
}

func (a *ageth) truncateHead(n uint64) {
	a.log.Warn("Truncating head", "number", n)
	var res bool
	err := a.client.Call(&res, "debug_setHead", hexutil.EncodeUint64(n)) // have to pass as string, annoying
	if err != nil {
		a.log.Error("debug_setHead", "error", err)
	}
	a.getHeadManually()
}

func (a *ageth) setHead(head *types.Header) {
	a.latestBlock = head
	if len(a.onHeadCallbacks) != 0 {
		for _, cb := range a.onHeadCallbacks {
			cb(a, head)
		}
	}
	a.latestBlockUpdatedAt = time.Now()
	if a.eventChan != nil {
		select {
		case a.eventChan <- eventNode{
			Node: Node{
				Name:      a.name,
				HeadHash:  a.block().hash.Hex()[2:8],
				HeadNum:   a.block().number,
				HeadMiner: a.block().coinbase == a.coinbase,
				HeadD:     a.block().difficulty,
				HeadTD:    a.td,
			},
			Up: true,
		}:
		default:
		}

	}
}

func (a *ageth) sameChainAs(b *ageth) bool {
	if a.block().hash == b.block().hash {
		return true
	}
	if a.block().number == b.block().number {
		return false
	}
	compare := func(x, y *ageth) bool {
		b, _ := y.eclient.BlockByNumber(context.Background(), big.NewInt(int64(x.block().number)))
		if b == nil {
			return false
		}
		return x.block().hash == b.Hash()
	}
	if a.block().number < b.block().number {
		return compare(a, b)
	}
	return compare(b, a)
}

func (a *ageth) onNewHead(head *types.Block) {

}

func (a *ageth) setMaxPeers(count int) {
	var result bool
	err := a.client.Call(&result, "admin_maxPeers", count)
	if err != nil {
		a.log.Error("admin_maxPeers errored", "error", err)
	}
}

func (a *ageth) getPeerCount() int64 {
	var result hexutil.Big
	a.client.CallContext(context.Background(), &result, "net_peerCount")
	return result.ToInt().Int64()
}

// refusePeers sets maxPeers to 0 (and drops all connected peers).
// It returns a function that has kept the original peers in a closure
// so that they may be added on resume peering.
func (a *ageth) refusePeers(resumeWithMaxPeers int) (resumePeers func()) {
	peers := []*p2p.PeerInfo{}
	a.client.Call(&peers, "admin_peers")
	var result interface{}
	a.setMaxPeers(0)
	return func() {
		a.setMaxPeers(resumeWithMaxPeers)
		for _, peer := range peers {
			a.client.Call(&result, "admin_addPeer", peer.Enode)
		}
	}
}
