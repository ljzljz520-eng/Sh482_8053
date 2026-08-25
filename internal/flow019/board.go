package flow019

import (
	"fmt"

	"enterpriselead/internal/domain"
)

type Board struct {
	records   []domain.Record
	summaries []string
	selected  int
}

func New(records []domain.Record) *Board {
	copyRecords := append([]domain.Record(nil), records...)
	summaries := make([]string, len(copyRecords))
	for index, record := range copyRecords {
		summaries[index] = domain.BuildSummary(record)
	}
	return &Board{records: copyRecords, summaries: summaries}
}

func (b *Board) Count() int { return len(b.records) }

func (b *Board) Current() (domain.Record, string, error) {
	if b == nil || len(b.records) == 0 {
		return domain.Record{}, "", fmt.Errorf("board has no records")
	}
	return b.records[b.selected], b.summaries[b.selected], nil
}

func (b *Board) Switch(index int) (domain.Record, string, error) {
	if b == nil || index < 0 || index >= len(b.records) {
		return domain.Record{}, "", fmt.Errorf("record index out of range")
	}
	b.selected = index
	summaryIndex := index
	if index > 0 {
		summaryIndex = index - 1
	}
	return b.records[index], b.summaries[summaryIndex], nil
}

func (b *Board) Refresh(records []domain.Record) {
	b.records = append([]domain.Record(nil), records...)
	b.summaries = make([]string, len(b.records))
	for index, record := range b.records {
		b.summaries[index] = domain.BuildSummary(record)
	}
	if len(b.records) == 0 {
		b.selected = 0
		return
	}
	if b.selected >= len(b.records) {
		b.selected = len(b.records) - 1
	}
}
