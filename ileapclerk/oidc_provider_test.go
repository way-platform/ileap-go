package ileapclerk

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/way-platform/ileap-go/ileapauthserver"
)

func testJWKS() ileapauthserver.JWKSet {
	return ileapauthserver.JWKSet{
		Keys: []ileapauthserver.JWK{{
			KeyType: "RSA",
			Use:     "sig",
			KeyID:   testKID,
			N:       "somenvalue",
			E:       "AQAB",
		}},
	}
}

func TestOIDCProvider_OpenIDConfiguration_JWKSURL(t *testing.T) {
	c := NewClient("unused")
	p := NewOIDCProvider(c)
	cfg := p.OpenIDConfiguration("https://example.com")
	want := "https://example.com/jwks"
	if cfg.JWKSURL != want {
		t.Errorf("JWKSURL = %q, want %q", cfg.JWKSURL, want)
	}
}

func TestOIDCProvider_JWKS_Caches(t *testing.T) {
	var callCount atomic.Int32
	jwks := testJWKS()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		callCount.Add(1)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(jwks)
	}))
	defer srv.Close()
	c := NewClient("unused", WithHTTPClient(&http.Client{
		Transport: &testTransport{target: srv},
	}))
	p := NewOIDCProvider(c)

	_ = p.JWKS()
	_ = p.JWKS()

	if got := callCount.Load(); got != 1 {
		t.Errorf("JWKS endpoint called %d times, want 1", got)
	}
}

func TestOIDCProvider_JWKS_RefreshesAfterTTL(t *testing.T) {
	var callCount atomic.Int32
	jwks := testJWKS()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		callCount.Add(1)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(jwks)
	}))
	defer srv.Close()
	c := NewClient("unused", WithHTTPClient(&http.Client{
		Transport: &testTransport{target: srv},
	}))
	p := NewOIDCProvider(c)

	_ = p.JWKS()
	if got := callCount.Load(); got != 1 {
		t.Errorf("after first call: endpoint called %d times, want 1", got)
	}

	// Simulate TTL expiry by backdating cachedAt.
	p.cachedAt = time.Now().Add(-(jwksCacheTTL + time.Second))

	_ = p.JWKS()
	if got := callCount.Load(); got != 2 {
		t.Errorf("after TTL expiry: endpoint called %d times, want 2", got)
	}
}

func TestOIDCProvider_JWKS_StaleOnError(t *testing.T) {
	jwks := testJWKS()
	var returnError atomic.Bool
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if returnError.Load() {
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
	p := NewOIDCProvider(c)

	// Prime the cache.
	got := p.JWKS()
	if len(got.Keys) != 1 {
		t.Fatalf("expected 1 key from initial fetch, got %d", len(got.Keys))
	}

	// Expire the cache and make the endpoint fail.
	p.cachedAt = time.Time{}
	returnError.Store(true)

	got = p.JWKS()
	if len(got.Keys) != 1 {
		t.Errorf("expected stale cache (1 key) on fetch error, got %d keys", len(got.Keys))
	}
}

func TestOIDCProvider_JWKS_EmptyOnErrorWithNoCache(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()
	c := NewClient("unused", WithHTTPClient(&http.Client{
		Transport: &testTransport{target: srv},
	}))
	p := NewOIDCProvider(c)

	got := p.JWKS()
	if len(got.Keys) != 0 {
		t.Errorf("expected empty JWKSet on error with no cache, got %d keys", len(got.Keys))
	}
}
