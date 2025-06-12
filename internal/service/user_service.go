package service

import (
	"context"
	"errors"

	"secretary/alpha/internal/domain"
)

type userService struct {
	repo domain.UserRepository
}

func NewUserService(repo domain.UserRepository) domain.UserService {
	return &userService{repo: repo}
}

func (s *userService) Register(ctx context.Context, user *domain.User) error {
	// Check if username already exists
	if _, err := s.repo.FindByUsername(user.Name); err == nil {
		return errors.New("username already exists")
	}

	// Check if email already exists
	if _, err := s.repo.FindByEmail(user.Email); err == nil {
		return errors.New("email already exists")
	}

	// Save user to database
	return s.repo.Create(user)
}

func (s *userService) Login(ctx context.Context, email, password string) (string, error) {
	// Find user by email
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// Validate password
	if !user.ValidatePassword(password) {
		return "", errors.New("invalid credentials")
	}

	return user.ID, nil
}

func (s *userService) GetByID(ctx context.Context, id string) (*domain.User, error) {
	return s.repo.FindByID(id)
}

func (s *userService) Update(ctx context.Context, user *domain.User) error {
	return s.repo.Update(user)
}

func (s *userService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(id)
}
