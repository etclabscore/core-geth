// Copyright 2019 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package rawdb

import (
	"log"
	"sync"

	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/metrics"
)

var (
// errUnknownTable is returned if the user attempts to read from a table that is
// not tracked by the freezer.
//errUnknownTable = errors.New("unknown table")

// errOutOrderInsertion is returned if the user attempts to inject out-of-order
// binary blobs into the freezer.
//errOutOrderInsertion = errors.New("the append operation is out-order")

// errSymlinkDatadir is returned if the ancient directory specified by user
// is a symbolic link.
//	errSymlinkDatadir = errors.New("symbolic link datadir is not supported")
)

const (
// freezerRecheckInterval is the frequency to check the key-value database for
// chain progression that might permit new blocks to be frozen into immutable
// storage.
//	freezerRecheckInterval = time.Minute

// freezerBatchLimit is the maximum number of blocks to freeze in one batch
// before doing an fsync and deleting it from the key-value store.
//freezerBatchLimit = 30000
)

// freezer is an memory mapped append-only database to store immutable chain data
// into flat files:
//
// - The append only nature ensures that disk writes are minimized.
// - The memory mapping ensures we can max out system memory for caching without
//   reserving it for go-ethereum. This would also reduce the memory requirements
//   of Geth, and thus also GC overhead.
type freezerRemote struct {
	// WARNING: The `frozen` field is accessed atomically. On 32 bit platforms, only
	// 64-bit aligned fields can be atomic. The struct is guaranteed to be so aligned,
	// so take advantage of that (https://golang.org/pkg/sync/atomic/#pkg-note-BUG).
	// frozen uint64 // Number of blocks already frozen, cached from index marker on remote

	/*
		tables       map[string]*freezerTable // Data tables for storing everything
		instanceLock fileutil.Releaser        // File-system lock to prevent double opens
	*/
	service ethdb.AncientStore
	mu sync.Mutex

	quit chan struct{}
}

// newFreezer creates a chain freezer that moves ancient chain data into
// append-only flat file containers.
func newFreezerRemote(freezerStr string, namespace string) (*freezerRemote, error) {

	// Create the initial freezer object
	// TODO
	var (
		err        error
		readMeter  = metrics.NewRegisteredMeter("ancient.remote."+freezerStr+"/read", nil)
		writeMeter = metrics.NewRegisteredMeter("ancient.remote."+freezerStr+"/write", nil)
		sizeGauge  = metrics.NewRegisteredGauge("ancient.remote."+freezerStr+"/size", nil)
	)

	// Open all the supported data tables
	freezer := &freezerRemote{
		quit: make(chan struct{}),
	}

	switch freezerStr {
	case "s3":
		freezer.service, err = newFreezerRemoteS3(namespace, readMeter, writeMeter, sizeGauge)
		if err != nil {
			return nil, err
		}

		_, err := freezer.service.Ancients()
		if err != nil {
			return freezer, err
		}

	default:
		log.Fatalf("unsupported remote service provider: %s", freezerStr)
	}

	return freezer, nil
}

// Close terminates the chain freezer, unmapping all the data files.
func (f *freezerRemote) Close() error {
	return f.service.Close()
}

// HasAncient returns an indicator whether the specified ancient data exists
// in the freezer.
func (f *freezerRemote) HasAncient(kind string, number uint64) (bool, error) {
	return f.service.HasAncient(kind, number)
}

// Ancient retrieves an ancient binary blob from the append-only immutable files.
func (f *freezerRemote) Ancient(kind string, number uint64) ([]byte, error) {
	return f.service.Ancient(kind, number)
}

// Ancients returns the length of the frozen items.
func (f *freezerRemote) Ancients() (uint64, error) {
	return f.service.Ancients()
}

// AncientSize returns the ancient size of the specified category.
func (f *freezerRemote) AncientSize(kind string) (uint64, error) {
	return f.service.AncientSize(kind)
}

// AppendAncient injects all binary blobs belong to block at the end of the
// append-only immutable table files.
//
// Notably, this function is lock free but kind of thread-safe. All out-of-order
// injection will be rejected. But if two injections with same number happen at
// the same time, we can get into the trouble.
//
// Note that the frozen marker is updated outside of the service calls.
func (f *freezerRemote) AppendAncient(number uint64, hash, header, body, receipts, td []byte) (err error) {
	return f.service.AppendAncient(number, hash, header, body, receipts, td)
}

// Truncate discards any recent data above the provided threshold number.
func (f *freezerRemote) TruncateAncients(items uint64) error {
	return f.service.TruncateAncients(items)
}

// sync flushes all data tables to disk.
func (f *freezerRemote) Sync() error {
	return f.service.Sync()
}

// repair truncates all data tables to the same length.
func (f *freezerRemote) repair() error {
	/*min := uint64(math.MaxUint64)
	for _, table := range f.tables {
		items := atomic.LoadUint64(&table.items)
		if min > items {
			min = items
		}
	}
	for _, table := range f.tables {
		if err := table.truncate(min); err != nil {
			return err
		}
	}
	atomic.StoreUint64(&f.frozen, min)
	*/
	return nil
}
