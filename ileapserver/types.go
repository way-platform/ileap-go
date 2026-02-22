package ileapserver

import (
	"encoding/json"
	"time"

	"github.com/way-platform/ileap-go/openapi/ileapv1"
)

// ListFootprintsRequest is the request for listing footprints.
type ListFootprintsRequest struct {
	// Limit is the maximum number of footprints to return. 0 means no limit.
	Limit int
	// Offset is the starting index for pagination.
	Offset int
	// Filter is the raw OData $filter query string.
	Filter string
}

// ListFootprintsResponse is the response for listing footprints.
type ListFootprintsResponse struct {
	// Data is the list of footprints.
	Data []ileapv1.ProductFootprintForILeapType
	// Total is the total number of footprints matching the filter (before pagination).
	Total int
}

// ListTADsRequest is the request for listing transport activity data.
type ListTADsRequest struct {
	// Limit is the maximum number of TADs to return. 0 means no limit.
	Limit int
	// Offset is the starting index for pagination.
	Offset int
	// Filter contains query parameter filters (key → values, case-insensitive match).
	Filter map[string][]string
}

// ListTADsResponse is the response for listing transport activity data.
type ListTADsResponse struct {
	// Data is the list of TADs.
	Data []ileapv1.TAD
	// Total is the total number of TADs matching the filter (before pagination).
	Total int
}

// Event is a PACT CloudEvent received by the server.
type Event struct {
	// Type is the type of the event.
	Type string `json:"type"`
	// Specversion is the version of the CloudEvents specification.
	Specversion string `json:"specversion"`
	// ID is a unique identifier for the event.
	ID string `json:"id"`
	// Source is the source of the event.
	Source string `json:"source"`
	// Time is the time the event occurred.
	Time time.Time `json:"time"`
	// Data is the event data as raw JSON.
	Data json.RawMessage `json:"data"`
}

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
