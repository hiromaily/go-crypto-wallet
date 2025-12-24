# Refactoring Implementation Checklist

## Quick Wins (Immediate Actions)

### 1. Dependency Management

- [x] Go version updated to 1.25.5
- [ ] Run `go mod tidy` regularly
- [ ] Run `go mod verify` for integrity checks
- [ ] Identify outdated dependencies (`go list -m -u all`)
- [ ] Security vulnerability scanning (`govulncheck ./...`)

### 2. Code Quality Immediate Improvements

- [ ] Remove `log.Fatal` usage (except in `main` functions)

  Target files:
  - [ ] `cmd/keygen/main.go` (multiple occurrences)
  - [ ] `cmd/sign/main.go` (multiple occurrences)
  - [ ] `cmd/watch/main.go` (multiple occurrences)
  - [ ] `cmd/tools/eth-key/main.go` (multiple occurrences)
  - [ ] `pkg/testutil/repository.go` (multiple occurrences)
  - [ ] `pkg/repository/watchrepo/*_test.go` (multiple test files)
  - [ ] `pkg/account/*_test.go` (test files)

  Replacement pattern:

  ```go
  // Before
  if err != nil {
      log.Fatal(err)
  }

  // After
  if err != nil {
      return fmt.Errorf("context: %w", err)
  }
  ```

- [ ] Remove unused imports (`goimports`)
- [ ] Remove unused variables (`go vet`)
- [ ] Remove commented-out code

### 3. Linter Configuration

- [x] `.golangci.yml` exists and configured
- [x] Recommended rules enabled (errcheck, staticcheck, gosec, gocritic)
- [ ] Fix remaining linter errors
- [ ] Enable `errcheck` (currently disabled in `.golangci.yml` line 71)

---

## Phase 1: Foundation Setup

### Go Version Update

- [x] `go.mod` `go` directive updated to 1.25.5
- [ ] Verify local Go version matches
- [x] CI environment Go version updated (`.github/workflows/lint-test.yml` uses 1.25)
- [ ] Fix any build errors
- [ ] Run and fix tests

### Major Dependency Updates

- [x] `ethereum/go-ethereum` (currently v1.16.7)
  - [ ] Check for newer versions
  - [ ] Update if security fixes available
  - [ ] Test after update

- [x] `btcsuite/btcd` (currently v0.25.0)
  - [ ] Check for newer versions
  - [ ] Update if security fixes available
  - [ ] Test after update

- [ ] Run tests after each update
- [ ] Verify and fix breaking changes

### Security Scanning

- [x] `govulncheck` available in `go.mod` tools
- [ ] Run vulnerability scan regularly

  ```bash
  govulncheck ./...
  ```

- [ ] Fix identified vulnerabilities
- [ ] Add Dependabot configuration (`.github/dependabot.yml`)
- [ ] Add security scanning to CI workflow

### CI/CD Pipeline

- [x] `.github/workflows/lint-test.yml` exists
  - [x] Test execution
  - [x] Lint checks
  - [ ] Security scanning (add `govulncheck`)
  - [ ] Build verification
- [ ] Create `.github/workflows/release.yml` (if needed)
- [ ] Verify test execution

---

## Phase 2: Code Quality Improvements

### Error Handling Improvements

#### Remove `log.Fatal`

Target files (as identified above):

- [ ] `cmd/keygen/main.go`
- [ ] `cmd/sign/main.go`
- [ ] `cmd/watch/main.go`
- [ ] `cmd/tools/eth-key/main.go`
- [ ] `pkg/testutil/repository.go` (multiple occurrences)
- [ ] Test files in `pkg/repository/watchrepo/`
- [ ] Test files in `pkg/account/`

Replacement pattern:

```go
// Before
if err != nil {
    log.Fatal(err)
}

// After
if err != nil {
    return fmt.Errorf("context: %w", err)
}
```

#### Error Wrapping Standardization

- [x] Using `fmt.Errorf` + `%w` (no `pkg/errors` found)
- [ ] Add context information to error messages
- [ ] Define custom error types
  - [ ] `ErrInvalidInput`
  - [ ] `ErrNotFound`
  - [ ] `ErrUnauthorized`
  - [ ] Domain-specific errors

#### Error Check Coverage

- [ ] Enable `errcheck` in `.golangci.yml` (currently disabled)
- [ ] Fix unchecked errors detected by `errcheck`

### Context Management Introduction

#### Add `context` to API Calls

Target packages:

- [x] `pkg/wallet/api/ethgrp` (already uses `context.Context`)
- [ ] `pkg/wallet/api/btcgrp` (verify all methods use context)
- [ ] `pkg/wallet/api/xrpgrp` (verify all methods use context)

Implementation example:

```go
// Before
func (b *Bitcoiner) GetBalance() (float64, error)

// After
func (b *Bitcoiner) GetBalance(ctx context.Context) (float64, error)
```

#### Timeout Implementation

- [ ] Add timeout settings to all API calls
- [ ] Define default timeout values
- [ ] Make timeout values configurable

#### Graceful Shutdown

- [ ] Implement signal handling in `main` functions
- [ ] Add resource cleanup
- [ ] Wait for in-progress operations to complete

Current status: `main` functions in `cmd/` use `os.Exit()` directly. Need to implement graceful shutdown.

### Logging Standardization

#### Structured Logging Consistency

- [x] Logger package exists (`pkg/logger/`)
- [ ] Unify log levels
- [ ] Unify log fields
- [ ] Mask sensitive information

#### Log Output Improvements

- [ ] Remove or conditionally output debug logs
- [ ] Add stack traces to error logs
- [ ] Add request IDs

---

## Phase 3: Architecture Improvements

### Dependency Injection Improvements

#### Interface Organization

- [ ] Identify unused interfaces
- [ ] Split interfaces appropriately
- [ ] Unify interface naming conventions

#### Mock Generation

- [ ] Introduce `mockgen`

  ```bash
  go install github.com/golang/mock/mockgen@latest
  ```

- [ ] Generate mocks for interfaces
- [ ] Use mocks in tests

### Layer Separation

#### Domain Logic Separation

- [ ] Extract business logic
- [ ] Define domain models
- [ ] Create domain services

#### Infrastructure Layer Separation

- [ ] Abstract database access
- [ ] Abstract external API calls
- [ ] Abstract file I/O

---

## Phase 4: Security Enhancements

### Private Key Management

#### Encrypted Storage

- [ ] Implement private key encryption
- [ ] Manage encryption keys
- [ ] Implement decryption process

#### Memory Protection

- [ ] Use `mlock` (Linux)
- [ ] Zero-clear private keys
- [ ] Prevent memory dumps

#### Password Input

- [ ] Introduce `gopass` or similar

  ```bash
  go get github.com/howeyc/gopass
  ```

- [ ] Improve password input
- [ ] Add password strength checks

### Authentication & Authorization

- [ ] Implement API authentication
- [ ] Implement rate limiting
- [ ] Implement audit logging

---

## Phase 5: Test Improvements

### Test Structure Organization

#### Build Tag Usage

- [x] `//go:build integration` tags added (35+ test files)
- [x] Integration tests separated
- [ ] Organize unit tests

#### Test Helpers

- [ ] Organize test helper functions
- [ ] Manage test data
- [ ] Organize mocks

### Test Coverage

- [ ] Generate coverage reports

  ```bash
  go test -coverprofile=coverage.out ./...
  go tool cover -html=coverage.out
  ```

- [ ] Set coverage goals
- [ ] Integrate coverage into CI

---

## Phase 6: Monitoring and Observability

### Migration to OpenTelemetry

- [ ] Introduce OpenTelemetry SDK
- [ ] Migrate from Jaeger (if used)
- [ ] Implement tracing
- [ ] Implement metrics

### Log Management

- [ ] Unify structured logging
- [ ] Configure log aggregation
- [ ] Set up alerts

---

## Phase 7: Performance Optimization

### Database

- [ ] Profile queries
- [ ] Review indexes
- [ ] Optimize connection pools

### Memory Management

- [ ] Memory profiling

  ```bash
  go test -memprofile=mem.prof ./...
  go tool pprof mem.prof
  ```

- [ ] Detect memory leaks
- [ ] Optimize memory usage

---

## Phase 8: Documentation

### Code Documentation

- [ ] Add package comments
- [ ] Add comments to exported functions and methods
- [ ] Add usage examples for complex logic

### Architecture Documentation

- [ ] Update architecture diagrams
- [ ] Record design decisions (ADR)
- [ ] Maintain API specifications

---

## Verification Methods

### Verification Items for Each Phase Completion

- [ ] All tests pass
- [ ] Zero linter errors
- [ ] Zero Critical/High security vulnerabilities
- [ ] Documentation updated
- [ ] Code review completed

### Performance Verification

- [ ] Run benchmark tests

  ```bash
  go test -bench=. -benchmem ./...
  ```

- [ ] Compare performance
- [ ] Compare memory usage

---

## Notes

### Managing Breaking Changes

- Verify impact scope of each change
- Gradual release
- Prepare rollback plans

### Importance of Testing

- Run tests before refactoring
- Run tests after refactoring
- Conduct regression tests

### Communication

- Record changes
- Share with team members
- Conduct code reviews

---

## Current Status Summary

### Completed ✅

- Go version: 1.25.5
- Linter configuration: `.golangci.yml` configured
- CI/CD: Basic workflow exists (`.github/workflows/lint-test.yml`)
- Integration tests: Build tags implemented
- Error wrapping: Using `fmt.Errorf` + `%w`
- Context usage: Partially implemented (ethgrp)
- Dependencies: Relatively modern versions

### In Progress / Needs Work 🔄

- `log.Fatal` removal: Many occurrences remain
- Context management: Not all APIs use context
- Graceful shutdown: Not implemented
- Security scanning: Not in CI
- Error types: Need standardization
- Test coverage: Needs measurement and improvement

### Priority Actions 🎯

1. Remove `log.Fatal` from non-main packages
2. Add `context.Context` to all API methods
3. Implement graceful shutdown
4. Add security scanning to CI
5. Enable `errcheck` in linter config
6. Standardize error types
