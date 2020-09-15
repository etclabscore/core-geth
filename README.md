### Develop and Build

Update static assets packaging (so cmd/ageth/index.html ships as a part of `ageth` executable).
If you don't update index.html, you don't need to do this.

```sh
statik -src ./cmd/ageth
```

Clean out all vestiges from previous ageth runs using local geths.
If you're using local geths, these will be data directories, sockets, and keystores.
If you're using remote geths, this is unnecessary.
```sh
rm -fr /tmp/ageth*
```

Build everything.
```sh
make all
```

### Run

Run `ageth` using stdin to provide an endpoint list definition.
```sh
printf './build/bin/geth\n./build/bin/geth\n./build/bin/geth\n./build/bin/geth\n./build/bin/geth\n./build/bin/geth' | ./build/bin/ageth > ageth.log
```

Run `ageth` using a file from which to read newline-delimited endpoint definitions.

```
./build/bin/ageth -f fake-endpoints.txt > ageth.log
``` 

This example redirects `ageth`'s stdout to a logfile. If ageth is running local geths, these logs
will be the aggregated logs of all the geths she's running.
If `ageth` is using remote geths, this is unnecessary; she doesn't have access to remote logs and her stdout will be empty.

- [ ] TODO: Handle log placement better. The test `scenarios` will produce their own logs, which probably should go to
stdout. This will collide with the local geth logs, if any.

### "Endpoint" definitions

__First__: It's a little obtuse, but endpoints can _either_ be
- `ws://` endpoints (websockets), OR
- `./build/bin/geth` "endpoints".

This is where the scare-quotes come in.

If the value doesn't have a valid URL `scheme`, she'll consider it as a path to a local `geth` exectuable. She'll attempt to start it, and connect over IPC.
If, on the other hand, the value _does_ have valid URL `scheme`, she'll just connect with it by dialing the RPC URL. 

__If remote, the endpoint MUST be a ws:// (websocket) URL.__ `ageth` needs a websocket so she can _subscribe_ to the node's head events.

`ageth` eats "endpoints" for geths, and

- if they're _local_, she can turn them on and aggregate their logs
- if they're _remote_, she'll just wrap them up in RPC clients and connections

An endpoint file may looks as follows. This will yield an `ageth` managing 4 remote nodes and 2 local ones.

```txt
ws://somewhere.com:8546
ws://somewhereelse.com:8546
ws://wherethough.com:8546
ws://anywhere.com:8546
./build/bin/geth
./build/bin/geth
```

As above, if you have this :point_up: in a file `ageth-endpoints.txt`, you'll use:

```sh
./build/bin/ageth -f ageth-endpoints.txt
```

OR

```sh
cat ageth-endpoints | ./build/bin/ageth
```

### Writing scenarios

The "Scenario" interface is defined as

```
type scenario func(nodes *agethSet)
```

You write these, and then add them to the root scenario loop:

```go
for {
    for i, s := range []scenario{
        scenario5,
    }{
        log.Info("Running scenario", "index", i)
        globalTick = 0
        s(world)
    }
}
```

An example scenario using this pattern is provided in [./cmd/ageth/cmd/scenario_5.go](./cmd/ageth/cmd/scenario_5.go).


### Observer

From startup, `ageth` hosts an "Observer" website at `:8008`. 
