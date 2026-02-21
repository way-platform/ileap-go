// Package ileapauthserver provides a reusable HTTP server adapter for
// iLEAP-compliant OAuth2 authentication endpoints.
package ileapauthserver

import (
	"context"

	"github.com/way-platform/ileap-go"
)

// TokenIssuer issues access tokens for valid credentials.
type TokenIssuer interface {
	// IssueToken validates credentials and returns client credentials.
	IssueToken(ctx context.Context, clientID, clientSecret string) (*ileap.ClientCredentials, error)
}

// OIDCProvider provides OpenID Connect discovery information.
type OIDCProvider interface {
	// OpenIDConfiguration returns the OIDC configuration for the given base URL.
	OpenIDConfiguration(baseURL string) *OpenIDConfiguration
	// JWKS returns the JSON Web Key Set.
	JWKS() *JWKSet
}
