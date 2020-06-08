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
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/metrics"
	"github.com/ethereum/go-ethereum/params"
)

const (
	awsDefaultRegion  = "us-west-1"
)

var (
	awsMainBucketName = fmt.Sprintf("%s-v%d", params.VersionName, params.VersionMajor)
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

	backlogUploads []s3manager.BatchUploadObject
	frozen         *uint64

	log log.Logger
}

func awsKeyRLPJSON(number uint64) string {
	return fmt.Sprintf("%09d.rlp.json", number)
}

func (f *freezerRemoteS3) bucketName() string {
	return fmt.Sprintf("%s-v%d-%s", params.VersionName, params.VersionMajor, f.namespace)
}

type AncientObjectS3 struct {
	Hash       []byte `json:"hash"`
	Header     []byte `json:"header"`
	Bodies     []byte `json:"bodies"`
	Receipt    []byte `json:"receipt"`
	Difficulty []byte `json:"difficulty"`
}

func (o *AncientObjectS3) GetKind(kind string) []byte {
	switch kind {
	case freezerHashTable:
		return o.Hash
	case freezerHeaderTable:
		return o.Header
	case freezerBodiesTable:
		return o.Bodies
	case freezerReceiptTable:
		return o.Receipt
	case freezerDifficultyTable:
		return o.Difficulty
	default:
		panic("unknown kind")
	}
}

// newFreezer creates a chain freezer that moves ancient chain data into
// append-only flat file containers.
func newFreezerRemoteS3(namespace string, readMeter, writeMeter metrics.Meter, sizeGauge metrics.Gauge) (*freezerRemoteS3, error) {
	var err error

	f := &freezerRemoteS3{
		namespace:      namespace,
		quit:           make(chan struct{}),
		readMeter:      readMeter,
		writeMeter:     writeMeter,
		sizeGauge:      sizeGauge,
		backlogUploads: []s3manager.BatchUploadObject{},
		log:            log.New("remote", "s3"),
	}

	f.log.Info("New session", "region", awsDefaultRegion)
	f.session, err = session.NewSession(&aws.Config{Region: aws.String(awsDefaultRegion)})
	if err != nil {
		f.log.Info("Session", "err", err)
		return nil, err
	}
	f.service = s3.New(f.session)

	// Create buckets per the schema, where each bucket is prefixed with the namespace
	// and suffixed with the schema Kind.
	for _, kind := range []string{
		awsMainBucketName, // All 'table' data is stored together in coerced JSON blobs.
	} {
		start := time.Now()
		f.log.Info("Creating bucket if not exists", "bucket", kind)
		result, err := f.service.CreateBucket(&s3.CreateBucketInput{
			Bucket: aws.String(f.bucketName()),
		})
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case s3.ErrCodeBucketAlreadyExists, s3.ErrCodeBucketAlreadyOwnedByYou:
					f.log.Debug("Bucket exists", "kind", kind)
					continue
				}
			}
			return f, err
		}
		err = f.service.WaitUntilBucketExists(&s3.HeadBucketInput{
			Bucket: aws.String(f.bucketName()),
		})
		if err != nil {
			return f, err
		}
		f.log.Info("Bucket created", "kind", kind, "bucket", result.Location, "elapsed", time.Since(start))
	}

	n, _ := f.Ancients()
	f.frozen = &n

	return f, nil
}

// Close terminates the chain freezer, unmapping all the data files.
func (f *freezerRemoteS3) Close() error {
	f.quit <- struct{}{}
	// I don't see any Close, Stop, or Quit methods for the AWS service.
	return nil
}

// HasAncient returns an indicator whether the specified ancient data exists
// in the freezer.
func (f *freezerRemoteS3) HasAncient(kind string, number uint64) (bool, error) {
	if atomic.LoadUint64(f.frozen) <= number {
		return false, nil
	}
	key := awsKeyRLPJSON(number)
	result, err := f.service.ListObjects(&s3.ListObjectsInput{
		Bucket:  aws.String(f.bucketName()),
		MaxKeys: aws.Int64(1),
		Prefix:  aws.String(key),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchKey:
				return false, nil
			}
		}
		f.log.Error("ListObjects error", "method", "HasAncient", "error", err, "key", key)
		return false, err
	}
	return len(result.Contents) > 0, nil
}

// Ancient retrieves an ancient binary blob from the append-only immutable files.
func (f *freezerRemoteS3) Ancient(kind string, number uint64) ([]byte, error) {
	if atomic.LoadUint64(f.frozen) <= number {
		return nil, nil
	}
	backlogLen := uint64(len(f.backlogUploads))
	if remoteHeight := atomic.LoadUint64(f.frozen) - backlogLen; remoteHeight <= number {
		// Take from backlog
		backlogIndex := number - remoteHeight
		obj := f.backlogUploads[backlogIndex]
		b, err := ioutil.ReadAll(obj.Object.Body)
		obj.Object.Body = bytes.NewReader(b)
		if err != nil {
			return nil, err
		}
		o := &AncientObjectS3{}
		err = json.Unmarshal(b, o)
		if err != nil {
			return nil, err
		}
		return o.GetKind(kind), nil
	}
	// Take from remote
	key := awsKeyRLPJSON(number)
	buf := aws.NewWriteAtBuffer([]byte{})
	downloader := s3manager.NewDownloader(f.session)
	_, err := downloader.Download(buf, &s3.GetObjectInput{
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
		f.log.Error("Download error", "method", "Ancient", "error", err, "kind", kind, "key", key, "number", number)
		return nil, err
	}
	obj := &AncientObjectS3{}
	err = json.Unmarshal(buf.Bytes(), obj)
	if err != nil {
		return nil, err
	}
	return obj.GetKind(kind), nil
}

// Ancients returns the length of the frozen items.
func (f *freezerRemoteS3) Ancients() (uint64, error) {
	if f.frozen != nil {
		return atomic.LoadUint64(f.frozen), nil
	}
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
	return strconv.ParseUint(string(contents), 10, 64)
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

	obj := &AncientObjectS3{
		Hash:       hash,
		Header:     header,
		Bodies:     body,
		Receipt:    receipts,
		Difficulty: td,
	}
	b, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	uploadObj := s3manager.BatchUploadObject{Object: &s3manager.UploadInput{
		Bucket: aws.String(f.bucketName()),
		Key:    aws.String(awsKeyRLPJSON(number)),
		Body:   bytes.NewReader(b),
	}}
	f.backlogUploads = append(f.backlogUploads, uploadObj)
	atomic.AddUint64(f.frozen, 1)

	return nil
}

// Truncate discards any recent data above the provided threshold number.
// TODO@meowsbits: handle pagination.
//   ListObjects will only return the first 1000. Need to implement pagination.
//   Also make sure that the Marker is working as expected.
func (f *freezerRemoteS3) TruncateAncients(items uint64) error {

	f.mu.Lock()
	defer f.mu.Unlock()

	n := atomic.LoadUint64(f.frozen)
	backlogLen := uint64(len(f.backlogUploads))
	if n - backlogLen <= items {
		index := items - (n-backlogLen)
		f.backlogUploads = f.backlogUploads[:index]
		atomic.StoreUint64(f.frozen, items)
		return nil
	}

	f.log.Info("Truncating ancients", "ancients", n, "target", items, "delta", n-items)
	start := time.Now()
	list := &s3.ListObjectsInput{
		Bucket:  aws.String(f.bucketName()),
		Marker: aws.String(awsKeyRLPJSON(items)),
	}
	iter := s3manager.NewDeleteListIterator(f.service, list)
	batcher := s3manager.NewBatchDeleteWithClient(f.service)
	if err := batcher.Delete(aws.BackgroundContext(), iter); err != nil {
		return err
	}

	err := f.setIndexMarker(items)
	if err != nil {
		return err
	}
	atomic.StoreUint64(f.frozen, items)
	f.log.Info("Finished truncating ancients", "elapsed", time.Since(start))
	return nil
}

// sync flushes all data tables to disk.
func (f *freezerRemoteS3) Sync() error {
	lenBacklog := len(f.backlogUploads)
	if lenBacklog == 0 {
		return nil
	}
	f.log.Info("Syncing ancients", "backlog", lenBacklog)
	start := time.Now()
	uploader := s3manager.NewUploader(f.session)
	uploader.Concurrency = 10
	iter := &s3manager.UploadObjectsIterator{Objects: f.backlogUploads}
	err := uploader.UploadWithIterator(aws.BackgroundContext(), iter)
	if err != nil {
		return err
	}
	f.backlogUploads = []s3manager.BatchUploadObject{}
	elapsed := time.Since(start)
	blocksPerSecond := fmt.Sprintf("%0.2f", float64(lenBacklog) / elapsed.Seconds())
	f.log.Info("Finished syncing ancients", "backlog", lenBacklog, "elapsed", elapsed, "bps", blocksPerSecond)
	err = f.setIndexMarker(atomic.LoadUint64(f.frozen))
	if err != nil {
		return err
	}

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
