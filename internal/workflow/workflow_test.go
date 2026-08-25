package workflow

import (
	"testing"
	"time"

	"enterpriselead/internal/domain"
)

func TestReviewWorkflowLifecycle(t *testing.T) {
	record := domain.Record{ID: "r1", Status: domain.StatusReview}
	wf := NewReview(record, "operator", "wf-1", time.Unix(1, 0))
	if err := ValidateSteps(wf); err != nil {
		t.Fatal(err)
	}
	wf, err := MarkReviewed(wf, "reviewer")
	if err != nil {
		t.Fatal(err)
	}
	wf, err = MarkApproved(wf, "reviewer", time.Unix(2, 0))
	if err != nil {
		t.Fatal(err)
	}
	wf, err = MarkArchived(wf, time.Unix(3, 0))
	if err != nil {
		t.Fatal(err)
	}
	if !IsComplete(wf) || wf.CompletedAt == nil {
		t.Fatalf("unexpected workflow %#v", wf)
	}
}
