# Filter Examples

The `parse` command supports powerful filtering capabilities to select specific rows from CSV files.

## Basic Usage

```bash
# Single filter on numeric column
go run main.go parse --file data.csv --has-header --filter "Price > 100"

# Single filter on string column
go run main.go parse --file data.csv --has-header --filter "Symbol = 'AAPL'"

# Multiple filters (AND logic)
go run main.go parse --file data.csv --has-header \
  --filter "Price > 100" \
  --filter "Exchange = 'NYSE'"
```

## Supported Operators

### Comparison Operators (Numeric and String)

- `=` - Equal to
- `!=` - Not equal to
- `>` - Greater than (numeric only)
- `>=` - Greater than or equal (numeric only)
- `<` - Less than (numeric only)
- `<=` - Less than or equal (numeric only)

### String Operators

- `contains` - Substring match
- `startswith` - Prefix match
- `endswith` - Suffix match
- `~=` - Regular expression match

## Filter Syntax

Filters use the format: `<column> <operator> <value>`

**Column specification:**

- By name: `Price > 100`
- By index: `8 > 100` (0-based index)

**Value specification:**

- Quoted: `Symbol = 'AAPL'`
- Unquoted: `Symbol = AAPL`
- Numeric: `Price > 100.50`

**Important:** Use spaces around operators: `Price > 100` (not `Price>100`)

## Examples

### Numeric Filtering

```bash
# Trades with price greater than 200
go run main.go parse --file data.csv --has-header --filter "Price > 200"

# Trades with price between 100 and 500
go run main.go parse --file data.csv --has-header \
  --filter "Price >= 100" \
  --filter "Price <= 500"

# Trades with negative change percent
go run main.go parse --file data.csv --has-header --filter "ChangePercent < 0"
```

### String Filtering

```bash
# Specific exchange
go run main.go parse --file data.csv --has-header --filter "Exchange = 'NYSE'"

# Multiple exchanges (using startswith for prefix)
go run main.go parse --file data.csv --has-header --filter "Exchange startswith N"

# Symbols containing 'AA'
go run main.go parse --file data.csv --has-header --filter "Symbol contains AA"

# Buy orders only
go run main.go parse --file data.csv --has-header --filter "TradeType = 'Buy'"
```

### Regular Expression Filtering

```bash
# Symbols starting with A, M, or G
go run main.go parse --file data.csv --has-header --filter "Symbol ~= ^[AMG]"

# Timestamps in morning hours (09:30 to 09:59)
go run main.go parse --file data.csv --has-header --filter "Timestamp ~= 09:3[0-9]"
```

### Combining Filters with Group-by

```bash
# NYSE trades over $300, grouped by symbol
go run main.go parse --file data.csv --has-header \
  --filter "Exchange = 'NYSE'" \
  --filter "Price >= 300" \
  --group-by 2

# Buy orders in Technology sector, grouped by exchange
go run main.go parse --file data.csv --has-header \
  --filter "TradeType = 'Buy'" \
  --filter "Sector = 'Technology'" \
  --group-by 3
```

### Combining with Validation

```bash
# Validate schema AND filter
go run main.go parse --file data.csv --has-header \
  --validate \
  --filter "Price > 100" \
  --filter "Exchange = 'NYSE'"
```

## Filter Logic

- Multiple filters use **AND** logic (all must match)
- To achieve OR logic, run separate commands or use regex
- Filters are applied AFTER schema validation
- Invalid rows (parsing errors) are excluded before filtering

## Performance Tips

1. Use numeric comparisons when possible (faster than string matching)
2. Filter early in the pipeline (before group-by)
3. Use column indices instead of names for slightly better performance
4. Avoid complex regex patterns on large datasets
