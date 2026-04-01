package helpers

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"sync"
)

var (
	cssVersion       string
	cssVersionV2     string
	jsVersion        string
	cssVersionOnce   sync.Once
	cssVersionV2Once sync.Once
	jsVersionOnce    sync.Once
)

// GetCSSVersion returns a stable version hash for style.css (legacy naming)
// The hash is calculated from the file content and cached for the lifetime of the process
// This prevents unnecessary cache invalidation when the file hasn't actually changed
func GetCSSVersion() string {
	cssVersionOnce.Do(func() {
		cssVersion = calculateFileHash("web/static/css/style.css")
	})
	return cssVersion
}

// GetCSSVersionV2 returns a stable version hash for tailwind-v2.css
// The hash is calculated from the file content and cached for the lifetime of the process
func GetCSSVersionV2() string {
	cssVersionV2Once.Do(func() {
		cssVersionV2 = calculateFileHash("web/static/css/tailwind-v2.css")
	})
	return cssVersionV2
}

// GetJSVersion returns a stable version hash for htmx-handlers.js
// The hash is calculated from the file content and cached for the lifetime of the process
func GetJSVersion() string {
	jsVersionOnce.Do(func() {
		jsVersion = calculateFileHash("web/static/js/htmx-handlers.js")
	})
	return jsVersion
}

// calculateFileHash reads the file and returns the first 8 characters of its SHA256 hash
// Falls back to a dev version if the file cannot be read
func calculateFileHash(filepath string) string {
	content, err := os.ReadFile(filepath)
	if err != nil {
		// File not found or error - use a consistent dev version
		// This won't change between requests, preventing CSS thrashing
		return "dev_notfound"
	}

	hash := sha256.Sum256(content)
	// Use first 8 characters of hash for a short but unique version string
	return hex.EncodeToString(hash[:4])
}

// ResetVersionCache resets all cached version strings
// This is useful for testing or when files are regenerated during runtime
func ResetVersionCache() {
	cssVersionOnce = sync.Once{}
	cssVersionV2Once = sync.Once{}
	jsVersionOnce = sync.Once{}
	cssVersion = ""
	cssVersionV2 = ""
	jsVersion = ""
}

// ResetCSSVersionCache is an alias for ResetVersionCache (backward compatibility)
// Deprecated: Use ResetVersionCache instead
func ResetCSSVersionCache() {
	ResetVersionCache()
}
