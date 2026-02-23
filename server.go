package ileap

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"mime"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"connectrpc.com/connect"
	ileapv1 "github.com/way-platform/ileap-go/proto/gen/wayplatform/connect/ileap/v1"
	"github.com/way-platform/ileap-go/proto/gen/wayplatform/connect/ileap/v1/ileapv1connect"
	"google.golang.org/protobuf/encoding/protojson"
)

// Server is an iLEAP data server HTTP handler.
//
// It translates the iLEAP HTTP protocol (JSON envelopes, OData filtering,
// Link header pagination, OAuth2 error formats) into calls on a standard
// Connect RPC service handler.
type Server struct {
	service    ileapv1connect.ILeapServiceHandler
	auth       AuthHandler
	pathPrefix string
	serveMux   *http.ServeMux
}

// ServerOption configures the server.
type ServerOption func(*Server)

// WithServiceHandler sets the ILeapService handler for footprints, TAD, and events.
func WithServiceHandler(h ileapv1connect.ILeapServiceHandler) ServerOption {
	return func(s *Server) { s.service = h }
}

// WithAuthHandler sets the auth handler for token issuance, validation, and OIDC discovery.
func WithAuthHandler(a AuthHandler) ServerOption {
	return func(s *Server) { s.auth = a }
}

// WithPathPrefix sets the path prefix for the service (e.g. "/ileap").
// Leading slashes are added if missing, and trailing slashes are trimmed.
func WithPathPrefix(p string) ServerOption {
	return func(s *Server) {
		if p == "" {
			s.pathPrefix = ""
			return
		}
		if !strings.HasPrefix(p, "/") {
			p = "/" + p
		}
		s.pathPrefix = strings.TrimRight(p, "/")
	}
}

// NewServer creates a new iLEAP data server with the given options.
func NewServer(opts ...ServerOption) *Server {
	s := &Server{
		serveMux: http.NewServeMux(),
	}
	for _, opt := range opts {
		opt(s)
	}
	s.setDefaults()
	s.registerRoutes()
	return s
}

// ServeHTTP implements http.Handler.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.serveMux.ServeHTTP(w, r)
}

func (s *Server) registerRoutes() {
	s.serveMux.Handle(
		"GET "+s.pathPrefix+"/2/footprints",
		s.pactAuthMiddleware(http.HandlerFunc(s.listFootprints)),
	)
	s.serveMux.Handle(
		"GET "+s.pathPrefix+"/2/footprints/{id}",
		s.pactAuthMiddleware(http.HandlerFunc(s.getFootprint)),
	)
	s.serveMux.Handle(
		"GET "+s.pathPrefix+"/2/ileap/tad",
		s.ileapAuthMiddleware(http.HandlerFunc(s.listTADs)),
	)
	s.serveMux.Handle(
		"POST "+s.pathPrefix+"/2/events",
		s.pactAuthMiddleware(http.HandlerFunc(s.events)),
	)
	s.serveMux.HandleFunc("POST "+s.pathPrefix+"/auth/token", s.authToken)
	// Workaround for ACT bug: PACT TC18/19 (OpenID Connect flow) mistakenly POSTs
	// to the base URL (/) instead of the token_endpoint advertised in
	// /.well-known/openid-configuration.
	s.serveMux.HandleFunc("POST "+s.pathPrefix+"/", s.authToken)
	s.serveMux.HandleFunc(
		"GET "+s.pathPrefix+"/.well-known/openid-configuration",
		s.openIDConfig,
	)
	s.serveMux.HandleFunc("GET "+s.pathPrefix+"/jwks", s.jwks)
}

// pactAuthMiddleware validates the bearer token using the server's TokenValidator
// and returns PACT-spec-formatted errors on failure (400 for missing/malformed,
// 401 for invalid). On success it calls next.
func (s *Server) pactAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			writeError(w, http.StatusBadRequest, ErrorCodeBadRequest, "missing authorization")
			return
		}
		if !strings.HasPrefix(authHeader, "Bearer ") {
			writeError(
				w,
				http.StatusBadRequest,
				ErrorCodeBadRequest,
				"unsupported authentication scheme",
			)
			return
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			writeError(w, http.StatusBadRequest, ErrorCodeBadRequest, "missing access token")
			return
		}
		if _, err := s.auth.ValidateToken(r.Context(), token); err != nil {
			if connect.CodeOf(err) == connect.CodeUnimplemented {
				writeError(w, http.StatusNotImplemented, ErrorCodeNotImplemented, "not implemented")
				return
			}
			slog.WarnContext(r.Context(), "token validation failed", "error", err)
			writeError(w, http.StatusUnauthorized, ErrorCodeAccessDenied, "invalid access token")
			return
		}
		ctx := WithAuthToken(r.Context(), token)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// ileapAuthMiddleware validates the bearer token using the server's AuthHandler
// and returns iLEAP-spec-formatted errors on failure (403 for invalid, 401 for
// expired). On success it calls next.
func (s *Server) ileapAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			writeError(w, http.StatusForbidden, ErrorCodeAccessDenied, "missing authorization")
			return
		}
		if !strings.HasPrefix(authHeader, "Bearer ") {
			writeError(
				w,
				http.StatusForbidden,
				ErrorCodeAccessDenied,
				"unsupported authentication scheme",
			)
			return
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			writeError(w, http.StatusForbidden, ErrorCodeAccessDenied, "missing access token")
			return
		}
		if _, err := s.auth.ValidateToken(r.Context(), token); err != nil {
			switch connect.CodeOf(err) {
			case connect.CodeUnimplemented:
				writeError(w, http.StatusNotImplemented, ErrorCodeNotImplemented, "not implemented")
			case connect.CodeUnauthenticated:
				slog.WarnContext(r.Context(), "token expired", "error", err)
				writeError(w, http.StatusUnauthorized, ErrorCodeTokenExpired, "token expired")
			default:
				slog.WarnContext(r.Context(), "token validation failed", "error", err)
				writeError(w, http.StatusForbidden, ErrorCodeAccessDenied, "invalid access token")
			}
			return
		}
		ctx := WithAuthToken(r.Context(), token)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *Server) resolveBaseURL(r *http.Request) string {
	scheme := r.Header.Get("X-Forwarded-Proto")
	if scheme == "" {
		if r.TLS != nil {
			scheme = "https"
		} else {
			scheme = "http"
		}
	}
	return scheme + "://" + r.Host + s.pathPrefix
}

func (s *Server) authToken(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
		writeOAuthError(
			w,
			http.StatusBadRequest,
			OAuthErrorCodeInvalidRequest,
			"invalid content type",
		)
		return
	}
	if err := r.ParseForm(); err != nil {
		writeOAuthError(
			w,
			http.StatusBadRequest,
			OAuthErrorCodeInvalidRequest,
			"invalid request body",
		)
		return
	}
	if grantType := r.Form.Get("grant_type"); grantType != "client_credentials" {
		writeOAuthError(
			w,
			http.StatusBadRequest,
			OAuthErrorCodeUnsupportedGrantType,
			"unsupported grant type",
		)
		return
	}
	username, password, ok := r.BasicAuth()
	if !ok {
		writeOAuthError(
			w,
			http.StatusBadRequest,
			OAuthErrorCodeInvalidRequest,
			"missing HTTP basic authorization",
		)
		return
	}

	clientID, err := url.QueryUnescape(username)
	if err != nil {
		clientID = username
	}
	clientSecret, err := url.QueryUnescape(password)
	if err != nil {
		clientSecret = password
	}

	creds, err := s.auth.IssueToken(r.Context(), clientID, clientSecret)
	if err != nil {
		switch connect.CodeOf(err) {
		case connect.CodeUnimplemented:
			writeError(w, http.StatusNotImplemented, ErrorCodeNotImplemented, "not implemented")
		case connect.CodePermissionDenied:
			slog.WarnContext(r.Context(), "failed to issue token", "error", err)
			writeOAuthError(
				w,
				http.StatusBadRequest,
				OAuthErrorCodeInvalidRequest,
				"invalid HTTP basic auth",
			)
		default:
			slog.ErrorContext(r.Context(), "failed to issue token", "error", err)
			writeOAuthError(
				w,
				http.StatusInternalServerError,
				OAuthErrorCodeServerError,
				"failed to issue token",
			)
		}
		return
	}
	writeJSON(w, creds)
}

func (s *Server) openIDConfig(w http.ResponseWriter, r *http.Request) {
	cfg := s.auth.OpenIDConfiguration(s.resolveBaseURL(r))
	if cfg == nil {
		writeError(w, http.StatusNotImplemented, ErrorCodeNotImplemented, "not implemented")
		return
	}
	writeJSON(w, cfg)
}

func (s *Server) jwks(w http.ResponseWriter, _ *http.Request) {
	jwks := s.auth.JWKS()
	if jwks == nil {
		writeError(w, http.StatusNotImplemented, ErrorCodeNotImplemented, "not implemented")
		return
	}
	writeJSON(w, jwks)
}

func (s *Server) listFootprints(w http.ResponseWriter, r *http.Request) {
	limit, err := parseLimit(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, ErrorCodeBadRequest, "invalid limit: %v", err)
		return
	}
	offset, err := parseOffset(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, ErrorCodeBadRequest, "invalid offset: %v", err)
		return
	}
	req := new(ileapv1.ListFootprintsRequest)
	req.SetLimit(int32(limit))
	req.SetOffset(int32(offset))
	req.SetFilter(r.URL.Query().Get("$filter"))
	resp, err := s.service.ListFootprints(r.Context(), req)
	if err != nil {
		writeHandlerError(w, err)
		return
	}
	data := resp.GetData()
	total := int(resp.GetTotal())
	next := offset + len(data)
	if next < total {
		linkLimit := limit
		if linkLimit == 0 {
			linkLimit = len(data)
		}
		base := s.resolveBaseURL(r)
		linkURL := fmt.Sprintf("%s/2/footprints?limit=%d&offset=%d", base, linkLimit, next)
		w.Header().Set("Link", fmt.Sprintf("<%s>; rel=\"next\"", linkURL))
	}
	writeListFootprintsResponse(w, data)
}

func (s *Server) getFootprint(w http.ResponseWriter, r *http.Request) {
	req := new(ileapv1.GetFootprintRequest)
	req.SetId(r.PathValue("id"))
	resp, err := s.service.GetFootprint(r.Context(), req)
	if err != nil {
		writeHandlerError(w, err)
		return
	}
	writeGetFootprintResponse(w, resp.GetData())
}

func (s *Server) listTADs(w http.ResponseWriter, r *http.Request) {
	limit, err := parseLimit(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, ErrorCodeBadRequest, "invalid limit: %v", err)
		return
	}
	offset, err := parseOffset(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, ErrorCodeBadRequest, "invalid offset: %v", err)
		return
	}
	req := new(ileapv1.ListTransportActivityDataRequest)
	req.SetLimit(int32(limit))
	req.SetOffset(int32(offset))
	q := r.URL.Query()
	if v := q.Get("mode"); v != "" {
		req.SetMode(v)
	}
	if v := q.Get("feedstock"); v != "" {
		req.SetFeedstock(v)
	}
	if v := q.Get("packagingOrTrEqType"); v != "" {
		req.SetPackagingOrTrEqType(v)
	}
	resp, err := s.service.ListTransportActivityData(r.Context(), req)
	if err != nil {
		writeHandlerError(w, err)
		return
	}
	data := resp.GetData()
	total := int(resp.GetTotal())
	next := offset + len(data)
	if next < total {
		linkLimit := limit
		if linkLimit == 0 {
			linkLimit = len(data)
		}
		base := s.resolveBaseURL(r)
		linkURL := fmt.Sprintf("%s/2/ileap/tad?offset=%d&limit=%d", base, next, linkLimit)
		w.Header().Set("Link", fmt.Sprintf("<%s>; rel=\"next\"", linkURL))
	}
	writeListTADsResponse(w, data)
}

func (s *Server) events(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") == "" {
		writeError(w, http.StatusBadRequest, ErrorCodeBadRequest, "missing content type")
		return
	}
	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		writeError(w, http.StatusBadRequest, ErrorCodeBadRequest, "invalid content type")
		return
	}
	if mediaType != "application/cloudevents+json" && mediaType != "application/json" {
		writeError(
			w,
			http.StatusBadRequest,
			ErrorCodeBadRequest,
			"invalid content type: %s",
			mediaType,
		)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(w, http.StatusBadRequest, ErrorCodeBadRequest, "failed to read request body")
		return
	}
	event := new(ileapv1.Event)
	if err := protojson.Unmarshal(body, event); err != nil {
		writeError(w, http.StatusBadRequest, ErrorCodeBadRequest, "invalid request body")
		return
	}
	if event.GetSpecversion() != "1.0" || event.GetId() == "" || event.GetSource() == "" {
		writeError(
			w,
			http.StatusBadRequest,
			ErrorCodeBadRequest,
			"missing required event fields",
		)
		return
	}
	if len(event.GetData()) == 0 {
		writeError(
			w,
			http.StatusBadRequest,
			ErrorCodeBadRequest,
			"missing event data",
		)
		return
	}
	req := new(ileapv1.HandleEventRequest)
	req.SetEvent(event)
	if _, err := s.service.HandleEvent(r.Context(), req); err != nil {
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
	switch connect.CodeOf(err) {
	case connect.CodeUnimplemented:
		writeError(w, http.StatusNotImplemented, ErrorCodeNotImplemented, "not implemented")
	case connect.CodeNotFound:
		writeError(w, http.StatusNotFound, ErrorCodeNoSuchFootprint, "%s", err)
	case connect.CodeInvalidArgument:
		writeError(w, http.StatusBadRequest, ErrorCodeBadRequest, "%s", err)
	default:
		slog.Error("handler error", "error", err)
		writeError(
			w,
			http.StatusInternalServerError,
			ErrorCodeInternalError,
			"internal error",
		)
	}
}

// writeOAuthError writes an OAuth 2.0 error response.
func writeOAuthError(
	w http.ResponseWriter,
	status int,
	code OAuthErrorCode,
	description string,
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(OAuthError{
		Code:        code,
		Description: description,
	}); err != nil {
		slog.Error("failed to encode OAuth error response", "error", err)
	}
}

// writeError writes a PACT-formatted JSON error response.
func writeError(
	w http.ResponseWriter,
	status int,
	code ErrorCode,
	format string,
	args ...any,
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(Error{
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

func writeGetFootprintResponse(w http.ResponseWriter, fp *ileapv1.ProductFootprint) {
	w.Header().Set("Content-Type", "application/json")
	data, err := protojson.Marshal(fp)
	if err != nil {
		slog.Error("failed to marshal footprint", "error", err)
		writeError(w, http.StatusInternalServerError, ErrorCodeInternalError, "internal error")
		return
	}
	if _, err := w.Write([]byte(`{"data":`)); err != nil {
		slog.Error("failed to write response", "error", err)
		return
	}
	if _, err := w.Write(data); err != nil {
		slog.Error("failed to write response", "error", err)
		return
	}
	if _, err := w.Write([]byte(`}`)); err != nil {
		slog.Error("failed to write response", "error", err)
		return
	}
}

func writeListFootprintsResponse(w http.ResponseWriter, fps []*ileapv1.ProductFootprint) {
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write([]byte(`{"data":[`)); err != nil {
		slog.Error("failed to write response", "error", err)
		return
	}
	for i, fp := range fps {
		if i > 0 {
			if _, err := w.Write([]byte(",")); err != nil {
				slog.Error("failed to write response", "error", err)
				return
			}
		}
		data, err := protojson.Marshal(fp)
		if err != nil {
			slog.Error("failed to marshal footprint", "error", err)
			writeError(w, http.StatusInternalServerError, ErrorCodeInternalError, "internal error")
			return
		}
		if _, err := w.Write(data); err != nil {
			slog.Error("failed to write response", "error", err)
			return
		}
	}
	if _, err := w.Write([]byte(`]}`)); err != nil {
		slog.Error("failed to write response", "error", err)
		return
	}
}

func writeListTADsResponse(w http.ResponseWriter, tads []*ileapv1.TAD) {
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write([]byte(`{"data":[`)); err != nil {
		slog.Error("failed to write response", "error", err)
		return
	}
	for i, tad := range tads {
		if i > 0 {
			if _, err := w.Write([]byte(",")); err != nil {
				slog.Error("failed to write response", "error", err)
				return
			}
		}
		data, err := protojson.Marshal(tad)
		if err != nil {
			slog.Error("failed to marshal TAD", "error", err)
			writeError(w, http.StatusInternalServerError, ErrorCodeInternalError, "internal error")
			return
		}
		if _, err := w.Write(data); err != nil {
			slog.Error("failed to write response", "error", err)
			return
		}
	}
	if _, err := w.Write([]byte(`]}`)); err != nil {
		slog.Error("failed to write response", "error", err)
		return
	}
}
