package largedataset

import (
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// ColumnStats holds statistics gathered during schema inference
type ColumnStats struct {
	Index          int
	Name           string
	TotalValues    int
	EmptyValues    int
	UniqueValues   map[string]bool
	NumericValues  int
	IntegerValues  int
	FloatValues    int
	BoolValues     int
	DateValues     int
	DateTimeValues int
	EmailValues    int
	MinNumeric     float64
	MaxNumeric     float64
	MinLength      int
	MaxLength      int
	SampleValues   []string // Keep first few samples
}

// SchemaInferenceConfig configures the schema inference process
type SchemaInferenceConfig struct {
	SampleSize       int     // Number of rows to analyze (0 = all)
	MinConfidence    float64 // Minimum confidence to infer type (0.0-1.0)
	MaxUniqueForEnum int     // Max unique values to consider as enum
	SampleCount      int     // Number of sample values to keep
}

// DefaultInferenceConfig returns default configuration for schema inference
func DefaultInferenceConfig() *SchemaInferenceConfig {
	return &SchemaInferenceConfig{
		SampleSize:       1000,
		MinConfidence:    0.8,
		MaxUniqueForEnum: 20,
		SampleCount:      5,
	}
}

// InferSchemaFromCSV analyzes a CSV file and infers its schema
func InferSchemaFromCSV(reader io.Reader, sep rune, hasHeader bool, config *SchemaInferenceConfig) (*CSVSchema, error) {
	if config == nil {
		config = DefaultInferenceConfig()
	}

	csvReader := csv.NewReader(reader)
	csvReader.Comma = sep
	csvReader.LazyQuotes = true
	csvReader.TrimLeadingSpace = true

	var header []string
	var columnStats []*ColumnStats
	rowCount := 0

	// Read header if present
	if hasHeader {
		var err error
		header, err = csvReader.Read()
		if err != nil {
			return nil, fmt.Errorf("failed to read header: %w", err)
		}

		// Initialize column stats
		for i, name := range header {
			columnStats = append(columnStats, &ColumnStats{
				Index:        i,
				Name:         name,
				UniqueValues: make(map[string]bool),
				MinNumeric:   math.Inf(1),
				MaxNumeric:   math.Inf(-1),
				MinLength:    math.MaxInt32,
				MaxLength:    0,
				SampleValues: make([]string, 0, config.SampleCount),
			})
		}
	}

	// Read and analyze rows
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue // Skip malformed rows
		}

		rowCount++

		// Initialize column stats on first row if no header
		if !hasHeader && len(columnStats) == 0 {
			for i := range record {
				columnStats = append(columnStats, &ColumnStats{
					Index:        i,
					Name:         fmt.Sprintf("Column%d", i),
					UniqueValues: make(map[string]bool),
					MinNumeric:   math.Inf(1),
					MaxNumeric:   math.Inf(-1),
					MinLength:    math.MaxInt32,
					MaxLength:    0,
					SampleValues: make([]string, 0, config.SampleCount),
				})
			}
		}

		// Analyze each column
		for i, value := range record {
			if i >= len(columnStats) {
				break // Skip extra columns
			}
			analyzeValue(columnStats[i], value, config)
		}

		// Stop if we've reached sample size
		if config.SampleSize > 0 && rowCount >= config.SampleSize {
			break
		}
	}

	// Build schema from statistics
	schema := &CSVSchema{
		MinColumns:    len(columnStats),
		StrictColumns: false,
		Columns:       make([]ColumnDef, 0, len(columnStats)),
	}

	for _, stats := range columnStats {
		colDef := inferColumnDef(stats, config)
		schema.Columns = append(schema.Columns, colDef)
	}

	return schema, nil
}

// analyzeValue updates column statistics with a single value
func analyzeValue(stats *ColumnStats, value string, config *SchemaInferenceConfig) {
	stats.TotalValues++

	if value == "" {
		stats.EmptyValues++
		return
	}

	// Track unique values (up to a limit)
	if len(stats.UniqueValues) < config.MaxUniqueForEnum*2 {
		stats.UniqueValues[value] = true
	}

	// Track string length
	length := len(value)
	if length < stats.MinLength {
		stats.MinLength = length
	}
	if length > stats.MaxLength {
		stats.MaxLength = length
	}

	// Sample values
	if len(stats.SampleValues) < config.SampleCount {
		stats.SampleValues = append(stats.SampleValues, value)
	}

	// Try to parse as numeric
	if numVal, err := strconv.ParseFloat(value, 64); err == nil {
		stats.NumericValues++
		if numVal < stats.MinNumeric {
			stats.MinNumeric = numVal
		}
		if numVal > stats.MaxNumeric {
			stats.MaxNumeric = numVal
		}

		// Check if it's an integer
		if _, err := strconv.ParseInt(value, 10, 64); err == nil {
			stats.IntegerValues++
		} else {
			stats.FloatValues++
		}
	}

	// Check if it's a boolean
	lowerVal := strings.ToLower(strings.TrimSpace(value))
	if lowerVal == "true" || lowerVal == "false" || lowerVal == "1" || lowerVal == "0" ||
		lowerVal == "yes" || lowerVal == "no" || lowerVal == "y" || lowerVal == "n" {
		stats.BoolValues++
	}

	// Check if it's a date/datetime
	if isDateTime(value) {
		stats.DateTimeValues++
	} else if isDate(value) {
		stats.DateValues++
	}

	// Check if it's an email
	if isEmail(value) {
		stats.EmailValues++
	}
}

// inferColumnDef infers a ColumnDef from gathered statistics
func inferColumnDef(stats *ColumnStats, config *SchemaInferenceConfig) ColumnDef {
	colDef := ColumnDef{
		Index: stats.Index,
		Name:  stats.Name,
	}

	nonEmptyCount := stats.TotalValues - stats.EmptyValues
	if nonEmptyCount == 0 {
		// All empty - default to string
		colDef.Type = TypeString
		colDef.Required = false
		return colDef
	}

	// Calculate confidence for each type
	numericConfidence := float64(stats.NumericValues) / float64(nonEmptyCount)
	integerConfidence := float64(stats.IntegerValues) / float64(nonEmptyCount)
	boolConfidence := float64(stats.BoolValues) / float64(nonEmptyCount)
	dateConfidence := float64(stats.DateValues) / float64(nonEmptyCount)
	dateTimeConfidence := float64(stats.DateTimeValues) / float64(nonEmptyCount)
	emailConfidence := float64(stats.EmailValues) / float64(nonEmptyCount)

	// Determine type based on highest confidence above threshold
	switch {
	case dateTimeConfidence >= config.MinConfidence:
		colDef.Type = TypeDateTime
		colDef.DateFormat = detectDateTimeFormat(stats.SampleValues)

	case dateConfidence >= config.MinConfidence:
		colDef.Type = TypeDate
		colDef.DateFormat = detectDateFormat(stats.SampleValues)

	case emailConfidence >= config.MinConfidence:
		colDef.Type = TypeEmail

	case boolConfidence >= config.MinConfidence:
		colDef.Type = TypeBool

	case integerConfidence >= config.MinConfidence:
		colDef.Type = TypeInt
		if !math.IsInf(stats.MinNumeric, 0) {
			min := stats.MinNumeric
			colDef.Min = &min
		}
		if !math.IsInf(stats.MaxNumeric, 0) {
			max := stats.MaxNumeric
			colDef.Max = &max
		}

	case numericConfidence >= config.MinConfidence:
		colDef.Type = TypeFloat
		if !math.IsInf(stats.MinNumeric, 0) {
			min := stats.MinNumeric
			colDef.Min = &min
		}
		if !math.IsInf(stats.MaxNumeric, 0) {
			max := stats.MaxNumeric
			colDef.Max = &max
		}

	default:
		colDef.Type = TypeString
		if stats.MinLength < math.MaxInt32 && stats.MinLength > 0 {
			colDef.MinLength = stats.MinLength
		}
		if stats.MaxLength > 0 {
			colDef.MaxLength = stats.MaxLength
		}
	}

	// Check if it should be an enum (limited unique values)
	uniqueCount := len(stats.UniqueValues)
	if colDef.Type == TypeString && uniqueCount > 0 && uniqueCount <= config.MaxUniqueForEnum {
		colDef.AllowedVals = make([]string, 0, uniqueCount)
		for val := range stats.UniqueValues {
			if val != "" {
				colDef.AllowedVals = append(colDef.AllowedVals, val)
			}
		}
	}

	// Determine if required
	emptyRatio := float64(stats.EmptyValues) / float64(stats.TotalValues)
	colDef.Required = emptyRatio < 0.05 // Less than 5% empty

	return colDef
}

// Helper functions for type detection
func isDate(value string) bool {
	dateFormats := []string{
		"2006-01-02",
		"01/02/2006",
		"02/01/2006",
		"2006/01/02",
		"01-02-2006",
		"02-01-2006",
	}
	for _, format := range dateFormats {
		if _, err := time.Parse(format, value); err == nil {
			return true
		}
	}
	return false
}

func isDateTime(value string) bool {
	dateTimeFormats := []string{
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05Z",
		"2006-01-02T15:04:05Z",
		"01/02/2006 15:04:05",
		"02/01/2006 15:04:05",
	}
	for _, format := range dateTimeFormats {
		if _, err := time.Parse(format, value); err == nil {
			return true
		}
	}
	return false
}

func isEmail(value string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(value)
}

func detectDateFormat(samples []string) string {
	formats := []string{
		"2006-01-02",
		"01/02/2006",
		"02/01/2006",
		"2006/01/02",
	}
	for _, format := range formats {
		for _, sample := range samples {
			if _, err := time.Parse(format, sample); err == nil {
				return format
			}
		}
	}
	return "2006-01-02" // default
}

func detectDateTimeFormat(samples []string) string {
	formats := []string{
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05Z",
		"01/02/2006 15:04:05",
	}
	for _, format := range formats {
		for _, sample := range samples {
			if _, err := time.Parse(format, sample); err == nil {
				return format
			}
		}
	}
	return "2006-01-02 15:04:05" // default
}

// PrintSchemaAsCode generates Go code for the inferred schema
func PrintSchemaAsCode(schema *CSVSchema, schemaName string) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("func New%sSchema() *CSVSchema {\n", schemaName))

	// Find if we need min/max variables and create mapping
	minMaxVars := make(map[float64]string)
	varCount := 0

	for _, col := range schema.Columns {
		if col.Min != nil {
			if _, exists := minMaxVars[*col.Min]; !exists {
				varName := fmt.Sprintf("v%d", varCount)
				minMaxVars[*col.Min] = varName
				varCount++
			}
		}
		if col.Max != nil {
			if _, exists := minMaxVars[*col.Max]; !exists {
				varName := fmt.Sprintf("v%d", varCount)
				minMaxVars[*col.Max] = varName
				varCount++
			}
		}
	}

	if len(minMaxVars) > 0 {
		sb.WriteString("\t// Define min/max constraints\n")
		for val, varName := range minMaxVars {
			sb.WriteString(fmt.Sprintf("\t%s := %.2f\n", varName, val))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("\treturn &CSVSchema{\n")
	sb.WriteString(fmt.Sprintf("\t\tMinColumns:    %d,\n", schema.MinColumns))
	sb.WriteString(fmt.Sprintf("\t\tStrictColumns: %v,\n", schema.StrictColumns))
	sb.WriteString("\t\tColumns: []ColumnDef{\n")

	for _, col := range schema.Columns {
		sb.WriteString("\t\t\t{")
		sb.WriteString(fmt.Sprintf("Index: %d, ", col.Index))
		sb.WriteString(fmt.Sprintf("Name: %q, ", col.Name))
		sb.WriteString(fmt.Sprintf("Type: %s, ", getTypeName(col.Type)))
		sb.WriteString(fmt.Sprintf("Required: %v", col.Required))

		if col.MinLength > 0 {
			sb.WriteString(fmt.Sprintf(", MinLength: %d", col.MinLength))
		}
		if col.MaxLength > 0 {
			sb.WriteString(fmt.Sprintf(", MaxLength: %d", col.MaxLength))
		}
		if col.Min != nil {
			if varName, exists := minMaxVars[*col.Min]; exists {
				sb.WriteString(fmt.Sprintf(", Min: &%s", varName))
			}
		}
		if col.Max != nil {
			if varName, exists := minMaxVars[*col.Max]; exists {
				sb.WriteString(fmt.Sprintf(", Max: &%s", varName))
			}
		}
		if col.DateFormat != "" {
			sb.WriteString(fmt.Sprintf(", DateFormat: %q", col.DateFormat))
		}
		if len(col.AllowedVals) > 0 {
			sb.WriteString(", AllowedVals: []string{")
			for i, val := range col.AllowedVals {
				if i > 0 {
					sb.WriteString(", ")
				}
				sb.WriteString(fmt.Sprintf("%q", val))
			}
			sb.WriteString("}")
		}

		sb.WriteString("},\n")
	}

	sb.WriteString("\t\t},\n")
	sb.WriteString("\t}\n")
	sb.WriteString("}\n")

	return sb.String()
}

func getTypeName(t ColumnType) string {
	switch t {
	case TypeString:
		return "TypeString"
	case TypeInt:
		return "TypeInt"
	case TypeFloat:
		return "TypeFloat"
	case TypeBool:
		return "TypeBool"
	case TypeDate:
		return "TypeDate"
	case TypeDateTime:
		return "TypeDateTime"
	case TypeEmail:
		return "TypeEmail"
	case TypeRegex:
		return "TypeRegex"
	default:
		return "TypeString"
	}
}
