package ileapdemo

import (
	"fmt"
	"net/http"

	"github.com/way-platform/ileap-go/ileapauthserver"
	"github.com/way-platform/ileap-go/ileapserver"
)

// Server is an example iLEAP server composing data and auth adapters.
type Server struct {
	serveMux *http.ServeMux
}

// NewServer creates a new example iLEAP server.
func NewServer(baseURL string) (*Server, error) {
	dataHandler, err := NewDataHandler()
	if err != nil {
		return nil, fmt.Errorf("create data handler: %w", err)
	}
	authProvider, err := NewAuthProvider()
	if err != nil {
		return nil, fmt.Errorf("create auth provider: %w", err)
	}
	return NewServerWith(baseURL, authProvider, dataHandler, authProvider), nil
}

// NewServerWith creates a new example iLEAP server with explicit dependencies.
func NewServerWith(
	baseURL string,
	auth interface {
		ileapauthserver.TokenIssuer
		ileapauthserver.OIDCProvider
		ileapserver.TokenValidator
	},
	data interface {
		ileapserver.FootprintHandler
		ileapserver.TADHandler
	},
	tokenValidator ileapserver.TokenValidator,
) *Server {
	dataSrv := ileapserver.NewServer(
		ileapserver.WithFootprintHandler(data),
		ileapserver.WithTADHandler(data),
		ileapserver.WithEventHandler(&EventHandler{}),
		ileapserver.WithTokenValidator(tokenValidator),
	)
	authSrv := ileapauthserver.NewServer(baseURL, auth, auth)
	mux := http.NewServeMux()
	mux.Handle("/auth/", authSrv)
	mux.Handle("/.well-known/", authSrv)
	mux.Handle("/jwks", authSrv)
	mux.Handle("/", dataSrv)
	return &Server{serveMux: mux}
}

// Handler returns the HTTP handler for the server.
func (s *Server) Handler() http.Handler {
	return s.serveMux
}
