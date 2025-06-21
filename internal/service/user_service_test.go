package service

import (
	"context"
	"testing"
	"time"

	"secretary/alpha/internal/domain"
	"secretary/alpha/internal/repository"
)

func setupTestUserService(t *testing.T) (domain.UserService, domain.UserRepository) {
	db, err := repository.InitDB("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to initialize test database: %v", err)
	}

	userRepo := repository.NewUserRepository(db)
	userService := NewUserService(userRepo)

	return userService, userRepo
}

func TestUserService_CreateUser(t *testing.T) {
	userService, _ := setupTestUserService(t)

	user := &domain.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "testpassword123",
		Name:     "Test User",
		Role:     "user",
	}

	err := userService.CreateUser(context.Background(), user)
	if err != nil {
		t.Errorf("CreateUser() error = %v", err)
	}

	if user.ID == "" {
		t.Error("CreateUser() should set user ID")
	}

	if user.CreatedAt.IsZero() {
		t.Error("CreateUser() should set CreatedAt")
	}

	if user.UpdatedAt.IsZero() {
		t.Error("CreateUser() should set UpdatedAt")
	}

	// Verify password was hashed
	if user.Password == "testpassword123" {
		t.Error("CreateUser() should hash the password")
	}
}

func TestUserService_CreateUser_DuplicateUsername(t *testing.T) {
	userService, _ := setupTestUserService(t)

	user1 := &domain.User{
		Username: "testuser",
		Email:    "test1@example.com",
		Password: "testpassword123",
		Name:     "Test User 1",
		Role:     "user",
	}

	user2 := &domain.User{
		Username: "testuser", // Same username
		Email:    "test2@example.com",
		Password: "testpassword123",
		Name:     "Test User 2",
		Role:     "user",
	}

	err := userService.CreateUser(context.Background(), user1)
	if err != nil {
		t.Errorf("First CreateUser() error = %v", err)
	}

	err = userService.CreateUser(context.Background(), user2)
	if err == nil {
		t.Error("CreateUser() should fail with duplicate username")
	}
}

func TestUserService_CreateUser_DuplicateEmail(t *testing.T) {
	userService, _ := setupTestUserService(t)

	user1 := &domain.User{
		Username: "testuser1",
		Email:    "test@example.com",
		Password: "testpassword123",
		Name:     "Test User 1",
		Role:     "user",
	}

	user2 := &domain.User{
		Username: "testuser2",
		Email:    "test@example.com", // Same email
		Password: "testpassword123",
		Name:     "Test User 2",
		Role:     "user",
	}

	err := userService.CreateUser(context.Background(), user1)
	if err != nil {
		t.Errorf("First CreateUser() error = %v", err)
	}

	err = userService.CreateUser(context.Background(), user2)
	if err == nil {
		t.Error("CreateUser() should fail with duplicate email")
	}
}

func TestUserService_GetByID(t *testing.T) {
	userService, _ := setupTestUserService(t)

	user := &domain.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "testpassword123",
		Name:     "Test User",
		Role:     "user",
	}

	err := userService.CreateUser(context.Background(), user)
	if err != nil {
		t.Fatalf("CreateUser() error = %v", err)
	}

	found, err := userService.GetByID(context.Background(), user.ID)
	if err != nil {
		t.Errorf("GetByID() error = %v", err)
	}

	if found == nil {
		t.Error("GetByID() should return user")
	}

	if found.Username != user.Username {
		t.Errorf("GetByID() username = %v, want %v", found.Username, user.Username)
	}

	if found.Email != user.Email {
		t.Errorf("GetByID() email = %v, want %v", found.Email, user.Email)
	}
}

func TestUserService_GetByID_NotFound(t *testing.T) {
	userService, _ := setupTestUserService(t)

	found, err := userService.GetByID(context.Background(), "nonexistent-id")
	if err == nil {
		t.Error("GetByID() should return error for non-existent user")
	}

	if found != nil {
		t.Error("GetByID() should return nil for non-existent user")
	}
}

func TestUserService_Update(t *testing.T) {
	userService, _ := setupTestUserService(t)

	user := &domain.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "testpassword123",
		Name:     "Test User",
		Role:     "user",
	}

	err := userService.CreateUser(context.Background(), user)
	if err != nil {
		t.Fatalf("CreateUser() error = %v", err)
	}

	originalUpdatedAt := user.UpdatedAt
	time.Sleep(1 * time.Millisecond) // Ensure time difference

	user.Name = "Updated Test User"
	user.Email = "updated@example.com"

	err = userService.Update(context.Background(), user)
	if err != nil {
		t.Errorf("Update() error = %v", err)
	}

	if user.UpdatedAt.Equal(originalUpdatedAt) {
		t.Error("Update() should update UpdatedAt timestamp")
	}

	// Verify the update
	found, err := userService.GetByID(context.Background(), user.ID)
	if err != nil {
		t.Errorf("GetByID() error = %v", err)
	}

	if found.Name != "Updated Test User" {
		t.Errorf("Update() name = %v, want %v", found.Name, "Updated Test User")
	}

	if found.Email != "updated@example.com" {
		t.Errorf("Update() email = %v, want %v", found.Email, "updated@example.com")
	}
}

func TestUserService_Update_NotFound(t *testing.T) {
	userService, _ := setupTestUserService(t)

	user := &domain.User{
		ID:       "nonexistent-id",
		Username: "testuser",
		Email:    "test@example.com",
		Password: "testpassword123",
		Name:     "Test User",
		Role:     "user",
	}

	err := userService.Update(context.Background(), user)
	if err == nil {
		t.Error("Update() should return error for non-existent user")
	}
}

func TestUserService_Delete(t *testing.T) {
	userService, _ := setupTestUserService(t)

	user := &domain.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "testpassword123",
		Name:     "Test User",
		Role:     "user",
	}

	err := userService.CreateUser(context.Background(), user)
	if err != nil {
		t.Fatalf("CreateUser() error = %v", err)
	}

	err = userService.Delete(context.Background(), user.ID)
	if err != nil {
		t.Errorf("Delete() error = %v", err)
	}

	// Verify deletion
	found, err := userService.GetByID(context.Background(), user.ID)
	if err == nil {
		t.Error("Delete() should remove user")
	}

	if found != nil {
		t.Error("Delete() should return nil for deleted user")
	}
}

func TestUserService_Delete_NotFound(t *testing.T) {
	userService, _ := setupTestUserService(t)

	err := userService.Delete(context.Background(), "nonexistent-id")
	if err == nil {
		t.Error("Delete() should return error for non-existent user")
	}
}

func TestUserService_Authenticate(t *testing.T) {
	userService, _ := setupTestUserService(t)

	user := &domain.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "testpassword123",
		Name:     "Test User",
		Role:     "user",
	}

	err := userService.CreateUser(context.Background(), user)
	if err != nil {
		t.Fatalf("CreateUser() error = %v", err)
	}

	// Test valid credentials
	authenticatedUser, err := userService.Authenticate(context.Background(), "testuser", "testpassword123")
	if err != nil {
		t.Errorf("Authenticate() error = %v", err)
	}

	if authenticatedUser == nil {
		t.Error("Authenticate() should return user for valid credentials")
	}

	if authenticatedUser.ID != user.ID {
		t.Errorf("Authenticate() returned wrong user ID: %v, want %v", authenticatedUser.ID, user.ID)
	}

	// Test invalid password
	authenticatedUser, err = userService.Authenticate(context.Background(), "testuser", "wrongpassword")
	if err == nil {
		t.Error("Authenticate() should return error for invalid password")
	}

	if authenticatedUser != nil {
		t.Error("Authenticate() should return nil for invalid password")
	}

	// Test non-existent user
	authenticatedUser, err = userService.Authenticate(context.Background(), "nonexistent", "testpassword123")
	if err == nil {
		t.Error("Authenticate() should return error for non-existent user")
	}

	if authenticatedUser != nil {
		t.Error("Authenticate() should return nil for non-existent user")
	}
}
