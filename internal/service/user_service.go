package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"secretary/alpha/internal/domain"
	"secretary/alpha/internal/utils"
)

type userService struct {
	repo domain.UserRepository
}

func NewUserService(repo domain.UserRepository) domain.UserService {
	return &userService{repo: repo}
}

func (s *userService) CreateUser(ctx context.Context, user *domain.User) error {
	// Check if username already exists
	if _, err := s.repo.FindByUsername(user.Username); err == nil {
		return errors.New("username already exists")
	}

	// Check if email already exists
	if _, err := s.repo.FindByEmail(user.Email); err == nil {
		return errors.New("email already exists")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return errors.New("failed to hash password")
	}
	user.Password = hashedPassword

	// Set default values
	user.ID = uuid.New().String()
	user.Role = "user"
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	// Save user to database
	return s.repo.Create(user)
}

func (s *userService) Authenticate(ctx context.Context, username, password string) (*domain.User, error) {
	// Find user by username
	user, err := s.repo.FindByUsername(username)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Validate password
	if !utils.CheckPasswordHash(password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

func (s *userService) GetByID(ctx context.Context, id string) (*domain.User, error) {
	return s.repo.FindByID(id)
}

func (s *userService) Update(ctx context.Context, user *domain.User) error {
	user.UpdatedAt = time.Now()
	return s.repo.Update(user)
}

func (s *userService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(id)
}
