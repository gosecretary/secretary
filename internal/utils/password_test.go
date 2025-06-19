package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "valid password",
			password: "testpass123",
			wantErr:  false,
		},
		{
			name:     "empty password",
			password: "",
			wantErr:  false, // bcrypt can hash empty strings
		},
		{
			name:     "long password",
			password: "this_is_a_very_long_password_that_should_work_fine_with_bcrypt",
			wantErr:  false,
		},
		{
			name:     "password with special characters",
			password: "p@ssw0rd!@#$%^&*()",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := HashPassword(tt.password)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, hash)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, hash)
				// Hash should be different from original password
				assert.NotEqual(t, tt.password, hash)
				// Hash should start with bcrypt prefix
				assert.Contains(t, hash, "$2a$")
			}
		})
	}
}

func TestCheckPasswordHash(t *testing.T) {
	// Pre-generate some hashes for testing
	validHash, _ := HashPassword("testpass123")

	tests := []struct {
		name     string
		password string
		hash     string
		want     bool
	}{
		{
			name:     "correct password",
			password: "testpass123",
			hash:     validHash,
			want:     true,
		},
		{
			name:     "incorrect password",
			password: "wrongpassword",
			hash:     validHash,
			want:     false,
		},
		{
			name:     "empty password with valid hash",
			password: "",
			hash:     validHash,
			want:     false,
		},
		{
			name:     "valid password with empty hash",
			password: "testpass123",
			hash:     "",
			want:     false,
		},
		{
			name:     "invalid hash format",
			password: "testpass123",
			hash:     "invalid_hash",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CheckPasswordHash(tt.password, tt.hash)
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestHashPassword_ConsistentLength(t *testing.T) {
	// Test that hash length is consistent
	password := "testpass123"

	hash1, err1 := HashPassword(password)
	hash2, err2 := HashPassword(password)

	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.NotEmpty(t, hash1)
	assert.NotEmpty(t, hash2)

	// Hashes should be different (due to salt) but same length
	assert.NotEqual(t, hash1, hash2)
	assert.Equal(t, len(hash1), len(hash2))
}

func TestPasswordRoundTrip(t *testing.T) {
	// Test that we can hash a password and then verify it
	passwords := []string{
		"simple",
		"Complex@Pass123!",
		"with spaces and symbols !@#$%^&*()",
		"unicode_测试_password",
		"",
	}

	for _, password := range passwords {
		t.Run("password_"+password, func(t *testing.T) {
			// Hash the password
			hash, err := HashPassword(password)
			assert.NoError(t, err)
			assert.NotEmpty(t, hash)

			// Verify the password
			isValid := CheckPasswordHash(password, hash)
			assert.True(t, isValid)

			// Verify wrong password fails
			isInvalid := CheckPasswordHash(password+"wrong", hash)
			assert.False(t, isInvalid)
		})
	}
}
