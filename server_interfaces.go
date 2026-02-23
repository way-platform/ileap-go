package ileap

import (
	"context"

	"golang.org/x/oauth2"
)

// AuthHandler handles all authentication and OIDC discovery operations.
type AuthHandler interface {
	// IssueToken validates credentials and returns an access token.
	IssueToken(ctx context.Context, clientID, clientSecret string) (*oauth2.Token, error)
	// ValidateToken validates an access token and returns token info.
	ValidateToken(ctx context.Context, token string) (*TokenInfo, error)
	// OpenIDConfiguration returns the OIDC configuration for the given base URL.
	OpenIDConfiguration(baseURL string) *OpenIDConfiguration
	// JWKS returns the JSON Web Key Set.
	JWKS() *JWKSet
}
