package flow019

import (
	"testing"

	"enterpriselead/internal/domain"
)

func Test482BusinessRegression(t *testing.T) {
	first := domain.Record{ID: "first", Company: "First Corp", ContactName: "A", ContactEmail: "a@first.test", Need: "contracts", Owner: "Ops", Status: domain.StatusDraft, Priority: domain.PriorityHigh}
	second := domain.Record{ID: "second", Company: "Second Corp", ContactName: "B", ContactEmail: "b@second.test", Need: "analytics", Owner: "Sales", Status: domain.StatusDraft, Priority: domain.PriorityNormal}
	board := New([]domain.Record{first, second})
	view, err := board.SelectByID(second.ID)
	if err != nil {
		t.Fatal(err)
	}
	if view.Record.ID != second.ID {
		t.Fatalf("unexpected selected record %s", view.Record.ID)
	}
	if view.Summary != domain.BuildSummary(second) {
		t.Fatalf("unexpected summary %q", view.Summary)
	}
}
