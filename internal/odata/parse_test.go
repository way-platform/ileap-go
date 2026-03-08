package odata

import (
	"testing"

	ileapv1 "github.com/way-platform/ileap-go/proto/gen/wayplatform/connect/ileap/v1"
)

func TestParseFilter(t *testing.T) {
	testCases := []struct {
		name string
		in   string
		want []string
	}{
		{
			name: "empty",
			in:   "",
			want: nil,
		},
		{
			name: "whitespace",
			in:   "  ",
			want: nil,
		},
		{
			name: "supported and unsupported mixed",
			in: "(pcf/geographyCountry eq 'DE') and " +
				"(productIds/any(productId:(productId eq 'urn:test:1'))) and " +
				"(updated weird '2024-01-01T00:00:00.000Z')",
			want: []string{
				"pcf.geographyCountry|EQ|DE",
				"productIds|EQ|urn:test:1",
			},
		},
		{
			name: "mixed keyword casing",
			in:   "PCF/GeographyCountry EQ 'de' AnD companyIds/AnY(id:(id Eq '123'))",
			want: []string{
				"PCF.GeographyCountry|EQ|de",
				"companyIds|EQ|123",
			},
		},
		{
			name: "any eq with nested alias path",
			in:   "tces/any(t:(t/origin/city eq 'Berlin'))",
			want: []string{
				"tces.origin.city|EQ|Berlin",
			},
		},
		{
			name: "invalid fragments ignored",
			in:   "bad syntax and still bad",
			want: nil,
		},
		{
			name: "companyIds any eq",
			in:   "companyIds/any(companyId:(companyId eq '12345'))",
			want: []string{
				"companyIds|EQ|12345",
			},
		},
		{
			name: "all comparison operators",
			in:   "a eq '1' and b ne '2' and c lt '3' and d le '4' and e gt '5' and f ge '6'",
			want: []string{
				"a|EQ|1",
				"b|NE|2",
				"c|LT|3",
				"d|LE|4",
				"e|GT|5",
				"f|GE|6",
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := ParseFilter(testCase.in)
			assertFilterSet(t, got, testCase.want...)
		})
	}
}

func TestParseClause(t *testing.T) {
	testCases := []struct {
		name string
		in   string
		want string
		ok   bool
	}{
		{
			name: "simple eq",
			in:   "productCategoryCpc eq '83117'",
			want: "productCategoryCpc|EQ|83117",
			ok:   true,
		},
		{
			name: "nested any eq",
			in:   "tces/any(t:(t/origin/city eq 'Berlin'))",
			want: "tces.origin.city|EQ|Berlin",
			ok:   true,
		},
		{
			name: "gt operator",
			in:   "updated gt '2024-01-01T00:00:00.000Z'",
			want: "updated|GT|2024-01-01T00:00:00.000Z",
			ok:   true,
		},
		{
			name: "unsupported operator",
			in:   "updated weird '2024-01-01T00:00:00.000Z'",
			ok:   false,
		},
		{
			name: "unsupported any body",
			in:   "productIds/any(p:(p weird 'x'))",
			ok:   false,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got, ok := parseClause(testCase.in)
			if ok != testCase.ok {
				t.Fatalf("ok=%v, want %v", ok, testCase.ok)
			}
			if !testCase.ok {
				return
			}
			assertFilterSet(t, []*ileapv1.Filter{got}, testCase.want)
		})
	}
}

func TestSplitTopLevelAndClauses(t *testing.T) {
	testCases := []struct {
		name string
		in   string
		want []string
	}{
		{
			name: "simple",
			in:   "a eq '1' and b eq '2'",
			want: []string{"a eq '1' ", " b eq '2'"},
		},
		{
			name: "uppercase and nested parens",
			in:   "(a eq '1') AnD (x/any(t:(t/id eq '2')))",
			want: []string{"(a eq '1') ", " (x/any(t:(t/id eq '2')))"},
		},
		{
			name: "and inside string literal",
			in:   "a eq 'A and B' and b eq '2'",
			want: []string{"a eq 'A and B' ", " b eq '2'"},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := splitTopLevelAndClauses(testCase.in)
			if len(got) != len(testCase.want) {
				t.Fatalf("got %v, want %v", got, testCase.want)
			}
			for i := range got {
				if got[i] != testCase.want[i] {
					t.Fatalf("got %v, want %v", got, testCase.want)
				}
			}
		})
	}
}

func assertFilterSet(t *testing.T, got []*ileapv1.Filter, want ...string) {
	t.Helper()
	gotCounts := make(map[string]int, len(got))
	for _, filter := range got {
		key := filter.GetFieldPath() + "|" + filter.GetOperator().String() + "|" + filter.GetValue()
		gotCounts[key]++
	}
	wantCounts := make(map[string]int, len(want))
	for _, key := range want {
		wantCounts[key]++
	}
	if len(gotCounts) != len(wantCounts) {
		t.Fatalf("got %v, want %v", gotCounts, wantCounts)
	}
	for key, count := range wantCounts {
		if gotCounts[key] != count {
			t.Fatalf("got %v, want %v", gotCounts, wantCounts)
		}
	}
}
