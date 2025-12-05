package largedataset

import (
	"fmt"
	"strconv"
)

// LogicalRow represents a cleaned / typed version of a CSV row.
// For now we only care about a single numeric column: amount.
type LogicalRow struct {
	RawRecord []string // keep the original record if needed
	Amount    float64  // parsed numeric column
	GroupKey  string   // optional group-by key
}

func ParseLogicalRow(record []string) (*LogicalRow, error) {
	return ParseLogicalRowWithGroupBy(record, -1)
}

func ParseLogicalRowWithGroupBy(record []string, groupByIndex int) (*LogicalRow, error) {
	const amountIndex = 8

	if len(record) <= amountIndex {
		return nil, fmt.Errorf("not enough columns, expected index %d", amountIndex)
	}

	rawAmount := record[amountIndex]

	// Try to parse the amount as float64
	amount, err := strconv.ParseFloat(rawAmount, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid amount %q: %w", rawAmount, err)
	}

	logicalRow := &LogicalRow{
		RawRecord: record,
		Amount:    amount,
	}

	// Extract group-by key if enabled
	if groupByIndex >= 0 && groupByIndex < len(record) {
		logicalRow.GroupKey = record[groupByIndex]
	}

	return logicalRow, nil
}
