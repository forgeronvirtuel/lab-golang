package largedataset

import "math"

// AmountStats holds simple streaming statistics for a numeric field.
type AmountStats struct {
	Count int64
	Sum   float64
	Min   float64
	Max   float64
}

// NewAmountStats creates a new AmountStats with proper initial values.
func NewAmountStats() *AmountStats {
	return &AmountStats{
		Count: 0,
		Sum:   0,
		Min:   math.Inf(1),  // +Inf so that any real value will be smaller
		Max:   math.Inf(-1), // -Inf so that any real value will be larger
	}
}

// Add updates the stats with a new logical row.
func (s *AmountStats) Add(row *LogicalRow) {
	v := row.Amount

	if s.Count == 0 {
		// First value, we can also directly set min/max to v
		s.Min = v
		s.Max = v
	} else {
		if v < s.Min {
			s.Min = v
		}
		if v > s.Max {
			s.Max = v
		}
	}

	s.Count++
	s.Sum += v
}

// HasData returns true if at least one value has been added.
func (s *AmountStats) HasData() bool {
	return s.Count > 0
}

// Average returns the average amount or NaN if there is no data.
func (s *AmountStats) Average() float64 {
	if s.Count == 0 {
		return math.NaN()
	}
	return s.Sum / float64(s.Count)
}
