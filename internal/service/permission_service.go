package service

import (
	"time"

	"secretary/alpha/internal/domain"

	"github.com/google/uuid"
)

type permissionService struct {
	repo domain.PermissionRepository
}

func NewPermissionService(repo domain.PermissionRepository) domain.PermissionService {
	return &permissionService{repo: repo}
}

func (s *permissionService) Create(userID, resourceID uuid.UUID, action string) (*domain.Permission, error) {
	permission := domain.NewPermission(userID, resourceID, action)
	if err := s.repo.Create(permission); err != nil {
		return nil, err
	}
	return permission, nil
}

func (s *permissionService) GetByID(id uuid.UUID) (*domain.Permission, error) {
	return s.repo.FindByID(id)
}

func (s *permissionService) GetByUserID(userID uuid.UUID) ([]*domain.Permission, error) {
	return s.repo.FindByUserID(userID)
}

func (s *permissionService) GetByResourceID(resourceID uuid.UUID) ([]*domain.Permission, error) {
	return s.repo.FindByResourceID(resourceID)
}

func (s *permissionService) Update(id uuid.UUID, action string) (*domain.Permission, error) {
	permission, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	permission.Action = action
	permission.UpdatedAt = time.Now()

	if err := s.repo.Update(permission); err != nil {
		return nil, err
	}

	return permission, nil
}

func (s *permissionService) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}

func (s *permissionService) DeleteByUserID(userID uuid.UUID) error {
	return s.repo.DeleteByUserID(userID)
}

func (s *permissionService) DeleteByResourceID(resourceID uuid.UUID) error {
	return s.repo.DeleteByResourceID(resourceID)
} 