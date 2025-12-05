package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "lab-golang",
	Short: "A CLI tool for working with CSV files and data generation",
	Long: `Lab Golang is a CLI application that provides utilities for:
  - Generating large CSV files with stock market data
  - Reading and analyzing CSV files
  - Anonymizing IP addresses`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Flags globaux peuvent être ajoutés ici si nécessaire
}
