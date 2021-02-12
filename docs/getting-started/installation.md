## Pre-built executable

If you just want to download and run `geth` or any of the other tools here, this is the quickest and simplest way.

Binary archives are published at https://github.com/etclabscore/core-geth/releases. Find the latest one for your OS, download it, (check the SHA sum), unarchive it, and run!

## With Docker

All runnable examples below are for images limited to `geth`. For images including the full suite of
tools available from this source, use the Docker Hub tag prefix `alltools.`, like `etclabscore/core-geth:alltools.latest`, or the associated Docker file directly `./Dockerfile.alltools`.

### `docker run`

One of the quickest ways to get Ethereum Classic up and running on your machine is by using Docker:

```shell
$ docker run -d \
    --name core-geth \
    -v $LOCAL_DATADIR:/root \
    -p 30303:30303 \
    -p 8545:8545 \
    etclabscore/core-geth \
    --classic \
    --rpc --rpcport 8545
```

This will start `geth` in fast-sync mode with a DB memory allowance of 1GB just as the
above command does.  It will also create a persistent volume in your `$LOCAL_DATADIR` for
saving your blockchain, as well as map the default devp2p and JSON-RPC API ports.

Do not forget `--http.addr 0.0.0.0`, if you want to access RPC from other containers
and/or hosts. By default, `geth` binds to the local interface and RPC endpoints is not
accessible from the outside.


### `docker pull`

Docker images are automatically [published on Docker Hub](https://hub.docker.com/r/etclabscore/core-geth/tags).

#### Image: `latest`

Image `latest` is built automatically from the `master` branch whenever it's updated.

```shell
$ docker pull etclabscore/core-geth:latest
```

#### Image: `<tag>`

Repository tags like `v1.2.3` correspond to Docker tags like __`version-1.2.3`__.

!!! Example

    ```shell
    $ docker pull etclabscore/core-geth:version-1.11.1
    ```
