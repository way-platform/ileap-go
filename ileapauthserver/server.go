package ileapauthserver

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/way-platform/ileap-go"
)

// Server is an iLEAP auth server HTTP handler.
type Server struct {
	issuer   TokenIssuer
	oidc     OIDCProvider
	serveMux *http.ServeMux
}

// NewServer creates a new iLEAP auth server.
func NewServer(issuer TokenIssuer, oidc OIDCProvider) *Server {
	s := &Server{
		issuer:   issuer,
		oidc:     oidc,
		serveMux: http.NewServeMux(),
	}
	s.serveMux.HandleFunc("POST /auth/token", s.authToken)
	s.serveMux.HandleFunc("GET /.well-known/openid-configuration", s.openIDConfig)
	s.serveMux.HandleFunc("GET /jwks", s.jwks)
	return s
}

// ServeHTTP implements http.Handler.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.serveMux.ServeHTTP(w, r)
}

func (s *Server) authToken(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
		writeOAuthError(
			w,
			http.StatusBadRequest,
			ileap.OAuthErrorCodeInvalidRequest,
			"invalid content type",
		)
		return
	}
	if err := r.ParseForm(); err != nil {
		writeOAuthError(
			w,
			http.StatusBadRequest,
			ileap.OAuthErrorCodeInvalidRequest,
			"invalid request body",
		)
		return
	}
	if grantType := r.Form.Get("grant_type"); grantType != "client_credentials" {
		writeOAuthError(
			w,
			http.StatusBadRequest,
			ileap.OAuthErrorCodeUnsupportedGrantType,
			"unsupported grant type",
		)
		return
	}
	username, password, ok := r.BasicAuth()
	if !ok {
		writeOAuthError(
			w,
			http.StatusBadRequest,
			ileap.OAuthErrorCodeInvalidRequest,
			"missing HTTP basic authorization",
		)
		return
	}
	creds, err := s.issuer.IssueToken(r.Context(), username, password)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			slog.WarnContext(r.Context(), "failed to issue token", "error", err)
			// ACT conformance test requires 400.
			writeOAuthError(
				w,
				http.StatusBadRequest,
				ileap.OAuthErrorCodeInvalidRequest,
				"invalid HTTP basic auth",
			)
			return
		}
		slog.ErrorContext(r.Context(), "failed to issue token", "error", err)
		writeOAuthError(
			w,
			http.StatusInternalServerError,
			ileap.OAuthErrorCodeServerError,
			"failed to issue token",
		)
		return
	}
	writeJSON(w, creds)
}

func (s *Server) openIDConfig(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, s.oidc.OpenIDConfiguration(baseURLFromRequest(r)))
}

func baseURLFromRequest(r *http.Request) string {
	scheme := r.Header.Get("X-Forwarded-Proto")
	if scheme == "" {
		scheme = "http"
	}
	return scheme + "://" + r.Host
}

func (s *Server) jwks(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, s.oidc.JWKS())
}

func writeOAuthError(
	w http.ResponseWriter,
	status int,
	code ileap.OAuthErrorCode,
	description string,
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(ileap.OAuthError{
		Code:        code,
		Description: description,
	}); err != nil {
		slog.Error("failed to encode OAuth error response", "error", err)
	}
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Error("failed to encode response", "error", err)
	}
}
