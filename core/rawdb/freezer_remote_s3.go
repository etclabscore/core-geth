// Copyright 2019 The go-ethereum Authors
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
	"encoding/json"
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

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/metrics"
	"github.com/ethereum/go-ethereum/rlp"
)

type freezerRemoteS3 struct {
	session *session.Session
	service *s3.S3

	namespace string
	quit      chan struct{}
	mu        sync.Mutex

	readMeter  metrics.Meter // Meter for measuring the effective amount of data read
	writeMeter metrics.Meter // Meter for measuring the effective amount of data written
	sizeGauge  metrics.Gauge // Gauge for tracking the combined size of all freezer tables

	uploader   *s3manager.Uploader
	downloader *s3manager.Downloader

	frozen               *uint64 // the length of the frozen blocks (next appended must == val)
	blockObjectGroupSize uint64  // how many blocks to include in a single S3 object
	hashObjectGroupSize  uint64

	appendCacheBlocks *Cache
	appendCacheHashes *Cache
	getCacheBlocks    *Cache
	getCacheHashes    *Cache

	retrievedBlocks    map[uint64]AncientObjectS3
	retrievedBlockLock sync.Mutex

	blockCache  map[uint64]AncientObjectS3
	blockCacheS []uint64

	hashCache map[uint64]common.Hash

	cacheLock sync.Mutex

	log log.Logger
}

type AncientObjectS3 struct {
	Hash       common.Hash                `json:"hash"`
	Header     *types.Header              `json:"header"`
	Body       *types.Body                `json:"body"`
	Receipts   []*types.ReceiptForStorage `json:"receipts"`
	Difficulty *big.Int                   `json:"difficulty"`
}

func NewAncientObjectS3(hashB, headerB, bodyB, receiptsB, difficultyB []byte) (*AncientObjectS3, error) {
	var err error

	hash := common.BytesToHash(hashB)

	header := &types.Header{}
	err = rlp.DecodeBytes(headerB, header)
	if err != nil {
		return nil, err
	}
	body := &types.Body{}
	err = rlp.DecodeBytes(bodyB, body)
	if err != nil {
		return nil, err
	}
	receipts := []*types.ReceiptForStorage{}
	err = rlp.DecodeBytes(receiptsB, &receipts)
	if err != nil {
		return nil, err
	}
	difficulty := new(big.Int)
	err = rlp.DecodeBytes(difficultyB, difficulty)
	if err != nil {
		return nil, err
	}
	return &AncientObjectS3{
		Hash:       hash,
		Header:     header,
		Body:       body,
		Receipts:   receipts,
		Difficulty: difficulty,
	}, nil
}

func (o *AncientObjectS3) RLPBytesForKind(kind string) []byte {
	switch kind {
	case freezerHashTable:
		return o.Hash.Bytes()
	case freezerHeaderTable:
		b, err := rlp.EncodeToBytes(o.Header)
		if err != nil {
			log.Crit("Failed to RLP encode block header", "err", err)
		}
		return b
	case freezerBodiesTable:
		b, err := rlp.EncodeToBytes(o.Body)
		if err != nil {
			log.Crit("Failed to RLP encode block body", "err", err)
		}
		return b
	case freezerReceiptTable:
		b, err := rlp.EncodeToBytes(o.Receipts)
		if err != nil {
			log.Crit("Failed to RLP encode block receipts", "err", err)
		}
		return b
	case freezerDifficultyTable:
		b, err := rlp.EncodeToBytes(o.Difficulty)
		if err != nil {
			log.Crit("Failed to RLP encode block difficulty", "err", err)
		}
		return b
	default:
		panic(fmt.Sprintf("unknown kind: %s", kind))
	}
}

type Cache struct {
	mu sync.Mutex
	m  map[uint64]interface{}
	sl []uint64
}

func NewCache() *Cache {
	return &Cache{
		m:  make(map[uint64]interface{}),
		sl: []uint64{},
	}
}

func (c *Cache) Reset() {
	c.mu.Lock()
	c.m = make(map[uint64]interface{})
	c.sl = []uint64{}
	c.mu.Unlock()
}

func (c *Cache) Set(start uint64, items []interface{}) {
	c.mu.Lock()
	c.sl = make([]uint64, len(items))
	for i, v := range items {
		n := start + uint64(i)
		c.m[n] = v
		c.sl[i] = n
	}
	c.mu.Unlock()
}

func (c *Cache) Add(n uint64, item interface{}) {
	c.mu.Lock()
	if _, ok := c.m[n]; ok {
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
	c.mu.Unlock()
}

func (c *Cache) Batch(offset, size uint64) interface{} {
	c.mu.Lock()
	s := []interface{}{}
	for _, n := range c.sl[offset : offset+size] {
		s = append(s, c.m[n])
	}
	c.mu.Unlock()
	return s
}

func (c *Cache) TruncateAbove(n uint64) {
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
	c.sl = c.sl[:index]
	c.mu.Unlock()
}

func (c *Cache) Slice(first int) {
	c.mu.Lock()
	for _, v := range c.sl {
		delete(c.m, v)
	}
	c.sl = c.sl[first:]
	c.mu.Unlock()
}

func (c *Cache) Get(n uint64) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	v, ok := c.m[n]
	return v, ok
}

func awsKeyBlock(number uint64) string {
	// Keep blocks in a dir.
	// This namespaces the resource, separating it from the 'index-marker' object.
	return fmt.Sprintf("blocks/%09d.json", number)
}

func awsKeyHash(number uint64) string {
	return fmt.Sprintf("hashes/%09d.json", number)
}

func (f *freezerRemoteS3) blockObjectKeyForN(n uint64) string {
	return awsKeyBlock((n / f.blockObjectGroupSize) * f.blockObjectGroupSize) // 0, 32, 64, 96, ...
}

func (f *freezerRemoteS3) hashObjectKeyForN(n uint64) string {
	return awsKeyHash((n / f.hashObjectGroupSize) * f.hashObjectGroupSize)
}

// TODO: this is superfluous now; bucket names must be user-configured
func (f *freezerRemoteS3) bucketName() string {
	return fmt.Sprintf("%s", f.namespace)
}

func (f *freezerRemoteS3) initializeBucket() error {
	bucketName := f.bucketName()
	start := time.Now()
	f.log.Info("Creating bucket if not exists", "name", bucketName)
	result, err := f.service.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(f.bucketName()),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeBucketAlreadyExists, s3.ErrCodeBucketAlreadyOwnedByYou:
				f.log.Debug("Bucket exists", "name", bucketName)
				return nil
			}
		}
		return err
	}
	err = f.service.WaitUntilBucketExists(&s3.HeadBucketInput{
		Bucket: aws.String(f.bucketName()),
	})
	if err != nil {
		return err
	}
	f.log.Info("Bucket created", "name", bucketName, "result", result.String(), "elapsed", time.Since(start))
	return nil
}

func (f *freezerRemoteS3) downloadBlocksObject(n uint64) ([]AncientObjectS3, error) {
	key := f.blockObjectKeyForN(n)
	buf := aws.NewWriteAtBuffer([]byte{})
	_, err := f.downloader.Download(buf, &s3.GetObjectInput{
		Bucket: aws.String(f.bucketName()),
		Key:    aws.String(key),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchKey:
				return nil, errOutOfBounds
			}
		}
		f.log.Error("Download error", "method", "pullCache", "error", err, "key", key)
		return nil, err
	}
	target := []AncientObjectS3{}
	err = json.Unmarshal(buf.Bytes(), &target)
	if err != nil {
		return nil, err
	}
	// sanity
	if len(target) > 0 && target[0].Header.Number.Uint64()%f.blockObjectGroupSize != 0 {
		panic(fmt.Sprintf("object does not begin at mod: n=%d", target[0].Header.Number.Uint64()))
	}
	f.retrievedBlockLock.Lock()
	f.retrievedBlocks = map[uint64]AncientObjectS3{}
	for _, v := range target {
		n := v.Header.Number.Uint64()
		f.retrievedBlocks[n] = v
	}
	f.retrievedBlockLock.Unlock()
	return target, nil
}

func (f *freezerRemoteS3) downloadHashesObject(n uint64) ([]common.Hash, error) {
	key := f.hashObjectKeyForN(n)
	buf := aws.NewWriteAtBuffer([]byte{})
	_, err := f.downloader.Download(buf, &s3.GetObjectInput{
		Bucket: aws.String(f.bucketName()),
		Key:    aws.String(key),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchKey:
				return nil, errOutOfBounds
			}
		}
		f.log.Error("Download error", "method", "pullCache", "error", err, "key", key)
		return nil, err
	}
	target := []common.Hash{}
	err = json.Unmarshal(buf.Bytes(), &target)
	if err != nil {
		return nil, err
	}
	return target, nil
}

func (f *freezerRemoteS3) appendCaches(a AncientObjectS3) {
	n := a.Header.Number.Uint64()
	f.cacheLock.Lock()
	if _, ok := f.hashCache[n]; !ok {
		f.hashCache[n] = a.Hash
	}
	if sliceIndexOf(f.blockCacheS, n) < 0 {
		f.blockCacheS = append(f.blockCacheS, n)
		f.blockCache[n] = a
	}
	// This is an insane level of sanity.
	sort.Slice(f.blockCacheS, func(i, j int) bool {
		return f.blockCacheS[i] < f.blockCacheS[j]
	})
	f.cacheLock.Unlock()
}

func (f *freezerRemoteS3) findCached(n uint64, kind string) ([]byte, bool) {
	if kind == freezerHashTable {
		f.cacheLock.Lock()
		v, ok := f.hashCache[n]
		f.cacheLock.Unlock()
		if ok {
			return v.Bytes(), ok
		}
	}
	f.cacheLock.Lock()
	v, ok := f.blockCache[n]
	f.cacheLock.Unlock()
	if ok {
		return v.RLPBytesForKind(kind), ok
	}
	f.retrievedBlockLock.Lock()
	v, ok = f.retrievedBlocks[n]
	f.retrievedBlockLock.Unlock()
	if ok {
		return v.RLPBytesForKind(kind), ok
	}
	return nil, false
}

func (f *freezerRemoteS3) truncateCaches(n uint64) {
	f.cacheLock.Lock()
	hashesS := make([]uint64, len(f.hashCache))
	for k := range f.hashCache {
		hashesS = append(hashesS, k)
	}
	sort.Slice(hashesS, func(i, j int) bool {
		return hashesS[i] < hashesS[j]
	})
	for _, v := range hashesS {
		if v <= n {
			delete(f.hashCache, v)
		} else {
			break
		}
	}
	cachedAt := -1
	for i, v := range f.blockCacheS {
		if v >= n {
			delete(f.blockCache, v)
		}
		if v == n {
			cachedAt = i
		}
	}
	if cachedAt >= 0 {
		f.blockCacheS = f.blockCacheS[:cachedAt]
	}
	f.cacheLock.Unlock()
}

func (f *freezerRemoteS3) spliceBlockCacheLeaving(remainder []uint64) {
	f.cacheLock.Lock()
	// splice first n groups, leaving mod leftovers
	for _, n := range f.blockCacheS {
		if sliceIndexOf(remainder, n) < 0 {
			delete(f.blockCache, n)
		}
	}
	f.blockCacheS = remainder
	f.cacheLock.Unlock()
}

func (f *freezerRemoteS3) pullCache(n uint64) error {
	f.log.Info("Pulling cache", "n", n)

	target, err := f.downloadBlocksObject(n)
	if err != nil {
		return err
	}
	for _, v := range target {
		f.appendCaches(v)
	}
	f.log.Info("Finished pulling cache", "n", n, "size", len(f.blockCache))
	return nil
}

func (f *freezerRemoteS3) pullHashCache(n uint64) error {
	f.log.Info("Pulling hashcache", "n", n)

	target, err := f.downloadHashesObject(n)
	if err != nil {
		return err
	}
	f.cacheLock.Lock()
	start := (n / f.hashObjectGroupSize) * f.hashObjectGroupSize
	for i, v := range target {
		f.hashCache[uint64(i)+start] = v
	}
	f.cacheLock.Unlock()
	f.log.Info("Finished pulling hashcache", "n", n, "size", len(f.hashCache))
	return nil
}

// newFreezer creates a chain freezer that moves ancient chain data into
// append-only flat file containers.
func newFreezerRemoteS3(namespace string, readMeter, writeMeter metrics.Meter, sizeGauge metrics.Gauge) (*freezerRemoteS3, error) {
	var err error

	freezerBlockGroupSize := uint64(32)
	freezerHashGroupSize := uint64(32 * 32 * 32)
	if v := os.Getenv("GETH_FREEZER_S3_BLOCK_GROUP_SIZE"); v != "" {
		i, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return nil, err
		}
		freezerBlockGroupSize = i
	}
	if v := os.Getenv("GETH_FREEZER_S3_HASH_GROUP_SIZE"); v != "" {
		i, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return nil, err
		}
		freezerHashGroupSize = i
	}
	f := &freezerRemoteS3{
		namespace:            namespace,
		quit:                 make(chan struct{}),
		readMeter:            readMeter,
		writeMeter:           writeMeter,
		sizeGauge:            sizeGauge,
		blockObjectGroupSize: freezerBlockGroupSize,
		hashObjectGroupSize:  freezerHashGroupSize,
		appendCacheBlocks:    NewCache(),
		appendCacheHashes:    NewCache(),
		getCacheBlocks:       NewCache(),
		getCacheHashes:       NewCache(),
		retrievedBlocks:      make(map[uint64]AncientObjectS3),
		blockCache:           make(map[uint64]AncientObjectS3),
		blockCacheS:          []uint64{},
		hashCache:            make(map[uint64]common.Hash),
		log:                  log.New("remote", "s3"),
	}

	/*
		By default NewSession will only load credentials from the shared credentials file (~/.aws/credentials).
		If the AWS_SDK_LOAD_CONFIG environment variable is set to a truthy value the Session will be created from the
		configuration values from the shared config (~/.aws/config) and shared credentials (~/.aws/credentials) files.
		Using the NewSessionWithOptions with SharedConfigState set to SharedConfigEnable will create the session as if the
		AWS_SDK_LOAD_CONFIG environment variable was set.
		> https://docs.aws.amazon.com/sdk-for-go/api/aws/session/
	*/
	f.session, err = session.NewSession()
	if err != nil {
		f.log.Info("Session", "err", err)
		return nil, err
	}
	f.log.Info("New session", "region", f.session.Config.Region)
	f.service = s3.New(f.session)

	// Create buckets per the schema, where each bucket is prefixed with the namespace
	// and suffixed with the schema Kind.
	err = f.initializeBucket()
	if err != nil {
		return f, err
	}

	f.uploader = s3manager.NewUploader(f.session)
	f.uploader.Concurrency = 10

	f.downloader = s3manager.NewDownloader(f.session)

	n, _ := f.Ancients()
	f.frozen = &n

	if n > 0 {
		err = f.pullCache(n)
		if err != nil {
			return f, err
		}

		err = f.pullHashCache(n)
		if err != nil {
			return f, err
		}
	}

	return f, nil
}

// Close terminates the chain freezer, unmapping all the data files.
func (f *freezerRemoteS3) Close() error {
	// err := f.pushHashCache()
	// if err != nil {
	// 	return err
	// }
	f.quit <- struct{}{}
	// I don't see any Close, Stop, or Quit methods for the AWS service.
	return nil
}

// HasAncient returns an indicator whether the specified ancient data exists
// in the freezer.
func (f *freezerRemoteS3) HasAncient(kind string, number uint64) (bool, error) {
	v, err := f.Ancient(kind, number)
	if err != nil {
		return false, err
	}

	return v != nil, nil
}

// Ancient retrieves an ancient binary blob from the append-only immutable files.
func (f *freezerRemoteS3) Ancient(kind string, number uint64) ([]byte, error) {
	if atomic.LoadUint64(f.frozen) <= number {
		return nil, nil
	}

	if v, ok := f.findCached(number, kind); ok {
		return v, nil
	}

	if kind == freezerHashTable {
		hashes, err := f.downloadHashesObject(number)
		if err != nil {
			return nil, err
		}
		start := (number / f.hashObjectGroupSize) * f.hashObjectGroupSize
		f.cacheLock.Lock()
		for i, v := range hashes {
			f.hashCache[start+uint64(i)] = v
		}
		v, ok := f.hashCache[number]
		f.cacheLock.Unlock()
		if ok {
			return v.Bytes(), nil
		} else {
			return nil, errOutOfBounds
		}
	}

	// Take from remote
	_, err := f.downloadBlocksObject(number)
	if err != nil {
		return nil, err
	}

	f.retrievedBlockLock.Lock()
	o := f.retrievedBlocks[number]
	f.retrievedBlockLock.Unlock()

	if o.Hash == (common.Hash{}) {
		fmt.Println("number", number, "kind", kind)
		panic("bad")
	}
	return o.RLPBytesForKind(kind), nil
}

// Ancients returns the length of the frozen items.
func (f *freezerRemoteS3) Ancients() (uint64, error) {
	if f.frozen != nil {
		return atomic.LoadUint64(f.frozen), nil
	}
	f.log.Info("Retrieving ancients number")
	result, err := f.service.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(f.bucketName()),
		Key:    aws.String("index-marker"),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchKey:
				return 0, nil
			}
		}
		f.log.Error("GetObject error", "method", "Ancients", "error", err)
		return 0, err
	}
	contents, err := ioutil.ReadAll(result.Body)
	if err != nil {
		return 0, err
	}
	s := strings.TrimSpace(string(contents))
	i, err := strconv.ParseUint(s, 10, 64)
	f.log.Info("Finished retrieving ancients num", "s", s, "n", i, "err?", err)
	return i, err
}

// AncientSize returns the ancient size of the specified category.
func (f *freezerRemoteS3) AncientSize(kind string) (uint64, error) {
	// AWS Go-SDK doesn't support this in a convenient way.
	// This would require listing all objects in the bucket and summing their sizes.
	// This method is only used in the InspectDatabase function, which isn't that
	// important.
	return 0, errNotSupported
}

func (f *freezerRemoteS3) setIndexMarker(number uint64) error {
	f.log.Info("Setting index marker", "number", number)
	numberStr := strconv.FormatUint(number, 10)
	reader := bytes.NewReader([]byte(numberStr))
	_, err := f.service.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(f.bucketName()),
		Key:    aws.String("index-marker"),
		Body:   reader,
	})
	return err
}

// AppendAncient injects all binary blobs belong to block at the end of the
// append-only immutable table files.
//
// Notably, this function is lock free but kind of thread-safe. All out-of-order
// injection will be rejected. But if two injections with same number happen at
// the same time, we can get into the trouble.
func (f *freezerRemoteS3) AppendAncient(number uint64, hash, header, body, receipts, td []byte) (err error) {

	f.mu.Lock()
	defer f.mu.Unlock()

	o, err := NewAncientObjectS3(hash, header, body, receipts, td)
	if err != nil {
		return err
	}
	f.appendCaches(*o)

	atomic.AddUint64(f.frozen, 1)

	return nil
}

// Truncate discards any recent data above the provided threshold number.
// TODO@meowsbits: handle pagination.
//   ListObjects will only (dubiously? might return millions?) return the first 1000. Need to implement pagination.
//   Also make sure that the Marker is working as expected.
func (f *freezerRemoteS3) TruncateAncients(items uint64) error {

	f.mu.Lock()
	defer f.mu.Unlock()

	var err error

	log.Info("Truncating ancients", "cacheS.len", len(f.blockCacheS), "items", items)

	f.cacheLock.Lock()
	_, blockok := f.blockCache[items]
	f.cacheLock.Unlock()
	if !blockok {
		err = f.pullCache(items)
		if err != nil {
			return err
		}
	}
	f.cacheLock.Lock()
	_, hashok := f.hashCache[items]
	f.cacheLock.Unlock()
	if !hashok {
		err = f.pullHashCache(items)
		if err != nil {
			return err
		}
	}

	// Now truncate all remote data from the latest grouped object and beyond.
	// Noting that this can remove data from the remote that is actually antecedent to the
	// desired truncation level, since we have to use the groups.
	// That's why we pulled the object into the cache first; on the next Sync, the un-truncated
	// blocks will be pushed back up to the remote.
	n := atomic.LoadUint64(f.frozen)
	f.log.Info("Truncating ancients", "frozen", n, "target", items, "delta", n-items)
	start := time.Now()

	list := &s3.ListObjectsInput{
		Bucket: aws.String(f.bucketName()),
		Marker: aws.String(f.blockObjectKeyForN(items)),
	}
	iter := s3manager.NewDeleteListIterator(f.service, list)
	batcher := s3manager.NewBatchDeleteWithClient(f.service)
	if err := batcher.Delete(aws.BackgroundContext(), iter); err != nil {
		return err
	}

	if !hashok {
		// We didn't have it in the cache, so we need to explicitly overwrite upstream.
		// Otherwise, it will be automatically handled by the push on the next Sync
		list := &s3.ListObjectsInput{
			Bucket: aws.String(f.bucketName()),
			Marker: aws.String(f.hashObjectKeyForN(items)),
		}
		iter := s3manager.NewDeleteListIterator(f.service, list)
		batcher := s3manager.NewBatchDeleteWithClient(f.service)
		if err := batcher.Delete(aws.BackgroundContext(), iter); err != nil {
			return err
		}
	}

	f.truncateCaches(items)

	// "On the next Sync" should do it.
	// if len(f.cache) > 0 {
	// 	err = f.pushBlockCache()
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	err = f.setIndexMarker(items)
	if err != nil {
		return err
	}
	atomic.StoreUint64(f.frozen, items)

	f.log.Info("Finished truncating ancients", "elapsed", time.Since(start))
	return nil
}

func sliceIndexOf(sl []uint64, n uint64) int {
	for i, s := range sl {
		if s == n {
			return i
		}
	}
	return -1
}

func (f *freezerRemoteS3) pushHashCache() error {
	if len(f.hashCache) == 0 {
		return nil
	}

	set := []common.Hash{}
	uploads := []s3manager.BatchUploadObject{}
	dict := make([]uint64, len(f.hashCache)) // optimize with prealloc
	remainders := []uint64{}

	f.cacheLock.Lock()
	i := 0
	for k := range f.hashCache {
		dict[i] = k
		i++
	}
	f.cacheLock.Unlock()
	sort.Slice(dict, func(i, j int) bool {
		return dict[i] < dict[j]
	})
	for i, n := range dict {
		f.cacheLock.Lock()
		v := f.hashCache[n]
		f.cacheLock.Unlock()
		set = append(set, v)
		remainders = append(remainders, n)

		endGroup := (n+1)%f.hashObjectGroupSize == 0
		if endGroup || i == len(dict)-1 {
			b, err := json.Marshal(set)
			if err != nil {
				return err
			}
			// panic(fmt.Sprintf("hashes: %s, dict: %v (len=%d), hcache: %v (len=%d)", string(b), dict, len(dict), f.hashCache, len(f.hashCache)))
			set = []common.Hash{}
			uploads = append(uploads, s3manager.BatchUploadObject{
				Object: &s3manager.UploadInput{
					Bucket: aws.String(f.bucketName()),
					Key:    aws.String(f.hashObjectKeyForN(n)),
					Body:   bytes.NewReader(b),
				},
			})
		}
		if endGroup {
			remainders = remainders[:0]
		}
	}

	iter := &s3manager.UploadObjectsIterator{Objects: uploads}
	err := f.uploader.UploadWithIterator(aws.BackgroundContext(), iter)
	if err != nil {
		return err
	}

	for _, v := range dict {
		if sliceIndexOf(remainders, v) < 0 {
			f.cacheLock.Lock()
			delete(f.hashCache, v)
			f.cacheLock.Unlock()
		}
	}
	return nil
}

func (f *freezerRemoteS3) pushBlockCache() error {
	if len(f.blockCacheS) == 0 {
		return nil
	}

	if f.blockCacheS[0]%f.blockObjectGroupSize != 0 {
		err := f.pullCache(f.blockCacheS[0])
		if err != nil {
			return err
		}
	}

	set := []AncientObjectS3{}
	uploads := []s3manager.BatchUploadObject{}
	remainders := []uint64{}
	for i, n := range f.blockCacheS {

		f.cacheLock.Lock()
		v := f.blockCache[n]
		f.cacheLock.Unlock()

		set = append(set, v)
		remainders = append(remainders, n)

		// finalize upload object if we have the group-by number in the set, or if the item is the last
		endGroup := (n+1)%f.blockObjectGroupSize == 0
		if endGroup || i == len(f.blockCacheS)-1 {
			// seal upload object
			b, err := json.Marshal(set)
			if err != nil {
				return err
			}
			set = []AncientObjectS3{}
			uploads = append(uploads, s3manager.BatchUploadObject{
				Object: &s3manager.UploadInput{
					Bucket: aws.String(f.bucketName()),
					Key:    aws.String(f.blockObjectKeyForN(n)),
					Body:   bytes.NewReader(b),
				},
			})
		}
		if endGroup {
			remainders = remainders[:0]
		}
	}

	iter := &s3manager.UploadObjectsIterator{Objects: uploads}
	err := f.uploader.UploadWithIterator(aws.BackgroundContext(), iter)
	if err != nil {
		return err
	}

	f.spliceBlockCacheLeaving(remainders)

	f.retrievedBlockLock.Lock()
	f.retrievedBlocks = map[uint64]AncientObjectS3{}
	f.retrievedBlockLock.Unlock()

	return nil
}

// sync flushes all data tables to disk.
func (f *freezerRemoteS3) Sync() error {
	lenCache := len(f.blockCache)
	if lenCache == 0 {
		return nil
	}

	f.mu.Lock()
	defer f.mu.Unlock()

	n := atomic.LoadUint64(f.frozen)

	var err error

	f.log.Info("Syncing ancients", "frozen", n, "blocks", lenCache)
	start := time.Now()

	err = f.pushBlockCache()
	if err != nil {
		return err
	}

	// if uint64(len(f.hashCache)) >= f.hashObjectGroupSize {
	err = f.pushHashCache()
	if err != nil {
		return err
	}
	// }

	elapsed := time.Since(start)
	blocksPerSecond := fmt.Sprintf("%0.2f", float64(lenCache)/elapsed.Seconds())

	err = f.setIndexMarker(atomic.LoadUint64(f.frozen))
	if err != nil {
		return err
	}

	f.log.Info("Finished syncing ancients", "frozen", n, "blocks", lenCache, "elapsed", elapsed, "bps", blocksPerSecond)
	return err
}

// repair truncates all data tables to the same length.
func (f *freezerRemoteS3) repair() error {
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
