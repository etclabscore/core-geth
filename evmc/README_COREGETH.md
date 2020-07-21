
This is a vendor fork of github.com/ethereum/evmc@v6.3.1 (e9d4648200d73c1ce63b74a33a50758db1fd48db) with only
very slight tweaks to the tests.  This version is the latest and last of the v6 series.

Changes made locally ("slight tweaks") are in d1a16da676812c95d2f2cd901ef16024a665ba29, which causes
tests to be skipped if the required `example_vm.so` is not available. This artifact
is generated, and can be generated with the `go:generate` command found in 
`binding/go/evmc/evmc_test.go`. If it is available, the tests will run as originally.

The rationale for this fork and the skipped tests is:
- As drafted in [this PR comment](https://github.com/etclabscore/core-geth/pull/57#issuecomment-656170790),
  installing the `evmc` package a submodule added undesirable complexity for build and test steps, namely
  requiring the submodule for build (extra `git` steps) and requiring the submodule and a generated artifact
  for testing -- yet another required extra step. In preference of conventional and "works-out-of-the-box"
  experiences and dependencies for developers, consumers, and CI jobs, it is better to simply include the
  vendored dependency as a slightly-modified and static copy of the source.
  + For precedent in this pattern (and as noted in the comment linked above), the local `log` package
    is a vendored source mostly-copy of `github.com/inconshreveable/log15`. 
- Tests may be skipped since these are tests incumbent to the vendored package. Their affirmation is (or, was)
  not the responsibility of core-geth originally, and
- Since this vendor version is expected to be eternally pegged (immutable), we can expect the outcome
  of the presumed (and affirmed in development) passing tests would be static and unresponsive to changes
  in the core-geth project.
- Further, since the tests are only run if the generated artifact is _not_ available, we can
  and do expect to generate the artifact and actually run the tests in CI, since the overhead is just
  mostly configuration and dependency install code and that's already written and available.

Easier for developers and consumers, all the tests for CI... win win!





