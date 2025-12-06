package largedataset

import (
	"testing"
)

func TestParseFilter(t *testing.T) {
	header := []string{"ID", "Name", "Age", "Country", "Amount"}

	tests := []struct {
		name        string
		expr        string
		header      []string
		wantErr     bool
		checkFilter func(*Filter) bool
	}{
		{
			name:    "numeric greater than",
			expr:    "Amount > 100",
			header:  header,
			wantErr: false,
			checkFilter: func(f *Filter) bool {
				return f.ColumnIndex == 4 && f.Operator == OpGreaterThan && f.IsNumeric && f.NumValue == 100
			},
		},
		{
			name:    "string equality",
			expr:    "Country = 'FR'",
			header:  header,
			wantErr: false,
			checkFilter: func(f *Filter) bool {
				return f.ColumnIndex == 3 && f.Operator == OpEqual && f.Value == "FR"
			},
		},
		{
			name:    "string equality without quotes",
			expr:    "Country = FR",
			header:  header,
			wantErr: false,
			checkFilter: func(f *Filter) bool {
				return f.ColumnIndex == 3 && f.Operator == OpEqual && f.Value == "FR"
			},
		},
		{
			name:    "column by index",
			expr:    "2 >= 18",
			header:  header,
			wantErr: false,
			checkFilter: func(f *Filter) bool {
				return f.ColumnIndex == 2 && f.Operator == OpGreaterThanOrEqual && f.IsNumeric
			},
		},
		{
			name:    "contains operator",
			expr:    "Name contains John",
			header:  header,
			wantErr: false,
			checkFilter: func(f *Filter) bool {
				return f.ColumnIndex == 1 && f.Operator == OpContains && f.Value == "John"
			},
		},
		{
			name:    "startswith operator",
			expr:    "Country startswith A",
			header:  header,
			wantErr: false,
			checkFilter: func(f *Filter) bool {
				return f.ColumnIndex == 3 && f.Operator == OpStartsWith && f.Value == "A"
			},
		},
		{
			name:    "regex operator",
			expr:    "Name ~= ^[A-Z].*",
			header:  header,
			wantErr: false,
			checkFilter: func(f *Filter) bool {
				return f.ColumnIndex == 1 && f.Operator == OpRegex && f.RegexPat != nil
			},
		},
		{
			name:    "invalid: no operator",
			expr:    "Amount 100",
			header:  header,
			wantErr: true,
		},
		{
			name:    "invalid: column not found",
			expr:    "Unknown > 100",
			header:  header,
			wantErr: true,
		},
		{
			name:    "invalid: bad regex",
			expr:    "Name ~= [invalid",
			header:  header,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter, err := ParseFilter(tt.expr, tt.header)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseFilter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.checkFilter != nil && !tt.checkFilter(filter) {
				t.Errorf("ParseFilter() filter validation failed for %q", tt.expr)
			}
		})
	}
}

func TestFilterEvaluate(t *testing.T) {
	header := []string{"ID", "Name", "Age", "Country", "Amount"}

	tests := []struct {
		name       string
		filterExpr string
		record     []string
		wantMatch  bool
		wantErr    bool
	}{
		{
			name:       "numeric greater than - match",
			filterExpr: "Amount > 100",
			record:     []string{"1", "John", "25", "US", "150.50"},
			wantMatch:  true,
		},
		{
			name:       "numeric greater than - no match",
			filterExpr: "Amount > 100",
			record:     []string{"1", "John", "25", "US", "50.00"},
			wantMatch:  false,
		},
		{
			name:       "string equality - match",
			filterExpr: "Country = 'FR'",
			record:     []string{"1", "Pierre", "30", "FR", "100"},
			wantMatch:  true,
		},
		{
			name:       "string equality - no match",
			filterExpr: "Country = 'FR'",
			record:     []string{"1", "John", "25", "US", "100"},
			wantMatch:  false,
		},
		{
			name:       "less than or equal - match",
			filterExpr: "Age <= 30",
			record:     []string{"1", "John", "25", "US", "100"},
			wantMatch:  true,
		},
		{
			name:       "contains - match",
			filterExpr: "Name contains oh",
			record:     []string{"1", "John", "25", "US", "100"},
			wantMatch:  true,
		},
		{
			name:       "contains - no match",
			filterExpr: "Name contains xyz",
			record:     []string{"1", "John", "25", "US", "100"},
			wantMatch:  false,
		},
		{
			name:       "startswith - match",
			filterExpr: "Name startswith Jo",
			record:     []string{"1", "John", "25", "US", "100"},
			wantMatch:  true,
		},
		{
			name:       "regex - match",
			filterExpr: "Name ~= ^J.*n$",
			record:     []string{"1", "John", "25", "US", "100"},
			wantMatch:  true,
		},
		{
			name:       "not equal - match",
			filterExpr: "Country != 'US'",
			record:     []string{"1", "Pierre", "30", "FR", "100"},
			wantMatch:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter, err := ParseFilter(tt.filterExpr, header)
			if err != nil {
				t.Fatalf("Failed to parse filter: %v", err)
			}

			match, err := filter.Evaluate(tt.record)
			if (err != nil) != tt.wantErr {
				t.Errorf("Evaluate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if match != tt.wantMatch {
				t.Errorf("Evaluate() match = %v, want %v", match, tt.wantMatch)
			}
		})
	}
}

func TestFilterSetEvaluate(t *testing.T) {
	header := []string{"ID", "Name", "Age", "Country", "Amount"}

	tests := []struct {
		name        string
		filterExprs []string
		record      []string
		wantMatch   bool
	}{
		{
			name: "multiple filters - all match",
			filterExprs: []string{
				"Amount > 100",
				"Country = 'US'",
				"Age >= 18",
			},
			record:    []string{"1", "John", "25", "US", "150.50"},
			wantMatch: true,
		},
		{
			name: "multiple filters - one fails",
			filterExprs: []string{
				"Amount > 100",
				"Country = 'FR'",
				"Age >= 18",
			},
			record:    []string{"1", "John", "25", "US", "150.50"},
			wantMatch: false,
		},
		{
			name: "single filter - match",
			filterExprs: []string{
				"Name contains John",
			},
			record:    []string{"1", "Johnson", "25", "US", "100"},
			wantMatch: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filterSet, err := NewFilterSet(tt.filterExprs, header)
			if err != nil {
				t.Fatalf("Failed to create filter set: %v", err)
			}

			match, err := filterSet.Evaluate(tt.record)
			if err != nil {
				t.Errorf("Evaluate() error = %v", err)
				return
			}
			if match != tt.wantMatch {
				t.Errorf("Evaluate() match = %v, want %v", match, tt.wantMatch)
			}
		})
	}
}
