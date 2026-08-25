package service

import (
	"fmt"

	"enterpriselead/internal/domain"
)

type BatchChange struct {
	IDs      []string       `json:"ids"`
	Expected map[string]int `json:"expected"`
	Change   domain.Change  `json:"change"`
	Actor    string         `json:"actor"`
}

type BatchResult struct {
	Updated []domain.Record `json:"updated"`
	Failed  []BatchFailure  `json:"failed"`
}

type BatchFailure struct {
	ID    string `json:"id"`
	Error string `json:"error"`
}

func (s *Service) BatchUpdate(batch BatchChange) BatchResult {
	result := BatchResult{Updated: make([]domain.Record, 0, len(batch.IDs)), Failed: make([]BatchFailure, 0)}
	seen := make(map[string]bool, len(batch.IDs))
	for _, id := range batch.IDs {
		if seen[id] {
			result.Failed = append(result.Failed, BatchFailure{ID: id, Error: "duplicate id"})
			continue
		}
		seen[id] = true
		expected := batch.Expected[id]
		record, err := s.Update(id, expected, batch.Change, batch.Actor)
		if err != nil {
			result.Failed = append(result.Failed, BatchFailure{ID: id, Error: err.Error()})
			continue
		}
		result.Updated = append(result.Updated, record)
	}
	return result
}

func (s *Service) Reassign(ids []string, owner, actor string) BatchResult {
	if owner == "" {
		return BatchResult{Failed: []BatchFailure{{ID: "", Error: "owner is required"}}}
	}
	results := make(map[string]int, len(ids))
	for _, id := range ids {
		record, err := s.Get(id)
		if err == nil {
			results[id] = record.Version
		}
	}
	return s.BatchUpdate(BatchChange{IDs: ids, Expected: results, Change: domain.Change{Owner: &owner}, Actor: actor})
}

func ValidateBatch(batch BatchChange) error {
	if len(batch.IDs) == 0 {
		return fmt.Errorf("batch is empty")
	}
	if len(batch.IDs) > 100 {
		return fmt.Errorf("batch exceeds 100 records")
	}
	if batch.Actor == "" {
		return fmt.Errorf("batch actor is required")
	}
	return domain.ValidateChange(batch.Change)
}
