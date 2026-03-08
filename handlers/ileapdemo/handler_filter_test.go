package ileapdemo

import (
	"testing"
	"time"

	ileapv1 "github.com/way-platform/ileap-go/proto/gen/wayplatform/connect/ileap/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestFootprintMatchesFilters_ORSemantics(t *testing.T) {
	fp := new(ileapv1.ProductFootprint)
	fp.SetProductCategoryCpc("83117")
	fp.SetCompanyIds([]string{"acme"})
	fp.SetCreated(timestamppb.New(time.Date(2022, 3, 1, 9, 32, 20, 0, time.UTC)))

	tests := []struct {
		name    string
		filters []*ileapv1.Filter
		want    bool
	}{
		{
			name: "single eq match",
			filters: []*ileapv1.Filter{
				newFilter("productCategoryCpc", ileapv1.Filter_EQ, "83117"),
			},
			want: true,
		},
		{
			name: "single eq no match",
			filters: []*ileapv1.Filter{
				newFilter("productCategoryCpc", ileapv1.Filter_EQ, "99999"),
			},
			want: false,
		},
		{
			name: "multiple filters are OR",
			filters: []*ileapv1.Filter{
				newFilter("productCategoryCpc", ileapv1.Filter_EQ, "99999"),
				newFilter("companyIds", ileapv1.Filter_EQ, "acme"),
			},
			want: true,
		},
		{
			name: "unsupported-only filters are ignored",
			filters: []*ileapv1.Filter{
				newFilter("unknownField", ileapv1.Filter_EQ, "x"),
			},
			want: true,
		},
		{
			name: "created ge supported",
			filters: []*ileapv1.Filter{
				newFilter("created", ileapv1.Filter_GE, "2022-03-01T09:32:20Z"),
			},
			want: true,
		},
		{
			name: "created lt supported",
			filters: []*ileapv1.Filter{
				newFilter("created", ileapv1.Filter_LT, "2022-03-01T09:32:20Z"),
			},
			want: false,
		},
		{
			name: "invalid timestamp filter ignored",
			filters: []*ileapv1.Filter{
				newFilter("created", ileapv1.Filter_GE, "not-a-date"),
			},
			want: true,
		},
		{
			name: "unsupported operator for string field ignored",
			filters: []*ileapv1.Filter{
				newFilter("productCategoryCpc", ileapv1.Filter_GT, "83117"),
			},
			want: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := footprintMatchesFilters(fp, tc.filters)
			if got != tc.want {
				t.Fatalf("footprintMatchesFilters() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestTADMatchesFilters_ORSemantics(t *testing.T) {
	tad := new(ileapv1.TAD)
	tad.SetActivityId("a-1")
	tad.SetMode("Road")
	tad.SetPackagingOrTrEqType("Pallet")
	tad.SetEnergyCarriers([]*ileapv1.EnergyCarrier{
		func() *ileapv1.EnergyCarrier {
			ec := new(ileapv1.EnergyCarrier)
			ec.SetFeedstocks([]*ileapv1.Feedstock{
				func() *ileapv1.Feedstock {
					fs := new(ileapv1.Feedstock)
					fs.SetFeedstock("Fossil")
					return fs
				}(),
			})
			return ec
		}(),
	})

	tests := []struct {
		name    string
		filters []*ileapv1.Filter
		want    bool
	}{
		{
			name: "single eq match",
			filters: []*ileapv1.Filter{
				newFilter("mode", ileapv1.Filter_EQ, "Road"),
			},
			want: true,
		},
		{
			name: "single ne no match",
			filters: []*ileapv1.Filter{
				newFilter("mode", ileapv1.Filter_NE, "Road"),
			},
			want: false,
		},
		{
			name: "multiple filters are OR",
			filters: []*ileapv1.Filter{
				newFilter("mode", ileapv1.Filter_EQ, "Rail"),
				newFilter("feedstock", ileapv1.Filter_EQ, "Fossil"),
			},
			want: true,
		},
		{
			name: "unsupported-only filters are ignored",
			filters: []*ileapv1.Filter{
				newFilter("distance", ileapv1.Filter_EQ, "100"),
				newFilter("mode", ileapv1.Filter_GT, "Road"),
			},
			want: true,
		},
		{
			name: "supported filters with no matches",
			filters: []*ileapv1.Filter{
				newFilter("mode", ileapv1.Filter_EQ, "Rail"),
				newFilter("feedstock", ileapv1.Filter_EQ, "Grid"),
			},
			want: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := tadMatchesFilter(tad, tc.filters)
			if got != tc.want {
				t.Fatalf("tadMatchesFilter() = %v, want %v", got, tc.want)
			}
		})
	}
}

func newFilter(fieldPath string, op ileapv1.Filter_Operator, value string) *ileapv1.Filter {
	filter := new(ileapv1.Filter)
	filter.SetFieldPath(fieldPath)
	filter.SetOperator(op)
	filter.SetValue(value)
	return filter
}
