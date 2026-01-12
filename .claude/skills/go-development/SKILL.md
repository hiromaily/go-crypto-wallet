---
name: go-development
description: Common workflow for Go development tasks including branch management, verification commands, and self-review. Use this skill for any code changes in the go-crypto-wallet repository.
---

# Go Development Workflow

Standard workflow for all Go code changes in this repository.

## Branch Management

### Creating a New Branch

Always create a new branch from the latest `main`:

```bash
git fetch origin
git checkout main
git reset --hard origin/main
git checkout -b {branch-type}/issue-{number}-{brief-description}
```

### Branch Naming Convention

| Type | Format | Example |
|------|--------|---------|
| Feature | `feature/issue-{n}-{desc}` | `feature/issue-123-add-taproot` |
| Bug fix | `fix/issue-{n}-{desc}` | `fix/issue-456-db-connection` |
| Refactor | `refactor/issue-{n}-{desc}` | `refactor/issue-789-clean-arch` |
| Docs | `docs/{n}-{desc}` | `docs/311-commands-update` |

## Verification Commands

**ALWAYS run these commands before committing:**

```bash
make go-lint      # Lint check
make tidy         # go mod tidy
make check-build  # Build verification
make gotest       # Unit tests
```

### Additional Commands (When Applicable)

| Situation | Command |
|-----------|---------|
| Security changes | `make go-check-vuln` |
| Integration tests | `make gotest-integration` |
| DB schema changes | `make atlas-fmt && make atlas-lint` |
| SQLC regeneration | `make sqlc` |

## Self-Review Checklist

Before creating a commit, verify:

### Code Quality

- [ ] Follows Clean Architecture (domain has ZERO infrastructure deps)
- [ ] Error handling uses `fmt.Errorf("context: %w", err)`
- [ ] Import order: standard → third-party → local
- [ ] Exported functions have godoc comments

### Security

- [ ] No private keys or sensitive data logged
- [ ] No hardcoded secrets
- [ ] Security-critical changes reviewed

### Auto-Generated Files

- [ ] NOT editing files with `DO NOT EDIT` comments:
  - `internal/infrastructure/database/sqlc/*.go`
  - `internal/infrastructure/api/ripple/xrp/*.pb.go`
  - `internal/infrastructure/contract/token-abi.go`

## Commit Message Format

Use conventional commits:

```
{type}: {brief description}

- {detail 1}
- {detail 2}

Closes #{issue_number}
```

### Types

- `feat:` - New feature
- `fix:` - Bug fix
- `refactor:` - Code refactoring
- `docs:` - Documentation
- `test:` - Tests
- `chore:` - Maintenance

## Pull Request Creation

```bash
git push -u origin {branch-name}

gh pr create --title "{type}: {description}" --body "$(cat <<'EOF'
## Summary
- {change 1}
- {change 2}

## Test plan
- [ ] Verification commands pass
- [ ] Manual testing completed

Closes #{issue_number}
EOF
)"
```

## Quick Reference

### Standard Workflow

1. **Branch**: Create from latest `main`
2. **Implement**: Follow Clean Architecture
3. **Verify**: Run all verification commands
4. **Review**: Complete self-review checklist
5. **Commit**: Use conventional commit format
6. **PR**: Create with clear description

### Verification Summary

```bash
make go-lint && make tidy && make check-build && make gotest
```
