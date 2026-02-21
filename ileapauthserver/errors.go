package ileapauthserver

import "errors"

// Sentinel errors returned by TokenIssuer implementations.
var (
	// ErrInvalidCredentials indicates the provided credentials are invalid.
	ErrInvalidCredentials = errors.New("invalid credentials")
)
