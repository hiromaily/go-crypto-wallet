# Shared Command Definitions

Minimal command definitions. Detailed workflows are in agent Skills.

## Available Commands

| Command | Description |
|---------|-------------|
| [fix-issue](fix-issue.md) | Fix GitHub issues |
| [fix-linter](fix-linter.md) | Fix linter errors |
| [fix-pr-review](fix-pr-review.md) | Address PR review comments |

## Workflow Reference

All commands follow the common Go development workflow:

1. **Branch**: Create from latest `main`
2. **Implement**: Follow Clean Architecture
3. **Verify**: `make go-lint && make tidy && make check-build && make gotest`
4. **Review**: Self-review checklist
5. **Commit**: Conventional commit format
6. **PR**: Create with description

## Related Skills

| Agent | Skill Location |
|-------|---------------|
| Claude | `.claude/skills/go-development/SKILL.md` |
| Cursor | `.cursor/rules/general.mdc` |
| Codex | `.codex/rules/general.md` |
