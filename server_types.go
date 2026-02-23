package ileap

// TokenInfo contains information extracted from a validated token.
type TokenInfo struct {
	// Subject is the subject (user) of the token.
	Subject string
}

// OpenIDConfiguration is an OpenID Connect discovery document.
type OpenIDConfiguration struct {
	IssuerURL              string   `json:"issuer"`
	AuthURL                string   `json:"authorization_endpoint"`
	TokenURL               string   `json:"token_endpoint"`
	DeviceAuthURL          string   `json:"device_authorization_endpoint,omitempty"`
	UserInfoURL            string   `json:"userinfo_endpoint,omitempty"`
	JWKSURL                string   `json:"jwks_uri"`
	Algorithms             []string `json:"id_token_signing_alg_values_supported"`
	ResponseTypesSupported []string `json:"response_types_supported"`
	SubjectTypesSupported  []string `json:"subject_types_supported"`
}

// JWKSet is a JSON Web Key Set.
type JWKSet struct {
	Keys []JWK `json:"keys"`
}

// JWK is a JSON Web Key.
type JWK struct {
	KeyType   string `json:"kty"`
	Use       string `json:"use,omitempty"`
	Algorithm string `json:"alg,omitempty"`
	KeyID     string `json:"kid,omitempty"`
	N         string `json:"n"`
	E         string `json:"e"`
}
