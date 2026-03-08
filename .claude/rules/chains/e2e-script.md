---
paths: ["scripts/operation/btc/e2e/**", "scripts/operation/bch/e2e/**", "scripts/operation/eth/e2e/**", "scripts/operation/xrp/e2e/**"]
---

# E2E Script Development Rules (All Chains)

Common rules for creating or modifying E2E scripts across all chains.
For chain-specific rules, see the per-chain files:

- BTC: `.claude/rules/chains/btc/e2e-script.md`
- BCH: `.claude/rules/chains/bch/e2e-script.md`
- ETH: `.claude/rules/chains/eth/e2e-script.md`
- XRP: `.claude/rules/chains/xrp/e2e-script.md`

## Standard Script Options

Every E2E script MUST support these options:

```bash
# Options:
#   --cleanup          Stop containers and cleanup state
#   --reset            Full reset and run from scratch
#   --verbose          Enable verbose output (set -x)
#   --non-interactive  Run without prompts (for CI/CD)
#   -h, --help         Display help message
```

## Database Configuration

E2E scripts support multiple database backends via the `DB_TYPE` environment variable:

| DB_TYPE                | Description        | Docker Required | Use Case              |
| ---------------------- | ------------------ | --------------- | --------------------- |
| `sqlite` (**default**) | Local SQLite files | No              | Fast testing, CI/CD   |
| `mysql`                | Docker MySQL       | Yes             | Full integration test |
| `postgres`             | Docker PostgreSQL  | Yes             | Full integration test |

**Always use SQLite for local development and CI** — it requires no Docker database container and is faster.

### When `DB_TYPE=sqlite`

1. `WALLET_DATABASE_TYPE=sqlite` is exported automatically by the common setup function
2. All wallet commands use SQLite instead of MySQL/PostgreSQL
3. No Docker database container is started

## Configuration File Policy

### ❌ Do NOT Edit Config Files Directly

Do **not** edit chain config files (e.g., `btc/watch.yaml`, `eth/keygen.yaml`) directly.
Use **environment variables** to override settings when different values are needed.

### ✅ Override via Environment Variables

```bash
# Priority order (applies to all chains):
# 1. Environment Variables  (highest priority)
# 2. Config File
# 3. Default Values         (lowest priority)
```

Config keys use the `WALLET_` prefix and map to config file fields.
See `pkg/config/README.md` for the full list of overridable keys.

## ⚠️ MANDATORY: Always Use Makefile Targets

**AI Agents and developers MUST use Makefile targets to run E2E tests.**
Do NOT execute E2E scripts directly.

```bash
# ✅ CORRECT: Use Makefile target
make btc-e2e-reset P=1
make eth-e2e-p1-reset

# ❌ WRONG: Do not run scripts directly
./scripts/operation/btc/e2e/e2e-p1-p2pkh-singlesig.sh --reset
./scripts/operation/eth/e2e/e2e-p1.sh --reset
```

### Why Makefile Targets?

1. **Automatic Build**: Targets depend on `build-all` (incremental — only rebuilds when Go sources change)
2. **Consistent Environment**: Properly sets `DB_TYPE`, `NODE_TYPE`, and other variables
3. **Validated Patterns**: Pattern numbers are validated before execution

## Verification Commands

### After Go Code Changes

```bash
# Required: lint, build, test
make go-lint && make check-build && make go-test

# Then run E2E (build is automatic via dependency)
make <chain>-e2e-<pattern>-reset
```

### After Shell Script Changes

```bash
make shfmt
```

## Retry Limit

**If the fix-test cycle exceeds 5 iterations, organize progress and report.**

### Escalation Conditions

- Same error occurs repeatedly
- Deep understanding of chain-specific protocol required
- Large-scale Go code changes needed

### Progress Report Format

```markdown
## Progress Report

### Error Details

[Error message that occurred]

### Attempted Fixes

1. [Fix attempt 1]
2. [Fix attempt 2]

### Current State

[Description of current state]

### Next Steps

[Required next actions]
```

## Security Rules

- ❌ Do NOT log private keys or passphrases
- ❌ Do NOT use test credentials in production configs
- Reference: `docs/guidelines/security.md`

## Avoiding Impact on Other Patterns

### When Modifying Shared `common.sh` or `{chain}_common.sh`

- Always verify impact on **all existing E2E patterns** for that chain
- Do not break existing pattern behavior
- Confirm regression with unit tests when modifying shared code

### Pattern-Specific Changes

- Set environment variables locally within each script
- Use isolated databases per pattern (e.g., pattern suffix in SQLite file path)

## Related Skills

| Skill             | Use Case                     |
| ----------------- | ---------------------------- |
| `shell-scripts`   | Script creation/modification |
| `go-development`  | Go code changes              |
| `makefile-update` | Makefile updates             |
| `git-workflow`    | Branch management            |
