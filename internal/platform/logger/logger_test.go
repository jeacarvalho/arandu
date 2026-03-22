package logger

import (
	"context"
	"log/slog"
	"testing"

	appcontext "arandu/internal/platform/context"
	"arandu/internal/platform/version"
)

func TestString(t *testing.T) {
	attr := String("key", "value")
	if attr.Key != "key" {
		t.Errorf("Expected key 'key', got %s", attr.Key)
	}
	if attr.Value.String() != "value" {
		t.Errorf("Expected value 'value', got %v", attr.Value)
	}
}

func TestInt(t *testing.T) {
	attr := Int("count", 42)
	if attr.Key != "count" {
		t.Errorf("Expected key 'count', got %s", attr.Key)
	}
	if attr.Value.Int64() != 42 {
		t.Errorf("Expected value 42, got %v", attr.Value)
	}
}

func TestBool(t *testing.T) {
	attr := Bool("enabled", true)
	if attr.Key != "enabled" {
		t.Errorf("Expected key 'enabled', got %s", attr.Key)
	}
	if !attr.Value.Bool() {
		t.Errorf("Expected value true, got %v", attr.Value)
	}
}

func TestFromContextWithTenantID(t *testing.T) {
	ctx := appcontext.WithTenantID(context.Background(), "tenant-123")
	logger := FromContext(ctx)

	if logger == nil {
		t.Error("Expected non-nil logger")
	}

	// Verifica que o logger tem os atributos globais
	if version.AppName != "arandu" {
		t.Errorf("Expected app name 'arandu', got %s", version.AppName)
	}
}

func TestFromContextWithRequestID(t *testing.T) {
	ctx := appcontext.WithRequestID(context.Background(), "req-456")
	logger := FromContext(ctx)

	if logger == nil {
		t.Error("Expected non-nil logger")
	}
}

func TestFromContextWithUserID(t *testing.T) {
	ctx := appcontext.WithUserID(context.Background(), "user-789")
	logger := FromContext(ctx)

	if logger == nil {
		t.Error("Expected non-nil logger")
	}
}

func TestFromContextWithAllContext(t *testing.T) {
	ctx := context.Background()
	ctx = appcontext.WithTenantID(ctx, "tenant-abc")
	ctx = appcontext.WithRequestID(ctx, "req-xyz")
	ctx = appcontext.WithUserID(ctx, "user-123")

	logger := FromContext(ctx)

	if logger == nil {
		t.Error("Expected non-nil logger")
	}
}

func TestAttrsToAny(t *testing.T) {
	attrs := []slog.Attr{
		String("key1", "value1"),
		Int("key2", 42),
		Bool("key3", true),
	}

	result := attrsToAny(attrs)

	if len(result) != 3 {
		t.Errorf("Expected 3 attributes, got %d", len(result))
	}
}
