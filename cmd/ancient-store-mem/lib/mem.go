// Copyright 2020 The core-geth Authors
// This file is part of the core-geth library.
//
// The core-geth library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The core-geth library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the core-geth library. If not, see <http://www.gnu.org/licenses/>.

package lib

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/common"
)

const (
	freezerRemoteHashTable       = "hashes"
	freezerRemoteHeaderTable     = "headers"
	freezerRemoteBodiesTable     = "bodies"
	freezerRemoteReceiptTable    = "receipts"
	freezerRemoteDifficultyTable = "diffs"
)

var (
	errOutOfBounds = errors.New("out of bounds")
	errOutOfOrder  = errors.New("out of order")
)

// MemFreezerRemoteServerAPI is a mock freezer server implementation.
type MemFreezerRemoteServerAPI struct {
	store map[string][]byte
	count uint64
	mu    sync.Mutex
}

func NewMemFreezerRemoteServerAPI() *MemFreezerRemoteServerAPI {
	return &MemFreezerRemoteServerAPI{
		store: make(map[string][]byte),
	}
}

func (r *MemFreezerRemoteServerAPI) storeKey(kind string, number uint64) string {
	return fmt.Sprintf("%s-%d", kind, number)
}

func (f *MemFreezerRemoteServerAPI) Reset() {
	f.count = 0
	f.mu.Lock()
	f.store = make(map[string][]byte)
	f.mu.Unlock()
}

func (f *MemFreezerRemoteServerAPI) HasAncient(kind string, number uint64) (bool, error) {
	// fmt.Println("mock server called", "method=HasAncient")
	f.mu.Lock()
	defer f.mu.Unlock()
	_, ok := f.store[f.storeKey(kind, number)]
	return ok, nil
}

func (f *MemFreezerRemoteServerAPI) Ancient(kind string, number uint64) ([]byte, error) {
	// fmt.Println("mock server called", "method=Ancient")
	f.mu.Lock()
	defer f.mu.Unlock()
	v, ok := f.store[f.storeKey(kind, number)]
	if !ok {
		return nil, errOutOfBounds
	}

	return v, nil
}

func (f *MemFreezerRemoteServerAPI) Ancients() (uint64, error) {
	// fmt.Println("mock server called", "method=Ancients")
	return f.count, nil
}

func (f *MemFreezerRemoteServerAPI) AncientRange(kind string, start, count, maxBytes uint64) ([][]byte, error) {
	res := make([][]byte, 0)
	// fmt.Println("mock server called", "method=Ancients")
	for i := uint64(0); i < count; i++ {
		item := f.store[f.storeKey(kind, start+i)]
		if len(item) > int(maxBytes) {
			item = item[:int(maxBytes)]
		}
		res = append(res, item)
	}
	return res, nil
}

func (f *MemFreezerRemoteServerAPI) AncientSize(kind string) (uint64, error) {
	// fmt.Println("mock server called", "method=AncientSize")
	sum := uint64(0)
	for k, v := range f.store {
		if strings.HasPrefix(k, kind) {
			sum += uint64(len(v))
		}
	}
	return sum, nil
}

var fieldNames = []string{
	freezerRemoteHashTable,
	freezerRemoteHeaderTable,
	freezerRemoteBodiesTable,
	freezerRemoteReceiptTable,
	freezerRemoteDifficultyTable,
}

func (f *MemFreezerRemoteServerAPI) AppendAncient(number uint64, hash, header, body, receipt, td []byte) error {
	// fmt.Println("mock server called", "method=AppendAncient", "number=", number, "header", fmt.Sprintf("%x", header))
	fields := [][]byte{hash, header, body, receipt, td}
	if number != f.count {
		return errOutOfOrder
	}
	f.count = number + 1
	f.mu.Lock()
	defer f.mu.Unlock()
	for i, fv := range fields {
		kind := fieldNames[i]
		f.store[f.storeKey(kind, number)] = fv
	}
	return nil
}

func (f *MemFreezerRemoteServerAPI) Append(kind string, num uint64, item interface{}) error {
	if f.count != num {
		return fmt.Errorf("%w: num=%d, count=%d", errOutOfOrder, num, f.count)
	}

	// This is a really crufty thing.
	// We need to increment the freezer counter when all of a block's data have been written.
	// This is harder to do than AppendAncient because we're only handling one field at a time.
	// So we look at the ModifyAncients method and use the implementation for the
	// current AncientWriteOperator closure. In accessors_chain.go, this
	// writes one block at time (via independent fields), and the last independent
	// field called is 'diffs'. (See methods: WriteAncientBlocks and writeAncientBlock).
	// As long as we assume that 'diffs' are the last-called field when writing
	// a block, we can use it as the trigger for the count incrementing.
	if kind == freezerRemoteDifficultyTable {
		f.count = num + 1
	}

	str := item.(string)

	f.mu.Lock()
	f.store[f.storeKey(kind, num)] = common.Hex2Bytes(str)
	f.mu.Unlock()

	return nil
}

func (f *MemFreezerRemoteServerAPI) AppendRaw(kind string, num uint64, item []byte) error {
	if f.count != num {
		return fmt.Errorf("%w: num=%d, count=%d", errOutOfOrder, num, f.count)
	}

	f.mu.Lock()
	f.store[f.storeKey(kind, num)] = item
	f.mu.Unlock()

	return nil
}

func (f *MemFreezerRemoteServerAPI) TruncateTail(n uint64) error {
	// fmt.Println("mock server called", "method=TruncateAncients")
	if f.count <= n {
		return nil
	}
	f.count = n
	f.mu.Lock()
	defer f.mu.Unlock()
	for k := range f.store {
		spl := strings.Split(k, "-")
		num, err := strconv.ParseUint(spl[1], 10, 64)
		if err != nil {
			return err
		}
		if num <= n {
			delete(f.store, k)
		}
	}
	return nil
}

func (f *MemFreezerRemoteServerAPI) TruncateHead(n uint64) error {
	// fmt.Println("mock server called", "method=TruncateAncients")
	if f.count <= n {
		return nil
	}
	f.count = n
	f.mu.Lock()
	defer f.mu.Unlock()
	for k := range f.store {
		spl := strings.Split(k, "-")
		num, err := strconv.ParseUint(spl[1], 10, 64)
		if err != nil {
			return err
		}
		if num >= n {
			delete(f.store, k)
		}
	}
	return nil
}

func (f *MemFreezerRemoteServerAPI) Sync() error {
	// fmt.Println("mock server called", "method=Sync")
	return nil
}

func (f *MemFreezerRemoteServerAPI) Close() error {
	// fmt.Println("mock server called", "method=Close")
	return nil
}
