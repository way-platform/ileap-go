package ileapauthserver

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/way-platform/ileap-go"
	"golang.org/x/oauth2"
)

type mockTokenIssuer struct{}

func (m *mockTokenIssuer) IssueToken(
	_ context.Context, clientID, clientSecret string,
) (*oauth2.Token, error) {
	if clientID == "hello" && clientSecret == "pathfinder" {
		return &oauth2.Token{AccessToken: "mock-token", TokenType: "bearer"}, nil
	}
	return nil, ErrInvalidCredentials
}

type mockOIDCProvider struct{}

func (m *mockOIDCProvider) OpenIDConfiguration(baseURL string) *OpenIDConfiguration {
	return &OpenIDConfiguration{
		IssuerURL:              baseURL,
		AuthURL:                baseURL + "/auth/token",
		TokenURL:               baseURL + "/auth/token",
		JWKSURL:                baseURL + "/jwks",
		Algorithms:             []string{"RS256"},
		ResponseTypesSupported: []string{"token"},
		SubjectTypesSupported:  []string{"public"},
	}
}

func (m *mockOIDCProvider) JWKS() *JWKSet {
	return &JWKSet{
		Keys: []JWK{{
			KeyType:   "RSA",
			Use:       "sig",
			Algorithm: "RS256",
			KeyID:     "test",
			N:         "abc",
			E:         "AQAB",
		}},
	}
}

func newTestServer() *Server {
	return NewServer(&mockTokenIssuer{}, &mockOIDCProvider{})
}

func TestAuthToken(t *testing.T) {
	srv := newTestServer()

	t.Run("success", func(t *testing.T) {
		req := httptest.NewRequest(
			"POST",
			"/auth/token",
			strings.NewReader("grant_type=client_credentials"),
		)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.SetBasicAuth("hello", "pathfinder")
		w := httptest.NewRecorder()
		srv.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
		}
		var creds oauth2.Token
		if err := json.NewDecoder(w.Body).Decode(&creds); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if creds.AccessToken != "mock-token" {
			t.Errorf("expected mock-token, got %s", creds.AccessToken)
		}
		if creds.TokenType != "bearer" {
			t.Errorf("expected bearer, got %s", creds.TokenType)
		}
	})

	t.Run("invalid content type", func(t *testing.T) {
		req := httptest.NewRequest(
			"POST",
			"/auth/token",
			strings.NewReader("grant_type=client_credentials"),
		)
		req.Header.Set("Content-Type", "application/json")
		req.SetBasicAuth("hello", "pathfinder")
		w := httptest.NewRecorder()
		srv.ServeHTTP(w, req)
		checkOAuthError(t, w, http.StatusBadRequest, ileap.OAuthErrorCodeInvalidRequest)
	})

	t.Run("unsupported grant type", func(t *testing.T) {
		req := httptest.NewRequest(
			"POST",
			"/auth/token",
			strings.NewReader("grant_type=authorization_code"),
		)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.SetBasicAuth("hello", "pathfinder")
		w := httptest.NewRecorder()
		srv.ServeHTTP(w, req)
		checkOAuthError(t, w, http.StatusBadRequest, ileap.OAuthErrorCodeUnsupportedGrantType)
	})

	t.Run("missing basic auth", func(t *testing.T) {
		req := httptest.NewRequest(
			"POST",
			"/auth/token",
			strings.NewReader("grant_type=client_credentials"),
		)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()
		srv.ServeHTTP(w, req)
		checkOAuthError(t, w, http.StatusBadRequest, ileap.OAuthErrorCodeInvalidRequest)
	})

	t.Run("invalid credentials", func(t *testing.T) {
		req := httptest.NewRequest(
			"POST",
			"/auth/token",
			strings.NewReader("grant_type=client_credentials"),
		)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.SetBasicAuth("bad", "creds")
		w := httptest.NewRecorder()
		srv.ServeHTTP(w, req)
		checkOAuthError(t, w, http.StatusBadRequest, ileap.OAuthErrorCodeInvalidRequest)
	})
}

func TestOpenIDConfig(t *testing.T) {
	srv := newTestServer()
	req := httptest.NewRequest("GET", "/.well-known/openid-configuration", nil)
	req.Host = "localhost:8080"
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var cfg OpenIDConfiguration
	if err := json.NewDecoder(w.Body).Decode(&cfg); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if cfg.IssuerURL != "http://localhost:8080" {
		t.Errorf("expected issuer http://localhost:8080, got %s", cfg.IssuerURL)
	}
	if cfg.JWKSURL != "http://localhost:8080/jwks" {
		t.Errorf("expected jwks_uri http://localhost:8080/jwks, got %s", cfg.JWKSURL)
	}
}

func TestJWKS(t *testing.T) {
	srv := newTestServer()
	req := httptest.NewRequest("GET", "/jwks", nil)
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var jwks JWKSet
	if err := json.NewDecoder(w.Body).Decode(&jwks); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(jwks.Keys) != 1 {
		t.Fatalf("expected 1 key, got %d", len(jwks.Keys))
	}
	if jwks.Keys[0].KeyID != "test" {
		t.Errorf("expected kid test, got %s", jwks.Keys[0].KeyID)
	}
}

func checkOAuthError(
	t *testing.T,
	w *httptest.ResponseRecorder,
	expectedStatus int,
	expectedCode ileap.OAuthErrorCode,
) {
	t.Helper()
	if w.Code != expectedStatus {
		t.Fatalf("expected status %d, got %d: %s", expectedStatus, w.Code, w.Body.String())
	}
	ct := w.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", ct)
	}
	var errResp ileap.OAuthError
	if err := json.NewDecoder(w.Body).Decode(&errResp); err != nil {
		t.Fatalf("decode OAuth error: %v", err)
	}
	if errResp.Code != expectedCode {
		t.Errorf("expected OAuth error code %s, got %s", expectedCode, errResp.Code)
	}
	if errResp.Description == "" {
		t.Error("expected non-empty OAuth error description")
	}
}
