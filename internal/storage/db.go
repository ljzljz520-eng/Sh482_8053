package storage

import (
	"errors"
	"fmt"
	"path/filepath"
	"time"

	"go.etcd.io/bbolt"
)

var (
	ErrNotFound = errors.New("record not found")
	ErrConflict = errors.New("version conflict")
)

var bucketRecords = []byte("records")
var bucketEvents = []byte("events")
var bucketWorkflows = []byte("workflows")
var bucketAttachments = []byte("attachments")

type DB struct {
	db *bbolt.DB
}

func Open(path string) (*DB, error) {
	if filepath.Clean(path) == "." || path == "" {
		return nil, fmt.Errorf("database path is required")
	}
	db, err := bbolt.Open(path, 0600, &bbolt.Options{Timeout: 2 * time.Second, NoSync: true})
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}
	store := &DB{db: db}
	if err := store.ensureBuckets(); err != nil {
		db.Close()
		return nil, err
	}
	return store, nil
}

func (s *DB) ensureBuckets() error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		for _, name := range [][]byte{bucketRecords, bucketEvents, bucketWorkflows, bucketAttachments} {
			if _, err := tx.CreateBucketIfNotExists(name); err != nil {
				return fmt.Errorf("create bucket %s: %w", name, err)
			}
		}
		return nil
	})
}

func (s *DB) Close() error {
	if s == nil || s.db == nil {
		return nil
	}
	return s.db.Close()
}

func (s *DB) Sync() error {
	if s == nil || s.db == nil {
		return fmt.Errorf("database is closed")
	}
	return s.db.Sync()
}

func (s *DB) View(fn func(*bbolt.Tx) error) error {
	if s == nil || s.db == nil {
		return fmt.Errorf("database is closed")
	}
	return s.db.View(fn)
}

func (s *DB) Update(fn func(*bbolt.Tx) error) error {
	if s == nil || s.db == nil {
		return fmt.Errorf("database is closed")
	}
	return s.db.Update(fn)
}
