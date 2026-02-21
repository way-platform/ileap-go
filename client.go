package ileap

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"runtime/debug"
)

// Client is an iLEAP API client.
type Client struct {
	config     ClientConfig
	httpClient *http.Client
}

// NewClient creates a new [Client] with the given base URL and options.
func NewClient(opts ...ClientOption) *Client {
	config := newClientConfig()
	for _, opt := range opts {
		opt(&config)
	}
	transport := http.RoundTripper(http.DefaultTransport)
	if config.debug {
		transport = &debugTransport{next: transport}
	}
	if config.auth != nil {
		transport = config.auth(transport)
	}
	if len(config.interceptors) > 0 {
		transport = &interceptorTransport{interceptors: config.interceptors, next: transport}
	}
	if config.retryCount > 0 {
		transport = &retryTransport{maxRetries: config.retryCount, next: transport}
	}
	return &Client{
		config:     config,
		httpClient: &http.Client{Transport: transport},
	}
}

func (c *Client) newRequest(
	ctx context.Context,
	method, requestPath string,
	body io.Reader,
) (_ *http.Request, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("new request: %w", err)
		}
	}()
	requestURL, err := url.JoinPath(c.config.baseURL, requestPath)
	if err != nil {
		return nil, fmt.Errorf("invalid request URL: %w", err)
	}
	request, err := http.NewRequestWithContext(ctx, method, requestURL, body)
	if err != nil {
		return nil, err
	}
	request.Header.Set("User-Agent", getUserAgent())
	return request, nil
}

func getUserAgent() string {
	userAgent := "wayplatform.com"
	if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" {
		userAgent += "/" + info.Main.Version
	}
	return userAgent
}
