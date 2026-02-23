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

func (f *fakeBackend) HandleEvent(
	_ context.Context,
	_ *ileapv1.HandleEventRequest,
) (*ileapv1.HandleEventResponse, error) {
	return &ileapv1.HandleEventResponse{}, nil
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

func newTestFixtures(t *testing.T) (ileapv1connect.ILeapServiceClient, *headerCapture) {
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
	client := NewClient(server.URL)
	return client, capture
}

func TestGetFootprint(t *testing.T) {
	client, _ := newTestFixtures(t)

	t.Run("found", func(t *testing.T) {
		req := new(ileapv1.GetFootprintRequest)
		req.SetId("fp-1")
		resp, err := client.GetFootprint(context.Background(), req)
		if err != nil {
			t.Fatalf("GetFootprint() error: %v", err)
		}
		if resp.GetData().GetId() != "fp-1" {
			t.Errorf("GetFootprint() id = %q, want %q", resp.GetData().GetId(), "fp-1")
		}
	})

	t.Run("not found", func(t *testing.T) {
		req := new(ileapv1.GetFootprintRequest)
		req.SetId("nonexistent")
		_, err := client.GetFootprint(context.Background(), req)
		if err == nil {
			t.Fatal("GetFootprint() expected error for nonexistent ID")
		}
		if connect.CodeOf(err) != connect.CodeNotFound {
			t.Errorf("GetFootprint() error code = %v, want CodeNotFound", connect.CodeOf(err))
		}
	})
}

func TestListFootprints(t *testing.T) {
	client, _ := newTestFixtures(t)

	t.Run("all", func(t *testing.T) {
		resp, err := client.ListFootprints(context.Background(), new(ileapv1.ListFootprintsRequest))
		if err != nil {
			t.Fatalf("ListFootprints() error: %v", err)
		}
		if resp.GetTotal() != 2 {
			t.Errorf("ListFootprints() total = %d, want 2", resp.GetTotal())
		}
		if len(resp.GetData()) != 2 {
			t.Errorf("ListFootprints() len(data) = %d, want 2", len(resp.GetData()))
		}
	})

	t.Run("with limit", func(t *testing.T) {
		req := new(ileapv1.ListFootprintsRequest)
		req.SetLimit(1)
		resp, err := client.ListFootprints(context.Background(), req)
		if err != nil {
			t.Fatalf("ListFootprints() error: %v", err)
		}
		if resp.GetTotal() != 2 {
			t.Errorf("ListFootprints() total = %d, want 2", resp.GetTotal())
		}
		if len(resp.GetData()) != 1 {
			t.Errorf("ListFootprints() len(data) = %d, want 1", len(resp.GetData()))
		}
		if resp.GetData()[0].GetId() != "fp-1" {
			t.Errorf("ListFootprints() data[0].id = %q, want %q", resp.GetData()[0].GetId(), "fp-1")
		}
	})

	t.Run("with offset", func(t *testing.T) {
		req := new(ileapv1.ListFootprintsRequest)
		req.SetOffset(1)
		resp, err := client.ListFootprints(context.Background(), req)
		if err != nil {
			t.Fatalf("ListFootprints() error: %v", err)
		}
		if resp.GetTotal() != 2 {
			t.Errorf("ListFootprints() total = %d, want 2", resp.GetTotal())
		}
		if len(resp.GetData()) != 1 {
			t.Errorf("ListFootprints() len(data) = %d, want 1", len(resp.GetData()))
		}
		if resp.GetData()[0].GetId() != "fp-2" {
			t.Errorf("ListFootprints() data[0].id = %q, want %q", resp.GetData()[0].GetId(), "fp-2")
		}
	})
}

func TestListTransportActivityData(t *testing.T) {
	client, _ := newTestFixtures(t)

	t.Run("all", func(t *testing.T) {
		resp, err := client.ListTransportActivityData(
			context.Background(),
			new(ileapv1.ListTransportActivityDataRequest),
		)
		if err != nil {
			t.Fatalf("ListTransportActivityData() error: %v", err)
		}
		if resp.GetTotal() != 2 {
			t.Errorf("ListTransportActivityData() total = %d, want 2", resp.GetTotal())
		}
		if len(resp.GetData()) != 2 {
			t.Errorf("ListTransportActivityData() len(data) = %d, want 2", len(resp.GetData()))
		}
	})

	t.Run("with limit and offset", func(t *testing.T) {
		req := new(ileapv1.ListTransportActivityDataRequest)
		req.SetLimit(1)
		req.SetOffset(1)
		resp, err := client.ListTransportActivityData(context.Background(), req)
		if err != nil {
			t.Fatalf("ListTransportActivityData() error: %v", err)
		}
		if resp.GetTotal() != 2 {
			t.Errorf("ListTransportActivityData() total = %d, want 2", resp.GetTotal())
		}
		if len(resp.GetData()) != 1 {
			t.Errorf("ListTransportActivityData() len(data) = %d, want 1", len(resp.GetData()))
		}
	})
}

func TestAuthForwarding(t *testing.T) {
	client, capture := newTestFixtures(t)
	ctx := ileap.WithAuthToken(context.Background(), "test-token-123")
	_, err := client.ListFootprints(ctx, new(ileapv1.ListFootprintsRequest))
	if err != nil {
		t.Fatalf("ListFootprints() error: %v", err)
	}
	want := "Bearer test-token-123"
	if got := capture.lastAuthHeader(); got != want {
		t.Errorf("auth header = %q, want %q", got, want)
	}
}

func TestGetFootprintResponse(t *testing.T) {
	client, _ := newTestFixtures(t)
	req := new(ileapv1.GetFootprintRequest)
	req.SetId("fp-1")
	resp, err := client.GetFootprint(context.Background(), req)
	if err != nil {
		t.Fatalf("GetFootprint() error: %v", err)
	}
	want := new(ileapv1.ProductFootprint)
	want.SetId("fp-1")
	want.SetVersion(1)
	if diff := cmp.Diff(want, resp.GetData(), protocmp.Transform()); diff != "" {
		t.Errorf("GetFootprint() mismatch (-want +got):\n%s", diff)
	}
}

func TestHandleEvent(t *testing.T) {
	client, _ := newTestFixtures(t)
	event := new(ileapv1.Event)
	event.SetType("org.wbcsd.pathfinder.ProductFootprint.Published.v1")
	event.SetSpecversion("1.0")
	event.SetId("evt-1")
	event.SetSource("test")
	event.SetData([]byte(`{}`))
	req := new(ileapv1.HandleEventRequest)
	req.SetEvent(event)
	_, err := client.HandleEvent(context.Background(), req)
	if err != nil {
		t.Fatalf("HandleEvent() error: %v", err)
	}
}
