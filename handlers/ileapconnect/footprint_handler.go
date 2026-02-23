package ileapconnect

import (
	"context"

	ileap "github.com/way-platform/ileap-go"
	ileapv1 "github.com/way-platform/ileap-go/proto/gen/wayplatform/connect/ileap/v1"
)

// GetFootprint retrieves a single footprint by ID from the Connect backend.
func (h *Handler) GetFootprint(ctx context.Context, id string) (*ileapv1.ProductFootprint, error) {
	req := new(ileapv1.GetFootprintRequest)
	req.SetId(id)
	resp, err := h.client.GetFootprint(ctx, req)
	if err != nil {
		return nil, mapConnectError(err)
	}
	return resp.GetData(), nil
}

// ListFootprints lists footprints from the Connect backend.
func (h *Handler) ListFootprints(
	ctx context.Context,
	req ileap.ListFootprintsRequest,
) (*ileap.ListFootprintsResponse, error) {
	protoReq := new(ileapv1.ListFootprintsRequest)
	protoReq.SetFilter(req.Filter)
	protoReq.SetLimit(int32(req.Limit))
	protoReq.SetOffset(int32(req.Offset))
	resp, err := h.client.ListFootprints(ctx, protoReq)
	if err != nil {
		return nil, mapConnectError(err)
	}
	return &ileap.ListFootprintsResponse{
		Data:  resp.GetData(),
		Total: int(resp.GetTotal()),
	}, nil
}
