// Package evidence verifies the signature envelope used by the Interviewtopia
// evidence producers. It deliberately does not decide whether the issuer is
// authorised for a particular property transaction; that is a business rule.
package evidence

import (
	"bytes"
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

var (
	ErrMalformedEnvelope = errors.New("malformed evidence envelope")
	ErrUnknownKey        = errors.New("unknown evidence key")
	ErrIssuerMismatch    = errors.New("evidence issuer does not own key")
	ErrInvalidSignature  = errors.New("invalid evidence signature")
)

type envelope struct {
	IssuerID  string          `json:"issuerId"`
	KeyID     string          `json:"keyId"`
	Payload   json.RawMessage `json:"payload"`
	Signature string          `json:"signature"`
}

type VerifiedEnvelope struct {
	IssuerID string
	KeyID    string
	Payload  json.RawMessage
}

type TrustedKey struct {
	IssuerID  string
	KeyID     string
	PublicKey ed25519.PublicKey
}

type TrustStore map[string]TrustedKey

type trustFile struct {
	Keys []struct {
		IssuerID        string `json:"issuerId"`
		KeyID           string `json:"keyId"`
		PublicKeyBase64 string `json:"publicKeyBase64"`
	} `json:"keys"`
}

func LoadTrustStore(data []byte) (TrustStore, error) {
	var file trustFile
	if err := decodeOne(data, &file); err != nil {
		return nil, fmt.Errorf("decode trust store: %w", err)
	}

	store := make(TrustStore, len(file.Keys))
	for _, item := range file.Keys {
		decoded, err := base64.StdEncoding.DecodeString(item.PublicKeyBase64)
		if err != nil || len(decoded) != ed25519.PublicKeySize {
			return nil, fmt.Errorf("key %q: invalid public key", item.KeyID)
		}
		if item.IssuerID == "" || item.KeyID == "" {
			return nil, errors.New("trust store contains an empty issuer or key ID")
		}
		if _, exists := store[item.KeyID]; exists {
			return nil, fmt.Errorf("duplicate key ID %q", item.KeyID)
		}
		store[item.KeyID] = TrustedKey{
			IssuerID:  item.IssuerID,
			KeyID:     item.KeyID,
			PublicKey: append(ed25519.PublicKey(nil), decoded...),
		}
	}
	return store, nil
}

func Verify(data []byte, store TrustStore) (VerifiedEnvelope, error) {
	var candidate envelope
	if err := decodeOne(data, &candidate); err != nil {
		return VerifiedEnvelope{}, fmt.Errorf("%w: %v", ErrMalformedEnvelope, err)
	}
	if candidate.IssuerID == "" || candidate.KeyID == "" || len(candidate.Payload) == 0 || candidate.Signature == "" {
		return VerifiedEnvelope{}, fmt.Errorf("%w: required field missing", ErrMalformedEnvelope)
	}

	key, ok := store[candidate.KeyID]
	if !ok {
		return VerifiedEnvelope{}, ErrUnknownKey
	}
	if key.IssuerID != candidate.IssuerID {
		return VerifiedEnvelope{}, ErrIssuerMismatch
	}

	signature, err := base64.StdEncoding.DecodeString(candidate.Signature)
	if err != nil || len(signature) != ed25519.SignatureSize {
		return VerifiedEnvelope{}, ErrInvalidSignature
	}
	var signedPayload bytes.Buffer
	if err := json.Compact(&signedPayload, candidate.Payload); err != nil {
		return VerifiedEnvelope{}, fmt.Errorf("%w: invalid payload JSON", ErrMalformedEnvelope)
	}
	if !ed25519.Verify(key.PublicKey, signedPayload.Bytes(), signature) {
		return VerifiedEnvelope{}, ErrInvalidSignature
	}

	return VerifiedEnvelope{
		IssuerID: candidate.IssuerID,
		KeyID:    candidate.KeyID,
		Payload:  append(json.RawMessage(nil), candidate.Payload...),
	}, nil
}

func decodeOne(data []byte, destination any) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(destination); err != nil {
		return err
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		if err == nil {
			return errors.New("multiple JSON values")
		}
		return err
	}
	return nil
}
