package ileapclerk

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/way-platform/ileap-go"
)

const testKID = "test-key-id"

func generateTestKey(t *testing.T) *rsa.PrivateKey {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate RSA key: %v", err)
	}
	return key
}

func makeTestJWT(t *testing.T, key *rsa.PrivateKey, claims map[string]any) string {
	t.Helper()
	header := map[string]string{"typ": "JWT", "alg": "RS256", "kid": testKID}
	headerBytes, _ := json.Marshal(header)
	payloadBytes, _ := json.Marshal(claims)
	h := base64.RawURLEncoding.EncodeToString(headerBytes)
	p := base64.RawURLEncoding.EncodeToString(payloadBytes)
	message := h + "." + p
	digest := sha256.Sum256([]byte(message))
	sig, err := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, digest[:])
	if err != nil {
		t.Fatalf("sign JWT: %v", err)
	}
	return message + "." + base64.RawURLEncoding.EncodeToString(sig)
}

func jwksServerForKey(t *testing.T, key *rsa.PrivateKey) *httptest.Server {
	t.Helper()
	jwks := ileap.JWKSet{
		Keys: []ileap.JWK{{
			KeyType: "RSA",
			Use:     "sig",
			KeyID:   testKID,
			N: base64.RawURLEncoding.EncodeToString(
				key.N.Bytes(),
			),
			E: base64.RawURLEncoding.EncodeToString(
				new(big.Int).SetInt64(int64(key.PublicKey.E)).Bytes(),
			),
		}},
	}
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(jwks)
	}))
}

func TestTokenValidator(t *testing.T) {
	key := generateTestKey(t)

	t.Run("valid token", func(t *testing.T) {
		srv := jwksServerForKey(t, key)
		defer srv.Close()
		c := NewClient("unused", WithHTTPClient(&http.Client{
			Transport: &testTransport{target: srv},
		}))
		validator := NewTokenValidator(c)
		claims := map[string]any{
			"sub": "user@example.com",
			"exp": float64(time.Now().Add(time.Hour).Unix()),
		}
		token := makeTestJWT(t, key, claims)
		info, err := validator.ValidateToken(context.Background(), token)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if info.Subject != "user@example.com" {
			t.Errorf("expected subject user@example.com, got %s", info.Subject)
		}
	})

	t.Run("expired token", func(t *testing.T) {
		srv := jwksServerForKey(t, key)
		defer srv.Close()
		c := NewClient("unused", WithHTTPClient(&http.Client{
			Transport: &testTransport{target: srv},
		}))
		validator := NewTokenValidator(c)
		claims := map[string]any{
			"sub": "user@example.com",
			"exp": float64(time.Now().Add(-time.Hour).Unix()),
		}
		token := makeTestJWT(t, key, claims)
		_, err := validator.ValidateToken(context.Background(), token)
		if err == nil {
			t.Fatal("expected error for expired token")
		}
	})

	t.Run("tampered signature", func(t *testing.T) {
		srv := jwksServerForKey(t, key)
		defer srv.Close()
		c := NewClient("unused", WithHTTPClient(&http.Client{
			Transport: &testTransport{target: srv},
		}))
		validator := NewTokenValidator(c)
		claims := map[string]any{
			"sub": "user@example.com",
			"exp": float64(time.Now().Add(time.Hour).Unix()),
		}
		token := makeTestJWT(t, key, claims)
		// Tamper with the signature (last part).
		parts := splitJWT(token)
		tampered := parts[0] + "." + parts[1] + ".invalidsignature"
		_, err := validator.ValidateToken(context.Background(), tampered)
		if err == nil {
			t.Fatal("expected error for tampered signature")
		}
	})
}

func splitJWT(token string) [3]string {
	var parts [3]string
	i := 0
	start := 0
	for j := 0; j < len(token); j++ {
		if token[j] == '.' {
			parts[i] = token[start:j]
			i++
			start = j + 1
			if i == 2 {
				parts[2] = token[start:]
				break
			}
		}
	}
	return parts
}
