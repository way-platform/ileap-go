package ileapconnect

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"connectrpc.com/connect"
	"github.com/google/go-cmp/cmp"
	ileap "github.com/way-platform/ileap-go"
	ileapv1 "github.com/way-platform/ileap-go/proto/gen/wayplatform/connect/ileap/v1"
	"github.com/way-platform/ileap-go/proto/gen/wayplatform/connect/ileap/v1/ileapv1connect"
	"google.golang.org/protobuf/testing/protocmp"
)

// fakeBackend implements ILeapServiceHandler for testing.
type fakeBackend struct {
	ileapv1connect.UnimplementedILeapServiceHandler

	footprints []*ileapv1.ProductFootprint
	tads       []*ileapv1.TAD
}

func (f *fakeBackend) ListFootprints(
	_ context.Context,
	req *ileapv1.ListFootprintsRequest,
) (*ileapv1.ListFootprintsResponse, error) {
	data := f.footprints
	offset := int(req.GetOffset())
	if offset > len(data) {
		offset = len(data)
	}
	data = data[offset:]
	limit := int(req.GetLimit())
	if limit > 0 && limit < len(data) {
		data = data[:limit]
	}
	resp := new(ileapv1.ListFootprintsResponse)
	resp.SetData(data)
	resp.SetTotal(int32(len(f.footprints)))
	return resp, nil
}

func (f *fakeBackend) GetFootprint(
	_ context.Context,
	req *ileapv1.GetFootprintRequest,
) (*ileapv1.GetFootprintResponse, error) {
	for _, fp := range f.footprints {
		if fp.GetId() == req.GetId() {
			resp := new(ileapv1.GetFootprintResponse)
			resp.SetData(fp)
			return resp, nil
		}
	}
	return nil, connect.NewError(connect.CodeNotFound, nil)
}

func (f *fakeBackend) ListTransportActivityData(
	_ context.Context,
	req *ileapv1.ListTransportActivityDataRequest,
) (*ileapv1.ListTransportActivityDataResponse, error) {
	data := f.tads
	offset := int(req.GetOffset())
	if offset > len(data) {
		offset = len(data)
	}
	data = data[offset:]
	limit := int(req.GetLimit())
	if limit > 0 && limit < len(data) {
		data = data[:limit]
	}
	resp := new(ileapv1.ListTransportActivityDataResponse)
	resp.SetData(data)
	resp.SetTotal(int32(len(f.tads)))
	return resp, nil
}

// headerCapture is an HTTP middleware that records the last Authorization header.
type headerCapture struct {
	mu      sync.Mutex
	handler http.Handler
	last    string
}

func (h *headerCapture) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mu.Lock()
	h.last = r.Header.Get("Authorization")
	h.mu.Unlock()
	h.handler.ServeHTTP(w, r)
}

func (h *headerCapture) lastAuthHeader() string {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.last
}

func newTestFixtures(t *testing.T) (*Handler, *headerCapture) {
	t.Helper()
	fp1 := new(ileapv1.ProductFootprint)
	fp1.SetId("fp-1")
	fp1.SetVersion(1)
	fp2 := new(ileapv1.ProductFootprint)
	fp2.SetId("fp-2")
	fp2.SetVersion(1)
	tad1 := new(ileapv1.TAD)
	tad1.SetActivityId("tad-1")
	tad2 := new(ileapv1.TAD)
	tad2.SetActivityId("tad-2")
	backend := &fakeBackend{
		footprints: []*ileapv1.ProductFootprint{fp1, fp2},
		tads:       []*ileapv1.TAD{tad1, tad2},
	}
	path, connectHandler := ileapv1connect.NewILeapServiceHandler(backend)
	mux := http.NewServeMux()
	mux.Handle(path, connectHandler)
	capture := &headerCapture{handler: mux}
	server := httptest.NewServer(capture)
	t.Cleanup(server.Close)
	h := New(server.URL)
	return h, capture
}

func TestGetFootprint(t *testing.T) {
	h, _ := newTestFixtures(t)

	t.Run("found", func(t *testing.T) {
		fp, err := h.GetFootprint(context.Background(), "fp-1")
		if err != nil {
			t.Fatalf("GetFootprint() error: %v", err)
		}
		if fp.GetId() != "fp-1" {
			t.Errorf("GetFootprint() id = %q, want %q", fp.GetId(), "fp-1")
		}
	})

	t.Run("not found", func(t *testing.T) {
		_, err := h.GetFootprint(context.Background(), "nonexistent")
		if err == nil {
			t.Fatal("GetFootprint() expected error for nonexistent ID")
		}
		if err != ileap.ErrNotFound {
			t.Errorf("GetFootprint() error = %v, want ErrNotFound", err)
		}
	})
}

func TestListFootprints(t *testing.T) {
	h, _ := newTestFixtures(t)

	t.Run("all", func(t *testing.T) {
		resp, err := h.ListFootprints(context.Background(), ileap.ListFootprintsRequest{})
		if err != nil {
			t.Fatalf("ListFootprints() error: %v", err)
		}
		if resp.Total != 2 {
			t.Errorf("ListFootprints() total = %d, want 2", resp.Total)
		}
		if len(resp.Data) != 2 {
			t.Errorf("ListFootprints() len(data) = %d, want 2", len(resp.Data))
		}
	})

	t.Run("with limit", func(t *testing.T) {
		resp, err := h.ListFootprints(context.Background(), ileap.ListFootprintsRequest{
			Limit: 1,
		})
		if err != nil {
			t.Fatalf("ListFootprints() error: %v", err)
		}
		if resp.Total != 2 {
			t.Errorf("ListFootprints() total = %d, want 2", resp.Total)
		}
		if len(resp.Data) != 1 {
			t.Errorf("ListFootprints() len(data) = %d, want 1", len(resp.Data))
		}
		if resp.Data[0].GetId() != "fp-1" {
			t.Errorf("ListFootprints() data[0].id = %q, want %q", resp.Data[0].GetId(), "fp-1")
		}
	})

	t.Run("with offset", func(t *testing.T) {
		resp, err := h.ListFootprints(context.Background(), ileap.ListFootprintsRequest{
			Offset: 1,
		})
		if err != nil {
			t.Fatalf("ListFootprints() error: %v", err)
		}
		if resp.Total != 2 {
			t.Errorf("ListFootprints() total = %d, want 2", resp.Total)
		}
		if len(resp.Data) != 1 {
			t.Errorf("ListFootprints() len(data) = %d, want 1", len(resp.Data))
		}
		if resp.Data[0].GetId() != "fp-2" {
			t.Errorf("ListFootprints() data[0].id = %q, want %q", resp.Data[0].GetId(), "fp-2")
		}
	})
}

func TestListTADs(t *testing.T) {
	h, _ := newTestFixtures(t)

	t.Run("all", func(t *testing.T) {
		resp, err := h.ListTADs(context.Background(), ileap.ListTADsRequest{})
		if err != nil {
			t.Fatalf("ListTADs() error: %v", err)
		}
		if resp.Total != 2 {
			t.Errorf("ListTADs() total = %d, want 2", resp.Total)
		}
		if len(resp.Data) != 2 {
			t.Errorf("ListTADs() len(data) = %d, want 2", len(resp.Data))
		}
	})

	t.Run("with limit and offset", func(t *testing.T) {
		resp, err := h.ListTADs(context.Background(), ileap.ListTADsRequest{
			Limit:  1,
			Offset: 1,
		})
		if err != nil {
			t.Fatalf("ListTADs() error: %v", err)
		}
		if resp.Total != 2 {
			t.Errorf("ListTADs() total = %d, want 2", resp.Total)
		}
		if len(resp.Data) != 1 {
			t.Errorf("ListTADs() len(data) = %d, want 1", len(resp.Data))
		}
	})
}

func TestAuthForwarding(t *testing.T) {
	h, capture := newTestFixtures(t)
	ctx := ileap.WithAuthToken(context.Background(), "test-token-123")
	_, err := h.ListFootprints(ctx, ileap.ListFootprintsRequest{})
	if err != nil {
		t.Fatalf("ListFootprints() error: %v", err)
	}
	want := "Bearer test-token-123"
	if got := capture.lastAuthHeader(); got != want {
		t.Errorf("auth header = %q, want %q", got, want)
	}
}

func TestErrorMapping(t *testing.T) {
	tests := []struct {
		name    string
		code    connect.Code
		wantErr error
	}{
		{name: "NotFound", code: connect.CodeNotFound, wantErr: ileap.ErrNotFound},
		{name: "InvalidArgument", code: connect.CodeInvalidArgument, wantErr: ileap.ErrBadRequest},
		{
			name:    "Unauthenticated",
			code:    connect.CodeUnauthenticated,
			wantErr: ileap.ErrTokenExpired,
		},
		{
			name:    "PermissionDenied",
			code:    connect.CodePermissionDenied,
			wantErr: ileap.ErrInvalidCredentials,
		},
		{name: "Unimplemented", code: connect.CodeUnimplemented, wantErr: ileap.ErrNotImplemented},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mapConnectError(connect.NewError(tt.code, nil))
			if got != tt.wantErr {
				t.Errorf("mapConnectError(Code%s) = %v, want %v", tt.name, got, tt.wantErr)
			}
		})
	}
}

func TestGetFootprintResponse(t *testing.T) {
	h, _ := newTestFixtures(t)
	fp, err := h.GetFootprint(context.Background(), "fp-1")
	if err != nil {
		t.Fatalf("GetFootprint() error: %v", err)
	}
	want := new(ileapv1.ProductFootprint)
	want.SetId("fp-1")
	want.SetVersion(1)
	if diff := cmp.Diff(want, fp, protocmp.Transform()); diff != "" {
		t.Errorf("GetFootprint() mismatch (-want +got):\n%s", diff)
	}
}
