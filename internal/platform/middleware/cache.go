package middleware

import "net/http"

// HTMXCacheMiddleware sets appropriate cache headers for HTMX requests
// Fragments should not be cached by the browser
func HTMXCacheMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if this is an HTMX request
		isHTMX := r.Header.Get("HX-Request") == "true"

		if isHTMX {
			// HTMX fragments should not be cached
			w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
			w.Header().Set("Pragma", "no-cache")
			w.Header().Set("Expires", "0")
			w.Header().Set("Vary", "HX-Request")
		} else {
			// For full page requests, allow caching but with validation
			w.Header().Set("Cache-Control", "private, no-cache")
			w.Header().Set("Vary", "HX-Request")
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
