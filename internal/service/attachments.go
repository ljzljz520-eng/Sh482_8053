package service

import (
	"fmt"

	"enterpriselead/internal/domain"
)

func (s *Service) SaveAttachment(recordID, name, mediaType string, content []byte) (domain.Attachment, error) {
	if err := s.ensureStore(); err != nil {
		return domain.Attachment{}, err
	}
	if _, err := s.store.GetRecord(recordID); err != nil {
		return domain.Attachment{}, err
	}
	if name == "" || mediaType == "" {
		return domain.Attachment{}, fmt.Errorf("attachment name and media type are required")
	}
	attachment := domain.Attachment{ID: s.next("att"), RecordID: recordID, Name: name, MediaType: mediaType, Size: int64(len(content)), Content: append([]byte(nil), content...), CreatedAt: s.now()}
	if err := s.store.PutAttachment(attachment); err != nil {
		return domain.Attachment{}, err
	}
	return s.store.GetAttachment(attachment.ID)
}

func (s *Service) GetAttachment(id string) (domain.Attachment, error) {
	if err := s.ensureStore(); err != nil {
		return domain.Attachment{}, err
	}
	return s.store.GetAttachment(id)
}
