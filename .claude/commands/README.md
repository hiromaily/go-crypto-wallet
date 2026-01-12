# Claude Code Commands

## Task Workflow

```
1. Create Issue                    2. Work on Issue
   ↓                                  ↓
   github-issue-creation           git-workflow + language skill
   - Type label                    - Branch from main
   - Language/Scope label          - Implement
   - Chain label (if needed)       - Verify (language-specific)
                                   - Commit & PR
```

## Commands

| Command | Description |
|---------|-------------|
| `/fix-issue #123` | Work on a GitHub issue |
| `/fix-linter` | Fix linter errors |
| `/fix-pr-review #123` | Address PR review comments |

## Skills

| Skill | Purpose |
|-------|---------|
| `github-issue-creation` | Create issues with proper classification |
| `git-workflow` | Branch, commit, PR workflow (all tasks) |
| `go-development` | Go verification & review |
| `typescript-development` | TypeScript verification & review |
| `solidity-development` | Solidity verification & review |

## Skill Composition

Most tasks use multiple skills:

```
git-workflow (common)
     +
language skill (based on label)
     =
complete workflow
```

### Examples

| Task | Skills Used |
|------|-------------|
| Go bug fix | `git-workflow` + `go-development` |
| TypeScript feature | `git-workflow` + `typescript-development` |
| Documentation update | `git-workflow` only |
| DevOps/CI change | `git-workflow` only |

## Task Classification (Labels)

### Type (required)

`bug`, `enhancement`, `refactoring`, `documentation`, `security`, `technical-debt`

### Language (for code tasks)

| Label | Skill |
|-------|-------|
| `lang:go` | `go-development` |
| `lang:typescript` | `typescript-development` |
| `lang:solidity` | `solidity-development` |

### Scope (for non-code tasks)

`scope:docs`, `scope:devops`, `scope:scripts`, `scope:makefile`, `scope:config`, `scope:db`

### Chain (if applicable)

`chain:btc`, `chain:bch`, `chain:eth`, `chain:erc20`, `chain:xrp`, `chain:all`
