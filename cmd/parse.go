package cmd

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/forgeronvirtuel/lab-golang/internal/largedataset"
	"github.com/spf13/cobra"
)

var ()

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

		processCSV(filePath, rune(separator[0]), hasHeader, showFirst)
	},
}

func init() {
	rootCmd.AddCommand(parseCmd)

	parseCmd.Flags().StringVarP(&filePath, "file", "f", "", "Path to the CSV file (required)")
	parseCmd.Flags().StringVarP(&separator, "sep", "s", ",", "CSV separator (single character)")
	parseCmd.Flags().IntVar(&showFirst, "show-first", 5, "Show first N rows for debugging (0 to disable)")
	parseCmd.Flags().BoolVar(&hasHeader, "has-header", false, "Specify if the CSV file has a header row")

	parseCmd.MarkFlagRequired("file")
}

// processCSV opens the file, streams CSV rows, parses them into LogicalRow,
// and counts valid / invalid rows.
func processCSV(path string, sep rune, hasHeader bool, showFirst int) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	reader := bufio.NewReader(f)
	r := csv.NewReader(reader)
	r.Comma = sep

	var (
		totalRows   int
		validRows   int
		invalidRows int
	)

	// Optionally read and ignore the header row
	if hasHeader {
		header, err := r.Read()
		if err != nil {
			if err == io.EOF {
				return fmt.Errorf("file contains only a header and no data rows")
			}
			return fmt.Errorf("failed to read header: %w", err)
		}
		fmt.Printf("Header: %v\n", header)
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

		logical, err := largedataset.ParseLogicalRow(record)
		if err != nil {
			invalidRows++
			// For now, we just log the error and continue.
			// Later, we can introduce --limit-errors, etc.
			log.Printf("skipping invalid row %d: %v", totalRows, err)
			continue
		}

		validRows++

		// Optionally show first N valid logical rows
		if showFirst > 0 && validRows <= showFirst {
			fmt.Printf("Valid row %d: amount=%.2f raw=%v\n", totalRows, logical.Amount, logical.RawRecord)
		}

		// Later, this is where we will send `logical` to an aggregator.
	}

	fmt.Printf("\n=== Summary ===\n")
	fmt.Printf("Total rows read:    %d\n", totalRows)
	fmt.Printf("Valid logical rows: %d\n", validRows)
	fmt.Printf("Invalid rows:       %d\n", invalidRows)

	return nil
}
