package service

import (
	"errors"

	"secretary/alpha/internal/domain"

	"github.com/google/uuid"
)

type userService struct {
	repo domain.UserRepository
}

func NewUserService(repo domain.UserRepository) domain.UserService {
	return &userService{repo: repo}
}

func (s *userService) Register(username, email, password string) (*domain.User, error) {
	// Check if username already exists
	if _, err := s.repo.FindByUsername(username); err == nil {
		return nil, errors.New("username already exists")
	}

	// Check if email already exists
	if _, err := s.repo.FindByEmail(email); err == nil {
		return nil, errors.New("email already exists")
	}

	// Create new user
	user, err := domain.NewUser(username, email, password)
	if err != nil {
		return nil, err
	}

	// Save user to database
	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) Login(username, password string) (*domain.User, error) {
	// Find user by username
	user, err := s.repo.FindByUsername(username)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Validate password
	if !user.ValidatePassword(password) {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

func (s *userService) GetByID(id uuid.UUID) (*domain.User, error) {
	return s.repo.FindByID(id)
}

func (s *userService) Update(user *domain.User) error {
	// Check if user exists
	existingUser, err := s.repo.FindByID(user.ID)
	if err != nil {
		return err
	}

	// Update user
	return s.repo.Update(user)
}

func (s *userService) Delete(id uuid.UUID) error {
	// Check if user exists
	if _, err := s.repo.FindByID(id); err != nil {
		return err
	}

	return s.repo.Delete(id)
} 