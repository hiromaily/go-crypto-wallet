# Shared Command Definitions

This directory contains **shared command definitions** used across all AI agents.
Each agent's command configuration references these definitions to ensure consistency.

## Available Commands

| Command | Description | Usage |
|---------|-------------|-------|
| [fix-issue](fix-issue.md) | Fix GitHub issues | `/fix-issue #123` |
| [fix-linter](fix-linter.md) | Fix linter errors | `/fix-linter` |
| [fix-pr-review](fix-pr-review.md) | Address PR review comments | `/fix-pr-review #123` |
| [task-bug-fix](task-bug-fix.md) | Bug fix workflow | `/task-bug-fix {description}` |
| [task-feature-add](task-feature-add.md) | Feature addition workflow | `/task-feature-add {description}` |
| [task-refactoring](task-refactoring.md) | Refactoring workflow | `/task-refactoring {description}` |
| [task-db-change](task-db-change.md) | Database change workflow | `/task-db-change {description}` |

## Command Structure

Each command follows this template:

```markdown
# Command Name

## Purpose
What this command does

## Parameters
- `{param1}`: Description
- `{param2}`: Description (optional)

## Process
1. Step 1
2. Step 2
3. Step 3

## Verification
Commands to run after completion

## Related Documents
- Link to relevant standards
```

## For AI Agents

- Commands are **human shortcuts** for common operations
- Always reference `docs/standards/` for detailed conventions
- Follow `AGENTS.md` for behavior guidelines

## Agent-Specific Configurations

| Agent | Location | Format |
|-------|----------|--------|
| Claude Code | `.claude/commands/` | Markdown (`.md`) |
| Cursor | `.cursor/commands/` | Markdown (`.md`) |
| Codex | `.codex/prompts/` | Markdown (`.md`) |
