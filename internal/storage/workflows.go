package storage

import (
	"enterpriselead/internal/domain"
	"go.etcd.io/bbolt"
)

func (s *DB) PutWorkflow(workflow domain.Workflow) error {
	data, err := encode(workflow)
	if err != nil {
		return err
	}
	return s.Update(func(tx *bbolt.Tx) error { return tx.Bucket(bucketWorkflows).Put([]byte(workflow.ID), data) })
}

func (s *DB) GetWorkflow(id string) (domain.Workflow, error) {
	var workflow domain.Workflow
	err := s.View(func(tx *bbolt.Tx) error {
		value := tx.Bucket(bucketWorkflows).Get([]byte(id))
		if value == nil {
			return ErrNotFound
		}
		return decode(value, &workflow)
	})
	return workflow, err
}

func (s *DB) ListWorkflows(recordID string) ([]domain.Workflow, error) {
	var workflows []domain.Workflow
	err := s.View(func(tx *bbolt.Tx) error {
		return tx.Bucket(bucketWorkflows).ForEach(func(_, value []byte) error {
			var workflow domain.Workflow
			if err := decode(value, &workflow); err != nil {
				return err
			}
			if recordID == "" || workflow.RecordID == recordID {
				workflows = append(workflows, workflow)
			}
			return nil
		})
	})
	return workflows, err
}
