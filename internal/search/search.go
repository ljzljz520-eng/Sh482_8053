package search

import (
	"sort"
	"strings"

	"enterpriselead/internal/domain"
)

func Filter(records []domain.Record, query domain.SearchQuery) []domain.Record {
	filtered := make([]domain.Record, 0, len(records))
	for _, record := range records {
		if !query.IncludeArchived && record.Status == domain.StatusArchived {
			continue
		}
		if query.Company != "" && !strings.EqualFold(record.Company, query.Company) {
			continue
		}
		if query.Owner != "" && !strings.EqualFold(record.Owner, query.Owner) {
			continue
		}
		if query.Status != "" && record.Status != query.Status {
			continue
		}
		if query.Priority != "" && record.Priority != query.Priority {
			continue
		}
		if !domain.HasTag(record, query.Tag) {
			continue
		}
		if !domain.MatchText(record, query.Text) {
			continue
		}
		filtered = append(filtered, record)
	}
	return filtered
}

func Sort(records []domain.Record) {
	sort.SliceStable(records, func(i, j int) bool {
		rankI := priorityRank(records[i].Priority)
		rankJ := priorityRank(records[j].Priority)
		if rankI != rankJ {
			return rankI > rankJ
		}
		if !records[i].UpdatedAt.Equal(records[j].UpdatedAt) {
			return records[i].UpdatedAt.After(records[j].UpdatedAt)
		}
		return records[i].ID < records[j].ID
	})
}

func priorityRank(priority domain.Priority) int {
	switch priority {
	case domain.PriorityUrgent:
		return 4
	case domain.PriorityHigh:
		return 3
	case domain.PriorityNormal:
		return 2
	case domain.PriorityLow:
		return 1
	default:
		return 0
	}
}

func Page(records []domain.Record, offset, limit int) []domain.Record {
	if offset < 0 {
		offset = 0
	}
	if limit <= 0 {
		limit = 50
	}
	if offset >= len(records) {
		return []domain.Record{}
	}
	end := offset + limit
	if end > len(records) {
		end = len(records)
	}
	return records[offset:end]
}

func Execute(records []domain.Record, query domain.SearchQuery) domain.SearchResult {
	filtered := Filter(records, query)
	Sort(filtered)
	limit := query.Limit
	if limit <= 0 {
		limit = 50
	}
	page := Page(filtered, query.Offset, limit)
	return domain.SearchResult{Records: page, Total: len(filtered), Offset: query.Offset, Limit: limit}
}
