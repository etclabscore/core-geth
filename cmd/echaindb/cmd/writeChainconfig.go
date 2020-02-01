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
	"io/ioutil"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/params/confp/generic"
	"github.com/spf13/cobra"
)

// writeChainconfigCmd represents the writeChainconfig command
var writeChainconfigCmd = &cobra.Command{
	Use:   "write-chainconfig",
	Short: "Write a chain config to the database",
	Long: `Chain configuration is stored in the chain database.
This command writes that value.

Value to write is taken from standard input (stdin).

Use:

	echaindb --chaindb <chaindata/path> write-chainconfig <0xgenesisHash>

Example:

	cat myconfig.json | echaindb --chaindb ./path/to/chaindata write-chainconfig 0xd4e56740f876aef8c010b86a40d5f56745a118d0906a34e69aec8c0db1cb8fa3

	echaindb --chaindb ./path/to/chaindata write-chainconfig 0xd4e56740f876aef8c010b86a40d5f56745a118d0906a34e69aec8c0db1cb8fa3 < myconfig.json 
	
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Fatal("need canonical hash key for config")
		}

		bs, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			log.Fatal(err)
		}

		conf, err := generic.UnmarshalChainConfigurator(bs)
		if err != nil {
			log.Fatal(err)
		}


		log.Println("Opening database...")
		db, err := rawdb.NewLevelDBDatabase(chainDBPath, 256, 16, "")
		if err != nil {
			log.Fatal(err)
		}

		rawdb.WriteChainConfig(db, common.HexToHash(args[0]), conf)
	},
}

func init() {
	rootCmd.AddCommand(writeChainconfigCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// writeChainconfigCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// writeChainconfigCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
