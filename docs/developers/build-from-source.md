## Dependencies

- Make sure your system has __Go__ installed. Version 1.15+ is recommended. https://golang.org/doc/install
- Make sure your system has a C compiler installed. For example, with Linux Ubuntu:

```shell
$ sudo apt-get install -y build-essential
```

## Source

Once the dependencies have been installed, it's time to clone and build the source:

```shell
$ git clone https://github.com/etclabscore/core-geth.git
$ cd core-geth
$ make all
$ ./build/bin/geth --help
```

## Build docker image

You can build a local docker image directly from the source:

```shell
$ git clone https://github.com/etclabscore/core-geth.git
$ cd core-geth
$ docker build -t=core-geth .
```

Or with all tools:

```shell
$ docker build -t core-geth-alltools -f Dockerfile.alltools .
```