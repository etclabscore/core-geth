# Running Geth with an External VM

Geth supports the [EVMC](https://github.com/ethereum/evmc/) VM connector API version 6 as an experimental feature.

From [PR #57](https://github.com/etclabscore/core-geth/pull/57) geth enables an externally defined VM, either EVM or EWASM, via
a `--vm.`-prefixed flag for normal instantiation, and `--evmc.` for testing.

Options include EWASM and EVM shared object libraries, as follows:

- `--vm.ewasm=<path/to/interpreter.so`
- `--vm.evm=<path/to/interpreter.so`

Only EVMC __Version 6__ is supported, which is compatible with the versions of Hera EWASM <=v0.2.5 and SSVM EWASM >=0.5.0.

This implementation may be tested by following the command defined in the Makefile as `evmc-test`, which
tests the implementation against both of these mentioned EWASM libraries against the `/tests/` StateTest suite.

While core-geth supports highly granular EIP/ECIP/xIP chain feature configuration (ie fork feature configs),
EVMC does not. EVMC only supports the Fork configurations supported by ethereum/go-ethereum (eg. Byzantium, Constantinople, &c). 
Thus, the implementation at core-geth of EVMC requires a somewhat arbitrary mapping of granular features as keys toggling
entire Ethereum fork configurations.

The following code snippet, taken from `./core/vm/evmc.go`, handles this translation.

```go
// getRevision translates ChainConfig's HF block information into EVMC revision.
func getRevision(env *EVM) evmc.Revision {
	n := env.BlockNumber
	conf := env.ChainConfig()
	switch {
	// This is an example of choosing to use an "abstracted" idea
	// about chain config, where I'm choosing to prioritize "indicative" features
	// as identifiers for Fork-Feature-Groups. Note that this is very different
	// than using Feature-complete sets to assert "did Forkage."
	case conf.IsEnabled(conf.GetEIP1884Transition, n):
		return evmc.Istanbul
	case conf.IsEnabled(conf.GetEIP1283DisableTransition, n):
		return evmc.Petersburg
	case conf.IsEnabled(conf.GetEIP145Transition, n):
		return evmc.Constantinople
	case conf.IsEnabled(conf.GetEIP198Transition, n):
		return evmc.Byzantium
	case conf.IsEnabled(conf.GetEIP155Transition, n):
		return evmc.SpuriousDragon
	case conf.IsEnabled(conf.GetEIP150Transition, n):
		return evmc.TangerineWhistle
	case conf.IsEnabled(conf.GetEIP7Transition, n):
		return evmc.Homestead
	default:
		return evmc.Frontier
	}
}
```

As you can see, individual features, like EIP1884, are translated as proxy signifiers for entire fork configurations
(in this case, an Istanbul-featured VM revision).
This approach, rather than requiring a complete set of the compositional features for any of these given Ethereum forks,
trades a descriptive 1:1 mapping for application flexibility. Pursuing a necessarily complete feature-set -> fork
map would presume chain features that are not necessarily relevant to the virtual machine, like block reward configurations,
or difficulty configurations, for example. This approach allows applications to use advanced opcodes with the fewest
number of incidental restrictions.

This approach is not without risk or nuance however; without a solid understanding of customizations here,
experiments in customization can result in foot shooting.
