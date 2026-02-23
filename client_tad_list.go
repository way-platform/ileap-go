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

// ListTADsParams is the request parameters for the [Client.ListTADs] method.
type ListTADsParams struct {
	// Limit is the maximum number of TADs to return.
	Limit int `json:"limit,omitempty"`
}

// ListTADsResult is the response for the [Client.ListTADs] method.
type ListTADsResult struct {
	// TADs is the list of transport activity data in the current page.
	TADs []*ileapv1.TAD
}

// ListTADs lists transport activity data.
func (c *Client) ListTADs(
	ctx context.Context,
	request *ListTADsParams,
) (_ *ListTADsResult, err error) {
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
	var response struct {
		Data []json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("unmarshal response body: %w", err)
	}
	tads := make([]*ileapv1.TAD, 0, len(response.Data))
	opts := protojson.UnmarshalOptions{DiscardUnknown: true}
	for _, raw := range response.Data {
		tad := &ileapv1.TAD{}
		if err := opts.Unmarshal(raw, tad); err != nil {
			return nil, fmt.Errorf("unmarshal TAD: %w", err)
		}
		tads = append(tads, tad)
	}
	return &ListTADsResult{
		TADs: tads,
	}, nil
}
