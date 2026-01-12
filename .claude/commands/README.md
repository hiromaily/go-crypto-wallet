# Claude Code Commands

Commands for Claude Code. These reference the shared command definitions in `docs/commands/`.

## Available Commands

| Command | Description |
|---------|-------------|
| `/fix-issue` | Fix GitHub issues |
| `/fix-linter` | Fix linter errors |
| `/fix-pr-review` | Address PR review comments |
| `/task-bug-fix` | Bug fix workflow |
| `/task-feature-add` | Feature addition workflow |
| `/task-refactoring` | Refactoring workflow |
| `/task-db-change` | Database change workflow |

## Claude-Specific Commands

| Command | Description |
|---------|-------------|
| `/convert-custom-slash-for-codex` | Convert Claude commands to Codex format |
| `/create-github-issue` | Create GitHub issue (migrated to Skills) |

## Shared Definitions

All common command definitions are maintained in @docs/commands/.

## Claude-Specific Notes

- Use `@file-name` syntax for file references
- Commands use Markdown format
- Reference @docs/standards/ for coding conventions
- Follow @AGENTS.md for behavior guidelines
