package main

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"storj.io/uplink"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/metrics"
	"github.com/ethereum/go-ethereum/rlp"
	lru "github.com/hashicorp/golang-lru"
)

const (
	storjEncodeJSON   = ".json"
	storjEncodeJSONGZ = ".json.gz"
	indexMarker       = "index-marker"
)

var (
	storjROnly           = false
	storjBlocksGroupSize = uint64(32 * 32)
	storjHashesGroupSize = uint64(32 * 32 * 32)
	storjEncoding        = storjEncodeJSONGZ
	errOutOfBounds       = errors.New("out of bounds")
	errNotSupported      = errors.New("this operation is not supported")
	errOutOrderInsertion = errors.New("the append operation is out-order")
)

func init() {
	if v := os.Getenv("GETH_FREEZER_STORJ_BLOCK_GROUP_SIZE"); v != "" {
		i, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			panic(err)
		}
		storjBlocksGroupSize = i
	}
	if v := os.Getenv("GETH_FREEZER_STORJ_HASH_GROUP_SIZE"); v != "" {
		i, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			panic(err)
		}
		storjHashesGroupSize = i
	}
	if v := os.Getenv("GETH_FREEZER_STORJ_ENCODING"); v != "" {
		storjEncoding = v
	}
	if v := os.Getenv("GETH_FREEZER_STORJ_R_ONLY"); v != "" {
		storjROnly = true
	}
}

type logMetrics struct {
	readMeter  metrics.Meter // Meter for measuring the effective amount of data read
	writeMeter metrics.Meter // Meter for measuring the effective amount of data written
	sizeGauge  metrics.Gauge // Gauge for tracking the combined size of all freezer tables
}

type storjAccess struct {
	passphrase string
	apiKey     string
	satellite  string
}

type freezerRemoteStorj struct {
	session *uplink.Access
	service *uplink.Project

	namespace string
	quit      chan struct{}
	mu        sync.Mutex

	readOnly bool

	readMeter  metrics.Meter // Meter for measuring the effective amount of data read
	writeMeter metrics.Meter // Meter for measuring the effective amount of data written
	sizeGauge  metrics.Gauge // Gauge for tracking the combined size of all freezer tables

	// uploader   *s3manager.Uploader
	// downloader *s3manager.Downloader

	frozen *uint64 // the length of the frozen blocks (next appended must == val)

	// TODO: Reusable gzip r/w
	encoding             string
	blockObjectGroupSize uint64 // how many blocks to include in a single S3 object
	hashObjectGroupSize  uint64

	wCacheBlocks *lru.Cache
	wCacheHashes *lru.Cache
	rCacheBlocks *lru.Cache
	rCacheHashes *lru.Cache

	log log.Logger
}

// AncientObjectStorj describes the storage encoding unit for a 'block'.
// These objects are grouped in an array for storage with size determined by storjBlocksGroupSize.
type AncientObjectStorj struct {
	Hash       common.Hash                `json:"hash"`
	Header     *types.Header              `json:"header"`
	Body       *types.Body                `json:"body"`
	Receipts   []*types.ReceiptForStorage `json:"receipts"`
	Difficulty *big.Int                   `json:"difficulty"`
}

// NewAncientObjectStorj reverses decodes the incoming RLP bytes, applying the values
// to its own type definition's fields.
func NewAncientObjectStorj(hashB, headerB, bodyB, receiptsB, difficultyB []byte) (AncientObjectStorj, error) {
	var err error

	hash := common.BytesToHash(hashB)

	header := &types.Header{}
	err = rlp.DecodeBytes(headerB, header)
	if err != nil {
		return AncientObjectStorj{}, err
	}
	body := &types.Body{}
	err = rlp.DecodeBytes(bodyB, body)
	if err != nil {
		return AncientObjectStorj{}, err
	}
	receipts := []*types.ReceiptForStorage{}
	err = rlp.DecodeBytes(receiptsB, &receipts)
	if err != nil {
		return AncientObjectStorj{}, err
	}
	difficulty := new(big.Int)
	err = rlp.DecodeBytes(difficultyB, difficulty)
	if err != nil {
		return AncientObjectStorj{}, err
	}
	return AncientObjectStorj{
		Hash:       hash,
		Header:     header,
		Body:       body,
		Receipts:   receipts,
		Difficulty: difficulty,
	}, nil
}

// RLPBytesForKind return the RLP encoded bytes expected by the AncientStore interface
// for a given 'kind' on the block object.
func (o AncientObjectStorj) RLPBytesForKind(kind string) []byte {
	switch kind {
	case rawdb.FreezerRemoteHashTable:
		return o.Hash.Bytes()
	case rawdb.FreezerRemoteHeaderTable:
		b, err := rlp.EncodeToBytes(o.Header)
		if err != nil {
			log.Crit("Failed to RLP encode block header", "err", err)
		}
		return b
	case rawdb.FreezerRemoteBodiesTable:
		b, err := rlp.EncodeToBytes(o.Body)
		if err != nil {
			log.Crit("Failed to RLP encode block body", "err", err)
		}
		return b
	case rawdb.FreezerRemoteReceiptTable:
		b, err := rlp.EncodeToBytes(o.Receipts)
		if err != nil {
			log.Crit("Failed to RLP encode block receipts", "err", err)
		}
		return b
	case rawdb.FreezerRemoteDifficultyTable:
		b, err := rlp.EncodeToBytes(o.Difficulty)
		if err != nil {
			log.Crit("Failed to RLP encode block difficulty", "err", err)
		}
		return b
	default:
		panic(fmt.Sprintf("unknown kind: %s", kind))
	}
}

func storjKeyBlock(number uint64) string {
	// Keep blocks in a dir.
	// This namespaces the resource, separating it from the 'index-marker' object.
	return fmt.Sprintf("blocks/%09d%s", number, storjEncoding)
}

func storjKeyHash(number uint64) string {
	return fmt.Sprintf("hashes/%09d%s", number, storjEncoding)
}

func (f *freezerRemoteStorj) blockObjectKeyForN(n uint64) string {
	return storjKeyBlock((n / f.blockObjectGroupSize) * f.blockObjectGroupSize) // 0, 32, 64, 96, ...
}

func (f *freezerRemoteStorj) hashObjectKeyForN(n uint64) string {
	return storjKeyHash((n / f.hashObjectGroupSize) * f.hashObjectGroupSize)
}

// TODO: this is superfluous now; bucket names must be user-configured
func (f *freezerRemoteStorj) bucketName() string {
	return fmt.Sprintf("%s", f.namespace)
}

func (f *freezerRemoteStorj) initializeBucket(ctx context.Context) error {
	start := time.Now()
	_, err := f.service.EnsureBucket(ctx, f.bucketName())
	if err != nil {
		return fmt.Errorf("Could not ensure bucket %v", err)
	}
	f.log.Info("Bucket ensured", "name", f.bucketName(), "elapsed", common.PrettyDuration(time.Since(start)))
	return nil
}

func (f *freezerRemoteStorj) downloadBlocksObject(ctx context.Context, n uint64) error {
	key := f.blockObjectKeyForN(n)
	result, err := f.service.DownloadObject(ctx, f.bucketName(), key, &uplink.DownloadOptions{Length: -1})
	if err != nil {
		f.log.Error("Download error", "method", "downloadBlocksObject", "error", err, "key", key)
		return err
	}
	target := []AncientObjectStorj{}
	buf, err := ioutil.ReadAll(result)
	if err != nil {
		f.log.Error("Failed to read buffer", "method", "downloadBlocksObject", "error", err, "key", key)
	}
	err = f.decodeObject(buf, &target)
	if err != nil {
		return err
	}

	if len(target) > 0 {
		for _, v := range target {
			// Ignore any persisted data above current frozen level.
			// This is truncated data that hasn't been overwritten yet.
			if v.Header.Number.Uint64() >= atomic.LoadUint64(f.frozen) {
				continue
			}
			f.rCacheBlocks.Add(v.Header.Number.Uint64(), v)
		}
	}
	return nil
}

func (f *freezerRemoteStorj) downloadHashesObject(ctx context.Context, n uint64) error {
	var buf []byte
	key := f.hashObjectKeyForN(n)
	result, err := f.service.DownloadObject(ctx, f.bucketName(), key, &uplink.DownloadOptions{Length: -1})
	if err != nil {
		f.log.Error("Download error", "method", "downloadHashesObject", "error", err, "key", key)
		return err
	}
	target := []common.Hash{}
	buf, err = ioutil.ReadAll(result)
	if err != nil {
		return err
	}
	err = f.decodeObject(buf, &target)
	if err != nil {
		return err
	}
	if len(target) > 0 {
		first := n - (n % f.hashObjectGroupSize)
		for i, v := range target {
			n := first + uint64(i)
			// Ignore any persisted data above current frozen level.
			// This is truncated data that hasn't been overwritten yet.
			if n >= atomic.LoadUint64(f.frozen) {
				continue
			}
			f.rCacheHashes.Add(n, v)
		}
	}
	return nil
}

func (f *freezerRemoteStorj) pullWCacheBlocks(ctx context.Context, n uint64) error {
	f.log.Info("Pulling write blocks cache", "n", n)
	err := f.downloadBlocksObject(ctx, n)
	if err != nil {
		return err
	}
	for _, k := range f.rCacheBlocks.Keys() {
		v, _ := f.rCacheBlocks.Get(k)
		f.wCacheBlocks.Add(k, v)
	}
	f.log.Info("Finished pulling blocks cache", "n", n, "size", f.wCacheBlocks.Len())
	return nil
}

func (f *freezerRemoteStorj) pullWCacheHashes(ctx context.Context, n uint64) error {
	f.log.Info("Pulling write hashes cache", "n", n)
	err := f.downloadHashesObject(ctx, n)
	if err != nil {
		return err
	}
	for _, k := range f.rCacheHashes.Keys() {
		v, _ := f.rCacheHashes.Get(k)
		f.wCacheHashes.Add(k, v)
	}
	f.log.Info("Finished pulling hashes cache", "n", n, "size", f.wCacheHashes.Len())
	return nil
}

func (f *freezerRemoteStorj) findCached(n uint64, kind string) ([]byte, bool) {
	if kind == rawdb.FreezerRemoteHashTable {
		if v, ok := f.wCacheHashes.Get(n); ok {
			return v.(common.Hash).Bytes(), ok
		}
		if v, ok := f.rCacheHashes.Get(n); ok {
			return v.(common.Hash).Bytes(), ok
		}
	}
	if v, ok := f.wCacheBlocks.Get(n); ok {
		return v.(AncientObjectStorj).RLPBytesForKind(kind), ok
	}
	if v, ok := f.rCacheBlocks.Get(n); ok {
		return v.(AncientObjectStorj).RLPBytesForKind(kind), ok
	}
	return nil, false
}

// newFreezer creates a chain freezer that moves ancient chain data into
// append-only flat file containers.
func newFreezerRemoteStorj(ctx context.Context, namespace string, access storjAccess, meter logMetrics) (*freezerRemoteStorj, error) {
	var err error

	// Set default cache sizes.
	// Sizes reflect max number of entries in the cache.
	// Cache size minimum must be greater than or equal to the group size * 2,
	// and should not be lower than 2048 * 2 because of how the ancient store rhythm during sync is.
	blockCacheSize := int(storjBlocksGroupSize) * 2
	hashCacheSize := int(storjHashesGroupSize) * 2
	if blockCacheSize < 2048*2 {
		blockCacheSize = 2048 * 2
	}
	if hashCacheSize < 2048*2 {
		hashCacheSize = 2048 * 2
	}

	rBlockCache, err := lru.New(blockCacheSize)
	if err != nil {
		return nil, err
	}
	rHashCache, err := lru.New(hashCacheSize)
	if err != nil {
		return nil, err
	}
	wBlockCache, err := lru.New(blockCacheSize)
	if err != nil {
		return nil, err
	}
	wHashCache, err := lru.New(hashCacheSize)
	if err != nil {
		return nil, err
	}

	f := &freezerRemoteStorj{
		namespace:  namespace,
		quit:       make(chan struct{}),
		readMeter:  meter.readMeter,
		writeMeter: meter.writeMeter,
		sizeGauge:  meter.sizeGauge,

		readOnly: storjROnly,

		// Globals for now. Should probably become CLI flags.
		// Maybe Remote Freezers need a config struct.
		blockObjectGroupSize: storjBlocksGroupSize,
		hashObjectGroupSize:  storjHashesGroupSize,

		wCacheBlocks: wBlockCache,
		wCacheHashes: wHashCache,
		rCacheBlocks: rBlockCache,
		rCacheHashes: rHashCache,
		log:          log.New("remote", "storj"),
	}

	switch storjEncoding {
	case storjEncodeJSONGZ:
		f.encoding = storjEncodeJSONGZ
	case storjEncodeJSON:
		f.encoding = storjEncodeJSON
	default:
		return nil, fmt.Errorf("unknown encoding: %s", storjEncoding)
	}

	f.session, err = uplink.RequestAccessWithPassphrase(ctx, access.satellite, access.apiKey, access.passphrase)
	if err != nil {
		fmt.Println(access.apiKey)
		return nil, fmt.Errorf("Access request failed %w", err)
	}
	f.service, err = uplink.OpenProject(ctx, f.session)
	if err != nil {
		return nil, fmt.Errorf("could not open project: %v", err)
	}
	_, err = f.service.EnsureBucket(ctx, namespace)
	if err != nil {
		return nil, fmt.Errorf("could not ensure bucket existence %v", err)
	}

	n, _ := f.Ancients(ctx)
	f.frozen = &n

	if n > 0 {
		err = f.pullWCacheBlocks(ctx, n)
		if err != nil {
			return f, err
		}

		err = f.pullWCacheHashes(ctx, n)
		if err != nil {
			return f, err
		}
	}

	return f, nil
}

// Close terminates the chain freezer, unmapping all the data files.
func (f *freezerRemoteStorj) Close(ctx context.Context) error {
	// I don't see any Close, Stop, or Quit methods for the AWS service.
	f.quit <- struct{}{}
	return nil
}

// HasAncient returns an indicator whether the specified ancient data exists
// in the freezer.
func (f *freezerRemoteStorj) HasAncient(ctx context.Context, kind string, number uint64) (bool, error) {
	v, err := f.Ancient(ctx, kind, number)
	if err != nil {
		return false, err
	}
	return v != nil, nil
}

// Ancient retrieves an ancient binary blob from the append-only immutable files.
func (f *freezerRemoteStorj) Ancient(ctx context.Context, kind string, number uint64) ([]byte, error) {
	if atomic.LoadUint64(f.frozen) <= number {
		return nil, nil
	}

	if v, ok := f.findCached(number, kind); ok {
		return v, nil
	}

	if kind == rawdb.FreezerRemoteHashTable {
		err := f.downloadHashesObject(ctx, number)
		if err != nil {
			return nil, err
		}
		if v, ok := f.rCacheHashes.Get(number); ok {
			return v.(common.Hash).Bytes(), nil
		}
		return nil, fmt.Errorf("%w: #%d (%s)", errOutOfBounds, number, kind)
	}

	// Take from remote
	err := f.downloadBlocksObject(ctx, number)
	if err != nil {
		return nil, err
	}
	if v, ok := f.rCacheBlocks.Get(number); ok {
		return v.(AncientObjectStorj).RLPBytesForKind(kind), nil
	}
	return nil, fmt.Errorf("%w: #%d (%s)", errOutOfBounds, number, kind)
}

// Ancients returns the length of the frozen items.
func (f *freezerRemoteStorj) Ancients(ctx context.Context) (uint64, error) {
	if f.frozen != nil {
		return atomic.LoadUint64(f.frozen), nil
	}
	f.log.Info("Retrieving ancients number")
	options := &uplink.DownloadOptions{Offset: 0, Length: -1}
	result, err := f.service.DownloadObject(ctx, f.bucketName(), indexMarker, options)
	if err != nil {
		f.log.Error("DownloadObject error", "method", "Ancients", "error", err)
		return 0, err
	}
	contents, err := ioutil.ReadAll(result)
	if err != nil {
		return 0, err
	}
	s := strings.TrimSpace(string(contents))
	i, err := strconv.ParseUint(s, 10, 64)
	f.log.Info("Finished retrieving ancients num", "s", s, "n", i, "err", err)
	return i, err
}

// AncientSize returns the ancient size of the specified category.
func (f *freezerRemoteStorj) AncientSize(ctx context.Context, kind string) (uint64, error) {
	// TODO for Storj Go-SDK doesn't support this in a convenient way.
	// This would might require listing all objects in the bucket and summing their sizes.
	// This method is only used in the InspectDatabase function, which isn't that
	// important.
	return 0, errNotSupported
}

func (f *freezerRemoteStorj) uploadObject(ctx context.Context, key string, object []byte) error {
	upload, err := f.service.UploadObject(ctx, f.bucketName(), key, &uplink.UploadOptions{})
	if err != nil {
		f.log.Error("UploadObject error", "method", "uploadObject", "error", err)
	}
	_, err = upload.Write(object)
	if err != nil {
		f.log.Error("Write error", "method", "uploadObject", "error", err)
	}
	err = upload.Commit()
	if err != nil {
		f.log.Error("Commit error", "method", "uploadObject", "error", err)
	}
	return err

}
func (f *freezerRemoteStorj) setIndexMarker(ctx context.Context, number uint64) error {
	f.log.Info("Setting index marker", "number", number)
	numberStr := strconv.FormatUint(number, 10)
	// TODO nice expiry here might be good for upload options
	upload, err := f.service.UploadObject(ctx, f.bucketName(), indexMarker, &uplink.UploadOptions{})
	if err != nil {
		f.log.Error("UploadObject error", "method", "setIndexMarker", "error", err)
	}
	_, err = upload.Write([]byte(numberStr))
	if err != nil {
		f.log.Error("Write error", "method", "setIndexMarker", "error", err)
	}
	err = upload.Commit()
	if err != nil {
		f.log.Error("Commit error", "method", "setIndexMarker", "error", err)
	}
	return err
}

// AppendAncient injects all binary blobs belong to block at the end of the
// append-only immutable table files.
//
// Notably, this function is lock free but kind of thread-safe. All out-of-order
// injection will be rejected. But if two injections with same number happen at
// the same time, we can get into the trouble.
func (f *freezerRemoteStorj) AppendAncient(ctx context.Context, number uint64, hash, header, body, receipts, td []byte) (err error) {
	// Ensure the binary blobs we are appending is continuous with freezer.
	if atomic.LoadUint64(f.frozen) != number {
		return errOutOrderInsertion
	}
	// f.log.Info("Appending ancient", "frozen", atomic.LoadUint64(f.frozen), "number", number)
	f.mu.Lock()
	defer f.mu.Unlock()
	o, err := NewAncientObjectStorj(hash, header, body, receipts, td)
	if err != nil {
		return err
	}

	// Append to both read and write caches.
	f.rCacheHashes.Add(number, common.BytesToHash(hash))
	f.rCacheBlocks.Add(number, o)
	f.wCacheHashes.Add(number, common.BytesToHash(hash))
	f.wCacheBlocks.Add(number, o)
	atomic.AddUint64(f.frozen, 1)
	f.log.Trace("Finished appending ancient", "frozen", atomic.LoadUint64(f.frozen), "number", number)
	return nil
}

// TruncateAncients discards any recent data above the provided threshold number.
func (f *freezerRemoteStorj) TruncateAncients(ctx context.Context, items uint64) error {

	var err error
	f.mu.Lock()
	defer f.mu.Unlock()

	n := atomic.LoadUint64(f.frozen)
	f.log.Info("Truncating ancients", "frozen", n, "target", items, "delta", n-items)
	start := time.Now()

	if !f.readOnly {
		err = f.setIndexMarker(ctx, items)
		if err != nil {
			return err
		}
	}

	for _, c := range []*lru.Cache{
		f.rCacheBlocks,
		f.rCacheHashes,
		f.wCacheBlocks,
		f.wCacheHashes,
	} {
		keys := c.Keys()
		for _, k := range keys {
			if k.(uint64) >= items {
				c.Remove(k)
			}
		}
	}

	atomic.StoreUint64(f.frozen, items)
	f.log.Info("Finished truncating ancients", "elapsed", common.PrettyDuration(time.Since(start)))
	return nil
}

func (f *freezerRemoteStorj) pushWCaches(ctx context.Context) error {
	var err error
	err = f.pushCacheGroups(ctx, f.wCacheBlocks, f.blockObjectGroupSize, f.blockObjectKeyForN, f.blockCacheGroupObjectFn)
	if err != nil {
		return err
	}
	err = f.pushCacheGroups(ctx, f.wCacheHashes, f.hashObjectGroupSize, f.hashObjectKeyForN, f.hashCacheGroupObjectFn)
	if err != nil {
		return err
	}
	return nil
}

// cacheSortUint64Keys assumes that the cache uses exclusively uint64 keys,
// and returns them in an ascending order.
func cacheSortUint64Keys(cache *lru.Cache) []interface{} {
	keys := cache.Keys()
	sort.Slice(keys, func(i, j int) bool {
		return keys[i].(uint64) < keys[j].(uint64)
	})
	return keys
}

// cacheKeyGroups returns groups of sorted keys, where each
// group contains a maximum groupSize elements.
// This is used for grouping write objects into respective Storj-object groups.
func cacheKeyGroups(c *lru.Cache, groupSize uint64) (groups [][]uint64) {
	keys := cacheSortUint64Keys(c)
	group := []uint64{}
	for _, k := range keys {
		if uint64(len(group)) == groupSize {
			groups = append(groups, group)
			group = []uint64{}
		}
		group = append(group, k.(uint64))
	}
	groups = append(groups, group)
	return groups

}

// TODO here is where we would implement batching for writes
func (f *freezerRemoteStorj) pushCacheGroups(ctx context.Context, cache *lru.Cache, size uint64, keyFn func(uint64) string, groupObjectFn func([]interface{}) interface{}) error {
	if cache.Len() == 0 {
		return nil
	}
	groups := cacheKeyGroups(cache, size)
	for _, keyGroup := range groups {
		if len(keyGroup) == 0 {
			continue
		}
		n := keyGroup[0]
		if n%size != 0 {
			log.Crit("Non-mod group leader", "n", n, "object.len", len(keyGroup))
		}
		object := []interface{}{}
		for _, key := range keyGroup {
			v, _ := cache.Get(key)
			object = append(object, v)
		}
		b, err := f.encodeObject(groupObjectFn(object))
		if err != nil {
			return err
		}

		// TODO might be slow would implement this as a part of batch
		err = f.uploadObject(ctx, keyFn(n), b)
		if err != nil {
			return err
		}
	}

	for _, keyGroup := range groups {
		// Not a complete group (storj-object), don't delete.
		if uint64(len(keyGroup)) != size {
			continue
		}
		for _, key := range keyGroup {
			cache.Remove(key)
		}
	}

	return nil
}

func (f *freezerRemoteStorj) encodeObject(any interface{}) ([]byte, error) {
	b, err := json.Marshal(any)
	if err != nil {
		return nil, err
	}
	if f.encoding == storjEncodeJSONGZ {
		w := bytes.NewBuffer([]byte{})
		gzW, _ := gzip.NewWriterLevel(w, gzip.BestCompression)
		_, err = gzW.Write(b)
		if err != nil {
			gzW.Close()
			return nil, err
		}
		gzW.Close()
		b = w.Bytes()
	}
	return b, nil
}

func (f *freezerRemoteStorj) decodeObject(input []byte, target interface{}) error {
	if f.encoding == storjEncodeJSONGZ {
		r, err := gzip.NewReader(bytes.NewBuffer(input))
		if err != nil {
			return err
		}
		defer r.Close()
		input, err = ioutil.ReadAll(r)
		if err != nil {
			return err
		}
	}
	err := json.Unmarshal(input, target)
	if err != nil {
		return err
	}
	return nil
}

func (f *freezerRemoteStorj) blockCacheGroupObjectFn(items []interface{}) interface{} {
	batchSet := make([]AncientObjectStorj, len(items))
	for i, v := range items {
		batchSet[i] = v.(AncientObjectStorj)
	}
	return batchSet
}

func (f *freezerRemoteStorj) hashCacheGroupObjectFn(items []interface{}) interface{} {
	batchSet := make([]common.Hash, len(items))
	for i, v := range items {
		batchSet[i] = v.(common.Hash)
	}
	return batchSet
}

// sync flushes all data tables to disk.
func (f *freezerRemoteStorj) Sync(ctx context.Context) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	n := atomic.LoadUint64(f.frozen)
	lenBlocks := f.wCacheBlocks.Len()
	f.log.Info("Syncing ancients", "frozen", n, "blocks", lenBlocks)
	start := time.Now()

	var err error

	if !f.readOnly {

		err = f.pushWCaches(ctx)
		if err != nil {
			return err
		}

		err = f.setIndexMarker(ctx, atomic.LoadUint64(f.frozen))
		if err != nil {
			return err
		}
	}

	elapsed := time.Since(start)
	blocksPerSecond := fmt.Sprintf("%0.2f", float64(lenBlocks)/elapsed.Seconds())

	f.log.Info("Finished syncing ancients", "frozen", n, "blocks", lenBlocks, "elapsed", common.PrettyDuration(elapsed), "bps", blocksPerSecond)
	return err
}
