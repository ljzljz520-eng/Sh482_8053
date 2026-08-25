package policy

import (
	"fmt"
	"strings"

	"enterpriselead/internal/domain"
)

type Finding struct {
	Code     string `json:"code"`
	Severity string `json:"severity"`
	Message  string `json:"message"`
}

type Decision struct {
	Allowed  bool      `json:"allowed"`
	Findings []Finding `json:"findings"`
}

func EvaluateCreate(record domain.Record) Decision {
	findings := make([]Finding, 0)
	if strings.EqualFold(record.ContactEmail, "unknown@example.com") {
		findings = append(findings, Finding{Code: "placeholder-email", Severity: "warning", Message: "contact email is a placeholder"})
	}
	if record.Priority == domain.PriorityUrgent && strings.TrimSpace(record.Owner) == "" {
		findings = append(findings, Finding{Code: "urgent-owner", Severity: "error", Message: "urgent leads require an owner"})
	}
	if len(record.Need) < 12 {
		findings = append(findings, Finding{Code: "thin-need", Severity: "warning", Message: "need description is brief"})
	}
	return Decision{Allowed: !hasError(findings), Findings: findings}
}

func EvaluateChange(before, after domain.Record) Decision {
	findings := make([]Finding, 0)
	if before.Company != after.Company && before.Status != domain.StatusDraft {
		findings = append(findings, Finding{Code: "company-lock", Severity: "error", Message: "company cannot change after submission"})
	}
	if before.Priority != after.Priority && after.Priority == domain.PriorityUrgent && after.Owner == "" {
		findings = append(findings, Finding{Code: "urgent-owner", Severity: "error", Message: "urgent leads require an owner"})
	}
	if before.ContactEmail != after.ContactEmail {
		findings = append(findings, Finding{Code: "contact-changed", Severity: "warning", Message: "contact identity changed"})
	}
	return Decision{Allowed: !hasError(findings), Findings: findings}
}

func EvaluateTransition(record domain.Record, target domain.LeadStatus, actor string) Decision {
	findings := make([]Finding, 0)
	if !domain.CanTransition(record.Status, target) {
		findings = append(findings, Finding{Code: "transition", Severity: "error", Message: fmt.Sprintf("cannot move %s to %s", record.Status, target)})
	}
	if strings.TrimSpace(actor) == "" {
		findings = append(findings, Finding{Code: "actor", Severity: "error", Message: "actor is required"})
	}
	if target == domain.StatusApproved && record.Priority == domain.PriorityUrgent && len(record.Tags) == 0 {
		findings = append(findings, Finding{Code: "urgent-context", Severity: "warning", Message: "urgent approval has no tags"})
	}
	return Decision{Allowed: !hasError(findings), Findings: findings}
}

func hasError(findings []Finding) bool {
	for _, finding := range findings {
		if finding.Severity == "error" {
			return true
		}
	}
	return false
}

func Explain(decision Decision) string {
	if decision.Allowed && len(decision.Findings) == 0 {
		return "allowed"
	}
	parts := make([]string, 0, len(decision.Findings))
	for _, finding := range decision.Findings {
		parts = append(parts, finding.Code+": "+finding.Message)
	}
	if decision.Allowed {
		return "allowed with findings: " + strings.Join(parts, "; ")
	}
	return "blocked: " + strings.Join(parts, "; ")
}
