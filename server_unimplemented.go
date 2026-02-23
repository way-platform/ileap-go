package ileap

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	"github.com/way-platform/ileap-go/proto/gen/wayplatform/connect/ileap/v1/ileapv1connect"
	"golang.org/x/oauth2"
)

type unimplementedAuthHandler struct{}

func (unimplementedAuthHandler) IssueToken(
	_ context.Context, _, _ string,
) (*oauth2.Token, error) {
	return nil, connect.NewError(
		connect.CodeUnimplemented,
		errors.New("token issuance not implemented"),
	)
}

func (unimplementedAuthHandler) ValidateToken(_ context.Context, _ string) (*TokenInfo, error) {
	return nil, connect.NewError(
		connect.CodeUnimplemented,
		errors.New("token validation not implemented"),
	)
}

func (unimplementedAuthHandler) OpenIDConfiguration(_ string) *OpenIDConfiguration {
	return nil
}

func (unimplementedAuthHandler) JWKS() *JWKSet {
	return nil
}

func (s *Server) setDefaults() {
	if s.service == nil {
		s.service = ileapv1connect.UnimplementedILeapServiceHandler{}
	}
	if s.auth == nil {
		s.auth = unimplementedAuthHandler{}
	}
}
