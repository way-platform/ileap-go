package ileap

import (
	"reflect"
	"testing"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/way-platform/ileap-go/ileapv1pb"
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

		{
			name: "productIds, any eq",
			data: "productIds/any(productId:(productId eq 'urn:gtin:5695872369587'))",
			want: FilterV2{
				Conjuctions: []FilterPredicateV2{
					{
						LHS:      "productIds",
						Operator: "any/eq",
						RHS:      "'urn:gtin:5695872369587'",
					},
				},
			},
		},

		{
			name: "companyIds, any eq",
			data: "companyIds/any(companyId:(companyId eq '12345'))",
			want: FilterV2{
				Conjuctions: []FilterPredicateV2{
					{
						LHS:      "companyIds",
						Operator: "any/eq",
						RHS:      "'12345'",
					},
				},
			},
		},

		{
			name: "productCategoryCpc and created",
			data: "(productCategoryCpc eq '6398') and (created gt '1900-01-01T00:00:00.000Z')",
			want: FilterV2{
				Conjuctions: []FilterPredicateV2{
					{
						LHS:      "productCategoryCpc",
						Operator: "eq",
						RHS:      "'6398'",
					},
					{
						LHS:      "created",
						Operator: "gt",
						RHS:      "'1900-01-01T00:00:00.000Z'",
					},
				},
			},
		},

		{
			name: "updated lt",
			data: "updated lt '2023-06-27T13:00:00.000Z'",
			want: FilterV2{
				Conjuctions: []FilterPredicateV2{
					{
						LHS:      "updated",
						Operator: "lt",
						RHS:      "'2023-06-27T13:00:00.000Z'",
					},
				},
			},
		},

		{
			name: "updated combined with geographyCountry",
			data: "pcf/geographyCountry eq 'FR' and updated lt '2023-06-27T13:00:00.000Z'",
			want: FilterV2{
				Conjuctions: []FilterPredicateV2{
					{
						LHS:      "pcf/geographyCountry",
						Operator: "eq",
						RHS:      "'FR'",
					},
					{
						LHS:      "updated",
						Operator: "lt",
						RHS:      "'2023-06-27T13:00:00.000Z'",
					},
				},
			},
		},

		{
			// ACT sends +00:00 timezone; URL-decoding converts + to space.
			// The filter RHS becomes '2023-06-27T13:00:00.000 00:00'.
			name: "updated lt with timezone offset decoded as space",
			data: "updated lt '2023-06-27T13:00:00.000 00:00'",
			want: FilterV2{
				Conjuctions: []FilterPredicateV2{
					{
						LHS:      "updated",
						Operator: "lt",
						RHS:      "'2023-06-27T13:00:00.000 00:00'",
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
		footprint *ileapv1pb.ProductFootprint
		filter    FilterV2
		want      bool
	}{
		{
			name:      "empty",
			footprint: &ileapv1pb.ProductFootprint{},
			filter:    FilterV2{},
			want:      true,
		},

		{
			name: "single predicate",
			footprint: func() *ileapv1pb.ProductFootprint {
				pf := &ileapv1pb.ProductFootprint{}
				pcf := &ileapv1pb.CarbonFootprint{}
				pcf.SetGeographyCountry("US")
				pf.SetPcf(pcf)
				return pf
			}(),
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
			footprint: func() *ileapv1pb.ProductFootprint {
				pf := &ileapv1pb.ProductFootprint{}
				pcf := &ileapv1pb.CarbonFootprint{}
				pcf.SetGeographyCountry("FR")
				pf.SetPcf(pcf)
				return pf
			}(),
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
			name: "created gt",
			footprint: func() *ileapv1pb.ProductFootprint {
				pf := &ileapv1pb.ProductFootprint{}
				pf.SetCreated(timestamppb.New(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)))
				return pf
			}(),
			filter: FilterV2{
				Conjuctions: []FilterPredicateV2{
					{
						LHS:      "created",
						Operator: "gt",
						RHS:      "'2024-01-01T00:00:00.000Z'",
					},
				},
			},
			want: true,
		},

		{
			name: "created gt, no match",
			footprint: func() *ileapv1pb.ProductFootprint {
				pf := &ileapv1pb.ProductFootprint{}
				pf.SetCreated(timestamppb.New(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)))
				return pf
			}(),
			filter: FilterV2{
				Conjuctions: []FilterPredicateV2{
					{
						LHS:      "created",
						Operator: "gt",
						RHS:      "'2024-01-01T00:00:00.000Z'",
					},
				},
			},
			want: false,
		},

		{
			name: "multiple predicates",
			footprint: func() *ileapv1pb.ProductFootprint {
				pf := &ileapv1pb.ProductFootprint{}
				pf.SetProductCategoryCpc("6398")
				pcf := &ileapv1pb.CarbonFootprint{}
				pcf.SetGeographyCountry("US")
				pf.SetPcf(pcf)
				return pf
			}(),
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

		{
			name: "productIds, any eq",
			footprint: func() *ileapv1pb.ProductFootprint {
				pf := &ileapv1pb.ProductFootprint{}
				pf.SetProductIds([]string{"urn:gtin:1234"})
				return pf
			}(),
			filter: FilterV2{
				Conjuctions: []FilterPredicateV2{
					{
						LHS:      "productIds",
						Operator: "any/eq",
						RHS:      "'urn:gtin:1234'",
					},
				},
			},
			want: true,
		},

		{
			name: "productIds, any eq, no match",
			footprint: func() *ileapv1pb.ProductFootprint {
				pf := &ileapv1pb.ProductFootprint{}
				pf.SetProductIds([]string{"urn:gtin:1234"})
				return pf
			}(),
			filter: FilterV2{
				Conjuctions: []FilterPredicateV2{
					{
						LHS:      "productIds",
						Operator: "any/eq",
						RHS:      "'urn:gtin:5678'",
					},
				},
			},
			want: false,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := testCase.filter.MatchesFootprint(testCase.footprint)
			if got != testCase.want {
				t.Fatalf("got %v, want %v", got, testCase.want)
			}
		})
	}
}

