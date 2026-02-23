package ileap

import (
	"context"
	"net/http"

	"connectrpc.com/connect"
)

type interceptorTransport struct {
	interceptors []func(http.RoundTripper) http.RoundTripper
	next         http.RoundTripper
}

var _ http.RoundTripper = &interceptorTransport{}

func (t *interceptorTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	rt := t.next
	for _, interceptor := range t.interceptors {
		rt = interceptor(rt)
	}
	return rt.RoundTrip(req)
}

// AuthForwardInterceptor returns a Connect client interceptor that reads the
// bearer token from the request context (set by Server's auth middleware) and
// attaches it to outgoing Connect requests.
func AuthForwardInterceptor() connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			if token, ok := AuthTokenFromContext(ctx); ok {
				req.Header().Set("Authorization", "Bearer "+token)
			}
			return next(ctx, req)
		}
	}
}
