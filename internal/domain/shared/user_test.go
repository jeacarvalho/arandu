package shared

import (
	"testing"
)

func TestNewUser(t *testing.T) {
	user := NewUser("user-uuid", "test@example.com", "hash123", "tenant-uuid")

	if user.ID != "user-uuid" {
		t.Errorf("expected ID 'user-uuid', got '%s'", user.ID)
	}
	if user.Email != "test@example.com" {
		t.Errorf("expected Email 'test@example.com', got '%s'", user.Email)
	}
	if user.PasswordHash != "hash123" {
		t.Errorf("expected PasswordHash 'hash123', got '%s'", user.PasswordHash)
	}
	if user.TenantID != "tenant-uuid" {
		t.Errorf("expected TenantID 'tenant-uuid', got '%s'", user.TenantID)
	}
	if user.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
}

func TestUserValidatePassword(t *testing.T) {
	hash, err := HashPassword("correct-password")
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	user := NewUser("id", "test@example.com", hash, "tenant")

	if !user.ValidatePassword("correct-password") {
		t.Error("expected password validation to succeed")
	}

	if user.ValidatePassword("wrong-password") {
		t.Error("expected password validation to fail with wrong password")
	}
}

func TestUserValidatePasswordEmptyHash(t *testing.T) {
	user := NewUser("id", "test@example.com", "", "tenant")

	if user.ValidatePassword("any-password") {
		t.Error("expected validation to fail with empty hash")
	}
}

func TestHashPassword(t *testing.T) {
	hash1, err := HashPassword("password123")
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	if hash1 == "" {
		t.Error("expected non-empty hash")
	}

	hash2, err := HashPassword("password123")
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	if hash1 == hash2 {
		t.Error("expected different hashes for same password (salt)")
	}
}

func TestUserUpdateLastLogin(t *testing.T) {
	user := NewUser("id", "test@example.com", "hash", "tenant")

	if user.LastLogin != nil {
		t.Error("expected LastLogin to be nil initially")
	}

	user.UpdateLastLogin()

	if user.LastLogin == nil {
		t.Error("expected LastLogin to be set after UpdateLastLogin")
	}
}

func TestUserHasPassword(t *testing.T) {
	userWithPassword := NewUser("id", "test@example.com", "hash", "tenant")
	if !userWithPassword.HasPassword() {
		t.Error("expected user with password to return true")
	}

	userWithoutPassword := NewUser("id", "test@example.com", "", "tenant")
	if userWithoutPassword.HasPassword() {
		t.Error("expected user without password to return false")
	}
}

func TestUser_ValidatePassword_InvalidHash(t *testing.T) {
	user := NewUser("id", "test@example.com", "not-a-valid-bcrypt-hash", "tenant")
	if user.ValidatePassword("anypassword") {
		t.Error("expected validation to fail with invalid hash")
	}
}

func TestUser_UpdateLastLogin_Twice(t *testing.T) {
	user := NewUser("id", "test@example.com", "hash", "tenant")
	user.UpdateLastLogin()
	first := user.LastLogin

	user.UpdateLastLogin()
	if user.LastLogin.Before(*first) {
		t.Error("LastLogin should be updated to a later time")
	}
}
