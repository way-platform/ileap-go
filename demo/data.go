// Package demo provides a demo server implementation.
package demo

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/way-platform/ileap-go/openapi/ileapv0"
)

//go:embed data/footprints.json
var footprintsJSON []byte

// LoadFootprints loads example iLEAP footprints.
func LoadFootprints() ([]ileapv0.ProductFootprintForILeapType, error) {
	var data struct {
		Footprints []ileapv0.ProductFootprintForILeapType `json:"footprints"`
	}
	if err := json.Unmarshal(footprintsJSON, &data); err != nil {
		return nil, fmt.Errorf("unmarshal footprints: %w", err)
	}
	for _, footprint := range data.Footprints {
		slog.Debug("loaded demo footprint", "id", footprint.ID)
	}
	slog.Debug("loaded demo footprints", "count", len(data.Footprints))
	return data.Footprints, nil
}

//go:embed data/tad.json
var tadJSON []byte

// LoadTADs loads example iLEAP LoadTADs.
func LoadTADs() ([]ileapv0.TAD, error) {
	var data struct {
		TADs []ileapv0.TAD `json:"tads"`
	}
	if err := json.Unmarshal(tadJSON, &data); err != nil {
		return nil, fmt.Errorf("unmarshal tad: %w", err)
	}
	for _, tad := range data.TADs {
		slog.Debug("loaded demo tad", "id", tad.ActivityID)
	}
	slog.Debug("loaded demo tads", "count", len(data.TADs))
	return data.TADs, nil
}
