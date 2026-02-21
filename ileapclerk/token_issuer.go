package ileapclerk

import (
	"context"
	"fmt"

	"github.com/way-platform/ileap-go/ileapauthserver"
	"golang.org/x/oauth2"
)

// TokenIssuer implements ileapauthserver.TokenIssuer using Clerk FAPI.
type TokenIssuer struct {
	client *Client
}

// NewTokenIssuer creates a new Clerk-backed token issuer.
func NewTokenIssuer(client *Client) *TokenIssuer {
	return &TokenIssuer{client: client}
}

// IssueToken validates credentials via Clerk and returns Clerk's session JWT.
func (t *TokenIssuer) IssueToken(
	_ context.Context, clientID, clientSecret string,
) (*oauth2.Token, error) {
	jwt, err := t.client.SignIn(clientID, clientSecret)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ileapauthserver.ErrInvalidCredentials, err)
	}
	return &oauth2.Token{
		AccessToken: jwt,
		TokenType:   "bearer",
	}, nil
}
