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
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/spf13/cobra"
	"github.com/tidwall/pretty"
)

// readChainconfigCmd represents the readChainconfig command
var readChainconfigCmd = &cobra.Command{
	Use:   "read-chainconfig",
	Short: "Print stored chain config",
	Long: `Chain configs are stored in the database and can be printed (dumped).
This command does that.

Use:

	read-chainconfig <0xgenesisHash>

Example:

	read-chainconfig 0xd4e56740f876aef8c010b86a40d5f56745a118d0906a34e69aec8c0db1cb8fa3

`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Fatal("need canonical hash key for config")
		}

		log.Println("Opening database...")
		db, err := rawdb.NewLevelDBDatabase(chainDBPath, 256, 16, "")
		if err != nil {
			log.Fatal(err)
		}

		data, err := db.Get(rawdb.ConfigKey(common.HexToHash(args[0])))
		if err != nil {
			log.Fatal(err)
		}

		data = pretty.Pretty(data)
		fmt.Println(string(data))
	},
}

func init() {
	rootCmd.AddCommand(readChainconfigCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// readChainconfigCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// readChainconfigCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
