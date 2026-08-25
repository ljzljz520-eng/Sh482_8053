package ingest

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"

	"enterpriselead/internal/domain"
	"enterpriselead/internal/service"
)

type Row struct {
	Line         int
	Company      string
	ContactName  string
	ContactEmail string
	Source       string
	Need         string
	Owner        string
	Priority     domain.Priority
	Tags         []string
}

type Importer struct{ app *service.Service }

func New(app *service.Service) *Importer { return &Importer{app: app} }

func Parse(reader io.Reader) ([]Row, []string, error) {
	csvReader := csv.NewReader(reader)
	csvReader.TrimLeadingSpace = true
	rows := make([]Row, 0)
	errors := make([]string, 0)
	line := 0
	for {
		fields, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		line++
		if err != nil {
			errors = append(errors, fmt.Sprintf("line %d: %v", line, err))
			continue
		}
		if line == 1 && strings.EqualFold(strings.TrimSpace(fields[0]), "company") {
			continue
		}
		if len(fields) < 8 {
			errors = append(errors, fmt.Sprintf("line %d: expected 8 columns", line))
			continue
		}
		rows = append(rows, Row{Line: line, Company: fields[0], ContactName: fields[1], ContactEmail: fields[2], Source: fields[3], Need: fields[4], Owner: fields[5], Priority: domain.Priority(strings.ToLower(strings.TrimSpace(fields[6]))), Tags: splitTags(fields[7])})
	}
	return rows, errors, nil
}

func splitTags(value string) []string {
	parts := strings.Split(value, "|")
	tags := make([]string, 0, len(parts))
	for _, part := range parts {
		if strings.TrimSpace(part) != "" {
			tags = append(tags, strings.TrimSpace(part))
		}
	}
	return tags
}

type RowResult struct {
	Line  int
	ID    string
	Error string
}

func (i *Importer) Import(reader io.Reader, actor string) (Report, error) {
	rows, parseErrors, err := Parse(reader)
	if err != nil {
		return Report{}, err
	}
	report := Report{Total: len(rows) + len(parseErrors), Errors: append([]string(nil), parseErrors...)}
	for _, row := range rows {
		row = CanonicalizeRow(row)
		if findings := ValidateRow(row); len(findings) > 0 {
			report.Failed++
			report.Errors = append(report.Errors, FindingsError(findings).Error())
			continue
		}
		record, createErr := i.app.Create(service.CreateInput{Company: row.Company, ContactName: row.ContactName, ContactEmail: row.ContactEmail, Source: row.Source, Need: row.Need, Owner: row.Owner, Priority: row.Priority, Tags: row.Tags}, actor)
		if createErr != nil {
			report.Failed++
			report.Errors = append(report.Errors, fmt.Sprintf("line %d: %v", row.Line, createErr))
			continue
		}
		report.Succeeded++
		report.Results = append(report.Results, RowResult{Line: row.Line, ID: record.ID})
	}
	return report, nil
}
