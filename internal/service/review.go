package service

import (
	"fmt"

	"enterpriselead/internal/audit"
	"enterpriselead/internal/domain"
	"enterpriselead/internal/policy"
	"enterpriselead/internal/workflow"
)

func (s *Service) Review(id, actor string) (domain.Record, domain.Workflow, error) {
	if err := s.ensureStore(); err != nil {
		return domain.Record{}, domain.Workflow{}, err
	}
	record, err := s.store.GetRecord(id)
	if err != nil {
		return domain.Record{}, domain.Workflow{}, err
	}
	if err := authorize(actor, "review", record); err != nil {
		return domain.Record{}, domain.Workflow{}, err
	}
	changed, err := domain.Transition(record, domain.StatusReview)
	if err != nil {
		return domain.Record{}, domain.Workflow{}, err
	}
	changed.UpdatedAt = s.now()
	changed.Version++
	wf := workflow.NewReview(changed, actor, s.next("wf"), changed.UpdatedAt)
	if err := workflow.ValidateSteps(wf); err != nil {
		return domain.Record{}, domain.Workflow{}, err
	}
	if err := s.store.PutRecord(changed); err != nil {
		return domain.Record{}, domain.Workflow{}, err
	}
	if err := s.store.PutWorkflow(wf); err != nil {
		return domain.Record{}, domain.Workflow{}, err
	}
	event := audit.NewEvent(record, "review_requested", actor, "lead sent for review", s.next("evt"), changed.UpdatedAt, changed.Status)
	if err := audit.Append(s.store, event); err != nil {
		return domain.Record{}, domain.Workflow{}, err
	}
	return changed, wf, nil
}

func (s *Service) Approve(id, actor string) (domain.Record, domain.Workflow, error) {
	if err := s.ensureStore(); err != nil {
		return domain.Record{}, domain.Workflow{}, err
	}
	record, err := s.store.GetRecord(id)
	if err != nil {
		return domain.Record{}, domain.Workflow{}, err
	}
	if err := authorize(actor, "approve", record); err != nil {
		return domain.Record{}, domain.Workflow{}, err
	}
	if err := workflow.ValidateForApproval(record); err != nil {
		return domain.Record{}, domain.Workflow{}, err
	}
	if decision := policy.EvaluateTransition(record, domain.StatusApproved, actor); !decision.Allowed {
		return domain.Record{}, domain.Workflow{}, fmt.Errorf("approval policy: %s", policy.Explain(decision))
	}
	changed, err := domain.Transition(record, domain.StatusApproved)
	if err != nil {
		return domain.Record{}, domain.Workflow{}, err
	}
	workflows, err := s.store.ListWorkflows(id)
	if err != nil || len(workflows) == 0 {
		return domain.Record{}, domain.Workflow{}, fmt.Errorf("review workflow missing")
	}
	wf := workflows[len(workflows)-1]
	wf, err = workflow.MarkReviewed(wf, actor)
	if err != nil {
		return domain.Record{}, domain.Workflow{}, err
	}
	wf, err = workflow.MarkApproved(wf, actor, s.now())
	if err != nil {
		return domain.Record{}, domain.Workflow{}, err
	}
	changed.UpdatedAt = s.now()
	changed.Version++
	changed.LastWorkflowID = wf.ID
	changed.Summary = domain.BuildSummary(changed)
	if err := s.store.PutRecord(changed); err != nil {
		return domain.Record{}, domain.Workflow{}, err
	}
	if err := s.store.PutWorkflow(wf); err != nil {
		return domain.Record{}, domain.Workflow{}, err
	}
	event := audit.NewEvent(record, "approved", actor, "lead approved and published", s.next("evt"), changed.UpdatedAt, changed.Status)
	if err := audit.Append(s.store, event); err != nil {
		return domain.Record{}, domain.Workflow{}, err
	}
	return changed, wf, nil
}

func (s *Service) Archive(id, actor string) (domain.Record, domain.Workflow, error) {
	if err := s.ensureStore(); err != nil {
		return domain.Record{}, domain.Workflow{}, err
	}
	record, err := s.store.GetRecord(id)
	if err != nil {
		return domain.Record{}, domain.Workflow{}, err
	}
	if err := authorize(actor, "archive", record); err != nil {
		return domain.Record{}, domain.Workflow{}, err
	}
	changed, err := domain.Transition(record, domain.StatusArchived)
	if err != nil {
		return domain.Record{}, domain.Workflow{}, err
	}
	workflows, err := s.store.ListWorkflows(id)
	if err != nil || len(workflows) == 0 {
		return domain.Record{}, domain.Workflow{}, fmt.Errorf("review workflow missing")
	}
	wf := workflows[len(workflows)-1]
	wf, err = workflow.MarkArchived(wf, s.now())
	if err != nil {
		return domain.Record{}, domain.Workflow{}, err
	}
	changed.UpdatedAt = s.now()
	when := changed.UpdatedAt
	changed.ArchivedAt = &when
	changed.Version++
	changed.LastWorkflowID = wf.ID
	if err := s.store.PutRecord(changed); err != nil {
		return domain.Record{}, domain.Workflow{}, err
	}
	if err := s.store.PutWorkflow(wf); err != nil {
		return domain.Record{}, domain.Workflow{}, err
	}
	event := audit.NewEvent(record, "archived", actor, "lead archived", s.next("evt"), changed.UpdatedAt, changed.Status)
	if err := audit.Append(s.store, event); err != nil {
		return domain.Record{}, domain.Workflow{}, err
	}
	return changed, wf, nil
}

func (s *Service) Reject(id, actor, note string) (domain.Record, error) {
	if err := s.ensureStore(); err != nil {
		return domain.Record{}, err
	}
	record, err := s.store.GetRecord(id)
	if err != nil {
		return domain.Record{}, err
	}
	changed, err := domain.Transition(record, domain.StatusRejected)
	if err != nil {
		return domain.Record{}, err
	}
	changed.UpdatedAt = s.now()
	changed.Version++
	if err := s.store.PutRecord(changed); err != nil {
		return domain.Record{}, err
	}
	event := audit.NewEvent(record, "rejected", actor, note, s.next("evt"), changed.UpdatedAt, changed.Status)
	if err := audit.Append(s.store, event); err != nil {
		return domain.Record{}, err
	}
	return changed, nil
}
