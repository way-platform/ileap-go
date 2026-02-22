package ileapclerk

import (
	"log/slog"
	"sync"
	"time"

	"github.com/way-platform/ileap-go/ileapserver"
)

const jwksCacheTTL = 15 * time.Minute

// OIDCProvider implements ileapserver.OIDCProvider using Clerk's JWKS.
type OIDCProvider struct {
	client     *Client
	mu         sync.RWMutex
	cachedAt   time.Time
	cachedJWKS *ileapserver.JWKSet
}

// NewOIDCProvider creates a new OIDC provider backed by the Clerk client.
func NewOIDCProvider(client *Client) *OIDCProvider {
	return &OIDCProvider{client: client}
}

// OpenIDConfiguration returns the OIDC configuration for the given base URL.
func (p *OIDCProvider) OpenIDConfiguration(baseURL string) *ileapserver.OpenIDConfiguration {
	return &ileapserver.OpenIDConfiguration{
		IssuerURL:              baseURL,
		AuthURL:                baseURL + "/auth/token",
		TokenURL:               baseURL + "/auth/token",
		JWKSURL:                baseURL + "/jwks",
		Algorithms:             []string{"RS256"},
		ResponseTypesSupported: []string{"token"},
		SubjectTypesSupported:  []string{"public"},
	}
}

// JWKS fetches and returns Clerk's JSON Web Key Set, with TTL-based caching.
func (p *OIDCProvider) JWKS() *ileapserver.JWKSet {
	// Fast path: serve from cache if fresh.
	p.mu.RLock()
	if p.cachedJWKS != nil && time.Since(p.cachedAt) < jwksCacheTTL {
		cached := p.cachedJWKS
		p.mu.RUnlock()
		return cached
	}
	p.mu.RUnlock()

	// Slow path: fetch and update cache.
	p.mu.Lock()
	defer p.mu.Unlock()
	// Re-check after acquiring write lock (another goroutine may have
	// refreshed between our RUnlock and Lock).
	if p.cachedJWKS != nil && time.Since(p.cachedAt) < jwksCacheTTL {
		return p.cachedJWKS
	}
	jwks, err := p.client.FetchJWKS()
	if err != nil {
		slog.Warn("failed to fetch JWKS from Clerk", "error", err)
		if p.cachedJWKS != nil {
			return p.cachedJWKS // serve stale on error
		}
		return &ileapserver.JWKSet{}
	}
	p.cachedJWKS = jwks
	p.cachedAt = time.Now()
	return p.cachedJWKS
}
