package ileap

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/way-platform/ileap-go/openapi/ileapv0"
)

// GetFootprintRequest is the request for the [Client.GetFootprint] method.
type GetFootprintRequest struct {
	// ID is the ID of the footprint to get.
	ID string
}

// GetFootprint fetches a product carbon footprint by ID.
func (c *Client) GetFootprint(
	ctx context.Context,
	request *GetFootprintRequest,
) (_ *ileapv0.ProductFootprintForILeapType, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("get iLEAP footprint: %w", err)
		}
	}()
	httpRequest, err := c.newRequest(ctx, http.MethodGet, "/2/footprints/"+request.ID, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
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
	var response ileapv0.ProductFootprintResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("unmarshal response body: %w", err)
	}
	return &response.Data, nil
}
