package appconfig

import "testing"

func TestFromEnvDefaults(t *testing.T) {
	t.Setenv("HTTP_ADDR", "")
	t.Setenv("REGISTRY_URL", "")
	t.Setenv("TRUST_STORE_PATH", "")

	got := FromEnv()
	if got.HTTPAddr != "127.0.0.1:8080" {
		t.Fatalf("HTTPAddr = %q", got.HTTPAddr)
	}
	if got.RegistryURL != "" {
		t.Fatalf("RegistryURL = %q, want empty", got.RegistryURL)
	}
	if got.TrustStorePath != "config/trusted-issuers.json" {
		t.Fatalf("TrustStorePath = %q", got.TrustStorePath)
	}
}

func TestFromEnvOverridesDefaults(t *testing.T) {
	t.Setenv("HTTP_ADDR", "127.0.0.1:9090")
	t.Setenv("REGISTRY_URL", "http://registry.example.test")
	t.Setenv("TRUST_STORE_PATH", "/tmp/test-trust-store.json")

	got := FromEnv()
	if got.HTTPAddr != "127.0.0.1:9090" {
		t.Fatalf("HTTPAddr = %q", got.HTTPAddr)
	}
	if got.RegistryURL != "http://registry.example.test" {
		t.Fatalf("RegistryURL = %q", got.RegistryURL)
	}
	if got.TrustStorePath != "/tmp/test-trust-store.json" {
		t.Fatalf("TrustStorePath = %q", got.TrustStorePath)
	}
}
