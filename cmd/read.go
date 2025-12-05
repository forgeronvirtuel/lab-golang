package cmd

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var readCmd = &cobra.Command{
	Use:   "read",
	Short: "Read and display CSV file contents",
	Long: `Read a CSV file and display its contents.
You can specify a custom separator and show the first N rows for debugging.`,
	Example: `  lab-golang read --file data.csv
  lab-golang read --file data.csv --sep ";" --show-first 10`,
	Run: func(cmd *cobra.Command, args []string) {
		if filePath == "" {
			log.Fatal("you must provide --file path to a CSV file")
		}

		if len(separator) != 1 {
			log.Fatal("separator must be a single character")
		}

		readCSV(filePath, rune(separator[0]), showFirst)
	},
}

func init() {
	rootCmd.AddCommand(readCmd)

	readCmd.Flags().StringVarP(&filePath, "file", "f", "", "Path to the CSV file (required)")
	readCmd.Flags().StringVarP(&separator, "sep", "s", ",", "CSV separator (single character)")
	readCmd.Flags().IntVar(&showFirst, "show-first", 5, "Show first N rows for debugging (0 to disable)")
	readCmd.Flags().BoolVar(&hasHeader, "has-header", false, "Specify if the CSV file has a header row")

	readCmd.MarkFlagRequired("file")
}

func readCSV(path string, sep rune, showN int) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()

	// Wrap in a buffered reader for efficient streaming
	bufferedReader := bufio.NewReader(file)

	reader := csv.NewReader(bufferedReader)
	reader.Comma = sep
	reader.LazyQuotes = true
	reader.TrimLeadingSpace = true

	var totalRows int

	for {
		record, err := reader.Read()
		if err != nil {
			if err.Error() != "EOF" {
				log.Fatalf("Error reading CSV: %v", err)
			}
			break
		}
		totalRows++

		if showN > 0 && totalRows <= showN {
			fmt.Printf("Row %d: %v\n", totalRows, record)
		}
	}

	fmt.Printf("Total rows read: %d\n", totalRows)
}
