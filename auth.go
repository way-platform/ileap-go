package ileap

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// TokenCredentials for an authenticated iLEAP API client.
type TokenCredentials struct {
	// Token is the bearer token for the authenticated client.
	AccessToken string `json:"access_token"`
	// TokenType is the type of token.
	TokenType string `json:"token_type"`
	// ExpireTime is the time when the token expires.
	ExpireTime time.Time `json:"expires_in,omitzero"`
}

// TokenAuthenticator is a pluggable interface for authenticating requests to an iLEAP API.
type TokenAuthenticator interface {
	// Authenticate the client and return a set of [TokenCredentials].
	Authenticate(ctx context.Context) (TokenCredentials, error)
}

type tokenAuthenticatorTransport struct {
	tokenAuthenticator TokenAuthenticator
	transport          http.RoundTripper
	mu                 sync.Mutex
	credentials        TokenCredentials
}

func (t *tokenAuthenticatorTransport) RoundTrip(req *http.Request) (_ *http.Response, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("token authenticator transport: %w", err)
		}
	}()
	token, err := t.getToken(req.Context())
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	return t.transport.RoundTrip(req)
}

func (t *tokenAuthenticatorTransport) getToken(ctx context.Context) (string, error) {
	t.mu.Lock()
	credentials := t.credentials
	t.mu.Unlock()
	if credentials.ExpireTime.IsZero() || credentials.ExpireTime.After(time.Now()) {
		return credentials.AccessToken, nil
	}
	newCredentials, err := t.tokenAuthenticator.Authenticate(ctx)
	if err != nil {
		return "", fmt.Errorf("authenticate: %w", err)
	}
	t.mu.Lock()
	t.credentials = newCredentials
	t.mu.Unlock()
	return newCredentials.AccessToken, nil
}

type reuseTokenCredentialsTransport struct {
	transport   http.RoundTripper
	credentials TokenCredentials
}

func (t *reuseTokenCredentialsTransport) RoundTrip(req *http.Request) (_ *http.Response, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("reuse token credentials transport: %w", err)
		}
	}()
	req.Header.Set("Authorization", "Bearer "+t.credentials.AccessToken)
	return t.transport.RoundTrip(req)
}

func NewOAuth2TokenAuthenticator(clientID, clientSecret, baseURL string) TokenAuthenticator {
	return &oauth2TokenAuthenticator{
		clientID:     clientID,
		clientSecret: clientSecret,
		baseURL:      baseURL,
		httpClient:   http.DefaultClient,
	}
}

type oauth2TokenAuthenticator struct {
	clientID     string
	clientSecret string
	baseURL      string
	httpClient   *http.Client
}

func (t *oauth2TokenAuthenticator) newRequest(ctx context.Context, method, path string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, t.baseURL+path, body)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}
	return req, nil
}

func (t *oauth2TokenAuthenticator) Authenticate(ctx context.Context) (_ TokenCredentials, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("authenticate: %w", err)
		}
	}()
	req, err := t.newRequest(ctx, http.MethodPost, "/auth/token", nil)
	if err != nil {
		return TokenCredentials{}, err
	}
	req.SetBasicAuth(t.clientID, t.clientSecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	body := url.Values{}
	body.Set("grant_type", "client_credentials")
	req.Body = io.NopCloser(bytes.NewBufferString(body.Encode()))
	res, err := t.httpClient.Do(req)
	if err != nil {
		return TokenCredentials{}, fmt.Errorf("send request: %w", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return TokenCredentials{}, fmt.Errorf("response status code: %d", res.StatusCode)
	}
	var response TokenCredentials
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return TokenCredentials{}, fmt.Errorf("decode response: %w", err)
	}
	return response, nil
}
