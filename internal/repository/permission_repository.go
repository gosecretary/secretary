package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"

	"secretary/alpha/internal/domain"
)

type permissionRepository struct {
	db *sql.DB
}

func NewPermissionRepository(db *sql.DB) domain.PermissionRepository {
	return &permissionRepository{db: db}
}

func (r *permissionRepository) Create(permission *domain.Permission) error {
	// Generate UUID if not provided
	if permission.ID == "" {
		permission.ID = uuid.New().String()
	}

	// Set timestamps
	permission.CreatedAt = time.Now()
	permission.UpdatedAt = time.Now()

	query := `
		INSERT INTO permissions (id, user_id, resource_id, role, action, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.Exec(query,
		permission.ID,
		permission.UserID,
		permission.ResourceID,
		permission.Role,
		permission.Action,
		permission.CreatedAt,
		permission.UpdatedAt,
	)
	return err
}

func (r *permissionRepository) FindByID(id string) (*domain.Permission, error) {
	query := `
		SELECT id, user_id, resource_id, role, action, created_at, updated_at
		FROM permissions
		WHERE id = ?
	`
	permission := &domain.Permission{}
	err := r.db.QueryRow(query, id).Scan(
		&permission.ID,
		&permission.UserID,
		&permission.ResourceID,
		&permission.Role,
		&permission.Action,
		&permission.CreatedAt,
		&permission.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New("permission not found")
	}
	return permission, err
}

func (r *permissionRepository) FindByUserID(userID string) ([]*domain.Permission, error) {
	query := `
		SELECT id, user_id, resource_id, role, action, created_at, updated_at
		FROM permissions
		WHERE user_id = ?
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []*domain.Permission
	for rows.Next() {
		permission := &domain.Permission{}
		err := rows.Scan(
			&permission.ID,
			&permission.UserID,
			&permission.ResourceID,
			&permission.Role,
			&permission.Action,
			&permission.CreatedAt,
			&permission.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, permission)
	}
	return permissions, nil
}

func (r *permissionRepository) FindByResourceID(resourceID string) ([]*domain.Permission, error) {
	query := `
		SELECT id, user_id, resource_id, role, action, created_at, updated_at
		FROM permissions
		WHERE resource_id = ?
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query, resourceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []*domain.Permission
	for rows.Next() {
		permission := &domain.Permission{}
		err := rows.Scan(
			&permission.ID,
			&permission.UserID,
			&permission.ResourceID,
			&permission.Role,
			&permission.Action,
			&permission.CreatedAt,
			&permission.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, permission)
	}
	return permissions, nil
}

func (r *permissionRepository) Update(permission *domain.Permission) error {
	permission.UpdatedAt = time.Now()
	query := `
		UPDATE permissions
		SET user_id = ?, resource_id = ?, role = ?, action = ?, updated_at = ?
		WHERE id = ?
	`
	result, err := r.db.Exec(query,
		permission.UserID,
		permission.ResourceID,
		permission.Role,
		permission.Action,
		permission.UpdatedAt,
		permission.ID,
	)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("permission not found")
	}
	return nil
}

func (r *permissionRepository) Delete(id string) error {
	query := `DELETE FROM permissions WHERE id = ?`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("permission not found")
	}
	return nil
}

func (r *permissionRepository) DeleteByUserID(userID string) error {
	query := `DELETE FROM permissions WHERE user_id = ?`
	_, err := r.db.Exec(query, userID)
	return err
}

func (r *permissionRepository) DeleteByResourceID(resourceID string) error {
	query := `DELETE FROM permissions WHERE resource_id = ?`
	_, err := r.db.Exec(query, resourceID)
	return err
}
