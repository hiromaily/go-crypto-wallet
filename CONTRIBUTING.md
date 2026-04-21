<!--
⚠️ AUTO-GENERATED FILE — DO NOT EDIT
Source: template/pages/CONTRIBUTING.tpl.md · Run `make docs` to regenerate.
-->

# Contributing to go-crypto-wallet

Thank you for your interest in contributing to go-crypto-wallet!

## Getting Started

### Option 1: Using VS Code

1. **Open the Project**

   ```bash
   cd go-crypto-wallet
   code .
   ```

2. **Reopen in Container**
   - VS Code will detect the `.devcontainer` configuration
   - Click "Reopen in Container" when prompted
   - **OR** press `F1` → "Dev Containers: Reopen in Container"

3. **Wait for Setup**
   - First time: 5-10 minutes (downloads image, installs tools)
   - Subsequent times: 30-60 seconds (uses cached image)

4. **Start Developing**
   - Terminal opens inside the container
   - All tools are pre-installed and ready to use

### Option 2: Using Cursor

1. **Open the Project**

   ```bash
   cd go-crypto-wallet
   cursor .
   ```

2. **Reopen in Container**
   - Cursor automatically detects DevContainer configuration
   - Click "Reopen in Container" when prompted
   - **OR** use Command Palette → "Reopen in Container"

3. **Start Developing**
   - Cursor AI works seamlessly inside the container
   - All project tools are available

### Option 3: Using Claude Code (CLI)

1. **Navigate to Project**

   ```bash
   cd go-crypto-wallet
   ```

2. **Start Claude Code**

   ```bash
   claude-code
   ```

3. **Claude Code will:**
   - Detect the DevContainer configuration
   - Ask if you want to use it (optional)
   - Work inside the container if you choose

4. **Benefits with Claude Code:**
   - Claude can safely execute commands inside the container
   - All tools (make, golangci-lint, Atlas) are available
   - Git operations work seamlessly
   - Docker Compose commands work through host socket

---

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
make go-test        # Run tests
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

## Branch Naming

**Format**: `{type}/{brief-description}-{issue-number}`

| Type | Prefix | Example |
|------|--------|---------|
| Feature | `feature/` | `feature/add-taproot-support-123` |
| Bug fix | `fix/` | `fix/fee-calculation-error-456` |
| Refactor | `refactor/` | `refactor/clean-arch-layer-789` |
| Docs | `docs/` | `docs/update-readme-101` |
| CI/CD | `ci/` | `ci/add-workflow-200` |

## Pull Request Guidelines

### PR Title Format

```
<type>: <description>
```

Example: `feat: add taproot address support`

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
- [ ] Run `make go-test` - all tests pass
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

- [ARCHITECTURE.md](./ARCHITECTURE.md) - System architecture
- [docs/guidelines/](./docs/guidelines) - Project guidelines and standards

## Questions?

If you have questions, please open an issue for discussion before starting work on large changes.
