package service

import (
	"fmt"
	"time"

	"enterpriselead/internal/audit"
	"enterpriselead/internal/domain"
	"enterpriselead/internal/policy"
	"enterpriselead/internal/storage"
	"enterpriselead/internal/workflow"
)

type TimeSource interface{ Now() time.Time }

type FixedTime struct{ Value time.Time }

func (t FixedTime) Now() time.Time { return t.Value }

type Service struct {
	store *storage.DB
	clock TimeSource
	ids   *workflow.Sequence
}

func New(store *storage.DB, clock TimeSource) *Service {
	if clock == nil {
		clock = FixedTime{Value: time.Unix(0, 0).UTC()}
	}
	return &Service{store: store, clock: clock, ids: &workflow.Sequence{}}
}

func (s *Service) next(prefix string) string { return s.ids.Next(prefix) }

func (s *Service) now() time.Time { return s.clock.Now().UTC() }

func (s *Service) ensureStore() error {
	if s == nil || s.store == nil {
		return fmt.Errorf("service store is unavailable")
	}
	return nil
}

func authorize(actor, action string, record domain.Record) error {
	role := policy.ParseRole(actor)
	if !policy.CanPerform(role, action, record) {
		return fmt.Errorf("actor %q cannot perform %s", actor, action)
	}
	return nil
}

type CreateInput struct {
	Company      string
	ContactName  string
	ContactEmail string
	Source       string
	Need         string
	Owner        string
	Priority     domain.Priority
	Tags         []string
}

func (s *Service) Create(input CreateInput, actor string) (domain.Record, error) {
	if err := s.ensureStore(); err != nil {
		return domain.Record{}, err
	}
	now := s.now()
	record := domain.Record{ID: s.next("rec"), Company: input.Company, ContactName: input.ContactName, ContactEmail: input.ContactEmail, Source: input.Source, Need: input.Need, Owner: input.Owner, Priority: input.Priority, Status: domain.StatusDraft, Tags: append([]string(nil), input.Tags...), CreatedAt: now, UpdatedAt: now, Version: 1}
	record = domain.ApplyDefaults(domain.NormalizeRecord(record))
	if err := authorize(actor, "create", record); err != nil {
		return domain.Record{}, err
	}
	record.Summary = domain.BuildSummary(record)
	if err := domain.ValidateRecord(record); err != nil {
		return domain.Record{}, err
	}
	decision := policy.EvaluateCreate(record)
	if !decision.Allowed {
		return domain.Record{}, fmt.Errorf("create policy: %s", policy.Explain(decision))
	}
	if err := s.store.CreateRecord(record); err != nil {
		return domain.Record{}, err
	}
	event := audit.NewEvent(record, "created", actor, "lead registered", s.next("evt"), now, record.Status)
	if err := audit.Append(s.store, event); err != nil {
		return domain.Record{}, err
	}
	return record, nil
}

func (s *Service) Get(id string) (domain.Record, error) {
	if err := s.ensureStore(); err != nil {
		return domain.Record{}, err
	}
	return s.store.GetRecord(id)
}
