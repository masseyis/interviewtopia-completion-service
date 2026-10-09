// Package appconfig contains the small amount of runtime configuration supplied
// by the starter. Keeping it in a value makes tests independent of process-wide
// environment variables.
package appconfig

import "os"

const (
	defaultHTTPAddr       = "127.0.0.1:8080"
	defaultTrustStorePath = "config/trusted-issuers.json"
	defaultDatabasePath   = "completion.db"
)

// Config contains runtime values that the completion service may need.
// RegistryBaseURL is intentionally empty until the interviewer supplies it.
type Config struct {
	HTTPAddr        string
	RegistryBaseURL string
	TrustStorePath  string
	DatabasePath    string
}

// FromEnv constructs the executable's configuration. Tests can construct a
// Config directly instead of mutating the environment.
func FromEnv() Config {
	return Config{
		HTTPAddr:        envOrDefault("HTTP_ADDR", defaultHTTPAddr),
		RegistryBaseURL: os.Getenv("REGISTRY_BASE_URL"),
		TrustStorePath:  envOrDefault("TRUST_STORE_PATH", defaultTrustStorePath),
		DatabasePath:    envOrDefault("COMPLETION_DB_PATH", defaultDatabasePath),
	}
}

func envOrDefault(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}
