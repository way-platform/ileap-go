package clerk

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSignIn(t *testing.T) {
	t.Run("success", func(t *testing.T) {
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
			_ = json.NewEncoder(w).Encode(signInResponse{Status: "complete"})
		}))
		defer srv.Close()
		// Use the test server's host as the FAPI domain.
		c := NewClient("unused", WithHTTPClient(srv.Client()))
		// Override the endpoint by creating a custom client that routes to test server.
		c.fapiDomain = srv.URL[len("http://"):]
		// The client constructs https:// URLs, but test server is http.
		// We need a different approach: use a custom HTTP transport.
		c.httpClient = srv.Client()
		// Actually, let's just test via the httptest server approach with a custom transport.
		// The test server URL won't work with https:// prefix. Let's use a custom RoundTripper.
		c.httpClient = &http.Client{
			Transport: &testTransport{target: srv},
		}
		err := c.SignIn("user@example.com", "secret")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
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
		err := c.SignIn("bad@example.com", "wrong")
		if err == nil {
			t.Fatal("expected error for invalid credentials")
		}
	})

	t.Run("incomplete status", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(signInResponse{Status: "needs_second_factor"})
		}))
		defer srv.Close()
		c := NewClient("unused", WithHTTPClient(&http.Client{
			Transport: &testTransport{target: srv},
		}))
		err := c.SignIn("user@example.com", "secret")
		if err == nil {
			t.Fatal("expected error for incomplete sign-in")
		}
	})
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
