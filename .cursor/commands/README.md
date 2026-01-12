# Cursor Commands

> **Note**: Cursor automatically loads Claude Code configurations from `.claude/` directory.
> Therefore, dedicated Cursor commands are not necessary.

## How It Works

Cursor uses Claude as its AI backend, which means:

- `.claude/commands/` - Automatically available in Cursor
- `.claude/skills/` - Automatically available in Cursor
- `CLAUDE.md` - Automatically loaded

## What Cursor Uses

| Source | Purpose |
|--------|---------|
| `.claude/commands/` | Slash commands (`/fix-issue`, etc.) |
| `.claude/skills/` | Skills (go-development, git-workflow, etc.) |
| `.cursor/rules/` | Cursor-specific rules (if needed) |

## Cursor-Specific Configuration

Only `.cursor/rules/` is Cursor-specific:

- `general.mdc` - General rules
- `security.mdc` - Security rules
- `task-context-loading.mdc` - Context loading rules

These rules supplement (not replace) the Claude configuration.

## Summary

**Do not duplicate commands here.** Use `.claude/commands/` instead.
