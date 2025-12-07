package largedataset

import (
	"strings"
	"testing"
)

func TestComplexFilterExpressions(t *testing.T) {
	header := []string{"Symbol", "Price", "Volume", "Exchange"}

	tests := []struct {
		name       string
		filterExpr string
		record     []string
		expected   bool
		shouldErr  bool
	}{
		// Simple AND tests
		{
			name:       "Simple AND - both true",
			filterExpr: "Price > 100 AND Volume > 1000",
			record:     []string{"AAPL", "150", "2000", "NASDAQ"},
			expected:   true,
		},
		{
			name:       "Simple AND - first false",
			filterExpr: "Price > 100 AND Volume > 1000",
			record:     []string{"AAPL", "50", "2000", "NASDAQ"},
			expected:   false,
		},
		{
			name:       "Simple AND - second false",
			filterExpr: "Price > 100 AND Volume > 1000",
			record:     []string{"AAPL", "150", "500", "NASDAQ"},
			expected:   false,
		},
		{
			name:       "Simple AND - both false",
			filterExpr: "Price > 100 AND Volume > 1000",
			record:     []string{"AAPL", "50", "500", "NASDAQ"},
			expected:   false,
		},

		// Simple OR tests
		{
			name:       "Simple OR - both true",
			filterExpr: "Price > 100 OR Volume > 1000",
			record:     []string{"AAPL", "150", "2000", "NASDAQ"},
			expected:   true,
		},
		{
			name:       "Simple OR - first true",
			filterExpr: "Price > 100 OR Volume > 1000",
			record:     []string{"AAPL", "150", "500", "NASDAQ"},
			expected:   true,
		},
		{
			name:       "Simple OR - second true",
			filterExpr: "Price > 100 OR Volume > 1000",
			record:     []string{"AAPL", "50", "2000", "NASDAQ"},
			expected:   true,
		},
		{
			name:       "Simple OR - both false",
			filterExpr: "Price > 100 OR Volume > 1000",
			record:     []string{"AAPL", "50", "500", "NASDAQ"},
			expected:   false,
		},

		// Mixed AND/OR tests (OR has lower precedence)
		{
			name:       "Mixed AND/OR - (A AND B) OR C",
			filterExpr: "Price > 100 AND Volume > 1000 OR Exchange = 'NYSE'",
			record:     []string{"AAPL", "150", "2000", "NASDAQ"},
			expected:   true, // Price AND Volume are true
		},
		{
			name:       "Mixed AND/OR - C true makes OR true",
			filterExpr: "Price > 100 AND Volume > 1000 OR Exchange = 'NYSE'",
			record:     []string{"AAPL", "50", "500", "NYSE"},
			expected:   true, // Exchange matches
		},
		{
			name:       "Mixed AND/OR - all false",
			filterExpr: "Price > 100 AND Volume > 1000 OR Exchange = 'NYSE'",
			record:     []string{"AAPL", "50", "500", "NASDAQ"},
			expected:   false,
		},

		// Parentheses tests
		{
			name:       "Parentheses - (A OR B) AND C",
			filterExpr: "(Price > 100 OR Volume > 5000) AND Exchange = 'NASDAQ'",
			record:     []string{"AAPL", "150", "1000", "NASDAQ"},
			expected:   true, // Price > 100 is true, Exchange matches
		},
		{
			name:       "Parentheses - (A OR B) AND C - C false",
			filterExpr: "(Price > 100 OR Volume > 5000) AND Exchange = 'NASDAQ'",
			record:     []string{"AAPL", "150", "1000", "NYSE"},
			expected:   false, // Exchange doesn't match
		},
		{
			name:       "Parentheses - (A OR B) AND C - both OR false",
			filterExpr: "(Price > 100 OR Volume > 5000) AND Exchange = 'NASDAQ'",
			record:     []string{"AAPL", "50", "1000", "NASDAQ"},
			expected:   false, // Neither Price nor Volume condition is true
		},

		// Complex nested tests
		{
			name:       "Complex nested - ((A OR B) AND C) OR D",
			filterExpr: "((Price > 100 OR Volume > 5000) AND Exchange = 'NASDAQ') OR Symbol = 'MSFT'",
			record:     []string{"MSFT", "50", "100", "NYSE"},
			expected:   true, // Symbol matches
		},
		{
			name:       "Complex nested - nested true",
			filterExpr: "((Price > 100 OR Volume > 5000) AND Exchange = 'NASDAQ') OR Symbol = 'MSFT'",
			record:     []string{"AAPL", "150", "1000", "NASDAQ"},
			expected:   true, // Nested condition is true
		},

		// String operators with logical operators
		{
			name:       "String contains with OR",
			filterExpr: "Symbol contains 'AA' OR Symbol contains 'MS'",
			record:     []string{"MSFT", "150", "1000", "NASDAQ"},
			expected:   true,
		},
		{
			name:       "String startswith with AND",
			filterExpr: "Symbol startswith 'A' AND Exchange = 'NASDAQ'",
			record:     []string{"AAPL", "150", "1000", "NASDAQ"},
			expected:   true,
		},

		// Multiple conditions
		{
			name:       "Three conditions with AND",
			filterExpr: "Price > 100 AND Volume > 1000 AND Exchange = 'NASDAQ'",
			record:     []string{"AAPL", "150", "2000", "NASDAQ"},
			expected:   true,
		},
		{
			name:       "Three conditions with OR",
			filterExpr: "Price > 200 OR Volume > 5000 OR Symbol = 'AAPL'",
			record:     []string{"AAPL", "50", "500", "NASDAQ"},
			expected:   true, // Symbol matches
		},
		{
			name:       "Multiple mixed conditions",
			filterExpr: "Price > 100 AND Volume > 1000 OR Symbol = 'MSFT' AND Exchange = 'NYSE'",
			record:     []string{"MSFT", "50", "500", "NYSE"},
			expected:   true, // MSFT AND NYSE
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsed, err := parseComplexFilter(tt.filterExpr, header)
			if tt.shouldErr {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			result, err := parsed.Evaluate(tt.record)
			if err != nil {
				t.Errorf("evaluation error: %v", err)
				return
			}

			if result != tt.expected {
				t.Errorf("expected %v but got %v for filter %q with record %v",
					tt.expected, result, tt.filterExpr, tt.record)
			}
		})
	}
}

func TestFilterSetWithComplexExpressions(t *testing.T) {
	header := []string{"Symbol", "Price", "Volume", "Exchange"}

	tests := []struct {
		name        string
		filterExprs []string
		record      []string
		expected    bool
	}{
		{
			name:        "Single complex filter with OR",
			filterExprs: []string{"Price > 100 OR Volume > 5000"},
			record:      []string{"AAPL", "150", "1000", "NASDAQ"},
			expected:    true,
		},
		{
			name:        "Multiple filters combined with AND (backward compatible)",
			filterExprs: []string{"Price > 100", "Volume > 1000"},
			record:      []string{"AAPL", "150", "2000", "NASDAQ"},
			expected:    true,
		},
		{
			name:        "Multiple filters - one fails",
			filterExprs: []string{"Price > 100", "Volume > 1000"},
			record:      []string{"AAPL", "50", "2000", "NASDAQ"},
			expected:    false,
		},
		{
			name:        "Multiple complex filters",
			filterExprs: []string{"Price > 100 OR Volume > 5000", "Exchange = 'NASDAQ'"},
			record:      []string{"AAPL", "150", "1000", "NASDAQ"},
			expected:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs, err := NewFilterSet(tt.filterExprs, header)
			if err != nil {
				t.Errorf("unexpected error creating filter set: %v", err)
				return
			}

			result, err := fs.Evaluate(tt.record)
			if err != nil {
				t.Errorf("evaluation error: %v", err)
				return
			}

			if result != tt.expected {
				t.Errorf("expected %v but got %v for filters %v with record %v",
					tt.expected, result, tt.filterExprs, tt.record)
			}
		})
	}
}

func TestFilterGroupString(t *testing.T) {
	header := []string{"Symbol", "Price", "Volume"}

	tests := []struct {
		name       string
		filterExpr string
		expected   string
	}{
		{
			name:       "Simple AND",
			filterExpr: "Price > 100 AND Volume > 1000",
			expected:   `Price > "100" AND Volume > "1000"`,
		},
		{
			name:       "Simple OR",
			filterExpr: "Price > 100 OR Volume > 1000",
			expected:   `Price > "100" OR Volume > "1000"`,
		},
		{
			name:       "Parentheses",
			filterExpr: "(Price > 100 OR Volume > 5000) AND Symbol = 'AAPL'",
			expected:   `(Price > "100" OR Volume > "5000") AND Symbol = "AAPL"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsed, err := parseComplexFilter(tt.filterExpr, header)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			result := parsed.String()
			if result != tt.expected {
				t.Errorf("expected %q but got %q", tt.expected, result)
			}
		})
	}
}

func TestBackwardCompatibility(t *testing.T) {
	header := []string{"Symbol", "Price", "Volume"}

	// Test that old behavior still works
	filterExprs := []string{"Price > 100", "Volume > 1000", "Symbol = 'AAPL'"}
	fs, err := NewFilterSet(filterExprs, header)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should combine with AND
	record := []string{"AAPL", "150", "2000"}
	result, err := fs.Evaluate(record)
	if err != nil {
		t.Fatalf("evaluation error: %v", err)
	}
	if !result {
		t.Errorf("expected true for matching record")
	}

	// Test String() output (order doesn't matter for AND logic)
	strResult := fs.String()
	if !strings.Contains(strResult, `Symbol = "AAPL"`) ||
		!strings.Contains(strResult, `Price > "100"`) ||
		!strings.Contains(strResult, `Volume > "1000"`) ||
		!strings.Contains(strResult, " AND ") {
		t.Errorf("expected all conditions with AND but got %q", strResult)
	}
}
