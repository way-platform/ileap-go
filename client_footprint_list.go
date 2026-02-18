package ileap

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/way-platform/ileap-go/openapi/ileapv0"
)

// ListFootprintsRequest is the request for the [Client.ListFootprints] method.
type ListFootprintsRequest struct {
	// Limit is the maximum number of footprints to return.
	Limit int `json:"limit,omitempty"`
	// Filter is the OData filter to apply to the request.
	Filter string `json:"$filter,omitempty"`
}

// ListFootprintsResponse is the response for the [Client.ListFootprints] method.
type ListFootprintsResponse struct {
	// Footprints is the list of footprints in the current page.
	Footprints []ileapv0.ProductFootprintForILeapType `json:"footprints"`
}

// ListFootprints fetches a list of product carbon footprints.
func (c *Client) ListFootprints(
	ctx context.Context,
	request *ListFootprintsRequest,
) (_ *ListFootprintsResponse, err error) {
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
	// TODO: Parse next page link.
	var response ileapv0.PfListingResponseInner
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("unmarshal response body: %w", err)
	}
	return &ListFootprintsResponse{
		Footprints: response.Data,
	}, nil
}
