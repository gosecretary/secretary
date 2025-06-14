package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"

	"secretary/alpha/internal/domain"
)

type accessRequestRepository struct {
	db *sql.DB
}

func NewAccessRequestRepository(db *sql.DB) domain.AccessRequestRepository {
	return &accessRequestRepository{db: db}
}

func (r *accessRequestRepository) Create(request *domain.AccessRequest) error {
	// Generate UUID if not provided
	if request.ID == "" {
		request.ID = uuid.New().String()
	}

	// Set timestamps
	request.CreatedAt = time.Now()
	request.UpdatedAt = time.Now()

	query := `
		INSERT INTO access_requests (
			id, user_id, resource_id, reason, status, reviewer_id, review_notes,
			requested_at, reviewed_at, expires_at, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.Exec(query,
		request.ID,
		request.UserID,
		request.ResourceID,
		request.Reason,
		request.Status,
		request.ReviewerID,
		request.ReviewNotes,
		request.RequestedAt,
		request.ReviewedAt,
		request.ExpiresAt,
		request.CreatedAt,
		request.UpdatedAt,
	)
	return err
}

func (r *accessRequestRepository) FindByID(id string) (*domain.AccessRequest, error) {
	query := `
		SELECT id, user_id, resource_id, reason, status, reviewer_id, review_notes,
			requested_at, reviewed_at, expires_at, created_at, updated_at
		FROM access_requests
		WHERE id = ?
	`
	request := &domain.AccessRequest{}
	var reviewedAt, expiresAt sql.NullTime

	err := r.db.QueryRow(query, id).Scan(
		&request.ID,
		&request.UserID,
		&request.ResourceID,
		&request.Reason,
		&request.Status,
		&request.ReviewerID,
		&request.ReviewNotes,
		&request.RequestedAt,
		&reviewedAt,
		&expiresAt,
		&request.CreatedAt,
		&request.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("access request not found")
	}

	if reviewedAt.Valid {
		request.ReviewedAt = reviewedAt.Time
	}
	if expiresAt.Valid {
		request.ExpiresAt = expiresAt.Time
	}

	return request, err
}

func (r *accessRequestRepository) FindByUserID(userID string) ([]*domain.AccessRequest, error) {
	query := `
		SELECT id, user_id, resource_id, reason, status, reviewer_id, review_notes,
			requested_at, reviewed_at, expires_at, created_at, updated_at
		FROM access_requests
		WHERE user_id = ?
		ORDER BY created_at DESC
	`
	return r.queryAccessRequests(query, userID)
}

func (r *accessRequestRepository) FindByResourceID(resourceID string) ([]*domain.AccessRequest, error) {
	query := `
		SELECT id, user_id, resource_id, reason, status, reviewer_id, review_notes,
			requested_at, reviewed_at, expires_at, created_at, updated_at
		FROM access_requests
		WHERE resource_id = ?
		ORDER BY created_at DESC
	`
	return r.queryAccessRequests(query, resourceID)
}

func (r *accessRequestRepository) FindByStatus(status string) ([]*domain.AccessRequest, error) {
	query := `
		SELECT id, user_id, resource_id, reason, status, reviewer_id, review_notes,
			requested_at, reviewed_at, expires_at, created_at, updated_at
		FROM access_requests
		WHERE status = ?
		ORDER BY created_at DESC
	`
	return r.queryAccessRequests(query, status)
}

func (r *accessRequestRepository) Update(request *domain.AccessRequest) error {
	request.UpdatedAt = time.Now()
	query := `
		UPDATE access_requests
		SET user_id = ?, resource_id = ?, reason = ?, status = ?, reviewer_id = ?,
			review_notes = ?, requested_at = ?, reviewed_at = ?, expires_at = ?, updated_at = ?
		WHERE id = ?
	`
	_, err := r.db.Exec(query,
		request.UserID,
		request.ResourceID,
		request.Reason,
		request.Status,
		request.ReviewerID,
		request.ReviewNotes,
		request.RequestedAt,
		request.ReviewedAt,
		request.ExpiresAt,
		request.UpdatedAt,
		request.ID,
	)
	return err
}

func (r *accessRequestRepository) Delete(id string) error {
	query := `DELETE FROM access_requests WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *accessRequestRepository) queryAccessRequests(query string, args ...interface{}) ([]*domain.AccessRequest, error) {
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []*domain.AccessRequest
	for rows.Next() {
		request := &domain.AccessRequest{}
		var reviewedAt, expiresAt sql.NullTime

		err := rows.Scan(
			&request.ID,
			&request.UserID,
			&request.ResourceID,
			&request.Reason,
			&request.Status,
			&request.ReviewerID,
			&request.ReviewNotes,
			&request.RequestedAt,
			&reviewedAt,
			&expiresAt,
			&request.CreatedAt,
			&request.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if reviewedAt.Valid {
			request.ReviewedAt = reviewedAt.Time
		}
		if expiresAt.Valid {
			request.ExpiresAt = expiresAt.Time
		}

		requests = append(requests, request)
	}
	return requests, nil
}
