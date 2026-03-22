package logger

import (
	"context"
	"log/slog"
	"os"
	"time"

	appcontext "arandu/internal/platform/context"
	"arandu/internal/platform/version"
)

var defaultLogger *slog.Logger

func init() {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	defaultLogger = slog.New(handler).With(
		slog.String("app", version.AppName),
		slog.String("version", version.Version),
		slog.String("commit", version.Commit),
	)
}

// FromContext retorna um logger com contexto enriquecido com tenant_id e request_id
// extraídos do context da requisição
func FromContext(ctx context.Context) *slog.Logger {
	logger := defaultLogger

	// Tenta extrair tenant_id do contexto
	if tenantID, err := appcontext.GetTenantID(ctx); err == nil {
		logger = logger.With(slog.String("tenant_id", tenantID))
	}

	// Tenta extrair request_id do contexto
	if requestID, err := appcontext.GetRequestID(ctx); err == nil {
		logger = logger.With(slog.String("request_id", requestID))
	}

	// Tenta extrair user_id do contexto
	if userID, err := appcontext.GetUserID(ctx); err == nil {
		logger = logger.With(slog.String("user_id", userID))
	}

	return logger
}

// Info logs uma mensagem de nível INFO
func Info(msg string, attrs ...slog.Attr) {
	defaultLogger.Info(msg, attrsToAny(attrs)...)
}

// Error logs uma mensagem de nível ERROR
func Error(msg string, attrs ...slog.Attr) {
	defaultLogger.Error(msg, attrsToAny(attrs)...)
}

// Warn logs uma mensagem de nível WARN
func Warn(msg string, attrs ...slog.Attr) {
	defaultLogger.Warn(msg, attrsToAny(attrs)...)
}

// Debug logs uma mensagem de nível DEBUG
func Debug(msg string, attrs ...slog.Attr) {
	defaultLogger.Debug(msg, attrsToAny(attrs)...)
}

// InfoContext logs uma mensagem de nível INFO com contexto
func InfoContext(ctx context.Context, msg string, attrs ...slog.Attr) {
	FromContext(ctx).Info(msg, attrsToAny(attrs)...)
}

// ErrorContext logs uma mensagem de nível ERROR com contexto
func ErrorContext(ctx context.Context, msg string, attrs ...slog.Attr) {
	FromContext(ctx).Error(msg, attrsToAny(attrs)...)
}

// WarnContext logs uma mensagem de nível WARN com contexto
func WarnContext(ctx context.Context, msg string, attrs ...slog.Attr) {
	FromContext(ctx).Warn(msg, attrsToAny(attrs)...)
}

// DebugContext logs uma mensagem de nível DEBUG com contexto
func DebugContext(ctx context.Context, msg string, attrs ...slog.Attr) {
	FromContext(ctx).Debug(msg, attrsToAny(attrs)...)
}

func attrsToAny(attrs []slog.Attr) []any {
	result := make([]any, len(attrs))
	for i, attr := range attrs {
		result[i] = attr
	}
	return result
}

// String cria um atributo de string para uso nos logs
func String(key, value string) slog.Attr {
	return slog.String(key, value)
}

// Int cria um atributo de int para uso nos logs
func Int(key string, value int) slog.Attr {
	return slog.Int(key, value)
}

// Int64 cria um atributo de int64 para uso nos logs
func Int64(key string, value int64) slog.Attr {
	return slog.Int64(key, value)
}

// Float64 cria um atributo de float64 para uso nos logs
func Float64(key string, value float64) slog.Attr {
	return slog.Float64(key, value)
}

// Bool cria um atributo de bool para uso nos logs
func Bool(key string, value bool) slog.Attr {
	return slog.Bool(key, value)
}

// Time cria um atributo de time.Time para uso nos logs
func Time(key string, value time.Time) slog.Attr {
	return slog.Time(key, value)
}

// Duration cria um atributo de time.Duration para uso nos logs
func Duration(key string, value time.Duration) slog.Attr {
	return slog.Duration(key, value)
}

// Any cria um atributo genérico para uso nos logs
func Any(key string, value any) slog.Attr {
	return slog.Any(key, value)
}
