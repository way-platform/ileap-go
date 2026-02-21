// Package main runs the demo server.
package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"

	"github.com/way-platform/ileap-go/clerk"
	"github.com/way-platform/ileap-go/demo"
	"github.com/way-platform/ileap-go/ileapauthserver"
	"github.com/way-platform/ileap-go/ileapserver"
)

func main() {
	ctx := context.Background()
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})))
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := run(ctx, port); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(ctx context.Context, port string) error {
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "https://ileap-demo-server-504882905500.europe-north1.run.app"
	}
	handler, err := buildHandler(baseURL)
	if err != nil {
		return err
	}
	address := fmt.Sprintf(":%s", port)
	slog.InfoContext(ctx, "iLEAP demo server listening", "address", address)
	lis, err := (&net.ListenConfig{}).Listen(ctx, "tcp", address)
	if err != nil {
		return err
	}
	if err := (&http.Server{Handler: handler}).Serve(lis); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}
	return nil
}

func buildHandler(baseURL string) (http.Handler, error) {
	authBackend := os.Getenv("AUTH_BACKEND")
	if authBackend == "" {
		authBackend = "demo"
	}
	switch authBackend {
	case "demo":
		server, err := demo.NewServer(baseURL)
		if err != nil {
			return nil, err
		}
		return server.Handler(), nil
	case "clerk":
		return buildClerkHandler(baseURL)
	default:
		return nil, fmt.Errorf("unknown AUTH_BACKEND: %s", authBackend)
	}
}

func buildClerkHandler(baseURL string) (http.Handler, error) {
	fapiDomain := os.Getenv("CLERK_FAPI_DOMAIN")
	if fapiDomain == "" {
		return nil, fmt.Errorf("CLERK_FAPI_DOMAIN required when AUTH_BACKEND=clerk")
	}
	keypair, err := demo.LoadKeyPair()
	if err != nil {
		return nil, fmt.Errorf("load keypair: %w", err)
	}
	dataHandler, err := demo.NewDataHandler()
	if err != nil {
		return nil, fmt.Errorf("create data handler: %w", err)
	}
	clerkClient := clerk.NewClient(fapiDomain)
	tokenIssuer := clerk.NewTokenIssuer(clerkClient, keypair)
	oidcProvider := clerk.NewOIDCProvider(keypair)
	// Reuse demo auth provider for JWT validation (same keypair).
	tokenValidator, err := demo.NewAuthProvider()
	if err != nil {
		return nil, fmt.Errorf("create auth provider: %w", err)
	}
	dataSrv := ileapserver.NewServer(
		ileapserver.WithFootprintHandler(dataHandler),
		ileapserver.WithTADHandler(dataHandler),
		ileapserver.WithEventHandler(&demo.EventHandler{}),
		ileapserver.WithTokenValidator(tokenValidator),
	)
	authSrv := ileapauthserver.NewServer(baseURL, tokenIssuer, oidcProvider)
	mux := http.NewServeMux()
	mux.Handle("/auth/", authSrv)
	mux.Handle("/.well-known/", authSrv)
	mux.Handle("/jwks", authSrv)
	mux.Handle("/", dataSrv)
	return mux, nil
}
