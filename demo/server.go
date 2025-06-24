package demo

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"mime"
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
	server.registerAuthenticatedRoute(server.eventsRoute())
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
			s.errorf(w, http.StatusUnauthorized, ileap.ErrorCodeBadRequest, "missing authorization")
			return
		}
		if !strings.HasPrefix(auth, "Bearer ") {
			s.errorf(w, http.StatusUnauthorized, ileap.ErrorCodeBadRequest, "unsupported authentication scheme")
			return
		}
		token := strings.TrimPrefix(auth, "Bearer ")
		if token == "" {
			s.errorf(w, http.StatusUnauthorized, ileap.ErrorCodeBadRequest, "missing access token")
			return
		}
		if _, err := s.keypair.ValidateJWT(token); err != nil {
			// TODO: ACT conformance test requires 400, but semantically this should be 401.
			s.errorf(w, http.StatusBadRequest, ileap.ErrorCodeBadRequest, "invalid access token")
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

func (s *Server) errorf(w http.ResponseWriter, status int, errorCode ileap.ErrorCode, format string, args ...any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(ileap.Error{
		Code:    errorCode,
		Message: fmt.Sprintf(format, args...),
	}); err != nil {
		slog.Error("failed to encode error response", "error", err)
		return
	}
}

func (s *Server) oauthError(w http.ResponseWriter, status int, code ileap.OAuthErrorCode, description string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(ileap.OAuthError{
		Code:        code,
		Description: description,
	}); err != nil {
		slog.Error("failed to encode OAuth error response", "error", err)
		return
	}
}

func (s *Server) authTokenRoute() (string, http.HandlerFunc) {
	return "POST /auth/token", func(w http.ResponseWriter, r *http.Request) {
		// Validate content type.
		if r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
			s.oauthError(w, http.StatusBadRequest, ileap.OAuthErrorCodeInvalidRequest, "invalid content type")
			return
		}
		// Parse URL values from request body.
		if err := r.ParseForm(); err != nil {
			s.oauthError(w, http.StatusBadRequest, ileap.OAuthErrorCodeInvalidRequest, "invalid request body")
			return
		}
		// Validate grant type.
		grantType := r.Form.Get("grant_type")
		if grantType != "client_credentials" {
			s.oauthError(w, http.StatusBadRequest, ileap.OAuthErrorCodeUnsupportedGrantType, "unsupported grant type")
			return
		}
		// Validate HTTP Basic Auth credentials.
		username, password, ok := r.BasicAuth()
		if !ok {
			s.oauthError(w, http.StatusBadRequest, ileap.OAuthErrorCodeInvalidRequest, "missing HTTP basic authorization")
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
			// TODO: ACT conformance test requires 400, but semantically this should be 401.
			s.oauthError(w, http.StatusBadRequest, ileap.OAuthErrorCodeInvalidRequest, "invalid HTTP basic auth")
			return
		}
		accessToken, err := s.keypair.CreateJWT(JWTClaims{
			Username: username,
			IssuedAt: time.Now().Unix(),
		})
		if err != nil {
			s.oauthError(w, http.StatusInternalServerError, ileap.OAuthErrorCodeServerError, "failed to create JWT")
			return
		}
		clientCredentials := ileap.ClientCredentials{
			AccessToken: accessToken,
			TokenType:   "bearer",
		}
		if err := json.NewEncoder(w).Encode(clientCredentials); err != nil {
			s.oauthError(w, http.StatusInternalServerError, ileap.OAuthErrorCodeServerError, "failed to encode response")
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
			s.errorf(w, http.StatusInternalServerError, ileap.ErrorCodeInternalError, "failed to encode response")
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
			s.errorf(w, http.StatusNotFound, ileap.ErrorCodeNotFound, "not found")
			return
		}
		w.Header().Set("Content-Type", "application/json")
		response := ileapv0.ProductFootprintResponse{
			Data: *footprint,
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			s.errorf(w, http.StatusInternalServerError, ileap.ErrorCodeInternalError, "failed to encode response")
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
			s.errorf(w, http.StatusInternalServerError, ileap.ErrorCodeInternalError, "failed to encode response")
			return
		}
	}
}

func (s *Server) openIDConnectConfigRoute() (string, http.HandlerFunc) {
	return "GET /.well-known/openid-configuration", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		response := OpenIDConfiguration{
			IssuerURL:              s.baseURL,
			AuthURL:                s.baseURL + "/auth/token",
			TokenURL:               s.baseURL + "/auth/token",
			JWKSURL:                s.baseURL + "/jwks",
			Algorithms:             []string{"RS256"},
			ResponseTypesSupported: []string{"token"},
			SubjectTypesSupported:  []string{"public"},
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			s.errorf(w, http.StatusInternalServerError, ileap.ErrorCodeInternalError, "failed to encode response")
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
			s.errorf(w, http.StatusInternalServerError, ileap.ErrorCodeInternalError, "failed to encode response")
			return
		}
	}
}

func (s *Server) eventsRoute() (string, http.HandlerFunc) {
	return "POST /2/events", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") == "" {
			s.errorf(w, http.StatusBadRequest, ileap.ErrorCodeBadRequest, "missing content type")
			return
		}
		if mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type")); err == nil {
			if mediaType != "application/cloudevents+json" {
				s.errorf(w, http.StatusBadRequest, ileap.ErrorCodeBadRequest, "invalid content type: %s", mediaType)
				return
			}
		} else {
			s.errorf(w, http.StatusBadRequest, ileap.ErrorCodeBadRequest, "invalid content type")
		}
		// Parse the event from request body.
		var event ileap.Event
		if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
			s.errorf(w, http.StatusBadRequest, ileap.ErrorCodeBadRequest, "invalid request body")
			return
		}
		switch event.Type {
		case ileap.EventTypeRequestCreated:
			// TODO: Handle RequestCreated.
		case ileap.EventTypeRequestFulfilled:
			// TODO: Handle RequestFulfilled.
		case ileap.EventTypeRequestRejected:
			// TODO: Handle RequestRejected.
		case ileap.EventTypePublished:
			// TODO: Handle Published.
		default:
			s.errorf(w, http.StatusBadRequest, ileap.ErrorCodeBadRequest, "invalid event type")
		}
		var response struct {
			Status  string `json:"status"`
			Message string `json:"message"`
		}
		response.Status = "accepted"
		response.Message = "Event successfully processed"
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			s.errorf(w, http.StatusInternalServerError, ileap.ErrorCodeInternalError, "failed to encode response")
			return
		}
	}
}
