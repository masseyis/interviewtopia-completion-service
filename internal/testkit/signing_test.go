package testkit

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/masseyis/interviewtopia-completion-service/pkg/evidence"
)

func TestSignProducesVerifiableEnvelope(t *testing.T) {
	payload := map[string]any{
		"messageId": "example-message",
		"type":      "example",
	}

	storeJSON, err := os.ReadFile(filepath.Join("..", "..", "config", "trusted-issuers.json"))
	if err != nil {
		t.Fatal(err)
	}
	store, err := evidence.LoadTrustStore(storeJSON)
	if err != nil {
		t.Fatal(err)
	}

	verified, err := evidence.Verify(Sign(t, "bank-7", payload), store)
	if err != nil {
		t.Fatalf("Verify(Sign()) error = %v", err)
	}
	if verified.IssuerID != "bank-7" {
		t.Fatalf("IssuerID = %q, want bank-7", verified.IssuerID)
	}
}
