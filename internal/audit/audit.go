package audit

import (
	"fmt"
	"sort"
	"time"

	"enterpriselead/internal/domain"
	"enterpriselead/internal/storage"
)

type Clock interface{ Now() time.Time }

type FixedClock struct{ Current time.Time }

func (c FixedClock) Now() time.Time { return c.Current }

func NewEvent(record domain.Record, action, actor, note, id string, now time.Time, to domain.LeadStatus) domain.AuditEvent {
	return domain.AuditEvent{ID: id, RecordID: record.ID, Action: action, Actor: actor, FromStatus: record.Status, ToStatus: to, Note: note, At: now.UTC(), Metadata: map[string]string{"version": fmt.Sprintf("%d", record.Version)}}
}

func Append(store *storage.DB, event domain.AuditEvent) error {
	if event.At.IsZero() {
		event.At = time.Unix(0, 0).UTC()
	}
	return store.PutEvent(event)
}

func Timeline(store *storage.DB, recordID string) ([]domain.AuditEvent, error) {
	events, err := store.ListEvents(recordID)
	if err != nil {
		return nil, err
	}
	sort.SliceStable(events, func(i, j int) bool {
		if events[i].At.Equal(events[j].At) {
			return events[i].ID < events[j].ID
		}
		return events[i].At.Before(events[j].At)
	})
	return events, nil
}

func Latest(events []domain.AuditEvent) (domain.AuditEvent, bool) {
	if len(events) == 0 {
		return domain.AuditEvent{}, false
	}
	latest := events[0]
	for _, event := range events[1:] {
		if event.At.After(latest.At) || (event.At.Equal(latest.At) && event.ID > latest.ID) {
			latest = event
		}
	}
	return latest, true
}
