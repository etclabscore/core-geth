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
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/metrics"
)

const (
	awsMetaBucketName = "meta"
	awsDefaultRegion  = "us-west-1"
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

	log log.Logger
}

func awsKeyRLP(number uint64) string {
	return fmt.Sprintf("%d.rlp", number)
}

func (f *freezerRemoteS3) bucketName(kind string) string {
	return fmt.Sprintf("%s-%s", f.namespace, kind)
}

// newFreezer creates a chain freezer that moves ancient chain data into
// append-only flat file containers.
func newFreezerRemoteS3(namespace string, readMeter, writeMeter metrics.Meter, sizeGauge metrics.Gauge) (*freezerRemoteS3, error) {
	var err error

	f := &freezerRemoteS3{
		namespace: namespace,
		quit:      make(chan struct{}),
		readMeter: readMeter,
		writeMeter: writeMeter,
		sizeGauge: sizeGauge,
		log:       log.New("remote", "s3"),
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
		awsMetaBucketName,
		freezerHashTable,
		freezerHeaderTable,
		freezerBodiesTable,
		freezerReceiptTable,
		freezerDifficultyTable,
	} {
		start := time.Now()
		f.log.Info("Creating bucket if not exists", "bucket", kind)
		result, err := f.service.CreateBucket(&s3.CreateBucketInput{
			Bucket: aws.String(f.bucketName(kind)),
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
			Bucket: aws.String(f.bucketName(kind)),
		})
		if err != nil {
			return f, err
		}
		f.log.Info("Bucket created", "kind", kind, "bucket", result.Location, "elapsed", time.Since(start))
	}

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
	key := awsKeyRLP(number)
	result, err := f.service.ListObjects(&s3.ListObjectsInput{
		Bucket:  aws.String(f.bucketName(kind)),
		MaxKeys: aws.Int64(1),
		Prefix:  aws.String(key),
	})
	if err != nil {
		f.log.Error("ListObjects error", "method", "HasAncient", "error", err, "key", key)
		return false, nil
	}
	return len(result.Contents) > 0, nil
}

// Ancient retrieves an ancient binary blob from the append-only immutable files.
func (f *freezerRemoteS3) Ancient(kind string, number uint64) ([]byte, error) {
	key := awsKeyRLP(number)
	buf := aws.NewWriteAtBuffer([]byte{})
	downloader := s3manager.NewDownloader(f.session)
	_, err := downloader.Download(buf, &s3.GetObjectInput{
		Bucket: aws.String(f.bucketName(kind)),
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
	result, err := f.service.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(f.bucketName(awsMetaBucketName)),
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
	numberStr := strconv.FormatUint(number, 10)
	reader := bytes.NewReader([]byte(numberStr))
	_, err := f.service.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(f.bucketName(awsMetaBucketName)),
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

	puts := map[string][]byte{
		freezerHashTable:       hash,
		freezerHeaderTable:     header,
		freezerBodiesTable:     body,
		freezerReceiptTable:    receipts,
		freezerDifficultyTable: td,
	}
	for k, v := range puts {
		reader := bytes.NewReader(v)
		_, err := f.service.PutObject(&s3.PutObjectInput{
			Bucket: aws.String(f.bucketName(k)),
			Key:    aws.String(awsKeyRLP(number)),
			Body:   reader,
		})
		if err != nil {
			f.log.Error("PutObject error", "method", "AppendAncient", "error", err)
			return err
		}

		f.writeMeter.Mark(int64(len(v)))
		f.sizeGauge.Inc(int64(len(v)))
	}

	err = f.setIndexMarker(number)
	if err != nil {
		f.log.Error("Append ancient", "key", "index-marker", "number", number, "error", err)
		return err
	}
	f.log.Info("Append ancient", "number", number)
	return nil
}

// Truncate discards any recent data above the provided threshold number.
// TODO@meowsbits: handle pagination.
//   ListObjects will only return the first 1000. Need to implement pagination.
//   Also make sure that the Marker is working as expected.
func (f *freezerRemoteS3) TruncateAncients(items uint64) error {

	f.mu.Lock()
	defer f.mu.Unlock()

	tables := []string{
		freezerHashTable,
		freezerHeaderTable,
		freezerBodiesTable,
		freezerReceiptTable,
		freezerDifficultyTable,
	}
	for _, v := range tables {
		result, err := f.service.ListObjects(&s3.ListObjectsInput{
			Bucket: aws.String(f.bucketName(v)),
			Marker: aws.String(awsKeyRLP(items)),
		})
		if err != nil {
			return err
		}
		for _, c := range result.Contents {
			_, err := f.service.DeleteObject(&s3.DeleteObjectInput{
				Bucket: aws.String(f.bucketName(v)),
				Key:    c.Key,
			})
			if err != nil {
				return err
			}
		}
	}
	return f.setIndexMarker(items)
}

// sync flushes all data tables to disk.
func (f *freezerRemoteS3) Sync() error {
	// TODO: Noop for now.
	//   We'll more than likely want to implement a caching strategy, but
	//   since we're writing+deleting directly on the AppendAncient+Truncate calls,
	//   there's nothing to do here.
	return nil
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
