package ileap

import "errors"

// Sentinel errors returned by handler implementations.
var (
	// ErrNotFound indicates the requested resource was not found.
	ErrNotFound = errors.New("not found")
	// ErrBadRequest indicates a malformed or invalid request.
	ErrBadRequest = errors.New("bad request")
	// ErrTokenExpired indicates the access token has expired.
	ErrTokenExpired = errors.New("token expired")
	// ErrInvalidCredentials indicates the provided credentials are invalid.
	ErrInvalidCredentials = errors.New("invalid credentials")
	// ErrNotImplemented indicates the handler is not configured.
	ErrNotImplemented = errors.New("not implemented")
)
