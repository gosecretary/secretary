package service

import (
	"context"

	"secretary/alpha/internal/domain"
)

type resourceService struct {
	repo domain.ResourceRepository
}

func NewResourceService(repo domain.ResourceRepository) domain.ResourceService {
	return &resourceService{repo: repo}
}

func (s *resourceService) Create(ctx context.Context, resource *domain.Resource) error {
	return s.repo.Create(resource)
}

func (s *resourceService) GetByID(ctx context.Context, id string) (*domain.Resource, error) {
	return s.repo.FindByID(id)
}

func (s *resourceService) GetAll(ctx context.Context) ([]*domain.Resource, error) {
	return s.repo.FindAll()
}

func (s *resourceService) Update(ctx context.Context, resource *domain.Resource) error {
	return s.repo.Update(resource)
}

func (s *resourceService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(id)
}
