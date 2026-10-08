// Package testkit provides synthetic signing and HTTP helpers for tests in this
// repository. Its deterministic private keys are public test data and must not
// be reused for any real system.
package testkit

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"testing"
)

type testIssuer struct {
	keyID string
	seed  string
}

var testIssuers = map[string]testIssuer{
	"agent-9":   {keyID: "agent-9-2026", seed: "MRHaMbcYK+aze47LNUcipwrHFa6iCU7E0m2zRphbT2M="},
	"bank-7":    {keyID: "bank-7-2026", seed: "oq7WdfqSuetUl4dHCDRl9rB8rhcGdtlDI2t/I4lVbAw="},
	"lawyer-17": {keyID: "lawyer-17-2026", seed: "E+kQBtHjFhmwDvr3rB05fUwU4IQ0eVdEq+6X4QcXExE="},
	"lawyer-31": {keyID: "lawyer-31-2026", seed: "bUNXjwM0Z4AZkMLTG/gfodYwD9lkeF/7OMFcINMh9ag="},
	"lawyer-99": {keyID: "lawyer-99-2026", seed: "pOM31QtV+ARb6AjKNHeeQkM94832AW9KKT0nFQbiJdg="},
}

type signedEnvelope struct {
	IssuerID  string          `json:"issuerId"`
	KeyID     string          `json:"keyId"`
	Payload   json.RawMessage `json:"payload"`
	Signature string          `json:"signature"`
}

// Sign returns a JSON evidence envelope signed by one of the synthetic issuers
// in config/trusted-issuers.json. It fails the calling test for an unknown
// issuer or a payload that cannot be encoded.
func Sign(t testing.TB, issuerID string, payload any) []byte {
	t.Helper()

	issuer, ok := testIssuers[issuerID]
	if !ok {
		t.Fatalf("testkit.Sign: unknown issuer %q", issuerID)
	}
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("testkit.Sign: encode payload: %v", err)
	}
	seed, err := base64.StdEncoding.DecodeString(issuer.seed)
	if err != nil {
		t.Fatalf("testkit.Sign: decode test seed: %v", err)
	}
	privateKey := ed25519.NewKeyFromSeed(seed)
	envelopeJSON, err := json.Marshal(signedEnvelope{
		IssuerID:  issuerID,
		KeyID:     issuer.keyID,
		Payload:   payloadJSON,
		Signature: base64.StdEncoding.EncodeToString(ed25519.Sign(privateKey, payloadJSON)),
	})
	if err != nil {
		t.Fatalf("testkit.Sign: encode envelope: %v", err)
	}
	return envelopeJSON
}
