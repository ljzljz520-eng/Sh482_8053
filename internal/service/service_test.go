package service_test

import (
	"strings"
	"testing"
	"time"

	"enterpriselead/internal/domain"
	"enterpriselead/internal/ingest"
	"enterpriselead/internal/service"
	"enterpriselead/internal/storage"
)

func testService(t *testing.T) (*service.Service, *storage.DB) {
	t.Helper()
	store, err := storage.Open(t.TempDir() + "/service.db")
	if err != nil {
		t.Fatal(err)
	}
	return service.New(store, service.FixedTime{Value: time.Unix(100, 0)}), store
}

func createInput() service.CreateInput {
	return service.CreateInput{Company: "Acme", ContactName: "Lin", ContactEmail: "lin@acme.test", Source: "event", Need: "contract workflow", Owner: "Ming", Priority: domain.PriorityHigh, Tags: []string{"cloud"}}
}

func TestWorkflowCreateReviewArchive(t *testing.T) {
	app, store := testService(t)
	defer store.Close()
	record, err := app.Create(createInput(), "operator")
	if err != nil {
		t.Fatal(err)
	}
	record, _, err = app.Review(record.ID, "reviewer")
	if err != nil {
		t.Fatal(err)
	}
	if record.Status != domain.StatusReview {
		t.Fatalf("unexpected status %s", record.Status)
	}
	record, _, err = app.Approve(record.ID, "reviewer")
	if err != nil {
		t.Fatal(err)
	}
	if record.Status != domain.StatusApproved || !strings.Contains(record.Summary, "Acme") {
		t.Fatalf("unexpected approved record %#v", record)
	}
	record, _, err = app.Archive(record.ID, "operator")
	if err != nil {
		t.Fatal(err)
	}
	if record.Status != domain.StatusArchived || record.ArchivedAt == nil {
		t.Fatalf("unexpected archived record %#v", record)
	}
	events, err := app.Timeline(record.ID)
	if err != nil || len(events) != 4 {
		t.Fatalf("timeline: %v %#v", err, events)
	}
}

func TestWorkflowSearchUpdatePublish(t *testing.T) {
	app, store := testService(t)
	defer store.Close()
	first, err := app.Create(createInput(), "operator")
	if err != nil {
		t.Fatal(err)
	}
	secondInput := createInput()
	secondInput.Company = "Beta"
	secondInput.ContactEmail = "b@beta.test"
	secondInput.Priority = domain.PriorityUrgent
	second, err := app.Create(secondInput, "operator")
	if err != nil {
		t.Fatal(err)
	}
	result, err := app.Search(domain.SearchQuery{Text: "beta", Limit: 10})
	if err != nil {
		t.Fatal(err)
	}
	if result.Total != 1 || result.Records[0].ID != second.ID {
		t.Fatalf("unexpected search %#v", result)
	}
	owner := "Sales"
	updated, err := app.Update(first.ID, first.Version, domain.Change{Owner: &owner}, "operator")
	if err != nil {
		t.Fatal(err)
	}
	if updated.Owner != owner || !strings.Contains(updated.Summary, owner) {
		t.Fatalf("summary was not published %#v", updated)
	}
	if _, err := app.Update(first.ID, first.Version, domain.Change{Owner: &owner}, "operator"); err == nil {
		t.Fatal("expected version conflict")
	}
}

func TestWorkflowImportReport(t *testing.T) {
	app, store := testService(t)
	defer store.Close()
	input := "company,contact,email,source,need,owner,priority,tags\nAcme,Lin,lin@acme.test,event,contracts,Ming,high,cloud|renewal\nBroken,No,bad,event,contracts,Ming,high,cloud\n"
	report, err := ingest.New(app).Import(strings.NewReader(input), "importer")
	if err != nil {
		t.Fatal(err)
	}
	if report.Total != 2 || report.Succeeded != 1 || report.Failed != 1 || len(report.Results) != 1 {
		t.Fatalf("unexpected report %#v", report)
	}
	result, err := app.Search(domain.SearchQuery{Limit: 10})
	if err != nil || result.Total != 1 {
		t.Fatalf("unexpected imported data: %v %#v", err, result)
	}
}

func TestAttachmentLifecycle(t *testing.T) {
	app, store := testService(t)
	defer store.Close()
	record, err := app.Create(createInput(), "operator")
	if err != nil {
		t.Fatal(err)
	}
	attachment, err := app.SaveAttachment(record.ID, "brief.txt", "text/plain", []byte("brief"))
	if err != nil {
		t.Fatal(err)
	}
	got, err := app.GetAttachment(attachment.ID)
	if err != nil || string(got.Content) != "brief" {
		t.Fatalf("unexpected attachment: %v %#v", err, got)
	}
}
