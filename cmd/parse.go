package cmd

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/forgeronvirtuel/lab-golang/internal/largedataset"
	"github.com/spf13/cobra"
)

var parseCmd = &cobra.Command{
	Use:   "parse",
	Short: "Read and display CSV file contents",
	Long: `Read a CSV file and display its contents.
You can specify a custom separator and show the first N rows for debugging.`,
	Example: `  lab-golang parse --file data.csv
  lab-golang parse --file data.csv --sep ";" --show-first 10 --has-header`,
	Run: func(cmd *cobra.Command, args []string) {
		if filePath == "" {
			log.Fatal("you must provide --file path to a CSV file")
		}

		if len(separator) != 1 {
			log.Fatal("separator must be a single character")
		}

		cfg := ProcessConfig{
			Path:       filePath,
			Sep:        rune(separator[0]),
			HasHeader:  hasHeader,
			ShowFirst:  showFirst,
			GroupByCol: groupByCol,
		}

		start := time.Now()

		// Build aggregators
		globalAgg := largedataset.NewGlobalAmountAggregator()
		var aggs []largedataset.Aggregator
		aggs = append(aggs, globalAgg)

		if cfg.GroupByCol >= 0 {
			aggs = append(aggs, largedataset.NewGroupByAggregator())
		}

		if cfg.ShowFirst > 0 {
			aggs = append(aggs, largedataset.NewDebugAggregator(cfg.ShowFirst))
		}

		composite := largedataset.NewCompositeAggregator(aggs...)

		if err := processCSV(cfg, composite); err != nil {
			log.Fatalf("Error processing CSV: %v", err)
		}
		elapsed := time.Since(start)
		fmt.Printf("\nProcessing took %s\n", elapsed)
	},
}

func init() {
	rootCmd.AddCommand(parseCmd)

	parseCmd.Flags().StringVarP(&filePath, "file", "f", "", "Path to the CSV file (required)")
	parseCmd.Flags().StringVarP(&separator, "sep", "s", ",", "CSV separator (single character)")
	parseCmd.Flags().IntVar(&showFirst, "show-first", 5, "Show first N rows for debugging (0 to disable)")
	parseCmd.Flags().BoolVar(&hasHeader, "has-header", false, "Specify if the CSV file has a header row")
	parseCmd.Flags().IntVar(&groupByCol, "group-by", -1, "Column index (0-based) for group-by statistics (-1 to disable)")

	parseCmd.MarkFlagRequired("file")
}

type ProcessConfig struct {
	Path       string
	Sep        rune
	HasHeader  bool
	ShowFirst  int
	GroupByCol int // -1 means no group-by
}

// processCSV opens the file, streams CSV rows, parses them into LogicalRow,
// and counts valid / invalid rows.
func processCSV(cfg ProcessConfig, composite largedataset.Aggregator) error {
	f, err := os.Open(cfg.Path)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	reader := bufio.NewReader(f)
	r := csv.NewReader(reader)
	r.Comma = cfg.Sep

	var (
		totalRows      int
		validRows      int
		invalidRows    int
		groupByEnabled = cfg.GroupByCol >= 0
	)

	// Optionally read and ignore the header row
	var header []string
	if cfg.HasHeader {
		header, err = r.Read()
		if err != nil {
			if err == io.EOF {
				return fmt.Errorf("file contains only a header and no data rows")
			}
			return fmt.Errorf("failed to read header: %w", err)
		}
		fmt.Printf("Header: %v\n", header)
		if groupByEnabled && groupByCol < len(header) {
			fmt.Printf("Group-by column: %s (index %d)\n", header[groupByCol], groupByCol)
		}
	}

	for {
		record, err := r.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			// At this stage, we consider this a "fatal" CSV error (bad format)
			return fmt.Errorf("error while reading CSV: %w", err)
		}

		totalRows++

		logical, err := largedataset.ParseLogicalRowWithGroupBy(record, groupByCol)
		if err != nil {
			invalidRows++
			// For now, we just log the error and continue.
			// Later, we can introduce --limit-errors, etc.
			log.Printf("skipping invalid row %d: %v", totalRows, err)
			continue
		}

		validRows++
		composite.Consume(logical)
	}

	fmt.Printf("\n=== Summary ===\n")
	fmt.Printf("Total rows read:    %d\n", totalRows)
	fmt.Printf("Valid logical rows: %d\n", validRows)
	fmt.Printf("Invalid rows:       %d\n", invalidRows)

	// Detailed reports
	composite.Report(os.Stdout)

	return nil
}
