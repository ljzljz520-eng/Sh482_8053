package storage

import (
	"strings"

	"enterpriselead/internal/domain"
	"go.etcd.io/bbolt"
)

func (s *DB) FindByCompany(company string) ([]domain.Record, error) {
	needle := strings.ToLower(strings.TrimSpace(company))
	result := make([]domain.Record, 0)
	err := s.View(func(tx *bbolt.Tx) error {
		return tx.Bucket(bucketRecords).ForEach(func(_, value []byte) error {
			var record domain.Record
			if err := decode(value, &record); err != nil {
				return err
			}
			if strings.EqualFold(record.Company, needle) {
				result = append(result, record)
			}
			return nil
		})
	})
	return result, err
}

func (s *DB) FindIDs(ids []string) ([]domain.Record, error) {
	result := make([]domain.Record, 0, len(ids))
	err := s.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(bucketRecords)
		for _, id := range ids {
			value := bucket.Get([]byte(id))
			if value == nil {
				continue
			}
			var record domain.Record
			if err := decode(value, &record); err != nil {
				return err
			}
			result = append(result, record)
		}
		return nil
	})
	return result, err
}

func (s *DB) ExportRecords() ([]domain.Record, error) { return s.ListRecords() }
