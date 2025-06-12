package repository

import (
	"database/sql"
	"errors"

	"secretary/alpha/internal/domain"
)

type credentialRepository struct {
	db *sql.DB
}

func NewCredentialRepository(db *sql.DB) domain.CredentialRepository {
	return &credentialRepository{db: db}
}

func (r *credentialRepository) Create(credential *domain.Credential) error {
	query := `
		INSERT INTO credentials (id, resource_id, username, password, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.Exec(query,
		credential.ID,
		credential.ResourceID,
		credential.Username,
		credential.Password,
		credential.CreatedAt,
		credential.UpdatedAt,
	)
	return err
}

func (r *credentialRepository) FindByID(id string) (*domain.Credential, error) {
	query := `
		SELECT id, resource_id, username, password, created_at, updated_at
		FROM credentials
		WHERE id = ?
	`
	credential := &domain.Credential{}
	err := r.db.QueryRow(query, id).Scan(
		&credential.ID,
		&credential.ResourceID,
		&credential.Username,
		&credential.Password,
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
		SELECT id, resource_id, username, password, created_at, updated_at
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
			&credential.Username,
			&credential.Password,
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
	query := `
		UPDATE credentials
		SET username = ?, password = ?, updated_at = ?
		WHERE id = ?
	`
	_, err := r.db.Exec(query,
		credential.Username,
		credential.Password,
		credential.UpdatedAt,
		credential.ID,
	)
	return err
}

func (r *credentialRepository) Delete(id string) error {
	query := `DELETE FROM credentials WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}
