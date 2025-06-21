package repository

import (
	"testing"
	"time"

	"secretary/alpha/internal/domain"
)

func TestResourceRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewResourceRepository(db)

	resource := &domain.Resource{
		Name:        "test-db-create",
		Description: "Test database",
		Type:        "postgresql",
	}

	err := repo.Create(resource)
	if err != nil {
		t.Errorf("Create() error = %v", err)
	}

	if resource.ID == "" {
		t.Error("Create() should set resource ID")
	}

	if resource.CreatedAt.IsZero() {
		t.Error("Create() should set CreatedAt")
	}

	if resource.UpdatedAt.IsZero() {
		t.Error("Create() should set UpdatedAt")
	}
}

func TestResourceRepository_Create_DuplicateName(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewResourceRepository(db)

	resource1 := &domain.Resource{
		Name:        "test-db-duplicate",
		Description: "Test database 1",
		Type:        "postgresql",
	}

	resource2 := &domain.Resource{
		Name:        "test-db-duplicate", // Same name
		Description: "Test database 2",
		Type:        "mysql",
	}

	err := repo.Create(resource1)
	if err != nil {
		t.Errorf("First Create() error = %v", err)
	}

	err = repo.Create(resource2)
	if err == nil {
		t.Error("Create() should fail with duplicate name")
	}
}

func TestResourceRepository_FindByID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewResourceRepository(db)

	resource := &domain.Resource{
		Name:        "test-db-findbyid",
		Description: "Test database",
		Type:        "postgresql",
	}

	err := repo.Create(resource)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	found, err := repo.FindByID(resource.ID)
	if err != nil {
		t.Errorf("FindByID() error = %v", err)
	}

	if found == nil {
		t.Error("FindByID() should return resource")
	}

	if found.Name != resource.Name {
		t.Errorf("FindByID() name = %v, want %v", found.Name, resource.Name)
	}

	if found.Description != resource.Description {
		t.Errorf("FindByID() description = %v, want %v", found.Description, resource.Description)
	}

	if found.Type != resource.Type {
		t.Errorf("FindByID() type = %v, want %v", found.Type, resource.Type)
	}
}

func TestResourceRepository_FindByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewResourceRepository(db)

	found, err := repo.FindByID("nonexistent-id")
	if err == nil {
		t.Error("FindByID() should return error for non-existent resource")
	}

	if found != nil {
		t.Error("FindByID() should return nil for non-existent resource")
	}
}

func TestResourceRepository_FindAll(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewResourceRepository(db)

	// Create multiple resources
	resources := []*domain.Resource{
		{
			Name:        "db-1",
			Description: "Database 1",
			Type:        "postgresql",
		},
		{
			Name:        "db-2",
			Description: "Database 2",
			Type:        "mysql",
		},
		{
			Name:        "ssh-server",
			Description: "SSH Server",
			Type:        "ssh",
		},
	}

	for _, resource := range resources {
		err := repo.Create(resource)
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}
	}

	found, err := repo.FindAll()
	if err != nil {
		t.Errorf("FindAll() error = %v", err)
	}

	if len(found) != len(resources) {
		t.Errorf("FindAll() returned %d resources, want %d", len(found), len(resources))
	}

	// Verify all resources are found
	foundNames := make(map[string]bool)
	for _, r := range found {
		foundNames[r.Name] = true
	}

	for _, resource := range resources {
		if !foundNames[resource.Name] {
			t.Errorf("FindAll() missing resource: %s", resource.Name)
		}
	}
}

func TestResourceRepository_FindAll_Empty(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewResourceRepository(db)

	found, err := repo.FindAll()
	if err != nil {
		t.Errorf("FindAll() error = %v", err)
	}

	if len(found) != 0 {
		t.Errorf("FindAll() returned %d resources, want 0", len(found))
	}
}

func TestResourceRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewResourceRepository(db)

	resource := &domain.Resource{
		Name:        "test-db-update",
		Description: "Test database",
		Type:        "postgresql",
	}

	err := repo.Create(resource)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	originalUpdatedAt := resource.UpdatedAt
	time.Sleep(1 * time.Millisecond) // Ensure time difference

	resource.Description = "Updated test database"
	resource.Type = "mysql"

	err = repo.Update(resource)
	if err != nil {
		t.Errorf("Update() error = %v", err)
	}

	if resource.UpdatedAt.Equal(originalUpdatedAt) {
		t.Error("Update() should update UpdatedAt timestamp")
	}

	// Verify the update
	found, err := repo.FindByID(resource.ID)
	if err != nil {
		t.Errorf("FindByID() error = %v", err)
	}

	if found.Description != "Updated test database" {
		t.Errorf("Update() description = %v, want %v", found.Description, "Updated test database")
	}

	if found.Type != "mysql" {
		t.Errorf("Update() type = %v, want %v", found.Type, "mysql")
	}
}

func TestResourceRepository_Update_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewResourceRepository(db)

	resource := &domain.Resource{
		ID:          "nonexistent-id",
		Name:        "test-db",
		Description: "Test database",
		Type:        "postgresql",
	}

	err := repo.Update(resource)
	if err == nil {
		t.Error("Update() should return error for non-existent resource")
	}
}

func TestResourceRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewResourceRepository(db)

	resource := &domain.Resource{
		Name:        "test-db",
		Description: "Test database",
		Type:        "postgresql",
	}

	err := repo.Create(resource)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	err = repo.Delete(resource.ID)
	if err != nil {
		t.Errorf("Delete() error = %v", err)
	}

	// Verify deletion
	found, err := repo.FindByID(resource.ID)
	if err == nil {
		t.Error("Delete() should remove resource")
	}

	if found != nil {
		t.Error("Delete() should return nil for deleted resource")
	}
}

func TestResourceRepository_Delete_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewResourceRepository(db)

	err := repo.Delete("nonexistent-id")
	if err == nil {
		t.Error("Delete() should return error for non-existent resource")
	}
}
