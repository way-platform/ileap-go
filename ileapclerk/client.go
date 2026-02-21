// Package clerk provides a Clerk FAPI authentication backend for iLEAP servers.
package ileapclerk

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/way-platform/ileap-go/ileapauthserver"
)

// Client is an HTTP client for the Clerk Frontend API.
type Client struct {
	fapiDomain string
	httpClient *http.Client
}

// NewClient creates a new Clerk FAPI client.
func NewClient(fapiDomain string, opts ...ClientOption) *Client {
	c := &Client{
		fapiDomain: fapiDomain,
		httpClient: http.DefaultClient,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// ClientOption configures the Clerk client.
type ClientOption func(*Client)

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(hc *http.Client) ClientOption {
	return func(c *Client) { c.httpClient = hc }
}

// signInResponse is the response from Clerk's sign_in endpoint.
type signInResponse struct {
	Response struct {
		Status string `json:"status"`
	} `json:"response"`
	Client struct {
		Sessions []struct {
			LastActiveToken struct {
				JWT string `json:"jwt"`
			} `json:"last_active_token"`
		} `json:"sessions"`
	} `json:"client"`
}

// SignIn authenticates a user via Clerk's password strategy and returns the session JWT.
func (c *Client) SignIn(identifier, password string) (string, error) {
	endpoint := fmt.Sprintf("https://%s/v1/client/sign_ins", c.fapiDomain)
	form := url.Values{}
	form.Set("strategy", "password")
	form.Set("identifier", identifier)
	form.Set("password", password)
	req, err := http.NewRequest("POST", endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("send request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("clerk sign_in failed: HTTP %d", resp.StatusCode)
	}
	var result signInResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}
	if result.Response.Status != "complete" {
		return "", fmt.Errorf("clerk sign_in not complete: %s", result.Response.Status)
	}
	sessions := result.Client.Sessions
	if len(sessions) == 0 || sessions[0].LastActiveToken.JWT == "" {
		return "", fmt.Errorf("clerk sign_in: no session JWT in response")
	}
	return sessions[0].LastActiveToken.JWT, nil
}

// FetchJWKS fetches the JSON Web Key Set from Clerk's JWKS endpoint.
func (c *Client) FetchJWKS() (*ileapauthserver.JWKSet, error) {
	endpoint := fmt.Sprintf("https://%s/.well-known/jwks.json", c.fapiDomain)
	resp, err := c.httpClient.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("fetch JWKS: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("fetch JWKS: HTTP %d", resp.StatusCode)
	}
	var jwks ileapauthserver.JWKSet
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return nil, fmt.Errorf("decode JWKS: %w", err)
	}
	return &jwks, nil
}
