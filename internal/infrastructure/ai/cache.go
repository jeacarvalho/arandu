package ai

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync"
	"time"
)

// CacheEntry represents a cached AI response
type CacheEntry struct {
	Synthesis   string
	GeneratedAt time.Time
	ExpiresAt   time.Time
}

// Cache implements a simple in-memory cache for AI responses
type Cache struct {
	mu    sync.RWMutex
	store map[string]CacheEntry
	ttl   time.Duration
}

// NewCache creates a new cache with the specified TTL
func NewCache(ttl time.Duration) *Cache {
	return &Cache{
		store: make(map[string]CacheEntry),
		ttl:   ttl,
	}
}

// generateKey creates a cache key from patient ID and timeframe
func (c *Cache) generateKey(patientID, timeframe string) string {
	hash := sha256.Sum256([]byte(fmt.Sprintf("%s:%s", patientID, timeframe)))
	return hex.EncodeToString(hash[:])
}

// Get retrieves a cached response if it exists and is not expired
func (c *Cache) Get(patientID, timeframe string) (*CacheEntry, bool) {
	key := c.generateKey(patientID, timeframe)

	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, exists := c.store[key]
	if !exists {
		return nil, false
	}

	// Check if entry has expired
	if time.Now().After(entry.ExpiresAt) {
		return nil, false
	}

	return &entry, true
}

// Set stores a response in the cache
func (c *Cache) Set(patientID, timeframe, synthesis string) {
	key := c.generateKey(patientID, timeframe)
	now := time.Now()

	entry := CacheEntry{
		Synthesis:   synthesis,
		GeneratedAt: now,
		ExpiresAt:   now.Add(c.ttl),
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	c.store[key] = entry
}

// Clear removes all entries from the cache
func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store = make(map[string]CacheEntry)
}

// Size returns the number of entries in the cache
func (c *Cache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.store)
}

// Cleanup removes expired entries from the cache
func (c *Cache) Cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for key, entry := range c.store {
		if now.After(entry.ExpiresAt) {
			delete(c.store, key)
		}
	}
}

// CachedGeminiClient wraps a GeminiClient with caching
type CachedGeminiClient struct {
	client GeminiClient
	cache  *Cache
}

// NewCachedGeminiClient creates a new cached Gemini client
func NewCachedGeminiClient(client GeminiClient, cache *Cache) *CachedGeminiClient {
	return &CachedGeminiClient{
		client: client,
		cache:  cache,
	}
}

// GenerateWithRetry generates content with caching
func (c *CachedGeminiClient) GenerateWithRetry(ctx context.Context, prompt string, maxRetries int) (string, error) {
	// For caching, we need to extract patient ID and timeframe from the prompt
	// This is a simplified approach - in production, we'd need a better way
	// to extract these parameters from the prompt
	patientID, timeframe := extractCacheParams(prompt)

	if patientID != "" && timeframe != "" {
		// Try to get from cache first
		if entry, found := c.cache.Get(patientID, timeframe); found {
			return entry.Synthesis, nil
		}
	}

	// Not in cache or can't extract params, call the actual API
	result, err := c.client.GenerateWithRetry(ctx, prompt, maxRetries)
	if err != nil {
		return "", err
	}

	// Store in cache if we have the params
	if patientID != "" && timeframe != "" {
		c.cache.Set(patientID, timeframe, result)
	}

	return result, nil
}

// Close closes the underlying client
func (c *CachedGeminiClient) Close() error {
	return c.client.Close()
}

// extractCacheParams attempts to extract patient ID and timeframe from prompt
// This is a simplified implementation - in production, we'd pass these as separate parameters
func extractCacheParams(prompt string) (patientID, timeframe string) {
	// Look for patterns in the prompt
	// This assumes the prompt contains patient data in a specific format
	// For now, return empty strings - we'll need to modify the AIService to pass these explicitly
	return "", ""
}
