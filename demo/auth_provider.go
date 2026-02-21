package demo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/way-platform/ileap-go"
	"github.com/way-platform/ileap-go/ileapauthserver"
	"github.com/way-platform/ileap-go/ileapserver"
)

// AuthProvider implements ileapauthserver.TokenIssuer, ileapauthserver.OIDCProvider,
// and ileapserver.TokenValidator using demo credentials and a local RSA keypair.
type AuthProvider struct {
	keypair *KeyPair
}

// NewAuthProvider creates a new AuthProvider with the embedded demo keypair.
func NewAuthProvider() (*AuthProvider, error) {
	kp, err := LoadKeyPair()
	if err != nil {
		return nil, err
	}
	return &AuthProvider{keypair: kp}, nil
}

// IssueToken validates demo credentials and returns a signed JWT.
func (a *AuthProvider) IssueToken(
	_ context.Context, clientID, clientSecret string,
) (*ileap.ClientCredentials, error) {
	var authorized bool
	for _, user := range Users() {
		if clientID == user.Username && clientSecret == user.Password {
			authorized = true
			break
		}
	}
	if !authorized {
		return nil, ileapauthserver.ErrInvalidCredentials
	}
	accessToken, err := a.keypair.CreateJWT(JWTClaims{
		Username: clientID,
		IssuedAt: time.Now().Unix(),
	})
	if err != nil {
		return nil, err
	}
	return &ileap.ClientCredentials{
		AccessToken: accessToken,
		TokenType:   "bearer",
	}, nil
}

// ValidateToken validates a JWT and returns token info.
func (a *AuthProvider) ValidateToken(_ context.Context, token string) (*ileapserver.TokenInfo, error) {
	claims, err := a.keypair.ValidateJWT(token)
	if err != nil {
		if errors.Is(err, ileapserver.ErrTokenExpired) {
			return nil, fmt.Errorf("validate token: %w", err)
		}
		return nil, err
	}
	return &ileapserver.TokenInfo{Subject: claims.Username}, nil
}

// OpenIDConfiguration returns the OIDC configuration for the given base URL.
func (a *AuthProvider) OpenIDConfiguration(baseURL string) *ileapauthserver.OpenIDConfiguration {
	return &ileapauthserver.OpenIDConfiguration{
		IssuerURL:              baseURL,
		AuthURL:                baseURL + "/auth/token",
		TokenURL:               baseURL + "/auth/token",
		JWKSURL:                baseURL + "/jwks",
		Algorithms:             []string{"RS256"},
		ResponseTypesSupported: []string{"token"},
		SubjectTypesSupported:  []string{"public"},
	}
}

// JWKS returns the JSON Web Key Set containing the public key.
func (a *AuthProvider) JWKS() *ileapauthserver.JWKSet {
	jwk := a.keypair.JWK()
	return &ileapauthserver.JWKSet{
		Keys: []ileapauthserver.JWK{{
			KeyType:   jwk.KeyType,
			Use:       jwk.Use,
			Algorithm: jwk.Algorithm,
			KeyID:     jwk.KeyID,
			N:         jwk.N,
			E:         jwk.E,
		}},
	}
}
