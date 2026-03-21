package ai

import (
	"testing"
	"time"
)

func TestCache_Get_Set(t *testing.T) {
	cache := NewCache(1 * time.Hour)

	entry, found := cache.Get("patient-1", "30d")
	if found {
		t.Error("expected entry not to be found before Set")
	}

	cache.Set("patient-1", "30d", "Synthesis for patient 1")

	entry, found = cache.Get("patient-1", "30d")
	if !found {
		t.Fatal("expected entry to be found after Set")
	}
	if entry.Synthesis != "Synthesis for patient 1" {
		t.Errorf("expected synthesis text, got %s", entry.Synthesis)
	}
}

func TestCache_Get_Expired(t *testing.T) {
	cache := NewCache(1 * time.Millisecond)

	cache.Set("patient-1", "30d", "Synthesis")

	time.Sleep(5 * time.Millisecond)

	_, found := cache.Get("patient-1", "30d")
	if found {
		t.Error("expected entry to be expired and not found")
	}
}

func TestCache_Get_DifferentKeys(t *testing.T) {
	cache := NewCache(1 * time.Hour)

	cache.Set("patient-1", "30d", "Synthesis 1")
	cache.Set("patient-2", "30d", "Synthesis 2")
	cache.Set("patient-1", "7d", "Synthesis 3")

	entry, found := cache.Get("patient-1", "30d")
	if !found || entry.Synthesis != "Synthesis 1" {
		t.Errorf("patient-1/30d: expected 'Synthesis 1', got %v", entry)
	}

	entry, found = cache.Get("patient-2", "30d")
	if !found || entry.Synthesis != "Synthesis 2" {
		t.Errorf("patient-2/30d: expected 'Synthesis 2', got %v", entry)
	}

	entry, found = cache.Get("patient-1", "7d")
	if !found || entry.Synthesis != "Synthesis 3" {
		t.Errorf("patient-1/7d: expected 'Synthesis 3', got %v", entry)
	}
}

func TestCache_Clear(t *testing.T) {
	cache := NewCache(1 * time.Hour)

	cache.Set("patient-1", "30d", "Synthesis 1")
	cache.Set("patient-2", "30d", "Synthesis 2")

	if cache.Size() != 2 {
		t.Errorf("expected size 2, got %d", cache.Size())
	}

	cache.Clear()

	if cache.Size() != 0 {
		t.Errorf("expected size 0 after Clear, got %d", cache.Size())
	}

	_, found := cache.Get("patient-1", "30d")
	if found {
		t.Error("expected entry not to be found after Clear")
	}
}

func TestCache_Size(t *testing.T) {
	cache := NewCache(1 * time.Hour)

	if cache.Size() != 0 {
		t.Errorf("expected size 0 initially, got %d", cache.Size())
	}

	cache.Set("p1", "30d", "s1")
	if cache.Size() != 1 {
		t.Errorf("expected size 1, got %d", cache.Size())
	}

	cache.Set("p2", "30d", "s2")
	if cache.Size() != 2 {
		t.Errorf("expected size 2, got %d", cache.Size())
	}

	cache.Set("p1", "30d", "s1-updated")
	if cache.Size() != 2 {
		t.Errorf("expected size 2 (same key), got %d", cache.Size())
	}
}

func TestCache_Cleanup(t *testing.T) {
	cache := NewCache(1 * time.Millisecond)

	cache.Set("patient-1", "30d", "Synthesis 1")
	cache.Set("patient-2", "30d", "Synthesis 2")

	time.Sleep(5 * time.Millisecond)

	cache.Set("patient-3", "30d", "Synthesis 3")

	if cache.Size() != 3 {
		t.Errorf("expected size 3 before cleanup, got %d", cache.Size())
	}

	cache.Cleanup()

	if cache.Size() != 1 {
		t.Errorf("expected size 1 after cleanup, got %d", cache.Size())
	}

	_, found := cache.Get("patient-3", "30d")
	if !found {
		t.Error("expected patient-3 to survive cleanup")
	}

	_, found = cache.Get("patient-1", "30d")
	if found {
		t.Error("expected patient-1 to be removed by cleanup")
	}
}

func TestCache_Cleanup_NoExpired(t *testing.T) {
	cache := NewCache(1 * time.Hour)

	cache.Set("patient-1", "30d", "Synthesis 1")

	before := cache.Size()
	cache.Cleanup()

	if cache.Size() != before {
		t.Errorf("expected size unchanged after cleanup with no expired items")
	}
}

func TestCache_GenerateKey(t *testing.T) {
	cache := NewCache(1 * time.Hour)

	key1 := cache.generateKey("patient-1", "30d")
	key2 := cache.generateKey("patient-1", "30d")
	key3 := cache.generateKey("patient-2", "30d")

	if key1 != key2 {
		t.Error("expected same key for same params")
	}
	if key1 == key3 {
		t.Error("expected different key for different patient ID")
	}

	if len(key1) != 64 {
		t.Errorf("expected SHA256 hex length 64, got %d", len(key1))
	}
}

func TestNewCache(t *testing.T) {
	ttl := 2 * time.Hour
	cache := NewCache(ttl)

	if cache.Size() != 0 {
		t.Errorf("expected empty cache, got size %d", cache.Size())
	}
}

func TestCacheEntry_Times(t *testing.T) {
	cache := NewCache(1 * time.Hour)
	before := time.Now()

	cache.Set("patient-1", "30d", "Synthesis")

	entry, found := cache.Get("patient-1", "30d")
	if !found {
		t.Fatal("entry not found")
	}

	if entry.GeneratedAt.Before(before) {
		t.Error("GeneratedAt should be >= before")
	}
	if entry.ExpiresAt.Before(entry.GeneratedAt) {
		t.Error("ExpiresAt should be after GeneratedAt")
	}
	if entry.ExpiresAt.Before(time.Now()) {
		t.Error("ExpiresAt should be in the future")
	}
}
