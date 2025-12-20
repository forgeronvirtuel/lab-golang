# Logical Operators in Filters - Documentation

## Overview

The filter system now supports complex logical expressions with `AND`, `OR` operators and parentheses for grouping.

## Basic Syntax

### AND Operator

All conditions must be true:

```bash
--filter "Price > 100 AND Volume > 1000"
```

### OR Operator

At least one condition must be true:

```bash
--filter "Price > 400 OR Volume > 9000"
```

### Parentheses

Control evaluation order:

```bash
--filter "(Price > 100 OR Volume > 5000) AND Exchange = 'NASDAQ'"
```

## Operator Precedence

1. **Parentheses** `()` - highest precedence
2. **AND** - higher precedence than OR
3. **OR** - lowest precedence

### Examples of Precedence

```bash
# Without parentheses: (A AND B) OR C
--filter "Price > 100 AND Volume > 1000 OR Symbol = 'AAPL'"

# With parentheses: A AND (B OR C)
--filter "Price > 100 AND (Volume > 1000 OR Symbol = 'AAPL')"
```

## Real-World Examples

### Example 1: High Value or Large Volume Trades

Find trades with either high price OR high volume:

```bash
go run main.go parse --file data.csv --has-header \
  --filter "Price > 400 OR Volume > 9000000"
```

### Example 2: Specific Exchange with Conditions

Find trades on TSE that meet price or volume criteria:

```bash
go run main.go parse --file data.csv --has-header \
  --filter "(Price > 400 OR Volume > 9000000) AND Exchange = 'TSE'"
```

### Example 3: Multiple Symbols on Specific Exchange

Find Apple or Microsoft stocks on LSE:

```bash
go run main.go parse --file data.csv --has-header \
  --filter "(Symbol = 'AAPL' OR Symbol = 'MSFT') AND Exchange = 'LSE'"
```

### Example 4: Complex Multi-Condition Filter

Find high-value trades OR specific symbol on specific exchange:

```bash
go run main.go parse --file data.csv --has-header \
  --filter "Price > 300 AND Volume > 5000000 OR Symbol = 'AAPL' AND Exchange = 'LSE'"
```

This is interpreted as: `(Price > 300 AND Volume > 5000000) OR (Symbol = 'AAPL' AND Exchange = 'LSE')`

### Example 5: Price Range with String Operators

Find trades in price range with symbol pattern:

```bash
go run main.go parse --file data.csv --has-header \
  --filter "Price >= 100 AND Price <= 500 AND Symbol startswith 'A'"
```

### Example 6: Multiple Exchanges

Find trades on NYSE or NASDAQ:

```bash
go run main.go parse --file data.csv --has-header \
  --filter "Exchange = 'NYSE' OR Exchange = 'NASDAQ'"
```

### Example 7: Nested Conditions

Complex nested logic:

```bash
go run main.go parse --file data.csv --has-header \
  --filter "((Price > 100 OR Volume > 5000) AND Exchange = 'NASDAQ') OR (Symbol = 'MSFT' AND TradeType = 'Buy')"
```

## Available Comparison Operators

All standard comparison operators work with logical operators:

- `=` - Equal
- `!=` - Not equal
- `>` - Greater than
- `>=` - Greater than or equal
- `<` - Less than
- `<=` - Less than or equal
- `contains` - String contains
- `startswith` - String starts with
- `endswith` - String ends with
- `~=` - Regular expression match

## Combining with Other Flags

Logical operators work seamlessly with other parse command features:

### With Group-By Statistics

```bash
go run main.go parse --file data.csv --has-header \
  --filter "Price > 200 OR Volume > 8000000" \
  --group-by 3 \
  --show-first 10
```

### With Validation

```bash
go run main.go parse --file data.csv --has-header \
  --filter "(Price > 100 AND Volume > 1000) OR Exchange = 'TSE'" \
  --validate
```

## Backward Compatibility

The system maintains full backward compatibility:

### Old Style (Multiple --filter flags)

These are automatically combined with AND:

```bash
go run main.go parse --file data.csv --has-header \
  --filter "Price > 100" \
  --filter "Volume > 1000" \
  --filter "Symbol = 'AAPL'"
```

Equivalent to: `Price > 100 AND Volume > 1000 AND Symbol = 'AAPL'`

### New Style (Single filter with logical operators)

```bash
go run main.go parse --file data.csv --has-header \
  --filter "Price > 100 AND Volume > 1000 AND Symbol = 'AAPL'"
```

## Tips and Best Practices

1. **Use Parentheses for Clarity**
   Even when not strictly needed, parentheses improve readability:

   ```bash
   --filter "(Price > 100) AND (Volume > 1000)"
   ```

2. **Understand Precedence**
   Remember that AND binds tighter than OR:

   - `A OR B AND C` means `A OR (B AND C)`
   - Use parentheses to get `(A OR B) AND C`

3. **Space Around Operators**
   Always use spaces around logical operators:

   - ✅ `Price > 100 AND Volume > 1000`
   - ❌ `Price > 100AND Volume > 1000`

4. **Quote the Entire Filter**
   Use single or double quotes around the filter expression:

   ```bash
   --filter "Price > 100 AND Volume > 1000"
   --filter '(Price > 100 OR Volume > 5000) AND Exchange = "NASDAQ"'
   ```

5. **Combining Many Conditions**
   Break complex filters into logical groups:
   ```bash
   --filter "(high_price_condition OR high_volume_condition) AND exchange_condition"
   ```

## Performance Notes

- **Short-Circuit Evaluation**:
  - AND stops at the first `false` condition
  - OR stops at the first `true` condition
- **Order Matters**: Put most selective conditions first in AND expressions
- **Complexity**: Deeply nested expressions may impact performance on very large datasets

## Implementation Details

The filter system uses an expression tree:

- `FilterExpression` - interface for any expression
- `FilterCondition` - leaf node (single condition)
- `FilterGroup` - composite node (AND/OR group)

This allows for arbitrary nesting and complex logical expressions while maintaining clean evaluation semantics.
