package clerk

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/way-platform/ileap-go/ileapauthserver"
)

func TestTokenIssuer(t *testing.T) {
	const wantJWT = "header.payload.signature"

	t.Run("success", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			resp := signInResponse{}
			resp.Response.Status = "complete"
			resp.Client.Sessions = []struct {
				LastActiveToken struct {
					JWT string `json:"jwt"`
				} `json:"last_active_token"`
			}{{}}
			resp.Client.Sessions[0].LastActiveToken.JWT = wantJWT
			_ = json.NewEncoder(w).Encode(resp)
		}))
		defer srv.Close()
		client := NewClient("unused", WithHTTPClient(&http.Client{
			Transport: &testTransport{target: srv},
		}))
		issuer := NewTokenIssuer(client)
		creds, err := issuer.IssueToken(context.Background(), "user@example.com", "password")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if creds.TokenType != "bearer" {
			t.Errorf("expected bearer, got %s", creds.TokenType)
		}
		if creds.AccessToken != wantJWT {
			t.Errorf("expected access token %q, got %q", wantJWT, creds.AccessToken)
		}
	})

	t.Run("invalid credentials", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusUnprocessableEntity)
			_, _ = w.Write([]byte(`{"errors":[{"message":"Invalid credentials"}]}`))
		}))
		defer srv.Close()
		client := NewClient("unused", WithHTTPClient(&http.Client{
			Transport: &testTransport{target: srv},
		}))
		issuer := NewTokenIssuer(client)
		_, err := issuer.IssueToken(context.Background(), "bad", "wrong")
		if err == nil {
			t.Fatal("expected error")
		}
		if !errors.Is(err, ileapauthserver.ErrInvalidCredentials) {
			t.Errorf("expected ErrInvalidCredentials, got: %v", err)
		}
	})
}
