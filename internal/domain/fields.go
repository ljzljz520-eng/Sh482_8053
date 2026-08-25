package domain

import (
	"sort"
	"strings"
)

func NormalizeTags(tags []string) []string {
	seen := make(map[string]bool, len(tags))
	result := make([]string, 0, len(tags))
	for _, raw := range tags {
		tag := strings.ToLower(strings.TrimSpace(raw))
		if tag == "" || seen[tag] {
			continue
		}
		seen[tag] = true
		result = append(result, tag)
	}
	sort.Strings(result)
	return result
}

func ApplyDefaults(record Record) Record {
	if record.Status == "" {
		record.Status = StatusDraft
	}
	if record.Priority == "" {
		record.Priority = PriorityNormal
	}
	if record.Source == "" {
		record.Source = "direct"
	}
	record.Tags = NormalizeTags(record.Tags)
	if record.Summary == "" {
		record.Summary = BuildSummary(record)
	}
	return record
}

func Diff(before, after Record) map[string]string {
	diff := make(map[string]string)
	if before.Company != after.Company {
		diff["company"] = before.Company + " -> " + after.Company
	}
	if before.ContactName != after.ContactName {
		diff["contact_name"] = before.ContactName + " -> " + after.ContactName
	}
	if before.ContactEmail != after.ContactEmail {
		diff["contact_email"] = before.ContactEmail + " -> " + after.ContactEmail
	}
	if before.Need != after.Need {
		diff["need"] = before.Need + " -> " + after.Need
	}
	if before.Owner != after.Owner {
		diff["owner"] = before.Owner + " -> " + after.Owner
	}
	if before.Priority != after.Priority {
		diff["priority"] = string(before.Priority) + " -> " + string(after.Priority)
	}
	if before.Summary != after.Summary {
		diff["summary"] = before.Summary + " -> " + after.Summary
	}
	if strings.Join(before.Tags, ",") != strings.Join(after.Tags, ",") {
		diff["tags"] = strings.Join(before.Tags, ",") + " -> " + strings.Join(after.Tags, ",")
	}
	return diff
}
