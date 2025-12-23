# Migrate from ericlagergren/decimal to quagmt/udecimal

## Background

Currently, the project uses `github.com/ericlagergren/decimal` for decimal number handling in financial calculations. We want to migrate to `github.com/quagmt/udecimal` for the following benefits:

- **High Performance**: 5x~20x faster than shopspring/decimal and ericlagergren/decimal
- **Zero Memory Allocation**: Designed for almost 99% zero memory allocation
- **High Precision**: Supports up to 19 decimal places with no precision loss
- **Panic-Free**: All errors are returned as values, ensuring no unexpected panics
- **Concurrent-Safe**: All arithmetic operations return a new `Decimal` value
- **Financial Application Focus**: Specifically designed for financial applications

## API Differences

### Type Changes
- **Old**: `*decimal.Big` (pointer type)
- **New**: `udecimal.Decimal` (value type)

### Creation Methods
- **Old**: 
  ```go
  dAmt := new(decimal.Big)
  dAmt, _ = dAmt.SetString("123.456")
  ```
- **New**:
  ```go
  dAmt, err := udecimal.Parse("123.456")
  // or
  dAmt, err := udecimal.NewFromFloat64(123.456)
  // or
  dAmt, err := udecimal.NewFromInt64(123456, 3) // 123.456
  ```

### String Conversion
- **Old**: `dAmt.String()`
- **New**: `dAmt.String()` (same, but also supports `StringFixed(precision)`)

### Error Handling
- **Old**: `SetString()` returns `(*Big, bool)` - errors are ignored
- **New**: `Parse()`, `NewFromFloat64()`, etc. return `(Decimal, error)` - must handle errors

## Files to Modify

### 1. Dependency Management
- [ ] Update `go.mod` to replace `github.com/ericlagergren/decimal` with `github.com/quagmt/udecimal`
- [ ] Remove the `replace` directive for `ericlagergren/decimal` if present
- [ ] Run `go mod tidy` to clean up dependencies

### 2. Model Definitions
- [ ] `pkg/models/rdb/models.go`
  - Change `*decimal.Big` to `udecimal.Decimal` in:
    - `BTCTX.TotalInputAmount`
    - `BTCTX.TotalOutputAmount`
    - `BTCTX.Fee`
    - `BTCTXInput.InputAmount`
    - `BTCTXOutput.OutputAmount`
    - `PaymentRequest.Amount`

### 3. API Interface
- [ ] `pkg/wallet/api/btcgrp/api-interface.go`
  - Change return types from `*decimal.Big` to `udecimal.Decimal`:
    - `AmountToDecimal(amt btcutil.Amount) udecimal.Decimal`
    - `FloatToDecimal(f float64) udecimal.Decimal`

### 4. Implementation
- [ ] `pkg/wallet/api/btcgrp/btc/amount.go`
  - Update `AmountToDecimal()` to use `udecimal.Parse()`
  - Update `FloatToDecimal()` to use `udecimal.NewFromFloat64()`
  - Handle errors properly (return error or use zero value)

### 5. Converter Package
- [ ] `pkg/converter/converter.go`
  - Update interface: `FloatToDecimal(f float64) (udecimal.Decimal, error)`
  - Update implementation to use `udecimal.NewFromFloat64()`
  - Return error instead of ignoring it

### 6. Repository Layer
- [ ] `pkg/repository/watchrepo/btc_tx_sqlc.go`
  - Update `convertSqlcBtcTxToModel()` to use `udecimal.Parse()`
  - Handle errors appropriately
- [ ] `pkg/repository/watchrepo/btc_tx_input_sqlc.go`
  - Update `convertSqlcBtcTxInputToModel()` to use `udecimal.Parse()`
- [ ] `pkg/repository/watchrepo/btc_tx_output_sqlc.go`
  - Update `convertSqlcBtcTxOutputToModel()` to use `udecimal.Parse()`
- [ ] `pkg/repository/watchrepo/payment_request_sqlc.go`
  - Update `convertSqlcPaymentRequestToModel()` to use `udecimal.Parse()`

### 7. Test Files
- [ ] `pkg/wallet/api/btcgrp/btc/amount_test.go`
  - Update tests to use `udecimal.Decimal` and handle errors
- [ ] `pkg/repository/watchrepo/btc_tx_sqlc_test.go`
  - Update test data creation to use `udecimal.Parse()` or `udecimal.NewFromFloat64()`
- [ ] `pkg/repository/watchrepo/btc_tx_input_sqlc_test.go`
  - Update test data creation
- [ ] `pkg/repository/watchrepo/btc_tx_output_sqlc_test.go`
  - Update test data creation
- [ ] `pkg/repository/watchrepo/payment_request_sqlc_test.go`
  - Update test data creation

### 8. Any Other Usage
- [ ] Search for any other files using `decimal.Big` or `ericlagergren/decimal`
- [ ] Update all occurrences

## Migration Steps

1. **Add new dependency**:
   ```bash
   go get github.com/quagmt/udecimal
   ```

2. **Update imports**: Replace all imports from:
   ```go
   "github.com/ericlagergren/decimal"
   ```
   to:
   ```go
   "github.com/quagmt/udecimal"
   ```

3. **Update type definitions**: Change `*decimal.Big` to `udecimal.Decimal` in model structs

4. **Update creation code**: Replace:
   ```go
   dAmt := new(decimal.Big)
   dAmt, _ = dAmt.SetString(str)
   ```
   with:
   ```go
   dAmt, err := udecimal.Parse(str)
   if err != nil {
       // handle error
   }
   ```

5. **Update float conversion**: Replace:
   ```go
   dAmt := new(decimal.Big)
   strAmt := fmt.Sprintf("%f", f)
   dAmt, _ = dAmt.SetString(strAmt)
   ```
   with:
   ```go
   dAmt, err := udecimal.NewFromFloat64(f)
   if err != nil {
       // handle error
   }
   ```

6. **Update interface methods**: Change return types and add error returns where appropriate

7. **Update all callers**: Ensure all code calling these methods handles the new types and errors

8. **Remove old dependency**:
   ```bash
   go mod edit -droprequire github.com/ericlagergren/decimal
   go mod tidy
   ```

## Testing Requirements

- [ ] Run all existing tests to ensure they pass with the new library
- [ ] Verify decimal precision is maintained correctly
- [ ] Test edge cases (very large numbers, very small numbers, zero, negative numbers)
- [ ] Run integration tests to ensure database operations work correctly
- [ ] Verify that string representations match expected formats
- [ ] Run `make lint-fix` to ensure code quality
- [ ] Run `make check-build` to verify the code builds successfully
- [ ] Run `make gotest` to run all tests

## Important Notes

1. **Error Handling**: Unlike `ericlagergren/decimal`, `udecimal` returns errors. All error returns must be properly handled. Do not ignore errors.

2. **Value vs Pointer**: `udecimal.Decimal` is a value type, not a pointer. This means:
   - No need for `new()` or pointer allocation
   - Methods return new values (immutable operations)
   - Can be safely used in concurrent contexts

3. **Precision**: `udecimal` supports up to 19 decimal places. Ensure this meets the project's requirements.

4. **Database Compatibility**: Verify that the database schema and SQL queries are compatible with the new decimal type representation.

5. **Backward Compatibility**: If there are any external APIs or serialization formats that depend on the decimal representation, ensure they remain compatible.

## References

- [udecimal GitHub Repository](https://github.com/quagmt/udecimal)
- [udecimal Documentation](https://quagmt.github.io/udecimal/)

## Acceptance Criteria

- [ ] All code has been migrated from `ericlagergren/decimal` to `udecimal`
- [ ] All tests pass
- [ ] No linter errors
- [ ] Code builds successfully
- [ ] All error handling is properly implemented
- [ ] Documentation is updated if necessary
- [ ] Performance improvements are verified (optional but recommended)

