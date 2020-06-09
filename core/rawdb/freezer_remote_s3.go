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
	"fmt"
	"io/ioutil"
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
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/metrics"
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
	pending        map[string][]byte
	pendingMu      sync.Mutex

	pendingCopy map[string][]byte

	frozen *uint64

	log log.Logger
}

func awsKeyRLP(kind string, number uint64) string {
	return fmt.Sprintf("%s/%09d.rlp", kind, number)
}

func fromAwsKey(key string) (kind string, number uint64) {
	k := strings.TrimSuffix(key, ".rlp")
	spl := strings.Split(k, "/")
	kind = spl[0]
	var err error
	number, _ = strconv.ParseUint(spl[1], 10, 64)
	if err != nil {
		panic(err)
	}
	return
}

// TODO: this is superfluous now; bucket names must be user-configured
func (f *freezerRemoteS3) bucketName() string {
	return fmt.Sprintf("%s", f.namespace)
}

func (f *freezerRemoteS3) initializeBucket() error {
	bucketName := f.bucketName()
	start := time.Now()
	f.log.Info("Creating bucket if not exists", "bucket", bucketName)
	result, err := f.service.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(f.bucketName()),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeBucketAlreadyExists, s3.ErrCodeBucketAlreadyOwnedByYou:
				f.log.Debug("Bucket exists", "kind", bucketName)
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
	f.log.Info("Bucket created", "kind", bucketName, "bucket", result.Location, "elapsed", time.Since(start))
	return nil
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
		pending:        make(map[string][]byte),
		pendingCopy:    make(map[string][]byte),
		log:            log.New("remote", "s3"),
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
	key := awsKeyRLP(kind, number)
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

	f.pendingMu.Lock()
	v, ok := f.pending[awsKeyRLP(kind, number)]
	f.pendingMu.Unlock()
	if ok {
		return v, nil
	}

	// Take from remote
	key := awsKeyRLP(kind, number)
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
	return buf.Bytes(), nil
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

	for _, v := range []struct {
		kind string
		val  []byte
	}{
		{freezerHashTable, hash},
		{freezerHeaderTable, header},
		{freezerBodiesTable, body},
		{freezerReceiptTable, receipts},
		{freezerDifficultyTable, td},
	} {
		f.pendingMu.Lock()
		f.pending[awsKeyRLP(v.kind, number)] = v.val
		f.pendingMu.Unlock()
		uploadObj := s3manager.BatchUploadObject{Object: &s3manager.UploadInput{
			Bucket: aws.String(f.bucketName()),
			Key:    aws.String(awsKeyRLP(v.kind, number)),
			Body:   bytes.NewReader(v.val),
		}}
		f.backlogUploads = append(f.backlogUploads, uploadObj)
	}

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

	foundPending := false
	f.pendingMu.Lock()

	for k := range f.pending {
		_, num := fromAwsKey(k)
		// Since items is length of items, we compare == to items as well; numbers at or above
		// this value indicate outside of range since length is 0 indexed.
		if num >= items {
			foundPending = true
			delete(f.pending, k)
		}
	}
	for i, v := range f.backlogUploads {
		_, num := fromAwsKey(*v.Object.Key)
		if num >= items {
			f.backlogUploads = append(f.backlogUploads[:i], f.backlogUploads[i+1:]...)
		}
	}
	// If we found the len in the pending map, that means it hasn't been uploaded yet.
	if foundPending {
		atomic.StoreUint64(f.frozen, items)
		f.pendingMu.Unlock()
		return nil
	}
	f.pendingMu.Unlock()

	f.log.Info("Truncating ancients", "ancients", n, "target", items, "delta", n-items)
	start := time.Now()

	for _, v := range []struct {
		kind string
	}{
		{"hash"},
		{"header"},
		{"body"},
		{"receipts"},
		{"td"},
	} {
		list := &s3.ListObjectsInput{
			Bucket: aws.String(f.bucketName()),
			Marker: aws.String(awsKeyRLP(v.kind, items)),
		}
		iter := s3manager.NewDeleteListIterator(f.service, list)
		batcher := s3manager.NewBatchDeleteWithClient(f.service)
		if err := batcher.Delete(aws.BackgroundContext(), iter); err != nil {
			return err
		}
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

	f.pendingCopy = f.pending
	f.pending = make(map[string][]byte)

	f.log.Info("Syncing ancients", "backlog", lenBacklog)
	start := time.Now()
	uploader := s3manager.NewUploader(f.session)
	uploader.Concurrency = 10
	iter := &s3manager.UploadObjectsIterator{Objects: f.backlogUploads}
	err := uploader.UploadWithIterator(aws.BackgroundContext(), iter)
	if err != nil {
		f.pending = f.pendingCopy
		f.pendingCopy = make(map[string][]byte)
		return err
	}
	f.backlogUploads = []s3manager.BatchUploadObject{}
	elapsed := time.Since(start)
	blocksPerSecond := fmt.Sprintf("%0.2f", float64(lenBacklog)/elapsed.Seconds())
	err = f.setIndexMarker(atomic.LoadUint64(f.frozen))
	if err != nil {
		f.pending = f.pendingCopy
		f.pendingCopy = make(map[string][]byte)
		return err
	}
	f.pendingCopy = make(map[string][]byte)
	f.pending = make(map[string][]byte)

	f.log.Info("Finished syncing ancients", "backlog", lenBacklog, "elapsed", elapsed, "bps", blocksPerSecond)

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
