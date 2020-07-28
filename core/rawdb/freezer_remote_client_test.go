package rawdb

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/ethereum/go-ethereum/rpc"
)

// FreezerRemoteServerAPI is essentially a mock struct on which methods
// describing what the client expects the server to do.
// These methods are 1:1 with ethdb.AncientStore, with the only modification(s)
// being that all []byte and uint64 types are replaced with their hexutil counterparts.
// This means that the client expects the remote ancient store server API to use
// hex encoding in these cases.
// Further, and what cannot be desribed with a skeleton mock like this,
// is that the client should expect the server API to hex encode EVERYTHING,
// including the string values. Package hexutil does not have a corresponding type for this,
// so it has to be done adhoc.
type FreezerRemoteServerAPI struct {
	store map[string][]byte
	count uint64
	mu    sync.Mutex
}

func NewFreezerRemoteServerAPI() *FreezerRemoteServerAPI {
	return &FreezerRemoteServerAPI{store: make(map[string][]byte)}
}

func (r *FreezerRemoteServerAPI) storeKey(kind string, number uint64) string {
	return fmt.Sprintf("%s-%d", kind, number)
}

func (f *FreezerRemoteServerAPI) HasAncient(kind string, number uint64) (bool, error) {
	fmt.Println("mock server called", "method=HasAncient")
	f.mu.Lock()
	defer f.mu.Unlock()
	_, ok := f.store[f.storeKey(kind, number)]
	return ok, nil
}

func (f *FreezerRemoteServerAPI) Ancient(kind string, number uint64) ([]byte, error) {
	fmt.Println("mock server called", "method=Ancient")
	f.mu.Lock()
	defer f.mu.Unlock()
	v, ok := f.store[f.storeKey(kind, number)]
	if !ok {
		return nil, errOutOfBounds
	}
	return v, nil
}

func (f *FreezerRemoteServerAPI) Ancients() (uint64, error) {
	fmt.Println("mock server called", "method=Ancients")
	return f.count, nil
}

func (f *FreezerRemoteServerAPI) AncientSize(kind string) (uint64, error) {
	fmt.Println("mock server called", "method=AncientSize")
	sum := uint64(0)
	for k, v := range f.store {
		if strings.HasPrefix(k, kind) {
			sum += uint64(len(v))
		}
	}
	return sum, nil
}

func (f *FreezerRemoteServerAPI) AppendAncient(number uint64, hash, header, body, receipt, td []byte) error {
	fmt.Println("mock server called", "method=AppendAncient", "number=", number, "header", fmt.Sprintf("%x", header))
	fieldNames := []string{FreezerRemoteHashTable, FreezerRemoteHeaderTable,
		FreezerRemoteBodiesTable, FreezerRemoteReceiptTable, FreezerRemoteDifficultyTable}
	fields := [][]byte{hash, header, body, receipt, td}
	if number != f.count {
		return errOutOrderInsertion
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

func (f *FreezerRemoteServerAPI) TruncateAncients(n uint64) error {
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

func (f *FreezerRemoteServerAPI) Sync() error {
	fmt.Println("mock server called", "method=Sync")
	return nil
}

func (f *FreezerRemoteServerAPI) Close() error {
	fmt.Println("mock server called", "method=Close")
	return nil
}

func newTestServer(t *testing.T) *rpc.Server {
	server := rpc.NewServer()
	mockFreezerServer := NewFreezerRemoteServerAPI()
	err := server.RegisterName("freezer", mockFreezerServer)
	if err != nil {
		t.Fatal(err)
	}
	return server
}

func TestClient1(t *testing.T) {
	server := newTestServer(t)
	client := rpc.DialInProc(server)

	frClient := &FreezerRemoteClient{
		client: client,
		quit:   make(chan struct{}),
	}

	ancientTestProgram := func(head *uint64, i int) {

		kinds := []string{FreezerRemoteHashTable, FreezerRemoteHeaderTable,
			FreezerRemoteBodiesTable, FreezerRemoteReceiptTable, FreezerRemoteDifficultyTable}
		testKind := kinds[i%len(kinds)]

		err := frClient.AppendAncient(*head, []byte{uint8(*head)}, []byte{uint8(*head)}, []byte{uint8(*head)}, []byte{uint8(*head)}, []byte{uint8(*head)})
		if err != nil {
			t.Fatalf("append: %v", err)
		}
		*head++

		ok, err := frClient.HasAncient(testKind, *head-1)
		if err != nil || !ok {
			t.Fatalf("has ancient: %v %v", err, ok)
		}

		serverHead, err := frClient.Ancients()
		if err != nil {
			t.Fatalf("ancients: %v", err)
		}
		if serverHead != *head {
			t.Fatalf("mismatch server/local *head: local=%d server=%d", *head, serverHead)
		}

		v, err := frClient.Ancient(testKind, *head-1)
		if err != nil {
			t.Fatalf("ancient: %v", err)
		}
		if bytes.Compare([]byte{uint8(*head - 1)}, v) != 0 {
			t.Fatalf("mismatch store value: want: %x, got: %x", []byte{uint8(*head)}, v)
		}

		got, err := frClient.AncientSize(testKind)
		if err != nil {
			t.Fatalf("ancient size: %v", err)
		}
		if got < 1 {
			t.Fatalf("low ancient size: %d", got)
		}

		if i > 23 && i%23 == 0 {
			err := frClient.Sync()
			if err != nil {
				t.Fatalf("sync: %v", err)
			}
		}
		if i > 42 && i%42 == 0 {
			*head -= 15
			err := frClient.TruncateAncients(*head)
			if err != nil {
				t.Fatalf("truncate: %v", err)
			}
		}
		if i > 88 && i%88 == 0 {

		}
	}

	head := uint64(0)
	for i := 0; i < 1000; i++ {
		ancientTestProgram(&head, i)
	}
}
