package audit

import (
	"encoding/json"
	"fmt"
	"strings"

	"enterpriselead/internal/domain"
)

type Entry struct {
	At     string `json:"at"`
	Actor  string `json:"actor"`
	Action string `json:"action"`
	Note   string `json:"note"`
}

func ToEntries(events []domain.AuditEvent) []Entry {
	entries := make([]Entry, 0, len(events))
	for _, event := range events {
		entries = append(entries, Entry{At: event.At.Format("2006-01-02T15:04:05Z07:00"), Actor: event.Actor, Action: event.Action, Note: event.Note})
	}
	return entries
}

func Render(events []domain.AuditEvent) string {
	entries := ToEntries(events)
	parts := make([]string, 0, len(entries))
	for _, entry := range entries {
		parts = append(parts, fmt.Sprintf("%s %s %s", entry.At, entry.Actor, entry.Action))
	}
	return strings.Join(parts, "\n")
}

func Marshal(events []domain.AuditEvent) ([]byte, error) { return json.Marshal(ToEntries(events)) }
