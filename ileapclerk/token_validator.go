package ileapclerk

import (
	"context"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/way-platform/ileap-go/ileapserver"
)

// TokenValidator implements ileapserver.TokenValidator using Clerk's JWKS.
type TokenValidator struct {
	client *Client
	mu     sync.Mutex
	jwks   *ileapserver.JWKSet
}

// NewTokenValidator creates a new token validator backed by Clerk's JWKS.
func NewTokenValidator(client *Client) *TokenValidator {
	return &TokenValidator{client: client}
}

// ValidateToken validates a Clerk RS256 JWT against Clerk's JWKS.
func (v *TokenValidator) ValidateToken(
	_ context.Context, token string,
) (*ileapserver.TokenInfo, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid JWT format")
	}
	headerBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, fmt.Errorf("decode JWT header: %w", err)
	}
	var header struct {
		Kid string `json:"kid"`
		Alg string `json:"alg"`
	}
	if err := json.Unmarshal(headerBytes, &header); err != nil {
		return nil, fmt.Errorf("parse JWT header: %w", err)
	}
	pub, err := v.findKey(header.Kid)
	if err != nil {
		return nil, fmt.Errorf("find signing key: %w", err)
	}
	sigBytes, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return nil, fmt.Errorf("decode JWT signature: %w", err)
	}
	message := parts[0] + "." + parts[1]
	digest := sha256.Sum256([]byte(message))
	if err := rsa.VerifyPKCS1v15(pub, crypto.SHA256, digest[:], sigBytes); err != nil {
		return nil, fmt.Errorf("invalid JWT signature: %w", err)
	}
	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("decode JWT payload: %w", err)
	}
	var claims map[string]any
	if err := json.Unmarshal(payloadBytes, &claims); err != nil {
		return nil, fmt.Errorf("parse JWT claims: %w", err)
	}
	if exp, ok := claims["exp"].(float64); ok {
		if time.Now().Unix() > int64(exp) {
			return nil, fmt.Errorf("JWT expired")
		}
	}
	sub, _ := claims["sub"].(string)
	return &ileapserver.TokenInfo{Subject: sub}, nil
}

func (v *TokenValidator) findKey(kid string) (*rsa.PublicKey, error) {
	v.mu.Lock()
	defer v.mu.Unlock()
	if v.jwks != nil {
		if pub := findKeyInSet(v.jwks, kid); pub != nil {
			return pub, nil
		}
	}
	jwks, err := v.client.FetchJWKS()
	if err != nil {
		return nil, fmt.Errorf("fetch JWKS: %w", err)
	}
	v.jwks = jwks
	pub := findKeyInSet(jwks, kid)
	if pub == nil {
		return nil, fmt.Errorf("key %q not found in JWKS", kid)
	}
	return pub, nil
}

func findKeyInSet(jwks *ileapserver.JWKSet, kid string) *rsa.PublicKey {
	for _, jwk := range jwks.Keys {
		if jwk.KeyID != kid {
			continue
		}
		nBytes, err := base64.RawURLEncoding.DecodeString(jwk.N)
		if err != nil {
			continue
		}
		eBytes, err := base64.RawURLEncoding.DecodeString(jwk.E)
		if err != nil {
			continue
		}
		return &rsa.PublicKey{
			N: new(big.Int).SetBytes(nBytes),
			E: int(new(big.Int).SetBytes(eBytes).Int64()),
		}
	}
	return nil
}
