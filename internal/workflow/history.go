package workflow

import (
	"sort"
	"time"

	"enterpriselead/internal/domain"
)

type Milestone struct {
	Name  string    `json:"name"`
	At    time.Time `json:"at"`
	Actor string    `json:"actor"`
}

func BuildMilestones(events []domain.AuditEvent) []Milestone {
	result := make([]Milestone, 0, len(events))
	for _, event := range events {
		result = append(result, Milestone{Name: event.Action, At: event.At, Actor: event.Actor})
	}
	sort.SliceStable(result, func(i, j int) bool { return result[i].At.Before(result[j].At) })
	return result
}

func DurationBetween(events []domain.AuditEvent, first, last string) (time.Duration, bool) {
	var start, end time.Time
	for _, event := range events {
		if event.Action == first && (start.IsZero() || event.At.Before(start)) {
			start = event.At
		}
		if event.Action == last && event.At.After(end) {
			end = event.At
		}
	}
	if start.IsZero() || end.IsZero() || end.Before(start) {
		return 0, false
	}
	return end.Sub(start), true
}

func Summarize(events []domain.AuditEvent) string {
	if len(events) == 0 {
		return "no activity"
	}
	latest := events[0]
	for _, event := range events[1:] {
		if event.At.After(latest.At) {
			latest = event
		}
	}
	return latest.Action + " by " + latest.Actor
}
