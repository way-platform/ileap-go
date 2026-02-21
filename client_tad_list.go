package ileap

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/way-platform/ileap-go/openapi/ileapv1"
)

// ListTADsRequest is the request for the [Client.ListTADs] method.
type ListTADsRequest struct {
	// Limit is the maximum number of TADs to return.
	Limit int `json:"limit,omitempty"`
}

// ListTADsResponse is the response for the [Client.ListTADs] method.
type ListTADsResponse struct {
	// TADs is the list of transport activity data in the current page.
	TADs []ileapv1.TAD `json:"tads"`
}

// ListTADs lists transport activity data.
func (c *Client) ListTADs(
	ctx context.Context,
	request *ListTADsRequest,
) (_ *ListTADsResponse, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("list iLEAP TADs: %w", err)
		}
	}()
	httpRequest, err := c.newRequest(ctx, http.MethodGet, "/2/ileap/tad", nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	query := url.Values{}
	if request.Limit > 0 {
		query.Set("limit", strconv.Itoa(request.Limit))
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
	var response ileapv1.TadListingResponseInner
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("unmarshal response body: %w", err)
	}
	return &ListTADsResponse{
		TADs: response.Data,
	}, nil
}
