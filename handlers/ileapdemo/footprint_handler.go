package ileapdemo

import (
	"context"

	"github.com/way-platform/ileap-go"
	"github.com/way-platform/ileap-go/ileapv1pb"
)

// FootprintHandler implements ileap.FootprintHandler using embedded demo data.
type FootprintHandler struct {
	footprints []*ileapv1pb.ProductFootprint
}

// NewFootprintHandler creates a new FootprintHandler with the embedded demo data.
func NewFootprintHandler() (*FootprintHandler, error) {
	footprints, err := LoadFootprints()
	if err != nil {
		return nil, err
	}
	return &FootprintHandler{
		footprints: footprints,
	}, nil
}

// GetFootprint returns a single footprint by ID.
func (h *FootprintHandler) GetFootprint(
	_ context.Context,
	id string,
) (*ileapv1pb.ProductFootprint, error) {
	for _, fp := range h.footprints {
		if fp.GetId() == id {
			return fp, nil
		}
	}
	return nil, ileap.ErrNotFound
}

// ListFootprints returns a filtered, limited list of footprints.
func (h *FootprintHandler) ListFootprints(
	_ context.Context, req ileap.ListFootprintsRequest,
) (*ileap.ListFootprintsResponse, error) {
	var filter ileap.FilterV2
	if err := filter.UnmarshalString(req.Filter); err != nil {
		return nil, ileap.ErrBadRequest
	}
	filtered := make([]*ileapv1pb.ProductFootprint, 0, len(h.footprints))
	for _, fp := range h.footprints {
		if filter.MatchesFootprint(fp) {
			filtered = append(filtered, fp)
		}
	}
	total := len(filtered)
	// Apply offset.
	if req.Offset > 0 {
		if req.Offset >= len(filtered) {
			filtered = nil
		} else {
			filtered = filtered[req.Offset:]
		}
	}
	// Apply limit.
	if req.Limit > 0 && len(filtered) > req.Limit {
		filtered = filtered[:req.Limit]
	}
	return &ileap.ListFootprintsResponse{Data: filtered, Total: total}, nil
}
