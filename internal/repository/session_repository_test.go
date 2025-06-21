package repository

import (
	"testing"
	"time"

	"secretary/alpha/internal/domain"
)

func TestSessionRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSessionRepository(db)

	session := &domain.Session{
		UserID:         "user-123",
		Username:       "testuser",
		ResourceID:     "resource-123",
		StartTime:      time.Now(),
		Status:         "active",
		ClientIP:       "192.168.1.1",
		ClientMetadata: "test metadata",
		ExpiresAt:      time.Now().Add(1 * time.Hour),
	}

	err := repo.Create(session)
	if err != nil {
		t.Errorf("Create() error = %v", err)
	}

	if session.ID == "" {
		t.Error("Create() should set session ID")
	}

	if session.CreatedAt.IsZero() {
		t.Error("Create() should set CreatedAt")
	}

	if session.UpdatedAt.IsZero() {
		t.Error("Create() should set UpdatedAt")
	}
}

func TestSessionRepository_FindByID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSessionRepository(db)

	session := &domain.Session{
		UserID:         "user-123",
		Username:       "testuser",
		ResourceID:     "resource-123",
		StartTime:      time.Now(),
		Status:         "active",
		ClientIP:       "192.168.1.1",
		ClientMetadata: "test metadata",
		ExpiresAt:      time.Now().Add(1 * time.Hour),
	}

	err := repo.Create(session)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	found, err := repo.FindByID(session.ID)
	if err != nil {
		t.Errorf("FindByID() error = %v", err)
	}

	if found == nil {
		t.Error("FindByID() should return session")
	}

	if found.UserID != session.UserID {
		t.Errorf("FindByID() UserID = %v, want %v", found.UserID, session.UserID)
	}

	if found.ResourceID != session.ResourceID {
		t.Errorf("FindByID() ResourceID = %v, want %v", found.ResourceID, session.ResourceID)
	}

	if found.Status != session.Status {
		t.Errorf("FindByID() Status = %v, want %v", found.Status, session.Status)
	}
}

func TestSessionRepository_FindByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSessionRepository(db)

	found, err := repo.FindByID("nonexistent-id")
	if err == nil {
		t.Error("FindByID() should return error for non-existent session")
	}

	if found != nil {
		t.Error("FindByID() should return nil for non-existent session")
	}
}

func TestSessionRepository_FindByUserID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSessionRepository(db)

	userID := "user-123"
	sessions := []*domain.Session{
		{
			UserID:     userID,
			Username:   "testuser",
			ResourceID: "resource-1",
			StartTime:  time.Now(),
			Status:     "active",
			ClientIP:   "192.168.1.1",
			ExpiresAt:  time.Now().Add(1 * time.Hour),
		},
		{
			UserID:     userID,
			Username:   "testuser",
			ResourceID: "resource-2",
			StartTime:  time.Now(),
			Status:     "active",
			ClientIP:   "192.168.1.2",
			ExpiresAt:  time.Now().Add(1 * time.Hour),
		},
		{
			UserID:     "other-user",
			Username:   "otheruser",
			ResourceID: "resource-3",
			StartTime:  time.Now(),
			Status:     "active",
			ClientIP:   "192.168.1.3",
			ExpiresAt:  time.Now().Add(1 * time.Hour),
		},
	}

	for _, session := range sessions {
		err := repo.Create(session)
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}
	}

	found, err := repo.FindByUserID(userID)
	if err != nil {
		t.Errorf("FindByUserID() error = %v", err)
	}

	if len(found) != 2 {
		t.Errorf("FindByUserID() returned %d sessions, want 2", len(found))
	}

	for _, session := range found {
		if session.UserID != userID {
			t.Errorf("FindByUserID() returned session with wrong UserID: %v", session.UserID)
		}
	}
}

func TestSessionRepository_FindByResourceID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSessionRepository(db)

	resourceID := "resource-123"
	sessions := []*domain.Session{
		{
			UserID:     "user-1",
			Username:   "user1",
			ResourceID: resourceID,
			StartTime:  time.Now(),
			Status:     "active",
			ClientIP:   "192.168.1.1",
			ExpiresAt:  time.Now().Add(1 * time.Hour),
		},
		{
			UserID:     "user-2",
			Username:   "user2",
			ResourceID: resourceID,
			StartTime:  time.Now(),
			Status:     "active",
			ClientIP:   "192.168.1.2",
			ExpiresAt:  time.Now().Add(1 * time.Hour),
		},
		{
			UserID:     "user-3",
			Username:   "user3",
			ResourceID: "other-resource",
			StartTime:  time.Now(),
			Status:     "active",
			ClientIP:   "192.168.1.3",
			ExpiresAt:  time.Now().Add(1 * time.Hour),
		},
	}

	for _, session := range sessions {
		err := repo.Create(session)
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}
	}

	found, err := repo.FindByResourceID(resourceID)
	if err != nil {
		t.Errorf("FindByResourceID() error = %v", err)
	}

	if len(found) != 2 {
		t.Errorf("FindByResourceID() returned %d sessions, want 2", len(found))
	}

	for _, session := range found {
		if session.ResourceID != resourceID {
			t.Errorf("FindByResourceID() returned session with wrong ResourceID: %v", session.ResourceID)
		}
	}
}

func TestSessionRepository_FindActive(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSessionRepository(db)

	sessions := []*domain.Session{
		{
			UserID:     "user-1",
			Username:   "user1",
			ResourceID: "resource-1",
			StartTime:  time.Now(),
			Status:     "active",
			ClientIP:   "192.168.1.1",
			ExpiresAt:  time.Now().Add(1 * time.Hour),
		},
		{
			UserID:     "user-2",
			Username:   "user2",
			ResourceID: "resource-2",
			StartTime:  time.Now(),
			Status:     "completed",
			ClientIP:   "192.168.1.2",
			EndTime:    time.Now(),
			ExpiresAt:  time.Now().Add(1 * time.Hour),
		},
		{
			UserID:     "user-3",
			Username:   "user3",
			ResourceID: "resource-3",
			StartTime:  time.Now(),
			Status:     "active",
			ClientIP:   "192.168.1.3",
			ExpiresAt:  time.Now().Add(1 * time.Hour),
		},
	}

	for _, session := range sessions {
		err := repo.Create(session)
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}
	}

	found, err := repo.FindActive()
	if err != nil {
		t.Errorf("FindActive() error = %v", err)
	}

	if len(found) != 2 {
		t.Errorf("FindActive() returned %d sessions, want 2", len(found))
	}

	for _, session := range found {
		if session.Status != "active" {
			t.Errorf("FindActive() returned session with wrong status: %v", session.Status)
		}
	}
}

func TestSessionRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSessionRepository(db)

	session := &domain.Session{
		UserID:     "user-123",
		Username:   "testuser",
		ResourceID: "resource-123",
		StartTime:  time.Now(),
		Status:     "active",
		ClientIP:   "192.168.1.1",
		ExpiresAt:  time.Now().Add(1 * time.Hour),
	}

	err := repo.Create(session)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	originalUpdatedAt := session.UpdatedAt
	time.Sleep(1 * time.Millisecond) // Ensure time difference

	session.Status = "completed"
	session.EndTime = time.Now()

	err = repo.Update(session)
	if err != nil {
		t.Errorf("Update() error = %v", err)
	}

	if session.UpdatedAt.Equal(originalUpdatedAt) {
		t.Error("Update() should update UpdatedAt timestamp")
	}

	// Verify the update
	found, err := repo.FindByID(session.ID)
	if err != nil {
		t.Errorf("FindByID() error = %v", err)
	}

	if found.Status != "completed" {
		t.Errorf("Update() status = %v, want %v", found.Status, "completed")
	}

	if found.EndTime.IsZero() {
		t.Error("Update() should set EndTime")
	}
}

func TestSessionRepository_Update_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSessionRepository(db)

	session := &domain.Session{
		ID:         "nonexistent-id",
		UserID:     "user-123",
		Username:   "testuser",
		ResourceID: "resource-123",
		StartTime:  time.Now(),
		Status:     "active",
		ClientIP:   "192.168.1.1",
		ExpiresAt:  time.Now().Add(1 * time.Hour),
	}

	err := repo.Update(session)
	if err == nil {
		t.Error("Update() should return error for non-existent session")
	}
}

func TestSessionRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSessionRepository(db)

	session := &domain.Session{
		UserID:     "user-123",
		Username:   "testuser",
		ResourceID: "resource-123",
		StartTime:  time.Now(),
		Status:     "active",
		ClientIP:   "192.168.1.1",
		ExpiresAt:  time.Now().Add(1 * time.Hour),
	}

	err := repo.Create(session)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	err = repo.Delete(session.ID)
	if err != nil {
		t.Errorf("Delete() error = %v", err)
	}

	// Verify deletion
	found, err := repo.FindByID(session.ID)
	if err == nil {
		t.Error("Delete() should remove session")
	}

	if found != nil {
		t.Error("Delete() should return nil for deleted session")
	}
}

func TestSessionRepository_Delete_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSessionRepository(db)

	err := repo.Delete("nonexistent-id")
	if err == nil {
		t.Error("Delete() should return error for non-existent session")
	}
}
