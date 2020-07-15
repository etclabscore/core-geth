package rawdb

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
)

// ExternalRemoteFreezer is rpc client functionality for a freezer
type ExternalRemoteFreezer struct {
	client   *rpc.Client
	status   string
}

// NewFreezerRemoteClient constructs a rpc client to connect to a remote freezer
func NewFreezerRemoteClient(endpoint string, ipc bool) (*ExternalRemoteFreezer, error) {
	client, err := rpc.Dial(endpoint)
	if err != nil {
		return nil, err
	}

	extfreezer := &ExternalRemoteFreezer{
		client:   client,
	}

	// Check if reachable
	version, err := extfreezer.pingVersion()
	if err != nil {
		return nil, err
	}
	extfreezer.status = fmt.Sprintf("ok [version=%v]", version)
	return extfreezer, nil
}

func (api *ExternalRemoteFreezer) pingVersion() (string, error) {

	return "version 1", nil
}

// Close terminates the chain freezer, unmapping all the data files.
func (api *ExternalRemoteFreezer) Close() error {
	var res string
	err := api.client.Call(&res, "freezer_close")
	return err
}

// HasAncient returns an indicator whether the specified ancient data exists
// in the freezer.
func (api *ExternalRemoteFreezer) HasAncient(kind string, number uint64) (bool, error) {
	var res bool
	err := api.client.Call(&res, "freezer_hasAncient", kind, number)
	return res, err
}

// Ancient retrieves an ancient binary blob from the append-only immutable files.
func (api *ExternalRemoteFreezer) Ancient(kind string, number uint64) ([]byte, error) {
	var res string
	if err := api.client.Call(&res, "freezer_ancient", kind, number); err != nil {
		return nil, err
	}
	return hexutil.Decode(res)
}

// Ancients returns the length of the frozen items.
func (api *ExternalRemoteFreezer) Ancients() (uint64, error) {
	var res uint64
	err := api.client.Call(&res, "freezer_ancients")
	return res, err
}

// AncientSize returns the ancient size of the specified category.
func (api *ExternalRemoteFreezer) AncientSize(kind string) (uint64, error) {
	var res uint64
	err := api.client.Call(&res, "freezer_ancientSize", kind)
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
func (api *ExternalRemoteFreezer) AppendAncient(number uint64, hash, header, body, receipts, td []byte) (err error) {
	var res string
	hexHash := hexutil.Encode(hash)
	hexHeader := hexutil.Encode(header)
	hexBody := hexutil.Encode(body)
	hexReceipts := hexutil.Encode(receipts)
	hexTd := hexutil.Encode(td)
	err = api.client.Call(&res, "freezer_appendAncient", number, hexHash, hexHeader, hexBody, hexReceipts, hexTd)
	return
}

// TruncateAncients discards any recent data above the provided threshold number.
func (api *ExternalRemoteFreezer) TruncateAncients(items uint64) error {
	var res string
	return api.client.Call(&res, "freezer_truncateAncients", items)
}

// Sync flushes all data tables to disk.
func (api *ExternalRemoteFreezer) Sync() error {
	var res string
	return api.client.Call(&res, "freezer_sync")
}

// repair truncates all data tables to the same length.
func (api *ExternalRemoteFreezer) repair() error {
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
