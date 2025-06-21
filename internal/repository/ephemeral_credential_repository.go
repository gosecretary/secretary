package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"

	"secretary/alpha/internal/domain"
)

type ephemeralCredentialRepository struct {
	db *sql.DB
}

func NewEphemeralCredentialRepository(db *sql.DB) domain.EphemeralCredentialRepository {
	return &ephemeralCredentialRepository{db: db}
}

func (r *ephemeralCredentialRepository) Create(credential *domain.EphemeralCredential) error {
	// Set default values if not provided
	if credential.ID == "" {
		credential.ID = uuid.New().String()
	}
	if credential.CreatedAt.IsZero() {
		credential.CreatedAt = time.Now()
	}

	query := `
		INSERT INTO ephemeral_credentials (
			id, user_id, resource_id, username, password, token, expires_at, created_at, used_at, duration, used
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
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
		credential.Duration,
		credential.Used,
	)
	return err
}

func (r *ephemeralCredentialRepository) FindByID(id string) (*domain.EphemeralCredential, error) {
	query := `
		SELECT id, user_id, resource_id, username, password, token, expires_at, created_at, used_at, duration, used
		FROM ephemeral_credentials
		WHERE id = ?
	`
	credential := &domain.EphemeralCredential{}
	var usedAt sql.NullTime
	var used bool

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
		&credential.Duration,
		&used,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("ephemeral credential not found")
	}

	if usedAt.Valid {
		credential.UsedAt = usedAt.Time
	}
	credential.Used = used

	return credential, err
}

func (r *ephemeralCredentialRepository) FindByToken(token string) (*domain.EphemeralCredential, error) {
	query := `
		SELECT id, user_id, resource_id, username, password, token, expires_at, created_at, used_at, duration, used
		FROM ephemeral_credentials
		WHERE token = ?
	`
	credential := &domain.EphemeralCredential{}
	var usedAt sql.NullTime
	var used bool

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
		&credential.Duration,
		&used,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("ephemeral credential not found")
	}

	if usedAt.Valid {
		credential.UsedAt = usedAt.Time
	}
	credential.Used = used

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
		SET used = ?, used_at = ?
		WHERE id = ?
	`
	result, err := r.db.Exec(query, credential.Used, credential.UsedAt, credential.ID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("ephemeral credential not found")
	}
	// Update in-memory struct
	credential.Used = true
	if credential.UsedAt.IsZero() {
		credential.UsedAt = time.Now()
	}
	return nil
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
	result, err := r.db.Exec(query, userID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("ephemeral credential not found for user")
	}
	return nil
}

func (r *ephemeralCredentialRepository) DeleteByResourceID(resourceID string) error {
	query := `
		DELETE FROM ephemeral_credentials
		WHERE resource_id = ?
	`
	result, err := r.db.Exec(query, resourceID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("ephemeral credential not found for resource")
	}
	return nil
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
