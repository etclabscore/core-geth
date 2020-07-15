package rawdb

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/log"
)

// FreezerRemoteAPI exposes a JSONRPC related API
type FreezerRemoteAPI struct {
	freezer *FreezerRemote
}

// NewFreezerRemoteAPI exposes an endpoint to create a remote service
func NewFreezerRemoteAPI(service ethdb.AncientStore) (*FreezerRemoteAPI, error) {
	log.Info("constructing new freezer")
	f, err := newFreezerRemoteService(service)
	if err != nil {
		return nil, err
	}

	freezerAPI := FreezerRemoteAPI{
		freezer: f,
	}
	return &freezerAPI, nil
}

// Close terminates the chain freezer.
func (freezerRemoteAPI *FreezerRemoteAPI) Close() error {
	return freezerRemoteAPI.freezer.Close()
}

// HasAncient returns an indicator whether the specified ancient data exists
// in the freezer.
func (freezerRemoteAPI *FreezerRemoteAPI) HasAncient(kind string, number uint64) (bool, error) {
	return freezerRemoteAPI.freezer.HasAncient(kind, number)
}

// Ancient retrieves an ancient block value as a string.
func (freezerRemoteAPI *FreezerRemoteAPI) Ancient(kind string, number uint64) (string, error) {
	ancient, err := freezerRemoteAPI.freezer.Ancient(kind, number)
	if err != nil {
		return "0x", err
	}
	return hexutil.Encode(ancient), err
}

// Ancients returns the length of the frozen items.
func (freezerRemoteAPI *FreezerRemoteAPI) Ancients() (uint64, error) {
	numAncients, err := freezerRemoteAPI.freezer.Ancients()
	if err != nil {
		return 0, err
	}
	return numAncients, err
}

// AncientSize returns the ancient size of the specified category.
func (freezerRemoteAPI *FreezerRemoteAPI) AncientSize(kind string) (uint64, error) {
	size, err := freezerRemoteAPI.freezer.AncientSize(kind)
	if err != nil {
		return 0, err
	}
	return size, err
}

func (freezerRemoteAPI *FreezerRemoteAPI) AppendAncient(number uint64, hash, header, body, receipts, td string) (err error) {
	var bHash, bHeader, bBody, bReceipts, bTd []byte
	bHash, err = hexutil.Decode(hash)
	if err != nil {
		return err
	}
	bHeader, err = hexutil.Decode(header)
	if err != nil {
		return err
	}
	bBody, err = hexutil.Decode(body)
	if err != nil {
		return err
	}
	bReceipts, err = hexutil.Decode(receipts)
	if err != nil {
		return err
	}
	bTd, err = hexutil.Decode(td)
	return freezerRemoteAPI.freezer.AppendAncient(number, bHash, bHeader, bBody, bReceipts, bTd)
}

// Truncate discards any recent data above the provided threshold number.
func (freezerRemoteAPI *FreezerRemoteAPI) TruncateAncients(items uint64) error {
	return freezerRemoteAPI.freezer.TruncateAncients(items)
}

// sync flushes data, either writing mem to disk, or sending remote requests.
func (freezerRemoteAPI *FreezerRemoteAPI) Sync() error {
	return freezerRemoteAPI.freezer.Sync()
}

// repair truncates all data tables to the same length.
func (freezerRemoteAPI *FreezerRemoteAPI) repair() error {
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
