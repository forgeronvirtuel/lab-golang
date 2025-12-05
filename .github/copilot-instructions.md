# Lab Golang - AI Coding Assistant Guide

## Project Overview

This is a CLI tool for working with large CSV datasets, specifically focused on stock market data generation and analysis. Built with Go 1.25.3 and Cobra for command handling.

## Architecture

### Command Structure (Cobra-based)

- **Entry point**: `main.go` → `cmd.Execute()`
- **Root command**: `cmd/root.go` - defines the base CLI structure
- **Subcommands**: Each in separate files under `cmd/`:
  - `generate` - creates synthetic stock market CSV files
  - `read` - simple CSV reading with streaming
  - `parse` - CSV parsing with validation, stats, and optional group-by functionality

### Flag Management Pattern

All CLI flags are **shared global variables** declared in `cmd/flagsvar.go`:

```go
var (
    outputFile string
    numRows    int
    filePath   string
    separator  string
    showFirst  int
    hasHeader  bool
    groupByCol int
)
```

Individual commands bind these variables using Cobra's `Flags().StringVarP()` pattern in their `init()` functions. When adding new flags, declare them in `flagsvar.go` first.

### Internal Packages

- `internal/largedataset/` - CSV parsing and streaming statistics
  - `csv.go`: `ParseLogicalRow()` and `ParseLogicalRowWithGroupBy()` - extract amount from column index 8 and optional group key
  - `stats.go`: `AmountStats` struct for streaming aggregations (count, sum, min, max, avg)
- `internal/fgip/` - IP address anonymization utilities (minimal usage)

## Key Patterns

### CSV Processing Philosophy

All commands use **streaming processing** with `bufio.NewReader()` wrapped around file handles - never load entire files into memory. Example pattern:

```go
f, _ := os.Open(path)
defer f.Close()
reader := bufio.NewReader(f)
csvReader := csv.NewReader(reader)
for {
    record, err := csvReader.Read()
    if err == io.EOF { break }
    // process record
}
```

### Group-by Statistics

The `parse` command supports optional group-by statistics via the `--group-by` flag:
- Statistics are calculated both globally and per group
- Groups are displayed sorted by total sum (descending)
- Group keys are extracted from the specified column index in `ParseLogicalRowWithGroupBy()`

### Generated CSV Format

The `generate` command creates stock data with 22 columns (see header in `generate.go`):

- Column 8 is `Price` - this is the hardcoded amount column used by parsing
- Column 2 is `Symbol` - commonly used for group-by operations
- Column 3 is `Exchange` - alternative group-by field

### Error Handling Convention

- Use `log.Fatal()` for user-facing errors (missing required flags, file not found)
- Return errors from internal functions for programmatic handling
- Track invalid rows separately but continue processing (don't fail entire job)

## Common Development Tasks

### Adding a New Command

1. Create `cmd/newcommand.go`
2. Define `newCmd = &cobra.Command{...}` with Use, Short, Long, Example, Run
3. Add flags in `init()` using variables from `cmd/flagsvar.go` (or add new ones there)
4. Register with `rootCmd.AddCommand(newCmd)` in the `init()` function
5. Test with: `go run main.go newcommand --help`

### Working with Test Data

Pre-generated CSV files exist with row counts indicated by filename:

- `100_bourse.csv`, `1000_bourse.csv`, ..., `10000000_bourse.csv`
- Use smaller files for development/testing
- Generate new data: `go run main.go generate -o test.csv -r 1000`

### Running the CLI

Development usage (see `memo.md` for examples):

```bash
go run main.go <command> --help
go run main.go generate -o data.csv -r 100000
go run main.go parse --file data.csv --has-header --group-by 2
```

### Column Indices for Group-by

When using the `parse` command with the generated stock data:

- `--group-by 2` for Symbol
- `--group-by 3` for Exchange
- `--group-by 4` for Sector

Note: The amount column (Price) is hardcoded to index 8 in `ParseLogicalRow()`

## Project Quirks

- No tests exist in the codebase currently
- `internal/fgip/anonimize.go` exists but isn't integrated into any commands
- The `read` and `parse` commands overlap in functionality - `parse` is more feature-complete with stats and group-by
- Progress reporting in `generate` happens every 100k rows
- Timestamps in generated data increment by 1 second per row starting from 2020-01-01 09:30:00 UTC
- Amount column index is hardcoded to 8 in `internal/largedataset/csv.go` (not configurable)
