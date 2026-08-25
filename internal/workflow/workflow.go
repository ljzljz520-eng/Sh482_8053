package workflow

import (
	"fmt"
	"time"

	"enterpriselead/internal/domain"
)

type IDSource interface{ Next(prefix string) string }

type Sequence struct{ Value int }

func (s *Sequence) Next(prefix string) string {
	s.Value++
	return fmt.Sprintf("%s-%04d", prefix, s.Value)
}

func NewReview(record domain.Record, requester, id string, now time.Time) domain.Workflow {
	return domain.Workflow{ID: id, RecordID: record.ID, Kind: "review", RequestedBy: requester, State: "requested", Steps: []string{"submitted", "reviewed", "approved", "archived"}, CreatedAt: now.UTC()}
}

func MarkReviewed(w domain.Workflow, actor string) (domain.Workflow, error) {
	if w.State != "requested" {
		return w, fmt.Errorf("workflow is not requested")
	}
	w.State = "reviewed"
	w.ApprovedBy = actor
	return w, nil
}

func MarkApproved(w domain.Workflow, actor string, now time.Time) (domain.Workflow, error) {
	if w.State != "reviewed" {
		return w, fmt.Errorf("workflow is not reviewed")
	}
	w.State = "approved"
	w.ApprovedBy = actor
	w.CompletedAt = nil
	return w, nil
}

func MarkArchived(w domain.Workflow, now time.Time) (domain.Workflow, error) {
	if w.State != "approved" {
		return w, fmt.Errorf("workflow is not approved")
	}
	w.State = "archived"
	when := now.UTC()
	w.CompletedAt = &when
	return w, nil
}

func ValidateSteps(w domain.Workflow) error {
	if len(w.Steps) < 4 {
		return fmt.Errorf("workflow needs at least four steps")
	}
	for _, step := range w.Steps {
		if step == "" {
			return fmt.Errorf("workflow step cannot be empty")
		}
	}
	return nil
}

func IsComplete(w domain.Workflow) bool { return w.State == "archived" }
