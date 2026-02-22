package ileapclerk

import (
	"context"
	"fmt"

	"github.com/way-platform/ileap-go"
	"golang.org/x/oauth2"
)

// TokenIssuer implements ileap.TokenIssuer using Clerk FAPI.
type TokenIssuer struct {
	client      *Client
	activeOrgID string
}

// TokenIssuerOption configures the Clerk token issuer.
type TokenIssuerOption func(*TokenIssuer)

// WithActiveOrganization sets the active organization ID for issued tokens.
func WithActiveOrganization(orgID string) TokenIssuerOption {
	return func(t *TokenIssuer) { t.activeOrgID = orgID }
}

// NewTokenIssuer creates a new Clerk-backed token issuer.
func NewTokenIssuer(client *Client, opts ...TokenIssuerOption) *TokenIssuer {
	t := &TokenIssuer{client: client}
	for _, opt := range opts {
		opt(t)
	}
	return t
}

// IssueToken validates credentials via Clerk and returns Clerk's session JWT.
func (t *TokenIssuer) IssueToken(
	_ context.Context, clientID, clientSecret string,
) (*oauth2.Token, error) {
	jwt, err := t.client.SignIn(clientID, clientSecret, t.activeOrgID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ileap.ErrInvalidCredentials, err)
	}
	return &oauth2.Token{
		AccessToken: jwt,
		TokenType:   "bearer",
	}, nil
}
