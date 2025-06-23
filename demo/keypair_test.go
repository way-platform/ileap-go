package demo

import (
	"reflect"
	"strings"
	"testing"
)

func TestKeyPair(t *testing.T) {
	keypair, err := LoadKeyPair()
	if err != nil {
		t.Fatalf("Failed to load keypair: %v", err)
	}
	t.Run("create and validate JWT", func(t *testing.T) {
		exectedClaims := JWTClaims{Username: "hello", IssuedAt: 1234567890}
		token, err := keypair.CreateJWT(exectedClaims)
		if err != nil {
			t.Fatalf("failed to create JWT: %v", err)
		}
		if token == "" {
			t.Fatal("generated JWT is empty")
		}
		actualClaims, err := keypair.ValidateJWT(token)
		if err != nil {
			t.Fatalf("failed to validate JWT: %v", err)
		}
		if !reflect.DeepEqual(exectedClaims, *actualClaims) {
			t.Fatalf("invalid JWT claims: expected %v, got %v", exectedClaims, *actualClaims)
		}
	})
}

func TestKeyPair_ValidateJWT(t *testing.T) {
	keypair, err := LoadKeyPair()
	if err != nil {
		t.Fatalf("failed to load keypair: %v", err)
	}
	testCases := []struct {
		name          string
		token         string
		expectedError string
	}{
		{
			name:          "empty token",
			token:         "",
			expectedError: "invalid JWT format",
		},
		{
			name:          "too few parts",
			token:         "header.payload",
			expectedError: "invalid JWT format",
		},
		{
			name:          "too many parts",
			token:         "header.payload.signature.extra",
			expectedError: "invalid JWT format",
		},
		{
			name:          "invalid base64 encoding",
			token:         "invalid_base64.payload.signature",
			expectedError: "decode signature: illegal base64 data",
		},
		{
			name:          "tampered payload",
			token:         "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJ1c2VybmFtZSI6InRhbXBlcmVkIn0.eyJ1c2VybmFtZSI6InRhbXBlcmVkIn0",
			expectedError: "verify signature: crypto/rsa: verification error",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := keypair.ValidateJWT(tc.token)
			switch {
			case tc.expectedError == "" && err != nil:
				t.Errorf("expected validation to succeed for %s, but it failed: %v", tc.name, err)
			case tc.expectedError != "" && err == nil:
				t.Errorf("expected validation to fail for %s, but it succeeded", tc.name)
			case tc.expectedError != "" && !strings.Contains(err.Error(), tc.expectedError):
				t.Errorf("expected error to contain %q, got %q", tc.expectedError, err.Error())
			}
		})
	}
}
