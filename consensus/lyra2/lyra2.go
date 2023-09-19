package lyra2

/*
#cgo CFLAGS: -std=gnu99
#include "Lyra2.h"
#include <stdlib.h>
*/
import "C"
import (
	"encoding/binary"
	"math/big"
	"math/rand"
	"sync"
	"time"
	"unsafe"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/metrics"
	"github.com/ethereum/go-ethereum/rpc"
)

type Lyra2 struct {
	fakeMode  bool
	fakeFail  uint64
	fakeDelay time.Duration

	log  log.Logger
	lock sync.Mutex

	rand     *rand.Rand
	hashrate metrics.Meter
	update   chan struct{}
	threads  int
	remote   *remoteSealer
}

type Config struct {
	FakeMode  bool
	FakeFail  uint64
	FakeDelay time.Duration

	Log  log.Logger
	Rand *rand.Rand
}

func New(config *Config, notify []string, noverify bool) *Lyra2 {
	if config == nil {
		config = &Config{}
	}
	lyra2 := &Lyra2{
		fakeMode:  config.FakeMode,
		fakeFail:  config.FakeFail,
		fakeDelay: config.FakeDelay,
		log:       log.Root(),
		hashrate:  metrics.NewMeter(),
		update:    make(chan struct{}),
	}
	if config.Log != nil {
		lyra2.log = config.Log
	}
	if config.Rand != nil {
		lyra2.rand = config.Rand
	}
	lyra2.remote = startRemoteSealer(lyra2, notify, noverify)
	return lyra2
}

func NewTester(notify []string, noverify bool) *Lyra2 {
	lyra2 := &Lyra2{
		fakeMode:  false,
		fakeFail:  0,
		fakeDelay: 0,
		log:       log.Root(),
		hashrate:  metrics.NewMeter(),
		update:    make(chan struct{}),
	}
	lyra2.remote = startRemoteSealer(lyra2, notify, noverify)
	return lyra2
}

func (lyra2 *Lyra2) calcHash(headerBytes []byte, nonce uint64, tcost int) *big.Int {
	var lyra2_ctx unsafe.Pointer = C.LYRA2_create()
	result := lyra2.compute(lyra2_ctx, headerBytes, nonce, tcost)
	C.LYRA2_destroy(lyra2_ctx)

	return result.Big()
}

func bytesToHash(in unsafe.Pointer) common.Hash {
	return *(*common.Hash)(in)
}

func (lyra2 *Lyra2) compute(ctx unsafe.Pointer, blockBytes []byte, nonce uint64, tcost int) common.Hash {
	binary.BigEndian.PutUint64(blockBytes[len(blockBytes)-8:], nonce)

	var in unsafe.Pointer = C.CBytes(blockBytes)
	var out unsafe.Pointer = C.malloc(common.HashLength)

	C.LYRA2(ctx, out, common.HashLength, in, C.int32_t(len(blockBytes)), C.int32_t(tcost))

	hash := bytesToHash(out)

	C.free(in)
	C.free(out)

	return hash
}

func (lyra2 *Lyra2) Close() error {
	return nil
}

// APIs implements consensus.Engine, returning the user facing RPC APIs.
func (lyra2 *Lyra2) APIs(chain consensus.ChainHeaderReader) []rpc.API {
	return []rpc.API{
		{
			Namespace: "eth",
			Version:   "1.0",
			Service:   &API{lyra2},
			Public:    true,
		},
	}
}

func (lyra2 *Lyra2) Hashrate() float64 {
	return lyra2.hashrate.Snapshot().Rate1()
}

// Threads returns the number of mining threads currently enabled. This doesn't
// necessarily mean that mining is running!
func (lyra2 *Lyra2) Threads() int {
	lyra2.lock.Lock()
	defer lyra2.lock.Unlock()

	return lyra2.threads
}

// SetThreads updates the number of mining threads currently enabled. Calling
// this method does not start mining, only sets the thread count. If zero is
// specified, the miner will use all cores of the machine. Setting a thread
// count below zero is allowed and will cause the miner to idle, without any
// work being done.
func (lyra2 *Lyra2) SetThreads(threads int) {
	lyra2.lock.Lock()
	defer lyra2.lock.Unlock()

	/*// If we're running a shared PoW, set the thread count on that instead
	  if ethash.shared != nil {
	      ethash.shared.SetThreads(threads)
	      return
	  }*/
	// Update the threads and ping any running seal to pull in any changes
	lyra2.threads = threads
	select {
	case lyra2.update <- struct{}{}:
	default:
	}
}
