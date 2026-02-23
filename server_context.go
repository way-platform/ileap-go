package ileap

import "context"

type contextKey int

const authTokenKey contextKey = iota

// WithAuthToken returns a new context with the given bearer token stored.
// This is used by the auth middleware to propagate the validated token to
// downstream handlers, enabling them to forward it to backend services.
func WithAuthToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, authTokenKey, token)
}

// AuthTokenFromContext retrieves the bearer token from the context.
// Returns the token and true if present, or empty string and false otherwise.
func AuthTokenFromContext(ctx context.Context) (string, bool) {
	token, ok := ctx.Value(authTokenKey).(string)
	return token, ok
}
