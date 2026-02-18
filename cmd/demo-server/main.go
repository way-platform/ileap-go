// Package main runs the demo server.
package main

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"

	"github.com/way-platform/ileap-go/demo"
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
	server, err := demo.NewServer("https://ileap-demo-server-504882905500.europe-north1.run.app")
	if err != nil {
		return err
	}
	address := fmt.Sprintf(":%s", port)
	slog.InfoContext(ctx, "iLEAP demo server listening", "address", address)
	lis, err := (&net.ListenConfig{}).Listen(ctx, "tcp", address)
	if err != nil {
		return err
	}
	if err := (&http.Server{Handler: server.Handler()}).Serve(lis); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}
	return nil
}
