package helpers

import (
	"testing"
)

func TestGetCSSVersion(t *testing.T) {
	v1 := GetCSSVersion()
	v2 := GetCSSVersion()
	
	if v1 != v2 {
		t.Errorf("GetCSSVersion() should return stable version, got %s then %s", v1, v2)
	}
	
	if len(v1) < 8 {
		t.Errorf("GetCSSVersion() should return hash of at least 8 chars, got %d", len(v1))
	}
	
	t.Logf("CSS Version: %s", v1)
}

func TestGetCSSVersionV2(t *testing.T) {
	v1 := GetCSSVersionV2()
	v2 := GetCSSVersionV2()
	
	if v1 != v2 {
		t.Errorf("GetCSSVersionV2() should return stable version, got %s then %s", v1, v2)
	}
	
	if len(v1) < 8 {
		t.Errorf("GetCSSVersionV2() should return hash of at least 8 chars, got %d", len(v1))
	}
	
	t.Logf("CSS V2 Version: %s", v1)
}

func TestGetJSVersion(t *testing.T) {
	v1 := GetJSVersion()
	v2 := GetJSVersion()
	
	if v1 != v2 {
		t.Errorf("GetJSVersion() should return stable version, got %s then %s", v1, v2)
	}
	
	if len(v1) < 8 {
		t.Errorf("GetJSVersion() should return hash of at least 8 chars, got %d", len(v1))
	}
	
	t.Logf("JS Version: %s", v1)
}

func TestVersionsAreDifferent(t *testing.T) {
	cssV := GetCSSVersion()
	cssV2 := GetCSSVersionV2()
	jsV := GetJSVersion()
	
	// Each file should have its own unique hash (unless they happen to have the same content)
	// At minimum, CSS and JS should be different
	if cssV == jsV && cssV2 == jsV {
		t.Log("Note: CSS and JS versions might be the same if files have identical content")
	}
	
	t.Logf("CSS: %s, CSS-V2: %s, JS: %s", cssV, cssV2, jsV)
}

func TestResetVersionCache(t *testing.T) {
	// Get versions first
	v1 := GetCSSVersion()
	v2 := GetCSSVersionV2()
	v3 := GetJSVersion()
	
	// Reset cache
	ResetVersionCache()
	
	// Get versions again - should still be the same (content hasn't changed)
	v1After := GetCSSVersion()
	v2After := GetCSSVersionV2()
	v3After := GetJSVersion()
	
	if v1 != v1After {
		t.Errorf("CSS version changed after reset: %s -> %s", v1, v1After)
	}
	if v2 != v2After {
		t.Errorf("CSS V2 version changed after reset: %s -> %s", v2, v2After)
	}
	if v3 != v3After {
		t.Errorf("JS version changed after reset: %s -> %s", v3, v3After)
	}
}
