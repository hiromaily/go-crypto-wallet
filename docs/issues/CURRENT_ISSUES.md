# Current Codebase Issues

## Overview

This document summarizes issues and technical debt found in the current codebase. Use it as a reference when determining refactoring priorities.

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

**Impact**: Potential for known security vulnerabilities
**Priority**: Highest
**Action**:

- [ ] Add `govulncheck` to CI pipeline
- [ ] Set up Dependabot for automated dependency updates
- [ ] Run regular security scans

### 2. `panic` Usage

**Location**: Multiple files
**Current State**:

- `pkg/di/container.go` (17 occurrences)
  - Lines: 108, 167, 175, 187, 189, 379, 433, 444, 455, 466, 485, 501, 517, 545, 577, 715, 797
- `cmd/tools/get-eth-key/main.go` (3 occurrences)
  - Lines: 35, 45, 50

**Problem**: Runtime errors cause program crashes
**Impact**: Application stops unexpectedly on invalid coin types or configuration errors
**Priority**: High
**Action**: Replace `panic` with error returns

**Example**:

```go
// Before
panic(fmt.Sprintf("coinType[%s] is not implemented yet.", c.conf.CoinTypeCode))

// After
return nil, fmt.Errorf("coinType[%s] is not implemented yet", c.conf.CoinTypeCode)
```

### 3. Inappropriate `log.Fatal` Usage

**Location**: Multiple files
**Current State**:

- `cmd/keygen/main.go` (multiple occurrences)
- `cmd/sign/main.go` (multiple occurrences)
- `cmd/watch/main.go` (multiple occurrences)
- `cmd/tools/eth-key/main.go` (multiple occurrences)
- `pkg/testutil/repository.go` (multiple occurrences)
- Test files in `pkg/repository/watchrepo/` (multiple files)
- Test files in `pkg/account/` (multiple files)

**Note**: `log.Fatal` in `main` functions is acceptable, but should be replaced with proper error handling and graceful shutdown.

**Problem**: Using `log.Fatal` in non-main packages prevents cleanup processing
**Impact**: Resource leaks, potential data inconsistency
**Priority**: High
**Action**: Replace with error returns (except in `main` functions where graceful shutdown should be implemented)

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

- ✅ `pkg/wallet/api/ethgrp` uses `context.Context`
- ⚠️ `pkg/wallet/api/btcgrp` - needs verification
- ⚠️ `pkg/wallet/api/xrpgrp` - needs verification
- ⚠️ Timeout and cancellation not implemented
- ⚠️ Tracing information propagation insufficient

**Impact**:

- Cannot control request timeouts
- Distributed tracing doesn't work properly
- Potential resource leaks

**Priority**: High
**Action**:

- [ ] Add `context.Context` to all API calls
- [ ] Implement timeout settings
- [ ] Implement cancellation
- [ ] Propagate tracing information

### 6. Logging Inconsistency

**Location**: Codebase-wide
**Current State**:

- ✅ Structured logging package exists (`pkg/logger/`)
- ✅ Logger interface defined
- ⚠️ One commented `log.Printf` found in `pkg/wallet/key/random_wallet.go:31`
- ⚠️ Log levels may not be unified
- ⚠️ Potential for sensitive information in logs

**Impact**: Difficult log analysis, security risk
**Priority**: High
**Action**:

- [ ] Remove commented `log.Printf` usage
- [ ] Unify log levels
- [ ] Ensure sensitive information is masked
- [ ] Standardize log format

---

## Medium Priority Issues

### 7. Test Separation

**Location**: Test files
**Current State**:

- ✅ Integration tests use build tags (`//go:build integration`)
- ✅ 35+ test files properly tagged
- ⚠️ Test data management could be improved
- ⚠️ Test helpers could be better organized

**Impact**: Long test execution time, difficult test maintenance
**Priority**: Medium
**Action**:

- [x] Test separation with build tags (completed)
- [ ] Improve test data management
- [ ] Organize test helpers

### 8. Interface Overuse

**Location**: `pkg/wallet/service/`
**Problem**:

- Many small interfaces defined
- Interface usage is inconsistent

**Impact**: Increased code complexity
**Priority**: Medium
**Action**: Review and consolidate interfaces

### 9. Type Safety Issues

**Location**:

- `pkg/repository/watchrepo/` (sqlc is used, not sqlboiler)
- Various places using `map[string]interface{}`

**Current State**:

- ✅ Using `sqlc` for type-safe database queries (not sqlboiler)
- ⚠️ Some `interface{}` usage remains
- ⚠️ Type assertions exist in multiple places

**Impact**: Runtime error risk, lack of type safety
**Priority**: Medium
**Action**: Change to type-safe implementation where possible

### 10. Commented-Out Code

**Location**: Multiple locations
**Problem**:

- Commented-out code remains
- Debug code remains

**Impact**: Reduced code readability, source of confusion
**Priority**: Medium
**Action**: Remove unnecessary code

**Examples found**:

- `pkg/testutil/xrp.go` (commented `log.Fatalf`)
- `pkg/config/config_test.go` (commented `log.Fatalf`)
- `pkg/wallet/key/random_wallet.go:31` (commented `log.Printf`)

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

- ✅ Major dependencies updated to recent versions
- ⚠️ Regular update process not automated
- ⚠️ Security scanning not in CI

**Problem**:

- Security patches may not be applied promptly
- New features may not be available

**Action**:

- [ ] Set up Dependabot
- [ ] Add security scanning to CI
- [ ] Establish regular update schedule

### 15. Architecture Complexity

**Current State**:

- ✅ Container-based DI pattern exists (`pkg/di/container.go`)
- ⚠️ Container implementation has many `panic` calls
- ⚠️ Layer separation could be improved

**Problem**:

- Container pattern implementation is complex
- Dependency injection may be overly complex
- Layer separation insufficient

**Action**: Review and simplify architecture

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

### 18. Test Separation

**Status**: ✅ Completed

- Integration tests now use build tags

### 19. Private Key Management

**Location**: README.md
**Problem**: Private key handling is inappropriate
**Action**:

- [ ] Improve password input
- [ ] Implement encrypted storage
- [ ] Implement memory protection

### 20. Fee Issue

**Location**: README.md
**Problem**: Overpaying fees on Signet
**Action**: Review fee calculation

---

## Priority Matrix

| Issue | Priority | Impact | Effort | Priority Order |
|-------|----------|--------|--------|----------------|
| Security vulnerabilities | Critical | High | Medium | 1 |
| `panic` usage | High | High | Low | 2 |
| `log.Fatal` usage | High | High | Low | 3 |
| Error handling | High | Medium | High | 4 |
| Context management | High | Medium | High | 5 |
| Logging unification | High | Medium | Medium | 6 |
| Test separation | Medium | Low | Medium | 7 ✅ |
| Interface organization | Medium | Low | Medium | 8 |
| Documentation | Low | Low | High | 9 |

---

## Action Plan

### Phase 1 (Immediate)

1. ✅ Security vulnerability scanning setup (govulncheck available)
   - [ ] Add to CI pipeline
   - [ ] Set up Dependabot
2. [ ] Remove `panic` usage (replace with error returns)
3. [ ] Remove `log.Fatal` from non-main packages
4. [ ] Implement graceful shutdown in `main` functions

### Phase 2 (Within 1-2 weeks)

1. [ ] Standardize error handling
2. [ ] Complete context management implementation
3. [ ] Unify logging

### Phase 3 (Within 1-2 months)

1. [ ] Improve test data management
2. [ ] Organize interfaces
3. [ ] Improve documentation

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
