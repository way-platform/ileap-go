package ileapclerk

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/way-platform/ileap-go/ileapauthserver"
)

func TestSignIn(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		const wantJWT = "header.payload.signature"
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "POST" {
				t.Errorf("expected POST, got %s", r.Method)
			}
			if err := r.ParseForm(); err != nil {
				t.Fatal(err)
			}
			if r.Form.Get("strategy") != "password" {
				t.Errorf("expected strategy=password, got %s", r.Form.Get("strategy"))
			}
			if r.Form.Get("identifier") != "user@example.com" {
				t.Errorf("expected identifier=user@example.com, got %s", r.Form.Get("identifier"))
			}
			if r.Form.Get("password") != "secret" {
				t.Errorf("expected password=secret, got %s", r.Form.Get("password"))
			}
			w.Header().Set("Content-Type", "application/json")
			resp := signInResponse{}
			resp.Response.Status = "complete"
			resp.Client.Sessions = []struct {
				ID string `json:"id"`
				LastActiveToken struct {
					JWT string `json:"jwt"`
				} `json:"last_active_token"`
			}{{}}
			resp.Client.Sessions[0].LastActiveToken.JWT = wantJWT
			_ = json.NewEncoder(w).Encode(resp)
		}))
		defer srv.Close()
		c := NewClient("unused", WithHTTPClient(&http.Client{
			Transport: &testTransport{target: srv},
		}))
		jwt, err := c.SignIn("user@example.com", "secret", "")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if jwt != wantJWT {
			t.Errorf("expected JWT %q, got %q", wantJWT, jwt)
		}
	})

	t.Run("invalid credentials", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusUnprocessableEntity)
			_, _ = w.Write([]byte(`{"errors":[{"message":"Invalid credentials"}]}`))
		}))
		defer srv.Close()
		c := NewClient("unused", WithHTTPClient(&http.Client{
			Transport: &testTransport{target: srv},
		}))
		jwt, err := c.SignIn("bad@example.com", "wrong", "")
		if err == nil {
			t.Fatal("expected error for invalid credentials")
		}
		if jwt != "" {
			t.Errorf("expected empty JWT, got %q", jwt)
		}
	})

	t.Run("incomplete status", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			resp := signInResponse{}
			resp.Response.Status = "needs_second_factor"
			_ = json.NewEncoder(w).Encode(resp)
		}))
		defer srv.Close()
		c := NewClient("unused", WithHTTPClient(&http.Client{
			Transport: &testTransport{target: srv},
		}))
		jwt, err := c.SignIn("user@example.com", "secret", "")
		if err == nil {
			t.Fatal("expected error for incomplete sign-in")
		}
		if jwt != "" {
			t.Errorf("expected empty JWT, got %q", jwt)
		}
	})
}

func TestFetchJWKS(t *testing.T) {
	wantJWKS := ileapauthserver.JWKSet{
		Keys: []ileapauthserver.JWK{{
			KeyType: "RSA",
			Use:     "sig",
			KeyID:   "test-key-id",
			N:       "somenvalue",
			E:       "AQAB",
		}},
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/.well-known/jwks.json" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(wantJWKS)
	}))
	defer srv.Close()
	c := NewClient("unused", WithHTTPClient(&http.Client{
		Transport: &testTransport{target: srv},
	}))
	jwks, err := c.FetchJWKS()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(jwks.Keys) != 1 {
		t.Fatalf("expected 1 key, got %d", len(jwks.Keys))
	}
	if jwks.Keys[0].KeyID != "test-key-id" {
		t.Errorf("expected kid=test-key-id, got %s", jwks.Keys[0].KeyID)
	}
}

// testTransport redirects all HTTPS requests to the httptest server.
type testTransport struct {
	target *httptest.Server
}

func (t *testTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Rewrite the URL to point to the test server.
	req.URL.Scheme = "http"
	req.URL.Host = t.target.URL[len("http://"):]
	return http.DefaultTransport.RoundTrip(req)
}
