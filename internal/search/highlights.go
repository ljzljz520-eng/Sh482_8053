package search

import (
	"strings"

	"enterpriselead/internal/domain"
)

type Highlight struct {
	Field string `json:"field"`
	Value string `json:"value"`
}

func HighlightMatches(record domain.Record, text string) []Highlight {
	needle := strings.ToLower(strings.TrimSpace(text))
	if needle == "" {
		return []Highlight{}
	}
	fields := []struct{ name, value string }{{"company", record.Company}, {"contact", record.ContactName}, {"need", record.Need}, {"summary", record.Summary}}
	result := make([]Highlight, 0)
	for _, field := range fields {
		if strings.Contains(strings.ToLower(field.value), needle) {
			result = append(result, Highlight{Field: field.name, Value: field.value})
		}
	}
	return result
}

func SearchOne(records []domain.Record, query domain.SearchQuery) (domain.Record, []Highlight, bool) {
	for _, record := range Filter(records, query) {
		highlights := HighlightMatches(record, query.Text)
		return record, highlights, true
	}
	return domain.Record{}, nil, false
}
