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

func TestServer(t *testing.T) {
	server, err := NewServer("http://localhost:8080")
	if err != nil {
		t.Fatalf("failed to create server: %v", err)
	}

	t.Run("POST /auth/token", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			req := httptest.NewRequest(
				"POST",
				"/auth/token",
				strings.NewReader("grant_type=client_credentials"),
			)
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
		})

		t.Run("invalid content type", func(t *testing.T) {
			req := httptest.NewRequest(
				"POST",
				"/auth/token",
				strings.NewReader("grant_type=client_credentials"),
			)
			req.Header.Set("Content-Type", "application/json") // Wrong content type
			req.SetBasicAuth("hello", "pathfinder")
			w := httptest.NewRecorder()
			server.Handler().ServeHTTP(w, req)
			checkOAuthErrorResponse(t, w, http.StatusBadRequest, ileap.OAuthErrorCodeInvalidRequest)
		})

		t.Run("invalid request body", func(t *testing.T) {
			req := httptest.NewRequest("POST", "/auth/token", strings.NewReader("invalid%body"))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			req.SetBasicAuth("hello", "pathfinder")
			w := httptest.NewRecorder()
			server.Handler().ServeHTTP(w, req)
			checkOAuthErrorResponse(t, w, http.StatusBadRequest, ileap.OAuthErrorCodeInvalidRequest)
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
			server.Handler().ServeHTTP(w, req)
			checkOAuthErrorResponse(
				t,
				w,
				http.StatusBadRequest,
				ileap.OAuthErrorCodeUnsupportedGrantType,
			)
		})

		t.Run("missing basic auth", func(t *testing.T) {
			req := httptest.NewRequest(
				"POST",
				"/auth/token",
				strings.NewReader("grant_type=client_credentials"),
			)
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			// No basic auth set
			w := httptest.NewRecorder()
			server.Handler().ServeHTTP(w, req)
			checkOAuthErrorResponse(t, w, http.StatusBadRequest, ileap.OAuthErrorCodeInvalidRequest)
		})

		t.Run("invalid credentials", func(t *testing.T) {
			req := httptest.NewRequest(
				"POST",
				"/auth/token",
				strings.NewReader("grant_type=client_credentials"),
			)
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			req.SetBasicAuth("invalid", "credentials")
			w := httptest.NewRecorder()
			server.Handler().ServeHTTP(w, req)
			checkOAuthErrorResponse(t, w, http.StatusBadRequest, ileap.OAuthErrorCodeInvalidRequest)
		})
	})

	t.Run("GET /2/footprints", func(t *testing.T) {
		t.Run("authenticated", func(t *testing.T) {
			req := httptest.NewRequest("GET", "/2/footprints", nil)
			req.Header.Set("Authorization", "Bearer "+getAccessToken(t, server))
			w := httptest.NewRecorder()
			server.Handler().ServeHTTP(w, req)
			if w.Code != http.StatusOK {
				t.Fatalf("expected status 200, got %d: %s", w.Code, w.Body.String())
			}
			contentType := w.Header().Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("expected Content-Type 'application/json', got '%s'", contentType)
			}
		})

		t.Run("unauthenticated", func(t *testing.T) {
			req := httptest.NewRequest("GET", "/2/footprints", nil)
			w := httptest.NewRecorder()
			server.Handler().ServeHTTP(w, req)
			checkErrorResponse(t, w, http.StatusBadRequest, ileap.ErrorCodeBadRequest)
		})

		t.Run("invalid token", func(t *testing.T) {
			req := httptest.NewRequest("GET", "/2/footprints", nil)
			req.Header.Set("Authorization", "Bearer invalid.token.here")
			w := httptest.NewRecorder()
			server.Handler().ServeHTTP(w, req)
			checkErrorResponse(t, w, http.StatusBadRequest, ileap.ErrorCodeBadRequest)
		})

		t.Run("missing bearer prefix", func(t *testing.T) {
			req := httptest.NewRequest("GET", "/2/footprints", nil)
			req.Header.Set("Authorization", "Basic invalid-auth")
			w := httptest.NewRecorder()
			server.Handler().ServeHTTP(w, req)
			checkErrorResponse(t, w, http.StatusBadRequest, ileap.ErrorCodeBadRequest)
		})

		t.Run("empty token", func(t *testing.T) {
			req := httptest.NewRequest("GET", "/2/footprints", nil)
			req.Header.Set("Authorization", "Bearer ")
			w := httptest.NewRecorder()
			server.Handler().ServeHTTP(w, req)
			checkErrorResponse(t, w, http.StatusBadRequest, ileap.ErrorCodeBadRequest)
		})
	})

	t.Run("GET /2/footprints/{id}", func(t *testing.T) {
		t.Run("not found", func(t *testing.T) {
			req := httptest.NewRequest("GET", "/2/footprints/nonexistent-id", nil)
			req.Header.Set("Authorization", "Bearer "+getAccessToken(t, server))
			w := httptest.NewRecorder()
			server.Handler().ServeHTTP(w, req)
			checkErrorResponse(t, w, http.StatusNotFound, ileap.ErrorCodeNoSuchFootprint)
		})

		t.Run("unauthenticated", func(t *testing.T) {
			req := httptest.NewRequest("GET", "/2/footprints/some-id", nil)
			w := httptest.NewRecorder()
			server.Handler().ServeHTTP(w, req)
			checkErrorResponse(t, w, http.StatusBadRequest, ileap.ErrorCodeBadRequest)
		})
	})

	t.Run("GET /2/ileap/tad", func(t *testing.T) {
		t.Run("authenticated", func(t *testing.T) {
			req := httptest.NewRequest("GET", "/2/ileap/tad", nil)
			req.Header.Set("Authorization", "Bearer "+getAccessToken(t, server))
			w := httptest.NewRecorder()
			server.Handler().ServeHTTP(w, req)
			if w.Code != http.StatusOK {
				t.Fatalf("expected status 200, got %d: %s", w.Code, w.Body.String())
			}
			contentType := w.Header().Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("expected Content-Type 'application/json', got '%s'", contentType)
			}
		})

		t.Run("unauthenticated returns 403", func(t *testing.T) {
			req := httptest.NewRequest("GET", "/2/ileap/tad", nil)
			w := httptest.NewRecorder()
			server.Handler().ServeHTTP(w, req)
			checkErrorResponse(t, w, http.StatusForbidden, ileap.ErrorCodeAccessDenied)
		})

		t.Run("invalid token returns 403", func(t *testing.T) {
			req := httptest.NewRequest("GET", "/2/ileap/tad", nil)
			req.Header.Set("Authorization", "Bearer invalid.token.here")
			w := httptest.NewRecorder()
			server.Handler().ServeHTTP(w, req)
			checkErrorResponse(t, w, http.StatusForbidden, ileap.ErrorCodeAccessDenied)
		})
	})
}

func checkErrorResponse(
	t *testing.T,
	w *httptest.ResponseRecorder,
	expectedStatus int,
	expectedCode ileap.ErrorCode,
) {
	t.Helper()
	if w.Code != expectedStatus {
		t.Fatalf("expected status %d, got %d: %s", expectedStatus, w.Code, w.Body.String())
	}
	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected Content-Type 'application/json', got '%s'", contentType)
	}
	var errorResp ileap.Error
	if err := json.NewDecoder(w.Body).Decode(&errorResp); err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}
	if errorResp.Code != expectedCode {
		t.Errorf("expected error code '%s', got '%s'", expectedCode, errorResp.Code)
	}
	if errorResp.Message == "" {
		t.Error("expected non-empty error message")
	}
}

func checkOAuthErrorResponse(
	t *testing.T,
	w *httptest.ResponseRecorder,
	expectedStatus int,
	expectedCode ileap.OAuthErrorCode,
) {
	t.Helper()
	if w.Code != expectedStatus {
		t.Fatalf("expected status %d, got %d: %s", expectedStatus, w.Code, w.Body.String())
	}
	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected Content-Type 'application/json', got '%s'", contentType)
	}
	var errorResp ileap.OAuthError
	if err := json.NewDecoder(w.Body).Decode(&errorResp); err != nil {
		t.Fatalf("failed to decode OAuth error response: %v", err)
	}
	if errorResp.Code != expectedCode {
		t.Errorf("expected OAuth error code '%s', got '%s'", expectedCode, errorResp.Code)
	}
	if errorResp.Description == "" {
		t.Error("expected non-empty OAuth error description")
	}
}

func getAccessToken(t *testing.T, server *Server) string {
	t.Helper()
	tokenRequest := httptest.NewRequest(
		"POST",
		"/auth/token",
		strings.NewReader("grant_type=client_credentials"),
	)
	tokenRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	tokenRequest.SetBasicAuth("hello", "pathfinder")
	tokenResponse := httptest.NewRecorder()
	server.Handler().ServeHTTP(tokenResponse, tokenRequest)
	if tokenResponse.Code != http.StatusOK {
		t.Fatalf(
			"failed to get access token: %d: %s",
			tokenResponse.Code,
			tokenResponse.Body.String(),
		)
	}
	var credentials ileap.ClientCredentials
	if err := json.NewDecoder(tokenResponse.Body).Decode(&credentials); err != nil {
		t.Fatalf("failed to decode token response: %v", err)
	}
	return credentials.AccessToken
}
