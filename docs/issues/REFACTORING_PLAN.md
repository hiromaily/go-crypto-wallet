# Refactoring Plan

## Overview

This project is a cryptocurrency wallet CLI application implemented in Go. This document outlines the refactoring and update plan to align with current best practices and address technical debt.

## Current State Analysis

### Technology Stack

- **Go**: 1.25.5 (updated from 1.20)
- **Major Dependencies**:
  - `ethereum/go-ethereum`: v1.16.7 (updated from v1.12.0)
  - `btcsuite/btcd`: v0.25.0 (updated from v0.23.4)
  - `sqlc`: Used for database code generation (not sqlboiler)
  - Tracing: Jaeger configuration exists but is optional (currently set to "none" in configs)

### Architecture

- **3 Wallet Types**: `watch`, `keygen`, `sign`
- **Multi-Coin Support**: BTC, BCH, ETH, ERC20, XRP
- **Dependency Injection**: Container-based DI pattern (`pkg/di/container.go`)
- **Repository Pattern**: Using sqlc for type-safe database queries
- **Package Layout**: Following `pkg` layout pattern

### Main Issues

1. **Error Handling**: Inappropriate use of `log.Fatal` in non-main packages
2. **Context Management**: Incomplete `context.Context` usage across API calls
3. **Testing**: Integration tests separated with build tags, but coverage needs improvement
4. **Logging**: Structured logging exists but needs consistency
5. **Security**: Private key and password management improvements needed
6. **Graceful Shutdown**: Not implemented in main functions

### Current Status

#### Completed ✅

- Go version updated to 1.25.5
- Major dependencies updated to recent versions
- Linter configuration (`.golangci.yml`) exists and is well-configured
- CI/CD workflow exists (`.github/workflows/lint-test.yml`)
- Integration tests use build tags (`//go:build integration`)
- Error wrapping uses `fmt.Errorf` + `%w`
- Context usage partially implemented (ethgrp uses it)
- Logger package exists (`pkg/logger/`)

#### In Progress / Needs Work 🔄

- `log.Fatal` removal: Many occurrences remain in test files and some packages
- Context management: Not all APIs use context (btcgrp, xrpgrp need verification)
- Graceful shutdown: Not implemented
- Security scanning: Not integrated into CI
- Error types: Need standardization
- Test coverage: Needs measurement and improvement

---

## Refactoring Phases

### Phase 1: Foundation Setup (Priority: High)

#### 1.1 Go Version and Dependency Updates

- [x] Go upgraded to 1.25.5
- [x] Major dependencies updated
  - [x] `ethereum/go-ethereum` → v1.16.7
  - [x] `btcsuite/btcd` → v0.25.0
- [x] Regular dependency updates
  - [x] Check for newer versions periodically
  - [x] Update if security fixes available
- [ ] Security vulnerability scanning and fixes
  - [x] `govulncheck` available in tools
  - [ ] Regular vulnerability checks
  - [ ] Add Dependabot configuration (`.github/dependabot.yml`)
  - [ ] Integrate security scanning into CI

#### 1.2 Build System Improvements

- [ ] Consider `go.work` introduction (if multi-module needed)
- [x] Review and optimize `Makefile`
- [x] CI/CD pipeline setup (GitHub Actions)
  - [x] Automated tests
  - [x] Lint checks
  - [ ] Security scanning (add `govulncheck`)
  - [ ] Build verification
  - [ ] Release workflow (if needed)

#### 1.3 Development Environment Setup

- [x] `.golangci.yml` configuration optimized
- [ ] Consider `pre-commit` hooks setup
- [ ] Optimize Docker Compose

---

### Phase 2: Code Quality Improvements (Priority: High)

#### 2.1 Error Handling Improvements

- [ ] Remove `log.Fatal` (except in `main` functions)

  Target files:
  - [ ] `cmd/keygen/main.go` (multiple occurrences)
  - [ ] `cmd/sign/main.go` (multiple occurrences)
  - [ ] `cmd/watch/main.go` (multiple occurrences)
  - [ ] `cmd/tools/eth-key/main.go` (multiple occurrences)
  - [ ] `pkg/testutil/repository.go` (multiple occurrences)
  - [ ] Test files in `pkg/repository/watchrepo/`
  - [ ] Test files in `pkg/account/`

- [x] Error wrapping standardization
  - [x] Using `fmt.Errorf` + `%w`
  - [ ] Add context information to error messages
  - [ ] Define custom error types
    - [ ] `ErrInvalidInput`
    - [ ] `ErrNotFound`
    - [ ] `ErrUnauthorized`
    - [ ] Domain-specific errors

- [ ] Error check coverage
  - [ ] Enable `errcheck` in `.golangci.yml` (currently disabled)
  - [ ] Fix unchecked errors

#### 2.2 Context Management

- [ ] Introduce `context.Context` comprehensively
  - [x] API calls in `ethgrp` use context
  - [ ] Verify all API calls in `btcgrp` use context
  - [ ] Verify all API calls in `xrpgrp` use context
  - [ ] Implement timeout and cancellation
  - [ ] Propagate tracing information

- [ ] Implement graceful shutdown
  - [ ] Signal handling in `main` functions
  - [ ] Resource cleanup
  - [ ] Wait for in-progress operations

#### 2.3 Logging Standardization

- [x] Structured logging package exists (`pkg/logger/`)
- [ ] Ensure logging consistency
  - [ ] Unify log levels
  - [ ] Unify log fields
  - [ ] Prevent logging sensitive information
  - [ ] Unify log format

#### 2.4 Code Style Standardization

- [x] `gofmt`/`goimports` automatic application (via linter)
- [ ] Unify naming conventions
- [ ] Add comments (godoc compliant)
- [ ] Remove unused code

---

### Phase 3: Architecture Improvements (Priority: Medium)

#### 3.1 Dependency Injection Improvements

- [x] Container-based DI pattern exists (`pkg/di/container.go`)
- [ ] Review container pattern
- [ ] Organize interfaces
  - [ ] Remove unused interfaces
  - [ ] Split interfaces appropriately
- [ ] Automate mock generation
  - [ ] Introduce `mockgen`
  - [ ] Generate mocks for interfaces
  - [ ] Use mocks in tests

#### 3.2 Layer Separation Clarification

- [x] Separate domain logic
- [x] Separate infrastructure layer
- [x] Organize application layer
  - See `ISSUE_ORGANIZE_APPLICATION_LAYER.md` for detailed implementation plan

#### 3.3 Configuration Management Improvements

- [ ] Prioritize environment variable configuration
- [ ] Strengthen configuration validation
- [ ] Improve configuration type safety

---

### Phase 4: Security Enhancements (Priority: High)

#### 4.1 Private Key Management

- [ ] Implement encrypted private key storage
- [ ] Protect private keys in memory (`mlock` usage)
- [ ] Zero-clear private keys (clear after use)
- [ ] Improve password input (`gopass` usage)

#### 4.2 Authentication & Authorization

- [ ] Strengthen API authentication
- [ ] Implement rate limiting
- [ ] Implement audit logging

#### 4.3 Security Best Practices

- [ ] Automate dependency security scanning
- [ ] Improve secret management (environment variables, secret managers)
- [ ] Strengthen input validation

---

### Phase 5: Test Improvements (Priority: Medium)

#### 5.1 Test Structure Organization

- [x] Unit and integration test separation (build tags used)
  - [x] `//go:build integration` tags added (35+ test files)
  - [x] Integration tests separated
- [ ] Organize test helpers
- [ ] Manage test data

#### 5.2 Test Coverage Improvement

- [ ] Generate coverage reports

  ```bash
  go test -coverprofile=coverage.out ./...
  go tool cover -html=coverage.out
  ```

- [ ] Set coverage goals (80%+)
- [ ] Add tests for critical paths

#### 5.3 Test Automation

- [x] Automated test execution in CI
- [ ] Build integration test environment (Docker Compose)
- [ ] Add performance tests

---

### Phase 6: Monitoring and Observability (Priority: Medium)

#### 6.1 Tracing Improvements

- [ ] Consider OpenTelemetry migration (Jaeger config exists but unused)
- [ ] Implement distributed tracing
- [ ] Propagate trace context

#### 6.2 Metrics

- [ ] Add Prometheus metrics
- [ ] Define business metrics
- [ ] Create dashboards

#### 6.3 Log Management

- [ ] Unify structured logging
- [ ] Configure log aggregation
- [ ] Set up alerts

---

### Phase 7: Performance Optimization (Priority: Low)

#### 7.1 Database

- [ ] Optimize queries
- [ ] Review indexes
- [ ] Optimize connection pools

#### 7.2 Memory Management

- [ ] Detect and fix memory leaks
- [ ] Optimize garbage collection
- [ ] Conduct profiling

#### 7.3 Concurrency

- [ ] Review concurrency patterns
- [ ] Detect and fix deadlocks
- [ ] Leverage context cancellation

---

### Phase 8: Documentation (Priority: Medium)

#### 8.1 Code Documentation

- [ ] Add package documentation
- [ ] Add function and method comments
- [ ] Add usage examples for complex logic

#### 8.2 Architecture Documentation

- [ ] Update architecture diagrams
- [ ] Record design decisions (ADR)
- [ ] Maintain API specifications

#### 8.3 Operational Documentation

- [ ] Update deployment procedures
- [ ] Create troubleshooting guide
- [ ] Maintain operational manual

---

## Implementation Priority

### Highest Priority (Immediate)

1. **Security vulnerability fixes**
2. **Error handling improvements** (`log.Fatal` removal)
3. **Context management** (complete implementation)
4. **Enable `errcheck` in linter**

### High Priority (Within 1-2 months)

1. **Logging standardization**
2. **Test structure organization**
3. **CI/CD pipeline enhancement** (add security scanning)
4. **Graceful shutdown implementation**

### Medium Priority (Within 3-6 months)

1. **Architecture improvements**
2. **Monitoring improvements**
3. **Documentation**

### Low Priority (6+ months)

1. **Performance optimization**
2. **Feature additions**

---

## Risk Management

### Technical Risks

- **Breaking Changes**: Compatibility issues from dependency updates
  - **Mitigation**: Gradual updates, sufficient testing
- **Performance Degradation**: Performance decline from refactoring
  - **Mitigation**: Benchmark tests, performance monitoring

### Operational Risks

- **Downtime**: Service interruption during releases
  - **Mitigation**: Gradual rollout, rollback plans
- **Data Migration**: Database schema changes
  - **Mitigation**: Migration script preparation, backups

---

## Success Metrics

### Code Quality

- [ ] Test coverage: 80%+
- [ ] Linter errors: 0
- [ ] Security vulnerabilities: 0 (Critical/High)

### Performance

- [ ] Response time: Maintain or improve
- [ ] Memory usage: Maintain or improve

### Maintainability

- [ ] Code review time: 20% reduction
- [ ] Bug fix time: 30% reduction
- [ ] New feature addition time: 25% reduction

---

## Timeline (Recommended)

### Q1 (Months 1-3)

- Phase 1: Foundation setup (complete)
- Phase 2: Code quality improvements (partial)
- Phase 4: Security enhancements (partial)

### Q2 (Months 4-6)

- Phase 2: Code quality improvements (complete)
- Phase 3: Architecture improvements
- Phase 5: Test improvements

### Q3 (Months 7-9)

- Phase 4: Security enhancements (complete)
- Phase 6: Monitoring and observability
- Phase 8: Documentation

### Q4 (Months 10-12)

- Phase 7: Performance optimization
- Final adjustments and review

---

## Reference Materials

### Go Related

- [Effective Go](https://go.dev/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Go Best Practices](https://golang.org/doc/effective_go)

### Security

- [OWASP Top 10](https://owasp.org/www-project-top-ten/)
- [Go Security Best Practices](https://go.dev/doc/security/best-practices)

### Architecture

- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Domain-Driven Design](https://martinfowler.com/bliki/DomainDrivenDesign.html)

---

## Update History

- 2024-XX-XX: Initial version created
- 2025-01-XX: Updated to reflect current state (Go 1.25.5, modern dependencies, current architecture)
