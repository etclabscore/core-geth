# Running Geth with an External VM

Geth supports the [EVMC](https://github.com/ethereum/evmc/) VM connector API version 6 as an experimental feature.

From [PR #57](https://github.com/etclabscore/core-geth/pull/57) geth enables an externally defined VM, either EVM or EWASM, via
a `--vm.`-prefixed flag for normal instantiation, and `--evmc.` for testing.

Options include EWASM and EVM shared object libraries, as follows:

- `--vm.ewasm=<path/to/interpreter.so`
- `--vm.evm=<path/to/interpreter.so`

Only EVMC version __6__ is supported, which is compatible with the latest versions of Hera EWASM v0.2.5 and SSVM EWASM 0.5.0.

This implementation may be tested by following the command defined in the Makefile as `evmc-test`, which
tests the implementation against both of these mentioned EWASM libraries against the `/tests/` StateTest suite.
