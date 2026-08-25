package audit

import (
	"testing"
	"time"

	"enterpriselead/internal/domain"
	"enterpriselead/internal/storage"
)

func TestTimelineOrdersEvents(t *testing.T) {
	store, err := storage.Open(t.TempDir() + "/audit.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	record := domain.Record{ID: "rec-1", Status: domain.StatusDraft}
	if err := Append(store, NewEvent(record, "second", "u", "", "evt-2", time.Unix(20, 0), domain.StatusReview)); err != nil {
		t.Fatal(err)
	}
	if err := Append(store, NewEvent(record, "first", "u", "", "evt-1", time.Unix(10, 0), domain.StatusDraft)); err != nil {
		t.Fatal(err)
	}
	events, err := Timeline(store, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != 2 || events[0].Action != "first" {
		t.Fatalf("unexpected timeline %#v", events)
	}
	latest, ok := Latest(events)
	if !ok || latest.Action != "second" {
		t.Fatalf("unexpected latest %#v", latest)
	}
}
