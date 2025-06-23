package demo

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/way-platform/ileap-go"
	"github.com/way-platform/ileap-go/openapi/ileapv0"
)

// Server is an example iLEAP server.
type Server struct {
	baseURL    string
	footprints []ileapv0.ProductFootprintForILeapType
	tads       []ileapv0.TAD
	keypair    *KeyPair
	serveMux   *http.ServeMux
}

// NewServer creates a new example iLEAP server.
func NewServer(baseURL string) (*Server, error) {
	footprints, err := LoadFootprints()
	if err != nil {
		return nil, fmt.Errorf("load footprints: %w", err)
	}
	tads, err := LoadTADs()
	if err != nil {
		return nil, fmt.Errorf("load tads: %w", err)
	}
	keypair, err := LoadKeyPair()
	if err != nil {
		return nil, fmt.Errorf("load keypair: %w", err)
	}
	server := &Server{
		baseURL:    baseURL,
		footprints: footprints,
		tads:       tads,
		keypair:    keypair,
		serveMux:   http.NewServeMux(),
	}
	server.registerRoute(server.authTokenRoute())
	server.registerRoute(server.openIDConnectConfigRoute())
	server.registerRoute(server.jwksRoute())
	server.registerAuthenticatedRoute(server.listFootprintsRoute())
	server.registerAuthenticatedRoute(server.getFootprintRoute())
	server.registerAuthenticatedRoute(server.listTADsRoute())
	return server, nil
}

// Handler returns the HTTP handler for the server.
func (s *Server) Handler() http.Handler {
	return s.serveMux
}

func (s *Server) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" {
			http.Error(w, "missing authorization", http.StatusUnauthorized)
			return
		}
		if !strings.HasPrefix(auth, "Bearer ") {
			http.Error(w, "unsupported authorization scheme", http.StatusUnauthorized)
			return
		}
		token := strings.TrimPrefix(auth, "Bearer ")
		if token == "" {
			http.Error(w, "missing access token", http.StatusUnauthorized)
			return
		}
		if _, err := s.keypair.ValidateJWT(token); err != nil {
			http.Error(w, "invalid access token", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s *Server) registerAuthenticatedRoute(pattern string, handler http.Handler) {
	slog.Debug("registering authenticated route", "pattern", pattern)
	s.serveMux.Handle(pattern, s.authMiddleware(handler))
}

func (s *Server) registerRoute(pattern string, handler http.Handler) {
	slog.Debug("registering route", "pattern", pattern)
	s.serveMux.Handle(pattern, handler)
}

func (s *Server) authTokenRoute() (string, http.HandlerFunc) {
	return "POST /auth/token", func(w http.ResponseWriter, r *http.Request) {
		// Validate content type.
		if r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
			http.Error(w, "invalid content type", http.StatusBadRequest)
			return
		}
		// Parse URL values from request body.
		if err := r.ParseForm(); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}
		// Validate grant type.
		grantType := r.Form.Get("grant_type")
		if grantType != "client_credentials" {
			http.Error(w, "unsupported grant type", http.StatusBadRequest)
			return
		}
		// Validate HTTP Basic Auth credentials.
		username, password, ok := r.BasicAuth()
		if !ok {
			http.Error(w, "missing HTTP basic authorization", http.StatusBadRequest)
			return
		}
		var authorized bool
		for _, user := range Users() {
			if username == user.Username && password == user.Password {
				authorized = true
				break
			}
		}
		if !authorized {
			http.Error(w, "invalid HTTP basic authorization", http.StatusUnauthorized)
			return
		}
		accessToken, err := s.keypair.CreateJWT(JWTClaims{
			Username: username,
			IssuedAt: time.Now().Unix(),
		})
		if err != nil {
			http.Error(w, "failed to create JWT", http.StatusInternalServerError)
			return
		}
		clientCredentials := ileap.ClientCredentials{
			AccessToken: accessToken,
			TokenType:   "bearer",
		}
		if err := json.NewEncoder(w).Encode(clientCredentials); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) listFootprintsRoute() (string, http.HandlerFunc) {
	return "GET /2/footprints", func(w http.ResponseWriter, r *http.Request) {
		// TODO: Handle limit.
		// TODO: Handle filter.
		w.Header().Set("Content-Type", "application/json")
		response := ileapv0.PfListingResponseInner{
			Data: s.footprints,
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) getFootprintRoute() (string, http.HandlerFunc) {
	return "GET /2/footprints/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		var footprint *ileapv0.ProductFootprintForILeapType
		for _, needle := range s.footprints {
			if needle.ID == id {
				footprint = &needle
				break
			}
		}
		if footprint == nil {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		response := ileapv0.ProductFootprintResponse{
			Data: *footprint,
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) listTADsRoute() (string, http.HandlerFunc) {
	return "GET /2/ileap/tad", func(w http.ResponseWriter, r *http.Request) {
		// TODO: Handle limit.
		w.Header().Set("Content-Type", "application/json")
		response := ileapv0.TadListingResponseInner{
			Data: s.tads,
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) openIDConnectConfigRoute() (string, http.HandlerFunc) {
	return "GET /.well-known/openid-configuration", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		response := OpenIDConfiguration{
			IssuerURL:  s.baseURL,
			AuthURL:    s.baseURL + "/auth/token",
			TokenURL:   s.baseURL + "/auth/token",
			JWKSURL:    s.baseURL + "/jwks",
			Algorithms: []string{"RS256"},
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) jwksRoute() (string, http.HandlerFunc) {
	return "GET /jwks", func(w http.ResponseWriter, r *http.Request) {
		jwks := JWKSet{
			Keys: []JWK{s.keypair.JWK()},
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(jwks); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}
