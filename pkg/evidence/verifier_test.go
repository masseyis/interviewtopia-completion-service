package evidence

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestVerifySuppliedExample(t *testing.T) {
	t.Parallel()

	storeBytes, err := os.ReadFile(filepath.Join("..", "..", "config", "trusted-issuers.json"))
	if err != nil {
		t.Fatal(err)
	}
	store, err := LoadTrustStore(storeBytes)
	if err != nil {
		t.Fatal(err)
	}

	envelopeBytes, err := os.ReadFile(filepath.Join("..", "..", "examples", "evidence", "identity-confirmed.json"))
	if err != nil {
		t.Fatal(err)
	}
	verified, err := Verify(envelopeBytes, store)
	if err != nil {
		t.Fatalf("Verify() error = %v", err)
	}
	if verified.IssuerID != "agent-9" {
		t.Fatalf("IssuerID = %q, want agent-9", verified.IssuerID)
	}
}

func TestVerifyRejectsTampering(t *testing.T) {
	t.Parallel()

	storeBytes, err := os.ReadFile(filepath.Join("..", "..", "config", "trusted-issuers.json"))
	if err != nil {
		t.Fatal(err)
	}
	store, err := LoadTrustStore(storeBytes)
	if err != nil {
		t.Fatal(err)
	}

	envelopeBytes, err := os.ReadFile(filepath.Join("..", "..", "examples", "evidence", "identity-confirmed.json"))
	if err != nil {
		t.Fatal(err)
	}
	tampered := append([]byte(nil), envelopeBytes...)
	for index := range tampered {
		if tampered[index] == '4' {
			tampered[index] = '5'
			break
		}
	}
	_, err = Verify(tampered, store)
	if !errors.Is(err, ErrInvalidSignature) {
		t.Fatalf("Verify() error = %v, want ErrInvalidSignature", err)
	}
}
