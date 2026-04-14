### Git Operations

**Allowed Operations:**

- `git add`: Stage changes for commit
- `git commit`: Commit staged changes
- `git push`: Push commits to GitHub (remote repository)

**NOT Allowed Operations:**

- ❌ `git merge` operations
- ❌ `gh` command merge operations (e.g., `gh pr merge`)
- ❌ `git commit` and `git push` to `main` branch

**Workflow:**

1. Create a feature branch for your work
2. Make changes and commit to the feature branch
3. Push the feature branch to GitHub
4. Create a pull request
5. Wait for review and approval
6. Let maintainers merge the pull request

**Commit Messages:**

This project uses [Conventional Commits](https://www.conventionalcommits.org/) format, enforced by `lefthook` pre-commit hooks.

Format:

```
<type>(<scope>): <description>

[optional body]

[optional footer]
```

**Types:**

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

**Scopes (Optional):**

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

**Examples:**

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

**Validation:**

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
