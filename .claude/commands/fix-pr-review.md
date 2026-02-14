# Fix PR Review #{pr_number}

Address PR review comments using the fix-pr-review skill.

## Process

1. **Fetch PR**: `gh pr view {pr_number}`
2. **Load fix-pr-review skill**: Determine what to fix from comments
   - See [fix-pr-review](../skills/fix-pr-review/SKILL.md)
3. **Classify by modified files**: Load appropriate development skills
4. **Follow loaded skill workflows**

## How It Works

```
fix-pr-review command
    │
    ├─ 1. Fetch PR details and review comments
    │
    ├─ 2. Load fix-pr-review skill
    │      │
    │      ├─ Classify by file type
    │      ├─ Determine development skills needed
    │      └─ Prioritize: security > functionality > quality
    │
    └─ 3. Execute workflow from loaded skills
```

## Key References

| Reference | Purpose |
|-----------|---------|
| [fix-pr-review](../skills/fix-pr-review/SKILL.md) | PR review workflow |
| [git-workflow](../skills/git-workflow/SKILL.md) | Git commit conventions |

## Example

```
/fix-pr-review #123
```

If PR #123 modifies Go files and has review comments:

1. Load `fix-pr-review` skill
2. Skill determines:
   - Modified files: `*.go` → Load `go-development` skill
   - Priority order: security → functionality → quality
3. Load `git-workflow` skill (always)
4. Load `go-development` skill
5. Follow workflow: checkout PR → fix comments → verify → commit → push
