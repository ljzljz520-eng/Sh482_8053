package workflow

import (
	"fmt"
	"strings"

	"enterpriselead/internal/domain"
)

type ChecklistItem struct {
	Key      string `json:"key"`
	Label    string `json:"label"`
	Required bool   `json:"required"`
	Complete bool   `json:"complete"`
}

func DefaultChecklist(record domain.Record) []ChecklistItem {
	return []ChecklistItem{
		{Key: "company", Label: "company identified", Required: true, Complete: record.Company != ""},
		{Key: "contact", Label: "contact verified", Required: true, Complete: record.ContactName != "" && record.ContactEmail != ""},
		{Key: "need", Label: "procurement need described", Required: true, Complete: len(strings.TrimSpace(record.Need)) >= 8},
		{Key: "owner", Label: "owner assigned", Required: true, Complete: record.Owner != ""},
		{Key: "source", Label: "source recorded", Required: false, Complete: record.Source != ""},
		{Key: "tags", Label: "classification tags added", Required: false, Complete: len(record.Tags) > 0},
	}
}

func ChecklistReady(items []ChecklistItem) bool {
	for _, item := range items {
		if item.Required && !item.Complete {
			return false
		}
	}
	return true
}

func ValidateForApproval(record domain.Record) error {
	items := DefaultChecklist(record)
	if !ChecklistReady(items) {
		return fmt.Errorf("approval checklist is incomplete")
	}
	if record.Status != domain.StatusReview {
		return fmt.Errorf("record is not in review")
	}
	return nil
}

func CompleteItem(items []ChecklistItem, key string) ([]ChecklistItem, error) {
	updated := append([]ChecklistItem(nil), items...)
	for index := range updated {
		if updated[index].Key == key {
			updated[index].Complete = true
			return updated, nil
		}
	}
	return items, fmt.Errorf("unknown checklist item %s", key)
}

func MissingRequired(items []ChecklistItem) []string {
	missing := make([]string, 0)
	for _, item := range items {
		if item.Required && !item.Complete {
			missing = append(missing, item.Key)
		}
	}
	return missing
}
