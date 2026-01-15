---
name: github-issue-creation
description: Create GitHub issues with proper task classification. Classification determines which Skills will be used when working on the issue.
---

# GitHub Issue Creation

Create issues with proper classification. **Labels determine which Skills are used.**

**SSOT Reference**: See [Task Classification](../../../docs/standards/task-classification.md) for the authoritative definition of all labels and mappings.

## Label Categories

### Type Labels (Required)

| Label | Description | Context File |
|-------|-------------|--------------|
| `bug` | Something isn't working | bug-fix.md |
| `enhancement` | New feature or request | feature-add.md |
| `refactoring` | Code improvement | refactoring.md |
| `documentation` | Documentation updates | documentation.md |
| `security` | Security-related | security.md |
| `technical-debt` | Code quality improvements | refactoring.md |
| `test` | Test additions or fixes | test.md |

### Language Labels (For Code Tasks)

| Label | Skill | Verification |
|-------|-------|--------------|
| `lang:go` | `go-development` | `make go-lint && make gotest` |
| `lang:typescript` | `typescript-development` | `npm run lint && npm run build` |
| `lang:solidity` | `solidity-development` | `truffle compile && truffle test` |

### Scope Labels (For Non-Code Tasks)

| Label | Skill | Verification |
|-------|-------|--------------|
| `scope:docs` | `docs-update` | Markdown formatting |
| `scope:devops` | `devops` | `yamllint`, workflow test |
| `scope:scripts` | `shell-scripts` | `make shfmt` |
| `scope:makefile` | `makefile-update` | `make mk-lint` |
| `scope:db` | `db-migration` | `make atlas-lint && make sqlc` |
| `scope:config` | - | Syntax validation |
| `scope:proto` | - | `protoc` compilation |

### Chain Labels (For Cryptocurrency Tasks)

| Label | Description |
|-------|-------------|
| `chain:btc` | Bitcoin-specific |
| `chain:bch` | Bitcoin Cash-specific |
| `chain:eth` | Ethereum-specific |
| `chain:erc20` | ERC-20 token-specific |
| `chain:xrp` | XRP/Ripple-specific |
| `chain:all` | Cross-chain |

### Test Scope Labels (When Type is `test`)

| Label | Description | Verification |
|-------|-------------|--------------|
| `unit-test` | Unit test additions or fixes | `make gotest` |
| `integration-test` | Integration test additions or fixes | `make gotest-integration` |
| `e2e-test` | End-to-end test additions or fixes | `make btc-e2e-*` |

## Label Selection Rules

### Required Labels

Every issue must have:
1. **One Type label** (`bug`, `enhancement`, `refactoring`, `documentation`, `security`, `technical-debt`, `test`)
2. **One Language OR Scope label** (`lang:*` or `scope:*`)

### Optional Labels

- **Chain label**: Only for cryptocurrency-specific code
- **Test Scope label**: Only when Type is `test`

## Label → Skill Mapping Examples

| Task | Labels | Skills Used |
|------|--------|-------------|
| Fix Go bug in BTC | `bug`, `lang:go`, `chain:btc` | `git-workflow` + `go-development` |
| Add TS feature for XRP | `enhancement`, `lang:typescript`, `chain:xrp` | `git-workflow` + `typescript-development` |
| Update README | `documentation`, `scope:docs` | `git-workflow` + `docs-update` |
| Add GitHub Action | `enhancement`, `scope:devops` | `git-workflow` + `devops` |
| Fix shell script | `bug`, `scope:scripts` | `git-workflow` + `shell-scripts` |
| Add Makefile target | `enhancement`, `scope:makefile` | `git-workflow` + `makefile-update` |
| Add DB migration | `enhancement`, `scope:db`, `lang:go` | `git-workflow` + `db-migration` + `go-development` |
| Add unit tests | `test`, `lang:go`, `unit-test` | `git-workflow` + `go-development` |
| Add E2E test for BTC | `test`, `lang:go`, `e2e-test`, `chain:btc` | `git-workflow` + `go-development` |
| Fix integration test | `test`, `lang:go`, `integration-test` | `git-workflow` + `go-development` |

## Issue Creation Process

### 1. Classify Task

From user request, determine:
- Type (bug, enhancement, refactoring, documentation, security, technical-debt, test)
- Language OR Scope
- Test scope (if applicable)
- Chain (if applicable)

### 2. Create Proposal

```markdown
## Proposed Issue

**Title**: [Clear title - 50-72 chars]

**Labels**: [type], [lang/scope], [test scope if applicable], [chain if applicable]

**Skills**: [git-workflow] + [skill based on label]

**Body**:
## Description
[What needs to be done]

## Acceptance Criteria
- [ ] Criterion 1
- [ ] Criterion 2
```

### 3. Create Issue (After Approval)

```bash
gh issue create \
  --title "Title" \
  --body "Body" \
  --label "type,lang/scope,chain"
```

## Quick Reference

```
Required: [Type] + [Language OR Scope]
Optional: [Test Scope] + [Chain]

→ Labels determine Skills
→ Skills determine workflow
```

## Related Documents

- [Task Classification SSOT](../../../docs/standards/task-classification.md) - Authoritative label definitions
- [Task Context Loading](../../rules/task-context-loading.md) - Context loading rules
