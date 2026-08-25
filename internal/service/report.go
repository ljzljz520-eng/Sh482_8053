package service

import (
	"enterpriselead/internal/domain"
	"enterpriselead/internal/storage"
)

type Snapshot struct {
	Records  int `json:"records"`
	Open     int `json:"open"`
	InReview int `json:"in_review"`
	Approved int `json:"approved"`
	Archived int `json:"archived"`
	Rejected int `json:"rejected"`
}

func (s *Service) Snapshot() (Snapshot, error) {
	if err := s.ensureStore(); err != nil {
		return Snapshot{}, err
	}
	records, err := s.store.ListRecords()
	if err != nil {
		return Snapshot{}, err
	}
	result := Snapshot{Records: len(records)}
	for _, record := range records {
		switch record.Status {
		case domain.StatusDraft:
			result.Open++
		case domain.StatusReview:
			result.InReview++
		case domain.StatusApproved:
			result.Approved++
		case domain.StatusArchived:
			result.Archived++
		case domain.StatusRejected:
			result.Rejected++
		}
	}
	return result, nil
}

func NewTestStore(path string) (*storage.DB, error) { return storage.Open(path) }
