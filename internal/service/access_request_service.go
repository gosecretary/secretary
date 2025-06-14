package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"secretary/alpha/internal/domain"
)

type accessRequestService struct {
	repo domain.AccessRequestRepository
}

func NewAccessRequestService(repo domain.AccessRequestRepository) domain.AccessRequestService {
	return &accessRequestService{repo: repo}
}

func (s *accessRequestService) Create(ctx context.Context, request *domain.AccessRequest) error {
	// Generate UUID if not provided
	if request.ID == "" {
		request.ID = uuid.New().String()
	}

	// Set default values
	request.Status = "pending"
	request.RequestedAt = time.Now()
	request.CreatedAt = time.Now()
	request.UpdatedAt = time.Now()

	return s.repo.Create(request)
}

func (s *accessRequestService) GetByID(ctx context.Context, id string) (*domain.AccessRequest, error) {
	return s.repo.FindByID(id)
}

func (s *accessRequestService) GetByUserID(ctx context.Context, userID string) ([]*domain.AccessRequest, error) {
	return s.repo.FindByUserID(userID)
}

func (s *accessRequestService) GetByResourceID(ctx context.Context, resourceID string) ([]*domain.AccessRequest, error) {
	return s.repo.FindByResourceID(resourceID)
}

func (s *accessRequestService) GetPending(ctx context.Context) ([]*domain.AccessRequest, error) {
	return s.repo.FindByStatus("pending")
}

func (s *accessRequestService) Approve(ctx context.Context, id string, reviewerID string, notes string, expiresAt time.Time) error {
	// Get the current request
	request, err := s.repo.FindByID(id)
	if err != nil {
		return fmt.Errorf("failed to find access request: %w", err)
	}

	// Only allow approving pending requests
	if request.Status != "pending" {
		return errors.New("only pending requests can be approved")
	}

	// Update request with approval details
	request.Status = "approved"
	request.ReviewerID = reviewerID
	request.ReviewNotes = notes
	request.ReviewedAt = time.Now()
	request.ExpiresAt = expiresAt
	request.UpdatedAt = time.Now()

	return s.repo.Update(request)
}

func (s *accessRequestService) Deny(ctx context.Context, id string, reviewerID string, notes string) error {
	// Get the current request
	request, err := s.repo.FindByID(id)
	if err != nil {
		return fmt.Errorf("failed to find access request: %w", err)
	}

	// Only allow denying pending requests
	if request.Status != "pending" {
		return errors.New("only pending requests can be denied")
	}

	// Update request with denial details
	request.Status = "denied"
	request.ReviewerID = reviewerID
	request.ReviewNotes = notes
	request.ReviewedAt = time.Now()
	request.UpdatedAt = time.Now()

	return s.repo.Update(request)
}

func (s *accessRequestService) CreateAccessRequest(ctx context.Context, request *domain.AccessRequest) error {
	// Generate UUID if not provided
	if request.ID == "" {
		request.ID = uuid.New().String()
	}

	// Set default values
	request.Status = "pending"
	request.RequestedAt = time.Now()
	request.CreatedAt = time.Now()
	request.UpdatedAt = time.Now()

	return s.repo.Create(request)
}

func (s *accessRequestService) GetAccessRequest(ctx context.Context, id string) (*domain.AccessRequest, error) {
	return s.repo.FindByID(id)
}

func (s *accessRequestService) ListAccessRequests(ctx context.Context) ([]*domain.AccessRequest, error) {
	return s.repo.FindByStatus("pending")
}

func (s *accessRequestService) UpdateAccessRequest(ctx context.Context, request *domain.AccessRequest) error {
	request.UpdatedAt = time.Now()
	return s.repo.Update(request)
}

func (s *accessRequestService) GetAccessRequestByResourceID(ctx context.Context, resourceID string) ([]*domain.AccessRequest, error) {
	return s.repo.FindByResourceID(resourceID)
}

func (s *accessRequestService) GetAccessRequestByUserID(ctx context.Context, userID string) ([]*domain.AccessRequest, error) {
	return s.repo.FindByUserID(userID)
}

func (s *accessRequestService) GetPendingAccessRequests(ctx context.Context) ([]*domain.AccessRequest, error) {
	return s.repo.FindByStatus("pending")
}
