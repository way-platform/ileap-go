package ileap

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/way-platform/ileap-go/ileapv1pb"
)

// FilterV2 is a limited implementation of PACT v2 filters.
type FilterV2 struct {
	// Conjuctions are multiple filter predicates that are ANDed together.
	Conjuctions []FilterPredicateV2 `json:"conjuctions"`
}

// UnmarshalString unmarshals a filter from a string.
func (f *FilterV2) UnmarshalString(filter string) error {
	f.Conjuctions = f.Conjuctions[:0]
	data := strings.TrimSpace(filter)
	for conjuction := range strings.SplitSeq(data, " and ") {
		conjuction = strings.TrimSpace(conjuction)
		if conjuction == "" {
			continue
		}
		var predicate FilterPredicateV2
		if err := predicate.UnmarshalString(conjuction); err != nil {
			return fmt.Errorf("invalid filter: `%s`: %w", data, err)
		}
		f.Conjuctions = append(f.Conjuctions, predicate)
	}
	return nil
}

// MatchesFootprint returns true if all predicates in the filter match the provided footprint.
func (f *FilterV2) MatchesFootprint(footprint *ileapv1pb.ProductFootprint) bool {
	for _, predicate := range f.Conjuctions {
		if !predicate.MatchesFootprint(footprint) {
			return false
		}
	}
	return true
}

// FilterPredicateV2 is a single filter predicate.
type FilterPredicateV2 struct {
	LHS string `json:"lhs"`
	// Operator is the operator of the predicate.
	Operator string `json:"operator"`
	// RHS is the right hand side of the predicate.
	RHS string `json:"rhs"`
}

// UnmarshalString unmarshals a filter predicate from a string.
func (f *FilterPredicateV2) UnmarshalString(predicate string) error {
	data := strings.TrimSpace(predicate)
	if strings.HasPrefix(data, "(") && strings.HasSuffix(data, ")") {
		data = data[1 : len(data)-1]
	}
	if strings.HasPrefix(data, "productIds/any(productId:(productId eq ") &&
		strings.HasSuffix(data, "))") {
		f.LHS = "productIds"
		f.Operator = "any/eq"
		f.RHS = data[len("productIds/any(productId:(productId eq ") : len(data)-len("))")]
		return nil
	}
	if strings.HasPrefix(data, "companyIds/any(companyId:(companyId eq ") &&
		strings.HasSuffix(data, "))") {
		f.LHS = "companyIds"
		f.Operator = "any/eq"
		f.RHS = data[len("companyIds/any(companyId:(companyId eq ") : len(data)-len("))")]
		return nil
	}
	// Use SplitN to limit splits to 3, so RHS can contain spaces.
	// This handles timestamps with +HH:MM timezone offsets that URL-decode
	// '+' as space (e.g., '2023-06-27T13:00:00.000+00:00' becomes
	// '2023-06-27T13:00:00.000 00:00' after URL decoding).
	fields := strings.SplitN(data, " ", 3)
	if len(fields) != 3 {
		return fmt.Errorf("invalid predicate: `%s`", data)
	}
	switch lhs := fields[0]; lhs {
	case "pcf/geographyCountry", "productCategoryCpc", "created", "updated":
		f.LHS = lhs
	default:
		return fmt.Errorf("invalid predicate LHS: `%s`", lhs)
	}
	switch operator := fields[1]; operator {
	case "eq", "gt", "lt", "ge", "le":
		// ACT conformance TC20 uses "ge" for timestamp filtering.
		f.Operator = operator
	default:
		return fmt.Errorf("invalid predicate operator: `%s`", operator)
	}
	if !strings.HasPrefix(fields[2], "'") || !strings.HasSuffix(fields[2], "'") {
		return fmt.Errorf("invalid predicate RHS: `%s`", fields[2])
	}
	f.RHS = fields[2]
	return nil
}

// MatchesFootprint returns true if the predicate matches the provided footprint.
func (f *FilterPredicateV2) MatchesFootprint(footprint *ileapv1pb.ProductFootprint) bool {
	if f.Operator == "any/eq" {
		return f.matchesAnyEq(footprint)
	}
	var lhsValue string
	switch f.LHS {
	case "pcf/geographyCountry":
		if pcf := footprint.GetPcf(); pcf != nil && pcf.HasGeographyCountry() {
			lhsValue = pcf.GetGeographyCountry()
		}
	case "productCategoryCpc":
		lhsValue = footprint.GetProductCategoryCpc()
	case "created":
		if footprint.HasCreated() {
			lhsValue = footprint.GetCreated().AsTime().Format(time.RFC3339)
		}
	case "updated":
		if footprint.HasUpdated() {
			lhsValue = footprint.GetUpdated().AsTime().Format(time.RFC3339)
		}
	default:
		return false
	}
	if !strings.HasPrefix(f.RHS, "'") || !strings.HasSuffix(f.RHS, "'") {
		return false
	}
	rhsValue := strings.Trim(f.RHS, "'")
	switch f.Operator {
	case "eq":
		return lhsValue == rhsValue
	case "gt":
		return lhsValue > rhsValue
	case "ge":
		// Added to satisfy ACT conformance test case 20.
		return lhsValue >= rhsValue
	case "lt":
		return lhsValue < rhsValue
	case "le":
		return lhsValue <= rhsValue
	default:
		return false
	}
}

func (f *FilterPredicateV2) matchesAnyEq(footprint *ileapv1pb.ProductFootprint) bool {
	var lhsValue []string
	switch f.LHS {
	case "productIds":
		lhsValue = footprint.GetProductIds()
	case "companyIds":
		lhsValue = footprint.GetCompanyIds()
	default:
		return false
	}
	if !strings.HasPrefix(f.RHS, "'") || !strings.HasSuffix(f.RHS, "'") {
		return false
	}
	rhsValue := strings.Trim(f.RHS, "'")
	return slices.Contains(lhsValue, rhsValue)
}
