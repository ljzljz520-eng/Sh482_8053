package search

import (
	"sort"
	"strings"

	"enterpriselead/internal/domain"
)

type Facets struct {
	Companies map[string]int `json:"companies"`
	Owners    map[string]int `json:"owners"`
	Statuses  map[string]int `json:"statuses"`
	Tags      map[string]int `json:"tags"`
}

func BuildFacets(records []domain.Record) Facets {
	facets := Facets{Companies: map[string]int{}, Owners: map[string]int{}, Statuses: map[string]int{}, Tags: map[string]int{}}
	for _, record := range records {
		facets.Companies[record.Company]++
		if record.Owner != "" {
			facets.Owners[record.Owner]++
		}
		facets.Statuses[string(record.Status)]++
		for _, tag := range record.Tags {
			facets.Tags[tag]++
		}
	}
	return facets
}

type FacetValue struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

func TopTags(facets Facets, limit int) []FacetValue {
	values := make([]FacetValue, 0, len(facets.Tags))
	for name, count := range facets.Tags {
		values = append(values, FacetValue{Name: name, Count: count})
	}
	sort.Slice(values, func(i, j int) bool {
		if values[i].Count != values[j].Count {
			return values[i].Count > values[j].Count
		}
		return values[i].Name < values[j].Name
	})
	if limit <= 0 || limit > len(values) {
		limit = len(values)
	}
	return values[:limit]
}

func NormalizeQuery(query domain.SearchQuery) domain.SearchQuery {
	query.Text = strings.TrimSpace(query.Text)
	query.Company = strings.TrimSpace(query.Company)
	query.Owner = strings.TrimSpace(query.Owner)
	query.Tag = strings.ToLower(strings.TrimSpace(query.Tag))
	if query.Limit > 500 {
		query.Limit = 500
	}
	if query.Offset < 0 {
		query.Offset = 0
	}
	return query
}

func GroupByStatus(records []domain.Record) map[domain.LeadStatus][]domain.Record {
	groups := make(map[domain.LeadStatus][]domain.Record)
	for _, record := range records {
		groups[record.Status] = append(groups[record.Status], record)
	}
	return groups
}
