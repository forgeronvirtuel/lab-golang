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

// LogicalOperator represents AND or OR
type LogicalOperator int

const (
	LogicalAND LogicalOperator = iota
	LogicalOR
)

// FilterExpression represents a filter expression tree
type FilterExpression interface {
	Evaluate(record []string) (bool, error)
	String() string
}

// FilterCondition is a leaf node (single filter)
type FilterCondition struct {
	Filter *Filter
}

// FilterGroup is a composite node (multiple expressions with AND/OR)
type FilterGroup struct {
	Operator    LogicalOperator
	Expressions []FilterExpression
}

// FilterSet represents a collection of filters with support for AND/OR logic
type FilterSet struct {
	Filters []*Filter        // Deprecated: kept for backward compatibility
	Root    FilterExpression // Root of the expression tree
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

// Evaluate for FilterCondition
func (fc *FilterCondition) Evaluate(record []string) (bool, error) {
	return fc.Filter.Evaluate(record)
}

// String for FilterCondition
func (fc *FilterCondition) String() string {
	return fc.Filter.String()
}

// Evaluate for FilterGroup
func (fg *FilterGroup) Evaluate(record []string) (bool, error) {
	if len(fg.Expressions) == 0 {
		return true, nil
	}

	if fg.Operator == LogicalAND {
		// All expressions must be true
		for _, expr := range fg.Expressions {
			match, err := expr.Evaluate(record)
			if err != nil {
				return false, err
			}
			if !match {
				return false, nil // Short-circuit on first false
			}
		}
		return true, nil
	} else { // LogicalOR
		// At least one expression must be true
		for _, expr := range fg.Expressions {
			match, err := expr.Evaluate(record)
			if err != nil {
				return false, err
			}
			if match {
				return true, nil // Short-circuit on first true
			}
		}
		return false, nil
	}
}

// String for FilterGroup
func (fg *FilterGroup) String() string {
	if len(fg.Expressions) == 0 {
		return ""
	}
	if len(fg.Expressions) == 1 {
		return fg.Expressions[0].String()
	}

	operator := " AND "
	if fg.Operator == LogicalOR {
		operator = " OR "
	}

	var parts []string
	for _, expr := range fg.Expressions {
		exprStr := expr.String()
		// Add parentheses if it's a group
		if _, isGroup := expr.(*FilterGroup); isGroup {
			exprStr = "(" + exprStr + ")"
		}
		parts = append(parts, exprStr)
	}
	return strings.Join(parts, operator)
}

// NewFilterSet creates a new filter set from multiple filter expressions
// Supports both simple (backward compatible) and complex (with AND/OR) expressions
func NewFilterSet(filterExprs []string, header []string) (*FilterSet, error) {
	fs := &FilterSet{
		Filters: make([]*Filter, 0, len(filterExprs)),
	}

	// Backward compatibility: if multiple expressions, treat as AND by default
	if len(filterExprs) == 0 {
		fs.Root = &FilterGroup{Operator: LogicalAND, Expressions: []FilterExpression{}}
		return fs, nil
	}

	// Parse expressions
	var expressions []FilterExpression
	for _, expr := range filterExprs {
		parsed, err := parseComplexFilter(expr, header)
		if err != nil {
			return nil, fmt.Errorf("failed to parse filter %q: %v", expr, err)
		}
		expressions = append(expressions, parsed)

		// Also populate Filters array for backward compatibility
		if cond, ok := parsed.(*FilterCondition); ok {
			fs.Filters = append(fs.Filters, cond.Filter)
		}
	}

	// If single expression, use it directly
	if len(expressions) == 1 {
		fs.Root = expressions[0]
	} else {
		// Multiple expressions: combine with AND (backward compatible behavior)
		fs.Root = &FilterGroup{
			Operator:    LogicalAND,
			Expressions: expressions,
		}
	}

	return fs, nil
}

// parseComplexFilter parses a filter expression with support for AND/OR/parentheses
func parseComplexFilter(expr string, header []string) (FilterExpression, error) {
	expr = strings.TrimSpace(expr)

	// Check for parentheses
	if strings.Contains(expr, "(") || strings.Contains(expr, ")") {
		return parseGroupedFilter(expr, header)
	}

	// Check for OR operator (lower precedence)
	if orParts := splitByLogicalOperator(expr, " OR "); len(orParts) > 1 {
		var expressions []FilterExpression
		for _, part := range orParts {
			subExpr, err := parseComplexFilter(part, header)
			if err != nil {
				return nil, err
			}
			expressions = append(expressions, subExpr)
		}
		return &FilterGroup{Operator: LogicalOR, Expressions: expressions}, nil
	}

	// Check for AND operator (higher precedence)
	if andParts := splitByLogicalOperator(expr, " AND "); len(andParts) > 1 {
		var expressions []FilterExpression
		for _, part := range andParts {
			subExpr, err := parseComplexFilter(part, header)
			if err != nil {
				return nil, err
			}
			expressions = append(expressions, subExpr)
		}
		return &FilterGroup{Operator: LogicalAND, Expressions: expressions}, nil
	}

	// Single filter condition
	filter, err := ParseFilter(expr, header)
	if err != nil {
		return nil, err
	}
	return &FilterCondition{Filter: filter}, nil
}

// parseGroupedFilter handles expressions with parentheses
func parseGroupedFilter(expr string, header []string) (FilterExpression, error) {
	expr = strings.TrimSpace(expr)

	// Find matching parentheses and split by logical operators
	var tokens []string
	var current strings.Builder
	depth := 0

	for i := 0; i < len(expr); i++ {
		ch := expr[i]
		switch ch {
		case '(':
			depth++
			current.WriteByte(ch)
		case ')':
			depth--
			current.WriteByte(ch)
		default:
			current.WriteByte(ch)
		}

		// Check for logical operators at depth 0
		if depth == 0 {
			rest := expr[i:]
			if strings.HasPrefix(rest, " OR ") {
				tokens = append(tokens, strings.TrimSpace(current.String()))
				tokens = append(tokens, "OR")
				current.Reset()
				i += 3 // Skip " OR"
			} else if strings.HasPrefix(rest, " AND ") {
				tokens = append(tokens, strings.TrimSpace(current.String()))
				tokens = append(tokens, "AND")
				current.Reset()
				i += 4 // Skip " AND"
			}
		}
	}

	if current.Len() > 0 {
		tokens = append(tokens, strings.TrimSpace(current.String()))
	}

	// Parse tokens
	return parseTokens(tokens, header)
}

// parseTokens converts a list of tokens into an expression tree
func parseTokens(tokens []string, header []string) (FilterExpression, error) {
	if len(tokens) == 0 {
		return nil, fmt.Errorf("empty expression")
	}

	if len(tokens) == 1 {
		token := tokens[0]
		// Remove outer parentheses if present
		if strings.HasPrefix(token, "(") && strings.HasSuffix(token, ")") {
			token = strings.TrimSpace(token[1 : len(token)-1])
			return parseComplexFilter(token, header)
		}
		return parseComplexFilter(token, header)
	}

	// Find OR operators first (lower precedence)
	for i, token := range tokens {
		if token == "OR" {
			left, err := parseTokens(tokens[:i], header)
			if err != nil {
				return nil, err
			}
			right, err := parseTokens(tokens[i+1:], header)
			if err != nil {
				return nil, err
			}

			// Merge with existing OR group if possible
			if leftGroup, ok := left.(*FilterGroup); ok && leftGroup.Operator == LogicalOR {
				leftGroup.Expressions = append(leftGroup.Expressions, right)
				return leftGroup, nil
			}

			return &FilterGroup{
				Operator:    LogicalOR,
				Expressions: []FilterExpression{left, right},
			}, nil
		}
	}

	// Find AND operators (higher precedence)
	for i, token := range tokens {
		if token == "AND" {
			left, err := parseTokens(tokens[:i], header)
			if err != nil {
				return nil, err
			}
			right, err := parseTokens(tokens[i+1:], header)
			if err != nil {
				return nil, err
			}

			// Merge with existing AND group if possible
			if leftGroup, ok := left.(*FilterGroup); ok && leftGroup.Operator == LogicalAND {
				leftGroup.Expressions = append(leftGroup.Expressions, right)
				return leftGroup, nil
			}

			return &FilterGroup{
				Operator:    LogicalAND,
				Expressions: []FilterExpression{left, right},
			}, nil
		}
	}

	return nil, fmt.Errorf("invalid token sequence")
}

// splitByLogicalOperator splits an expression by a logical operator (outside parentheses)
func splitByLogicalOperator(expr string, operator string) []string {
	var parts []string
	var current strings.Builder
	depth := 0

	for i := 0; i < len(expr); i++ {
		ch := expr[i]
		if ch == '(' {
			depth++
			current.WriteByte(ch)
		} else if ch == ')' {
			depth--
			current.WriteByte(ch)
		} else if depth == 0 && strings.HasPrefix(expr[i:], operator) {
			parts = append(parts, strings.TrimSpace(current.String()))
			current.Reset()
			i += len(operator) - 1
		} else {
			current.WriteByte(ch)
		}
	}

	if current.Len() > 0 {
		parts = append(parts, strings.TrimSpace(current.String()))
	}

	return parts
}

// Evaluate checks if a record matches the filter expression tree
func (fs *FilterSet) Evaluate(record []string) (bool, error) {
	if fs.Root != nil {
		return fs.Root.Evaluate(record)
	}

	// Backward compatibility: use Filters array with AND logic
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
	if fs.Root != nil {
		return fs.Root.String()
	}

	// Backward compatibility
	if len(fs.Filters) == 0 {
		return "no filters"
	}
	var parts []string
	for _, f := range fs.Filters {
		parts = append(parts, f.String())
	}
	return strings.Join(parts, " AND ")
}
