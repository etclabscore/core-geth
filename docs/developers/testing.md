# Testing

## Lint

- `make lint` runs the module-standard linter config. Under the hood, this uses `golangci-lint` and its configuration at `/.golangci.yml`.

## Run tests

The `Makefile` at the root of the project includes many commands that can be useful for testing.

!!! note "_make_ is not _required_"
    Please keep in mind that while you _can_ use `make` to run tests, you don't have to.
    You can also use plain old `go test` (and `go build`) as you would in any Go project.


- `make test` runs all package tests which are, for the most part, _shared with go-ethereum_.
- `make test-coregeth` runs a suite of tests that are specific to CoreGeth.


## Test generation

CoreGeth is capable of generating some sets of tests used in the `tests` package, which are originally (and still largely)
driven by the [ethereum/tests](https://github.com/ethereum/tests) suite.

- `make tests-generate` runs test(s) generation for the `state` and `difficulty` subsections of this suite, extending the ethereum/tests version
  of the controls to include configurations for Ethereum Classic chain configurations at various points in Ethereum Classic hardfork history.

## Flaky (spuriously erroring) tests

Especially when run in CI environments, some tests can be expected to fail more-or-less randomly.


