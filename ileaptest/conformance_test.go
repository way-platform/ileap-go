package ileaptest_test

import (
	"os"
	"strings"
	"testing"

	"github.com/way-platform/ileap-go/ileaptest"
)

func TestConformance(t *testing.T) {
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
