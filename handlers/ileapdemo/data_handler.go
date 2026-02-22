package ileapdemo

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/way-platform/ileap-go"
	"github.com/way-platform/ileap-go/openapi/ileapv1"
)

// DataHandler implements ileap.FootprintHandler and ileap.TADHandler
// using embedded demo data.
type DataHandler struct {
	footprints []ileapv1.ProductFootprintForILeapType
	tads       []ileapv1.TAD
}

// NewDataHandler creates a new DataHandler with the embedded demo data.
func NewDataHandler() (*DataHandler, error) {
	footprints, err := LoadFootprints()
	if err != nil {
		return nil, err
	}
	tads, err := LoadTADs()
	if err != nil {
		return nil, err
	}
	return &DataHandler{
		footprints: footprints,
		tads:       tads,
	}, nil
}

// GetFootprint returns a single footprint by ID.
func (h *DataHandler) GetFootprint(
	_ context.Context,
	id string,
) (*ileapv1.ProductFootprintForILeapType, error) {
	for _, fp := range h.footprints {
		if fp.ID == id {
			return &fp, nil
		}
	}
	return nil, ileap.ErrNotFound
}

// ListFootprints returns a filtered, limited list of footprints.
func (h *DataHandler) ListFootprints(
	_ context.Context, req ileap.ListFootprintsRequest,
) (*ileap.ListFootprintsResponse, error) {
	var filter ileap.FilterV2
	if err := filter.UnmarshalString(req.Filter); err != nil {
		return nil, ileap.ErrBadRequest
	}
	filtered := make([]ileapv1.ProductFootprintForILeapType, 0, len(h.footprints))
	for _, fp := range h.footprints {
		if filter.MatchesFootprint(&fp) {
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

// ListTADs returns a filtered, paginated list of transport activity data.
func (h *DataHandler) ListTADs(
	_ context.Context, req ileap.ListTADsRequest,
) (*ileap.ListTADsResponse, error) {
	filtered := h.tads
	if len(req.Filter) > 0 {
		filtered = filterTADs(h.tads, req.Filter)
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
	return &ileap.ListTADsResponse{Data: filtered, Total: total}, nil
}

// filterTADs returns TADs matching all filter criteria.
// Filter matching: serialize TAD to JSON, flatten to key→values,
// then check each filter key/value pair (case-insensitive value match).
func filterTADs(tads []ileapv1.TAD, filters map[string][]string) []ileapv1.TAD {
	result := make([]ileapv1.TAD, 0, len(tads))
	for _, tad := range tads {
		if tadMatchesFilters(tad, filters) {
			result = append(result, tad)
		}
	}
	return result
}

func tadMatchesFilters(tad ileapv1.TAD, filters map[string][]string) bool {
	data, err := json.Marshal(tad)
	if err != nil {
		return false
	}
	var flat map[string]any
	if err := json.Unmarshal(data, &flat); err != nil {
		return false
	}
	flatValues := flattenJSON("", flat)
	for key, wantValues := range filters {
		gotValues, ok := flatValues[key]
		if !ok {
			return false
		}
		for _, want := range wantValues {
			if !containsCaseInsensitive(gotValues, want) {
				return false
			}
		}
	}
	return true
}

// flattenJSON recursively flattens a JSON object to leaf-key → string values.
// For nested objects, the key is the leaf field name (not dot-joined).
// For arrays, each element's values are collected under the element's keys.
func flattenJSON(prefix string, data map[string]any) map[string][]string {
	result := make(map[string][]string)
	for key, value := range data {
		switch v := value.(type) {
		case string:
			result[key] = append(result[key], v)
		case float64:
			// skip numeric fields for string matching
		case bool:
			// skip
		case nil:
			// skip
		case map[string]any:
			for k, vals := range flattenJSON(key, v) {
				result[k] = append(result[k], vals...)
			}
		case []any:
			for _, elem := range v {
				if m, ok := elem.(map[string]any); ok {
					for k, vals := range flattenJSON(key, m) {
						result[k] = append(result[k], vals...)
					}
				} else if s, ok := elem.(string); ok {
					result[key] = append(result[key], s)
				}
			}
		}
	}
	return result
}

func containsCaseInsensitive(haystack []string, needle string) bool {
	for _, s := range haystack {
		if strings.EqualFold(s, needle) {
			return true
		}
	}
	return false
}
