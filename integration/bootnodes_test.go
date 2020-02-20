package integration

import (
	"fmt"
	"net"
	"os"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/p2p/discover"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/ethereum/go-ethereum/params"
)

// startV4 starts an ephemeral discovery V4 node.
func startV4(t *testing.T) *discover.UDPv4 {
	socket, ln, cfg, err := listen()
	if err != nil {
		t.Fatal(err)
	}
	disc, err := discover.ListenV4(socket, ln, cfg)
	if err != nil {
		t.Fatal(err)
	}
	return disc
}

func listen() (*net.UDPConn, *enode.LocalNode, discover.Config, error) {
	var cfg discover.Config
	cfg.PrivateKey, _ = crypto.GenerateKey()
	db, _ := enode.OpenDB("")
	ln := enode.NewLocalNode(db, cfg.PrivateKey)

	socket, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.IP{0, 0, 0, 0}})
	if err != nil {
		db.Close()
		return nil, nil, cfg, err
	}
	addr := socket.LocalAddr().(*net.UDPAddr)
	ln.SetFallbackIP(net.IP{127, 0, 0, 1})
	ln.SetFallbackUDP(addr.Port)
	return socket, ln, cfg, nil
}

func testBootnodes(t *testing.T, nodes []string, minPassRate float64, maxTrials int) {
	if maxTrials == 0 {
		t.Skip("trials disabled")
	}

	total := len(nodes)
	failed := 0
	minPassN := float64(total) * minPassRate // Minimum number of nodes that must be reachable for the test not to fail.

	// Case where pass rate is epsilon (but non-zero) and rounding causes n nodes == 0; infer that at least just 1 node must pass.
	if minPassN == 0 && minPassRate > 0 {
		minPassN = 1
	}

	disc := startV4(t)

nodesloop:
	for _, n := range nodes {
		en, err := enode.ParseV4(n)
		if err != nil {
			t.Fatal(err)
		}
		for i := 1; i <= maxTrials; i++ {
			start := time.Now()
			err = disc.Ping(en)
			if err == nil {
				t.Logf("OK (RTT %v): enode=%s", time.Since(start), en.String())
				continue nodesloop
			}
		}
		// Max trial attempts were reached, all with errors.
		t.Logf("FAIL (%d/%d): enode=%s err=%v", maxTrials, maxTrials, en.String(), err)
		failed++
	}

	okCount := total - failed
	line := fmt.Sprintf("=> %.0f%% (%d / %d) nodes responded to ping [min pass rate = %.02f, max trials = %d]", float64(okCount)/float64(total)*100, okCount, total, minPassRate, maxTrials)
	if okCount < int(minPassN) {
		t.Error(line)
	} else {
		t.Log(line)
	}
}

func TestBootnodesDiscV4Ping(t *testing.T) {
	if os.Getenv("MULTIGETH_TEST_BOOTNODE_AVAILABILITY") != "on" {
		t.Skip("Skipping bootnode availability tests.")
	}

	// MinPassRate defines the minimum tolerance for node OK rate
	// 1.0 would require all nodes to pass, 0.0 would require none to pass.
	zero := 0.0
	epsilon := 0.01 // An epsilon pass rate (eg 0.01) will mean >= 1 node must succeed.
	few := 0.3
	some := 0.5
	//most := 0.7
	defaultMinPassRate := some

	// MaxTrials.
	defaultMaxTrials := 3

	for _, c := range []struct {
		Name        string
		Bootnodes   []string
		MinPassRate *float64
		MaxTrials   *int
	}{
		{Name: "classic", Bootnodes: params.ClassicBootnodes},
		{Name: "foundation", Bootnodes: params.MainnetBootnodes},
		{Name: "kotti", Bootnodes: params.KottiBootnodes},
		{Name: "goerli", Bootnodes: params.GoerliBootnodes, MinPassRate: &few},
		{Name: "mordor", Bootnodes: params.MordorBootnodes, MinPassRate: &epsilon},
		{Name: "ropsten", Bootnodes: params.TestnetBootnodes, MinPassRate: &epsilon},
		{Name: "rinkeby", Bootnodes: params.RinkebyBootnodes},
		{Name: "social", Bootnodes: params.SocialBootnodes, MinPassRate: &zero},
		{Name: "ethersocial", Bootnodes: params.EthersocialBootnodes},
		{Name: "mix", Bootnodes: params.MixBootnodes},
	} {
		t.Run(c.Name, func(t *testing.T) {
			rate := defaultMinPassRate
			if c.MinPassRate != nil {
				rate = *c.MinPassRate
			}
			trials := defaultMaxTrials
			if c.MaxTrials != nil {
				trials = *c.MaxTrials
			}

			testBootnodes(t, c.Bootnodes, rate, trials)
		})
	}
}
