# Fix PR Review #{pr_number}

Address review comments by selecting appropriate skills based on modified files.

## Task Classification

### Identify Language from PR Files

| Files Modified | Skill to Use |
|----------------|--------------|
| `internal/`, `pkg/`, `cmd/`, `*.go` | `go-development` |
| `apps/ripple-lib-server/`, `*.ts` | `typescript-development` |
| `apps/erc20-token/contracts/`, `*.sol` | `solidity-development` |

## Process

1. **Fetch PR**: `gh pr view {pr_number}`
2. **Get comments**: `gh api repos/{owner}/{repo}/pulls/{pr_number}/comments`
3. **Classify**: Identify language from modified files
4. **Load Skill**: Use appropriate `{lang}-development` Skill
5. **Fix**: Address each comment (priority: security > functionality > quality)
6. **Verify**: Run Skill-specific verification commands
7. **Push**: `git push` (updates existing PR)

## Skills Reference

| Skill | Path |
|-------|------|
| Go | `.claude/skills/go-development/SKILL.md` |
| TypeScript | `.claude/skills/typescript-development/SKILL.md` |
| Solidity | `.claude/skills/solidity-development/SKILL.md` |

## Example

```
/fix-pr-review #123
```
