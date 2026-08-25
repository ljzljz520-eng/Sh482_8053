package domain

import (
	"fmt"
	"regexp"
	"strings"
)

var emailPattern = regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)

func ValidateRecord(r Record) error {
	if strings.TrimSpace(r.Company) == "" {
		return fmt.Errorf("company is required")
	}
	if len([]rune(r.Company)) > 160 {
		return fmt.Errorf("company is too long")
	}
	if strings.TrimSpace(r.ContactName) == "" {
		return fmt.Errorf("contact name is required")
	}
	if !emailPattern.MatchString(strings.TrimSpace(r.ContactEmail)) {
		return fmt.Errorf("contact email is invalid")
	}
	if strings.TrimSpace(r.Need) == "" {
		return fmt.Errorf("need is required")
	}
	if r.Priority != PriorityLow && r.Priority != PriorityNormal && r.Priority != PriorityHigh && r.Priority != PriorityUrgent {
		return fmt.Errorf("unknown priority %q", r.Priority)
	}
	if r.Status == "" {
		return fmt.Errorf("status is required")
	}
	if len(r.Tags) > 20 {
		return fmt.Errorf("too many tags")
	}
	seen := make(map[string]struct{}, len(r.Tags))
	for _, tag := range r.Tags {
		normalized := strings.ToLower(strings.TrimSpace(tag))
		if normalized == "" {
			return fmt.Errorf("tag cannot be empty")
		}
		if len(normalized) > 48 {
			return fmt.Errorf("tag %q is too long", tag)
		}
		if _, exists := seen[normalized]; exists {
			return fmt.Errorf("duplicate tag %q", tag)
		}
		seen[normalized] = struct{}{}
	}
	return nil
}

func ValidateChange(c Change) error {
	if c.Company != nil && strings.TrimSpace(*c.Company) == "" {
		return fmt.Errorf("company cannot be empty")
	}
	if c.ContactName != nil && strings.TrimSpace(*c.ContactName) == "" {
		return fmt.Errorf("contact name cannot be empty")
	}
	if c.ContactEmail != nil && !emailPattern.MatchString(strings.TrimSpace(*c.ContactEmail)) {
		return fmt.Errorf("contact email is invalid")
	}
	if c.Need != nil && strings.TrimSpace(*c.Need) == "" {
		return fmt.Errorf("need cannot be empty")
	}
	if c.Priority != nil {
		if *c.Priority != PriorityLow && *c.Priority != PriorityNormal && *c.Priority != PriorityHigh && *c.Priority != PriorityUrgent {
			return fmt.Errorf("unknown priority %q", *c.Priority)
		}
	}
	if c.Tags != nil {
		if len(c.Tags) > 20 {
			return fmt.Errorf("too many tags")
		}
		for _, tag := range c.Tags {
			if strings.TrimSpace(tag) == "" {
				return fmt.Errorf("tag cannot be empty")
			}
		}
	}
	return nil
}

func NormalizeRecord(r Record) Record {
	r.Company = strings.TrimSpace(r.Company)
	r.ContactName = strings.TrimSpace(r.ContactName)
	r.ContactEmail = strings.ToLower(strings.TrimSpace(r.ContactEmail))
	r.Source = strings.TrimSpace(r.Source)
	r.Need = strings.TrimSpace(r.Need)
	r.Owner = strings.TrimSpace(r.Owner)
	r.Summary = strings.TrimSpace(r.Summary)
	for i, tag := range r.Tags {
		r.Tags[i] = strings.ToLower(strings.TrimSpace(tag))
	}
	return r
}
