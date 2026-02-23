// Package ileapconnect implements ileap.FootprintHandler and ileap.TADHandler
// by forwarding requests to a Connect RPC backend implementing ILeapService.
//
// This enables the ileap.Server to act as a protocol translator: it handles
// all iLEAP/PACT HTTP protocol specifics (JSON envelope, Link header
// pagination, OData filtering, OAuth2 error formats) while delegating data
// retrieval to an internal Connect service.
package ileapconnect

import (
	"context"
	"errors"
	"net/http"

	"connectrpc.com/connect"
	ileap "github.com/way-platform/ileap-go"
	"github.com/way-platform/ileap-go/proto/gen/wayplatform/connect/ileap/v1/ileapv1connect"
)

var (
	_ ileap.FootprintHandler = (*Handler)(nil)
	_ ileap.TADHandler       = (*Handler)(nil)
)

// Handler implements ileap.FootprintHandler and ileap.TADHandler by forwarding
// to a Connect RPC backend.
type Handler struct {
	client ileapv1connect.ILeapServiceClient
}

// Option configures the Handler.
type Option func(*options)

type options struct {
	httpClient connect.HTTPClient
	clientOpts []connect.ClientOption
}

// WithHTTPClient sets the HTTP client used for Connect RPC calls.
// Defaults to http.DefaultClient.
func WithHTTPClient(c connect.HTTPClient) Option {
	return func(o *options) { o.httpClient = c }
}

// WithClientOptions appends Connect client options (e.g. interceptors).
func WithClientOptions(opts ...connect.ClientOption) Option {
	return func(o *options) { o.clientOpts = append(o.clientOpts, opts...) }
}

// New creates a Handler that forwards to the Connect backend at backendURL.
// The incoming Authorization header is automatically forwarded on all
// outgoing Connect calls.
func New(backendURL string, opts ...Option) *Handler {
	o := options{
		httpClient: http.DefaultClient,
	}
	for _, opt := range opts {
		opt(&o)
	}
	clientOpts := append([]connect.ClientOption{
		connect.WithInterceptors(authForwardInterceptor()),
	}, o.clientOpts...)
	return &Handler{
		client: ileapv1connect.NewILeapServiceClient(o.httpClient, backendURL, clientOpts...),
	}
}

// authForwardInterceptor reads the bearer token from the request context
// (set by ileap.Server's auth middleware) and attaches it to outgoing
// Connect requests.
func authForwardInterceptor() connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			if token, ok := ileap.AuthTokenFromContext(ctx); ok {
				req.Header().Set("Authorization", "Bearer "+token)
			}
			return next(ctx, req)
		}
	}
}

// mapConnectError translates Connect error codes to ileap sentinel errors.
func mapConnectError(err error) error {
	var connectErr *connect.Error
	if !errors.As(err, &connectErr) {
		return err
	}
	switch connectErr.Code() {
	case connect.CodeNotFound:
		return ileap.ErrNotFound
	case connect.CodeInvalidArgument:
		return ileap.ErrBadRequest
	case connect.CodeUnauthenticated:
		return ileap.ErrTokenExpired
	case connect.CodePermissionDenied:
		return ileap.ErrInvalidCredentials
	case connect.CodeUnimplemented:
		return ileap.ErrNotImplemented
	default:
		return err
	}
}
