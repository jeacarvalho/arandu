package env

import (
	"os"
	"strings"
)

// IsDev reports whether the application is running in development mode.
// Set APP_ENV=development or APP_ENV=dev to enable dev mode.
func IsDev() bool {
	v := strings.ToLower(strings.TrimSpace(os.Getenv("APP_ENV")))
	return v == "development" || v == "dev"
}
