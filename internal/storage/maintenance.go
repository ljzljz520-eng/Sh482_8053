package storage

import (
	"fmt"
	"time"

	"enterpriselead/internal/domain"
	"go.etcd.io/bbolt"
)

type Health struct {
	Records     int       `json:"records"`
	Events      int       `json:"events"`
	Workflows   int       `json:"workflows"`
	Attachments int       `json:"attachments"`
	CheckedAt   time.Time `json:"checked_at"`
}

func (s *DB) Health(at time.Time) (Health, error) {
	health := Health{CheckedAt: at.UTC()}
	err := s.View(func(tx *bbolt.Tx) error {
		var err error
		health.Records, err = countBucket(tx.Bucket(bucketRecords))
		if err != nil {
			return err
		}
		health.Events, err = countBucket(tx.Bucket(bucketEvents))
		if err != nil {
			return err
		}
		health.Workflows, err = countBucket(tx.Bucket(bucketWorkflows))
		if err != nil {
			return err
		}
		health.Attachments, err = countBucket(tx.Bucket(bucketAttachments))
		return err
	})
	return health, err
}

func countBucket(bucket *bbolt.Bucket) (int, error) {
	if bucket == nil {
		return 0, fmt.Errorf("bucket missing")
	}
	count := 0
	err := bucket.ForEach(func(_, _ []byte) error { count++; return nil })
	return count, err
}

func (s *DB) ReplaceRecord(record domain.Record, expectedVersion int) error {
	data, err := encode(record)
	if err != nil {
		return err
	}
	return s.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(bucketRecords)
		current := bucket.Get([]byte(record.ID))
		if current == nil {
			return ErrNotFound
		}
		var existing domain.Record
		if err := decode(current, &existing); err != nil {
			return err
		}
		if existing.Version != expectedVersion {
			return ErrConflict
		}
		return bucket.Put([]byte(record.ID), data)
	})
}

func (s *DB) CountByStatus(status domain.LeadStatus) (int, error) {
	count := 0
	err := s.View(func(tx *bbolt.Tx) error {
		return tx.Bucket(bucketRecords).ForEach(func(_, value []byte) error {
			var record domain.Record
			if err := decode(value, &record); err != nil {
				return err
			}
			if record.Status == status {
				count++
			}
			return nil
		})
	})
	return count, err
}
