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

Use conventional commits:

```
type: brief description

- Detail 1
- Detail 2

Closes #123
```

Types: `feat`, `fix`, `refactor`, `docs`, `test`, `chore`

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
