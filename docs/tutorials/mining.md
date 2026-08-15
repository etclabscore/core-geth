# GPU Mining with Core-Geth



```sh
./build/bin/geth version
CoreGeth
Version: 1.12.9-unstable
Git Commit: 5bcbb9dae506d82d1cb7af9d53c66986501d0f57
Git Commit Date: 20220801
Architecture: amd64
Go Version: go1.18.3
Operating System: linux
GOPATH=/home/ia/go
GOROOT=/home/ia/go1.18.3.linux-amd64

./build/bin/geth \
    --mordor \
    --datadir ~/SAMSUNG_T5_3/chaindata/mordor \
    --ipcpath /tmp/geth.ipc  \
    --http \
    --mine \
    --miner.etherbase 0x31C488d65472a8c0345e5bBd928730a54d29bF37 \
    --miner.notify http://localhost:8888
```

- https://github.com/cyberpoolorg/etc-stratum
- [etc-stratum_config_mordor_ia.json ](https://gist.github.com/meowsbits/f1a28a8b11534c0a2536af871796bc5a)
```sh
./build/bin/etc-stratum etc-stratum_config_mordor_ia.json 
```


- https://trex-miner.com/ T-Rex worked on Mordor, while GMiner and lolMiner both got the epoch wrong (eg. 230 instead of 115). Having the wrong epoch caused them to produce invalid shares and blocks. This issue was only with Mordor; I've used both GMiner and lolMiner to mine ETC mainnet successfully.

```sh
./t-rex-0.26.4-linux/t-rex \
    -i 9 \
    -a etchash \
    -o stratum+tcp://0.0.0.0:8008 \
    -u 0x31c488d65472a8c0345e5bbd928730a54d29bf37 \
    -p x \
    -w rig0
```
