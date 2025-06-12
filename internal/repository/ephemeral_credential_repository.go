package repository

import (
	"database/sql"
	"errors"
	"time"

	"secretary/alpha/internal/domain"
)

type ephemeralCredentialRepository struct {
	db *sql.DB
}

func NewEphemeralCredentialRepository(db *sql.DB) domain.EphemeralCredentialRepository {
	return &ephemeralCredentialRepository{db: db}
}

func (r *ephemeralCredentialRepository) Create(credential *domain.EphemeralCredential) error {
	query := `
		INSERT INTO ephemeral_credentials (
			id, user_id, resource_id, username, password, token, expires_at, created_at, used_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.Exec(query,
		credential.ID,
		credential.UserID,
		credential.ResourceID,
		credential.Username,
		credential.Password,
		credential.Token,
		credential.ExpiresAt,
		credential.CreatedAt,
		credential.UsedAt,
	)
	return err
}

func (r *ephemeralCredentialRepository) FindByID(id string) (*domain.EphemeralCredential, error) {
	query := `
		SELECT id, user_id, resource_id, username, password, token, expires_at, created_at, used_at
		FROM ephemeral_credentials
		WHERE id = ?
	`
	credential := &domain.EphemeralCredential{}
	var usedAt sql.NullTime

	err := r.db.QueryRow(query, id).Scan(
		&credential.ID,
		&credential.UserID,
		&credential.ResourceID,
		&credential.Username,
		&credential.Password,
		&credential.Token,
		&credential.ExpiresAt,
		&credential.CreatedAt,
		&usedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("ephemeral credential not found")
	}

	if usedAt.Valid {
		credential.UsedAt = usedAt.Time
	}

	return credential, err
}

func (r *ephemeralCredentialRepository) FindByToken(token string) (*domain.EphemeralCredential, error) {
	query := `
		SELECT id, user_id, resource_id, username, password, token, expires_at, created_at, used_at
		FROM ephemeral_credentials
		WHERE token = ?
	`
	credential := &domain.EphemeralCredential{}
	var usedAt sql.NullTime

	err := r.db.QueryRow(query, token).Scan(
		&credential.ID,
		&credential.UserID,
		&credential.ResourceID,
		&credential.Username,
		&credential.Password,
		&credential.Token,
		&credential.ExpiresAt,
		&credential.CreatedAt,
		&usedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("ephemeral credential not found")
	}

	if usedAt.Valid {
		credential.UsedAt = usedAt.Time
	}

	return credential, err
}

func (r *ephemeralCredentialRepository) FindByUserID(userID string) ([]*domain.EphemeralCredential, error) {
	query := `
		SELECT id, user_id, resource_id, username, password, token, expires_at, created_at, used_at
		FROM ephemeral_credentials
		WHERE user_id = ? AND expires_at > ?
		ORDER BY created_at DESC
	`
	return r.queryEphemeralCredentials(query, userID, time.Now())
}

func (r *ephemeralCredentialRepository) Update(credential *domain.EphemeralCredential) error {
	query := `
		UPDATE ephemeral_credentials
		SET used_at = ?
		WHERE id = ?
	`
	_, err := r.db.Exec(query, credential.UsedAt, credential.ID)
	return err
}

func (r *ephemeralCredentialRepository) DeleteExpired() error {
	query := `
		DELETE FROM ephemeral_credentials
		WHERE expires_at < ?
	`
	_, err := r.db.Exec(query, time.Now())
	return err
}

func (r *ephemeralCredentialRepository) DeleteByUserID(userID string) error {
	query := `
		DELETE FROM ephemeral_credentials
		WHERE user_id = ?
	`
	_, err := r.db.Exec(query, userID)
	return err
}

func (r *ephemeralCredentialRepository) DeleteByResourceID(resourceID string) error {
	query := `
		DELETE FROM ephemeral_credentials
		WHERE resource_id = ?
	`
	_, err := r.db.Exec(query, resourceID)
	return err
}

func (r *ephemeralCredentialRepository) queryEphemeralCredentials(query string, args ...interface{}) ([]*domain.EphemeralCredential, error) {
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var credentials []*domain.EphemeralCredential
	for rows.Next() {
		credential := &domain.EphemeralCredential{}
		var usedAt sql.NullTime

		err := rows.Scan(
			&credential.ID,
			&credential.UserID,
			&credential.ResourceID,
			&credential.Username,
			&credential.Password,
			&credential.Token,
			&credential.ExpiresAt,
			&credential.CreatedAt,
			&usedAt,
		)
		if err != nil {
			return nil, err
		}

		if usedAt.Valid {
			credential.UsedAt = usedAt.Time
		}

		credentials = append(credentials, credential)
	}
	return credentials, nil
}
