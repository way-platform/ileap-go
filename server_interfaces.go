package ileap

import (
	"context"

	ileapv1 "github.com/way-platform/ileap-go/proto/gen/wayplatform/connect/ileap/v1"
	"golang.org/x/oauth2"
)

// TokenIssuer issues access tokens for valid credentials.
type TokenIssuer interface {
	// IssueToken validates credentials and returns client credentials.
	IssueToken(ctx context.Context, clientID, clientSecret string) (*oauth2.Token, error)
}

// OIDCProvider provides OpenID Connect discovery information.
type OIDCProvider interface {
	// OpenIDConfiguration returns the OIDC configuration for the given base URL.
	OpenIDConfiguration(baseURL string) *OpenIDConfiguration
	// JWKS returns the JSON Web Key Set.
	JWKS() *JWKSet
}

// FootprintHandler handles product footprint requests.
type FootprintHandler interface {
	// GetFootprint returns a single footprint by ID.
	GetFootprint(ctx context.Context, id string) (*ileapv1.ProductFootprint, error)
	// ListFootprints returns a filtered, limited list of footprints.
	ListFootprints(ctx context.Context, req ListFootprintsRequest) (*ListFootprintsResponse, error)
}

// TADHandler handles transport activity data requests.
type TADHandler interface {
	// ListTADs returns a limited list of transport activity data.
	ListTADs(ctx context.Context, req ListTADsRequest) (*ListTADsResponse, error)
}

// EventHandler handles PACT CloudEvents.
type EventHandler interface {
	// HandleEvent processes an incoming event.
	HandleEvent(ctx context.Context, event Event) error
}

// TokenValidator validates bearer tokens.
type TokenValidator interface {
	// ValidateToken validates an access token and returns token info.
	ValidateToken(ctx context.Context, token string) (*TokenInfo, error)
}
