package ileapdemo

import (
	"context"
	"slices"
	"strings"

	"connectrpc.com/connect"
	ileapv1 "github.com/way-platform/ileap-go/proto/gen/wayplatform/connect/ileap/v1"
	"github.com/way-platform/ileap-go/proto/gen/wayplatform/connect/ileap/v1/ileapv1connect"
)

var _ ileapv1connect.ILeapServiceHandler = (*Handler)(nil)

// Handler implements ILeapServiceHandler using embedded demo data.
type Handler struct {
	ileapv1connect.UnimplementedILeapServiceHandler
	footprints []*ileapv1.ProductFootprint
	tads       []*ileapv1.TAD
}

// NewHandler creates a new Handler with the embedded demo data.
func NewHandler() (*Handler, error) {
	footprints, err := LoadFootprints()
	if err != nil {
		return nil, err
	}
	tads, err := LoadTADs()
	if err != nil {
		return nil, err
	}
	return &Handler{
		footprints: footprints,
		tads:       tads,
	}, nil
}

// GetFootprint returns a single footprint by ID.
func (h *Handler) GetFootprint(
	_ context.Context, req *ileapv1.GetFootprintRequest,
) (*ileapv1.GetFootprintResponse, error) {
	for _, fp := range h.footprints {
		if fp.GetId() == req.GetId() {
			resp := new(ileapv1.GetFootprintResponse)
			resp.SetData(fp)
			return resp, nil
		}
	}
	return nil, connect.NewError(connect.CodeNotFound, nil)
}

// ListFootprints returns a filtered, limited list of footprints.
func (h *Handler) ListFootprints(
	_ context.Context, req *ileapv1.ListFootprintsRequest,
) (*ileapv1.ListFootprintsResponse, error) {
	filtered := make([]*ileapv1.ProductFootprint, 0, len(h.footprints))
	for _, fp := range h.footprints {
		if footprintMatchesFilters(fp, req.GetFilters()) {
			filtered = append(filtered, fp)
		}
	}
	total := len(filtered)
	offset := int(req.GetOffset())
	limit := int(req.GetLimit())
	if offset > 0 {
		if offset >= len(filtered) {
			filtered = nil
		} else {
			filtered = filtered[offset:]
		}
	}
	if limit > 0 && len(filtered) > limit {
		filtered = filtered[:limit]
	}
	resp := new(ileapv1.ListFootprintsResponse)
	resp.SetData(filtered)
	resp.SetTotal(int32(total))
	return resp, nil
}

// ListTransportActivityData returns a filtered, paginated list of transport activity data.
func (h *Handler) ListTransportActivityData(
	_ context.Context, req *ileapv1.ListTransportActivityDataRequest,
) (*ileapv1.ListTransportActivityDataResponse, error) {
	filtered := h.tads
	if hasFilter(req) {
		filtered = filterTADs(h.tads, req.GetFilters())
	}
	total := len(filtered)
	offset := int(req.GetOffset())
	limit := int(req.GetLimit())
	if offset > 0 {
		if offset >= len(filtered) {
			filtered = nil
		} else {
			filtered = filtered[offset:]
		}
	}
	if limit > 0 && len(filtered) > limit {
		filtered = filtered[:limit]
	}
	resp := new(ileapv1.ListTransportActivityDataResponse)
	resp.SetData(filtered)
	resp.SetTotal(int32(total))
	return resp, nil
}

func hasFilter(req *ileapv1.ListTransportActivityDataRequest) bool {
	return len(req.GetFilters()) > 0
}

func filterTADs(
	tads []*ileapv1.TAD,
	filters []*ileapv1.ListTransportActivityDataRequest_Filter,
) []*ileapv1.TAD {
	result := make([]*ileapv1.TAD, 0, len(tads))
	for _, tad := range tads {
		if tadMatchesFilter(tad, filters) {
			result = append(result, tad)
		}
	}
	return result
}

func tadMatchesFilter(
	tad *ileapv1.TAD,
	filters []*ileapv1.ListTransportActivityDataRequest_Filter,
) bool {
	for _, filter := range filters {
		name := strings.ToLower(filter.GetFieldPath())
		value := filter.GetValue()
		switch name {
		case "activityid":
			if !strings.EqualFold(tad.GetActivityId(), value) {
				return false
			}
		case "mode":
			if !strings.EqualFold(tad.GetMode(), value) {
				return false
			}
		case "packagingortreqtype":
			if !strings.EqualFold(tad.GetPackagingOrTrEqType(), value) {
				return false
			}
		case "feedstock", "energycarriers.feedstocks.feedstock":
			if !tadHasFeedstock(tad, value) {
				return false
			}
		default:
			continue
		}
	}
	return true
}

func tadHasFeedstock(tad *ileapv1.TAD, feedstock string) bool {
	for _, ec := range tad.GetEnergyCarriers() {
		for _, fs := range ec.GetFeedstocks() {
			if strings.EqualFold(fs.GetFeedstock(), feedstock) {
				return true
			}
		}
	}
	return false
}

func footprintMatchesFilters(
	fp *ileapv1.ProductFootprint,
	filters []*ileapv1.ListFootprintsRequest_Filter,
) bool {
	for _, filter := range filters {
		name := strings.ToLower(filter.GetFieldPath())
		value := filter.GetValue()
		switch name {
		case "productcategorycpc":
			if !strings.EqualFold(fp.GetProductCategoryCpc(), value) {
				return false
			}
		case "pcf.geographycountry":
			pcf := fp.GetPcf()
			if pcf == nil || !strings.EqualFold(pcf.GetGeographyCountry(), value) {
				return false
			}
		case "productids":
			if !containsFold(fp.GetProductIds(), value) {
				return false
			}
		case "companyids":
			if !containsFold(fp.GetCompanyIds(), value) {
				return false
			}
		default:
			continue
		}
	}
	return true
}

func containsFold(values []string, value string) bool {
	return slices.ContainsFunc(values, func(candidate string) bool {
		return strings.EqualFold(candidate, value)
	})
}
