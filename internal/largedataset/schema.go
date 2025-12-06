package largedataset

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

// ColumnType represents the expected data type of a CSV column
type ColumnType int

const (
	TypeString ColumnType = iota
	TypeInt
	TypeFloat
	TypeBool
	TypeDate
	TypeDateTime
	TypeEmail
	TypeRegex
)

// ColumnDef defines the schema for a single CSV column
type ColumnDef struct {
	Index       int        // Column index (0-based)
	Name        string     // Column name (for documentation/errors)
	Type        ColumnType // Expected data type
	Required    bool       // Whether the column must be non-empty
	MinLength   int        // Minimum string length (for TypeString)
	MaxLength   int        // Maximum string length (for TypeString)
	Min         *float64   // Minimum value (for TypeInt/TypeFloat)
	Max         *float64   // Maximum value (for TypeInt/TypeFloat)
	Pattern     string     // Regex pattern (for TypeRegex)
	DateFormat  string     // Date format (for TypeDate/TypeDateTime)
	AllowedVals []string   // Whitelist of allowed values
}

// CSVSchema represents the complete schema for a CSV file
type CSVSchema struct {
	Columns       []ColumnDef
	MinColumns    int  // Minimum number of columns required
	StrictColumns bool // If true, reject rows with extra columns
}

// ValidationError represents a validation error for a specific column
type ValidationError struct {
	Column   int
	ColName  string
	Value    string
	Expected string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("column %d (%s): invalid value %q - %s", e.Column, e.ColName, e.Value, e.Expected)
}

// ValidateRecord validates a CSV record against the schema
func (s *CSVSchema) ValidateRecord(record []string) error {
	// Check minimum column count
	if len(record) < s.MinColumns {
		return fmt.Errorf("expected at least %d columns, got %d", s.MinColumns, len(record))
	}

	// Check strict column count
	if s.StrictColumns && len(s.Columns) > 0 {
		maxIndex := 0
		for _, col := range s.Columns {
			if col.Index > maxIndex {
				maxIndex = col.Index
			}
		}
		if len(record) > maxIndex+1 {
			return fmt.Errorf("expected exactly %d columns, got %d", maxIndex+1, len(record))
		}
	}

	// Validate each defined column
	for _, colDef := range s.Columns {
		if colDef.Index >= len(record) {
			if colDef.Required {
				return &ValidationError{
					Column:   colDef.Index,
					ColName:  colDef.Name,
					Value:    "",
					Expected: "column is required but missing",
				}
			}
			continue
		}

		value := record[colDef.Index]

		// Check if required
		if colDef.Required && value == "" {
			return &ValidationError{
				Column:   colDef.Index,
				ColName:  colDef.Name,
				Value:    value,
				Expected: "non-empty value required",
			}
		}

		// Skip validation for empty optional fields
		if !colDef.Required && value == "" {
			continue
		}

		// Validate based on type
		if err := validateColumnValue(value, &colDef); err != nil {
			return &ValidationError{
				Column:   colDef.Index,
				ColName:  colDef.Name,
				Value:    value,
				Expected: err.Error(),
			}
		}
	}

	return nil
}

// validateColumnValue validates a single column value against its definition
func validateColumnValue(value string, colDef *ColumnDef) error {
	// Check allowed values whitelist
	if len(colDef.AllowedVals) > 0 {
		found := false
		for _, allowed := range colDef.AllowedVals {
			if value == allowed {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("value must be one of: %v", colDef.AllowedVals)
		}
	}

	// Type-specific validation
	switch colDef.Type {
	case TypeString:
		if colDef.MinLength > 0 && len(value) < colDef.MinLength {
			return fmt.Errorf("string length must be at least %d", colDef.MinLength)
		}
		if colDef.MaxLength > 0 && len(value) > colDef.MaxLength {
			return fmt.Errorf("string length must be at most %d", colDef.MaxLength)
		}

	case TypeInt:
		intVal, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fmt.Errorf("expected integer")
		}
		if colDef.Min != nil && float64(intVal) < *colDef.Min {
			return fmt.Errorf("value must be >= %.0f", *colDef.Min)
		}
		if colDef.Max != nil && float64(intVal) > *colDef.Max {
			return fmt.Errorf("value must be <= %.0f", *colDef.Max)
		}

	case TypeFloat:
		floatVal, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("expected float")
		}
		if colDef.Min != nil && floatVal < *colDef.Min {
			return fmt.Errorf("value must be >= %.2f", *colDef.Min)
		}
		if colDef.Max != nil && floatVal > *colDef.Max {
			return fmt.Errorf("value must be <= %.2f", *colDef.Max)
		}

	case TypeBool:
		lowerVal := value
		if lowerVal != "true" && lowerVal != "false" && lowerVal != "1" && lowerVal != "0" &&
			lowerVal != "yes" && lowerVal != "no" && lowerVal != "y" && lowerVal != "n" {
			return fmt.Errorf("expected boolean (true/false, 1/0, yes/no, y/n)")
		}

	case TypeDate:
		dateFormat := colDef.DateFormat
		if dateFormat == "" {
			dateFormat = "2006-01-02"
		}
		if _, err := time.Parse(dateFormat, value); err != nil {
			return fmt.Errorf("expected date in format %s", dateFormat)
		}

	case TypeDateTime:
		dateFormat := colDef.DateFormat
		if dateFormat == "" {
			dateFormat = "2006-01-02 15:04:05"
		}
		if _, err := time.Parse(dateFormat, value); err != nil {
			return fmt.Errorf("expected datetime in format %s", dateFormat)
		}

	case TypeEmail:
		emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
		if !emailRegex.MatchString(value) {
			return fmt.Errorf("expected valid email address")
		}

	case TypeRegex:
		if colDef.Pattern == "" {
			return fmt.Errorf("regex pattern not specified")
		}
		matched, err := regexp.MatchString(colDef.Pattern, value)
		if err != nil {
			return fmt.Errorf("invalid regex pattern: %v", err)
		}
		if !matched {
			return fmt.Errorf("value must match pattern: %s", colDef.Pattern)
		}
	}

	return nil
}

// NewStockDataSchema creates a schema for the generated stock market CSV
func NewStockDataSchema() *CSVSchema {
	minZero := 0.0
	minOne := 1.0

	return &CSVSchema{
		MinColumns:    22,
		StrictColumns: false,
		Columns: []ColumnDef{
			{Index: 0, Name: "TradeID", Type: TypeInt, Required: true, Min: &minOne},
			{Index: 1, Name: "Timestamp", Type: TypeDateTime, Required: true, DateFormat: "2006-01-02 15:04:05"},
			{Index: 2, Name: "Symbol", Type: TypeString, Required: true, MinLength: 1, MaxLength: 10},
			{Index: 3, Name: "Exchange", Type: TypeString, Required: true, AllowedVals: []string{"NYSE", "NASDAQ", "EURONEXT", "LSE", "TSE"}},
			{Index: 4, Name: "Sector", Type: TypeString, Required: true},
			{Index: 5, Name: "TradeType", Type: TypeString, Required: true, AllowedVals: []string{"Buy", "Sell"}},
			{Index: 6, Name: "OrderType", Type: TypeString, Required: true, AllowedVals: []string{"Market", "Limit", "Stop", "Stop-Limit"}},
			{Index: 7, Name: "Quantity", Type: TypeInt, Required: true, Min: &minOne},
			{Index: 8, Name: "Price", Type: TypeFloat, Required: true, Min: &minZero},
			{Index: 9, Name: "TotalValue", Type: TypeFloat, Required: true, Min: &minZero},
			{Index: 10, Name: "OpenPrice", Type: TypeFloat, Required: true, Min: &minZero},
			{Index: 11, Name: "ClosePrice", Type: TypeFloat, Required: true, Min: &minZero},
			{Index: 12, Name: "HighPrice", Type: TypeFloat, Required: true, Min: &minZero},
			{Index: 13, Name: "LowPrice", Type: TypeFloat, Required: true, Min: &minZero},
			{Index: 14, Name: "Volume", Type: TypeInt, Required: true, Min: &minZero},
			{Index: 15, Name: "MarketCap", Type: TypeFloat, Required: true, Min: &minZero},
			{Index: 16, Name: "PERatio", Type: TypeFloat, Required: true, Min: &minZero},
			{Index: 17, Name: "DividendYield", Type: TypeFloat, Required: true, Min: &minZero},
			{Index: 18, Name: "Beta", Type: TypeFloat, Required: true, Min: &minZero},
			{Index: 19, Name: "52WeekHigh", Type: TypeFloat, Required: true, Min: &minZero},
			{Index: 20, Name: "52WeekLow", Type: TypeFloat, Required: true, Min: &minZero},
			{Index: 21, Name: "ChangePercent", Type: TypeFloat, Required: true},
		},
	}
}
