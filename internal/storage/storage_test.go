package storage

import (
	"testing"
	"time"

	"enterpriselead/internal/domain"
)

func TestPersistenceSurvivesReopen(t *testing.T) {
	path := t.TempDir() + "/leads.db"
	record := domain.Record{ID: "rec-1", Company: "Northwind", ContactName: "A", ContactEmail: "a@northwind.test", Need: "contracts", Owner: "ops", Status: domain.StatusDraft, Priority: domain.PriorityNormal, CreatedAt: time.Unix(10, 0), UpdatedAt: time.Unix(10, 0), Version: 1}
	first, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := first.CreateRecord(record); err != nil {
		t.Fatal(err)
	}
	if err := first.PutEvent(domain.AuditEvent{ID: "evt-1", RecordID: record.ID, Action: "created", At: time.Unix(10, 0)}); err != nil {
		t.Fatal(err)
	}
	if err := first.PutWorkflow(domain.Workflow{ID: "wf-1", RecordID: record.ID, Kind: "review", State: "requested", Steps: []string{"a", "b", "c", "d"}}); err != nil {
		t.Fatal(err)
	}
	if err := first.PutAttachment(domain.Attachment{ID: "att-1", RecordID: record.ID, Name: "brief.txt", MediaType: "text/plain", Size: 4, Content: []byte("brief")[:4]}); err != nil {
		t.Fatal(err)
	}
	if err := first.Close(); err != nil {
		t.Fatal(err)
	}
	second, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer second.Close()
	got, err := second.GetRecord(record.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Company != record.Company || got.Version != 1 {
		t.Fatalf("unexpected record after reopen: %#v", got)
	}
	events, err := second.ListEvents(record.ID)
	if err != nil || len(events) != 1 {
		t.Fatalf("events: %v %#v", err, events)
	}
	workflows, err := second.ListWorkflows(record.ID)
	if err != nil || len(workflows) != 1 {
		t.Fatalf("workflows: %v %#v", err, workflows)
	}
	attachment, err := second.GetAttachment("att-1")
	if err != nil || string(attachment.Content) != "brie" {
		t.Fatalf("attachment: %v %#v", err, attachment)
	}
}

func TestAttachmentSizeLimit(t *testing.T) {
	if _, err := PrepareAttachment(domain.Attachment{Size: MaxAttachmentSize + 1, Content: make([]byte, MaxAttachmentSize+1)}); err == nil {
		t.Fatal("expected size error")
	}
}
