package clerk

import (
	"fmt"

	"github.com/way-platform/ileap-go/ileapauthserver"
)

// OIDCProvider implements ileapauthserver.OIDCProvider using Clerk's JWKS.
type OIDCProvider struct {
	client *Client
}

// NewOIDCProvider creates a new OIDC provider backed by the Clerk client.
func NewOIDCProvider(client *Client) *OIDCProvider {
	return &OIDCProvider{client: client}
}

// OpenIDConfiguration returns the OIDC configuration for the given base URL.
func (p *OIDCProvider) OpenIDConfiguration(baseURL string) *ileapauthserver.OpenIDConfiguration {
	return &ileapauthserver.OpenIDConfiguration{
		IssuerURL:              baseURL,
		AuthURL:                baseURL + "/auth/token",
		TokenURL:               baseURL + "/auth/token",
		JWKSURL:                fmt.Sprintf("https://%s/.well-known/jwks.json", p.client.fapiDomain),
		Algorithms:             []string{"RS256"},
		ResponseTypesSupported: []string{"token"},
		SubjectTypesSupported:  []string{"public"},
	}
}

// JWKS fetches and returns Clerk's JSON Web Key Set.
func (p *OIDCProvider) JWKS() *ileapauthserver.JWKSet {
	jwks, err := p.client.FetchJWKS()
	if err != nil {
		return &ileapauthserver.JWKSet{}
	}
	return jwks
}
