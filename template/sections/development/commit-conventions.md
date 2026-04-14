## Commit Message Conventions

This project uses [Conventional Commits](https://www.conventionalcommits.org/) format.
Commit messages are validated by `lefthook` on every commit.

### Format

```
<type>(<scope>): <description>

[optional body]

[optional footer]
```

### Types

| Type | Description | Release Impact |
|------|-------------|----------------|
| `feat` | New feature | MINOR version |
| `fix` | Bug fix | PATCH version |
| `docs` | Documentation only | - |
| `refactor` | Code refactoring (no feature/fix) | - |
| `test` | Adding or updating tests | - |
| `ci` | CI/CD changes | - |
| `chore` | Maintenance tasks | - |
| `build` | Build system changes | - |
| `perf` | Performance improvements | PATCH version |
| `style` | Code style (formatting, etc.) | - |
| `revert` | Revert a previous commit | - |

### Scopes (Optional)

The following are suggested scopes, but other alphanumeric scopes (including hyphens and underscores) are also permitted.

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

### Breaking Changes

Add `!` after type/scope to indicate breaking changes:

```bash
feat(btc)!: change address format to bech32m only
fix!: remove deprecated API endpoint
```

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

# Breaking change
feat(btc)!: require bech32m for all taproot addresses

BREAKING CHANGE: Legacy address format no longer supported
```

### Validation Error

If your commit message doesn't follow the format, you'll see:

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
