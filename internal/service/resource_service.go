package service

import (
	"time"

	"secretary/alpha/internal/domain"

	"github.com/google/uuid"
)

type resourceService struct {
	repo domain.ResourceRepository
}

func NewResourceService(repo domain.ResourceRepository) domain.ResourceService {
	return &resourceService{repo: repo}
}

func (s *resourceService) Create(name, description string) (*domain.Resource, error) {
	resource := domain.NewResource(name, description)
	if err := s.repo.Create(resource); err != nil {
		return nil, err
	}
	return resource, nil
}

func (s *resourceService) GetByID(id uuid.UUID) (*domain.Resource, error) {
	return s.repo.FindByID(id)
}

func (s *resourceService) GetAll() ([]*domain.Resource, error) {
	return s.repo.FindAll()
}

func (s *resourceService) Update(id uuid.UUID, name, description string) (*domain.Resource, error) {
	resource, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	resource.Name = name
	resource.Description = description
	resource.UpdatedAt = time.Now()

	if err := s.repo.Update(resource); err != nil {
		return nil, err
	}

	return resource, nil
}

func (s *resourceService) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
} 