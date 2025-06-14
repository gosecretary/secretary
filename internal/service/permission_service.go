package service

import (
	"context"

	"secretary/alpha/internal/domain"
)

type permissionService struct {
	repo domain.PermissionRepository
}

func NewPermissionService(repo domain.PermissionRepository) domain.PermissionService {
	return &permissionService{repo: repo}
}

func (s *permissionService) Create(ctx context.Context, permission *domain.Permission) error {
	return s.repo.Create(permission)
}

func (s *permissionService) GetByID(ctx context.Context, id string) (*domain.Permission, error) {
	return s.repo.FindByID(id)
}

func (s *permissionService) GetByUserID(ctx context.Context, userID string) ([]*domain.Permission, error) {
	return s.repo.FindByUserID(userID)
}

func (s *permissionService) GetByResourceID(ctx context.Context, resourceID string) ([]*domain.Permission, error) {
	return s.repo.FindByResourceID(resourceID)
}

func (s *permissionService) Update(ctx context.Context, permission *domain.Permission) error {
	return s.repo.Update(permission)
}

func (s *permissionService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(id)
}

func (s *permissionService) DeleteByUserID(ctx context.Context, userID string) error {
	return s.repo.DeleteByUserID(userID)
}

func (s *permissionService) DeleteByResourceID(ctx context.Context, resourceID string) error {
	return s.repo.DeleteByResourceID(resourceID)
}

func (s *permissionService) CreatePermission(ctx context.Context, permission *domain.Permission) error {
	return s.repo.Create(permission)
}

func (s *permissionService) DeletePermission(ctx context.Context, id string) error {
	return s.repo.Delete(id)
}

func (s *permissionService) GetPermission(ctx context.Context, id string) (*domain.Permission, error) {
	return s.repo.FindByID(id)
}

func (s *permissionService) GetPermissionByResourceID(ctx context.Context, resourceID string) ([]*domain.Permission, error) {
	return s.repo.FindByResourceID(resourceID)
}

func (s *permissionService) GetPermissionByUserID(ctx context.Context, userID string) ([]*domain.Permission, error) {
	return s.repo.FindByUserID(userID)
}

func (s *permissionService) ListPermissions(ctx context.Context) ([]*domain.Permission, error) {
	return s.repo.FindByUserID("") // TODO: Implement proper listing
}

func (s *permissionService) UpdatePermission(ctx context.Context, permission *domain.Permission) error {
	return s.repo.Update(permission)
}
