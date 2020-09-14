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
	"encoding/json"
	"io"
	"io/ioutil"
	llog "log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"github.com/rakyll/statik/fs"
	"github.com/spf13/cobra"
)

// playCmd represents the play command
var playCmd = &cobra.Command{
	Use:   "play",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		var target io.Reader

		// HTTP/WS stuff.
		http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			ws, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				log.Error("WS errored", "error", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			defer ws.Close()

			scanner := bufio.NewScanner(target)
			//adjust the capacity to your need (max characters in line)
			const maxCapacity = 1024*1024*8
			buf := make([]byte, maxCapacity)
			scanner.Buffer(buf, maxCapacity)

			lineN := 0
			for scanner.Scan() {
				n := NetworkGraphData{}
				line := scanner.Bytes()
				err = json.Unmarshal(line, &n)
				if err != nil {
					log.Warn("Scan failed to parse network graphic data from line", "line", string(line))
				}

				log.Info("Scanned", "line", lineN)
				if n.Tick == 0 {
					n.Tick = lineN
				}

				err := ws.WriteJSON(Event{
					Typ:     "state",
					Payload: n,
				})
				if err != nil {
					log.Warn("Write WS event errored", "error", err)
				}

				time.Sleep(50*time.Millisecond)
				lineN++
			}
		})

		// set up index.html handler
		statikFS, err := fs.New()
		if err != nil {
			llog.Fatal(err)
		}
		http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
			enableCors(writer)
			m, _ := url.ParseQuery(request.URL.RawQuery)
			if len(m["target"]) == 0 {
				log.Info("URL query param '?target=' is empty", "reading from", "stdin")
				target = os.Stdin
			} else {
				// check file exists and is readable at arg 0
				fi, err := os.Open(m["target"][0])
				if err != nil {
					llog.Fatal(err)
				}
				target = fi
			}
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

		// start server
		go func() {
			log.Info("Ready to play at http://localhost:8008 ▶")
			if err := http.ListenAndServe(":8008", nil); err != nil {
				llog.Fatal(err)
			}
		}()

		for {}
	},
}

func init() {
	rootCmd.AddCommand(playCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// playCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// playCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
