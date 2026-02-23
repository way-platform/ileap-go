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

	"connectrpc.com/connect"
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
		Keys: []ileap.JWK{
			{
				KeyType: "RSA",
				Use:     "sig",
				KeyID:   testKID,
				N:       base64.RawURLEncoding.EncodeToString(key.N.Bytes()),
				E: base64.RawURLEncoding.EncodeToString(
					new(big.Int).SetInt64(int64(key.PublicKey.E)).Bytes(),
				),
			},
		},
	}
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(jwks)
	}))
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

func testJWKS() ileap.JWKSet {
	return ileap.JWKSet{
		Keys: []ileap.JWK{{
			KeyType: "RSA",
			Use:     "sig",
			KeyID:   testKID,
			N:       "somenvalue",
			E:       "AQAB",
		}},
	}
}

func TestAuthHandler_IssueToken(t *testing.T) {
	const wantJWT = "header.payload.signature"

	t.Run("success", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			resp := signInResponse{}
			resp.Response.Status = "complete"
			resp.Client.Sessions = []struct {
				ID              string `json:"id"`
				LastActiveToken struct {
					JWT string `json:"jwt"`
				} `json:"last_active_token"`
			}{{}}
			resp.Client.Sessions[0].LastActiveToken.JWT = wantJWT
			_ = json.NewEncoder(w).Encode(resp)
		}))
		defer srv.Close()
		client := NewClient("unused", WithHTTPClient(&http.Client{
			Transport: &testTransport{target: srv},
		}))
		auth := NewAuthHandler(client)
		creds, err := auth.IssueToken(context.Background(), "user@example.com", "password")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if creds.TokenType != "bearer" {
			t.Errorf("expected bearer, got %s", creds.TokenType)
		}
		if creds.AccessToken != wantJWT {
			t.Errorf("expected access token %q, got %q", wantJWT, creds.AccessToken)
		}
	})

	t.Run("invalid credentials", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusUnprocessableEntity)
			_, _ = w.Write([]byte(`{"errors":[{"message":"Invalid credentials"}]}`))
		}))
		defer srv.Close()
		client := NewClient("unused", WithHTTPClient(&http.Client{
			Transport: &testTransport{target: srv},
		}))
		auth := NewAuthHandler(client)
		_, err := auth.IssueToken(context.Background(), "bad", "wrong")
		if err == nil {
			t.Fatal("expected error")
		}
		if connect.CodeOf(err) != connect.CodePermissionDenied {
			t.Errorf("expected CodePermissionDenied, got: %v", err)
		}
	})
}

func TestAuthHandler_ValidateToken(t *testing.T) {
	key := generateTestKey(t)

	t.Run("valid token", func(t *testing.T) {
		srv := jwksServerForKey(t, key)
		defer srv.Close()
		c := NewClient("unused", WithHTTPClient(&http.Client{
			Transport: &testTransport{target: srv},
		}))
		auth := NewAuthHandler(c)
		claims := map[string]any{
			"sub": "user@example.com",
			"exp": float64(time.Now().Add(time.Hour).Unix()),
		}
		token := makeTestJWT(t, key, claims)
		info, err := auth.ValidateToken(context.Background(), token)
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
		auth := NewAuthHandler(c)
		claims := map[string]any{
			"sub": "user@example.com",
			"exp": float64(time.Now().Add(-time.Hour).Unix()),
		}
		token := makeTestJWT(t, key, claims)
		_, err := auth.ValidateToken(context.Background(), token)
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
		auth := NewAuthHandler(c)
		claims := map[string]any{
			"sub": "user@example.com",
			"exp": float64(time.Now().Add(time.Hour).Unix()),
		}
		token := makeTestJWT(t, key, claims)
		parts := splitJWT(token)
		tampered := parts[0] + "." + parts[1] + ".invalidsignature"
		_, err := auth.ValidateToken(context.Background(), tampered)
		if err == nil {
			t.Fatal("expected error for tampered signature")
		}
	})
}

func TestAuthHandler_OpenIDConfiguration_JWKSURL(t *testing.T) {
	c := NewClient("unused")
	auth := NewAuthHandler(c)
	cfg := auth.OpenIDConfiguration("https://example.com")
	want := "https://example.com/jwks"
	if cfg.JWKSURL != want {
		t.Errorf("JWKSURL = %q, want %q", cfg.JWKSURL, want)
	}
}

func TestAuthHandler_JWKS_Caches(t *testing.T) {
	var callCount int
	jwks := testJWKS()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		callCount++
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(jwks)
	}))
	defer srv.Close()
	c := NewClient("unused", WithHTTPClient(&http.Client{
		Transport: &testTransport{target: srv},
	}))
	auth := NewAuthHandler(c)

	_ = auth.JWKS()
	_ = auth.JWKS()

	if callCount != 1 {
		t.Errorf("JWKS endpoint called %d times, want 1", callCount)
	}
}

func TestAuthHandler_JWKS_RefreshesAfterTTL(t *testing.T) {
	var callCount int
	jwks := testJWKS()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		callCount++
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(jwks)
	}))
	defer srv.Close()
	c := NewClient("unused", WithHTTPClient(&http.Client{
		Transport: &testTransport{target: srv},
	}))
	auth := NewAuthHandler(c, WithJWKSCacheTTL(1*time.Millisecond))

	_ = auth.JWKS()
	if callCount != 1 {
		t.Errorf("after first call: endpoint called %d times, want 1", callCount)
	}

	time.Sleep(2 * time.Millisecond)

	_ = auth.JWKS()
	if callCount != 2 {
		t.Errorf("after TTL expiry: endpoint called %d times, want 2", callCount)
	}
}

func TestAuthHandler_JWKS_StaleOnError(t *testing.T) {
	jwks := testJWKS()
	var returnError bool
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if returnError {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(jwks)
	}))
	defer srv.Close()
	c := NewClient("unused", WithHTTPClient(&http.Client{
		Transport: &testTransport{target: srv},
	}))
	auth := NewAuthHandler(c)

	// Prime the cache.
	got := auth.JWKS()
	if len(got.Keys) != 1 {
		t.Fatalf("expected 1 key from initial fetch, got %d", len(got.Keys))
	}

	// Expire the cache and make the endpoint fail.
	returnError = true
	auth.mu.Lock()
	auth.cachedAt = time.Time{}
	auth.mu.Unlock()

	got = auth.JWKS()
	if len(got.Keys) != 1 {
		t.Errorf("expected stale cache (1 key) on fetch error, got %d keys", len(got.Keys))
	}
}

func TestAuthHandler_JWKS_EmptyOnErrorWithNoCache(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()
	c := NewClient("unused", WithHTTPClient(&http.Client{
		Transport: &testTransport{target: srv},
	}))
	auth := NewAuthHandler(c)

	got := auth.JWKS()
	if len(got.Keys) != 0 {
		t.Errorf("expected empty JWKSet on error with no cache, got %d keys", len(got.Keys))
	}
}
