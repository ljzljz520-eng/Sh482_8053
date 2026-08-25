package service

import (
	"fmt"

	"enterpriselead/internal/audit"
	"enterpriselead/internal/domain"
	"enterpriselead/internal/policy"
	"enterpriselead/internal/search"
)

func (s *Service) Search(query domain.SearchQuery) (domain.SearchResult, error) {
	if err := s.ensureStore(); err != nil {
		return domain.SearchResult{}, err
	}
	records, err := s.store.ListRecords()
	if err != nil {
		return domain.SearchResult{}, err
	}
	return search.Execute(records, search.NormalizeQuery(query)), nil
}

func (s *Service) Update(id string, expectedVersion int, change domain.Change, actor string) (domain.Record, error) {
	if err := s.ensureStore(); err != nil {
		return domain.Record{}, err
	}
	if err := domain.ValidateChange(change); err != nil {
		return domain.Record{}, err
	}
	record, err := s.store.GetRecord(id)
	if err != nil {
		return domain.Record{}, err
	}
	if expectedVersion != record.Version {
		return domain.Record{}, fmt.Errorf("%w: expected %d got %d", storageConflict(), expectedVersion, record.Version)
	}
	if err := authorize(actor, "update", record); err != nil {
		return domain.Record{}, err
	}
	if domain.IsTerminal(record.Status) {
		return domain.Record{}, fmt.Errorf("archived records cannot be changed")
	}
	original := record
	applyChange(&record, change)
	record = domain.NormalizeRecord(record)
	record.Summary = domain.BuildSummary(record)
	record.Version++
	record.UpdatedAt = s.now()
	if err := domain.ValidateRecord(record); err != nil {
		return domain.Record{}, err
	}
	decision := policy.EvaluateChange(original, record)
	if !decision.Allowed {
		return domain.Record{}, fmt.Errorf("update policy: %s", policy.Explain(decision))
	}
	if err := s.store.PutRecord(record); err != nil {
		return domain.Record{}, err
	}
	event := audit.NewEvent(record, "updated", actor, "lead details changed", s.next("evt"), record.UpdatedAt, record.Status)
	if err := audit.Append(s.store, event); err != nil {
		return domain.Record{}, err
	}
	return record, nil
}

func storageConflict() error { return fmt.Errorf("version conflict") }

func applyChange(record *domain.Record, change domain.Change) {
	if change.Company != nil {
		record.Company = *change.Company
	}
	if change.ContactName != nil {
		record.ContactName = *change.ContactName
	}
	if change.ContactEmail != nil {
		record.ContactEmail = *change.ContactEmail
	}
	if change.Source != nil {
		record.Source = *change.Source
	}
	if change.Need != nil {
		record.Need = *change.Need
	}
	if change.Owner != nil {
		record.Owner = *change.Owner
	}
	if change.Priority != nil {
		record.Priority = *change.Priority
	}
	if change.Tags != nil {
		record.Tags = append([]string(nil), change.Tags...)
	}
	if change.Summary != nil {
		record.Summary = *change.Summary
	}
}

func (s *Service) Timeline(id string) ([]domain.AuditEvent, error) {
	if err := s.ensureStore(); err != nil {
		return nil, err
	}
	return audit.Timeline(s.store, id)
}
