package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"

	"secretary/alpha/internal/domain"
)

type credentialRepository struct {
	db *sql.DB
}

func NewCredentialRepository(db *sql.DB) domain.CredentialRepository {
	return &credentialRepository{db: db}
}

func (r *credentialRepository) Create(credential *domain.Credential) error {
	// Generate UUID if not provided
	if credential.ID == "" {
		credential.ID = uuid.New().String()
	}

	// Set timestamps
	credential.CreatedAt = time.Now()
	credential.UpdatedAt = time.Now()

	query := `
		INSERT INTO credentials (id, resource_id, type, secret, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.Exec(query,
		credential.ID,
		credential.ResourceID,
		credential.Type,
		credential.Secret,
		credential.CreatedAt,
		credential.UpdatedAt,
	)
	return err
}

func (r *credentialRepository) FindByID(id string) (*domain.Credential, error) {
	query := `
		SELECT id, resource_id, type, secret, created_at, updated_at
		FROM credentials
		WHERE id = ?
	`
	credential := &domain.Credential{}
	err := r.db.QueryRow(query, id).Scan(
		&credential.ID,
		&credential.ResourceID,
		&credential.Type,
		&credential.Secret,
		&credential.CreatedAt,
		&credential.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New("credential not found")
	}
	return credential, err
}

func (r *credentialRepository) FindByResourceID(resourceID string) ([]*domain.Credential, error) {
	query := `
		SELECT id, resource_id, type, secret, created_at, updated_at
		FROM credentials
		WHERE resource_id = ?
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query, resourceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var credentials []*domain.Credential
	for rows.Next() {
		credential := &domain.Credential{}
		err := rows.Scan(
			&credential.ID,
			&credential.ResourceID,
			&credential.Type,
			&credential.Secret,
			&credential.CreatedAt,
			&credential.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		credentials = append(credentials, credential)
	}
	return credentials, nil
}

func (r *credentialRepository) Update(credential *domain.Credential) error {
	credential.UpdatedAt = time.Now()
	query := `
		UPDATE credentials
		SET resource_id = ?, type = ?, secret = ?, updated_at = ?
		WHERE id = ?
	`
	result, err := r.db.Exec(query,
		credential.ResourceID,
		credential.Type,
		credential.Secret,
		credential.UpdatedAt,
		credential.ID,
	)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("credential not found")
	}
	return nil
}

func (r *credentialRepository) Delete(id string) error {
	query := `DELETE FROM credentials WHERE id = ?`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("credential not found")
	}
	return nil
}
