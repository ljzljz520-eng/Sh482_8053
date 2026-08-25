package storage

import (
	"bytes"
	"sort"

	"enterpriselead/internal/domain"
	"go.etcd.io/bbolt"
)

func (s *DB) PutRecord(record domain.Record) error {
	data, err := encode(record)
	if err != nil {
		return err
	}
	return s.Update(func(tx *bbolt.Tx) error {
		return tx.Bucket(bucketRecords).Put([]byte(record.ID), data)
	})
}

func (s *DB) CreateRecord(record domain.Record) error {
	data, err := encode(record)
	if err != nil {
		return err
	}
	return s.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(bucketRecords)
		if bucket.Get([]byte(record.ID)) != nil {
			return domainErr("record already exists")
		}
		return bucket.Put([]byte(record.ID), data)
	})
}

func (s *DB) GetRecord(id string) (domain.Record, error) {
	var record domain.Record
	err := s.View(func(tx *bbolt.Tx) error {
		value := tx.Bucket(bucketRecords).Get([]byte(id))
		if value == nil {
			return ErrNotFound
		}
		return decode(value, &record)
	})
	return record, err
}

func (s *DB) ListRecords() ([]domain.Record, error) {
	result := make([]domain.Record, 0)
	err := s.View(func(tx *bbolt.Tx) error {
		return tx.Bucket(bucketRecords).ForEach(func(_, value []byte) error {
			var record domain.Record
			if err := decode(value, &record); err != nil {
				return err
			}
			record.Tags = append([]string(nil), record.Tags...)
			result = append(result, record)
			return nil
		})
	})
	sort.Slice(result, func(i, j int) bool { return bytes.Compare([]byte(result[i].ID), []byte(result[j].ID)) < 0 })
	return result, err
}

func (s *DB) DeleteRecord(id string) error {
	return s.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(bucketRecords)
		if bucket.Get([]byte(id)) == nil {
			return ErrNotFound
		}
		return bucket.Delete([]byte(id))
	})
}

type storageError string

func (e storageError) Error() string { return string(e) }

func domainErr(message string) error { return storageError(message) }
