package ileap

import (
	"context"
	"fmt"

	ileapv1 "github.com/way-platform/ileap-go/proto/gen/wayplatform/connect/ileap/v1"
	"golang.org/x/oauth2"
)

type unimplementedFootprintHandler struct{}

func (unimplementedFootprintHandler) GetFootprint(
	_ context.Context, _ string,
) (*ileapv1.ProductFootprint, error) {
	return nil, fmt.Errorf("footprints: %w", ErrNotImplemented)
}

func (unimplementedFootprintHandler) ListFootprints(
	_ context.Context, _ ListFootprintsRequest,
) (*ListFootprintsResponse, error) {
	return nil, fmt.Errorf("footprints: %w", ErrNotImplemented)
}

type unimplementedTADHandler struct{}

func (unimplementedTADHandler) ListTADs(
	_ context.Context, _ ListTADsRequest,
) (*ListTADsResponse, error) {
	return nil, fmt.Errorf("tad: %w", ErrNotImplemented)
}

type unimplementedEventHandler struct{}

func (unimplementedEventHandler) HandleEvent(_ context.Context, _ Event) error {
	return fmt.Errorf("events: %w", ErrNotImplemented)
}

type unimplementedTokenValidator struct{}

func (unimplementedTokenValidator) ValidateToken(_ context.Context, _ string) (*TokenInfo, error) {
	return nil, fmt.Errorf("token validation: %w", ErrNotImplemented)
}

type unimplementedTokenIssuer struct{}

func (unimplementedTokenIssuer) IssueToken(
	_ context.Context, _, _ string,
) (*oauth2.Token, error) {
	return nil, fmt.Errorf("token issuance: %w", ErrNotImplemented)
}

type unimplementedOIDCProvider struct{}

func (unimplementedOIDCProvider) OpenIDConfiguration(_ string) *OpenIDConfiguration {
	return nil
}

func (unimplementedOIDCProvider) JWKS() *JWKSet {
	return nil
}

func (s *Server) setDefaults() {
	if s.footprintHandler == nil {
		s.footprintHandler = unimplementedFootprintHandler{}
	}
	if s.tadHandler == nil {
		s.tadHandler = unimplementedTADHandler{}
	}
	if s.eventHandler == nil {
		s.eventHandler = unimplementedEventHandler{}
	}
	if s.tokenValidator == nil {
		s.tokenValidator = unimplementedTokenValidator{}
	}
	if s.issuer == nil {
		s.issuer = unimplementedTokenIssuer{}
	}
	if s.oidc == nil {
		s.oidc = unimplementedOIDCProvider{}
	}
}
