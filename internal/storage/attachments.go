package storage

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"enterpriselead/internal/domain"
	"go.etcd.io/bbolt"
)

const MaxAttachmentSize int64 = 5 * 1024 * 1024

func PrepareAttachment(attachment domain.Attachment) (domain.Attachment, error) {
	if attachment.Size < 0 || int64(len(attachment.Content)) != attachment.Size {
		return attachment, fmt.Errorf("attachment size mismatch")
	}
	if attachment.Size > MaxAttachmentSize {
		return attachment, fmt.Errorf("attachment exceeds %d bytes", MaxAttachmentSize)
	}
	hash := sha256.Sum256(attachment.Content)
	attachment.Checksum = hex.EncodeToString(hash[:])
	attachment.Content = cloneBytes(attachment.Content)
	return attachment, nil
}

func (s *DB) PutAttachment(attachment domain.Attachment) error {
	prepared, err := PrepareAttachment(attachment)
	if err != nil {
		return err
	}
	data, err := encode(prepared)
	if err != nil {
		return err
	}
	return s.Update(func(tx *bbolt.Tx) error { return tx.Bucket(bucketAttachments).Put([]byte(prepared.ID), data) })
}

func (s *DB) GetAttachment(id string) (domain.Attachment, error) {
	var attachment domain.Attachment
	err := s.View(func(tx *bbolt.Tx) error {
		value := tx.Bucket(bucketAttachments).Get([]byte(id))
		if value == nil {
			return ErrNotFound
		}
		return decode(value, &attachment)
	})
	attachment.Content = cloneBytes(attachment.Content)
	return attachment, err
}
