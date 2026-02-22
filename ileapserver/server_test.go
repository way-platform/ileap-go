package ileapserver

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/way-platform/ileap-go"
	"github.com/way-platform/ileap-go/openapi/ileapv1"
)

type mockTokenValidator struct {
	valid   bool
	expired bool
}

func (m *mockTokenValidator) ValidateToken(_ context.Context, _ string) (*TokenInfo, error) {
	if m.expired {
		return nil, fmt.Errorf("token expired: %w", ErrTokenExpired)
	}
	if !m.valid {
		return nil, fmt.Errorf("invalid token")
	}
	return &TokenInfo{Subject: "test-user"}, nil
}

type mockFootprintHandler struct {
	footprints []ileapv1.ProductFootprintForILeapType
}

func (m *mockFootprintHandler) GetFootprint(
	_ context.Context, id string,
) (*ileapv1.ProductFootprintForILeapType, error) {
	for _, fp := range m.footprints {
		if fp.ID == id {
			return &fp, nil
		}
	}
	return nil, ErrNotFound
}

func (m *mockFootprintHandler) ListFootprints(
	_ context.Context, req ListFootprintsRequest,
) (*ListFootprintsResponse, error) {
	result := m.footprints
	total := len(result)
	if req.Offset > 0 {
		if req.Offset >= len(result) {
			result = nil
		} else {
			result = result[req.Offset:]
		}
	}
	if req.Limit > 0 && len(result) > req.Limit {
		result = result[:req.Limit]
	}
	return &ListFootprintsResponse{Data: result, Total: total}, nil
}

type mockTADHandler struct {
	tads []ileapv1.TAD
}

func (m *mockTADHandler) ListTADs(
	_ context.Context,
	req ListTADsRequest,
) (*ListTADsResponse, error) {
	result := m.tads
	total := len(result)
	if req.Offset > 0 {
		if req.Offset >= len(result) {
			result = nil
		} else {
			result = result[req.Offset:]
		}
	}
	if req.Limit > 0 && len(result) > req.Limit {
		result = result[:req.Limit]
	}
	return &ListTADsResponse{Data: result, Total: total}, nil
}

type mockEventHandler struct {
	lastEvent *Event
}

func (m *mockEventHandler) HandleEvent(_ context.Context, event Event) error {
	m.lastEvent = &event
	return nil
}

func newTestServer() *Server {
	return NewServer(
		WithTokenValidator(&mockTokenValidator{valid: true}),
		WithFootprintHandler(&mockFootprintHandler{
			footprints: []ileapv1.ProductFootprintForILeapType{
				{ID: "fp-1"},
				{ID: "fp-2"},
			},
		}),
		WithTADHandler(&mockTADHandler{
			tads: []ileapv1.TAD{
				{ActivityID: "tad-1"},
			},
		}),
		WithEventHandler(&mockEventHandler{}),
	)
}

func TestPACTAuthMiddleware(t *testing.T) {
	srv := newTestServer()

	t.Run("missing authorization", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/2/footprints", nil)
		w := httptest.NewRecorder()
		srv.ServeHTTP(w, req)
		checkErrorResponse(t, w, http.StatusBadRequest, ileap.ErrorCodeBadRequest)
	})

	t.Run("non-bearer scheme", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/2/footprints", nil)
		req.Header.Set("Authorization", "Basic abc")
		w := httptest.NewRecorder()
		srv.ServeHTTP(w, req)
		checkErrorResponse(t, w, http.StatusBadRequest, ileap.ErrorCodeBadRequest)
	})

	t.Run("empty bearer token", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/2/footprints", nil)
		req.Header.Set("Authorization", "Bearer ")
		w := httptest.NewRecorder()
		srv.ServeHTTP(w, req)
		checkErrorResponse(t, w, http.StatusBadRequest, ileap.ErrorCodeBadRequest)
	})

	t.Run("invalid token returns 400", func(t *testing.T) {
		srv := NewServer(
			WithTokenValidator(&mockTokenValidator{valid: false}),
			WithFootprintHandler(&mockFootprintHandler{}),
		)
		req := httptest.NewRequest("GET", "/2/footprints", nil)
		req.Header.Set("Authorization", "Bearer bad-token")
		w := httptest.NewRecorder()
		srv.ServeHTTP(w, req)
		checkErrorResponse(t, w, http.StatusBadRequest, ileap.ErrorCodeBadRequest)
	})
}

func TestILeapAuthMiddleware(t *testing.T) {
	srv := newTestServer()

	t.Run("missing authorization returns 403", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/2/ileap/tad", nil)
		w := httptest.NewRecorder()
		srv.ServeHTTP(w, req)
		checkErrorResponse(t, w, http.StatusForbidden, ileap.ErrorCodeAccessDenied)
	})

	t.Run("invalid token returns 403", func(t *testing.T) {
		srv := NewServer(
			WithTokenValidator(&mockTokenValidator{valid: false}),
			WithTADHandler(&mockTADHandler{}),
		)
		req := httptest.NewRequest("GET", "/2/ileap/tad", nil)
		req.Header.Set("Authorization", "Bearer bad-token")
		w := httptest.NewRecorder()
		srv.ServeHTTP(w, req)
		checkErrorResponse(t, w, http.StatusForbidden, ileap.ErrorCodeAccessDenied)
	})

	t.Run("expired token returns 401", func(t *testing.T) {
		srv := NewServer(
			WithTokenValidator(&mockTokenValidator{expired: true}),
			WithTADHandler(&mockTADHandler{}),
		)
		req := httptest.NewRequest("GET", "/2/ileap/tad", nil)
		req.Header.Set("Authorization", "Bearer expired-token")
		w := httptest.NewRecorder()
		srv.ServeHTTP(w, req)
		checkErrorResponse(t, w, http.StatusUnauthorized, ileap.ErrorCodeTokenExpired)
	})
}

func TestListFootprints(t *testing.T) {
	srv := newTestServer()

	t.Run("success", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/2/footprints", nil)
		req.Header.Set("Authorization", "Bearer valid")
		w := httptest.NewRecorder()
		srv.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
		}
		var resp ileapv1.PfListingResponseInner
		if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if len(resp.Data) != 2 {
			t.Errorf("expected 2 footprints, got %d", len(resp.Data))
		}
	})

	t.Run("invalid limit", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/2/footprints?limit=abc", nil)
		req.Header.Set("Authorization", "Bearer valid")
		w := httptest.NewRecorder()
		srv.ServeHTTP(w, req)
		checkErrorResponse(t, w, http.StatusBadRequest, ileap.ErrorCodeBadRequest)
	})
}

func TestListFootprintsPagination(t *testing.T) {
	srv := NewServer(
		WithTokenValidator(&mockTokenValidator{valid: true}),
		WithFootprintHandler(&mockFootprintHandler{
			footprints: []ileapv1.ProductFootprintForILeapType{
				{ID: "fp-1"}, {ID: "fp-2"}, {ID: "fp-3"},
			},
		}),
	)

	t.Run("link header on first page", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/2/footprints?limit=2", nil)
		req.Header.Set("Authorization", "Bearer valid")
		req.Host = "example.com"
		req.Header.Set("X-Forwarded-Proto", "https")
		w := httptest.NewRecorder()
		srv.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
		}
		got := w.Header().Get("Link")
		want := `<https://example.com/2/footprints?limit=2&offset=2>; rel="next"`
		if got != want {
			t.Errorf("Link = %q, want %q", got, want)
		}
	})

	t.Run("no link header on last page", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/2/footprints?limit=2&offset=2", nil)
		req.Header.Set("Authorization", "Bearer valid")
		req.Host = "example.com"
		req.Header.Set("X-Forwarded-Proto", "https")
		w := httptest.NewRecorder()
		srv.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
		}
		if w.Header().Get("Link") != "" {
			t.Errorf("expected no Link header on last page")
		}
	})
}

func TestGetFootprint(t *testing.T) {
	srv := newTestServer()

	t.Run("found", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/2/footprints/fp-1", nil)
		req.Header.Set("Authorization", "Bearer valid")
		w := httptest.NewRecorder()
		srv.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
		}
		var resp ileapv1.ProductFootprintResponse
		if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if resp.Data.ID != "fp-1" {
			t.Errorf("expected fp-1, got %s", resp.Data.ID)
		}
	})

	t.Run("not found", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/2/footprints/nonexistent", nil)
		req.Header.Set("Authorization", "Bearer valid")
		w := httptest.NewRecorder()
		srv.ServeHTTP(w, req)
		checkErrorResponse(t, w, http.StatusNotFound, ileap.ErrorCodeNoSuchFootprint)
	})
}

func TestListTADs(t *testing.T) {
	srv := NewServer(
		WithTokenValidator(&mockTokenValidator{valid: true}),
		WithTADHandler(&mockTADHandler{
			tads: []ileapv1.TAD{
				{ActivityID: "tad-1"},
				{ActivityID: "tad-2"},
				{ActivityID: "tad-3"},
			},
		}),
	)

	t.Run("success", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/2/ileap/tad", nil)
		req.Header.Set("Authorization", "Bearer valid")
		w := httptest.NewRecorder()
		srv.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
		}
		var resp ileapv1.TadListingResponseInner
		if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if len(resp.Data) != 3 {
			t.Errorf("expected 3 TADs, got %d", len(resp.Data))
		}
	})

	t.Run("pagination link header", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/2/ileap/tad?limit=1", nil)
		req.Header.Set("Authorization", "Bearer valid")
		w := httptest.NewRecorder()
		srv.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
		}
		var resp ileapv1.TadListingResponseInner
		if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if len(resp.Data) != 1 {
			t.Errorf("expected 1 TAD, got %d", len(resp.Data))
		}
		link := w.Header().Get("Link")
		if link == "" {
			t.Fatal("expected Link header")
		}
		if !strings.Contains(link, `rel="next"`) {
			t.Errorf("expected rel=next in Link header, got %s", link)
		}
		if !strings.Contains(link, "offset=1") {
			t.Errorf("expected offset=1 in Link header, got %s", link)
		}
	})

	t.Run("no link header when all data returned", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/2/ileap/tad?limit=10", nil)
		req.Header.Set("Authorization", "Bearer valid")
		w := httptest.NewRecorder()
		srv.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
		}
		link := w.Header().Get("Link")
		if link != "" {
			t.Errorf("expected no Link header, got %s", link)
		}
	})

	t.Run("query params passed as filter", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/2/ileap/tad?mode=Road", nil)
		req.Header.Set("Authorization", "Bearer valid")
		w := httptest.NewRecorder()
		srv.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
		}
	})
}

func TestEvents(t *testing.T) {
	eh := &mockEventHandler{}
	srv := NewServer(
		WithTokenValidator(&mockTokenValidator{valid: true}),
		WithEventHandler(eh),
	)

	t.Run("cloudevents content type", func(t *testing.T) {
		body := `{"type":"org.wbcsd.pathfinder.ProductFootprint.Published.v1","specversion":"1.0","id":"evt-1","source":"test","data":{"pfIds":[]}}`
		req := httptest.NewRequest("POST", "/2/events", strings.NewReader(body))
		req.Header.Set("Authorization", "Bearer valid")
		req.Header.Set("Content-Type", "application/cloudevents+json")
		w := httptest.NewRecorder()
		srv.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
		}
		if eh.lastEvent == nil {
			t.Fatal("expected event to be handled")
		}
		if eh.lastEvent.ID != "evt-1" {
			t.Errorf("expected event ID evt-1, got %s", eh.lastEvent.ID)
		}
	})

	t.Run("application/json content type", func(t *testing.T) {
		body := `{"type":"org.wbcsd.pathfinder.ProductFootprint.Published.v1","specversion":"1.0","id":"evt-2","source":"test","data":{"pfIds":[]}}`
		req := httptest.NewRequest("POST", "/2/events", strings.NewReader(body))
		req.Header.Set("Authorization", "Bearer valid")
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		srv.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
		}
	})

	t.Run("missing content type", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/2/events", strings.NewReader("{}"))
		req.Header.Set("Authorization", "Bearer valid")
		w := httptest.NewRecorder()
		srv.ServeHTTP(w, req)
		checkErrorResponse(t, w, http.StatusBadRequest, ileap.ErrorCodeBadRequest)
	})

	t.Run("invalid content type", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/2/events", strings.NewReader("{}"))
		req.Header.Set("Authorization", "Bearer valid")
		req.Header.Set("Content-Type", "text/plain")
		w := httptest.NewRecorder()
		srv.ServeHTTP(w, req)
		checkErrorResponse(t, w, http.StatusBadRequest, ileap.ErrorCodeBadRequest)
	})
}

func TestEventsValidationMissingFields(t *testing.T) {
	srv := newTestServer()
	cases := []struct {
		name string
		body string
	}{
		{
			"missing specversion",
			`{"id":"1","source":"x","type":"org.wbcsd.pathfinder.ProductFootprint.Published.v1"}`,
		},
		{
			"missing id",
			`{"specversion":"1.0","source":"x","type":"org.wbcsd.pathfinder.ProductFootprint.Published.v1"}`,
		},
		{
			"missing source",
			`{"specversion":"1.0","id":"1","type":"org.wbcsd.pathfinder.ProductFootprint.Published.v1"}`,
		},
		{
			"missing data",
			`{"specversion":"1.0","id":"1","source":"x","type":"org.wbcsd.pathfinder.ProductFootprint.Published.v1"}`,
		},
		{
			"null data",
			`{"specversion":"1.0","id":"1","source":"x","type":"org.wbcsd.pathfinder.ProductFootprint.Published.v1","data":null}`,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/2/events", strings.NewReader(tc.body))
			req.Header.Set("Authorization", "Bearer valid")
			req.Header.Set("Content-Type", "application/cloudevents+json")
			w := httptest.NewRecorder()
			srv.ServeHTTP(w, req)
			checkErrorResponse(t, w, http.StatusBadRequest, ileap.ErrorCodeBadRequest)
		})
	}
}

func TestNotImplemented(t *testing.T) {
	srv := NewServer(
		WithTokenValidator(&mockTokenValidator{valid: true}),
	)

	t.Run("footprints", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/2/footprints", nil)
		req.Header.Set("Authorization", "Bearer valid")
		w := httptest.NewRecorder()
		srv.ServeHTTP(w, req)
		checkErrorResponse(t, w, http.StatusNotImplemented, ileap.ErrorCodeNotImplemented)
	})

	t.Run("tads", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/2/ileap/tad", nil)
		req.Header.Set("Authorization", "Bearer valid")
		w := httptest.NewRecorder()
		srv.ServeHTTP(w, req)
		checkErrorResponse(t, w, http.StatusNotImplemented, ileap.ErrorCodeNotImplemented)
	})

	t.Run("events", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/2/events", strings.NewReader("{}"))
		req.Header.Set("Authorization", "Bearer valid")
		req.Header.Set("Content-Type", "application/cloudevents+json")
		w := httptest.NewRecorder()
		srv.ServeHTTP(w, req)
		checkErrorResponse(t, w, http.StatusNotImplemented, ileap.ErrorCodeNotImplemented)
	})
}

func checkErrorResponse(
	t *testing.T,
	w *httptest.ResponseRecorder,
	expectedStatus int,
	expectedCode ileap.ErrorCode,
) {
	t.Helper()
	if w.Code != expectedStatus {
		t.Fatalf("expected status %d, got %d: %s", expectedStatus, w.Code, w.Body.String())
	}
	ct := w.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", ct)
	}
	var errResp ileap.Error
	if err := json.NewDecoder(w.Body).Decode(&errResp); err != nil {
		t.Fatalf("decode error response: %v", err)
	}
	if errResp.Code != expectedCode {
		t.Errorf("expected error code %s, got %s", expectedCode, errResp.Code)
	}
	if errResp.Message == "" {
		t.Error("expected non-empty error message")
	}
}
