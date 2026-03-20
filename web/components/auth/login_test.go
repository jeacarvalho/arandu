package auth

import (
	"testing"
)

func TestLoginData(t *testing.T) {
	data := LoginData{
		Error: "Test error",
	}

	if data.Error != "Test error" {
		t.Errorf("expected 'Test error', got '%s'", data.Error)
	}
}

func TestLoginData_EmptyError(t *testing.T) {
	data := LoginData{
		Error: "",
	}

	if data.Error != "" {
		t.Errorf("expected empty string, got '%s'", data.Error)
	}
}
