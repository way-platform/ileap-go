package ileap

import "net/http"

// DemoBaseURL is the base URL for the SINE Foundation's demo iLEAP API.
const DemoBaseURL = "https://api.ileap.sine.dev"

// ClientConfig is the configuration for a [Client].
type ClientConfig struct {
	baseURL      string
	retryCount   int
	debug        bool
	interceptors []func(http.RoundTripper) http.RoundTripper
	auth         func(http.RoundTripper) http.RoundTripper
}

// newClientConfig creates a new default [ClientConfig].
func newClientConfig() ClientConfig {
	return ClientConfig{}
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
		cc.auth = func(next http.RoundTripper) http.RoundTripper {
			return &tokenAuthenticatorTransport{
				tokenAuthenticator: &oauth2TokenAuthenticator{
					baseURL:      cc.baseURL,
					clientID:     clientID,
					clientSecret: clientSecret,
					httpClient:   http.DefaultClient,
				},
				transport: next,
			}
		}
	}
}

// WithReuseTokenAuth authenticates requests by re-using existing [ClientCredentials].
func WithReuseTokenAuth(credentials ClientCredentials) ClientOption {
	return func(cc *ClientConfig) {
		cc.auth = func(next http.RoundTripper) http.RoundTripper {
			return &reuseTokenCredentialsTransport{
				transport:   next,
				credentials: credentials,
			}
		}
	}
}

// WithRetryCount sets the maximum number of times to retry a request.
func WithRetryCount(retryCount int) ClientOption {
	return func(cc *ClientConfig) {
		cc.retryCount = retryCount
	}
}

// WithDebug toggles debug mode (request/response dumps to stderr).
func WithDebug(debug bool) ClientOption {
	return func(cc *ClientConfig) {
		cc.debug = debug
	}
}

// WithInterceptor adds a request interceptor for the [Client].
func WithInterceptor(interceptor func(http.RoundTripper) http.RoundTripper) ClientOption {
	return func(cc *ClientConfig) {
		cc.interceptors = append(cc.interceptors, interceptor)
	}
}
