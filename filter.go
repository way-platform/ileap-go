package ileap

import (
	"fmt"
	"strings"

	"github.com/way-platform/ileap-go/openapi/ileapv0"
)

// FilterV2 is a limited implementation of PACT v2 filters.
type FilterV2 struct {
	// Conjuctions are multiple filter predicates that are ANDed together.
	Conjuctions []FilterPredicateV2 `json:"conjuctions"`
}

// UnmarshalString unmarshals a filter from a string.
func (f *FilterV2) UnmarshalString(data string) error {
	data = strings.Trim(data, "() ")
	f.Conjuctions = f.Conjuctions[:0]
	for conjuction := range strings.SplitSeq(data, " and ") {
		conjuction = strings.Trim(conjuction, "() ")
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
func (f *FilterV2) MatchesFootprint(footprint *ileapv0.ProductFootprintForILeapType) bool {
	for _, predicate := range f.Conjuctions {
		if !predicate.MatchesFootprint(footprint) {
			return false
		}
	}
	return true
}

// FilterPredicateV2 is a single filter predicate.
type FilterPredicateV2 struct {
	// LHS is the left hand side of the predicate.
	LHS string `json:"lhs"`
	// Operator is the operator of the predicate.
	Operator string `json:"operator"`
	// RHS is the right hand side of the predicate.
	RHS string `json:"rhs"`
}

// UnmarshalString unmarshals a filter predicate from a string.
func (f *FilterPredicateV2) UnmarshalString(data string) error {
	fields := strings.Fields(data)
	if len(fields) != 3 {
		return fmt.Errorf("invalid predicate: `%s`", data)
	}
	switch lhs := fields[0]; lhs {
	case "pcf/geographyCountry", "productCategoryCpc":
		f.LHS = lhs
	default:
		return fmt.Errorf("invalid predicate LHS: `%s`", lhs)
	}
	switch operator := fields[1]; operator {
	case "eq":
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
func (f *FilterPredicateV2) MatchesFootprint(footprint *ileapv0.ProductFootprintForILeapType) bool {
	var lhsValue string
	switch f.LHS {
	case "pcf/geographyCountry":
		if footprint.Pcf.GeographyCountry != nil {
			lhsValue = *footprint.Pcf.GeographyCountry
		}
	case "productCategoryCpc":
		lhsValue = footprint.ProductCategoryCpc
	default:
		return false
	}
	if !strings.HasPrefix(f.RHS, "'") || !strings.HasSuffix(f.RHS, "'") {
		return false
	}
	rhsValue := strings.Trim(f.RHS, "'")
	if f.Operator == "eq" {
		return lhsValue == rhsValue
	}
	return false
}
