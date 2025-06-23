package demo

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"log/slog"
	"math/big"
	"strings"

	_ "embed"
)

//go:embed testdata/keypair.pem
var keypairData []byte

// KeyPair represents an RSA keypair for JWT operations.
type KeyPair struct {
	// PrivateKey is the RSA private key.
	PrivateKey *rsa.PrivateKey
	// PublicKey is the RSA public key.
	PublicKey *rsa.PublicKey
}

// JWK returns the public key in JWK format.
func (k *KeyPair) JWK() JWK {
	return JWK{
		KeyType:   "RSA",
		Use:       "sig",
		Algorithm: "RS256",
		KeyID:     "Public key",
		N: base64.RawURLEncoding.EncodeToString(
			k.PublicKey.N.Bytes(),
		),
		E: base64.RawURLEncoding.EncodeToString(
			big.NewInt(int64(k.PublicKey.E)).Bytes(),
		),
	}
}

// CreateJWT creates a JSON Web Token with the given claims.
func (k *KeyPair) CreateJWT(claims JWTClaims) (string, error) {
	header := JWTHeader{
		Type:      "JWT",
		Algorithm: "RS256",
	}
	headerBytes, err := json.Marshal(header)
	if err != nil {
		return "", fmt.Errorf("marshal header: %w", err)
	}
	headerEncoded := base64.RawURLEncoding.EncodeToString(headerBytes)
	payloadBytes, err := json.Marshal(claims)
	if err != nil {
		return "", fmt.Errorf("marshal payload: %w", err)
	}
	payloadEncoded := base64.RawURLEncoding.EncodeToString(payloadBytes)
	// Create signing input
	signingInput := headerEncoded + "." + payloadEncoded
	// Sign with RSA-SHA256
	hash := sha256.Sum256([]byte(signingInput))
	signature, err := rsa.SignPKCS1v15(rand.Reader, k.PrivateKey, crypto.SHA256, hash[:])
	if err != nil {
		return "", fmt.Errorf("sign JWT: %w", err)
	}
	signatureEncoded := base64.RawURLEncoding.EncodeToString(signature)
	// Combine all parts
	token := signingInput + "." + signatureEncoded
	return token, nil
}

// ParseKeyPair parses the embedded PEM data and returns a KeyPair.
func LoadKeyPair() (*KeyPair, error) {
	defer slog.Debug("loaded demo keypair")
	block, _ := pem.Decode(keypairData)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}
	if block.Type != "PRIVATE KEY" {
		return nil, fmt.Errorf("unsupported PEM block type: %s", block.Type)
	}
	parsedKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse PKCS8 private key: %w", err)
	}
	privateKey, ok := parsedKey.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("parsed key is not an RSA private key")
	}
	return &KeyPair{
		PrivateKey: privateKey,
		PublicKey:  &privateKey.PublicKey,
	}, nil
}

func (k *KeyPair) ValidateJWT(token string) (*JWTClaims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid JWT format")
	}
	headerPart, payloadPart, signaturePart := parts[0], parts[1], parts[2]
	signingInput := headerPart + "." + payloadPart
	hash := sha256.Sum256([]byte(signingInput))
	signature, err := base64.RawURLEncoding.DecodeString(signaturePart)
	if err != nil {
		return nil, fmt.Errorf("decode signature: %w", err)
	}
	if err := rsa.VerifyPKCS1v15(k.PublicKey, crypto.SHA256, hash[:], signature); err != nil {
		return nil, fmt.Errorf("verify signature: %w", err)
	}
	payloadBytes, err := base64.RawURLEncoding.DecodeString(payloadPart)
	if err != nil {
		return nil, fmt.Errorf("decode payload: %w", err)
	}
	var claims JWTClaims
	if err := json.Unmarshal(payloadBytes, &claims); err != nil {
		return nil, fmt.Errorf("unmarshal payload: %w", err)
	}
	return &claims, nil
}
