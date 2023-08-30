// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package node

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/xeipuuv/gojsonschema"

	"github.com/stretchr/testify/assert"
)

var (
	testNodeKey, _ = crypto.GenerateKey()
)

func testNodeConfig() *Config {
	return &Config{
		Name: "test node",
		P2P:  p2p.Config{PrivateKey: testNodeKey},
	}
}

// Tests that an empty protocol stack can be closed more than once.
func TestNodeCloseMultipleTimes(t *testing.T) {
	stack, err := New(testNodeConfig())
	if err != nil {
		t.Fatalf("failed to create protocol stack: %v", err)
	}
	stack.Close()

	// Ensure that a stopped node can be stopped again
	for i := 0; i < 3; i++ {
		if err := stack.Close(); err != ErrNodeStopped {
			t.Fatalf("iter %d: stop failure mismatch: have %v, want %v", i, err, ErrNodeStopped)
		}
	}
}

func TestNodeStartMultipleTimes(t *testing.T) {
	stack, err := New(testNodeConfig())
	if err != nil {
		t.Fatalf("failed to create protocol stack: %v", err)
	}

	// Ensure that a node can be successfully started, but only once
	if err := stack.Start(); err != nil {
		t.Fatalf("failed to start node: %v", err)
	}
	if err := stack.Start(); err != ErrNodeRunning {
		t.Fatalf("start failure mismatch: have %v, want %v ", err, ErrNodeRunning)
	}
	// Ensure that a node can be stopped, but only once
	if err := stack.Close(); err != nil {
		t.Fatalf("failed to stop node: %v", err)
	}
	if err := stack.Close(); err != ErrNodeStopped {
		t.Fatalf("stop failure mismatch: have %v, want %v ", err, ErrNodeStopped)
	}
}

// Tests that if the data dir is already in use, an appropriate error is returned.
func TestNodeUsedDataDir(t *testing.T) {
	// Create a temporary folder to use as the data directory
	dir := t.TempDir()

	// Create a new node based on the data directory
	original, err := New(&Config{DataDir: dir})
	if err != nil {
		t.Fatalf("failed to create original protocol stack: %v", err)
	}
	defer original.Close()
	if err := original.Start(); err != nil {
		t.Fatalf("failed to start original protocol stack: %v", err)
	}

	// Create a second node based on the same data directory and ensure failure
	_, err = New(&Config{DataDir: dir})
	if err != ErrDatadirUsed {
		t.Fatalf("duplicate datadir failure mismatch: have %v, want %v", err, ErrDatadirUsed)
	}
}

// Tests whether a Lifecycle can be registered.
func TestLifecycleRegistry_Successful(t *testing.T) {
	stack, err := New(testNodeConfig())
	if err != nil {
		t.Fatalf("failed to create protocol stack: %v", err)
	}
	defer stack.Close()

	noop := NewNoop()
	stack.RegisterLifecycle(noop)

	if !containsLifecycle(stack.lifecycles, noop) {
		t.Fatalf("lifecycle was not properly registered on the node, %v", err)
	}
}

// Tests whether a service's protocols can be registered properly on the node's p2p server.
func TestRegisterProtocols(t *testing.T) {
	stack, err := New(testNodeConfig())
	if err != nil {
		t.Fatalf("failed to create protocol stack: %v", err)
	}
	defer stack.Close()

	fs, err := NewFullService(stack)
	if err != nil {
		t.Fatalf("could not create full service: %v", err)
	}

	for _, protocol := range fs.Protocols() {
		if !containsProtocol(stack.server.Protocols, protocol) {
			t.Fatalf("protocol %v was not successfully registered", protocol)
		}
	}

	for _, api := range fs.APIs() {
		if !containsAPI(stack.rpcAPIs, api) {
			t.Fatalf("api %v was not successfully registered", api)
		}
	}
}

// TestRegisterProtocols_OpenRPC tests whether a running stack adequately responds to rpc_discover.
func TestRegisterProtocols_OpenRPC_HTTP(t *testing.T) {
	stack, err := New(testNodeConfig())
	if err != nil {
		t.Fatalf("failed to create protocol stack: %v", err)
	}
	defer stack.Close()

	datadir := filepath.Join(os.TempDir(), "node-test")
	defer os.RemoveAll(datadir)

	stack.config = &DefaultConfig
	stack.config.HTTPHost = DefaultHTTPHost
	stack.config.DataDir = datadir

	if err := stack.Start(); err != nil {
		t.Fatal(err)
	}
	defer stack.Close()

	client, err := rpc.Dial("http://localhost:8545")
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	res := make(map[string]interface{})
	err = client.Call(&res, "rpc_discover")
	if err != nil {
		t.Errorf("client call: %v", err)
	}

	// b, _ := json.MarshalIndent(res, "", "    ")
	// t.Log("response", string(b))

	if res["info"].(map[string]interface{})["title"].(string) != "Core-Geth RPC API" {
		t.Fatal("bad")
	}
}

// TestRegisterProtocols_OpenRPC tests whether a running stack adequately responds to rpc_discover.
func TestRegisterProtocols_OpenRPC_WS(t *testing.T) {
	stack, err := New(testNodeConfig())
	if err != nil {
		t.Fatalf("failed to create protocol stack: %v", err)
	}
	defer stack.Close()

	datadir := filepath.Join(os.TempDir(), "node-test")
	defer os.RemoveAll(datadir)

	stack.config = &DefaultConfig
	stack.config.WSHost = DefaultWSHost
	stack.config.WSPort = DefaultWSPort
	stack.config.DataDir = datadir

	if err := stack.Start(); err != nil {
		t.Fatal(err)
	}
	defer stack.Close()

	client, err := rpc.Dial("ws://localhost:8546")
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	res := make(map[string]interface{})
	ctx, cancelFn := context.WithTimeout(context.Background(), time.Second)
	defer cancelFn()
	err = client.CallContext(ctx, &res, "rpc_discover")
	if err != nil {
		t.Errorf("client call: %v", err)
	}

	// b, _ := json.MarshalIndent(res, "", "    ")
	// t.Log("response", string(b))

	if _, ok := res["info"]; !ok {
		t.Fatal("no response")
		return
	}

	if res["info"].(map[string]interface{})["title"].(string) != "Core-Geth RPC API" {
		t.Fatal("bad")
	}
}

// TestOpenRPCMetaSchemaSpec tests compliance of discovery endpoint with the OpenRPC spec stored locally
func TestOpenRPCMetaSchemaSpec(t *testing.T) {
	stack, err := New(testNodeConfig())
	if err != nil {
		t.Fatalf("failed to create protocol stack: %v", err)
	}
	defer stack.Close()

	datadir := filepath.Join(os.TempDir(), "node-test")
	defer os.RemoveAll(datadir)

	stack.config = &DefaultConfig
	stack.config.HTTPHost = DefaultHTTPHost
	// enable all available APIs
	stack.config.HTTPModules = []string{"admin", "debug", "eth", "ethash", "miner", "net", "personal", "rpc", "trace", "txpool", "web3"}
	stack.config.DataDir = datadir

	if err := stack.Start(); err != nil {
		t.Fatal(err)
	}
	defer stack.Close()

	client, err := rpc.Dial("http://localhost:8545")
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	res := make(map[string]interface{})
	err = client.Call(&res, "rpc_discover")
	if err != nil {
		t.Errorf("client call: %v", err)
	}

	schemaLoader := gojsonschema.NewSchemaLoader()
	schemaLoader.Draft = gojsonschema.Draft7 // force schema version
	schemaLoader.AutoDetect = false

	openRPCSpecFile := "file://./testdata/open-rpc-meta-schema-1.14.0.json"
	schema, err := schemaLoader.Compile(gojsonschema.NewReferenceLoader(openRPCSpecFile))
	if err != nil {
		t.Fatal(err)
	}

	documentLoader := gojsonschema.NewGoLoader(res)

	result, err := schema.Validate(documentLoader)
	if err != nil {
		t.Errorf("failed to validate OpenRPC schema: %v", err)
	}

	if !result.Valid() {
		for _, desc := range result.Errors() {
			fmt.Printf("- %s\n", desc)
		}
		t.Fatal("the OpenRPC discovery document is not compliant with spec")
	}
}

// Suggested latest endpoint to be used is
// COREGETH_OPENRPC_LATEST_SPEC_ENDPOINT=https://github.com/open-rpc/meta-schema/releases/latest/download/open-rpc-meta-schema.json
var openRPCSpecLatestTestEnvVarName = "COREGETH_OPENRPC_LATEST_SPEC_ENDPOINT"

// TestOpenRPCMetaSchemaSpecLatest compares local JSON spec file with the latest JSON spec published remotely,
// so as we keep the local file in sync
func TestOpenRPCMetaSchemaSpecLatest(t *testing.T) {
	openRPCSpecLatestEndpoint := os.Getenv(openRPCSpecLatestTestEnvVarName)
	if os.Getenv("CI") == "" || openRPCSpecLatestEndpoint == "" {
		t.Skip()
	}

	client := &http.Client{}
	resp, err := client.Get(openRPCSpecLatestEndpoint)
	if err != nil {
		t.Fatalf("failed to fetch latest spec JSON: %v", err)
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)

	var latestSpec interface{}
	if err := dec.Decode(&latestSpec); err != nil {
		t.Fatalf("failed to decode latest spec JSON: %v", err)
	}

	localJSONFilePath := "./testdata/open-rpc-meta-schema-1.14.0.json"
	localJSONFile, err := os.ReadFile(localJSONFilePath)
	if err != nil {
		t.Fatalf("failed to load local spec JSON file: %v", err)
	}

	var localSpec interface{}
	if err := json.Unmarshal(localJSONFile, &localSpec); err != nil {
		t.Fatalf("failed to unmarshal local JSON spec: %v", err)
	}

	if !reflect.DeepEqual(latestSpec, localSpec) {
		t.Fatalf("Local OpenRPC Schema (%v) is not updated to the latest one (%v)", localJSONFilePath, openRPCSpecLatestEndpoint)
	}
}

// This test checks that open databases are closed with node.
func TestNodeCloseClosesDB(t *testing.T) {
	stack, _ := New(testNodeConfig())
	defer stack.Close()

	db, err := stack.OpenDatabase("mydb", 0, 0, "", false)
	if err != nil {
		t.Fatal("can't open DB:", err)
	}
	if err = db.Put([]byte{}, []byte{}); err != nil {
		t.Fatal("can't Put on open DB:", err)
	}

	stack.Close()
	if err = db.Put([]byte{}, []byte{}); err == nil {
		t.Fatal("Put succeeded after node is closed")
	}
}

// This test checks that OpenDatabase can be used from within a Lifecycle Start method.
func TestNodeOpenDatabaseFromLifecycleStart(t *testing.T) {
	stack, _ := New(testNodeConfig())
	defer stack.Close()

	var db ethdb.Database
	var err error
	stack.RegisterLifecycle(&InstrumentedService{
		startHook: func() {
			db, err = stack.OpenDatabase("mydb", 0, 0, "", false)
			if err != nil {
				t.Fatal("can't open DB:", err)
			}
		},
		stopHook: func() {
			db.Close()
		},
	})

	stack.Start()
	stack.Close()
}

// This test checks that OpenDatabase can be used from within a Lifecycle Stop method.
func TestNodeOpenDatabaseFromLifecycleStop(t *testing.T) {
	stack, _ := New(testNodeConfig())
	defer stack.Close()

	stack.RegisterLifecycle(&InstrumentedService{
		stopHook: func() {
			db, err := stack.OpenDatabase("mydb", 0, 0, "", false)
			if err != nil {
				t.Fatal("can't open DB:", err)
			}
			db.Close()
		},
	})

	stack.Start()
	stack.Close()
}

// Tests that registered Lifecycles get started and stopped correctly.
func TestLifecycleLifeCycle(t *testing.T) {
	stack, _ := New(testNodeConfig())
	defer stack.Close()

	started := make(map[string]bool)
	stopped := make(map[string]bool)

	// Create a batch of instrumented services
	lifecycles := map[string]Lifecycle{
		"A": &InstrumentedService{
			startHook: func() { started["A"] = true },
			stopHook:  func() { stopped["A"] = true },
		},
		"B": &InstrumentedService{
			startHook: func() { started["B"] = true },
			stopHook:  func() { stopped["B"] = true },
		},
		"C": &InstrumentedService{
			startHook: func() { started["C"] = true },
			stopHook:  func() { stopped["C"] = true },
		},
	}
	// register lifecycles on node
	for _, lifecycle := range lifecycles {
		stack.RegisterLifecycle(lifecycle)
	}
	// Start the node and check that all services are running
	if err := stack.Start(); err != nil {
		t.Fatalf("failed to start protocol stack: %v", err)
	}
	for id := range lifecycles {
		if !started[id] {
			t.Fatalf("service %s: freshly started service not running", id)
		}
		if stopped[id] {
			t.Fatalf("service %s: freshly started service already stopped", id)
		}
	}
	// Stop the node and check that all services have been stopped
	if err := stack.Close(); err != nil {
		t.Fatalf("failed to stop protocol stack: %v", err)
	}
	for id := range lifecycles {
		if !stopped[id] {
			t.Fatalf("service %s: freshly terminated service still running", id)
		}
	}
}

// Tests that if a Lifecycle fails to start, all others started before it will be
// shut down.
func TestLifecycleStartupError(t *testing.T) {
	stack, err := New(testNodeConfig())
	if err != nil {
		t.Fatalf("failed to create protocol stack: %v", err)
	}
	defer stack.Close()

	started := make(map[string]bool)
	stopped := make(map[string]bool)

	// Create a batch of instrumented services
	lifecycles := map[string]Lifecycle{
		"A": &InstrumentedService{
			startHook: func() { started["A"] = true },
			stopHook:  func() { stopped["A"] = true },
		},
		"B": &InstrumentedService{
			startHook: func() { started["B"] = true },
			stopHook:  func() { stopped["B"] = true },
		},
		"C": &InstrumentedService{
			startHook: func() { started["C"] = true },
			stopHook:  func() { stopped["C"] = true },
		},
	}
	// register lifecycles on node
	for _, lifecycle := range lifecycles {
		stack.RegisterLifecycle(lifecycle)
	}

	// Register a service that fails to construct itself
	failure := errors.New("fail")
	failer := &InstrumentedService{start: failure}
	stack.RegisterLifecycle(failer)

	// Start the protocol stack and ensure all started services stop
	if err := stack.Start(); err != failure {
		t.Fatalf("stack startup failure mismatch: have %v, want %v", err, failure)
	}
	for id := range lifecycles {
		if started[id] && !stopped[id] {
			t.Fatalf("service %s: started but not stopped", id)
		}
		delete(started, id)
		delete(stopped, id)
	}
}

// Tests that even if a registered Lifecycle fails to shut down cleanly, it does
// not influence the rest of the shutdown invocations.
func TestLifecycleTerminationGuarantee(t *testing.T) {
	stack, err := New(testNodeConfig())
	if err != nil {
		t.Fatalf("failed to create protocol stack: %v", err)
	}
	defer stack.Close()

	started := make(map[string]bool)
	stopped := make(map[string]bool)

	// Create a batch of instrumented services
	lifecycles := map[string]Lifecycle{
		"A": &InstrumentedService{
			startHook: func() { started["A"] = true },
			stopHook:  func() { stopped["A"] = true },
		},
		"B": &InstrumentedService{
			startHook: func() { started["B"] = true },
			stopHook:  func() { stopped["B"] = true },
		},
		"C": &InstrumentedService{
			startHook: func() { started["C"] = true },
			stopHook:  func() { stopped["C"] = true },
		},
	}
	// register lifecycles on node
	for _, lifecycle := range lifecycles {
		stack.RegisterLifecycle(lifecycle)
	}

	// Register a service that fails to shot down cleanly
	failure := errors.New("fail")
	failer := &InstrumentedService{stop: failure}
	stack.RegisterLifecycle(failer)

	// Start the protocol stack, and ensure that a failing shut down terminates all
	// Start the stack and make sure all is online
	if err := stack.Start(); err != nil {
		t.Fatalf("failed to start protocol stack: %v", err)
	}
	for id := range lifecycles {
		if !started[id] {
			t.Fatalf("service %s: service not running", id)
		}
		if stopped[id] {
			t.Fatalf("service %s: service already stopped", id)
		}
	}
	// Stop the stack, verify failure and check all terminations
	err = stack.Close()
	if err, ok := err.(*StopError); !ok {
		t.Fatalf("termination failure mismatch: have %v, want StopError", err)
	} else {
		failer := reflect.TypeOf(&InstrumentedService{})
		if err.Services[failer] != failure {
			t.Fatalf("failer termination failure mismatch: have %v, want %v", err.Services[failer], failure)
		}
		if len(err.Services) != 1 {
			t.Fatalf("failure count mismatch: have %d, want %d", len(err.Services), 1)
		}
	}
	for id := range lifecycles {
		if !stopped[id] {
			t.Fatalf("service %s: service not terminated", id)
		}
		delete(started, id)
		delete(stopped, id)
	}

	stack.server = &p2p.Server{}
	stack.server.PrivateKey = testNodeKey
}

// Tests whether a handler can be successfully mounted on the canonical HTTP server
// on the given prefix
func TestRegisterHandler_Successful(t *testing.T) {
	node := createNode(t, 7878, 7979)
	defer node.Close()
	// create and mount handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("success"))
	})
	node.RegisterHandler("test", "/test", handler)

	// start node
	if err := node.Start(); err != nil {
		t.Fatalf("could not start node: %v", err)
	}

	// create HTTP request
	httpReq, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:7878/test", nil)
	if err != nil {
		t.Error("could not issue new http request ", err)
	}

	// check response
	resp := doHTTPRequest(t, httpReq)
	buf := make([]byte, 7)
	_, err = io.ReadFull(resp.Body, buf)
	if err != nil {
		t.Fatalf("could not read response: %v", err)
	}
	assert.Equal(t, "success", string(buf))
}

// Tests that the given handler will not be successfully mounted since no HTTP server
// is enabled for RPC
func TestRegisterHandler_Unsuccessful(t *testing.T) {
	node, err := New(&DefaultConfig)
	if err != nil {
		t.Fatalf("could not create new node: %v", err)
	}

	// create and mount handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("success"))
	})
	node.RegisterHandler("test", "/test", handler)
}

// Tests whether websocket requests can be handled on the same port as a regular http server.
func TestWebsocketHTTPOnSamePort_WebsocketRequest(t *testing.T) {
	node := startHTTP(t, 0, 0)
	defer node.Close()

	ws := strings.Replace(node.HTTPEndpoint(), "http://", "ws://", 1)

	if node.WSEndpoint() != ws {
		t.Fatalf("endpoints should be the same")
	}
	if !checkRPC(ws) {
		t.Fatalf("ws request failed")
	}
	if !checkRPC(node.HTTPEndpoint()) {
		t.Fatalf("http request failed")
	}
}

func TestWebsocketHTTPOnSeparatePort_WSRequest(t *testing.T) {
	// try and get a free port
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal("can't listen:", err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	listener.Close()

	node := startHTTP(t, 0, port)
	defer node.Close()

	wsOnHTTP := strings.Replace(node.HTTPEndpoint(), "http://", "ws://", 1)
	ws := fmt.Sprintf("ws://127.0.0.1:%d", port)

	if node.WSEndpoint() == wsOnHTTP {
		t.Fatalf("endpoints should not be the same")
	}
	// ensure ws endpoint matches the expected endpoint
	if node.WSEndpoint() != ws {
		t.Fatalf("ws endpoint is incorrect: expected %s, got %s", ws, node.WSEndpoint())
	}

	if !checkRPC(ws) {
		t.Fatalf("ws request failed")
	}
	if !checkRPC(node.HTTPEndpoint()) {
		t.Fatalf("http request failed")
	}
}

type rpcPrefixTest struct {
	httpPrefix, wsPrefix string
	// These lists paths on which JSON-RPC should be served / not served.
	wantHTTP   []string
	wantNoHTTP []string
	wantWS     []string
	wantNoWS   []string
}

func TestNodeRPCPrefix(t *testing.T) {
	t.Parallel()

	tests := []rpcPrefixTest{
		// both off
		{
			httpPrefix: "", wsPrefix: "",
			wantHTTP:   []string{"/", "/?p=1"},
			wantNoHTTP: []string{"/test", "/test?p=1"},
			wantWS:     []string{"/", "/?p=1"},
			wantNoWS:   []string{"/test", "/test?p=1"},
		},
		// only http prefix
		{
			httpPrefix: "/testprefix", wsPrefix: "",
			wantHTTP:   []string{"/testprefix", "/testprefix?p=1", "/testprefix/x", "/testprefix/x?p=1"},
			wantNoHTTP: []string{"/", "/?p=1", "/test", "/test?p=1"},
			wantWS:     []string{"/", "/?p=1"},
			wantNoWS:   []string{"/testprefix", "/testprefix?p=1", "/test", "/test?p=1"},
		},
		// only ws prefix
		{
			httpPrefix: "", wsPrefix: "/testprefix",
			wantHTTP:   []string{"/", "/?p=1"},
			wantNoHTTP: []string{"/testprefix", "/testprefix?p=1", "/test", "/test?p=1"},
			wantWS:     []string{"/testprefix", "/testprefix?p=1", "/testprefix/x", "/testprefix/x?p=1"},
			wantNoWS:   []string{"/", "/?p=1", "/test", "/test?p=1"},
		},
		// both set
		{
			httpPrefix: "/testprefix", wsPrefix: "/testprefix",
			wantHTTP:   []string{"/testprefix", "/testprefix?p=1", "/testprefix/x", "/testprefix/x?p=1"},
			wantNoHTTP: []string{"/", "/?p=1", "/test", "/test?p=1"},
			wantWS:     []string{"/testprefix", "/testprefix?p=1", "/testprefix/x", "/testprefix/x?p=1"},
			wantNoWS:   []string{"/", "/?p=1", "/test", "/test?p=1"},
		},
	}

	for _, test := range tests {
		test := test
		name := fmt.Sprintf("http=%s ws=%s", test.httpPrefix, test.wsPrefix)
		t.Run(name, func(t *testing.T) {
			cfg := &Config{
				HTTPHost:       "127.0.0.1",
				HTTPPathPrefix: test.httpPrefix,
				WSHost:         "127.0.0.1",
				WSPathPrefix:   test.wsPrefix,
			}
			node, err := New(cfg)
			if err != nil {
				t.Fatal("can't create node:", err)
			}
			defer node.Close()
			if err := node.Start(); err != nil {
				t.Fatal("can't start node:", err)
			}
			test.check(t, node)
		})
	}
}

func (test rpcPrefixTest) check(t *testing.T, node *Node) {
	t.Helper()
	httpBase := "http://" + node.http.listenAddr()
	wsBase := "ws://" + node.http.listenAddr()

	if node.WSEndpoint() != wsBase+test.wsPrefix {
		t.Errorf("Error: node has wrong WSEndpoint %q", node.WSEndpoint())
	}

	for _, path := range test.wantHTTP {
		resp := rpcRequest(t, httpBase+path, testMethod)
		if resp.StatusCode != 200 {
			t.Errorf("Error: %s: bad status code %d, want 200", path, resp.StatusCode)
		}
	}
	for _, path := range test.wantNoHTTP {
		resp := rpcRequest(t, httpBase+path, testMethod)
		if resp.StatusCode != 404 {
			t.Errorf("Error: %s: bad status code %d, want 404", path, resp.StatusCode)
		}
	}
	for _, path := range test.wantWS {
		err := wsRequest(t, wsBase+path)
		if err != nil {
			t.Errorf("Error: %s: WebSocket connection failed: %v", path, err)
		}
	}
	for _, path := range test.wantNoWS {
		err := wsRequest(t, wsBase+path)
		if err == nil {
			t.Errorf("Error: %s: WebSocket connection succeeded for path in wantNoWS", path)
		}
	}
}

func createNode(t *testing.T, httpPort, wsPort int) *Node {
	conf := &Config{
		HTTPHost:     "127.0.0.1",
		HTTPPort:     httpPort,
		WSHost:       "127.0.0.1",
		WSPort:       wsPort,
		HTTPTimeouts: rpc.DefaultHTTPTimeouts,
	}
	node, err := New(conf)
	if err != nil {
		t.Fatalf("could not create a new node: %v", err)
	}
	return node
}

func startHTTP(t *testing.T, httpPort, wsPort int) *Node {
	node := createNode(t, httpPort, wsPort)
	err := node.Start()
	if err != nil {
		t.Fatalf("could not start http service on node: %v", err)
	}

	return node
}

func doHTTPRequest(t *testing.T, req *http.Request) *http.Response {
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("could not issue a GET request to the given endpoint: %v", err)
	}
	t.Cleanup(func() { resp.Body.Close() })
	return resp
}

func containsProtocol(stackProtocols []p2p.Protocol, protocol p2p.Protocol) bool {
	for _, a := range stackProtocols {
		if reflect.DeepEqual(a, protocol) {
			return true
		}
	}
	return false
}

func containsAPI(stackAPIs []rpc.API, api rpc.API) bool {
	for _, a := range stackAPIs {
		if reflect.DeepEqual(a, api) {
			return true
		}
	}
	return false
}
