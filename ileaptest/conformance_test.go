package ileaptest_test

import (
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/way-platform/ileap-go"
	"github.com/way-platform/ileap-go/handlers/ileapdemo"
	"github.com/way-platform/ileap-go/ileaptest"
)

func TestConformance(t *testing.T) {
	handler, err := ileapdemo.NewHandler()
	if err != nil {
		t.Fatalf("create demo handler: %v", err)
	}
	auth, err := ileapdemo.NewAuthProvider()
	if err != nil {
		t.Fatalf("create demo auth provider: %v", err)
	}
	server := httptest.NewServer(ileap.NewServer(
		ileap.WithServiceHandler(handler),
		ileap.WithAuthHandler(auth),
	))
	t.Cleanup(server.Close)

	ileaptest.RunConformanceTests(t, ileaptest.ConformanceTestConfig{
		ServerURL: server.URL,
		Username:  "hello",
		Password:  "pathfinder",
	})
}

func TestConformanceRemote(t *testing.T) {
	serverURL := os.Getenv("ILEAP_SERVER_URL")
	username := os.Getenv("ILEAP_USERNAME")
	password := os.Getenv("ILEAP_PASSWORD")
	if serverURL == "" || username == "" || password == "" {
		t.Skip("set ILEAP_SERVER_URL, ILEAP_USERNAME, ILEAP_PASSWORD to run")
	}
	ileaptest.RunConformanceTests(t, ileaptest.ConformanceTestConfig{
		ServerURL: strings.TrimRight(serverURL, "/"),
		Username:  username,
		Password:  password,
	})
}
