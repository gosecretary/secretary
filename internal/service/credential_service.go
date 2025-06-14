package service

import (
	"context"
	"secretary/alpha/internal/domain"
)

type credentialService struct {
	repo domain.CredentialRepository
}

func NewCredentialService(repo domain.CredentialRepository) domain.CredentialService {
	return &credentialService{repo: repo}
}

func (s *credentialService) Create(ctx context.Context, credential *domain.Credential) error {
	return s.repo.Create(credential)
}

func (s *credentialService) GetByID(ctx context.Context, id string) (*domain.Credential, error) {
	return s.repo.FindByID(id)
}

func (s *credentialService) GetByResourceID(ctx context.Context, resourceID string) ([]*domain.Credential, error) {
	return s.repo.FindByResourceID(resourceID)
}

func (s *credentialService) Update(ctx context.Context, credential *domain.Credential) error {
	return s.repo.Update(credential)
}

func (s *credentialService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(id)
}

func (s *credentialService) CreateCredential(ctx context.Context, credential *domain.Credential) error {
	return s.repo.Create(credential)
}

func (s *credentialService) ListCredentials(ctx context.Context) ([]*domain.Credential, error) {
	return s.repo.FindByResourceID("") // TODO: Implement proper listing
}

func (s *credentialService) GetCredential(ctx context.Context, id string) (*domain.Credential, error) {
	return s.repo.FindByID(id)
}

func (s *credentialService) UpdateCredential(ctx context.Context, credential *domain.Credential) error {
	return s.repo.Update(credential)
}

func (s *credentialService) DeleteCredential(ctx context.Context, id string) error {
	return s.repo.Delete(id)
}

func (s *credentialService) RotateCredential(ctx context.Context, id string) error {
	// TODO: Implement credential rotation logic
	return nil
}

func (s *credentialService) GetCredentialByResourceID(ctx context.Context, resourceID string) ([]*domain.Credential, error) {
	return s.repo.FindByResourceID(resourceID)
}
