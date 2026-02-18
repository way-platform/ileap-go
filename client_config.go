package ileap

import "net/http"

// DemoBaseURL is the base URL for the SINE Foundation's demo iLEAP API.
const DemoBaseURL = "https://api.ileap.sine.dev"

// ClientConfig is the configuration for a [Client].
type ClientConfig struct {
	baseURL    string
	transport  http.RoundTripper
	retryCount int
	logger     Logger
}

// newClientConfig creates a new default [ClientConfig].
func newClientConfig() ClientConfig {
	return ClientConfig{
		transport: http.DefaultTransport,
	}
}

// ClientOption is an option that configures a [Client].
type ClientOption func(*ClientConfig)

// WithBaseURL sets the API base URL for the [Client].
func WithBaseURL(baseURL string) ClientOption {
	return func(cc *ClientConfig) {
		cc.baseURL = baseURL
	}
}

// WithOAuth2 authenticates requests using OAuth 2.0.
func WithOAuth2(clientID, clientSecret string) ClientOption {
	return func(cc *ClientConfig) {
		cc.transport = &tokenAuthenticatorTransport{
			tokenAuthenticator: &oauth2TokenAuthenticator{
				baseURL:      cc.baseURL,
				clientID:     clientID,
				clientSecret: clientSecret,
				httpClient:   http.DefaultClient,
			},
			transport: cc.transport,
		}
	}
}

// WithReuseTokenAuth authenticates requests by re-using existing [ClientCredentials].
func WithReuseTokenAuth(credentials ClientCredentials) ClientOption {
	return func(cc *ClientConfig) {
		cc.transport = &reuseTokenCredentialsTransport{
			transport:   cc.transport,
			credentials: credentials,
		}
	}
}

// WithRetryCount sets the maximum number of times to retry a request.
func WithRetryCount(retryCount int) ClientOption {
	return func(cc *ClientConfig) {
		cc.retryCount = retryCount
	}
}

// Logger is a leveled logger interface.
type Logger interface {
	Debug(msg string, keysAndValues ...any)
	Info(msg string, keysAndValues ...any)
	Warn(msg string, keysAndValues ...any)
	Error(msg string, keysAndValues ...any)
}

// WithLogger sets the [Logger] for the [Client].
func WithLogger(logger Logger) ClientOption {
	return func(cc *ClientConfig) {
		cc.logger = logger
	}
}
