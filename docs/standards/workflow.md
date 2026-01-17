# Development Workflow

Git operations and development workflow for go-crypto-wallet.

## Git Operations

### Allowed

```bash
git add .
git commit -m "message"
git push origin branch-name
```

### NOT Allowed

- ❌ `git merge`
- ❌ `gh pr merge`
- ❌ Push to `main` branch directly

## Branch Naming

| Type | Pattern | Example |
|------|---------|---------|
| Feature | `feature/{issue}-{desc}` | `feature/123-add-taproot` |
| Bug fix | `fix/{issue}-{desc}` | `fix/456-fee-calculation` |
| Refactor | `refactor/{issue}-{desc}` | `refactor/789-clean-arch` |
| Docs | `docs/{issue}-{desc}` | `docs/101-update-readme` |

## Commit Messages

This project uses [Conventional Commits](https://www.conventionalcommits.org/) format, enforced by `lefthook` pre-commit hooks.

### Format

```
<type>(<scope>): <description>

[optional body]

[optional footer]
```

### Types

| Type | Description | Release |
|------|-------------|---------|
| `feat` | New feature | MINOR |
| `fix` | Bug fix | PATCH |
| `docs` | Documentation only | - |
| `refactor` | Code refactoring (no feature/fix) | - |
| `test` | Adding or updating tests | - |
| `ci` | CI/CD changes | - |
| `chore` | Maintenance tasks | - |
| `build` | Build system changes | - |
| `perf` | Performance improvements | PATCH |
| `style` | Code style (formatting, etc.) | - |
| `revert` | Revert a previous commit | - |

### Scopes (Optional)

The following are suggested scopes, but other alphanumeric scopes are also permitted.

| Scope | Description |
|-------|-------------|
| `btc` | Bitcoin-related |
| `bch` | Bitcoin Cash-related |
| `eth` | Ethereum-related |
| `xrp` | XRP-related |
| `db` | Database-related |
| `api` | API-related |
| `cli` | CLI-related |
| `pr` | PR review fixes |

### Examples

```bash
# Feature with scope
feat(btc): add taproot address support

# Bug fix without scope
fix: resolve database connection timeout

# Documentation
docs: update architecture guide

# Refactoring with scope
refactor(api): reorganize endpoint handlers

# Breaking change (add ! after type/scope)
feat(btc)!: change address format to bech32m only
```

### Validation

Commit messages are validated by `lefthook` on every commit. If validation fails:

```
ERROR: Commit message does not follow Conventional Commits format.

Expected format: <type>(<scope>): <description>

Types: feat, fix, docs, refactor, test, ci, chore, build, perf, style, revert

Examples:
  feat(btc): add taproot address support
  fix: resolve database connection timeout
  docs: update architecture guide
  refactor(api): reorganize endpoint handlers

Your commit message: <your-invalid-commit-message>
```

## Verification Checklist

Before creating a PR:

```bash
# Go code
make go-lint       # ✅ No errors
make tidy          # ✅ No changes
make check-build   # ✅ Build succeeds
make gotest        # ✅ Tests pass

# Security (if applicable)
make go-check-vuln # ✅ No vulnerabilities
```

## Pull Request Template

```markdown
## Summary
Brief description of changes

## Changes
- Change 1
- Change 2

## Testing
- [ ] Unit tests pass
- [ ] Integration tests pass (if applicable)

Closes #XXX
```

## Auto-Generated Files

**DO NOT EDIT** files with `DO NOT EDIT` comments:

- `internal/infrastructure/database/sqlc/*.go`
- `internal/infrastructure/api/xrp/xrp/*.pb.go`
- `internal/infrastructure/contract/token-abi.go`

## Detailed Guidelines

See [docs/guidelines/workflow.md](../guidelines/workflow.md) for full workflow documentation.
