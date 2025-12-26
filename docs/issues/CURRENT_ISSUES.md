# Current Codebase Issues

## Overview

This document summarizes issues and technical debt found in the current codebase.
Use it as a reference when determining refactoring priorities.

---

## Critical Issues

### 1. Security Vulnerability Risk

**Location**: Dependencies
**Status**: Partially addressed
**Current State**:

- ✅ Go 1.25.5 (updated from 1.20)
- ✅ Major dependencies updated
  - `ethereum/go-ethereum`: v1.16.7 (updated from v1.12.0)
  - `btcsuite/btcd`: v0.25.0 (updated from v0.23.4)
- ⚠️ Regular vulnerability scanning not automated in CI
- ⚠️ CI pipeline not configured (no `.github/workflows/` found)

**Impact**: Potential for known security vulnerabilities
**Priority**: Highest
**Action**:

- [ ] Set up CI pipeline (GitHub Actions)
- [ ] Add `govulncheck` to CI pipeline
- [ ] Set up Dependabot for automated dependency updates
- [ ] Run regular security scans

### 2. Inappropriate `log.Fatal` Usage in Test Utilities

**Location**: `pkg/testutil/repository.go`
**Current State**:

- `pkg/testutil/repository.go` (18 occurrences)
  - Multiple test helper functions using `log.Fatalf`
  - Prevents proper test cleanup on failure

**Problem**: Using `log.Fatal` in test utilities prevents cleanup processing
**Impact**: Resource leaks, potential data inconsistency, test failures don't provide proper cleanup
**Priority**: High
**Action**:

- [ ] Replace `log.Fatal` in `pkg/testutil/repository.go` with error returns
- [ ] Update all test helper functions to return errors instead of calling `log.Fatal`

### 3. `panic` and `log.Fatal` Usage in Main Functions

**Location**: Command tools
**Current State**:

- `cmd/tools/get-eth-key/main.go` (3 `panic` occurrences at lines 35, 45, 50)
- `cmd/tools/eth-key/main.go` (5 `log.Fatal` occurrences at lines 23, 26, 29, 38, 44)

**Note**: `panic` and `log.Fatal` in `main` functions are technically acceptable,
  but graceful shutdown should be implemented for better error handling.

**Problem**: Runtime errors cause program crashes without cleanup
**Impact**: Application stops unexpectedly, no graceful shutdown
**Priority**: Medium
**Action**:

- [ ] Implement graceful shutdown in `main` functions instead of `panic`/`log.Fatal`
- [ ] Add proper error handling and cleanup before exit

---

## High Priority Issues

### 4. Error Handling Inconsistency

**Location**: Codebase-wide
**Current State**:

- ✅ Using `fmt.Errorf` + `%w` (no `pkg/errors` found)
- ⚠️ Error messages lack context information
- ⚠️ Error wrapping is inconsistent
- ⚠️ Custom error types not standardized

**Impact**: Difficult debugging, insufficient error tracing
**Priority**: High
**Action**:

- [ ] Create error handling guidelines
- [ ] Add context information to error messages
- [ ] Define custom error types (`ErrInvalidInput`, `ErrNotFound`, etc.)
- [ ] Standardize error wrapping patterns

### 5. Insufficient Context Management

**Location**: API calls
**Current State**:

- ✅ `internal/infrastructure/api/ethereum/` - Most methods use `context.Context`
- ✅ `internal/infrastructure/api/ripple/` - All gRPC methods use `context.Context`
- ⚠️ `internal/infrastructure/api/bitcoin/` - Methods do NOT use `context.Context`
- ✅ `internal/infrastructure/network/websocket/` - Uses `context.Context`
- ⚠️ Timeout and cancellation not consistently implemented
- ⚠️ Tracing information propagation insufficient

**Impact**:

- Cannot control request timeouts consistently
- Distributed tracing doesn't work properly
- Potential resource leaks

**Priority**: Medium
**Action**:

- [ ] Add `context.Context` parameter to all `internal/infrastructure/api/bitcoin/` methods
- [ ] Implement consistent timeout settings across all API calls
- [ ] Implement cancellation handling
- [ ] Propagate tracing information

### 6. Logging Inconsistency

**Location**: Codebase-wide
**Current State**:

- ✅ Structured logging package exists (`pkg/logger/`)
- ✅ Logger interface defined
- ⚠️ One commented `log.Printf` found in `internal/infrastructure/wallet/key/random_wallet.go:31`
- ⚠️ Log levels may not be unified
- ⚠️ Potential for sensitive information in logs

**Impact**: Difficult log analysis, security risk
**Priority**: High
**Action**:

- [ ] Remove commented `log.Printf` usage in `internal/infrastructure/wallet/key/random_wallet.go:31`
- [ ] Unify log levels
- [ ] Ensure sensitive information is masked
- [ ] Standardize log format

---

## Medium Priority Issues

### 7. Test Infrastructure Improvements

**Location**: Test files
**Current State**:

- ✅ Integration tests use build tags (`//go:build integration`)
- ✅ 35 test files properly tagged with integration build tag
- ⚠️ Test data management could be improved
- ⚠️ Test helpers could be better organized
- ⚠️ `pkg/testutil/repository.go` uses `log.Fatal` which prevents proper test cleanup (see Issue #2)

**Impact**: Long test execution time, difficult test maintenance, improper test failure handling
**Priority**: Medium
**Action**:

- [x] Test separation with build tags (completed)
- [ ] Improve test data management
- [ ] Organize test helpers

### 8. Architecture Migration

**Location**: Application layer
**Current State**:

- ✅ New application layer (`internal/application/usecase/`) exists and is being used
- ✅ DI container (`internal/di/container.go`) uses new use case layer
- ⚠️ Legacy application layer (`internal/wallet/service/`) may still exist in parallel
- ⚠️ Migration strategy unclear

**Impact**: Increased code complexity, unclear which layer to use
**Priority**: Medium
**Action**:

- [ ] Verify if legacy `internal/wallet/service/` is still in use
- [ ] Complete migration from legacy layer to `internal/application/usecase/` if needed
- [ ] Document architecture migration strategy
- [ ] Review and consolidate interfaces

### 9. Type Safety Issues

**Location**: Codebase-wide
**Current State**:

- ✅ Using `sqlc` for type-safe database queries in `internal/infrastructure/repository/`
- ⚠️ Some `interface{}` usage remains
- ⚠️ Type assertions exist in multiple places

**Impact**: Runtime error risk, lack of type safety
**Priority**: Medium
**Action**:

- [ ] Identify remaining `interface{}` usage
- [ ] Change to type-safe implementation where possible
- [ ] Reduce type assertions where feasible

### 10. Commented-Out Code

**Location**: Multiple locations
**Problem**:

- Commented-out code remains
- Debug code remains

**Impact**: Reduced code readability, source of confusion
**Priority**: Medium
**Action**: Remove unnecessary code

**Examples found**:

- `pkg/testutil/xrp.go` (3 commented `log.Fatalf` at lines 61, 68, 71 - entire function is commented out)
- `internal/infrastructure/wallet/key/random_wallet.go:31` (commented `log.Printf`)

---

## Low Priority Issues

### 11. Insufficient Documentation

**Location**: Codebase-wide
**Problem**:

- Package documentation missing
- Function and method comments missing
- Usage examples missing

**Impact**: Difficult code understanding, difficult onboarding
**Priority**: Low
**Action**: Improve documentation

### 12. Naming Convention Inconsistency

**Location**: Codebase-wide
**Problem**:

- Naming conventions not unified
- Abbreviation usage inconsistent

**Impact**: Reduced code readability
**Priority**: Low
**Action**: Unify naming conventions

### 13. Unused Code

**Location**: Codebase-wide
**Problem**:

- Unused imports
- Unused variables
- Unused functions

**Impact**: Code bloat, reduced maintainability
**Priority**: Low
**Action**: Remove unused code (linter can help detect)

---

## Technical Debt

### 14. Dependency Management

**Current State**:

- ✅ Major dependencies updated to recent versions (Go 1.25.5, btcsuite/btcd v0.25.0, ethereum/go-ethereum v1.16.7)
- ⚠️ Regular update process not automated
- ⚠️ Security scanning not in CI (see Issue #1)

**Problem**:

- Security patches may not be applied promptly
- New features may not be available

**Action**:

- [ ] Set up Dependabot (see Issue #1)
- [ ] Add security scanning to CI (see Issue #1)
- [ ] Establish regular update schedule

### 15. Architecture Complexity

**Current State**:

- ✅ Container-based DI pattern exists (`internal/di/container.go`)
- ✅ Domain layer separated (`internal/domain/`)
- ✅ Infrastructure layer separated (`internal/infrastructure/`)
- ✅ New application layer (`internal/application/usecase/`) exists and is being used
- ✅ Container implementation uses `panic` during instance construction (acceptable per project guidelines)
- ⚠️ Legacy application layer may still exist (see Issue #8)

**Problem**:

- Container pattern implementation is complex
- Dependency injection may be overly complex

**Action**:

- [ ] Complete architecture migration (see Issue #8)
- [ ] Document architecture migration strategy
- [ ] Review and simplify architecture where possible

### 16. Insufficient Test Coverage

**Problem**:

- Test coverage unknown
- Critical path tests may be insufficient

**Action**:

- [ ] Measure test coverage
- [ ] Set coverage goals (80%+)
- [ ] Add tests for critical paths

---

## Known Bugs (from README TODO)

### 17. BCH Dependency Issue

**Location**: README.md
**Problem**: `github.com/cpacia/bchutil` is outdated
**Action**: Replace with `github.com/gcash/bchd`

### 18. Private Key Management

**Location**: README.md
**Problem**: Private key handling is inappropriate
**Action**:

- [ ] Improve password input
- [ ] Implement encrypted storage
- [ ] Implement memory protection

### 19. Fee Issue

**Location**: README.md
**Problem**: Overpaying fees on Signet
**Action**: Review fee calculation

---

## Priority Matrix

| Issue | Priority | Impact | Effort | Priority Order |
| ----- | -------- | ------ | ------ | -------------- |
| Security vulnerabilities | Critical | High | Medium | 1 |
| `log.Fatal` usage in testutil | High | High | Low | 2 |
| Error handling | High | Medium | High | 3 |
| Context management | Medium | Medium | Medium | 4 |
| Logging unification | High | Medium | Low | 5 |
| Test infrastructure | Medium | Low | Medium | 6 |
| Architecture migration | Medium | Medium | High | 7 |
| Commented-out code | Medium | Low | Low | 8 |
| Documentation | Low | Low | High | 9 |

---

## Action Plan

### Phase 1 (Immediate)

1. [ ] Set up CI pipeline (GitHub Actions)
2. [ ] Add `govulncheck` to CI pipeline
3. [ ] Set up Dependabot for automated dependency updates
4. [ ] Remove `log.Fatal` from `pkg/testutil/repository.go` (replace with error returns)
5. [ ] Remove commented-out code

### Phase 2 (Within 1-2 weeks)

1. [ ] Implement graceful shutdown in `main` functions
2. [ ] Standardize error handling
3. [ ] Add `context.Context` to all `internal/infrastructure/api/bitcoin/` methods
4. [ ] Implement consistent timeout settings across all API calls
5. [ ] Remove commented `log.Printf` in `internal/infrastructure/wallet/key/random_wallet.go`
6. [ ] Unify logging levels and format

### Phase 3 (Within 1-2 months)

1. [ ] Verify and complete migration from legacy layer to `internal/application/usecase/`
2. [ ] Document architecture migration strategy
3. [ ] Improve test data management
4. [ ] Improve documentation

---

## Notes

- Verify impact scope before fixing each issue
- Run sufficient tests before committing
- Make gradual fixes
- Always conduct code reviews

---

## Update History

- 2024-XX-XX: Initial version created
- 2025-01-XX: Updated to reflect current state (Go 1.25.5, updated dependencies, current architecture)
- 2025-01-XX: Updated panic usage status (acceptable in pkg/di per project guidelines),
  updated context management status (most APIs already use context),
  added architecture migration note
- 2025-01-XX: Reorganized issues, removed completed items, updated file paths to reflect actual structure,
  consolidated duplicate issues, updated priority matrix

## Additional Notes

The codebase is not production-ready for a cryptocurrency wallet without:

- ❌ HD key derivation validation
- ❌ Multisig transaction testing
- ❌ RPC integration testing
- ❌ Comprehensive security audit
- ❌ Private key encryption and secure storage
- ❌ Graceful shutdown implementation
- ❌ Complete context management
- ❌ Error handling standardization
