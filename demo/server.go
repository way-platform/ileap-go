package demo

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/way-platform/ileap-go"
	"github.com/way-platform/ileap-go/openapi/ileapv0"
)

// Server is an example iLEAP server.
type Server struct {
	footprints []ileapv0.ProductFootprintForILeapType
	tads       []ileapv0.TAD
	serveMux   *http.ServeMux
}

// NewServer creates a new example iLEAP server.
func NewServer() (*Server, error) {
	footprints, err := Footprints()
	if err != nil {
		return nil, fmt.Errorf("load footprints: %w", err)
	}
	tads, err := TADs()
	if err != nil {
		return nil, fmt.Errorf("load tads: %w", err)
	}
	server := &Server{
		footprints: footprints,
		tads:       tads,
		serveMux:   http.NewServeMux(),
	}
	server.registerRoute(server.authTokenRoute())
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
		// TODO: Validate access token.
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
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}
		if username != "hello" || password != "pathfinder" {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}
		// Return client credentials.
		clientCredentials := ileap.ClientCredentials{
			AccessToken: "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJ1c2VybmFtZSI6ImhlbGxvIn0.p3qamb1KAOwE3clLDykok1fL7WDI5bQYziYUIP_XYRY_ysZC0SeKw_n3EO7sgxB26Bh33UdJDGcchZKw3oM6NfjCbT4lG8tECIPoZC1Vg-2RUl-wCLhNQzcBXX3i3UKWMn9z9TBv-KBeAq7Y6mwaeD4DSnlYAFIo0r_iON_9emMmTUf09iOnOlXjzRrbIfLOlrrc5wi4rIz0tRR-563yBGVpSokXag8LfSj5S_Nj7LOgFIGtRDUKoqZvZuyRMz_wQYXC5T9pz2h58D2ImkpoANLH7669FdO71dX_cCMgju7vTp9UjZuj_Xi73HU4mJQaFh2_g6iSCZl-wxOWJpo_Qw",
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
