# This Makefile is meant to be used by people that do not usually work
# with Go source code. If you know what GOPATH is then you probably
# don't need to bother with make.

.PHONY: all test clean
.PHONY: evmc
.PHONY: geth android ios geth-cross
.PHONY: geth-linux geth-linux-386 geth-linux-amd64 geth-linux-mips64 geth-linux-mips64le
.PHONY: geth-linux-arm geth-linux-arm-5 geth-linux-arm-6 geth-linux-arm-7 geth-linux-arm64
.PHONY: geth-darwin geth-darwin-386 geth-darwin-amd64
.PHONY: geth-windows geth-windows-386 geth-windows-amd64
.PHONY: mkdocs-serve

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

sync-parity-chainspecs:
	./params/parity.json.d/sync-parity-remote.sh

test-coregeth: \
 test-coregeth-features \
 test-coregeth-chainspecs \
 test-coregeth-consensus \
 test-coregeth-regression-condensed ## Runs all tests specific to core-geth.

# Generate the necessary shared object for EVMC unit tests.
evmc/bindings/go/evmc/example_vm.so:
	./build/evmc-example_vm.so.sh

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
test-evmc: evmc/bindings/go/evmc/example_vm.so hera ssvm evmone aleth-interpreter
	go test -count 1 ./evmc/...
	go test -count 1 ./tests -run TestState -evmc.ewasm=$(ROOT_DIR)/build/_workspace/hera/build/src/libhera.so
	go test -count 1 ./tests -run TestState -evmc.ewasm=$(ROOT_DIR)/build/_workspace/SSVM/build/tools/ssvm-evmc/libssvmEVMC.so
	go test -count 1 ./tests -run TestState -evmc.evm=$(ROOT_DIR)/build/_workspace/evmone/lib/libevmone.so
	go test -count 1 ./tests -run TestState -evmc.evm=$(ROOT_DIR)/build/_workspace/aleth/lib/libaleth-interpreter.so

clean-evmc:
	rm -rf evmc/bindings/go/evmc/example_vm.so ./build/_workspace/hera ./build/_workspace/SSVM ./build/_workspace/evmone ./build/_workspace/aleth

test-coregeth-features: test-coregeth-features-parity test-coregeth-features-coregeth test-coregeth-features-multigethv0 ## Runs tests specific to multi-geth using Fork/Feature configs.

test-coregeth-consensus: test-coregeth-features-clique-consensus

test-coregeth-features-parity:
	@echo "Testing fork/feature/datatype implementation; equivalence - OPENETHEREUM."
	env COREGETH_TESTS_CHAINCONFIG_FEATURE_EQUIVALENCE_OPENETHEREUM=on go test -count=1 ./tests

test-coregeth-features-coregeth:
	@echo "Testing fork/feature/datatype implementation; equivalence - COREGETH."
	env COREGETH_TESTS_CHAINCONFIG_FEATURE_EQUIVALENCE_COREGETH=on go test -count=1 ./tests

test-coregeth-features-multigethv0:
	@echo "Testing fork/feature/datatype implementation; equivalence - MULTIGETHv0."
	env COREGETH_TESTS_CHAINCONFIG_FEATURE_EQUIVALENCE_MULTIGETHV0=on go test -count=1 ./tests

test-coregeth-features-clique-consensus:
	@echo "Testing fork/feature/datatype implementation; equivalence - Clique consensus"
	env COREGETH_TESTS_CHAINCONFIG_CONSENSUS_EQUIVALENCE_CLIQUE=on go test -count=1 -run TestState ./tests ## Only run state tests here, since Blockchain tests will care about rewards, etc.

test-coregeth-chainspecs: ## Run tests specific to multi-geth using chainspec file configs.
	@echo "Testing Parity JSON chainspec equivalence."
	env COREGETH_TESTS_CHAINCONFIG_OPENETHEREUM_SPECS=on go test -count=1 ./tests

test-coregeth-regression-condensed: geth
	@echo "Running condensed regression tests (imports) against simulated canonical blockchains."
	./tests/regression/simulated/test.sh ./tests/regression/simulated/classic-condense-state/classic.conf.json ./tests/regression/simulated/classic-condense-state/export.rlp.gz
	./tests/regression/simulated/test.sh ./tests/regression/simulated/foundation-condense-state/foundation.conf.json ./tests/regression/simulated/foundation-condense-state/export.rlp.gz
	./tests/regression/simulated/test.sh ./tests/regression/simulated/foundation-condense-state-2/foundation.conf.json ./tests/regression/simulated/foundation-condense-state-2/export.rlp.gz

tests-generate: tests-generate-state tests-generate-difficulty ## Generate all tests.

tests-generate-state: ## Generate state tests.
	@echo "Generating state tests."
	env COREGETH_TESTS_CHAINCONFIG_OPENETHEREUM_SPECS=on \
	env COREGETH_TESTS_GENERATE_STATE_TESTS=on \
	go run build/ci.go test -v ./tests -run TestGenState

tests-generate-difficulty: ## Generate difficulty tests.
	@echo "Generating difficulty tests."
	env COREGETH_TESTS_CHAINCONFIG_OPENETHEREUM_SPECS=on \
	env COREGETH_TESTS_GENERATE_DIFFICULTY_TESTS=on \
	go run build/ci.go test -v ./tests -run TestDifficultyGen

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
	env GOBIN= go get -u golang.org/x/tools/cmd/stringer
	env GOBIN= go get -u github.com/kevinburke/go-bindata/go-bindata
	env GOBIN= go get -u github.com/fjl/gencodec
	env GOBIN= go get -u github.com/golang/protobuf/protoc-gen-go
	env GOBIN= go install ./cmd/abigen
	@type "npm" 2> /dev/null || echo 'Please install node.js and npm'
	@type "solc" 2> /dev/null || echo 'Please install solc'
	@type "protoc" 2> /dev/null || echo 'Please install protoc'

# Cross Compilation Targets (xgo)

geth-cross: geth-linux geth-darwin geth-windows geth-android geth-ios
	@echo "Full cross compilation done:"
	@ls -ld $(GOBIN)/geth-*

geth-linux: geth-linux-386 geth-linux-amd64 geth-linux-arm geth-linux-mips64 geth-linux-mips64le
	@echo "Linux cross compilation done:"
	@ls -ld $(GOBIN)/geth-linux-*

geth-linux-386:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/386 -v ./cmd/geth
	@echo "Linux 386 cross compilation done:"
	@ls -ld $(GOBIN)/geth-linux-* | grep 386

geth-linux-amd64:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/amd64 -v ./cmd/geth
	@echo "Linux amd64 cross compilation done:"
	@ls -ld $(GOBIN)/geth-linux-* | grep amd64

geth-linux-arm: geth-linux-arm-5 geth-linux-arm-6 geth-linux-arm-7 geth-linux-arm64
	@echo "Linux ARM cross compilation done:"
	@ls -ld $(GOBIN)/geth-linux-* | grep arm

geth-linux-arm-5:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/arm-5 -v ./cmd/geth
	@echo "Linux ARMv5 cross compilation done:"
	@ls -ld $(GOBIN)/geth-linux-* | grep arm-5

geth-linux-arm-6:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/arm-6 -v ./cmd/geth
	@echo "Linux ARMv6 cross compilation done:"
	@ls -ld $(GOBIN)/geth-linux-* | grep arm-6

geth-linux-arm-7:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/arm-7 -v ./cmd/geth
	@echo "Linux ARMv7 cross compilation done:"
	@ls -ld $(GOBIN)/geth-linux-* | grep arm-7

geth-linux-arm64:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/arm64 -v ./cmd/geth
	@echo "Linux ARM64 cross compilation done:"
	@ls -ld $(GOBIN)/geth-linux-* | grep arm64

geth-linux-mips:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/mips --ldflags '-extldflags "-static"' -v ./cmd/geth
	@echo "Linux MIPS cross compilation done:"
	@ls -ld $(GOBIN)/geth-linux-* | grep mips

geth-linux-mipsle:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/mipsle --ldflags '-extldflags "-static"' -v ./cmd/geth
	@echo "Linux MIPSle cross compilation done:"
	@ls -ld $(GOBIN)/geth-linux-* | grep mipsle

geth-linux-mips64:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/mips64 --ldflags '-extldflags "-static"' -v ./cmd/geth
	@echo "Linux MIPS64 cross compilation done:"
	@ls -ld $(GOBIN)/geth-linux-* | grep mips64

geth-linux-mips64le:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/mips64le --ldflags '-extldflags "-static"' -v ./cmd/geth
	@echo "Linux MIPS64le cross compilation done:"
	@ls -ld $(GOBIN)/geth-linux-* | grep mips64le

geth-darwin: geth-darwin-386 geth-darwin-amd64
	@echo "Darwin cross compilation done:"
	@ls -ld $(GOBIN)/geth-darwin-*

geth-darwin-386:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=darwin/386 -v ./cmd/geth
	@echo "Darwin 386 cross compilation done:"
	@ls -ld $(GOBIN)/geth-darwin-* | grep 386

geth-darwin-amd64:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=darwin/amd64 -v ./cmd/geth
	@echo "Darwin amd64 cross compilation done:"
	@ls -ld $(GOBIN)/geth-darwin-* | grep amd64

geth-windows: geth-windows-386 geth-windows-amd64
	@echo "Windows cross compilation done:"
	@ls -ld $(GOBIN)/geth-windows-*

geth-windows-386:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=windows/386 -v ./cmd/geth
	@echo "Windows 386 cross compilation done:"
	@ls -ld $(GOBIN)/geth-windows-* | grep 386

geth-windows-amd64:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=windows/amd64 -v ./cmd/geth
	@echo "Windows amd64 cross compilation done:"
	@ls -ld $(GOBIN)/geth-windows-* | grep amd64
