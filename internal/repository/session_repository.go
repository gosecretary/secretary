package repository

import (
	"database/sql"
	"errors"
	"time"

	"secretary/alpha/internal/domain"
)

type sessionRepository struct {
	db *sql.DB
}

func NewSessionRepository(db *sql.DB) domain.SessionRepository {
	return &sessionRepository{db: db}
}

func (r *sessionRepository) Create(session *domain.Session) error {
	query := `
		INSERT INTO sessions (
			id, user_id, resource_id, start_time, end_time, status, 
			client_ip, client_metadata, audit_path, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.Exec(query,
		session.ID,
		session.UserID,
		session.ResourceID,
		session.StartTime,
		session.EndTime,
		session.Status,
		session.ClientIP,
		session.ClientMetadata,
		session.AuditPath,
		session.CreatedAt,
		session.UpdatedAt,
	)
	return err
}

func (r *sessionRepository) FindByID(id string) (*domain.Session, error) {
	query := `
		SELECT id, user_id, resource_id, start_time, end_time, status, 
			client_ip, client_metadata, audit_path, created_at, updated_at
		FROM sessions
		WHERE id = ?
	`
	session := &domain.Session{}
	var endTime sql.NullTime

	err := r.db.QueryRow(query, id).Scan(
		&session.ID,
		&session.UserID,
		&session.ResourceID,
		&session.StartTime,
		&endTime,
		&session.Status,
		&session.ClientIP,
		&session.ClientMetadata,
		&session.AuditPath,
		&session.CreatedAt,
		&session.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("session not found")
	}

	if endTime.Valid {
		session.EndTime = endTime.Time
	}

	return session, err
}

func (r *sessionRepository) FindByUserID(userID string) ([]*domain.Session, error) {
	query := `
		SELECT id, user_id, resource_id, start_time, end_time, status, 
			client_ip, client_metadata, audit_path, created_at, updated_at
		FROM sessions
		WHERE user_id = ?
		ORDER BY created_at DESC
	`
	return r.querySessions(query, userID)
}

func (r *sessionRepository) FindByResourceID(resourceID string) ([]*domain.Session, error) {
	query := `
		SELECT id, user_id, resource_id, start_time, end_time, status, 
			client_ip, client_metadata, audit_path, created_at, updated_at
		FROM sessions
		WHERE resource_id = ?
		ORDER BY created_at DESC
	`
	return r.querySessions(query, resourceID)
}

func (r *sessionRepository) FindActive() ([]*domain.Session, error) {
	query := `
		SELECT id, user_id, resource_id, start_time, end_time, status, 
			client_ip, client_metadata, audit_path, created_at, updated_at
		FROM sessions
		WHERE status = 'active'
		ORDER BY created_at DESC
	`
	return r.querySessions(query)
}

func (r *sessionRepository) Update(session *domain.Session) error {
	query := `
		UPDATE sessions
		SET end_time = ?, status = ?, audit_path = ?, updated_at = ?
		WHERE id = ?
	`
	_, err := r.db.Exec(query,
		session.EndTime,
		session.Status,
		session.AuditPath,
		time.Now(),
		session.ID,
	)
	return err
}

func (r *sessionRepository) Delete(id string) error {
	query := `DELETE FROM sessions WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *sessionRepository) querySessions(query string, args ...interface{}) ([]*domain.Session, error) {
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []*domain.Session
	for rows.Next() {
		session := &domain.Session{}
		var endTime sql.NullTime

		err := rows.Scan(
			&session.ID,
			&session.UserID,
			&session.ResourceID,
			&session.StartTime,
			&endTime,
			&session.Status,
			&session.ClientIP,
			&session.ClientMetadata,
			&session.AuditPath,
			&session.CreatedAt,
			&session.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if endTime.Valid {
			session.EndTime = endTime.Time
		}

		sessions = append(sessions, session)
	}
	return sessions, nil
}
