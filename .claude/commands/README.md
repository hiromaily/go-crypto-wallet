# Claude Code Commands

## Task Workflow

```
1. Create Issue (classify task)    2. Work on Issue (use appropriate skill)
         ↓                                    ↓
   github-issue-creation              language/scope-based skill
   - Determine type                   - go-development
   - Assign labels                    - typescript-development
   - Set scope/chain                  - solidity-development
```

## Commands

| Command | Description |
|---------|-------------|
| `/fix-issue #123` | Work on a GitHub issue |
| `/fix-linter` | Fix linter errors |
| `/fix-pr-review #123` | Address PR review comments |

## Task Classification (at Issue Creation)

### Type Labels (required - pick one)

| Label | Description |
|-------|-------------|
| `bug` | Something isn't working |
| `enhancement` | New feature |
| `refactoring` | Code improvement |
| `documentation` | Docs updates |
| `security` | Security-related |
| `technical-debt` | Code quality |

### Language Labels (for code tasks)

| Label | Skill | Directories |
|-------|-------|-------------|
| `lang:go` | `go-development` | `internal/`, `pkg/`, `cmd/` |
| `lang:typescript` | `typescript-development` | `apps/ripple-lib-server/` |
| `lang:solidity` | `solidity-development` | `apps/erc20-token/contracts/` |

### Scope Labels (for non-code tasks)

| Label | Directories |
|-------|-------------|
| `scope:docs` | `docs/`, `*.md` |
| `scope:devops` | `.github/workflows/`, `docker/` |
| `scope:scripts` | `scripts/`, `*.sh` |
| `scope:makefile` | `Makefile`, `make/` |
| `scope:config` | `config/`, `*.toml` |
| `scope:db` | `tools/atlas/`, `tools/sqlc/` |

### Chain Labels (if applicable)

| Label | Chain |
|-------|-------|
| `chain:btc` | Bitcoin |
| `chain:bch` | Bitcoin Cash |
| `chain:eth` | Ethereum |
| `chain:erc20` | ERC-20 tokens |
| `chain:xrp` | XRP/Ripple |
| `chain:all` | Cross-chain |

## Skills

| Skill | Purpose |
|-------|---------|
| `github-issue-creation` | Create issues with proper classification |
| `go-development` | Go code workflow |
| `typescript-development` | TypeScript/JS workflow |
| `solidity-development` | Solidity contract workflow |
