package largedataset

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// FilterOperator represents a comparison operator
type FilterOperator int

const (
	OpEqual FilterOperator = iota
	OpNotEqual
	OpGreaterThan
	OpGreaterThanOrEqual
	OpLessThan
	OpLessThanOrEqual
	OpContains
	OpStartsWith
	OpEndsWith
	OpRegex
)

// Filter represents a single filter condition
type Filter struct {
	ColumnIndex int            // Column index to filter on
	ColumnName  string         // Column name (for display)
	Operator    FilterOperator // Comparison operator
	Value       string         // Value to compare against
	NumValue    float64        // Parsed numeric value (for numeric comparisons)
	IsNumeric   bool           // Whether this is a numeric comparison
	RegexPat    *regexp.Regexp // Compiled regex pattern (for regex operator)
}

// FilterSet represents a collection of filters (AND logic)
type FilterSet struct {
	Filters []*Filter
}

// ParseFilter parses a filter expression like "amount > 100" or "symbol = 'AAPL'"
func ParseFilter(expr string, header []string) (*Filter, error) {
	expr = strings.TrimSpace(expr)

	// Match filter pattern: column operator value
	// Operators: =, !=, >, >=, <, <=, ~= (regex), contains, startswith, endswith
	operators := []struct {
		op      FilterOperator
		pattern string
	}{
		{OpGreaterThanOrEqual, " >= "},
		{OpLessThanOrEqual, " <= "},
		{OpNotEqual, " != "},
		{OpRegex, " ~= "},
		{OpEqual, " = "},
		{OpGreaterThan, " > "},
		{OpLessThan, " < "},
		{OpContains, " contains "},
		{OpStartsWith, " startswith "},
		{OpEndsWith, " endswith "},
	}

	var filter Filter
	var foundOp bool
	var leftPart, rightPart string

	// Try to find an operator (with spaces around it for proper tokenization)
	for _, op := range operators {
		idx := strings.Index(expr, op.pattern)
		if idx >= 0 {
			leftPart = strings.TrimSpace(expr[:idx])
			rightPart = strings.TrimSpace(expr[idx+len(op.pattern):])
			filter.Operator = op.op
			foundOp = true
			break
		}
	}

	if !foundOp {
		return nil, fmt.Errorf("invalid filter expression: no operator found in %q", expr)
	}

	// Parse column name or index
	if idx, err := strconv.Atoi(leftPart); err == nil {
		// It's a column index
		filter.ColumnIndex = idx
		if header != nil && idx < len(header) {
			filter.ColumnName = header[idx]
		} else {
			filter.ColumnName = fmt.Sprintf("col_%d", idx)
		}
	} else {
		// It's a column name - find its index
		if header == nil {
			return nil, fmt.Errorf("column name %q used but no header provided", leftPart)
		}
		found := false
		for i, name := range header {
			if strings.EqualFold(name, leftPart) {
				filter.ColumnIndex = i
				filter.ColumnName = name
				found = true
				break
			}
		}
		if !found {
			return nil, fmt.Errorf("column %q not found in header", leftPart)
		}
	}

	// Parse value (remove quotes if present)
	filter.Value = rightPart
	if (strings.HasPrefix(rightPart, "'") && strings.HasSuffix(rightPart, "'")) ||
		(strings.HasPrefix(rightPart, "\"") && strings.HasSuffix(rightPart, "\"")) {
		filter.Value = rightPart[1 : len(rightPart)-1]
	}

	// Try to parse as numeric for numeric operators
	if filter.Operator >= OpEqual && filter.Operator <= OpLessThanOrEqual {
		if numVal, err := strconv.ParseFloat(filter.Value, 64); err == nil {
			filter.IsNumeric = true
			filter.NumValue = numVal
		}
	}

	// Compile regex pattern if needed
	if filter.Operator == OpRegex {
		pat, err := regexp.Compile(filter.Value)
		if err != nil {
			return nil, fmt.Errorf("invalid regex pattern %q: %v", filter.Value, err)
		}
		filter.RegexPat = pat
	}

	return &filter, nil
}

// Evaluate checks if a record matches the filter
func (f *Filter) Evaluate(record []string) (bool, error) {
	if f.ColumnIndex >= len(record) {
		return false, fmt.Errorf("column index %d out of range (record has %d columns)", f.ColumnIndex, len(record))
	}

	cellValue := record[f.ColumnIndex]

	switch f.Operator {
	case OpEqual:
		if f.IsNumeric {
			numVal, err := strconv.ParseFloat(cellValue, 64)
			if err != nil {
				return false, nil // Can't parse as number, doesn't match
			}
			return numVal == f.NumValue, nil
		}
		return cellValue == f.Value, nil

	case OpNotEqual:
		if f.IsNumeric {
			numVal, err := strconv.ParseFloat(cellValue, 64)
			if err != nil {
				return true, nil // Can't parse as number, considered not equal
			}
			return numVal != f.NumValue, nil
		}
		return cellValue != f.Value, nil

	case OpGreaterThan:
		numVal, err := strconv.ParseFloat(cellValue, 64)
		if err != nil {
			return false, nil
		}
		return numVal > f.NumValue, nil

	case OpGreaterThanOrEqual:
		numVal, err := strconv.ParseFloat(cellValue, 64)
		if err != nil {
			return false, nil
		}
		return numVal >= f.NumValue, nil

	case OpLessThan:
		numVal, err := strconv.ParseFloat(cellValue, 64)
		if err != nil {
			return false, nil
		}
		return numVal < f.NumValue, nil

	case OpLessThanOrEqual:
		numVal, err := strconv.ParseFloat(cellValue, 64)
		if err != nil {
			return false, nil
		}
		return numVal <= f.NumValue, nil

	case OpContains:
		return strings.Contains(cellValue, f.Value), nil

	case OpStartsWith:
		return strings.HasPrefix(cellValue, f.Value), nil

	case OpEndsWith:
		return strings.HasSuffix(cellValue, f.Value), nil

	case OpRegex:
		return f.RegexPat.MatchString(cellValue), nil

	default:
		return false, fmt.Errorf("unknown operator: %d", f.Operator)
	}
}

// String returns a human-readable representation of the filter
func (f *Filter) String() string {
	opStr := ""
	switch f.Operator {
	case OpEqual:
		opStr = "="
	case OpNotEqual:
		opStr = "!="
	case OpGreaterThan:
		opStr = ">"
	case OpGreaterThanOrEqual:
		opStr = ">="
	case OpLessThan:
		opStr = "<"
	case OpLessThanOrEqual:
		opStr = "<="
	case OpContains:
		opStr = "contains"
	case OpStartsWith:
		opStr = "startswith"
	case OpEndsWith:
		opStr = "endswith"
	case OpRegex:
		opStr = "~="
	}
	return fmt.Sprintf("%s %s %q", f.ColumnName, opStr, f.Value)
}

// NewFilterSet creates a new filter set from multiple filter expressions
func NewFilterSet(filterExprs []string, header []string) (*FilterSet, error) {
	fs := &FilterSet{
		Filters: make([]*Filter, 0, len(filterExprs)),
	}

	for _, expr := range filterExprs {
		filter, err := ParseFilter(expr, header)
		if err != nil {
			return nil, fmt.Errorf("failed to parse filter %q: %v", expr, err)
		}
		fs.Filters = append(fs.Filters, filter)
	}

	return fs, nil
}

// Evaluate checks if a record matches all filters (AND logic)
func (fs *FilterSet) Evaluate(record []string) (bool, error) {
	for _, filter := range fs.Filters {
		match, err := filter.Evaluate(record)
		if err != nil {
			return false, err
		}
		if !match {
			return false, nil // Short-circuit: if any filter fails, return false
		}
	}
	return true, nil
}

// String returns a human-readable representation of the filter set
func (fs *FilterSet) String() string {
	if len(fs.Filters) == 0 {
		return "no filters"
	}
	var parts []string
	for _, f := range fs.Filters {
		parts = append(parts, f.String())
	}
	return strings.Join(parts, " AND ")
}
