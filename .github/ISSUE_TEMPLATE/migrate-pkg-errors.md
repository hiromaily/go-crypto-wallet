# Migrate from `github.com/pkg/errors` to Standard Library

## Background

The project currently uses `github.com/pkg/errors` v0.9.1, which is an outdated package. Since Go 1.13, error wrapping and unwrapping have been standardized in the standard library through the `errors` package and the `%w` verb in `fmt.Errorf`.

As of 2025, the best practice for Go error handling is to use only the standard library. Migrating from `github.com/pkg/errors` to the standard library provides the following benefits:

- Reduced external dependencies
- Adoption of standard error handling patterns
- Improved code readability and maintainability
- Alignment with guidelines in `AGENTS.md` and `REFACTORING_CHECKLIST.md`

## Current State Analysis

### Usage Statistics

- **Scope of Impact**: `github.com/pkg/errors` is used in approximately 108 files
- **Main Usage Patterns**:
  - `errors.Wrap(err, "message")` - Most common (100+ occurrences)
  - `errors.Wrapf(err, "format %s", arg)` - Formatted wrapping (~10 occurrences)
  - `errors.Errorf("format %s", arg)` - Error creation (several occurrences)
  - `errors.New("message")` - Error creation (several occurrences)

### Unused Features

The following `github.com/pkg/errors` features are not used:

- `errors.Cause()` - Root cause retrieval
- `errors.WithMessage()` - Message addition
- `errors.WithStack()` - Stack trace addition

### Main Affected Packages

- `pkg/wallet/api/` - API client layer (btcgrp, ethgrp, xrpgrp)
- `pkg/wallet/service/` - Business logic layer
- `pkg/repository/` - Repository layer (watchrepo, coldrepo)
- `pkg/db/` - Database layer
- `pkg/wallet/key/` - Key management layer
- Many other packages

## Migration Plan

### Phase 1: Preparation and Analysis (1-2 hours)

1. **Complete Usage Identification**

   ```bash
   # Search for usage locations
   grep -r "github.com/pkg/errors" --include="*.go" | wc -l
   grep -r "errors\.Wrap" --include="*.go" | wc -l
   grep -r "errors\.Wrapf" --include="*.go" | wc -l
   grep -r "errors\.Errorf" --include="*.go" | wc -l
   ```

2. **Test Execution and Baseline Establishment**

   ```bash
   make gotest
   make check-build
   ```

### Phase 2: Incremental Migration (4-6 hours)

#### Step 1: Bottom-Up Approach

Migration will be performed in the following order:

1. **Foundation Layer** (`pkg/db/`, `pkg/repository/`)
2. **API Layer** (`pkg/wallet/api/`)
3. **Service Layer** (`pkg/wallet/service/`)
4. **Other Packages**

#### Step 2: Package-by-Package Migration

For each package, follow these steps:

1. Update imports
2. Replace error wrapping
3. Run tests
4. Lint check

### Phase 3: Verification and Cleanup (1-2 hours)

1. Run all tests
2. Lint check
3. Remove dependencies
4. Update documentation

## Implementation Steps

### 1. Update Imports

**Before:**

```go
import (
    "github.com/pkg/errors"
)
```

**After:**

```go
import (
    "errors"
    "fmt"
)
```

### 2. Replace Error Wrapping

#### Replacing `errors.Wrap(err, "message")`

**Before:**

```go
if err != nil {
    return nil, errors.Wrap(err, "failed to call GetAllPaymentRequests()")
}
```

**After:**

```go
if err != nil {
    return nil, fmt.Errorf("failed to call GetAllPaymentRequests(): %w", err)
}
```

#### Replacing `errors.Wrapf(err, "format %s", arg)`

**Before:**

```go
if err != nil {
    return 0, errors.Wrapf(err, "fail to call btcutil.NewAmount(%f)", f)
}
```

**After:**

```go
if err != nil {
    return 0, fmt.Errorf("fail to call btcutil.NewAmount(%f): %w", f, err)
}
```

**Note**: With `fmt.Errorf`, append `: %w` to the format string and pass the error as the last argument.

#### Replacing `errors.Errorf("format %s", arg)`

**Before:**

```go
return nil, errors.Errorf("Connection(): error: %v", err)
```

**After:**

```go
return nil, fmt.Errorf("Connection(): error: %v", err)
```

#### Replacing `errors.New("message")`

**Before:**

```go
return 0, errors.New("server is closed")
```

**After:**

```go
import "errors"  // Standard library errors package

return 0, errors.New("server is closed")
```

**Note**: `errors.New()` uses the standard library `errors` package. Be careful not to confuse it with `fmt.Errorf()`.

### 3. Improve Error Checking (Optional)

It is recommended to add error checking using `errors.Is()` or `errors.As()` as part of the migration.

**Example:**

```go
// Check for io.EOF
if err != nil {
    if errors.Is(err, io.EOF) {
        // Handle EOF
    }
    return fmt.Errorf("operation failed: %w", err)
}
```

### 4. Bulk Replacement Considerations

**Not Recommended:**

- Bulk replacement using regex (may lose context)
- Manual bulk replacement (risk of missing cases)

**Recommended Approach:**

- Migrate incrementally, file-by-file or package-by-package
- Run tests after each migration to verify
- Conduct code reviews

## Verification Methods

### 1. Build Check

```bash
make check-build
```

### 2. Lint Check

```bash
make lint-fix
```

### 3. Test Execution

```bash
make gotest
```

### 4. Dependency Verification

```bash
go mod tidy
go mod verify
```

### 5. Usage Verification

After migration is complete, verify with:

```bash
# Verify that github.com/pkg/errors usage is zero
grep -r "github.com/pkg/errors" --include="*.go" | wc -l
# Expected: 0

# Verify that standard library errors and fmt are used appropriately
grep -r "fmt.Errorf.*%w" --include="*.go" | wc -l
```

## Important Considerations

### Security

- Ensure error messages do not contain sensitive information (private keys, passwords, etc.)
- Follow security guidelines in `AGENTS.md`

### Compatibility

- Be aware that error message formats may change, affecting code that depends on external APIs or logs
- If using error type checking (`errors.Is()`, `errors.As()`), verify behavior

### Performance

- Performance impact from migrating to standard library is minimal
- Stack trace information is lost, but errors can be tracked using standard Go 1.13+ methods

### Testing

- Run tests for each package after migration
- Run integration tests to verify error propagation works correctly

## Completion Criteria

- [ ] All `github.com/pkg/errors` usage has been replaced with standard library
- [ ] `github.com/pkg/errors` dependency has been removed from `go.mod`
- [ ] All tests are passing
- [ ] Zero lint errors
- [ ] `make check-build` succeeds
- [ ] `go mod tidy` completes successfully
- [ ] Code review is complete

## References

- [Go 1.13 Error Handling](https://go.dev/blog/go1.13-errors)
- [Working with Errors in Go 1.13](https://go.dev/doc/go1.13#error_wrapping)
- [AGENTS.md](../AGENTS.md) - Error handling guidelines
- [REFACTORING_CHECKLIST.md](../REFACTORING_CHECKLIST.md) - Refactoring checklist

## Example Commands for Implementation

```bash
# 1. Check current usage
grep -r "errors\.Wrap" --include="*.go" | head -20

# 2. Migrate specific file
# After editing the file...

# 3. Build check
make check-build

# 4. Lint check
make lint-fix

# 5. Run tests
make gotest

# 6. Organize dependencies
go mod tidy

# 7. Verify after migration
grep -r "github.com/pkg/errors" --include="*.go"
```
