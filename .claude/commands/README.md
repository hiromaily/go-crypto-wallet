# Claude Code Commands

Commands that select appropriate Skills based on task classification.

## Available Commands

| Command | Description |
|---------|-------------|
| `/fix-issue #123` | Fix a GitHub issue |
| `/fix-linter` | Fix linter errors |
| `/fix-pr-review #123` | Address PR review comments |

## Task Classification

Commands automatically select Skills based on:

### Language (from labels or files)

| Classification | Skill |
|----------------|-------|
| `lang:go` / Go files | `go-development` |
| `lang:typescript` / TS files | `typescript-development` |
| `lang:solidity` / Solidity files | `solidity-development` |

### Chain (from labels)

| Classification | Context |
|----------------|---------|
| `chain:btc` | Bitcoin considerations |
| `chain:bch` | Bitcoin Cash considerations |
| `chain:eth` | Ethereum considerations |
| `chain:erc20` | ERC-20 token considerations |
| `chain:xrp` | XRP/Ripple considerations |
| `chain:all` | Cross-chain considerations |

## Skills

| Skill | Description | Path |
|-------|-------------|------|
| `go-development` | Go workflow | `.claude/skills/go-development/SKILL.md` |
| `typescript-development` | TypeScript/JS workflow | `.claude/skills/typescript-development/SKILL.md` |
| `solidity-development` | Solidity workflow | `.claude/skills/solidity-development/SKILL.md` |
| `github-issue-creation` | Issue creation | `.claude/skills/github-issue-creation/SKILL.md` |
