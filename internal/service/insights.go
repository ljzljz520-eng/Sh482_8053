package service

import (
	"sort"

	"enterpriselead/internal/domain"
	"enterpriselead/internal/search"
)

type Insight struct {
	RecordID string       `json:"record_id"`
	Company  string       `json:"company"`
	Score    domain.Score `json:"score"`
	Summary  string       `json:"summary"`
}

type Insights struct {
	Items       []Insight     `json:"items"`
	Facets      search.Facets `json:"facets"`
	Actionable  int           `json:"actionable"`
	GeneratedAt string        `json:"generated_at"`
}

func (s *Service) Insights(query domain.SearchQuery) (Insights, error) {
	if err := s.ensureStore(); err != nil {
		return Insights{}, err
	}
	records, err := s.store.ListRecords()
	if err != nil {
		return Insights{}, err
	}
	filtered := search.Filter(records, search.NormalizeQuery(query))
	items := make([]Insight, 0, len(filtered))
	actionable := 0
	for _, record := range filtered {
		score := domain.ScoreRecord(record)
		if domain.IsActionable(score) {
			actionable++
		}
		items = append(items, Insight{RecordID: record.ID, Company: record.Company, Score: score, Summary: domain.BuildSummary(record)})
	}
	sort.SliceStable(items, func(i, j int) bool {
		if items[i].Score.Value != items[j].Score.Value {
			return items[i].Score.Value > items[j].Score.Value
		}
		return items[i].RecordID < items[j].RecordID
	})
	return Insights{Items: items, Facets: search.BuildFacets(filtered), Actionable: actionable, GeneratedAt: s.now().Format("2006-01-02T15:04:05Z07:00")}, nil
}

func (s *Service) TopActionable(limit int) ([]Insight, error) {
	insights, err := s.Insights(domain.SearchQuery{IncludeArchived: false, Limit: 500})
	if err != nil {
		return nil, err
	}
	if limit <= 0 || limit > len(insights.Items) {
		limit = len(insights.Items)
	}
	return insights.Items[:limit], nil
}
