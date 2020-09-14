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
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	llog "log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/p2p/enode"
	_ "github.com/ethereum/go-ethereum/statik"
	"github.com/gorilla/websocket"
	"github.com/mattn/go-colorable"
	"github.com/mattn/go-isatty"
	"github.com/montanaflynn/stats"
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

var gethPath = "./build/bin/geth" // "/home/ia/go/src/github.com/ethereum/go-ethereum/build/bin/geth"

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

		if len(args) > 0 {
			gethPath = args[0]
		}

		var reportsDir string
		if reportToFS {
			reportsDir = filepath.Join("reports", fmt.Sprintf("%d", time.Now().Unix()))
			err := os.MkdirAll(reportsDir, os.ModePerm)
			if err != nil {
				llog.Fatal(err)
			}
			gethVersionCmd := exec.Command(gethPath, "version")
			gethVersionBytes, err := gethVersionCmd.CombinedOutput()
			if err != nil {
				llog.Fatal(err)
			}
			err = ioutil.WriteFile(filepath.Join(reportsDir, "metadata.txt"), gethVersionBytes, os.ModePerm)
			if err != nil {
				llog.Fatal(err)
			}
		}

		// "Global"s, don't touch.
		wsEventChan := make(chan interface{}, 10000)
		reportEventChan := make(chan interface{})
		defer close(wsEventChan)

		go func() {
			var fi *os.File
			if reportToFS {
				fi, err := os.OpenFile(filepath.Join(reportsDir, "log.txt"), os.O_CREATE|os.O_RDWR, os.ModePerm)
				if err != nil {
					llog.Fatal(err)
				}
				defer fi.Close()
			}
			for {
				select {
				case event := <-reportEventChan:
					select {
					case wsEventChan <- event: // forward to ws
					default:
						// log.Warn("failed to write event")
					}

					if reportToFS {
						// write to stable storage
						state := getWorldView(world)
						b, err := json.Marshal(state)
						if err != nil {
							llog.Fatal(err)
						}
						_, err = fi.Write(b)
						if err != nil {
							llog.Fatal(err)
						}
					}

				default:
				}
			}
		}()

		// HTTP/WS stuff.
		http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			ws, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				log.Error("WS errored", "error", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			debounce := time.NewTicker(300 * time.Millisecond)
			defer debounce.Stop()
			didEvent := false

			tookData := []float64{}
			quit := false
			go func() {
				for !quit {
					took := struct {
						Took int32 `json:"took"`
					}{}
					err = ws.ReadJSON(&took)
					if err != nil {
						log.Debug("WS read took message errored", "error", err)
						time.Sleep(time.Second)
						continue
					}
					tookData = append(tookData, float64(took.Took))
					if len(tookData) < 10 {
						continue
					}
					mean, _ := stats.Mean(tookData)
					if mean < 10 {
						mean = 10
					}
					tookData = []float64{}
					log.Debug("Update ticker interval", "interval.ms", mean)
					debounce.Stop()
					debounce = time.NewTicker(time.Duration(mean) * time.Millisecond)
				}
			}()

			defer func() {
				quit = true
				ws.Close()
			}()
			for {
				select {
				case <-debounce.C:
					if didEvent {
						payload := getWorldView(world)
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
				case <-wsEventChan:
					didEvent = true
					// On any event, just send out the whole global state.
					// Surely this isn't as efficient as can be, but saves me headaches,
					// and it's not _that_ much data.
					// payload := getWorldView(world)
					// payload.Tick = globalTick
					// globalTick++
					// err := ws.WriteJSON(Event{
					// 	Typ:     "state",
					// 	Payload: payload,
					// })
					// if err != nil {
					// 	log.Debug("Write WS event errored", "error", err)
					// }
				default:
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
			if err := http.ListenAndServe(":8008", nil); err != nil {
				llog.Fatal(err)
			}
		}()

		scenario1(reportEventChan)

		for {
			globalTick = 0
			// scenario2(reportEventChan)
			// scenario3(reportEventChan)
			// scenario4(reportEventChan)
			for _, g := range world.all() {
				g.stop()
				g = nil // KILLL
			}
			world = newAgethSet()
		}

	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var reportToFS bool

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ageth.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolVarP(&reportToFS, "report", "p", false, "Write a report to a timestamped directory")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

/*
	// alice := newAgeth()
	// alice.run()
	// alice.startMining(1)
	// go func() {
	// 	for alice.latestBlock == nil || alice.latestBlock.NumberU64() < 100 {
	// 		time.Sleep(1 * time.Second)
	// 	}
	// 	alice.stopMining()
	// }()
	//
	// bob := newAgeth()
	// bob.run()
	//
	// time.Sleep(15 * time.Second)
	// bob.startMining(1)
	// time.Sleep(5*time.Second)
	//
	// bob.addPeer(alice)
	//
	// time.Sleep(100 * time.Second)
	//
	// alice.stop()
	// bob.stop()

*/
