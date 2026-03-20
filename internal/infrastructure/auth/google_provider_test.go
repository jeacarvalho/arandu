package auth

import (
	"testing"
	"time"
)

func TestGenerateState(t *testing.T) {
	state1, err := GenerateState()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if state1 == "" {
		t.Error("expected non-empty state")
	}

	state2, err := GenerateState()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if state1 == state2 {
		t.Error("expected different states (random)")
	}

	if len(state1) < 20 {
		t.Errorf("expected state to be longer, got %d", len(state1))
	}
}

func TestOAuthSession_IsValid(t *testing.T) {
	session := &OAuthSession{
		State:     "test-state",
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}

	if !session.IsValid() {
		t.Error("expected session to be valid")
	}
}

func TestOAuthSession_Expired(t *testing.T) {
	session := &OAuthSession{
		State:     "test-state",
		ExpiresAt: time.Now().Add(-10 * time.Minute),
	}

	if session.IsValid() {
		t.Error("expected session to be expired")
	}
}

func TestGoogleProvider_ValidateState(t *testing.T) {
	gp, err := NewGoogleProvider()
	if err != nil {
		t.Skip("Google provider not configured")
	}

	validState := "test-state"

	if !gp.ValidateState(validState, validState) {
		t.Error("expected state validation to pass")
	}

	if gp.ValidateState("state1", "state2") {
		t.Error("expected state validation to fail for different states")
	}

	if gp.ValidateState("", "state") {
		t.Error("expected state validation to fail for empty state")
	}

	if gp.ValidateState("state", "") {
		t.Error("expected state validation to fail for empty expected")
	}

	if gp.ValidateState("", "") {
		t.Error("expected state validation to fail for both empty")
	}
}
