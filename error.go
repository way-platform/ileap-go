package ileap

import (
	"fmt"
	"net/http"
)

// Error is an error returned by the rFMS API.
type Error struct {
	// Method is the HTTP method used to make the request.
	Method string
	// URL is the URL of the request.
	URL string
	// Status is the HTTP status.
	Status string
	// StatusCode is the HTTP status code.
	StatusCode int
}

func newHTTPError(resp *http.Response) *Error {
	return &Error{
		Method:     resp.Request.Method,
		URL:        resp.Request.URL.String(),
		Status:     resp.Status,
		StatusCode: resp.StatusCode,
	}
}

// Error implements the error interface.
func (e *Error) Error() string {
	return fmt.Sprintf("%s %s: %s", e.Method, e.URL, e.Status)
}
