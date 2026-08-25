package ingest

import (
	"encoding/csv"
	"io"
	"strings"

	"enterpriselead/internal/domain"
)

func WriteTemplate(writer io.Writer) error {
	w := csv.NewWriter(writer)
	if err := w.Write([]string{"company", "contact", "email", "source", "need", "owner", "priority", "tags"}); err != nil {
		return err
	}
	w.Flush()
	return w.Error()
}

func CanonicalizeRow(row Row) Row {
	row.Company = strings.TrimSpace(row.Company)
	row.ContactName = strings.TrimSpace(row.ContactName)
	row.ContactEmail = strings.ToLower(strings.TrimSpace(row.ContactEmail))
	row.Source = strings.TrimSpace(row.Source)
	row.Need = strings.TrimSpace(row.Need)
	row.Owner = strings.TrimSpace(row.Owner)
	row.Priority = domain.Priority(strings.ToLower(string(row.Priority)))
	return row
}
