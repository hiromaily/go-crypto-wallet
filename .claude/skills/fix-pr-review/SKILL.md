---
name: fix-pr-review
description: Address PR review comments by selecting appropriate skills based on modified files. Use when fixing review feedback on pull requests.
---

# Fix PR Review Workflow

Workflow for addressing review comments on pull requests.

## Input Formats

User may provide either:
- **PR number**: `563`, `#563`
- **Full review URL**: `https://github.com/{owner}/{repo}/pull/{pr_number}#pullrequestreview-{review_id}`

Extract `{pr_number}` (and optionally `{review_id}`) from the input.

## Prerequisites

- PR number or review URL is required
- Use `git-workflow` Skill for commit conventions

## Process Overview

```
1. Parse input (PR number or review URL)
2. Fetch PR details and review comments
3. Classify by file type
4. Load appropriate development skill
5. Fix each comment
6. Run verification
7. Display summary table
8. Commit and push changes
```

## Step 1: Fetch PR Information

```bash
# Get PR details
gh pr view {pr_number}

# Get PR diff to see modified files
gh pr diff {pr_number} --name-only

# Get review comments (all reviews)
gh api repos/{owner}/{repo}/pulls/{pr_number}/reviews

# Get review comments for a specific review (if review_id is provided)
gh api repos/{owner}/{repo}/pulls/{pr_number}/reviews/{review_id}/comments

# Get all inline comments on the PR
gh api repos/{owner}/{repo}/pulls/{pr_number}/comments
```

## Step 2: Classify by Modified Files

| Files Modified                         | Development Skill        |
| -------------------------------------- | ------------------------ |
| `internal/`, `pkg/`, `cmd/`, `*.go`    | `go-development`         |
| `apps/xrpl-grpc-server/`, `*.ts`       | `typescript-development` |
| `apps/erc20-token/contracts/`, `*.sol` | `solidity-development`   |
| `scripts/`, `*.sh`                     | `shell-scripts`          |
| `Makefile`, `make/*.mk`                | `makefile-update`        |
| `tools/atlas/`, `*.sql`, `*.hcl`       | `db-migration`           |
| `docs/`, `*.md`                        | `docs-update`            |

## Step 3: Load Development Skill

Based on classification, load the appropriate skill:

| Skill      | Path                                             |
| ---------- | ------------------------------------------------ |
| Go         | `.claude/skills/go-development/SKILL.md`         |
| TypeScript | `.claude/skills/typescript-development/SKILL.md` |
| Solidity   | `.claude/skills/solidity-development/SKILL.md`   |
| Shell      | `.claude/skills/shell-scripts/SKILL.md`          |
| Makefile   | `.claude/skills/makefile-update/SKILL.md`        |
| Database   | `.claude/skills/db-migration/SKILL.md`           |
| Docs       | `.claude/skills/docs-update/SKILL.md`            |

## Step 4: Address ALL Comments

**IMPORTANT**: Every review comment must be addressed. Do not skip or miss any comment.

### Procedure

1. **List all comments** - Enumerate every review comment from the API response
2. **Sort by priority** - Security > Functionality > Code Quality
3. **Fix each one** - Address every comment, one by one
4. **Track status** - Record the outcome for each comment (FIXED, ALREADY APPLIED, or SKIPPED)

### Priority Order

1. **Security** - Address security concerns first
2. **Functionality** - Fix bugs or logic issues
3. **Code Quality** - Style, naming, documentation

### Comment Categories

| Category               | Action                                |
| ---------------------- | ------------------------------------- |
| Bug fix requested      | Fix the issue, add test if applicable |
| Refactoring suggestion | Apply if improves readability         |
| Style/naming           | Follow project conventions            |
| Documentation          | Add/update comments                   |
| Question               | Respond in code comment or PR         |

## Step 5: Run Verification

Run verification commands from the loaded development skill:

### Go Files

```bash
make go-lint && make tidy && make check-build && make gotest
```

### TypeScript Files

```bash
cd apps/xrpl-grpc-server && bun run lint && bun run build && bun test
```

### Shell Scripts

```bash
make shfmt
```

### Multiple Languages

If PR modifies multiple file types, run all applicable verification commands.

## Step 6: Display Summary

After fixing all comments, display a summary table:

```
Summary

┌─────────────────────────┬────────────────────┬──────────────────────────────┐
│     Review Comment      │       Status       │           Details            │
├─────────────────────────┼────────────────────┼──────────────────────────────┤
│ [comment description]   │ ✅ FIXED           │ [what was changed]           │
├─────────────────────────┼────────────────────┼──────────────────────────────┤
│ [comment description]   │ ✅ ALREADY APPLIED │ Applied in commit [hash]     │
├─────────────────────────┼────────────────────┼──────────────────────────────┤
│ [comment description]   │ ⏭️ SKIPPED         │ [reason for skipping]        │
└─────────────────────────┴────────────────────┴──────────────────────────────┘
```

### Status Values

| Status             | Meaning                                          |
| ------------------ | ------------------------------------------------ |
| ✅ FIXED           | Comment addressed with new code changes           |
| ✅ ALREADY APPLIED | Already fixed in a previous commit                |
| ⏭️ SKIPPED         | Not applicable or intentionally not addressed     |

## Step 7: Commit and Push

```bash
# Stage changes
git add .

# Commit with descriptive message
git commit -m "$(cat <<'EOF'
fix(scope): address PR review comments

- Fix: [specific fix 1]
- Fix: [specific fix 2]
- Update: [update description]
EOF
)"

# Push to update PR
git push
```

## Common Review Patterns

### Error Handling (Go)

**Review comment**: "Missing error context"

```go
// Before
return err

// After
return fmt.Errorf("failed to create wallet: %w", err)
```

### Logging Security

**Review comment**: "Don't log sensitive data"

```go
// Before
log.Info("Processing key", "key", privateKey)

// After
log.Info("Processing key", "keyID", keyID)
```

### Code Duplication

**Review comment**: "Extract to helper function"

```go
// Before: Duplicated code in multiple places

// After: Create helper function
func helperFunction(params) result {
    // extracted logic
}
```

## Related Skills

- `go-development` - Go code changes
- `typescript-development` - TypeScript changes
- `shell-scripts` - Shell script changes
- `git-workflow` - Branch and commit workflow
