package rawdb

import (
	"sync"

	"github.com/ethereum/go-ethereum/rpc"
)

// FreezerRemoteClient is an RPC client implementing the interface of ethdb.AncientStore.
// The struct's methods delegate the business logic to an external server
// that is responsible for managing an actual ancient store.
type FreezerRemoteClient struct {
	client *rpc.Client
	quit   chan struct{}
	mu     sync.Mutex
}

const (
	FreezerMethodClose            = "freezer_close"
	FreezerMethodHasAncient       = "freezer_hasAncient"
	FreezerMethodAncient          = "freezer_ancient"
	FreezerMethodAncients         = "freezer_ancients"
	FreezerMethodAncientSize      = "freezer_ancientSize"
	FreezerMethodAppendAncient    = "freezer_appendAncient"
	FreezerMethodTruncateAncients = "freezer_truncateAncients"
	FreezerMethodSync             = "freezer_sync"
)

// newFreezerRemoteClient constructs a rpc client to connect to a remote freezer
func newFreezerRemoteClient(endpoint string) (*FreezerRemoteClient, error) {
	client, err := rpc.Dial(endpoint)
	if err != nil {
		return nil, err
	}
	return &FreezerRemoteClient{
		client: client,
	}, nil
}

// Close terminates the chain freezer, unmapping all the data files.
func (api *FreezerRemoteClient) Close() error {
	api.mu.Lock()
	defer api.mu.Unlock()
	return api.client.Call(nil, FreezerMethodClose)
}

// HasAncient returns an indicator whether the specified ancient data exists
// in the freezer.
func (api *FreezerRemoteClient) HasAncient(kind string, number uint64) (bool, error) {
	api.mu.Lock()
	defer api.mu.Unlock()
	var res bool
	err := api.client.Call(&res, FreezerMethodHasAncient, kind, number)
	return res, err
}

// Ancient retrieves an ancient binary blob from the append-only immutable files.
func (api *FreezerRemoteClient) Ancient(kind string, number uint64) ([]byte, error) {
	api.mu.Lock()
	defer api.mu.Unlock()
	res := []byte{}
	if err := api.client.Call(&res, FreezerMethodAncient, kind, number); err != nil {
		return nil, err
	}
	return res, nil
}

// Ancients returns the length of the frozen items.
func (api *FreezerRemoteClient) Ancients() (uint64, error) {
	api.mu.Lock()
	defer api.mu.Unlock()
	var res uint64
	err := api.client.Call(&res, FreezerMethodAncients)
	return res, err
}

// AncientSize returns the ancient size of the specified category.
func (api *FreezerRemoteClient) AncientSize(kind string) (uint64, error) {
	api.mu.Lock()
	defer api.mu.Unlock()
	var res uint64
	err := api.client.Call(&res, FreezerMethodAncientSize, kind)
	return res, err
}

// AppendAncient injects all binary blobs belong to block at the end of the
// append-only immutable table files.
//
// Notably, this function is lock free but kind of thread-safe. All out-of-order
// injection will be rejected. But if two injections with same number happen at
// the same time, we can get into the trouble.
//
// Note that the frozen marker is updated outside of the service calls.
func (api *FreezerRemoteClient) AppendAncient(number uint64, hash, header, body, receipts, td []byte) (err error) {
	api.mu.Lock()
	defer api.mu.Unlock()
	return api.client.Call(nil, FreezerMethodAppendAncient, number, hash, header, body, receipts, td)
}

// TruncateAncients discards any recent data above the provided threshold number.
func (api *FreezerRemoteClient) TruncateAncients(items uint64) error {
	api.mu.Lock()
	defer api.mu.Unlock()
	return api.client.Call(nil, FreezerMethodTruncateAncients, items)
}

// Sync flushes all data tables to disk.
func (api *FreezerRemoteClient) Sync() error {
	api.mu.Lock()
	defer api.mu.Unlock()
	return api.client.Call(nil, FreezerMethodSync)
}
