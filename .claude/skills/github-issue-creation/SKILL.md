---
name: github-issue-creation
description: Create GitHub issues with proper task classification. Classification determines which Skills will be used when working on the issue.
---

# GitHub Issue Creation

Create issues with proper classification. **Labels determine which Skills are used.**

## Label → Skill Mapping

### Language Labels (code tasks)

| Label | Skill | Verification |
|-------|-------|--------------|
| `lang:go` | `go-development` | `make go-lint && make gotest` |
| `lang:typescript` | `typescript-development` | `npm run lint && npm run build` |
| `lang:solidity` | `solidity-development` | `truffle compile && truffle test` |

### Scope Labels (non-code tasks)

| Label | Skill | Verification |
|-------|-------|--------------|
| `scope:docs` | `docs-update` | Markdown formatting |
| `scope:devops` | `devops` | `yamllint`, workflow test |
| `scope:scripts` | `shell-scripts` | `make shfmt` |
| `scope:makefile` | `makefile-update` | `make mk-lint` |
| `scope:db` | `db-migration` | `make atlas-lint && make sqlc` |
| `scope:config` | (no specific skill) | Syntax validation |
| `scope:proto` | (no specific skill) | `protoc` compilation |

### Chain Labels (additional context)

| Label | Context |
|-------|---------|
| `chain:btc` | Bitcoin-specific considerations |
| `chain:bch` | Bitcoin Cash-specific considerations |
| `chain:eth` | Ethereum-specific considerations |
| `chain:erc20` | ERC-20 token considerations |
| `chain:xrp` | XRP/Ripple considerations |
| `chain:all` | Cross-chain considerations |

### Test Labels (test scope)

| Label | Description | Verification |
|-------|-------------|--------------|
| `unit-test` | Unit test additions or fixes | `make gotest` |
| `integration-test` | Integration test additions or fixes | `make gotest-integration` |
| `e2e-test` | End-to-end test additions or fixes | `make btc-e2e-*` |

## Task Classification

### Step 1: Type Label (required)

| Type | Label | Description |
|------|-------|-------------|
| Bug | `bug` | Something isn't working |
| Feature | `enhancement` | New feature |
| Refactoring | `refactoring` | Code improvement |
| Documentation | `documentation` | Docs updates |
| Security | `security` | Security-related |
| Technical Debt | `technical-debt` | Code quality |

### Step 2: Language OR Scope Label (required)

**Code tasks** → Language label (`lang:*`)
**Non-code tasks** → Scope label (`scope:*`)

### Step 3: Test Label (if applicable)

| Label | When to use |
|-------|-------------|
| `unit-test` | Adding or fixing unit tests |
| `integration-test` | Adding or fixing integration tests |
| `e2e-test` | Adding or fixing E2E tests |

### Step 4: Chain Label (if applicable)

Only for cryptocurrency-specific code.

## Label Examples

| Task | Labels | Skills Used |
|------|--------|-------------|
| Fix Go bug in BTC | `bug`, `lang:go`, `chain:btc` | `git-workflow` + `go-development` |
| Add TS feature for XRP | `enhancement`, `lang:typescript`, `chain:xrp` | `git-workflow` + `typescript-development` |
| Update README | `documentation`, `scope:docs` | `git-workflow` + `docs-update` |
| Add GitHub Action | `enhancement`, `scope:devops` | `git-workflow` + `devops` |
| Fix shell script | `bug`, `scope:scripts` | `git-workflow` + `shell-scripts` |
| Add Makefile target | `enhancement`, `scope:makefile` | `git-workflow` + `makefile-update` |
| Add DB migration | `enhancement`, `scope:db`, `lang:go` | `git-workflow` + `db-migration` + `go-development` |
| Fix unit test | `bug`, `lang:go`, `unit-test` | `git-workflow` + `go-development` |
| Add E2E test for BTC | `enhancement`, `lang:go`, `e2e-test`, `chain:btc` | `git-workflow` + `go-development` |
| Fix integration test | `bug`, `lang:go`, `integration-test` | `git-workflow` + `go-development` |

## Issue Creation Process

### 1. Classify Task

From user request, determine:
- Type (bug, feature, etc.)
- Language OR Scope
- Test scope (if applicable)
- Chain (if applicable)

### 2. Create Proposal

```markdown
## Proposed Issue

**Title**: [Clear title - 50-72 chars]

**Labels**: [type], [lang/scope], [test if applicable], [chain if applicable]

**Skills**: [git-workflow] + [skill based on label]

**Body**:
## Description
[What needs to be done]

## Acceptance Criteria
- [ ] Criterion 1
- [ ] Criterion 2
```

### 3. Create Issue (after approval)

```bash
gh issue create \
  --title "Title" \
  --body "Body" \
  --label "type,lang/scope,chain"
```

## Quick Reference

```
Required: [Type] + [Language OR Scope]
Optional: [Test] + [Chain]

→ Labels determine Skills
→ Skills determine workflow
```
