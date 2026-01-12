# Fix Issue #{issue_number}

Work on a GitHub issue using Skills determined by its labels.

## Process

1. **Fetch issue**: `gh issue view {issue_number}`
2. **Check labels**: Identify `lang:*` or `scope:*` label
3. **Load Skills**:
   - Always: `git-workflow`
   - Based on label: See mapping below
4. **Follow Skill workflows**

## Label → Skill Mapping

### Language Labels

| Label | Skill |
|-------|-------|
| `lang:go` | `go-development` |
| `lang:typescript` | `typescript-development` |
| `lang:solidity` | `solidity-development` |

### Scope Labels

| Label | Skill |
|-------|-------|
| `scope:docs` | `docs-update` |
| `scope:devops` | `devops` |
| `scope:scripts` | `shell-scripts` |
| `scope:makefile` | `makefile-update` |
| `scope:db` | `db-migration` (+ `go-development` for SQLC) |

## Verification by Skill

| Skill | Commands |
|-------|----------|
| `go-development` | `make go-lint && make tidy && make check-build && make gotest` |
| `typescript-development` | `npm run lint && npm run build && npm test` |
| `solidity-development` | `truffle compile && truffle test` |
| `docs-update` | Check markdown formatting |
| `devops` | Validate YAML, test workflow |
| `shell-scripts` | `make shfmt` |
| `makefile-update` | `make mk-lint` |
| `db-migration` | `make atlas-lint && make sqlc && make check-build` |

## Example

```
/fix-issue #123
```

If issue #123 has labels `bug`, `lang:go`, `chain:btc`:
1. Load `git-workflow` Skill
2. Load `go-development` Skill
3. Create branch from main
4. Implement fix
5. Run Go verification commands
6. Commit and create PR
