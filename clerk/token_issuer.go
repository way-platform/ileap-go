package clerk

import (
	"context"
	"fmt"
	"time"

	"github.com/way-platform/ileap-go"
	"github.com/way-platform/ileap-go/demo"
	"github.com/way-platform/ileap-go/ileapauthserver"
)

// TokenIssuer implements ileapauthserver.TokenIssuer using Clerk FAPI
// for credential validation and a local keypair for JWT issuance.
type TokenIssuer struct {
	client  *Client
	keypair *demo.KeyPair
}

// NewTokenIssuer creates a new Clerk-backed token issuer.
func NewTokenIssuer(client *Client, keypair *demo.KeyPair) *TokenIssuer {
	return &TokenIssuer{
		client:  client,
		keypair: keypair,
	}
}

// IssueToken validates credentials via Clerk and issues a local JWT.
func (t *TokenIssuer) IssueToken(
	_ context.Context, clientID, clientSecret string,
) (*ileap.ClientCredentials, error) {
	if err := t.client.SignIn(clientID, clientSecret); err != nil {
		return nil, fmt.Errorf("%w: %w", ileapauthserver.ErrInvalidCredentials, err)
	}
	accessToken, err := t.keypair.CreateJWT(demo.JWTClaims{
		Username: clientID,
		IssuedAt: time.Now().Unix(),
	})
	if err != nil {
		return nil, fmt.Errorf("create JWT: %w", err)
	}
	return &ileap.ClientCredentials{
		AccessToken: accessToken,
		TokenType:   "bearer",
	}, nil
}
