# Contributing to go-crypto-wallet

Thank you for your interest in contributing to go-crypto-wallet!

## Getting Started

1. Fork the repository
2. Clone your fork: `git clone https://github.com/YOUR_USERNAME/go-crypto-wallet.git`
3. Create a feature branch: `git checkout -b feature/issue-123-description`
4. Make your changes
5. Run verification commands
6. Commit with conventional commit message
7. Push and create a Pull Request

## Development Setup

### Prerequisites

- Go 1.25+
- Docker & Docker Compose
- Make

### Install Dependencies

```bash
make setup
```

### Run Verification

```bash
make go-lint       # Lint check
make check-build   # Build verification
make gotest        # Run tests
```

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

## Branch Naming

| Type | Pattern | Example |
|------|---------|---------|
| Feature | `feature/issue-{number}-{desc}` | `feature/issue-123-add-taproot` |
| Bug fix | `fix/issue-{number}-{desc}` | `fix/issue-456-fee-calculation` |
| Refactor | `refactor/issue-{number}-{desc}` | `refactor/issue-789-clean-arch` |
| Docs | `docs/issue-{number}-{desc}` | `docs/issue-101-update-readme` |
| CI/CD | `ci/issue-{number}-{desc}` | `ci/issue-200-add-workflow` |

## Pull Request Guidelines

### PR Title Format

```
<type>: <description> (Closes #<issue_number>)
```

Example: `feat: add taproot address support (Closes #123)`

### PR Description Template

```markdown
## Summary
Brief description of changes

## Changes
- Change 1
- Change 2

## Test Plan
- [ ] Unit tests pass
- [ ] Integration tests pass (if applicable)
- [ ] Manual testing completed

Closes #XXX
```

### Before Submitting

- [ ] Run `make go-lint` - no errors
- [ ] Run `make check-build` - build succeeds
- [ ] Run `make gotest` - all tests pass
- [ ] Commit messages follow Conventional Commits format
- [ ] PR title follows the format above

## Code Style

- Follow existing code patterns
- Use `gofmt` for formatting (handled by linter)
- Write meaningful comments for complex logic
- Keep functions focused and small

## Security Considerations

This project handles private keys and cryptocurrency transactions.

- **NEVER** log private keys or sensitive information
- **NEVER** commit secrets or credentials
- Consider offline wallet implications for keygen/sign operations
- Security-related changes require careful review

## Auto-Generated Files

**DO NOT EDIT** files containing `DO NOT EDIT` comments:

- `internal/infrastructure/database/mysql/sqlcgen/*.go`
- `internal/infrastructure/database/sqlite/sqlcgen/*.go`
- `internal/infrastructure/api/xrp/xrp/*.pb.go`
- `internal/infrastructure/contract/token-abi.go`

## Documentation

- [ARCHITECTURE.md](ARCHITECTURE.md) - System architecture
- [docs/standards/](docs/standards/) - Coding standards
- [docs/guidelines/](docs/guidelines/) - Detailed guidelines

## Questions?

If you have questions, please open an issue for discussion before starting work on large changes.
