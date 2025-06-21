package repository

import (
	"testing"
	"time"

	"secretary/alpha/internal/domain"
)

func TestEphemeralCredentialRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewEphemeralCredentialRepository(db)

	credential := &domain.EphemeralCredential{
		UserID:     "user-123",
		ResourceID: "resource-123",
		Username:   "sec_user_abc123",
		Password:   "securepassword123",
		Token:      "abc123def456",
		ExpiresAt:  time.Now().Add(1 * time.Hour),
		Duration:   "1h",
		Used:       false,
	}

	err := repo.Create(credential)
	if err != nil {
		t.Errorf("Create() error = %v", err)
	}

	if credential.ID == "" {
		t.Error("Create() should set credential ID")
	}

	if credential.CreatedAt.IsZero() {
		t.Error("Create() should set CreatedAt")
	}
}

func TestEphemeralCredentialRepository_FindByID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewEphemeralCredentialRepository(db)

	credential := &domain.EphemeralCredential{
		UserID:     "user-123",
		ResourceID: "resource-123",
		Username:   "sec_user_abc123",
		Password:   "securepassword123",
		Token:      "abc123def456",
		ExpiresAt:  time.Now().Add(1 * time.Hour),
		Duration:   "1h",
		Used:       false,
	}

	err := repo.Create(credential)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	found, err := repo.FindByID(credential.ID)
	if err != nil {
		t.Errorf("FindByID() error = %v", err)
	}

	if found == nil {
		t.Error("FindByID() should return credential")
	}

	if found.UserID != credential.UserID {
		t.Errorf("FindByID() UserID = %v, want %v", found.UserID, credential.UserID)
	}

	if found.ResourceID != credential.ResourceID {
		t.Errorf("FindByID() ResourceID = %v, want %v", found.ResourceID, credential.ResourceID)
	}

	if found.Username != credential.Username {
		t.Errorf("FindByID() Username = %v, want %v", found.Username, credential.Username)
	}

	if found.Token != credential.Token {
		t.Errorf("FindByID() Token = %v, want %v", found.Token, credential.Token)
	}
}

func TestEphemeralCredentialRepository_FindByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewEphemeralCredentialRepository(db)

	found, err := repo.FindByID("nonexistent-id")
	if err == nil {
		t.Error("FindByID() should return error for non-existent credential")
	}

	if found != nil {
		t.Error("FindByID() should return nil for non-existent credential")
	}
}

func TestEphemeralCredentialRepository_FindByToken(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewEphemeralCredentialRepository(db)

	token := "abc123def456"
	credential := &domain.EphemeralCredential{
		UserID:     "user-123",
		ResourceID: "resource-123",
		Username:   "sec_user_abc123",
		Password:   "securepassword123",
		Token:      token,
		ExpiresAt:  time.Now().Add(1 * time.Hour),
		Duration:   "1h",
		Used:       false,
	}

	err := repo.Create(credential)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	found, err := repo.FindByToken(token)
	if err != nil {
		t.Errorf("FindByToken() error = %v", err)
	}

	if found == nil {
		t.Error("FindByToken() should return credential")
	}

	if found.ID != credential.ID {
		t.Errorf("FindByToken() ID = %v, want %v", found.ID, credential.ID)
	}
}

func TestEphemeralCredentialRepository_FindByToken_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewEphemeralCredentialRepository(db)

	found, err := repo.FindByToken("nonexistent-token")
	if err == nil {
		t.Error("FindByToken() should return error for non-existent token")
	}

	if found != nil {
		t.Error("FindByToken() should return nil for non-existent token")
	}
}

func TestEphemeralCredentialRepository_FindByUserID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewEphemeralCredentialRepository(db)

	userID := "user-123"
	credentials := []*domain.EphemeralCredential{
		{
			UserID:     userID,
			ResourceID: "resource-1",
			Username:   "sec_user_abc123",
			Password:   "securepassword123",
			Token:      "token1",
			ExpiresAt:  time.Now().Add(1 * time.Hour),
			Duration:   "1h",
			Used:       false,
		},
		{
			UserID:     userID,
			ResourceID: "resource-2",
			Username:   "sec_user_def456",
			Password:   "securepassword456",
			Token:      "token2",
			ExpiresAt:  time.Now().Add(1 * time.Hour),
			Duration:   "1h",
			Used:       false,
		},
		{
			UserID:     "other-user",
			ResourceID: "resource-3",
			Username:   "sec_user_ghi789",
			Password:   "securepassword789",
			Token:      "token3",
			ExpiresAt:  time.Now().Add(1 * time.Hour),
			Duration:   "1h",
			Used:       false,
		},
	}

	for _, credential := range credentials {
		err := repo.Create(credential)
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}
	}

	found, err := repo.FindByUserID(userID)
	if err != nil {
		t.Errorf("FindByUserID() error = %v", err)
	}

	if len(found) != 2 {
		t.Errorf("FindByUserID() returned %d credentials, want 2", len(found))
	}

	for _, credential := range found {
		if credential.UserID != userID {
			t.Errorf("FindByUserID() returned credential with wrong UserID: %v", credential.UserID)
		}
	}
}

func TestEphemeralCredentialRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewEphemeralCredentialRepository(db)

	credential := &domain.EphemeralCredential{
		UserID:     "user-123",
		ResourceID: "resource-123",
		Username:   "sec_user_abc123",
		Password:   "securepassword123",
		Token:      "abc123def456",
		ExpiresAt:  time.Now().Add(1 * time.Hour),
		Duration:   "1h",
		Used:       false,
	}

	err := repo.Create(credential)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	credential.Used = true
	credential.UsedAt = time.Now()

	err = repo.Update(credential)
	if err != nil {
		t.Errorf("Update() error = %v", err)
	}

	// Verify the update
	found, err := repo.FindByID(credential.ID)
	if err != nil {
		t.Errorf("FindByID() error = %v", err)
	}

	if !found.Used {
		t.Error("Update() should set Used to true")
	}

	if found.UsedAt.IsZero() {
		t.Error("Update() should set UsedAt")
	}
}

func TestEphemeralCredentialRepository_DeleteExpired(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewEphemeralCredentialRepository(db)

	credentials := []*domain.EphemeralCredential{
		{
			UserID:     "user-1",
			ResourceID: "resource-1",
			Username:   "sec_user_abc123",
			Password:   "securepassword123",
			Token:      "token1",
			ExpiresAt:  time.Now().Add(1 * time.Hour), // Not expired
			Duration:   "1h",
			Used:       false,
		},
		{
			UserID:     "user-2",
			ResourceID: "resource-2",
			Username:   "sec_user_def456",
			Password:   "securepassword456",
			Token:      "token2",
			ExpiresAt:  time.Now().Add(-1 * time.Hour), // Expired
			Duration:   "1h",
			Used:       false,
		},
		{
			UserID:     "user-3",
			ResourceID: "resource-3",
			Username:   "sec_user_ghi789",
			Password:   "securepassword789",
			Token:      "token3",
			ExpiresAt:  time.Now().Add(-2 * time.Hour), // Expired
			Duration:   "1h",
			Used:       false,
		},
	}

	for _, credential := range credentials {
		err := repo.Create(credential)
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}
	}

	err := repo.DeleteExpired()
	if err != nil {
		t.Errorf("DeleteExpired() error = %v", err)
	}

	// Verify only non-expired credentials remain
	found, err := repo.FindByUserID("user-1")
	if err != nil {
		t.Errorf("FindByUserID() error = %v", err)
	}

	if len(found) != 1 {
		t.Errorf("DeleteExpired() should keep 1 credential, got %d", len(found))
	}

	// Verify expired credentials are deleted
	found, err = repo.FindByUserID("user-2")
	if err != nil {
		t.Errorf("FindByUserID() error = %v", err)
	}

	if len(found) != 0 {
		t.Errorf("DeleteExpired() should delete expired credentials, got %d", len(found))
	}
}

func TestEphemeralCredentialRepository_DeleteByUserID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewEphemeralCredentialRepository(db)

	userID := "user-123"
	credentials := []*domain.EphemeralCredential{
		{
			UserID:     userID,
			ResourceID: "resource-1",
			Username:   "sec_user_abc123",
			Password:   "securepassword123",
			Token:      "token1",
			ExpiresAt:  time.Now().Add(1 * time.Hour),
			Duration:   "1h",
			Used:       false,
		},
		{
			UserID:     userID,
			ResourceID: "resource-2",
			Username:   "sec_user_def456",
			Password:   "securepassword456",
			Token:      "token2",
			ExpiresAt:  time.Now().Add(1 * time.Hour),
			Duration:   "1h",
			Used:       false,
		},
		{
			UserID:     "other-user",
			ResourceID: "resource-3",
			Username:   "sec_user_ghi789",
			Password:   "securepassword789",
			Token:      "token3",
			ExpiresAt:  time.Now().Add(1 * time.Hour),
			Duration:   "1h",
			Used:       false,
		},
	}

	for _, credential := range credentials {
		err := repo.Create(credential)
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}
	}

	err := repo.DeleteByUserID(userID)
	if err != nil {
		t.Errorf("DeleteByUserID() error = %v", err)
	}

	// Verify user's credentials are deleted
	found, err := repo.FindByUserID(userID)
	if err != nil {
		t.Errorf("FindByUserID() error = %v", err)
	}

	if len(found) != 0 {
		t.Errorf("DeleteByUserID() should delete all user credentials, got %d", len(found))
	}

	// Verify other user's credentials remain
	found, err = repo.FindByUserID("other-user")
	if err != nil {
		t.Errorf("FindByUserID() error = %v", err)
	}

	if len(found) != 1 {
		t.Errorf("DeleteByUserID() should keep other user's credentials, got %d", len(found))
	}
}

func TestEphemeralCredentialRepository_DeleteByResourceID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewEphemeralCredentialRepository(db)

	resourceID := "resource-123"
	credentials := []*domain.EphemeralCredential{
		{
			UserID:     "user-1",
			ResourceID: resourceID,
			Username:   "sec_user_abc123",
			Password:   "securepassword123",
			Token:      "token1",
			ExpiresAt:  time.Now().Add(1 * time.Hour),
			Duration:   "1h",
			Used:       false,
		},
		{
			UserID:     "user-2",
			ResourceID: resourceID,
			Username:   "sec_user_def456",
			Password:   "securepassword456",
			Token:      "token2",
			ExpiresAt:  time.Now().Add(1 * time.Hour),
			Duration:   "1h",
			Used:       false,
		},
		{
			UserID:     "user-3",
			ResourceID: "other-resource",
			Username:   "sec_user_ghi789",
			Password:   "securepassword789",
			Token:      "token3",
			ExpiresAt:  time.Now().Add(1 * time.Hour),
			Duration:   "1h",
			Used:       false,
		},
	}

	for _, credential := range credentials {
		err := repo.Create(credential)
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}
	}

	err := repo.DeleteByResourceID(resourceID)
	if err != nil {
		t.Errorf("DeleteByResourceID() error = %v", err)
	}

	// Verify resource's credentials are deleted
	found, err := repo.FindByUserID("user-1")
	if err != nil {
		t.Errorf("FindByUserID() error = %v", err)
	}

	if len(found) != 0 {
		t.Errorf("DeleteByResourceID() should delete all resource credentials, got %d", len(found))
	}

	// Verify other resource's credentials remain
	found, err = repo.FindByUserID("user-3")
	if err != nil {
		t.Errorf("FindByUserID() error = %v", err)
	}

	if len(found) != 1 {
		t.Errorf("DeleteByResourceID() should keep other resource's credentials, got %d", len(found))
	}
}
