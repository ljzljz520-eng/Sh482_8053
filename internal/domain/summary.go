package domain

import (
	"fmt"
	"strings"
)

func BuildSummary(r Record) string {
	company := r.Company
	if company == "" {
		company = "Unknown company"
	}
	need := r.Need
	if need == "" {
		need = "unspecified need"
	}
	owner := r.Owner
	if owner == "" {
		owner = "unassigned"
	}
	return fmt.Sprintf("%s seeks %s; owner %s; priority %s; status %s", company, need, owner, r.Priority, r.Status)
}

func MatchText(r Record, text string) bool {
	needle := strings.ToLower(strings.TrimSpace(text))
	if needle == "" {
		return true
	}
	fields := []string{r.Company, r.ContactName, r.ContactEmail, r.Source, r.Need, r.Owner, r.Summary}
	for _, field := range fields {
		if strings.Contains(strings.ToLower(field), needle) {
			return true
		}
	}
	for _, tag := range r.Tags {
		if strings.Contains(strings.ToLower(tag), needle) {
			return true
		}
	}
	return false
}

func HasTag(r Record, tag string) bool {
	needle := strings.ToLower(strings.TrimSpace(tag))
	if needle == "" {
		return true
	}
	for _, candidate := range r.Tags {
		if strings.EqualFold(candidate, needle) {
			return true
		}
	}
	return false
}
