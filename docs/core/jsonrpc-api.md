### OpenRPC Discovery

CoreGeth supports [OpenRPC's Service Discovery method](https://spec.open-rpc.org/#service-discovery-method), enabling efficient and well-spec'd JSON RPC interfacing and tooling. This method follows the established JSON RPC patterns, and is accessible via HTTP, WebSocket, IPC, and console servers. To use this method:
```shell
$ curl -X POST -H 'Content-Type: application/json' --data '{"jsonrpc":"2.0","method":"rpc.discover","params":[],"id":1}'
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
