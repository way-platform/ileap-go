package ileapclerk

import (
	"context"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"math/big"
	"strings"
	"sync"
	"time"

	"connectrpc.com/connect"
	"github.com/way-platform/ileap-go"
	"golang.org/x/oauth2"
)

const defaultJWKSCacheTTL = 15 * time.Minute

var _ ileap.AuthHandler = (*AuthHandler)(nil)

// AuthHandler implements ileap.AuthHandler using Clerk for token issuance,
// validation, and OIDC discovery.
type AuthHandler struct {
	client       *Client
	activeOrgID  string
	mu           sync.RWMutex
	cachedJWKS   *ileap.JWKSet
	cachedAt     time.Time
	jwksCacheTTL time.Duration
}

// AuthHandlerOption configures the AuthHandler.
type AuthHandlerOption func(*AuthHandler)

// WithActiveOrganization sets the active organization ID for issued tokens.
func WithActiveOrganization(orgID string) AuthHandlerOption {
	return func(a *AuthHandler) { a.activeOrgID = orgID }
}

// WithJWKSCacheTTL sets the JWKS cache TTL. Used for testing. Default is 15 minutes.
func WithJWKSCacheTTL(d time.Duration) AuthHandlerOption {
	return func(a *AuthHandler) { a.jwksCacheTTL = d }
}

// NewAuthHandler creates an AuthHandler backed by the given Clerk client.
func NewAuthHandler(client *Client, opts ...AuthHandlerOption) *AuthHandler {
	a := &AuthHandler{
		client:       client,
		jwksCacheTTL: defaultJWKSCacheTTL,
	}
	for _, opt := range opts {
		opt(a)
	}
	return a
}

func (a *AuthHandler) IssueToken(
	_ context.Context, clientID, clientSecret string,
) (*oauth2.Token, error) {
	jwt, err := a.client.SignIn(clientID, clientSecret, a.activeOrgID)
	if err != nil {
		return nil, connect.NewError(
			connect.CodePermissionDenied,
			fmt.Errorf("invalid credentials: %w", err),
		)
	}
	return &oauth2.Token{
		AccessToken: jwt,
		TokenType:   "bearer",
	}, nil
}

func (a *AuthHandler) ValidateToken(
	_ context.Context, token string,
) (*ileap.TokenInfo, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid JWT format")
	}
	headerBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, fmt.Errorf("decode JWT header: %w", err)
	}
	var header struct {
		Kid string `json:"kid"`
		Alg string `json:"alg"`
	}
	if err := json.Unmarshal(headerBytes, &header); err != nil {
		return nil, fmt.Errorf("parse JWT header: %w", err)
	}
	pub, err := a.findKey(header.Kid)
	if err != nil {
		return nil, fmt.Errorf("find signing key: %w", err)
	}
	sigBytes, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return nil, fmt.Errorf("decode JWT signature: %w", err)
	}
	message := parts[0] + "." + parts[1]
	digest := sha256.Sum256([]byte(message))
	if err := rsa.VerifyPKCS1v15(pub, crypto.SHA256, digest[:], sigBytes); err != nil {
		return nil, fmt.Errorf("invalid JWT signature: %w", err)
	}
	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("decode JWT payload: %w", err)
	}
	var claims map[string]any
	if err := json.Unmarshal(payloadBytes, &claims); err != nil {
		return nil, fmt.Errorf("parse JWT claims: %w", err)
	}
	if exp, ok := claims["exp"].(float64); ok {
		if time.Now().Unix() > int64(exp) {
			return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("token expired"))
		}
	}
	sub, _ := claims["sub"].(string)
	return &ileap.TokenInfo{Subject: sub}, nil
}

func (a *AuthHandler) OpenIDConfiguration(baseURL string) *ileap.OpenIDConfiguration {
	return &ileap.OpenIDConfiguration{
		IssuerURL:              baseURL,
		AuthURL:                baseURL + "/auth/token",
		TokenURL:               baseURL + "/auth/token",
		JWKSURL:                baseURL + "/jwks",
		Algorithms:             []string{"RS256"},
		ResponseTypesSupported: []string{"token"},
		SubjectTypesSupported:  []string{"public"},
	}
}

func (a *AuthHandler) JWKS() *ileap.JWKSet {
	// Fast path: serve from cache if fresh.
	a.mu.RLock()
	if a.cachedJWKS != nil && time.Since(a.cachedAt) < a.jwksCacheTTL {
		cached := a.cachedJWKS
		a.mu.RUnlock()
		return cached
	}
	a.mu.RUnlock()

	// Slow path: fetch and update cache.
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.cachedJWKS != nil && time.Since(a.cachedAt) < a.jwksCacheTTL {
		return a.cachedJWKS
	}
	jwks, err := a.client.FetchJWKS()
	if err != nil {
		slog.Warn("failed to fetch JWKS from Clerk", "error", err)
		if a.cachedJWKS != nil {
			return a.cachedJWKS // serve stale on error
		}
		return &ileap.JWKSet{}
	}
	a.cachedJWKS = jwks
	a.cachedAt = time.Now()
	return a.cachedJWKS
}

func (a *AuthHandler) findKey(kid string) (*rsa.PublicKey, error) {
	jwks := a.JWKS()
	if pub := findKeyInSet(jwks, kid); pub != nil {
		return pub, nil
	}
	// Key not in cache (e.g. key rotation); fetch and retry.
	a.mu.Lock()
	jwks, err := a.client.FetchJWKS()
	if err != nil {
		a.mu.Unlock()
		return nil, fmt.Errorf("fetch JWKS: %w", err)
	}
	a.cachedJWKS = jwks
	a.cachedAt = time.Now()
	pub := findKeyInSet(jwks, kid)
	a.mu.Unlock()
	if pub == nil {
		return nil, fmt.Errorf("key %q not found in JWKS", kid)
	}
	return pub, nil
}

func findKeyInSet(jwks *ileap.JWKSet, kid string) *rsa.PublicKey {
	for _, jwk := range jwks.Keys {
		if jwk.KeyID != kid {
			continue
		}
		nBytes, err := base64.RawURLEncoding.DecodeString(jwk.N)
		if err != nil {
			continue
		}
		eBytes, err := base64.RawURLEncoding.DecodeString(jwk.E)
		if err != nil {
			continue
		}
		return &rsa.PublicKey{
			N: new(big.Int).SetBytes(nBytes),
			E: int(new(big.Int).SetBytes(eBytes).Int64()),
		}
	}
	return nil
}
