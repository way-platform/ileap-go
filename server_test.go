package ileap

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	ileapv1 "github.com/way-platform/ileap-go/proto/gen/wayplatform/connect/ileap/v1"
	"golang.org/x/oauth2"
	"google.golang.org/protobuf/encoding/protojson"
)

type mockTokenValidator struct {
	valid   bool
	expired bool
}

func (m *mockTokenValidator) ValidateToken(_ context.Context, _ string) (*TokenInfo, error) {
	if m.expired {
		return nil, fmt.Errorf("token expired: %w", ErrTokenExpired)
	}
	if !m.valid {
		return nil, fmt.Errorf("invalid token")
	}
	return &TokenInfo{Subject: "test-user"}, nil
}

type mockFootprintHandler struct {
	footprints []*ileapv1.ProductFootprint
}

func (m *mockFootprintHandler) GetFootprint(
	_ context.Context, id string,
) (*ileapv1.ProductFootprint, error) {
	for _, fp := range m.footprints {
		if fp.GetId() == id {
			return fp, nil
		}
	}
	return nil, ErrNotFound
}

func (m *mockFootprintHandler) ListFootprints(
	_ context.Context, req ListFootprintsRequest,
) (*ListFootprintsResponse, error) {
	result := m.footprints
	total := len(result)
	if req.Offset > 0 {
		if req.Offset >= len(result) {
			result = nil
		} else {
			result = result[req.Offset:]
		}
	}
	if req.Limit > 0 && len(result) > req.Limit {
		result = result[:req.Limit]
	}
	return &ListFootprintsResponse{Data: result, Total: total}, nil
}

type mockTADHandler struct {
	tads []*ileapv1.TAD
}

func (m *mockTADHandler) ListTADs(
	_ context.Context,
	req ListTADsRequest,
) (*ListTADsResponse, error) {
	result := m.tads
	total := len(result)
	if req.Offset > 0 {
		if req.Offset >= len(result) {
			result = nil
		} else {
			result = result[req.Offset:]
		}
	}
	if req.Limit > 0 && len(result) > req.Limit {
		result = result[:req.Limit]
	}
	return &ListTADsResponse{Data: result, Total: total}, nil
}

type mockEventHandler struct {
	lastEvent *Event
}

func (m *mockEventHandler) HandleEvent(_ context.Context, event Event) error {
	m.lastEvent = &event
	return nil
}

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
	return NewServer(
		WithTokenValidator(&mockTokenValidator{valid: true}),
		WithFootprintHandler(&mockFootprintHandler{
			footprints: []*ileapv1.ProductFootprint{
				func() *ileapv1.ProductFootprint { p := &ileapv1.ProductFootprint{}; p.SetId("fp-1"); return p }(),
				func() *ileapv1.ProductFootprint { p := &ileapv1.ProductFootprint{}; p.SetId("fp-2"); return p }(),
			},
		}),
		WithTADHandler(&mockTADHandler{
			tads: []*ileapv1.TAD{
				func() *ileapv1.TAD { t := &ileapv1.TAD{}; t.SetActivityId("tad-1"); return t }(),
			},
		}),
		WithEventHandler(&mockEventHandler{}),
	)
}

func authTestServer(opts ...ServerOption) *Server {
	base := []ServerOption{
		WithTokenIssuer(&mockTokenIssuer{}),
		WithOIDCProvider(&mockOIDCProvider{}),
	}
	return NewServer(append(base, opts...)...)
}

func TestAuthToken(t *testing.T) {
	srv := authTestServer()

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
		checkOAuthError(t, w, http.StatusBadRequest, OAuthErrorCodeInvalidRequest)
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
		checkOAuthError(t, w, http.StatusBadRequest, OAuthErrorCodeUnsupportedGrantType)
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
		checkOAuthError(t, w, http.StatusBadRequest, OAuthErrorCodeInvalidRequest)
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
		checkOAuthError(t, w, http.StatusBadRequest, OAuthErrorCodeInvalidRequest)
	})
}

func TestOpenIDConfig(t *testing.T) {
	srv := authTestServer()
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
	srv := authTestServer()
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

func TestWithPathPrefix(t *testing.T) {
	t.Run("OIDC discovery uses configured path prefix", func(t *testing.T) {
		srv := authTestServer(WithPathPrefix("/ileap"))
		req := httptest.NewRequest("GET", "/ileap/.well-known/openid-configuration", nil)
		req.Host = "api.example.com"
		req.Header.Set("X-Forwarded-Proto", "https")
		w := httptest.NewRecorder()
		srv.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
		}
		var cfg OpenIDConfiguration
		if err := json.NewDecoder(w.Body).Decode(&cfg); err != nil {
			t.Fatalf("decode: %v", err)
		}
		wantToken := "https://api.example.com/ileap/auth/token"
		if cfg.TokenURL != wantToken {
			t.Errorf("TokenURL = %q, want %q", cfg.TokenURL, wantToken)
		}
		wantJWKS := "https://api.example.com/ileap/jwks"
		if cfg.JWKSURL != wantJWKS {
			t.Errorf("JWKSURL = %q, want %q", cfg.JWKSURL, wantJWKS)
		}
	})

	t.Run("normalization", func(t *testing.T) {
		srv := authTestServer(WithPathPrefix("ileap/"))
		req := httptest.NewRequest("GET", "/ileap/.well-known/openid-configuration", nil)
		req.Host = "api.example.com"
		w := httptest.NewRecorder()
		srv.ServeHTTP(w, req)
		var cfg OpenIDConfiguration
		_ = json.NewDecoder(w.Body).Decode(&cfg)
		if cfg.TokenURL != "http://api.example.com/ileap/auth/token" {
			t.Errorf("TokenURL = %q, want http://api.example.com/ileap/auth/token", cfg.TokenURL)
		}
	})

	t.Run("pagination Link uses prefix", func(t *testing.T) {
		srv := NewServer(
			WithTokenValidator(&mockTokenValidator{valid: true}),
			WithFootprintHandler(&mockFootprintHandler{
				footprints: []*ileapv1.ProductFootprint{
					func() *ileapv1.ProductFootprint { p := &ileapv1.ProductFootprint{}; p.SetId("fp-1"); return p }(),
					func() *ileapv1.ProductFootprint { p := &ileapv1.ProductFootprint{}; p.SetId("fp-2"); return p }(),
					func() *ileapv1.ProductFootprint { p := &ileapv1.ProductFootprint{}; p.SetId("fp-3"); return p }(),
				},
			}),
			WithPathPrefix("/ileap"),
		)
		req := httptest.NewRequest("GET", "/ileap/2/footprints?limit=2", nil)
		req.Header.Set("Authorization", "Bearer valid")
		req.Host = "api.example.com"
		req.Header.Set("X-Forwarded-Proto", "https")
		w := httptest.NewRecorder()
		srv.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
		}
		got := w.Header().Get("Link")
		want := `<https://api.example.com/ileap/2/footprints?limit=2&offset=2>; rel="next"`
		if got != want {
			t.Errorf("Link = %q, want %q", got, want)
		}
	})
}

func TestPACTAuthMiddleware(t *testing.T) {
	srv := newTestServer()

	t.Run("missing authorization", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/2/footprints", nil)
		w := httptest.NewRecorder()
		srv.ServeHTTP(w, req)
		checkErrorResponse(t, w, http.StatusBadRequest, ErrorCodeBadRequest)
	})

	t.Run("non-bearer scheme", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/2/footprints", nil)
		req.Header.Set("Authorization", "Basic abc")
		w := httptest.NewRecorder()
		srv.ServeHTTP(w, req)
		checkErrorResponse(t, w, http.StatusBadRequest, ErrorCodeBadRequest)
	})

	t.Run("empty bearer token", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/2/footprints", nil)
		req.Header.Set("Authorization", "Bearer ")
		w := httptest.NewRecorder()
		srv.ServeHTTP(w, req)
		checkErrorResponse(t, w, http.StatusBadRequest, ErrorCodeBadRequest)
	})

	t.Run("invalid token returns 401", func(t *testing.T) {
		srv := NewServer(
			WithTokenValidator(&mockTokenValidator{valid: false}),
			WithFootprintHandler(&mockFootprintHandler{}),
		)
		req := httptest.NewRequest("GET", "/2/footprints", nil)
		req.Header.Set("Authorization", "Bearer bad-token")
		w := httptest.NewRecorder()
		srv.ServeHTTP(w, req)
		checkErrorResponse(t, w, http.StatusUnauthorized, ErrorCodeAccessDenied)
	})
}

func TestILeapAuthMiddleware(t *testing.T) {
	srv := newTestServer()

	t.Run("missing authorization returns 403", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/2/ileap/tad", nil)
		w := httptest.NewRecorder()
		srv.ServeHTTP(w, req)
		checkErrorResponse(t, w, http.StatusForbidden, ErrorCodeAccessDenied)
	})

	t.Run("invalid token returns 403", func(t *testing.T) {
		srv := NewServer(
			WithTokenValidator(&mockTokenValidator{valid: false}),
			WithTADHandler(&mockTADHandler{}),
		)
		req := httptest.NewRequest("GET", "/2/ileap/tad", nil)
		req.Header.Set("Authorization", "Bearer bad-token")
		w := httptest.NewRecorder()
		srv.ServeHTTP(w, req)
		checkErrorResponse(t, w, http.StatusForbidden, ErrorCodeAccessDenied)
	})

	t.Run("expired token returns 401", func(t *testing.T) {
		srv := NewServer(
			WithTokenValidator(&mockTokenValidator{expired: true}),
			WithTADHandler(&mockTADHandler{}),
		)
		req := httptest.NewRequest("GET", "/2/ileap/tad", nil)
		req.Header.Set("Authorization", "Bearer expired-token")
		w := httptest.NewRecorder()
		srv.ServeHTTP(w, req)
		checkErrorResponse(t, w, http.StatusUnauthorized, ErrorCodeTokenExpired)
	})
}

func TestListFootprints(t *testing.T) {
	srv := newTestServer()

	t.Run("success", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/2/footprints", nil)
		req.Header.Set("Authorization", "Bearer valid")
		w := httptest.NewRecorder()
		srv.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
		}
		var resp struct {
			Data []json.RawMessage `json:"data"`
		}
		if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if len(resp.Data) != 2 {
			t.Errorf("expected 2 footprints, got %d", len(resp.Data))
		}
	})

	t.Run("invalid limit", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/2/footprints?limit=abc", nil)
		req.Header.Set("Authorization", "Bearer valid")
		w := httptest.NewRecorder()
		srv.ServeHTTP(w, req)
		checkErrorResponse(t, w, http.StatusBadRequest, ErrorCodeBadRequest)
	})
}

func TestListFootprintsPagination(t *testing.T) {
	srv := NewServer(
		WithTokenValidator(&mockTokenValidator{valid: true}),
		WithFootprintHandler(&mockFootprintHandler{
			footprints: []*ileapv1.ProductFootprint{
				func() *ileapv1.ProductFootprint { p := &ileapv1.ProductFootprint{}; p.SetId("fp-1"); return p }(),
				func() *ileapv1.ProductFootprint { p := &ileapv1.ProductFootprint{}; p.SetId("fp-2"); return p }(),
				func() *ileapv1.ProductFootprint { p := &ileapv1.ProductFootprint{}; p.SetId("fp-3"); return p }(),
			},
		}),
	)

	t.Run("link header on first page", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/2/footprints?limit=2", nil)
		req.Header.Set("Authorization", "Bearer valid")
		req.Host = "example.com"
		req.Header.Set("X-Forwarded-Proto", "https")
		w := httptest.NewRecorder()
		srv.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
		}
		got := w.Header().Get("Link")
		want := `<https://example.com/2/footprints?limit=2&offset=2>; rel="next"`
		if got != want {
			t.Errorf("Link = %q, want %q", got, want)
		}
	})

	t.Run("no link header on last page", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/2/footprints?limit=2&offset=2", nil)
		req.Header.Set("Authorization", "Bearer valid")
		req.Host = "example.com"
		req.Header.Set("X-Forwarded-Proto", "https")
		w := httptest.NewRecorder()
		srv.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
		}
		if w.Header().Get("Link") != "" {
			t.Errorf("expected no Link header on last page")
		}
	})
}

func TestGetFootprint(t *testing.T) {
	srv := newTestServer()

	t.Run("found", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/2/footprints/fp-1", nil)
		req.Header.Set("Authorization", "Bearer valid")
		w := httptest.NewRecorder()
		srv.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
		}
		var resp struct {
			Data json.RawMessage `json:"data"`
		}
		if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
			t.Fatalf("decode: %v", err)
		}
		pf := &ileapv1.ProductFootprint{}
		if err := protojson.Unmarshal(resp.Data, pf); err != nil {
			t.Fatalf("unmarshal footprint: %v", err)
		}
		if pf.GetId() != "fp-1" {
			t.Errorf("expected fp-1, got %s", pf.GetId())
		}
	})

	t.Run("not found", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/2/footprints/nonexistent", nil)
		req.Header.Set("Authorization", "Bearer valid")
		w := httptest.NewRecorder()
		srv.ServeHTTP(w, req)
		checkErrorResponse(t, w, http.StatusNotFound, ErrorCodeNoSuchFootprint)
	})
}

func TestListTads(t *testing.T) {
	srv := NewServer(
		WithTokenValidator(&mockTokenValidator{valid: true}),
		WithTADHandler(&mockTADHandler{
			tads: []*ileapv1.TAD{
				func() *ileapv1.TAD { t := &ileapv1.TAD{}; t.SetActivityId("tad-1"); return t }(),
				func() *ileapv1.TAD { t := &ileapv1.TAD{}; t.SetActivityId("tad-2"); return t }(),
				func() *ileapv1.TAD { t := &ileapv1.TAD{}; t.SetActivityId("tad-3"); return t }(),
			},
		}),
	)

	t.Run("success", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/2/ileap/tad", nil)
		req.Header.Set("Authorization", "Bearer valid")
		w := httptest.NewRecorder()
		srv.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
		}
		var resp struct {
			Data []json.RawMessage `json:"data"`
		}
		if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if len(resp.Data) != 3 {
			t.Errorf("expected 3 TADs, got %d", len(resp.Data))
		}
	})

	t.Run("pagination link header", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/2/ileap/tad?limit=1", nil)
		req.Header.Set("Authorization", "Bearer valid")
		w := httptest.NewRecorder()
		srv.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
		}
		var resp struct {
			Data []json.RawMessage `json:"data"`
		}
		if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if len(resp.Data) != 1 {
			t.Errorf("expected 1 TAD, got %d", len(resp.Data))
		}
		link := w.Header().Get("Link")
		if link == "" {
			t.Fatal("expected Link header")
		}
		if !strings.Contains(link, `rel="next"`) {
			t.Errorf("expected rel=next in Link header, got %s", link)
		}
		if !strings.Contains(link, "offset=1") {
			t.Errorf("expected offset=1 in Link header, got %s", link)
		}
	})

	t.Run("no link header when all data returned", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/2/ileap/tad?limit=10", nil)
		req.Header.Set("Authorization", "Bearer valid")
		w := httptest.NewRecorder()
		srv.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
		}
		link := w.Header().Get("Link")
		if link != "" {
			t.Errorf("expected no Link header, got %s", link)
		}
	})

	t.Run("query params passed as filter", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/2/ileap/tad?mode=Road", nil)
		req.Header.Set("Authorization", "Bearer valid")
		w := httptest.NewRecorder()
		srv.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
		}
	})
}

func TestEvents(t *testing.T) {
	eh := &mockEventHandler{}
	srv := NewServer(
		WithTokenValidator(&mockTokenValidator{valid: true}),
		WithEventHandler(eh),
	)

	t.Run("cloudevents content type", func(t *testing.T) {
		body := `{"type":"org.wbcsd.pathfinder.ProductFootprint.Published.v1","specversion":"1.0","id":"evt-1","source":"test","data":{"pfIds":[]}}`
		req := httptest.NewRequest("POST", "/2/events", strings.NewReader(body))
		req.Header.Set("Authorization", "Bearer valid")
		req.Header.Set("Content-Type", "application/cloudevents+json")
		w := httptest.NewRecorder()
		srv.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
		}
		if eh.lastEvent == nil {
			t.Fatal("expected event to be handled")
		}
		if eh.lastEvent.ID != "evt-1" {
			t.Errorf("expected event ID evt-1, got %s", eh.lastEvent.ID)
		}
	})

	t.Run("application/json content type", func(t *testing.T) {
		body := `{"type":"org.wbcsd.pathfinder.ProductFootprint.Published.v1","specversion":"1.0","id":"evt-2","source":"test","data":{"pfIds":[]}}`
		req := httptest.NewRequest("POST", "/2/events", strings.NewReader(body))
		req.Header.Set("Authorization", "Bearer valid")
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		srv.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
		}
	})

	t.Run("missing content type", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/2/events", strings.NewReader("{}"))
		req.Header.Set("Authorization", "Bearer valid")
		w := httptest.NewRecorder()
		srv.ServeHTTP(w, req)
		checkErrorResponse(t, w, http.StatusBadRequest, ErrorCodeBadRequest)
	})

	t.Run("invalid content type", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/2/events", strings.NewReader("{}"))
		req.Header.Set("Authorization", "Bearer valid")
		req.Header.Set("Content-Type", "text/plain")
		w := httptest.NewRecorder()
		srv.ServeHTTP(w, req)
		checkErrorResponse(t, w, http.StatusBadRequest, ErrorCodeBadRequest)
	})
}

func TestEventsValidationMissingFields(t *testing.T) {
	srv := newTestServer()
	cases := []struct {
		name string
		body string
	}{
		{
			"missing specversion",
			`{"id":"1","source":"x","type":"org.wbcsd.pathfinder.ProductFootprint.Published.v1"}`,
		},
		{
			"missing id",
			`{"specversion":"1.0","source":"x","type":"org.wbcsd.pathfinder.ProductFootprint.Published.v1"}`,
		},
		{
			"missing source",
			`{"specversion":"1.0","id":"1","type":"org.wbcsd.pathfinder.ProductFootprint.Published.v1"}`,
		},
		{
			"missing data",
			`{"specversion":"1.0","id":"1","source":"x","type":"org.wbcsd.pathfinder.ProductFootprint.Published.v1"}`,
		},
		{
			"null data",
			`{"specversion":"1.0","id":"1","source":"x","type":"org.wbcsd.pathfinder.ProductFootprint.Published.v1","data":null}`,
		},
		{
			"wrong specversion",
			`{"specversion":"0.3","id":"1","source":"x","type":"org.wbcsd.pathfinder.ProductFootprint.Published.v1","data":{"pfIds":[]}}`,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/2/events", strings.NewReader(tc.body))
			req.Header.Set("Authorization", "Bearer valid")
			req.Header.Set("Content-Type", "application/cloudevents+json")
			w := httptest.NewRecorder()
			srv.ServeHTTP(w, req)
			checkErrorResponse(t, w, http.StatusBadRequest, ErrorCodeBadRequest)
		})
	}
}

func TestNotImplemented(t *testing.T) {
	t.Run("data handlers with auth configured", func(t *testing.T) {
		srv := NewServer(
			WithTokenValidator(&mockTokenValidator{valid: true}),
		)

		t.Run("footprints", func(t *testing.T) {
			req := httptest.NewRequest("GET", "/2/footprints", nil)
			req.Header.Set("Authorization", "Bearer valid")
			w := httptest.NewRecorder()
			srv.ServeHTTP(w, req)
			checkErrorResponse(t, w, http.StatusNotImplemented, ErrorCodeNotImplemented)
		})

		t.Run("tads", func(t *testing.T) {
			req := httptest.NewRequest("GET", "/2/ileap/tad", nil)
			req.Header.Set("Authorization", "Bearer valid")
			w := httptest.NewRecorder()
			srv.ServeHTTP(w, req)
			checkErrorResponse(t, w, http.StatusNotImplemented, ErrorCodeNotImplemented)
		})

		t.Run("events", func(t *testing.T) {
			body := `{"type":"org.wbcsd.pathfinder.ProductFootprint.Published.v1","specversion":"1.0","id":"evt-1","source":"test","data":{"pfIds":[]}}`
			req := httptest.NewRequest("POST", "/2/events", strings.NewReader(body))
			req.Header.Set("Authorization", "Bearer valid")
			req.Header.Set("Content-Type", "application/cloudevents+json")
			w := httptest.NewRecorder()
			srv.ServeHTTP(w, req)
			checkErrorResponse(t, w, http.StatusNotImplemented, ErrorCodeNotImplemented)
		})
	})

	t.Run("bare server returns 501 for all endpoints", func(t *testing.T) {
		srv := NewServer()

		t.Run("footprints with token", func(t *testing.T) {
			req := httptest.NewRequest("GET", "/2/footprints", nil)
			req.Header.Set("Authorization", "Bearer any-token")
			w := httptest.NewRecorder()
			srv.ServeHTTP(w, req)
			checkErrorResponse(t, w, http.StatusNotImplemented, ErrorCodeNotImplemented)
		})

		t.Run("tad with token", func(t *testing.T) {
			req := httptest.NewRequest("GET", "/2/ileap/tad", nil)
			req.Header.Set("Authorization", "Bearer any-token")
			w := httptest.NewRecorder()
			srv.ServeHTTP(w, req)
			checkErrorResponse(t, w, http.StatusNotImplemented, ErrorCodeNotImplemented)
		})

		t.Run("auth token", func(t *testing.T) {
			req := httptest.NewRequest(
				"POST",
				"/auth/token",
				strings.NewReader("grant_type=client_credentials"),
			)
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			req.SetBasicAuth("user", "pass")
			w := httptest.NewRecorder()
			srv.ServeHTTP(w, req)
			checkErrorResponse(t, w, http.StatusNotImplemented, ErrorCodeNotImplemented)
		})

		t.Run("openid configuration", func(t *testing.T) {
			req := httptest.NewRequest("GET", "/.well-known/openid-configuration", nil)
			w := httptest.NewRecorder()
			srv.ServeHTTP(w, req)
			checkErrorResponse(t, w, http.StatusNotImplemented, ErrorCodeNotImplemented)
		})

		t.Run("jwks", func(t *testing.T) {
			req := httptest.NewRequest("GET", "/jwks", nil)
			w := httptest.NewRecorder()
			srv.ServeHTTP(w, req)
			checkErrorResponse(t, w, http.StatusNotImplemented, ErrorCodeNotImplemented)
		})
	})
}

func checkErrorResponse(
	t *testing.T,
	w *httptest.ResponseRecorder,
	expectedStatus int,
	expectedCode ErrorCode,
) {
	t.Helper()
	if w.Code != expectedStatus {
		t.Fatalf("expected status %d, got %d: %s", expectedStatus, w.Code, w.Body.String())
	}
	ct := w.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", ct)
	}
	var errResp Error
	if err := json.NewDecoder(w.Body).Decode(&errResp); err != nil {
		t.Fatalf("decode error response: %v", err)
	}
	if errResp.Code != expectedCode {
		t.Errorf("expected error code %s, got %s", expectedCode, errResp.Code)
	}
	if errResp.Message == "" {
		t.Error("expected non-empty error message")
	}
}

func checkOAuthError(
	t *testing.T,
	w *httptest.ResponseRecorder,
	expectedStatus int,
	expectedCode OAuthErrorCode,
) {
	t.Helper()
	if w.Code != expectedStatus {
		t.Fatalf("expected status %d, got %d: %s", expectedStatus, w.Code, w.Body.String())
	}
	ct := w.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", ct)
	}
	var errResp OAuthError
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
