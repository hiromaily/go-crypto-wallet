---
name: git-workflow
description: Git branch management, commit conventions, and PR creation workflow. Use for all tasks that require code changes, regardless of language or scope.
---

# Git Workflow

Standard Git workflow for all tasks in this repository.

## Branch Management

### Creating a New Branch

Always create from latest `main`:

```bash
git fetch origin
git checkout main
git reset --hard origin/main
git checkout -b {type}/issue-{number}-{brief-description}
```

### Branch Naming Convention

| Type | Prefix | Example |
|------|--------|---------|
| Feature | `feature/` | `feature/issue-123-add-taproot` |
| Bug fix | `fix/` | `fix/issue-456-db-connection` |
| Refactoring | `refactor/` | `refactor/issue-789-clean-arch` |
| Documentation | `docs/` | `docs/311-update-readme` |
| DevOps/CI | `ci/` | `ci/issue-100-add-workflow` |
| Chore | `chore/` | `chore/issue-200-update-deps` |

### Branch Rules

- **Always from `main`**: Never branch from feature branches
- **One issue = One branch**: 1つのissueに対して1つのbranchのみ作成
- **Short-lived**: Merge within days, not weeks
- **Delete after merge**: Keep repository clean

### ⚠️ 重要: 複数ブランチ禁止ルール

**1つのissueに対して複数のbranchを作成しないでください。**

```
❌ 禁止パターン:
  issue-123 → fix/issue-123-first-attempt
            → fix/issue-123-second-attempt  ← これは作らない
            → fix/issue-123-another-fix     ← これも作らない

✅ 正しいパターン:
  issue-123 → fix/issue-123-description
            → PR作成 → レビュー → マージ
            → (必要なら) 新しいissueで新しいbranch
```

### 作業開始前の確認

**新しいbranchを作成する前に、必ず以下を確認:**

```bash
# 1. 既存のブランチを確認
git branch -a | grep "issue-{number}"

# 2. 既存のPRを確認
gh pr list --search "issue-{number}"
```

- **既存branchがある場合**: そのbranchで作業を継続
- **既存PRがある場合**: PRをmergeしてから新しい作業を開始
- **何もない場合**: 新しいbranchを作成してOK

## Commit Conventions

### Format

```
{type}({scope}): {brief description}

- {detail 1}
- {detail 2}

Closes #{issue_number}
```

### Types

| Type | Description |
|------|-------------|
| `feat` | New feature |
| `fix` | Bug fix |
| `refactor` | Code refactoring |
| `docs` | Documentation |
| `test` | Tests |
| `ci` | CI/CD changes |
| `chore` | Maintenance |

### Scope (optional)

| Scope | Description |
|-------|-------------|
| `btc` | Bitcoin-related |
| `eth` | Ethereum-related |
| `xrp` | XRP-related |
| `db` | Database |
| `api` | API changes |
| `cli` | CLI changes |
| `pr` | PR review fixes |

### Examples

```bash
# Feature
git commit -m "feat(btc): add taproot address support

- Add bech32m encoding
- Update address validation

Closes #123"

# Bug fix
git commit -m "fix(db): resolve connection timeout

- Increase pool timeout
- Add retry logic

Closes #456"

# Documentation
git commit -m "docs: update architecture guide

- Add layer diagram
- Clarify dependency rules

Closes #789"
```

## Pull Request Creation

### Create PR

```bash
git push -u origin {branch-name}

gh pr create --title "{type}: {description}" --body "$(cat <<'EOF'
## Summary
- {change 1}
- {change 2}

## Test plan
- [ ] Verification commands pass
- [ ] Manual testing completed

Closes #{issue_number}
EOF
)"
```

### PR Title Format

```
{type}: {description} (Closes #{issue_number})
```

Examples:

- `feat: add taproot address support (Closes #123)`
- `fix: resolve database connection timeout (Closes #456)`
- `docs: update architecture guide (Closes #789)`

### PR Description Template

```markdown
## Summary
- Brief description of changes

## Changes
- Change 1
- Change 2

## Test plan
- [ ] Unit tests pass
- [ ] Integration tests pass (if applicable)
- [ ] Manual testing completed

## Related
- Closes #{issue_number}
- Related to #{other_issue}
```

## Safety Rules

### Allowed Operations

- `git add` - Stage changes
- `git commit` - Create commits
- `git push` - Push to remote
- `git checkout -b` - Create branches

### NOT Allowed Operations

- `git push --force` - Never force push
- `git merge` - Don't merge locally
- `git rebase` on shared branches - Avoid rebasing
- Direct commits to `main` - Always use PRs

## Quick Reference

### New Issue Workflow

```bash
# 0. 既存ブランチ/PRを確認（必須）
git branch -a | grep "issue-{number}"
gh pr list --search "issue-{number}"
# → 既存があれば、そのbranchで作業を継続

# 1. Update main
git fetch origin && git checkout main && git reset --hard origin/main

# 2. Create branch (既存がない場合のみ)
git checkout -b {type}/issue-{number}-{description}

# 3. Make changes...

# 4. Commit
git add <files>
git commit -m "{type}: {description}

Closes #{number}"

# 5. Push and create PR
git push -u origin {branch-name}
gh pr create --title "{type}: {description}"
```

### PR Review Fixes

```bash
# Already on PR branch
git add <files>
git commit -m "fix(pr): address review comments"
git push
```
