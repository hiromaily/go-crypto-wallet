# Claude Code Commands

## Workflow Overview

```
1. Create Issue                     2. Work on Issue
   github-issue-creation               git-workflow + task skill
         ↓                                    ↓
   Classify: Type + Lang/Scope        Use skill based on label
   Assign labels                      (see mapping below)
```

## Commands

| Command | Description |
|---------|-------------|
| `/fix-issue #123` | Work on a GitHub issue |
| `/fix-linter` | Fix linter errors |
| `/fix-pr-review #123` | Address PR review comments |
| `/fix-btc-e2e` | Fix BTC E2E test (Pattern 8: P2SH-P2WSH 3-of-3) |

## Label → Skill Mapping

### Language Labels (code tasks)

| Label | Skill |
|-------|-------|
| `lang:go` | `go-development` |
| `lang:typescript` | `typescript-development` |
| `lang:solidity` | `solidity-development` |

### Scope Labels (non-code tasks)

| Label | Skill |
|-------|-------|
| `scope:docs` | `docs-update` |
| `scope:devops` | `devops` |
| `scope:scripts` | `shell-scripts` |
| `scope:makefile` | `makefile-update` |
| `scope:db` | `db-migration` |

## All Skills

| Skill | Purpose |
|-------|---------|
| `github-issue-creation` | Task classification |
| `git-workflow` | Branch/commit/PR (all tasks) |
| `go-development` | Go verification |
| `typescript-development` | TypeScript verification |
| `solidity-development` | Solidity verification |
| `docs-update` | Documentation workflow |
| `devops` | CI/CD workflow |
| `shell-scripts` | Shell script workflow |
| `makefile-update` | Makefile workflow |
| `db-migration` | Database change workflow |

## Skill Composition

Every task uses:
```
git-workflow (common)
     +
task-specific skill (based on label)
```

### Examples

| Task | Labels | Skills |
|------|--------|--------|
| Go bug fix | `bug`, `lang:go` | `git-workflow` + `go-development` |
| Docs update | `docs`, `scope:docs` | `git-workflow` + `docs-update` |
| Add CI workflow | `enhancement`, `scope:devops` | `git-workflow` + `devops` |
| DB migration | `enhancement`, `scope:db` | `git-workflow` + `db-migration` + `go-development` |
