package middleware

import (
	"net/http"
	"strings"
	"time"

	"arandu/internal/platform/logger"
)

// TelemetryMiddleware captura métricas HTTP e loga informações de telemetria
type TelemetryMiddleware struct {
	skipPaths []string
}

// NewTelemetryMiddleware cria um novo middleware de telemetria
func NewTelemetryMiddleware(skipPaths ...string) *TelemetryMiddleware {
	return &TelemetryMiddleware{
		skipPaths: skipPaths,
	}
}

// Middleware retorna o handler HTTP com telemetria
func (tm *TelemetryMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verifica se o path deve ser ignorado
		if tm.shouldSkip(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		// Captura tempo de início
		start := time.Now()

		// Wrapper para capturar status code
		wrappedWriter := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Executa o handler
		next.ServeHTTP(wrappedWriter, r)

		// Calcula duração
		duration := time.Since(start)

		// Extrai path seguro (sem query params sensíveis)
		safePath := tm.sanitizePath(r.URL.Path, r.URL.RawQuery)

		// Extrai IP do cliente
		clientIP := tm.getClientIP(r)

		// Loga a requisição usando o logger estruturado
		logger.InfoContext(r.Context(), "HTTP Request",
			logger.String("method", r.Method),
			logger.String("path", safePath),
			logger.Int("status", wrappedWriter.statusCode),
			logger.Int64("duration_ms", duration.Milliseconds()),
			logger.String("ip", clientIP),
		)
	})
}

// shouldSkip verifica se o path deve ser ignorado
func (tm *TelemetryMiddleware) shouldSkip(path string) bool {
	for _, skipPath := range tm.skipPaths {
		if strings.HasPrefix(path, skipPath) {
			return true
		}
	}
	return false
}

// sanitizePath remove parâmetros sensíveis da URL
func (tm *TelemetryMiddleware) sanitizePath(path, query string) string {
	// Se não houver query, retorna o path
	if query == "" {
		return path
	}

	// Lista de parâmetros sensíveis que devem ser removidos
	sensitiveParams := []string{"code", "token", "password", "secret", "api_key", "access_token", "refresh_token"}

	// Parse da query string
	params := make(map[string]string)
	for _, part := range strings.Split(query, "&") {
		if idx := strings.Index(part, "="); idx > 0 {
			key := part[:idx]
			value := part[idx+1:]
			params[key] = value
		}
	}

	// Constrói query string filtrada
	var filteredParams []string
	for key, value := range params {
		isSensitive := false
		for _, sensitive := range sensitiveParams {
			if strings.EqualFold(key, sensitive) {
				isSensitive = true
				break
			}
		}

		if isSensitive {
			filteredParams = append(filteredParams, key+"=[REDACTED]")
		} else {
			filteredParams = append(filteredParams, key+"="+value)
		}
	}

	// Reconstrói a URL
	if len(filteredParams) > 0 {
		return path + "?" + strings.Join(filteredParams, "&")
	}
	return path
}

// getClientIP extrai o IP do cliente da requisição
func (tm *TelemetryMiddleware) getClientIP(r *http.Request) string {
	// Verifica headers de proxy
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	if xForwardedFor != "" {
		// Pega o primeiro IP da lista
		ips := strings.Split(xForwardedFor, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	xRealIP := r.Header.Get("X-Real-Ip")
	if xRealIP != "" {
		return xRealIP
	}

	// Fallback para RemoteAddr
	ip := r.RemoteAddr
	if idx := strings.LastIndex(ip, ":"); idx != -1 {
		return ip[:idx]
	}
	return ip
}

// responseWriter wrapper para capturar status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
	written    bool
}

func (rw *responseWriter) WriteHeader(code int) {
	if !rw.written {
		rw.statusCode = code
		rw.written = true
		rw.ResponseWriter.WriteHeader(code)
	}
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.written {
		rw.WriteHeader(http.StatusOK)
	}
	return rw.ResponseWriter.Write(b)
}
