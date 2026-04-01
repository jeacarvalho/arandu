package middleware

import "net/http"

// IsHTMXRequest checks if the request is from HTMX
func IsHTMXRequest(r *http.Request) bool {
	return r.Header.Get("HX-Request") == "true"
}

// HTMXCacheMiddleware sets appropriate cache headers for HTMX requests
// Fragments should not be cached by the browser
// All HTML responses use strict no-cache to prevent CSS desync issues
func HTMXCacheMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if this is an HTMX request
		isHTMX := IsHTMXRequest(r)

		// Always prevent caching of HTML responses to avoid CSS desync
		// This is the safest approach for HTMX applications with dynamic CSS
		if isHTMX {
			// HTMX fragments should never be cached
			w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
			w.Header().Set("Pragma", "no-cache")
			w.Header().Set("Expires", "0")
			w.Header().Set("Vary", "HX-Request")
		} else {
			// For full page requests, use strict no-cache to ensure CSS is always fresh
			// This prevents the "CSS works on first load but breaks on navigation" issue
			w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate, max-age=0")
			w.Header().Set("Pragma", "no-cache")
			w.Header().Set("Expires", "0")
			w.Header().Set("Vary", "HX-Request, Accept-Encoding")
		}

		next.ServeHTTP(w, r)
	})
}

// StaticAssetsCacheMiddleware sets aggressive caching for static assets with version query strings
func StaticAssetsCacheMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if request has version query parameter
		if r.URL.Query().Has("v") {
			// Assets with version can be cached aggressively (1 year)
			w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		} else {
			// Assets without version - cache but revalidate
			w.Header().Set("Cache-Control", "public, max-age=3600, must-revalidate")
		}

		next.ServeHTTP(w, r)
	})
}
