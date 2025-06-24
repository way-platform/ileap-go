package ileap

import (
	"reflect"
	"testing"

	"github.com/way-platform/ileap-go/openapi/ileapv0"
)

func TestFilterV2_UnmarshalString(t *testing.T) {
	testCases := []struct {
		name string
		data string
		want FilterV2
	}{
		{
			name: "empty",
			data: "",
			want: FilterV2{},
		},

		{
			name: "single predicate",
			data: "pcf/geographyCountry eq 'US'",
			want: FilterV2{
				Conjuctions: []FilterPredicateV2{
					{
						LHS:      "pcf/geographyCountry",
						Operator: "eq",
						RHS:      "'US'",
					},
				},
			},
		},

		{
			name: "multiple predicates",
			data: "(pcf/geographyCountry eq 'FR') and (pcf/geographyCountry eq 'DE')",
			want: FilterV2{
				Conjuctions: []FilterPredicateV2{
					{
						LHS:      "pcf/geographyCountry",
						Operator: "eq",
						RHS:      "'FR'",
					},
					{
						LHS:      "pcf/geographyCountry",
						Operator: "eq",
						RHS:      "'DE'",
					},
				},
			},
		},

		{
			name: "productCategoryCpc",
			data: "productCategoryCpc eq '6398'",
			want: FilterV2{
				Conjuctions: []FilterPredicateV2{
					{
						LHS:      "productCategoryCpc",
						Operator: "eq",
						RHS:      "'6398'",
					},
				},
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			var filter FilterV2
			if err := filter.UnmarshalString(testCase.data); err != nil {
				t.Fatalf("UnmarshalString(%q) = %v", testCase.data, err)
			}
			if !reflect.DeepEqual(filter, testCase.want) {
				t.Fatalf("UnmarshalString(%q) = %v, want %v", testCase.data, filter, testCase.want)
			}
		})
	}
}

func TestFilterV2_MatchesFootprint(t *testing.T) {
	testCases := []struct {
		name      string
		footprint *ileapv0.ProductFootprintForILeapType
		filter    FilterV2
		want      bool
	}{
		{
			name:      "empty",
			footprint: &ileapv0.ProductFootprintForILeapType{},
			filter:    FilterV2{},
			want:      true,
		},

		{
			name: "single predicate",
			footprint: &ileapv0.ProductFootprintForILeapType{
				Pcf: ileapv0.CarbonFootprint{
					GeographyCountry: ptr("US"),
				},
			},
			filter: FilterV2{
				Conjuctions: []FilterPredicateV2{
					{
						LHS:      "pcf/geographyCountry",
						Operator: "eq",
						RHS:      "'US'",
					},
				},
			},
			want: true,
		},

		{
			name: "single predicate, no match",
			footprint: &ileapv0.ProductFootprintForILeapType{
				Pcf: ileapv0.CarbonFootprint{
					GeographyCountry: ptr("FR"),
				},
			},
			filter: FilterV2{
				Conjuctions: []FilterPredicateV2{
					{
						LHS:      "pcf/geographyCountry",
						Operator: "eq",
						RHS:      "'US'",
					},
				},
			},
			want: false,
		},

		{
			name: "multiple predicates",
			footprint: &ileapv0.ProductFootprintForILeapType{
				ProductCategoryCpc: "6398",
				Pcf: ileapv0.CarbonFootprint{
					GeographyCountry: ptr("US"),
				},
			},
			filter: FilterV2{
				Conjuctions: []FilterPredicateV2{
					{
						LHS:      "pcf/geographyCountry",
						Operator: "eq",
						RHS:      "'US'",
					},
					{
						LHS:      "productCategoryCpc",
						Operator: "eq",
						RHS:      "'6398'",
					},
				},
			},
			want: true,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := testCase.filter.MatchesFootprint(testCase.footprint)
			if got != testCase.want {
				t.Fatalf("MatchesFootprint(%v) = %v, want %v", testCase.footprint, got, testCase.want)
			}
		})
	}
}

func ptr[T any](t T) *T {
	return &t
}
