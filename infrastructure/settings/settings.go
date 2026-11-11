package settings

import "github.com/vitor-franz/book-log/entities"

// Provider supplies configuration for the application and its components.
type Provider interface {
	// GetDatabaseSettings returns the cached connection and migration configuration.
	GetDatabaseSettings() entities.DatabaseSettings

	// GetJWTSettings returns the cached token signing and cookie configuration.
	GetJWTSettings() entities.JWTSettings

	// GetLoggerSettings returns the cached logging output configuration.
	GetLoggerSettings() entities.LoggerSettings

	// GetRuntimeSettings returns the cached application execution environment.
	GetRuntimeSettings() entities.RuntimeSettings
}
