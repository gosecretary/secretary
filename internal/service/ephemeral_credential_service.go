package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/google/uuid"

	"secretary/alpha/internal/domain"
)

const (
	defaultCredentialLifetime = 8 * time.Hour
	usernamePrefix            = "sec"
	tokenLength               = 48
	passwordLength            = 24
	passwordCharset           = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()-_=+,.?;:[]{}|"
	maxTokenRetries           = 10
)

type ephemeralCredentialService struct {
	repo domain.EphemeralCredentialRepository
}

func NewEphemeralCredentialService(repo domain.EphemeralCredentialRepository) domain.EphemeralCredentialService {
	return &ephemeralCredentialService{repo: repo}
}

func (s *ephemeralCredentialService) Generate(ctx context.Context, userID string, resourceID string, duration time.Duration) (*domain.EphemeralCredential, error) {
	if duration == 0 {
		duration = defaultCredentialLifetime
	}

	// Generate secure random elements
	username, err := s.generateUsername()
	if err != nil {
		return nil, err
	}

	password, err := s.generatePassword(passwordLength)
	if err != nil {
		return nil, err
	}

	token, err := s.generateToken(tokenLength)
	if err != nil {
		return nil, err
	}

	// Create credential
	credential := &domain.EphemeralCredential{
		ID:         uuid.New().String(),
		UserID:     userID,
		ResourceID: resourceID,
		Username:   username,
		Password:   password,
		Token:      token,
		ExpiresAt:  time.Now().Add(duration),
		CreatedAt:  time.Now(),
	}

	err = s.repo.Create(credential)
	if err != nil {
		return nil, fmt.Errorf("failed to store credential: %w", err)
	}

	return credential, nil
}

func (s *ephemeralCredentialService) GetByID(ctx context.Context, id string) (*domain.EphemeralCredential, error) {
	return s.repo.FindByID(id)
}

func (s *ephemeralCredentialService) GetByToken(ctx context.Context, token string) (*domain.EphemeralCredential, error) {
	return s.repo.FindByToken(token)
}

func (s *ephemeralCredentialService) MarkAsUsed(ctx context.Context, id string) error {
	// Get the current credential
	credential, err := s.repo.FindByID(id)
	if err != nil {
		return fmt.Errorf("failed to find credential: %w", err)
	}

	// Check if the credential is expired
	if credential.ExpiresAt.Before(time.Now()) {
		return fmt.Errorf("credential has expired")
	}

	// Check if the credential has already been used
	if !credential.UsedAt.IsZero() {
		return fmt.Errorf("credential has already been used")
	}

	// Mark as used
	credential.UsedAt = time.Now()

	return s.repo.Update(credential)
}

func (s *ephemeralCredentialService) RevokeByUserID(ctx context.Context, userID string) error {
	return s.repo.DeleteByUserID(userID)
}

func (s *ephemeralCredentialService) RevokeByResourceID(ctx context.Context, resourceID string) error {
	return s.repo.DeleteByResourceID(resourceID)
}

func (s *ephemeralCredentialService) Create(ctx context.Context, credential *domain.EphemeralCredential) (*domain.EphemeralCredential, error) {
	if credential.ID == "" {
		credential.ID = uuid.New().String()
	}
	if credential.CreatedAt.IsZero() {
		credential.CreatedAt = time.Now()
	}
	if credential.ExpiresAt.IsZero() {
		credential.ExpiresAt = time.Now().Add(defaultCredentialLifetime)
	}

	// Generate unique token with retry mechanism
	var token string
	var err error
	for i := 0; i < maxTokenRetries; i++ {
		token, err = s.generateToken(tokenLength)
		if err != nil {
			continue
		}

		// Check if token already exists
		existingCred, err := s.repo.FindByToken(token)
		if err != nil && existingCred == nil {
			// Token doesn't exist, we can use it
			credential.Token = token
			break
		}
	}

	if credential.Token == "" {
		return nil, fmt.Errorf("failed to generate unique token after %d attempts", maxTokenRetries)
	}

	if err := s.repo.Create(credential); err != nil {
		return nil, err
	}
	return credential, nil
}

func (s *ephemeralCredentialService) GetEphemeralCredential(ctx context.Context, id string) (*domain.EphemeralCredential, error) {
	return s.repo.FindByID(id)
}

func (s *ephemeralCredentialService) DeleteEphemeralCredential(ctx context.Context, id string) error {
	return s.repo.DeleteByUserID(id)
}

func (s *ephemeralCredentialService) List(ctx context.Context) ([]*domain.EphemeralCredential, error) {
	return s.repo.FindByUserID("")
}

func (s *ephemeralCredentialService) MarkAsUsedEphemeralCredential(ctx context.Context, id string) error {
	credential, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	credential.UsedAt = time.Now()
	return s.repo.Update(credential)
}

// Helper methods for secure random generation

func (s *ephemeralCredentialService) generateUsername() (string, error) {
	// Create a unique username with a UUID
	uid, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s_%s", usernamePrefix, strings.ReplaceAll(uid.String(), "-", "")[:12]), nil
}

func (s *ephemeralCredentialService) generatePassword(length int) (string, error) {
	password := make([]byte, length)
	charsetLength := big.NewInt(int64(len(passwordCharset)))

	for i := 0; i < length; i++ {
		idx, err := rand.Int(rand.Reader, charsetLength)
		if err != nil {
			return "", err
		}
		password[i] = passwordCharset[idx.Int64()]
	}

	return string(password), nil
}

func (s *ephemeralCredentialService) generateToken(length int) (string, error) {
	tokenBytes := make([]byte, length)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(tokenBytes), nil
}
