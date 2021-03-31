# "trace" Module Overview

The trace module is for getting a deeper insight into transaction processing. It includes two sets of calls; the transaction trace filtering API and the ad-hoc tracing API. You can find the documentation for the supported methods [here](/JSON-RPC-API/modules/trace/).

It's good to mention that `trace_*` methods are nothing more than aliases to some existing `debug_*` methods. The reason for creating those aliases, was to reach compatibility with OpenEthereum's (aka Parity) trace module, which has been requested by the community in order they can fully use core-geth. For achieving this, the `trace_*` methods set the default tracer to `callTracerParity` if none is set.

!!! Note "Full sync"
    In order to use the Transaction-Trace Filtering API, core-geth must be fully synced using `--syncmode=full --gcmode=archive`. Otherwise, you can set the number of blocks to `reexec` back for rebuilding the state, though taking longer for a trace call to finish.

## JSON-RPC methods

### Ad-hoc Tracing

The ad-hoc tracing API allows you to perform a number of different diagnostics on calls or transactions, either historical ones from the chain or hypothetical ones not yet mined.

- [x] trace_call *(alias to debug_traceCall)*
- [x] trace_callMany
- [ ] trace_rawTransaction
- [ ] trace_replayBlockTransactions
- [ ] trace_replayTransaction

### Transaction-Trace Filtering

These APIs allow you to get a full externality trace on any transaction executed throughout the blockchain.

- [x] trace_block *(alias to debug_traceBlock)*
- [x] trace_transaction *(alias to debug_traceTransaction)*
- [x] trace_filter (doesn't support address filtering yet)
- [ ] trace_get

## Available tracers

- `callTracerParity` Transaction trace returning a response equivalent to OpenEthereum's (aka Parity) response schema. For documentation on this response value see [here](#calltracerparity).
- `vmTrace` Virtual Machine execution trace. Provides a full trace of the VM’s state throughout the execution of the transaction, including for any subcalls. *(Not implemented yet)*
- `stateDiffTracer` State difference. Provides information detailing all altered portions of the Ethereum state made due to the execution of the transaction. For documentation on this response value see [here](#statedifftracer).

!!! Example "Example trace_* API method config (last method argument)"

    ```js
    {
        "tracer": "stateDiffTracer",
        "timeout: "10s",
        "reexec: "10000",               // number of block to reexec back for calculating state
        "nestedTraceOutput": true  // in Ad-hoc Tracing methods the response is nested similar to OpenEthereum's output
    }
    ```

### Tracers' output documentation

#### callTracerParity

The output result is an array including the outer transaction (first object in the array), as well the internal transactions (next objects in the array) that were being triggered.

Each object that represents an internal transaction consists of:

* the `action` object with all the call args,
* the `resutls` object with the outcome as well the gas used,
* the `subtraces` field, representing the number of internal transactions that were being triggered by the current transaction,
* the `traceAddress` field, representing the exact nesting location in the call trace *[index in root, index in first CALL, index in second CALL, …]*.

```js
[
    {
        "action": {
            "callType": "call",
            "from": "0x877bd459c9b7d8576b44e59e09d076c25946f443",
            "gas": "0x1e30e8",
            "input": "0xb595b8b50000000000000000000000000000000000000000000000000000000000000000",
            "to": "0x5e0fddd49e4bfd02d03f2cefa0ea3a3740d1bb3d",
            "value": "0xde0b6b3a7640000"
        },
        "result": {
            "gasUsed": "0x25e9",
            "output": "0x00000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000031436173696e6f2068617320696e73756666696369656e742066756e647320666f7220746869732062657420616d6f756e74000000000000000000000000000000"
        },
        "subtraces": 1,
        "traceAddress": [],
        "type": "call"
    },
    {
        "action": {
            "callType": "call",
            "from": "0x5e0fddd49e4bfd02d03f2cefa0ea3a3740d1bb3d",
            "gas": "0x8fc",
            "input": "0x",
            "to": "0x877bd459c9b7d8576b44e59e09d076c25946f443",
            "value": "0xde0b6b3a7640000"
        },
        "result": {
            "gasUsed": "0x0",
            "output": "0x"
        },
        "subtraces": 0,
        "traceAddress": [
            0
        ],
        "type": "call"
    }
]
```

#### stateDiffTracer

Provides information detailing all **altered portions** of the Ethereum state made due to the execution of the transaction.

Each address object provides the state differences for `balance`, `nonce`, `code` and `storage`.
Actually, under the `storage` object, we can find the state differences for each contract's storage key.

**Special symbols** explanation:

* `+`, when we have a new entry in the state DB,
* `-`, when we have a removal from the state DB,
* `*`, when existing data have changed in the state DB, providing the `from` (old) and the `to` (new) values,
* `=`, when the data remained the same.

```js
{
    "0x877bd459c9b7d8576b44e59e09d076c25946f443": {
        "balance": {
            "*": {
                "from": "0xd062abd70db4255a296",
                "to": "0xd062ac59cb1bd516296"
            }
        },
        "nonce": {
            "*": {
                "from": "0x1c7ff",
                "to": "0x1c800"
            }
        },
        "code": "=",
        "storage": {
            "0x0000000000000000000000000000000000000000000000000000000000000001": {
                "*": {
                    "from": "0x0000000000000000000000000000000000000000000000000000000000000000",
                    "to": "0x0000000000000000000000000000000000000000000000000000000000000061"
                }
              },
        }
    },
    ...
}
```

## "stateDiff" tracer differences with OpenEthereum

1. **SSTORE** in some edge cases persists data in state but are not being returned on stateDiff storage results on OpenEthereum output.
   > Happens only on 2 transactions on **Mordor** testnet, as of **block 2,519,999**. (TX hashes: *0xab73afe7b92ad9b537df3f168de0d06f275ed34edf9e19b36362ac6fa304c0bf*, *0x15a7c727a9bbfdd43d09805288668cc4a0ec647772d717957e882a71ace80b1a*)
2. When error **ErrInsufficientFundsForTransfer** happens, **OpenEthereum** leaves the tracer run producing negative balances, though using safe math for overflows it **returns 0 balance**, on the other hand the `to` account **receives the full amount**.
**Core-geth removes only the gas cost** from the sender and **adds it to the coinbase balance**.
3. Same as in 2, but on top of that, the **sender account doesn't have to pay for the gas cost** even. In this case, **core-geth returns an empty JSON**, as in reality this transaction will remain in the tx_pool and never be executed, neither change the state.
4. On **OpenEthereum the block gasLimit is set to be U256::max()**, which leads into problems on contracts using it for pseudo-randomness. On **core-geth**, we believe that the user utilising the trace_* wants to **see what will happen in reality**, though we **leave the block untouched to its true values**.
5. When an internal call fails with out of gas, and its state is not being persisted, we don't add it in stateDiff output, as it happens on OpenEthereum.