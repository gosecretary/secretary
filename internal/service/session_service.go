package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"secretary/alpha/internal/domain"
)

type sessionService struct {
	repo domain.SessionRepository
}

func NewSessionService(repo domain.SessionRepository) domain.SessionService {
	return &sessionService{repo: repo}
}

func (s *sessionService) Create(ctx context.Context, session *domain.Session) error {
	// Generate UUID if not provided
	if session.ID == "" {
		session.ID = uuid.New().String()
	}

	// Set initial values
	session.StartTime = time.Now()
	session.Status = "active"
	session.CreatedAt = time.Now()
	session.UpdatedAt = time.Now()

	return s.repo.Create(session)
}

func (s *sessionService) GetByID(ctx context.Context, id string) (*domain.Session, error) {
	return s.repo.FindByID(id)
}

func (s *sessionService) GetByUserID(ctx context.Context, userID string) ([]*domain.Session, error) {
	return s.repo.FindByUserID(userID)
}

func (s *sessionService) GetByResourceID(ctx context.Context, resourceID string) ([]*domain.Session, error) {
	return s.repo.FindByResourceID(resourceID)
}

func (s *sessionService) GetActive(ctx context.Context) ([]*domain.Session, error) {
	return s.repo.FindActive()
}

func (s *sessionService) Update(ctx context.Context, session *domain.Session) error {
	// Get the current session to validate state
	currentSession, err := s.repo.FindByID(session.ID)
	if err != nil {
		return fmt.Errorf("failed to find session: %w", err)
	}

	// Only allow updating active sessions
	if currentSession.Status != "active" {
		return fmt.Errorf("cannot update a non-active session")
	}

	// Update the timestamp
	session.UpdatedAt = time.Now()

	return s.repo.Update(session)
}

func (s *sessionService) Terminate(ctx context.Context, id string) error {
	// Get the current session
	session, err := s.repo.FindByID(id)
	if err != nil {
		return fmt.Errorf("failed to find session: %w", err)
	}

	// Only terminate active sessions
	if session.Status != "active" {
		return fmt.Errorf("session is not active")
	}

	// Update session state
	session.Status = "terminated"
	session.EndTime = time.Now()
	session.UpdatedAt = time.Now()

	return s.repo.Update(session)
}
