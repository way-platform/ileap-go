package clerk

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/way-platform/ileap-go/demo"
	"github.com/way-platform/ileap-go/ileapauthserver"
)

func TestTokenIssuer(t *testing.T) {
	keypair, err := demo.LoadKeyPair()
	if err != nil {
		t.Fatalf("load keypair: %v", err)
	}

	t.Run("success", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(signInResponse{Status: "complete"})
		}))
		defer srv.Close()
		client := NewClient("unused", WithHTTPClient(&http.Client{
			Transport: &testTransport{target: srv},
		}))
		issuer := NewTokenIssuer(client, keypair)
		creds, err := issuer.IssueToken(context.Background(), "user@example.com", "password")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if creds.TokenType != "bearer" {
			t.Errorf("expected bearer, got %s", creds.TokenType)
		}
		if creds.AccessToken == "" {
			t.Error("expected non-empty access token")
		}
		// Validate the JWT.
		claims, err := keypair.ValidateJWT(creds.AccessToken)
		if err != nil {
			t.Fatalf("validate JWT: %v", err)
		}
		if claims.Username != "user@example.com" {
			t.Errorf("expected username user@example.com, got %s", claims.Username)
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
		issuer := NewTokenIssuer(client, keypair)
		_, err := issuer.IssueToken(context.Background(), "bad", "wrong")
		if err == nil {
			t.Fatal("expected error")
		}
		if !errors.Is(err, ileapauthserver.ErrInvalidCredentials) {
			t.Errorf("expected ErrInvalidCredentials, got: %v", err)
		}
	})
}
