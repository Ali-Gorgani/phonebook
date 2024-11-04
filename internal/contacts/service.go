package contacts

import (
	"database/sql"
)

type Service struct {
	repo IRepository
}

func NewService(db *sql.DB) *Service {
	return &Service{
		repo: NewRepository(db),
	}
}

func (s *Service) CreateContact(contact *Contact) (int, error) {
	return s.repo.CreateContact(contact)
}

func (s *Service) UpdateContact(contact *Contact) error {
	return s.repo.UpdateContact(contact)
}

func (s *Service) SearchContacts(query string) ([]Contact, error) {
	return s.repo.SearchContacts(query)
}
