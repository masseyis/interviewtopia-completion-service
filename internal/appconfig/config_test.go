package appconfig

import "testing"

func TestFromEnvDefaults(t *testing.T) {
	t.Setenv("HTTP_ADDR", "")
	t.Setenv("REGISTRY_BASE_URL", "")
	t.Setenv("TRUST_STORE_PATH", "")
	t.Setenv("COMPLETION_DB_PATH", "")

	got := FromEnv()
	if got.HTTPAddr != "127.0.0.1:8080" {
		t.Fatalf("HTTPAddr = %q", got.HTTPAddr)
	}
	if got.RegistryBaseURL != "" {
		t.Fatalf("RegistryBaseURL = %q, want empty", got.RegistryBaseURL)
	}
	if got.TrustStorePath != "config/trusted-issuers.json" {
		t.Fatalf("TrustStorePath = %q", got.TrustStorePath)
	}
	if got.DatabasePath != "completion.db" {
		t.Fatalf("DatabasePath = %q", got.DatabasePath)
	}
}

func TestFromEnvOverridesDefaults(t *testing.T) {
	t.Setenv("HTTP_ADDR", "127.0.0.1:9090")
	t.Setenv("REGISTRY_BASE_URL", "http://registry.example.test")
	t.Setenv("TRUST_STORE_PATH", "/tmp/test-trust-store.json")
	t.Setenv("COMPLETION_DB_PATH", "/tmp/test-completion.db")

	got := FromEnv()
	if got.HTTPAddr != "127.0.0.1:9090" {
		t.Fatalf("HTTPAddr = %q", got.HTTPAddr)
	}
	if got.RegistryBaseURL != "http://registry.example.test" {
		t.Fatalf("RegistryBaseURL = %q", got.RegistryBaseURL)
	}
	if got.TrustStorePath != "/tmp/test-trust-store.json" {
		t.Fatalf("TrustStorePath = %q", got.TrustStorePath)
	}
	if got.DatabasePath != "/tmp/test-completion.db" {
		t.Fatalf("DatabasePath = %q", got.DatabasePath)
	}
}
