// Package ileapserver provides a reusable HTTP server adapter for
// iLEAP-compliant data endpoints.
package ileapserver

import (
	"context"

	"github.com/way-platform/ileap-go/openapi/ileapv1"
)

// FootprintHandler handles product footprint requests.
type FootprintHandler interface {
	// GetFootprint returns a single footprint by ID.
	GetFootprint(ctx context.Context, id string) (*ileapv1.ProductFootprintForILeapType, error)
	// ListFootprints returns a filtered, limited list of footprints.
	ListFootprints(ctx context.Context, req ListFootprintsRequest) (*ListFootprintsResponse, error)
}

// TADHandler handles transport activity data requests.
type TADHandler interface {
	// ListTADs returns a limited list of transport activity data.
	ListTADs(ctx context.Context, req ListTADsRequest) (*ListTADsResponse, error)
}

// EventHandler handles PACT CloudEvents.
type EventHandler interface {
	// HandleEvent processes an incoming event.
	HandleEvent(ctx context.Context, event Event) error
}

// TokenValidator validates bearer tokens.
type TokenValidator interface {
	// ValidateToken validates an access token and returns token info.
	ValidateToken(ctx context.Context, token string) (*TokenInfo, error)
}
