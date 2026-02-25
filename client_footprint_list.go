package ileap

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	ileapv1 "github.com/way-platform/ileap-go/proto/gen/wayplatform/connect/ileap/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

// ListFootprintsParams is the request parameters for the [Client.ListFootprints] method.
type ListFootprintsParams struct {
	// Limit is the maximum number of footprints to return.
	Limit int `json:"limit,omitempty"`
	// Filter is the OData filter to apply to the request.
	Filter string `json:"$filter,omitempty"`
}

// ListFootprints fetches a list of product carbon footprints.
func (c *Client) ListFootprints(
	ctx context.Context,
	request *ListFootprintsParams,
) (_ *ileapv1.ListFootprintsResponse, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("get iLEAP footprint: %w", err)
		}
	}()
	httpRequest, err := c.newRequest(ctx, http.MethodGet, "/2/footprints", nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	query := url.Values{}
	if request.Limit > 0 {
		query.Set("limit", strconv.Itoa(request.Limit))
	}
	if request.Filter != "" {
		query.Set("$filter", request.Filter)
	}
	httpRequest.URL.RawQuery = query.Encode()
	httpResponse, err := c.httpClient.Do(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer func() { _ = httpResponse.Body.Close() }()
	if httpResponse.StatusCode != http.StatusOK {
		return nil, newClientError(httpResponse)
	}
	data, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}
	var response struct {
		Data []json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("unmarshal response body: %w", err)
	}
	footprints := make([]*ileapv1.ProductFootprint, 0, len(response.Data))
	opts := protojson.UnmarshalOptions{DiscardUnknown: true}
	for _, raw := range response.Data {
		pf := &ileapv1.ProductFootprint{}
		if err := opts.Unmarshal(raw, pf); err != nil {
			return nil, fmt.Errorf("unmarshal footprint: %w", err)
		}
		footprints = append(footprints, pf)
	}
	resp := &ileapv1.ListFootprintsResponse{}
	resp.SetData(footprints)
	return resp, nil
}
