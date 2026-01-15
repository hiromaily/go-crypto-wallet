# Recreate PR #{pr_number}

Create a copy of an existing PR with review comments summarized as fixes.
This allows Gemini (or other reviewers) to re-review the changes.

## Purpose

- Copy an existing PR to get a fresh review
- Consolidate all commits into one
- Summarize review comments as actionable fixes in the new PR description

## Process

### Step 1: Fetch Original PR Information

```bash
# Get PR details (title, body, branch)
gh pr view {pr_number} --json title,body,headRefName,baseRefName,files

# Get review comments
gh api repos/{owner}/{repo}/pulls/{pr_number}/comments

# Get PR reviews
gh api repos/{owner}/{repo}/pulls/{pr_number}/reviews

# Get PR diff
gh pr diff {pr_number}
```

### Step 2: Checkout Original Branch and Create New Branch

```bash
# Fetch the original PR branch
gh pr checkout {pr_number}

# Get the original branch name
ORIGINAL_BRANCH=$(git branch --show-current)

# Switch to main and update
git checkout main
git fetch origin
git reset --hard origin/main

# Create new branch with -v2 suffix
NEW_BRANCH="${ORIGINAL_BRANCH}-v2"
git checkout -b ${NEW_BRANCH}
```

### Step 3: Apply Changes as Single Commit

Option A: Cherry-pick and squash (if original has multiple commits)
```bash
# Get commits from original PR
gh pr view {pr_number} --json commits

# Cherry-pick all commits
git cherry-pick <commit1> <commit2> ...

# Squash into single commit
git reset --soft HEAD~N  # N = number of commits
git commit -m "squashed commit message"
```

Option B: Apply diff directly
```bash
# Apply the PR diff to the new branch
gh pr diff {pr_number} | git apply

# Commit all changes
git add -A
git commit -m "consolidated changes from PR #{pr_number}"
```

### Step 4: Analyze Review Comments

Extract and categorize review comments:

| Category | Action |
|----------|--------|
| Bug/Error | Must fix |
| Improvement | Should address |
| Style/Naming | Apply convention |
| Question | Document or clarify |

### Step 5: Create New PR

Create PR with original message + review summary:

```bash
gh pr create --title "{original_title} (v2)" --body "$(cat <<'EOF'
## Summary

{original_pr_body}

---

## Review Feedback from PR #{pr_number}

### Fixes Applied

- [ ] {review_comment_1_summary}
- [ ] {review_comment_2_summary}
- [ ] {review_comment_3_summary}

### Original Review Comments

{formatted_review_comments}

---

## Test Plan

- [ ] Verification commands pass
- [ ] Review comments addressed
- [ ] Manual testing completed

Related: #{pr_number}
EOF
)"
```

## Output Format

After creating the new PR, report:

```
Created PR #{new_pr_number} from PR #{pr_number}

Original PR: #{pr_number}
New PR: #{new_pr_number}
New Branch: {branch_name}

Review Comments Summarized:
- {comment_count} comments from {reviewer_count} reviewers
- Key fixes:
  1. {fix_1}
  2. {fix_2}
  ...
```

## Example Usage

```
/recreate-pr #123
```

This will:
1. Fetch PR #123 details and review comments
2. Create new branch `{original-branch}-v2`
3. Apply all changes as single commit
4. Create new PR with:
   - Original PR title + " (v2)"
   - Original PR description
   - Summary of review comments as checklist

## Notes

- The new PR will reference the original PR
- Review comments are formatted as actionable items
- All commits are squashed into one for cleaner history
- Use this when you want a fresh review on the same changes
