package repository

import (
	"database/sql"
	"errors"

	"secretary/alpha/internal/domain"

	"github.com/google/uuid"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) domain.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *domain.User) error {
	query := `
		INSERT INTO users (id, username, email, password, role, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.Exec(query,
		user.ID,
		user.Username,
		user.Email,
		user.Password,
		user.Role,
		user.CreatedAt,
		user.UpdatedAt,
	)
	return err
}

func (r *userRepository) FindByID(id uuid.UUID) (*domain.User, error) {
	query := `
		SELECT id, username, email, password, role, created_at, updated_at
		FROM users
		WHERE id = ?
	`
	user := &domain.User{}
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
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
		SELECT id, username, email, password, role, created_at, updated_at
		FROM users
		WHERE username = ?
	`
	user := &domain.User{}
	err := r.db.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
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
		SELECT id, username, email, password, role, created_at, updated_at
		FROM users
		WHERE email = ?
	`
	user := &domain.User{}
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
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
	query := `
		UPDATE users
		SET username = ?, email = ?, password = ?, role = ?, updated_at = ?
		WHERE id = ?
	`
	_, err := r.db.Exec(query,
		user.Username,
		user.Email,
		user.Password,
		user.Role,
		user.UpdatedAt,
		user.ID,
	)
	return err
}

func (r *userRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
} 