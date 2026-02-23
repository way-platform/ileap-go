package ileapdemo

import (
	"context"
	"fmt"
	"strings"

	"connectrpc.com/connect"
	"github.com/way-platform/ileap-go"
	ileapv1 "github.com/way-platform/ileap-go/proto/gen/wayplatform/connect/ileap/v1"
	"github.com/way-platform/ileap-go/proto/gen/wayplatform/connect/ileap/v1/ileapv1connect"
)

var _ ileapv1connect.ILeapServiceHandler = (*Handler)(nil)

// Handler implements ILeapServiceHandler using embedded demo data.
type Handler struct {
	ileapv1connect.UnimplementedILeapServiceHandler
	footprints []*ileapv1.ProductFootprint
	tads       []*ileapv1.TAD
}

// NewHandler creates a new Handler with the embedded demo data.
func NewHandler() (*Handler, error) {
	footprints, err := LoadFootprints()
	if err != nil {
		return nil, err
	}
	tads, err := LoadTADs()
	if err != nil {
		return nil, err
	}
	return &Handler{
		footprints: footprints,
		tads:       tads,
	}, nil
}

// GetFootprint returns a single footprint by ID.
func (h *Handler) GetFootprint(
	_ context.Context, req *ileapv1.GetFootprintRequest,
) (*ileapv1.GetFootprintResponse, error) {
	for _, fp := range h.footprints {
		if fp.GetId() == req.GetId() {
			resp := new(ileapv1.GetFootprintResponse)
			resp.SetData(fp)
			return resp, nil
		}
	}
	return nil, connect.NewError(connect.CodeNotFound, nil)
}

// ListFootprints returns a filtered, limited list of footprints.
func (h *Handler) ListFootprints(
	_ context.Context, req *ileapv1.ListFootprintsRequest,
) (*ileapv1.ListFootprintsResponse, error) {
	var filter ileap.FilterV2
	if err := filter.UnmarshalString(req.GetFilter()); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, nil)
	}
	filtered := make([]*ileapv1.ProductFootprint, 0, len(h.footprints))
	for _, fp := range h.footprints {
		if filter.MatchesFootprint(fp) {
			filtered = append(filtered, fp)
		}
	}
	total := len(filtered)
	offset := int(req.GetOffset())
	limit := int(req.GetLimit())
	if offset > 0 {
		if offset >= len(filtered) {
			filtered = nil
		} else {
			filtered = filtered[offset:]
		}
	}
	if limit > 0 && len(filtered) > limit {
		filtered = filtered[:limit]
	}
	resp := new(ileapv1.ListFootprintsResponse)
	resp.SetData(filtered)
	resp.SetTotal(int32(total))
	return resp, nil
}

// ListTransportActivityData returns a filtered, paginated list of transport activity data.
func (h *Handler) ListTransportActivityData(
	_ context.Context, req *ileapv1.ListTransportActivityDataRequest,
) (*ileapv1.ListTransportActivityDataResponse, error) {
	filtered := h.tads
	if hasFilter(req) {
		filtered = filterTADs(h.tads, req)
	}
	total := len(filtered)
	offset := int(req.GetOffset())
	limit := int(req.GetLimit())
	if offset > 0 {
		if offset >= len(filtered) {
			filtered = nil
		} else {
			filtered = filtered[offset:]
		}
	}
	if limit > 0 && len(filtered) > limit {
		filtered = filtered[:limit]
	}
	resp := new(ileapv1.ListTransportActivityDataResponse)
	resp.SetData(filtered)
	resp.SetTotal(int32(total))
	return resp, nil
}

// HandleEvent processes an incoming PACT event.
func (h *Handler) HandleEvent(
	_ context.Context, req *ileapv1.HandleEventRequest,
) (*ileapv1.HandleEventResponse, error) {
	event := req.GetEvent()
	switch ileap.EventType(event.GetType()) {
	case ileap.EventTypeRequestCreatedV1:
	case ileap.EventTypeRequestFulfilledV1:
	case ileap.EventTypeRequestRejectedV1:
	case ileap.EventTypePublishedV1:
	default:
		return nil, connect.NewError(
			connect.CodeInvalidArgument,
			fmt.Errorf("invalid event type: %s", event.GetType()),
		)
	}
	return &ileapv1.HandleEventResponse{}, nil
}

func hasFilter(req *ileapv1.ListTransportActivityDataRequest) bool {
	return req.GetMode() != "" || req.GetFeedstock() != "" || req.GetPackagingOrTrEqType() != ""
}

func filterTADs(tads []*ileapv1.TAD, req *ileapv1.ListTransportActivityDataRequest) []*ileapv1.TAD {
	result := make([]*ileapv1.TAD, 0, len(tads))
	for _, tad := range tads {
		if tadMatchesFilter(tad, req) {
			result = append(result, tad)
		}
	}
	return result
}

func tadMatchesFilter(tad *ileapv1.TAD, req *ileapv1.ListTransportActivityDataRequest) bool {
	if mode := req.GetMode(); mode != "" {
		if !strings.EqualFold(tad.GetMode(), mode) {
			return false
		}
	}
	if packagingOrTrEqType := req.GetPackagingOrTrEqType(); packagingOrTrEqType != "" {
		if !strings.EqualFold(tad.GetPackagingOrTrEqType(), packagingOrTrEqType) {
			return false
		}
	}
	if feedstock := req.GetFeedstock(); feedstock != "" {
		if !tadHasFeedstock(tad, feedstock) {
			return false
		}
	}
	return true
}

func tadHasFeedstock(tad *ileapv1.TAD, feedstock string) bool {
	for _, ec := range tad.GetEnergyCarriers() {
		for _, fs := range ec.GetFeedstocks() {
			if strings.EqualFold(fs.GetFeedstock(), feedstock) {
				return true
			}
		}
	}
	return false
}
