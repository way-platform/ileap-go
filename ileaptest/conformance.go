package ileaptest

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/way-platform/ileap-go"
	"github.com/way-platform/ileap-go/openapi/ileapv1"
	"golang.org/x/oauth2"
)

const (
	schemaShipmentFootprint = "https://api.ileap.sine.dev/shipment-footprint.json"
	schemaTOC               = "https://api.ileap.sine.dev/toc.json"
	schemaHOC               = "https://api.ileap.sine.dev/hoc.json"
)

// ConformanceTestConfig configures the conformance test suite.
type ConformanceTestConfig struct {
	ServerURL string
	Username  string
	Password  string
}

// RunConformanceTests runs the full iLEAP/PACT conformance test suite
// against the server specified in cfg.
func RunConformanceTests(t *testing.T, cfg ConformanceTestConfig) {
	t.Helper()
	token := authenticate(t, cfg.ServerURL, cfg.Username, cfg.Password)

	t.Run("TC001_ShipmentFootprint", func(t *testing.T) {
		body := getJSON(t, cfg.ServerURL, "/2/footprints", token)
		var resp ileapv1.PfListingResponseInner
		if err := json.Unmarshal(body, &resp); err != nil {
			t.Fatalf("decode response: %v", err)
		}
		for _, fp := range resp.Data {
			for _, ext := range fp.Extensions {
				if ext.DataSchema != schemaShipmentFootprint {
					continue
				}
				if ext.SpecVersion != "2.0.0" {
					t.Errorf("extension specVersion: got %q, want 2.0.0", ext.SpecVersion)
				}
				if fp.ProductCategoryCpc != "83117" {
					t.Errorf("productCategoryCpc: got %q, want 83117", fp.ProductCategoryCpc)
				}
				if fp.Pcf.PackagingEmissionsIncluded {
					t.Error("packagingEmissionsIncluded: got true, want false")
				}
				if len(fp.Extensions) != 1 {
					t.Errorf("extensions count: got %d, want exactly 1", len(fp.Extensions))
				}
				var sf map[string]any
				if err := json.Unmarshal(ext.Data, &sf); err != nil {
					t.Fatalf("unmarshal SF data: %v", err)
				}
				if _, ok := sf["mass"]; !ok {
					t.Error("ShipmentFootprint missing required field: mass")
				}
				if _, ok := sf["shipmentId"]; !ok {
					t.Error("ShipmentFootprint missing required field: shipmentId")
				}
				tces, ok := sf["tces"].([]any)
				if !ok || len(tces) == 0 {
					t.Error("ShipmentFootprint must have non-empty tces array")
				}
				for i, tceAny := range tces {
					tce, ok := tceAny.(map[string]any)
					if !ok {
						t.Errorf("tce[%d]: not an object", i)
						continue
					}
					for _, key := range []string{"tceId", "shipmentId", "mass", "co2eWTW", "co2eTTW"} {
						if _, ok := tce[key]; !ok {
							t.Errorf("TCE missing required field: %s", key)
						}
					}
				}
				return
			}
		}
		t.Fatal("no footprint with ShipmentFootprint extension found")
	})

	t.Run("TC002_TOC", func(t *testing.T) {
		body := getJSON(t, cfg.ServerURL, "/2/footprints", token)
		var resp ileapv1.PfListingResponseInner
		if err := json.Unmarshal(body, &resp); err != nil {
			t.Fatalf("decode response: %v", err)
		}
		for _, fp := range resp.Data {
			for _, ext := range fp.Extensions {
				if ext.DataSchema != schemaTOC {
					continue
				}
				if ext.SpecVersion != "2.0.0" {
					t.Errorf("extension specVersion: got %q, want 2.0.0", ext.SpecVersion)
				}
				if fp.ProductCategoryCpc != "83117" {
					t.Errorf("productCategoryCpc: got %q, want 83117", fp.ProductCategoryCpc)
				}
				if fp.Pcf.PackagingEmissionsIncluded {
					t.Error("packagingEmissionsIncluded: got true, want false")
				}
				if len(fp.Extensions) != 1 {
					t.Errorf("extensions count: got %d, want exactly 1", len(fp.Extensions))
				}
				var toc map[string]any
				if err := json.Unmarshal(ext.Data, &toc); err != nil {
					t.Fatalf("unmarshal TOC data: %v", err)
				}
				for _, key := range []string{"tocId", "mode", "co2eIntensityWTW", "co2eIntensityTTW", "transportActivityUnit"} {
					if _, ok := toc[key]; !ok {
						t.Errorf("TOC missing required field: %s", key)
					}
				}
				ec, ok := toc["energyCarriers"].([]any)
				if !ok || len(ec) == 0 {
					t.Error("TOC must have non-empty energyCarriers array")
				}
				_ = ec
				return
			}
		}
		t.Fatal("no footprint with TOC extension found")
	})

	t.Run("TC003_HOC", func(t *testing.T) {
		body := getJSON(t, cfg.ServerURL, "/2/footprints", token)
		var resp ileapv1.PfListingResponseInner
		if err := json.Unmarshal(body, &resp); err != nil {
			t.Fatalf("decode response: %v", err)
		}
		for _, fp := range resp.Data {
			for _, ext := range fp.Extensions {
				if ext.DataSchema != schemaHOC {
					continue
				}
				if ext.SpecVersion != "2.0.0" {
					t.Errorf("extension specVersion: got %q, want 2.0.0", ext.SpecVersion)
				}
				if fp.ProductCategoryCpc != "83117" {
					t.Errorf("productCategoryCpc: got %q, want 83117", fp.ProductCategoryCpc)
				}
				if fp.Pcf.PackagingEmissionsIncluded {
					t.Error("packagingEmissionsIncluded: got true, want false")
				}
				if len(fp.Extensions) != 1 {
					t.Errorf("extensions count: got %d, want exactly 1", len(fp.Extensions))
				}
				var hoc map[string]any
				if err := json.Unmarshal(ext.Data, &hoc); err != nil {
					t.Fatalf("unmarshal HOC data: %v", err)
				}
				for _, key := range []string{"hocId", "hubType", "co2eIntensityWTW", "co2eIntensityTTW", "hubActivityUnit"} {
					if _, ok := hoc[key]; !ok {
						t.Errorf("HOC missing required field: %s", key)
					}
				}
				ec, ok := hoc["energyCarriers"].([]any)
				if !ok || len(ec) == 0 {
					t.Error("HOC must have non-empty energyCarriers array")
				}
				_ = ec
				return
			}
		}
		t.Fatal("no footprint with HOC extension found")
	})

	t.Run("TC004_ListAllTAD", func(t *testing.T) {
		body := getJSON(t, cfg.ServerURL, "/2/ileap/tad", token)
		var resp ileapv1.TadListingResponseInner
		if err := json.Unmarshal(body, &resp); err != nil {
			t.Fatalf("decode response: %v", err)
		}
		if len(resp.Data) == 0 {
			t.Fatal("TAD list is empty")
		}
		for i, tad := range resp.Data {
			if tad.ActivityID == "" {
				t.Errorf("TAD[%d]: missing activityId", i)
			}
			if len(tad.ConsignmentIds) == 0 {
				t.Errorf("TAD[%d]: consignmentIds must be non-empty", i)
			}
			if tad.Origin.Country == "" {
				t.Errorf("TAD[%d]: origin.country missing", i)
			}
			if tad.Destination.Country == "" {
				t.Errorf("TAD[%d]: destination.country missing", i)
			}
			if tad.Mode == "" {
				t.Errorf("TAD[%d]: mode missing", i)
			}
			if tad.DepartureAt.IsZero() {
				t.Errorf("TAD[%d]: departureAt missing", i)
			}
			if tad.ArrivalAt.IsZero() {
				t.Errorf("TAD[%d]: arrivalAt missing", i)
			}
		}
	})

	t.Run("TC005_FilteredTAD", func(t *testing.T) {
		body := getJSON(t, cfg.ServerURL, "/2/ileap/tad?mode=Road", token)
		var resp ileapv1.TadListingResponseInner
		if err := json.Unmarshal(body, &resp); err != nil {
			t.Fatalf("decode response: %v", err)
		}
		for i, tad := range resp.Data {
			if tad.Mode != "Road" {
				t.Errorf("TAD[%d]: mode = %q, want Road", i, tad.Mode)
			}
		}
	})

	t.Run("TC006_LimitedTAD", func(t *testing.T) {
		body := getJSON(t, cfg.ServerURL, "/2/ileap/tad?limit=1", token)
		var resp ileapv1.TadListingResponseInner
		if err := json.Unmarshal(body, &resp); err != nil {
			t.Fatalf("decode response: %v", err)
		}
		if len(resp.Data) > 1 {
			t.Errorf("limit=1: got %d results, want at most 1", len(resp.Data))
		}
	})

	t.Run("TC007_TADInvalidToken", func(t *testing.T) {
		resp := getResponse(t, cfg.ServerURL, "/2/ileap/tad", "invalid-token")
		if resp.StatusCode != http.StatusForbidden {
			t.Errorf("status: got %d, want 403", resp.StatusCode)
		}
		body := readBody(resp)
		var errResp ileap.Error
		if err := json.Unmarshal([]byte(body), &errResp); err != nil {
			t.Fatalf("decode error response: %v", err)
		}
		if errResp.Code != ileap.ErrorCodeAccessDenied {
			t.Errorf("error code: got %q, want AccessDenied", errResp.Code)
		}
	})

	t.Run("PACT_TC01_AuthValidCredentials", func(t *testing.T) {
		resp := postAuthToken(t, cfg.ServerURL, cfg.Username, cfg.Password)
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("status: got %d, want 200", resp.StatusCode)
		}
		var tok struct {
			AccessToken string `json:"access_token"`
			TokenType   string `json:"token_type"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&tok); err != nil {
			t.Fatalf("decode response: %v", err)
		}
		if tok.AccessToken == "" {
			t.Error("access_token is empty")
		}
		if !strings.EqualFold(tok.TokenType, "bearer") {
			t.Errorf("token_type: got %q, want bearer", tok.TokenType)
		}
	})

	t.Run("PACT_TC02_AuthInvalidCredentials", func(t *testing.T) {
		resp := postAuthToken(t, cfg.ServerURL, "wrong-user", "wrong-password")
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("status: got %d, want 400", resp.StatusCode)
		}
		var oauthErr ileap.OAuthError
		if err := json.NewDecoder(resp.Body).Decode(&oauthErr); err != nil {
			t.Fatalf("decode OAuth error: %v", err)
		}
		if oauthErr.Code != ileap.OAuthErrorCodeInvalidRequest {
			t.Errorf("OAuth error code: got %q, want invalid_request", oauthErr.Code)
		}
	})

	t.Run("PACT_TC03_GetFootprint", func(t *testing.T) {
		listBody := getJSON(t, cfg.ServerURL, "/2/footprints", token)
		var listResp ileapv1.PfListingResponseInner
		if err := json.Unmarshal(listBody, &listResp); err != nil {
			t.Fatalf("decode list response: %v", err)
		}
		if len(listResp.Data) == 0 {
			t.Fatal("footprint list is empty")
		}
		fpID := listResp.Data[0].ID
		getBody := getJSON(t, cfg.ServerURL, "/2/footprints/"+fpID, token)
		var getResp ileapv1.ProductFootprintResponse
		if err := json.Unmarshal(getBody, &getResp); err != nil {
			t.Fatalf("decode get response: %v", err)
		}
		if getResp.Data.ID != fpID {
			t.Errorf("footprint ID: got %q, want %q", getResp.Data.ID, fpID)
		}
	})

	t.Run("PACT_TC04_ListFootprints", func(t *testing.T) {
		body := getJSON(t, cfg.ServerURL, "/2/footprints", token)
		var resp ileapv1.PfListingResponseInner
		if err := json.Unmarshal(body, &resp); err != nil {
			t.Fatalf("decode response: %v", err)
		}
		if len(resp.Data) == 0 {
			t.Fatal("footprint list is empty")
		}
		for i, fp := range resp.Data {
			if fp.ID == "" {
				t.Errorf("footprint[%d]: missing id", i)
			}
			if fp.SpecVersion == "" {
				t.Errorf("footprint[%d]: missing specVersion", i)
			}
			if fp.Created.IsZero() {
				t.Errorf("footprint[%d]: missing created", i)
			}
			if fp.Status == "" {
				t.Errorf("footprint[%d]: missing status", i)
			}
			if fp.CompanyName == "" {
				t.Errorf("footprint[%d]: missing companyName", i)
			}
			if len(fp.CompanyIds) == 0 {
				t.Errorf("footprint[%d]: companyIds must be non-empty", i)
			}
			if len(fp.ProductIds) == 0 {
				t.Errorf("footprint[%d]: productIds must be non-empty", i)
			}
		}
	})

	t.Run("PACT_TC05_Pagination", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, cfg.ServerURL+"/2/footprints?limit=1", nil)
		if err != nil {
			t.Fatalf("create request: %v", err)
		}
		req.Header.Set("Authorization", "Bearer "+token)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("request: %v", err)
		}
		defer func() { _ = resp.Body.Close() }()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("status: got %d, want 200", resp.StatusCode)
		}
		var listResp ileapv1.PfListingResponseInner
		if err := json.NewDecoder(resp.Body).Decode(&listResp); err != nil {
			t.Fatalf("decode response: %v", err)
		}
		if len(listResp.Data) > 1 {
			t.Errorf("limit=1: got %d results, want at most 1", len(listResp.Data))
		}

		linkHeader := resp.Header.Get("Link")
		if linkHeader == "" {
			t.Skip("no Link header; server has 1 or fewer footprints")
		}
		if !strings.Contains(linkHeader, `rel="next"`) {
			t.Fatalf("Link header missing rel=next: %s", linkHeader)
		}
		nextURL := strings.TrimRight(strings.TrimLeft(
			strings.Split(linkHeader, ";")[0], "<"), ">")

		nextReq, err := http.NewRequest(http.MethodGet, nextURL, nil)
		if err != nil {
			t.Fatalf("create next request: %v", err)
		}
		nextReq.Header.Set("Authorization", "Bearer "+token)
		nextResp, err := http.DefaultClient.Do(nextReq)
		if err != nil {
			t.Fatalf("next request: %v", err)
		}
		defer func() { _ = nextResp.Body.Close() }()

		if nextResp.StatusCode != http.StatusOK {
			t.Fatalf("next page status: got %d, want 200", nextResp.StatusCode)
		}
	})

	t.Run("PACT_TC06_ListFootprintsInvalidToken", func(t *testing.T) {
		resp := getResponse(t, cfg.ServerURL, "/2/footprints", "invalid-token")
		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("status: got %d, want 401", resp.StatusCode)
		}
		body := readBody(resp)
		var errResp ileap.Error
		if err := json.Unmarshal([]byte(body), &errResp); err != nil {
			t.Fatalf("decode error response: %v", err)
		}
		if errResp.Code != ileap.ErrorCodeAccessDenied {
			t.Errorf("error code: got %q, want AccessDenied", errResp.Code)
		}
	})

	t.Run("PACT_TC07_GetFootprintInvalidToken", func(t *testing.T) {
		resp := getResponse(t, cfg.ServerURL, "/2/footprints/some-id", "invalid-token")
		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("status: got %d, want 401", resp.StatusCode)
		}
		body := readBody(resp)
		var errResp ileap.Error
		if err := json.Unmarshal([]byte(body), &errResp); err != nil {
			t.Fatalf("decode error response: %v", err)
		}
		if errResp.Code != ileap.ErrorCodeAccessDenied {
			t.Errorf("error code: got %q, want AccessDenied", errResp.Code)
		}
	})

	t.Run("PACT_TC08_GetFootprintNotFound", func(t *testing.T) {
		resp := getResponse(t, cfg.ServerURL, "/2/footprints/non-existent-id", token)
		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("status: got %d, want 404", resp.StatusCode)
		}
		body := readBody(resp)
		var errResp ileap.Error
		if err := json.Unmarshal([]byte(body), &errResp); err != nil {
			t.Fatalf("decode error response: %v", err)
		}
		if errResp.Code != ileap.ErrorCodeNoSuchFootprint {
			t.Errorf("error code: got %q, want NoSuchFootprint", errResp.Code)
		}
	})

	t.Run("PACT_TC15_ReceivePublishedEvent", func(t *testing.T) {
		event := `{
		"type": "org.wbcsd.pathfinder.ProductFootprint.Published.v1",
		"specversion": "1.0",
		"id": "test-event-001",
		"source": "//test.example.com",
		"time": "2024-01-15T10:00:00Z",
		"data": {
			"pfIds": ["91715e5e-fd0b-4d1c-8fab-76290c46e6ed"]
		}
	}`
		resp := postEvent(t, cfg.ServerURL, token, event)
		if resp.StatusCode != http.StatusOK {
			body := readBody(resp)
			t.Errorf("status: got %d, want 200: %s", resp.StatusCode, body)
		}
	})

	t.Run("PACT_TC16_EventsInvalidToken", func(t *testing.T) {
		event := `{
		"type": "org.wbcsd.pathfinder.ProductFootprint.Published.v1",
		"specversion": "1.0",
		"id": "test-event-002",
		"source": "//test.example.com",
		"time": "2024-01-15T10:00:00Z",
		"data": {
			"pfIds": ["91715e5e-fd0b-4d1c-8fab-76290c46e6ed"]
		}
	}`
		resp := postEvent(t, cfg.ServerURL, "invalid-token", event)
		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("status: got %d, want 401", resp.StatusCode)
		}
	})

	t.Run("PACT_TC18_OIDCAuthFlow", func(t *testing.T) {
		body := getJSON(t, cfg.ServerURL, "/.well-known/openid-configuration", "")
		var oidc struct {
			Issuer        string `json:"issuer"`
			TokenEndpoint string `json:"token_endpoint"`
			JWKSURI       string `json:"jwks_uri"`
		}
		if err := json.Unmarshal(body, &oidc); err != nil {
			t.Fatalf("decode OIDC config: %v", err)
		}
		if oidc.TokenEndpoint == "" {
			t.Fatal("OIDC config missing token_endpoint")
		}
		if oidc.JWKSURI == "" {
			t.Error("OIDC config missing jwks_uri")
		}

		req, err := http.NewRequest(
			http.MethodPost,
			oidc.TokenEndpoint,
			strings.NewReader("grant_type=client_credentials"),
		)
		if err != nil {
			t.Fatalf("create token request: %v", err)
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.SetBasicAuth(cfg.Username, cfg.Password)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("token request: %v", err)
		}
		defer func() { _ = resp.Body.Close() }()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("OIDC auth: got status %d, want 200", resp.StatusCode)
		}
		var tok struct {
			AccessToken string `json:"access_token"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&tok); err != nil {
			t.Fatalf("decode token response: %v", err)
		}
		if tok.AccessToken == "" {
			t.Error("access_token is empty")
		}
	})

	t.Run("PACT_TC19_OIDCAuthFlowInvalidCredentials", func(t *testing.T) {
		body := getJSON(t, cfg.ServerURL, "/.well-known/openid-configuration", "")
		var oidc struct {
			TokenEndpoint string `json:"token_endpoint"`
		}
		if err := json.Unmarshal(body, &oidc); err != nil {
			t.Fatalf("decode OIDC config: %v", err)
		}
		if oidc.TokenEndpoint == "" {
			t.Fatal("OIDC config missing token_endpoint")
		}

		req, err := http.NewRequest(
			http.MethodPost,
			oidc.TokenEndpoint,
			strings.NewReader("grant_type=client_credentials"),
		)
		if err != nil {
			t.Fatalf("create token request: %v", err)
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.SetBasicAuth("wrong-user", "wrong-password")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("token request: %v", err)
		}
		defer func() { _ = resp.Body.Close() }()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("OIDC auth with invalid creds: got status %d, want 400", resp.StatusCode)
		}
	})

	t.Run("PACT_TC20_FilteredListFootprints", func(t *testing.T) {
		body := getJSON(
			t,
			cfg.ServerURL,
			"/2/footprints?$filter=productCategoryCpc+eq+'83117'",
			token,
		)
		var resp ileapv1.PfListingResponseInner
		if err := json.Unmarshal(body, &resp); err != nil {
			t.Fatalf("decode response: %v", err)
		}
		for i, fp := range resp.Data {
			if fp.ProductCategoryCpc != "83117" {
				t.Errorf(
					"footprint[%d]: productCategoryCpc = %q, want 83117",
					i,
					fp.ProductCategoryCpc,
				)
			}
		}
	})
}

func authenticate(t *testing.T, serverURL, username, password string) string {
	t.Helper()
	req, err := http.NewRequest(
		http.MethodPost,
		serverURL+"/auth/token",
		strings.NewReader("grant_type=client_credentials"),
	)
	if err != nil {
		t.Fatalf("create auth request: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(username, password)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("auth request: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("auth: got status %d, want 200: %s", resp.StatusCode, readBody(resp))
	}
	var tok oauth2.Token
	if err := json.NewDecoder(resp.Body).Decode(&tok); err != nil {
		t.Fatalf("decode auth response: %v", err)
	}
	return tok.AccessToken
}

func getJSON(t *testing.T, serverURL, path, token string) []byte {
	t.Helper()
	req, err := http.NewRequest(http.MethodGet, serverURL+path, nil)
	if err != nil {
		t.Fatalf("create request: %v", err)
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()
	body := readBody(resp)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("GET %s: got status %d, want 200: %s", path, resp.StatusCode, body)
	}
	return []byte(body)
}

func getResponse(t *testing.T, serverURL, path, token string) *http.Response {
	t.Helper()
	req, err := http.NewRequest(http.MethodGet, serverURL+path, nil)
	if err != nil {
		t.Fatalf("create request: %v", err)
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	t.Cleanup(func() { _ = resp.Body.Close() })
	return resp
}

func postEvent(t *testing.T, serverURL, token, body string) *http.Response {
	t.Helper()
	req, err := http.NewRequest(http.MethodPost, serverURL+"/2/events", strings.NewReader(body))
	if err != nil {
		t.Fatalf("create event request: %v", err)
	}
	req.Header.Set("Content-Type", "application/cloudevents+json; charset=UTF-8")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("event request: %v", err)
	}
	t.Cleanup(func() { _ = resp.Body.Close() })
	return resp
}

func postAuthToken(t *testing.T, serverURL, username, password string) *http.Response {
	t.Helper()
	req, err := http.NewRequest(
		http.MethodPost,
		serverURL+"/auth/token",
		strings.NewReader("grant_type=client_credentials"),
	)
	if err != nil {
		t.Fatalf("create auth request: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(username, password)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("auth request: %v", err)
	}
	t.Cleanup(func() { _ = resp.Body.Close() })
	return resp
}

func readBody(resp *http.Response) string {
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(resp.Body)
	return buf.String()
}
