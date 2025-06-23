package demo

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log/slog"

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
