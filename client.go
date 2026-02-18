package ileap

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"runtime/debug"

	"github.com/hashicorp/go-retryablehttp"
)

// Client is an iLEAP API client.
type Client struct {
	config     ClientConfig
	httpClient *retryablehttp.Client
}

// NewClient creates a new [Client] with the given base URL and options.
func NewClient(opts ...ClientOption) *Client {
	config := newClientConfig()
	for _, opt := range opts {
		opt(&config)
	}
	httpClient := retryablehttp.NewClient()
	if config.transport != nil {
		httpClient.HTTPClient.Transport = config.transport
	}
	httpClient.RetryMax = config.retryCount
	if config.logger != nil {
		httpClient.Logger = config.logger
	}
	return &Client{
		config:     config,
		httpClient: httpClient,
	}
}

func (c *Client) newRequest(
	ctx context.Context,
	method, requestPath string,
	body io.Reader,
) (_ *retryablehttp.Request, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("new request: %w", err)
		}
	}()
	requestURL, err := url.JoinPath(c.config.baseURL, requestPath)
	if err != nil {
		return nil, fmt.Errorf("invalid request URL: %w", err)
	}
	request, err := retryablehttp.NewRequestWithContext(ctx, method, requestURL, body)
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
