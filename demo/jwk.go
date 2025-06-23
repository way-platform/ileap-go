package demo

// JWKSet is a JSON Web Key Set.
type JWKSet struct {
	// Keys is the list of JSON Web Keys.
	Keys []JWK `json:"keys"`
}

// JWK is a JSON Web Key.
type JWK struct {
	// KeyType specifies the cryptographic algorithm family used with the key
	KeyType string `json:"kty"`
	// Use identifies the intended use of the public key
	Use string `json:"use,omitempty"`
	// Algorithm identifies the algorithm intended for use with the key
	Algorithm string `json:"alg,omitempty"`
	// KeyID is a hint indicating which key was used to secure the JWS
	KeyID string `json:"kid,omitempty"`
	// N is the modulus for the RSA public key (base64url encoded).
	N string `json:"n"`
	// E is the exponent for the RSA public key (base64url encoded).
	E string `json:"e"`
}
