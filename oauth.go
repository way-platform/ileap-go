package ileap

import "fmt"

// OAuthError is an OAuth 2.0 error response body returned by a PACT-compliant server.
//
// See: https://datatracker.ietf.org/doc/html/rfc6749#section-5.2
type OAuthError struct {
	// Code is the OAuth 2.0 error code identifier.
	Code OAuthErrorCode `json:"error"`
	// Description is a human readable description of the error.
	Description string `json:"error_description,omitempty"`
}

// Error implements the error interface.
func (e *OAuthError) Error() string {
	return fmt.Sprintf("OAuth error %s: %s", e.Code, e.Description)
}

// OAuthErrorCode is an OAuth 2.0 error code identifier.
type OAuthErrorCode string

const (
	// OAuthErrorCodeInvalidRequest means the request is missing a required parameter,
	// includes an invalid parameter value, includes a parameter more than once,
	// or is otherwise malformed.
	OAuthErrorCodeInvalidRequest OAuthErrorCode = "invalid_request"

	// OAuthErrorCodeUnauthorizedClient means the client is not authorized to request an access
	// token using this method.
	OAuthErrorCodeUnauthorizedClient OAuthErrorCode = "unauthorized_client"

	// OAuthErrorCodeAccessDenied means the resource owner or authorization server denied the request.
	OAuthErrorCodeAccessDenied OAuthErrorCode = "access_denied"

	// OAuthErrorCodeUnsupportedResponseType means the authorization server does not support
	// obtaining an access token using this method.
	OAuthErrorCodeUnsupportedResponseType OAuthErrorCode = "unsupported_response_type"

	// OAuthErrorCodeInvalidScope means the requested scope is invalid, unknown, or malformed.
	OAuthErrorCodeInvalidScope OAuthErrorCode = "invalid_scope"

	// OAuthErrorCodeServerError means the authorization server encountered an unexpected condition
	// that prevented it from fulfilling the request.
	OAuthErrorCodeServerError OAuthErrorCode = "server_error"

	// OAuthErrorCodeTemporarilyUnavailable means the authorization server is currently unable
	// to handle the request due to a temporary overloading or maintenance of the server.
	OAuthErrorCodeTemporarilyUnavailable OAuthErrorCode = "temporarily_unavailable"

	// OAuthErrorCodeUnsupportedGrantType means the authorization server does not support
	// the authorization grant type used.
	OAuthErrorCodeUnsupportedGrantType OAuthErrorCode = "unsupported_grant_type"
)
