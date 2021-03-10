# Adding a Network to Core-Geth

CoreGeth currently supports a handful of networks out of the box, 
and can readily be configured to support others.

This guide will show you how to add a network to CoreGeth.

For the context of this tutorial, I'm going to use __AlphaBeta Coin__:registered:
as the name of my new network. 
 - AlphaBeta Coin (ABC) will use Proof-of-Work for block issuance, 
   namely ETHash. (Just like Ethereum.)
 - AlphaBeta Coin will have some arbitrary pre-mine funds allocated to a single address.
 - AlphaBeta Coin will have the "Istanbul" (aka "ETC's Phoenix") protocol upgrades and
   EVM features activated from genesis (the very first block (number `0`)).

### Define the configuration.

Commit: `a901f84401`

```
> git --no-pager diff --name-only v1.11.22
docs/developers/add-network.md
params/bootnodes_abc.go
params/config_abc.go
params/config_abc_test.go
params/genesis_abc.go
```

Commit `2233ec80cb` writes an `Example`-style test showing and validating the
JSON encoding of our new configuration.

---

At this point the hard work is done.

We can now pursue two paths:

1. Use the JSON configuration to initialize a chaindata database and start our node(s), and/or
2. Expose the configuration as a core-geth default through `geth`'s CLI flags via `--abc`.

This tutorial won't cover (2) (yet). 

#### Initialize geth's database from the JSON configuration.

Build `geth`.
```
> make geth
```

Create a file containing the JSON encoding of ABC network's configuration (JSON data taken from the example test above).
```
> cat <<EOF > abc_genesis.json
{
  "config": {
    "networkId": 4269, 
    "chainId": 4269, 
    "eip2FBlock": 0, 
    "eip7FBlock": 0, 
    "eip150Block": 0, 
    "eip155Block": 0, 
    "eip160Block": 0, 
    "eip161FBlock": 0, 
    "eip170FBlock": 0, 
    "eip100FBlock": 0, 
    "eip140FBlock": 0, 
    "eip198FBlock": 0, 
    "eip211FBlock": 0, 
    "eip212FBlock": 0, 
    "eip213FBlock": 0, 
    "eip214FBlock": 0, 
    "eip658FBlock": 0, 
    "eip145FBlock": 0, 
    "eip1014FBlock": 0, 
    "eip1052FBlock": 0, 
    "eip152FBlock": 0, 
    "eip1108FBlock": 0, 
    "eip1344FBlock": 0, 
    "eip1884FBlock": 0, 
    "eip2028FBlock": 0, 
    "eip2200FBlock": 0, 
    "disposalBlock": 0, 
    "ethash": {}, 
    "requireBlockHashes": {}
  }, 
  "nonce": "0x0", 
  "timestamp": "0x6048d57c", 
  "extraData": "0x42", 
  "gasLimit": "0x2fefd8", 
  "difficulty": "0x20000", 
  "mixHash": "0x0000000000000000000000000000000000000000000000000000000000000000", 
  "coinbase": "0x0000000000000000000000000000000000000000", 
  "alloc": {
    "366ae7da62294427c764870bd2a460d7ded29d30": {
      "balance": "0x2a"
    }
  }, 
  "number": "0x0", 
  "gasUsed": "0x0", 
  "parentHash": "0x0000000000000000000000000000000000000000000000000000000000000000"
}
EOF

```

Initialize geth with this configuration (using a custom data directory).
```
./build/bin/geth --datadir=./abc-datadir init abc_genesis.json 
INFO [03-10|09:00:25.710] Maximum peer count                       ETH=50 LES=0 total=50
INFO [03-10|09:00:25.710] Smartcard socket not found, disabling    err="stat /run/pcscd/pcscd.comm: no such file or directory"
INFO [03-10|09:00:25.711] Set global gas cap                       cap=25000000
INFO [03-10|09:00:25.711] Allocated cache and file handles         database=/home/ia/go/src/github.com/ethereum/go-ethereum/abc-datadir/geth/chaindata cache=16.00MiB handles=16
INFO [03-10|09:00:25.728] Writing custom genesis block 
INFO [03-10|09:00:25.729] Persisted trie from memory database      nodes=1 size=139.00B time="178.942µs" gcnodes=0 gcsize=0.00B gctime=0s livenodes=1 livesize=0.00B
INFO [03-10|09:00:25.729] Wrote custom genesis block OK            config="NetworkID: 4269, ChainID: 4269 Engine: ethash EIP1014: 0 EIP1052: 0 EIP1108: 0 EIP1344: 0 EIP140: 0 EIP145: 0 EIP150: 0 EIP152: 0 EIP155: 0 EIP160: 0 EIP161abc: 0 EIP161d: 0 EIP170: 0 EIP1884: 0 EIP198: 0 EIP2028: 0 EIP211: 0 EIP212: 0 EIP213: 0 EIP214: 0 EIP2200: 0 EIP2: 0 EIP658: 0 EIP7: 0 EthashECIP1041: 0 EthashEIP100B: 0 EthashHomestead: 0 "
INFO [03-10|09:00:25.730] Successfully wrote genesis state         database=chaindata hash="5f32ce…1fe582"
INFO [03-10|09:00:25.730] Allocated cache and file handles         database=/home/ia/go/src/github.com/ethereum/go-ethereum/abc-datadir/geth/lightchaindata cache=16.00MiB handles=16
INFO [03-10|09:00:25.746] Writing custom genesis block 
INFO [03-10|09:00:25.747] Persisted trie from memory database      nodes=1 size=139.00B time="91.084µs"  gcnodes=0 gcsize=0.00B gctime=0s livenodes=1 livesize=0.00B
INFO [03-10|09:00:25.748] Wrote custom genesis block OK            config="NetworkID: 4269, ChainID: 4269 Engine: ethash EIP1014: 0 EIP1052: 0 EIP1108: 0 EIP1344: 0 EIP140: 0 EIP145: 0 EIP150: 0 EIP152: 0 EIP155: 0 EIP160: 0 EIP161abc: 0 EIP161d: 0 EIP170: 0 EIP1884: 0 EIP198: 0 EIP2028: 0 EIP211: 0 EIP212: 0 EIP213: 0 EIP214: 0 EIP2200: 0 EIP2: 0 EIP658: 0 EIP7: 0 EthashECIP1041: 0 EthashEIP100B: 0 EthashHomestead: 0 "
INFO [03-10|09:00:25.749] Successfully wrote genesis state         database=lightchaindata hash="5f32ce…1fe582"
```

Start geth, reusing our initialized database.
Since geth won't have default bootnodes for this configuration (only available when using CLI flags), we'll need to
use geth's `--bootnodes` flag.
```
./build/bin/geth --datadir=./abc-datadir --bootnodes=enode://3e12c4c633157ae52e7e05c168f4b1aa91685a36ba33a0901aa8a83cfeb84c3633226e3dd2eaf59bfc83492139e1d68918bf5b60ba93e2deaedb4e6a2ded5d32@42.152.120.98:30303
INFO [03-10|09:07:52.762] Starting Geth on Ethereum mainnet... 
INFO [03-10|09:07:52.762] Bumping default cache on mainnet         provided=1024 updated=4096
INFO [03-10|09:07:52.763] Maximum peer count                       ETH=50 LES=0 total=50
INFO [03-10|09:07:52.763] Smartcard socket not found, disabling    err="stat /run/pcscd/pcscd.comm: no such file or directory"
INFO [03-10|09:07:52.764] Set global gas cap                       cap=25000000
INFO [03-10|09:07:52.764] Allocated trie memory caches             clean=1023.00MiB dirty=1024.00MiB
INFO [03-10|09:07:52.765] Allocated cache and file handles         database=/home/ia/go/src/github.com/ethereum/go-ethereum/abc-datadir/geth/chaindata cache=2.00GiB handles=524288
INFO [03-10|09:07:52.853] Opened ancient database                  database=/home/ia/go/src/github.com/ethereum/go-ethereum/abc-datadir/geth/chaindata/ancient
INFO [03-10|09:07:52.854] Found stored genesis block               config="NetworkID: 4269, ChainID: 4269 Engine: ethash EIP1014: 0 EIP1052: 0 EIP1108: 0 EIP1344: 0 EIP140: 0 EIP145: 0 EIP150: 0 EIP152: 0 EIP155: 0 EIP160: 0 EIP161abc: 0 EIP161d: 0 EIP170: 0 EIP1884: 0 EIP198: 0 EIP2028: 0 EIP211: 0 EIP212: 0 EIP213: 0 EIP214: 0 EIP2200: 0 EIP2: 0 EIP658: 0 EIP7: 0 EthashECIP1041: 0 EthashEIP100B: 0 EthashHomestead: 0 "
INFO [03-10|09:07:52.854] Found non-defaulty stored config, using it. 
INFO [03-10|09:07:52.854] Initialised chain configuration          config="NetworkID: 4269, ChainID: 4269 Engine: ethash EIP1014: 0 EIP1052: 0 EIP1108: 0 EIP1344: 0 EIP140: 0 EIP145: 0 EIP150: 0 EIP152: 0 EIP155: 0 EIP160: 0 EIP161abc: 0 EIP161d: 0 EIP170: 0 EIP1884: 0 EIP198: 0 EIP2028: 0 EIP211: 0 EIP212: 0 EIP213: 0 EIP214: 0 EIP2200: 0 EIP2: 0 EIP658: 0 EIP7: 0 EthashECIP1041: 0 EthashEIP100B: 0 EthashHomestead: 0 "
INFO [03-10|09:07:52.854] Disk storage enabled for ethash caches   dir=/home/ia/go/src/github.com/ethereum/go-ethereum/abc-datadir/geth/ethash count=3
INFO [03-10|09:07:52.854] Disk storage enabled for ethash DAGs     dir=/home/ia/.ethash count=2
INFO [03-10|09:07:52.854] Initialising Ethereum protocol           versions="[65 64 63]" network=1 dbversion=8
INFO [03-10|09:07:52.855] Loaded most recent local header          number=0 hash="5f32ce…1fe582" td=131072 age=48m12s
INFO [03-10|09:07:52.855] Loaded most recent local full block      number=0 hash="5f32ce…1fe582" td=131072 age=48m12s
INFO [03-10|09:07:52.855] Loaded most recent local fast block      number=0 hash="5f32ce…1fe582" td=131072 age=48m12s
INFO [03-10|09:07:52.855] Loaded local transaction journal         transactions=0 dropped=0
INFO [03-10|09:07:52.856] Regenerated local transaction journal    transactions=0 accounts=0
INFO [03-10|09:07:52.877] Allocated fast sync bloom                size=2.00GiB
INFO [03-10|09:07:52.878] Initialized fast sync bloom              items=1 errorrate=0.000 elapsed="520.228µs"
INFO [03-10|09:07:52.879] Starting peer-to-peer node               instance=CoreGeth/v1.11.22-stable-72df266d/linux-amd64/go1.16
INFO [03-10|09:07:52.898] New local node record                    seq=3 id=0a86440c3ab5e22c ip=127.0.0.1 udp=30303 tcp=30303
INFO [03-10|09:07:52.899] Started P2P networking                   self=enode://5256dcfe7725a98f38cf15b702847fabcaf59bbaa733a6ae5ea68e1089fdd1d274192e17593dc20df00a45ea91372f7c1ca97c8d186fa9e779167240fde15338@127.0.0.1:30303
INFO [03-10|09:07:52.900] IPC endpoint opened                      url=/home/ia/go/src/github.com/ethereum/go-ethereum/abc-datadir/geth.ipc
INFO [03-10|09:07:52.905] Mapped network port                      proto=udp extport=30303 intport=30303 interface=NAT-PMP(192.168.86.1)
INFO [03-10|09:07:52.909] Mapped network port                      proto=tcp extport=30303 intport=30303 interface=NAT-PMP(192.168.86.1)
INFO [03-10|09:07:54.403] New local node record                    seq=4 id=0a86440c3ab5e22c ip=75.134.144.252 udp=30303 tcp=30303
INFO [03-10|09:08:06.969] Looking for peers                        peercount=0 tried=5 static=0
```

#### Establish a network.

In order to establish your network, you'll want to make sure you have a bootnode
available that new nodes coming online can use to query for their peers.

##### Set up a bootnode.

Initialize the bootnode's database and get its self-reported `enode` value. 
```
./build/bin/geth --datadir=./abc-datadir init abc_genesis.json

2>/dev/null ./build/bin/geth --datadir=./abc-datadir --exec 'admin.nodeInfo.enode' console
"enode://5256dcfe7725a98f38cf15b702847fabcaf59bbaa733a6ae5ea68e1089fdd1d274192e17593dc20df00a45ea91372f7c1ca97c8d186fa9e779167240fde15338@75.134.144.252:30303"
```

This (`enode://5256dcfe7725a98f38cf15b702847fabcaf59bbaa733a6ae5ea68e1089fdd1d274192e17593dc20df00a45ea91372f7c1ca97c8d186fa9e779167240fde15338@75.134.144.252:30303`)
will be the bootnode `enode` value for the other nodes.

Then turn the bootnode on.
```
./build/bin/geth --datadir=./abc-datadir
```

##### Start up a few nodes.

```
./build/bin/geth --datadir=./abc-datadir-1 init abc_genesis.json
./build/bin/geth --datadir=./abc-datadir-2 init abc_genesis.json
./build/bin/geth --datadir=./abc-datadir-3 init abc_genesis.json
```

```
./build/bin/geth --datadir=./abc-datadir-1 --bootnodes=enode://5256dcfe7725a98f38cf15b702847fabcaf59bbaa733a6ae5ea68e1089fdd1d274192e17593dc20df00a45ea91372f7c1ca97c8d186fa9e779167240fde15338@75.134.144.252:30303
```

```
./build/bin/geth --datadir=./abc-datadir-2 --bootnodes=enode://5256dcfe7725a98f38cf15b702847fabcaf59bbaa733a6ae5ea68e1089fdd1d274192e17593dc20df00a45ea91372f7c1ca97c8d186fa9e779167240fde15338@75.134.144.252:30303
```

```
./build/bin/geth --datadir=./abc-datadir-3 --bootnodes=enode://5256dcfe7725a98f38cf15b702847fabcaf59bbaa733a6ae5ea68e1089fdd1d274192e17593dc20df00a45ea91372f7c1ca97c8d186fa9e779167240fde15338@75.134.144.252:30303
```
