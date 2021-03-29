# "trace" Module Overview

The trace module is for getting a deeper insight into transaction processing. It includes two sets of calls; the transaction trace filtering API and the ad-hoc tracing API. You can find the documentation for the supported methods [here](/JSON-RPC-API/modules/trace/).

!!! Note "Full sync"
    In order to use the Transaction-Trace Filtering API, core-geth must be fully synced using `--syncmode=full --gcmode=archive`. Otherwise, you can set the number of blocks to `reexec` back for rebuilding the state, though taking longer for a trace call to finish.

## JSON-RPC methods

### Ad-hoc Tracing

The ad-hoc tracing API allows you to perform a number of different diagnostics on calls or transactions, either historical ones from the chain or hypothetical ones not yet mined.

- [x] trace_call
- [x] trace_callMany
- [ ] trace_rawTransaction
- [ ] trace_replayBlockTransactions
- [ ] trace_replayTransaction

### Transaction-Trace Filtering

These APIs allow you to get a full externality trace on any transaction executed throughout the blockchain.

- [x] trace_block
- [x] trace_transaction
- [x] trace_filter (doesn't support address filtering yet)
- [ ] trace_get

### Available tracers

- `callTracerParity` Transaction trace.
- `vmTrace` Virtual Machine execution trace. Provides a full trace of the VM’s state throughout the execution of the transaction, including for any subcalls. *(Not implemented yet)*
- `stateDiffTracer` State difference. Provides information detailing all altered portions of the Ethereum state made due to the execution of the transaction.

!!! Example "Example trace_* API method config (last method argument)"

    ```js
    {
        "tracer": "stateDiffTracer",
        "timeout: "10s",
        "reexec: "10000",               // number of block to reexec back for calculating state
        "parityCompatibleOutput": true  // for Ad-hoc Tracing only in order the response is nested as in OpenEthereum output
    }
    ```

## "stateDiff" tracer differences with OpenEthereum

1. **SSTORE** in some edge cases persists data in state but are not being returned on stateDiff storage results on OpenEthereum output. *(Happens only on 2 transactions on Mordor testnet)*.
2. When error **ErrInsufficientFundsForTransfer** happens, **OpenEthereum** leaves the tracer run producing negative balances, though using safe math for overflows it **returns 0 balance**, on the other hand the `to` account **receives the full amount**.
**Core-geth removes only the gas cost** from the sender and **adds it to the coinbase balance**.
3. Same as in 2, but on top of that, the **sender account doesn't have to pay for the gas cost** even. In this case, **core-geth returns an empty JSON**, as in reality this transaction will remain in the tx_pool and never be executed, neither change the state.
4. On **OpenEthereum the block gasLimit is set to be U256::max()**, which leads into problems on contracts using it for pseudo-randomness. On **core-geth**, we believe that the user utilising the trace_* wants to **see what will happen in reality**, though we **leave the block untouched to its true values**.
5. When an internal call fails with out of gas, and its state is not being persisted, we don't add it in stateDiff output, as it happens on OpenEthereum.