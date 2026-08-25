package service

import (
	"encoding/csv"
	"io"

	"enterpriselead/internal/domain"
)

func (s *Service) Export(writer io.Writer, query domain.SearchQuery) error {
	result, err := s.Search(query)
	if err != nil {
		return err
	}
	csvWriter := csv.NewWriter(writer)
	if err := csvWriter.Write([]string{"id", "company", "contact", "email", "need", "owner", "status", "priority", "summary"}); err != nil {
		return err
	}
	for _, record := range result.Records {
		row := []string{record.ID, record.Company, record.ContactName, record.ContactEmail, record.Need, record.Owner, string(record.Status), string(record.Priority), record.Summary}
		if err := csvWriter.Write(row); err != nil {
			return err
		}
	}
	csvWriter.Flush()
	return csvWriter.Error()
}
