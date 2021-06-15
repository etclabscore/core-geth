// Copyright 2018 The go-ethereum Authors
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
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/ethdb/leveldb"
	"github.com/ethereum/go-ethereum/ethdb/memorydb"
	"github.com/ethereum/go-ethereum/log"
	"github.com/olekukonko/tablewriter"
)

// freezerdb is a database wrapper that enabled freezer data retrievals.
type freezerdb struct {
	ethdb.KeyValueStore
	ethdb.AncientStore
}

// Close implements io.Closer, closing both the fast key-value store as well as
// the slow ancient tables.
func (frdb *freezerdb) Close() error {
	var errs []error
	if err := frdb.AncientStore.Close(); err != nil {
		errs = append(errs, err)
	}
	if err := frdb.KeyValueStore.Close(); err != nil {
		errs = append(errs, err)
	}
	if len(errs) != 0 {
		return fmt.Errorf("%v", errs)
	}
	return nil
}

// Freeze is a helper method used for external testing to trigger and block until
// a freeze cycle completes, without having to sleep for a minute to trigger the
// automatic background run.
func (frdb *freezerdb) Freeze(threshold uint64) error {
	if frdb.AncientStore.(*freezer).readonly {
		return errReadOnly
	}
	// Set the freezer threshold to a temporary value
	defer func(old uint64) {
		atomic.StoreUint64(&frdb.AncientStore.(*freezer).threshold, old)
	}(atomic.LoadUint64(&frdb.AncientStore.(*freezer).threshold))
	atomic.StoreUint64(&frdb.AncientStore.(*freezer).threshold, threshold)

	// Trigger a freeze cycle and block until it's done
	trigger := make(chan struct{}, 1)
	frdb.AncientStore.(*freezer).trigger <- trigger
	<-trigger
	return nil
}

// nofreezedb is a database wrapper that disables freezer data retrievals.
type nofreezedb struct {
	ethdb.KeyValueStore
}

// HasAncient returns an error as we don't have a backing chain freezer.
func (db *nofreezedb) HasAncient(kind string, number uint64) (bool, error) {
	return false, errNotSupported
}

// Ancient returns an error as we don't have a backing chain freezer.
func (db *nofreezedb) Ancient(kind string, number uint64) ([]byte, error) {
	return nil, errNotSupported
}

// Ancients returns an error as we don't have a backing chain freezer.
func (db *nofreezedb) Ancients() (uint64, error) {
	return 0, errNotSupported
}

// AncientSize returns an error as we don't have a backing chain freezer.
func (db *nofreezedb) AncientSize(kind string) (uint64, error) {
	return 0, errNotSupported
}

// AppendAncient returns an error as we don't have a backing chain freezer.
func (db *nofreezedb) AppendAncient(number uint64, hash, header, body, receipts, td []byte) error {
	return errNotSupported
}

// TruncateAncients returns an error as we don't have a backing chain freezer.
func (db *nofreezedb) TruncateAncients(items uint64) error {
	return errNotSupported
}

// Sync returns an error as we don't have a backing chain freezer.
func (db *nofreezedb) Sync() error {
	return errNotSupported
}

// NewDatabase creates a high level database on top of a given key-value data
// store without a freezer moving immutable chain segments into cold storage.
func NewDatabase(db ethdb.KeyValueStore) ethdb.Database {
	return &nofreezedb{
		KeyValueStore: db,
	}
}

// NewDatabaseWithFreezerRemote creates a high level database on top of a given key-
// value data store with a freezer moving immutable chain segments into cold
// storage.
func NewDatabaseWithFreezerRemote(db ethdb.KeyValueStore, freezerURL string, readonly bool) (ethdb.Database, error) {
	// Create the idle freezer instance
	log.Info("New remote freezer", "freezer", freezerURL)

	frdb, err := newFreezerRemoteClient(freezerURL, readonly)
	if err != nil {
		log.Error("NewDatabaseWithFreezerRemote error", "error", err)
		return nil, err
	}
	// Core-Geth: The validation below is the original and contemporary upstream
	// ethereum/go-ethereum implementation of validations in NewDatabaseWithFreezer. Core-Geth's
	// implementation of the "standard" (built-in FS) freezer initialization has been
	// upgraded to attempt recovery in case an ancient/kv gap is detected, as you can see below.
	// This recovery-if-invalid logic is not echoed below in the case of a remote freezer
	// to limit complexity and stresses on the external implementation, while providing at a loud
	// complaint if db states are determined invalid.
	//
	// Original comment:
	//
	// Since the freezer can be stored separately from the user's key-value database,
	// there's a fairly high probability that the user requests invalid combinations
	// of the freezer and database. Ensure that we don't shoot ourselves in the foot
	// by serving up conflicting data, leading to both datastores getting corrupted.
	//
	//   - If both the freezer and key-value store is empty (no genesis), we just
	//     initialized a new empty freezer, so everything's fine.
	//   - If the key-value store is empty, but the freezer is not, we need to make
	//     sure the user's genesis matches the freezer. That will be checked in the
	//     blockchain, since we don't have the genesis block here (nor should we at
	//     this point care, the key-value/freezer combo is valid).
	//   - If neither the key-value store nor the freezer is empty, cross validate
	//     the genesis hashes to make sure they are compatible. If they are, also
	//     ensure that there's no gap between the freezer and sunsequently leveldb.
	//   - If the key-value store is not empty, but the freezer is we might just be
	//     upgrading to the freezer release, or we might have had a small chain and
	//     not frozen anything yet. Ensure that no blocks are missing yet from the
	//     key-value store, since that would mean we already had an old freezer.
	if kvgenesis, _ := db.Get(headerHashKey(0)); len(kvgenesis) > 0 {
		if frozen, _ := frdb.Ancients(); frozen > 0 {
			// If the freezer already contains something, ensure that the genesis blocks
			// match, otherwise we might mix up freezers across chains and destroy both
			// the freezer and the key-value store.
			frgenesis, err := frdb.Ancient(freezerHashTable, 0)
			if err != nil {
				return nil, fmt.Errorf("failed to retrieve genesis from ancient %v", err)
			} else if !bytes.Equal(kvgenesis, frgenesis) {
				return nil, fmt.Errorf("genesis mismatch: %#x (leveldb) != %#x (ancients)", kvgenesis, frgenesis)
			}
			// Key-value store and freezer belong to the same network. Ensure that they
			// are contiguous, otherwise we might end up with a non-functional freezer.
			if kvhash, _ := db.Get(headerHashKey(frozen)); len(kvhash) == 0 {
				// Subsequent header after the freezer limit is missing from the database.
				// Reject startup is the database has a more recent head.
				if *ReadHeaderNumber(db, ReadHeadHeaderHash(db)) > frozen-1 {
					return nil, fmt.Errorf("gap (#%d) in the chain between ancients and leveldb", frozen)
				}
				// Database contains only older data than the freezer, this happens if the
				// state was wiped and reinited from an existing freezer.
			}
			// Otherwise, key-value store continues where the freezer left off, all is fine.
			// We might have duplicate blocks (crash after freezer write but before key-value
			// store deletion, but that's fine).
		} else {
			// If the freezer is empty, ensure nothing was moved yet from the key-value
			// store, otherwise we'll end up missing data. We check block #1 to decide
			// if we froze anything previously or not, but do take care of databases with
			// only the genesis block.
			if ReadHeadHeaderHash(db) != common.BytesToHash(kvgenesis) {
				// Key-value store contains more data than the genesis block, make sure we
				// didn't freeze anything yet.
				if kvblob, _ := db.Get(headerHashKey(1)); len(kvblob) == 0 {
					return nil, errors.New(`ancient chain segments already extracted from KV, but freezer is empty.
Please set --ancient.rpc to the correct path, and/or review the remote freezer's contents and state.`)
				}
				// Block #1 is still in the database, we're allowed to init a new feezer
			}
			// Otherwise, the head header is still the genesis, we're allowed to init a new
			// feezer.
		}
	}
	// Freezer is consistent with the key-value database, permit combining the two
	if !frdb.readonly {
		go freezeRemote(db, frdb, frdb.threshold, frdb.quit, frdb.trigger)
	}
	return &freezerdb{
		KeyValueStore: db,
		AncientStore:  frdb,
	}, nil
}

// NewDatabaseWithFreezer creates a high level database on top of a given key-
// value data store with a freezer moving immutable chain segments into cold
// storage.
func NewDatabaseWithFreezer(db ethdb.KeyValueStore, freezerStr string, namespace string, readonly bool) (ethdb.Database, error) {
	// Create the idle freezer instance
	frdb, err := newFreezer(freezerStr, namespace, readonly)
	if err != nil {
		return nil, err
	}
	// Since the freezer can be stored separately from the user's key-value database,
	// there's a fairly high probability that the user requests invalid combinations
	// of the freezer and database. Ensure that we don't shoot ourselves in the foot
	// by serving up conflicting data, leading to both datastores getting corrupted.
	//
	//   - If both the freezer and key-value store is empty (no genesis), we just
	//     initialized a new empty freezer, so everything's fine.
	//   - If the key-value store is empty, but the freezer is not, we need to make
	//     sure the user's genesis matches the freezer. That will be checked in the
	//     blockchain, since we don't have the genesis block here (nor should we at
	//     this point care, the key-value/freezer combo is valid).
	//   - If neither the key-value store nor the freezer is empty, cross validate
	//     the genesis hashes to make sure they are compatible. If they are, also
	//     ensure that there's no gap between the freezer and sunsequently leveldb.
	//   - If the key-value store is not empty, but the freezer is we might just be
	//     upgrading to the freezer release, or we might have had a small chain and
	//     not frozen anything yet. Ensure that no blocks are missing yet from the
	//     key-value store, since that would mean we already had an old freezer.
	validateErr := validateFreezerVsKV(frdb, db)
	if validateErr != nil {

		log.Warn("Freezer/KV validation error, attempting freezer repair", "error", validateErr)
		if reperr := frdb.repair(); reperr != nil {
			log.Warn("Freezer repair errored", "error", reperr)

			// Repair did error, AND the validation errored, so return both together because that's double bad.
			return nil, fmt.Errorf("freezer/kv error=%v freezer repair error=%v", validateErr, reperr)
		}

		log.Warn("Freezer repair OK")

		truncateKVtoFreezer(frdb, db)

		// Re-validate the ancient/kv dbs.
		// If still a gap, try removing the kv data back to the ancient level.
		validateErr = validateFreezerVsKV(frdb, db)
	}

	if validateErr != nil && strings.Contains(validateErr.Error(), "gap") {
		// Re-validate again.
		validateErr = validateFreezerVsKV(frdb, db)
	} else if validateErr != nil {
		return nil, validateErr
	}

	if validateErr != nil {
		// If this fails, there's nothing left for us to do.
		log.Warn("KV truncation failed to resuscitate Freezer/KV db gap.")
		return nil, validateErr
	}

	// Freezer is consistent with the key-value database, permit combining the two
	if !frdb.readonly {
		go frdb.freeze(db)
	}
	return &freezerdb{
		KeyValueStore: db,
		AncientStore:  frdb,
	}, nil
}

// NewMemoryDatabase creates an ephemeral in-memory key-value database without a
// freezer moving immutable chain segments into cold storage.
func NewMemoryDatabase() ethdb.Database {
	return NewDatabase(memorydb.New())
}

// NewMemoryDatabaseWithCap creates an ephemeral in-memory key-value database
// with an initial starting capacity, but without a freezer moving immutable
// chain segments into cold storage.
func NewMemoryDatabaseWithCap(size int) ethdb.Database {
	return NewDatabase(memorydb.NewWithCap(size))
}

// NewLevelDBDatabase creates a persistent key-value database without a freezer
// moving immutable chain segments into cold storage.
func NewLevelDBDatabase(file string, cache int, handles int, namespace string, readonly bool) (ethdb.Database, error) {
	db, err := leveldb.New(file, cache, handles, namespace, readonly)
	if err != nil {
		return nil, err
	}
	return NewDatabase(db), nil
}

// NewLevelDBDatabaseWithFreezer creates a persistent key-value database with a
// freezer moving immutable chain segments into cold storage.
func NewLevelDBDatabaseWithFreezerRemote(file string, cache int, handles int, freezerURL string, readonly bool) (ethdb.Database, error) {
	kvdb, err := leveldb.New(file, cache, handles, "eth/db/chaindata", readonly)
	if err != nil {
		return nil, err
	}
	frdb, err := NewDatabaseWithFreezerRemote(kvdb, freezerURL, readonly)
	if err != nil {
		kvdb.Close()
		return nil, err
	}
	return frdb, nil
}

// NewLevelDBDatabaseWithFreezer creates a persistent key-value database with a
// freezer moving immutable chain segments into cold storage.
func NewLevelDBDatabaseWithFreezer(file string, cache int, handles int, freezer string, namespace string, readonly bool) (ethdb.Database, error) {
	kvdb, err := leveldb.New(file, cache, handles, namespace, readonly)
	if err != nil {
		return nil, err
	}
	frdb, err := NewDatabaseWithFreezer(kvdb, freezer, namespace, readonly)
	if err != nil {
		kvdb.Close()
		return nil, err
	}
	return frdb, nil
}

type counter uint64

func (c counter) String() string {
	return fmt.Sprintf("%d", c)
}

func (c counter) Percentage(current uint64) string {
	return fmt.Sprintf("%d", current*100/uint64(c))
}

// stat stores sizes and count for a parameter
type stat struct {
	size  common.StorageSize
	count counter
}

// Add size to the stat and increase the counter by 1
func (s *stat) Add(size common.StorageSize) {
	s.size += size
	s.count++
}

func (s *stat) Size() string {
	return s.size.String()
}

func (s *stat) Count() string {
	return s.count.String()
}

// InspectDatabase traverses the entire database and checks the size
// of all different categories of data.
func InspectDatabase(db ethdb.Database, keyPrefix, keyStart []byte) error {
	it := db.NewIterator(keyPrefix, keyStart)
	defer it.Release()

	var (
		count  int64
		start  = time.Now()
		logged = time.Now()

		// Key-value store statistics
		headers         stat
		bodies          stat
		receipts        stat
		tds             stat
		numHashPairings stat
		hashNumPairings stat
		tries           stat
		codes           stat
		txLookups       stat
		accountSnaps    stat
		storageSnaps    stat
		preimages       stat
		bloomBits       stat
		cliqueSnaps     stat

		// Ancient store statistics
		ancientHeadersSize  common.StorageSize
		ancientBodiesSize   common.StorageSize
		ancientReceiptsSize common.StorageSize
		ancientTdsSize      common.StorageSize
		ancientHashesSize   common.StorageSize

		// Les statistic
		chtTrieNodes   stat
		bloomTrieNodes stat

		// Meta- and unaccounted data
		metadata     stat
		unaccounted  stat
		shutdownInfo stat

		// Totals
		total common.StorageSize
	)
	// Inspect key-value database first.
	for it.Next() {
		var (
			key  = it.Key()
			size = common.StorageSize(len(key) + len(it.Value()))
		)
		total += size
		switch {
		case bytes.HasPrefix(key, headerPrefix) && len(key) == (len(headerPrefix)+8+common.HashLength):
			headers.Add(size)
		case bytes.HasPrefix(key, blockBodyPrefix) && len(key) == (len(blockBodyPrefix)+8+common.HashLength):
			bodies.Add(size)
		case bytes.HasPrefix(key, blockReceiptsPrefix) && len(key) == (len(blockReceiptsPrefix)+8+common.HashLength):
			receipts.Add(size)
		case bytes.HasPrefix(key, headerPrefix) && bytes.HasSuffix(key, headerTDSuffix):
			tds.Add(size)
		case bytes.HasPrefix(key, headerPrefix) && bytes.HasSuffix(key, headerHashSuffix):
			numHashPairings.Add(size)
		case bytes.HasPrefix(key, headerNumberPrefix) && len(key) == (len(headerNumberPrefix)+common.HashLength):
			hashNumPairings.Add(size)
		case len(key) == common.HashLength:
			tries.Add(size)
		case bytes.HasPrefix(key, CodePrefix) && len(key) == len(CodePrefix)+common.HashLength:
			codes.Add(size)
		case bytes.HasPrefix(key, txLookupPrefix) && len(key) == (len(txLookupPrefix)+common.HashLength):
			txLookups.Add(size)
		case bytes.HasPrefix(key, SnapshotAccountPrefix) && len(key) == (len(SnapshotAccountPrefix)+common.HashLength):
			accountSnaps.Add(size)
		case bytes.HasPrefix(key, SnapshotStoragePrefix) && len(key) == (len(SnapshotStoragePrefix)+2*common.HashLength):
			storageSnaps.Add(size)
		case bytes.HasPrefix(key, preimagePrefix) && len(key) == (len(preimagePrefix)+common.HashLength):
			preimages.Add(size)
		case bytes.HasPrefix(key, bloomBitsPrefix) && len(key) == (len(bloomBitsPrefix)+10+common.HashLength):
			bloomBits.Add(size)
		case bytes.HasPrefix(key, BloomBitsIndexPrefix):
			bloomBits.Add(size)
		case bytes.HasPrefix(key, []byte("clique-")) && len(key) == 7+common.HashLength:
			cliqueSnaps.Add(size)
		case bytes.HasPrefix(key, []byte("cht-")) ||
			bytes.HasPrefix(key, []byte("chtIndexV2-")) ||
			bytes.HasPrefix(key, []byte("chtRootV2-")): // Canonical hash trie
			chtTrieNodes.Add(size)
		case bytes.HasPrefix(key, []byte("blt-")) ||
			bytes.HasPrefix(key, []byte("bltIndex-")) ||
			bytes.HasPrefix(key, []byte("bltRoot-")): // Bloomtrie sub
			bloomTrieNodes.Add(size)
		case bytes.Equal(key, uncleanShutdownKey):
			shutdownInfo.Add(size)
		default:
			var accounted bool
			for _, meta := range [][]byte{
				databaseVersionKey, headHeaderKey, headBlockKey, headFastBlockKey, lastPivotKey,
				fastTrieProgressKey, snapshotDisabledKey, snapshotRootKey, snapshotJournalKey,
				snapshotGeneratorKey, snapshotRecoveryKey, txIndexTailKey, fastTxLookupLimitKey,
				uncleanShutdownKey, badBlockKey,
			} {
				if bytes.Equal(key, meta) {
					metadata.Add(size)
					accounted = true
					break
				}
			}
			if !accounted {
				unaccounted.Add(size)
			}
		}
		count++
		if count%1000 == 0 && time.Since(logged) > 8*time.Second {
			log.Info("Inspecting database", "count", count, "elapsed", common.PrettyDuration(time.Since(start)))
			logged = time.Now()
		}
	}
	// Inspect append-only file store then.
	ancientSizes := []*common.StorageSize{&ancientHeadersSize, &ancientBodiesSize, &ancientReceiptsSize, &ancientHashesSize, &ancientTdsSize}
	for i, category := range []string{freezerHeaderTable, freezerBodiesTable, freezerReceiptTable, freezerHashTable, freezerDifficultyTable} {
		if size, err := db.AncientSize(category); err == nil {
			*ancientSizes[i] += common.StorageSize(size)
			total += common.StorageSize(size)
		}
	}
	// Get number of ancient rows inside the freezer
	ancients := counter(0)
	if count, err := db.Ancients(); err == nil {
		ancients = counter(count)
	}
	// Display the database statistic.
	stats := [][]string{
		{"Key-Value store", "Headers", headers.Size(), headers.Count()},
		{"Key-Value store", "Bodies", bodies.Size(), bodies.Count()},
		{"Key-Value store", "Receipt lists", receipts.Size(), receipts.Count()},
		{"Key-Value store", "Difficulties", tds.Size(), tds.Count()},
		{"Key-Value store", "Block number->hash", numHashPairings.Size(), numHashPairings.Count()},
		{"Key-Value store", "Block hash->number", hashNumPairings.Size(), hashNumPairings.Count()},
		{"Key-Value store", "Transaction index", txLookups.Size(), txLookups.Count()},
		{"Key-Value store", "Bloombit index", bloomBits.Size(), bloomBits.Count()},
		{"Key-Value store", "Contract codes", codes.Size(), codes.Count()},
		{"Key-Value store", "Trie nodes", tries.Size(), tries.Count()},
		{"Key-Value store", "Trie preimages", preimages.Size(), preimages.Count()},
		{"Key-Value store", "Account snapshot", accountSnaps.Size(), accountSnaps.Count()},
		{"Key-Value store", "Storage snapshot", storageSnaps.Size(), storageSnaps.Count()},
		{"Key-Value store", "Clique snapshots", cliqueSnaps.Size(), cliqueSnaps.Count()},
		{"Key-Value store", "Singleton metadata", metadata.Size(), metadata.Count()},
		{"Key-Value store", "Shutdown metadata", shutdownInfo.Size(), shutdownInfo.Count()},
		{"Ancient store", "Headers", ancientHeadersSize.String(), ancients.String()},
		{"Ancient store", "Bodies", ancientBodiesSize.String(), ancients.String()},
		{"Ancient store", "Receipt lists", ancientReceiptsSize.String(), ancients.String()},
		{"Ancient store", "Difficulties", ancientTdsSize.String(), ancients.String()},
		{"Ancient store", "Block number->hash", ancientHashesSize.String(), ancients.String()},
		{"Light client", "CHT trie nodes", chtTrieNodes.Size(), chtTrieNodes.Count()},
		{"Light client", "Bloom trie nodes", bloomTrieNodes.Size(), bloomTrieNodes.Count()},
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Database", "Category", "Size", "Items"})
	table.SetFooter([]string{"", "Total", total.String(), " "})
	table.AppendBulk(stats)
	table.Render()

	if unaccounted.size > 0 {
		log.Error("Database contains unaccounted data", "size", unaccounted.size, "count", unaccounted.count)
	}

	return nil
}

func validateFreezerVsKV(freezerdb *freezer, db ethdb.KeyValueStore) error {
	// If the genesis hash is empty, we have a new key-value store, so nothing to
	// validate in this method. If, however, the genesis hash is not nil, compare
	// it to the freezer content.
	if kvgenesis, _ := db.Get(headerHashKey(0)); len(kvgenesis) > 0 {
		if frozen, _ := freezerdb.Ancients(); frozen > 0 {
			// If the freezer already contains something, ensure that the genesis blocks
			// match, otherwise we might mix up freezers across chains and destroy both
			// the freezer and the key-value store.
			frgenesis, err := freezerdb.Ancient(freezerHashTable, 0)
			if err != nil {
				return fmt.Errorf("failed to retrieve genesis from ancient %v", err)
			} else if !bytes.Equal(kvgenesis, frgenesis) {
				return fmt.Errorf("genesis mismatch: %#x (leveldb) != %#x (ancients)", kvgenesis, frgenesis)
			}
			// Key-value store and freezer belong to the same network. Ensure that they
			// are contiguous, otherwise we might end up with a non-functional freezer.
			if kvhash, _ := db.Get(headerHashKey(frozen)); len(kvhash) == 0 {
				// Subsequent header after the freezer limit is missing from the database.
				// Reject startup is the database has a more recent head.
				if headHeaderN := *ReadHeaderNumber(db, ReadHeadHeaderHash(db)); headHeaderN > frozen-1 {
					return fmt.Errorf("gap (chaindb=#%d frozen=#%d) in the chain between ancients and leveldb", headHeaderN, frozen)
				}
				// Database contains only older data than the freezer, this happens if the
				// state was wiped and reinited from an existing freezer.
			}
			// Otherwise, key-value store continues where the freezer left off, all is fine.
			// We might have duplicate blocks (crash after freezer write but before key-value
			// store deletion, but that's fine).
		} else {
			// If the freezer is empty, ensure nothing was moved yet from the key-value
			// store, otherwise we'll end up missing data. We check block #1 to decide
			// if we froze anything previously or not, but do take care of databases with
			// only the genesis block.
			if ReadHeadHeaderHash(db) != common.BytesToHash(kvgenesis) {
				// Key-value store contains more data than the genesis block, make sure we
				// didn't freeze anything yet.
				if kvblob, _ := db.Get(headerHashKey(1)); len(kvblob) == 0 {
					return errors.New("ancient chain segments already extracted, please set --datadir.ancient to the correct path")
				}
				// Block #1 is still in the database, we're allowed to init a new freezer.
			}
			// Otherwise, the head header is still the genesis, we're allowed to init a new
			// freezer.
		}
	}
	return nil
}

func truncateKVtoFreezer(freezerdb *freezer, db ethdb.KeyValueStore) {
	hhh := ReadHeadHeaderHash(db)
	n := *ReadHeaderNumber(db, hhh)
	frozen, _ := freezerdb.Ancients()

	normalize := func(n *uint64) uint64 {
		if n == nil {
			return 0
		}
		return *n
	}

	headHeaderHash := ReadHeadHeaderHash(db)
	headHeader := normalize(ReadHeaderNumber(db, headHeaderHash))

	headFastHash := ReadHeadFastBlockHash(db)
	headFast := normalize(ReadHeaderNumber(db, headFastHash))

	headFullHash := ReadHeadBlockHash(db)
	headFull := normalize(ReadHeaderNumber(db, headFullHash))

	log.Warn("Head header", "number", headHeader, "hash", headHeaderHash)
	log.Warn("Head fast", "number", headFast, "hash", headFastHash)
	log.Warn("Head full", "number", headFull, "hash", headFullHash)

	log.Warn("Persistent Freezer/KV gap: Truncating KV database to freezer height", "ancients", frozen, "kv.head_header_number", n, "kv.head_header_hash", hhh)

	for ; n > frozen-1 && n != 0; n-- {
		for _, hash := range ReadAllHashes(db, n) {
			if n%10000 == 0 {
				log.Warn("Removing KV block data", "n", n, "hash", hash.String())
			}
			DeleteBlock(db, hash, n)
			DeleteCanonicalHash(db, n)
		}
	}
	log.Warn("Finished KV truncation")

	data, _ := freezerdb.Ancient(freezerHashTable, n)
	h := common.BytesToHash(data)

	// If h is the empty common hash, then when the headHeaderHash gets read, whoever's reading it isn't going to like that.
	// This logic doesn't check for that because there's really nothing that can be sensibly done in this scope,
	// and it seems reasonable to think that when a higher level function like `loadLastState` finds an empty hash in the
	// headHeaderHash value, it's going to bark pretty loudly and probably just roll the whole thing (database(s)) back since the
	// ancient database would appear to be screwy beyond repair since it lied about what frozen headers it had.
	// So we're just gonna write this sucker.
	log.Warn("Writing KV head header", "hash", h.String())
	WriteHeadHeaderHash(db, h)

	// If we had nonzero values for full and/or fast blocks, infer that preceding states will still be valid.
	if headFast != 0 {
		WriteHeadFastBlockHash(db, h)
	}
	if headFull != 0 {
		WriteHeadBlockHash(db, h)
	}
}
