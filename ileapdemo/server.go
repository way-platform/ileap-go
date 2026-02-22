package ileapdemo

import (
	"fmt"
	"net/http"

	"github.com/way-platform/ileap-go/ileapserver"
)

// Server is an example iLEAP server composing data and auth adapters.
type Server struct {
	handler http.Handler
}

// NewServer creates a new example iLEAP server.
func NewServer() (*Server, error) {
	dataHandler, err := NewDataHandler()
	if err != nil {
		return nil, fmt.Errorf("create data handler: %w", err)
	}
	authProvider, err := NewAuthProvider()
	if err != nil {
		return nil, fmt.Errorf("create auth provider: %w", err)
	}
	return NewServerWith(authProvider, dataHandler, authProvider), nil
}

// NewServerWith creates a new example iLEAP server with explicit dependencies.
func NewServerWith(
	auth interface {
		ileapserver.TokenIssuer
		ileapserver.OIDCProvider
		ileapserver.TokenValidator
	},
	data interface {
		ileapserver.FootprintHandler
		ileapserver.TADHandler
	},
	tokenValidator ileapserver.TokenValidator,
) *Server {
	srv := ileapserver.NewServer(
		ileapserver.WithFootprintHandler(data),
		ileapserver.WithTADHandler(data),
		ileapserver.WithEventHandler(&EventHandler{}),
		ileapserver.WithTokenValidator(tokenValidator),
		ileapserver.WithTokenIssuer(auth),
		ileapserver.WithOIDCProvider(auth),
	)
	return &Server{handler: srv}
}

// Handler returns the HTTP handler for the server.
func (s *Server) Handler() http.Handler {
	return s.handler
}
