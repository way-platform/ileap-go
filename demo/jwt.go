package demo

// JWTHeader is the header of a JWT.
type JWTHeader struct {
	// Type is the type of the JWT.
	Type string `json:"typ"`
	// Algorithm is the algorithm used to sign the JWT.
	Algorithm string `json:"alg"`
}

// JWTClaims are the claims of a JWT.
type JWTClaims struct {
	// Username is the username of the user.
	Username string `json:"username"`
	// IssuedAt is the Unix timestamp of the JWT issuance.
	IssuedAt int64 `json:"iat,omitempty"`
}
