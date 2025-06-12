package repository

import (
	"database/sql"
	"errors"

	"secretary/alpha/internal/domain"
)

type accessRequestRepository struct {
	db *sql.DB
}

func NewAccessRequestRepository(db *sql.DB) domain.AccessRequestRepository {
	return &accessRequestRepository{db: db}
}

func (r *accessRequestRepository) Create(request *domain.AccessRequest) error {
	query := `
		INSERT INTO access_requests (
			id, user_id, resource_id, reason, status, reviewer_id, 
			review_notes, requested_at, reviewed_at, expires_at, created_at, updated_at
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
		SELECT id, user_id, resource_id, reason, status, reviewer_id, 
			review_notes, requested_at, reviewed_at, expires_at, created_at, updated_at
		FROM access_requests
		WHERE id = ?
	`
	request := &domain.AccessRequest{}
	var reviewerID, reviewNotes sql.NullString
	var reviewedAt, expiresAt sql.NullTime

	err := r.db.QueryRow(query, id).Scan(
		&request.ID,
		&request.UserID,
		&request.ResourceID,
		&request.Reason,
		&request.Status,
		&reviewerID,
		&reviewNotes,
		&request.RequestedAt,
		&reviewedAt,
		&expiresAt,
		&request.CreatedAt,
		&request.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("access request not found")
	}

	if reviewerID.Valid {
		request.ReviewerID = reviewerID.String
	}

	if reviewNotes.Valid {
		request.ReviewNotes = reviewNotes.String
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
		SELECT id, user_id, resource_id, reason, status, reviewer_id, 
			review_notes, requested_at, reviewed_at, expires_at, created_at, updated_at
		FROM access_requests
		WHERE user_id = ?
		ORDER BY created_at DESC
	`
	return r.queryAccessRequests(query, userID)
}

func (r *accessRequestRepository) FindByResourceID(resourceID string) ([]*domain.AccessRequest, error) {
	query := `
		SELECT id, user_id, resource_id, reason, status, reviewer_id, 
			review_notes, requested_at, reviewed_at, expires_at, created_at, updated_at
		FROM access_requests
		WHERE resource_id = ?
		ORDER BY created_at DESC
	`
	return r.queryAccessRequests(query, resourceID)
}

func (r *accessRequestRepository) FindByStatus(status string) ([]*domain.AccessRequest, error) {
	query := `
		SELECT id, user_id, resource_id, reason, status, reviewer_id, 
			review_notes, requested_at, reviewed_at, expires_at, created_at, updated_at
		FROM access_requests
		WHERE status = ?
		ORDER BY created_at ASC
	`
	return r.queryAccessRequests(query, status)
}

func (r *accessRequestRepository) Update(request *domain.AccessRequest) error {
	query := `
		UPDATE access_requests
		SET status = ?, reviewer_id = ?, review_notes = ?, 
			reviewed_at = ?, expires_at = ?, updated_at = ?
		WHERE id = ?
	`
	_, err := r.db.Exec(query,
		request.Status,
		request.ReviewerID,
		request.ReviewNotes,
		request.ReviewedAt,
		request.ExpiresAt,
		request.UpdatedAt,
		request.ID,
	)
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
		var reviewerID, reviewNotes sql.NullString
		var reviewedAt, expiresAt sql.NullTime

		err := rows.Scan(
			&request.ID,
			&request.UserID,
			&request.ResourceID,
			&request.Reason,
			&request.Status,
			&reviewerID,
			&reviewNotes,
			&request.RequestedAt,
			&reviewedAt,
			&expiresAt,
			&request.CreatedAt,
			&request.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if reviewerID.Valid {
			request.ReviewerID = reviewerID.String
		}

		if reviewNotes.Valid {
			request.ReviewNotes = reviewNotes.String
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
