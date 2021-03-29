---
title: OpenRPC Discovery
hide:
  - toc        # Hide table of contents
---

# OpenRPC

!!! tldr "TLDR"
    The `rpc.discover` method returns an API service description structured per the [OpenRPC specification](https://github.com/open-rpc/spec/blob/master/spec.md). 

## Discovery

CoreGeth supports [OpenRPC's Service Discovery method](https://spec.open-rpc.org/#service-discovery-method), enabling efficient and well-spec'd JSON RPC interfacing and tooling. This method follows the established JSON RPC patterns, and is accessible via HTTP, WebSocket, IPC, and console servers. To use this method:

!!! Example

    ```shell
    $ curl -X POST -H 'Content-Type: application/json' --data '{"jsonrpc":"2.0","method":"rpc_discover","params":[],"id":1}'
    {
      "jsonrpc": "2.0",
      "id": 1,
      "result": {
        "openrpc": "1.0.10",
        "info": {
          "description": "This API lets you interact with an EVM-based client via JSON-RPC",
          "license": {
            "name": "Apache 2.0",
            "url": "https://www.apache.org/licenses/LICENSE-2.0.html"
          },
          "title": "Ethereum JSON-RPC",
          "version": "1.0.0"
        },
        "servers": [],
        "methods": [
          {
            "description": "Returns the version of the current client",
            "name": "web3_clientVersion",
            "params": [],
            "result": {
              "description": "client version",
              "name": "clientVersion",
              "schema": {
                "type": "string"
              }
            },
            "summary": "current client version"
          },

    [...]
    ```

!!! Tip "Better representation of the discovery result at the OpenRPC playground"

    You can see an example use case of the discovery service in the respective [OpenRPC Playground](https://playground.open-rpc.org/?schemaUrl=https://gist.githubusercontent.com/ziogaschr/c51916d70ca5304bb3e3abf4dcd518ca/raw/8079eafd8de6436bd3e4ab6c9df0db64c25cd1a6/core-geth_rpc-discovery_1.11.21-unstable.json).
