package storage

import (
	"fmt"

	"enterpriselead/internal/domain"
	"go.etcd.io/bbolt"
)

func (s *DB) PutEvent(event domain.AuditEvent) error {
	data, err := encode(event)
	if err != nil {
		return err
	}
	key := []byte(fmt.Sprintf("%s:%s", event.RecordID, event.ID))
	return s.Update(func(tx *bbolt.Tx) error { return tx.Bucket(bucketEvents).Put(key, data) })
}

func (s *DB) ListEvents(recordID string) ([]domain.AuditEvent, error) {
	var events []domain.AuditEvent
	prefix := []byte(recordID + ":")
	err := s.View(func(tx *bbolt.Tx) error {
		cursor := tx.Bucket(bucketEvents).Cursor()
		for key, value := cursor.Seek(prefix); key != nil && len(key) >= len(prefix) && string(key[:len(prefix)]) == string(prefix); key, value = cursor.Next() {
			var event domain.AuditEvent
			if err := decode(value, &event); err != nil {
				return err
			}
			events = append(events, event)
		}
		return nil
	})
	return events, err
}
