# Claude Code Commands

## Commands

| Command | Description |
|---------|-------------|
| `/fix-issue #123` | Work on a GitHub issue |
| `/fix-linter` | Fix linter errors |
| `/fix-pr-review #123` | Address PR review comments |
| `/recreate-pr #123` | Copy PR with review comments as fix checklist |
| `/e2e/fix-btc-e2e-p1` | Fix BTC E2E test (Pattern 1) |
| `/e2e/fix-btc-e2e-p2` | Fix BTC E2E test (Pattern 2) |
| `/e2e/fix-btc-e2e-p8` | Fix BTC E2E test (Pattern 8: P2SH-P2WSH 3-of-3) |

## Workflow

```
1. Create Issue                     2. Work on Issue
   github-issue-creation               git-workflow + task skill
         ↓                                    ↓
   Classify: Type + Lang/Scope        Use skill based on label
   Assign labels
```

## Related

- Label → Skill mapping: See [github-issue-creation](../skills/github-issue-creation/SKILL.md)
- Git workflow: See [git-workflow](../skills/git-workflow/SKILL.md)
