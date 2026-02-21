// Package demoserver provides the demo-server subcommand.
package demoserver

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/way-platform/ileap-go/clerk"
	"github.com/way-platform/ileap-go/demo"
	"github.com/way-platform/ileap-go/ileapauthserver"
	"github.com/way-platform/ileap-go/ileapserver"
)

// NewCommand returns the demo-server cobra command.
func NewCommand() *cobra.Command {
	v := viper.New()
	v.SetEnvPrefix("ILEAP")
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.BindEnv("port", "PORT") //nolint:errcheck
	v.AutomaticEnv()
	cmd := &cobra.Command{
		Use:   "demo-server",
		Short: "Start the iLEAP demo server",
	}
	cmd.Flags().Int("port", 8080, "port to listen on")
	cmd.Flags().String(
		"base-url",
		"https://ileap-demo-server-504882905500.europe-north1.run.app",
		"base URL of the server",
	)
	cmd.Flags().String("auth-backend", "demo", "auth backend to use (demo, clerk)")
	cmd.Flags().String("clerk-fapi-domain", "", "Clerk FAPI domain (required when auth-backend=clerk)")
	v.BindPFlag("port", cmd.Flags().Lookup("port"))                         //nolint:errcheck
	v.BindPFlag("base-url", cmd.Flags().Lookup("base-url"))                 //nolint:errcheck
	v.BindPFlag("auth-backend", cmd.Flags().Lookup("auth-backend"))         //nolint:errcheck
	v.BindPFlag("clerk-fapi-domain", cmd.Flags().Lookup("clerk-fapi-domain")) //nolint:errcheck
	cmd.RunE = func(cmd *cobra.Command, _ []string) error {
		slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})))
		return run(cmd.Context(), v)
	}
	return cmd
}

func run(ctx context.Context, v *viper.Viper) error {
	port := v.GetInt("port")
	baseURL := v.GetString("base-url")
	handler, err := buildHandler(v, baseURL)
	if err != nil {
		return err
	}
	address := fmt.Sprintf(":%d", port)
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

func buildHandler(v *viper.Viper, baseURL string) (http.Handler, error) {
	authBackend := v.GetString("auth-backend")
	switch authBackend {
	case "demo":
		server, err := demo.NewServer(baseURL)
		if err != nil {
			return nil, err
		}
		return server.Handler(), nil
	case "clerk":
		return buildClerkHandler(v, baseURL)
	default:
		return nil, fmt.Errorf("unknown auth-backend: %s", authBackend)
	}
}

func buildClerkHandler(v *viper.Viper, baseURL string) (http.Handler, error) {
	fapiDomain := v.GetString("clerk-fapi-domain")
	if fapiDomain == "" {
		return nil, fmt.Errorf("--clerk-fapi-domain required when --auth-backend=clerk")
	}
	dataHandler, err := demo.NewDataHandler()
	if err != nil {
		return nil, fmt.Errorf("create data handler: %w", err)
	}
	clerkClient := clerk.NewClient(fapiDomain)
	tokenIssuer := clerk.NewTokenIssuer(clerkClient)
	oidcProvider := clerk.NewOIDCProvider(clerkClient)
	tokenValidator := clerk.NewTokenValidator(clerkClient)
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
