package service

import (
	"secretary/alpha/internal/domain"

	"github.com/google/uuid"
)

type credentialService struct {
	repo domain.CredentialRepository
}

func NewCredentialService(repo domain.CredentialRepository) domain.CredentialService {
	return &credentialService{repo: repo}
}

func (s *credentialService) Create(resourceID uuid.UUID, username, password string) (*domain.Credential, error) {
	credential, err := domain.NewCredential(resourceID, username, password)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Create(credential); err != nil {
		return nil, err
	}
	return credential, nil
}

func (s *credentialService) GetByID(id uuid.UUID) (*domain.Credential, error) {
	return s.repo.FindByID(id)
}

func (s *credentialService) GetByResourceID(resourceID uuid.UUID) ([]*domain.Credential, error) {
	return s.repo.FindByResourceID(resourceID)
}

func (s *credentialService) Update(id uuid.UUID, username, password string) (*domain.Credential, error) {
	credential, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	credential.Username = username
	if err := credential.UpdatePassword(password); err != nil {
		return nil, err
	}

	if err := s.repo.Update(credential); err != nil {
		return nil, err
	}

	return credential, nil
}

func (s *credentialService) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}
