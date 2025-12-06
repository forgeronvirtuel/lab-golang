package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/forgeronvirtuel/lab-golang/internal/largedataset"
	"github.com/spf13/cobra"
)

var (
	inferSampleSize int
	inferConfidence float64
	inferMaxEnum    int
	inferOutputCode bool
	inferSchemaName string
)

var inferCmd = &cobra.Command{
	Use:   "infer",
	Short: "Infer CSV schema from data",
	Long: `Analyze a CSV file and automatically infer its schema including column types,
constraints, and allowed values. The inferred schema can be output as Go code.`,
	Example: `  # Infer schema from first 1000 rows
  lab-golang infer --file data.csv --has-header

  # Generate Go code for the schema
  lab-golang infer --file data.csv --has-header --code --schema-name MyData

  # Analyze all rows with custom confidence threshold
  lab-golang infer --file data.csv --has-header --sample 0 --confidence 0.9`,
	Run: func(cmd *cobra.Command, args []string) {
		if filePath == "" {
			log.Fatal("you must provide --file path to a CSV file")
		}

		if len(separator) != 1 {
			log.Fatal("separator must be a single character")
		}

		inferSchema()
	},
}

func init() {
	rootCmd.AddCommand(inferCmd)

	inferCmd.Flags().StringVarP(&filePath, "file", "f", "", "Path to the CSV file (required)")
	inferCmd.Flags().StringVarP(&separator, "sep", "s", ",", "CSV separator (single character)")
	inferCmd.Flags().BoolVar(&hasHeader, "has-header", false, "Specify if the CSV file has a header row")
	inferCmd.Flags().IntVar(&inferSampleSize, "sample", 1000, "Number of rows to analyze (0 for all)")
	inferCmd.Flags().Float64Var(&inferConfidence, "confidence", 0.8, "Minimum confidence threshold for type inference (0.0-1.0)")
	inferCmd.Flags().IntVar(&inferMaxEnum, "max-enum", 20, "Maximum unique values to consider as enum")
	inferCmd.Flags().BoolVar(&inferOutputCode, "code", false, "Output inferred schema as Go code")
	inferCmd.Flags().StringVar(&inferSchemaName, "schema-name", "Inferred", "Name for the generated schema (used with --code)")

	inferCmd.MarkFlagRequired("file")
}

func inferSchema() {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer f.Close()

	reader := bufio.NewReader(f)

	config := &largedataset.SchemaInferenceConfig{
		SampleSize:       inferSampleSize,
		MinConfidence:    inferConfidence,
		MaxUniqueForEnum: inferMaxEnum,
		SampleCount:      5,
	}

	fmt.Printf("🔍 Inferring schema from CSV file...\n")
	fmt.Printf("Configuration:\n")
	fmt.Printf("  Sample size: %d rows", config.SampleSize)
	if config.SampleSize == 0 {
		fmt.Printf(" (all)")
	}
	fmt.Printf("\n")
	fmt.Printf("  Confidence threshold: %.1f%%\n", config.MinConfidence*100)
	fmt.Printf("  Max enum values: %d\n", config.MaxUniqueForEnum)
	fmt.Printf("  Has header: %v\n", hasHeader)
	fmt.Println()

	schema, err := largedataset.InferSchemaFromCSV(reader, rune(separator[0]), hasHeader, config)
	if err != nil {
		log.Fatalf("Failed to infer schema: %v", err)
	}

	fmt.Printf("✅ Schema inferred successfully!\n\n")
	fmt.Printf("═══════════════════════════════════════════════════════════════\n")
	fmt.Printf("INFERRED SCHEMA\n")
	fmt.Printf("═══════════════════════════════════════════════════════════════\n\n")

	fmt.Printf("Total columns: %d\n", len(schema.Columns))
	fmt.Printf("Minimum columns required: %d\n", schema.MinColumns)
	fmt.Printf("Strict column count: %v\n\n", schema.StrictColumns)

	// Print column details
	for i, col := range schema.Columns {
		fmt.Printf("[%d] %s\n", i, col.Name)
		fmt.Printf("    Type: %s\n", getTypeDescription(col.Type))
		fmt.Printf("    Required: %v\n", col.Required)

		if col.MinLength > 0 || col.MaxLength > 0 {
			fmt.Printf("    Length: ")
			if col.MinLength > 0 {
				fmt.Printf("min=%d ", col.MinLength)
			}
			if col.MaxLength > 0 {
				fmt.Printf("max=%d", col.MaxLength)
			}
			fmt.Println()
		}

		if col.Min != nil || col.Max != nil {
			fmt.Printf("    Range: ")
			if col.Min != nil {
				fmt.Printf("min=%.2f ", *col.Min)
			}
			if col.Max != nil {
				fmt.Printf("max=%.2f", *col.Max)
			}
			fmt.Println()
		}

		if col.DateFormat != "" {
			fmt.Printf("    Date format: %s\n", col.DateFormat)
		}

		if len(col.AllowedVals) > 0 {
			fmt.Printf("    Allowed values (%d): ", len(col.AllowedVals))
			for j, val := range col.AllowedVals {
				if j > 0 {
					fmt.Print(", ")
				}
				if j >= 10 {
					fmt.Printf("... (%d more)", len(col.AllowedVals)-10)
					break
				}
				fmt.Printf("%q", val)
			}
			fmt.Println()
		}

		fmt.Println()
	}

	// Output as Go code if requested
	if inferOutputCode {
		fmt.Printf("═══════════════════════════════════════════════════════════════\n")
		fmt.Printf("GO CODE\n")
		fmt.Printf("═══════════════════════════════════════════════════════════════\n\n")
		code := largedataset.PrintSchemaAsCode(schema, inferSchemaName)
		fmt.Println(code)
	}
}

func getTypeDescription(t largedataset.ColumnType) string {
	switch t {
	case largedataset.TypeString:
		return "String"
	case largedataset.TypeInt:
		return "Integer"
	case largedataset.TypeFloat:
		return "Float"
	case largedataset.TypeBool:
		return "Boolean"
	case largedataset.TypeDate:
		return "Date"
	case largedataset.TypeDateTime:
		return "DateTime"
	case largedataset.TypeEmail:
		return "Email"
	case largedataset.TypeRegex:
		return "Regex"
	default:
		return "Unknown"
	}
}
