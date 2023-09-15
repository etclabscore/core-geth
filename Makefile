# This Makefile is meant to be used by people that do not usually work
# with Go source code. If you know what GOPATH is then you probably
# don't need to bother with make.

.PHONY: geth android ios geth-cross evm evmc mkdocs-serve all test clean

GOBIN = ./build/bin
GO ?= latest
GORUN = go run
ROOT_DIR := $(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

geth:
	$(GORUN) build/ci.go install ./cmd/geth
	@echo "Done building."
	@echo "Run \"$(GOBIN)/geth\" to launch geth."

all:
	$(GORUN) build/ci.go install

test: all
	$(GORUN) build/ci.go test -timeout 20m

# DEPRECATED.
# No attempt will be made after the Istanbul fork to maintain
# Parity configuration support.
sync-parity-chainspecs:
	./params/parity.json.d/sync-parity-remote.sh

test-coregeth: \
 test-coregeth-features \
 test-coregeth-consensus \
 test-coregeth-regression-condensed ## Runs all tests specific to core-geth.

# The following commands acquire external EWASM and EVM interpreter shared objects for
# testing EVMC support.
hera:
	./build/hera.sh

evmone:
	./build/evmone.sh

# Test EVMC support against various external interpreters.
test-evmc: hera evmone
	go test -count 1 ./tests -run TestState -evmc.ewasm=$(ROOT_DIR)/build/_workspace/hera/build/src/libhera.so
	go test -count 1 ./tests -run TestState -evmc.evm=$(ROOT_DIR)/build/_workspace/evmone/lib/libevmone.so

clean-evmc:
	rm -rf ./build/_workspace/hera ./build/_workspace/evmone

test-coregeth-features: \
	test-coregeth-features-coregeth ## Runs tests specific to multi-geth using Fork/Feature configs.

test-coregeth-consensus: test-coregeth-features-clique-consensus

test-coregeth-features-coregeth:
	@echo "Testing fork/feature/datatype implementation; equivalence - COREGETH."
	env COREGETH_TESTS_CHAINCONFIG_FEATURE_EQUIVALENCE_COREGETH=on go test -count=1 -timeout 60m ./tests

test-coregeth-features-clique-consensus:
	@echo "Testing fork/feature/datatype implementation; equivalence - Clique consensus"
	env COREGETH_TESTS_CHAINCONFIG_CONSENSUS_EQUIVALENCE_CLIQUE=on go test -count=1 -timeout 60m -run TestState ./tests ## Only run state tests here, since Blockchain tests will care about rewards, etc.

test-coregeth-chainspecs-coregeth: ## Run tests specific to core-geth using coregeth chainspec file configs.
	@echo "Testing CoreGeth JSON chainspec equivalence."
	env COREGETH_TESTS_CHAINCONFIG_COREGETH_SPECS=on go test -count=1 ./tests

test-coregeth-regression-condensed: geth
	@echo "Running condensed regression tests (imports) against simulated canonical blockchains."
	./tests/regression/simulated/test.sh ./tests/regression/simulated/classic-condense-state/classic.conf.json ./tests/regression/simulated/classic-condense-state/export.rlp.gz
	./tests/regression/simulated/test.sh ./tests/regression/simulated/foundation-condense-state/foundation.conf.json ./tests/regression/simulated/foundation-condense-state/export.rlp.gz
	./tests/regression/simulated/test.sh ./tests/regression/simulated/foundation-condense-state-2/foundation.conf.json ./tests/regression/simulated/foundation-condense-state-2/export.rlp.gz

tests-generate: tests-generate-state tests-generate-difficulty ## Generate all tests.

tests-generate-state: ## Generate state tests.
	@echo "Generating state tests."
	env COREGETH_TESTS_GENERATE_STATE_TESTS=on \
	env COREGETH_TESTS_CHAINCONFIG_FEATURE_EQUIVALENCE_COREGETH=on \
	go test -p 1 -v -timeout 60m ./tests -run TestGenStateAll
	rm -rf ./tests/testdata-etc/GeneralStateTests
	mv ./tests/testdata_generated/GeneralStateTests ./tests/testdata-etc/GeneralStateTests
	rm -rf ./tests/testdata-etc/LegacyTests
	mv ./tests/testdata_generated/LegacyTests ./tests/testdata-etc/LegacyTests
	rm -rf ./tests/testdata_generated

tests-generate-difficulty: ## Generate difficulty tests.
	@echo "Generating difficulty tests."
	env COREGETH_TESTS_GENERATE_DIFFICULTY_TESTS=on \
	go run build/ci.go test -v -timeout 10m ./tests -run TestDifficultyGen
	rm -rf ./tests/testdata-etc/DifficultyTests
	mv ./tests/testdata_generated/DifficultyTests ./tests/testdata-etc/DifficultyTests
	rm -rf ./tests/testdata_generated

lint: ## Run linters.
	$(GORUN) build/ci.go lint

clean: clean-evmc
	go clean -cache
	rm -fr build/_workspace/pkg/ $(GOBIN)/*

mkdocs-serve: ## Serve generated documentation (during development)
	@build/mkdocs-serve.sh

docs-generate: ## Generate JSON RPC API documentation from the OpenRPC service discovery document.
	env COREGETH_GEN_OPENRPC_DOCS=on go test -count=1 -run BuildStatic ./ethclient

# The devtools target installs tools required for 'go generate'.
# You need to put $GOBIN (or $GOPATH/bin) in your PATH to use 'go generate'.

devtools:
	env GOBIN= go install golang.org/x/tools/cmd/stringer@latest
	env GOBIN= go install github.com/fjl/gencodec@latest
	env GOBIN= go install github.com/golang/protobuf/protoc-gen-go@latest
	env GOBIN= go install ./cmd/abigen
	@type "solc" 2> /dev/null || echo 'Please install solc'
	@type "protoc" 2> /dev/null || echo 'Please install protoc'
