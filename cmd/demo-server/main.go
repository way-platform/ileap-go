package main

import (
	"context"
	_ "embed"
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
	if err := run(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	server, err := demo.NewServer()
	if err != nil {
		return err
	}
	slog.InfoContext(ctx, "iLEAP demo server listening", "address", ":8080")
	lis, err := (&net.ListenConfig{}).Listen(ctx, "tcp", ":8080")
	if err != nil {
		return err
	}
	return (&http.Server{Handler: server.Handler()}).Serve(lis)
}
