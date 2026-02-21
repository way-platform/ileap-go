package ileapserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"mime"
	"net/http"
	"strconv"
	"strings"

	"github.com/way-platform/ileap-go"
	"github.com/way-platform/ileap-go/openapi/ileapv0"
)

// Server is an iLEAP data server HTTP handler.
type Server struct {
	footprintHandler FootprintHandler
	tadHandler       TADHandler
	eventHandler     EventHandler
	tokenValidator   TokenValidator
	serveMux         *http.ServeMux
}

// Option configures the server.
type Option func(*Server)

// WithFootprintHandler sets the footprint handler.
func WithFootprintHandler(h FootprintHandler) Option {
	return func(s *Server) { s.footprintHandler = h }
}

// WithTADHandler sets the TAD handler.
func WithTADHandler(h TADHandler) Option {
	return func(s *Server) { s.tadHandler = h }
}

// WithEventHandler sets the event handler.
func WithEventHandler(h EventHandler) Option {
	return func(s *Server) { s.eventHandler = h }
}

// WithTokenValidator sets the token validator.
func WithTokenValidator(v TokenValidator) Option {
	return func(s *Server) { s.tokenValidator = v }
}

// NewServer creates a new iLEAP data server with the given options.
func NewServer(opts ...Option) *Server {
	s := &Server{
		serveMux: http.NewServeMux(),
	}
	for _, opt := range opts {
		opt(s)
	}
	s.registerRoutes()
	return s
}

// ServeHTTP implements http.Handler.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.serveMux.ServeHTTP(w, r)
}

func (s *Server) registerRoutes() {
	s.serveMux.Handle("GET /2/footprints", s.pactAuthMiddleware(http.HandlerFunc(s.listFootprints)))
	s.serveMux.Handle("GET /2/footprints/{id}", s.pactAuthMiddleware(http.HandlerFunc(s.getFootprint)))
	s.serveMux.Handle("GET /2/ileap/tad", s.ileapAuthMiddleware(http.HandlerFunc(s.listTADs)))
	s.serveMux.Handle("POST /2/events", s.pactAuthMiddleware(http.HandlerFunc(s.events)))
}

// pactAuthMiddleware returns 400 BadRequest for all auth failures (PACT spec).
func (s *Server) pactAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if s.tokenValidator == nil {
			writeError(w, http.StatusBadRequest, ileap.ErrorCodeBadRequest, "no token validator configured")
			return
		}
		auth := r.Header.Get("Authorization")
		if auth == "" {
			writeError(w, http.StatusBadRequest, ileap.ErrorCodeBadRequest, "missing authorization")
			return
		}
		if !strings.HasPrefix(auth, "Bearer ") {
			writeError(w, http.StatusBadRequest, ileap.ErrorCodeBadRequest, "unsupported authentication scheme")
			return
		}
		token := strings.TrimPrefix(auth, "Bearer ")
		if token == "" {
			writeError(w, http.StatusBadRequest, ileap.ErrorCodeBadRequest, "missing access token")
			return
		}
		if _, err := s.tokenValidator.ValidateToken(r.Context(), token); err != nil {
			writeError(w, http.StatusBadRequest, ileap.ErrorCodeBadRequest, "invalid access token")
			return
		}
		next.ServeHTTP(w, r)
	})
}

// ileapAuthMiddleware returns 403 AccessDenied for invalid tokens and
// 401 TokenExpired for expired tokens (iLEAP spec).
func (s *Server) ileapAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if s.tokenValidator == nil {
			writeError(w, http.StatusForbidden, ileap.ErrorCodeAccessDenied, "no token validator configured")
			return
		}
		auth := r.Header.Get("Authorization")
		if auth == "" {
			writeError(w, http.StatusForbidden, ileap.ErrorCodeAccessDenied, "missing authorization")
			return
		}
		if !strings.HasPrefix(auth, "Bearer ") {
			writeError(w, http.StatusForbidden, ileap.ErrorCodeAccessDenied, "unsupported authentication scheme")
			return
		}
		token := strings.TrimPrefix(auth, "Bearer ")
		if token == "" {
			writeError(w, http.StatusForbidden, ileap.ErrorCodeAccessDenied, "missing access token")
			return
		}
		if _, err := s.tokenValidator.ValidateToken(r.Context(), token); err != nil {
			if errors.Is(err, ErrTokenExpired) {
				writeError(w, http.StatusUnauthorized, ileap.ErrorCodeTokenExpired, "token expired")
				return
			}
			writeError(w, http.StatusForbidden, ileap.ErrorCodeAccessDenied, "invalid access token")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s *Server) listFootprints(w http.ResponseWriter, r *http.Request) {
	if s.footprintHandler == nil {
		writeError(w, http.StatusNotImplemented, ileap.ErrorCodeNotImplemented, "not implemented")
		return
	}
	limit, err := parseLimit(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, ileap.ErrorCodeBadRequest, "invalid limit: %v", err)
		return
	}
	req := ListFootprintsRequest{
		Limit:  limit,
		Filter: r.URL.Query().Get("$filter"),
	}
	resp, err := s.footprintHandler.ListFootprints(r.Context(), req)
	if err != nil {
		writeHandlerError(w, err)
		return
	}
	writeJSON(w, ileapv0.PfListingResponseInner{Data: resp.Data})
}

func (s *Server) getFootprint(w http.ResponseWriter, r *http.Request) {
	if s.footprintHandler == nil {
		writeError(w, http.StatusNotImplemented, ileap.ErrorCodeNotImplemented, "not implemented")
		return
	}
	id := r.PathValue("id")
	fp, err := s.footprintHandler.GetFootprint(r.Context(), id)
	if err != nil {
		writeHandlerError(w, err)
		return
	}
	writeJSON(w, ileapv0.ProductFootprintResponse{Data: *fp})
}

func (s *Server) listTADs(w http.ResponseWriter, r *http.Request) {
	if s.tadHandler == nil {
		writeError(w, http.StatusNotImplemented, ileap.ErrorCodeNotImplemented, "not implemented")
		return
	}
	limit, err := parseLimit(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, ileap.ErrorCodeBadRequest, "invalid limit: %v", err)
		return
	}
	offset, err := parseOffset(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, ileap.ErrorCodeBadRequest, "invalid offset: %v", err)
		return
	}
	// All query params except limit and offset are TAD filters.
	filters := make(map[string][]string)
	for key, values := range r.URL.Query() {
		if key == "limit" || key == "offset" {
			continue
		}
		filters[key] = values
	}
	req := ListTADsRequest{
		Limit:  limit,
		Offset: offset,
		Filter: filters,
	}
	resp, err := s.tadHandler.ListTADs(r.Context(), req)
	if err != nil {
		writeHandlerError(w, err)
		return
	}
	// Emit Link header for pagination when more data is available.
	next := offset + len(resp.Data)
	if next < resp.Total {
		host := r.Host
		scheme := "https"
		if r.TLS == nil {
			scheme = "http"
		}
		if fwd := r.Header.Get("X-Forwarded-Proto"); fwd != "" {
			scheme = fwd
		}
		linkLimit := limit
		if linkLimit == 0 {
			linkLimit = len(resp.Data)
		}
		linkURL := fmt.Sprintf("%s://%s/2/ileap/tad?offset=%d&limit=%d", scheme, host, next, linkLimit)
		w.Header().Set("Link", fmt.Sprintf("<%s>; rel=\"next\"", linkURL))
	}
	writeJSON(w, ileapv0.TadListingResponseInner{Data: resp.Data})
}

func (s *Server) events(w http.ResponseWriter, r *http.Request) {
	if s.eventHandler == nil {
		writeError(w, http.StatusNotImplemented, ileap.ErrorCodeNotImplemented, "not implemented")
		return
	}
	if r.Header.Get("Content-Type") == "" {
		writeError(w, http.StatusBadRequest, ileap.ErrorCodeBadRequest, "missing content type")
		return
	}
	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		writeError(w, http.StatusBadRequest, ileap.ErrorCodeBadRequest, "invalid content type")
		return
	}
	// PACT specification requires "application/cloudevents+json",
	// but the conformance checker also sends application/json.
	if mediaType != "application/cloudevents+json" && mediaType != "application/json" {
		writeError(w, http.StatusBadRequest, ileap.ErrorCodeBadRequest, "invalid content type: %s", mediaType)
		return
	}
	var event Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		writeError(w, http.StatusBadRequest, ileap.ErrorCodeBadRequest, "invalid request body")
		return
	}
	if err := s.eventHandler.HandleEvent(r.Context(), event); err != nil {
		writeHandlerError(w, err)
		return
	}
}

func parseLimit(r *http.Request) (int, error) {
	limitStr := r.URL.Query().Get("limit")
	if limitStr == "" {
		return 0, nil
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return 0, err
	}
	if limit <= 0 {
		return 0, fmt.Errorf("limit must be positive")
	}
	return limit, nil
}

func parseOffset(r *http.Request) (int, error) {
	offsetStr := r.URL.Query().Get("offset")
	if offsetStr == "" {
		return 0, nil
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		return 0, err
	}
	if offset < 0 {
		return 0, fmt.Errorf("offset must be non-negative")
	}
	return offset, nil
}

func writeHandlerError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrNotFound):
		writeError(w, http.StatusNotFound, ileap.ErrorCodeNoSuchFootprint, "%s", err)
	case errors.Is(err, ErrBadRequest):
		writeError(w, http.StatusBadRequest, ileap.ErrorCodeBadRequest, "%s", err)
	default:
		writeError(w, http.StatusInternalServerError, ileap.ErrorCodeInternalError, "internal error")
	}
}

func writeError(w http.ResponseWriter, status int, code ileap.ErrorCode, format string, args ...any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(ileap.Error{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
	}); err != nil {
		slog.Error("failed to encode error response", "error", err)
	}
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Error("failed to encode response", "error", err)
	}
}
