package demo

// OpenIDConfiguration is an OpenID Connect configuration.
type OpenIDConfiguration struct {
	// IssuerURL is the identity of the provider, and the string it uses to sign
	// ID tokens with. For example "https://accounts.google.com". This value MUST
	// match ID tokens exactly.
	IssuerURL string `json:"issuer"`

	// AuthURL is the endpoint used by the provider to support the OAuth 2.0
	// authorization endpoint.
	AuthURL string `json:"authorization_endpoint"`

	// TokenURL is the endpoint used by the provider to support the OAuth 2.0
	// token endpoint.
	TokenURL string `json:"token_endpoint"`

	// DeviceAuthURL is the endpoint used by the provider to support the OAuth 2.0
	// device authorization endpoint.
	DeviceAuthURL string `json:"device_authorization_endpoint"`

	// UserInfoURL is the endpoint used by the provider to support the OpenID
	// Connect UserInfo flow.
	//
	// https://openid.net/specs/openid-connect-core-1_0.html#UserInfo
	UserInfoURL string `json:"userinfo_endpoint"`

	// JWKSURL is the endpoint used by the provider to advertise public keys to
	// verify issued ID tokens. This endpoint is polled as new keys are made
	// available.
	JWKSURL string `json:"jwks_uri"`

	// Algorithms, if provided, indicate a list of JWT algorithms allowed to sign
	// ID tokens. If not provided, this defaults to the algorithms advertised by
	// the JWK endpoint, then the set of algorithms supported by this package.
	Algorithms []string `json:"id_token_signing_alg_values_supported"`
}
