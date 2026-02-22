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
	"github.com/way-platform/ileap-go"
	"github.com/way-platform/ileap-go/handlers/ileapclerk"
	"github.com/way-platform/ileap-go/handlers/ileapdemo"
)

// NewCommand returns the demo-server cobra command.
func NewCommand() *cobra.Command {
	v := viper.New()
	v.SetEnvPrefix("ILEAP")
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	_ = v.BindEnv("port", "PORT")
	v.AutomaticEnv()
	cmd := &cobra.Command{
		Use:   "demo-server",
		Short: "Start the iLEAP demo server",
	}
	cmd.Flags().Int("port", 8080, "port to listen on")
	cmd.Flags().String("auth-backend", "demo", "auth backend to use (demo, clerk)")
	cmd.Flags().
		String("clerk-fapi-domain", "", "Clerk FAPI domain (required when auth-backend=clerk)")
	cmd.Flags().
		String("clerk-organization-id", "", "Clerk organization ID to activate upon login (optional)")
	_ = v.BindPFlag(
		"port",
		cmd.Flags().Lookup("port"),
	)
	_ = v.BindPFlag(
		"auth-backend",
		cmd.Flags().Lookup("auth-backend"),
	)
	_ = v.BindPFlag(
		"clerk-fapi-domain",
		cmd.Flags().Lookup("clerk-fapi-domain"),
	)
	_ = v.BindPFlag(
		"clerk-organization-id",
		cmd.Flags().Lookup("clerk-organization-id"),
	)
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
	handler, err := buildHandler(v)
	if err != nil {
		return err
	}
	address := fmt.Sprintf(":%d", port)
	slog.InfoContext(ctx, "iLEAP demo server listening", "address", address)
	lis, err := (&net.ListenConfig{}).Listen(ctx, "tcp", address)
	if err != nil {
		return err
	}
	if err := (&http.Server{Handler: logRequests(handler)}).Serve(lis); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}
	return nil
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r)
		slog.InfoContext(r.Context(), "request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", rec.status,
		)
	})
}

func buildHandler(v *viper.Viper) (http.Handler, error) {
	authBackend := v.GetString("auth-backend")
	slog.Info("starting demo server", "auth-backend", authBackend)
	switch authBackend {
	case "demo":
		authProvider, err := ileapdemo.NewAuthProvider()
		if err != nil {
			return nil, err
		}
		dataHandler, err := ileapdemo.NewDataHandler()
		if err != nil {
			return nil, err
		}
		return ileap.NewServer(
			ileap.WithFootprintHandler(dataHandler),
			ileap.WithTADHandler(dataHandler),
			ileap.WithEventHandler(&ileapdemo.EventHandler{}),
			ileap.WithTokenValidator(authProvider),
			ileap.WithTokenIssuer(authProvider),
			ileap.WithOIDCProvider(authProvider),
		), nil
	case "clerk":
		return buildClerkHandler(v)
	default:
		return nil, fmt.Errorf("unknown auth-backend: %s", authBackend)
	}
}

func buildClerkHandler(v *viper.Viper) (http.Handler, error) {
	fapiDomain := v.GetString("clerk-fapi-domain")
	if fapiDomain == "" {
		return nil, fmt.Errorf("--clerk-fapi-domain required when --auth-backend=clerk")
	}
	slog.Info("using clerk auth backend", "fapi-domain", fapiDomain)

	activeOrgID := v.GetString("clerk-organization-id")

	dataHandler, err := ileapdemo.NewDataHandler()
	if err != nil {
		return nil, fmt.Errorf("create data handler: %w", err)
	}

	clerkClient := ileapclerk.NewClient(fapiDomain)
	tokenIssuer := ileapclerk.NewTokenIssuer(
		clerkClient,
		ileapclerk.WithActiveOrganization(activeOrgID),
	)
	oidcProvider := ileapclerk.NewOIDCProvider(clerkClient)
	tokenValidator := ileapclerk.NewTokenValidator(clerkClient)
	srv := ileap.NewServer(
		ileap.WithFootprintHandler(dataHandler),
		ileap.WithTADHandler(dataHandler),
		ileap.WithEventHandler(&ileapdemo.EventHandler{}),
		ileap.WithTokenValidator(tokenValidator),
		ileap.WithTokenIssuer(tokenIssuer),
		ileap.WithOIDCProvider(oidcProvider),
	)
	return srv, nil
}
