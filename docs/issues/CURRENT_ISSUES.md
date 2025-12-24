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

**Impact**: Potential for known security vulnerabilities
**Priority**: Highest
**Action**:

- [ ] Add `govulncheck` to CI pipeline
- [ ] Set up Dependabot for automated dependency updates
- [ ] Run regular security scans

### 2. `panic` Usage

**Location**: Multiple files
**Current State**:

- `pkg/di/container.go` (22 occurrences)
  - Lines: 150, 209, 217, 229, 231, 421, 475, 486, 497, 508, 527, 543, 559, 587, 619, 757, 839, 973, 986, 999, 1032, 1059
  - **Note**: According to AGENTS.md, `panic` is acceptable in `pkg/di` package during instance construction phase
- `cmd/tools/get-eth-key/main.go` (3 occurrences)
  - Lines: 35, 45, 50
  - **Note**: `panic` in `main.go` files is acceptable, but graceful shutdown should be implemented

**Problem**: Runtime errors cause program crashes
**Impact**: Application stops unexpectedly on invalid coin types or configuration errors
**Priority**: Medium (for `pkg/di` - acceptable per project guidelines),
  High (for `main.go` - should implement graceful shutdown)
**Action**:

- `pkg/di`: Acceptable per project guidelines (instance construction phase)
- `main.go` files: Implement graceful shutdown instead of `panic`

### 3. Inappropriate `log.Fatal` Usage

**Location**: Multiple files
**Current State**:

- `cmd/tools/eth-key/main.go` (5 occurrences)
  - Lines: 23, 26, 29, 38, 44
- `pkg/testutil/repository.go` (18 occurrences)
  - Multiple test helper functions using `log.Fatalf`
- `pkg/testutil/xrp.go` (3 commented occurrences)
  - Lines: 61, 68, 71 (commented out)

**Note**: `log.Fatal` in `main` functions is acceptable,
  but should be replaced with proper error handling and graceful shutdown.

**Problem**: Using `log.Fatal` in non-main packages (especially `pkg/testutil/`)
  prevents cleanup processing
**Impact**: Resource leaks, potential data inconsistency,
  test failures don't provide proper cleanup
**Priority**: High
**Action**:

- Replace `log.Fatal` in `pkg/testutil/` with error returns
- Implement graceful shutdown in `main` functions instead of `log.Fatal`

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

- ✅ `pkg/infrastructure/api/ethereum/` - Most methods use `context.Context`
- ✅ `pkg/infrastructure/api/ripple/` - All gRPC methods use `context.Context`
- ✅ `pkg/infrastructure/api/bitcoin/` - Needs verification
- ✅ `pkg/infrastructure/network/websocket/` - Uses `context.Context`
- ⚠️ Timeout and cancellation not consistently implemented
- ⚠️ Tracing information propagation insufficient

**Impact**:

- Cannot control request timeouts consistently
- Distributed tracing doesn't work properly
- Potential resource leaks

**Priority**: Medium (most APIs already use context)
**Action**:

- [ ] Verify all `pkg/infrastructure/api/bitcoin/` methods use `context.Context`
- [ ] Implement consistent timeout settings across all API calls
- [ ] Implement cancellation handling
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
- ✅ 35 test files properly tagged with integration build tag
- ⚠️ Test data management could be improved
- ⚠️ Test helpers could be better organized
- ⚠️ `pkg/testutil/repository.go` uses `log.Fatal` which prevents proper test cleanup

**Impact**: Long test execution time, difficult test maintenance, improper test failure handling
**Priority**: Medium
**Action**:

- [x] Test separation with build tags (completed)
- [ ] Replace `log.Fatal` in `pkg/testutil/` with error returns
- [ ] Improve test data management
- [ ] Organize test helpers

### 8. Interface Overuse

**Location**: `pkg/wallet/service/` and `pkg/application/usecase/`
**Problem**:

- Many small interfaces defined
- Interface usage is inconsistent
- Both `pkg/wallet/service/` (legacy) and `pkg/application/usecase/` (new) exist in parallel

**Impact**: Increased code complexity, unclear which layer to use
**Priority**: Medium
**Action**:

- Review and consolidate interfaces
- Complete migration from `pkg/wallet/service/` to `pkg/application/usecase/`
- Document migration strategy

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

- `pkg/testutil/xrp.go` (3 commented `log.Fatalf` at lines 61, 68, 71)
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
- ✅ Domain layer separated (`pkg/domain/`)
- ✅ Infrastructure layer separated (`pkg/infrastructure/`)
- ✅ New application layer (`pkg/application/usecase/`) being introduced
- ⚠️ Container implementation has many `panic` calls (acceptable per project guidelines)
- ⚠️ Legacy application layer (`pkg/wallet/service/`) and new layer (`pkg/application/usecase/`) coexist
- ⚠️ Migration strategy unclear

**Problem**:

- Container pattern implementation is complex
- Dependency injection may be overly complex
- Two application layer implementations exist in parallel

**Action**:

- Complete migration from `pkg/wallet/service/` to `pkg/application/usecase/`
- Document architecture migration strategy
- Review and simplify architecture where possible

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
| ----- | -------- | ------ | ------ | -------------- |
| Security vulnerabilities | Critical | High | Medium | 1 |
| `log.Fatal` usage in testutil | High | High | Low | 2 |
| Error handling | High | Medium | High | 3 |
| Context management | Medium | Medium | Medium | 4 |
| Logging unification | High | Medium | Medium | 5 |
| Test separation | Medium | Low | Medium | 6 ✅ |
| Architecture migration | Medium | Medium | High | 7 |
| Interface organization | Medium | Low | Medium | 8 |
| Documentation | Low | Low | High | 9 |

---

## Action Plan

### Phase 1 (Immediate)

1. ✅ Security vulnerability scanning setup (govulncheck available)
   - [ ] Add to CI pipeline
   - [ ] Set up Dependabot
2. [ ] Remove `log.Fatal` from `pkg/testutil/` (replace with error returns)
3. [ ] Implement graceful shutdown in `main` functions
4. [ ] Remove commented-out code

### Phase 2 (Within 1-2 weeks)

1. [ ] Standardize error handling
2. [ ] Complete context management implementation (verify bitcoin API, add timeouts)
3. [ ] Unify logging
4. [ ] Document architecture migration strategy

### Phase 3 (Within 1-2 months)

1. [ ] Complete migration from `pkg/wallet/service/` to `pkg/application/usecase/`
2. [ ] Improve test data management
3. [ ] Organize interfaces
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
