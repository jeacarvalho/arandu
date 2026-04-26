package env_test

import (
	"testing"

	"arandu/internal/platform/env"
)

func TestIsDev_ReturnsFalseByDefault(t *testing.T) {
	t.Setenv("APP_ENV", "")
	if env.IsDev() {
		t.Error("IsDev() should return false when APP_ENV is empty")
	}
}

func TestIsDev_TrueWhenDevelopment(t *testing.T) {
	t.Setenv("APP_ENV", "development")
	if !env.IsDev() {
		t.Error("IsDev() should return true when APP_ENV=development")
	}
}

func TestIsDev_TrueWhenDev(t *testing.T) {
	t.Setenv("APP_ENV", "dev")
	if !env.IsDev() {
		t.Error("IsDev() should return true when APP_ENV=dev")
	}
}

func TestIsDev_CaseInsensitive(t *testing.T) {
	cases := []string{"DEV", "Dev", "DEVELOPMENT", "Development"}
	for _, v := range cases {
		t.Run(v, func(t *testing.T) {
			t.Setenv("APP_ENV", v)
			if !env.IsDev() {
				t.Errorf("IsDev() should return true when APP_ENV=%q", v)
			}
		})
	}
}

func TestIsDev_FalseForProductionValues(t *testing.T) {
	cases := []string{"production", "prod", "staging", "test", "release"}
	for _, v := range cases {
		t.Run(v, func(t *testing.T) {
			t.Setenv("APP_ENV", v)
			if env.IsDev() {
				t.Errorf("IsDev() should return false when APP_ENV=%q", v)
			}
		})
	}
}

func TestIsDev_TrimsWhitespace(t *testing.T) {
	t.Setenv("APP_ENV", "  dev  ")
	if !env.IsDev() {
		t.Error("IsDev() should trim whitespace before comparing")
	}
}
