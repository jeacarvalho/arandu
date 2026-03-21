package shared

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           string     `json:"id"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"-"`
	TenantID     string     `json:"tenant_id"`
	CreatedAt    time.Time  `json:"created_at"`
	LastLogin    *time.Time `json:"last_login,omitempty"`
}

func NewUser(id, email, passwordHash, tenantID string) *User {
	return &User{
		ID:           id,
		Email:        email,
		PasswordHash: passwordHash,
		TenantID:     tenantID,
		CreatedAt:    time.Now(),
	}
}

func (u *User) ValidatePassword(password string) bool {
	if u.PasswordHash == "" {
		return false
	}
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.New("failed to hash password")
	}
	return string(bytes), nil
}

func (u *User) UpdateLastLogin() {
	now := time.Now()
	u.LastLogin = &now
}

func (u *User) HasPassword() bool {
	return u.PasswordHash != ""
}
