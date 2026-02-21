package clerk

import (
	"github.com/way-platform/ileap-go/demo"
	"github.com/way-platform/ileap-go/ileapauthserver"
)

// OIDCProvider implements ileapauthserver.OIDCProvider using a local keypair.
type OIDCProvider struct {
	keypair *demo.KeyPair
}

// NewOIDCProvider creates a new OIDC provider backed by the given keypair.
func NewOIDCProvider(keypair *demo.KeyPair) *OIDCProvider {
	return &OIDCProvider{keypair: keypair}
}

// OpenIDConfiguration returns the OIDC configuration for the given base URL.
func (p *OIDCProvider) OpenIDConfiguration(baseURL string) *ileapauthserver.OpenIDConfiguration {
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
func (p *OIDCProvider) JWKS() *ileapauthserver.JWKSet {
	jwk := p.keypair.JWK()
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
