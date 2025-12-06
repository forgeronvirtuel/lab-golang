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
	Long: `Read a CSV file and display its contents with statistics.
You can specify a custom separator, show the first N rows, enable group-by, validate schema, and apply filters.`,
	Example: `  lab-golang parse --file data.csv
  lab-golang parse --file data.csv --sep ";" --show-first 10 --has-header
  lab-golang parse --file data.csv --has-header --group-by 2
  lab-golang parse --file data.csv --has-header --validate
  lab-golang parse --file data.csv --has-header --filter "Price > 100" --filter "Symbol = 'AAPL'"`,
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
			Filters:    filters,
		}

		start := time.Now()

		// Initialize schema if validation is enabled
		var schema *largedataset.CSVSchema
		if validate {
			schema = largedataset.NewStockDataSchema()
			fmt.Printf("Schema validation: ENABLED\n")
		}

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

		if err := processCSV(cfg, schema, composite); err != nil {
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
	parseCmd.Flags().BoolVar(&validate, "validate", false, "Enable CSV schema validation for stock market data")
	parseCmd.Flags().StringArrayVar(&filters, "filter", []string{}, "Filter expression (can be specified multiple times, e.g., --filter \"Price > 100\")")

	parseCmd.MarkFlagRequired("file")
}

type ProcessConfig struct {
	Path       string
	Sep        rune
	HasHeader  bool
	ShowFirst  int
	GroupByCol int      // -1 means no group-by
	Filters    []string // Filter expressions
}

// processCSV opens the file, streams CSV rows, parses them into LogicalRow,
// and counts valid / invalid rows.
func processCSV(cfg ProcessConfig, schema *largedataset.CSVSchema, composite largedataset.Aggregator) error {
	f, err := os.Open(cfg.Path)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	reader := bufio.NewReader(f)
	r := csv.NewReader(reader)
	r.Comma = cfg.Sep

	var (
		totalRows        int
		validRows        int
		invalidRows      int
		validationErrors int
		filteredRows     int
		groupByEnabled   = cfg.GroupByCol >= 0
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

	// Parse filters if any
	var filterSet *largedataset.FilterSet
	if len(cfg.Filters) > 0 {
		filterSet, err = largedataset.NewFilterSet(cfg.Filters, header)
		if err != nil {
			return fmt.Errorf("failed to parse filters: %w", err)
		}
		fmt.Printf("Filters: %s\n", filterSet.String())
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

		// Validate schema if enabled
		if schema != nil {
			if err := schema.ValidateRecord(record); err != nil {
				validationErrors++
				log.Printf("row %d validation error: %v", totalRows, err)
				continue
			}
		}

		// Apply filters if any
		if filterSet != nil {
			match, err := filterSet.Evaluate(record)
			if err != nil {
				log.Printf("row %d filter error: %v", totalRows, err)
				continue
			}
			if !match {
				filteredRows++
				continue // Skip this row
			}
		}

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
	if len(cfg.Filters) > 0 {
		fmt.Printf("Filtered out:       %d\n", filteredRows)
	}
	fmt.Printf("Valid logical rows: %d\n", validRows)
	if schema != nil {
		fmt.Printf("Validation errors:  %d\n", validationErrors)
	} else {
		fmt.Printf("Invalid rows:       %d\n", invalidRows)
	}

	// Detailed reports
	composite.Report(os.Stdout)

	return nil
}
