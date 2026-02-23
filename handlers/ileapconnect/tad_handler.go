package ileapconnect

import (
	"context"

	ileap "github.com/way-platform/ileap-go"
	ileapv1 "github.com/way-platform/ileap-go/proto/gen/wayplatform/connect/ileap/v1"
)

// ListTADs lists transport activity data from the Connect backend.
func (h *Handler) ListTADs(
	ctx context.Context,
	req ileap.ListTADsRequest,
) (*ileap.ListTADsResponse, error) {
	protoReq := new(ileapv1.ListTransportActivityDataRequest)
	protoReq.SetLimit(int32(req.Limit))
	protoReq.SetOffset(int32(req.Offset))
	if v := firstFilter(req.Filter, "mode"); v != "" {
		protoReq.SetMode(v)
	}
	if v := firstFilter(req.Filter, "feedstock"); v != "" {
		protoReq.SetFeedstock(v)
	}
	if v := firstFilter(req.Filter, "packagingOrTrEqType"); v != "" {
		protoReq.SetPackagingOrTrEqType(v)
	}
	resp, err := h.client.ListTransportActivityData(ctx, protoReq)
	if err != nil {
		return nil, mapConnectError(err)
	}
	return &ileap.ListTADsResponse{
		Data:  resp.GetData(),
		Total: int(resp.GetTotal()),
	}, nil
}

// firstFilter returns the first value for a key in the filter map, or "".
func firstFilter(filters map[string][]string, key string) string {
	if vals := filters[key]; len(vals) > 0 {
		return vals[0]
	}
	return ""
}
