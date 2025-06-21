package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"

	"secretary/alpha/internal/domain"
)

type resourceRepository struct {
	db *sql.DB
}

func NewResourceRepository(db *sql.DB) domain.ResourceRepository {
	return &resourceRepository{db: db}
}

func (r *resourceRepository) Create(resource *domain.Resource) error {
	if resource.ID == "" {
		resource.ID = uuid.New().String()
	}
	if resource.CreatedAt.IsZero() {
		resource.CreatedAt = time.Now()
	}
	if resource.UpdatedAt.IsZero() {
		resource.UpdatedAt = time.Now()
	}
	query := `
		INSERT INTO resources (id, name, description, type, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.Exec(query, resource.ID, resource.Name, resource.Description, resource.Type, resource.CreatedAt, resource.UpdatedAt)
	return err
}

func (r *resourceRepository) FindByID(id string) (*domain.Resource, error) {
	query := `
		SELECT id, name, description, type, created_at, updated_at
		FROM resources
		WHERE id = ?
	`
	resource := &domain.Resource{}
	err := r.db.QueryRow(query, id).Scan(
		&resource.ID,
		&resource.Name,
		&resource.Description,
		&resource.Type,
		&resource.CreatedAt,
		&resource.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New("resource not found")
	}
	return resource, err
}

func (r *resourceRepository) FindAll() ([]*domain.Resource, error) {
	query := `
		SELECT id, name, description, type, created_at, updated_at
		FROM resources
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var resources []*domain.Resource
	for rows.Next() {
		resource := &domain.Resource{}
		err := rows.Scan(
			&resource.ID,
			&resource.Name,
			&resource.Description,
			&resource.Type,
			&resource.CreatedAt,
			&resource.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		resources = append(resources, resource)
	}
	return resources, nil
}

func (r *resourceRepository) Update(resource *domain.Resource) error {
	resource.UpdatedAt = time.Now()
	query := `
		UPDATE resources
		SET name = ?, description = ?, type = ?, updated_at = ?
		WHERE id = ?
	`
	result, err := r.db.Exec(query,
		resource.Name,
		resource.Description,
		resource.Type,
		resource.UpdatedAt,
		resource.ID,
	)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("resource not found")
	}
	return nil
}

func (r *resourceRepository) Delete(id string) error {
	query := `DELETE FROM resources WHERE id = ?`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("resource not found")
	}
	return nil
}
