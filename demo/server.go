package demo

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/big"
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
func NewServer() (*Server, error) {
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
		baseURL:    "https://ileap-demo-server-504882905500.europe-north1.run.app",
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
		// Validate access token.
		payload, err := s.validateJWT(token)
		if err != nil {
			slog.Debug("JWT validation failed", "error", err)
			http.Error(w, "invalid access token", http.StatusUnauthorized)
			return
		}
		slog.Debug("JWT validated successfully", "username", payload.Username)
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
		// TODO: Generate and sign a JWT for the user.
		accessToken, err := s.createJWT(username)
		if err != nil {
			http.Error(w, "failed to create JWT", http.StatusInternalServerError)
			return
		}
		// Return client credentials.
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
		jwk := JWK{
			KeyType:   "RSA",
			Use:       "sig",
			Algorithm: "RS256",
			KeyID:     "Public key",
			N: base64.RawURLEncoding.EncodeToString(
				s.keypair.PublicKey.N.Bytes(),
			),
			E: base64.RawURLEncoding.EncodeToString(
				big.NewInt(int64(s.keypair.PublicKey.E)).Bytes(),
			),
		}
		jwks := JWKSet{
			Keys: []JWK{jwk},
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(jwks); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}

// createJWT creates and signs a JWT token using RS256 algorithm
func (s *Server) createJWT(username string) (string, error) {
	// Create header
	header := JWTHeader{
		Type:      "JWT",
		Algorithm: "RS256",
	}
	headerBytes, err := json.Marshal(header)
	if err != nil {
		return "", fmt.Errorf("marshal header: %w", err)
	}
	headerEncoded := base64.RawURLEncoding.EncodeToString(headerBytes)
	// Create payload
	payload := Claims{
		Username: username,
		IssuedAt: time.Now().Unix(),
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("marshal payload: %w", err)
	}
	payloadEncoded := base64.RawURLEncoding.EncodeToString(payloadBytes)
	// Create signing input
	signingInput := headerEncoded + "." + payloadEncoded
	// Sign with RSA-SHA256
	hash := sha256.Sum256([]byte(signingInput))
	signature, err := rsa.SignPKCS1v15(rand.Reader, s.keypair.PrivateKey, crypto.SHA256, hash[:])
	if err != nil {
		return "", fmt.Errorf("sign JWT: %w", err)
	}
	signatureEncoded := base64.RawURLEncoding.EncodeToString(signature)
	// Combine all parts
	token := signingInput + "." + signatureEncoded
	return token, nil
}

type JWKSet struct {
	Keys []JWK `json:"keys"`
}

type JWK struct {
	// KeyType specifies the cryptographic algorithm family used with the key
	KeyType string `json:"kty"`
	// Use identifies the intended use of the public key
	Use string `json:"use,omitempty"`
	// Algorithm identifies the algorithm intended for use with the key
	Algorithm string `json:"alg,omitempty"`
	// KeyID is a hint indicating which key was used to secure the JWS
	KeyID string `json:"kid,omitempty"`
	// N is the modulus for the RSA public key (base64url encoded).
	N string `json:"n"`
	// E is the exponent for the RSA public key (base64url encoded).
	E string `json:"e"`
}

type OpenIDConfiguration struct {
	// IssuerURL is the identity of the provider, and the string it uses to sign
	// ID tokens with. For example "https://accounts.google.com". This value MUST
	// match ID tokens exactly.
	IssuerURL string `json:"issuer"`

	// AuthURL is the endpoint used by the provider to support the OAuth 2.0
	// authorization endpoint.
	AuthURL string `json:"authorization_endpoint"`

	// TokenURL is the endpoint used by the provider to support the OAuth 2.0
	// token endpoint.
	TokenURL string `json:"token_endpoint"`

	// DeviceAuthURL is the endpoint used by the provider to support the OAuth 2.0
	// device authorization endpoint.
	DeviceAuthURL string `json:"device_authorization_endpoint"`

	// UserInfoURL is the endpoint used by the provider to support the OpenID
	// Connect UserInfo flow.
	//
	// https://openid.net/specs/openid-connect-core-1_0.html#UserInfo
	UserInfoURL string `json:"userinfo_endpoint"`

	// JWKSURL is the endpoint used by the provider to advertise public keys to
	// verify issued ID tokens. This endpoint is polled as new keys are made
	// available.
	JWKSURL string `json:"jwks_uri"`

	// Algorithms, if provided, indicate a list of JWT algorithms allowed to sign
	// ID tokens. If not provided, this defaults to the algorithms advertised by
	// the JWK endpoint, then the set of algorithms supported by this package.
	Algorithms []string `json:"id_token_signing_alg_values_supported"`
}

type JWTHeader struct {
	Type      string `json:"typ"`
	Algorithm string `json:"alg"`
}

type Claims struct {
	Username string `json:"username"`
	IssuedAt int64  `json:"iat,omitempty"`
}
