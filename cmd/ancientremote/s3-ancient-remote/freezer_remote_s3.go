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

package main

import (
	"bytes"
	"compress/gzip"
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

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/metrics"
	"github.com/ethereum/go-ethereum/rlp"
	lru "github.com/hashicorp/golang-lru"
)

const (
	s3EncodeJSON   = ".json"
	s3EncodeJSONGZ = ".json.gz"
)

var (
	s3ROnly              = false
	s3BlocksGroupSize    = uint64(32 * 32)
	s3HashesGroupSize    = uint64(32 * 32 * 32)
	s3Encoding           = s3EncodeJSONGZ
	errOutOfBounds       = errors.New("out of bounds")
	errNotSupported      = errors.New("this operation is not supported")
	errOutOrderInsertion = errors.New("the append operation is out-order")
)

func init() {
	if v := os.Getenv("GETH_FREEZER_S3_BLOCK_GROUP_SIZE"); v != "" {
		i, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			panic(err)
		}
		s3BlocksGroupSize = i
	}
	if v := os.Getenv("GETH_FREEZER_S3_HASH_GROUP_SIZE"); v != "" {
		i, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			panic(err)
		}
		s3HashesGroupSize = i
	}
	if v := os.Getenv("GETH_FREEZER_S3_ENCODING"); v != "" {
		s3Encoding = v
	}
	if v := os.Getenv("GETH_FREEZER_S3_R_ONLY"); v != "" {
		s3ROnly = true
	}
}

type freezerRemoteS3 struct {
	session *session.Session
	service *s3.S3

	namespace string
	quit      chan struct{}
	mu        sync.Mutex

	readOnly bool

	readMeter  metrics.Meter // Meter for measuring the effective amount of data read
	writeMeter metrics.Meter // Meter for measuring the effective amount of data written
	sizeGauge  metrics.Gauge // Gauge for tracking the combined size of all freezer tables

	uploader   *s3manager.Uploader
	downloader *s3manager.Downloader

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

// AncientObjectS3 describes the storage encoding unit for a 'block'.
// These objects are grouped in an array for storage with size determined by s3BlocksGroupSize.
type AncientObjectS3 struct {
	Hash       common.Hash                `json:"hash"`
	Header     *types.Header              `json:"header"`
	Body       *types.Body                `json:"body"`
	Receipts   []*types.ReceiptForStorage `json:"receipts"`
	Difficulty *big.Int                   `json:"difficulty"`
}

// NewAncientObjectS3 reverses decodes the incoming RLP bytes, applying the values
// to its own type definition's fields.
func NewAncientObjectS3(hashB, headerB, bodyB, receiptsB, difficultyB []byte) (AncientObjectS3, error) {
	var err error

	hash := common.BytesToHash(hashB)

	header := &types.Header{}
	err = rlp.DecodeBytes(headerB, header)
	if err != nil {
		return AncientObjectS3{}, err
	}
	body := &types.Body{}
	err = rlp.DecodeBytes(bodyB, body)
	if err != nil {
		return AncientObjectS3{}, err
	}
	receipts := []*types.ReceiptForStorage{}
	err = rlp.DecodeBytes(receiptsB, &receipts)
	if err != nil {
		return AncientObjectS3{}, err
	}
	difficulty := new(big.Int)
	err = rlp.DecodeBytes(difficultyB, difficulty)
	if err != nil {
		return AncientObjectS3{}, err
	}
	return AncientObjectS3{
		Hash:       hash,
		Header:     header,
		Body:       body,
		Receipts:   receipts,
		Difficulty: difficulty,
	}, nil
}

// RLPBytesForKind return the RLP encoded bytes expected by the AncientStore interface
// for a given 'kind' on the block object.
func (o AncientObjectS3) RLPBytesForKind(kind string) []byte {
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

func awsKeyBlock(number uint64) string {
	// Keep blocks in a dir.
	// This namespaces the resource, separating it from the 'index-marker' object.
	return fmt.Sprintf("blocks/%09d%s", number, s3Encoding)
}

func awsKeyHash(number uint64) string {
	return fmt.Sprintf("hashes/%09d%s", number, s3Encoding)
}

func (f *freezerRemoteS3) blockObjectKeyForN(n uint64) string {
	return awsKeyBlock((n / f.blockObjectGroupSize) * f.blockObjectGroupSize) // 0, 32, 64, 96, ...
}

func (f *freezerRemoteS3) hashObjectKeyForN(n uint64) string {
	return awsKeyHash((n / f.hashObjectGroupSize) * f.hashObjectGroupSize)
}

// TODO: this is superfluous now; bucket names must be user-configured
func (f *freezerRemoteS3) bucketName() string {
	return f.namespace
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
	f.log.Info("Bucket created", "name", bucketName, "result", result.String(), "elapsed", common.PrettyDuration(time.Since(start)))
	return nil
}

func (f *freezerRemoteS3) downloadBlocksObject(n uint64) error {
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
				return errOutOfBounds
			}
		}
		f.log.Error("Download error", "method", "pullWCacheBlocks", "error", err, "key", key)
		return err
	}
	target := []AncientObjectS3{}
	err = f.decodeObject(buf.Bytes(), &target)
	if err != nil {
		return err
	}
	// sanity
	if len(target) > 0 {
		first := target[0].Header.Number.Uint64()
		if first%f.blockObjectGroupSize != 0 {
			panic(fmt.Sprintf("object does not begin at mod: n=%d", target[0].Header.Number.Uint64()))
		}
		f.rCacheBlocks.Purge()
		for i, v := range target {
			f.rCacheBlocks.Add(first+uint64(i), v)
		}
	}
	return nil
}

func (f *freezerRemoteS3) downloadHashesObject(n uint64) error {
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
				return errOutOfBounds
			}
		}
		f.log.Error("Download error", "method", "pullWCacheBlocks", "error", err, "key", key)
		return err
	}
	target := []common.Hash{}
	err = f.decodeObject(buf.Bytes(), &target)
	if err != nil {
		return err
	}
	if len(target) > 0 {
		f.rCacheHashes.Purge()
		first := (n / f.hashObjectGroupSize) * f.hashObjectGroupSize
		for i, v := range target {
			f.rCacheHashes.Add(first+uint64(i), v)
		}
	}
	return nil
}

func (f *freezerRemoteS3) pullWCacheBlocks(n uint64) error {
	f.log.Info("Pulling write blocks cache", "n", n)
	err := f.downloadBlocksObject(n)
	if err != nil {
		return err
	}
	for k := range f.rCacheBlocks.Keys() {
		v, _ := f.rCacheBlocks.Get(k)
		f.wCacheBlocks.Add(k, v)
	}
	f.log.Info("Finished pulling blocks cache", "n", n, "size", f.wCacheBlocks.Len())
	return nil
}

func (f *freezerRemoteS3) pullWCacheHashes(n uint64) error {
	f.log.Info("Pulling write hashes cache", "n", n)
	err := f.downloadHashesObject(n)
	if err != nil {
		return err
	}
	for k := range f.rCacheHashes.Keys() {
		v, _ := f.rCacheHashes.Get(k)
		f.wCacheHashes.Add(k, v)
	}
	f.log.Info("Finished pulling hashes cache", "n", n, "size", f.wCacheHashes.Len())
	return nil
}

func (f *freezerRemoteS3) findCached(n uint64, kind string) ([]byte, bool) {
	if kind == rawdb.FreezerRemoteHashTable {
		if v, ok := f.wCacheHashes.Get(n); ok {
			return v.(common.Hash).Bytes(), ok
		}
		if v, ok := f.rCacheHashes.Get(n); ok {
			return v.(common.Hash).Bytes(), ok
		}
	}
	if v, ok := f.wCacheBlocks.Get(n); ok {
		return v.(AncientObjectS3).RLPBytesForKind(kind), ok
	}
	if v, ok := f.rCacheBlocks.Get(n); ok {
		return v.(AncientObjectS3).RLPBytesForKind(kind), ok
	}
	return nil, false
}

// newFreezer creates a chain freezer that moves ancient chain data into
// append-only flat file containers.
func newFreezerRemoteS3(namespace string, readMeter, writeMeter metrics.Meter, sizeGauge metrics.Gauge) (*freezerRemoteS3, error) {
	var err error

	rBlockCache, err := lru.New(int(s3BlocksGroupSize) / 2)
	if err != nil {
		return nil, err
	}
	rHashCache, err := lru.New(int(s3HashesGroupSize) / 2)
	if err != nil {
		return nil, err
	}
	wBlockCache, err := lru.New(int(s3BlocksGroupSize) / 2)
	if err != nil {
		return nil, err
	}
	wHashCache, err := lru.New(int(s3HashesGroupSize) / 2)
	if err != nil {
		return nil, err
	}

	f := &freezerRemoteS3{
		namespace:  namespace,
		quit:       make(chan struct{}),
		readMeter:  readMeter,
		writeMeter: writeMeter,
		sizeGauge:  sizeGauge,

		readOnly: s3ROnly,

		// Globals for now. Should probably become CLI flags.
		// Maybe Remote Freezers need a config struct.
		blockObjectGroupSize: s3BlocksGroupSize,
		hashObjectGroupSize:  s3HashesGroupSize,

		wCacheBlocks: wBlockCache,
		wCacheHashes: wHashCache,
		rCacheBlocks: rBlockCache,
		rCacheHashes: rHashCache,
		log:          log.New("remote", "s3"),
	}

	switch s3Encoding {
	case s3EncodeJSONGZ:
		f.encoding = s3EncodeJSONGZ
	case s3EncodeJSON:
		f.encoding = s3EncodeJSON
	default:
		return nil, fmt.Errorf("unknown encoding: %s", s3Encoding)
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
	f.log.Info("New session", "region", *f.session.Config.Region)
	f.service = s3.New(f.session)

	// Create buckets per the schema, where each bucket is prefixed with the namespace
	// and suffixed with the schema Kind.
	err = f.initializeBucket()
	if err != nil {
		return f, err
	}

	f.uploader = s3manager.NewUploader(f.session)
	f.downloader = s3manager.NewDownloader(f.session)

	n, _ := f.Ancients()
	f.frozen = &n

	if n > 0 {
		err = f.pullWCacheBlocks(n)
		if err != nil {
			return f, err
		}

		err = f.pullWCacheHashes(n)
		if err != nil {
			return f, err
		}
	}

	return f, nil
}

// Close terminates the chain freezer, unmapping all the data files.
func (f *freezerRemoteS3) Close() error {
	// I don't see any Close, Stop, or Quit methods for the AWS service.
	f.quit <- struct{}{}
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

	if kind == rawdb.FreezerRemoteHashTable {
		err := f.downloadHashesObject(number)
		if err != nil {
			return nil, err
		}
		if v, ok := f.rCacheHashes.Get(number); ok {
			return v.(common.Hash).Bytes(), nil
		}
		return nil, fmt.Errorf("%w: #%d (%s)", errOutOfBounds, number, kind)
	}

	// Take from remote
	err := f.downloadBlocksObject(number)
	if err != nil {
		return nil, err
	}
	if v, ok := f.rCacheBlocks.Get(number); ok {
		return v.(AncientObjectS3).RLPBytesForKind(kind), nil
	}
	return nil, fmt.Errorf("%w: #%d (%s)", errOutOfBounds, number, kind)
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
	f.log.Info("Finished retrieving ancients num", "s", s, "n", i, "err", err)
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
	// Ensure the binary blobs we are appending is continuous with freezer.
	if atomic.LoadUint64(f.frozen) != number {
		return errOutOrderInsertion
	}
	// f.log.Info("Appending ancient", "frozen", atomic.LoadUint64(f.frozen), "number", number)
	f.mu.Lock()
	defer f.mu.Unlock()
	o, err := NewAncientObjectS3(hash, header, body, receipts, td)
	if err != nil {
		return err
	}
	f.wCacheHashes.Add(number, common.BytesToHash(hash))
	f.wCacheBlocks.Add(number, o)
	atomic.AddUint64(f.frozen, 1)
	// f.log.Info("Finished appending ancient", "frozen", atomic.LoadUint64(f.frozen), "number", number)
	return nil
}

// Truncate discards any recent data above the provided threshold number.
// TODO@meowsbits: handle pagination.
//   ListObjects will only (dubiously? might return millions?) return the first 1000. Need to implement pagination.
//   Also make sure that the Marker is working as expected.
func (f *freezerRemoteS3) TruncateAncients(items uint64) error {

	f.mu.Lock()
	defer f.mu.Unlock()

	n := atomic.LoadUint64(f.frozen)
	f.log.Info("Truncating ancients", "frozen", n, "target", items, "delta", n-items)
	start := time.Now()

	// How this works:
	// 0. Push the new index marker.
	// All data above the truncation limit is allowed to persist, but will eventually be overwritten.

	var err error

	if !f.readOnly {
		err = f.setIndexMarker(items)
		if err != nil {
			return err
		}
	}

	f.rCacheBlocks.Purge()
	f.rCacheHashes.Purge()
	f.wCacheBlocks.Purge()
	f.wCacheHashes.Purge()

	atomic.StoreUint64(f.frozen, items)
	f.log.Info("Finished truncating ancients", "elapsed", common.PrettyDuration(time.Since(start)))
	return nil
}

func (f *freezerRemoteS3) pushWCaches() error {
	var err error
	err = f.pushCacheBatch(f.wCacheBlocks, f.blockObjectGroupSize, f.blockObjectKeyForN, f.blockCacheBatchObjectFn)
	if err != nil {
		return err
	}
	err = f.pushCacheBatch(f.wCacheHashes, f.hashObjectGroupSize, f.hashObjectKeyForN, f.hashCacheBatchObjectFn)
	if err != nil {
		return err
	}
	return nil
}

func firstN(cache *lru.Cache) uint64 {
	keys := cache.Keys()
	sort.Slice(keys, func(i, j int) bool {
		return keys[i].(uint64) < keys[j].(uint64)
	})
	return keys[0].(uint64)
}

func (f *freezerRemoteS3) pushCacheBatch(cache *lru.Cache, size uint64, keyFn func(uint64) string, batchObjFn func([]interface{}) interface{}) error {
	if cache.Len() == 0 {
		return nil
	}
	uploads := []s3manager.BatchUploadObject{}
	for {
		if cache.Len() == 0 {
			break
		}
		keys := cache.Keys()
		var n = firstN(cache)
		if n%size != 0 {
			panic(fmt.Sprintf("bad mod: n=%d r=%d mod=%d len=%d", n, n%size, size, cache.Len()))
		}
		batch := []interface{}{}
		batchKeys := []interface{}{}
		for i := uint64(0); i < size; i++ {
			k := keys[i]
			v, ok := cache.Get(k)
			if !ok {
				break
			}
			batch = append(batch, v)
			batchKeys = append(batchKeys, k)
		}
		b, err := f.encodeObject(batchObjFn(batch))
		if err != nil {
			return err
		}
		uploads = append(uploads, s3manager.BatchUploadObject{
			Object: &s3manager.UploadInput{
				Bucket: aws.String(f.bucketName()),
				Key:    aws.String(keyFn(n)),
				Body:   bytes.NewReader(b),
			},
		})
		batchLen := uint64(len(batch))
		if batchLen%size == 0 {
			for k := range batchKeys {
				cache.Remove(k)
			}
		} else {
			break
		}
	}
	iter := &s3manager.UploadObjectsIterator{Objects: uploads}
	err := f.uploader.UploadWithIterator(aws.BackgroundContext(), iter)
	if err != nil {
		return err
	}
	return nil
}

func (f *freezerRemoteS3) encodeObject(any interface{}) ([]byte, error) {
	b, err := json.Marshal(any)
	if err != nil {
		return nil, err
	}
	if f.encoding == s3EncodeJSONGZ {
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

func (f *freezerRemoteS3) decodeObject(input []byte, target interface{}) error {
	if f.encoding == s3EncodeJSONGZ {
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

func (f *freezerRemoteS3) blockCacheBatchObjectFn(items []interface{}) interface{} {
	batchSet := make([]AncientObjectS3, len(items))
	for i, v := range items {
		batchSet[i] = v.(AncientObjectS3)
	}
	return batchSet
}

func (f *freezerRemoteS3) hashCacheBatchObjectFn(items []interface{}) interface{} {
	batchSet := make([]common.Hash, len(items))
	for i, v := range items {
		batchSet[i] = v.(common.Hash)
	}
	return batchSet
}

// sync flushes all data tables to disk.
func (f *freezerRemoteS3) Sync() error {
	f.mu.Lock()

	n := atomic.LoadUint64(f.frozen)
	lenBlocks := f.wCacheBlocks.Len()
	f.log.Info("Syncing ancients", "frozen", n, "blocks", lenBlocks)
	start := time.Now()

	var err error

	if !f.readOnly {
		if lenBlocks > 0 {
			wCacheBlocksFirstN := firstN(f.wCacheBlocks)
			if r := wCacheBlocksFirstN % f.blockObjectGroupSize; r != 0 {
				f.log.Warn("Found out-of-order block cache", "n", wCacheBlocksFirstN)
				err = f.pullWCacheBlocks(wCacheBlocksFirstN)
				if err != nil {
					f.mu.Unlock()
					return err
				}
				f.mu.Unlock()
				return f.Sync()
			}
			wCacheHashesFirstN := firstN(f.wCacheHashes)
			if r := wCacheHashesFirstN % f.hashObjectGroupSize; r != 0 {
				f.log.Warn("Found out-of-order hash cache", "n", wCacheHashesFirstN)
				err = f.pullWCacheHashes(wCacheHashesFirstN)
				if err != nil {
					f.mu.Unlock()
					return err
				}
				f.mu.Unlock()
				return f.Sync()
			}
		}

		err = f.pushWCaches()
		if err != nil {
			f.mu.Unlock()
			return err
		}

		err = f.setIndexMarker(atomic.LoadUint64(f.frozen))
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
