// Package clerk provides a Clerk FAPI authentication backend for iLEAP servers.
package ileapclerk

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/way-platform/ileap-go/ileapserver"
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
		Status           string `json:"status"`
		CreatedSessionID string `json:"created_session_id"`
	} `json:"response"`
	Client struct {
		Sessions []struct {
			ID              string `json:"id"`
			LastActiveToken struct {
				JWT string `json:"jwt"`
			} `json:"last_active_token"`
		} `json:"sessions"`
	} `json:"client"`
}

// SignIn authenticates a user via Clerk's password strategy and returns the session JWT.
func (c *Client) SignIn(identifier, password, activeOrgID string) (string, error) {
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
		var body []byte
		if resp.Body != nil {
			body, _ = io.ReadAll(resp.Body)
		}
		return "", fmt.Errorf("clerk sign_in failed: HTTP %d: %s", resp.StatusCode, string(body))
	}
	var result signInResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}
	if result.Response.Status != "complete" {
		return "", fmt.Errorf("clerk sign_in not complete: %s", result.Response.Status)
	}

	sessionID := result.Response.CreatedSessionID
	if sessionID == "" && len(result.Client.Sessions) > 0 {
		sessionID = result.Client.Sessions[0].ID
	}

	authHeader := resp.Header.Get("Authorization")

	if activeOrgID != "" {
		if sessionID == "" {
			return "", fmt.Errorf("clerk sign_in: missing session ID for organization activation")
		}
		return c.TouchSession(sessionID, activeOrgID, authHeader)
	}

	sessions := result.Client.Sessions
	if len(sessions) == 0 || sessions[0].LastActiveToken.JWT == "" {
		return "", fmt.Errorf("clerk sign_in: no session JWT in response")
	}
	return sessions[0].LastActiveToken.JWT, nil
}

// TouchSession activates an organization for the given session and returns the new session JWT.
func (c *Client) TouchSession(sessionID, activeOrgID, authHeader string) (string, error) {
	endpoint := fmt.Sprintf("https://%s/v1/client/sessions/%s/touch", c.fapiDomain, sessionID)
	form := url.Values{}
	form.Set("active_organization_id", activeOrgID)
	req, err := http.NewRequest("POST", endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return "", fmt.Errorf("create touch request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("send touch request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("clerk touch_session failed: HTTP %d", resp.StatusCode)
	}
	var result signInResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("decode touch response: %w", err)
	}
	sessions := result.Client.Sessions
	if len(sessions) == 0 || sessions[0].LastActiveToken.JWT == "" {
		return "", fmt.Errorf("clerk touch_session: no session JWT in response")
	}
	return sessions[0].LastActiveToken.JWT, nil
}

// FetchJWKS fetches the JSON Web Key Set from Clerk's JWKS endpoint.
func (c *Client) FetchJWKS() (*ileapserver.JWKSet, error) {
	endpoint := fmt.Sprintf("https://%s/.well-known/jwks.json", c.fapiDomain)
	resp, err := c.httpClient.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("fetch JWKS: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("fetch JWKS: HTTP %d", resp.StatusCode)
	}
	var jwks ileapserver.JWKSet
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return nil, fmt.Errorf("decode JWKS: %w", err)
	}
	return &jwks, nil
}
