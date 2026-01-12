---
name: go-development
description: Go development workflow including verification commands and self-review checklist. Use when modifying Go code in internal/, pkg/, or cmd/ directories.
---

# Go Development Workflow

Workflow for Go code changes in this repository.

## Prerequisites

**Use `git-workflow` Skill** for branch management, commit conventions, and PR creation.

## Applicable Directories

| Directory | Description |
|-----------|-------------|
| `internal/` | Core application code (Clean Architecture) |
| `pkg/` | Reusable shared packages |
| `cmd/` | Application entry points |

## Verification Commands

**ALWAYS run before committing:**

```bash
make go-lint      # Lint check
make tidy         # go mod tidy
make check-build  # Build verification
make gotest       # Unit tests
```

### Quick One-Liner

```bash
make go-lint && make tidy && make check-build && make gotest
```

### Additional Commands

| Situation | Command |
|-----------|---------|
| Security changes | `make go-check-vuln` |
| Integration tests | `make gotest-integration` |
| Format code | `make go-fmt` |

## Self-Review Checklist

### Clean Architecture

- [ ] Domain layer has ZERO infrastructure dependencies
- [ ] Dependencies flow inward (interface-adapters → application → domain)
- [ ] Use cases are in `internal/application/usecase/`
- [ ] Interfaces defined in application layer, implemented in infrastructure

### Code Quality

- [ ] Error handling uses `fmt.Errorf("context: %w", err)`
- [ ] Import order: standard → third-party → local
- [ ] Exported functions have godoc comments
- [ ] No `panic()` except for truly unrecoverable situations

### Security

- [ ] No private keys or sensitive data logged
- [ ] No hardcoded secrets
- [ ] Security-critical changes reviewed
- [ ] Consider offline wallet impact (keygen, sign)

### Auto-Generated Files

**DO NOT EDIT** files with `DO NOT EDIT` comments:

- `internal/infrastructure/database/sqlc/*.go`
- `internal/infrastructure/api/ripple/xrp/*.pb.go`
- `internal/infrastructure/contract/token-abi.go`

## Database Changes

If modifying database schema:

```bash
# 1. Modify HCL schema
# tools/atlas/schemas/

# 2. Format and lint
make atlas-fmt && make atlas-lint

# 3. Generate migrations
make atlas-dev-reset

# 4. Regenerate SQLC
make sqlc

# 5. Verify
make check-build && make gotest
```

## Related Skills

- `git-workflow` - Branch, commit, PR workflow
- `github-issue-creation` - Task classification
