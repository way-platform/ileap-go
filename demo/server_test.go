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

func TestJWTCreationAndValidation(t *testing.T) {
	// Create a new server instance
	server, err := NewServer()
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	// Test JWT creation
	username := "hello"
	token, err := server.createJWT(username)
	if err != nil {
		t.Fatalf("Failed to create JWT: %v", err)
	}

	if token == "" {
		t.Fatal("Generated JWT is empty")
	}

	t.Logf("Generated JWT: %s", token)

	// Verify JWT format (should have 3 parts separated by dots)
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		t.Fatalf("Invalid JWT format: expected 3 parts, got %d", len(parts))
	}

	// Decode and verify header
	headerBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		t.Fatalf("Failed to decode JWT header: %v", err)
	}

	var header JWTHeader
	if err := json.Unmarshal(headerBytes, &header); err != nil {
		t.Fatalf("Failed to unmarshal JWT header: %v", err)
	}

	if header.Type != "JWT" {
		t.Errorf("Expected header type 'JWT', got '%s'", header.Type)
	}
	if header.Algorithm != "RS256" {
		t.Errorf("Expected algorithm 'RS256', got '%s'", header.Algorithm)
	}

	t.Logf("Header: %s", string(headerBytes))

	// Decode and verify payload
	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		t.Fatalf("Failed to decode JWT payload: %v", err)
	}

	var payload Claims
	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		t.Fatalf("Failed to unmarshal JWT payload: %v", err)
	}

	if payload.Username != username {
		t.Errorf("Expected username '%s', got '%s'", username, payload.Username)
	}
	if payload.IssuedAt == 0 {
		t.Error("Expected IssuedAt to be set")
	}

	t.Logf("Payload: %s", string(payloadBytes))

	// Test JWT validation
	validatedPayload, err := server.validateJWT(token)
	if err != nil {
		t.Fatalf("Failed to validate JWT: %v", err)
	}

	if validatedPayload.Username != username {
		t.Errorf("Expected validated username '%s', got '%s'", username, validatedPayload.Username)
	}

	t.Logf("Validated username: %s", validatedPayload.Username)
}

func TestJWTValidationFailures(t *testing.T) {
	server, err := NewServer()
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}
	testCases := []struct {
		name  string
		token string
	}{
		{
			name:  "empty token",
			token: "",
		},
		{
			name:  "invalid format - too few parts",
			token: "header.payload",
		},
		{
			name:  "invalid format - too many parts",
			token: "header.payload.signature.extra",
		},
		{
			name:  "invalid base64 encoding",
			token: "invalid_base64.payload.signature",
		},
		{
			name:  "tampered payload",
			token: "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJ1c2VybmFtZSI6InRhbXBlcmVkIn0.signature",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := server.validateJWT(tc.token)
			if err == nil {
				t.Errorf("Expected validation to fail for %s, but it succeeded", tc.name)
			}
			t.Logf("Validation correctly failed for %s: %v", tc.name, err)
		})
	}
}

func TestJWTRoundTrip(t *testing.T) {
	server, err := NewServer()
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	usernames := []string{"hello", "martin", "test_user", "admin"}

	for _, username := range usernames {
		t.Run("username_"+username, func(t *testing.T) {
			// Create JWT
			token, err := server.createJWT(username)
			if err != nil {
				t.Fatalf("Failed to create JWT for username '%s': %v", username, err)
			}

			// Validate JWT
			payload, err := server.validateJWT(token)
			if err != nil {
				t.Fatalf("Failed to validate JWT for username '%s': %v", username, err)
			}

			// Verify username matches
			if payload.Username != username {
				t.Errorf("Username mismatch: expected '%s', got '%s'", username, payload.Username)
			}
		})
	}
}

func TestAuthTokenEndpoint(t *testing.T) {
	server, err := NewServer()
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	// Test valid credentials
	req := httptest.NewRequest("POST", "/auth/token", strings.NewReader("grant_type=client_credentials"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth("hello", "pathfinder") // Using demo credentials from users.go

	w := httptest.NewRecorder()
	server.Handler().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var credentials ileap.ClientCredentials
	if err := json.NewDecoder(w.Body).Decode(&credentials); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if credentials.AccessToken == "" {
		t.Fatal("Access token is empty")
	}
	if credentials.TokenType != "bearer" {
		t.Errorf("Expected token type 'bearer', got '%s'", credentials.TokenType)
	}

	t.Logf("Received access token: %s", credentials.AccessToken)

	// Verify the token can be validated
	payload, err := server.validateJWT(credentials.AccessToken)
	if err != nil {
		t.Fatalf("Failed to validate generated token: %v", err)
	}
	if payload.Username != "hello" {
		t.Errorf("Expected username 'hello', got '%s'", payload.Username)
	}
}

func TestAuthenticatedEndpoints(t *testing.T) {
	server, err := NewServer()
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	// First get a valid token
	token, err := server.createJWT("hello")
	if err != nil {
		t.Fatalf("Failed to create JWT: %v", err)
	}

	// Test authenticated endpoint with valid token
	req := httptest.NewRequest("GET", "/2/footprints", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	server.Handler().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	// Test authenticated endpoint without token
	req = httptest.NewRequest("GET", "/2/footprints", nil)
	w = httptest.NewRecorder()
	server.Handler().ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("Expected status 401, got %d", w.Code)
	}

	// Test authenticated endpoint with invalid token
	req = httptest.NewRequest("GET", "/2/footprints", nil)
	req.Header.Set("Authorization", "Bearer invalid.token.here")
	w = httptest.NewRecorder()
	server.Handler().ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("Expected status 401, got %d", w.Code)
	}
}
