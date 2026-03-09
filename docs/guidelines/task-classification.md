# Task Classification (SSOT)

This document is the **Single Source of Truth (SSOT)** for task classification in this project.
All labels, task types, and their mappings are defined here.

## Overview

Task classification determines:

1. Which **GitHub labels** to apply to issues
2. Which **context files** to load when working on tasks
3. Which **Skills** to use for implementation

## Label Categories

### Type Labels (Required)

Define the nature of the task.

| Label | Description | Context File | Keywords |
|-------|-------------|--------------|----------|
| `bug` | Something isn't working | `docs/task-contexts/bug-fix.md` | bug, fix, error, issue |
| `enhancement` | New feature or request | `docs/task-contexts/feature-add.md` | add, implement, feature, new |
| `refactoring` | Code improvement without changing functionality | `docs/task-contexts/refactoring.md` | refactor, reorganize, move, cleanup |
| `documentation` | Documentation updates | `docs/task-contexts/documentation.md` | document, README, docs, comment |
| `security` | Security-related issues | `docs/task-contexts/security.md` | security, vulnerability, CVE |
| `technical-debt` | Code quality improvements | `docs/task-contexts/refactoring.md` | debt, cleanup, improve |
| `test` | Test additions or fixes | `docs/task-contexts/test.md` | test, coverage, spec |

### Language Labels (For Code Tasks)

Specify the programming language involved.

| Label | Description | Skill | Verification |
|-------|-------------|-------|--------------|
| `lang:go` | Go language code | `go-development` | `make go-lint && make go-test` |
| `lang:typescript` | TypeScript/JavaScript code | `typescript-development` | `npm run lint && npm run build` |
| `lang:solidity` | Solidity smart contracts | `solidity-development` | `truffle compile && truffle test` |

### Scope Labels (For Non-Code Tasks)

Specify the scope when not primarily code changes.

| Label | Description | Skill | Verification |
|-------|-------------|-------|--------------|
| `scope:docs` | Documentation files (docs/, *.md) | `docs-update` | Markdown formatting |
| `scope:devops` | CI/CD, GitHub Actions, Docker | `devops` | `yamllint`, workflow test |
| `scope:scripts` | Shell scripts (scripts/, *.sh) | `shell-scripts` | `make shfmt` |
| `scope:makefile` | Makefile and make/ directory | `makefile-update` | `make mk-lint` |
| `scope:db` | Database schema and migrations | `db-migration` | `make atlas-lint && make sqlc` |
| `scope:config` | Configuration files (config/, *.toml,*.yaml) | - | Syntax validation |
| `scope:proto` | Protocol Buffer definitions (proto/) | - | `protoc` compilation |

### Chain Labels (For Cryptocurrency Tasks)

Specify the blockchain when applicable.

| Label | Description | Keywords |
|-------|-------------|----------|
| `chain:btc` | Bitcoin (BTC) related | Bitcoin, BTC, Taproot, Descriptor, PSBT, MuSig2 |
| `chain:bch` | Bitcoin Cash (BCH) related | Bitcoin Cash, BCH, CashAddr |
| `chain:eth` | Ethereum (ETH) related | Ethereum, ETH, Gas, Nonce |
| `chain:erc20` | ERC-20 token related | ERC-20, token |
| `chain:xrp` | XRP/Ripple related | Ripple, XRP, Destination Tag |
| `chain:all` | Affects all supported chains | multi-chain, all chains |

### Test Scope Labels (For Test Tasks)

Specify the test granularity when Type is `test`.

| Label | Description | Verification |
|-------|-------------|--------------|
| `unit-test` | Unit test additions or fixes | `make go-test` |
| `integration-test` | Integration test additions or fixes | `make go-test-integration` |
| `e2e-test` | End-to-end test additions or fixes | `make btc-e2e-*` |

## Label Selection Rules

### Required Labels

Every issue must have:

1. **One Type label** (bug, enhancement, refactoring, documentation, security, technical-debt, test)
2. **One Language OR Scope label** (lang:*or scope:*)

### Optional Labels

- **Chain label**: Only for cryptocurrency-specific code
- **Test Scope label**: Only when Type is `test`

## Task Detection from Labels

When working on a GitHub Issue, determine task context from labels:

```
Priority Order:
1. Type label → Context file
2. Language label → Skill + Verification commands
3. Scope label → Skill + Verification commands (if no language label)
4. Chain label → Chain-specific context
5. Test Scope label → Test-specific verification
```

### Example Mappings

| Labels | Skills Used | Context Files |
|--------|-------------|---------------|
| `bug`, `lang:go`, `chain:btc` | `git-workflow` + `go-development` | bug-fix.md, chains/btc.md |
| `enhancement`, `lang:typescript`, `chain:xrp` | `git-workflow` + `typescript-development` | feature-add.md, chains/xrp.md |
| `test`, `lang:go`, `unit-test` | `git-workflow` + `go-development` | test.md |
| `documentation`, `scope:docs` | `git-workflow` + `docs-update` | documentation.md |
| `refactoring`, `scope:db`, `lang:go` | `git-workflow` + `db-migration` + `go-development` | refactoring.md, db-change.md |

## Context File Mapping

| Task Type | Primary Context | Additional Contexts |
|-----------|-----------------|---------------------|
| bug | bug-fix.md | - |
| enhancement | feature-add.md | - |
| refactoring | refactoring.md | - |
| documentation | documentation.md | - |
| security | security.md | - |
| technical-debt | refactoring.md | - |
| test | test.md | - |
| (with chain:*) | chain-specific.md | chains/{chain}.md |
| (with scope:db) | db-change.md | - |

## Verification Commands by File Type

| File Extension | Required | Optional |
|----------------|----------|----------|
| `*.go` | `make go-lint`, `make check-build` | `make go-test` |
| `*.md`, `*.mdc` | - | markdownlint |
| `*.sql`, `*.hcl` | `make atlas-fmt`, `make atlas-lint` | `make sqlc` |
| `*.yaml`, `*.toml` | - | yamllint |
| `*.sh` | - | shellcheck, `make shfmt` |
| `*.proto` | `make proto` | - |
| `*.ts`, `*.js` | `npm run lint`, `npm run build` | `npm test` |

## Related Documents

- [GitHub Labels](https://github.com/hiromaily/go-crypto-wallet/blob/main/.github/labels.yml) - Label definitions with colors
- `.claude/rules/task-context-loading.md` - Auto-loading rules
- `.claude/skills/github-issue-creation/SKILL.md` - Issue creation workflow
- [Task Contexts](../task-contexts/README.md) - Context file details
