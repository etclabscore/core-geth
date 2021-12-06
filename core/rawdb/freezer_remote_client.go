package rawdb

import (
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params/vars"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/rpc"
)

// FreezerRemoteClient is an RPC client implementing the interface of ethdb.AncientStore.
// The struct's methods delegate the business logic to an external server
// that is responsible for managing an actual ancient store.
type FreezerRemoteClient struct {
	client     *rpc.Client
	quit       chan struct{}
	threshold  uint64             // Number of recent blocks not to freeze (params.FullImmutabilityThreshold apart from tests)
	trigger    chan chan struct{} // Manual blocking freeze trigger, test determinism
	closeOnce  sync.Once
	readonly   bool
	wg         sync.WaitGroup
	writeBatch *freezerBatchRemote
	writeMu    sync.Mutex
}

const (
	FreezerMethodClose            = "freezer_close"
	FreezerMethodHasAncient       = "freezer_hasAncient"
	FreezerMethodAncient          = "freezer_ancient"
	FreezerMethodAncients         = "freezer_ancients"
	FreezerMethodAncientSize      = "freezer_ancientSize"
	FreezerMethodAppendAncient    = "freezer_appendAncient"
	FreezerMethodModifyAncients   = "freezer_modifyAncients"
	FreezerMethodTruncateAncients = "freezer_truncateAncients"
	FreezerMethodSync             = "freezer_sync"

	// FreezerMethodWriteAppend and FreezerMethodWriteAppendRaw are
	// methods for re-written (get it?) freezer design with write batching.
	FreezerMethodWriteAppend    = "freezer_append"
	FreezerMethodWriteAppendRaw = "freezer_appendRaw"
)

// newFreezerRemoteClient constructs a rpc client to connect to a remote freezer
func newFreezerRemoteClient(endpoint string, readonly bool) (*FreezerRemoteClient, error) {
	client, err := rpc.Dial(endpoint)
	if err != nil {
		return nil, err
	}
	return &FreezerRemoteClient{
		client:     client,
		threshold:  vars.FullImmutabilityThreshold,
		quit:       make(chan struct{}),
		trigger:    make(chan chan struct{}),
		readonly:   readonly,
		writeBatch: &freezerBatchRemote{client: client},
	}, nil
}

// Close terminates the chain freezer, unmapping all the data files.
func (api *FreezerRemoteClient) Close() error {
	return api.client.Call(nil, FreezerMethodClose)
}

// HasAncient returns an indicator whether the specified ancient data exists
// in the freezer.
func (api *FreezerRemoteClient) HasAncient(kind string, number uint64) (bool, error) {
	var res bool
	err := api.client.Call(&res, FreezerMethodHasAncient, kind, number)
	return res, err
}

// Ancient retrieves an ancient binary blob from the append-only immutable files.
func (api *FreezerRemoteClient) Ancient(kind string, number uint64) ([]byte, error) {
	res := []byte{}
	if err := api.client.Call(&res, FreezerMethodAncient, kind, number); err != nil {
		// llog.Println("frc.Ancient", kind, number, string(res))
		return nil, err
	}
	return res, nil
}

// Ancients returns the length of the frozen items.
func (api *FreezerRemoteClient) Ancients() (uint64, error) {
	var res uint64
	err := api.client.Call(&res, FreezerMethodAncients)
	return res, err
}

// 	// ReadAncients runs the given read operation while ensuring that no writes take place
//	// on the underlying freezer.
func (api *FreezerRemoteClient) ReadAncients(fn func(ethdb.AncientReader) error) (err error) {
	api.writeMu.Lock()
	defer api.writeMu.Unlock()
	return fn(api)
}

// AncientSize returns the ancient size of the specified category.
func (api *FreezerRemoteClient) AncientSize(kind string) (uint64, error) {
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
	api.writeMu.Lock()
	defer api.writeMu.Unlock()
	return api.client.Call(nil, FreezerMethodAppendAncient, number, hash, header, body, receipts, td)
}

type freezerBatchRemote struct {
	client *rpc.Client

	// writeSize is managed here (at the client, rather than the server)
	// because the write-operations are granular (per header, receipt, hash, etcetera), and the size written value
	// is an aggregate at the ModifyAncients-level.
	writeSize int64
}

func (b *freezerBatchRemote) Append(kind string, num uint64, item interface{}) error {
	var res int64
	bs, err := rlp.EncodeToBytes(item)
	if err != nil {
		return err
	}

	err = b.client.Call(&res, FreezerMethodWriteAppend, kind, num, common.Bytes2Hex(bs))
	if err != nil {
		b.writeSize += int64(len(bs))
	}
	return err
}

func (b *freezerBatchRemote) AppendRaw(kind string, num uint64, item []byte) error {
	err := b.client.Call(nil, FreezerMethodWriteAppendRaw, kind, num, item)
	if err != nil {
		b.writeSize += int64(len(item))
	}
	return err
}

func (b *freezerBatchRemote) reset() {
	b.writeSize = 0
}

// ModifyAncients runs the given write operation.
func (api *FreezerRemoteClient) ModifyAncients(fn func(ethdb.AncientWriteOperator) error) (writeSize int64, err error) {
	if api.readonly {
		return 0, errReadOnly
	}
	prev, err := api.Ancients()
	if err != nil {
		return 0, err
	}
	api.writeMu.Lock()
	defer api.writeMu.Unlock()

	// Roll back all tables to the starting position in case of error.
	defer func() {
		if err != nil {
			log.Warn("Rolling back ancients", "target(previous)", prev)
			if err := api.TruncateAncients(prev); err != nil {
				log.Error("Freezer table roll-back failed", "index", prev, "err", err)
			}
		}
	}()

	// Do the writing.
	api.writeBatch.reset()
	if err := fn(api.writeBatch); err != nil {
		return 0, err
	}

	return api.writeBatch.writeSize, nil
}

// TruncateAncients discards any recent data above the provided threshold number.
func (api *FreezerRemoteClient) TruncateAncients(items uint64) error {
	api.writeMu.Lock()
	defer api.writeMu.Unlock()
	return api.client.Call(nil, FreezerMethodTruncateAncients, items)
}

// Sync flushes all data tables to disk.
func (api *FreezerRemoteClient) Sync() error {
	api.writeMu.Lock()
	defer api.writeMu.Unlock()
	return api.client.Call(nil, FreezerMethodSync)
}

// freezeRemote is a background thread that periodically checks the blockchain for any
// import progress and moves ancient data from the fast database into the freezer.
//
// This functionality is deliberately broken off from block importing to avoid
// incurring additional data shuffling delays on block propagation.
//
// This function is a near-duplicate of the logic implemented by *freezer.freeze,
// but is used instead by the remote freezer client rather than the builtin FS ancient
// store. Code is near-duplicated to permit the default FS ancient store logic
// to exist unmodified and untouched by the remote freezer client, which demands
// a slightly different signature, and uses the freezer.Ancients() method instead
// of direct access to the atomic freezer.frozen field.
func freezeRemote(db ethdb.KeyValueStore, f ethdb.AncientStore, threshold uint64, quitChan chan struct{}, triggerChanChan chan chan struct{}) {
	nfdb := &nofreezedb{KeyValueStore: db}

	var (
		backoff   bool
		triggered chan struct{} // Used in tests
	)
	for {
		select {
		case <-quitChan:
			log.Info("Freezer shutting down")
			return
		default:
		}
		if backoff {
			// If we were doing a manual trigger, notify it
			if triggered != nil {
				triggered <- struct{}{}
				triggered = nil
			}
			select {
			case <-time.NewTimer(freezerRecheckInterval).C:
				backoff = false
			case triggered = <-triggerChanChan:
				backoff = false
			case <-quitChan:
				return
			}
		}
		// Retrieve the freezing threshold.
		hash := ReadHeadBlockHash(nfdb)
		if hash == (common.Hash{}) {
			log.Debug("Current full block hash unavailable") // new chain, empty database
			backoff = true
			continue
		}
		numFrozen, err := f.Ancients()
		if err != nil {
			log.Crit("ancient db freeze", "error", err)
		}
		number := ReadHeaderNumber(nfdb, hash)
		// threshold := atomic.LoadUint64(&f.threshold)

		switch {
		case number == nil:
			log.Error("Current full block number unavailable", "hash", hash)
			backoff = true
			continue

		case *number < threshold:
			log.Debug("Current full block not old enough", "number", *number, "hash", hash, "delay", threshold)
			backoff = true
			continue

		case *number-threshold <= numFrozen:
			log.Debug("Ancient blocks frozen already", "number", *number, "hash", hash, "frozen", numFrozen)
			backoff = true
			continue
		}
		head := ReadHeader(nfdb, hash, *number)
		if head == nil {
			log.Error("Current full block unavailable", "number", *number, "hash", hash)
			backoff = true
			continue
		}
		// Seems we have data ready to be frozen, process in usable batches
		limit := *number - threshold
		if limit-numFrozen > freezerBatchLimit {
			limit = numFrozen + freezerBatchLimit
		}
		var (
			start    = time.Now()
			first    = numFrozen
			ancients = make([]common.Hash, 0, limit-numFrozen)
		)
		for numFrozen <= limit {
			// Retrieves all the components of the canonical block
			hash := ReadCanonicalHash(nfdb, numFrozen)
			if hash == (common.Hash{}) {
				log.Error("Canonical hash missing, can't freeze", "number", numFrozen)
				break
			}
			header := ReadHeaderRLP(nfdb, hash, numFrozen)
			if len(header) == 0 {
				log.Error("Block header missing, can't freeze", "number", numFrozen, "hash", hash)
				break
			}
			body := ReadBodyRLP(nfdb, hash, numFrozen)
			if len(body) == 0 {
				log.Error("Block body missing, can't freeze", "number", numFrozen, "hash", hash)
				break
			}
			receipts := ReadReceiptsRLP(nfdb, hash, numFrozen)
			if len(receipts) == 0 {
				log.Error("Block receipts missing, can't freeze", "number", numFrozen, "hash", hash)
				break
			}
			td := ReadTdRLP(nfdb, hash, numFrozen)
			if len(td) == 0 {
				log.Error("Total difficulty missing, can't freeze", "number", numFrozen, "hash", hash)
				break
			}
			log.Trace("Deep froze ancient block", "number", numFrozen, "hash", hash)
			// Inject all the components into the relevant data tables
			if err := f.(*FreezerRemoteClient).AppendAncient(numFrozen, hash[:], header, body, receipts, td); err != nil {
				break
			}
			ancients = append(ancients, hash)
		}
		// Batch of blocks have been frozen, flush them before wiping from leveldb
		if err := f.Sync(); err != nil {
			log.Crit("Failed to flush frozen tables", "err", err)
		}
		// Wipe out all data from the active database
		batch := db.NewBatch()
		for i := 0; i < len(ancients); i++ {
			// Always keep the genesis block in active database
			if first+uint64(i) != 0 {
				DeleteBlockWithoutNumber(batch, ancients[i], first+uint64(i))
				DeleteCanonicalHash(batch, first+uint64(i))
			}
		}
		if err := batch.Write(); err != nil {
			log.Crit("Failed to delete frozen canonical blocks", "err", err)
		}
		batch.Reset()

		// Wipe out side chains also and track dangling side chians
		var dangling []common.Hash
		for number := first; number < numFrozen; number++ {
			// Always keep the genesis block in active database
			if number != 0 {
				dangling = ReadAllHashes(db, number)
				for _, hash := range dangling {
					log.Trace("Deleting side chain", "number", number, "hash", hash)
					DeleteBlock(batch, hash, number)
				}
			}
		}
		if err := batch.Write(); err != nil {
			log.Crit("Failed to delete frozen side blocks", "err", err)
		}
		batch.Reset()

		// Step into the future and delete and dangling side chains
		if numFrozen > 0 {
			tip := numFrozen
			for len(dangling) > 0 {
				drop := make(map[common.Hash]struct{})
				for _, hash := range dangling {
					log.Debug("Dangling parent from freezer", "number", tip-1, "hash", hash)
					drop[hash] = struct{}{}
				}
				children := ReadAllHashes(db, tip)
				for i := 0; i < len(children); i++ {
					// Dig up the child and ensure it's dangling
					child := ReadHeader(nfdb, children[i], tip)
					if child == nil {
						log.Error("Missing dangling header", "number", tip, "hash", children[i])
						continue
					}
					if _, ok := drop[child.ParentHash]; !ok {
						children = append(children[:i], children[i+1:]...)
						i--
						continue
					}
					// Delete all block data associated with the child
					log.Debug("Deleting dangling block", "number", tip, "hash", children[i], "parent", child.ParentHash)
					DeleteBlock(batch, children[i], tip)
				}
				dangling = children
				tip++
			}
			if err := batch.Write(); err != nil {
				log.Crit("Failed to delete dangling side blocks", "err", err)
			}
		}
		// Log something friendly for the user
		context := []interface{}{
			"blocks", numFrozen - first, "elapsed", common.PrettyDuration(time.Since(start)), "number", numFrozen - 1,
		}
		if n := len(ancients); n > 0 {
			context = append(context, []interface{}{"hash", ancients[n-1]}...)
		}
		log.Info("Deep froze chain segment", context...)

		// Avoid database thrashing with tiny writes
		if numFrozen-first < freezerBatchLimit {
			backoff = true
		}
	}
}
