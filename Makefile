# This Makefile is meant to be used by people that do not usually work
# with Go source code. If you know what GOPATH is then you probably
# don't need to bother with make.

.PHONY: geth android ios geth-cross evm evmc mkdocs-serve all test clean

GOBIN = ./build/bin
GO ?= latest
GORUN = env GO111MODULE=on go run
ROOT_DIR := $(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

geth:
	$(GORUN) build/ci.go install ./cmd/geth
	@echo "Done building."
	@echo "Run \"$(GOBIN)/geth\" to launch geth."

all:
	$(GORUN) build/ci.go install

android:
	$(GORUN) build/ci.go aar --local
	@echo "Done building."
	@echo "Import \"$(GOBIN)/geth.aar\" to use the library."
	@echo "Import \"$(GOBIN)/geth-sources.jar\" to add javadocs"
	@echo "For more info see https://stackoverflow.com/questions/20994336/android-studio-how-to-attach-javadoc"

ios:
	$(GORUN) build/ci.go xcode --local
	@echo "Done building."
	@echo "Import \"$(GOBIN)/Geth.framework\" to use the library."

test:
	$(GORUN) build/ci.go test

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

ssvm:
	./build/ssvm.sh

evmone:
	./build/evmone.sh

aleth-interpreter:
	./build/aleth-interpreter.sh

# Test EVMC support against various external interpreters.
test-evmc: hera ssvm evmone aleth-interpreter
	go test -count 1 ./tests -run TestState -evmc.ewasm=$(ROOT_DIR)/build/_workspace/hera/build/src/libhera.so
	go test -count 1 ./tests -run TestState -evmc.ewasm=$(ROOT_DIR)/build/_workspace/SSVM/build/tools/ssvm-evmc/libssvmEVMC.so
	go test -count 1 ./tests -run TestState -evmc.evm=$(ROOT_DIR)/build/_workspace/evmone/lib/libevmone.so
	go test -count 1 ./tests -run TestState -evmc.evm=$(ROOT_DIR)/build/_workspace/aleth/lib/libaleth-interpreter.so

clean-evmc:
	rm -rf ./build/_workspace/hera ./build/_workspace/SSVM ./build/_workspace/evmone ./build/_workspace/aleth

test-coregeth-features: \
	test-coregeth-features-coregeth \
	test-coregeth-features-multigethv0 ## Runs tests specific to multi-geth using Fork/Feature configs.

test-coregeth-consensus: test-coregeth-features-clique-consensus

test-coregeth-features-coregeth:
	@echo "Testing fork/feature/datatype implementation; equivalence - COREGETH."
	env COREGETH_TESTS_CHAINCONFIG_FEATURE_EQUIVALENCE_COREGETH=on go test -count=1 -timeout 60m ./tests

test-coregeth-features-multigethv0:
	@echo "Testing fork/feature/datatype implementation; equivalence - MULTIGETHv0."
	env COREGETH_TESTS_CHAINCONFIG_FEATURE_EQUIVALENCE_MULTIGETHV0=on go test -count=1 -timeout 60m ./tests

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

tests-generate-difficulty: ## Generate difficulty tests.
	@echo "Generating difficulty tests configs."
	env COREGETH_TESTS_GENERATE_DIFFICULTY_TESTS_CONFIGS=on \
	go run build/ci.go test -v -timeout 10m ./tests -run TestDifficultyTestConfigGen

	@echo "Generating difficulty tests."
	env COREGETH_TESTS_GENERATE_DIFFICULTY_TESTS=on \
	go run build/ci.go test -v -timeout 10m ./tests -run TestDifficultyGen

lint: ## Run linters.
	$(GORUN) build/ci.go lint

clean: clean-evmc
	env GO111MODULE=on go clean -cache
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
