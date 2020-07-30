package lib

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
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
	return &MemFreezerRemoteServerAPI{store: make(map[string][]byte)}
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
	fmt.Println("mock server called", "method=HasAncient")
	f.mu.Lock()
	defer f.mu.Unlock()
	_, ok := f.store[f.storeKey(kind, number)]
	return ok, nil
}

func (f *MemFreezerRemoteServerAPI) Ancient(kind string, number uint64) ([]byte, error) {
	fmt.Println("mock server called", "method=Ancient")
	f.mu.Lock()
	defer f.mu.Unlock()
	v, ok := f.store[f.storeKey(kind, number)]
	if !ok {
		return nil, errOutOfBounds
	}
	return v, nil
}

func (f *MemFreezerRemoteServerAPI) Ancients() (uint64, error) {
	fmt.Println("mock server called", "method=Ancients")
	return f.count, nil
}

func (f *MemFreezerRemoteServerAPI) AncientSize(kind string) (uint64, error) {
	fmt.Println("mock server called", "method=AncientSize")
	sum := uint64(0)
	for k, v := range f.store {
		if strings.HasPrefix(k, kind) {
			sum += uint64(len(v))
		}
	}
	return sum, nil
}

func (f *MemFreezerRemoteServerAPI) AppendAncient(number uint64, hash, header, body, receipt, td []byte) error {
	fmt.Println("mock server called", "method=AppendAncient", "number=", number, "header", fmt.Sprintf("%x", header))
	fieldNames := []string{
		freezerRemoteHashTable,
		freezerRemoteHeaderTable,
		freezerRemoteBodiesTable,
		freezerRemoteReceiptTable,
		freezerRemoteDifficultyTable,
	}
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

func (f *MemFreezerRemoteServerAPI) TruncateAncients(n uint64) error {
	fmt.Println("mock server called", "method=TruncateAncients")
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
	fmt.Println("mock server called", "method=Sync")
	return nil
}

func (f *MemFreezerRemoteServerAPI) Close() error {
	fmt.Println("mock server called", "method=Close")
	return nil
}
