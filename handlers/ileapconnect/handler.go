// Package ileapconnect creates a Connect RPC client that can be used directly
// as an ILeapServiceHandler with the ileap.Server.
//
// The generated ILeapServiceClient and ILeapServiceHandler interfaces have
// identical method signatures, so the client satisfies the handler interface
// directly. Auth forwarding is handled by the ileap.AuthForwardInterceptor.
package ileapconnect

import (
	"net/http"

	"connectrpc.com/connect"
	ileap "github.com/way-platform/ileap-go"
	"github.com/way-platform/ileap-go/proto/gen/wayplatform/connect/ileap/v1/ileapv1connect"
)

// Option configures the client created by NewClient.
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

// NewClient creates an ILeapServiceClient that forwards to the Connect backend
// at backendURL. The incoming Authorization header is automatically forwarded
// on all outgoing Connect calls via ileap.AuthForwardInterceptor.
//
// The returned client satisfies ileapv1connect.ILeapServiceHandler and can be
// passed directly to ileap.WithServiceHandler.
func NewClient(backendURL string, opts ...Option) ileapv1connect.ILeapServiceClient {
	o := options{
		httpClient: http.DefaultClient,
	}
	for _, opt := range opts {
		opt(&o)
	}
	clientOpts := append([]connect.ClientOption{
		connect.WithInterceptors(ileap.AuthForwardInterceptor()),
	}, o.clientOpts...)
	return ileapv1connect.NewILeapServiceClient(o.httpClient, backendURL, clientOpts...)
}
