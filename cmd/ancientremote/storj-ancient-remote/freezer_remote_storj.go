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

	wCacheBlocks *cache
	wCacheHashes *cache
	rCacheBlocks *cache
	rCacheHashes *cache

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

type cache struct {
	max int
	mu  sync.Mutex
	m   map[uint64]interface{}
	sl  []uint64
}

func newCache(max int) *cache {
	return &cache{
		max: max,
		m:   make(map[uint64]interface{}),
		sl:  []uint64{},
	}
}

func (c *cache) reset() {
	c.mu.Lock()
	c.m = make(map[uint64]interface{})
	c.sl = []uint64{}
	c.mu.Unlock()
}

func (c *cache) onto(c2 *cache) {
	if len(c.sl) == 0 {
		return
	}
	c.mu.Lock()
	for k, v := range c.m {
		c2.add(k, v)
	}
	c.mu.Unlock()
}

func (c *cache) len() int {
	return len(c.sl)
}

// FirstN assumes the cache is not empty. Callers must ensure this.
func (c *cache) firstN() uint64 {
	return c.sl[0]
}

func (c *cache) add(n uint64, item interface{}) {
	c.mu.Lock()
	if _, ok := c.m[n]; ok {
		c.mu.Unlock()
		return
	}
	c.m[n] = item
	c.sl = append(c.sl, n)
	if len(c.sl) > 1 {
		// if out of order
		if c.sl[len(c.sl)-2]+1 != c.sl[len(c.sl)-1] {
			sort.Slice(c.sl, func(i, j int) bool {
				return c.sl[i] < c.sl[j]
			})
		}
	}
	if c.max > 0 && c.len() > c.max {
		c.mu.Unlock()
		c.splice(1)
		return
	}
	c.mu.Unlock()
}

func (c *cache) batch(offset, size uint64) []interface{} {
	c.mu.Lock()
	defer c.mu.Unlock()
	max := uint64(c.len())
	if offset >= max {
		return nil
	}
	s := []interface{}{}
	end := offset + size
	if end > max {
		end = max
	}
	for _, n := range c.sl[offset:end] {
		s = append(s, c.m[n])
	}
	return s
}

func (c *cache) truncateFrom(n uint64) {
	c.mu.Lock()
	index := -1
	for i, v := range c.sl {
		if v >= n {
			if index < 0 {
				index = i
			}
			delete(c.m, v)
		}
	}
	if index >= 0 {
		c.sl = c.sl[:index]
	}
	c.mu.Unlock()
}

func (c *cache) splice(firstN uint64) {
	c.mu.Lock()
	if firstN == 0 {
		c.mu.Unlock()
		return
	}
	for _, v := range c.sl[:firstN] {
		delete(c.m, v)
	}
	c.sl = c.sl[firstN:]
	c.mu.Unlock()
}

func (c *cache) get(n uint64) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	v, ok := c.m[n]
	return v, ok
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
	// sanity
	if len(target) > 0 {
		first := target[0].Header.Number.Uint64()
		if first%f.blockObjectGroupSize != 0 {
			panic(fmt.Sprintf("object does not begin at mod: n=%d", target[0].Header.Number.Uint64()))
		}
		f.rCacheBlocks.reset()
		for i, v := range target {
			f.rCacheBlocks.add(first+uint64(i), v)
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
		f.rCacheHashes.reset()
		first := (n / f.hashObjectGroupSize) * f.hashObjectGroupSize
		for i, v := range target {
			f.rCacheHashes.add(first+uint64(i), v)
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
	f.rCacheBlocks.onto(f.wCacheBlocks)
	f.log.Info("Finished pulling blocks cache", "n", n, "size", f.wCacheBlocks.len())
	return nil
}

func (f *freezerRemoteStorj) pullWCacheHashes(ctx context.Context, n uint64) error {
	f.log.Info("Pulling write hashes cache", "n", n)
	err := f.downloadHashesObject(ctx, n)
	if err != nil {
		return err
	}
	f.rCacheHashes.onto(f.wCacheHashes)
	f.log.Info("Finished pulling hashes cache", "n", n, "size", f.wCacheHashes.len())
	return nil
}

func (f *freezerRemoteStorj) findCached(n uint64, kind string) ([]byte, bool) {
	if kind == rawdb.FreezerRemoteHashTable {
		if v, ok := f.wCacheHashes.get(n); ok {
			return v.(common.Hash).Bytes(), ok
		}
		if v, ok := f.rCacheHashes.get(n); ok {
			return v.(common.Hash).Bytes(), ok
		}
	}
	if v, ok := f.wCacheBlocks.get(n); ok {
		return v.(AncientObjectStorj).RLPBytesForKind(kind), ok
	}
	if v, ok := f.rCacheBlocks.get(n); ok {
		return v.(AncientObjectStorj).RLPBytesForKind(kind), ok
	}
	return nil, false
}

// newFreezer creates a chain freezer that moves ancient chain data into
// append-only flat file containers.
func newFreezerRemoteStorj(ctx context.Context, namespace string, access storjAccess, meter logMetrics) (*freezerRemoteStorj, error) {
	var err error

	cacheBlocksMax := 0
	cacheHashesMax := 0
	if storjROnly {
		cacheBlocksMax = int(storjBlocksGroupSize * 2)
		cacheHashesMax = int(storjHashesGroupSize * 2)
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

		wCacheBlocks: newCache(cacheBlocksMax),
		wCacheHashes: newCache(cacheHashesMax),
		rCacheBlocks: newCache(cacheBlocksMax),
		rCacheHashes: newCache(cacheHashesMax),
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

	/*
		By default NewSession will only load credentials from the shared credentials file (~/.aws/credentials).
		If the AWS_SDK_LOAD_CONFIG environment variable is set to a truthy value the Session will be created from the
		configuration values from the shared config (~/.aws/config) and shared credentials (~/.aws/credentials) files.
		Using the NewSessionWithOptions with SharedConfigState set to SharedConfigEnable will create the session as if the
		AWS_SDK_LOAD_CONFIG environment variable was set.
		> https://docs.aws.amazon.com/sdk-for-go/api/aws/session/
	*/

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

	n, _ := f.Ancients(&ctx)
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

	fmt.Println("Internally good ")
	return f, nil
}

// Close terminates the chain freezer, unmapping all the data files.
func (f *freezerRemoteStorj) Close(ctx *context.Context) error {
	// I don't see any Close, Stop, or Quit methods for the AWS service.
	f.quit <- struct{}{}
	return nil
}

// HasAncient returns an indicator whether the specified ancient data exists
// in the freezer.
func (f *freezerRemoteStorj) HasAncient(ctx *context.Context, kind string, number uint64) (bool, error) {
	v, err := f.Ancient(ctx, kind, number)
	if err != nil {
		return false, err
	}
	return v != nil, nil
}

// Ancient retrieves an ancient binary blob from the append-only immutable files.
func (f *freezerRemoteStorj) Ancient(ctx *context.Context, kind string, number uint64) ([]byte, error) {
	if atomic.LoadUint64(f.frozen) <= number {
		return nil, nil
	}

	if v, ok := f.findCached(number, kind); ok {
		return v, nil
	}

	if kind == rawdb.FreezerRemoteHashTable {
		err := f.downloadHashesObject(*ctx, number)
		if err != nil {
			return nil, err
		}
		if v, ok := f.rCacheHashes.get(number); ok {
			return v.(common.Hash).Bytes(), nil
		}
		return nil, fmt.Errorf("%w: #%d (%s)", errOutOfBounds, number, kind)
	}

	// Take from remote
	err := f.downloadBlocksObject(*ctx, number)
	if err != nil {
		return nil, err
	}
	if v, ok := f.rCacheBlocks.get(number); ok {
		return v.(AncientObjectStorj).RLPBytesForKind(kind), nil
	}
	return nil, fmt.Errorf("%w: #%d (%s)", errOutOfBounds, number, kind)
}

// Ancients returns the length of the frozen items.
func (f *freezerRemoteStorj) Ancients(ctx *context.Context) (uint64, error) {
	if f.frozen != nil {
		return atomic.LoadUint64(f.frozen), nil
	}
	f.log.Info("Retrieving ancients number")
	options := &uplink.DownloadOptions{Offset: 0, Length: -1}
	result, err := f.service.DownloadObject(*ctx, f.bucketName(), indexMarker, options)
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
func (f *freezerRemoteStorj) AncientSize(ctx *context.Context, kind string) (uint64, error) {
	// TODO for Storj Go-SDK doesn't support this in a convenient way.
	// This would require listing all objects in the bucket and summing their sizes.
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
func (f *freezerRemoteStorj) AppendAncient(ctx *context.Context, number uint64, hash, header, body, receipts, td []byte) (err error) {
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
	f.wCacheHashes.add(number, common.BytesToHash(hash))
	f.wCacheBlocks.add(number, o)
	atomic.AddUint64(f.frozen, 1)
	// f.log.Info("Finished appending ancient", "frozen", atomic.LoadUint64(f.frozen), "number", number)
	return nil
}

// Truncate discards any recent data above the provided threshold number.
// TODO@meowsbits: handle pagination.
//   ListObjects will only (dubiously? might return millions?) return the first 1000. Need to implement pagination.
//   Also make sure that the Marker is working as expected.
func (f *freezerRemoteStorj) TruncateAncients(ctx *context.Context, items uint64) error {

	f.mu.Lock()
	defer f.mu.Unlock()

	n := atomic.LoadUint64(f.frozen)
	f.log.Info("Truncating ancients", "frozen", n, "target", items, "delta", n-items)
	start := time.Now()

	// How this works:
	// 0. Push the new index marker. If everything goes south from here, we'll just end up overwriting the truncated objects.
	// 1. Ensure or get in cache the object group corresponding to the truncate height.
	// 2. Truncate the cache.
	// 3. Push the cache, overwriting the corresponding object on the remote.
	// 4. Iteratively delete all subsequent objects above the truncated object (ie remove dangling objects).

	var err error

	if !f.readOnly {
		err = f.setIndexMarker(*ctx, items)
		if err != nil {
			return err
		}
	}

	f.rCacheBlocks.reset()
	f.rCacheHashes.reset()

	_, ok := f.wCacheBlocks.get(items)
	if !ok {
		// We DON'T have in the write cache.
		// This means that the target block is at least a batch below our current frozen level.
		// So we need to download the corresponding target batch (clearing the current) and truncate that.
		// The S3 delete iterator will delete all objects above the target.
		f.log.Warn("Target block below current batch")
		f.wCacheBlocks.reset()
		err = f.pullWCacheBlocks(*ctx, items)
		if err != nil {
			return err
		}
	}

	_, ok = f.wCacheHashes.get(items)
	if !ok {
		f.wCacheHashes.reset()
		err = f.pullWCacheHashes(*ctx, items)
		if err != nil {
			return err
		}
	}

	f.wCacheBlocks.truncateFrom(items)
	f.wCacheHashes.truncateFrom(items)

	if f.readOnly {
		atomic.StoreUint64(f.frozen, items)
		f.log.Info("Finished truncating ancients", "elapsed", common.PrettyDuration(time.Since(start)))
		return nil
	}

	err = f.pushWCaches(*ctx)
	if err != nil {
		return err
	}

	atomic.StoreUint64(f.frozen, items)

	// Iteratively delete any dangling _above_ the current object.
	marker := f.blockObjectKeyForN(items + f.blockObjectGroupSize)
	log.Info("Deleting block objects", "marker", marker, "target", items)
	iter := f.service.ListObjects(*ctx, f.bucketName(), &uplink.ListObjectsOptions{Cursor: marker})
	// TODO might be really slow
	for iter.Next() {
		item := iter.Item()
		_, err := f.service.DeleteObject(*ctx, f.bucketName(), item.Key)
		if err != nil {
			return fmt.Errorf("Failure deleting block objects %w", err)
		}
	}

	marker = f.hashObjectKeyForN(items + f.hashObjectGroupSize)
	log.Info("Deleting hash objects", "marker", marker, "target", items)
	iter = f.service.ListObjects(*ctx, f.bucketName(), &uplink.ListObjectsOptions{Cursor: marker})
	// TODO might be really slow
	for iter.Next() {
		item := iter.Item()
		_, err := f.service.DeleteObject(*ctx, f.bucketName(), item.Key)
		if err != nil {
			return fmt.Errorf("Failure deleting hash objects %w", err)
		}
	}

	f.log.Info("Finished truncating ancients", "elapsed", common.PrettyDuration(time.Since(start)))
	return nil
}

func (f *freezerRemoteStorj) pushWCaches(ctx context.Context) error {
	var err error
	err = f.pushCacheBatch(ctx, f.wCacheBlocks, f.blockObjectGroupSize, f.blockObjectKeyForN, f.blockCacheBatchObjectFn)
	if err != nil {
		return err
	}
	err = f.pushCacheBatch(ctx, f.wCacheHashes, f.hashObjectGroupSize, f.hashObjectKeyForN, f.hashCacheBatchObjectFn)
	if err != nil {
		return err
	}
	return nil
}

func (f *freezerRemoteStorj) pushCacheBatch(ctx context.Context, cache *cache, size uint64, keyFn func(uint64) string, batchObjFn func([]interface{}) interface{}) error {
	if cache.len() == 0 {
		return nil
	}
	/*	var uploads []struct {
		key  string
		body []byte
	}*/
	for {
		if cache.len() == 0 {
			break
		}
		var n = cache.firstN()
		if n%size != 0 {
			panic(fmt.Sprintf("bad mod: n=%d r=%d mod=%d len=%d", n, n%size, size, cache.len()))
		}
		batch := cache.batch(0, size)
		b, err := f.encodeObject(batchObjFn(batch))
		if err != nil {
			return err
		}

		// TODO might be really slow
		err = f.uploadObject(ctx, keyFn(n), b)
		if err != nil {
			return err
		}
		// uploads = append(uploads, struct{key: keyFn(n), body: b})
		batchLen := uint64(len(batch))
		if batchLen%size == 0 {
			cache.splice(size)
		} else {
			break
		}
	}
	// TODO might need an upload to be cancelled here if partial failure in batch upload if batched again
	//iter := &s3manager.UploadObjectsIterator{Objects: uploads}
	//err := f.uploader.UploadWithIterator(aws.BackgroundContext(), iter)

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

func (f *freezerRemoteStorj) blockCacheBatchObjectFn(items []interface{}) interface{} {
	batchSet := make([]AncientObjectStorj, len(items))
	for i, v := range items {
		batchSet[i] = v.(AncientObjectStorj)
	}
	return batchSet
}

func (f *freezerRemoteStorj) hashCacheBatchObjectFn(items []interface{}) interface{} {
	batchSet := make([]common.Hash, len(items))
	for i, v := range items {
		batchSet[i] = v.(common.Hash)
	}
	return batchSet
}

// sync flushes all data tables to disk.
func (f *freezerRemoteStorj) Sync(ctx context.Context) error {
	f.mu.Lock()

	n := atomic.LoadUint64(f.frozen)
	lenBlocks := f.wCacheBlocks.len()
	f.log.Info("Syncing ancients", "frozen", n, "blocks", lenBlocks)
	start := time.Now()

	var err error

	if !f.readOnly {
		if lenBlocks > 0 {
			if r := f.wCacheBlocks.firstN() % f.blockObjectGroupSize; r != 0 {
				f.log.Warn("Found out-of-order block cache", "n", f.wCacheBlocks.firstN())
				err = f.pullWCacheBlocks(ctx, f.wCacheBlocks.firstN())
				if err != nil {
					f.mu.Unlock()
					return err
				}
				f.mu.Unlock()
				return f.Sync(ctx)
			}
			if r := f.wCacheHashes.firstN() % f.hashObjectGroupSize; r != 0 {
				f.log.Warn("Found out-of-order hash cache", "n", f.wCacheHashes.firstN())
				err = f.pullWCacheHashes(ctx, f.wCacheHashes.firstN())
				if err != nil {
					f.mu.Unlock()
					return err
				}
				f.mu.Unlock()
				return f.Sync(ctx)
			}
		}

		err = f.pushWCaches(ctx)
		if err != nil {
			f.mu.Unlock()
			return err
		}

		err = f.setIndexMarker(ctx, atomic.LoadUint64(f.frozen))
		if err != nil {
			f.mu.Unlock()
			return err
		}
	}

	elapsed := time.Since(start)
	blocksPerSecond := fmt.Sprintf("%0.2f", float64(lenBlocks)/elapsed.Seconds())

	f.log.Info("Finished syncing ancients", "frozen", n, "blocks", lenBlocks, "elapsed", common.PrettyDuration(elapsed), "bps", blocksPerSecond)
	f.mu.Unlock()
	return err
}
