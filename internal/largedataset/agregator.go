package largedataset

import (
	"fmt"
	"io"
	"sort"
)

type Aggregator interface {
	Consume(row *LogicalRow)
	Report(w io.Writer)
}

type CompositeAggregator struct {
	aggs []Aggregator
}

func NewCompositeAggregator(aggs ...Aggregator) *CompositeAggregator {
	return &CompositeAggregator{aggs: aggs}
}

func (c *CompositeAggregator) Consume(row *LogicalRow) {
	for _, a := range c.aggs {
		a.Consume(row)
	}
}

func (c *CompositeAggregator) Report(w io.Writer) {
	for _, a := range c.aggs {
		a.Report(w)
	}
}

type GlobalAmountAggregator struct {
	Stats *AmountStats
}

func NewGlobalAmountAggregator() *GlobalAmountAggregator {
	stats := NewAmountStats()
	return &GlobalAmountAggregator{Stats: stats}
}

func (g *GlobalAmountAggregator) Consume(row *LogicalRow) {
	g.Stats.Add(row)
}

func (g *GlobalAmountAggregator) Report(w io.Writer) {
	if !g.Stats.HasData() {
		fmt.Fprintln(w, "\nNo valid amount data to compute stats.")
		return
	}

	fmt.Fprintln(w, "\n=== Amount stats (global) ===")
	fmt.Fprintf(w, "Count:   %d\n", g.Stats.Count)
	fmt.Fprintf(w, "Sum:     %.2f\n", g.Stats.Sum)
	fmt.Fprintf(w, "Min:     %.2f\n", g.Stats.Min)
	fmt.Fprintf(w, "Max:     %.2f\n", g.Stats.Max)
	fmt.Fprintf(w, "Average: %.2f\n", g.Stats.Average())
}

type GroupByAggregator struct {
	statsByKey map[string]*AmountStats
}

func NewGroupByAggregator() *GroupByAggregator {
	return &GroupByAggregator{
		statsByKey: make(map[string]*AmountStats),
	}
}

func (g *GroupByAggregator) Consume(row *LogicalRow) {
	if row.GroupKey == "" {
		return
	}

	s, ok := g.statsByKey[row.GroupKey]
	if !ok {
		s = NewAmountStats()
		g.statsByKey[row.GroupKey] = s
	}

	s.Add(row)
}

func (g *GroupByAggregator) Report(w io.Writer) {
	if len(g.statsByKey) == 0 {
		return
	}

	fmt.Fprintln(w, "\n=== Group-by statistics ===")
	fmt.Fprintf(w, "Number of groups: %d\n\n", len(g.statsByKey))

	type groupStat struct {
		key   string
		stats *AmountStats
	}

	groupList := make([]groupStat, 0, len(g.statsByKey))
	for key, stats := range g.statsByKey {
		groupList = append(groupList, groupStat{key: key, stats: stats})
	}

	sort.Slice(groupList, func(i, j int) bool {
		return groupList[i].stats.Sum > groupList[j].stats.Sum
	})

	fmt.Fprintln(w, "Groups sorted by total sum (descending):")
	for i, group := range groupList {
		fmt.Fprintf(w, "[%d] %s\n", i+1, group.key)
		fmt.Fprintf(w, "  Count:   %d\n", group.stats.Count)
		fmt.Fprintf(w, "  Sum:     %.2f\n", group.stats.Sum)
		fmt.Fprintf(w, "  Min:     %.2f\n", group.stats.Min)
		fmt.Fprintf(w, "  Max:     %.2f\n", group.stats.Max)
		fmt.Fprintf(w, "  Average: %.2f\n\n", group.stats.Average())
	}
}

type DebugAggregator struct {
	currentRow int
	maxRows    int
}

func NewDebugAggregator(maxRows int) *DebugAggregator {
	return &DebugAggregator{
		maxRows: maxRows,
	}
}

func (g *DebugAggregator) Consume(row *LogicalRow) {
	if g.currentRow < g.maxRows {
		fmt.Printf("Debug Row %d: %+v\n", g.currentRow+1, row.RawRecord)
		g.currentRow++
	}
}

func (g *DebugAggregator) Report(w io.Writer) {}
