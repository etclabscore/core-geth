/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	llog "log"
	"math/rand"
	"net/http"
	"os"
	"reflect"
	"runtime"
	"sync"
	"time"
	// _ "net/http/pprof"

	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/p2p/enode"
	_ "github.com/ethereum/go-ethereum/statik"
	"github.com/gorilla/websocket"
	"github.com/mattn/go-colorable"
	"github.com/mattn/go-isatty"
	"github.com/rakyll/statik/fs"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		EnableCompression: true,
	}
)

var (
	ostream log.Handler
	glogger *log.GlogHandler
)

func setupLogging() {
	usecolor := (isatty.IsTerminal(os.Stderr.Fd()) || isatty.IsCygwinTerminal(os.Stderr.Fd())) && os.Getenv("TERM") != "dumb"
	output := io.Writer(os.Stderr)
	if usecolor {
		output = colorable.NewColorableStderr()
	}
	ostream = log.StreamHandler(output, log.TerminalFormat(usecolor))
	glogger = log.NewGlogHandler(ostream)
}

func init() {
	rand.Seed(time.Now().UnixNano())
	setupLogging()
	log.PrintOrigins(true)
	glogger.Verbosity(3)
	glogger.Vmodule("")
	glogger.BacktraceAt("")
	log.Root().SetHandler(glogger)
}

// stupid little helpers to make sure that the random names
// assigned to each geth are unique
var enodeNames = make(map[string]string)
var enodeNamesMu = sync.Mutex{}

func nameIsValid(name string) bool {
	enodeNamesMu.Lock()
	defer enodeNamesMu.Unlock()
	_, ok := enodeNames[name]
	return !ok && name != ""
}

var longestName = 0

var runningRegistry = map[string]*ageth{}
var regMu = sync.Mutex{}

func getAgethByEnode(en string) *ageth {
	regMu.Lock()
	defer regMu.Unlock()
	want := enode.MustParse(en)
	for k, v := range runningRegistry {
		n := enode.MustParse(k)
		if want.ID() == n.ID() {
			return v
		}
	}
	return nil
}

// since geth can/should add and remove peers sovereignly, as well as manually,
// we'll only send notifications that some peer event happened. It is up
// to the reporter to provide a global state of connections.
type eventPeer struct{}

type eventNode struct {
	Node Node `json:"node"`
	Up   bool `json:"up"`
}

func endpoint(set *agethSet) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		enableCors(w)
		b, err := json.Marshal(getWorldView(set))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error("Network graph data query failed", "error", err)
			return
		}
		w.Write(b)
	}
}

func enableCors(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
}

type Event struct {
	Typ     string      `json:"type"`
	Payload interface{} `json:"payload"`
}

type scenario func(nodes *agethSet)

// "Global"s, don't touch.
var world = newAgethSet()
var globalTick = 0

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ageth",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {

		// "Global"s, don't touch.
		reportEventChan := make(chan interface{})    // this is what the geths get
		globalState := getWorldView(world)
		programQuitting := make(chan struct{})

		wsEventChan := make(chan interface{}, 10000) // it gets passed to wsEventChan for web ui
		defer close(wsEventChan)

		go func() {
			var writer io.Writer
			if reportToFS != "" {
				fi, err := os.OpenFile(reportToFS, os.O_CREATE|os.O_RDWR, os.ModePerm)
				if err != nil {
					llog.Fatal(err)
				}
				writer = fi
				defer fi.Close()
			} else if reportToStdout {
				writer = os.Stdout
			}
			lastReport := time.Now()
			for {
				select {
				case event := <-reportEventChan:
					select {
					case wsEventChan <- event: // forward to ws
					default:
						// log.Warn("failed to write event")
					}

					globalState = getWorldView(world)

					if writer != nil {
						// write to stable storage
						if time.Since(lastReport) < time.Second {
							continue
						}
						lastReport = time.Now()
						b, err := json.Marshal(globalState)
						if err != nil {
							llog.Fatal(err)
						}
						_, err = writer.Write(b)
						if err != nil {
							llog.Fatal(err)
						}
					}
				// default:
				}
			}
		}()

		// HTTP/WS stuff.
		if httpAddr != "" {
			go runWeb(wsEventChan, programQuitting, func() NetworkGraphData {
				return globalState
			})
		}

		listEndpoints := []string{}
		var readFrom io.Reader
		if endpointsFile != "" {
			log.Info("Reading endpoints from file", "file", endpointsFile)
			b, err := ioutil.ReadFile(endpointsFile)
			if err != nil {
				llog.Fatal(err)
			}
			readFrom = bytes.NewBuffer(b)
		} else {
			log.Info("Reading endpoints from stdin...")
			readFrom = os.Stdin
		}
		scanner := bufio.NewScanner(readFrom)
		for scanner.Scan() {
			ep := scanner.Text()
			if ep == "" {
				continue
			}
			listEndpoints = append(listEndpoints, ep)
		}
		if err := scanner.Err(); err != nil && !errors.Is(err, io.EOF) && !errors.Is(err, bufio.ErrTooLong) {
			llog.Fatal(err)
		}
		log.Info("Read endpoints", "length", len(listEndpoints))
		if len(listEndpoints) == 0 {
			log.Crit("No endpoints found")
		}

		agethEndpointCh := make(chan string)
		agethCh := make(chan *ageth)

		go func() {
			for _, e := range listEndpoints {
				agethEndpointCh <- e
			}
			close(agethEndpointCh)
		}()
		var wg sync.WaitGroup
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func(wg *sync.WaitGroup) {
				for e := range agethEndpointCh {
					g := newAgeth(e)
					g.eventChan = reportEventChan
					g.run()
					agethCh <- g
				}
				wg.Done()
			}(&wg)
		}
		go func(wg *sync.WaitGroup) {
			wg.Wait()
			close(agethCh)
		}(&wg)
		for g := range agethCh {
			world.push(g)
		}

		// If --read-only=true (which it IS, BY DEFAULT, this is where
		// the flow will stop.
		// You must use an explicit:
		//
		//   ageth --read-only=false
		//
		// to proceed, and to run the scenarios below.
		if readOnly {
			q := make(chan struct{})
			<-q
		}

		scenarios := []scenario{
			generateScenarioPartitioning(false, 3*time.Minute),
			generateScenarioPartitioning(false, 3*time.Minute),
			generateScenarioPartitioning(false, 3*time.Minute),
			generateScenarioPartitioning(true, 15*time.Minute),
			generateScenarioPartitioning(true, 15*time.Minute),

			// scenarioGenerator(13, 10 * time.Minute, 2 * time.Minute, 1.13, .666, true),
			// scenarioGenerator(13, 10 * time.Minute, 2 * time.Minute, 1.02, .666, false),

			// scenarioGenerator(13, 34 * time.Minute, 10 * time.Minute, 1.55, .666, true),
			// scenarioGenerator(13, 34 * time.Minute, 10 * time.Minute, 1.45, .666, false),
			//
			// scenarioGenerator(13, 49 * time.Minute, 10 * time.Minute, 2.1, .666, true),
			// scenarioGenerator(13, 49 * time.Minute, 10 * time.Minute, 1.9, .666, false),
			//
			// scenarioGenerator(13, 70 * time.Minute, 10 * time.Minute, 3.14, .666, true),
			// scenarioGenerator(13, 70 * time.Minute, 10 * time.Minute, 2.9, .666, false),
			//
			// scenarioGenerator(13, 86 * time.Minute, 10 * time.Minute, 4.15, .666, true),
			// scenarioGenerator(13, 86 * time.Minute, 10 * time.Minute, 3.9, .666, false),
			//
			// scenarioGenerator(13, 100 * time.Minute, 10 * time.Minute, 5.17, .666, true),
			// scenarioGenerator(13, 100 * time.Minute, 10 * time.Minute, 4.87, .666, false),

		}

		for i, s := range scenarios {
			log.Info("Running scenario", "index", i, "scenarios.len", len(scenarios),
				"name", runtime.FuncForPC(reflect.ValueOf(s).Pointer()).Name())
			globalTick = 0
			stabilize(world)
			s(world)
			// Note that the loop assumes no responsibility for tear down.
			// Each scenario needs to be responsible for getting the nodes
			// in the initial state they want them in without any assumptions
			// about what that might be.
			// This also means that any local geths left running at the end of a scenario
			// will still be running.
		}
		stabilize(world)
	},
}

func runWeb(reportEventChan chan interface{}, quitChan chan struct{}, globalState func() NetworkGraphData) {
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Error("WS errored", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		log.Info("Websocket connection", "remote.addr", r.RemoteAddr)

		payload := globalState()
		payload.Tick = globalTick
		globalTick++
		err = ws.WriteJSON(Event{
			Typ:     "state",
			Payload: payload,
		})
		if err != nil {
			log.Debug("Write WS event errored", "error", err)
		}

		debounce := time.NewTicker(300 * time.Millisecond)
		defer debounce.Stop()
		didEvent := false
		lastWSState := NetworkGraphData{}

		// tookData := []float64{}
		quit := false
		// go func() {
		// 	for !quit {
		// 		took := struct {
		// 			Took int32 `json:"took"`
		// 		}{}
		// 		err = ws.ReadJSON(&took)
		// 		if err != nil {
		// 			log.Debug("WS read took message errored", "error", err)
		// 			time.Sleep(time.Second)
		// 			continue
		// 		}
		// 		tookData = append(tookData, float64(took.Took))
		// 		if len(tookData) < 10 {
		// 			continue
		// 		}
		// 		mean, _ := stats.Mean(tookData)
		// 		if mean < 100 {
		// 			mean = 100
		// 		}
		// 		tookData = []float64{}
		// 		log.Debug("Update ticker interval", "interval.ms", mean)
		// 		debounce.Stop()
		// 		debounce = time.NewTicker(time.Duration(mean) * time.Millisecond)
		// 	}
		// }()

		defer func() {
			quit = true
			ws.Close()
		}()
		for {
			select {
			case <-debounce.C:
				if didEvent {
					payload := globalState()
					payload.Tick = globalTick
					globalTick++
					err := ws.WriteJSON(Event{
						Typ:     "state",
						Payload: payload,
					})
					if err != nil {
						log.Debug("Write WS event errored", "error", err)
					}
					didEvent = false
				}
			case <-reportEventChan:
				// The world's ugliest dedupery
				// to avoid sending websocket events for equivalent states.
				isDupe := true
				if len(lastWSState.Nodes) != len(globalState().Nodes) {
					isDupe = false
				}
				if isDupe && len(lastWSState.Links) != len(globalState().Links) {
					isDupe = false
				}
				if isDupe {
					b, _ := json.Marshal(lastWSState)
					b2, _ := json.Marshal(globalState)
					cb, cb2 := bytes.NewBuffer([]byte{}), bytes.NewBuffer([]byte{})
					json.Compact(cb, b)
					json.Compact(cb2, b2)
					isDupe = bytes.Equal(cb.Bytes(), cb2.Bytes())
				}
				didEvent = !isDupe
				lastWSState = globalState()
				// default:
			}
		}
	})
	http.HandleFunc("/state", func(writer http.ResponseWriter, request *http.Request) {
		endpoint(world)(writer, request)
	})
	statikFS, err := fs.New()
	if err != nil {
		llog.Fatal(err)
	}
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		enableCors(writer)
		r, err := statikFS.Open("/index.html")
		if err != nil {
			log.Error("Open index.html errored", "error", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		b, err := ioutil.ReadAll(r)
		if err != nil {
			log.Error("Read index.html errored", "error", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, err = writer.Write(b)
		if err != nil {
			log.Error("Write index.html errored", "error", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
	go func() {
		if !tls {
			log.Info("Serving HTTP", "addr", httpAddr)
			if err := http.ListenAndServe(httpAddr, nil); err != nil {
				llog.Fatal(err)
			}
		} else {
			log.Info("Serving HTTPS", "addr", httpAddr)
			if err := http.ListenAndServeTLS(httpAddr,
				"/root/.local/share/caddy/certificates/acme-v02.api.letsencrypt.org-directory/mess.canhaz.net/mess.canhaz.net.crt",
				"/root/.local/share/caddy/certificates/acme-v02.api.letsencrypt.org-directory/mess.canhaz.net/mess.canhaz.net.key",
				nil); err != nil {
				llog.Fatal(err)
			}
		}
	}()
	<-quitChan
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var tls bool
var reportToFS string
var reportToStdout bool
var endpointsFile string
var httpAddr string
var readOnly bool

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ageth.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().StringVarP(&reportToFS, "report", "r", "", "Write reporting logs to a given file")
	rootCmd.Flags().BoolVarP(&reportToStdout, "events-stdout", "e", false, "Write reporting logs to stdoutput")
	rootCmd.Flags().StringVarP(&endpointsFile, "endpoints", "f", "", "Read newline-deliminted endpoints from this file")
	rootCmd.Flags().BoolVarP(&readOnly, "read-only", "o", true, "Read only (dont run scenarios)")
	rootCmd.Flags().StringVarP(&httpAddr, "http", "p", "", "Serve http at endpoint")
	rootCmd.Flags().BoolVarP(&tls, "tls", "s", false, "Use HTTPS/TLS for websocket and http server")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
