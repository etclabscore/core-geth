package rawdb

import (
	"fmt"
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
type FreezerRemoteServerAPI struct {}

func (f *FreezerRemoteServerAPI) HasAncient(kind string, number uint64) (bool, error) {
	fmt.Println("mock server called", "method=HasAncient")
	return true, nil
}

func (f *FreezerRemoteServerAPI) Ancient(kind string, number uint64) ([]byte, error) {
	fmt.Println("mock server called", "method=Ancient")
	return nil, nil
}

func (f *FreezerRemoteServerAPI) Ancients() (uint64, error) {
	fmt.Println("mock server called", "method=Ancients")
	return 0, nil
}

func (f *FreezerRemoteServerAPI) AncientSize(kind string) (uint64, error) {
	fmt.Println("mock server called", "method=AncientSize")
	return 0, nil
}

func (f *FreezerRemoteServerAPI) AppendAncient(number uint64, hash, header, body, receipt, td []byte) error {
	fmt.Println("mock server called", "method=AppendAncient")
	return nil
}

func (f *FreezerRemoteServerAPI) TruncateAncients(n uint64) error {
	fmt.Println("mock server called", "method=TruncateAncients")
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
	mockFreezerServer := new(FreezerRemoteServerAPI)
	err := server.RegisterName("freezer",mockFreezerServer)
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
		quit: make(chan struct{}),
	}

	n, err := frClient.Ancients()
	if err != nil {
		t.Error(err)
	}
	t.Log(n)
}
