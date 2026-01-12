---
name: github-issue-creation
description: Create well-structured GitHub issues with proper task classification. This is the starting point for all tasks - proper classification here determines which Skills will be used when working on the issue.
---

# GitHub Issue Creation

Create GitHub issues with proper task classification. **This is the critical first step** that determines how the issue will be worked on.

## Task Classification System

### Step 1: Identify Task Type

| Type | Label | Description |
|------|-------|-------------|
| Bug Fix | `bug` | Something isn't working |
| Feature | `enhancement` | New feature or request |
| Refactoring | `refactoring` | Code improvement without behavior change |
| Documentation | `documentation` | Docs updates |
| Security | `security` | Security-related |
| Technical Debt | `technical-debt` | Code quality cleanup |

### Step 2: Identify Scope

**Code Tasks** - Select language:

| Scope | Label | Affected Directories |
|-------|-------|---------------------|
| Go Development | `lang:go` | `internal/`, `pkg/`, `cmd/` |
| TypeScript Development | `lang:typescript` | `apps/ripple-lib-server/` |
| Solidity Development | `lang:solidity` | `apps/erc20-token/contracts/` |

**Non-Code Tasks** - Select scope:

| Scope | Label | Affected Directories |
|-------|-------|---------------------|
| Documentation | `scope:docs` | `docs/`, `*.md` |
| DevOps/CI/CD | `scope:devops` | `.github/workflows/`, `docker/`, `compose.*.yaml` |
| Shell Scripts | `scope:scripts` | `scripts/`, `*.sh` |
| Makefile | `scope:makefile` | `Makefile`, `make/` |
| Configuration | `scope:config` | `config/`, `*.toml`, `*.yaml` |
| Protocol Buffers | `scope:proto` | `proto/` |
| Database | `scope:db` | `tools/atlas/`, `tools/sqlc/` |

### Step 3: Identify Chain (if applicable)

Only for code that involves specific cryptocurrency:

| Chain | Label | When to Use |
|-------|-------|-------------|
| Bitcoin | `chain:btc` | BTC-specific code |
| Bitcoin Cash | `chain:bch` | BCH-specific code |
| Ethereum | `chain:eth` | ETH-specific code |
| ERC-20 | `chain:erc20` | Token contract code |
| XRP | `chain:xrp` | Ripple-specific code |
| All Chains | `chain:all` | Cross-chain code |

## Label Combination Examples

| Task Description | Labels |
|------------------|--------|
| Fix bug in Bitcoin address generation | `bug`, `lang:go`, `chain:btc` |
| Add new ETH transaction type | `enhancement`, `lang:go`, `chain:eth` |
| Refactor XRP gRPC server | `refactoring`, `lang:typescript`, `chain:xrp` |
| Update ARCHITECTURE.md | `documentation`, `scope:docs` |
| Add new GitHub Action workflow | `enhancement`, `scope:devops` |
| Fix shell script permission issue | `bug`, `scope:scripts` |
| Update Makefile targets | `enhancement`, `scope:makefile` |
| Add new database migration | `enhancement`, `scope:db`, `lang:go` |

## Issue Creation Process

### 1. Gather Information

From user request, identify:
- **What**: Clear description of the task
- **Why**: Context and motivation
- **Type**: Bug, feature, refactoring, etc.
- **Scope**: Language or non-code scope
- **Chain**: If cryptocurrency-specific

### 2. Determine Labels

Apply classification:
1. One **Type** label (required)
2. One **Language** OR **Scope** label (required)
3. One **Chain** label (if applicable)

### 3. Create Issue Proposal

```markdown
## Proposed Issue

**Title**: [Clear, imperative title - 50-72 chars]

**Labels**: [type], [lang/scope], [chain if applicable]

**Body**:

## Description
[What needs to be done]

## Context
[Why this is needed]

## Acceptance Criteria
- [ ] Criterion 1
- [ ] Criterion 2

## Technical Notes
- Affected files: [list]
- Related docs: [links]

---
**Ready to create?** Confirm to proceed.
```

### 4. Create Issue (After Approval)

```bash
# Verify labels exist
gh label list

# Create issue
gh issue create \
  --title "Title" \
  --body "Body" \
  --label "type,lang/scope,chain"
```

## Skill Mapping

When the issue is worked on, these labels determine which Skill to use:

| Label | Skill to Use |
|-------|--------------|
| `lang:go` | `go-development` |
| `lang:typescript` | `typescript-development` |
| `lang:solidity` | `solidity-development` |
| `scope:docs` | No specific skill (follow docs standards) |
| `scope:devops` | No specific skill (follow DevOps practices) |
| `scope:scripts` | No specific skill (run `make shfmt`) |
| `scope:makefile` | No specific skill (run `make mk-lint`) |
| `scope:db` | `go-development` + database workflow |

## Quick Reference

### Required Labels Per Issue

```
[Type] + [Language OR Scope] + [Chain if applicable]

Examples:
- bug + lang:go + chain:btc
- enhancement + scope:devops
- documentation + scope:docs
- refactoring + lang:typescript + chain:xrp
```

### Label Sync Command

If labels are missing:

```bash
# Check existing labels
gh label list

# Create missing labels from labels.yml
# (Labels are defined in .github/labels.yml)
```

## Important Notes

1. **Classification is critical** - It determines workflow and verification commands
2. **One type label** - Don't mix bug and enhancement
3. **Language XOR Scope** - Code tasks use language, non-code use scope
4. **Chain is optional** - Only for cryptocurrency-specific code
5. **Wait for approval** - Always confirm before creating issue
