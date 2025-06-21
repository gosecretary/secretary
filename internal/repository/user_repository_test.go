package repository

import (
	"database/sql"
	"testing"
	"time"

	"secretary/alpha/internal/domain"
)

func setupTestDB(t *testing.T) *sql.DB {
	// Use a temporary database for testing
	dbPath := ":memory:"
	db, err := InitDB("sqlite3", dbPath)
	if err != nil {
		t.Fatalf("Failed to initialize test database: %v", err)
	}
	return db
}

func TestUserRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)

	user := &domain.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "hashedpassword",
		Name:     "Test User",
		Role:     "user",
	}

	err := repo.Create(user)
	if err != nil {
		t.Errorf("Create() error = %v", err)
	}

	if user.ID == "" {
		t.Error("Create() should set user ID")
	}

	if user.CreatedAt.IsZero() {
		t.Error("Create() should set CreatedAt")
	}

	if user.UpdatedAt.IsZero() {
		t.Error("Create() should set UpdatedAt")
	}
}

func TestUserRepository_Create_DuplicateUsername(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)

	user1 := &domain.User{
		Username: "testuser",
		Email:    "test1@example.com",
		Password: "hashedpassword",
		Name:     "Test User 1",
		Role:     "user",
	}

	user2 := &domain.User{
		Username: "testuser", // Same username
		Email:    "test2@example.com",
		Password: "hashedpassword",
		Name:     "Test User 2",
		Role:     "user",
	}

	err := repo.Create(user1)
	if err != nil {
		t.Errorf("First Create() error = %v", err)
	}

	err = repo.Create(user2)
	if err == nil {
		t.Error("Create() should fail with duplicate username")
	}
}

func TestUserRepository_Create_DuplicateEmail(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)

	user1 := &domain.User{
		Username: "testuser1",
		Email:    "test@example.com",
		Password: "hashedpassword",
		Name:     "Test User 1",
		Role:     "user",
	}

	user2 := &domain.User{
		Username: "testuser2",
		Email:    "test@example.com", // Same email
		Password: "hashedpassword",
		Name:     "Test User 2",
		Role:     "user",
	}

	err := repo.Create(user1)
	if err != nil {
		t.Errorf("First Create() error = %v", err)
	}

	err = repo.Create(user2)
	if err == nil {
		t.Error("Create() should fail with duplicate email")
	}
}

func TestUserRepository_FindByID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)

	user := &domain.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "hashedpassword",
		Name:     "Test User",
		Role:     "user",
	}

	err := repo.Create(user)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	found, err := repo.FindByID(user.ID)
	if err != nil {
		t.Errorf("FindByID() error = %v", err)
	}

	if found == nil {
		t.Error("FindByID() should return user")
	}

	if found.Username != user.Username {
		t.Errorf("FindByID() username = %v, want %v", found.Username, user.Username)
	}

	if found.Email != user.Email {
		t.Errorf("FindByID() email = %v, want %v", found.Email, user.Email)
	}
}

func TestUserRepository_FindByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)

	found, err := repo.FindByID("nonexistent-id")
	if err == nil {
		t.Error("FindByID() should return error for non-existent user")
	}

	if found != nil {
		t.Error("FindByID() should return nil for non-existent user")
	}
}

func TestUserRepository_FindByUsername(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)

	user := &domain.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "hashedpassword",
		Name:     "Test User",
		Role:     "user",
	}

	err := repo.Create(user)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	found, err := repo.FindByUsername("testuser")
	if err != nil {
		t.Errorf("FindByUsername() error = %v", err)
	}

	if found == nil {
		t.Error("FindByUsername() should return user")
	}

	if found.ID != user.ID {
		t.Errorf("FindByUsername() ID = %v, want %v", found.ID, user.ID)
	}
}

func TestUserRepository_FindByUsername_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)

	found, err := repo.FindByUsername("nonexistent")
	if err == nil {
		t.Error("FindByUsername() should return error for non-existent username")
	}

	if found != nil {
		t.Error("FindByUsername() should return nil for non-existent username")
	}
}

func TestUserRepository_FindByEmail(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)

	user := &domain.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "hashedpassword",
		Name:     "Test User",
		Role:     "user",
	}

	err := repo.Create(user)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	found, err := repo.FindByEmail("test@example.com")
	if err != nil {
		t.Errorf("FindByEmail() error = %v", err)
	}

	if found == nil {
		t.Error("FindByEmail() should return user")
	}

	if found.ID != user.ID {
		t.Errorf("FindByEmail() ID = %v, want %v", found.ID, user.ID)
	}
}

func TestUserRepository_FindByEmail_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)

	found, err := repo.FindByEmail("nonexistent@example.com")
	if err == nil {
		t.Error("FindByEmail() should return error for non-existent email")
	}

	if found != nil {
		t.Error("FindByEmail() should return nil for non-existent email")
	}
}

func TestUserRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)

	user := &domain.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "hashedpassword",
		Name:     "Test User",
		Role:     "user",
	}

	err := repo.Create(user)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	originalUpdatedAt := user.UpdatedAt
	time.Sleep(1 * time.Millisecond) // Ensure time difference

	user.Name = "Updated Test User"
	user.Email = "updated@example.com"

	err = repo.Update(user)
	if err != nil {
		t.Errorf("Update() error = %v", err)
	}

	if user.UpdatedAt.Equal(originalUpdatedAt) {
		t.Error("Update() should update UpdatedAt timestamp")
	}

	// Verify the update
	found, err := repo.FindByID(user.ID)
	if err != nil {
		t.Errorf("FindByID() error = %v", err)
	}

	if found.Name != "Updated Test User" {
		t.Errorf("Update() name = %v, want %v", found.Name, "Updated Test User")
	}

	if found.Email != "updated@example.com" {
		t.Errorf("Update() email = %v, want %v", found.Email, "updated@example.com")
	}
}

func TestUserRepository_Update_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)

	user := &domain.User{
		ID:       "nonexistent-id",
		Username: "testuser",
		Email:    "test@example.com",
		Password: "hashedpassword",
		Name:     "Test User",
		Role:     "user",
	}

	err := repo.Update(user)
	if err == nil {
		t.Error("Update() should return error for non-existent user")
	}
}

func TestUserRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)

	user := &domain.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "hashedpassword",
		Name:     "Test User",
		Role:     "user",
	}

	err := repo.Create(user)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	err = repo.Delete(user.ID)
	if err != nil {
		t.Errorf("Delete() error = %v", err)
	}

	// Verify deletion
	found, err := repo.FindByID(user.ID)
	if err == nil {
		t.Error("Delete() should remove user")
	}

	if found != nil {
		t.Error("Delete() should return nil for deleted user")
	}
}

func TestUserRepository_Delete_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)

	err := repo.Delete("nonexistent-id")
	if err == nil {
		t.Error("Delete() should return error for non-existent user")
	}
}
