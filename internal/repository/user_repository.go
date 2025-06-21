package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"

	"secretary/alpha/internal/domain"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) domain.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *domain.User) error {
	// Set default values if not provided
	if user.ID == "" {
		user.ID = uuid.New().String()
	}
	if user.CreatedAt.IsZero() {
		user.CreatedAt = time.Now()
	}
	if user.UpdatedAt.IsZero() {
		user.UpdatedAt = time.Now()
	}

	query := `
		INSERT INTO users (id, username, email, password, name, role, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.Exec(query,
		user.ID,
		user.Username,
		user.Email,
		user.Password,
		user.Name,
		user.Role,
		user.CreatedAt,
		user.UpdatedAt,
	)
	return err
}

func (r *userRepository) FindByID(id string) (*domain.User, error) {
	query := `
		SELECT id, username, email, password, name, role, created_at, updated_at
		FROM users
		WHERE id = ?
	`
	user := &domain.User{}
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Name,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	return user, err
}

func (r *userRepository) FindByUsername(username string) (*domain.User, error) {
	query := `
		SELECT id, username, email, password, name, role, created_at, updated_at
		FROM users
		WHERE username = ?
	`
	user := &domain.User{}
	err := r.db.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Name,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	return user, err
}

func (r *userRepository) FindByEmail(email string) (*domain.User, error) {
	query := `
		SELECT id, username, email, password, name, role, created_at, updated_at
		FROM users
		WHERE email = ?
	`
	user := &domain.User{}
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Name,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	return user, err
}

func (r *userRepository) Update(user *domain.User) error {
	// Check if user exists
	_, err := r.FindByID(user.ID)
	if err != nil {
		return err
	}

	// Update timestamp
	user.UpdatedAt = time.Now()

	query := `
		UPDATE users
		SET username = ?, email = ?, password = ?, name = ?, role = ?, updated_at = ?
		WHERE id = ?
	`
	result, err := r.db.Exec(query,
		user.Username,
		user.Email,
		user.Password,
		user.Name,
		user.Role,
		user.UpdatedAt,
		user.ID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

func (r *userRepository) Delete(id string) error {
	// Check if user exists
	_, err := r.FindByID(id)
	if err != nil {
		return err
	}

	query := `DELETE FROM users WHERE id = ?`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}
