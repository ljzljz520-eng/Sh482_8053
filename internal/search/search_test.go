package search

import (
	"testing"
	"time"

	"enterpriselead/internal/domain"
)

func TestExecuteFiltersSortsAndPages(t *testing.T) {
	records := []domain.Record{
		{ID: "r1", Company: "Acme", Owner: "a", Status: domain.StatusDraft, Priority: domain.PriorityNormal, Tags: []string{"cloud"}, UpdatedAt: time.Unix(1, 0)},
		{ID: "r2", Company: "Acme", Owner: "b", Status: domain.StatusApproved, Priority: domain.PriorityUrgent, Tags: []string{"cloud"}, UpdatedAt: time.Unix(2, 0)},
		{ID: "r3", Company: "Other", Owner: "c", Status: domain.StatusArchived, Priority: domain.PriorityHigh, Tags: []string{"cloud"}, UpdatedAt: time.Unix(3, 0)},
	}
	result := Execute(records, domain.SearchQuery{Company: "acme", Tag: "cloud", Limit: 1})
	if result.Total != 2 || len(result.Records) != 1 || result.Records[0].ID != "r2" {
		t.Fatalf("unexpected result %#v", result)
	}
	result = Execute(records, domain.SearchQuery{IncludeArchived: true, Offset: 10, Limit: 2})
	if len(result.Records) != 0 {
		t.Fatalf("expected empty page %#v", result.Records)
	}
}
