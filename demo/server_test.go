package demo

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/way-platform/ileap-go"
)

func TestServer_Route_AuthToken(t *testing.T) {
	server, err := NewServer("http://localhost:8080")
	if err != nil {
		t.Fatalf("failed to create server: %v", err)
	}
	req := httptest.NewRequest("POST", "/auth/token", strings.NewReader("grant_type=client_credentials"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth("hello", "pathfinder")
	w := httptest.NewRecorder()
	server.Handler().ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}
	var credentials ileap.ClientCredentials
	if err := json.NewDecoder(w.Body).Decode(&credentials); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if credentials.TokenType != "bearer" {
		t.Errorf("Expected token type 'bearer', got '%s'", credentials.TokenType)
	}
	parts := strings.Split(credentials.AccessToken, ".")
	if len(parts) != 3 {
		t.Fatalf("expected 3 parts in JWT, got %d", len(parts))
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		t.Fatalf("failed to decode payload: %v", err)
	}
	var claims JWTClaims
	if err := json.Unmarshal(payload, &claims); err != nil {
		t.Fatalf("failed to unmarshal payload: %v", err)
	}
	if claims.Username != "hello" {
		t.Errorf("expected username 'hello', got '%s'", claims.Username)
	}
}

func TestServer_Route_ListFootprints(t *testing.T) {
	server, err := NewServer("http://localhost:8080")
	if err != nil {
		t.Fatalf("failed to create server: %v", err)
	}

	t.Run("unauthenticated", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/2/footprints", nil)
		w := httptest.NewRecorder()
		server.Handler().ServeHTTP(w, req)
		expected := http.StatusUnauthorized
		if w.Code != expected {
			t.Fatalf("expected status %d, got %d", expected, w.Code)
		}
	})

	t.Run("invalid token", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/2/footprints", nil)
		req.Header.Set("Authorization", "Bearer invalid.token.here")
		w := httptest.NewRecorder()
		server.Handler().ServeHTTP(w, req)
		expected := http.StatusBadRequest
		if w.Code != expected {
			t.Fatalf("expected status %d, got %d", expected, w.Code)
		}
	})

	t.Run("authenticated", func(t *testing.T) {
		tokenRequest := httptest.NewRequest("POST", "/auth/token", strings.NewReader("grant_type=client_credentials"))
		tokenRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		tokenRequest.SetBasicAuth("hello", "pathfinder")
		tokenResponse := httptest.NewRecorder()
		server.Handler().ServeHTTP(tokenResponse, tokenRequest)
		if tokenResponse.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d: %s", tokenResponse.Code, tokenResponse.Body.String())
		}
		var credentials ileap.ClientCredentials
		if err := json.NewDecoder(tokenResponse.Body).Decode(&credentials); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		request := httptest.NewRequest("GET", "/2/footprints", nil)
		request.Header.Set("Authorization", "Bearer "+credentials.AccessToken)
		response := httptest.NewRecorder()
		server.Handler().ServeHTTP(response, request)
		if response.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d: %s", response.Code, response.Body.String())
		}
	})
}
