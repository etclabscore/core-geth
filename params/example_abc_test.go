package params

import (
	"encoding/json"
	"fmt"
	"log"
)

func ExampleABCGenesisJSON() {
	genesis := DefaultABCGenesisBlock()
	jsonBytes, err := json.MarshalIndent(genesis, "", "    ")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(jsonBytes))
	// Output:
	// {
	//     "config": {
	//         "networkId": 4269,
	//         "chainId": 4269,
	//         "eip2FBlock": 0,
	//         "eip7FBlock": 0,
	//         "eip150Block": 0,
	//         "eip155Block": 0,
	//         "eip160Block": 0,
	//         "eip161FBlock": 0,
	//         "eip170FBlock": 0,
	//         "eip100FBlock": 0,
	//         "eip140FBlock": 0,
	//         "eip198FBlock": 0,
	//         "eip211FBlock": 0,
	//         "eip212FBlock": 0,
	//         "eip213FBlock": 0,
	//         "eip214FBlock": 0,
	//         "eip658FBlock": 0,
	//         "eip145FBlock": 0,
	//         "eip1014FBlock": 0,
	//         "eip1052FBlock": 0,
	//         "eip152FBlock": 0,
	//         "eip1108FBlock": 0,
	//         "eip1344FBlock": 0,
	//         "eip1884FBlock": 0,
	//         "eip2028FBlock": 0,
	//         "eip2200FBlock": 0,
	//         "disposalBlock": 0,
	//         "ethash": {},
	//         "requireBlockHashes": {}
	//     },
	//     "nonce": "0x0",
	//     "timestamp": "0x6048d57c",
	//     "extraData": "0x42",
	//     "gasLimit": "0x2fefd8",
	//     "difficulty": "0x20000",
	//     "mixHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
	//     "coinbase": "0x0000000000000000000000000000000000000000",
	//     "alloc": {
	//         "366ae7da62294427c764870bd2a460d7ded29d30": {
	//             "balance": "0x2a"
	//         }
	//     },
	//     "number": "0x0",
	//     "gasUsed": "0x0",
	//     "parentHash": "0x0000000000000000000000000000000000000000000000000000000000000000"
	// }
	//
}
