package ileapdemo

import (
	"context"
	"slices"
	"strings"
	"time"

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
	filters []*ileapv1.Filter,
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
	filters []*ileapv1.Filter,
) bool {
	// Concatenated filters are evaluated disjunctively (OR). Unsupported
	// field/operator pairs are ignored.
	supported := false
	for _, filter := range filters {
		ok, match := tadMatchesSingleFilter(tad, filter)
		if !ok {
			continue
		}
		supported = true
		if match {
			return true
		}
	}
	return !supported
}

func tadMatchesSingleFilter(tad *ileapv1.TAD, filter *ileapv1.Filter) (bool, bool) {
	name := strings.ToLower(filter.GetFieldPath())
	value := filter.GetValue()
	switch name {
	case "activityid":
		return matchesStringFilter(tad.GetActivityId(), value, filter.GetOperator())
	case "mode":
		return matchesStringFilter(tad.GetMode(), value, filter.GetOperator())
	case "packagingortreqtype":
		return matchesStringFilter(tad.GetPackagingOrTrEqType(), value, filter.GetOperator())
	case "feedstock", "energycarriers.feedstocks.feedstock":
		return matchesFeedstockFilter(tad, value, filter.GetOperator())
	default:
		return false, false
	}
}

func matchesFeedstockFilter(
	tad *ileapv1.TAD,
	feedstock string,
	operator ileapv1.Filter_Operator,
) (bool, bool) {
	if operator != ileapv1.Filter_OPERATOR_UNSPECIFIED &&
		operator != ileapv1.Filter_EQ &&
		operator != ileapv1.Filter_NE {
		return false, false
	}
	found := false
	for _, ec := range tad.GetEnergyCarriers() {
		for _, fs := range ec.GetFeedstocks() {
			if strings.EqualFold(fs.GetFeedstock(), feedstock) {
				found = true
				break
			}
		}
	}
	if operator == ileapv1.Filter_NE {
		return true, !found
	}
	return true, found
}

func footprintMatchesFilters(
	fp *ileapv1.ProductFootprint,
	filters []*ileapv1.Filter,
) bool {
	// Concatenated filters are evaluated disjunctively (OR). Unsupported
	// field/operator pairs are ignored.
	supported := false
	for _, filter := range filters {
		ok, match := footprintMatchesSingleFilter(fp, filter)
		if !ok {
			continue
		}
		supported = true
		if match {
			return true
		}
	}
	return !supported
}

func footprintMatchesSingleFilter(
	fp *ileapv1.ProductFootprint,
	filter *ileapv1.Filter,
) (bool, bool) {
	name := strings.ToLower(filter.GetFieldPath())
	value := filter.GetValue()
	switch name {
	case "productcategorycpc":
		return matchesStringFilter(fp.GetProductCategoryCpc(), value, filter.GetOperator())
	case "pcf.geographycountry":
		pcf := fp.GetPcf()
		if pcf == nil {
			return true, false
		}
		return matchesStringFilter(
			pcf.GetGeographyCountry(),
			value,
			filter.GetOperator(),
		)
	case "productids":
		return containsFold(fp.GetProductIds(), value, filter.GetOperator())
	case "companyids":
		return containsFold(fp.GetCompanyIds(), value, filter.GetOperator())
	case "created":
		return matchesTimestampFilter(fp.GetCreated(), value, filter.GetOperator())
	case "updated":
		return matchesTimestampFilter(fp.GetUpdated(), value, filter.GetOperator())
	default:
		return false, false
	}
}

func containsFold(
	values []string,
	value string,
	operator ileapv1.Filter_Operator,
) (bool, bool) {
	contains := slices.ContainsFunc(values, func(candidate string) bool {
		return strings.EqualFold(candidate, value)
	})
	switch operator {
	case ileapv1.Filter_OPERATOR_UNSPECIFIED, ileapv1.Filter_EQ:
		return true, contains
	case ileapv1.Filter_NE:
		return true, !contains
	default:
		return false, false
	}
}

func matchesStringFilter(
	candidate string,
	value string,
	operator ileapv1.Filter_Operator,
) (bool, bool) {
	switch operator {
	case ileapv1.Filter_OPERATOR_UNSPECIFIED, ileapv1.Filter_EQ:
		return true, strings.EqualFold(candidate, value)
	case ileapv1.Filter_NE:
		return true, !strings.EqualFold(candidate, value)
	default:
		return false, false
	}
}

func matchesTimestampFilter(
	candidate interface{ AsTime() time.Time },
	value string,
	operator ileapv1.Filter_Operator,
) (bool, bool) {
	if candidate == nil {
		return false, false
	}
	want, err := parseRFC3339Value(value)
	if err != nil {
		return false, false
	}
	left := candidate.AsTime().UTC()
	switch operator {
	case ileapv1.Filter_OPERATOR_UNSPECIFIED, ileapv1.Filter_EQ:
		return true, left.Equal(want)
	case ileapv1.Filter_NE:
		return true, !left.Equal(want)
	case ileapv1.Filter_LT:
		return true, left.Before(want)
	case ileapv1.Filter_LE:
		return true, left.Before(want) || left.Equal(want)
	case ileapv1.Filter_GT:
		return true, left.After(want)
	case ileapv1.Filter_GE:
		return true, left.After(want) || left.Equal(want)
	default:
		return false, false
	}
}

func parseRFC3339Value(value string) (time.Time, error) {
	if parsed, err := time.Parse(time.RFC3339Nano, value); err == nil {
		return parsed.UTC(), nil
	}
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return time.Time{}, err
	}
	return parsed.UTC(), nil
}
