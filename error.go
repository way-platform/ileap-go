package ileap

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"mime"
	"net/http"
)

// Error is an error response body returned by a PACT-compliant server.
//
// See: https://docs.carbon-transparency.org/tr/data-exchange-protocol/latest/#error
type Error struct {
	// Code is the error code identifier. Required.
	Code ErrorCode `json:"code"`
	// Message is a human readable error message. Required.
	Message string `json:"message"`
}

// ErrorCode is an error code identifier.
type ErrorCode string

// Known PACT error codes.
const (
	ErrorCodeBadRequest      ErrorCode = "BadRequest"
	ErrorCodeAccessDenied    ErrorCode = "AccessDenied"
	ErrorCodeTokenExpired    ErrorCode = "TokenExpired"
	ErrorCodeNotFound        ErrorCode = "NotFound"
	ErrorCodeInternalError   ErrorCode = "InternalError"
	ErrorCodeNotImplemented  ErrorCode = "NotImplemented"
	ErrorCodeNoSuchFootprint ErrorCode = "NoSuchFootprint"
)

// Error implements the error interface.
func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// ClientError is an HTTP response error received by an iLEAP client.
type ClientError struct {
	// Method is the HTTP method used to make the request.
	Method string `json:"method"`
	// URL is the URL of the request.
	URL string `json:"url"`
	// Status is the HTTP status.
	Status string `json:"status"`
	// StatusCode is the HTTP status code.
	StatusCode int `json:"statusCode"`
	// Body is the PACT error body.
	Body Error `json:"error"`
}

// Error implements the error interface.
func (e *ClientError) Error() string {
	return fmt.Sprintf("%s %s: HTTP %s: (%s)", e.Method, e.URL, e.Status, e.Body)
}

func newClientError(resp *http.Response) *ClientError {
	var errorBody Error
	if body, err := io.ReadAll(resp.Body); err == nil {
		if mediaType, _, err := mime.ParseMediaType(resp.Header.Get("Content-Type")); err == nil {
			if mediaType == "application/json" {
				if err := json.Unmarshal(body, &errorBody); err != nil {
					slog.Debug("failed to unmarshal error body", "error", err)
				}
			}
		}
	}
	return &ClientError{
		Method:     resp.Request.Method,
		URL:        resp.Request.URL.String(),
		Status:     resp.Status,
		StatusCode: resp.StatusCode,
		Body:       errorBody,
	}
}
