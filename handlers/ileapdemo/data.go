// Package demo provides a demo server implementation.
package ileapdemo

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log/slog"

	"google.golang.org/protobuf/encoding/protojson"

	"github.com/way-platform/ileap-go/ileapv1pb"
)

//go:embed data/footprints.json
var footprintsJSON []byte

// LoadFootprints loads example iLEAP footprints.
func LoadFootprints() ([]*ileapv1pb.ProductFootprint, error) {
	var data struct {
		Footprints []json.RawMessage `json:"footprints"`
	}
	if err := json.Unmarshal(footprintsJSON, &data); err != nil {
		return nil, fmt.Errorf("unmarshal footprints: %w", err)
	}
	result := make([]*ileapv1pb.ProductFootprint, 0, len(data.Footprints))
	opts := protojson.UnmarshalOptions{DiscardUnknown: true}
	for _, raw := range data.Footprints {
		pf := &ileapv1pb.ProductFootprint{}
		if err := opts.Unmarshal(raw, pf); err != nil {
			return nil, fmt.Errorf("unmarshal footprint: %w", err)
		}
		slog.Debug("loaded demo footprint", "id", pf.GetId())
		result = append(result, pf)
	}
	slog.Debug("loaded demo footprints", "count", len(result))
	return result, nil
}

//go:embed data/tad.json
var tadJSON []byte

// LoadTADs loads example iLEAP TADs.
func LoadTADs() ([]*ileapv1pb.TAD, error) {
	var data struct {
		TADs []json.RawMessage `json:"tads"`
	}
	if err := json.Unmarshal(tadJSON, &data); err != nil {
		return nil, fmt.Errorf("unmarshal tad: %w", err)
	}
	result := make([]*ileapv1pb.TAD, 0, len(data.TADs))
	opts := protojson.UnmarshalOptions{DiscardUnknown: true}
	for _, raw := range data.TADs {
		tad := &ileapv1pb.TAD{}
		if err := opts.Unmarshal(raw, tad); err != nil {
			return nil, fmt.Errorf("unmarshal tad: %w", err)
		}
		slog.Debug("loaded demo tad", "id", tad.GetActivityId())
		result = append(result, tad)
	}
	slog.Debug("loaded demo tads", "count", len(result))
	return result, nil
}
