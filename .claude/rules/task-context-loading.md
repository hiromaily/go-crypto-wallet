# Task-Oriented Context Loading

Rules for automatically loading appropriate context documents when receiving a task.

**SSOT Reference**: See [Task Classification](../docs/standards/task-classification.md) for the authoritative definition of labels and task types.

## Target Files

This rule applies when editing the following file types:

- `**/*.go`
- `**/*.md`
- `**/*.toml`
- `**/*.yaml`
- `**/*.hcl`
- `**/*.sql`
- `**/*.sh`
- `**/*.proto`

## Task Detection Priority

1. **GitHub Issue Labels** (if working on an issue) → Highest priority
2. **Explicit User Specification** → Second priority
3. **Keyword Detection** → Fallback

## Task Type from GitHub Labels

When working on a GitHub Issue, use labels to determine task type:

| Label | Task Type | Context File |
|-------|-----------|--------------|
| `bug` | bug | `docs/task-contexts/bug-fix.md` |
| `enhancement` | feature-add | `docs/task-contexts/feature-add.md` |
| `refactoring` | refactoring | `docs/task-contexts/refactoring.md` |
| `documentation` | documentation | `docs/task-contexts/documentation.md` |
| `security` | security | `docs/task-contexts/security.md` |
| `technical-debt` | refactoring | `docs/task-contexts/refactoring.md` |
| `test` | test | `docs/task-contexts/test.md` |

## Task Type Detection (Keyword Fallback)

If no GitHub labels are available, determine task type from keywords:

| Keywords | Task Type | Context File |
|----------|-----------|--------------|
| bug, fix, error, issue | `bug` | `docs/task-contexts/bug-fix.md` |
| add, implement, feature, new | `enhancement` | `docs/task-contexts/feature-add.md` |
| refactor, reorganize, move, cleanup | `refactoring` | `docs/task-contexts/refactoring.md` |
| schema, DB, table, column, migration | `db-change` | `docs/task-contexts/db-change.md` |
| document, README, description, docs, comment | `documentation` | `docs/task-contexts/documentation.md` |
| test, coverage, spec, unit test, integration test | `test` | `docs/task-contexts/test.md` |
| security, vulnerability, CVE | `security` | `docs/task-contexts/security.md` |

## Chain Detection

Identify the chain from the following keywords:

| Keywords | Chain | Context Files |
|----------|-------|---------------|
| Bitcoin, BTC, Taproot, Descriptor, PSBT, MuSig2 | BTC | `docs/task-contexts/chain-specific.md`, `docs/task-contexts/chains/btc.md` |
| Bitcoin Cash, BCH, CashAddr | BCH | `docs/task-contexts/chain-specific.md`, `docs/task-contexts/chains/bch.md` |
| Ethereum, ETH, ERC-20, Gas, Nonce | ETH | `docs/task-contexts/chain-specific.md`, `docs/task-contexts/chains/eth.md` |
| Ripple, XRP, Destination Tag | XRP | `docs/task-contexts/chain-specific.md`, `docs/task-contexts/chains/xrp.md` |

## File Type Detection for Verification

Determine appropriate verification commands based on the edited file extension:

| File Extension | Required Verification | Optional |
|----------------|----------------------|----------|
| `*.go` | `make go-lint`, `make check-build` | `make gotest` |
| `*.md`, `*.mdc` | (none) | markdownlint |
| `*.sql`, `*.hcl` | `make atlas-fmt`, `make atlas-lint` | |
| `*.yaml`, `*.toml` | (none) | |
| `*.sh` | (none) | shellcheck |
| `*.proto` | `make proto` | |

**Important**: Do **NOT** run Go-related verification commands for documentation-only (`*.md`) changes.

See `docs/task-contexts/verification.md` for details.

## Context Loading Procedure

1. **Determine task type**: Identify task type from user request
2. **Identify chain**: Identify chain if cryptocurrency-related
3. **Load context**: Load the relevant documents
4. **Detect file type**: Check the extension of files to edit
5. **Determine verification commands**: Select verification commands based on file type

```
Order:
1. docs/task-contexts/{task-type}.md (by task type)
2. docs/task-contexts/chain-specific.md (if chain-related)
3. docs/task-contexts/chains/{chain}.md (for specific chain)
4. docs/task-contexts/verification.md (verification commands)
5. Additional related documents (specified in each context file)
```

## Explicit Task Type Specification

If the user explicitly specifies the task type, follow that type:

```
Task Type: bug-fix, Chain: BCH. {description}
Task Type: feature-add, Chain: BTC. {description}
Task Type: db-change. {description}
Task Type: documentation. {description}
```

## Chain-Specific Rules

### BCH (Bitcoin Cash) - Important

For BCH tasks, always check the following rules:

1. **Do NOT modify BTC code directly**: Handle BCH-specific issues by overriding methods on the BCH side
2. **BitcoinCash struct**: Embeds `btc.Bitcoin`
3. **No SegWit/Taproot support**: Related features are not used in BCH

```go
// BCH implementation location
internal/infrastructure/api/bitcoin/bch/

// Override example
func (b *BitcoinCash) GetAddressInfo(addr string) (*dtobtc.AddressInfo, error) {
    // BCH-specific implementation
}
```

## Directory Structure Reference

### Use Case Layer

```
internal/application/usecase/
├── keygen/{btc,eth,xrp}/
├── sign/{btc,eth,xrp}/
└── watch/{btc,eth,xrp}/
```

### Infrastructure Layer

```
internal/infrastructure/api/
├── bitcoin/{btc,bch}/
├── ethereum/{eth,erc20}/
└── ripple/xrp/
```

### CLI Layer

```
internal/interface-adapters/cli/
├── keygen/api/{btc,eth}/
├── sign/
└── watch/api/{btc,eth,xrp}/
```

## Verification Commands by File Type

### Go Files (*.go)

```bash
make go-lint      # Required: Linter check
make check-build  # Required: Build verification
make gotest       # Recommended: Run tests (when functionality changes)
make tidy         # Recommended: Tidy dependencies (when imports change)
```

### Documentation Files (*.md)

```bash
# Go-related commands not required
# Optional: markdownlint (if installed)
```

### Database Files (*.sql, *.hcl)

```bash
make atlas-fmt    # Required: Format
make atlas-lint   # Required: Lint
make sqlc         # On schema change
```

### Shell Scripts (*.sh)

```bash
# Optional: shellcheck (if installed)
```

## Quick Decision Tree

```
Task Received
    │
    ├─ Editing Go files?
    │   └─ Yes → make go-lint, make check-build, (make gotest)
    │
    ├─ Documentation only?
    │   └─ Yes → No verification commands needed
    │
    ├─ DB schema change?
    │   └─ Yes → make atlas-fmt, make atlas-lint, make sqlc
    │
    └─ Other
        └─ Decide based on file type
```

## Related Documents

- [Task-Oriented Context Management](../../docs/task-oriented-context.md)
- [Task Contexts](../../docs/task-contexts/README.md)
- [Verification Matrix](../../docs/task-contexts/verification.md)
- [AGENTS.md](../../AGENTS.md) - Project-wide guidelines
