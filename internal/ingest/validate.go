package ingest

import (
	"fmt"
	"strings"

	"enterpriselead/internal/domain"
)

type RowFinding struct {
	Line    int    `json:"line"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func ValidateRow(row Row) []RowFinding {
	findings := make([]RowFinding, 0)
	if strings.TrimSpace(row.Company) == "" {
		findings = append(findings, RowFinding{Line: row.Line, Code: "company", Message: "company is required"})
	}
	if strings.TrimSpace(row.ContactName) == "" {
		findings = append(findings, RowFinding{Line: row.Line, Code: "contact", Message: "contact is required"})
	}
	if !strings.Contains(row.ContactEmail, "@") {
		findings = append(findings, RowFinding{Line: row.Line, Code: "email", Message: "email is invalid"})
	}
	if row.Priority == "" {
		findings = append(findings, RowFinding{Line: row.Line, Code: "priority", Message: "priority is required"})
	}
	if row.Priority != domain.PriorityLow && row.Priority != domain.PriorityNormal && row.Priority != domain.PriorityHigh && row.Priority != domain.PriorityUrgent {
		findings = append(findings, RowFinding{Line: row.Line, Code: "priority", Message: "priority is unknown"})
	}
	if strings.TrimSpace(row.Need) == "" {
		findings = append(findings, RowFinding{Line: row.Line, Code: "need", Message: "need is required"})
	}
	return findings
}

func FindingsError(findings []RowFinding) error {
	if len(findings) == 0 {
		return nil
	}
	parts := make([]string, 0, len(findings))
	for _, finding := range findings {
		parts = append(parts, fmt.Sprintf("line %d %s: %s", finding.Line, finding.Code, finding.Message))
	}
	return fmt.Errorf("invalid import row: %s", strings.Join(parts, "; "))
}

func ValidateRows(rows []Row) []RowFinding {
	findings := make([]RowFinding, 0)
	for _, row := range rows {
		findings = append(findings, ValidateRow(row)...)
	}
	return findings
}
